package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/logs"
	"strings"
	"github.com/astaxie/beego"
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
	//状态：0 正常/1 已删除
	Status int 		`orm:"column(status);type(int);default(0)" json:"status"`
	//默认的编辑器.
	Editor string		`orm:"column(editor);size(50)" json:"editor"`
	// DocCount 包含文档数量.
	DocCount int		`orm:"column(doc_count);type(int)" json:"doc_count"`
	// CommentStatus 评论设置的状态:open 为允许所有人评论，closed 为不允许评论, group_only 仅允许参与者评论 ,registered_only 仅允许注册者评论.
	CommentStatus string	`orm:"column(comment_status);size(20);default(open)" json:"comment_status"`
	CommentCount int	`orm:"column(comment_count);type(int)" json:"comment_count"`
	//封面地址
	Cover string 		`orm:"column(cover);size(1000)" json:"cover"`
	//主题风格
	Theme string 		`orm:"column(theme);size(255);default(default)" json:"theme"`
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
//	o.Begin()

	_,err := o.Insert(m)

	if err == nil {
		relationship := NewRelationship()
		relationship.BookId = m.BookId
		relationship.RoleId = 0
		relationship.MemberId = m.MemberId
		err = relationship.Insert()
		if err != nil {
			logs.Error("插入项目与用户关联 => ",err)
			//o.Rollback()
			return err
		}
		document := NewDocument()
		document.BookId = m.BookId
		document.DocumentName = "空白文档"
		document.MemberId = m.MemberId
		err = document.InsertOrUpdate()
		if err != nil{
			//o.Rollback()
			return err
		}
		//o.Commit()
		return nil
	}
	//o.Rollback()
	return err
}

func (m *Book) Find(id int) (*Book,error) {
	if id <= 0 {
		return m,ErrInvalidParameter
	}
	o := orm.NewOrm()


	err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id",id).One(m)

	return m,err
}

func (m *Book) Update(cols... string) error  {
	o := orm.NewOrm()

	_,err := o.Update(m,cols...)
	return err
}

//根据指定字段查询结果集.
func (m *Book) FindByField(field string,value interface{}) ([]*Book,error)  {
	o := orm.NewOrm()

	var books []*Book
	_,err := o.QueryTable(m.TableNameWithPrefix()).Filter(field,value).All(&books)

	return books,err
}

//根据指定字段查询一个结果.
func (m *Book) FindByFieldFirst(field string,value interface{})(*Book,error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter(field,value).One(m)

	return m,err

}

func (m *Book) FindByIdentify(identify string) (*Book,error) {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("identify",identify).One(m)

	return m,err
}

