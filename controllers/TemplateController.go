package controllers

import (
	"errors"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
)

type TemplateController struct {
	BaseController
	BookId int
}

func (c *TemplateController) isPermission() (error) {
	c.Prepare()

	bookIdentify := c.GetString("identify", "")

	if bookIdentify == "" {
		return errors.New("参数错误")
	}

	if !c.Member.IsAdministrator() {
		book, err := models.NewBookResult().FindByIdentify(bookIdentify, c.Member.MemberId)
		if err != nil {
			if err == orm.ErrNoRows {
				return errors.New("项目不存在或没有权限")
			}
			return errors.New("查询项目模板失败")
		}
		c.BookId = book.BookId
	} else {
		book, err := models.NewBook().FindByIdentify(bookIdentify, "book_id")
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
//获取指定模板信息
func (c *TemplateController) Get() {
	if err := c.isPermission() ; err != nil {
		c.JsonResult(500,err.Error())
	}
	templateId, _ := c.GetInt("template_id",0)

	template,err := models.NewTemplate().Find(templateId)

	if err != nil {
		c.JsonResult(500,"读取模板失败")
	}
	if template.IsGlobal == 0 && template.BookId != c.BookId {
		c.JsonResult(404,"模板不存在或已删除")
	}
	c.JsonResult(0,"OK",template)
}

//获取模板列表
func (c *TemplateController) List() {
	c.TplName = "template/list.tpl"
	if err := c.isPermission() ; err != nil {
		c.Data["ErrorMessage"] = err.Error()
		return
	}

	templateList,err := models.NewTemplate().FindAllByBookId(c.BookId)

	if err != nil && err != orm.ErrNoRows{
		c.Data["ErrorMessage"] = "查询项目模板失败"
	}
	if templateList != nil {
		for i,t := range templateList {
			templateList[i] = t.Preload()
		}
	}
	//c.JsonResult(0,"OK",templateList)
	c.Data["List"] = templateList
}

//添加模板
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


	if templateId > 0 {
		t,err := template.Find(templateId);
		if err != nil {
			c.JsonResult(500,"模板不存在")
		}

		template = t
		template.ModifyAt = c.Member.MemberId
	}


	template.TemplateId = templateId
	template.BookId = c.BookId
	template.TemplateContent = content
	template.TemplateName = templateName

	//只有管理员才能设置全局模板
	if c.Member.IsAdministrator() {
		template.IsGlobal = isGlobal
	}else{
		template.IsGlobal = 0
		//如果不是管理员需要判断是否有项目权限
		rel,err := models.NewRelationship().FindByBookIdAndMemberId(c.BookId,c.Member.MemberId)
		if err != nil || rel.RoleId == conf.BookObserver {
			c.JsonResult(403,"没有权限")
		}
		//如果修改的模板不是本人创建的，并且又不是项目创建者则禁止修改
		if template.MemberId > 0 && template.MemberId != c.Member.MemberId && rel.RoleId != conf.BookFounder {
			c.JsonResult(403,"没有权限")
		}
	}
	template.MemberId = c.Member.MemberId

	var cols []string

	if templateId > 0 {
		cols = []string{ "template_content", "modify_time","modify_at","version" }
	}

	if err := template.Save(cols...); err != nil {
		c.JsonResult(500,"报错模板失败")
	}
	c.JsonResult(0,"OK",template)
}

//删除模板
func (c *TemplateController) Delete() {
	if err := c.isPermission() ; err != nil {
		c.JsonResult(500,err.Error())
	}
	templateId, _ := c.GetInt("template_id",0)

	template,err := models.NewTemplate().Find(templateId)

	if err != nil {
		c.JsonResult(404,"模板不存在")
	}

	if c.Member.IsAdministrator() {
		err := models.NewTemplate().Delete(templateId,0)
		if err != nil {
			c.JsonResult(500,"删除失败")
		}
	}else{
		//如果不是管理员需要判断是否有项目权限
		rel,err := models.NewRelationship().FindByBookIdAndMemberId(template.BookId,c.Member.MemberId)
		if err != nil || rel.RoleId == conf.BookObserver {
			c.JsonResult(403,"没有权限")
		}

		//如果是创始人或管理者可以直接删除模板
		if  rel.RoleId == conf.BookFounder || rel.RoleId == conf.BookAdmin {
			err := models.NewTemplate().Delete(templateId,0)
			if err != nil {
				c.JsonResult(500,"删除失败")
			}
		}

		if err := models.NewTemplate().Delete(templateId,c.Member.MemberId);err != nil {
			c.JsonResult(500,"删除失败")
		}
	}
	c.JsonResult(0,"OK")
}