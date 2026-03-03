package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/cryptil"
	"github.com/mindoc-org/mindoc/utils/filetil"
	"github.com/mindoc-org/mindoc/utils/gopool"
	"github.com/mindoc-org/mindoc/utils/pagination"
	"github.com/russross/blackfriday/v2"
)

// DocumentController struct
type DocumentController struct {
	BaseController
}

// Document prev&next
type DocumentTreeFlatten struct {
	DocumentId   int    `json:"id"`
	DocumentName string `json:"text"`
	// ParentId     interface{} `json:"parent"`
	Identify string `json:"identify"`
	// BookIdentify string      `json:"-"`
	// Version      int64       `json:"version"`
}

// 文档首页
func (c *DocumentController) Index() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")

	if identify == "" {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !c.isUserLoggedIn() {
		promptUserToLogIn(c)
		return
	}

	bookResult := c.isReadable(identify, token)

	c.TplName = "document/" + bookResult.Theme + "_read.tpl"

	selected := 0

	if bookResult.IsUseFirstDocument {
		doc, err := bookResult.FindFirstDocumentByBookId(bookResult.BookId)
		if err == nil {
			selected = doc.DocumentId
			c.Data["Title"] = doc.DocumentName
			c.Data["Content"] = template.HTML(doc.Release)
			c.Data["Description"] = utils.AutoSummary(doc.Release, 120)
			c.Data["FoldSetting"] = "first"

			if bookResult.Editor == EditorCherryMarkdown {
				c.Data["MarkdownTheme"] = doc.MarkdownTheme
			}

			if bookResult.IsDisplayComment {
				// 获取评论、分页
				comments, count, _ := models.NewComment().QueryCommentByDocumentId(doc.DocumentId, 1, conf.PageSize, c.Member)
				page := pagination.PageUtil(int(count), 1, conf.PageSize, comments)
				c.Data["Page"] = page
			}
		}
	} else {
		c.Data["Title"] = i18n.Tr(c.Lang, "blog.summary")
		c.Data["Content"] = template.HTML(blackfriday.Run([]byte(bookResult.Description)))
		c.Data["FoldSetting"] = "closed"
	}

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, selected)

	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.no_doc_in_cur_proj"))
		} else {
			logs.Error("生成项目文档树时出错 -> ", err)
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.build_doc_tree_error"))
		}
	}
	c.Data["IS_DOCUMENT_INDEX"] = true
	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)
}

