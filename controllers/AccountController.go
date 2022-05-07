package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/lifei6671/gocaptcha"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/mail"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/dingtalk"
	"github.com/mindoc-org/mindoc/utils/workweixin"
)

const (
	WorkWeixin_AuthorizeUrlBase = "https://open.weixin.qq.com/connect/oauth2/authorize"
	WorkWeixin_QRConnectUrlBase = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"
	SessionUserInfoKey          = "session-user-info-key"
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

// AccountController 用户登录与注册
type AccountController struct {
	BaseController
}

func (c *AccountController) referer() string {
	u, _ := url.PathUnescape(c.GetString("url"))
	if u == "" {
		u = conf.URLFor("HomeController.Index")
	}
	return u
}

func (c *AccountController) IsInWorkWeixin() (is_in_workweixin bool) {
	ua := c.Ctx.Input.UserAgent()
	var wechatRule = regexp.MustCompile(`\bMicroMessenger\/\d+(\.\d+)*\b`)
	var wxworkRule = regexp.MustCompile(`\bwxwork\/\d+(\.\d+)*\b`)
	return wechatRule.MatchString(ua) && wxworkRule.MatchString(ua)
}

func (c *AccountController) Prepare() {
	c.BaseController.Prepare()
	c.EnableXSRF = web.AppConfig.DefaultBool("enablexsrf", true)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	c.Data["CanLoginWorkWeixin"] = len(web.AppConfig.DefaultString("workweixin_corpid", "")) > 0

	c.Data["corpID"], _ = web.AppConfig.String("dingtalk_corpid")
	c.Data["CanLoginDingTalk"] = len(web.AppConfig.DefaultString("dingtalk_corpid", "")) > 0
	if reflect.ValueOf(c.Data["CanLoginDingTalk"]).Bool() {
		c.Data["ENABLE_QR_DINGTALK"] = true
	}
	c.Data["dingtalk_qr_key"], _ = web.AppConfig.String("dingtalk_qr_key")

	if !c.EnableXSRF {
		return
	}
	if c.Ctx.Input.IsPost() {
		token := c.Ctx.Input.Query("_xsrf")
		if token == "" {
			token = c.Ctx.Request.Header.Get("X-Xsrftoken")
		}
		if token == "" {
			token = c.Ctx.Request.Header.Get("X-Csrftoken")
		}
		if token == "" {
			if c.IsAjax() {
				c.JsonResult(403, i18n.Tr(c.Lang, "message.illegal_request"))
			} else {
				c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.illegal_request"))
			}
		}
		xsrfToken := c.XSRFToken()
		if xsrfToken != token {
			if c.IsAjax() {
				c.JsonResult(403, i18n.Tr(c.Lang, "message.illegal_request"))
			} else {
				c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.illegal_request"))
			}
		}
	}
}

// Login 用户登录
func (c *AccountController) Login() {
	c.Prepare()

	c.TplName = "account/login.tpl"

	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
		u := c.GetString("url")
		if u == "" {
			u = c.Ctx.Request.Header.Get("Referer")
		}
		if u == "" {
			u = conf.URLFor("HomeController.Index")
		}
		c.Redirect(u, 302)
	}
	var remember CookieRemember
	// 如果 Cookie 中存在登录信息
	if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
		if err := utils.Decode(cookie, &remember); err == nil {
			if member, err := models.NewMember().Find(remember.MemberId); err == nil {
				c.SetMember(*member)
				c.LoggedIn(false)
				c.StopRun()
			}
		}
	}

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password := c.GetString("password")
		captcha := c.GetString("code")
		isRemember := c.GetString("is_remember")

		// 如果开启了验证码
		if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
			v, ok := c.GetSession(conf.CaptchaSessionName).(string)
			if !ok || !strings.EqualFold(v, captcha) {
				c.JsonResult(6001, i18n.Tr(c.Lang, "message.captcha_wrong"))
			}
		}

		if account == "" || password == "" {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.account_or_password_empty"))
		}

		member, err := models.NewMember().Login(account, password)
		if err == nil {
			member.LastLoginTime = time.Now()
			_ = member.Update("last_login_time")

			c.SetMember(*member)

			if strings.EqualFold(isRemember, "yes") {
				remember.MemberId = member.MemberId
				remember.Account = member.Account
				remember.Time = time.Now()
				v, err := utils.Encode(remember)
				if err == nil {
					c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30).Unix())
				}
			}

			c.JsonResult(0, "ok", c.referer())
		} else {
			logs.Error("用户登录 ->", err)
			c.JsonResult(500, i18n.Tr(c.Lang, "message.wrong_account_password"), nil)
		}
	} else {
		// 默认登录方式
		login_method := "AccountController.Login"
		var redirect_uri string
		// 企业微信登录检查
		canLoginWorkWeixin := reflect.ValueOf(c.Data["CanLoginWorkWeixin"]).Bool()
		referer := c.referer()
		if canLoginWorkWeixin {
			// 企业微信登录方式
			login_method = "AccountController.WorkWeixinLogin"
			u := c.GetString("url")
			if u == "" {
				u = referer
				if u == "" {
					u = conf.BaseUrl
				}
			} else {
				var schemaRule = regexp.MustCompile(`^https?\:\/\/`)
				if !schemaRule.MatchString(u) {
					u = conf.BaseUrl + u
				}
			}
			redirect_uri = conf.URLFor(login_method, "url", url.PathEscape(u))
			// 是否在企业微信内部打开
			isInWorkWeixin := c.IsInWorkWeixin()
			c.Data["IsInWorkWeixin"] = isInWorkWeixin
			if isInWorkWeixin {
				// 客户端拥有微信标识和企业微信标识
				c.Redirect(redirect_uri, 302)
				return
			} else {
				c.Data["workweixin_login_url"] = redirect_uri
			}
		}
		c.Data["url"] = referer
	}
}

