package models

import "time"

// Book struct .
type Book struct {
	BookId int		`orm:"pk;auto;unique;column(book_id)" json:"book_id"`
	// BookName 项目名称.
	BookName string		`orm:"column(book_name);size(500)" json:"book_name"`
	// Identify 项目唯一标识.
	Identify string		`orm:"column(identify);size(100);unique" json:"identify"`
	// Description 项目描述.
	Description string	`orm:"column(description);size(2000)" json:"description"`
	// PrivatelyOwned 项目私有： 0 公开/ 1 私有
	PrivatelyOwned int	`orm:"column(privately_owned);type(int);default(0)" json:"privately_owned"`
	// 当项目是私有时的访问Token.
	PrivateToken string 	`orm:"column(private_token);size(500)" json:"private_token"`
	// DocCount 包含文档数量.
	DocCount int		`orm:"column(doc_count);type(int)" json:"doc_count"`
	// CommentStatus 评论设置的状态:open为允许所有人评论，closed为不允许评论, group_only 仅允许参与者评论 ,registered_only 仅允许注册者评论.
	CommentStatus string	`orm:"column(comment_status);size(20);default(open)" json:"comment_status"`
	CommentCount int	`orm:"column(comment_count);type(int)" json:"comment_count"`

	// CreateTime 创建时间 .
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
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

