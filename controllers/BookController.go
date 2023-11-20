package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/i18n"
	"github.com/mindoc-org/mindoc/utils/sqltil"

	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/graphics"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/pagination"
	"github.com/russross/blackfriday/v2"
)

type BookController struct {
	BaseController
}

func (c *BookController) Index() {
	c.Prepare()
	c.TplName = "book/index.tpl"

	pageIndex, _ := c.GetInt("page", 1)

	books, totalCount, err := models.NewBook().FindToPager(pageIndex, conf.PageSize, c.Member.MemberId, c.Lang)

	if err != nil {
		logs.Error("BookController.Index => ", err)
		c.Abort("500")
	}

	for i, book := range books {
		books[i].Description = utils.StripTags(string(blackfriday.Run([]byte(book.Description))))
		books[i].ModifyTime = book.ModifyTime.Local()
		books[i].CreateTime = book.CreateTime.Local()
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	b, err := json.Marshal(books)

	if err != nil || len(books) <= 0 {
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
	if itemsets, err := models.NewItemsets().First(1); err == nil {
		c.Data["Item"] = itemsets
	}
}

// Dashboard 项目概要 .
func (c *BookController) Dashboard() {
	c.Prepare()
	c.TplName = "book/dashboard.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == "" {
		c.Abort("404")
	}

	book, err := models.NewBookResult().SetLang(c.Lang).FindByIdentify(key, c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
		return
	}

	c.Data["Description"] = template.HTML(blackfriday.Run([]byte(book.Description)))
	c.Data["Model"] = *book
}

// Setting 项目设置 .
func (c *BookController) Setting() {
	c.Prepare()
	c.TplName = "book/setting.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == "" {
		c.Abort("404")
	}

	book, err := models.NewBookResult().FindByIdentify(key, c.Member.MemberId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Abort("404")
		}
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
		return
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}
	if book.PrivateToken != "" {
		book.PrivateToken = conf.URLFor("DocumentController.Index", ":key", book.Identify, "token", book.PrivateToken)
	}
	c.Data["Model"] = book

}

// 保存项目信息
func (c *BookController) SaveBook() {
	bookResult, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
	}
	book, err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		logs.Error("SaveBook => ", err)
		c.JsonResult(6002, err.Error())
	}

	bookName := strings.TrimSpace(c.GetString("book_name"))
	description := strings.TrimSpace(c.GetString("description", ""))
	commentStatus := c.GetString("comment_status")
	//tag := strings.TrimSpace(c.GetString("label"))
	editor := strings.TrimSpace(c.GetString("editor"))
	autoRelease := strings.TrimSpace(c.GetString("auto_release")) == "on"
	publisher := strings.TrimSpace(c.GetString("publisher"))
	historyCount, _ := c.GetInt("history_count", 0)
	isDownload := strings.TrimSpace(c.GetString("is_download")) == "on"
	enableShare := strings.TrimSpace(c.GetString("enable_share")) == "on"
	isUseFirstDocument := strings.TrimSpace(c.GetString("is_use_first_document")) == "on"
	autoSave := strings.TrimSpace(c.GetString("auto_save")) == "on"
	itemId, _ := c.GetInt("itemId")

	if strings.Count(description, "") > 500 {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.project_desc_tips"))
	}
	if commentStatus != "open" && commentStatus != "closed" && commentStatus != "group_only" && commentStatus != "registered_only" {
		commentStatus = "closed"
	}

	if !models.NewItemsets().Exist(itemId) {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.project_space_not_exist"))
	}
	if editor != EditorMarkdown && editor != EditorCherryMarkdown && editor != EditorHtml && editor != EditorNewHtml {
		editor = EditorMarkdown
	}

	book.BookName = bookName
	book.Description = description
	book.CommentStatus = commentStatus
	book.Publisher = publisher
	//book.Label = tag
	if book.Editor == EditorMarkdown && editor == EditorCherryMarkdown || book.Editor == EditorCherryMarkdown && editor == EditorMarkdown {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.editors_not_compatible"))
	}
	book.Editor = editor
	if editor == EditorCherryMarkdown {
		book.Theme = "cherry"
	}
	book.HistoryCount = historyCount
	book.IsDownload = 0
	book.BookPassword = strings.TrimSpace(c.GetString("bPassword"))
	book.ItemId = itemId

	if autoRelease {
		book.AutoRelease = 1
	} else {
		book.AutoRelease = 0
	}
	if isDownload {
		book.IsDownload = 0
	} else {
		book.IsDownload = 1
	}
	if enableShare {
		book.IsEnableShare = 0
	} else {
		book.IsEnableShare = 1
	}
	if isUseFirstDocument {
		book.IsUseFirstDocument = 1
	} else {
		book.IsUseFirstDocument = 0
	}
	if autoSave {
		book.AutoSave = 1
	} else {
		book.AutoSave = 0
	}
	if err := book.Update(); err != nil {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed"))
	}
	bookResult.BookName = bookName
	bookResult.Description = description
	bookResult.CommentStatus = commentStatus

	logs.Info("用户 [", c.Member.Account, "] 修改了项目 ->", book)

	c.JsonResult(0, "ok", bookResult)
}

