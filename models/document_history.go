package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
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
}

type DocumentHistorySimpleResult struct {
	HistoryId    int                `json:"history_id"`
	ActionName   string             `json:"action_name"`
	MemberId     int 		`json:"member_id"`
	Account      string 		`json:"account"`
	ModifyAt     int 		`json:"modify_at"`
	ModifyName   string 		`json:"modify_name"`
	ModifyTime   time.Time 		`json:"modify_time"`
	Version      int64 		`json:"version"`
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
func (m *DocumentHistory) Find(id int) (*DocumentHistory,error) {
	o := orm.NewOrm()
	err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id",id).One(m)

	return m,err
}
//清空指定文档的历史.
func (m *DocumentHistory) Clear(doc_id int) error {
	o := orm.NewOrm()

	_, err := o.Raw("DELETE md_document_history WHERE document_id = ?", doc_id).Exec()

	return err
}

//删除历史.
func (m *DocumentHistory) Delete(history_id,doc_id int) error {
	o := orm.NewOrm()

	_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id",history_id).Filter("document_id",doc_id).Delete()

	return err
}

//恢复指定历史的文档.
func (m *DocumentHistory) Restore(history_id,doc_id,uid int) error {
	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("history_id", history_id).Filter("document_id",doc_id).One(m)

	if err != nil {
		return err
	}
	doc, err := NewDocument().Find(m.DocumentId)

	if err != nil {
		return err
	}
	history := NewDocumentHistory()
	history.DocumentId = doc_id
	history.Content = doc.Content
	history.Markdown = doc.Markdown
	history.DocumentName = doc.DocumentName
	history.ModifyAt = uid
	history.MemberId = doc.MemberId
	history.ParentId = doc.ParentId
	history.Version = time.Now().Unix()
	history.Action = "restore"
	history.ActionName = "恢复文档"

	history.InsertOrUpdate()

	doc.DocumentName = m.DocumentName
	doc.Content = m.Content
	doc.Markdown = m.Markdown
	doc.Release = m.Content
	doc.Version = time.Now().Unix()

	_, err = o.Update(doc)

	return err
}

func (m *DocumentHistory) InsertOrUpdate() (history *DocumentHistory,err error)  {
	o := orm.NewOrm()
	history = m

	if m.HistoryId > 0 {
		_,err = o.Update(m)
	}else{
		_,err = o.Insert(m)
	}
	return
}
//分页查询指定文档的历史.
func (m *DocumentHistory) FindToPager(doc_id, page_index, page_size int) (docs []*DocumentHistorySimpleResult, totalCount int, err error) {

	o := orm.NewOrm()

	offset := (page_index - 1) * page_size

	totalCount = 0

	sql := `SELECT history.*,m1.account,m2.account as modify_name
FROM md_document_history AS history
LEFT JOIN md_members AS m1 ON history.member_id = m1.member_id
LEFT JOIN md_members AS m2 ON history.modify_at = m2.member_id
WHERE history.document_id = ? ORDER BY history.history_id DESC LIMIT ?,?;`

	_, err = o.Raw(sql,doc_id,offset,page_size).QueryRows(&docs)

	if err != nil {
		return
	}
	var count int64
	count, err = o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", doc_id).Count()

	if err != nil {
		return
	}
	totalCount = int(count)

	return
}
