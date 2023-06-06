package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mindoc-org/mindoc/cache"
	"github.com/mindoc-org/mindoc/utils/auth2"
	"github.com/mindoc-org/mindoc/utils/auth2/dingtalk"
	"github.com/mindoc-org/mindoc/utils/auth2/wecom"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
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
)

const (
	SessionUserInfoKey  = "session-user-info-key"
	AccessTokenCacheKey = "access-token-cache-key"
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

func (c *AccountController) IsInWorkWeixin() bool {
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
	c.Data["CanLoginDingTalk"] = len(web.AppConfig.DefaultString("dingtalk_app_key", "")) > 0

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
		return
	}

	referer := c.referer()
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
	c.Data["url"] = referer

	auth2Redirect := "AccountController.Auth2Redirect"
	if can, _ := c.Data["CanLoginWorkWeixin"].(bool); can {
		c.Data["workweixin_login_url"] = conf.URLFor(auth2Redirect, ":app", wecom.AppName, "url", url.PathEscape(u))
	}

	if can, _ := c.Data["CanLoginDingTalk"].(bool); can {
		c.Data["dingtalk_login_url"] = conf.URLFor(auth2Redirect, ":app", dingtalk.AppName, "url", url.PathEscape(u))

	}
	return
}

/*
Auth2.0 第三方对接思路:
1. Auth2Redirect: 点击相应第三方接口，路由重定向至第三方提供的Auth2.0地址
2. Auth2Callback: 第三方回调处理，接收回调的授权码，并获取用户信息
	已绑定: 则读取用户信息，直接登录
	未绑定: 则弹窗提示（需要敏感信息）
		a) Auth2BindAccount: 绑定已有账户（用户名+密码）
		b) Auth2AutoAccount: 自动创建账户，以第三方用户ID作为用户名，密码123456。
							 用该方式创建的账户，无法使用账号密码登录，需要修改一次密码后才可以进行账号密码登录。
*/

func (c *AccountController) getAuth2Client() (auth2.Client, error) {
	app := c.Ctx.Input.Param(":app")
	var client auth2.Client
	tokenKey := AccessTokenCacheKey + "-" + app

	switch app {
	case wecom.AppName:
		if can, _ := c.Data["CanLoginWorkWeixin"].(bool); !can {
			return nil, errors.New("auth2.client.wecom.disabled")
		}
		corpId, _ := web.AppConfig.String("workweixin_corpid")
		agentId, _ := web.AppConfig.String("workweixin_agentid")
		secret, _ := web.AppConfig.String("workweixin_secret")
		client = wecom.NewClient(corpId, agentId, secret)

	case dingtalk.AppName:
		if can, _ := c.Data["CanLoginDingTalk"].(bool); !can {
			return nil, errors.New("auth2.client.dingtalk.disabled")
		}

		appKey, _ := web.AppConfig.String("dingtalk_app_key")
		appSecret, _ := web.AppConfig.String("dingtalk_app_secret")
		client = dingtalk.NewClient(appSecret, appKey)

	default:
		return nil, errors.New("auth2.client.notsupported")
	}

	var tokenCache auth2.AccessTokenCache
	err := cache.Get(tokenKey, &tokenCache)
	if err != nil {
		logs.Info("AccessToken从缓存读取失败")
		token, err := client.GetAccessToken(context.Background())
		if err != nil {
			return client, nil
		}
		tokenCache = auth2.NewAccessToken(token)
		cache.Put(tokenKey, tokenCache, tokenCache.GetExpireIn())
	}

	// 处理过期Token
	if tokenCache.IsExpired() {
		token, err := client.GetAccessToken(context.Background())
		if err != nil {
			return client, nil
		}
		tokenCache = auth2.NewAccessToken(token)
		cache.Put(tokenKey, tokenCache, tokenCache.GetExpireIn())
	}

	client.SetAccessToken(tokenCache)
	return client, nil
}

func (c *AccountController) parseAuth2CallbackParam() (code, state string) {
	switch c.Ctx.Input.Param(":app") {
	case wecom.AppName:
		code = c.GetString("code")
		state = c.GetString("state")
	case dingtalk.AppName:
		code = c.GetString("authCode")
		state = c.GetString("state")
	}

	logs.Debug("code: ", code)
	logs.Debug("state: ", state)
	return
}

