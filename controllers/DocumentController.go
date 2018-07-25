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
	"net/url"
	"image/png"

	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
	"github.com/lifei6671/mindoc/utils/pagination"
	"gopkg.in/russross/blackfriday.v2"
	"github.com/lifei6671/mindoc/utils/cryptil"
	"fmt"
	"github.com/lifei6671/mindoc/utils/filetil"
	"github.com/lifei6671/mindoc/utils/gopool"
)

// DocumentController struct
type DocumentController struct {
	BaseController
}

// 文档首页
func (c *DocumentController) Index() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")

	if identify == "" {
		c.ShowErrorPage(404,"项目不存在或已删除")
	}

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !isUserLoggedIn(c) {
		promptUserToLogIn(c)
		return
	}

	bookResult := isReadable(identify, token, c)

	c.TplName = "document/" + bookResult.Theme + "_read.tpl"

	selected := 0

	if bookResult.IsUseFirstDocument {
		doc,err := bookResult.FindFirstDocumentByBookId(bookResult.BookId)
		if err == nil {
			if strings.TrimSpace(doc.Release) != "" {
				doc.Release += "<div class=\"wiki-bottom\">文档更新时间: " + doc.ModifyTime.Local().Format("2006-01-02 15:04") + "</div>";
			}
			selected = doc.DocumentId
			c.Data["Title"] = doc.DocumentName
			c.Data["Content"] = template.HTML(doc.Release)

			c.Data["Description"] = utils.AutoSummary(doc.Release,120)
		}
	}else {
		c.Data["Title"] = "概要"
		c.Data["Content"] = template.HTML(blackfriday.Run([]byte(bookResult.Description)))
	}

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, selected)

	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404,"生成项目文档树时出错")
		} else {
			beego.Error("生成项目文档树时出错 -> ",err)
			c.ShowErrorPage(500,"生成项目文档树时出错")
		}
	}
	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)

}

// 阅读文档
func (c *DocumentController) Read() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	id := c.GetString(":id")

	c.Data["DocumentId"] = id

	if identify == "" || id == "" {
		c.ShowErrorPage(404,"项目不存或已删除")
	}

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !isUserLoggedIn(c) {
		promptUserToLogIn(c)
		return
	}

	bookResult := isReadable(identify, token, c)

	c.TplName = fmt.Sprintf("document/%s_read.tpl",bookResult.Theme)

	doc := models.NewDocument()

	if docId, err := strconv.Atoi(id); err == nil {
		doc, err = doc.FromCacheById(docId)
		if err != nil {
			beego.Error("从缓存中读取文档时失败 ->",err)
			c.ShowErrorPage(500,"文档不存在或已删除")
		}
	} else {
		doc, err = doc.FromCacheByIdentify(id,bookResult.BookId)
		if err != nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(404,"文档不存在或已删除")
			}else{
				beego.Error("从缓存查询文档时出错 ->" ,err)
				c.ShowErrorPage(500,"未知异常")
			}
		}
	}

	if doc.BookId != bookResult.BookId {
		c.ShowErrorPage(404,"文档不存在或已删除")
	}

	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	cdnimg := beego.AppConfig.String("cdnimg")
	if doc.Release != "" && cdnimg != "" {
		query, err := goquery.NewDocumentFromReader(bytes.NewBufferString(doc.Release))
		if err != nil {
			beego.Error(err)
		} else {
			query.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
				if src, ok := contentSelection.Attr("src"); ok && strings.HasPrefix(src, "/uploads/") {
					contentSelection.SetAttr("src", utils.JoinURI(cdnimg, src))
				}
			})

			html, err := query.Html()
			if err != nil {
				beego.Error(err)
			} else {
				doc.Release = html
			}
		}
	}

	// assemble doc info, added by dandycheung, 2017-12-20
	docInfo := ""
	docCreator, err := models.NewMember().Find(doc.MemberId)
	if err == nil {
		docInfo += docCreator.Account
	}

	docInfo += " 创建于 "
	docInfo += doc.CreateTime.Local().Format("2006-01-02 15:04")

	if doc.ModifyTime != doc.CreateTime {
		docInfo += "；更新于 "
		docInfo += doc.ModifyTime.Local().Format("2006-01-02 15:04")
		if strings.TrimSpace(doc.Release) != "" {
			doc.Release += "<div class=\"wiki-bottom\">文档更新时间: " + doc.ModifyTime.Local().Format("2006-01-02 15:04") +" &nbsp;&nbsp;作者：";
			if docCreator != nil {
				if docCreator.RealName != "" {
					doc.Release += docCreator.RealName
				}else{
					doc.Release += docCreator.Account
				}
			}
			doc.Release += "</div>"
		}
	}



	if c.IsAjax() {
		var data struct {
			DocTitle string `json:"doc_title"`
			Body     string `json:"body"`
			Title    string `json:"title"`
			DocInfo  string `json:"doc_info"`
		}
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release
		data.Title = doc.DocumentName + " - Powered by MinDoc"
		data.DocInfo = docInfo

		c.JsonResult(0, "ok", data)
	}

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, doc.DocumentId)

	if err != nil && err != orm.ErrNoRows {
		beego.Error("生成项目文档树时出错 ->",err)

		c.ShowErrorPage(500,"生成项目文档树时出错")
	}


	c.Data["Description"] =  utils.AutoSummary(doc.Release,120)


	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = doc.DocumentName
	c.Data["Info"] = docInfo
	c.Data["Content"] = template.HTML(doc.Release)
}

