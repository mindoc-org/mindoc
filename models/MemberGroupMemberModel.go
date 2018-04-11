package models

import (
	"github.com/lifei6671/mindoc/conf"
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
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
func (m *MemberGroupMembers) TableUnique() [][]string {
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
//添加或更新用户组成员
func (m *MemberGroupMembers) InsertOrUpdate(cols ...string) error {
	o := orm.NewOrm()

	if m.GroupMemberId > 0 {
		_,err := o.Update(m,cols...)
		if err != nil {
			beego.Error("更新用户组成员失败 =>",err)
		}
		return err
	}else{
		if m.GroupId <= 0 {
			return errors.New("用户组不能为空")
		}
		_,err := o.Insert(m)
		if err != nil {
			beego.Error("添加用户组成员失败 =>",err)
		}else{
			o.Raw(fmt.Sprintf("UPDATE %s SET group_number=group_number+1 WHERE group_id=%d",NewMemberGroup().TableNameWithPrefix(), m.GroupId)).Exec()
		}

		return err
	}
}
//删除用户组成员
func (m *MemberGroupMembers) Delete(id int) error {
	o := orm.NewOrm()

	i,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_member_id",id).Delete()

	if err != nil {
		beego.Error("删除用户组成员失败 =>",err)
		return err
	}
	if i <= 0 {
		beego.Info("删除用户组成员返回行数 =>",i)
	}

	return nil
}

//分页获取用户组成员
func (m *MemberGroupMembers) FindToPager(pageIndex, pageSize, groupId int) ([]*MemberGroupMemberResult,int,error) {
	o := orm.NewOrm()

	if pageIndex <= 0 {
		pageIndex = 1
	}

	offset := (pageIndex - 1) * pageSize
	var memberGroupMembers []*MemberGroupMembers
	totalCount := 0
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).Offset(offset).Limit(pageSize).All(&memberGroupMembers)

	memberGroupMemberList := make([]*MemberGroupMemberResult,0)
	if err != nil {
		beego.Error("分页查询用户组成员失败 =>",err)
	}else{
		i,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).Count()
		if err != nil {
			beego.Error("分页查询用户组成员失败 =>",err)
		}else {
			totalCount = int(i)
		}
		for _,member := range memberGroupMembers {
			memberGroupMemberList = append(memberGroupMemberList,member.ToMemberRelationshipResult())
		}
	}

	return memberGroupMemberList,totalCount,err
}

//将用户组信息转换为完整的用户信息
func (m *MemberGroupMembers) ToMemberRelationshipResult() *MemberGroupMemberResult {
	memberGroupMemberResult := &MemberGroupMemberResult{}

	memberGroupMemberResult.GroupId = m.GroupId
	memberGroupMemberResult.GroupMemberId = m.GroupMemberId
	memberGroupMemberResult.MemberId = m.MemberId
	memberGroupMemberResult.CreateAt = m.CreateAt
	memberGroupMemberResult.CreateTime = m.CreateTime

	if m.MemberId > 0 {
		o := orm.NewOrm()
		member := NewMember()
		_,err := o.QueryTable(member.TableNameWithPrefix()).Filter("member_id", m.MemberId).All(&member)
		if err != nil {
			beego.Error("查询用户组成员信息时出错 =>",err)
		}else{
			memberGroupMemberResult.RealName = member.RealName
			memberGroupMemberResult.Avatar = member.Avatar
			memberGroupMemberResult.Account = member.Account
		}
	}

	return memberGroupMemberResult
}





