func (c *AccountController) getAuth2Account() (models.Auth2Account, error) {
	switch c.Ctx.Input.Param(":app") {
	case wecom.AppName:
		return models.NewWorkWeixinAccount(), nil

	case dingtalk.AppName:
		return models.NewDingTalkAccount(), nil
	}

	return nil, errors.New("auth2.account.notsupported")
}

// Auth2Redirect 第三方auth2.0登录: 钉钉、企业微信
func (c *AccountController) Auth2Redirect() {
	client, err := c.getAuth2Client()
	if err != nil {
		c.DelSession(conf.LoginSessionName)
		c.SetMember(models.Member{})
		c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		c.StopRun()
		return
	}

	app := c.Ctx.Input.Param(":app")
	var isAppBrowser bool
	switch app {
	case wecom.AppName:
		isAppBrowser = c.IsInWorkWeixin()
	}

	var callback string
	u := c.GetString("url")
	if u == "" {
		u = c.referer()
		callback = conf.URLFor("AccountController.Auth2Callback", ":app", app)
	}
	if u != "" {
		var schemaRule = regexp.MustCompile(`^https?\:\/\/`)
		if !schemaRule.MatchString(u) {
			u = strings.TrimRight(conf.BaseUrl, "/") + strings.TrimLeft(u, "/")
		}
		callback = conf.URLFor("AccountController.Auth2Callback", ":app", app, "url", url.PathEscape(u))
	}

	logs.Debug("callback: ", callback) // debug
	c.Redirect(client.BuildURL(callback, isAppBrowser), http.StatusFound)
}

// Auth2Callback 第三方auth2.0回调
func (c *AccountController) Auth2Callback() {
	client, err := c.getAuth2Client()
	if err != nil {
		c.DelSession(conf.LoginSessionName)
		c.SetMember(models.Member{})
		c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		c.StopRun()
		logs.Error(err)
		return
	}

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

	c.TplName = "account/auth2_callback.tpl"
	bindExisted := "false"
	errMsg := ""
	userInfoJson := "{}"
	defer func() {
		c.Data["bind_existed"] = template.JS(bindExisted)
		logs.Debug("bind_existed: ", bindExisted)
		c.Data["error_msg"] = template.JS(errMsg)
		c.Data["user_info_json"] = template.JS(userInfoJson)
		c.Data["app"] = template.JS(c.Ctx.Input.Param(":app"))
	}()

	// 请求参数获取
	code, state := c.parseAuth2CallbackParam()
	if err := client.ValidateCallback(state); err != nil {
		c.DelSession(conf.LoginSessionName)
		c.SetMember(models.Member{})
		c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		errMsg = err.Error()
		logs.Error(err)
		return
	}

	userInfo, err := client.GetUserInfo(context.Background(), code)
	if err != nil {
		c.DelSession(conf.LoginSessionName)
		c.SetMember(models.Member{})
		c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
		errMsg = err.Error()
		logs.Error(err)
		return
	}

	account, err := c.getAuth2Account()
	if err != nil {
		logs.Error("获取Auth2用户失败 ->", err)
		c.JsonResult(500, "不支持的第三方用户", nil)
		return
	}

	member, err := account.ExistedMember(userInfo.UserId)
	if err != nil {
		if err == orm.ErrNoRows {
			if userInfo.Mobile == "" {
				errMsg = "请到应用浏览器中登录，并授权获取敏感信息。"
			} else {
				jsonInfo, _ := json.Marshal(userInfo)
				userInfoJson = string(jsonInfo)
				errMsg = ""
				c.SetSession(SessionUserInfoKey, userInfo)
			}
		} else {
			logs.Error("Error: ", err)
			errMsg = "登录错误: " + err.Error()
		}
		return
	}

	bindExisted = "true"
	errMsg = ""

	member.LastLoginTime = time.Now()
	_ = member.Update("last_login_time")

	c.SetMember(*member)
	remember.MemberId = member.MemberId
	remember.Account = member.Account
	remember.Time = time.Now()
	v, err := utils.Encode(remember)
	if err == nil {
		c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
	}
	u := c.GetString("url")
	if u == "" {
		u = conf.URLFor("HomeController.Index")
	}
	c.Redirect(u, 302)

}