// 编辑文档
func (c *DocumentController) Edit() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.ShowErrorPage(404,"无法解析项目标识")
	}

	bookResult := models.NewBookResult()

	var err error
	// 如果是超级管理者，则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookResult = models.NewBookResult().ToBookResult(*book)
	} else {
		bookResult, err = models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(403, "项目不存在或没有权限")
			}else{
				beego.Error("查询项目时出错 -> ", err)
				c.ShowErrorPage(500,"查询项目时出错")
			}
		}
		if bookResult.RoleId == conf.BookObserver {
			c.JsonResult(6002, "项目不存在或权限不足")
		}
	}

	// 根据不同编辑器类型加载编辑器
	if bookResult.Editor == "markdown" {
		c.TplName = "document/markdown_edit_template.tpl"
	} else if bookResult.Editor == "html" {
		c.TplName = "document/new_html_edit_template.tpl"
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

// 创建一个文档
func (c *DocumentController) Create() {
	identify := c.GetString("identify")
	docIdentify := c.GetString("doc_identify")
	docName := c.GetString("doc_name")
	parentId, _ := c.GetInt("parent_id", 0)
	docId, _ := c.GetInt("doc_id", 0)

	if identify == "" {
		c.JsonResult(6001, "参数错误")
	}

	if docName == "" {
		c.JsonResult(6004, "文档名称不能为空")
	}


	bookId := 0

	// 如果是超级管理员则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error(err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = bookResult.BookId
	}

	if docIdentify != "" {
		if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, docIdentify); !ok || err != nil {
			c.JsonResult(6003, "文档标识只能包含小写字母、数字，以及“-”、“.”和“_”符号")
		}

		d, _ := models.NewDocument().FindByIdentityFirst(docIdentify,bookId)
		if d.DocumentId > 0 && d.DocumentId != docId {
			c.JsonResult(6006, "文档标识已被使用")
		}
	}
	if parentId > 0 {
		doc, err := models.NewDocument().Find(parentId)
		if err != nil || doc.BookId != bookId {
			c.JsonResult(6003, "父分类不存在")
		}
	}

	document, _ := models.NewDocument().Find(docId)

	document.MemberId = c.Member.MemberId
	document.BookId = bookId

	document.Identify = docIdentify

	document.Version = time.Now().Unix()
	document.DocumentName = docName
	document.ParentId = parentId

	if err := document.InsertOrUpdate(); err != nil {
		beego.Error("InsertOrUpdate => ", err)
		c.JsonResult(6005, "保存失败")
	} else {
		c.JsonResult(0, "ok", document)
	}
}

