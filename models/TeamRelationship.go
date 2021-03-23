package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
)

type TeamRelationship struct {
	TeamRelationshipId int       `orm:"column(team_relationship_id);pk;auto;unique;" json:"team_relationship_id"`
	BookId             int       `orm:"column(book_id)" json:"book_id"`
	TeamId             int       `orm:"column(team_id)" json:"team_id"`
	CreateTime         time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	TeamName           string    `orm:"-" json:"team_name"`
	MemberCount        int       `orm:"-" json:"member_count"`
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

func (m *TeamRelationship) First(teamId int, cols ...string) (*TeamRelationship, error) {
	if teamId <= 0 {
		return nil, ErrInvalidParameter
	}
	err := m.QueryTable().Filter("team_id", teamId).One(m, cols...)
	if err != nil {
		logs.Error("查询项目团队失败 ->", err)
	}
	return m, err
}

//查找指定项目的指定团队.
func (m *TeamRelationship) FindByBookId(bookId int, teamId int) (*TeamRelationship, error) {
	if teamId <= 0 || bookId <= 0 {
		return nil, ErrInvalidParameter
	}
	err := m.QueryTable().Filter("team_id", teamId).Filter("book_id", bookId).One(m)
	if err != nil {
		logs.Error("查询项目团队失败 ->", err)
	}
	return m, err
}

//删除指定项目的指定团队.
func (m *TeamRelationship) DeleteByBookId(bookId int, teamId int) error {
	err := m.QueryTable().Filter("team_id", teamId).Filter("book_id", bookId).One(m)
	if err != nil {
		logs.Error("查询项目团队失败 ->", err)
		return err
	}
	m.Include()
	return m.Delete(m.TeamRelationshipId)
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
		logs.Error("保存团队项目时出错 ->", err)
	}
	return
}

func (m *TeamRelationship) Delete(teamRelId int) (err error) {
	if teamRelId <= 0 {
		return ErrInvalidParameter
	}
	_, err = m.QueryTable().Filter("team_relationship_id", teamRelId).Delete()

	if err != nil {
		logs.Error("删除团队项目失败 ->", err)
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
		logs.Error("查询团队项目时出错 ->", err)
		return
	}
	count, err := m.QueryTable().Filter("team_id", teamId).Count()

	if err != nil {
		logs.Error("查询团队项目时出错 ->", err)
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
	if m.TeamId > 0 {
		team, err := NewTeam().First(m.TeamId)
		if err == nil {
			m.TeamName = team.TeamName
			m.MemberCount = team.MemberCount
		}
	}
	return m, nil
}

//查询未加入团队的项目.
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
		logs.Error("查询团队项目时出错 ->", err)
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

//查找指定项目中未加入的团队.
func (m *TeamRelationship) FindNotJoinBookByBookIdentify(bookId int, teamName string, limit int) (*SelectMemberResult, error) {
	if bookId <= 0 || teamName == "" {
		return nil, ErrInvalidParameter
	}

	o := orm.NewOrm()
	sql := `select *
from md_teams as team
where team.team_id not in (select rel.team_id from md_team_relationship as rel where rel.book_id = ?) 
and team.team_name like ? 
order by team.team_id desc limit ?;`
	teams := make([]*Team, 0)

	_, err := o.Raw(sql, bookId, "%"+teamName+"%", limit).QueryRows(&teams)

	if err != nil {
		logs.Error("查询团队项目时出错 ->", err)
		return nil, err
	}

	result := SelectMemberResult{}
	items := make([]KeyValueItem, 0)

	for _, team := range teams {
		item := KeyValueItem{}
		item.Id = team.TeamId
		item.Text = team.TeamName
		items = append(items, item)
	}
	result.Result = items

	return &result, err
}

//查询指定项目的团队.
func (m *TeamRelationship) FindByBookToPager(bookId, pageIndex, pageSize int) (list []*TeamRelationship, totalCount int, err error) {

	if bookId <= 0 {
		err = ErrInvalidParameter
		return
	}

	offset := (pageIndex - 1) * pageSize

	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", bookId).OrderBy("-team_relationship_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		logs.Error("查询团队项目时出错 ->", err)
		return
	}
	count, err := m.QueryTable().Filter("book_id", bookId).Count()

	if err != nil {
		logs.Error("查询团队项目时出错 ->", err)
		return
	}
	totalCount = int(count)
	for _, item := range list {
		item.Include()
	}
	return
}
