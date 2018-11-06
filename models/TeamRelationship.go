package models

import (
	"github.com/lifei6671/mindoc/conf"
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"errors"
)

type TeamRelationship struct {
	TeamRelationshipId int       `orm:"column(team_relationship_id);pk;auto;unique;" json:"team_relationship_id"`
	BookId             int       `orm:"column(book_id)"`
	TeamId             int       `orm:"column(team_id)"`
	CreateTime         time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	BookMemberId       int       `orm:"-" json:"book_member_id"`
	BookMemberName     string    `orm:"-" json:"book_member_name"`
	BookName           string    `orm:"-" json:"book_name"`
}

// TableName 获取对应数据库表名.
func (m *TeamRelationship) TableName() string {
	return "team_relationship"
}
func (m *TeamRelationship) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

// TableEngine 获取数据使用的引擎.
func (m *TeamRelationship) TableEngine() string {
	return "INNODB"
}

// 联合唯一键
func (m *TeamRelationship) TableUnique() [][]string {
	return [][]string{{"team_id", "book_id"}}
}
func (m *TeamRelationship) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m.TableNameWithPrefix())
}

func NewTeamRelationship() *TeamRelationship {
	return &TeamRelationship{}
}

//保存团队项目.
func (m *TeamRelationship) Save(cols ...string) (err error) {
	if m.TeamId <= 0 || m.BookId <= 0 {
		return ErrInvalidParameter
	}
	if (m.TeamRelationshipId > 0 && m.QueryTable().Filter("book_id", m.BookId).Filter("team_id", m.TeamId).Filter("team_relationship_id__ne", m.TeamRelationshipId).Exist()) ||
		m.TeamRelationshipId <= 0 && m.QueryTable().Filter("book_id", m.BookId).Filter("team_id", m.TeamId).Exist() {
		return errors.New("当前团队已加入该项目")
	}
	if m.TeamRelationshipId > 0 {
		_, err = orm.NewOrm().Update(m)
	} else {
		_, err = orm.NewOrm().Insert(m)
	}
	if err != nil {
		beego.Error("保存团队项目时出错 ->", err)
	}
	return
}

func (m *TeamRelationship) Delete(teamRelId int) (err error) {
	if teamRelId <= 0 {
		return ErrInvalidParameter
	}
	_, err = m.QueryTable().Filter("team_relationship_id", teamRelId).Delete()

	if err != nil {
		beego.Error("删除团队项目失败 ->", err)
	}
	return
}

//分页查询团队项目.
func (m *TeamRelationship) FindToPager(teamId, pageIndex, pageSize int) (list []*TeamRelationship, totalCount int, err error) {
	if teamId <= 0 {
		err = ErrInvalidParameter
		return
	}
	offset := (pageIndex - 1) * pageSize

	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("team_id", teamId).OrderBy("-team_relationship_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		beego.Error("查询团队项目时出错 ->", err)
		return
	}
	count, err := m.QueryTable().Filter("team_id", teamId).Count()

	if err != nil {
		beego.Error("查询团队项目时出错 ->", err)
		return
	}
	totalCount = int(count)
	for _, item := range list {
		item.Include()
	}
	return
}

//加载附加数据.
func (m *TeamRelationship) Include() (*TeamRelationship, error) {
	if m.BookId > 0 {
		b, err := NewBook().Find(m.BookId, "book_name", "identify", "member_id")
		if err != nil {
			return m, err
		}
		m.BookName = b.BookName
		m.BookMemberId = b.MemberId
		if b.MemberId > 0 {
			member, err := NewMember().Find(b.MemberId, "account", "real_name")
			if err != nil {
				return m, err
			}
			if member.RealName == "" {
				m.BookMemberName = member.Account
			} else {
				m.BookMemberName = member.RealName
			}
		}
	}
	return m, nil
}

//查询未加入团队的用户。
func (m *TeamRelationship) FindNotJoinBookByName(teamId int, bookName string, limit int) (*SelectMemberResult, error) {
	if teamId <= 0 {
		return nil, ErrInvalidParameter
	}
	o := orm.NewOrm()

	sql := `select book.book_id,book.book_name
from  md_books as book
where book.book_id not in (select team.book_id from md_team_relationship as team where team_id=?)
and book.book_name like ? order by book_id desc limit ?;`

	books := make([]*Book, 0)

	_, err := o.Raw(sql, teamId, "%"+bookName+"%", limit).QueryRows(&books)

	if err != nil {
		beego.Error("查询团队项目时出错 ->", err)
		return nil, err
	}

	result := SelectMemberResult{}
	items := make([]KeyValueItem, 0)

	for _, book := range books {
		item := KeyValueItem{}
		item.Id = book.BookId
		item.Text = book.BookName
		items = append(items, item)
	}
	result.Result = items

	return &result, err
}

//查询指定项目的团队.
func (m *TeamRelationship) FindByBookToPager(bookId, pageIndex,pageSize int) (list []*TeamRelationship,totalCount int,err error) {

	if bookId <= 0{
		err = ErrInvalidParameter
		return
	}

	offset := (pageIndex - 1) * pageSize

	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", bookId).OrderBy("-team_relationship_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		beego.Error("查询团队项目时出错 ->", err)
		return
	}
	count, err := m.QueryTable().Filter("book_id", bookId).Count()

	if err != nil {
		beego.Error("查询团队项目时出错 ->", err)
		return
	}
	totalCount = int(count)
	for _, item := range list {
		item.Include()
	}
	return
}





























