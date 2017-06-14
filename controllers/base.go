package controllers

import (

	"bytes"

	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego"
	"strings"
	"encoding/json"
	"io"
)


type BaseController struct {
	beego.Controller
	Member *models.Member
	Option map[string]string
	EnableAnonymous bool
	EnableDocumentHistory bool
}

// Prepare 预处理.
func (c *BaseController) Prepare (){
	c.Data["SiteName"] = "MinDoc"
	c.Data["Member"] = models.Member{}
	c.EnableAnonymous = false
	c.EnableDocumentHistory = false

	if member,ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0{
		c.Member = &member
		c.Data["Member"] = c.Member
	}else{
		//c.Member = models.NewMember()
		//c.Member.Find(1)
		//c.Data["Member"] = *c.Member
	}
	c.Data["BaseUrl"] = c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host

	if options,err := models.NewOption().All();err == nil {
		c.Option = make(map[string]string,len(options))
		for _,item := range options {
			c.Data[item.OptionName] = item.OptionValue
			c.Option[item.OptionName] = item.OptionValue
			if strings.EqualFold(item.OptionName,"ENABLE_ANONYMOUS") && item.OptionValue == "true" {
				c.EnableAnonymous = true
			}
			if strings.EqualFold(item.OptionName,"ENABLE_DOCUMENT_HISTORY") && item.OptionValue == "true" {
				c.EnableDocumentHistory = true
			}
		}
	}
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
func (c *BaseController) JsonResult(errCode int,errMsg string,data ...interface{}){
	jsonData := make(map[string]interface{},3)

	jsonData["errcode"] = errCode
	jsonData["message"] = errMsg

	if len(data) > 0 && data[0] != nil{
		jsonData["data"] = data[0]
	}

	returnJSON, err := json.Marshal(jsonData)

	if err != nil {
		beego.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	io.WriteString(c.Ctx.ResponseWriter,string(returnJSON))

	c.StopRun()
}

// ExecuteViewPathTemplate 执行指定的模板并返回执行结果.
func (c *BaseController) ExecuteViewPathTemplate(tplName string,data interface{}) (string,error){
	var buf bytes.Buffer

	viewPath := c.ViewPath

	if c.ViewPath == "" {
		viewPath = beego.BConfig.WebConfig.ViewsPath

	}

	if err := beego.ExecuteViewPathTemplate(&buf,tplName,viewPath,data); err != nil {
		return "",err
	}
	return buf.String(),nil
}

func (c *BaseController) BaseUrl() string {
	return c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
}

//显示错误信息页面.
func (c *BaseController) ShowErrorPage(errCode int,errMsg string)  {
	c.TplName = "errors/error.tpl"
	c.Data["ErrorMessage"] = errMsg
	c.Data["ErrorCode"] = errCode
	c.StopRun()
}