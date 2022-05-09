// Package models .
package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
)

type WorkWeixinAccount struct {
	MemberId          int    `orm:"column(member_id);type(int);default(-1);index" json:"member_id"`
	UserDbId          int    `orm:"pk;auto;unique;column(user_db_id)" json:"user_db_id"`
	WorkWeixin_UserId string `orm:"size(100);unique;column(workweixin_user_id)" json:"workweixin_user_id"`
	// WorkWeixin_Name   string    `orm:"size(255);column(workweixin_name)" json:"workweixin_name"`
	// WorkWeixin_Phone  string    `orm:"size(25);column(workweixin_phone)" json:"workweixin_phone"`
	// WorkWeixin_Email  string    `orm:"size(255);column(workweixin_email)" json:"workweixin_email"`
	// WorkWeixin_Status int       `orm:"type(int);column(status)" json:"status"`
	// WorkWeixin_Avatar string    `orm:"size(1024);column(avatar)" json:"avatar"`
	CreateTime    time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt      int       `orm:"type(int);column(create_at)" json:"create_at"`
	LastLoginTime time.Time `orm:"type(datetime);column(last_login_time);null" json:"last_login_time"`
}

// TableName 获取对应数据库表名.
func (m *WorkWeixinAccount) TableName() string {
	return "workweixin_accounts"
}

// TableEngine 获取数据使用的引擎.
func (m *WorkWeixinAccount) TableEngine() string {
	return "INNODB"
}

func (m *WorkWeixinAccount) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewWorkWeixinAccount() *WorkWeixinAccount {
	return &WorkWeixinAccount{}
}

func (a *WorkWeixinAccount) ExistedMember(workweixin_user_id string) (*Member, error) {
	o := orm.NewOrm()
	account := NewWorkWeixinAccount()
	member := NewMember()
	err := o.QueryTable(a.TableNameWithPrefix()).Filter("workweixin_user_id", workweixin_user_id).One(account)
	if err == nil {
		if member, err = member.Find(account.MemberId); err == nil {
			return member, nil
		} else {
			return member, err
		}
	} else {
		return member, err
	}
}

// Add 添加一个用户.
func (a *WorkWeixinAccount) AddBind(o orm.Ormer) error {
	if c, err := o.QueryTable(a.TableNameWithPrefix()).Filter("member_id", a.MemberId).Count(); err == nil && c > 0 {
		return errors.New("已绑定，不可重复绑定")
	}

	_, err := o.Insert(a)
	if err != nil {
		logs.Error("保存用户数据到数据时失败 =>", err)
		return errors.New("用户信息绑定失败, 数据库错误")
	}

	return nil
}
