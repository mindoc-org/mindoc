package controllers

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils/pagination"
)

type ItemsetsController struct {
	BaseController
}

func (c *ItemsetsController) Prepare() {
	c.BaseController.Prepare()

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		return
	}
}
func (c *ItemsetsController) Index() {
	c.Prepare()
	c.TplName = "items/index.tpl"
	pageSize := 18

	pageIndex, _ := c.GetInt("page", 0)

	items, totalCount, err := models.NewItemsets().FindToPager(pageIndex, pageSize)

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	c.Data["TotalPages"] = pageIndex
	if err == orm.ErrNoRows || len(items) <= 0 {
		c.Data["Lists"] = items
		c.Data["PageHtml"] = ""
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, pageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["Lists"] = items
}

func (c *ItemsetsController) List() {
	c.Prepare()
	c.TplName = "items/list.tpl"
	pageSize := 18
	itemKey := c.Ctx.Input.Param(":key")
	pageIndex, _ := c.GetInt("page", 1)

	if itemKey == "" {
		c.Abort("404")
	}
	item, err := models.NewItemsets().FindFirst(itemKey)

	if err != nil {
		if err == orm.ErrNoRows {
			c.Abort("404")
		} else {
			logs.Error(err)
			c.Abort("500")
		}
	}
	memberId := 0
	if c.Member != nil {
		memberId = c.Member.MemberId
	}
	searchResult, totalCount, err := models.NewItemsets().FindItemsetsByItemKey(itemKey, pageIndex, pageSize, memberId)

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, "查询文档列表时出错")
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, pageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Lists"] = searchResult

	c.Data["Model"] = item
}
