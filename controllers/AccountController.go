package controllers

import (
	"regexp"
	"strings"
	"time"
	"net/url"

	"github.com/lifei6671/mindoc/mail"
	"github.com/astaxie/beego"
	"github.com/lifei6671/gocaptcha"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
	"html/template"
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
	c.EnableXSRF = true
	c.Data["xsrfdata"]=template.HTML(c.XSRFFormHTML())
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
				c.JsonResult(403,"非法请求")
			} else {
				c.ShowErrorPage(403, "非法请求")
			}
		}
		xsrfToken := c.XSRFToken()
		if xsrfToken != token {
			if c.IsAjax() {
				c.JsonResult(403,"非法请求")
			} else {
				c.ShowErrorPage(403, "非法请求")
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
				c.JsonResult(6001, "验证码不正确")
			}
		}

		if account == "" || password == "" {
			c.JsonResult(6002, "账号或密码不能为空")
		}

		member, err := models.NewMember().Login(account, password)
		if err == nil {
			member.LastLoginTime = time.Now()
			member.Update()

			c.SetMember(*member)

			if strings.EqualFold(isRemember, "yes") {
				remember.MemberId = member.MemberId
				remember.Account = member.Account
				remember.Time = time.Now()
				v, err := utils.Encode(remember)
				if err == nil {
					c.SetSecureCookie(conf.GetAppKey(), "login", v, time.Now().Add(time.Hour * 24 * 30).Unix())
				}
			}

			c.JsonResult(0, "ok", c.referer())
		} else {
			beego.Error("用户登录 ->", err)
			c.JsonResult(500, "账号或密码错误", nil)
		}
	} else {
		c.Data["url"] = c.referer()
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
			c.JsonResult(6001, "账号只能由英文字母数字组成，且在3-50个字符")
		}
		if l := strings.Count(password1, ""); password1 == "" || l > 50 || l < 6 {
			c.JsonResult(6002, "密码必须在6-50个字符之间")
		}
		if password1 != password2 {
			c.JsonResult(6003, "确认密码不正确")
		}
		if ok, err := regexp.MatchString(conf.RegexpEmail, email); !ok || err != nil || email == "" {
			c.JsonResult(6004, "邮箱格式不正确")
		}
		// 如果开启了验证码
		if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
			v, ok := c.GetSession(conf.CaptchaSessionName).(string)
			if !ok || !strings.EqualFold(v, captcha) {
				c.JsonResult(6001, "验证码不正确")
			}
		}

		member := models.NewMember()

		if _, err := member.FindByAccount(account); err == nil && member.MemberId > 0 {
			c.JsonResult(6005, "账号已存在")
		}

		member.Account = account
		member.Password = password1
		member.Role = conf.MemberGeneralRole
		member.Avatar = conf.GetDefaultAvatar()
		member.CreateAt = 0
		member.Email = email
		member.Status = 0
		if err := member.Add(); err != nil {
			c.JsonResult(6006, "注册失败，请联系系统管理员处理")
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
			c.JsonResult(6005, "邮箱地址不能为空")
		}
		if !mailConf.EnableMail {
			c.JsonResult(6004, "未启用邮件服务")
		}

		// 如果开启了验证码
		if v, ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v, "true") {
			v, ok := c.GetSession(conf.CaptchaSessionName).(string)
			if !ok || !strings.EqualFold(v, captcha) {
				c.JsonResult(6001, "验证码不正确")
			}
		}

		member, err := models.NewMember().FindByFieldFirst("email", email)
		if err != nil {
			c.JsonResult(6006, "邮箱不存在")
		}
		if member.Status != 0 {
			c.JsonResult(6007, "账号已被禁用")
		}
		if member.AuthMethod == conf.AuthMethodLDAP {
			c.JsonResult(6011, "当前用户不支持找回密码")
		}

		count, err := models.NewMemberToken().FindSendCount(email, time.Now().Add(-1*time.Hour), time.Now())

		if err != nil {
			beego.Error(err)
			c.JsonResult(6008, "发送邮件失败")
		}
		if count > mailConf.MailNumber {
			c.JsonResult(6008, "发送次数太多，请稍候再试")
		}

		memberToken := models.NewMemberToken()

		memberToken.Token = string(utils.Krand(32, utils.KC_RAND_KIND_ALL))
		memberToken.Email = email
		memberToken.MemberId = member.MemberId
		memberToken.IsValid = false
		if _, err := memberToken.InsertOrUpdate(); err != nil {
			c.JsonResult(6009, "邮件发送失败")
		}

		data := map[string]interface{}{
			"SITE_NAME": c.Option["SITE_NAME"],
			"url":       conf.URLFor("AccountController.FindPassword", "token", memberToken.Token, "mail", email),
			"BaseUrl":   c.BaseUrl(),
		}

		body, err := c.ExecuteViewPathTemplate("account/mail_template.tpl", data)
		if err != nil {
			beego.Error(err)
			c.JsonResult(6003, "邮件发送失败")
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
			beego.Info(mailConfig)

			c := mail.NewSMTPClient(mailConfig)
			m := mail.NewMail()

			m.AddFrom(mailConf.FormUserName)
			m.AddFromName(mailConf.FormUserName)
			m.AddSubject("找回密码")
			m.AddHTML(body)
			m.AddTo(email)

			if e := c.Send(m); e != nil {
				beego.Error("发送邮件失败：" + e.Error())
			} else {
				beego.Info("邮件发送成功：" + email)
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
			//	beego.Error("邮件发送失败 => ", email, err)
			//}
		}(mailConf, email, body)

		c.JsonResult(0, "ok", conf.URLFor("AccountController.Login"))
	}

	token := c.GetString("token")
	email := c.GetString("mail")

	if token != "" && email != "" {
		memberToken, err := models.NewMemberToken().FindByFieldFirst("token", token)

		if err != nil {
			beego.Error(err)
			c.Data["ErrorMessage"] = "邮件已失效"
			c.TplName = "errors/error.tpl"
			return
		}
		subTime := memberToken.SendTime.Sub(time.Now())

		if !strings.EqualFold(memberToken.Email, email) || subTime.Minutes() > float64(mailConf.MailExpired) || !memberToken.ValidTime.IsZero() {
			c.Data["ErrorMessage"] = "验证码已过期，请重新操作。"
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
		c.JsonResult(6001, "密码不能为空")
	}
	if l := strings.Count(password1, ""); l < 6 || l > 50 {
		c.JsonResult(6001, "密码不能为空且必须在6-50个字符之间")
	}
	if password2 == "" {
		c.JsonResult(6002, "确认密码不能为空")
	}
	if password1 != password2 {
		c.JsonResult(6003, "确认密码输入不正确")
	}
	if captcha == "" {
		c.JsonResult(6004, "验证码不能为空")
	}
	v, ok := c.GetSession(conf.CaptchaSessionName).(string)
	if !ok || !strings.EqualFold(v, captcha) {
		c.JsonResult(6001, "验证码不正确")
	}

	mailConf := conf.GetMailConfig()
	memberToken, err := models.NewMemberToken().FindByFieldFirst("token", token)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6007, "邮件已失效")
	}
	subTime := memberToken.SendTime.Sub(time.Now())

	if !strings.EqualFold(memberToken.Email, email) || subTime.Minutes() > float64(mailConf.MailExpired) || !memberToken.ValidTime.IsZero() {

		c.JsonResult(6008, "验证码已过期，请重新操作。")
	}
	member, err := models.NewMember().Find(memberToken.MemberId)
	if err != nil {
		beego.Error(err)
		c.JsonResult(6005, "用户不存在")
	}
	hash, err := utils.PasswordHash(password1)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6006, "保存密码失败")
	}

	member.Password = hash

	err = member.Update("password")
	memberToken.ValidTime = time.Now()
	memberToken.IsValid = true
	memberToken.InsertOrUpdate()

	if err != nil {
		beego.Error(err)
		c.JsonResult(6006, "保存密码失败")
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

	captchaImage, err := gocaptcha.NewCaptchaImage(140, 40, gocaptcha.RandLightColor())

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}

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
