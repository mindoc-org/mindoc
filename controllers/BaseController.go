package controllers

import (
	"bytes"

	"encoding/json"
	"io"
	"strings"
	"time"

	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
)

type BaseController struct {
	web.Controller
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
	}
	conf.BaseUrl = c.BaseUrl()
	c.Data["BaseUrl"] = c.BaseUrl()

	if options, err := models.NewOption().All(); err == nil {
		c.Option = make(map[string]string, len(options))
		for _, item := range options {
			c.Data[item.OptionName] = item.OptionValue
			c.Option[item.OptionName] = item.OptionValue
		}
		c.EnableAnonymous = strings.EqualFold(c.Option["ENABLE_ANONYMOUS"], "true")
		c.EnableDocumentHistory = strings.EqualFold(c.Option["ENABLE_DOCUMENT_HISTORY"], "true")
	}
	c.Data["HighlightStyle"] = web.AppConfig.DefaultString("highlight_style", "github")

	if b, err := ioutil.ReadFile(filepath.Join(web.BConfig.WebConfig.ViewsPath, "widgets", "scripts.tpl")); err == nil {
		c.Data["Scripts"] = template.HTML(string(b))
	}

}

//判断用户是否登录.
func (c *BaseController) isUserLoggedIn() bool {
	return c.Member != nil && c.Member.MemberId > 0
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
		logs.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	c.StopRun()
}

//如果错误不为空，则响应错误信息到浏览器.
func (c *BaseController) CheckJsonError(code int, err error) {

	if err == nil {
		return
	}
	jsonData := make(map[string]interface{}, 3)

	jsonData["errcode"] = code
	jsonData["message"] = err.Error()

	returnJSON, err := json.Marshal(jsonData)

	if err != nil {
		logs.Error(err)
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
		viewPath = web.BConfig.WebConfig.ViewsPath

	}

	if err := adapter.ExecuteViewPathTemplate(&buf, tplName, viewPath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *BaseController) BaseUrl() string {
	baseUrl := web.AppConfig.DefaultString("baseurl", "")
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

	if err := adapter.ExecuteViewPathTemplate(&buf, "errors/error.tpl", web.BConfig.WebConfig.ViewsPath, map[string]interface{}{"ErrorMessage": errMsg, "ErrorCode": errCode, "BaseUrl": conf.BaseUrl}); err != nil {
		c.Abort("500")
	}
	if errCode >= 200 && errCode <= 510 {
		c.CustomAbort(errCode, buf.String())
	} else {
		c.CustomAbort(500, buf.String())
	}
}

func (c *BaseController) CheckErrorResult(code int, err error) {
	if err != nil {
		c.ShowErrorPage(code, err.Error())
	}
}