// Auth2BindAccount 第三方auth2.0绑定已有账号
func (c *AccountController) Auth2BindAccount() {
	userInfo, ok := c.GetSession(SessionUserInfoKey).(auth2.UserInfo)
	if !ok || len(userInfo.UserId) <= 0 {
		c.DelSession(SessionUserInfoKey)
		c.JsonResult(400, "请求错误, 请从首页重新登录")
		return
	}

	account := c.GetString("account")
	password := c.GetString("password")
	if account == "" || password == "" {
		c.JsonResult(400, "账号或密码不能为空")
		return
	}

	member, err := models.NewMember().Login(account, password)
	if err != nil {
		logs.Error("用户登录 ->", err)
		c.JsonResult(500, "账号或密码错误", nil)
		return
	}

	bindAccount, err := c.getAuth2Account()
	if err != nil {
		logs.Error("获取Auth2用户失败 ->", err)
		c.JsonResult(500, "不支持的第三方用户", nil)
		return
	}

	member.CreateAt = 0
	ormer := orm.NewOrm()
	o, err := ormer.Begin()
	if err != nil {
		logs.Error("开启事务时出错 -> ", err)
		c.JsonResult(500, "开启事务时出错: ", err.Error())
		return
	}
	if err := bindAccount.AddBind(ormer, userInfo, member); err != nil {
		logs.Error(err)
		o.Rollback()
		c.JsonResult(500, "绑定失败，数据库错误: "+err.Error())
		return
	}

	// 绑定成功之后修改用户信息
	member.LastLoginTime = time.Now()
	//member.RealName = user_info.Name
	//member.Avatar = user_info.Avatar
	if len(member.Avatar) < 1 {
		member.Avatar = conf.GetDefaultAvatar()
	}
	//member.Email = user_info.Email
	//member.Phone = user_info.Mobile
	if _, err := ormer.Update(member, "last_login_time", "real_name", "avatar", "email", "phone"); err != nil {
		o.Rollback()
		logs.Error("保存用户信息失败=>", err)
		c.JsonResult(500, "绑定失败，现有账户信息更新失败: "+err.Error())
		return

	}

	if err := o.Commit(); err != nil {
		logs.Error("开启事务时出错 -> ", err)
		c.JsonResult(500, "开启事务时出错: ", err.Error())
		return
	}

	c.DelSession(SessionUserInfoKey)
	c.SetMember(*member)

	var remember CookieRemember
	remember.MemberId = member.MemberId
	remember.Account = member.Account
	remember.Time = time.Now()
	v, err := utils.Encode(remember)
	if err != nil {
		c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
		return
	}

	c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
	c.JsonResult(0, "绑定成功", nil)
}