// 钉钉登录
func (c *AccountController) DingTalkLogin() {
	c.Prepare()

	code := c.GetString("dingtalk_code")
	if code == "" {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_obtain_user_info"), nil)
	}

	appKey, _ := web.AppConfig.String("dingtalk_app_key")
	appSecret, _ := web.AppConfig.String("dingtalk_app_secret")
	tmpReader, _ := web.AppConfig.String("dingtalk_tmp_reader")

	if appKey == "" || appSecret == "" || tmpReader == "" {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.dingtalk_auto_login_not_enable"), nil)
		c.StopRun()
	}

	dingtalkAgent := dingtalk.NewDingTalkAgent(appSecret, appKey)
	err := dingtalkAgent.GetAccesstoken()
	if err != nil {
		logs.Warn("获取钉钉临时Token失败 ->", err)
		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
		c.StopRun()
	}

	userid, err := dingtalkAgent.GetUserIDByCode(code)
	if err != nil {
		logs.Warn("获取钉钉用户ID失败 ->", err)
		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
		c.StopRun()
	}

	username, avatar, err := dingtalkAgent.GetUserNameAndAvatarByUserID(userid)
	if err != nil {
		logs.Warn("获取钉钉用户信息失败 ->", err)
		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
		c.StopRun()
	}

	member, err := models.NewMember().TmpLogin(tmpReader)
	if err == nil {
		member.LastLoginTime = time.Now()
		_ = member.Update("last_login_time")
		member.Account = username
		if avatar != "" {
			member.Avatar = avatar
		}

		c.SetMember(*member)
	}
	c.JsonResult(0, "ok", username)
}

