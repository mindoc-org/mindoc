package models

import (
	"time"

	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/cache"
	"github.com/lifei6671/mindoc/conf"
	"os"
	"path/filepath"
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"strings"
	"github.com/lifei6671/mindoc/utils"
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
	Content    string    `orm:"column(content);type(text);null" json:"content"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId   int       `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt   int       `orm:"column(modify_at);type(int)" json:"-"`
	Version    int64     `orm:"column(version);type(bigint);" json:"version"`
	//是否展开子目录：0 否/1 是
	IsOpen     int           `orm:"column(is_open);type(int);default(0)" json:"is_open"`
	AttachList []*Attachment `orm:"-" json:"attach"`
}

// 多字段唯一键
func (item *Document) TableUnique() [][]string {
	return [][]string{
		[]string{"book_id", "identify"},
	}
}

// TableName 获取对应数据库表名.
func (item *Document) TableName() string {
	return "documents"
}

// TableEngine 获取数据使用的引擎.
func (item *Document) TableEngine() string {
	return "INNODB"
}

func (item *Document) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + item.TableName()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}

//根据文档ID查询指定文档.
func (item *Document) Find(id int) (*Document, error) {
	if id <= 0 {
		return item, ErrInvalidParameter
	}

	o := orm.NewOrm()

	err := o.QueryTable(item.TableNameWithPrefix()).Filter("document_id", id).One(item)

	if err == orm.ErrNoRows {
		return item, ErrDataNotExist
	}

	return item, nil
}

//插入和更新文档.
func (item *Document) InsertOrUpdate(cols ...string) error {
	o := orm.NewOrm()
	var err error
	if item.DocumentId > 0 {
		_, err = o.Update(item, cols...)
	} else {
		if item.Identify == "" {
			book := NewBook()
			identify := "docs"
			if err := o.QueryTable(book.TableNameWithPrefix()).Filter("book_id", item.BookId).One(book, "identify"); err == nil {
				identify = book.Identify
			}

			item.Identify = fmt.Sprintf("%s-%s", identify, strconv.FormatInt(time.Now().UnixNano(), 32))
		}

		if item.OrderSort == 0 {
			sort, _ := o.QueryTable(item.TableNameWithPrefix()).Filter("book_id", item.BookId).Filter("parent_id", item.ParentId).Count()
			item.OrderSort = int(sort) + 1
		}
		_, err = o.Insert(item)
		NewBook().ResetDocumentNumber(item.BookId)
	}
	if err != nil {
		return err
	}

	return nil
}

//根据文档识别编号和项目id获取一篇文档
func (item *Document) FindByIdentityFirst(identify string, bookId int) (*Document, error) {
	o := orm.NewOrm()

	err := o.QueryTable(item.TableNameWithPrefix()).Filter("book_id", bookId).Filter("identify", identify).One(item)

	return item, err
}

//递归删除一个文档.
func (item *Document) RecursiveDocument(docId int) error {

	o := orm.NewOrm()

	if doc, err := item.Find(docId); err == nil {
		o.Delete(doc)
		NewDocumentHistory().Clear(doc.DocumentId)
	}
	var maps []orm.Params

	_, err := o.Raw("SELECT document_id FROM " + item.TableNameWithPrefix() + " WHERE parent_id=" + strconv.Itoa(docId)).Values(&maps)
	if err != nil {
		beego.Error("RecursiveDocument => ", err)
		return err
	}

	for _, param := range maps {
		if docId, ok := param["document_id"].(string); ok {
			id, _ := strconv.Atoi(docId)
			o.QueryTable(item.TableNameWithPrefix()).Filter("document_id", id).Delete()
			item.RecursiveDocument(id)
		}
	}

	return nil
}

//将文档写入缓存
func (item *Document) PutToCache() {
	go func(m Document) {

		if m.Identify == "" {

			if err := cache.Put("Document.Id."+strconv.Itoa(m.DocumentId), m, time.Second*3600); err != nil {
				beego.Info("文档缓存失败:", m.DocumentId)
			}
		} else {
			if err := cache.Put(fmt.Sprintf("Document.BookId.%d.Identify.%s", m.BookId, m.Identify), m, time.Second*3600); err != nil {
				beego.Info("文档缓存失败:", m.DocumentId)
			}
		}

	}(*item)
}

//清除缓存
func (item *Document) RemoveCache() {
	go func(m Document) {
		cache.Put("Document.Id."+strconv.Itoa(m.DocumentId), m, time.Second*3600)

		if m.Identify != "" {
			cache.Put(fmt.Sprintf("Document.BookId.%d.Identify.%s", m.BookId, m.Identify), m, time.Second*3600)
		}
	}(*item)
}

