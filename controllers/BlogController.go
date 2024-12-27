package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/pagination"
)

type BlogController struct {
	BaseController
}

func (c *BlogController) Prepare() {
	c.BaseController.Prepare()
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login")+"?url="+url.PathEscape(conf.BaseUrl+c.Ctx.Request.URL.RequestURI()), 302)
	}
}

// 文章阅读
func (c *BlogController) Index() {
	c.Prepare()
	c.TplName = "blog/index.tpl"
	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if blogId <= 0 {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.page_not_existed"))
	}
	blogReadSession := fmt.Sprintf("blog:read:%d", blogId)

	blog, err := models.NewBlog().FindFromCache(blogId)

	if err != nil {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.blog_not_existed"))
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		if blog.BlogStatus == "password" && password != blog.Password {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.blog_pwd_incorrect"))
		} else if blog.BlogStatus == "password" && password == blog.Password {
			// Store the session value for the next GET request.
			_ = c.CruSession.Set(context.TODO(), blogReadSession, blogId)
			c.JsonResult(0, "OK")
		} else {
			c.JsonResult(0, "OK")
		}
	} else if blog.BlogStatus == "password" && c.CruSession.Get(context.TODO(), blogReadSession) == nil && // Read session doesn't exist
		(c.Member == nil || (blog.MemberId != c.Member.MemberId && !c.Member.IsAdministrator())) { // User isn't author or administrator
		//如果不存在已输入密码的标记
		c.TplName = "blog/index_password.tpl"
	}
	if blog.BlogType != 1 {
		//加载文章附件
		_ = blog.LinkAttach()
	}

	c.Data["Model"] = blog
	c.Data["Content"] = template.HTML(blog.BlogRelease)

	if blog.BlogExcerpt == "" {
		c.Data["Description"] = utils.AutoSummary(blog.BlogRelease, 120)
	} else {
		c.Data["Description"] = blog.BlogExcerpt
	}

	if nextBlog, err := models.NewBlog().QueryNext(blogId); err == nil {
		c.Data["Next"] = nextBlog
	}
	if preBlog, err := models.NewBlog().QueryPrevious(blogId); err == nil {
		c.Data["Previous"] = preBlog
	}
}

// 文章列表
func (c *BlogController) List() {
	c.Prepare()
	c.TplName = "blog/list.tpl"
	pageIndex, _ := c.GetInt("page", 1)

	var blogList []*models.Blog
	var totalCount int
	var err error

	blogList, totalCount, err = models.NewBlog().FindToPager(pageIndex, conf.PageSize, 0, "")

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
		for _, blog := range blogList {
			//如果没有添加文章摘要，则自动提取
			if blog.BlogExcerpt == "" {
				blog.BlogExcerpt = utils.AutoSummary(blog.BlogRelease, 120)
			}
			blog.Link()
		}
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["Lists"] = blogList
}