// WorkWeixinLogin 用户企业微信登录
func (c *AccountController) WorkWeixinLogin() {
	c.Prepare()

	logs.Info("UserAgent: ", c.Ctx.Input.UserAgent()) // debug

	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
		u := c.GetString("url")
		if u == "" {
			u = c.Ctx.Request.Header.Get("Referer")
			if u == "" {
				u = conf.URLFor("HomeController.Index")
			}
		}
		// session自动登录时刷新session内容
		member, err := models.NewMember().Find(member.MemberId)
		if err != nil {
			c.DelSession(conf.LoginSessionName)
			c.SetMember(models.Member{})
			c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		} else {
			c.SetMember(*member)
		}
		c.Redirect(u, 302)
	}
	var remember CookieRemember
	// 如果 Cookie 中存在登录信息
	if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
		if err := utils.Decode(cookie, &remember); err == nil {
			if member, err := models.NewMember().Find(remember.MemberId); err == nil {
				c.SetMember(*member)
				c.LoggedIn(false)
				c.StopRun()
			}
		}
	}

	if c.Ctx.Input.IsPost() {
		// account := c.GetString("account")
		// password := c.GetString("password")
		// captcha := c.GetString("code")
		// isRemember := c.GetString("is_remember")
		c.JsonResult(400, "request method not allowed", nil)
	} else {
		var callback_u string
		u := c.GetString("url")
		if u == "" {
			u = c.referer()
		}
		if u != "" {
			var schemaRule = regexp.MustCompile(`^https?\:\/\/`)
			if !schemaRule.MatchString(u) {
				u = strings.TrimRight(conf.BaseUrl, "/") + strings.TrimLeft(u, "/")
			}
		}
		if u == "" {
			callback_u = conf.URLFor("AccountController.WorkWeixinLoginCallback")
		} else {
			callback_u = conf.URLFor("AccountController.WorkWeixinLoginCallback", "url", url.PathEscape(u))
		}
		logs.Info("callback_u: ", callback_u) // debug

		state := "mindoc"
		workweixinConf := conf.GetWorkWeixinConfig()
		appid := workweixinConf.CorpId
		agentid := workweixinConf.AgentId
		var redirect_uri string

		isInWorkWeixin := c.IsInWorkWeixin()
		c.Data["IsInWorkWeixin"] = isInWorkWeixin
		if isInWorkWeixin {
			// 企业微信内-网页授权登录
			urlFmt := "%s?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
			redirect_uri = fmt.Sprintf(urlFmt, WorkWeixin_AuthorizeUrlBase, appid, url.PathEscape(callback_u), state)
		} else {
			// 浏览器内-扫码授权登录
			urlFmt := "%s?appid=%s&agentid=%s&redirect_uri=%s&state=%s"
			redirect_uri = fmt.Sprintf(urlFmt, WorkWeixin_QRConnectUrlBase, appid, agentid, url.PathEscape(callback_u), state)
		}
		logs.Info("redirect_uri: ", redirect_uri) // debug
		c.Redirect(redirect_uri, 302)
	}
}

/*
思路:
1. 浏览器打开
        用户名+密码 登录 与企业微信没有交集
        手机企业微信登录->扫码页面->扫码后获取用户信息, 判断是否绑定了企业微信
            已绑定，则读取用户信息，直接登录
            未绑定，则弹窗提示[未绑定企业微信，请先在企业微信中打开，完成绑定]
2. 企业微信打开->自动登录->判断是否绑定了企业微信
        已绑定，则读取用户信息，直接登录
        未绑定，则弹窗提示
            是否已有账户(用户名+密码方式)
                有: 弹窗输入[用户名+密码+验证码]校验
                无: 直接以企业UserId作为用户名(小写)，创建随机密码
*/

