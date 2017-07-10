package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
	"math"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Prepare()
	c.TplName = "home/index.tpl"
	//如果没有开启匿名访问，则跳转到登录页面
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"),302)
	}
	pageIndex,_ := c.GetInt("page",1)
	pageSize := 18

	member_id := 0

	if c.Member != nil {
		member_id = c.Member.MemberId
	}
	books,totalCount,err := models.NewBook().FindForHomeToPager(pageIndex,pageSize,member_id)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, pageSize, totalCount)

		c.Data["PageHtml"] = html
	}else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	c.Data["Lists"] = books

	labels ,totalCount,err := models.NewLabel().FindToPager(1,10)

	if err != nil {
		c.Data["Labels"] = make([]*models.Label,0)
	}else{
		c.Data["Labels"] = labels
	}
}
