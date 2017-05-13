package controllers

import (
	"container/list"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"image/png"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/models"
	"github.com/lifei6671/godoc/utils/wkhtmltopdf"
)

//DocumentController struct.
type DocumentController struct {
	BaseController
}

//判断用户是否可以阅读文档.
func isReadable(identify, token string, c *DocumentController) *models.BookResult {
	book, err := models.NewBook().FindByFieldFirst("identify", identify)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	if c.Member != nil && c.Member.Role == conf.MemberSuperRole {
		bookResult := book.ToBookResult()
		return bookResult
	}
	//如果文档是私有的
	if book.PrivatelyOwned == 1 {

		is_ok := false

		if c.Member != nil {
			_, err := models.NewRelationship().FindForRoleId(book.BookId, c.Member.MemberId)
			if err == nil {
				is_ok = true
			}
		}
		if book.PrivateToken != "" && !is_ok {
			//如果有访问的Token，并且该项目设置了访问Token，并且和用户提供的相匹配，则记录到Session中.
			//如果用户未提供Token且用户登录了，则判断用户是否参与了该项目.
			//如果用户未登录，则从Session中读取Token.
			if token != "" && strings.EqualFold(token, book.PrivateToken) {
				c.SetSession(identify, token)

			} else if token, ok := c.GetSession(identify).(string); !ok || !strings.EqualFold(token, book.PrivateToken) {
				c.Abort("403")
			}
		} else if !is_ok {
			c.Abort("403")
		}

	}
	bookResult := book.ToBookResult()

	if c.Member != nil {

		rel, err := models.NewRelationship().FindByBookIdAndMemberId(bookResult.BookId, c.Member.MemberId)

		if err == nil {
			bookResult.MemberId = rel.MemberId
			bookResult.RoleId = rel.RoleId
			bookResult.RelationshipId = rel.RelationshipId
		}

	}
	//判断是否需要显示评论框
	if bookResult.CommentStatus == "closed" {
		bookResult.IsDisplayComment = false
	} else if bookResult.CommentStatus == "open" {
		bookResult.IsDisplayComment = true
	} else if bookResult.CommentStatus == "group_only" {
		bookResult.IsDisplayComment = bookResult.RelationshipId > 0
	} else if bookResult.CommentStatus == "registered_only" {
		bookResult.IsDisplayComment = true
	}

	return bookResult
}

//文档首页.
func (c *DocumentController) Index() {
	c.Prepare()
	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")

	if identify == "" {
		c.Abort("404")
	}
	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}
	bookResult := isReadable(identify, token, c)

	c.TplName = "document/" + bookResult.Theme + "_read.tpl"

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, 0)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}

	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = "概要"
	c.Data["Content"] = bookResult.Description
}

//阅读文档.
func (c *DocumentController) Read() {
	c.Prepare()
	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	id := c.GetString(":id")

	if identify == "" || id == "" {
		c.Abort("404")
	}

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}

	bookResult := isReadable(identify, token, c)

	c.TplName = "document/" + bookResult.Theme + "_read.tpl"

	doc := models.NewDocument()

	if doc_id, err := strconv.Atoi(id); err == nil {
		doc, err = doc.Find(doc_id)
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
	} else {
		doc, err = doc.FindByFieldFirst("identify", id)
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
	}

	if doc.BookId != bookResult.BookId {
		c.Abort("403")
	}
	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	if c.IsAjax() {
		var data struct {
			DocTitle string `json:"doc_title"`
			Body     string `json:"body"`
			Title    string `json:"title"`
		}
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release
		data.Title = doc.DocumentName + " - Powered by MinDoc"

		c.JsonResult(0, "ok", data)
	}

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, doc.DocumentId)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}

	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = doc.DocumentName
	c.Data["Content"] = template.HTML(doc.Release)
}

