package controllers

import (
	"github.com/lifei6671/mindoc/models"
	"github.com/astaxie/beego/orm"
	"github.com/qiniu/x/errors.v7"
	"strings"
)

type TemplateController struct {
	BaseController
	BookId int
}

func (c *TemplateController) isPermission() (error) {
	c.Prepare()

	if c.IsAjax() {
		bookIdentify := c.GetString("identify", "")

		if bookIdentify == "" {
			return  errors.New("参数错误")
		}

		if !c.Member.IsAdministrator() {
			book, err := models.NewBookResult().FindByIdentify(bookIdentify, c.Member.MemberId,"book_id")
			if err != nil {
				if err == orm.ErrNoRows {
					return errors.New("项目不存在或没有权限")
				}
				return errors.New("查询项目模板失败")
			}
			c.BookId = book.BookId
		}else{
			book,err := models.NewBook().FindByIdentify(bookIdentify,"book_id")
			if err != nil {
				if err == orm.ErrNoRows {
					return errors.New("项目不存在或没有权限")
				}
				return errors.New("查询项目模板失败")
			}
			c.BookId = book.BookId
		}
		return nil
	}
	return errors.New("请求方法不支持")
}

//获取模板列表
func (c *TemplateController) List() {
	if err := c.isPermission() ; err != nil {
		c.JsonResult(500,err.Error())
	}

	templateList,err := models.NewTemplate().FindAllByBookId(c.BookId)

	if err != nil {
		if err == orm.ErrNoRows {
			c.JsonResult(404,"没有模板")
		}
		c.JsonResult(500,"查询项目模板失败")
	}
	c.JsonResult(0,"OK",templateList)
}

func (c *TemplateController) Add() {
	if err := c.isPermission() ; err != nil {
		c.JsonResult(500,err.Error())
	}

	templateId, _ := c.GetInt("template_id",0)
	content := c.GetString("content")
	isGlobal,_ := c.GetInt("is_global",0)
	templateName := c.GetString("template_name","")

	if templateName == "" || strings.Count(templateName,"") > 300 {
		c.JsonResult(500,"模板名称不能为空且必须小于300字")
	}
	template := models.NewTemplate()
	template.TemplateId = templateId
	template.BookId = c.BookId
	template.TemplateContent = content
	template.MemberId = c.Member.MemberId

	if templateId > 0 {
		template.ModifyAt = c.Member.MemberId
	}
	//只有管理员才能设置全局模板
	if c.Member.IsAdministrator() {
		template.IsGlobal = isGlobal
	}else{
		template.IsGlobal = 0
	}

	var cols []string

	if templateId > 0 {
		cols = []string{ "template_content", "modify_time","modify_at","version" }
	}

	if err := template.Save(cols...); err != nil {
		c.JsonResult(500,"报错模板失败")
	}
	c.JsonResult(0,"OK",template)
}

func (c *TemplateController) Delete() {
	if err := c.isPermission() ; err != nil {
		c.JsonResult(500,err.Error())
	}
	templateId, _ := c.GetInt("template_id",0)

	if c.Member.IsAdministrator() {
		err := models.NewTemplate().Delete(templateId,0)
		if err != nil {
			c.JsonResult(500,"删除失败")
		}
	}else{
		err := models.NewTemplate().Delete(templateId,c.Member.MemberId)
		if err != nil {
			c.JsonResult(500,"删除失败")
		}
	}
	c.JsonResult(0,"OK")
}