package controllers

import (
	"github.com/lifei6671/godoc/models"
	"github.com/astaxie/beego/logs"
)

type DocumentController struct {
	BaseController
}

func (p *DocumentController) Index()  {
	p.TplName = "document/index.tpl"
}

func (p *DocumentController) Read() {
	p.TplName = "document/kancloud.tpl"
}

func (c *DocumentController) Edit()  {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")

	book,err := models.NewBook().FindByFieldFirst("identify",identify)

	if err != nil {
		logs.Error("DocumentController.Edit => ",err)
		c.Abort("500")
	}
	if book.Editor == "markdown" {
		c.TplName = "document/markdown_edit_template.tpl"
	}else{
		c.TplName = "document/html_edit_template.tpl"
	}

}