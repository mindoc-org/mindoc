package models

import (
	"html/template"
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
				li := fmt.Sprintf("<li><a href=\"%s\" target=\"_blank\" title=\"%s\">%s</a></li>", attach.HttpPath, attach.Description, attach.Description)

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

	book, err := NewBook().Find(book_id)
	if err != nil {
		return
	}
	if book.LinkId == 0 {
		_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", book_id).All(&docs)
	} else {
		rdoc, err := m.FindByFieldFirst("book_id", book_id)
		if err == nil {
			doclinks := rdoc.Markdown
			beego.Info(doclinks)
			sql2 := `SELECT * 
				FROM md_documents  
				WHERE book_id = ? AND FIND_IN_SET(document_id,?)> 0
				ORDER BY order_sort  ,document_id  `
			_, err = o.Raw(sql2, book.LinkId, doclinks).QueryRows(&docs)
		}
	}
	return
}

// GetLinkBookDocuments ...
func (m *Document) GetLinkBookDocuments(book_id int) (doclinks string, doclist string, err error) {

	book, err := NewBook().Find(book_id)
	if err != nil {
		return "", "", err
	}

	rdoc, err := m.FindByFieldFirst("book_id", book_id)
	if err != nil {
		return "", "", err
	}
	doclinks = rdoc.Markdown

	buf := bytes.NewBufferString("")

	GetLinkBookDocumentsInternal(doclinks, book.LinkId, 0, buf)

	doclist = buf.String()
	return doclinks, doclist, nil
}

// GetLinkBookDocumentsInternal ...
func GetLinkBookDocumentsInternal(doclinks string, book_id, parent_id int, buf *bytes.Buffer) {
	var docs []*Document
	o := orm.NewOrm()
	sql2 := `SELECT document_id,document_name,FIND_IN_SET(document_id,?) AS modify_at
		FROM md_documents  
		WHERE book_id = ? AND parent_id= ?  
		ORDER BY order_sort  ,document_id  `
	count, _ := o.Raw(sql2, doclinks, book_id, parent_id).QueryRows(&docs)
	if count > 0 {
		buf.WriteString("<ul>\r\n")
		for _, item := range docs {
			buf.WriteString("<li><input type=\"checkbox\" id=\"")
			buf.WriteString(fmt.Sprintf("%d", item.DocumentId))
			buf.WriteString("\"  checked=\"checked\" ")
			buf.WriteString("/><label><input type=\"checkbox\" ")
			if item.ModifyAt > 0 {
				buf.WriteString(" checked=\"checked\" ")
			}
			buf.WriteString("/><span></span></label><label for=\"")
			buf.WriteString(fmt.Sprintf("%d", item.DocumentId))
			buf.WriteString("\">")
			buf.WriteString(template.HTMLEscapeString(item.DocumentName))
			buf.WriteString("</label>")
			GetLinkBookDocumentsInternal(doclinks, book_id, item.DocumentId, buf)
			buf.WriteString("</li>\r\n")
		}
		buf.WriteString("</ul>\r\n")
	}

}