// Auth2AutoAccount auth2.0自动创建账号
func (c *AccountController) Auth2AutoAccount() {
	app := c.Ctx.Input.Param(":app")
	logs.Debug("app: ", app)

	userInfo, ok := c.GetSession(SessionUserInfoKey).(auth2.UserInfo)
	if !ok || len(userInfo.UserId) <= 0 {
		c.DelSession(SessionUserInfoKey)
		c.JsonResult(400, "请求错误, 请从首页重新登录")
		return
	}

	c.DelSession(SessionUserInfoKey)
	member := models.NewMember()

	if _, err := member.FindByAccount(userInfo.UserId); err == nil && member.MemberId > 0 {
		c.JsonResult(400, "账号已存在")
		return
	}

	ormer := orm.NewOrm()
	o, err := ormer.Begin()
	if err != nil {
		logs.Error("开启事务时出错 -> ", err)
		c.JsonResult(500, "开启事务时出错: ", err.Error())
		return
	}

	member.Account = userInfo.UserId
	member.RealName = userInfo.Name
	member.Password = "123456" // 强制设置默认密码，需修改一次密码后，才可以进行账号密码登录
	hash, err := utils.PasswordHash(member.Password)

	if err != nil {
		logs.Error("加密用户密码失败 =>", err)
		c.JsonResult(500, "加密用户密码失败"+err.Error())
		return
	}

	logs.Debug("member.Password: ", member.Password)
	logs.Debug("hash: ", hash)
	member.Password = hash

	member.Role = conf.MemberGeneralRole
	member.Avatar = userInfo.Avatar
	if len(member.Avatar) < 1 {
		member.Avatar = conf.GetDefaultAvatar()
	}
	member.CreateAt = 0
	member.Email = userInfo.Mail
	member.Phone = userInfo.Mobile
	member.Status = 0
	if _, err = ormer.Insert(member); err != nil {
		o.Rollback()
		c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
		return
	}

	account, err := c.getAuth2Account()
	if err != nil {
		logs.Error("获取Auth2用户失败 ->", err)
		c.JsonResult(500, "不支持的第三方用户", nil)
		return
	}

	member.CreateAt = 0
	if err := account.AddBind(ormer, userInfo, member); err != nil {
		logs.Error(err)
		o.Rollback()
		c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
		return
	}

	if err := o.Commit(); err != nil {
		logs.Error("提交事务时出错 -> ", err)
		c.JsonResult(500, "提交事务时出错: ", err.Error())
		return
	}

	member.LastLoginTime = time.Now()
	_ = member.Update("last_login_time")

	c.SetMember(*member)

	var remember CookieRemember
	remember.MemberId = member.MemberId
	remember.Account = member.Account
	remember.Time = time.Now()
	v, err := utils.Encode(remember)
	if err != nil {
		c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
		return
	}

	c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
	c.JsonResult(0, "绑定成功", nil)
}

// 钉钉登录
//func (c *AccountController) DingTalkLogin() {
//	code := c.GetString("dingtalk_code")
//	if code == "" {
//		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_obtain_user_info"), nil)
//	}
//
//	appKey, _ := web.AppConfig.String("dingtalk_app_key")
//	appSecret, _ := web.AppConfig.String("dingtalk_app_secret")
//	tmpReader, _ := web.AppConfig.String("dingtalk_tmp_reader")
//
//	if appKey == "" || appSecret == "" || tmpReader == "" {
//		c.JsonResult(500, i18n.Tr(c.Lang, "message.dingtalk_auto_login_not_enable"), nil)
//		c.StopRun()
//	}
//
//	dingtalkAgent := dingtalk.NewDingTalkAgent(appSecret, appKey)
//	err := dingtalkAgent.GetAccesstoken()
//	if err != nil {
//		logs.Warn("获取钉钉临时Token失败 ->", err)
//		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
//		c.StopRun()
//	}
//
//	userid, err := dingtalkAgent.GetUserIDByCode(code)
//	if err != nil {
//		logs.Warn("获取钉钉用户ID失败 ->", err)
//		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
//		c.StopRun()
//	}
//
//	username, avatar, err := dingtalkAgent.GetUserNameAndAvatarByUserID(userid)
//	if err != nil {
//		logs.Warn("获取钉钉用户信息失败 ->", err)
//		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_auto_login"), nil)
//		c.StopRun()
//	}
//
//	member, err := models.NewMember().TmpLogin(tmpReader)
//	if err == nil {
//		member.LastLoginTime = time.Now()
//		_ = member.Update("last_login_time")
//		member.Account = username
//		if avatar != "" {
//			member.Avatar = avatar
//		}
//
//		c.SetMember(*member)
//	}
//	c.JsonResult(0, "ok", username)
//}

