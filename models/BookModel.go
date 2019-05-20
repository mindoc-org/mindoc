package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/utils/cryptil"
	"github.com/lifei6671/mindoc/utils/filetil"
	"github.com/lifei6671/mindoc/utils/requests"
	"github.com/lifei6671/mindoc/utils/ziptil"
	"gopkg.in/russross/blackfriday.v2"
	"encoding/json"
	"github.com/lifei6671/mindoc/utils"
)

// Book struct .
type Book struct {
	BookId int `orm:"pk;auto;unique;column(book_id)" json:"book_id"`
	// BookName 项目名称.
	BookName string `orm:"column(book_name);size(500)" json:"book_name"`
	//所属项目空间
	ItemId   int    `orm:"column(item_id);type(int);default(1)" json:"item_id"`
	// Identify 项目唯一标识.
	Identify string `orm:"column(identify);size(100);unique" json:"identify"`
	//是否是自动发布 0 否/1 是
	AutoRelease int `orm:"column(auto_release);type(int);default(0)" json:"auto_release"`
	//是否开启下载功能 0 是/1 否
	IsDownload int `orm:"column(is_download);type(int);default(0)" json:"is_download"`
	OrderIndex int `orm:"column(order_index);type(int);default(0)" json:"order_index"`
	// Description 项目描述.
	Description string `orm:"column(description);size(2000)" json:"description"`
	//发行公司
	Publisher string `orm:"column(publisher);size(500)" json:"publisher"`
	Label     string `orm:"column(label);size(500)" json:"label"`
	// PrivatelyOwned 项目私有： 0 公开/ 1 私有
	PrivatelyOwned int `orm:"column(privately_owned);type(int);default(0)" json:"privately_owned"`
	// 当项目是私有时的访问Token.
	PrivateToken string `orm:"column(private_token);size(500);null" json:"private_token"`
	//访问密码.
	BookPassword string `orm:"column(book_password);size(500);null" json:"book_password"`
	//状态：0 正常/1 已删除
	Status int `orm:"column(status);type(int);default(0)" json:"status"`
	//默认的编辑器.
	Editor string `orm:"column(editor);size(50)" json:"editor"`
	// DocCount 包含文档数量.
	DocCount int `orm:"column(doc_count);type(int)" json:"doc_count"`
	// CommentStatus 评论设置的状态:open 为允许所有人评论，closed 为不允许评论, group_only 仅允许参与者评论 ,registered_only 仅允许注册者评论.
	CommentStatus string `orm:"column(comment_status);size(20);default(open)" json:"comment_status"`
	CommentCount  int    `orm:"column(comment_count);type(int)" json:"comment_count"`
	//封面地址
	Cover string `orm:"column(cover);size(1000)" json:"cover"`
	//主题风格
	Theme string `orm:"column(theme);size(255);default(default)" json:"theme"`
	// CreateTime 创建时间 .
	CreateTime time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	//每个文档保存的历史记录数量，0 为不限制
	HistoryCount int `orm:"column(history_count);type(int);default(0)" json:"history_count"`
	//是否启用分享，0启用/1不启用
	IsEnableShare int       `orm:"column(is_enable_share);type(int);default(0)" json:"is_enable_share"`
	MemberId      int       `orm:"column(member_id);size(100)" json:"member_id"`
	ModifyTime    time.Time `orm:"type(datetime);column(modify_time);null;auto_now" json:"modify_time"`
	Version       int64     `orm:"type(bigint);column(version)" json:"version"`
	//是否使用第一篇文章项目为默认首页,0 否/1 是
	IsUseFirstDocument int `orm:"column(is_use_first_document);type(int);default(0)" json:"is_use_first_document"`
	//是否开启自动保存：0 否/1 是
	AutoSave int `orm:"column(auto_save);type(tinyint);default(0)" json:"auto_save"`
}

func (book *Book) String() string {
	ret, err := json.Marshal(*book)

	if err != nil {
		return ""
	}
	return string(ret)
}

