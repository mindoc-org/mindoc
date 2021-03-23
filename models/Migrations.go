package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/mindoc-org/mindoc/conf"
)

type Migration struct {
	MigrationId int       `orm:"column(migration_id);pk;auto;unique;" json:"migration_id"`
	Name        string    `orm:"column(name);size(500)" json:"name"`
	Statements  string    `orm:"column(statements);type(text);null" json:"statements"`
	Status      string    `orm:"column(status);default(update)" json:"status"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	Version     int64     `orm:"type(bigint);column(version);unique" json:"version"`
}

// TableName 获取对应数据库表名.
func (m *Migration) TableName() string {
	return "migrations"
}

// TableEngine 获取数据使用的引擎.
func (m *Migration) TableEngine() string {
	return "INNODB"
}

func (m *Migration) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewMigration() *Migration {
	return &Migration{}
}

func (m *Migration) FindFirst() (*Migration, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).OrderBy("-migration_id").One(m)

	return m, err
}
