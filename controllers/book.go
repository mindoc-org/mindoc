package controllers

import "strings"

type BookController struct {
	BaseController
}

func (p *BookController) Index() {
	p.TplName = "book/index.tpl"
}

// Dashboard 项目概要 .
func (p *BookController) Dashboard() {
	p.TplName = "book/dashboard.tpl"
}

// Setting 项目设置 .
func (p *BookController) Setting()  {
	p.TplName = "book/setting.tpl"
}

func (p *BookController) Users() {
	p.TplName = "book/users.tpl"
}

func (c *BookController) Create() {

	if c.Ctx.Input.IsPost() {
		book_name := strings.TrimSpace(c.GetString("book_name",""))
		identify := strings.TrimSpace(c.GetString("identify",""))
		description := strings.TrimSpace(c.GetString("description",""))
		privately_owned := c.GetString("privately_owned")
		comment_status := c.GetString("comment_status")

	}
	c.JsonResult(6001,"error")
}

// Edit 编辑项目.
func (p *BookController) Edit() {
	p.TplName = "book/edit.tpl"
}

// Delete 删除项目.
func (p *BookController) Delete() {
	p.StopRun()
}

// Transfer 转让项目.
func (p *BookController)Transfer()  {

}