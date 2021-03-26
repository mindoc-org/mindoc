package controllers

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
)

type BookMemberController struct {
	BaseController
}

// AddMember 参加参与用户.
func (c *BookMemberController) AddMember() {
	identify := c.GetString("identify")
	account, _ := c.GetInt("account")
	roleId, _ := c.GetInt("role_id", 3)
	logs.Info(account)
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
	relationship.RoleId = conf.BookRole(roleId)

	if err := relationship.Insert(); err == nil {
		memberRelationshipResult := models.NewMemberRelationshipResult().FromMember(member)
		memberRelationshipResult.RoleId = conf.BookRole(roleId)
		memberRelationshipResult.RelationshipId = relationship.RelationshipId
		memberRelationshipResult.BookId = book.BookId
		memberRelationshipResult.ResolveRoleName()

		c.JsonResult(0, "ok", memberRelationshipResult)
	}
	c.JsonResult(500, err.Error())
}

// 变更指定用户在指定项目中的权限
func (c *BookMemberController) ChangeRole() {
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

	relationship, err := models.NewRelationship().UpdateRoleId(book.BookId, memberId, conf.BookRole(role))

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
	member_id, _ := c.GetInt("member_id", 0)

	if identify == "" || member_id <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	if member_id == c.Member.MemberId {
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
	err = models.NewRelationship().DeleteByBookIdAndMemberId(book.BookId, member_id)

	if err != nil {
		c.JsonResult(6007, err.Error())
	}
	c.JsonResult(0, "ok")
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
