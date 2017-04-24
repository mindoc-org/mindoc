package controllers

import (
	"strings"
	"regexp"
	"strconv"
	"time"
	"encoding/json"
	"html/template"

	"github.com/lifei6671/godoc/models"
	"github.com/lifei6671/godoc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/conf"
)

type BookController struct {
	BaseController
}

func (c *BookController) Index() {
	c.Prepare()
	c.TplName = "book/index.tpl"

	pageIndex, _ := c.GetInt("page", 1)

	books,totalCount,err := models.NewBook().FindToPager(pageIndex,conf.PageSize,c.Member.MemberId)

	if err != nil {
		c.Abort("500")
	}

	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, conf.PageSize, totalCount)

		c.Data["PageHtml"] = html
	}else {
		c.Data["PageHtml"] = ""
	}
	b,err := json.Marshal(books)

	if err != nil || len(books) <= 0{
		c.Data["Result"] = template.JS("[]")
	}else{
		c.Data["Result"] = template.JS(string(b))
	}
}

// Dashboard 项目概要 .
func (c *BookController) Dashboard() {
	c.Prepare()
	c.TplName = "book/dashboard.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
	}

	c.Data["Model"] = *book
}

// Setting 项目设置 .
func (c *BookController) Setting()  {
	c.Prepare()
	c.TplName = "book/setting.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
	}

	c.Data["Model"] = *book

}

// Users 用户列表.
func (c *BookController) Users() {
	c.Prepare()
	c.TplName = "book/users.tpl"

	key := c.Ctx.Input.Param(":key")
	pageIndex,_ := c.GetInt("page",1)

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
	}

	c.Data["Model"] = *book

	members,totalCount,err := models.NewMemberRelationshipResult().FindForUsersByBookId(book.BookId,pageIndex,15)

	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, 10, totalCount)

		c.Data["PageHtml"] = html
	}else{
		c.Data["PageHtml"] = ""
	}
	b,err := json.Marshal(members)

	if err != nil {
		c.Data["Result"] = template.JS("[]")
	}else{
		c.Data["Result"] = template.JS(string(b))
	}
}

// AddMember 参加参与用户.
func (c *BookController) AddMember()  {
	identify := c.GetString("identify")
	account := c.GetString("account")
	role_id,_ := c.GetInt("role_id",3)

	if identify == "" || account == ""{
		c.JsonResult(6001,"参数错误")
	}
	book ,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		c.JsonResult(6002,err.Error())
	}
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403,"权限不足")
	}

	member := models.NewMember()

	if err := member.FindByAccount(account) ; err != nil {
		c.JsonResult(404,"用户不存在")
	}
	if member.Status == 1 {
		c.JsonResult(6003,"用户已被禁用")
	}

	if _,err := models.NewRelationship().FindForRoleId(book.BookId,member.MemberId);err == nil {
		c.JsonResult(6003,"用户已存在该项目中")
	}

	relationship := models.NewRelationship()
	relationship.BookId = book.BookId
	relationship.MemberId = member.MemberId
	relationship.RoleId = role_id

	if err := relationship.Insert(); err == nil {
		memberRelationshipResult := models.NewMemberRelationshipResult().FromMember(member)
		memberRelationshipResult.RoleId = role_id
		memberRelationshipResult.RelationshipId = relationship.RelationshipId
		memberRelationshipResult.BookId = book.BookId
		memberRelationshipResult.ResolveRoleName()


		c.JsonResult(0,"ok",memberRelationshipResult)
	}
	c.JsonResult(500,err.Error())
}

// 变更指定用户在指定项目中的权限
func (c *BookController) ChangeRole() {
	identify := c.GetString("identify")
	member_id,_ := c.GetInt("member_id",0)
	role,_ := c.GetInt("role_id",0)

	if identify == "" || member_id <=0 {
		c.JsonResult(6001,"参数错误")
	}
	if member_id == c.Member.MemberId {
		c.JsonResult(6006,"不能变更自己的权限")
	}
	book ,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		c.JsonResult(6002,err.Error())
	}
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403,"权限不足")
	}

	member := models.NewMember()

	if err := member.Find(member_id); err != nil {
		c.JsonResult(6003,"用户不存在")
	}
	if member.Status == 1 {
		c.JsonResult(6004,"用户已被禁用")
	}

	relationship,err := models.NewRelationship().UpdateRoleId(book.BookId,member_id,role);

	if err != nil {
		logs.Error("变更用户在项目中的权限 => ",err)
		c.JsonResult(6005,err.Error())
	}

	memberRelationshipResult := models.NewMemberRelationshipResult().FromMember(member)
	memberRelationshipResult.RoleId = relationship.RoleId
	memberRelationshipResult.RelationshipId = relationship.RelationshipId
	memberRelationshipResult.BookId = book.BookId
	memberRelationshipResult.ResolveRoleName()

	c.JsonResult(0,"ok",memberRelationshipResult)
}

