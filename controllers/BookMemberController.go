package controllers

import (
	"errors"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/astaxie/beego"
	"strings"
)

type BookMemberController struct {
	BaseController
}

// AddMember 参加参与用户.
func (c *BookMemberController) AddMember() {
	c.Prepare()
	identify := c.GetString("identify")
	account,_ := c.GetInt("account")
	roleId, _ := c.GetInt("role_id", 3)
	beego.Info(account)
	if identify == "" || account <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	book, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
	}


	member := models.NewMember()

	if _, err := member.Find(account); err != nil {
		c.JsonResult(404, "用户不存在")
	}
	if member.Status == 1 {
		c.JsonResult(6003, "用户已被禁用")
	}

	if _, err := models.NewRelationship().FindForRoleId(book.BookId, member.MemberId); err == nil {
		c.JsonResult(6003, "用户已存在该项目中")
	}

	relationship := models.NewRelationship()
	relationship.BookId = book.BookId
	relationship.MemberId = member.MemberId
	relationship.RoleId = roleId

	if err := relationship.Insert(); err == nil {
		memberRelationshipResult := models.NewMemberRelationshipResult().FromMember(member)
		memberRelationshipResult.RoleId = roleId
		memberRelationshipResult.RelationshipId = relationship.RelationshipId
		memberRelationshipResult.BookId = book.BookId
		memberRelationshipResult.ResolveRoleName()

		c.JsonResult(0, "ok", memberRelationshipResult)
	}
	c.JsonResult(500, err.Error())
}

// 变更指定用户在指定项目中的权限
func (c *BookMemberController) ChangeRole() {
	c.Prepare()
	identify := c.GetString("identify")
	memberId, _ := c.GetInt("member_id", 0)
	role, _ := c.GetInt("role_id", 0)

	if identify == "" || memberId <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	if memberId == c.Member.MemberId {
		c.JsonResult(6006, "不能变更自己的权限")
	}
	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, "权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404, "项目不存在")
		}
		c.JsonResult(6002, err.Error())
	}
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403, "权限不足")
	}

	member := models.NewMember()

	if _, err := member.Find(memberId); err != nil {
		c.JsonResult(6003, "用户不存在")
	}
	if member.Status == 1 {
		c.JsonResult(6004, "用户已被禁用")
	}

	relationship, err := models.NewRelationship().UpdateRoleId(book.BookId, memberId, role)

	if err != nil {
		logs.Error("变更用户在项目中的权限 => ", err)
		c.JsonResult(6005, err.Error())
	}

	memberRelationshipResult := models.NewMemberRelationshipResult().FromMember(member)
	memberRelationshipResult.RoleId = relationship.RoleId
	memberRelationshipResult.RelationshipId = relationship.RelationshipId
	memberRelationshipResult.BookId = book.BookId
	memberRelationshipResult.ResolveRoleName()

	c.JsonResult(0, "ok", memberRelationshipResult)
}

// 删除参与者.
func (c *BookMemberController) RemoveMember() {
	identify := c.GetString("identify")
	memberId, _ := c.GetInt("member_id", 0)

	if identify == "" || memberId <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	if memberId == c.Member.MemberId {
		c.JsonResult(6006, "不能删除自己")
	}
	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, "权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404, "项目不存在")
		}
		c.JsonResult(6002, err.Error())
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.JsonResult(403, "权限不足")
	}
	err = models.NewRelationship().DeleteByBookIdAndMemberId(book.BookId, memberId)

	if err != nil {
		c.JsonResult(6007, err.Error())
	}
	c.JsonResult(0, "ok")
}

//添加用户组到项目
func (c *BookMemberController) AddMemberGroup() {
	c.Prepare()
	memberGroupId,err := c.GetInt("group_id")
	roleId, _ := c.GetInt("role_id", 3)

	if roleId != 1 && roleId != 2 {
		roleId = 3
	}
	if err != nil {
		beego.Error("解析用户组ID时失败 =>",err)
		c.JsonResult(6001,"参数异常")
	}
	if memberGroupId <= 0 {
		c.JsonResult(6002,"参数错误")
	}
	bookResult,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6003,err.Error())
	}

	if !models.NewMemberGroup().Exist(memberGroupId) {
		beego.Error("查询用户组时失败 =>",err)
		c.JsonResult(6004,"用户组不存在")
	}
	memberGroupMembers,err := models.NewMemberGroupMembers().FindByGroupId(memberGroupId)
	if err != nil {
		beego.Error("查询用户组用户时时失败 =>",err)
		c.JsonResult(6004,"用户组成员不存在不存在")
	}

	for _,item := range memberGroupMembers {
		member,err := models.NewMember().Find(item.MemberId)

		if err != nil {
			beego.Error("用户不存在 =>",item.MemberId)
			continue
		}
		if member.Status == 1 {
			beego.Error("用户被禁用 =>",item.MemberId)
			continue
		}

		if _, err := models.NewRelationship().FindForRoleId(bookResult.BookId, member.MemberId); err == nil {
			beego.Error("用户已存在该项目中 =>",item.MemberId)
			continue
		}

		relationship := models.NewRelationship()
		relationship.BookId = bookResult.BookId
		relationship.MemberId = member.MemberId
		relationship.RoleId = roleId
		if err := relationship.Insert();err != nil {
			beego.Error("添加用户失败 =>",err)
		}
	}
	c.JsonResult(0,"ok")
}

func (c *BookMemberController)  MemberGroupList() {
	c.Prepare()

	q := strings.TrimSpace(c.GetString("q"))

	members,err := models.NewMemberGroup().FindMemberGroupList(q)

	if err != nil {
		beego.Error("查询异常",err)
		c.JsonResult(6001, "查询错误")
	}
	result := models.SelectMemberResult{}
	items := make([]models.KeyValueItem, 0)

	for _, member := range members {
		item := models.KeyValueItem{}
		item.Id = member.GroupId
		item.Text = member.GroupName
		items = append(items, item)
	}

	result.Result = items

	c.JsonResult(0, "OK", result)
}

func (c *BookMemberController) IsPermission() (*models.BookResult, error) {
	identify := c.GetString("identify")
	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			return book, errors.New("权限不足")
		}
		if err == orm.ErrNoRows {
			return book, errors.New("项目不存在")
		}
		return book, err
	}
	if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
		return book, errors.New("权限不足")
	}
	return book, nil
}

































