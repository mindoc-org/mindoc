package controllers

import (
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils/pagination"
)

type CommentController struct {
	BaseController
}

func (c *CommentController) Lists() {

	docid, _ := c.GetInt("docid", 0)
	pageIndex, _ := c.GetInt("page", 1)

	beego.Info("CommentController.Lists", docid, pageIndex)

	// 获取评论、分页
	comments, count := models.NewComment().QueryCommentByDocumentId(docid, pageIndex, conf.PageSize)
	page := pagination.PageUtil(int(count), pageIndex, conf.PageSize, comments)
	beego.Info("docid=", docid, "Page", page)

	var data struct {
		DocId     int               `json:"doc_id"`
		Page      pagination.Page   `json:"page"`
	}
	data.DocId = docid
	data.Page = page

	c.JsonResult(0, "ok", data)
	return
}

func (c *CommentController) Create() {
	content := c.GetString("content")
	id, _ := c.GetInt("doc_id")

	m := models.NewComment()
	m.DocumentId = id
	if len(c.Member.RealName) != 0 {
		m.Author = c.Member.RealName
	} else {
		m.Author = c.Member.Account
	}
	m.MemberId = c.Member.MemberId
	m.IPAddress = c.Ctx.Request.RemoteAddr
	m.IPAddress = strings.Split(m.IPAddress, ":")[0]
	m.CommentDate = time.Now()
	m.Content = content
	beego.Info(m)
	m.Insert()

	c.JsonResult(0, "ok")
}

func (c *CommentController) Index() {

	c.Prepare()
	c.TplName = "comment/index.tpl"
}
