package models

import (
	"time"
	"github.com/lifei6671/godoc/conf"
)

// Document struct.
type Document struct {
	DocumentId int		`orm:"pk;auto;unique;column(document_id)" json:"document_id"`
	DocumentName string	`orm:"column(document_name);size(500)" json:"document_name"`
	// Identify 文档唯一标识
	Identify string		`orm:"column(identify);size(100);unique" json:"identify"`
	BookId int		`orm:"column(book_id);type(int);index" json:"book_id"`
	OrderSort int		`orm:"column(order_sort);default(0);type(int)" json:"order_sort"`
	// Markdown markdown格式文档.
	Markdown string		`orm:"column(markdown);type(longtext)" json:"markdown"`
	// Release 发布后的Html格式内容.
	Release string		`orm:"column(release);type(longtext)" json:"release"`
	// Content 未发布的 Html 格式内容.
	Content string		`orm:"column(content);type(longtext)" json:"content"`
	CreateTime time.Time	`orm:"column(create_time);type(datetime)" json:"create_time"`
	CreateAt int		`orm:"column(create_at);type(int)" json:"create_at"`
	ModifyTime time.Time	`orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt int		`orm:"column(modify_at);type(int)" json:"modify_at"`
	Version int64		`orm:"type(bigint);column(version)" json:"version"`
}


// TableName 获取对应数据库表名.
func (m *Document) TableName() string {
	return "documents"
}
// TableEngine 获取数据使用的引擎.
func (m *Document) TableEngine() string {
	return "INNODB"
}

func (m *Document) TableNameWithPrefix()  string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewDocument() *Document  {
	return &Document{}
}



















