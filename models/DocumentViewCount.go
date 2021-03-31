package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/mindoc-org/mindoc/conf"
)

type DocumentViewCount struct {
	DocumentId	int	`orm:"pk;column(document_id);type(int)" json:"doc_id"`
	ViewCount	int	`orm:"column(view_count);type(int)" json:"view_count"`
}

// TableName 获取对应数据库表名.
func (v *DocumentViewCount) TableName() string {
	return "document_viewcount"
}

// TableEngine 获取数据使用的引擎.
func (v *DocumentViewCount) TableEngine() string {
	return "INNODB"
}

func (v *DocumentViewCount) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + v.TableName()
}

func NewDocumentViewCount() *DocumentViewCount {
	return &DocumentViewCount{}
}

func (v *DocumentViewCount) IncrViewCount(id int) int {
	o := orm.NewOrm()
	num, _ := o.QueryTable(v.TableNameWithPrefix()).Filter("document_id", id).Update(orm.Params{
		"view_count": orm.ColValue(orm.ColAdd, 1),
	})
	if 0 == num {
		v.DocumentId = id
		v.ViewCount = 1
		num, _ = o.Insert(v)
	} else {
		o.QueryTable(v.TableNameWithPrefix()).Filter("document_id", id).One(v)
	}
	return v.ViewCount
}