// WorkWeixinLoginCallback 用户企业微信登录-回调
func (c *AccountController) WorkWeixinLoginCallback() {
	c.TplName = "account/workweixin-login-callback.tpl"

	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
		u := c.GetString("url")
		if u == "" {
			u = c.Ctx.Request.Header.Get("Referer")
		}
		if u == "" {
			u = conf.URLFor("HomeController.Index")
		}
		member, err := models.NewMember().Find(member.MemberId)
		if err != nil {
			c.DelSession(conf.LoginSessionName)
			c.SetMember(models.Member{})
			c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		} else {
			c.SetMember(*member)
		}
		c.Redirect(u, 302)
	}

	var remember CookieRemember
	// 如果 Cookie 中存在登录信息
	if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
		if err := utils.Decode(cookie, &remember); err == nil {
			if member, err := models.NewMember().Find(remember.MemberId); err == nil {
				c.SetMember(*member)
				c.LoggedIn(false)
				c.StopRun()
			}
		}
	}

	// 请求参数获取
	req_code := c.GetString("code")
	logs.Warning("req_code: ", req_code)
	req_state := c.GetString("state")
	logs.Warning("req_state: ", req_state)
	var user_info_json string
	var error_msg string
	var bind_existed string
	if len(req_code) > 0 && req_state == "mindoc" {
		// 获取当前应用的access_token
		access_token, ok := workweixin.GetAccessToken(false)
		if ok {
			logs.Warning("access_token: ", access_token)
			// 获取当前请求的userid
			user_id, ok := workweixin.RequestUserId(access_token, req_code)
			if ok {
				logs.Warning("user_id: ", user_id)
				// 获取通讯录应用的access_token
				contact_access_token, ok := workweixin.GetAccessToken(true)
				if ok {
					logs.Warning("contact_access_token: ", contact_access_token)
					user_info, err_msg, ok := workweixin.RequestUserInfo(contact_access_token, user_id)
					if ok {
						// [-------所有字段-Debug----------
						// user_info.UserId
						// user_info.Name
						// user_info.HideMobile
						// user_info.Mobile
						// user_info.Department
						// user_info.Email
						// user_info.IsLeaderInDept
						// user_info.IsLeader
						// user_info.Avatar
						// user_info.Alias
						// user_info.Status
						// user_info.MainDepartment
						// -----------------------------]
						// logs.Debug("user_info.UserId: ", user_info.UserId)
						// logs.Debug("user_info.Name: ", user_info.Name)
						json_info, _ := json.Marshal(user_info)
						user_info_json = string(json_info)
						// 查询系统现有数据，是否绑定了当前请求用户的企业微信
						member, err := models.NewWorkWeixinAccount().ExistedMember(user_info.UserId)
						if err == nil {
							member.LastLoginTime = time.Now()
							_ = member.Update("last_login_time")

							c.SetMember(*member)

							var remember CookieRemember
							remember.MemberId = member.MemberId
							remember.Account = member.Account
							remember.Time = time.Now()
							v, err := utils.Encode(remember)
							if err == nil {
								c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
							}
							bind_existed = "true"
							error_msg = ""
							u := c.GetString("url")
							if u == "" {
								u = conf.URLFor("HomeController.Index")
							}
							c.Redirect(u, 302)
						} else {
							if err == orm.ErrNoRows {
								c.SetSession(SessionUserInfoKey, user_info)
								bind_existed = "false"
								error_msg = ""
							} else {
								logs.Error("Error: ", err)
								error_msg = "数据库错误: " + err.Error()
							}
						}
						//
					} else {
						error_msg = "获取用户信息失败: " + err_msg
					}
				} else {
					error_msg = "通讯录访问凭据获取失败: " + contact_access_token
				}
			} else {
				error_msg = "获取用户Id失败: " + user_id
			}
		} else {
			error_msg = "应用凭据获取失败: " + access_token
		}
	} else {
		error_msg = "参数错误"
	}
	if user_info_json == "" {
		user_info_json = "{}"
	}
	if bind_existed == "" {
		bind_existed = "null"
	}
	// refer & doc:
	// - https://golang.org/pkg/html/template/#HTML
	// - https://stackoverflow.com/questions/24411880/go-html-templates-can-i-stop-the-templates-package-inserting-quotes-around-stri
	// - https://stackoverflow.com/questions/38035176/insert-javascript-snippet-inside-template-with-beego-golang
	c.Data["bind_existed"] = template.JS(bind_existed)
	logs.Debug("bind_existed: ", bind_existed)
	c.Data["error_msg"] = template.JS(error_msg)
	c.Data["user_info_json"] = template.JS(user_info_json)
	/*
		// 调试: 显示源码
		result, err := c.RenderString()
		if err != nil {
			logs.Error(err)
		} else {
			logs.Warning(result)
		}
	*/
}