// 管理后台文章列表
func (c *BlogController) ManageList() {
	c.Prepare()
	c.TplName = "blog/manage_list.tpl"

	pageIndex, _ := c.GetInt("page", 1)

	blogList, totalCount, err := models.NewBlog().FindToPager(pageIndex, conf.PageSize, c.Member.MemberId, "all")

	if err != nil {
		c.ShowErrorPage(500, err.Error())
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["ModelList"] = blogList

}

// 文章设置
func (c *BlogController) ManageSetting() {
	c.Prepare()
	c.TplName = "blog/manage_setting.tpl"
	//如果是post请求
	if c.Ctx.Input.IsPost() {
		blogId, _ := c.GetInt("id", 0)
		blogTitle := c.GetString("title")
		blogIdentify := c.GetString("identify")
		orderIndex, _ := c.GetInt("order_index", 0)
		blogType, _ := c.GetInt("blog_type", 0)
		blogExcerpt := c.GetString("excerpt", "")
		blogStatus := c.GetString("status", "publish")
		blogPassword := c.GetString("password", "")
		documentIdentify := strings.TrimSpace(c.GetString("documentIdentify"))
		bookIdentify := strings.TrimSpace(c.GetString("bookIdentify"))
		documentId := 0

		if c.Member.Role == conf.MemberReaderRole {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.no_permission"))
		}
		if blogTitle == "" {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.blog_title_empty"))
		}
		if strings.Count(blogExcerpt, "") > 500 {
			c.JsonResult(6008, i18n.Tr(c.Lang, "message.blog_digest_tips"))
		}
		if blogStatus != "private" && blogStatus != "public" && blogStatus != "password" && blogStatus != "draft" {
			blogStatus = "public"
		}
		if blogStatus == "password" && blogPassword == "" {
			c.JsonResult(6010, i18n.Tr(c.Lang, "message.set_pwd_pls"))
		}
		if blogType != 0 && blogType != 1 {
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.unknown_blog_type"))
		}
		if strings.Count(blogTitle, "") > 200 {
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.blog_title_tips"))
		}
		//如果是关联文章，需要同步关联的文档
		if blogType == 1 {
			book, err := models.NewBook().FindByIdentify(bookIdentify)

			if err != nil {
				c.JsonResult(6011, i18n.Tr(c.Lang, "message.ref_doc_not_exist_or_no_permit"))
			}
			doc, err := models.NewDocument().FindByIdentityFirst(documentIdentify, book.BookId)
			if err != nil {
				c.JsonResult(6003, i18n.Tr(c.Lang, "message.query_failed"))
			}
			documentId = doc.DocumentId

			// 如果不是超级管理员，则校验权限
			if !c.Member.IsAdministrator() {
				bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)

				if err != nil || bookResult.RoleId == conf.BookObserver {
					c.JsonResult(6002, i18n.Tr(c.Lang, "message.ref_doc_not_exist_or_no_permit"))
				}
			}
		}

		var blog *models.Blog
		var err error
		//如果文章ID存在，则从数据库中查询文章
		if blogId > 0 {
			if c.Member.IsAdministrator() {
				blog, err = models.NewBlog().Find(blogId)
			} else {
				blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
			}

			if err != nil {
				c.JsonResult(6003, i18n.Tr(c.Lang, "message.blog_not_exist"))
			}
			//如果设置了文章标识
			if blogIdentify != "" {
				//如果查询到的文章标识存在并且不是当前文章的id
				if b, err := models.NewBlog().FindByIdentify(blogIdentify); err == nil && b.BlogId != blogId {
					c.JsonResult(6004, i18n.Tr(c.Lang, "message.blog_id_existed"))
				}
			}
			blog.Modified = time.Now()
			blog.ModifyAt = c.Member.MemberId
		} else {
			//如果设置了文章标识
			if blogIdentify != "" {
				if models.NewBlog().IsExist(blogIdentify) {
					c.JsonResult(6004, i18n.Tr(c.Lang, "message.blog_id_existed"))
				}
			}

			blog = models.NewBlog()
			blog.MemberId = c.Member.MemberId
			blog.Created = time.Now()
		}
		if blogIdentify == "" {
			blog.BlogIdentify = fmt.Sprintf("%s-%d", "post", time.Now().UnixNano())
		} else {
			blog.BlogIdentify = blogIdentify
		}

		blog.BlogTitle = blogTitle

		blog.OrderIndex = orderIndex
		blog.BlogType = blogType
		if blogType == 1 {
			blog.DocumentId = documentId
		}
		blog.BlogExcerpt = blogExcerpt
		blog.BlogStatus = blogStatus
		blog.Password = blogPassword

		if err := blog.Save(); err != nil {
			logs.Error("保存文章失败 -> ", err)
			c.JsonResult(6011, i18n.Tr(c.Lang, "message.failed"))
		} else {
			c.JsonResult(0, "ok", blog)
		}
	}
	if c.Ctx.Input.Referer() == "" {
		c.Data["Referer"] = "javascript:history.back();"
	} else {
		c.Data["Referer"] = c.Ctx.Input.Referer()
	}
	blogId, err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	c.Data["DocumentIdentify"] = ""
	if err == nil {
		blog, err := models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
		if err != nil {
			c.ShowErrorPage(500, err.Error())
		}

		c.Data["Model"] = blog
	} else {
		c.Data["Model"] = models.NewBlog()
	}
}

