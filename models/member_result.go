package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
)

type MemberRelationshipResult struct {
	MemberId int		`json:"member_id"`
	Account string 		`json:"account"`
	Description string	`json:"description"`
	Email string 		`json:"email"`
	Phone string 		`json:"phone"`
	Avatar string 		`json:"avatar"`
	Role int		`json:"role"`	//用户角色：0 管理员/ 1 普通用户
	Status int 		`json:"status"`	//用户状态：0 正常/1 禁用
	CreateTime time.Time	`json:"create_time"`
	CreateAt int		`json:"create_at"`
	RelationshipId int      `json:"relationship_id"`
	BookId int		`json:"book_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId int		`json:"role_id"`
	RoleName string		`json:"role_name"`
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

	return m
}

func (m *MemberRelationshipResult) ResolveRoleName () *MemberRelationshipResult {
	if m.RoleId == conf.BookAdmin {
		m.RoleName = "管理者"
	}else if m.RoleId == conf.BookEditor {
		m.RoleName = "编辑者"
	}else if m.RoleId == conf.BookObserver {
		m.RoleName = "观察者"
	}
	return m
}

func (m *MemberRelationshipResult) FindForUsersByBookId(book_id ,pageIndex, pageSize int) ([]*MemberRelationshipResult,int,error) {
	o := orm.NewOrm()

	var members []*MemberRelationshipResult

	sql1 := "SELECT * FROM md_relationship AS rel LEFT JOIN md_members as member ON rel.member_id = member.member_id WHERE rel.book_id = ? ORDER BY rel.relationship_id DESC  LIMIT ?,?"

	sql2 := "SELECT count(*) AS total_count FROM md_relationship AS rel LEFT JOIN md_members as member ON rel.member_id = member.member_id WHERE rel.book_id = ?"

	var total_count int

	err := o.Raw(sql2,book_id).QueryRow(&total_count)

	if err != nil {
		return members,0,err
	}

	offset := (pageIndex-1) * pageSize

	_,err = o.Raw(sql1,book_id,offset,pageSize).QueryRows(&members)

	if err != nil {
		return members,0,err
	}

	for _,item := range members {
		item.ResolveRoleName()
	}
	return members,total_count,nil
}






