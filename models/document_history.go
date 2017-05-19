package models

import (
	"time"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/orm"
)

type DocumentHistory struct {
	HistoryId    int       `orm:"column(history_id);pk;auto;unique" json:"history_id"`
	DocumentId   int       `orm:"column(document_id);type(int);index" json:"doc_id"`
	DocumentName string    `orm:"column(document_name);size(500)" json:"doc_name"`
	ParentId     int       `orm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	Markdown     string    `orm:"column(markdown);type(text);null" json:"markdown"`
	Content      string    `orm:"column(content);type(text);null" json:"content"`
	MemberId     int       `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt     int       `orm:"column(modify_at);type(int)" json:"-"`
	Version      int64     `orm:"type(bigint);column(version)" json:"version"`
}

// TableName 获取对应数据库表名.
func (m *DocumentHistory) TableName() string {
	return "document_history"
}

// TableEngine 获取数据使用的引擎.
func (m *DocumentHistory) TableEngine() string {
	return "INNODB"
}

func (m *DocumentHistory) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}


func (m *DocumentHistory) FindToPager(doc_id,page_index,page_size int) (docs []*DocumentHistory,totalCount int,err error) {

	o := orm.NewOrm()

	offset := (page_index - 1) * page_size

	totalCount = 0
	_,err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id",doc_id).Offset(offset).Limit(page_size).All(docs)

	if err != nil {
		return
	}
	var count int64
	count,err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id",doc_id).Count()

	if err != nil {
		return
	}
	totalCount = int(count)

	return
}