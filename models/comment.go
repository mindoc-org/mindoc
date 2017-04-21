package models

import "time"

//Comment struct
type Comment struct {
	CommentId int			`orm:"pk;auto;unique;column(comment_id)" json:"comment_id"`
	// DocumentId 评论所属的文档.
	DocumentId int			`orm:"column(document_id);type(int)" json:"document_id"`
	// Author 评论作者.
	Author string			`orm:"column(author);size(100)" json:"author"`
	//MemberId 评论用户ID.
	MemberId int			`orm:"column(member_id);type(int)" json:"member_id"`
	// IPAddress 评论者的IP地址
	IPAddress string		`orm:"column(ip_address);size(100)" json:"ip_address"`
	// 评论日期.
	CommentDate time.Time		`orm:"type(datetime);column(comment_date);auto_now_add" json:"comment_date"`
	//Content 评论内容.
	Content string			`orm:"column(content);size(2000)" json:"content"`
	// Approved 评论状态：0 待审核/1 已审核/2 垃圾评论
	Approved int			`orm:"column(approved);type(int)" json:"approved"`
	// UserAgent 评论者浏览器内容
	UserAgent string		`orm:"column(user_agent);size(500)" json:"user_agent"`
	// Parent 评论所属父级
	ParentId int			`orm:"column(parent_id);type(int);default(0)" json:"parent_id"`
}

// TableName 获取对应数据库表名.
func (m *Comment) TableName() string {
	return "comments"
}
// TableEngine 获取数据使用的引擎.
func (m *Comment) TableEngine() string {
	return "INNODB"
}