// 删除参与者.
func (c *BookController) RemoveMember()  {
	identify := c.GetString("identify")
	member_id,_ := c.GetInt("member_id",0)

	if identify == "" || member_id <=0 {
		c.JsonResult(6001,"参数错误")
	}
	if member_id == c.Member.MemberId {
		c.JsonResult(6006,"不能变更自己的权限")
	}
	book ,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		c.JsonResult(6002,err.Error())
	}
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403,"权限不足")
	}
	err = models.NewRelationship().DeleteByBookIdAndMemberId(book.BookId,member_id)

	if err != nil {
		c.JsonResult(6007,err.Error())
	}
	c.JsonResult(0,"ok")
}

// Create 创建项目.
func (c *BookController) Create() {

	if c.Ctx.Input.IsPost() {
		book_name := strings.TrimSpace(c.GetString("book_name",""))
		identify := strings.TrimSpace(c.GetString("identify",""))
		description := strings.TrimSpace(c.GetString("description",""))
		privately_owned,_ := strconv.Atoi(c.GetString("privately_owned"))
		comment_status := c.GetString("comment_status")

		if book_name == "" {
			c.JsonResult(6001,"项目名称不能为空")
		}
		if identify == "" {
			c.JsonResult(6002,"项目标识不能为空")
		}
		if ok,err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`,identify); !ok || err != nil {
			c.JsonResult(6003,"文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头")
		}
		if strings.Count(identify,"") > 50 {
			c.JsonResult(6004,"文档标识不能超过50字")
		}
		if strings.Count(description,"") > 500 {
			c.JsonResult(6004,"项目描述不能大于500字")
		}
		if privately_owned !=0 && privately_owned != 1 {
			privately_owned = 1
		}
		if comment_status != "open" && comment_status != "closed" && comment_status != "group_only" && comment_status != "registered_only" {
			comment_status = "closed"
		}

		book := models.NewBook()

		if books,_ := book.FindByField("identify",identify); len(books) > 0 {
			c.JsonResult(6006,"项目标识已存在")
		}

		book.BookName = book_name
		book.Description = description
		book.CommentCount = 0
		book.PrivatelyOwned = privately_owned
		book.CommentStatus = comment_status
		book.Identify = identify
		book.DocCount = 0
		book.MemberId = c.Member.MemberId
		book.CommentCount = 0
		book.Version = time.Now().Unix()
		book.Cover = beego.AppConfig.String("cover")

		err := book.Insert()

		if err != nil {
			c.JsonResult(6005,err.Error())
		}
		bookResult := models.NewBookResult()
		bookResult.FindByIdentify(book.Identify,c.Member.MemberId)

		c.JsonResult(0,"ok",bookResult)
	}
	c.JsonResult(6001,"error")
}

// Edit 编辑项目.
func (p *BookController) Edit() {
	p.TplName = "book/edit.tpl"

}

// CreateToken 创建访问来令牌.
func (c *BookController) CreateToken() {
	book_id,_ := c.GetInt("book_id",0)

	if book_id <= 0{
		c.JsonResult(6001,"参数错误")
	}

	book := models.NewBook()

	if err := book.Find(book_id);err != nil {
		c.JsonResult(6001,"项目不存在")
	}
	bookResult ,err := models.NewBookResult().FindByIdentify("identify",c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		logs.Error("生成阅读令牌失败 =>",err)
		c.JsonResult(6002,err.Error())
	}
	//必须是管理员或创始人才能删除项目
	if bookResult.RoleId != 0 && bookResult.RoleId != 1 {
		c.JsonResult(403,"权限不足")
	}
	if bookResult.PrivatelyOwned == 0 {
		c.JsonResult(6001,"公开项目不能创建阅读令牌")
	}

	book.PrivateToken = string(utils.Krand(20,utils.KC_RAND_KIND_ALL))
	if err := book.Update(); err != nil {
		logs.Error("生成阅读令牌失败 => ",err)
		c.JsonResult(6003,"生成阅读令牌失败")
	}
	c.JsonResult(0,"ok", c.BaseUrl() + "?token="+ book.PrivateToken)
}

// Delete 删除项目.
func (c *BookController) Delete() {
	c.Prepare()

	book_id,_ := c.GetInt("book_id",0)

	if book_id <= 0{
		c.JsonResult(6001,"参数错误")
	}

	book ,err := models.NewBookResult().FindByIdentify("identify",c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		logs.Error("删除项目 =>",err)
		c.JsonResult(6002,err.Error())
	}
	//必须是管理员或创始人才能删除项目
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403,"权限不足")
	}

	err = models.NewBook().ThoroughDeleteBook(book_id)

	if err == orm.ErrNoRows {
		c.JsonResult(6002,"项目不存在")
	}
	if err != nil {
		logs.Error("删除项目 => ",err)
		c.JsonResult(6003,"删除失败")
	}
	c.JsonResult(0,"ok")
}

// Transfer 转让项目.
func (p *BookController)Transfer()  {

}