// WorkWeixinLogin 用户企业微信登录
//func (c *AccountController) WorkWeixinLogin() {
//	logs.Info("UserAgent: ", c.Ctx.Input.UserAgent()) // debug
//
//	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
//		u := c.GetString("url")
//		if u == "" {
//			u = c.Ctx.Request.Header.Get("Referer")
//			if u == "" {
//				u = conf.URLFor("HomeController.Index")
//			}
//		}
//		// session自动登录时刷新session内容
//		member, err := models.NewMember().Find(member.MemberId)
//		if err != nil {
//			c.DelSession(conf.LoginSessionName)
//			c.SetMember(models.Member{})
//			c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
//		} else {
//			c.SetMember(*member)
//		}
//		c.Redirect(u, 302)
//	}
//	var remember CookieRemember
//	// 如果 Cookie 中存在登录信息
//	if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
//		if err := utils.Decode(cookie, &remember); err == nil {
//			if member, err := models.NewMember().Find(remember.MemberId); err == nil {
//				c.SetMember(*member)
//				c.LoggedIn(false)
//				c.StopRun()
//			}
//		}
//	}
//
//	if c.Ctx.Input.IsPost() {
//		// account := c.GetString("account")
//		// password := c.GetString("password")
//		// captcha := c.GetString("code")
//		// isRemember := c.GetString("is_remember")
//		c.JsonResult(400, "request method not allowed", nil)
//	} else {
//		var callback_u string
//		u := c.GetString("url")
//		if u == "" {
//			u = c.referer()
//		}
//		if u != "" {
//			var schemaRule = regexp.MustCompile(`^https?\:\/\/`)
//			if !schemaRule.MatchString(u) {
//				u = strings.TrimRight(conf.BaseUrl, "/") + strings.TrimLeft(u, "/")
//			}
//		}
//		if u == "" {
//			callback_u = conf.URLFor("AccountController.WorkWeixinLoginCallback")
//		} else {
//			callback_u = conf.URLFor("AccountController.WorkWeixinLoginCallback", "url", url.PathEscape(u))
//		}
//		logs.Info("callback_u: ", callback_u) // debug
//
//		state := "mindoc"
//		workweixinConf := conf.GetWorkWeixinConfig()
//		appid := workweixinConf.CorpId
//		agentid := workweixinConf.AgentId
//		var redirect_uri string
//
//		isInWorkWeixin := c.IsInWorkWeixin()
//		c.Data["IsInWorkWeixin"] = isInWorkWeixin
//		if isInWorkWeixin {
//			// 企业微信内-网页授权登录
//			urlFmt := "%s?appid=%s&agentid=%s&redirect_uri=%s&response_type=code&scope=snsapi_privateinfo&state=%s#wechat_redirect"
//			redirect_uri = fmt.Sprintf(urlFmt, WorkWeixin_AuthorizeUrlBase, appid, agentid, url.PathEscape(callback_u), state)
//		} else {
//			// 浏览器内-扫码授权登录
//			urlFmt := "%s?login_type=CorpApp&appid=%s&agentid=%s&redirect_uri=%s&state=%s"
//			redirect_uri = fmt.Sprintf(urlFmt, WorkWeixin_QRConnectUrlBase, appid, agentid, url.PathEscape(callback_u), state)
//		}
//		logs.Info("redirect_uri: ", redirect_uri) // debug
//		c.Redirect(redirect_uri, 302)
//	}
//}

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
//func (c *AccountController) WorkWeixinLoginCallback() {
//	c.TplName = "account/auth2_callback.tpl"
//
//	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {
//		u := c.GetString("url")
//		if u == "" {
//			u = c.Ctx.Request.Header.Get("Referer")
//		}
//		if u == "" {
//			u = conf.URLFor("HomeController.Index")
//		}
//		member, err := models.NewMember().Find(member.MemberId)
//		if err != nil {
//			c.DelSession(conf.LoginSessionName)
//			c.SetMember(models.Member{})
//			c.SetSecureCookie(conf.GetAppKey(), "login", "", -3600)
//		} else {
//			c.SetMember(*member)
//		}
//		c.Redirect(u, 302)
//	}
//
//	var remember CookieRemember
//	// 如果 Cookie 中存在登录信息
//	if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
//		if err := utils.Decode(cookie, &remember); err == nil {
//			if member, err := models.NewMember().Find(remember.MemberId); err == nil {
//				c.SetMember(*member)
//				c.LoggedIn(false)
//				c.StopRun()
//			}
//		}
//	}
//
//	// 请求参数获取
//	req_code := c.GetString("code")
//	logs.Warning("req_code: ", req_code)
//	req_state := c.GetString("state")
//	logs.Warning("req_state: ", req_state)
//	var user_info_json string
//	var error_msg string
//	var bind_existed string
//	if len(req_code) > 0 && req_state == "mindoc" {
//		// 获取当前应用的access_token
//		access_token, ok := workweixin.GetAccessToken()
//		if ok {
//			logs.Warning("access_token: ", access_token)
//			// 获取当前请求的userid
//			user_id, ticket, ok := workweixin.RequestUserId(access_token, req_code)
//			if ok {
//				logs.Warning("user_id: ", user_id)
//				// 查询系统现有数据，是否绑定了当前请求用户的企业微信
//				member, err := models.NewWorkWeixinAccount().ExistedMember(user_id)
//				if err == nil {
//					member.LastLoginTime = time.Now()
//					_ = member.Update("last_login_time")
//
//					c.SetMember(*member)
//
//					var remember CookieRemember
//					remember.MemberId = member.MemberId
//					remember.Account = member.Account
//					remember.Time = time.Now()
//					v, err := utils.Encode(remember)
//					if err == nil {
//						c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
//					}
//					bind_existed = "true"
//					error_msg = ""
//					u := c.GetString("url")
//					if u == "" {
//						u = conf.URLFor("HomeController.Index")
//					}
//					c.Redirect(u, 302)
//				} else if err == orm.ErrNoRows {
//					bind_existed = "false"
//					if ticket == "" {
//						error_msg = "请到企业微信中登录，并授权获取敏感信息。"
//					} else {
//						user_info, err := workweixin.RequestUserPrivateInfo(access_token, user_id, ticket)
//						if err != nil {
//							error_msg = "获取敏感信息错误: " + err.Error()
//						} else {
//							json_info, _ := json.Marshal(user_info)
//							user_info_json = string(json_info)
//							error_msg = ""
//							c.SetSession(SessionUserInfoKey, user_info)
//						}
//					}
//				} else {
//					logs.Error("Error: ", err)
//					error_msg = "登录错误: " + err.Error()
//				}
//			} else {
//				error_msg = "获取用户Id失败: " + user_id
//			}
//		} else {
//			error_msg = "应用凭据获取失败: " + access_token
//		}
//	} else {
//		error_msg = "参数错误"
//	}
//	if user_info_json == "" {
//		user_info_json = "{}"
//	}
//	if bind_existed == "" {
//		bind_existed = "null"
//	}
//	// refer & doc:
//	// - https://golang.org/pkg/html/template/#HTML
//	// - https://stackoverflow.com/questions/24411880/go-html-templates-can-i-stop-the-templates-package-inserting-quotes-around-stri
//	// - https://stackoverflow.com/questions/38035176/insert-javascript-snippet-inside-template-with-beego-golang
//	c.Data["bind_existed"] = template.JS(bind_existed)
//	logs.Debug("bind_existed: ", bind_existed)
//	c.Data["error_msg"] = template.JS(error_msg)
//	c.Data["user_info_json"] = template.JS(user_info_json)
//	/*
//		// 调试: 显示源码
//		result, err := c.RenderString()
//		if err != nil {
//			logs.Error(err)
//		} else {
//			logs.Warning(result)
//		}
//	*/
//}

