// Package models .
package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/utils"
)

type Member struct {
	MemberId int		`orm:"pk;auto;unique;column(member_id)"`
	Account string 		`orm:"size(100);unique;column(account)"`
	Password string 	`orm:"size(1000);column(password)"`
	Description string	`orm:"column(description);size(2000)" json:"description"`
	Email string 		`orm:"size(255);column(email);null;default(null)"`
	Phone string 		`orm:"size(255);column(phone);null;default(null)"`
	Avatar string 		`orm:"size(1000);column(avatar)"`
	Role int		`orm:"column(role);type(int);default(1);index"`	//用户角色：0 管理员/ 1 普通用户
	Status int 		`orm:"column(status);type(int);default(0)"`	//用户状态：0 正常/1 禁用
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int		`orm:"type(int);column(create_at)"`
	LastLoginTime time.Time	`orm:"type(datetime);column(last_login_time);null"`
}

// TableName 获取对应数据库表名.
func (m *Member) TableName() string {
	return "members"
}
// TableEngine 获取数据使用的引擎.
func (m *Member) TableEngine() string {
	return "INNODB"
}

func NewMember() *Member {
	return &Member{}
}

// Login 用户登录.
func (m *Member) Login(account string,password string) (*Member,error) {
	o := orm.NewOrm()

	member := &Member{}

	err := o.QueryTable("md_" + m.TableName()).Filter("account",account).Filter("status",0).One(member);

	if err != nil {
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
	return nil
}