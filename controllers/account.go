package controllers

import (
	"time"
	"strings"

	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/models"
	"github.com/lifei6671/godoc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/gocaptcha"

	"regexp"
)

// AccountController 用户登录与注册.
type AccountController struct {
	BaseController
}

// Login 用户登录.
func (c *AccountController) Login()  {
	c.Prepare()

	var remember struct { MemberId int ; Account string; Time time.Time}

	//如果Cookie中存在登录信息
	if cookie,ok := c.GetSecureCookie(conf.GetAppKey(),"login");ok{

		if err := utils.Decode(cookie,&remember); err == nil {
			member := models.NewMember()
			member.MemberId = remember.MemberId

			if err := models.NewMember().Find(remember.MemberId); err == nil {
				c.SetMember(*member)

				c.Redirect(beego.URLFor("HomeController.Index"), 302)
				c.StopRun()
			}
		}
	}

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password := c.GetString("password")
		captcha := c.GetString("code")
		is_remember := c.GetString("is_remember")

		//如果开启了验证码
		if v,ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v,"true") {
			v,ok := c.GetSession(conf.CaptchaSessionName).(string);
			if !ok || !strings.EqualFold(v,captcha){
				c.JsonResult(6001,"验证码不正确")
			}
		}
		member,err := models.NewMember().Login(account,password)

		//如果没有数据
		if err == nil {
			c.SetMember(*member)
			if strings.EqualFold(is_remember,"yes") {
				remember.MemberId = member.MemberId
				remember.Account = member.Account
				remember.Time = time.Now()
				v ,err := utils.Encode(remember)
				if err == nil {
					c.SetSecureCookie(conf.GetAppKey(),"login",v)
				}

			}
			c.JsonResult(0,"ok")
		}else{
			logs.Error("用户登录 =>",err)
			c.JsonResult(500,"账号或密码错误",nil)
		}

		return
	}else{

		c.Layout = ""
		c.TplName = "account/login.tpl"
	}
}

func (c *AccountController) Register()  {
	c.TplName = "account/register.tpl"

	//如果没有开启用户注册
	if v,ok := c.Option["ENABLED_REGISTER"]; ok && !strings.EqualFold(v,"true") {
		c.Abort("404")
	}

	if c.Ctx.Input.IsPost() {
		account := c.GetString("account")
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		email := c.GetString("email")
		captcha := c.GetString("code")

		if ok,err := regexp.MatchString(conf.RegexpAccount,account); account == "" || !ok || err != nil {
			c.JsonResult(6001,"账号只能由英文字母数字组成，且在3-50个字符")
		}
		if  l := strings.Count(password1,"") ; password1 == "" || l > 50 || l < 6{
			c.JsonResult(6002,"密码必须在6-50个字符之间")
		}
		if password1 != password2 {
			c.JsonResult(6003,"确认密码不正确")
		}
		if  ok,err := regexp.MatchString(conf.RegexpEmail,email); !ok || err != nil || email == "" {
			c.JsonResult(6004,"邮箱格式不正确")
		}
		//如果开启了验证码
		if v,ok := c.Option["ENABLED_CAPTCHA"]; ok && strings.EqualFold(v,"true") {
			v,ok := c.GetSession(conf.CaptchaSessionName).(string);
			if !ok || !strings.EqualFold(v,captcha){
				c.JsonResult(6001,"验证码不正确")
			}
		}

		member := models.NewMember()

		if _,err := member.FindByAccount(account); err == nil && member.MemberId > 0 {
			c.JsonResult(6005,"账号已存在")
		}

		member.Account = account
		member.Password = password1
		member.Role = conf.MemberGeneralRole
		member.Avatar = conf.GetDefaultAvatar()
		member.CreateAt = 0
		member.Email = email
		member.Status = 0
		if err := member.Add(); err != nil {
			beego.Error(err)
			c.JsonResult(6006,"注册失败，请联系系统管理员处理")
		}

		c.JsonResult(0,"ok",member)
	}
}

func (p *AccountController) FindPassword()  {
	p.TplName = "account/find_password.tpl"
}
// Logout 退出登录.
func (c *AccountController) Logout(){
	c.SetMember(models.Member{});

	c.Redirect(beego.URLFor("AccountController.Login"),302)
}

func (c *AccountController) Captcha()  {
	c.Prepare()

	captchaImage, err := gocaptcha.NewCaptchaImage(140, 40, gocaptcha.RandLightColor())

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}

	captchaImage.DrawNoise(gocaptcha.CaptchaComplexLower)

	//captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexHigh)
	txt := gocaptcha.RandText(4)

	c.SetSession(conf.CaptchaSessionName,txt)

	captchaImage.DrawText(txt)
	//captchaImage.Drawline(3);
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))
	//captchaImage.DrawHollowLine()


	captchaImage.SaveImage(c.Ctx.ResponseWriter, gocaptcha.ImageFormatJpeg)
	c.StopRun()
}