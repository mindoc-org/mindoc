//数据库模型.
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/conf"
	"os"
	"github.com/astaxie/beego"
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

	_,err := o.Delete(m)

	if err == nil {
		if err1 := os.Remove(m.FilePath); err1 != nil {
			beego.Error(err1)
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

func (m *Attachment) FindListByDocumentId(doc_id int) (attaches []*Attachment, err error) {
	o := orm.NewOrm()

	_,err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id",doc_id).OrderBy("-attachment_id").All(&attaches)
	return
}