// 上传附件或图片
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

	type Size interface {
		Size() int64
	}

	if conf.GetUploadFileSize() > 0 && moreFile.Size > conf.GetUploadFileSize() {
		c.JsonResult(6009, "查过文件允许的上传最大值")
	}

	ext := filepath.Ext(moreFile.Filename)

	if ext == "" {
		c.JsonResult(6003, "无法解析文件的格式")
	}
	//如果文件类型设置为 * 标识不限制文件类型
	if beego.AppConfig.DefaultString("upload_file_ext", "") != "*" {
		if !conf.IsAllowUploadFileExt(ext) {
			c.JsonResult(6004, "不允许的文件类型")
		}
	}

	bookId := 0

	// 如果是超级管理员，则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)

		if err != nil {
			c.JsonResult(6006, "文档不存在或权限不足")
		}

		bookId = book.BookId
	} else {
		book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			beego.Error("DocumentController.Edit => ", err)
			if err == orm.ErrNoRows {
				c.JsonResult(6006, "权限不足")
			}

			c.JsonResult(6001, err.Error())
		}

		// 如果没有编辑权限
		if book.RoleId != conf.BookEditor && book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
			c.JsonResult(6006, "权限不足")
		}

		bookId = book.BookId
	}

	if doc_id > 0 {
		doc, err := models.NewDocument().Find(doc_id)
		if err != nil {
			c.JsonResult(6007, "文档不存在")
		}

		if doc.BookId != bookId {
			c.JsonResult(6008, "文档不属于指定的项目")
		}
	}

	fileName := "attach_" + strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"),identify, fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile(name, filePath)

	if err != nil {
		beego.Error("SaveToFile => ", err)
		c.JsonResult(6005, "保存文件失败")
	}

	attachment := models.NewAttachment()
	attachment.BookId = bookId
	attachment.FileName = moreFile.Filename
	attachment.CreateAt = c.Member.MemberId
	attachment.FileExt = ext
	attachment.FilePath = strings.TrimPrefix(filePath, conf.WorkingDirectory)
	attachment.DocumentId = doc_id

	if fileInfo, err := os.Stat(filePath); err == nil {
		attachment.FileSize = float64(fileInfo.Size())
	}

	if doc_id > 0 {
		attachment.DocumentId = doc_id
	}

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") || strings.EqualFold(ext, ".png") || strings.EqualFold(ext, ".gif") {
		attachment.HttpPath = "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)
		if strings.HasPrefix(attachment.HttpPath, "//") {
			attachment.HttpPath = conf.URLForWithCdnImage(string(attachment.HttpPath[1:]))
		}

		is_attach = false
	}

	err = attachment.Insert()

	if err != nil {
		os.Remove(filePath)
		beego.Error("Attachment Insert => ", err)
		c.JsonResult(6006, "文件保存失败")
	}

	if attachment.HttpPath == "" {
		attachment.HttpPath = conf.URLFor("DocumentController.DownloadAttachment", ":key", identify, ":attach_id", attachment.AttachmentId)

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

	c.Ctx.Output.JSON(result, true, false)
	c.StopRun()
}

// 下载附件
func (c *DocumentController) DownloadAttachment() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	attachId, _ := strconv.Atoi(c.Ctx.Input.Param(":attach_id"))
	token := c.GetString("token")

	memberId := 0

	if c.Member != nil {
		memberId = c.Member.MemberId
	}

	bookId := 0

	// 判断用户是否参与了项目
	bookResult, err := models.NewBookResult().FindByIdentify(identify, memberId)

	if err != nil {
		// 判断项目公开状态
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(404,"项目不存在或已删除")
			}else{
				beego.Error("查找项目时出错 ->",err)
				c.ShowErrorPage(500,"系统错误")
			}
		}

		// 如果不是超级管理员则判断权限
		if c.Member == nil || c.Member.Role != conf.MemberSuperRole {
			// 如果项目是私有的，并且 token 不正确
			if (book.PrivatelyOwned == 1 && token == "") || (book.PrivatelyOwned == 1 && book.PrivateToken != token) {
				c.ShowErrorPage(403,"权限不足")
			}
		}

		bookId = book.BookId
	} else {
		bookId = bookResult.BookId
	}

	// 查找附件
	attachment, err := models.NewAttachment().Find(attachId)

	if err != nil {
		beego.Error("查找附件时出错 -> ", err)
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404,"附件不存在或已删除")
		} else {
			c.ShowErrorPage(500,"查找附件时出错")
		}
	}

	if attachment.BookId != bookId {
		c.ShowErrorPage(404,"附件不存在或已删除")
	}

	c.Ctx.Output.Download(filepath.Join(conf.WorkingDirectory, attachment.FilePath), attachment.FileName)
	c.StopRun()
}

// 删除附件
func (c *DocumentController) RemoveAttachment() {
	c.Prepare()
	attachId, _ := c.GetInt("attach_id")

	if attachId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	attach, err := models.NewAttachment().Find(attachId)

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

	os.Remove(filepath.Join(conf.WorkingDirectory, attach.FilePath))

	c.JsonResult(0, "ok", attach)
}