// CheckPassword : Handles password verification for private documents,
// and front-end requests are made through Ajax.
func (c *DocumentController) CheckPassword() {
	identify := c.Ctx.Input.Param(":key")
	password := c.GetString("bPassword")

	if identify == "" || password == "" {
		c.JsonResult(http.StatusBadRequest, i18n.Tr(c.Lang, "message.param_error"))
	}

	// You have not logged in and need to log in again.
	if !c.EnableAnonymous && !c.isUserLoggedIn() {
		logs.Info("You have not logged in and need to log in again(SessionId: %s).",
			c.CruSession.SessionID(context.TODO()))
		c.JsonResult(6000, i18n.Tr(c.Lang, "message.need_relogin"))
		return
	}

	book, err := models.NewBook().FindByFieldFirst("identify", identify)

	if err != nil {
		logs.Error(err)
		c.JsonResult(500, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	if book.BookPassword != password {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.wrong_password"))
	} else {
		c.SetSession(identify, password)
		c.JsonResult(0, "OK")
	}
}

// 阅读文档
func (c *DocumentController) Read() {
	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	id := c.GetString(":id")

	if identify == "" || id == "" {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !c.isUserLoggedIn() {
		promptUserToLogIn(c)
		return
	}

	bookResult := c.isReadable(identify, token)

	c.TplName = fmt.Sprintf("document/%s_read.tpl", bookResult.Theme)

	doc := models.NewDocument()
	if docId, err := strconv.Atoi(id); err == nil {
		doc, err = doc.FromCacheById(docId)
		if err != nil || doc == nil {
			logs.Error("从缓存中读取文档时失败 ->", err)
			c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.doc_not_exist"))
			return
		}
	} else {
		doc, err = doc.FromCacheByIdentify(id, bookResult.BookId)
		if err != nil || doc == nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.doc_not_exist"))
			} else {
				logs.Error("从数据库查询文档时出错 ->", err)
				c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.unknown_exception"))
			}
			return
		}
	}

	if doc.BookId != bookResult.BookId {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.doc_not_exist"))
	}
	doc.Lang = c.Lang
	doc.Processor()

	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	// prev,next
	treeJson, err := models.NewDocument().FindDocumentTree2(bookResult.BookId)
	if err != nil {
		logs.Error("生成项目文档树时出错 ->", err)
	}

	res := getTreeRecursive(treeJson, 0)
	flat := make([]DocumentTreeFlatten, 0)
	Flatten(res, &flat)
	var index int
	for i, v := range flat {
		if v.Identify == id {
			index = i
		}
	}
	var PrevName, PrevPath, NextName, NextPath string
	if index == 0 {
		c.Data["PrevName"] = ""
		PrevName = ""
	} else {
		c.Data["PrevPath"] = identify + "/" + flat[index-1].Identify
		c.Data["PrevName"] = flat[index-1].DocumentName
		PrevPath = identify + "/" + flat[index-1].Identify
		PrevName = flat[index-1].DocumentName
	}
	if index == len(flat)-1 {
		c.Data["NextName"] = ""
		NextName = ""
	} else {
		c.Data["NextPath"] = identify + "/" + flat[index+1].Identify
		c.Data["NextName"] = flat[index+1].DocumentName
		NextPath = identify + "/" + flat[index+1].Identify
		NextName = flat[index+1].DocumentName
	}

	doc.IncrViewCount(doc.DocumentId)
	doc.ViewCount = doc.ViewCount + 1
	doc.PutToCache()

	if c.IsAjax() {
		var data struct {
			DocId         int    `json:"doc_id"`
			DocIdentify   string `json:"doc_identify"`
			DocTitle      string `json:"doc_title"`
			Body          string `json:"body"`
			Title         string `json:"title"`
			Version       int64  `json:"version"`
			ViewCount     int    `json:"view_count"`
			MarkdownTheme string `json:"markdown_theme"`
			IsMarkdown    bool   `json:"is_markdown"`
		}
		data.DocId = doc.DocumentId
		data.DocIdentify = doc.Identify
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release + "<div class='wiki-bottom-left'>"+ i18n.Tr(c.Lang, "doc.prev") + "： <a href='/docs/" + PrevPath + "' rel='prev'>" + PrevName + "</a><br />" + i18n.Tr(c.Lang, "doc.next") + "： <a href='/docs/" + NextPath + "' rel='next'>" + NextName + "</a><br /></div>"
		data.Title = doc.DocumentName + " - Powered by MinDoc"
		data.Version = doc.Version
		data.ViewCount = doc.ViewCount
		data.MarkdownTheme = doc.MarkdownTheme
		if bookResult.Editor == EditorCherryMarkdown {
			data.IsMarkdown = true
		}
		c.JsonResult(0, "ok", data)
	} else {
		c.Data["DocumentId"] = doc.DocumentId
		c.Data["DocIdentify"] = doc.Identify
		if bookResult.IsDisplayComment {
			// 获取评论、分页
			comments, count, _ := models.NewComment().QueryCommentByDocumentId(doc.DocumentId, 1, conf.PageSize, c.Member)
			page := pagination.PageUtil(int(count), 1, conf.PageSize, comments)
			c.Data["Page"] = page
		}
	}

	tree, err := models.NewDocument().CreateDocumentTreeForHtml(bookResult.BookId, doc.DocumentId)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("生成项目文档树时出错 ->", err)
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.build_doc_tree_error"))
	}

	c.Data["Description"] = utils.AutoSummary(doc.Release, 120)

	c.Data["Model"] = bookResult
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = doc.DocumentName
	c.Data["Content"] = template.HTML(doc.Release + "<div class='wiki-bottom-left'>"+ i18n.Tr(c.Lang, "doc.prev") + "： <a href='/docs/" + PrevPath + "' rel='prev'>" + PrevName + "</a><br />" + i18n.Tr(c.Lang, "doc.next") + "： <a href='/docs/" + NextPath + "' rel='next'>" + NextName + "</a><br /></div>")
	c.Data["ViewCount"] = doc.ViewCount
	c.Data["FoldSetting"] = "closed"
	if bookResult.Editor == EditorCherryMarkdown {
		c.Data["MarkdownTheme"] = doc.MarkdownTheme
	}
	if doc.IsOpen == 1 {
		c.Data["FoldSetting"] = "open"
	} else if doc.IsOpen == 2 {
		c.Data["FoldSetting"] = "empty"
	}
}

