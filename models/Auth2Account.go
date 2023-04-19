// Package models .
package models

import (
	"errors"
	"github.com/mindoc-org/mindoc/utils/auth2"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
)

var (
	_ Auth2Account = (*WorkWeixinAccount)(nil)
	_ Auth2Account = (*DingTalkAccount)(nil)
)

type Auth2Account interface {
	ExistedMember(id string) (*Member, error)
	AddBind(o orm.Ormer, userInfo auth2.UserInfo, member *Member) error
}

func NewWorkWeixinAccount() *WorkWeixinAccount {
	return &WorkWeixinAccount{}
}

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

func (m *WorkWeixinAccount) ExistedMember(workweixin_user_id string) (*Member, error) {
	o := orm.NewOrm()
	account := NewWorkWeixinAccount()
	member := NewMember()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("workweixin_user_id", workweixin_user_id).One(account)
	if err != nil {
		return member, err
	}

	member, err = member.Find(account.MemberId)
	if err != nil {
		return member, err
	}

	if member.Status != 0 {
		return member, errors.New("receive_account_disabled")
	}

	return member, nil

}

// AddBind 添加一个用户.
func (m *WorkWeixinAccount) AddBind(o orm.Ormer, userInfo auth2.UserInfo, member *Member) error {
	tmpM := NewWorkWeixinAccount()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("workweixin_user_id", userInfo.UserId).One(tmpM)
	if err == nil {
		tmpM.MemberId = member.MemberId
		_, err = o.Update(tmpM)
		if err != nil {
			logs.Error("保存用户数据到数据时失败 =>", err)
			return errors.New("用户信息绑定失败, 数据库错误")
		}
		return nil
	}

	m.MemberId = member.MemberId
	m.WorkWeixin_UserId = userInfo.UserId

	if c, err := o.QueryTable(m.TableNameWithPrefix()).Filter("member_id", m.MemberId).Count(); err == nil && c > 0 {
		return errors.New("已绑定，不可重复绑定")
	}

	_, err = o.Insert(m)
	if err != nil {
		logs.Error("保存用户数据到数据时失败 =>", err)
		return errors.New("用户信息绑定失败, 数据库错误")
	}

	return nil
}

func NewDingTalkAccount() *DingTalkAccount {
	return &DingTalkAccount{}
}

type DingTalkAccount struct {
	MemberId        int       `orm:"column(member_id);type(int);default(-1);index" json:"member_id"`
	UserDbId        int       `orm:"pk;auto;unique;column(user_db_id)" json:"user_db_id"`
	Dingtalk_UserId string    `orm:"size(100);unique;column(dingtalk_user_id)" json:"dingtalk_user_id"`
	CreateTime      time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt        int       `orm:"type(int);column(create_at)" json:"create_at"`
	LastLoginTime   time.Time `orm:"type(datetime);column(last_login_time);null" json:"last_login_time"`
}

// TableName 获取对应数据库表名.
func (m *DingTalkAccount) TableName() string {
	return "dingtalk_accounts"
}

// TableEngine 获取数据使用的引擎.
func (m *DingTalkAccount) TableEngine() string {
	return "INNODB"
}

func (m *DingTalkAccount) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func (m *DingTalkAccount) ExistedMember(userid string) (*Member, error) {
	o := orm.NewOrm()
	account := NewDingTalkAccount()
	member := NewMember()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("dingtalk_user_id", userid).One(account)
	if err != nil {
		return member, err
	}

	member, err = member.Find(account.MemberId)
	if err != nil {
		return member, err
	}

	if member.Status != 0 {
		return member, errors.New("receive_account_disabled")
	}

	return member, nil

}

// AddBind 添加一个用户.
func (m *DingTalkAccount) AddBind(o orm.Ormer, userInfo auth2.UserInfo, member *Member) error {
	tmpM := NewDingTalkAccount()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("dingtalk_user_id", userInfo.UserId).One(tmpM)
	if err == nil {
		tmpM.MemberId = member.MemberId
		_, err = o.Update(tmpM)
		if err != nil {
			logs.Error("保存用户数据到数据时失败 =>", err)
			return errors.New("用户信息绑定失败, 数据库错误")
		}
		return nil
	}

	m.Dingtalk_UserId = userInfo.UserId
	m.MemberId = member.MemberId

	if c, err := o.QueryTable(m.TableNameWithPrefix()).Filter("member_id", m.MemberId).Count(); err == nil && c > 0 {
		return errors.New("已绑定，不可重复绑定")
	}

	_, err = o.Insert(m)
	if err != nil {
		logs.Error("保存用户数据到数据时失败 =>", err)
		return errors.New("用户信息绑定失败, 数据库错误")
	}

	return nil
}