// WorkWeixinLoginBind 用户企业微信登录-绑定
//func (c *AccountController) WorkWeixinLoginBind() {
//	if user_info, ok := c.GetSession(SessionUserInfoKey).(workweixin.WorkWeixinUserPrivateInfo); ok && len(user_info.UserId) > 0 {
//		req_account := c.GetString("account")
//		req_password := c.GetString("password")
//		if req_account == "" || req_password == "" {
//			c.JsonResult(400, "账号或密码不能为空")
//		} else {
//			member, err := models.NewMember().Login(req_account, req_password)
//			if err == nil {
//				account := models.NewWorkWeixinAccount()
//				account.MemberId = member.MemberId
//				account.WorkWeixin_UserId = user_info.UserId
//				member.CreateAt = 0
//				ormer := orm.NewOrm()
//				o, err := ormer.Begin()
//				if err != nil {
//					logs.Error("开启事务时出错 -> ", err)
//					c.JsonResult(500, "开启事务时出错: ", err.Error())
//				}
//				if err := account.AddBind(ormer); err != nil {
//					o.Rollback()
//					c.JsonResult(500, "绑定失败，数据库错误: "+err.Error())
//				} else {
//					// 绑定成功之后修改用户信息
//					member.LastLoginTime = time.Now()
//					//member.RealName = user_info.Name
//					//member.Avatar = user_info.Avatar
//					if len(member.Avatar) < 1 {
//						member.Avatar = conf.GetDefaultAvatar()
//					}
//					//member.Email = user_info.Email
//					//member.Phone = user_info.Mobile
//					if _, err := ormer.Update(member, "last_login_time", "real_name", "avatar", "email", "phone"); err != nil {
//						o.Rollback()
//						logs.Error("保存用户信息失败=>", err)
//						c.JsonResult(500, "绑定失败，现有账户信息更新失败: "+err.Error())
//					} else {
//						if err := o.Commit(); err != nil {
//							logs.Error("开启事务时出错 -> ", err)
//							c.JsonResult(500, "开启事务时出错: ", err.Error())
//						} else {
//							c.DelSession(SessionUserInfoKey)
//							c.SetMember(*member)
//
//							var remember CookieRemember
//							remember.MemberId = member.MemberId
//							remember.Account = member.Account
//							remember.Time = time.Now()
//							v, err := utils.Encode(remember)
//							if err == nil {
//								c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
//								c.JsonResult(0, "绑定成功", nil)
//							} else {
//								c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
//							}
//						}
//					}
//
//				}
//
//			} else {
//				logs.Error("用户登录 ->", err)
//				c.JsonResult(500, "账号或密码错误", nil)
//			}
//			c.JsonResult(500, "TODO: 绑定以后账号功能开发中")
//		}
//	} else {
//		if ok {
//			c.DelSession(SessionUserInfoKey)
//		}
//		c.JsonResult(400, "请求错误, 请从首页重新登录")
//	}
//
//}

