package controllers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
)

type SearchController struct {
	BaseController
}

func (c *SearchController) Index() {
	c.Prepare()
	c.TplName = "search/index.tpl"
	c.Data["SIDEBAR_ID"] = "search"
	c.Data["SIDEBAR_BOOK"] = 0

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}
	keyword := c.GetString("keyword")

	sidebar_id := c.GetString("sidebar_id")
	book_identify := c.GetString("book_identify")

	if sidebar_id == "users" {
		c.Redirect(beego.URLFor("ManagerController.Users", "keyword", keyword), 302)
		return
	} else if sidebar_id == "books" {
		c.Redirect(beego.URLFor("ManagerController.Books", "keyword", keyword), 302)
		return
	} else if sidebar_id == "attach" {
		c.Redirect(beego.URLFor("ManagerController.AttachList", "keyword", keyword), 302)
		return
	} else if sidebar_id == "mybook" {
		c.Redirect(beego.URLFor("BookController.Index", "keyword", keyword), 302)
		return
	} else if sidebar_id == "bookuser" {
		c.Redirect(beego.URLFor("BookController.Users", ":key", book_identify, "keyword", keyword), 302)
		return
	} else if sidebar_id == "booklink" {
		c.Redirect(beego.URLFor("BookController.Links", ":key", book_identify, "keyword", keyword), 302)
		return
	} else if sidebar_id == "bookattach" {
		c.Redirect(beego.URLFor("BookController.Attach", ":key", book_identify, "keyword", keyword), 302)
		return
	}

	pageIndex, _ := c.GetInt("page", 1)

	c.Data["BaseUrl"] = c.BaseUrl()

	if keyword != "" {
		c.Data["Keyword"] = keyword
		member_id := 0
		if c.Member != nil {
			member_id = c.Member.MemberId
		}
		search_result, totalCount, err := models.NewDocumentSearchResult().FindToPager(keyword, pageIndex, conf.PageSize, member_id)

		if err != nil {
			beego.Error(err)
			return
		}
		if totalCount > 0 {
			html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, conf.PageSize, totalCount)

			c.Data["PageHtml"] = html
		} else {
			c.Data["PageHtml"] = ""
		}
		if len(search_result) > 0 {
			for _, item := range search_result {
				item.DocumentName = strings.Replace(item.DocumentName, keyword, "<em>"+keyword+"</em>", -1)

				if item.Description != "" {
					src := item.Description

					//将HTML标签全转换成小写
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					src = re.ReplaceAllStringFunc(src, strings.ToLower)

					//去除STYLE
					re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
					src = re.ReplaceAllString(src, "")

					//去除SCRIPT
					re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
					src = re.ReplaceAllString(src, "")

					//去除所有尖括号内的HTML代码，并换成换行符
					re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
					src = re.ReplaceAllString(src, "\n")

					//去除连续的换行符
					re, _ = regexp.Compile("\\s{2,}")
					src = re.ReplaceAllString(src, "\n")

					r := []rune(src)

					if len(r) > 100 {
						src = string(r[:100])
					} else {
						src = string(r)
					}
					item.Description = strings.Replace(src, keyword, "<em>"+keyword+"</em>", -1)
				}

				if item.Identify == "" {
					item.Identify = strconv.Itoa(item.DocumentId)
				}
				if item.ModifyTime.IsZero() {
					item.ModifyTime = item.CreateTime
				}
			}
		}
		c.Data["Lists"] = search_result
	}
}