// 删除文档
func (c *DocumentController) Delete() {
	c.Prepare()

	identify := c.GetString("identify")
	doc_id, err := c.GetInt("doc_id", 0)

	book_id := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
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
	// 如果文档所属项目错误
	if doc.BookId != book_id {
		c.JsonResult(6004, "参数错误")
	}

	// 递归删除项目下的文档以及子文档
	err = doc.RecursiveDocument(doc.DocumentId)
	if err != nil {
		c.JsonResult(6005, "删除失败")
	}

	// 重置文档数量统计
	models.NewBook().ResetDocumentNumber(doc.BookId)
	c.JsonResult(0, "ok")
}

// 获取文档内容
func (c *DocumentController) Content() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	docId, err := c.GetInt("doc_id")

	if err != nil {
		docId, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	}

	bookId := 0
	autoRelease := false

	// 如果是超级管理员，则忽略权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = book.BookId
		autoRelease = book.AutoRelease == 1
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = bookResult.BookId
		autoRelease = bookResult.AutoRelease
	}

	if docId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	if c.Ctx.Input.IsPost() {
		markdown := strings.TrimSpace(c.GetString("markdown", ""))
		content := c.GetString("html")
		version, _ := c.GetInt64("version", 0)
		isCover := c.GetString("cover")

		doc, err := models.NewDocument().Find(docId)

		if err != nil {
			c.JsonResult(6003, "读取文档错误")
		}

		if doc.BookId != bookId {
			c.JsonResult(6004, "保存的文档不属于指定项目")
		}

		if doc.Version != version && !strings.EqualFold(isCover, "yes") {
			beego.Info("%d|", version, doc.Version)
			c.JsonResult(6005, "文档已被修改确定要覆盖吗？")
		}

		history := models.NewDocumentHistory()
		history.DocumentId = docId
		history.Content = doc.Content
		history.Markdown = doc.Markdown
		history.DocumentName = doc.DocumentName
		history.ModifyAt = c.Member.MemberId
		history.MemberId = doc.MemberId
		history.ParentId = doc.ParentId
		history.Version = time.Now().Unix()
		history.Action = "modify"
		history.ActionName = "修改文档"

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

		// 如果启用了文档历史，则添加历史文档
		///如果两次保存的MD5值不同则保存为历史，否则忽略
		go func(history *models.DocumentHistory) {
			if c.EnableDocumentHistory && cryptil.Md5Crypt(history.Markdown) != cryptil.Md5Crypt(doc.Markdown) {
				_, err = history.InsertOrUpdate()
				if err != nil {
					beego.Error("DocumentHistory InsertOrUpdate => ", err)
				}
			}
		}(history)

		//如果启用了自动发布
		if autoRelease {
			go func(identify string) {
				models.NewBook().ReleaseContent(bookId)

			}(identify)
		}

		c.JsonResult(0, "ok", doc)
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		c.JsonResult(6003, "文档不存在")
	}

	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	c.JsonResult(0, "ok", doc)
}
//
//func (c *DocumentController) GetDocumentById(id string) (doc *models.Document, err error) {
//	doc = models.NewDocument()
//	if doc_id, err := strconv.Atoi(id); err == nil {
//		doc, err = doc.Find(doc_id)
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		doc, err = doc.FindByFieldFirst("identify", id)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return doc, nil
//}