// TableName 获取对应数据库表名.
func (book *Book) TableName() string {
	return "books"
}

// TableEngine 获取数据使用的引擎.
func (book *Book) TableEngine() string {
	return "INNODB"
}
func (book *Book) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + book.TableName()
}

func (book *Book) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(book.TableNameWithPrefix())
}

func NewBook() *Book {
	return &Book{}
}

//添加一个项目
func (book *Book) Insert() error {
	o := orm.NewOrm()
	//	o.Begin()
	book.BookName = utils.StripTags(book.BookName)
	if book.ItemId <= 0 {
		book.ItemId = 1
	}
	_, err := o.Insert(book)

	if err == nil {
		if book.Label != "" {
			NewLabel().InsertOrUpdateMulti(book.Label)
		}

		relationship := NewRelationship()
		relationship.BookId = book.BookId
		relationship.RoleId = 0
		relationship.MemberId = book.MemberId
		err = relationship.Insert()
		if err != nil {
			logs.Error("插入项目与用户关联 -> ", err)
			//o.Rollback()
			return err
		}
		document := NewDocument()
		document.BookId = book.BookId
		document.DocumentName = "空白文档"
		document.MemberId = book.MemberId
		err = document.InsertOrUpdate()
		if err != nil {
			//o.Rollback()
			return err
		}
		//o.Commit()
		return nil
	}
	//o.Rollback()
	return err
}

func (book *Book) Find(id int, cols ...string) (*Book, error) {
	if id <= 0 {
		return book, ErrInvalidParameter
	}
	o := orm.NewOrm()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("book_id", id).One(book, cols...)

	return book, err
}

//更新一个项目
func (book *Book) Update(cols ...string) error {
	o := orm.NewOrm()

	book.BookName = utils.StripTags(book.BookName)
	temp := NewBook()
	temp.BookId = book.BookId

	if err := o.Read(temp); err != nil {
		return err
	}

	if book.Label != "" || temp.Label != "" {

		go NewLabel().InsertOrUpdateMulti(book.Label + "," + temp.Label)
	}

	_, err := o.Update(book, cols...)
	return err
}