// WorkWeixinLoginBind 用户企业微信登录-绑定
func (c *AccountController) WorkWeixinLoginBind() {
	c.Prepare()

	if user_info, ok := c.GetSession(SessionUserInfoKey).(workweixin.WorkWeixinUserInfo); ok && len(user_info.UserId) > 0 {
		req_account := c.GetString("account")
		req_password := c.GetString("password")
		if req_account == "" || req_password == "" {
			c.JsonResult(400, "账号或密码不能为空")
		} else {
			member, err := models.NewMember().Login(req_account, req_password)
			if err == nil {
				account := models.NewWorkWeixinAccount()
				account.MemberId = member.MemberId
				account.WorkWeixin_UserId = user_info.UserId
				member.CreateAt = 0
				ormer := orm.NewOrm()
				o, err := ormer.Begin()
				if err != nil {
					logs.Error("开启事物时出错 -> ", err)
					c.JsonResult(500, "开启事物时出错: ", err.Error())
				}
				if err := account.AddBind(ormer); err != nil {
					o.Rollback()
					c.JsonResult(500, "绑定失败，数据库错误: "+err.Error())
				} else {
					member.LastLoginTime = time.Now()
					member.RealName = user_info.Name
					member.Avatar = user_info.Avatar
					if len(member.Avatar) < 1 {
						member.Avatar = conf.GetDefaultAvatar()
					}
					member.Email = user_info.Email
					member.Phone = user_info.Mobile
					if _, err := ormer.Update(member, "last_login_time", "real_name", "avatar", "email", "phone"); err != nil {
						o.Rollback()
						logs.Error("保存用户信息失败=>", err)
						c.JsonResult(500, "绑定失败，现有账户信息更新失败: "+err.Error())
					} else {
						if err := o.Commit(); err != nil {
							logs.Error("提交事物时出错 -> ", err)
							c.JsonResult(500, "提交事物时出错: ", err.Error())
						} else {
							c.DelSession(SessionUserInfoKey)
							c.SetMember(*member)

							var remember CookieRemember
							remember.MemberId = member.MemberId
							remember.Account = member.Account
							remember.Time = time.Now()
							v, err := utils.Encode(remember)
							if err == nil {
								c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
								c.JsonResult(0, "绑定成功", nil)
							} else {
								c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
							}
						}
					}

				}

			} else {
				logs.Error("用户登录 ->", err)
				c.JsonResult(500, "账号或密码错误", nil)
			}
			c.JsonResult(500, "TODO: 绑定以后账号功能开发中")
		}
	} else {
		if ok {
			c.DelSession(SessionUserInfoKey)
		}
		c.JsonResult(400, "请求错误, 请从首页重新登录")
	}

}

// WorkWeixinLoginIgnore 用户企业微信登录-忽略
func (c *AccountController) WorkWeixinLoginIgnore() {
	if user_info, ok := c.GetSession(SessionUserInfoKey).(workweixin.WorkWeixinUserInfo); ok && len(user_info.UserId) > 0 {
		c.DelSession(SessionUserInfoKey)
		member := models.NewMember()

		if _, err := member.FindByAccount(user_info.UserId); err == nil && member.MemberId > 0 {
			c.JsonResult(400, "账号已存在")
		}

		ormer := orm.NewOrm()
		o, err := ormer.Begin()
		if err != nil {
			logs.Error("开启事物时出错 -> ", err)
			c.JsonResult(500, "开启事物时出错: ", err.Error())
		}

		member.Account = user_info.UserId
		member.RealName = user_info.Name
		var rnd = rand.New(src)
		// fmt.Sprintf("%x", rnd.Uint64())
		// strconv.FormatUint(rnd.Uint64(), 16)
		member.Password = user_info.UserId + strconv.FormatUint(rnd.Uint64(), 16)
		member.Password = "pathea.2020" // 强制设置默认密码，不然无法修改密码(因为目前修改密码需要知道当前密码)
		hash, err := utils.PasswordHash(member.Password)
		if err != nil {
			logs.Error("加密用户密码失败 =>", err)
			c.JsonResult(500, "加密用户密码失败"+err.Error())
		} else {
			logs.Error("member.Password: ", member.Password)
			logs.Error("hash: ", hash)
			member.Password = hash
		}
		member.Role = conf.MemberGeneralRole
		member.Avatar = user_info.Avatar
		if len(member.Avatar) < 1 {
			member.Avatar = conf.GetDefaultAvatar()
		}
		member.CreateAt = 0
		member.Email = user_info.Email
		member.Phone = user_info.Mobile
		member.Status = 0
		if _, err = ormer.Insert(member); err != nil {
			o.Rollback()
			c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
		} else {
			account := models.NewWorkWeixinAccount()
			account.MemberId = member.MemberId
			account.WorkWeixin_UserId = user_info.UserId
			member.CreateAt = 0
			if err := account.AddBind(ormer); err != nil {
				o.Rollback()
				c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
			} else {
				if err := o.Commit(); err != nil {
					logs.Error("提交事物时出错 -> ", err)
					c.JsonResult(500, "提交事物时出错: ", err.Error())
				} else {
					member.LastLoginTime = time.Now()
					_ = member.Update("last_login_time")

					c.SetMember(*member)

					var remember CookieRemember
					remember.MemberId = member.MemberId
					remember.Account = member.Account
					remember.Time = time.Now()
					v, err := utils.Encode(remember)
					if err == nil {
						c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
						c.JsonResult(0, "绑定成功", nil)
					} else {
						c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
					}
				}
			}
		}
	} else {
		if ok {
			c.DelSession(SessionUserInfoKey)
		}
		c.JsonResult(400, "请求错误, 请从首页重新登录")
	}
}

