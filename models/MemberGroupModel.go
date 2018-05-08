package models

import (
	"time"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
	"errors"
	"strconv"
)

type MemberGroup struct {
	GroupId int 				`orm:"column(group_id);pk;auto;unique;" json:"group_id"`
	GroupName string			`orm:"column(group_name);size(255);" json:"group_name"`
	GroupNumber int				`orm:"column(group_number);default(0)" json:"group_number"`
	CreateTime    time.Time 	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt      int       	`orm:"type(int);column(create_at)" json:"create_at"`
	CreateName 	string 			`orm:"-" json:"create_name"`
	CreateRealName string 	 	`orm:"-" json:"create_real_name"`
	ModifyTime time.Time     	`orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	Resources string		 	`orm:"column(resources);type(text);null" json:"-"`
	IsEnableDelete bool			`orm:"column(is_enable_delete);type(bool);default(true)" json:"is_enable_delete"`
	ResourceList []*ResourceModel	`orm:"-" json:"resource_list"`
	ModifyAt   int           	`orm:"column(modify_at);type(int)" json:"-"`
	ModifyName string 		 	`orm:"-" json:"modify_name"`
	ModifyRealName string 	 	`orm:"-" json:"modify_real_name"`
}

// TableName 获取对应数据库表名.
func (m *MemberGroup) TableName() string {
	return "member_group"
}

// TableEngine 获取数据使用的引擎.
func (m *MemberGroup) TableEngine() string {
	return "INNODB"
}

func (m *MemberGroup) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewMemberGroup() *MemberGroup {
	return &MemberGroup{}
}

//根据id查询用户组
func (m *MemberGroup) FindFirst(id int) (*MemberGroup,error){
	o := orm.NewOrm()

	if err :=o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",id).One(m); err != nil {
		beego.Error("查询用户组时出错 =>",err)
		return m,err
	}
	createMember,err := NewMember().Find(m.CreateAt);
	if err != nil {
		beego.Error("查询用户组创建人失败 =>",err)
	}else{

		m.CreateName = createMember.Account
		m.CreateRealName = createMember.RealName
	}

	if m.ModifyAt > 0 {
		modifyMember, err := NewMember().Find(m.ModifyAt)
		if err != nil {
			beego.Error("查询用户组修改人失败 =>",err)
		}else{

			m.ModifyName = modifyMember.Account
			m.ModifyRealName = modifyMember.RealName
		}
	}
	return m,nil
}

//删除指定用户组
func (m *MemberGroup) Delete(id int) error {
	o := orm.NewOrm()

	o.Begin()
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",id).Delete()

	if err != nil {
		o.Rollback()
		beego.Error("删除用户组失败 =>",err)
	}
	_,err = o.QueryTable(NewMemberGroupMembers().TableNameWithPrefix()).Filter("group_id",id).Delete()
	if err != nil {
		o.Rollback()
		beego.Error("删除用户组失败 =>",err)
	}
	return o.Commit()
}

//分页查询用户组
func (m *MemberGroup) FindByPager(pageIndex, pageSize int) ([]*MemberGroup,int,error){
	o := orm.NewOrm()

	if pageIndex <= 0 {
		pageIndex = 1
	}

	offset := (pageIndex - 1) * pageSize
	var memberGroups []*MemberGroup
	totalCount := 0
	_,err := o.QueryTable(m.TableNameWithPrefix()).OrderBy("-group_id").Offset(offset).Limit(pageSize).All(&memberGroups)

	if err != nil {
		beego.Error("分页查询用户组失败 =>",err)
	}else{
		i,err := o.QueryTable(m.TableNameWithPrefix()).Count()
		if err != nil {
			beego.Error("分页查询用户组失败 =>",err)
		}else {
			totalCount = int(i)
		}
	}
	memberIds := make([]int,0)

	for _,memberGroup := range memberGroups {
		if memberGroup.CreateAt > 0 {
			memberIds = append(memberIds,memberGroup.CreateAt)
		}
		if memberGroup.ModifyAt > 0 {
			memberIds = append(memberIds,memberGroup.ModifyAt)
		}
	}

	var members []*Member

	_,err = o.QueryTable(NewMember().TableNameWithPrefix()).Filter("member_id__in",memberIds).All(&members,"member_id","account","real_name")

	if err != nil {
		beego.Error("查询用户组信息时出错 =>",err)
	}else {
		for _,memberGroup := range memberGroups {
			for _,member := range members {
				if memberGroup.ModifyAt == member.MemberId {
					memberGroup.ModifyRealName = member.RealName
					memberGroup.ModifyName = member.Account
				}
				if memberGroup.CreateAt == member.MemberId {
					memberGroup.CreateRealName = member.RealName
					memberGroup.CreateName = member.Account
				}
			}
		}
	}

	return memberGroups,totalCount,err
}

//添加或更新用户组信息
func (m *MemberGroup) InsertOrUpdate(cols...string) error {
	o := orm.NewOrm()

	if m.ResourceList != nil && len(m.ResourceList) > 0 {
		for _,resource := range m.ResourceList {
			m.Resources = m.Resources + strconv.Itoa(resource.ResourceId) + ","
		}
		m.Resources = strings.Trim(m.Resources,",")
	}

	var err error
	//如果用户组已存在
	if m.GroupId > 0 {
		group := &MemberGroup{}

		err = o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",m.GroupId).One(group)
	}
	if err == nil {
		_,err = o.Update(m, cols...)
	}else{
		_,err = o.Insert(m)
	}
	return err
}

//重置用户组数量
func (m *MemberGroup) ResetMemberGroupNumber(groupId int) error {
	o := orm.NewOrm()

	i,err := o.QueryTable(NewMemberGroupMembers().TableNameWithPrefix()).Filter("group_id",groupId).Count()
	if err != nil {
		beego.Error("重置用户组用户数量失败 =>",err)
	}else{
		err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).One(m)
		if err != nil {
			beego.Error("重置用户组用户数量失败 =>",err)
			return err
		}else{
			m.GroupNumber  = int(i)
			_,err = o.Update(m)
			return err
		}
	}
	return  nil
}

//判断用户组是否存在
func (m *MemberGroup) Exist(groupId int) bool {
	o := orm.NewOrm()
	return o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).Exist()
}

//模糊查询用户组列表
func (m *MemberGroup) FindMemberGroupList(keyword string) ([]*MemberGroup,error) {
	o := orm.NewOrm()
	var memberGroups []*MemberGroup
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_name__icontains",keyword).All(&memberGroups)

	return memberGroups,err
}

//查询指定用户组的资源列表
func (m *MemberGroup) FindMemberGroupResourceList(groupId int) ([]*ResourceModel,error) {
	o := orm.NewOrm()

	memberGroup := NewMemberGroup()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("group_id",groupId).One(memberGroup)

	if err != nil {
		beego.Error("查询用户组资源时出错 =>", err)
		return nil,err
	}

	if memberGroup.Resources != "" {
		resourceIds := strings.Split(strings.Trim(memberGroup.Resources,","),",")

		var resources []*ResourceModel
		_,err = o.QueryTable(NewResource().TableNameWithPrefix()).Filter("resource_id__in",resourceIds).All(resources)
		if err != nil {
			beego.Error("查询用户组资源时出错 =>", err)
		}
		return resources,err
	}
	return nil,errors.New("当前用户组未设置资源")
}










