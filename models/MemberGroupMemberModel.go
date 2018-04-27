package models

import (
	"github.com/lifei6671/mindoc/conf"
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
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
	GroupMemberId int		`json:"group_member_id"`
	MemberId int			`json:"member_id"`
	Account string			`json:"account"`
	RealName string			`json:"real_name"`
	RoleName string			`json:"role_name"`
	Avatar string			`json:"avatar"`
	GroupId int				`json:"group_id"`
	CreateTime    time.Time	`json:"create_time"`
	CreateAt      int		`json:"create_at"`
}

//搜索未加入当前用户组的用户信息
type MemberGroupMemberNoJoinSearchResult struct {
	MemberId int	`json:"member_id"`
	Account string	`json:"account"`
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
func (m *MemberGroupMembers) FindResultByGroupId(groupId int) ([]*MemberGroupMemberResult,error) {
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

func (m *MemberGroupMembers) FindByGroupId(groupId int) ([]*MemberGroupMembers,error) {
	o := orm.NewOrm()
	var groupMembers []*MemberGroupMembers
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).All(&groupMembers);
	if err != nil {
		beego.Error("获取用户组成员出错 =>",err)
		return nil,err
	}
	return groupMembers,nil
}

//判断一个用户是否加入了指定用户组
func (m *MemberGroupMembers) IsJoin(groupId ,memberId int) (bool) {
	o := orm.NewOrm()

	return o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).Filter("member_id",memberId).Exist()
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
			go NewMemberGroup().ResetMemberGroupNumber(m.GroupId)
		}

		return err
	}
}

//删除用户组成员
func (m *MemberGroupMembers) Delete(id int) error {
	o := orm.NewOrm()

	if err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_member_id",id).One(m);err != nil {
		beego.Error("删除用户组成员失败 =>",err)
		return err
	}

	i,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_member_id",id).Delete()

	if err != nil {
		beego.Error("删除用户组成员失败 =>",err)
		return err
	}
	if i <= 0 {
		beego.Info("删除用户组成员返回行数 =>",i)
	}

	go NewMemberGroup().ResetMemberGroupNumber(m.GroupId)

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
		member,err := NewMember().Find(memberGroupMemberResult.MemberId)
		if err != nil {
			beego.Error("查询用户组成员信息时出错 =>",err)
		}else{
			memberGroupMemberResult.RealName = member.RealName
			memberGroupMemberResult.Avatar = member.Avatar
			memberGroupMemberResult.Account = member.Account
			memberGroupMemberResult.RoleName = member.RoleName
		}
	}

	return memberGroupMemberResult
}

//获取未加入用户组的用户成员
func (m *MemberGroupMembers) FindMemberGroupMemberNoJoinSearchResult(groupId int,q string) ([]*MemberGroupMemberNoJoinSearchResult,error) {
	o := orm.NewOrm()

	sql := "select member.member_id,member.account from md_members as member left join md_member_group_members as member_group on member_group.member_id = member.member_id and member_group.group_id=? where member_group.member_id isnull  and account like ?limit 20;"
	var memberGroupMembers []*MemberGroupMemberNoJoinSearchResult

	_,err := o.Raw(sql,groupId,"%" + q + "%").QueryRows(&memberGroupMembers)

	if err != nil {
		beego.Error("获取未加入用户组的用户失败 =>",err)
	}
	return memberGroupMembers,err
}































