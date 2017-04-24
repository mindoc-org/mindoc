package models

import (
	"time"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/orm"
)

// Attachment struct .
type Attachment struct {
	AttachmentId int	`orm:"column(attachment_id);pk;auto;unique" json:"attachment_id"`
	DocumentId int		`orm:"column(document_id);type(int)" json:"document_id"`
	FileName string		`orm:"column(file_name);size(2000)" json:"file_name"`
	FileSize float64	`orm:"column(file_size);type(float)" json:"file_size"`
	FileExt string 		`orm:"column(file_ext);size(50)" json:"file_ext"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int		`orm:"column(create_at);type(int)" json:"create_at"`
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

func NewAttachment() *Attachment  {
	return &Attachment{}
}

func (m *Attachment) Insert() error  {
	o := orm.NewOrm()

	_,err := o.Insert(m)

	return err
}

func (m *Attachment) Find(id int) (*Attachment,error) {
	if id <= 0 {
		return m,ErrInvalidParameter
	}
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("attachment_id",id).One(m)

	return m,err
}