// 导出
func (c *DocumentController) Export() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.ShowErrorPage(500,"参数错误")
	}

	output := c.GetString("output")
	token := c.GetString("token")

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !isUserLoggedIn(c) {
		promptUserToLogIn(c)
		return
	}
	if !conf.GetEnableExport() {
		c.ShowErrorPage(500,"系统没有开启导出功能")
	}

	bookResult := models.NewBookResult()
	if c.Member != nil && c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByIdentify(identify)
		if err != nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(404,"项目不存在")
			} else {
				beego.Error("查找项目时出错 ->",err)
				c.ShowErrorPage(500,"查找项目时出错")
			}
		}
		bookResult = models.NewBookResult().ToBookResult(*book)
	} else {
		bookResult = isReadable(identify, token, c)
	}
	if !bookResult.IsDownload {
		c.ShowErrorPage(200,"当前项目没有开启导出功能")
	}

	if !strings.HasPrefix(bookResult.Cover, "http:://") && !strings.HasPrefix(bookResult.Cover, "https:://") {
		bookResult.Cover = conf.URLForWithCdnImage(bookResult.Cover)
	}

	if output == "markdown" {
		if bookResult.Editor != "markdown"{
			c.ShowErrorPage(500,"当前项目不支持Markdown编辑器")
		}
		p,err := bookResult.ExportMarkdown(c.CruSession.SessionID())

		if err != nil {
			c.ShowErrorPage(500,"导出文档失败")
		}
		c.Ctx.Output.Download(p, bookResult.BookName+".zip")

		c.StopRun()
		return
	}

	outputPath := filepath.Join(conf.GetExportOutputPath(), strconv.Itoa(bookResult.BookId))

	pdfpath := filepath.Join(outputPath, "book.pdf")
	epubpath := filepath.Join(outputPath, "book.epub")
	mobipath := filepath.Join(outputPath, "book.mobi")
	docxpath := filepath.Join(outputPath, "book.docx")

	if output == "pdf" && filetil.FileExists(pdfpath){
		c.Ctx.Output.Download(pdfpath, bookResult.BookName+".pdf")
		c.Abort("200")
	} else if output == "epub"  && filetil.FileExists(epubpath){
		c.Ctx.Output.Download(epubpath, bookResult.BookName+".epub")

		c.Abort("200")
	} else if output == "mobi" && filetil.FileExists(mobipath) {
		c.Ctx.Output.Download(mobipath, bookResult.BookName+".mobi")

		c.Abort("200")
	} else if output == "docx"  && filetil.FileExists(docxpath){
		c.Ctx.Output.Download(docxpath, bookResult.BookName+".docx")

		c.Abort("200")

	}else if output == "pdf" || output == "epub" || output == "docx" || output == "mobi"{
		if err := models.BackgroupConvert(c.CruSession.SessionID(),bookResult);err != nil && err != gopool.ErrHandlerIsExist{
			c.ShowErrorPage(500,"导出失败，请查看系统日志")
		}

		c.ShowErrorPage(200,"文档正在后台转换，请稍后再下载")
	}else{
		c.ShowErrorPage(200,"不支持的文件格式")
	}

	c.ShowErrorPage(404,"项目没有导出文件")
}

// 生成项目访问的二维码
func (c *DocumentController) QrCode() {
	c.Prepare()

	identify := c.GetString(":key")

	book, err := models.NewBook().FindByIdentify(identify)
	if err != nil || book.BookId <= 0 {
		c.ShowErrorPage(404,"项目不存在")
	}

	uri := conf.URLFor("DocumentController.Index", ":key", identify)
	code, err := qr.Encode(uri, qr.L, qr.Unicode)
	if err != nil {
		beego.Error("生成二维码失败 ->",err)
		c.ShowErrorPage(500,"生成二维码失败")
	}

	code, err = barcode.Scale(code, 150, 150)
	if err != nil {
		beego.Error("生成二维码失败 ->",err)
		c.ShowErrorPage(500,"生成二维码失败")
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")

	// imgpath := filepath.Join("cache","qrcode",identify + ".png")

	err = png.Encode(c.Ctx.ResponseWriter, code)
	if err != nil {
		beego.Error("生成二维码失败 ->",err)
		c.ShowErrorPage(500,"生成二维码失败")
	}
}

// 项目内搜索
func (c *DocumentController) Search() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	keyword := strings.TrimSpace(c.GetString("keyword"))

	if identify == "" {
		c.JsonResult(6001, "参数错误")
	}

	if !c.EnableAnonymous && !isUserLoggedIn(c) {
		promptUserToLogIn(c)
		return
	}

	bookResult := isReadable(identify, token, c)

	docs, err := models.NewDocumentSearchResult().SearchDocument(keyword, bookResult.BookId)
	if err != nil {
		beego.Error(err)
		c.JsonResult(6002, "搜索结果错误")
	}

	if len(docs) < 0 {
		c.JsonResult(404, "没有数据库")
	}

	for _, doc := range docs {
		doc.BookId = bookResult.BookId
		doc.BookName = bookResult.BookName
		doc.Description = bookResult.Description
		doc.BookIdentify = bookResult.Identify
	}

	c.JsonResult(0, "ok", docs)
}

