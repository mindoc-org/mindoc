package models

import (
	"github.com/lifei6671/mindoc/conf"
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type MemberGroupMembers struct {
	GroupMemberId int		`orm:"column(group_member_id);pk;auto;unique;" json:"group_member_id"`
	MemberId int			`orm:"column(member_id);index" json:"member_id"`
	GroupId int 			`orm:"column(group_id);index" json:"group_id"`
	CreateTime    time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt      int       `orm:"type(int);column(create_at)" json:"create_at"`
}

type MemberGroupMemberResult struct {
	GroupMemberId int
	MemberId int
	Account string
	RealName string
	Avatar string
	GroupId int
	CreateTime    time.Time
	CreateAt      int
}

// TableName 获取对应数据库表名.
func (m *MemberGroupMembers) TableName() string {
	return "member_group_members"
}

// TableEngine 获取数据使用的引擎.
func (m *MemberGroupMembers) TableEngine() string {
	return "INNODB"
}

func (m *MemberGroupMembers) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

// 多字段唯一键
func (u *MemberGroupMembers) TableUnique() [][]string {
	return [][]string{
		{"member_id", "group_id"},
	}
}

func NewMemberGroupMembers() *MemberGroupMembers {
	return &MemberGroupMembers{}
}

// 查询用户组成员
func (m *MemberGroupMembers) FindByGroupId(groupId int) ([]*MemberGroupMemberResult,error) {
	o := orm.NewOrm()
	var groupMembers []*MemberGroupMemberResult
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).All(&groupMembers);
	if err != nil {
		beego.Error("获取用户组成员出错 =>",err)
		return nil,err
	}
	ids := make([]int,0)

	for _,member := range groupMembers {
		ids = append(ids,member.MemberId)
	}

	var members []*Member

	_,err = o.QueryTable(NewMember().TableNameWithPrefix()).Filter("member_id__in",ids).All(&members)
	if err != nil {
		beego.Error("获取用户组成员出错 =>",err)
		return nil,err
	}
	for _,member := range members {
		for _,groupMember := range groupMembers {
			if groupMember.MemberId == member.MemberId {
				groupMember.Account = member.Account
				groupMember.Avatar = member.Avatar
				groupMember.RealName = member.RealName
			}
		}
	}
	return groupMembers,nil
}






























