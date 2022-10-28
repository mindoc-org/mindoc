package controllers

import (
	"strings"
	"time"

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

	// 获取评论、分页
	comments, count, pageIndex := models.NewComment().QueryCommentByDocumentId(docid, pageIndex, conf.PageSize, c.Member)
	page := pagination.PageUtil(int(count), pageIndex, conf.PageSize, comments)

	var data struct {
		DocId int             `json:"doc_id"`
		Page  pagination.Page `json:"page"`
	}
	data.DocId = docid
	data.Page = page

	c.JsonResult(0, "ok", data)
	return
}

func (c *CommentController) Create() {
	content := c.GetString("content")
	id, _ := c.GetInt("doc_id")

	_, err := models.NewDocument().Find(id)
	if err != nil {
		c.JsonResult(1, "文章不存在")
	}

	m := models.NewComment()
	m.DocumentId = id
	if c.Member == nil {
		c.JsonResult(1, "请先登录，再评论")
	}
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
	m.Insert()

	var data struct {
		DocId int `json:"doc_id"`
	}
	data.DocId = id

	c.JsonResult(0, "ok", data)
}

func (c *CommentController) Index() {
	c.Prepare()
	c.TplName = "comment/index.tpl"
}

func (c *CommentController) Delete() {
	if c.Ctx.Input.IsPost() {
		id, _ := c.GetInt("id", 0)
		m, err := models.NewComment().Find(id)
		if err != nil {
			c.JsonResult(1, "评论不存在")
		}

		doc, err := models.NewDocument().Find(m.DocumentId)
		if err != nil {
			c.JsonResult(1, "文章不存在")
		}

		// 判断是否有权限删除
		bookRole, _ := models.NewRelationship().FindForRoleId(doc.BookId, c.Member.MemberId)
		if m.CanDelete(c.Member.MemberId, bookRole) {
			err := m.Delete()
			if err != nil {
				c.JsonResult(1, "删除错误")
			} else {
				c.JsonResult(0, "ok")
			}
		} else {
			c.JsonResult(1, "没有权限删除")
		}
	}
}