// 文章创建或编辑
func (c *BlogController) ManageEdit() {
	c.Prepare()
	c.TplName = "blog/manage_edit.tpl"

	if c.Member.Role == conf.MemberReaderRole {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.no_permission"))
	}

	if c.Ctx.Input.IsPost() {
		blogId, _ := c.GetInt("blogId", 0)

		if blogId <= 0 {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
		}
		blogContent := c.GetString("content", "")
		blogHtml := c.GetString("htmlContent", "")
		version, _ := c.GetInt64("version", 0)
		cover := c.GetString("cover")

		var blog *models.Blog
		var err error

		if c.Member.IsAdministrator() {
			blog, err = models.NewBlog().Find(blogId)
		} else {
			blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
		}
		if err != nil {
			logs.Error("查询文章失败 ->", err)
			c.JsonResult(6002, i18n.Tr(c.Lang, "message.query_failed"))
		}
		if version > 0 && blog.Version != version && cover != "yes" {
			c.JsonResult(6005, i18n.Tr(c.Lang, "message.blog_has_modified"))
		}
		//如果是关联文章，需要同步关联的文档
		if blog.BlogType == 1 {
			doc, err := models.NewDocument().Find(blog.DocumentId)
			if err != nil {
				logs.Error("查询关联项目文档时出错 ->", err)
				c.JsonResult(6003, i18n.Tr(c.Lang, "message.query_failed"))
			}
			book, err := models.NewBook().Find(doc.BookId)
			if err != nil {
				c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist_or_no_permit"))
			}

			// 如果不是超级管理员，则校验权限
			if !c.Member.IsAdministrator() {
				bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)

				if err != nil || bookResult.RoleId == conf.BookObserver {
					logs.Error("FindByIdentify => ", err)
					c.JsonResult(6002, i18n.Tr(c.Lang, "message.ref_doc_not_exist_or_no_permit"))
				}
			}

			doc.Markdown = blogContent
			doc.Release = blogHtml
			doc.Content = blogHtml
			doc.ModifyTime = time.Now()
			doc.ModifyAt = c.Member.MemberId
			if err := doc.InsertOrUpdate("markdown", "release", "content", "modify_time", "modify_at"); err != nil {
				logs.Error("保存关联文档时出错 ->", err)
				c.JsonResult(6004, i18n.Tr(c.Lang, "message.failed"))
			}
		}

		blog.BlogContent = blogContent
		blog.BlogRelease = blogHtml
		blog.ModifyAt = c.Member.MemberId
		blog.Modified = time.Now()

		if err := blog.Save("blog_content", "blog_release", "modify_at", "modify_time", "version"); err != nil {
			logs.Error("保存文章失败 -> ", err)
			c.JsonResult(6011, i18n.Tr(c.Lang, "message.failed"))
		} else {
			c.JsonResult(0, "ok", blog)
		}
	}

	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if blogId <= 0 {
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.param_error"))
	}
	var blog *models.Blog
	var err error

	if c.Member.IsAdministrator() {
		blog, err = models.NewBlog().Find(blogId)
	} else {
		blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
	}
	if err != nil {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.blog_not_exist"))
	}
	blog.LinkAttach()

	if len(blog.AttachList) > 0 {
		returnJSON, err := json.Marshal(blog.AttachList)
		if err != nil {
			logs.Error("序列化文章附件时出错 ->", err)
		} else {
			c.Data["AttachList"] = template.JS(string(returnJSON))
		}
	} else {
		c.Data["AttachList"] = template.JS("[]")
	}
	if conf.GetUploadFileSize() > 0 {
		c.Data["UploadFileSize"] = conf.GetUploadFileSize()
	} else {
		c.Data["UploadFileSize"] = "undefined"
	}
	c.Data["Model"] = blog
}

// 删除文章
func (c *BlogController) ManageDelete() {
	c.Prepare()
	blogId, _ := c.GetInt("blog_id", 0)

	if blogId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	var blog *models.Blog
	var err error

	if c.Member.IsAdministrator() {
		blog, err = models.NewBlog().Find(blogId)
	} else {
		blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
	}
	if err != nil {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.blog_not_exist"))
	}

	if err := blog.Delete(blogId); err != nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
	} else {
		c.JsonResult(0, i18n.Tr(c.Lang, "message.success"))
	}

}

