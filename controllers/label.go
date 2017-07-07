package controllers

import (
	"github.com/lifei6671/mindoc/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/utils"
)

type LabelController struct {
	BaseController
}

func (c *LabelController) Index() {
	c.Prepare()
	c.TplName = "label/index.tpl"

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"),302)
		return
	}

	labelName := c.Ctx.Input.Param(":key")
	pageIndex,_ := c.GetInt("page",1)
	if labelName == "" {
		c.Abort("404")
	}
	_,err := models.NewLabel().FindFirst("label_name",labelName)

	if err != nil {
		if err == orm.ErrNoRows {
			c.Abort("404")
		}else{
			beego.Error(err)
			c.Abort("500")
		}
	}
	member_id := 0
	if c.Member != nil {
		member_id = c.Member.MemberId
	}
	search_result,totalCount,err := models.NewBook().FindForLabelToPager(labelName,pageIndex,conf.PageSize,member_id)

	if err != nil {
		beego.Error(err)
		return
	}
	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, conf.PageSize, totalCount)

		c.Data["PageHtml"] = html
	}else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Lists"] = search_result

	c.Data["LabelName"] = labelName
}
