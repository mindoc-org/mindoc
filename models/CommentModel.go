package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/mindoc-org/mindoc/conf"
)

//Comment struct
type Comment struct {
	CommentId int `orm:"pk;auto;unique;column(comment_id)" json:"comment_id"`
	Floor     int `orm:"column(floor);type(unsigned);default(0)" json:"floor"`
	BookId    int `orm:"column(book_id);type(int)" json:"book_id"`
	// DocumentId 评论所属的文档.
	DocumentId int `orm:"column(document_id);type(int)" json:"document_id"`
	// Author 评论作者.
	Author string `orm:"column(author);size(100)" json:"author"`
	//MemberId 评论用户ID.
	MemberId int `orm:"column(member_id);type(int)" json:"member_id"`
	// IPAddress 评论者的IP地址
	IPAddress string `orm:"column(ip_address);size(100)" json:"ip_address"`
	// 评论日期.
	CommentDate time.Time `orm:"type(datetime);column(comment_date);auto_now_add" json:"comment_date"`
	//Content 评论内容.
	Content string `orm:"column(content);size(2000)" json:"content"`
	// Approved 评论状态：0 待审核/1 已审核/2 垃圾评论/ 3 已删除
	Approved int `orm:"column(approved);type(int)" json:"approved"`
	// UserAgent 评论者浏览器内容
	UserAgent string `orm:"column(user_agent);size(500)" json:"user_agent"`
	// Parent 评论所属父级
	ParentId     int `orm:"column(parent_id);type(int);default(0)" json:"parent_id"`
	AgreeCount   int `orm:"column(agree_count);type(int);default(0)" json:"agree_count"`
	AgainstCount int `orm:"column(against_count);type(int);default(0)" json:"against_count"`
	Index        int `orm:"-" json:"index"`
	ShowDel      int `orm:"-" json:"show_del"`
}

// TableName 获取对应数据库表名.
func (m *Comment) TableName() string {
	return "comments"
}

// TableEngine 获取数据使用的引擎.
func (m *Comment) TableEngine() string {
	return "INNODB"
}

func (m *Comment) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewComment() *Comment {
	return &Comment{}
}

// 是否有权限删除
func (m *Comment) CanDelete(user_memberid int, user_bookrole conf.BookRole) bool {
	return user_memberid == m.MemberId || user_bookrole == conf.BookFounder || user_bookrole == conf.BookAdmin
}

// 根据文档id查询文档评论
func (m *Comment) QueryCommentByDocumentId(doc_id, page, pagesize int, member *Member) (comments []Comment, count int64, ret_page int) {
	doc, err := NewDocument().Find(doc_id)
	if err != nil {
		return
	}

	o := orm.NewOrm()
	count, _ = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc_id).Count()
	if -1 == page {     // 请求最后一页
		var total int = int(count)
		if total % pagesize == 0 {
			page = total / pagesize
		} else {
			page = total / pagesize + 1
		}
	}
	offset := (page - 1) * pagesize
	ret_page = page
	o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc_id).OrderBy("comment_date").Offset(offset).Limit(pagesize).All(&comments)

	bookRole, _ := NewRelationship().FindForRoleId(doc.BookId, member.MemberId)
	for i := 0; i < len(comments); i++ {
		comments[i].Index = (i + 1) + (page - 1) * pagesize
		if comments[i].CanDelete(member.MemberId, bookRole) {
			comments[i].ShowDel = 1
		}
	}
	return
}

func (m *Comment) Update(cols ...string) error {
	o := orm.NewOrm()

	_, err := o.Update(m, cols...)

	return err
}

//Insert 添加一条评论.
func (m *Comment) Insert() error {
	if m.DocumentId <= 0 {
		return errors.New("评论文档不存在")
	}
	if m.Content == "" {
		return ErrCommentContentNotEmpty
	}

	o := orm.NewOrm()

	if m.CommentId > 0 {
		comment := NewComment()
		//如果父评论不存在
		if err := o.Read(comment); err != nil {
			return err
		}
	}

	document := NewDocument()
	//如果评论的文档不存在
	if _, err := document.Find(m.DocumentId); err != nil {
		return err
	}
	book, err := NewBook().Find(document.BookId)
	//如果评论的项目不存在
	if err != nil {
		return err
	}
	//如果已关闭评论
	if book.CommentStatus == "closed" {
		return ErrCommentClosed
	}
	if book.CommentStatus == "registered_only" && m.MemberId <= 0 {
		return ErrPermissionDenied
	}
	//如果仅参与者评论
	if book.CommentStatus == "group_only" {
		if m.MemberId <= 0 {
			return ErrPermissionDenied
		}
		rel := NewRelationship()
		if _, err := rel.FindForRoleId(book.BookId, m.MemberId); err != nil {
			return ErrPermissionDenied
		}
	}

	if m.MemberId > 0 {
		member := NewMember()
		//如果用户不存在
		if _, err := member.Find(m.MemberId); err != nil {
			return ErrMemberNoExist
		}
		//如果用户被禁用
		if member.Status == 1 {
			return ErrMemberDisabled
		}
	} else if m.Author == "" {
		m.Author = "[匿名用户]"
	}
	m.BookId = book.BookId
	_, err = o.Insert(m)

	return err
}

// 删除一条评论
func (m *Comment) Delete() error {
	o := orm.NewOrm()
	_, err := o.Delete(m)
	return err
}

func (m *Comment) Find(id int, cols ...string) (*Comment, error) {
	o := orm.NewOrm()
	if err := o.QueryTable(m.TableNameWithPrefix()).Filter("comment_id", id).One(m, cols...); err != nil {
		return m, err
	}
	return m, nil
}