// QR二维码登录
func (c *AccountController) QRLogin() {
	c.Prepare()

	appName := c.Ctx.Input.Param(":app")

	switch appName {
	// 钉钉扫码登录
	case "dingtalk":
		code := c.GetString("code")
		state := c.GetString("state")
		if state != "1" || code == "" {
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}
		appKey, _ := web.AppConfig.String("dingtalk_qr_key")
		appSecret, _ := web.AppConfig.String("dingtalk_qr_secret")

		qrDingtalk := dingtalk.NewDingtalkQRLogin(appSecret, appKey)
		unionID, err := qrDingtalk.GetUnionIDByCode(code)
		if err != nil {
			logs.Warn("获取钉钉临时UnionID失败 ->", err)
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}

		appKey, _ = web.AppConfig.String("dingtalk_app_key")
		appSecret, _ = web.AppConfig.String("dingtalk_app_secret")
		tmpReader, _ := web.AppConfig.String("dingtalk_tmp_reader")

		dingtalkAgent := dingtalk.NewDingTalkAgent(appSecret, appKey)
		err = dingtalkAgent.GetAccesstoken()
		if err != nil {
			logs.Warn("获取钉钉临时Token失败 ->", err)
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}

		userid, err := dingtalkAgent.GetUserIDByUnionID(unionID)
		if err != nil {
			logs.Warn("获取钉钉用户ID失败 ->", err)
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}

		username, avatar, err := dingtalkAgent.GetUserNameAndAvatarByUserID(userid)
		if err != nil {
			logs.Warn("获取钉钉用户信息失败 ->", err)
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}

		member, err := models.NewMember().TmpLogin(tmpReader)
		if err == nil {
			member.LastLoginTime = time.Now()
			_ = member.Update("last_login_time")
			member.Account = username
			if avatar != "" {
				member.Avatar = avatar
			}

			c.SetMember(*member)
			c.LoggedIn(false)
			c.StopRun()
		}
		c.Redirect(conf.URLFor("AccountController.Login"), 302)

	// 企业微信扫码登录
	case "workweixin":
		//

	default:
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		c.StopRun()
	}
}

// 登录成功后的操作，如重定向到原始请求页面
func (c *AccountController) LoggedIn(isPost bool) interface{} {

	turl := c.referer()

	if !isPost {
		c.Redirect(turl, 302)
		return nil
	} else {
		var data struct {
			TURL string `json:"url"`
		}
		data.TURL = turl
		return data
	}
}

// 用户注册
func (c *AccountController) Register() {
	c.TplName = "account/register.tpl"

	//如果用户登录了，则跳转到网站首页
	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
		c.Redirect(conf.URLFor("HomeController.Index"), 302)
	}
	// 如果没有开启用户注册
	if v, ok := c.Option["ENABLED_REGISTER"]; ok && !strings.EqualFold(v, "true") {
		c.Abort("404")
	}

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		email := c.GetString("email")
		captcha := c.GetString("code")

		if ok, err := regexp.MatchString(conf.RegexpAccount, account); account == "" || !ok || err != nil {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.username_invalid_format"))
		}
		if l := strings.Count(password1, ""); password1 == "" || l > 50 || l < 6 {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.password_length_invalid"))
		}
		if password1 != password2 {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.incorrect_confirm_password"))
		}
		if ok, err := regexp.MatchString(conf.RegexpEmail, email); !ok || err != nil || email == "" {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.email_invalid_format"))
		}
		// 如果开启了验证码
		if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
			v, ok := c.GetSession(conf.CaptchaSessionName).(string)
			if !ok || !strings.EqualFold(v, captcha) {
				c.JsonResult(6001, i18n.Tr(c.Lang, "message.captcha_wrong"))
			}
		}

		member := models.NewMember()

		if _, err := member.FindByAccount(account); err == nil && member.MemberId > 0 {
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.account_existed"))
		}

		member.Account = account
		member.Password = password1
		member.Role = conf.MemberGeneralRole
		member.Avatar = conf.GetDefaultAvatar()
		member.CreateAt = 0
		member.Email = email
		member.Status = 0
		if err := member.Add(); err != nil {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed_register"))
		}

		c.JsonResult(0, "ok", member)
	}
}

