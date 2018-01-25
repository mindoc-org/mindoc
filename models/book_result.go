package models

import (
	"time"
	"bytes"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/mindoc/conf"
	"strings"
	"github.com/lifei6671/mindoc/converter"
	"strconv"
	"github.com/russross/blackfriday"
	"path/filepath"
	"github.com/astaxie/beego"
	"os"
	"github.com/PuerkitoBio/goquery"
	"github.com/lifei6671/mindoc/utils"
)

type BookResult struct {
	BookId int              `json:"book_id"`
	BookName string         `json:"book_name"`
	Identify string         `json:"identify"`
	OrderIndex int          `json:"order_index"`
	Description string      `json:"description"`
	Publisher string 		`json:"publisher"`
	IsCacheEBook bool		`json:"is_cache_ebook"`
	PrivatelyOwned int      `json:"privately_owned"`
	PrivateToken string     `json:"private_token"`
	DocCount int            `json:"doc_count"`
	CommentStatus string    `json:"comment_status"`
	CommentCount int        `json:"comment_count"`
	CreateTime time.Time    `json:"create_time"`
	CreateName string 		`json:"create_name"`
	ModifyTime time.Time	`json:"modify_time"`
	Cover string            `json:"cover"`
	Theme string			`json:"theme"`
	Label string			`json:"label"`
	MemberId int            `json:"member_id"`
	Editor string           `json:"editor"`
	AutoRelease bool		`json:"auto_release"`

	RelationshipId int	`json:"relationship_id"`
	RoleId int        	`json:"role_id"`
	RoleName string 	`json:"role_name"`
	Status int

	LastModifyText string 	`json:"last_modify_text"`
	IsDisplayComment bool	`json:"is_display_comment"`
}

func NewBookResult() *BookResult {
	return &BookResult{}
}


// 根据项目标识查询项目以及指定用户权限的信息.
func (m *BookResult) FindByIdentify(identify string,member_id int) (*BookResult,error) {
	if identify == "" || member_id <= 0 {
		return m,ErrInvalidParameter
	}
	o := orm.NewOrm()

	book := NewBook()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("identify", identify).One(book)

	if err != nil {
		return m, err
	}

	relationship := NewRelationship()

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id", book.BookId).Filter("member_id", member_id).One(relationship)

	if err != nil {
		return m, err
	}
	var relationship2 Relationship

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id", book.BookId).Filter("role_id", 0).One(&relationship2)

	if err != nil {
		logs.Error("根据项目标识查询项目以及指定用户权限的信息 => ", err)
		return m, ErrPermissionDenied
	}

	member, err := NewMember().Find(relationship2.MemberId)
	if err != nil {
		return m, err
	}

	m = NewBookResult().ToBookResult(*book)

	m.CreateName 		= member.Account
	m.MemberId 		= relationship.MemberId
	m.RoleId		= relationship.RoleId
	m.RelationshipId	= relationship.RelationshipId


	if m.RoleId == conf.BookFounder {
		m.RoleName = "创始人"
	} else if m.RoleId == conf.BookAdmin {
		m.RoleName = "管理员"
	} else if m.RoleId == conf.BookEditor {
		m.RoleName = "编辑者"
	} else if m.RoleId == conf.BookObserver {
		m.RoleName = "观察者"
	}

	doc := NewDocument()

	err = o.QueryTable(doc.TableNameWithPrefix()).Filter("book_id", book.BookId).OrderBy("modify_time").One(doc)

	if err == nil {
		member2 := NewMember()
		member2.Find(doc.ModifyAt)

		m.LastModifyText = member2.Account + " 于 " + doc.ModifyTime.Format("2006-01-02 15:04:05")
	}

	return m, nil
}

func (m *BookResult) FindToPager(pageIndex, pageSize int) (books []*BookResult,totalCount int,err error)  {
	o := orm.NewOrm()

	count,err := o.QueryTable(NewBook().TableNameWithPrefix()).Count()

	if err != nil {
		return
	}
	totalCount = int(count)

	sql := `SELECT
			book.*,rel.relationship_id,rel.role_id,m.account AS create_name
		FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS m ON rel.member_id = m.member_id
		ORDER BY book.order_index DESC ,book.book_id DESC  LIMIT ?,?`

	offset := (pageIndex -1 )* pageSize

	_,err = o.Raw(sql,offset,pageSize).QueryRows(&books)

	return
}

//实体转换
func (m *BookResult) ToBookResult(book Book) *BookResult {

	m.BookId = book.BookId
	m.BookName = book.BookName
	m.Identify = book.Identify
	m.OrderIndex = book.OrderIndex
	m.Description = strings.Replace(book.Description, "\r\n", "<br/>", -1)
	m.PrivatelyOwned = book.PrivatelyOwned
	m.PrivateToken = book.PrivateToken
	m.DocCount = book.DocCount
	m.CommentStatus = book.CommentStatus
	m.CommentCount = book.CommentCount
	m.CreateTime = book.CreateTime
	m.ModifyTime = book.ModifyTime
	m.Cover = book.Cover
	m.Label = book.Label
	m.Status = book.Status
	m.Editor = book.Editor
	m.Theme = book.Theme
	m.AutoRelease = book.AutoRelease == 1
	m.Publisher = book.Publisher
	m.IsCacheEBook = book.IsCacheEBook == 1

	if book.Theme == "" {
		m.Theme = "default"
	}
	if book.Editor == "" {
		m.Editor = "markdown"
	}
	return m
}