//编辑文档.
func (c *DocumentController) Edit() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	bookResult := models.NewBookResult()
	var err error
	//如果是超级管理者，则不判断权限
	if c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		bookResult = book.ToBookResult()

	} else {
		bookResult, err = models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			beego.Error("DocumentController.Edit => ", err)

			c.Abort("403")
		}
		if bookResult.RoleId == conf.BookObserver {

			c.JsonResult(6002, "项目不存在或权限不足")
		}
	}

	//根据不同编辑器类型加载编辑器
	if bookResult.Editor == "markdown" {
		c.TplName = "document/markdown_edit_template.tpl"
	} else if bookResult.Editor == "html" {
		c.TplName = "document/html_edit_template.tpl"
	} else {
		c.TplName = "document/" + bookResult.Editor + "_edit_template.tpl"
	}

	c.Data["Model"] = bookResult

	r, _ := json.Marshal(bookResult)

	c.Data["ModelResult"] = template.JS(string(r))

	c.Data["Result"] = template.JS("[]")

	trees, err := models.NewDocument().FindDocumentTree(bookResult.BookId)

	if err != nil {
		beego.Error("FindDocumentTree => ", err)
	} else {
		if len(trees) > 0 {
			if jtree, err := json.Marshal(trees); err == nil {
				c.Data["Result"] = template.JS(string(jtree))
			}
		} else {
			c.Data["Result"] = template.JS("[]")
		}
	}
	c.Data["BaiDuMapKey"] = beego.AppConfig.DefaultString("baidumapkey", "")

}

