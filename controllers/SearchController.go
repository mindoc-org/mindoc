package controllers

import (
	"strconv"
	"strings"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/pagination"
	"github.com/mindoc-org/mindoc/utils/sqltil"
)

type SearchController struct {
	BaseController
}

//搜索首页
func (c *SearchController) Index() {
	c.Prepare()
	c.TplName = "search/index.tpl"

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		return
	}

	keyword := c.GetString("keyword")
	pageIndex, _ := c.GetInt("page", 1)

	c.Data["BaseUrl"] = c.BaseUrl()

	if keyword != "" {
		c.Data["Keyword"] = keyword
		memberId := 0
		if c.Member != nil {
			memberId = c.Member.MemberId
		}
		searchResult, totalCount, err := models.NewDocumentSearchResult().FindToPager(sqltil.EscapeLike(keyword), pageIndex, conf.PageSize, memberId)

		if err != nil {
			logs.Error("搜索失败 ->", err)
			return
		}
		if totalCount > 0 {
			pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
			c.Data["PageHtml"] = pager.HtmlPages()
		} else {
			c.Data["PageHtml"] = ""
		}
		if len(searchResult) > 0 {
			keywords := strings.Split(keyword, " ")

			for _, item := range searchResult {
				for _, word := range keywords {
					item.DocumentName = strings.Replace(item.DocumentName, word, "<em>"+word+"</em>", -1)
					if item.Description != "" {
						src := item.Description

						r := []rune(utils.StripTags(item.Description))

						if len(r) > 100 {
							src = string(r[:100])
						} else {
							src = string(r)
						}
						item.Description = strings.Replace(src, word, "<em>"+word+"</em>", -1)
					}
				}
				if item.Identify == "" {
					item.Identify = strconv.Itoa(item.DocumentId)
				}
				if item.ModifyTime.IsZero() {
					item.ModifyTime = item.CreateTime
				}
			}
		}
		c.Data["Lists"] = searchResult
	}
}

//搜索用户
func (c *SearchController) User() {
	c.Prepare()
	key := c.Ctx.Input.Param(":key")
	keyword := strings.TrimSpace(c.GetString("q"))
	if key == "" || keyword == "" {
		c.JsonResult(404, "参数错误")
	}
	keyword = sqltil.EscapeLike(keyword)

	book, err := models.NewBookResult().FindByIdentify(key, c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, "没有权限")
		}
		c.JsonResult(500, "项目不存在")
	}

	//members, err := models.NewMemberRelationshipResult().FindNotJoinUsersByAccount(book.BookId, 10, "%"+keyword+"%")
	members, err := models.NewMemberRelationshipResult().FindNotJoinUsersByAccountOrRealName(book.BookId, 10, "%"+keyword+"%")
	if err != nil {
		logs.Error("查询用户列表出错：" + err.Error())
		c.JsonResult(500, err.Error())
	}
	result := models.SelectMemberResult{}
	items := make([]models.KeyValueItem, 0)

	for _, member := range members {
		item := models.KeyValueItem{}
		item.Id = member.MemberId
		item.Text = member.Account + "[" + member.RealName + "]"
		items = append(items, item)
	}

	result.Result = items

	c.JsonResult(0, "OK", result)
}