// 递归得到树状结构体
func getTreeRecursive(list []*models.DocumentTree, parentId int) (res []*models.DocumentTree) {
	for _, v := range list {
		if v.ParentId == parentId {
			v.Children = getTreeRecursive(list, v.DocumentId)
			res = append(res, v)
		}
	}
	return res
}

// 递归将树状结构体转换为扁平结构体数组
// func Flatten(list []*models.DocumentTree, flattened *[]DocumentTreeFlatten) (flatten *[]DocumentTreeFlatten) {
func Flatten(list []*models.DocumentTree, flattened *[]DocumentTreeFlatten) {
	// Treeslice := make([]*DocumentTreeFlatten, 0)
	for _, v := range list {
		tree := make([]DocumentTreeFlatten, 1)
		tree[0].DocumentId = v.DocumentId
		tree[0].DocumentName = v.DocumentName
		tree[0].Identify = v.Identify
		*flattened = append(*flattened, tree...)
		if len(v.Children) > 0 {
			Flatten(v.Children, flattened)
		}
	}
	return
}

// 编辑文档
func (c *DocumentController) Edit() {
	c.Prepare()

	if c.Member.Role == conf.MemberReaderRole {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.no_permission"))
	}

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.project_id_error"))
	}

	bookResult := models.NewBookResult()

	var err error
	// 如果是管理者，则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookResult = models.NewBookResult().ToBookResult(*book)
	} else {
		bookResult, err = models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			if err == orm.ErrNoRows || err == models.ErrPermissionDenied {
				c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
			} else {
				logs.Error("查询项目时出错 -> ", err)
				c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
			}
			return
		}
		if bookResult.RoleId == conf.BookObserver {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}
	}

	c.TplName = fmt.Sprintf("document/%s_edit_template.tpl", bookResult.Editor)

	c.Data["Model"] = bookResult

	r, _ := json.Marshal(bookResult)

	c.Data["ModelResult"] = template.JS(string(r))

	c.Data["Result"] = template.JS("[]")

	trees, err := models.NewDocument().FindDocumentTree(bookResult.BookId)

	if err != nil {
		logs.Error("FindDocumentTree => ", err)
	} else {
		if len(trees) > 0 {
			if jtree, err := json.Marshal(trees); err == nil {
				c.Data["Result"] = template.JS(string(jtree))
			}
		} else {
			c.Data["Result"] = template.JS("[]")
		}
	}

	c.Data["BaiDuMapKey"] = web.AppConfig.DefaultString("baidumapkey", "")

	if conf.GetUploadFileSize() > 0 {
		c.Data["UploadFileSize"] = conf.GetUploadFileSize()
	} else {
		c.Data["UploadFileSize"] = "undefined"
	}
}

// 创建一个文档
func (c *DocumentController) Create() {
	identify := c.GetString("identify")
	docIdentify := c.GetString("doc_identify")
	docName := c.GetString("doc_name")
	parentId, _ := c.GetInt("parent_id", 0)
	docId, _ := c.GetInt("doc_id", 0)
	isOpen, _ := c.GetInt("is_open", 0)

	if identify == "" {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	if docName == "" {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.doc_name_empty"))
	}

	bookId := 0

	// 如果是超级管理员则不判断权限
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error(err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_existed_or_no_permit"))
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("FindByIdentify => ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_existed_or_no_permit"))
		}

		bookId = bookResult.BookId
	}

	if docIdentify != "" {
		if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, docIdentify); !ok || err != nil {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.project_id_tips"))
		}

		d, _ := models.NewDocument().FindByIdentityFirst(docIdentify, bookId)
		if d.DocumentId > 0 && d.DocumentId != docId {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.project_id_existed"))
		}
	}
	if parentId > 0 {
		doc, err := models.NewDocument().Find(parentId)
		if err != nil || doc.BookId != bookId {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.parent_id_not_existed"))
		}
	}

	document, _ := models.NewDocument().Find(docId)

	document.MemberId = c.Member.MemberId
	document.BookId = bookId

	document.Identify = docIdentify

	document.Version = time.Now().Unix()
	document.DocumentName = docName
	document.ParentId = parentId

	if isOpen == 1 {
		document.IsOpen = 1
	} else if isOpen == 2 {
		document.IsOpen = 2
	} else {
		document.IsOpen = 0
	}

	if err := document.InsertOrUpdate(); err != nil {
		logs.Error("添加或更新文档时出错 -> ", err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
	} else {
		c.JsonResult(0, "ok", document)
	}
}