// WorkWeixinLoginIgnore 用户企业微信登录-忽略
//func (c *AccountController) WorkWeixinLoginIgnore() {
//	if user_info, ok := c.GetSession(SessionUserInfoKey).(workweixin.WorkWeixinUserPrivateInfo); ok && len(user_info.UserId) > 0 {
//		c.DelSession(SessionUserInfoKey)
//		member := models.NewMember()
//
//		if _, err := member.FindByAccount(user_info.UserId); err == nil && member.MemberId > 0 {
//			c.JsonResult(400, "账号已存在")
//		}
//
//		ormer := orm.NewOrm()
//		o, err := ormer.Begin()
//		if err != nil {
//			logs.Error("开启事务时出错 -> ", err)
//			c.JsonResult(500, "开启事务时出错: ", err.Error())
//		}
//
//		member.Account = user_info.UserId
//		member.RealName = user_info.Name
//		var rnd = rand.New(src)
//		// fmt.Sprintf("%x", rnd.Uint64())
//		// strconv.FormatUint(rnd.Uint64(), 16)
//		member.Password = user_info.UserId + strconv.FormatUint(rnd.Uint64(), 16)
//		member.Password = "123456" // 强制设置默认密码，需修改一次密码后，才可以进行账号密码登录
//		hash, err := utils.PasswordHash(member.Password)
//		if err != nil {
//			logs.Error("加密用户密码失败 =>", err)
//			c.JsonResult(500, "加密用户密码失败"+err.Error())
//		} else {
//			logs.Error("member.Password: ", member.Password)
//			logs.Error("hash: ", hash)
//			member.Password = hash
//		}
//		member.Role = conf.MemberGeneralRole
//		member.Avatar = user_info.Avatar
//		if len(member.Avatar) < 1 {
//			member.Avatar = conf.GetDefaultAvatar()
//		}
//		member.CreateAt = 0
//		member.Email = user_info.BizMail
//		member.Phone = user_info.Mobile
//		member.Status = 0
//		if _, err = ormer.Insert(member); err != nil {
//			o.Rollback()
//			c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
//		} else {
//			account := models.NewWorkWeixinAccount()
//			account.MemberId = member.MemberId
//			account.WorkWeixin_UserId = user_info.UserId
//			member.CreateAt = 0
//			if err := account.AddBind(ormer); err != nil {
//				o.Rollback()
//				c.JsonResult(500, "注册失败，数据库错误: "+err.Error())
//			} else {
//				if err := o.Commit(); err != nil {
//					logs.Error("提交事务时出错 -> ", err)
//					c.JsonResult(500, "提交事务时出错: ", err.Error())
//				} else {
//					member.LastLoginTime = time.Now()
//					_ = member.Update("last_login_time")
//
//					c.SetMember(*member)
//
//					var remember CookieRemember
//					remember.MemberId = member.MemberId
//					remember.Account = member.Account
//					remember.Time = time.Now()
//					v, err := utils.Encode(remember)
//					if err == nil {
//						c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour*24*30*5).Unix())
//						c.JsonResult(0, "绑定成功", nil)
//					} else {
//						c.JsonResult(500, "绑定成功, 但自动登录失败, 请返回首页重新登录", nil)
//					}
//				}
//			}
//		}
//	} else {
//		if ok {
//			c.DelSession(SessionUserInfoKey)
//		}
//		c.JsonResult(400, "请求错误, 请从首页重新登录")
//	}
//}

