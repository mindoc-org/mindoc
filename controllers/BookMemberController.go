package controllers

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
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
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	book, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
	}

	member := models.NewMember()

	if _, err := member.Find(account); err != nil {
		c.JsonResult(404, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	if member.Status == 1 {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.user_disable"))
	}

	if _, err := models.NewRelationship().FindForRoleId(book.BookId, member.MemberId); err == nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.user_exist_in_proj"))
	}
	//如果是只读用户，只能设置为观察者
	if member.Role == conf.MemberReaderRole && roleId != int(conf.BookObserver) {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.readusr_only_observer"))
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
		memberRelationshipResult.ResolveRoleName(c.Lang)

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
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	if memberId == c.Member.MemberId {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.cannot_change_own_priv"))
	}
	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404, i18n.Tr(c.Lang, "message.item_not_exist"))
		}
		c.JsonResult(6002, err.Error())
	}
	if book.RoleId != 0 && book.RoleId != 1 {
		c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
	}

	member := models.NewMember()

	if _, err := member.Find(memberId); err != nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	if member.Status == 1 {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.user_disable"))
	}
	//如果是只读用户，只能设置为观察者
	if member.Role == conf.MemberReaderRole && role != int(conf.BookObserver) {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.readusr_only_observer"))
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
	memberRelationshipResult.ResolveRoleName(c.Lang)

	c.JsonResult(0, "ok", memberRelationshipResult)
}

// 删除参与者.
func (c *BookMemberController) RemoveMember() {
	identify := c.GetString("identify")
	member_id, _ := c.GetInt("member_id", 0)

	if identify == "" || member_id <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	if member_id == c.Member.MemberId {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.cannot_delete_self"))
	}
	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404, i18n.Tr(c.Lang, "message.item_not_exist"))
		}
		c.JsonResult(6002, err.Error())
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
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
			return book, errors.New(i18n.Tr(c.Lang, "message.no_permission"))
		}
		if err == orm.ErrNoRows {
			return book, errors.New(i18n.Tr(c.Lang, "message.item_not_exist"))
		}
		return book, err
	}
	if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
		return book, errors.New(i18n.Tr(c.Lang, "message.no_permission"))
	}
	return book, nil
}