// 上传附件或图片
func (c *DocumentController) Upload() {
	identify := c.GetString("identify")
	docId, _ := c.GetInt("doc_id")
	isAttach := true

	if identify == "" {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	names := []string{"editormd-file-file", "editormd-image-file", "file", "editormd-resource-file"}
	var files []*multipart.FileHeader
	for _, name := range names {
		file, err := c.GetFiles(name)
		if err != nil {
			continue
		}
		if len(file) > 0 && err == nil {
			files = append(files, file...)
		}
	}

	if len(files) == 0 {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.upload_file_empty"))
		return
	}

	result2 := []map[string]interface{}{}
	var result map[string]interface{}
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()

		defer file.Close()

		if err != nil {
			c.JsonResult(6002, err.Error())
		}

		// defer file.Close()

		type Size interface {
			Size() int64
		}

		// if conf.GetUploadFileSize() > 0 && moreFile.Size > conf.GetUploadFileSize() {
		if conf.GetUploadFileSize() > 0 && files[i].Size > conf.GetUploadFileSize() {
			c.JsonResult(6009, i18n.Tr(c.Lang, "message.upload_file_size_limit"))
		}

		// ext := filepath.Ext(moreFile.Filename)
		ext := filepath.Ext(files[i].Filename)
		//文件必须带有后缀名
		if ext == "" {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.upload_file_type_error"))
		}
		//如果文件类型设置为 * 标识不限制文件类型
		if conf.IsAllowUploadFileExt(ext) == false {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.upload_file_type_error"))
		}

		bookId := 0

		// 如果是超级管理员，则不判断权限
		if c.Member.IsAdministrator() {
			book, err := models.NewBook().FindByFieldFirst("identify", identify)

			if err != nil {
				c.JsonResult(6006, i18n.Tr(c.Lang, "message.doc_not_exist_or_no_permit"))
			}

			bookId = book.BookId
		} else {
			book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

			if err != nil {
				logs.Error("DocumentController.Edit => ", err)
				if err == orm.ErrNoRows {
					c.JsonResult(6006, i18n.Tr(c.Lang, "message.no_permission"))
				}

				c.JsonResult(6001, err.Error())
			}

			// 如果没有编辑权限
			if book.RoleId != conf.BookEditor && book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
				c.JsonResult(6006, i18n.Tr(c.Lang, "message.no_permission"))
			}

			bookId = book.BookId
		}

		if docId > 0 {
			doc, err := models.NewDocument().Find(docId)
			if err != nil {
				c.JsonResult(6007, i18n.Tr(c.Lang, "message.doc_not_exist"))
			}

			if doc.BookId != bookId {
				c.JsonResult(6008, i18n.Tr(c.Lang, "message.doc_not_belong_project"))
			}
		}

		fileName := "m_" + cryptil.UniqueId() + "_r"
		filePath := filepath.Join(conf.WorkingDirectory, "uploads", identify)

		//将图片和文件分开存放
		attachment := models.NewAttachment()
		var strategy filetil.FileTypeStrategy
		if filetil.IsImageExt(files[i].Filename) {
			strategy = filetil.ImageStrategy{}
			attachment.ResourceType = "image"
		} else if filetil.IsVideoExt(files[i].Filename) {
			strategy = filetil.VideoStrategy{}
			attachment.ResourceType = "video"
		} else {
			strategy = filetil.DefaultStrategy{}
			attachment.ResourceType = "file"
		}

		filePath = strategy.GetFilePath(filePath, fileName, ext)

		path := filepath.Dir(filePath)

		_ = os.MkdirAll(path, os.ModePerm)

		//copy the uploaded file to the destination file
		dst, err := os.Create(filePath)
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			logs.Error("保存文件失败 -> ", err)
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
		}

		attachment.BookId = bookId
		// attachment.FileName = moreFile.Filename
		attachment.FileName = files[i].Filename
		attachment.CreateAt = c.Member.MemberId
		attachment.FileExt = ext
		attachment.FilePath = strings.TrimPrefix(filePath, conf.WorkingDirectory)
		attachment.DocumentId = docId

		if fileInfo, err := os.Stat(filePath); err == nil {
			attachment.FileSize = float64(fileInfo.Size())
		}

		if docId > 0 {
			attachment.DocumentId = docId
		}

		if filetil.IsImageExt(files[i].Filename) || filetil.IsVideoExt(files[i].Filename) {
			attachment.HttpPath = "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)
			if strings.HasPrefix(attachment.HttpPath, "//") {
				attachment.HttpPath = conf.URLForWithCdnImage(string(attachment.HttpPath[1:]))
			}

			isAttach = false
		}

		err = attachment.Insert()

		if err != nil {
			os.Remove(filePath)
			logs.Error("文件保存失败 ->", err)
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed"))
		}

		if attachment.HttpPath == "" {
			attachment.HttpPath = conf.URLForNotHost("DocumentController.DownloadAttachment", ":key", identify, ":attach_id", attachment.AttachmentId)

			if err := attachment.Update(); err != nil {
				logs.Error("保存文件失败 ->", err)
				c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
			}
		}
		result = map[string]interface{}{
			"errcode":       0,
			"success":       1,
			"message":       "ok",
			"url":           attachment.HttpPath,
			"link":          attachment.HttpPath,
			"alt":           attachment.FileName,
			"is_attach":     isAttach,
			"attach":        attachment,
			"resource_type": attachment.ResourceType,
		}
		result2 = append(result2, result)
	}
	if len(files) == 1 {
		// froala单文件上传
		c.Ctx.Output.JSON(result, true, false)
	} else {
		c.Ctx.Output.JSON(result2, true, false)
	}
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
				c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
			} else {
				logs.Error("查找项目时出错 ->", err)
				c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
			}
		}

		// 如果不是超级管理员则判断权限
		if c.Member == nil || c.Member.Role != conf.MemberSuperRole {
			// 如果项目是私有的，并且 token 不正确
			if (book.PrivatelyOwned == 1 && token == "") || (book.PrivatelyOwned == 1 && book.PrivateToken != token) {
				c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
			}
		}

		bookId = book.BookId
	} else {
		bookId = bookResult.BookId
	}

	// 查找附件
	attachment, err := models.NewAttachment().Find(attachId)

	if err != nil {
		logs.Error("查找附件时出错 -> ", err)
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
		} else {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
		}
	}

	if attachment.BookId != bookId {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	}

	c.Ctx.Output.Download(filepath.Join(conf.WorkingDirectory, attachment.FilePath), attachment.FileName)
	c.StopRun()
}

