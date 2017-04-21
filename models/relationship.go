package models

type Relationship struct {
	RelationshipId int	`orm:"pk;auto;unique;column(relationship_id)" json:"relationship_id"`
	MemberId int		`orm:"column(member_id);type(int)" json:"member_id"`
	BookId int		`orm:"column(book_id);type(int)" json:"book_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId int		`orm:"column(role_id);type(int)" json:"role_id"`
}


// TableName 获取对应数据库表名.
func (m *Relationship) TableName() string {
	return "relationship"
}
// TableEngine 获取数据使用的引擎.
func (m *Relationship) TableEngine() string {
	return "INNODB"
}

// 联合唯一键
func (u *Relationship) TableUnique() [][]string {
	return [][]string{
		[]string{"MemberId", "BookId"},
	}
}