//复制项目
func (book *Book) Copy(identify string) error {
	o := orm.NewOrm()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("identify", identify).One(book)

	if err != nil {
		beego.Error("查询项目时出错 -> ", err)
		return err
	}
	if err := o.Begin(); err != nil {
		beego.Error("开启事物时出错 -> ", err)
		return err
	}

	bookId := book.BookId
	book.BookId = 0
	book.Identify = book.Identify + fmt.Sprintf("%s-%s", identify, strconv.FormatInt(time.Now().UnixNano(), 32))
	book.BookName = book.BookName + "[副本]"
	book.CreateTime = time.Now()
	book.CommentCount = 0
	book.HistoryCount = 0

	if _, err := o.Insert(book); err != nil {
		beego.Error("复制项目时出错 -> ", err)
		o.Rollback()
		return err
	}
	var rels []*Relationship

	if _, err := o.QueryTable(NewRelationship().TableNameWithPrefix()).Filter("book_id", bookId).All(&rels); err != nil {
		beego.Error("复制项目关系时出错 -> ", err)
		o.Rollback()
		return err
	}

	for _, rel := range rels {
		rel.BookId = book.BookId
		rel.RelationshipId = 0
		if _, err := o.Insert(rel); err != nil {
			beego.Error("复制项目关系时出错 -> ", err)
			o.Rollback()
			return err
		}
	}

	var docs []*Document

	if _, err := o.QueryTable(NewDocument().TableNameWithPrefix()).Filter("book_id", bookId).Filter("parent_id", 0).All(&docs); err != nil && err != orm.ErrNoRows {
		beego.Error("读取项目文档时出错 -> ", err)
		o.Rollback()
		return err
	}
	if len(docs) > 0 {
		if err := recursiveInsertDocument(docs, o, book.BookId, 0); err != nil {
			beego.Error("复制项目时出错 -> ", err)
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}

//递归的复制文档
func recursiveInsertDocument(docs []*Document, o orm.Ormer, bookId int, parentId int) error {
	for _, doc := range docs {

		docId := doc.DocumentId
		doc.DocumentId = 0
		doc.ParentId = parentId
		doc.BookId = bookId
		doc.Version = time.Now().Unix()

		if _, err := o.Insert(doc); err != nil {
			beego.Error("插入项目时出错 -> ", err)
			return err
		}

		var attachList []*Attachment
		//读取所有附件列表
		if _, err := o.QueryTable(NewAttachment().TableNameWithPrefix()).Filter("document_id", docId).All(&attachList); err == nil {
			for _, attach := range attachList {
				attach.BookId = bookId
				attach.DocumentId = doc.DocumentId
				attach.AttachmentId = 0
				if _, err := o.Insert(attach); err != nil {
					return err
				}
			}
		}
		var subDocs []*Document

		if _, err := o.QueryTable(NewDocument().TableNameWithPrefix()).Filter("parent_id", docId).All(&subDocs); err != nil && err != orm.ErrNoRows {
			beego.Error("读取文档时出错 -> ", err)
			return err
		}
		if len(subDocs) > 0 {

			if err := recursiveInsertDocument(subDocs, o, bookId, doc.DocumentId); err != nil {
				return err
			}
		}
	}
	return nil
}

//根据指定字段查询结果集.
func (book *Book) FindByField(field string, value interface{}, cols ...string) ([]*Book, error) {
	o := orm.NewOrm()

	var books []*Book
	_, err := o.QueryTable(book.TableNameWithPrefix()).Filter(field, value).All(&books, cols...)

	return books, err
}

//根据指定字段查询一个结果.
func (book *Book) FindByFieldFirst(field string, value interface{}) (*Book, error) {
	o := orm.NewOrm()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter(field, value).One(book)

	return book, err

}

//根据项目标识查询项目
func (book *Book) FindByIdentify(identify string, cols ...string) (*Book, error) {
	o := orm.NewOrm()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("identify", identify).One(book, cols...)

	return book, err
}

//分页查询指定用户的项目
func (book *Book) FindToPager(pageIndex, pageSize, memberId int) (books []*BookResult, totalCount int, err error) {

	o := orm.NewOrm()

	//sql1 := "SELECT COUNT(book.book_id) AS total_count FROM " + book.TableNameWithPrefix() + " AS book LEFT JOIN " +
	//	relationship.TableNameWithPrefix() + " AS rel ON book.book_id=rel.book_id AND rel.member_id = ? WHERE rel.relationship_id > 0 "

	sql1 := `SELECT
count(*) AS total_count
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) as role_id
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
					as t group by t.book_id)
			as team on team.book_id=book.book_id WHERE rel.role_id >= 0 or team.role_id >= 0`

	err = o.Raw(sql1, memberId, memberId).QueryRow(&totalCount)

	if err != nil {
		return
	}

	offset := (pageIndex - 1) * pageSize

	//sql2 := "SELECT book.*,rel.member_id,rel.role_id,m.account as create_name FROM " + book.TableNameWithPrefix() + " AS book" +
	//	" LEFT JOIN " + relationship.TableNameWithPrefix() + " AS rel ON book.book_id=rel.book_id AND rel.member_id = ?" +
	//	" LEFT JOIN " + relationship.TableNameWithPrefix() + " AS rel1 ON book.book_id=rel1.book_id  AND rel1.role_id=0" +
	//	" LEFT JOIN " + NewMember().TableNameWithPrefix() + " AS m ON rel1.member_id=m.member_id " +
	//	" WHERE rel.relationship_id > 0 ORDER BY book.order_index DESC,book.book_id DESC LIMIT " + fmt.Sprintf("%d,%d", offset, pageSize)

	sql2 := `SELECT
  book.*,
  case when rel.relationship_id  is null then team.role_id else rel.role_id end as role_id,
  m.account as create_name
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) as role_id
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
					as t group by book_id) as team 
			on team.book_id=book.book_id
  LEFT JOIN md_relationship AS rel1 ON book.book_id = rel1.book_id AND rel1.role_id = 0
  LEFT JOIN md_members AS m ON rel1.member_id = m.member_id
WHERE rel.role_id >= 0 or team.role_id >= 0
ORDER BY book.order_index, book.book_id DESC limit ?,?`

	_, err = o.Raw(sql2, memberId, memberId, offset, pageSize).QueryRows(&books)
	if err != nil {
		logs.Error("分页查询项目列表 => ", err)
		return
	}
	sql := "SELECT m.account,doc.modify_time FROM md_documents AS doc LEFT JOIN md_members AS m ON doc.modify_at=m.member_id WHERE book_id = ? LIMIT 1 ORDER BY doc.modify_time DESC"

	if err == nil && len(books) > 0 {
		for index, book := range books {
			var text struct {
				Account    string
				ModifyTime time.Time
			}

			err1 := o.Raw(sql, book.BookId).QueryRow(&text)
			if err1 == nil {
				books[index].LastModifyText = text.Account + " 于 " + text.ModifyTime.Format("2006-01-02 15:04:05")
			}
			if book.RoleId == 0 {
				book.RoleName = "创始人"
			} else if book.RoleId == 1 {
				book.RoleName = "管理员"
			} else if book.RoleId == 2 {
				book.RoleName = "编辑者"
			} else if book.RoleId == 3 {
				book.RoleName = "观察者"
			}
		}
	}
	return
}

// 彻底删除项目.
func (book *Book) ThoroughDeleteBook(id int) error {
	if id <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	book, err := book.Find(id)
	if err != nil {
		return err
	}
	o.Begin()

	//删除附件,这里没有删除实际物理文件
	_, err = o.Raw("DELETE FROM "+NewAttachment().TableNameWithPrefix()+" WHERE book_id=?", book.BookId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	//删除文档
	_, err = o.Raw("DELETE FROM "+NewDocument().TableNameWithPrefix()+" WHERE book_id = ?", book.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}
	//删除项目
	_, err = o.Raw("DELETE FROM "+book.TableNameWithPrefix()+" WHERE book_id = ?", book.BookId).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	//删除关系
	_, err = o.Raw("DELETE FROM "+NewRelationship().TableNameWithPrefix()+" WHERE book_id = ?", book.BookId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Raw(fmt.Sprintf("DELETE FROM %s WHERE book_id=?", NewTeamRelationship().TableNameWithPrefix()), book.BookId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//删除模板
	_, err = o.Raw("DELETE FROM "+NewTemplate().TableNameWithPrefix()+" WHERE book_id = ?", book.BookId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	if book.Label != "" {
		NewLabel().InsertOrUpdateMulti(book.Label)
	}

	//删除导出缓存
	if err := os.RemoveAll(filepath.Join(conf.GetExportOutputPath(), strconv.Itoa(id))); err != nil {
		beego.Error("删除项目缓存失败 ->", err)
	}
	//删除附件和图片
	if err := os.RemoveAll(filepath.Join(conf.WorkingDirectory, "uploads", book.Identify)); err != nil {
		beego.Error("删除项目附件和图片失败 ->", err)
	}

	return o.Commit()

}

//分页查找系统首页数据.
func (book *Book) FindForHomeToPager(pageIndex, pageSize, memberId int) (books []*BookResult, totalCount int, err error) {
	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize
	//如果是登录用户
	if memberId > 0 {
		sql1 := `SELECT COUNT(*)
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) AS role_id
             from (select book_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
as t group by book_id) as team on team.book_id=book.book_id
WHERE book.privately_owned = 0 or rel.role_id >=0 or team.role_id >=0`
		err = o.Raw(sql1, memberId, memberId).QueryRow(&totalCount)
		if err != nil {
			return
		}
		sql2 := `SELECT book.*,rel1.*,member.account AS create_name,member.real_name FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) AS role_id
             from (select book_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
as t group by book_id) as team on team.book_id=book.book_id
  LEFT JOIN md_relationship AS rel1 ON rel1.book_id = book.book_id AND rel1.role_id = 0
  LEFT JOIN md_members AS member ON rel1.member_id = member.member_id
WHERE book.privately_owned = 0 or rel.role_id >=0 or team.role_id >=0 ORDER BY order_index desc,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql2, memberId, memberId, offset, pageSize).QueryRows(&books)

	} else {
		count, err1 := o.QueryTable(book.TableNameWithPrefix()).Filter("privately_owned", 0).Count()

		if err1 != nil {
			err = err1
			return
		}
		totalCount = int(count)

		sql := `SELECT book.*,rel.*,member.account AS create_name,member.real_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS member ON rel.member_id = member.member_id
			WHERE book.privately_owned = 0 ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql, offset, pageSize).QueryRows(&books)

	}
	return
}

//分页全局搜索.
func (book *Book) FindForLabelToPager(keyword string, pageIndex, pageSize, memberId int) (books []*BookResult, totalCount int, err error) {
	o := orm.NewOrm()

	keyword = "%" + keyword + "%"
	offset := (pageIndex - 1) * pageSize
	//如果是登录用户
	if memberId > 0 {
		sql1 := `SELECT COUNT(*)
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
  left join (select *
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.role_id,t.team_member_id,t.book_id) as team on team.book_id = book.book_id
WHERE (relationship_id > 0 OR book.privately_owned = 0 or team.team_member_id > 0) AND book.label LIKE ?`

		err = o.Raw(sql1, memberId, memberId, keyword).QueryRow(&totalCount)
		if err != nil {
			return
		}
		sql2 := `SELECT book.*,rel1.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
			left join (select * from (select book_id,team_member_id,role_id
                   	from md_team_relationship as mtr
					left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.role_id,t.team_member_id,t.book_id) as team 
					on team.book_id = book.book_id
			LEFT JOIN md_relationship AS rel1 ON rel1.book_id = book.book_id AND rel1.role_id = 0
			LEFT JOIN md_members AS member ON rel1.member_id = member.member_id
			WHERE (rel.relationship_id > 0 OR book.privately_owned = 0 or team.team_member_id > 0) 
			AND book.label LIKE ? ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql2, memberId, memberId, keyword, offset, pageSize).QueryRows(&books)

		return

	} else {
		count, err1 := o.QueryTable(NewBook().TableNameWithPrefix()).Filter("privately_owned", 0).Filter("label__icontains", keyword).Count()

		if err1 != nil {
			err = err1
			return
		}
		totalCount = int(count)

		sql := `SELECT book.*,rel.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS member ON rel.member_id = member.member_id
			WHERE book.privately_owned = 0 AND book.label LIKE ? ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql, keyword, offset, pageSize).QueryRows(&books)

		return

	}
}

//发布文档
func (book *Book) ReleaseContent(bookId int) {

	o := orm.NewOrm()

	var docs []*Document
	_, err := o.QueryTable(NewDocument().TableNameWithPrefix()).Filter("book_id", bookId).All(&docs)

	if err != nil {
		beego.Error("发布失败 =>", bookId, err)
		return
	}
	for _, item := range docs {
		item.BookId = bookId
		_ = item.ReleaseContent()
	}
}

//重置文档数量
func (book *Book) ResetDocumentNumber(bookId int) {
	o := orm.NewOrm()

	totalCount, err := o.QueryTable(NewDocument().TableNameWithPrefix()).Filter("book_id", bookId).Count()
	if err == nil {
		_, err = o.Raw("UPDATE md_books SET doc_count = ? WHERE book_id = ?", int(totalCount), bookId).Exec()
		if err != nil {
			beego.Error("重置文档数量失败 =>", bookId, err)
		}
	} else {
		beego.Error("获取文档数量失败 =>", bookId, err)
	}
}

//导入项目
func (book *Book) ImportBook(zipPath string) error {
	if !filetil.FileExists(zipPath) {
		return errors.New("文件不存在 => " + zipPath)
	}

	w := md5.New()
	io.WriteString(w, zipPath) //将str写入到w中
	io.WriteString(w, time.Now().String())
	io.WriteString(w, book.BookName)
	md5str := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式

	tempPath := filepath.Join(os.TempDir(), md5str)

	if err := os.MkdirAll(tempPath, 0766); err != nil {
		beego.Error("创建导入目录出错 => ", err)
	}
	//如果加压缩失败
	if err := ziptil.Unzip(zipPath, tempPath); err != nil {
		return err
	}
	//当导入结束后，删除临时文件
	//defer os.RemoveAll(tempPath)

	for {
		//如果当前目录下只有一个目录，则重置根目录
		if entries, err := ioutil.ReadDir(tempPath); err == nil && len(entries) == 1 {
			dir := entries[0]
			if dir.IsDir() && dir.Name() != "." && dir.Name() != ".." {
				tempPath = filepath.Join(tempPath, dir.Name())
				break
			}

		} else {
			break
		}
	}

	tempPath = strings.Replace(tempPath, "\\", "/", -1)

	docMap := make(map[string]int, 0)

	o := orm.NewOrm()

	o.Insert(book)
	relationship := NewRelationship()
	relationship.BookId = book.BookId
	relationship.RoleId = 0
	relationship.MemberId = book.MemberId
	relationship.Insert()

	err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
		path = strings.Replace(path, "\\", "/", -1)
		if path == tempPath {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			//如果是Markdown文件
			if strings.EqualFold(ext, ".md") || strings.EqualFold(ext, ".markdown") {
				beego.Info("正在处理 =>", path, info.Name())
				doc := NewDocument()
				doc.BookId = book.BookId
				doc.MemberId = book.MemberId
				docIdentify := strings.Replace(strings.TrimPrefix(path, tempPath+"/"), "/", "-", -1)

				if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, docIdentify); !ok || err != nil {
					docIdentify = "import-" + docIdentify
				}

				doc.Identify = docIdentify
				//匹配图片，如果图片语法是在代码块中，这里同样会处理
				re := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
				markdown, err := filetil.ReadFileAndIgnoreUTF8BOM(path)
				if err != nil {
					return err
				}

				//处理图片
				doc.Markdown = re.ReplaceAllStringFunc(string(markdown), func(image string) string {

					images := re.FindAllSubmatch([]byte(image), -1)
					if len(images) <= 0 || len(images[0]) < 3 {
						return image
					}
					originalImageUrl := string(images[0][2])
					imageUrl := strings.Replace(string(originalImageUrl), "\\", "/", -1)

					//如果是本地路径，则需要将图片复制到项目目录
					if !strings.HasPrefix(imageUrl, "http://") &&
						!strings.HasPrefix(imageUrl, "https://") &&
						!strings.HasPrefix(imageUrl, "ftp://") {
						//如果路径中存在参数
						if l := strings.Index(imageUrl, "?"); l > 0 {
							imageUrl = imageUrl[:l]
						}

						if strings.HasPrefix(imageUrl, "/") {
							imageUrl = filepath.Join(tempPath, imageUrl)
						} else if strings.HasPrefix(imageUrl, "./") {
							imageUrl = filepath.Join(filepath.Dir(path), strings.TrimPrefix(imageUrl, "./"))
						} else if strings.HasPrefix(imageUrl, "../") {
							imageUrl = filepath.Join(filepath.Dir(path), imageUrl)
						} else {
							imageUrl = filepath.Join(filepath.Dir(path), imageUrl)
						}
						imageUrl = strings.Replace(imageUrl, "\\", "/", -1)
						dstFile := filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), strings.TrimPrefix(imageUrl, tempPath))

						if filetil.FileExists(imageUrl) {
							filetil.CopyFile(imageUrl, dstFile)

							imageUrl = strings.TrimPrefix(strings.Replace(dstFile, "\\", "/", -1), strings.Replace(conf.WorkingDirectory, "\\", "/", -1))

							if !strings.HasPrefix(imageUrl, "/") && !strings.HasPrefix(imageUrl, "\\") {
								imageUrl = "/" + imageUrl
							}
						}

					} else {
						imageExt := cryptil.Md5Crypt(imageUrl) + filepath.Ext(imageUrl)

						dstFile := filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), imageExt)

						if err := requests.DownloadAndSaveFile(imageUrl, dstFile); err == nil {
							imageUrl = strings.TrimPrefix(strings.Replace(dstFile, "\\", "/", -1), strings.Replace(conf.WorkingDirectory, "\\", "/", -1))
							if !strings.HasPrefix(imageUrl, "/") && !strings.HasPrefix(imageUrl, "\\") {
								imageUrl = "/" + imageUrl
							}
						}
					}

					imageUrl = strings.Replace(strings.TrimSuffix(image, originalImageUrl+")")+conf.URLForWithCdnImage(imageUrl)+")", "\\", "/", -1)
					return imageUrl
				})

				linkRegexp := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

				//处理链接
				doc.Markdown = linkRegexp.ReplaceAllStringFunc(doc.Markdown, func(link string) string {
					links := linkRegexp.FindAllStringSubmatch(link, -1)
					originalLink := links[0][2]
					var linkPath string
					var err error
					if strings.HasPrefix(originalLink, "<") {
						originalLink = strings.TrimPrefix(originalLink, "<")
					}
					if strings.HasSuffix(originalLink, ">") {
						originalLink = strings.TrimSuffix(originalLink, ">")
					}
					//如果是从根目录开始，
					if strings.HasPrefix(originalLink, "/") {
						linkPath, err = filepath.Abs(filepath.Join(tempPath, originalLink))
					} else if strings.HasPrefix(originalLink, "./") {
						linkPath, err = filepath.Abs(filepath.Join(filepath.Dir(path), originalLink[1:]))
					} else {
						linkPath, err = filepath.Abs(filepath.Join(filepath.Dir(path), originalLink))
					}

					if err == nil {
						//如果本地存在该链接
						if filetil.FileExists(linkPath) {
							ext := filepath.Ext(linkPath)
							//beego.Info("当前后缀 -> ",ext)
							//如果链接是Markdown文件，则生成文档标识,否则，将目标文件复制到项目目录
							if strings.EqualFold(ext, ".md") || strings.EqualFold(ext, ".markdown") {
								docIdentify := strings.Replace(strings.TrimPrefix(strings.Replace(linkPath, "\\", "/", -1), tempPath+"/"), "/", "-", -1)
								//beego.Info(originalLink, "|", linkPath, "|", docIdentify)
								if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, docIdentify); !ok || err != nil {
									docIdentify = "import-" + docIdentify
								}
								docIdentify = strings.TrimSuffix(docIdentify, "-README.md")

								link = strings.TrimSuffix(link, originalLink+")") + conf.URLFor("DocumentController.Read", ":key", book.Identify, ":id", docIdentify) + ")"

							} else {
								dstPath := filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), originalLink)

								filetil.CopyFile(linkPath, dstPath)

								tempLink := conf.BaseUrl + strings.TrimPrefix(strings.Replace(dstPath, "\\", "/", -1), strings.Replace(conf.WorkingDirectory, "\\", "/", -1))

								link = strings.TrimSuffix(link, originalLink+")") + tempLink + ")"

							}
						} else {
							beego.Info("文件不存在 ->", linkPath)
						}
					}

					return link
				})

				//codeRe := regexp.MustCompile("```\\w+")

				//doc.Markdown = codeRe.ReplaceAllStringFunc(doc.Markdown, func(s string) string {
				//	//beego.Info(s)
				//	return strings.Replace(s,"```","``` ",-1)
				//})

				doc.Content = string(blackfriday.Run([]byte(doc.Markdown)))

				doc.Version = time.Now().Unix()

				//解析文档名称，默认使用第一个h标签为标题
				docName := strings.TrimSuffix(info.Name(), ext)

				for _, line := range strings.Split(doc.Markdown, "\n") {
					if strings.HasPrefix(line, "#") {
						docName = strings.TrimLeft(line, "#")
						break
					}
				}

				doc.DocumentName = strings.TrimSpace(docName)

				parentId := 0

				parentIdentify := strings.Replace(strings.Trim(strings.TrimSuffix(strings.TrimPrefix(path, tempPath), info.Name()), "/"), "/", "-", -1)

				if parentIdentify != "" {

					if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, parentIdentify); !ok || err != nil {
						parentIdentify = "import-" + parentIdentify
					}
					if id, ok := docMap[parentIdentify]; ok {
						parentId = id
					}
				}
				if strings.EqualFold(info.Name(), "README.md") {
					beego.Info(path, "|", info.Name(), "|", parentIdentify, "|", parentId)
				}
				isInsert := false
				//如果当前文件是README.md，则将内容更新到父级
				if strings.EqualFold(info.Name(), "README.md") && parentId != 0 {

					doc.DocumentId = parentId
					//beego.Info(path,"|",parentId)
				} else {
					//beego.Info(path,"|",parentIdentify)
					doc.ParentId = parentId
					isInsert = true
				}
				if err := doc.InsertOrUpdate("document_name", "markdown", "content"); err != nil {
					beego.Error(doc.DocumentId, err)
				}
				if isInsert {
					docMap[docIdentify] = doc.DocumentId
				}
			}
		} else {
			//如果当前目录下存在Markdown文件，则需要创建此节点
			if filetil.HasFileOfExt(path, []string{".md", ".markdown"}) {
				beego.Info("正在处理 =>", path, info.Name())
				identify := strings.Replace(strings.Trim(strings.TrimPrefix(path, tempPath), "/"), "/", "-", -1)
				if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, identify); !ok || err != nil {
					identify = "import-" + identify
				}

				parentDoc := NewDocument()

				parentDoc.MemberId = book.MemberId
				parentDoc.BookId = book.BookId
				parentDoc.Identify = identify
				parentDoc.Version = time.Now().Unix()
				parentDoc.DocumentName = "空白文档"

				parentId := 0

				parentIdentify := strings.TrimSuffix(identify, "-"+info.Name())

				if id, ok := docMap[parentIdentify]; ok {
					parentId = id
				}

				parentDoc.ParentId = parentId

				if err := parentDoc.InsertOrUpdate(); err != nil {
					beego.Error(err)
				}

				docMap[identify] = parentDoc.DocumentId
				//beego.Info(path,"|",parentDoc.DocumentId,"|",identify,"|",info.Name(),"|",parentIdentify)
			}
		}

		return nil
	})

	if err != nil {
		beego.Error("导入项目异常 => ", err)
		book.Description = "【项目导入存在错误：" + err.Error() + "】"
	}
	beego.Info("项目导入完毕 => ", book.BookName)
	book.ReleaseContent(book.BookId)
	return err
}

func (book *Book) FindForRoleId(bookId, memberId int) (conf.BookRole, error) {
	o := orm.NewOrm()

	var relationship Relationship

	err := NewRelationship().QueryTable().Filter("book_id", bookId).Filter("member_id", memberId).One(&relationship)

	if err != nil && err != orm.ErrNoRows {
		return 0, err
	}
	if err == nil {
		return relationship.RoleId, nil
	}
	sql := `select role_id
from md_team_relationship as mtr
left join md_team_member as mtm using (team_id)
where mtr.book_id = ? and mtm.member_id = ? order by mtm.role_id asc limit 1;`

	var roleId int
	err = o.Raw(sql, bookId, memberId).QueryRow(&roleId)

	if err != nil {
		beego.Error("查询用户项目角色出错 -> book_id=", bookId, " member_id=", memberId, err)
		return 0, err
	}
	return conf.BookRole(roleId), nil
}