// 设置项目私有状态.
func (c *BookController) PrivatelyOwned() {

	status := c.GetString("status")

	if status != "open" && status != "close" {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.param_error"))
	}
	state := 0
	if status == "open" {
		state = 0
	} else {
		state = 1
	}

	bookResult, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
		return
	}
	//只有创始人才能变更私有状态
	if bookResult.RoleId != conf.BookFounder {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.no_permission"))
	}

	book, err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.item_not_exist"))
		return
	}
	book.PrivatelyOwned = state

	err = book.Update()

	if err != nil {
		logs.Error("PrivatelyOwned => ", err)
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.failed"))
	}
	logs.Info("用户 【", c.Member.Account, "]修改了项目权限 ->", state)
	c.JsonResult(0, "ok")
}

// Transfer 转让项目.
func (c *BookController) Transfer() {
	c.Prepare()
	account := c.GetString("account")

	if account == "" {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.receive_account_empty"))
	}
	member, err := models.NewMember().FindByAccount(account)

	if err != nil {
		logs.Error("FindByAccount => ", err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.receive_account_not_exist"))
	}
	if member.Status != 0 {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.receive_account_disabled"))
	}
	if member.MemberId == c.Member.MemberId {
		c.JsonResult(6007, i18n.Tr(c.Lang, "message.cannot_handover_myself"))
	}
	bookResult, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
		return
	}

	err = models.NewRelationship().Transfer(bookResult.BookId, c.Member.MemberId, member.MemberId)

	if err != nil {
		logs.Error("转让项目失败 -> ", err)
		c.JsonResult(6008, err.Error())
	}
	c.JsonResult(0, "ok")
}

// 上传项目封面.
func (c *BookController) UploadCover() {

	bookResult, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
		return
	}
	book, err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		logs.Error("SaveBook => ", err)
		c.JsonResult(6002, err.Error())
	}

	file, moreFile, err := c.GetFile("image-file")

	if err != nil {
		logs.Error("获取上传文件失败 ->", err.Error())
		c.JsonResult(500, "读取文件异常")
		return
	}
	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		c.JsonResult(500, "不支持的图片格式")
	}

	x1, _ := strconv.ParseFloat(c.GetString("x"), 10)
	y1, _ := strconv.ParseFloat(c.GetString("y"), 10)
	w1, _ := strconv.ParseFloat(c.GetString("width"), 10)
	h1, _ := strconv.ParseFloat(c.GetString("height"), 10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	fileName := "cover_" + strconv.FormatInt(time.Now().UnixNano(), 16)

	//附件路径按照项目组织
	// 	filePath := filepath.Join("uploads", book.Identify, "images", fileName+ext)
	filePath := filepath.Join(conf.WorkingDirectory, "uploads", book.Identify, "images", fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file", filePath)

	if err != nil {
		logs.Error("", err)
		c.JsonResult(500, "图片保存失败")
	}
	defer func(filePath string) {
		os.Remove(filePath)
	}(filePath)

	//剪切图片
	subImg, err := graphics.ImageCopyFromFile(filePath, x, y, width, height)

	if err != nil {
		logs.Error("graphics.ImageCopyFromFile => ", err)
		c.JsonResult(500, "图片剪切")
	}

	filePath = filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName+"_small"+ext)

	//生成缩略图并保存到磁盘
	err = graphics.ImageResizeSaveFile(subImg, 350, 460, filePath)

	if err != nil {
		logs.Error("ImageResizeSaveFile => ", err.Error())
		c.JsonResult(500, "保存图片失败")
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)

	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}

	oldCover := book.Cover

	book.Cover = conf.URLForWithCdnImage(url)

	if err := book.Update(); err != nil {
		c.JsonResult(6001, "保存图片失败")
	}
	//如果原封面不是默认封面则删除
	if oldCover != conf.GetDefaultCover() {
		os.Remove("." + oldCover)
	}
	logs.Info("用户[", c.Member.Account, "]上传了项目封面 ->", book.BookName, book.BookId, book.Cover)

	c.JsonResult(0, "ok", url)
}