// QR二维码登录
//func (c *AccountController) QRLogin() {
//	appName := c.Ctx.Input.Param(":app")
//
//	switch appName {
//	// 钉钉扫码登录
//	case "dingtalk":
//		code := c.GetString("code")
//		state := c.GetString("state")
//		if state != "1" || code == "" {
//			c.Redirect(conf.URLFor("AccountController.Login"), 302)
//			c.StopRun()
//		}
//		appKey, _ := web.AppConfig.String("dingtalk_qr_key")
//		appSecret, _ := web.AppConfig.String("dingtalk_qr_secret")
//
//		qrDingtalk := dingtalk.NewDingtalkQRLogin(appSecret, appKey)
//		unionID, err := qrDingtalk.GetUnionIDByCode(code)
//		if err != nil {
//			logs.Warn("获取钉钉临时UnionID失败 ->", err)
//			c.Redirect(conf.URLFor("AccountController.Login"), 302)
//			c.StopRun()
//		}
//
//		appKey, _ = web.AppConfig.String("dingtalk_app_key")
//		appSecret, _ = web.AppConfig.String("dingtalk_app_secret")
//		tmpReader, _ := web.AppConfig.String("dingtalk_tmp_reader")
//
//		dingtalkAgent := dingtalk.NewDingTalkAgent(appSecret, appKey)
//		err = dingtalkAgent.GetAccesstoken()
//		if err != nil {
//			logs.Warn("获取钉钉临时Token失败 ->", err)
//			c.Redirect(conf.URLFor("AccountController.Login"), 302)
//			c.StopRun()
//		}
//
//		userid, err := dingtalkAgent.GetUserIDByUnionID(unionID)
//		if err != nil {
//			logs.Warn("获取钉钉用户ID失败 ->", err)
//			c.Redirect(conf.URLFor("AccountController.Login"), 302)
//			c.StopRun()
//		}
//
//		username, avatar, err := dingtalkAgent.GetUserNameAndAvatarByUserID(userid)
//		if err != nil {
//			logs.Warn("获取钉钉用户信息失败 ->", err)
//			c.Redirect(conf.URLFor("AccountController.Login"), 302)
//			c.StopRun()
//		}
//
//		member, err := models.NewMember().TmpLogin(tmpReader)
//		if err == nil {
//			member.LastLoginTime = time.Now()
//			_ = member.Update("last_login_time")
//			member.Account = username
//			if avatar != "" {
//				member.Avatar = avatar
//			}
//
//			c.SetMember(*member)
//			c.LoggedIn(false)
//			c.StopRun()
//		}
//		c.Redirect(conf.URLFor("AccountController.Login"), 302)
//
//	// 企业微信扫码登录
//	case "workweixin":
//		//
//
//	default:
//		c.Redirect(conf.URLFor("AccountController.Login"), 302)
//		c.StopRun()
//	}
//}

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
