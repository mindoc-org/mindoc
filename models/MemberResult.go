package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/i18n"
	"github.com/mindoc-org/mindoc/conf"
)

type MemberRelationshipResult struct {
	MemberId       int             `json:"member_id"`
	Account        string          `json:"account"`
	RealName       string          `json:"real_name"`
	Description    string          `json:"description"`
	Email          string          `json:"email"`
	Phone          string          `json:"phone"`
	Avatar         string          `json:"avatar"`
	Role           conf.SystemRole `json:"role"`   //用户角色：0 管理员/ 1 普通用户
	Status         int             `json:"status"` //用户状态：0 正常/1 禁用
	CreateTime     time.Time       `json:"create_time"`
	CreateAt       int             `json:"create_at"`
	RelationshipId int             `json:"relationship_id"`
	BookId         int             `json:"book_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId   conf.BookRole `json:"role_id"`
	RoleName string        `json:"role_name"`
}

type SelectMemberResult struct {
	Result []KeyValueItem `json:"results"`
}
type KeyValueItem struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

func NewMemberRelationshipResult() *MemberRelationshipResult {
	return &MemberRelationshipResult{}
}

func (m *MemberRelationshipResult) FromMember(member *Member) *MemberRelationshipResult {
	m.MemberId = member.MemberId
	m.Account = member.Account
	m.Description = member.Description
	m.Email = member.Email
	m.Phone = member.Phone
	m.Avatar = member.Avatar
	m.Role = member.Role
	m.Status = member.Status
	m.CreateTime = member.CreateTime
	m.CreateAt = member.CreateAt
	m.RealName = member.RealName

	return m
}

func (m *MemberRelationshipResult) ResolveRoleName(lang string) *MemberRelationshipResult {
	if m.RoleId == conf.BookAdmin {
		m.RoleName = i18n.Tr(lang, "common.administrator")
	} else if m.RoleId == conf.BookEditor {
		m.RoleName = i18n.Tr(lang, "common.editor")
	} else if m.RoleId == conf.BookObserver {
		m.RoleName = i18n.Tr(lang, "common.observer")
	}
	return m
}

// 根据项目ID查询用户
func (m *MemberRelationshipResult) FindForUsersByBookId(lang string, bookId, pageIndex, pageSize int) ([]*MemberRelationshipResult, int, error) {
	o := orm.NewOrm()

	var members []*MemberRelationshipResult

	sql1 := "SELECT * FROM md_relationship AS rel LEFT JOIN md_members as mdmb ON rel.member_id = mdmb.member_id WHERE rel.book_id = ? ORDER BY rel.relationship_id DESC  limit ? offset ?"

	sql2 := "SELECT count(*) AS total_count FROM md_relationship AS rel LEFT JOIN md_members as mdmb ON rel.member_id = mdmb.member_id WHERE rel.book_id = ?"

	var total_count int

	err := o.Raw(sql2, bookId).QueryRow(&total_count)

	if err != nil {
		return members, 0, err
	}

	offset := (pageIndex - 1) * pageSize

	_, err = o.Raw(sql1, bookId, pageSize, offset).QueryRows(&members)

	if err != nil {
		return members, 0, err
	}

	for _, item := range members {
		item.ResolveRoleName(lang)
	}
	return members, total_count, nil
}

// 查询指定文档中不存在的用户列表
func (m *MemberRelationshipResult) FindNotJoinUsersByAccount(bookId, limit int, account string) ([]*Member, error) {
	o := orm.NewOrm()

	sql := "SELECT m.* FROM md_members as m LEFT JOIN md_relationship as rel ON m.member_id=rel.member_id AND rel.book_id = ? WHERE rel.relationship_id IS NULL AND m.account LIKE ? LIMIT 0,?;"

	var members []*Member

	_, err := o.Raw(sql, bookId, account, limit).QueryRows(&members)

	return members, err
}

// 根据姓名以及用户名模糊查询指定文档中不存在的用户列表
func (m *MemberRelationshipResult) FindNotJoinUsersByAccountOrRealName(bookId, limit int, keyWord string) ([]*Member, error) {
	o := orm.NewOrm()

	sql := "SELECT m.* FROM md_members as m LEFT JOIN md_relationship as rel ON rel.member_id = m.member_id AND rel.book_id = ? WHERE rel.relationship_id IS NULL AND (m.real_name LIKE ? OR m.account LIKE ?) LIMIT ? OFFSET 0;"

	var members []*Member

	_, err := o.Raw(sql, bookId, keyWord, keyWord, limit).QueryRows(&members)

	return members, err
}
