package models

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/mindoc-org/mindoc/conf"
)

type DocumentHistory struct {
	HistoryId    int       `orm:"column(history_id);pk;auto;unique" json:"history_id"`
	Action       string    `orm:"column(action);size(255)" json:"action"`
	ActionName   string    `orm:"column(action_name);size(255)" json:"action_name"`
	DocumentId   int       `orm:"column(document_id);type(int);index" json:"doc_id"`
	DocumentName string    `orm:"column(document_name);size(500)" json:"doc_name"`
	ParentId     int       `orm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	Markdown     string    `orm:"column(markdown);type(text);null" json:"markdown"`
	Content      string    `orm:"column(content);type(text);null" json:"content"`
	MemberId     int       `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt     int       `orm:"column(modify_at);type(int)" json:"-"`
	Version      int64     `orm:"type(bigint);column(version)" json:"version"`
	IsOpen       int       `orm:"column(is_open);type(int);default(0)" json:"is_open"`
}

type DocumentHistorySimpleResult struct {
	HistoryId  int       `json:"history_id"`
	ActionName string    `json:"action_name"`
	MemberId   int       `json:"member_id"`
	Account    string    `json:"account"`
	ModifyAt   int       `json:"modify_at"`
	ModifyName string    `json:"modify_name"`
	ModifyTime time.Time `json:"modify_time"`
	Version    int64     `json:"version"`
}

// TableName 获取对应数据库表名.
func (m *DocumentHistory) TableName() string {
	return "document_history"
}

// TableEngine 获取数据使用的引擎.
func (m *DocumentHistory) TableEngine() string {
	return "INNODB"
}

func (m *DocumentHistory) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewDocumentHistory() *DocumentHistory {
	return &DocumentHistory{}
}
func (m *DocumentHistory) Find(id int) (*DocumentHistory, error) {
	o := orm.NewOrm()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id", id).One(m)

	return m, err
}

//清空指定文档的历史.
func (m *DocumentHistory) Clear(docId int) error {
	o := orm.NewOrm()

	_, err := o.Raw("DELETE md_document_history WHERE document_id = ?", docId).Exec()

	return err
}

//删除历史.
func (m *DocumentHistory) Delete(historyId, docId int) error {
	o := orm.NewOrm()

	_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id", historyId).Filter("document_id", docId).Delete()

	return err
}

//恢复指定历史的文档.
func (m *DocumentHistory) Restore(historyId, docId, uid int) error {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id", historyId).Filter("document_id", docId).One(m)

	if err != nil {
		return err
	}
	doc, err := NewDocument().Find(m.DocumentId)

	if err != nil {
		return err
	}
	history := NewDocumentHistory()
	history.DocumentId = docId
	history.Content = doc.Content
	history.Markdown = doc.Markdown
	history.DocumentName = doc.DocumentName
	history.ModifyAt = uid
	history.MemberId = doc.MemberId
	history.ParentId = doc.ParentId
	history.Version = time.Now().Unix()
	history.Action = "restore"
	history.ActionName = "恢复文档"
	history.IsOpen = doc.IsOpen

	history.InsertOrUpdate()

	doc.DocumentName = m.DocumentName
	doc.Content = m.Content
	doc.Markdown = m.Markdown
	doc.Release = m.Content
	doc.Version = time.Now().Unix()
	doc.IsOpen = m.IsOpen

	_, err = o.Update(doc)

	return err
}

func (m *DocumentHistory) InsertOrUpdate() (history *DocumentHistory, err error) {
	o := orm.NewOrm()
	history = m

	if m.HistoryId > 0 {
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
		if err == nil {
			if doc, e := NewDocument().Find(m.DocumentId); e == nil {
				if book, e := NewBook().Find(doc.BookId); e == nil && book.HistoryCount > 0 {
					//如果已存在的历史记录大于指定的记录，则清除旧记录
					if c, e := o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc.DocumentId).Count(); e == nil && c > int64(book.HistoryCount) {

						count := c - int64(book.HistoryCount)
						logs.Info("需要删除的历史文档数量：", count)
						var lists []DocumentHistory

						if _, e := o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc.DocumentId).OrderBy("history_id").Limit(count).All(&lists, "history_id"); e == nil {
							for _, d := range lists {
								o.Delete(&d)
							}
						}
					} else {
						logs.Info(book.HistoryCount)
					}
				}
			}

		}
	}
	return
}

//分页查询指定文档的历史.
func (m *DocumentHistory) FindToPager(docId, pageIndex, pageSize int) (docs []*DocumentHistorySimpleResult, totalCount int, err error) {

	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize

	totalCount = 0

	sql := `SELECT history.*,m1.account,m2.account as modify_name
FROM md_document_history AS history
LEFT JOIN md_members AS m1 ON history.member_id = m1.member_id
LEFT JOIN md_members AS m2 ON history.modify_at = m2.member_id
WHERE history.document_id = ? ORDER BY history.history_id DESC LIMIT ?,?;`

	_, err = o.Raw(sql, docId, offset, pageSize).QueryRows(&docs)

	if err != nil {
		return
	}
	var count int64
	count, err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", docId).Count()

	if err != nil {
		return
	}
	totalCount = int(count)

	return
}