// 删除附件
func (c *DocumentController) RemoveAttachment() {
	c.Prepare()
	attachId, _ := c.GetInt("attach_id")

	if attachId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	attach, err := models.NewAttachment().Find(attachId)

	if err != nil {
		logs.Error(err)
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	}

	document, err := models.NewDocument().Find(attach.DocumentId)

	if err != nil {
		logs.Error(err)
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.doc_not_exist"))
	}

	if c.Member.Role != conf.MemberSuperRole {
		rel, err := models.NewRelationship().FindByBookIdAndMemberId(document.BookId, c.Member.MemberId)
		if err != nil {
			logs.Error(err)
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.no_permission"))
		}

		if rel.RoleId == conf.BookObserver {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.no_permission"))
		}
	}

	err = attach.Delete()
	if err != nil {
		logs.Error(err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
	}

	os.Remove(filepath.Join(conf.WorkingDirectory, attach.FilePath))

	c.JsonResult(0, "ok", attach)
}

// 删除文档
func (c *DocumentController) Delete() {
	c.Prepare()

	identify := c.GetString("identify")
	docId, err := c.GetInt("doc_id", 0)

	bookId := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("FindByIdentify => ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("FindByIdentify => ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = bookResult.BookId
	}

	if docId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	doc, err := models.NewDocument().Find(docId)

	if err != nil {
		logs.Error("Delete => ", err)
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
	}
	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.param_error"))
	}

	// 递归删除项目下的文档以及子文档
	err = doc.RecursiveDocument(doc.DocumentId)
	if err != nil {
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
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
		if err != nil || book == nil {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
			return
		}

		bookId = book.BookId
		autoRelease = book.AutoRelease == 1
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("项目不存在或权限不足 -> ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = bookResult.BookId
		autoRelease = bookResult.AutoRelease
	}

	if docId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	if c.Ctx.Input.IsPost() {
		markdown := strings.TrimSpace(c.GetString("markdown", ""))
		content := c.GetString("html")
		markdownTheme := c.GetString("markdown_theme", "theme__light")
		version, _ := c.GetInt64("version", 0)
		isCover := c.GetString("cover")

		doc, err := models.NewDocument().Find(docId)
		if err != nil || doc == nil {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.read_file_error"))
			return
		}

		if doc.BookId != bookId {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.dock_not_belong_project"))
		}

		if doc.Version != version && !strings.EqualFold(isCover, "yes") {
			logs.Info("%d|", version, doc.Version)
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.confirm_override_doc"))
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
		history.ActionName = i18n.Tr(c.Lang, "doc.modify_doc")

		if markdown == "" && content != "" {
			doc.Markdown = content
		} else {
			doc.Markdown = markdown
			doc.MarkdownTheme = markdownTheme
		}

		doc.Version = time.Now().Unix()
		doc.Content = content
		doc.ModifyAt = c.Member.MemberId

		if err := doc.InsertOrUpdate(); err != nil {
			logs.Error("InsertOrUpdate => ", err)
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed"))
		}

		// 如果启用了文档历史，则添加历史文档
		///如果两次保存的MD5值不同则保存为历史，否则忽略
		go func(history *models.DocumentHistory) {
			if c.EnableDocumentHistory && cryptil.Md5Crypt(history.Markdown) != cryptil.Md5Crypt(doc.Markdown) {
				_, err = history.InsertOrUpdate()
				if err != nil {
					logs.Error("DocumentHistory InsertOrUpdate => ", err)
				}
			}
		}(history)

		//如果启用了自动发布
		if autoRelease {
			go func() {
				doc.Lang = c.Lang
				err := doc.ReleaseContent()
				if err == nil {
					logs.Informational(i18n.Tr(c.Lang, "message.doc_auto_published")+"-> document_id=%d;document_name=%s", doc.DocumentId, doc.DocumentName)
				}
			}()
		}

		c.JsonResult(0, "ok", doc)
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.doc_not_exist"))
		return
	}

	attach, err := models.NewAttachment().FindListByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	c.JsonResult(0, "ok", doc)
}

// Export 导出
func (c *DocumentController) Export() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")

	if identify == "" {
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.param_error"))
	}

	output := c.GetString("output")
	token := c.GetString("token")

	// 如果没有开启匿名访问则跳转到登录
	if !c.EnableAnonymous && !c.isUserLoggedIn() {
		promptUserToLogIn(c)
		return
	}
	if !conf.GetEnableExport() {
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "export_func_disable"))
	}

	bookResult := models.NewBookResult()
	if c.Member != nil && c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByIdentify(identify)
		if err != nil {
			if err == orm.ErrNoRows {
				c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
			} else {
				logs.Error("查找项目时出错 ->", err)
				c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
			}
		}
		bookResult = models.NewBookResult().ToBookResult(*book)
	} else {
		bookResult = c.isReadable(identify, token)
	}
	if !bookResult.IsDownload {
		c.ShowErrorPage(200, i18n.Tr(c.Lang, "message.cur_project_export_func_disable"))
	}

	if !strings.HasPrefix(bookResult.Cover, "http:://") && !strings.HasPrefix(bookResult.Cover, "https:://") {
		bookResult.Cover = conf.URLForWithCdnImage(bookResult.Cover)
	}
	if output == Markdown {
		if bookResult.Editor != EditorMarkdown && bookResult.Editor != EditorCherryMarkdown {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.cur_project_not_support_md"))
		}
		p, err := bookResult.ExportMarkdown(c.CruSession.SessionID(context.TODO()))

		if err != nil {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.failed"))
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

	if output == "pdf" && filetil.FileExists(pdfpath) {
		c.Ctx.Output.Download(pdfpath, bookResult.BookName+".pdf")
		c.Abort("200")
	} else if output == "epub" && filetil.FileExists(epubpath) {
		c.Ctx.Output.Download(epubpath, bookResult.BookName+".epub")

		c.Abort("200")
	} else if output == "mobi" && filetil.FileExists(mobipath) {
		c.Ctx.Output.Download(mobipath, bookResult.BookName+".mobi")

		c.Abort("200")
	} else if output == "docx" && filetil.FileExists(docxpath) {
		c.Ctx.Output.Download(docxpath, bookResult.BookName+".docx")

		c.Abort("200")

	} else if output == "pdf" || output == "epub" || output == "docx" || output == "mobi" {
		if err := models.BackgroundConvert(c.CruSession.SessionID(context.TODO()), bookResult); err != nil && err != gopool.ErrHandlerIsExist {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.export_failed"))
		}

		c.ShowErrorPage(200, i18n.Tr(c.Lang, "message.file_converting"))
	} else {
		c.ShowErrorPage(200, i18n.Tr(c.Lang, "message.unsupport_file_type"))
	}

	c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.no_exportable_file"))
}

// 生成项目访问的二维码
func (c *DocumentController) QrCode() {
	c.Prepare()

	identify := c.GetString(":key")

	book, err := models.NewBook().FindByIdentify(identify)
	if err != nil || book.BookId <= 0 {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	uri := conf.URLFor("DocumentController.Index", ":key", identify)
	code, err := qr.Encode(uri, qr.L, qr.Unicode)
	if err != nil {
		logs.Error("生成二维码失败 ->", err)
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.gen_qrcode_failed"))
	}

	code, err = barcode.Scale(code, 150, 150)
	if err != nil {
		logs.Error("生成二维码失败 ->", err)
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.gen_qrcode_failed"))
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")

	// imgpath := filepath.Join("cache","qrcode",identify + ".png")

	err = png.Encode(c.Ctx.ResponseWriter, code)
	if err != nil {
		logs.Error("生成二维码失败 ->", err)
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.gen_qrcode_failed"))
	}
}

// 项目内搜索
func (c *DocumentController) Search() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	token := c.GetString("token")
	keyword := strings.TrimSpace(c.GetString("keyword"))

	if identify == "" {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	if !c.EnableAnonymous && !c.isUserLoggedIn() {
		promptUserToLogIn(c)
		return
	}

	bookResult := c.isReadable(identify, token)

	docs, err := models.NewDocumentSearchResult().SearchDocument(keyword, bookResult.BookId)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.search_result_error"))
	}

	if len(docs) < 0 {
		c.JsonResult(404, i18n.Tr(c.Lang, "message.no_data"))
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

	bookId := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("查找项目失败 ->", err)
			c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit")
			return
		}

		bookId = book.BookId
		c.Data["Model"] = book
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("查找项目失败 ->", err)
			c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit")
			return
		}

		bookId = bookResult.BookId
		c.Data["Model"] = bookResult
	}

	if docId <= 0 {
		c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.param_error")
		return
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		logs.Error("Delete => ", err)
		c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.get_doc_his_failed")
		return
	}

	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.param_error")
		return
	}

	histories, totalCount, err := models.NewDocumentHistory().FindToPager(docId, pageIndex, conf.PageSize)
	if err != nil {
		logs.Error("分页查找文档历史失败 ->", err)
		c.Data["ErrorMessage"] = i18n.Tr(c.Lang, "message.get_doc_his_failed")
		return
	}

	c.Data["List"] = histories
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
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	bookId := 0

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("查找项目失败 ->", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("查找项目失败 ->", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = bookResult.BookId
	}

	if docId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		logs.Error("Delete => ", err)
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.get_doc_his_failed"))
	}

	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	err = models.NewDocumentHistory().Delete(historyId, docId)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.failed"))
	}

	c.JsonResult(0, "ok")
}

