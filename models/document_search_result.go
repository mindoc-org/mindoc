package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type DocumentSearchResult struct {
	DocumentId   int    `json:"doc_id"`
	DocumentName string `json:"doc_name"`
	// Identify 文档唯一标识
	Identify     string    `json:"identify"`
	Description  string    `json:"description"`
	Author       string    `json:"author"`
	ModifyTime   time.Time `json:"modify_time"`
	CreateTime   time.Time `json:"create_time"`
	BookId       int       `json:"book_id"`
	BookName     string    `json:"book_name"`
	BookIdentify string    `json:"book_identify"`
}

func NewDocumentSearchResult() *DocumentSearchResult {
	return &DocumentSearchResult{}
}

//分页全局搜索.
func (m *DocumentSearchResult) FindToPager(keyword string, page_index, page_size, member_id int) (search_result []*DocumentSearchResult, total_count int, err error) {
	o := orm.NewOrm()

	offset := (page_index - 1) * page_size
	keyword = "%" + keyword + "%"

	if member_id <= 0 {
		sql1 := `SELECT count(doc.document_id) as total_count FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
WHERE book.privately_owned = 0 AND (doc.document_name LIKE ? OR doc.release LIKE ?) `

		sql2 := `SELECT doc.document_id,doc.modify_time,doc.create_time,doc.document_name,doc.identify,doc.release as description,doc.modify_time,book.identify as book_identify,book.book_name,rel.member_id,member.account AS author FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
  LEFT JOIN md_members as member ON rel.member_id = member.member_id
WHERE book.privately_owned = 0 AND (doc.document_name LIKE ? OR doc.release LIKE ?)
 ORDER BY doc.document_id DESC LIMIT ?,? `

		err = o.Raw(sql1, keyword, keyword).QueryRow(&total_count)
		if err != nil {
			return
		}
		_, err = o.Raw(sql2, keyword, keyword, offset, page_size).QueryRows(&search_result)
		if err != nil {
			return
		}
	} else {
		sql1 := `SELECT count(doc.document_id) as total_count FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
  LEFT JOIN md_relationship AS rel ON doc.book_id = rel.book_id AND rel.role_id = 0
  LEFT JOIN md_relationship AS rel1 ON doc.book_id = rel1.book_id AND rel1.member_id = ?
WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0)  AND (doc.document_name LIKE ? OR doc.release LIKE ?) `

		sql2 := `SELECT doc.document_id,doc.modify_time,doc.create_time,doc.document_name,doc.identify,doc.release as description,doc.modify_time,book.identify as book_identify,book.book_name,rel.member_id,member.account AS author FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
  LEFT JOIN md_members as member ON rel.member_id = member.member_id
  LEFT JOIN md_relationship AS rel1 ON doc.book_id = rel1.book_id AND rel1.member_id = ?
WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0)  AND (doc.document_name LIKE ? OR doc.release LIKE ?)
 ORDER BY doc.document_id DESC LIMIT ?,? `

		err = o.Raw(sql1, member_id, keyword, keyword).QueryRow(&total_count)
		if err != nil {
			return
		}
		_, err = o.Raw(sql2, member_id, keyword, keyword, offset, page_size).QueryRows(&search_result)
		if err != nil {
			return
		}
	}
	return
}

//项目内搜索.
func (m *DocumentSearchResult) SearchDocument(keyword string, book_id int) (docs []*DocumentSearchResult, err error) {
	o := orm.NewOrm()

	sql := "SELECT * FROM md_documents WHERE book_id = ? AND (document_name LIKE ? OR `release` LIKE ?) "
	keyword = "%" + keyword + "%"

	_, err = o.Raw(sql, book_id, keyword, keyword).QueryRows(&docs)

	return
}
