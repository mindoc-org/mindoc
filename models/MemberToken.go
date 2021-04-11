package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/mindoc-org/mindoc/conf"
)

type MemberToken struct {
	TokenId   int       `orm:"column(token_id);pk;auto;unique" json:"token_id"`
	MemberId  int       `orm:"column(member_id);type(int)" json:"member_id"`
	Token     string    `orm:"column(token);size(150);index" json:"token"`
	Email     string    `orm:"column(email);size(255)" json:"email"`
	IsValid   bool      `orm:"column(is_valid)" json:"is_valid"`
	ValidTime time.Time `orm:"column(valid_time);null" json:"valid_time"`
	SendTime  time.Time `orm:"column(send_time);auto_now_add;type(datetime)" json:"send_time"`
}

// TableName 获取对应数据库表名.
func (m *MemberToken) TableName() string {
	return "member_token"
}

// TableEngine 获取数据使用的引擎.
func (m *MemberToken) TableEngine() string {
	return "INNODB"
}

func (m *MemberToken) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewMemberToken() *MemberToken {
	return &MemberToken{}
}

func (m *MemberToken) InsertOrUpdate() (*MemberToken, error) {
	o := orm.NewOrm()

	if m.TokenId > 0 {
		_, err := o.Update(m)
		return m, err
	}
	_, err := o.Insert(m)

	return m, err
}

func (m *MemberToken) FindByFieldFirst(field string, value interface{}) (*MemberToken, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter(field, value).OrderBy("-token_id").One(m)

	return m, err
}

func (m *MemberToken) FindSendCount(mail string, start_time time.Time, end_time time.Time) (int, error) {
	o := orm.NewOrm()

	c, err := o.QueryTable(m.TableNameWithPrefix()).Filter("send_time__gte", start_time.Format("2006-01-02 15:04:05")).Filter("send_time__lte", end_time.Format("2006-01-02 15:04:05")).Count()

	if err != nil {
		return 0, err
	}
	return int(c), nil
}
