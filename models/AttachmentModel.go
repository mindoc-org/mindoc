//数据库模型.
package models

import (
	"time"

	"os"

	"strings"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/utils/filetil"
)

// Attachment struct .
type Attachment struct {
	AttachmentId int       `orm:"column(attachment_id);pk;auto;unique" json:"attachment_id"`
	BookId       int       `orm:"column(book_id);type(int)" json:"book_id"`
	DocumentId   int       `orm:"column(document_id);type(int);null" json:"doc_id"`
	FileName     string    `orm:"column(file_name);size(255)" json:"file_name"`
	FilePath     string    `orm:"column(file_path);size(2000)" json:"file_path"`
	FileSize     float64   `orm:"column(file_size);type(float)" json:"file_size"`
	HttpPath     string    `orm:"column(http_path);size(2000)" json:"http_path"`
	FileExt      string    `orm:"column(file_ext);size(50)" json:"file_ext"`
	CreateTime   time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt     int       `orm:"column(create_at);type(int)" json:"create_at"`
}

// TableName 获取对应数据库表名.
func (m *Attachment) TableName() string {
	return "attachment"
}

// TableEngine 获取数据使用的引擎.
func (m *Attachment) TableEngine() string {
	return "INNODB"
}
func (m *Attachment) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewAttachment() *Attachment {
	return &Attachment{}
}

func (m *Attachment) Insert() error {
	o := orm.NewOrm()

	_, err := o.Insert(m)

	return err
}
func (m *Attachment) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(m)
	return err
}

func (m *Attachment) Delete() error {
	o := orm.NewOrm()

	_, err := o.Delete(m)

	if err == nil {
		if err1 := os.Remove(m.FilePath); err1 != nil {
			logs.Error(err1)
		}
	}

	return err
}

func (m *Attachment) Find(id int) (*Attachment, error) {
	if id <= 0 {
		return m, ErrInvalidParameter
	}
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("attachment_id", id).One(m)

	return m, err
}

//查询指定文档的附件列表
func (m *Attachment) FindListByDocumentId(docId int) (attaches []*Attachment, err error) {
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", docId).Filter("book_id__gt", 0).OrderBy("-attachment_id").All(&attaches)
	return
}

//分页查询附件
func (m *Attachment) FindToPager(pageIndex, pageSize int) (attachList []*AttachmentResult, totalCount int, err error) {
	o := orm.NewOrm()

	total, err := o.QueryTable(m.TableNameWithPrefix()).Count()

	if err != nil {

		return nil, 0, err
	}
	totalCount = int(total)
	offset := (pageIndex - 1) * pageSize

	var list []*Attachment

	_, err = o.QueryTable(m.TableNameWithPrefix()).OrderBy("-attachment_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		if err == orm.ErrNoRows {
			logs.Info("没有查到附件 ->", err)
			err = nil
		}
		return
	}

	for _, item := range list {
		attach := &AttachmentResult{}
		attach.Attachment = *item
		attach.FileShortSize = filetil.FormatBytes(int64(attach.FileSize))
		//当项目ID为0标识是文章的附件
		if item.BookId == 0 && item.DocumentId > 0 {
			blog := NewBlog()
			if err := o.QueryTable(blog.TableNameWithPrefix()).Filter("blog_id", item.DocumentId).One(blog, "blog_title"); err == nil {
				attach.BookName = blog.BlogTitle
			} else {
				attach.BookName = "[文章不存在]"
			}
		} else {
			book := NewBook()

			if e := o.QueryTable(book.TableNameWithPrefix()).Filter("book_id", item.BookId).One(book, "book_name"); e == nil {
				attach.BookName = book.BookName

				doc := NewDocument()

				if e := o.QueryTable(doc.TableNameWithPrefix()).Filter("document_id", item.DocumentId).One(doc, "document_name"); e == nil {
					attach.DocumentName = doc.DocumentName
				} else {
					attach.DocumentName = "[文档不存在]"
				}

			} else {
				attach.BookName = "[项目不存在]"
			}
		}
		attach.LocalHttpPath = strings.Replace(item.FilePath, "\\", "/", -1)

		attachList = append(attachList, attach)
	}

	return
}
