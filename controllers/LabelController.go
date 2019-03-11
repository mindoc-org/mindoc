package controllers

import (
	"math"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils/pagination"
)

type LabelController struct {
	BaseController
}

func (c *LabelController) Prepare() {
	c.BaseController.Prepare()

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		return
	}
}

//查看包含标签的文档列表.
func (c *LabelController) Index() {
	c.Prepare()
	c.TplName = "label/index.tpl"

	labelName := c.Ctx.Input.Param(":key")
	pageIndex, _ := c.GetInt("page", 1)
	if labelName == "" {
		c.Abort("404")
	}
	_, err := models.NewLabel().FindFirst("label_name", labelName)

	if err != nil {
		if err == orm.ErrNoRows {
			c.Abort("404")
		} else {
			beego.Error(err)
			c.Abort("500")
		}
	}
	memberId := 0
	if c.Member != nil {
		memberId = c.Member.MemberId
	}
	searchResult, totalCount, err := models.NewBook().FindForLabelToPager(labelName, pageIndex, conf.PageSize, memberId)

	if err != nil && err != orm.ErrNoRows {
		beego.Error("查询标签时出错 ->", err)
		c.ShowErrorPage(500, "查询文档列表时出错")
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Lists"] = searchResult

	c.Data["LabelName"] = labelName
}

func (c *LabelController) List() {
	c.Prepare()
	c.TplName = "label/list.tpl"

	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 200

	labels, totalCount, err := models.NewLabel().FindToPager(pageIndex, pageSize)

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	c.Data["Labels"] = labels
}
