package models

import (
	"time"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"strings"
	"os"
	"path/filepath"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/lifei6671/mindoc/cache"
	"encoding/json"
)

// Document struct.
type Document struct {
	DocumentId   int    `orm:"pk;auto;unique;column(document_id)" json:"doc_id"`
	DocumentName string `orm:"column(document_name);size(500)" json:"doc_name"`
	// Identify 文档唯一标识
	Identify  string `orm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId    int    `orm:"column(book_id);type(int);index" json:"book_id"`
	ParentId  int    `orm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	OrderSort int    `orm:"column(order_sort);default(0);type(int);index" json:"order_sort"`
	// Markdown markdown格式文档.
	Markdown string `orm:"column(markdown);type(text);null" json:"markdown"`
	// Release 发布后的Html格式内容.
	Release string `orm:"column(release);type(text);null" json:"release"`
	// Content 未发布的 Html 格式内容.
	Content    string        `orm:"column(content);type(text);null" json:"content"`
	CreateTime time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId   int           `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime time.Time     `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt   int           `orm:"column(modify_at);type(int)" json:"-"`
	Version    int64         `orm:"type(bigint);column(version)" json:"version"`
	AttachList []*Attachment `orm:"-" json:"attach"`
}

// TableName 获取对应数据库表名.
func (m *Document) TableName() string {
	return "documents"
}

// TableEngine 获取数据使用的引擎.
func (m *Document) TableEngine() string {
	return "INNODB"
}

func (m *Document) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}

//根据文档ID查询指定文档.
func (m *Document) Find(id int) (*Document, error) {
	if id <= 0 {
		return m, ErrInvalidParameter
	}

	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", id).One(m)

	if err == orm.ErrNoRows {
		return m, ErrDataNotExist
	}

	return m, nil
}

//插入和更新文档.
func (m *Document) InsertOrUpdate(cols ...string) error {
	o := orm.NewOrm()
	var err error
	if m.DocumentId > 0 {
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
		NewBook().ResetDocumentNumber(m.BookId)
	}
	if err != nil {
		return err
	}

	return nil
}

//根据指定字段查询一条文档.
func (m *Document) FindByFieldFirst(field string, v interface{}) (*Document, error) {

	o := orm.NewOrm()

	err := o.QueryTable(m.TableNameWithPrefix()).Filter(field, v).One(m)

	return m, err
}

//递归删除一个文档.
func (m *Document) RecursiveDocument(docId int) error {

	o := orm.NewOrm()

	if doc, err := m.Find(docId); err == nil {
		o.Delete(doc)
		NewDocumentHistory().Clear(doc.DocumentId)
	}
	//
	//var docs []*Document
	//
	//_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("parent_id", doc_id).All(&docs)

	var maps []orm.Params

	_, err := o.Raw("SELECT document_id FROM " + m.TableNameWithPrefix() + " WHERE parent_id=" + strconv.Itoa(docId)).Values(&maps)
	if err != nil {
		beego.Error("RecursiveDocument => ", err)
		return err
	}

	for _, item := range maps {
		if docId,ok := item["document_id"].(string); ok{
			id,_ := strconv.Atoi(docId)
			o.QueryTable(m.TableNameWithPrefix()).Filter("document_id", id).Delete()
			m.RecursiveDocument(id)
		}
	}

	return nil
}