// Users 用户列表.
func (c *BookController) Users() {
	c.Prepare()
	c.TplName = "book/users.tpl"

	key := c.Ctx.Input.Param(":key")
	pageIndex, _ := c.GetInt("page", 1)

	if key == "" {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	book, err := models.NewBookResult().FindByIdentify(key, c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
		return
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}
	c.Data["Model"] = *book

	members, totalCount, err := models.NewMemberRelationshipResult().FindForUsersByBookId(c.Lang, book.BookId, pageIndex, conf.PageSize)

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	b, err := json.Marshal(members)

	if err != nil {
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
}

// Create 创建项目.
func (c *BookController) Create() {

	if c.Ctx.Input.IsPost() {
		bookName := strings.TrimSpace(c.GetString("book_name", ""))
		identify := strings.TrimSpace(c.GetString("identify", ""))
		description := strings.TrimSpace(c.GetString("description", ""))
		privatelyOwned, _ := strconv.Atoi(c.GetString("privately_owned"))
		commentStatus := c.GetString("comment_status")
		itemId, _ := c.GetInt("itemId")

		if bookName == "" {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.project_name_empty"))
		}
		if identify == "" {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.project_id_empty"))
		}
		if ok, err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`, identify); !ok || err != nil {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.project_id_tips"))
		}
		if strings.Count(identify, "") > 50 {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.project_id_length"))
		}
		if strings.Count(description, "") > 500 {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.project_desc_tips"))
		}
		if privatelyOwned != 0 && privatelyOwned != 1 {
			privatelyOwned = 1
		}
		if !models.NewItemsets().Exist(itemId) {
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.project_space_not_exist"))
		}
		if commentStatus != "open" && commentStatus != "closed" && commentStatus != "group_only" && commentStatus != "registered_only" {
			commentStatus = "closed"
		}
		book := models.NewBook()
		book.Cover = conf.GetDefaultCover()

		//如果客户端上传了项目封面则直接保存
		if file, moreFile, err := c.GetFile("image-file"); err == nil {
			defer file.Close()

			ext := filepath.Ext(moreFile.Filename)

			//如果上传的是图片
			if strings.EqualFold(ext, ".png") || strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".gif") || strings.EqualFold(ext, ".jpeg") {

				fileName := "cover_" + strconv.FormatInt(time.Now().UnixNano(), 16)

				filePath := filepath.Join("uploads", time.Now().Format("200601"), fileName+ext)

				path := filepath.Dir(filePath)

				os.MkdirAll(path, os.ModePerm)

				if err := c.SaveToFile("image-file", filePath); err == nil {
					url := "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)

					if strings.HasPrefix(url, "//") {
						url = string(url[1:])
					}
					book.Cover = url
				}
			}
		}

		if books, _ := book.FindByField("identify", identify, "book_id"); len(books) > 0 {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.project_id_existed"))
		}

		book.BookName = bookName
		book.Description = description
		book.CommentCount = 0
		book.PrivatelyOwned = privatelyOwned
		book.CommentStatus = commentStatus
		book.Identify = identify
		book.DocCount = 0
		book.MemberId = c.Member.MemberId
		book.Version = time.Now().Unix()
		book.IsEnableShare = 0
		book.IsUseFirstDocument = 1
		book.IsDownload = 1
		book.AutoRelease = 0
		book.ItemId = itemId

		book.Editor = "markdown"
		book.Theme = "default"

		if err := book.Insert(c.Lang); err != nil {
			logs.Error("Insert => ", err)
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
		}
		bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)

		if err != nil {
			logs.Error(err)
		}

		logs.Info("用户[", c.Member.Account, "]创建了项目 ->", book)
		c.JsonResult(0, "ok", bookResult)
	}
	c.JsonResult(6001, "error")
}

// 复制项目
func (c *BookController) Copy() {
	if c.Ctx.Input.IsPost() {
		//检查是否有复制项目的权限
		if _, err := c.IsPermission(); err != nil {
			c.JsonResult(500, err.Error())
		}

		identify := strings.TrimSpace(c.GetString("identify", ""))
		if identify == "" {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
		}
		book := models.NewBook()
		err := book.Copy(identify)
		if err != nil {
			c.JsonResult(6002, "复制项目出错")
		} else {
			bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)
			if err != nil {
				logs.Error("查询失败")
			}
			c.JsonResult(0, "ok", bookResult)
		}
	}
}

// 导入zip压缩包或docx
func (c *BookController) Import() {

	file, moreFile, err := c.GetFile("import-file")
	if err == http.ErrMissingFile {
		c.JsonResult(6003, "没有发现需要上传的文件")
	}

	defer file.Close()

	bookName := strings.TrimSpace(c.GetString("book_name"))
	identify := strings.TrimSpace(c.GetString("identify"))
	description := strings.TrimSpace(c.GetString("description", ""))
	privatelyOwned, _ := strconv.Atoi(c.GetString("privately_owned"))
	itemId, _ := c.GetInt("itemId")

	if bookName == "" {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.project_name_empty"))
	}
	if len([]rune(bookName)) > 500 {
		c.JsonResult(6002, "项目名称不能大于500字")
	}
	if identify == "" {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.project_id_empty"))
	}
	if ok, err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`, identify); !ok || err != nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.project_id_tips"))
	}
	if !models.NewItemsets().Exist(itemId) {
		c.JsonResult(6007, i18n.Tr(c.Lang, "message.project_space_not_exist"))
	}
	if strings.Count(identify, "") > 50 {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.project_id_length"))
	}

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext, ".zip") && !strings.EqualFold(ext, ".docx") {
		c.JsonResult(6004, "不支持的文件类型")
	}

	if books, _ := models.NewBook().FindByField("identify", identify, "book_id"); len(books) > 0 {
		c.JsonResult(6006, i18n.Tr(c.Lang, "message.project_id_existed"))
	}

	tempPath := filepath.Join(os.TempDir(), c.CruSession.SessionID(context.TODO()))

	os.MkdirAll(tempPath, 0766)

	tempPath = filepath.Join(tempPath, moreFile.Filename)

	err = c.SaveToFile("import-file", tempPath)
	if err != nil {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.upload_failed"))
	}

	book := models.NewBook()

	book.MemberId = c.Member.MemberId
	book.Cover = conf.GetDefaultCover()
	book.BookName = bookName
	book.Description = description
	book.CommentCount = 0
	book.PrivatelyOwned = privatelyOwned
	book.CommentStatus = "closed"
	book.Identify = identify
	book.DocCount = 0
	book.MemberId = c.Member.MemberId
	book.Version = time.Now().Unix()
	book.ItemId = itemId

	book.Editor = "markdown"
	book.Theme = "default"

	if strings.EqualFold(ext, ".zip") {
		go book.ImportBook(tempPath, c.Lang)
	} else if strings.EqualFold(ext, ".docx") {
		go book.ImportWordBook(tempPath, c.Lang)
	}

	logs.Info("用户[", c.Member.Account, "]导入了项目 ->", book)

	c.JsonResult(0, "项目正在后台转换中，请稍后查看")
}

