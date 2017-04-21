package controllers

import (
	"fmt"
	"time"

	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/models"
	"github.com/lifei6671/godoc/utils"
	"github.com/astaxie/beego"
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
		account := c.GetString("inputAccount")
		password := c.GetString("inputPassword")

		member,err := models.NewMember().Login(account,password)

		//如果没有数据
		if err == nil {
			c.SetMember(*member)
			c.JsonResult(0,"ok")
			c.StopRun()
		}else{
			fmt.Println(err)
			c.JsonResult(500,"账号或密码错误",nil)
		}

		return
	}else{

		c.Layout = ""
		c.TplName = "account/login.tpl"
	}
}

func (p *AccountController) Register()  {
	p.TplName = "account/register.tpl"

	
}

func (p *AccountController) FindPassword()  {
	p.TplName = "account/find_password.tpl"
}
// Logout 退出登录.
func (c *AccountController) Logout(){
	c.SetMember(models.Member{});

	c.Redirect(beego.URLFor("AccountController.Login"),302)
}

