package models

import "github.com/beego/beego/v2/adapter/orm"

type CommentResult struct {
	Comment
	Author       string `json:"author"`
	ReplyAccount string `json:"reply_account"`
}

func (m *CommentResult) FindForDocumentToPager(doc_id, page_index, page_size int) (comments []*CommentResult, totalCount int, err error) {

	o := orm.NewOrm()

	sql1 := `
SELECT
  comment.* ,
  parent.* ,
  member.account AS author,
  p_member.account AS reply_account
FROM md_comments AS comment
  LEFT JOIN md_members AS member ON comment.member_id = member.member_id
  LEFT JOIN md_comments AS parent ON comment.parent_id = parent.comment_id
  LEFT JOIN md_members AS p_member ON p_member.member_id = parent.member_id

WHERE comment.document_id = ? ORDER BY comment.comment_id DESC LIMIT 0,10`

	offset := (page_index - 1) * page_size

	_, err = o.Raw(sql1, doc_id, offset, page_size).QueryRows(&comments)

	v, err := o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc_id).Count()

	if err == nil {
		totalCount = int(v)
	}

	return
}
