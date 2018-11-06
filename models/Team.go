package models

import (
	"time"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"errors"
	"github.com/astaxie/beego"
)

//团队.
type Team struct {
	TeamId      int       `orm:"column(team_id);pk;auto;unique;" json:"team_id"`
	TeamName    string    `orm:"column(team_name);size(255)" json:"team_name"`
	MemberId    int       `orm:"column(member_id);type(int);" json:"member_id"`
	IsDelete    bool      `orm:"column(is_delete);default(0)" json:"is_delete"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberCount int       `orm:"-" json:"member_count"`
	BookCount   int       `orm:"-" json:"book_count"`
	MemberName  string    `orm:"-" json:"member_name"`
}

// TableName 获取对应数据库表名.
func (t *Team) TableName() string {
	return "teams"
}

// TableEngine 获取数据使用的引擎.
func (t *Team) TableEngine() string {
	return "INNODB"
}

func (t *Team) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + t.TableName()
}

func NewTeam() *Team {
	return &Team{}
}

// 查询一个团队.
func (t *Team) First(id int, cols ...string) (*Team, error) {
	if id <= 0 {
		return nil, orm.ErrNoRows
	}
	o := orm.NewOrm()
	err := o.QueryTable(t.TableNameWithPrefix()).Filter("team_id", id).One(t, cols...)

	return t, err
}

func (t *Team) Delete(id int) (err error) {
	if id <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	err = o.Begin()

	if err != nil {
		beego.Error("开启事物时出错 ->",err)
		return
	}
	_,err = o.QueryTable(t.TableNameWithPrefix()).Filter("team_id",id).Delete()

	if err != nil {
		beego.Error("删除团队时出错 ->", err)
		o.Rollback()
		return
	}

	_,err = o.Raw("delete from md_team_member where team_id=?;", id).Exec()

	if err != nil {
		beego.Error("删除团队成员时出错 ->", err)
		o.Rollback()
		return
	}

	_,err = o.Raw("delete from md_team_relationship where team_id=?;",id).Exec()

	if err != nil {
		beego.Error("删除团队项目时出错 ->", err)
		o.Rollback()
		return err
	}

	err = o.Commit()
	return
}

//分页查询团队.
func (t *Team) FindToPager(pageIndex, pageSize int) (list []*Team, totalCount int, err error) {
	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize

	_, err = o.QueryTable(t.TableNameWithPrefix()).OrderBy("-team_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		return
	}

	c, err := o.QueryTable(t.TableNameWithPrefix()).Count()
	if err != nil {
		return
	}
	totalCount = int(c)

	for i,item := range list {
		if member,err := NewMember().Find(item.MemberId,"account","real_name"); err == nil {
			if member.RealName != "" {
				list[i].MemberName = member.RealName
			} else {
				list[i].MemberName = member.Account
			}
		}
		if c,err := o.QueryTable(NewTeamRelationship().TableNameWithPrefix()).Filter("team_id", item.TeamId).Count(); err == nil {
			list[i].BookCount = int(c)
		}
		if c,err := o.QueryTable(NewTeamMember().TableNameWithPrefix()).Filter("team_id", item.TeamId).Count(); err == nil {
			list[i].MemberCount = int(c)
		}
	}
	return
}

//更新或添加一个团队.
func (t *Team) Save(cols ... string) (err error) {
	if t.TeamName == "" {
		return NewError(5001, "团队名称不能为空")
	}

	o := orm.NewOrm()

	if t.TeamId <= 0 && o.QueryTable(t.TableNameWithPrefix()).Filter("team_name", t.TeamName).Exist() {
		return errors.New("团队名称已存在")
	}
	if t.TeamId <= 0 {
		_, err = o.Insert(t)
	} else {
		_, err = o.Update(t, cols...)
	}
	if err != nil {
		beego.Error("在保存团队时出错 ->", err)
	}
	return
}
