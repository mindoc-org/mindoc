package models

import (
	"time"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/mindoc-org/mindoc/conf"
)

type CommentVote struct {
	VoteId          int       `orm:"column(vote_id);pk;auto;unique" json:"vote_id"`
	CommentId       int       `orm:"column(comment_id);type(int);index" json:"comment_id"`
	CommentMemberId int       `orm:"column(comment_member_id);type(int);index;default(0)" json:"comment_member_id"`
	VoteMemberId    int       `orm:"column(vote_member_id);type(int);index" json:"vote_member_id"`
	VoteState       int       `orm:"column(vote_state);type(int)" json:"vote_state"`
	CreateTime      time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
}

// TableName 获取对应数据库表名.
func (m *CommentVote) TableName() string {
	return "comment_votes"
}

// TableEngine 获取数据使用的引擎.
func (m *CommentVote) TableEngine() string {
	return "INNODB"
}

func (m *CommentVote) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}
func (u *CommentVote) TableUnique() [][]string {
	return [][]string{
		[]string{"comment_id", "vote_member_id"},
	}
}
func NewCommentVote() *CommentVote {
	return &CommentVote{}
}
func (m *CommentVote) InsertOrUpdate() (*CommentVote, error) {
	o := orm.NewOrm()

	if m.VoteId > 0 {
		_, err := o.Update(m)
		return m, err
	} else {
		_, err := o.Insert(m)

		return m, err
	}
}