// 文档历史列表
func (c *DocumentController) History() {
	c.Prepare()

	c.TplName = "document/history.tpl"

	identify := c.GetString("identify")
	docId, err := c.GetInt("doc_id", 0)
	pageIndex, _ := c.GetInt("page", 1)

	book_id := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error("FindByIdentify => ", err)
			c.Data["ErrorMessage"] = "项目不存在或权限不足"
			return
		}

		book_id = book.BookId
		c.Data["Model"] = book
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.Data["ErrorMessage"] = "项目不存在或权限不足"
			return
		}

		book_id = bookResult.BookId
		c.Data["Model"] = bookResult
	}

	if docId <= 0 {
		c.Data["ErrorMessage"] = "参数错误"
		return
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		beego.Error("Delete => ", err)
		c.Data["ErrorMessage"] = "获取历史失败"
		return
	}

	// 如果文档所属项目错误
	if doc.BookId != book_id {
		c.Data["ErrorMessage"] = "参数错误"
		return
	}

	historis, totalCount, err := models.NewDocumentHistory().FindToPager(docId, pageIndex, conf.PageSize)
	if err != nil {
		beego.Error("FindToPager => ", err)
		c.Data["ErrorMessage"] = "获取历史失败"
		return
	}

	c.Data["List"] = historis
	c.Data["PageHtml"] = ""
	c.Data["Document"] = doc

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	}
}

func (c *DocumentController) DeleteHistory() {
	c.Prepare()

	c.TplName = "document/history.tpl"

	identify := c.GetString("identify")
	docId, err := c.GetInt("doc_id", 0)
	historyId, _ := c.GetInt("history_id", 0)

	if historyId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	bookId := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = bookResult.BookId
	}

	if docId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		beego.Error("Delete => ", err)
		c.JsonResult(6001, "获取历史失败")
	}

	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.JsonResult(6001, "参数错误")
	}

	err = models.NewDocumentHistory().Delete(historyId, docId)
	if err != nil {
		beego.Error(err)
		c.JsonResult(6002, "删除失败")
	}

	c.JsonResult(0, "ok")
}
//通过文档历史恢复文档
func (c *DocumentController) RestoreHistory() {
	c.Prepare()

	c.TplName = "document/history.tpl"

	identify := c.GetString("identify")
	docId, err := c.GetInt("doc_id", 0)
	historyId, _ := c.GetInt("history_id", 0)

	if historyId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	bookId := 0
	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.JsonResult(6002, "项目不存在或权限不足")
		}

		bookId = bookResult.BookId
	}

	if docId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		beego.Error("Delete => ", err)
		c.JsonResult(6001, "获取历史失败")
	}

	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.JsonResult(6001, "参数错误")
	}

	err = models.NewDocumentHistory().Restore(historyId, docId, c.Member.MemberId)
	if err != nil {
		beego.Error(err)
		c.JsonResult(6002, "删除失败")
	}

	c.JsonResult(0, "ok", doc)
}

func (c *DocumentController) Compare() {
	c.Prepare()

	c.TplName = "document/compare.tpl"

	historyId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	identify := c.Ctx.Input.Param(":key")

	bookId := 0
	editor := "markdown"

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			beego.Error("DocumentController.Compare => ", err)
			c.ShowErrorPage(403,"权限不足")
			return
		}

		bookId = book.BookId
		c.Data["Model"] = book
		editor = book.Editor
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			beego.Error("FindByIdentify => ", err)
			c.ShowErrorPage(403,"权限不足")
			return
		}

		bookId = bookResult.BookId
		c.Data["Model"] = bookResult
		editor = bookResult.Editor
	}

	if historyId <= 0 {
		c.ShowErrorPage(60002, "参数错误")
	}

	history, err := models.NewDocumentHistory().Find(historyId)
	if err != nil {
		beego.Error("DocumentController.Compare => ", err)
		c.ShowErrorPage(60003, err.Error())
	}

	doc, err := models.NewDocument().Find(history.DocumentId)
	if doc.BookId != bookId {
		c.ShowErrorPage(60002, "参数错误")
	}

	c.Data["HistoryId"] = historyId
	c.Data["DocumentId"] = doc.DocumentId

	if editor == "markdown" {
		c.Data["HistoryContent"] = history.Markdown
		c.Data["Content"] = doc.Markdown
	} else {
		c.Data["HistoryContent"] = template.HTML(history.Content)
		c.Data["Content"] = template.HTML(doc.Content)
	}
}