func (m *BookResult) Converter(sessionId string) (ConvertBookResult,error) {

	convertBookResult := ConvertBookResult{}
	outputPath := filepath.Join(beego.AppConfig.DefaultString("book_output_path","cache"),sessionId,strconv.Itoa(m.BookId))

	if m.IsCacheEBook {
		outputPath = filepath.Join(beego.AppConfig.DefaultString("book_output_path","cache"),strconv.Itoa(m.BookId))
	}

	if m.IsCacheEBook {
		pdfpath := filepath.Join(outputPath,"output","book.pdf")
		epubpath := filepath.Join(outputPath,"output","book.epub")
		mobipath := filepath.Join(outputPath,"output","book.mobi")

		if utils.FileExists(pdfpath) && utils.FileExists(epubpath) && utils.FileExists(mobipath){
			convertBookResult.EpubPath = epubpath
			convertBookResult.MobiPath = mobipath
			convertBookResult.PDFPath = pdfpath
			return convertBookResult,nil
		}
	}
	docs, err := NewDocument().FindListByBookId(m.BookId)
	if err != nil {
		return convertBookResult,err
	}

	tocList := make([]converter.Toc,0)

	for _, item := range docs {
		if item.ParentId == 0 {
			toc := converter.Toc{
				Id: item.DocumentId,
				Link: strconv.Itoa(item.DocumentId) + ".html",
				Pid: item.ParentId,
				Title: item.DocumentName,
			}

			tocList = append(tocList,toc)
		}
	}
	for _, item := range docs {
		if item.ParentId != 0 {
			toc := converter.Toc{
				Id: item.DocumentId,
				Link: strconv.Itoa(item.DocumentId) + ".html",
				Pid: item.ParentId,
				Title: item.DocumentName,
			}
			tocList = append(tocList,toc)
		}
	}

	ebookConfig := converter.Config{
		Charset :      "utf-8",
		Cover :        m.Cover,
		Timestamp :    time.Now().Format("2006-01-02 15:04:05"),
		Description :  string(blackfriday.MarkdownBasic([]byte(m.Description))),
		Footer :       "<p style='color:#8E8E8E;font-size:12px;'>本文档使用 <a href='https://www.iminho.me' style='text-decoration:none;color:#1abc9c;font-weight:bold;'>MinDoc</a> 构建 <span style='float:right'>- _PAGENUM_ -</span></p>",
		Header :       "<p style='color:#8E8E8E;font-size:12px;'>_SECTION_</p>",
		Identifier :   "",
		Language :     "zh-CN",
		Creator :      m.CreateName,
		Publisher :    m.Publisher,
		Contributor :  m.Publisher,
		Title :        m.BookName,
		Format:        []string{"epub", "mobi", "pdf"},
		FontSize :     "14",
		PaperSize :    "a4",
		MarginLeft :   "72",
		MarginRight :  "72",
		MarginTop :    "72",
		MarginBottom : "72",
		Toc :          tocList,
		More :         []string{},

	}


	os.MkdirAll(outputPath, 0766)
	if outputPath, err = filepath.Abs(outputPath); err != nil {
		beego.Error("导出目录配置错误：" + err.Error())
		return convertBookResult,err
	}

	viewPath := beego.BConfig.WebConfig.ViewsPath
	baseUrl := beego.AppConfig.DefaultString("baseurl","")

	for _,item := range docs {
		name := strconv.Itoa(item.DocumentId)
		fpath := filepath.Join(outputPath,name + ".html")

		f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return convertBookResult,err
		}
		var buf bytes.Buffer

		if err := beego.ExecuteViewPathTemplate(&buf,"document/export.tpl",viewPath,map[string]interface{}{"Model": m, "Lists": item, "BaseUrl": baseUrl}); err != nil {
			return convertBookResult,err
		}
		html := buf.String()


		if err != nil {

			f.Close()
			return convertBookResult,err
		}

		bufio := bytes.NewReader(buf.Bytes())

		doc, err := goquery.NewDocumentFromReader(bufio)
		doc.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
			if src, ok := contentSelection.Attr("src"); ok && strings.HasPrefix(src, "/uploads/") {
				contentSelection.SetAttr("src", baseUrl + src)
			}
		})

		html, err = doc.Html()
		if err != nil {
			f.Close()
			return convertBookResult,err
		}

		// html = strings.Replace(html, "<img src=\"/uploads", "<img src=\"" + c.BaseUrl() + "/uploads", -1)

		f.WriteString(html)
		f.Close()
	}
	eBookConverter := &converter.Converter{
		BasePath : outputPath,
		Config :   ebookConfig,
		Debug :    false,
	}

	if err := eBookConverter.Convert();err != nil {
		beego.Error("转换文件错误：" + m.BookName +" => "+ err.Error())
		return convertBookResult,err
	}
	convertBookResult.MobiPath = filepath.Join(outputPath,"output","book.mobi")
	convertBookResult.PDFPath = filepath.Join(outputPath,"output","book.pdf")
	convertBookResult.EpubPath = filepath.Join(outputPath,"output","book.epub")
	return  convertBookResult,nil
}


