// 通过文档历史恢复文档
func (c *DocumentController) RestoreHistory() {
	c.Prepare()

	c.TplName = "document/history.tpl"

	identify := c.GetString("identify")
	docId, err := c.GetInt("doc_id", 0)
	historyId, _ := c.GetInt("history_id", 0)

	if historyId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	bookId := 0
	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("FindByIdentify => ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("FindByIdentify => ", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}

		bookId = bookResult.BookId
	}

	if docId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	doc, err := models.NewDocument().Find(docId)
	if err != nil {
		logs.Error("Delete => ", err)
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.get_doc_his_failed"))
	}

	// 如果文档所属项目错误
	if doc.BookId != bookId {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	err = models.NewDocumentHistory().Restore(historyId, docId, c.Member.MemberId)
	if err != nil {
		logs.Error(err)
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.failed"))
	}

	c.JsonResult(0, "ok", doc)
}

func (c *DocumentController) Compare() {
	c.Prepare()

	c.TplName = "document/compare.tpl"

	historyId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	identify := c.Ctx.Input.Param(":key")

	bookId := 0
	editor := EditorMarkdown

	// 如果是超级管理员则忽略权限判断
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("DocumentController.Compare => ", err)
			c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
			return
		}

		bookId = book.BookId
		c.Data["Model"] = book
		editor = book.Editor
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil || bookResult.RoleId == conf.BookObserver {
			logs.Error("FindByIdentify => ", err)
			c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
			return
		}

		bookId = bookResult.BookId
		c.Data["Model"] = bookResult
		editor = bookResult.Editor
	}

	if historyId <= 0 {
		c.ShowErrorPage(60002, i18n.Tr(c.Lang, "message.param_error"))
	}

	history, err := models.NewDocumentHistory().Find(historyId)
	if err != nil {
		logs.Error("DocumentController.Compare => ", err)
		c.ShowErrorPage(60003, err.Error())
	}

	doc, err := models.NewDocument().Find(history.DocumentId)
	if err != nil || doc == nil || doc.BookId != bookId {
		c.ShowErrorPage(60002, i18n.Tr(c.Lang, "message.doc_not_exist"))
		return
	}

	c.Data["HistoryId"] = historyId
	c.Data["DocumentId"] = doc.DocumentId

	if editor == EditorMarkdown || editor == EditorCherryMarkdown {
		c.Data["HistoryContent"] = history.Markdown
		c.Data["Content"] = doc.Markdown
	} else {
		c.Data["HistoryContent"] = template.HTML(history.Content)
		c.Data["Content"] = template.HTML(doc.Content)
	}
}

