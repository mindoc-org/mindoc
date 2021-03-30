package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/beego/i18n"
	"net/url"
	"regexp"
	"strings"
	"time"

	"html/template"

	"github.com/astaxie/beego"
	"github.com/lifei6671/gocaptcha"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/mail"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/dingtalk"
)

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

func (c *AccountController) Prepare() {
	c.BaseController.Prepare()
	c.EnableXSRF = beego.AppConfig.DefaultBool("enablexsrf", true)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["corpID"] = beego.AppConfig.String("dingtalk_corpid")
	c.Data["ENABLE_QR_DINGTALK"] = (beego.AppConfig.String("dingtalk_corpid") != "")
	c.Data["dingtalk_qr_key"] = beego.AppConfig.String("dingtalk_qr_key")

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
		c.Data["url"] = c.referer()
	}
}

// 钉钉登录
func (c *AccountController) DingTalkLogin() {
	c.Prepare()

	code := c.GetString("dingtalk_code")
	if code == "" {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.failed_obtain_user_info"), nil)
	}

	appKey := beego.AppConfig.String("dingtalk_app_key")
	appSecret := beego.AppConfig.String("dingtalk_app_secret")
	tmpReader := beego.AppConfig.String("dingtalk_tmp_reader")

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
		appKey := beego.AppConfig.String("dingtalk_qr_key")
		appSecret := beego.AppConfig.String("dingtalk_qr_secret")

		qrDingtalk := dingtalk.NewDingtalkQRLogin(appSecret, appKey)
		unionID, err := qrDingtalk.GetUnionIDByCode(code)
		if err != nil {
			logs.Warn("获取钉钉临时UnionID失败 ->", err)
			c.Redirect(conf.URLFor("AccountController.Login"), 302)
			c.StopRun()
		}

		appKey = beego.AppConfig.String("dingtalk_app_key")
		appSecret = beego.AppConfig.String("dingtalk_app_secret")
		tmpReader := beego.AppConfig.String("dingtalk_tmp_reader")

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
		subTime := memberToken.SendTime.Sub(time.Now())

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
	subTime := memberToken.SendTime.Sub(time.Now())

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
