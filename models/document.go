package models

import (
	"time"

	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
)

// Document struct.
type Document struct {
	DocumentId   int    `orm:"pk;auto;unique;column(document_id)" json:"doc_id"`
	DocumentName string `orm:"column(document_name);size(500)" json:"doc_name"`
	// Identify 文档唯一标识
	Identify  string `orm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId    int    `orm:"column(book_id);type(int);index" json:"book_id"`
	ParentId  int    `orm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	OrderSort int    `orm:"column(order_sort);default(0);type(int);index" json:"order_sort"`
	// Markdown markdown格式文档.
	Markdown string `orm:"column(markdown);type(text);null" json:"markdown"`
	// Release 发布后的Html格式内容.
	Release string `orm:"column(release);type(text);null" json:"release"`
	// Content 未发布的 Html 格式内容.
	Content    string        `orm:"column(content);type(text);null" json:"content"`
	CreateTime time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId   int           `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime time.Time     `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt   int           `orm:"column(modify_at);type(int)" json:"-"`
	Version    int64         `orm:"type(bigint);column(version)" json:"version"`
	AttachList []*Attachment `orm:"-" json:"attach"`
}

// TableName 获取对应数据库表名.
func (m *Document) TableName() string {
	return "documents"
}

// TableEngine 获取数据使用的引擎.
func (m *Document) TableEngine() string {
	return "INNODB"
}

func (m *Document) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}

//根据文档ID查询指定文档.
func (m *Document) Find(id int) (*Document, error) {
	if id <= 0 {
		return m, ErrInvalidParameter
	}
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", id).One(m)

	if err == orm.ErrNoRows {
		return m, ErrDataNotExist
	}
	return m, nil
}

//插入和更新文档.
func (m *Document) InsertOrUpdate(cols ...string) error {
	o := orm.NewOrm()

	if m.DocumentId > 0 {
		_, err := o.Update(m)
		return err
	} else {
		_, err := o.Insert(m)
		NewBook().ResetDocumentNumber(m.BookId)
		return err
	}
}

//根据指定字段查询一条文档.
func (m *Document) FindByFieldFirst(field string, v interface{}) (*Document, error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter(field, v).One(m)

	return m, err
}

//递归删除一个文档.
func (m *Document) RecursiveDocument(doc_id int) error {

	o := orm.NewOrm()

	if doc, err := m.Find(doc_id); err == nil {
		o.Delete(doc)
		NewDocumentHistory().Clear(doc.DocumentId)
	}

	var docs []*Document

	_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("parent_id", doc_id).All(&docs)

	if err != nil {
		beego.Error("RecursiveDocument => ", err)
		return err
	}

	for _, item := range docs {
		doc_id := item.DocumentId
		o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc_id).Delete()
		m.RecursiveDocument(doc_id)
	}

	return nil
}

//发布文档
func (m *Document) ReleaseContent(book_id int) {

	o := orm.NewOrm()

	var docs []*Document
	_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).All(&docs, "document_id", "content")

	if err != nil {
		beego.Error("发布失败 => ", err)
		return
	}
	for _, item := range docs {
		item.Release = item.Content
		attach_list, err := NewAttachment().FindListByDocumentId(item.DocumentId)
		if err == nil && len(attach_list) > 0 {
			content := bytes.NewBufferString("<div class=\"attach-list\"><strong>附件</strong><ul>")
			for _, attach := range attach_list {
				li := fmt.Sprintf("<li><a href=\"%s\" target=\"_blank\" title=\"%s\">%s</a></li>", attach.HttpPath, attach.FileName, attach.FileName)

				content.WriteString(li)
			}
			content.WriteString("</ul></div>")
			item.Release += content.String()
		}
		_, err = o.Update(item, "release")
		if err != nil {
			beego.Error(fmt.Sprintf("发布失败 => %+v", item), err)
		}
	}
}

//根据项目ID查询文档列表.
func (m *Document) FindListByBookId(book_id int) (docs []*Document, err error) {
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).All(&docs)

	return
}
