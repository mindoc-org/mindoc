package controllers

import (
	"bytes"

	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
)

type BaseController struct {
	beego.Controller
	Member                *models.Member
	Option                map[string]string
	EnableAnonymous       bool
	EnableDocumentHistory bool
}

type CookieRemember struct {
	MemberId int
	Account  string
	Time     time.Time
}

// Prepare 预处理.
func (c *BaseController) Prepare() {
	c.Data["SiteName"] = "MinDoc"
	c.Data["Member"] = models.NewMember()
	controller, action := c.GetControllerAndAction()

	c.Data["ActionName"] = action
	c.Data["ControllerName"] = controller

	c.EnableAnonymous = false
	c.EnableDocumentHistory = false

	if member, ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0 {

		c.Member = &member
		c.Data["Member"] = c.Member
	} else {
		var remember CookieRemember
		// //如果Cookie中存在登录信息，从cookie中获取用户信息
		if cookie, ok := c.GetSecureCookie(conf.GetAppKey(), "login"); ok {
			if err := utils.Decode(cookie, &remember); err == nil {
				if member, err := models.NewMember().Find(remember.MemberId); err == nil {
					c.Member = member
					c.Data["Member"] = member
					c.SetMember(*member)
				}
			}
		}
		//c.Member = models.NewMember()
		//c.Member.Find(1)
		//c.Data["Member"] = *c.Member
	}
	conf.BaseUrl = c.BaseUrl()
	c.Data["BaseUrl"] = c.BaseUrl()

	if options, err := models.NewOption().All(); err == nil {
		c.Option = make(map[string]string, len(options))
		for _, item := range options {
			c.Data[item.OptionName] = item.OptionValue
			c.Option[item.OptionName] = item.OptionValue

			if strings.EqualFold(item.OptionName, "ENABLE_ANONYMOUS") && item.OptionValue == "true" {
				c.EnableAnonymous = true
			}
			if strings.EqualFold(item.OptionName, "ENABLE_DOCUMENT_HISTORY") && item.OptionValue == "true" {
				c.EnableDocumentHistory = true
			}
		}
	}
	c.Data["HighlightStyle"] = beego.AppConfig.DefaultString("highlight_style","github")
}

// SetMember 获取或设置当前登录用户信息,如果 MemberId 小于 0 则标识删除 Session
func (c *BaseController) SetMember(member models.Member) {

	if member.MemberId <= 0 {
		c.DelSession(conf.LoginSessionName)
		c.DelSession("uid")
		c.DestroySession()
	} else {
		c.SetSession(conf.LoginSessionName, member)
		c.SetSession("uid", member.MemberId)
	}
}

// JsonResult 响应 json 结果
func (c *BaseController) JsonResult(errCode int, errMsg string, data ...interface{}) {
	jsonData := make(map[string]interface{}, 3)

	jsonData["errcode"] = errCode
	jsonData["message"] = errMsg

	if len(data) > 0 && data[0] != nil {
		jsonData["data"] = data[0]
	}

	returnJSON, err := json.Marshal(jsonData)

	if err != nil {
		beego.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	c.StopRun()
}

// ExecuteViewPathTemplate 执行指定的模板并返回执行结果.
func (c *BaseController) ExecuteViewPathTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer

	viewPath := c.ViewPath

	if c.ViewPath == "" {
		viewPath = beego.BConfig.WebConfig.ViewsPath

	}

	if err := beego.ExecuteViewPathTemplate(&buf, tplName, viewPath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *BaseController) BaseUrl() string {
	baseUrl := beego.AppConfig.DefaultString("baseurl", "")
	if baseUrl != "" {
		if strings.HasSuffix(baseUrl, "/") {
			baseUrl = strings.TrimSuffix(baseUrl, "/")
		}
	} else {
		baseUrl = c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
	}
	return baseUrl
}

//显示错误信息页面.
func (c *BaseController) ShowErrorPage(errCode int, errMsg string) {
	c.TplName = "errors/error.tpl"

	c.Data["ErrorMessage"] = errMsg
	c.Data["ErrorCode"] = errCode

	var buf bytes.Buffer

	if err := beego.ExecuteViewPathTemplate(&buf, "errors/error.tpl", beego.BConfig.WebConfig.ViewsPath, map[string]interface{}{"ErrorMessage": errMsg, "ErrorCode": errCode, "BaseUrl": conf.BaseUrl}); err != nil {
		c.Abort("500")
	}
	if errCode >= 200 && errCode <= 510 {
		c.CustomAbort(errCode, buf.String())
	}else{
		c.CustomAbort(500, buf.String())
	}
}
