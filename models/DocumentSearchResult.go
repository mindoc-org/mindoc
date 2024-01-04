package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
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
	SearchType   string    `json:"search_type"`
}

var escape_re = regexp.MustCompile(`(?mi)(\bLIKE\s+\?)`)
var escape_replace = "${1} ESCAPE '\\'"

func need_escape(keyword string) bool {
	dbadapter, _ := web.AppConfig.String("db_adapter")
	if strings.EqualFold(dbadapter, "sqlite3") && (strings.Contains(keyword, "\\_") || strings.Contains(keyword, "\\%")) {
		return true
	}
	return false
}

func NewDocumentSearchResult() *DocumentSearchResult {
	return &DocumentSearchResult{}
}

// 分页全局搜索.
func (m *DocumentSearchResult) FindToPager(keyword string, pageIndex, pageSize, memberId int) (searchResult []*DocumentSearchResult, totalCount int, err error) {
	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize

	keyword = "%" + strings.Replace(keyword, " ", "%", -1) + "%"

	_need_escape := need_escape(keyword)
	escape_sql := func(sql string) string {
		if _need_escape {
			return escape_re.ReplaceAllString(sql, escape_replace)
		}
		return sql
	}

	if memberId <= 0 {
		sql1 := `SELECT count(doc.document_id) as total_count FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
WHERE book.privately_owned = 0 AND (doc.document_name LIKE ? OR doc.release LIKE ?) `

		sql2 := `SELECT *
FROM (
       SELECT
         doc.document_id,
         doc.modify_time,
         doc.create_time,
         doc.document_name,
         doc.identify,
         doc.release    AS description,
         book.identify  AS book_identify,
         book.book_name,
         rel.member_id,
         mdmb.account AS author,
         'document'     AS search_type
       FROM md_documents AS doc
         LEFT JOIN md_books AS book ON doc.book_id = book.book_id
         LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
         LEFT JOIN md_members AS mdmb ON rel.member_id = mdmb.member_id
       WHERE book.privately_owned = 0 AND (doc.document_name LIKE ? OR doc.release LIKE ?)
     UNION ALL
SELECT
  book.book_id AS document_id,
  book.modify_time,
  book.create_time,
  book.book_name AS document_name,
  book.identify,
  book.description,
  book.identify  AS book_identify,
  book.book_name,
  rel.member_id,
  mdmb.account AS author,
  'book'     AS search_type
FROM  md_books AS book
       LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
       LEFT JOIN md_members AS mdmb ON rel.member_id = mdmb.member_id
WHERE book.privately_owned = 0 AND (book.book_name LIKE ? OR book.description LIKE ?)

       UNION ALL
       SELECT
         blog.blog_id AS document_id,
         blog.modify_time,
         blog.create_time,
         blog.blog_title as document_name,
         blog.blog_identify,
         blog.blog_release,
         blog.blog_identify,
         blog.blog_title as book_name,
         blog.member_id,
         mdmb.account,
         'blog' AS search_type
       FROM md_blogs AS blog
         LEFT JOIN md_members AS mdmb ON blog.member_id = mdmb.member_id
       WHERE blog.blog_status = 'public' AND (blog.blog_release LIKE ? OR blog.blog_title LIKE ?)
     ) AS union_table
ORDER BY create_time DESC
LIMIT ? OFFSET ?;`

		err = o.Raw(escape_sql(sql1), keyword, keyword).QueryRow(&totalCount)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}
		sql3 := `       SELECT
         count(*)
       FROM md_blogs AS blog
       WHERE blog.blog_status = 'public' AND (blog.blog_release LIKE ? OR blog.blog_title LIKE ?);`

		c := 0
		err = o.Raw(escape_sql(sql3), keyword, keyword).QueryRow(&c)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}

		totalCount += c
		//查询项目的数量
		sql4 := `SELECT count(*) as total_count FROM md_books as book
WHERE book.privately_owned = 0 AND (book.book_name LIKE ? OR book.description LIKE ?);`

		err = o.Raw(escape_sql(sql4), keyword, keyword).QueryRow(&c)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}

		totalCount += c

		_, err = o.Raw(escape_sql(sql2), keyword, keyword, keyword, keyword, keyword, keyword, pageSize, offset).QueryRows(&searchResult)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}
	} else {
		sql1 := `SELECT count(doc.document_id) as total_count FROM md_documents AS doc
  LEFT JOIN md_books as book ON doc.book_id = book.book_id
  LEFT JOIN md_relationship AS rel ON doc.book_id = rel.book_id AND rel.role_id = 0
  LEFT JOIN md_relationship AS rel1 ON doc.book_id = rel1.book_id AND rel1.member_id = ?
			left join (select * from (select book_id,team_member_id,role_id
                   	from md_team_relationship as mtr
					left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.role_id,t.team_member_id,t.book_id) as team 
					on team.book_id = book.book_id
WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0 or team.team_member_id > 0)  AND (doc.document_name LIKE ? OR doc.release LIKE ?);`

		sql2 := `SELECT *
FROM (
       SELECT
         doc.document_id,
         doc.modify_time,
         doc.create_time,
         doc.document_name,
         doc.identify,
         doc.release    AS description,
         book.identify  AS book_identify,
         book.book_name,
         rel.member_id,
         mdmb.account AS author,
         'document'     AS search_type
       FROM md_documents AS doc
         LEFT JOIN md_books AS book ON doc.book_id = book.book_id
         LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
         LEFT JOIN md_members AS mdmb ON rel.member_id = mdmb.member_id
         LEFT JOIN md_relationship AS rel1 ON doc.book_id = rel1.book_id AND rel1.member_id = ?
         LEFT JOIN (SELECT *
                    FROM (SELECT
                            book_id,
                            team_member_id,
                            role_id
                          FROM md_team_relationship AS mtr
                            LEFT JOIN md_team_member AS mtm ON mtm.team_id = mtr.team_id AND mtm.member_id = ?
                          ORDER BY role_id DESC) AS t
                    GROUP BY t.role_id, t.team_member_id, t.book_id) AS team
           ON team.book_id = book.book_id
       WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0 OR team.team_member_id > 0) AND
             (doc.document_name LIKE ? OR doc.release LIKE ?)
       UNION ALL

       SELECT
         book.book_id AS document_id,
         book.modify_time,
         book.create_time,
         book.book_name AS document_name,
         book.identify,
         book.description AS description,
         book.identify  AS book_identify,
         book.book_name,
         rel.member_id,
         mdmb.account AS author,
         'book'     AS search_type
       FROM md_books AS book
         LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
         LEFT JOIN md_members AS mdmb ON rel.member_id = mdmb.member_id
         LEFT JOIN md_relationship AS rel1 ON book.book_id = rel1.book_id AND rel1.member_id = ?
         LEFT JOIN (SELECT *
                    FROM (SELECT
                            book_id,
                            team_member_id,
                            role_id
                          FROM md_team_relationship AS mtr
                            LEFT JOIN md_team_member AS mtm ON mtm.team_id = mtr.team_id AND mtm.member_id = ?
                          ORDER BY role_id DESC) AS t
                    GROUP BY t.role_id, t.team_member_id, t.book_id) AS team
           ON team.book_id = book.book_id
       WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0 OR team.team_member_id > 0) AND
             (book.book_name LIKE ? OR book.description LIKE ?)
 UNION ALL
       SELECT
         blog.blog_id AS document_id, 
         blog.modify_time,
         blog.create_time,
         blog.blog_title as document_name,
         blog.blog_identify as identify,
         blog.blog_release as description,
         blog.blog_identify  AS book_identify,
         blog.blog_title as book_name,
         blog.member_id,
         mdmb.account,
         'blog' AS search_type
       FROM md_blogs AS blog
         LEFT JOIN md_members AS mdmb ON blog.member_id = mdmb.member_id
       WHERE (blog.blog_status = 'public' OR blog.member_id = ?) AND blog.blog_type = 0 AND
             (blog.blog_release LIKE ? OR blog.blog_title LIKE ?)
     ) AS union_table
ORDER BY create_time DESC
LIMIT ? OFFSET ?;`

		err = o.Raw(escape_sql(sql1), memberId, memberId, keyword, keyword).QueryRow(&totalCount)
		if err != nil {
			return
		}
		sql3 := `       SELECT
         count(*)
       FROM md_blogs AS blog
       WHERE (blog.blog_status = 'public' OR blog.member_id = ?) AND blog.blog_type = 0 AND
             (blog.blog_release LIKE ? OR blog.blog_title LIKE ?);`

		c := 0
		err = o.Raw(escape_sql(sql3), memberId, keyword, keyword).QueryRow(&c)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}

		totalCount += c

		sql4 := `SELECT count(*) as total_count FROM md_books as book
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.role_id = 0
  LEFT JOIN md_relationship AS rel1 ON book.book_id = rel1.book_id AND rel1.member_id = ?
			left join (select * from (select book_id,team_member_id,role_id
                   	from md_team_relationship as mtr
					left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.role_id,t.team_member_id,t.book_id) as team
					on team.book_id = book.book_id
WHERE (book.privately_owned = 0 OR rel1.relationship_id > 0 or team.team_member_id > 0)  AND (book.book_name LIKE ? OR book.description LIKE ?);`

		err = o.Raw(escape_sql(sql4), memberId, memberId, keyword, keyword).QueryRow(&c)
		if err != nil {
			logs.Error("查询搜索结果失败 -> ", err)
			return
		}

		totalCount += c

		_, err = o.Raw(escape_sql(sql2), memberId, memberId, keyword, keyword, memberId, memberId, keyword, keyword, memberId, keyword, keyword, pageSize, offset).QueryRows(&searchResult)
		if err != nil {
			return
		}
	}
	return
}