//发布文档
func (m *Document) ReleaseContent(bookId int) {

	o := orm.NewOrm()

	var docs []*Document
	_, err := o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", bookId).All(&docs, "document_id","identify", "content")

	if err != nil {
		beego.Error("发布失败 => ", err)
		return
	}
	for _, item := range docs {
		if item.Content != "" {
			item.Release = item.Content
			bufio := bytes.NewReader([]byte(item.Content))
			//解析文档中非本站的链接，并设置为新窗口打开
			if content, err := goquery.NewDocumentFromReader(bufio);err == nil {

				content.Find("a").Each(func(i int, contentSelection *goquery.Selection) {
					if src, ok := contentSelection.Attr("href"); ok{
						if strings.HasPrefix(src, "http://") || strings.HasPrefix(src,"https://") {
							//beego.Info(src,conf.BaseUrl,strings.HasPrefix(src,conf.BaseUrl))
							if conf.BaseUrl != "" && !strings.HasPrefix(src,conf.BaseUrl) {
								contentSelection.SetAttr("target", "_blank")
								if html, err := content.Html();err == nil {
									item.Release = html
								}
							}
						}

					}
				})
			}
		}

		attachList, err := NewAttachment().FindListByDocumentId(item.DocumentId)
		if err == nil && len(attachList) > 0 {
			content := bytes.NewBufferString("<div class=\"attach-list\"><strong>附件</strong><ul>")
			for _, attach := range attachList {
				if strings.HasPrefix(attach.HttpPath, "/") {
					attach.HttpPath = strings.TrimSuffix(conf.BaseUrl, "/") + attach.HttpPath
				}
				li := fmt.Sprintf("<li><a href=\"%s\" target=\"_blank\" title=\"%s\">%s</a></li>", attach.HttpPath, attach.FileName, attach.FileName)

				content.WriteString(li)
			}
			content.WriteString("</ul></div>")
			item.Release += content.String()
		}
		_, err = o.Update(item, "release")
		if err != nil {
			beego.Error(fmt.Sprintf("发布失败 => %+v", item), err)
		}else {
			//当文档发布后，需要清除已缓存的转换文档和文档缓存
			if doc,err := NewDocument().Find(item.DocumentId); err == nil {
				doc.PutToCache()
			}else{
				doc.RemoveCache()
			}

			os.RemoveAll(filepath.Join(conf.WorkingDirectory,"uploads","books",strconv.Itoa(bookId)))
		}
	}
}

//将文档写入缓存
func (m *Document) PutToCache(){
	go func(m Document) {
		if v,err := json.Marshal(&m);err == nil {
			if m.Identify == "" {

				if err := cache.Put("Document.Id." + strconv.Itoa(m.DocumentId), v, time.Second*3600); err != nil {
					beego.Info("文档缓存失败:", m.DocumentId)
				}
			}else{
				if err := cache.Put("Document.Identify."+ m.Identify, v, time.Second*3600); err != nil {
					beego.Info("文档缓存失败:", m.DocumentId)
				}
			}
		}
	}(*m)
}
//清除缓存
func (m *Document) RemoveCache() {
	go func(m Document) {
		cache.Put("Document.Id." + strconv.Itoa(m.DocumentId), m, time.Second*3600);

		if m.Identify != "" {
			cache.Put("Document.Identify."+ m.Identify, m, time.Second*3600);
		}
	}(*m)
}

//从缓存获取
func (m *Document) FromCacheById(id int) (*Document,error) {
	b := cache.Get("Document.Id." + strconv.Itoa(id))
	if v,ok := b.([]byte); ok {

		if err := json.Unmarshal(v,m);err == nil{
			beego.Info("从缓存中获取文档信息成功",m.DocumentId)
			return m,nil
		}
	}
	defer func() {
		if m.DocumentId > 0 {
			m.PutToCache()
		}
	}()
	return m.Find(id)
}
//根据文档标识从缓存中查询文档
func (m *Document) FromCacheByIdentify(identify string) (*Document,error) {
	b := cache.Get("Document.Identify." + identify)
	if v,ok := b.([]byte); ok {
		if err := json.Unmarshal(v,m);err == nil{
			beego.Info("从缓存中获取文档信息成功",m.DocumentId,identify)
			return m,nil
		}
	}
	defer func() {
		if m.DocumentId > 0 {
			m.PutToCache()
		}
	}()
	return m.FindByFieldFirst("identify",identify)
}

//根据项目ID查询文档列表.
func (m *Document) FindListByBookId(bookId int) (docs []*Document, err error) {
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableNameWithPrefix()).Filter("book_id", bookId).OrderBy("order_sort").All(&docs)

	return
}
