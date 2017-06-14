package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/mindoc/conf"
)

type BookResult struct {
	BookId int              `json:"book_id"`
	BookName string         `json:"book_name"`
	Identify string         `json:"identify"`
	OrderIndex int          `json:"order_index"`
	Description string      `json:"description"`
	PrivatelyOwned int      `json:"privately_owned"`
	PrivateToken string     `json:"private_token"`
	DocCount int            `json:"doc_count"`
	CommentStatus string    `json:"comment_status"`
	CommentCount int        `json:"comment_count"`
	CreateTime time.Time    `json:"create_time"`
	CreateName string 	`json:"create_name"`
	ModifyTime time.Time	`json:"modify_time"`
	Cover string            `json:"cover"`
	Theme string		`json:"theme"`
	Label string		`json:"label"`
	MemberId int            `json:"member_id"`
	Editor string           `json:"editor"`

	RelationshipId int	`json:"relationship_id"`
	RoleId int        	`json:"role_id"`
	RoleName string 	`json:"role_name"`
	Status int

	LastModifyText string 	`json:"last_modify_text"`
	IsDisplayComment bool	`json:"is_display_comment"`
}

func NewBookResult() *BookResult {
	return &BookResult{}
}


// 根据项目标识查询项目以及指定用户权限的信息.
func (m *BookResult) FindByIdentify(identify string,member_id int) (*BookResult,error) {
	if identify == "" || member_id <= 0 {
		return m,ErrInvalidParameter
	}
	o := orm.NewOrm()

	book := NewBook()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("identify", identify).One(book)

	if err != nil {
		return m, err
	}

	relationship := NewRelationship()

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id", book.BookId).Filter("member_id", member_id).One(relationship)

	if err != nil {
		return m, err
	}
	var relationship2 Relationship

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id", book.BookId).Filter("role_id", 0).One(&relationship2)

	if err != nil {
		logs.Error("根据项目标识查询项目以及指定用户权限的信息 => ", err)
		return m, ErrPermissionDenied
	}

	member, err := NewMember().Find(relationship2.MemberId)
	if err != nil {
		return m, err
	}

	m = book.ToBookResult()

	m.CreateName 		= member.Account
	m.MemberId 		= relationship.MemberId
	m.RoleId		= relationship.RoleId
	m.RelationshipId	= relationship.RelationshipId

	if m.RoleId == conf.BookFounder {
		m.RoleName = "创始人"
	} else if m.RoleId == conf.BookAdmin {
		m.RoleName = "管理员"
	} else if m.RoleId == conf.BookEditor {
		m.RoleName = "编辑者"
	} else if m.RoleId == conf.BookObserver {
		m.RoleName = "观察者"
	}

	doc := NewDocument()

	err = o.QueryTable(doc.TableNameWithPrefix()).Filter("book_id", book.BookId).OrderBy("modify_time").One(doc)

	if err == nil {
		member2 := NewMember()
		member2.Find(doc.ModifyAt)

		m.LastModifyText = member2.Account + " 于 " + doc.ModifyTime.Format("2006-01-02 15:04:05")
	}

	return m, nil
}

func (m *BookResult) FindToPager(pageIndex, pageSize int) (books []*BookResult,totalCount int,err error)  {
	o := orm.NewOrm()

	count,err := o.QueryTable(NewBook().TableNameWithPrefix()).Count()

	if err != nil {
		return
	}
	totalCount = int(count)

	sql := `SELECT
			book.*,rel.relationship_id,rel.role_id,m.account AS create_name
		FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS m ON rel.member_id = m.member_id
		ORDER BY book.order_index DESC ,book.book_id DESC  LIMIT ?,?`

	offset := (pageIndex -1 )* pageSize

	_,err = o.Raw(sql,offset,pageSize).QueryRows(&books)

	return
}



