// 找回密码
func (c *AccountController) FindPassword() {
	c.TplName = "account/find_password_setp1.tpl"
	mailConf := conf.GetMailConfig()

	if c.Ctx.Input.IsPost() {

		email := c.GetString("email")
		captcha := c.GetString("code")

		if email == "" {
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.email_empty"))
		}
		if !mailConf.EnableMail {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.mail_service_not_enable"))
		}

		// 如果开启了验证码
		if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
			v, ok := c.GetSession(conf.CaptchaSessionName).(string)
			if !ok || !strings.EqualFold(v, captcha) {
				c.JsonResult(6001, i18n.Tr(c.Lang, "message.captcha_wrong"))
			}
		}

		member, err := models.NewMember().FindByFieldFirst("email", email)
		if err != nil {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.email_not_exist"))
		}
		if member == nil || member.Status != 0 {
			c.JsonResult(6007, i18n.Tr(c.Lang, "message.account_disable"))
		}
		if member == nil || member.AuthMethod == conf.AuthMethodLDAP {
			c.JsonResult(6011, i18n.Tr(c.Lang, "message.account_not_support_retrieval"))
		}

		count, err := models.NewMemberToken().FindSendCount(email, time.Now().Add(-1*time.Hour), time.Now())

		if err != nil {
			logs.Error(err)
			c.JsonResult(6008, i18n.Tr(c.Lang, "message.failed_send_mail"))
		}
		if count > mailConf.MailNumber {
			c.JsonResult(6008, i18n.Tr(c.Lang, "message.sent_too_many_times"))
		}

		memberToken := models.NewMemberToken()

		memberToken.Token = string(utils.Krand(32, utils.KC_RAND_KIND_ALL))
		memberToken.Email = email
		memberToken.MemberId = member.MemberId
		memberToken.IsValid = false
		if _, err := memberToken.InsertOrUpdate(); err != nil {
			c.JsonResult(6009, i18n.Tr(c.Lang, "message.failed_send_mail"))
		}

		data := map[string]interface{}{
			"SITE_NAME": c.Option["SITE_NAME"],
			"url":       conf.URLFor("AccountController.FindPassword", "token", memberToken.Token, "mail", email),
			"BaseUrl":   c.BaseUrl(),
		}

		body, err := c.ExecuteViewPathTemplate("account/mail_template.tpl", data)
		if err != nil {
			logs.Error(err)
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed_send_mail"))
		}

		go func(mailConf *conf.SmtpConf, email string, body string) {

			mailConfig := &mail.SMTPConfig{
				Username: mailConf.SmtpUserName,
				Password: mailConf.SmtpPassword,
				Host:     mailConf.SmtpHost,
				Port:     mailConf.SmtpPort,
				Secure:   mailConf.Secure,
				Identity: "",
			}
			logs.Info(mailConfig)

			c := mail.NewSMTPClient(mailConfig)
			m := mail.NewMail()

			m.AddFrom(mailConf.FormUserName)
			m.AddFromName(mailConf.FormUserName)
			m.AddSubject("找回密码")
			m.AddHTML(body)
			m.AddTo(email)

			if e := c.Send(m); e != nil {
				logs.Error("发送邮件失败：" + e.Error())
			} else {
				logs.Info("邮件发送成功：" + email)
			}
			//auth := smtp.PlainAuth(
			//	"",
			//	mail_conf.SmtpUserName,
			//	mail_conf.SmtpPassword,
			//	mail_conf.SmtpHost,
			//)
			//
			//mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
			//subject := "Subject: 找回密码!\n"
			//
			//err = smtp.SendMail(
			//	mail_conf.SmtpHost+":"+strconv.Itoa(mail_conf.SmtpPort),
			//	auth,
			//	mail_conf.FormUserName,
			//	[]string{email},
			//	[]byte(subject+mime+"\n"+body),
			//)
			//if err != nil {
			//	logs.Error("邮件发送失败 => ", email, err)
			//}
		}(mailConf, email, body)

		c.JsonResult(0, "ok", conf.URLFor("AccountController.Login"))
	}

	token := c.GetString("token")
	email := c.GetString("mail")

	if token != "" && email != "" {
		memberToken, err := models.NewMemberToken().FindByFieldFirst("token", token)
		if err != nil {
			logs.Error(err)
			c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.mail_expired")
			c.TplName = "errors/error.tpl"
			return
		}
		subTime := time.Until(memberToken.SendTime)

		if !strings.EqualFold(memberToken.Email, email) || subTime.Minutes() > float64(mailConf.MailExpired) || !memberToken.ValidTime.IsZero() {
			c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.captcha_expired")
			c.TplName = "errors/error.tpl"
			return
		}
		c.Data["Email"] = memberToken.Email
		c.Data["Token"] = memberToken.Token
		c.TplName = "account/find_password_setp2.tpl"

	}
}

