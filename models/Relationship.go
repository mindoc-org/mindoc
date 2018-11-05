package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
)

type Relationship struct {
	RelationshipId int `orm:"pk;auto;unique;column(relationship_id)" json:"relationship_id"`
	MemberId       int `orm:"column(member_id);type(int)" json:"member_id"`
	BookId         int `orm:"column(book_id);type(int)" json:"book_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId conf.BookRole `orm:"column(role_id);type(int)" json:"role_id"`
}

// TableName 获取对应数据库表名.
func (m *Relationship) TableName() string {
	return "relationship"
}
func (m *Relationship) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

// TableEngine 获取数据使用的引擎.
func (m *Relationship) TableEngine() string {
	return "INNODB"
}

// 联合唯一键
func (u *Relationship) TableUnique() [][]string {
	return [][]string{
		{"member_id", "book_id"},
	}
}

func NewRelationship() *Relationship {
	return &Relationship{}
}

func (m *Relationship) Find(id int) (*Relationship, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("relationship_id", id).One(m)
	return m, err
}

//查询指定项目的创始人.
func (m *Relationship) FindFounder(book_id int) (*Relationship, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("role_id", 0).One(m)

	return m, err
}

func (m *Relationship) UpdateRoleId(bookId, memberId int, roleId conf.BookRole) (*Relationship, error) {
	o := orm.NewOrm()
	book := NewBook()
	book.BookId = bookId

	if err := o.Read(book); err != nil {
		logs.Error("UpdateRoleId => ", err)
		return m, errors.New("项目不存在")
	}
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("member_id", memberId).Filter("book_id", bookId).One(m)

	if err == orm.ErrNoRows {
		m = NewRelationship()
		m.BookId = bookId
		m.MemberId = memberId
		m.RoleId = roleId
	} else if err != nil {
		return m, err
	} else if m.RoleId == conf.BookFounder {
		return m, errors.New("不能变更创始人的权限")
	}
	m.RoleId = roleId

	if m.RelationshipId > 0 {
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
	}

	return m, err

}

func (m *Relationship) FindForRoleId(book_id, member_id int) (conf.BookRole, error) {
	o := orm.NewOrm()

	relationship := NewRelationship()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("member_id", member_id).One(relationship)

	if err != nil {

		return 0, err
	}
	return relationship.RoleId, nil
}

func (m *Relationship) FindByBookIdAndMemberId(book_id, member_id int) (*Relationship, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("member_id", member_id).One(m)

	return m, err
}

func (m *Relationship) Insert() error {
	o := orm.NewOrm()

	_, err := o.Insert(m)

	return err
}

func (m *Relationship) Update() error {
	o := orm.NewOrm()

	_, err := o.Update(m)

	return err
}

func (m *Relationship) DeleteByBookIdAndMemberId(book_id, member_id int) error {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("member_id", member_id).One(m)

	if err == orm.ErrNoRows {
		return errors.New("用户未参与该项目")
	}
	if m.RoleId == conf.BookFounder {
		return errors.New("不能删除创始人")
	}
	_, err = o.Delete(m)

	if err != nil {
		logs.Error("删除项目参与者 => ", err)
		return errors.New("删除失败")
	}
	return nil

}

func (m *Relationship) Transfer(book_id, founder_id, receive_id int) error {
	o := orm.NewOrm()

	founder := NewRelationship()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("member_id", founder_id).One(founder)

	if err != nil {
		return err
	}
	if founder.RoleId != conf.BookFounder {
		return errors.New("转让者不是创始人")
	}
	receive := NewRelationship()

	err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).Filter("member_id", receive_id).One(receive)

	if err != orm.ErrNoRows && err != nil {
		return err
	}
	o.Begin()

	founder.RoleId = conf.BookAdmin

	receive.MemberId = receive_id
	receive.RoleId = conf.BookFounder
	receive.BookId = book_id

	if err := founder.Update(); err != nil {
		o.Rollback()
		return err
	}
	if receive.RelationshipId > 0 {
		if _, err := o.Update(receive); err != nil {
			o.Rollback()
			return err
		}
	} else {
		if _, err := o.Insert(receive); err != nil {
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}
