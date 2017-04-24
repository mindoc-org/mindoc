package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/logs"
)

// Book struct .
type Book struct {
	BookId int		`orm:"pk;auto;unique;column(book_id)" json:"book_id"`
	// BookName 项目名称.
	BookName string		`orm:"column(book_name);size(500)" json:"book_name"`
	// Identify 项目唯一标识.
	Identify string		`orm:"column(identify);size(100);unique" json:"identify"`
	OrderIndex int		`orm:"column(order_index);type(int);default(0)" json:"order_index"`
	// Description 项目描述.
	Description string	`orm:"column(description);size(2000)" json:"description"`
	Label string		`orm:"column(label);size(500)" json:"label"`
	// PrivatelyOwned 项目私有： 0 公开/ 1 私有
	PrivatelyOwned int	`orm:"column(privately_owned);type(int);default(0)" json:"privately_owned"`
	// 当项目是私有时的访问Token.
	PrivateToken string 	`orm:"column(private_token);size(500);null" json:"private_token"`
	//评论状态：0 正常/1 已删除
	Status int 		`orm:"column(status);type(int);default(0)" json:"status"`
	// DocCount 包含文档数量.
	DocCount int		`orm:"column(doc_count);type(int)" json:"doc_count"`
	// CommentStatus 评论设置的状态:open 为允许所有人评论，closed 为不允许评论, group_only 仅允许参与者评论 ,registered_only 仅允许注册者评论.
	CommentStatus string	`orm:"column(comment_status);size(20);default(open)" json:"comment_status"`
	CommentCount int	`orm:"column(comment_count);type(int)" json:"comment_count"`
	Cover string 		`orm:"column();size(1000)" json:"cover"`

	// CreateTime 创建时间 .
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	MemberId int		`orm:"column(member_id);size(100)" json:"member_id"`
	ModifyTime time.Time	`orm:"type(datetime);column(modify_time);null;auto_now" json:"modify_time"`
	Version int64		`orm:"type(bigint);column(version)" json:"version"`
}


// TableName 获取对应数据库表名.
func (m *Book) TableName() string {
	return "books"
}
// TableEngine 获取数据使用的引擎.
func (m *Book) TableEngine() string {
	return "INNODB"
}
func (m *Book) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewBook() *Book {
	return &Book{}
}

func (m *Book) Insert() error {
	o := orm.NewOrm()
	_,err := o.Insert(m)

	if err == nil {
		relationship := NewRelationship()
		relationship.BookId = m.BookId
		relationship.RoleId = 0
		relationship.MemberId = m.MemberId
		err = relationship.Insert()
	}

	return err
}

func (m *Book) Find(id int) error {
	if id <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	return o.Read(m)
}

func (m *Book) Update(cols... string) error  {
	o := orm.NewOrm()

	_,err := o.Update(m,cols...)
	return err
}

func (m *Book) FindByField(field string,value interface{}) ([]Book,error)  {
	o := orm.NewOrm()

	var books []Book
	_,err := o.QueryTable(conf.GetDatabasePrefix() + m.TableName()).Filter(field,value).All(&books)

	return books,err
}

func (m *Book) FindByFieldFirst(field string,value interface{})(*Book,error) {
	o := orm.NewOrm()

	err := o.QueryTable(conf.GetDatabasePrefix() + m.TableName()).Filter(field,value).One(m)

	return m,err

}
func (m *Book) FindToPager(pageIndex, pageSize ,memberId int) (books []*BookResult,totalCount int,err error){

	relationship := NewRelationship()

	o := orm.NewOrm()

	qb, _ := orm.NewQueryBuilder("mysql")

	qb.Select("COUNT(book.book_id) AS total_count").
		From(m.TableNameWithPrefix() + " AS book").
		LeftJoin(relationship.TableNameWithPrefix() + " AS rel").
		On("book.book_id=rel.book_id").
		Where("rel.member_id=?")

	err = o.Raw(qb.String(),memberId).QueryRow(&totalCount)

	if err != nil {
		return
	}

	offset := (pageIndex - 1) * pageSize
	qb2,_ := orm.NewQueryBuilder("mysql")

	qb2.Select("book.*,rel.member_id","rel.role_id","m.account as create_name").
		From(m.TableNameWithPrefix() + " AS book").
		LeftJoin(relationship.TableNameWithPrefix() + " AS rel").
		On("book.book_id=rel.book_id").
		LeftJoin(NewMember().TableNameWithPrefix() + " AS m").On("rel.member_id=m.member_id AND rel.role_id=0").
		Where("rel.member_id=?").
		OrderBy("book.order_index DESC ","book.book_id").Desc().
		Limit(pageSize).
		Offset(offset)

	//logs.Info("",qb2.String())
	_,err = o.Raw(qb2.String(),memberId).QueryRows(&books)
	if err != nil {
		logs.Error("分页查询项目列表 => ",err)
		return
	}
	sql := "SELECT m.account,doc.modify_time FROM md_documents AS doc LEFT JOIN md_members AS m ON doc.modify_at=m.member_id WHERE book_id = ? LIMIT 1 ORDER BY doc.modify_time DESC"

	if err == nil && len(books) > 0{
		for index,book := range books  {
			var text struct{
				Account string
				ModifyTime time.Time
			}


			err1 := o.Raw(sql,book.BookId).QueryRow(&text)
			if err1 == nil {
				books[index].LastModifyText = text.Account + " 于 " + text.ModifyTime.Format("2006-01-02 15:04:05")
			}
			if book.RoleId == 0{
				book.RoleName = "创始人"
			}else if book.RoleId == 1 {
				book.RoleName = "管理员"
			}else if book.RoleId == 2 {
				book.RoleName = "编辑者"
			}else if book.RoleId == 3 {
				book.RoleName = "观察者"
			}
		}
	}
	return
}

// 彻底删除项目.
func (m *Book) ThoroughDeleteBook(id int) error {
	if id <= 0{
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	m.BookId = id
	if err := o.Read(m); err != nil {
		return err
	}
	o.Begin()
	sql1 := "DELETE FROM " + NewComment().TableNameWithPrefix() + " WHERE book_id = ?"

	_,err := o.Raw(sql1,m.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	sql2 := "DELETE FROM " + NewDocument().TableNameWithPrefix() + " WHERE book_id = ?"

	_,err = o.Raw(sql2,m.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	sql3 := "DELETE FROM " + m.TableNameWithPrefix() + " WHERE book_id = ?"

	_,err = o.Raw(sql3).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	sql4 := "DELETE FROM " + NewRelationship().TableNameWithPrefix() + " WHERE book_id = ?"

	_,err = o.Raw(sql4).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()

}
























