//分页查询指定用户的项目
func (m *Book) FindToPager(pageIndex, pageSize ,memberId int) (books []*BookResult,totalCount int,err error){

	relationship := NewRelationship()

	o := orm.NewOrm()

	qb, _ := orm.NewQueryBuilder("mysql")

	qb.Select("COUNT(book.book_id) AS total_count").
		From(m.TableNameWithPrefix() + " AS book").
		LeftJoin(relationship.TableNameWithPrefix() + " AS rel").
		On("book.book_id=rel.book_id AND rel.member_id = ?").
		Where("rel.relationship_id > 0")

	err = o.Raw(qb.String(),memberId).QueryRow(&totalCount)

	if err != nil {
		return
	}

	offset := (pageIndex - 1) * pageSize
	qb2,_ := orm.NewQueryBuilder("mysql")

	qb2.Select("book.*,rel.member_id","rel.role_id","m.account as create_name").
		From(m.TableNameWithPrefix() + " AS book").
		LeftJoin(relationship.TableNameWithPrefix() + " AS rel").On("book.book_id=rel.book_id AND rel.member_id = ?").
		LeftJoin(relationship.TableNameWithPrefix() + " AS rel1").On("book.book_id=rel1.book_id  AND rel1.role_id=0").
		LeftJoin(NewMember().TableNameWithPrefix() + " AS m").On("rel1.member_id=m.member_id").
		Where("rel.relationship_id > 0").
		OrderBy("book.order_index DESC ","book.book_id").Desc().
		Limit(pageSize).
		Offset(offset)

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
	//sql1 := "DELETE FROM " + NewComment().TableNameWithPrefix() + " WHERE book_id = ?"
	//
	//_,err := o.Raw(sql1,m.BookId).Exec()
	//
	//if err != nil {
	//	o.Rollback()
	//	return err
	//}
	sql2 := "DELETE FROM " + NewDocument().TableNameWithPrefix() + " WHERE book_id = ?"

	_,err := o.Raw(sql2,m.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	sql3 := "DELETE FROM " + m.TableNameWithPrefix() + " WHERE book_id = ?"

	_,err = o.Raw(sql3,m.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	sql4 := "DELETE FROM " + NewRelationship().TableNameWithPrefix() + " WHERE book_id = ?"

	_,err = o.Raw(sql4,m.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()

}

//分页查找系统首页数据.
func (m *Book) FindForHomeToPager(pageIndex, pageSize ,member_id int) (books []*BookResult,totalCount int,err error) {
	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize
	//如果是登录用户
	if member_id > 0 {
		sql1 := "SELECT COUNT(*) FROM md_books AS book LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ? WHERE relationship_id > 0 OR book.privately_owned = 0"

		err = o.Raw(sql1,member_id).QueryRow(&totalCount)
		if err != nil {
			return
		}
		sql2 := `SELECT book.*,rel1.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
			LEFT JOIN md_relationship AS rel1 ON rel1.book_id = book.book_id AND rel1.role_id = 0
			LEFT JOIN md_members AS member ON rel1.member_id = member.member_id
			WHERE rel.relationship_id > 0 OR book.privately_owned = 0 ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_,err = o.Raw(sql2,member_id,offset,pageSize).QueryRows(&books)

		return

	}else{
		count,err1 := o.QueryTable(m.TableNameWithPrefix()).Filter("privately_owned",0).Count()

		if err1 != nil {
			err = err1
			return
		}
		totalCount = int(count)

		sql := `SELECT book.*,rel.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS member ON rel.member_id = member.member_id
			WHERE book.privately_owned = 0 ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_,err = o.Raw(sql,offset,pageSize).QueryRows(&books)

		return

	}

}

func (book *Book) ToBookResult() *BookResult {

	m := NewBookResult()

	m.BookId 	 	= book.BookId
	m.BookName 	 	= book.BookName
	m.Identify 	 	= book.Identify
	m.OrderIndex 	 	= book.OrderIndex
	m.Description 	 	= strings.Replace(book.Description, "\r\n", "<br/>", -1)
	m.PrivatelyOwned 	= book.PrivatelyOwned
	m.PrivateToken 		= book.PrivateToken
	m.DocCount 		= book.DocCount
	m.CommentStatus 	= book.CommentStatus
	m.CommentCount 		= book.CommentCount
	m.CreateTime 		= book.CreateTime
	m.ModifyTime 		= book.ModifyTime
	m.Cover 		= book.Cover
	m.Label 		= book.Label
	m.Status 		= book.Status
	m.Editor 		= book.Editor
	m.Theme			= book.Theme


	if book.Theme == ""{
		m.Theme = "default"
	}
	if book.Editor == "" {
		m.Editor = "markdown"
	}
	return m
}

//重置文档数量
func (m *Book) ResetDocumentNumber(book_id int)  {
	o := orm.NewOrm()

	totalCount,err := o.QueryTable(NewDocument().TableNameWithPrefix()).Filter("book_id",book_id).Count()
	if err == nil {
		o.Raw("UPDATE md_books SET doc_count = ? WHERE book_id = ?",int(totalCount),book_id).Exec()
	}else{
		beego.Error(err)
	}
}
















