//创建一个文档.
func (c *DocumentController) Create() {
	identify := c.GetString("identify")
	doc_identify := c.GetString("doc_identify")
	doc_name := c.GetString("doc_name")
	parent_id, _ := c.GetInt("parent_id", 0)
	doc_id, _ := c.GetInt("doc_id", 0)

	if identify == "" {
		c.JsonResult(6001, "参数错误")
	}
	if doc_name == "" {
		c.JsonResult(6004, "文档名称不能为空")
	}
	if doc_identify != "" {
		if ok, err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`, doc_identify); !ok || err != nil {

			c.JsonResult(6003, "文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头")
		}
		d, _ := models.NewDocument().FindByFieldFirst("identify", doc_identify)
		if d.DocumentId > 0 && d.DocumentId != doc_id {
			c.JsonResult(6006, "文档标识已被使用")
		}
	}
	book_id := 0
	//如果是超级管理员则不判断权限
	if c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error(err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = bookResult.BookId
	}
	if parent_id > 0 {
		doc, err := models.NewDocument().Find(parent_id)
		if err != nil || doc.BookId != book_id {
			c.JsonResult(6003, "父分类不存在")
		}
	}

	document, _ := models.NewDocument().Find(doc_id)

	document.MemberId = c.Member.MemberId
	document.BookId = book_id
	if doc_identify != "" {
		document.Identify = doc_identify
	}
	document.Version = time.Now().Unix()
	document.DocumentName = doc_name
	document.ParentId = parent_id

	if err := document.InsertOrUpdate(); err != nil {
		beego.Error("InsertOrUpdate => ", err)
		c.JsonResult(6005, "保存失败")
	} else {

		c.JsonResult(0, "ok", document)
	}
}

//上传附件或图片.
func (c *DocumentController) Upload() {

	identify := c.GetString("identify")
	doc_id, _ := c.GetInt("doc_id")
	is_attach := true

	if identify == "" {
		c.JsonResult(6001, "参数错误")
	}

	name := "editormd-file-file"

	file, moreFile, err := c.GetFile(name)
	if err == http.ErrMissingFile {
		name = "editormd-image-file"
		file, moreFile, err = c.GetFile(name)
		if err == http.ErrMissingFile {
			c.JsonResult(6003, "没有发现需要上传的文件")
		}
	}
	if err != nil {
		c.JsonResult(6002, err.Error())
	}

	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)

	if ext == "" {
		c.JsonResult(6003, "无法解析文件的格式")
	}

	if !conf.IsAllowUploadFileExt(ext) {
		c.JsonResult(6004, "不允许的文件类型")
	}
	book_id := 0
	//如果是超级管理员，则不判断权限
	if c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)

		if err != nil {
			c.JsonResult(6006, "文档不存在或权限不足")
		}
		book_id = book.BookId

	} else {
		book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			beego.Error("DocumentController.Edit => ", err)
			if err == orm.ErrNoRows {
				c.JsonResult(6006, "权限不足")
			}
			c.JsonResult(6001, err.Error())
		}
		//如果没有编辑权限
		if book.RoleId != conf.BookEditor && book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
			c.JsonResult(6006, "权限不足")
		}
		book_id = book.BookId
	}

	if doc_id > 0 {
		doc, err := models.NewDocument().Find(doc_id)
		if err != nil {
			c.JsonResult(6007, "文档不存在")
		}
		if doc.BookId != book_id {
			c.JsonResult(6008, "文档不属于指定的项目")
		}
	}

	fileName := "attach_" + strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := "uploads/" + time.Now().Format("200601") + "/" + fileName + ext

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile(name, filePath)

	if err != nil {
		beego.Error("SaveToFile => ", err)
		c.JsonResult(6005, "保存文件失败")
	}
	attachment := models.NewAttachment()
	attachment.BookId = book_id
	attachment.FileName = moreFile.Filename
	attachment.CreateAt = c.Member.MemberId
	attachment.FileExt = ext
	attachment.FilePath = filePath
	attachment.DocumentId = doc_id

	if fileInfo, err := os.Stat(filePath); err == nil {
		attachment.FileSize = float64(fileInfo.Size())
	}
	if doc_id > 0 {
		attachment.DocumentId = doc_id
	}

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") || strings.EqualFold(ext, "png") || strings.EqualFold(ext, "gif") {
		attachment.HttpPath = "/" + filePath
		is_attach = false
	}

	err = attachment.Insert()

	if err != nil {
		os.Remove(filePath)
		beego.Error("Attachment Insert => ", err)
		c.JsonResult(6006, "文件保存失败")
	}
	if attachment.HttpPath == "" {
		attachment.HttpPath = beego.URLFor("DocumentController.DownloadAttachment", ":key", identify, ":attach_id", attachment.AttachmentId)

		if err := attachment.Update(); err != nil {
			beego.Error("SaveToFile => ", err)
			c.JsonResult(6005, "保存文件失败")
		}
	}

	result := map[string]interface{}{
		"errcode":   0,
		"success":   1,
		"message":   "ok",
		"url":       attachment.HttpPath,
		"alt":       attachment.FileName,
		"is_attach": is_attach,
		"attach":    attachment,
	}

	//c.Data["json"] = result
	//c.ServeJSON(true)
	//c.StopRun()
	//
	//returnJSON, err := json.Marshal(result)
	//
	//if err != nil {
	//	beego.Error(err)
	//}
	//
	//c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	//fmt.Fprint(c.Ctx.ResponseWriter,string(returnJSON))
	c.Ctx.Output.JSON(result, true, false)
	c.StopRun()
}

//DownloadAttachment 下载附件.
func (c *DocumentController) DownloadAttachment() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	attach_id, _ := strconv.Atoi(c.Ctx.Input.Param(":attach_id"))
	token := c.GetString("token")

	member_id := 0

	if c.Member != nil {
		member_id = c.Member.MemberId
	}
	book_id := 0

	//判断用户是否参与了项目
	bookResult, err := models.NewBookResult().FindByIdentify(identify, member_id)

	if err != nil {
		//判断项目公开状态
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.Abort("404")
		}
		//如果不是超级管理员则判断权限
		if c.Member == nil || c.Member.Role != conf.MemberSuperRole {
			//如果项目是私有的，并且token不正确
			if (book.PrivatelyOwned == 1 && token == "") || (book.PrivatelyOwned == 1 && book.PrivateToken != token) {
				c.Abort("403")
			}
		}

		book_id = book.BookId
	} else {
		book_id = bookResult.BookId
	}
	//查找附件
	attachment, err := models.NewAttachment().Find(attach_id)

	if err != nil {
		beego.Error("DownloadAttachment => ", err)
		if err == orm.ErrNoRows {
			c.Abort("404")
		} else {
			c.Abort("500")
		}
	}
	if attachment.BookId != book_id {
		c.Abort("404")
	}
	c.Ctx.Output.Download(attachment.FilePath, attachment.FileName)

	c.StopRun()
}

//删除附件.
func (c *DocumentController) RemoveAttachment() {
	c.Prepare()
	attach_id, _ := c.GetInt("attach_id")

	if attach_id <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	attach, err := models.NewAttachment().Find(attach_id)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6002, "附件不存在")
	}
	document, err := models.NewDocument().Find(attach.DocumentId)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6003, "文档不存在")
	}
	if c.Member.Role != conf.MemberSuperRole {
		rel, err := models.NewRelationship().FindByBookIdAndMemberId(document.BookId, c.Member.MemberId)
		if err != nil {
			beego.Error(err)
			c.JsonResult(6004, "权限不足")
		}
		if rel.RoleId == conf.BookObserver {
			c.JsonResult(6004, "权限不足")
		}
	}
	err = attach.Delete()

	if err != nil {
		beego.Error(err)
		c.JsonResult(6005, "删除失败")
	}
	c.JsonResult(0, "ok", attach)
}

//删除文档.
func (c *DocumentController) Delete() {
	c.Prepare()

	identify := c.GetString("identify")
	doc_id, err := c.GetInt("doc_id", 0)

	book_id := 0
	//如果是超级管理员则忽略权限判断
	if c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = bookResult.BookId
	}

	if doc_id <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	doc, err := models.NewDocument().Find(doc_id)

	if err != nil {
		beego.Error("Delete => ", err)
		c.JsonResult(6003, "删除失败")
	}
	//如果文档所属项目错误
	if doc.BookId != book_id {
		c.JsonResult(6004, "参数错误")
	}
	//递归删除项目下的文档以及子文档
	err = doc.RecursiveDocument(doc.DocumentId)
	if err != nil {
		c.JsonResult(6005, "删除失败")
	}
	//重置文档数量统计
	models.NewBook().ResetDocumentNumber(doc.BookId)
	c.JsonResult(0, "ok")
}

//获取文档内容.
func (c *DocumentController) Content() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	doc_id, err := c.GetInt("doc_id")

	if err != nil {
		doc_id, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	}
	book_id := 0
	//如果是超级管理员，则忽略权限
	if c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}
		book_id = bookResult.BookId
	}

	if doc_id <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	if c.Ctx.Input.IsPost() {
		markdown := strings.TrimSpace(c.GetString("markdown", ""))
		content := c.GetString("html")
		version, _ := c.GetInt64("version", 0)
		is_cover := c.GetString("cover")

		doc, err := models.NewDocument().Find(doc_id)

		if err != nil {
			c.JsonResult(6003, "读取文档错误")
		}
		if doc.BookId != book_id {
			c.JsonResult(6004, "保存的文档不属于指定项目")
		}
		if doc.Version != version && !strings.EqualFold(is_cover, "yes") {
			beego.Info("%d|", version, doc.Version)
			c.JsonResult(6005, "文档已被修改确定要覆盖吗？")
		}
		if markdown == "" && content != "" {
			doc.Markdown = content
		} else {
			doc.Markdown = markdown
		}
		doc.Version = time.Now().Unix()
		doc.Content = content
		if err := doc.InsertOrUpdate(); err != nil {
			beego.Error("InsertOrUpdate => ", err)
			c.JsonResult(6006, "保存失败")
		}

		c.JsonResult(0, "ok", doc)
	}
	doc, err := models.NewDocument().Find(doc_id)

	if err != nil {
		c.JsonResult(6003, "文档不存在")
	}
	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}
	c.JsonResult(0, "ok", doc)
}

//导出文件
func (c *DocumentController) Export() {
	c.Prepare()
	c.TplName = "document/export.tpl"

	identify := c.Ctx.Input.Param(":key")

	output := c.GetString("output")

	token := c.GetString("token")

	if identify == "" {
		c.Abort("404")
	}
	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}
	bookResult := models.NewBookResult()
	if c.Member != nil && c.Member.Role == conf.MemberSuperRole {
		book, err := models.NewBook().FindByIdentify(identify)
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
		bookResult = book.ToBookResult()
	} else {
		bookResult = isReadable(identify, token, c)
	}
	docs, err := models.NewDocument().FindListByBookId(bookResult.BookId)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}

	if output == "pdf" {

		exe := beego.AppConfig.String("wkhtmltopdf")

		if exe == "" {
			c.TplName = "errors/error.tpl"
			c.Data["ErrorMessage"] = "没有配置PDF导出程序"
			c.Data["ErrorCode"] = 50010
			return
		}
		dpath := "cache/" + bookResult.Identify

		os.MkdirAll(dpath, 0766)

		pathList := list.New()

		RecursiveFun(0, "", dpath, c, bookResult, docs, pathList)

		defer os.RemoveAll(dpath)

		os.MkdirAll("./cache", 0766)
		pdfpath := "cache/" + identify + "_" + c.CruSession.SessionID() + ".pdf"

		if _, err := os.Stat(pdfpath); os.IsNotExist(err) {

			wkhtmltopdf.SetPath(beego.AppConfig.String("wkhtmltopdf"))
			pdfg, err := wkhtmltopdf.NewPDFGenerator()
			pdfg.MarginBottom.Set(35)

			if err != nil {
				beego.Error(err)
				c.Abort("500")
			}

			for e := pathList.Front(); e != nil; e = e.Next() {
				if page, ok := e.Value.(string); ok {
					pdfg.AddPage(wkhtmltopdf.NewPage(page))
				}
			}
			err = pdfg.Create()
			if err != nil {
				beego.Error(err)
				c.Abort("500")
			}

			err = pdfg.WriteFile(pdfpath)
			if err != nil {
				beego.Error(err)
			}
		}

		c.Ctx.Output.Download(pdfpath, identify+".pdf")

		defer os.Remove(pdfpath)

		c.StopRun()
	}

	c.Abort("404")
}

//生成项目访问的二维码.
func (c *DocumentController) QrCode() {
	c.Prepare()
	identify := c.GetString(":key")

	book, err := models.NewBook().FindByIdentify(identify)

	if err != nil || book.BookId <= 0 {
		c.Abort("404")
	}

	uri := c.BaseUrl() + beego.URLFor("DocumentController.Index", ":key", identify)
	code, err := qr.Encode(uri, qr.L, qr.Unicode)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	code, err = barcode.Scale(code, 150, 150)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")

	//imgpath := filepath.Join("cache","qrcode",identify + ".png")

	err = png.Encode(c.Ctx.ResponseWriter, code)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
}

func (c *DocumentController) Search()  {
	c.Prepare()
	
	c.JsonResult(0,"ok")
}

//递归生成文档序列数组.
func RecursiveFun(parent_id int, prefix, dpath string, c *DocumentController, book *models.BookResult, docs []*models.Document, paths *list.List) {
	for _, item := range docs {
		if item.ParentId == parent_id {
			name := prefix + strconv.Itoa(item.ParentId) + strconv.Itoa(item.OrderSort) + strconv.Itoa(item.DocumentId)
			fpath := dpath + "/" + name + ".html"
			paths.PushBack(fpath)

			f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0777)

			if err != nil {
				beego.Error(err)
				c.Abort("500")
			}

			html, err := c.ExecuteViewPathTemplate("document/export.tpl", map[string]interface{}{"Model": book, "Lists": item, "BaseUrl": c.BaseUrl()})
			if err != nil {
				f.Close()
				beego.Error(err)
				c.Abort("500")
			}
			//beego.Info(fpath,html)
			f.WriteString(html)
			f.Close()

			for _, sub := range docs {
				if sub.ParentId == item.DocumentId {
					RecursiveFun(item.DocumentId, name, dpath, c, book, docs, paths)
					break
				}
			}
		}
	}
}