//从缓存获取
func (item *Document) FromCacheById(id int) (*Document, error) {

	if err := cache.Get("Document.Id."+strconv.Itoa(id), &item); err == nil && item.DocumentId > 0 {
		beego.Info("从缓存中获取文档信息成功 ->", item.DocumentId)
		return item, nil
	}

	if item.DocumentId > 0 {
		item.PutToCache()
	}
	item, err := item.Find(id)

	if err == nil {
		item.PutToCache()
	}
	return item, err
}

//根据文档标识从缓存中查询文档
func (item *Document) FromCacheByIdentify(identify string, bookId int) (*Document, error) {

	key := fmt.Sprintf("Document.BookId.%d.Identify.%s", bookId, identify)

	if err := cache.Get(key, item); err == nil && item.DocumentId > 0 {
		beego.Info("从缓存中获取文档信息成功 ->", key)
		return item, nil
	}

	defer func() {
		if item.DocumentId > 0 {
			item.PutToCache()
		}
	}()
	return item.FindByIdentityFirst(identify, bookId)
}

//根据项目ID查询文档列表.
func (item *Document) FindListByBookId(bookId int) (docs []*Document, err error) {
	o := orm.NewOrm()

	_, err = o.QueryTable(item.TableNameWithPrefix()).Filter("book_id", bookId).OrderBy("order_sort").All(&docs)

	return
}

//判断文章是否存在
func (item *Document) IsExist(documentId int) bool {
	o := orm.NewOrm()

	return o.QueryTable(item.TableNameWithPrefix()).Filter("document_id", documentId).Exist()
}

//发布单篇文档
func (item *Document) ReleaseContent() error {

	item.Release = strings.TrimSpace(item.Content)

	err := item.Processor().InsertOrUpdate("release")

	if err != nil {
		beego.Error(fmt.Sprintf("发布失败 -> %+v", item), err)
		return err
	}
	//当文档发布后，需要清除已缓存的转换文档和文档缓存
	item.RemoveCache()

	if err := os.RemoveAll(filepath.Join(conf.WorkingDirectory, "uploads", "books", strconv.Itoa(item.BookId))); err != nil {
		beego.Error("删除已缓存的文档目录失败 -> ", filepath.Join(conf.WorkingDirectory, "uploads", "books", strconv.Itoa(item.BookId)))
		return err
	}

	return nil
}

//处理文档的外链，附件，底部编辑信息等.
func (item *Document) Processor() *Document {
	if item.Release != "" {
		item.Release = utils.SafetyProcessor(item.Release)

		//安全过滤，移除危险标签和属性
		if docQuery, err := goquery.NewDocumentFromReader(bytes.NewBufferString(item.Release)); err == nil {

			//处理附件
			if selector := docQuery.Find("div.attach-list").First(); selector.Size() <= 0 {
				//处理附件
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
					if docQuery == nil {
						docQuery, err = goquery.NewDocumentFromReader(content);
					} else {
						if selector := docQuery.Find("div.wiki-bottom").First(); selector.Size() > 0 {
							selector.BeforeHtml(content.String())
						} else if selector := docQuery.Find("div.markdown-article").First(); selector.Size() > 0 {
							selector.AppendHtml(content.String())
						} else if selector := docQuery.Find("article.markdown-article-inner").First(); selector.Size() > 0 {
							selector.AppendHtml(content.String())
						} else {
							docQuery.Children().WrapHtml("<article class=\"markdown-article-inner\"></article>").AppendHtml(content.String())
						}
					}
				}
			}

			//处理了文档底部信息
			if selector := docQuery.Find("div.wiki-bottom").First(); selector.Size() <= 0 && item.MemberId > 0 {
				//处理文档结尾信息
				docCreator, err := NewMember().Find(item.MemberId, "real_name", "account")
				release := "<div class=\"wiki-bottom\">文档更新时间: " + item.ModifyTime.Local().Format("2006-01-02 15:04") + " &nbsp;&nbsp;作者："
				if err == nil && docCreator != nil {
					if docCreator.RealName != "" {
						release += docCreator.RealName
					} else {
						release += docCreator.Account
					}
				}
				release += "</div>"

				if selector := docQuery.Find("div.markdown-article").First(); selector.Size() > 0 {
					selector.AppendHtml(release)
				} else if selector := docQuery.Find("article.markdown-article-inner").First(); selector.Size() > 0 {
					selector.First().AppendHtml(release)
				} else {
					docQuery.Children().WrapHtml("<article class=\"markdown-article-inner\"></article>").AppendHtml(release);
				}
			}

			//设置图片为CDN地址
			if cdnimg := beego.AppConfig.String("cdnimg"); cdnimg != "" {
				docQuery.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
					if src, ok := contentSelection.Attr("src"); ok && strings.HasPrefix(src, "/uploads/") {
						contentSelection.SetAttr("src", utils.JoinURI(cdnimg, src))
					}
				})
			}

			if html, err := docQuery.Html(); err == nil {
				item.Release = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(html), "<html><head></head><body>"), "</body></html>")
			}
		}
	}
	return item
}