// 递归生成文档序列数组
func RecursiveFun(parentId int, prefix, dpath string, c *DocumentController, book *models.BookResult, docs []*models.Document, paths *list.List) {
	for _, item := range docs {
		if item.ParentId == parentId {
			EachFun(prefix, dpath, c, book, item, paths)

			for _, sub := range docs {
				if sub.ParentId == item.DocumentId {
					prefix += strconv.Itoa(item.ParentId) + strconv.Itoa(item.OrderSort) + strconv.Itoa(item.DocumentId)
					RecursiveFun(item.DocumentId, prefix, dpath, c, book, docs, paths)
					break
				}
			}
		}
	}
}

func EachFun(prefix, dpath string, c *DocumentController, book *models.BookResult, item *models.Document, paths *list.List) {
	name := prefix + strconv.Itoa(item.ParentId) + strconv.Itoa(item.OrderSort) + strconv.Itoa(item.DocumentId)
	fpath := dpath + "/" + name + ".html"
	paths.PushBack(fpath)

	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		beego.Error(err)
		c.ShowErrorPage(500,"系统错误")
	}

	html, err := c.ExecuteViewPathTemplate("document/export.tpl", map[string]interface{}{"Model": book, "Lists": item, "BaseUrl": c.BaseUrl()})
	if err != nil {
		f.Close()
		beego.Error(err)
		c.ShowErrorPage(500,"系统错误")
	}

	buf := bytes.NewReader([]byte(html))
	doc, err := goquery.NewDocumentFromReader(buf)
	doc.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
		if src, ok := contentSelection.Attr("src"); ok && strings.HasPrefix(src, "/uploads/") {
			contentSelection.SetAttr("src", c.BaseUrl()+src)
		}
	})

	html, err = doc.Html()
	if err != nil {
		f.Close()
		beego.Error(err)
		c.ShowErrorPage(500,"系统错误")
	}

	// html = strings.Replace(html, "<img src=\"/uploads", "<img src=\"" + c.BaseUrl() + "/uploads", -1)

	f.WriteString(html)
	f.Close()
}


// 判断用户是否可以阅读文档
func isReadable(identify, token string, c *DocumentController) *models.BookResult {
	book, err := models.NewBook().FindByFieldFirst("identify", identify)

	if err != nil {
		beego.Error(err)
		c.ShowErrorPage(500,"项目不存在")
	}

	// 如果文档是私有的
	if book.PrivatelyOwned == 1 && !c.Member.IsAdministrator() {
		is_ok := false

		if c.Member != nil {
			_, err := models.NewRelationship().FindForRoleId(book.BookId, c.Member.MemberId)
			if err == nil {
				is_ok = true
			}
		}

		if book.PrivateToken != "" && !is_ok {
			// 如果有访问的 Token，并且该项目设置了访问 Token，并且和用户提供的相匹配，则记录到 Session 中。
			// 如果用户未提供 Token 且用户登录了，则判断用户是否参与了该项目。
			// 如果用户未登录，则从 Session 中读取 Token。
			if token != "" && strings.EqualFold(token, book.PrivateToken) {
				c.SetSession(identify, token)
			} else if token, ok := c.GetSession(identify).(string); !ok || !strings.EqualFold(token, book.PrivateToken) {
				c.ShowErrorPage(403,"权限不足")
			}
		} else if !is_ok {
			c.ShowErrorPage(403,"权限不足")
		}
	}

	bookResult := models.NewBookResult().ToBookResult(*book)

	if c.Member != nil {
		rel, err := models.NewRelationship().FindByBookIdAndMemberId(bookResult.BookId, c.Member.MemberId)

		if err == nil {
			bookResult.MemberId = rel.MemberId
			bookResult.RoleId = rel.RoleId
			bookResult.RelationshipId = rel.RelationshipId
		}
	}

	// 判断是否需要显示评论框
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

func isUserLoggedIn(c *DocumentController) bool {
	return c.Member != nil && c.Member.MemberId > 0
}

func promptUserToLogIn(c *DocumentController) {
	beego.Info("Access " + c.Ctx.Request.URL.RequestURI() + " not permitted.")
	beego.Info("  Access will be redirected to login page(SessionId: " + c.CruSession.SessionID() + ").")

	if c.IsAjax() {
		c.JsonResult(6000, "请重新登录。")
	} else {
		c.Redirect(conf.URLFor("AccountController.Login")+ "?url=" + url.PathEscape(conf.BaseUrl+ c.Ctx.Request.URL.RequestURI()), 302)
	}
}