// 项目内搜索.
func (m *DocumentSearchResult) SearchDocument(keyword string, bookId int) (docs []*DocumentSearchResult, err error) {
	o := orm.NewOrm()

	sql := "SELECT * FROM md_documents WHERE book_id = ? AND (document_name LIKE ? OR `release` LIKE ?) "
	keyword = "%" + keyword + "%"

	_need_escape := need_escape(keyword)
	escape_sql := func(sql string) string {
		if _need_escape {
			return escape_re.ReplaceAllString(sql, escape_replace)
		}
		return sql
	}

	_, err = o.Raw(escape_sql(sql), bookId, keyword, keyword).QueryRows(&docs)

	return
}

// 所有项目搜索.
func (m *DocumentSearchResult) SearchAllDocument(keyword string) (docs []*DocumentSearchResult, err error) {
	o := orm.NewOrm()

	sql := "SELECT * FROM md_documents WHERE (document_name LIKE ? OR `release` LIKE ?) "
	keyword = "%" + keyword + "%"

	_need_escape := need_escape(keyword)
	escape_sql := func(sql string) string {
		if _need_escape {
			return escape_re.ReplaceAllString(sql, escape_replace)
		}
		return sql
	}

	_, err = o.Raw(escape_sql(sql), keyword, keyword).QueryRows(&docs)

	return
}