// CreateToken 创建访问来令牌.
//func (c *BookController) CreateToken() {
//
//	action := c.GetString("action")
//
//	bookResult, err := c.IsPermission()
//
//	if err != nil {
//		if err == models.ErrPermissionDenied {
//			c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
//		}
//		if err == orm.ErrNoRows {
//			c.JsonResult(404, i18n.Tr(c.Lang, "message.item_not_exist"))
//		}
//		logs.Error("生成阅读令牌失败 =>", err)
//		c.JsonResult(6002, err.Error())
//	}
//	book := models.NewBook()
//
//	if _, err := book.Find(bookResult.BookId); err != nil {
//		c.JsonResult(6001, i18n.Tr(c.Lang, "message.item_not_exist"))
//	}
//	if action == "create" {
//		if bookResult.PrivatelyOwned == 0 {
//			c.JsonResult(6001, "公开项目不能创建阅读令牌")
//		}
//
//		book.PrivateToken = string(utils.Krand(conf.GetTokenSize(), utils.KC_RAND_KIND_ALL))
//		if err := book.Update(); err != nil {
//			logs.Error("生成阅读令牌失败 => ", err)
//			c.JsonResult(6003, "生成阅读令牌失败")
//		}
//		logs.Info("用户[", c.Member.Account, "]创建项目令牌 ->", book.PrivateToken)
//		c.JsonResult(0, "ok", conf.URLFor("DocumentController.Index", ":key", book.Identify, "token", book.PrivateToken))
//	} else {
//		book.PrivateToken = ""
//		if err := book.Update(); err != nil {
//			logs.Error("CreateToken => ", err)
//			c.JsonResult(6004, "删除令牌失败")
//		}
//		logs.Info("用户[", c.Member.Account, "]创建项目令牌 ->", book.PrivateToken)
//		c.JsonResult(0, "ok", "")
//	}
//}

