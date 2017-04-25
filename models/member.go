// Package models .
package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/utils"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/logs"
)

type Member struct {
	MemberId int		`orm:"pk;auto;unique;column(member_id)" json:"member_id"`
	Account string 		`orm:"size(100);unique;column(account)" json:"account"`
	Password string 	`orm:"size(1000);column(password)" json:"-"`
	Description string	`orm:"column(description);size(2000)" json:"description"`
	Email string 		`orm:"size(255);column(email);null;default(null)" json:"email"`
	Phone string 		`orm:"size(255);column(phone);null;default(null)" json:"phone"`
	Avatar string 		`orm:"size(1000);column(avatar)" json:"avatar"`
	//用户角色：0 超级管理员 /1 管理员/ 2 普通用户 .
	Role int		`orm:"column(role);type(int);default(1);index" json:"role"`
	RoleName string 	`orm:"-" json:"role_name"`
	Status int 		`orm:"column(status);type(int);default(0)" json:"status"`	//用户状态：0 正常/1 禁用
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt int		`orm:"type(int);column(create_at)" json:"create_at"`
	LastLoginTime time.Time	`orm:"type(datetime);column(last_login_time);null" json:"last_login_time"`

}

// TableName 获取对应数据库表名.
func (m *Member) TableName() string {
	return "members"
}
// TableEngine 获取数据使用的引擎.
func (m *Member) TableEngine() string {
	return "INNODB"
}

func (m *Member)TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() +  m.TableName()
}

func NewMember() *Member {
	return &Member{}
}

// Login 用户登录.
func (m *Member) Login(account string,password string) (*Member,error) {
	o := orm.NewOrm()

	member := &Member{}

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("account",account).Filter("status",0).One(member);

	if err != nil {
		logs.Error("用户登录 => ",err)
		return  member,ErrMemberNoExist
	}

	ok,err := utils.PasswordVerify(member.Password,password) ;

	if ok && err == nil {
		return member,nil
	}

	return member,ErrorMemberPasswordError
}

// Add 添加一个用户.
func (member *Member) Add () (error) {
	o := orm.NewOrm()

	hash ,err := utils.PasswordHash(member.Password);

	if  err != nil {
		return err
	}

	member.Password = hash

	_,err = o.Insert(member)

	if err != nil {
		return err
	}
	return  nil
}

// Update 更新用户信息.
func (m *Member) Update(cols... string) (error) {
	o := orm.NewOrm()

	if _,err := o.Update(m,cols...);err != nil {
		return err
	}
	return nil
}

func (m *Member) Find(id int) error{
	o := orm.NewOrm()

	m.MemberId = id
	if err := o.Read(m); err != nil {
		return  err
	}
	if m.Role == conf.MemberSuperRole {
		m.RoleName = "超级管理员"
	}else if m.Role == conf.MemberAdminRole {
		m.RoleName = "管理员"
	}else if m.Role == conf.MemberGeneralRole {
		m.RoleName = "普通用户"
	}
	return nil
}

func (m *Member) ResolveRoleName (){
	if m.Role == 0 {
		m.RoleName = "超级管理员"
	}else if m.Role == 1 {
		m.RoleName = "管理员"
	}else if m.Role == 2 {
		m.RoleName = "普通用户"
	}
}

func (m *Member) FindByAccount (account string) (*Member,error)  {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("account",account).One(m)

	if err == nil {
		m.ResolveRoleName()
	}
	return m,err
}

func (m *Member) FindToPager(pageIndex, pageSize int) ([]*Member,int64,error)  {
	o := orm.NewOrm()

	var members []*Member

	offset := (pageIndex - 1) * pageSize

	totalCount,err := o.QueryTable(m.TableNameWithPrefix()).Count()

	if err != nil {
		return members,0,err
	}

	_,err = o.QueryTable(m.TableNameWithPrefix()).OrderBy("-member_id").Offset(offset).Limit(pageSize).All(&members)

	if err != nil {
		return members,0,err
	}

	for _,m := range members {
		if m.Role == 0 {
			m.RoleName = "超级管理员"
		}else if m.Role == 1 {
			m.RoleName = "管理员"
		}else if m.Role == 2 {
			m.RoleName = "普通用户"
		}
	}
	return members,totalCount,nil
}


func (c *Member) IsAdministrator() bool {
	if c == nil || c.MemberId <= 0{
		return false
	}
	return c.Role == 0 || c.Role == 1
}



