// 上传附件或图片
func (c *BlogController) Upload() {
	c.Prepare()
	blogId, _ := c.GetInt("blogId")

	if blogId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	blog, err := models.NewBlog().Find(blogId)

	if err != nil {
		c.JsonResult(6010, i18n.Tr(c.Lang, "message.blog_not_exist"))
	}
	if !c.Member.IsAdministrator() && blog.MemberId != c.Member.MemberId {
		c.JsonResult(6011, i18n.Tr(c.Lang, "message.no_permission"))
	}

	name := "editormd-file-file"

	file, moreFile, err := c.GetFile(name)
	if err == http.ErrMissingFile {
		name = "editormd-image-file"
		file, moreFile, err = c.GetFile(name)
		if err == http.ErrMissingFile {
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.upload_file_empty"))
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
		c.JsonResult(6009, i18n.Tr(c.Lang, "message.upload_file_size_limit"))
	}

	ext := filepath.Ext(moreFile.Filename)

	if ext == "" {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.upload_file_type_error"))
	}
	//如果文件类型设置为 * 标识不限制文件类型
	if web.AppConfig.DefaultString("upload_file_ext", "") != "*" {
		if !conf.IsAllowUploadFileExt(ext) {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.upload_file_type_error"))
		}
	}

	// 如果是超级管理员，则不判断权限
	if c.Member.IsAdministrator() {
		_, err := models.NewBlog().Find(blogId)

		if err != nil {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.doc_not_exist_or_no_permit"))
		}

	} else {
		_, err := models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)

		if err != nil {
			logs.Error("查询文章时出错 -> ", err)
			if err == orm.ErrNoRows {
				c.JsonResult(6006, i18n.Tr(c.Lang, "message.no_permission"))
			}

			c.JsonResult(6001, err.Error())
		}
	}

	fileName := "attach_" + strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := filepath.Join(conf.WorkingDirectory, "uploads", "blog", time.Now().Format("200601"), fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile(name, filePath)

	if err != nil {
		logs.Error("SaveToFile => ", err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
	}

	var httpPath string
	result := make(map[string]interface{})
	//如果是图片，则当做内置图片处理，否则当做附件处理
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") || strings.EqualFold(ext, ".png") || strings.EqualFold(ext, ".gif") {
		httpPath = "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)
		if strings.HasPrefix(httpPath, "//") {
			httpPath = conf.URLForWithCdnImage(string(httpPath[1:]))
		}
	} else {
		attachment := models.NewAttachment()
		attachment.BookId = 0
		attachment.FileName = moreFile.Filename
		attachment.CreateAt = c.Member.MemberId
		attachment.FileExt = ext
		attachment.FilePath = strings.TrimPrefix(filePath, conf.WorkingDirectory)
		attachment.DocumentId = blogId
		//如果是关联文章，则将附件设置为关联文档的文档上
		if blog.BlogType == 1 {
			attachment.BookId = blog.BookId
			attachment.DocumentId = blog.DocumentId
		}
		if fileInfo, err := os.Stat(filePath); err == nil {
			attachment.FileSize = float64(fileInfo.Size())
		}

		attachment.HttpPath = httpPath

		if err := attachment.Insert(); err != nil {
			os.Remove(filePath)
			logs.Error("保存文件附件失败 -> ", err)
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed"))
		}
		if attachment.HttpPath == "" {
			attachment.HttpPath = conf.URLForNotHost("BlogController.Download", ":id", blogId, ":attach_id", attachment.AttachmentId)

			if err := attachment.Update(); err != nil {
				logs.Error("保存文件失败 -> ", attachment.FilePath, err)
				c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
			}
		}
		result["attach"] = attachment
	}

	result["errcode"] = 0
	result["success"] = 1
	result["message"] = "ok"
	result["url"] = httpPath
	result["alt"] = fileName

	c.Ctx.Output.JSON(result, true, false)
	c.StopRun()
}

// 删除附件
func (c *BlogController) RemoveAttachment() {
	c.Prepare()
	attachId, _ := c.GetInt("attach_id")
	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if attachId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	blog, err := models.NewBlog().Find(blogId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.doc_not_exist"))
		} else {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.query_failed"))
		}
	}
	attach, err := models.NewAttachment().Find(attachId)

	if err != nil {
		logs.Error(err)
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	}

	if !c.Member.IsAdministrator() {
		_, err := models.NewBlog().FindByIdAndMemberId(attach.DocumentId, c.Member.MemberId)

		if err != nil {
			logs.Error(err)
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.doc_not_exist"))
		}
	}

	if blog.BlogType == 1 && attach.BookId != blog.BookId && attach.DocumentId != blog.DocumentId {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	} else if attach.BookId != 0 || attach.DocumentId != blogId {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	}

	if err := attach.Delete(); err != nil {
		logs.Error(err)
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.failed"))
	}

	os.Remove(filepath.Join(conf.WorkingDirectory, attach.FilePath))

	c.JsonResult(0, "ok", attach)
}

// 下载附件
func (c *BlogController) Download() {
	c.Prepare()

	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	attachId, _ := strconv.Atoi(c.Ctx.Input.Param(":attach_id"))
	password := c.GetString("password")

	blog, err := models.NewBlog().Find(blogId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.doc_not_exist"))
		} else {
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.query_failed"))
		}
	}
	blogReadSession := fmt.Sprintf("blog:read:%d", blogId)
	//如果没有启动匿名访问，或者设置了访问密码
	if (c.Member == nil && !c.EnableAnonymous) || (blog.BlogStatus == "password" && password != blog.Password && c.CruSession.Get(context.TODO(), blogReadSession) == nil) {
		c.ShowErrorPage(403, i18n.Tr(c.Lang, "message.no_permission"))
	}

	// 查找附件
	attachment, err := models.NewAttachment().Find(attachId)

	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
		} else {
			logs.Error("查询附件时出现异常 -> ", err)
			c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.query_failed"))
		}
	}

	//如果是链接的文章，需要校验文档ID是否一致，如果不是，需要保证附件的项目ID为0且文档的ID等于博文ID
	if blog.BlogType == 1 && attachment.DocumentId != blog.DocumentId {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	} else if blog.BlogType != 1 && (attachment.BookId != 0 || attachment.DocumentId != blogId) {
		c.ShowErrorPage(404, i18n.Tr(c.Lang, "message.attachment_not_exist"))
	}

	c.Ctx.Output.Download(filepath.Join(conf.WorkingDirectory, attachment.FilePath), attachment.FileName)
	c.StopRun()
}