// Delete 删除项目.
func (c *BookController) Delete() {
	c.Prepare()

	bookResult, err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001, err.Error())
		return
	}

	if bookResult.RoleId != conf.BookFounder {
		c.JsonResult(6002, "只有创始人才能删除项目")
	}
	err = models.NewBook().ThoroughDeleteBook(bookResult.BookId)

	if err == orm.ErrNoRows {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist"))
	}
	if err != nil {
		logs.Error("删除项目 => ", err)
		c.JsonResult(6003, "删除失败")
	}
	logs.Info("用户[", c.Member.Account, "]删除了项目 ->", bookResult)
	c.JsonResult(0, "ok")
}

// 发布项目.
func (c *BookController) Release() {
	c.Prepare()

	identify := c.GetString("identify")

	bookId := 0

	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil {
			logs.Error("发布文档失败 ->", err)
			c.JsonResult(6003, "文档不存在")
			return
		}
		bookId = book.BookId
	} else {
		book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			if err == models.ErrPermissionDenied {
				c.JsonResult(6001, i18n.Tr(c.Lang, "message.no_permission"))
			}
			if err == orm.ErrNoRows {
				c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist"))
			}
			logs.Error(err)
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.unknown_exception"))
		}
		if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder && book.RoleId != conf.BookEditor {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.no_permission"))
		}
		bookId = book.BookId
	}
	go models.NewBook().ReleaseContent(bookId, c.Lang)

	c.JsonResult(0, i18n.Tr(c.Lang, "message.publish_to_queue"))
}

// 文档排序.
func (c *BookController) SaveSort() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	bookId := 0
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().FindByFieldFirst("identify", identify)
		if err != nil || book == nil {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.item_not_exist"))
			return
		}
		bookId = book.BookId
	} else {
		bookResult, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)
		if err != nil {
			logs.Error("DocumentController.Edit => ", err)

			c.Abort("403")
		}
		if bookResult.RoleId == conf.BookObserver {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
		}
		bookId = bookResult.BookId
	}

	content := c.Ctx.Input.RequestBody

	var docs []map[string]interface{}

	err := json.Unmarshal(content, &docs)

	if err != nil {
		logs.Error(err)
		c.JsonResult(6003, "数据错误")
	}

	for _, item := range docs {
		if docId, ok := item["id"].(float64); ok {
			doc, err := models.NewDocument().Find(int(docId))
			if err != nil {
				logs.Error(err)
				continue
			}
			if doc.BookId != bookId {
				logs.Info("%s", i18n.Tr(c.Lang, "message.no_permission"))
				continue
			}
			sort, ok := item["sort"].(float64)
			if !ok {
				logs.Info("排序数字转换失败 => ", item)
				continue
			}
			parentId, ok := item["parent"].(float64)
			if !ok {
				logs.Info("父分类转换失败 => ", item)
				continue
			}
			if parentId > 0 {
				if parent, err := models.NewDocument().Find(int(parentId)); err != nil || parent.BookId != bookId {
					continue
				}
			}
			doc.OrderSort = int(sort)
			doc.ParentId = int(parentId)
			if err := doc.InsertOrUpdate(); err != nil {
				fmt.Printf("%s", err.Error())
				logs.Error(err)
			}
		} else {
			fmt.Printf("文档ID转换失败 => %+v", item)
		}

	}
	c.JsonResult(0, "ok")
}

