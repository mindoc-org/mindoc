package models

import (
	"time"

	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/orm"
)

// Document struct.
type Document struct {
	DocumentId int		`orm:"pk;auto;unique;column(document_id)" json:"doc_id"`
	DocumentName string	`orm:"column(document_name);size(500)" json:"doc_name"`
	// Identify 文档唯一标识
	Identify string		`orm:"column(identify);size(100);unique;null;default(null)" json:"identify"`
	BookId int		`orm:"column(book_id);type(int);index" json:"book_id"`
	ParentId int 		`orm:"column(parent_id);type(int);index" json:"parent_id"`
	OrderSort int		`orm:"column(order_sort);default(0);type(int);index" json:"order_sort"`
	// Markdown markdown格式文档.
	Markdown string		`orm:"column(markdown);type(longtext)" json:"markdown"`
	// Release 发布后的Html格式内容.
	Release string		`orm:"column(release);type(longtext)" json:"release"`
	// Content 未发布的 Html 格式内容.
	Content string		`orm:"column(content);type(longtext)" json:"content"`
	CreateTime time.Time	`orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId int		`orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime time.Time	`orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt int		`orm:"column(modify_at);type(int)" json:"-"`
	Version int64		`orm:"type(bigint);column(version)" json:"-"`
}

type DocumentTree struct {
	DocumentId int               `json:"id,string"`
	DocumentName string 		`json:"text"`
	ParentId int                	`json:"parent_id,string"`
	State *DocumentSelected         `json:"state,omitempty"`
}
type DocumentSelected struct {
	Selected bool        	`json:"selected"`
	Opened bool        	`json:"opened"`
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

func (m *Document) Find(id int) (*Document,error) {
	if id <= 0 {
		return m,ErrInvalidParameter
	}
	o := orm.NewOrm()


	err := o.Read(m)

	if err == orm.ErrNoRows{
		return m,ErrDataNotExist
	}
	return m,nil
}


func (m *Document) FindDocumentTree(book_id int) ([]*DocumentTree,error){
	o := orm.NewOrm()

	trees := make([]*DocumentTree,0)

	var docs []*Document

	count ,err := o.QueryTable(m).Filter("book_id",book_id).OrderBy("-order_sort","document_id").All(&docs,"document_id","document_name","parent_id")

	if err != nil {
		return trees,err
	}

	trees = make([]*DocumentTree,count)

	for index,item := range docs {
		tree := &DocumentTree{}
		if index == 0{
			tree.State = &DocumentSelected{ Selected: true, Opened: true }
		}
		tree.DocumentId = item.DocumentId
		tree.ParentId = item.ParentId
		tree.DocumentName = item.DocumentName

		trees[index] = tree
	}

	return trees,nil
}

func (m *Document) InsertOrUpdate(cols... string) error {
	o := orm.NewOrm()

	if m.DocumentId > 0 {
		_,err := o.Update(m)
		return err
	}else{
		_,err := o.Insert(m)
		return err
	}
}

func (m *Document) FindByFieldFirst(field string,v interface{}) (*Document,error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter(field,v).One(m)

	return m,err
}