// 校验邮件并修改密码
func (c *AccountController) ValidEmail() {
	c.Prepare()
	password1 := c.GetString("password1")
	password2 := c.GetString("password2")
	captcha := c.GetString("code")
	token := c.GetString("token")
	email := c.GetString("mail")

	if password1 == "" {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.password_empty"))
	}
	if l := strings.Count(password1, ""); l < 6 || l > 50 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.password_length_invalid"))
	}
	if password2 == "" {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.confirm_password_empty"))
	}
	if password1 != password2 {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.incorrect_confirm_password"))
	}
	if captcha == "" {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.captcha_empty"))
	}
	v, ok := c.GetSession(conf.CaptchaSessionName).(string)
	if !ok || !strings.EqualFold(v, captcha) {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.captcha_wrong"))
	}

	mailConf := conf.GetMailConfig()
	memberToken, err := models.NewMemberToken().FindByFieldFirst("token", token)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6007, i18n.Tr(c.Lang, "message.mail_expired"))
	}
	subTime := time.Until(memberToken.SendTime)

	if !strings.EqualFold(memberToken.Email, email) || subTime.Minutes() > float64(mailConf.MailExpired) || !memberToken.ValidTime.IsZero() {

		c.JsonResult(6008, i18n.Tr(c.Lang, "message.captcha_expired"))
	}
	member, err := models.NewMember().Find(memberToken.MemberId)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	hash, err := utils.PasswordHash(password1)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed_save_password"))
	}

	member.Password = hash

	err = member.Update("password")
	memberToken.ValidTime = time.Now()
	memberToken.IsValid = true
	memberToken.InsertOrUpdate()

	if err != nil {
		logs.Error(err)
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed_save_password"))
	}
	c.JsonResult(0, "ok", conf.URLFor("AccountController.Login"))
}

// Logout 退出登录
func (c *AccountController) Logout() {
	c.SetMember(models.Member{})
	c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
	u := c.Ctx.Request.Header.Get("Referer")
	c.Redirect(conf.URLFor("AccountController.Login", "url", u), 302)
}

// 验证码
func (c *AccountController) Captcha() {
	c.Prepare()

	captchaImage := gocaptcha.NewCaptchaImage(140, 40, gocaptcha.RandLightColor())

	captchaImage.DrawNoise(gocaptcha.CaptchaComplexLower)

	// captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexHigh)
	txt := gocaptcha.RandText(4)

	c.SetSession(conf.CaptchaSessionName, txt)

	captchaImage.DrawText(txt)
	// captchaImage.Drawline(3);
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))
	// captchaImage.DrawHollowLine()

	captchaImage.SaveImage(c.Ctx.ResponseWriter, gocaptcha.ImageFormatJpeg)
	c.StopRun()
}