func (c *BookController) Team() {
	c.Prepare()
	c.TplName = "book/team.tpl"

	key := c.Ctx.Input.Param(":key")
	pageIndex, _ := c.GetInt("page", 1)

	if key == "" {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	book, err := models.NewBookResult().FindByIdentify(key, c.Member.MemberId)
	if err != nil || book == nil {
		if err == models.ErrPermissionDenied {
			c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
		}
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.system_error"))
		return
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}
	c.Data["Model"] = book

	members, totalCount, err := models.NewTeamRelationship().FindByBookToPager(book.BookId, pageIndex, conf.PageSize)

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	b, err := json.Marshal(members)

	if err != nil {
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
}

func (c *BookController) TeamAdd() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")

	book, err := c.IsPermission()

	if err != nil {
		c.JsonResult(500, err.Error())
		return
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}
	_, err = models.NewTeam().First(teamId, "team_id")
	if err != nil {
		if err == orm.ErrNoRows {
			c.JsonResult(500, "团队不存在")
		}
		c.JsonResult(5002, err.Error())
	}
	if _, err := models.NewTeamRelationship().FindByBookId(book.BookId, teamId); err == nil {
		c.JsonResult(5003, "团队已加入当前项目")
	}
	teamRel := models.NewTeamRelationship()
	teamRel.BookId = book.BookId
	teamRel.TeamId = teamId
	err = teamRel.Save()
	if err != nil {
		c.JsonResult(5004, "加入项目失败")
		return
	}
	teamRel.Include()

	c.JsonResult(0, "OK", teamRel)
}

// 删除项目的团队.
func (c *BookController) TeamDelete() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")

	if teamId <= 0 {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.param_error"))
	}
	book, err := c.IsPermission()

	if err != nil {
		c.JsonResult(5002, err.Error())
		return
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}

	err = models.NewTeamRelationship().DeleteByBookId(book.BookId, teamId)

	if err != nil {
		if err == orm.ErrNoRows {
			c.JsonResult(5003, "团队未加入项目")
		}
		c.JsonResult(5004, err.Error())
	}
	c.JsonResult(0, "OK")
}

// 团队搜索.
func (c *BookController) TeamSearch() {
	c.Prepare()

	keyword := strings.TrimSpace(c.GetString("q"))

	book, err := c.IsPermission()

	if err != nil {
		c.JsonResult(500, err.Error())
	}
	keyword = sqltil.EscapeLike(keyword)
	searchResult, err := models.NewTeamRelationship().FindNotJoinBookByBookIdentify(book.BookId, keyword, 10)

	if err != nil {
		c.JsonResult(500, err.Error(), searchResult)
	}
	c.JsonResult(0, "OK", searchResult)

}

// 项目空间搜索.
func (c *BookController) ItemsetsSearch() {
	c.Prepare()

	keyword := strings.TrimSpace(c.GetString("q"))
	keyword = sqltil.EscapeLike(keyword)

	searchResult, err := models.NewItemsets().FindItemsetsByName(keyword, 10)

	if err != nil {
		c.JsonResult(500, err.Error(), searchResult)
	}
	c.JsonResult(0, "OK", searchResult)

}

func (c *BookController) IsPermission() (*models.BookResult, error) {
	identify := c.GetString("identify")

	if identify == "" {
		return nil, errors.New(i18n.Tr(c.Lang, "message.param_error"))
	}

	book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			return book, errors.New(i18n.Tr(c.Lang, "message.no_permission"))
		}
		if err == orm.ErrNoRows {
			return book, errors.New(i18n.Tr(c.Lang, "message.item_not_exist"))
		}
		return book, err
	}
	if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
		return book, errors.New(i18n.Tr(c.Lang, "message.no_permission"))
	}
	return book, nil
}