// 判断用户是否可以阅读文档
func (c *DocumentController) isReadable(identify, token string) *models.BookResult {
	book, err := models.NewBook().FindByFieldFirst("identify", identify)

	if err != nil {
		logs.Error(err)
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.item_not_exist"))
	}
	bookResult := models.NewBookResult().ToBookResult(*book)
	isOk := false

	if c.isUserLoggedIn() {
		roleId, err := models.NewBook().FindForRoleId(book.BookId, c.Member.MemberId)
		if err == nil {
			isOk = true
			bookResult.MemberId = c.Member.MemberId
			bookResult.RoleId = roleId
		}
	}

	/* 	私有项目：
	 *   管理员可以直接访问
	 *   参与者可以直接访问
	 *   其他用户（支持匿名访问）
	 *   	token设置情况
	 *   		已设置：可以通过token访问
	 *   		未设置：不可以通过token访问
	 *   	password设置情况
	 *   		已设置：可以通过password访问
	 *   		未设置：不可以通过password访问
	 *   注意：
	 *   1. 第一次访问需要存session
	 *   2. 有session优先使用session中的token或者password，再使用携带的token或者password
	 *   3. 私有项目如果token和password都没有设置，则除管理员和参与者的其他用户不可以访问
	 *   4. 使用token访问如果不通过，则提示输入密码
	 */
	if book.PrivatelyOwned == 1 {
		if c.isUserLoggedIn() && c.Member.IsAdministrator() {
			return bookResult
		}
		if isOk { // Project participant.
			return bookResult
		}

		// Use session in preference.
		if tokenOrPassword, ok := c.GetSession(identify).(string); ok {
			if strings.EqualFold(book.PrivateToken, tokenOrPassword) || strings.EqualFold(book.BookPassword, tokenOrPassword) {
				return bookResult
			}
		}

		// Next: Session not exist or not correct.
		if book.PrivateToken != "" && book.PrivateToken == token {
			c.SetSession(identify, token)
			return bookResult
		} else if book.BookPassword != "" {
			// Send a page for inputting password.
			// For verification, see function DocumentController.CheckPassword
			body, err := c.ExecuteViewPathTemplate("document/document_password.tpl",
				map[string]string{"Identify": book.Identify, "Lang": c.Lang})
			if err != nil {
				logs.Error("显示密码页面失败 ->", err)
				c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
			}
			c.CustomAbort(200, body)
		} else {
			// No permission to access this book.
			logs.Info("尝试访问文档但权限不足 ->", identify, token)
			c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
		}
	}

	return bookResult
}

func promptUserToLogIn(c *DocumentController) {
	logs.Info("Access " + c.Ctx.Request.URL.RequestURI() + " not permitted.")
	logs.Info("  Access will be redirected to login page(SessionId: " + c.CruSession.SessionID(context.TODO()) + ").")

	if c.IsAjax() {
		c.JsonResult(6000, i18n.Tr(c.Lang, "message.need_relogin"))
	} else {
		c.Redirect(conf.URLFor("AccountController.Login")+"?url="+url.PathEscape(conf.BaseUrl+c.Ctx.Request.URL.RequestURI()), 302)
	}
}
