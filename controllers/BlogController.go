package controllers

import (
	"strings"
	"github.com/lifei6671/mindoc/models"
	"time"
	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/utils/pagination"
	"strconv"
	"fmt"
	"os"
	"net/http"
	"path/filepath"
	"github.com/astaxie/beego/orm"
	"html/template"
	"encoding/json"
	"github.com/lifei6671/mindoc/utils"
	"net/url"
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

//文章阅读
func (c *BlogController) Index() {
	c.Prepare()
	c.TplName = "blog/index.tpl"
	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if blogId <= 0 {
		c.ShowErrorPage(404, "页面不存在")
	}
	blogReadSession := fmt.Sprintf("blog:read:%d", blogId)

	blog, err := models.NewBlog().FindFromCache(blogId)

	if err != nil {
		c.ShowErrorPage(404, "文章不存在")
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password");
		if blog.BlogStatus == "password" && password != blog.Password {
			c.JsonResult(6001, "文章密码不正确")
		} else if blog.BlogStatus == "password" && password == blog.Password {
			//如果密码输入正确，则存入session中
			c.CruSession.Set(blogReadSession, blogId)
			c.JsonResult(0, "OK")
		}
		c.JsonResult(0, "OK")
	} else if blog.BlogStatus == "password" && (c.CruSession.Get(blogReadSession) == nil || (c.Member != nil && blog.MemberId != c.Member.MemberId && !c.Member.IsAdministrator())) {
		//如果不存在已输入密码的标记
		c.TplName = "blog/index_password.tpl"
	}
	//加载文章附件
	blog.LinkAttach();
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

//文章列表
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

//管理后台文章列表
func (c *BlogController) ManageList() {
	c.Prepare()
	c.TplName = "blog/manage_list.tpl"

	pageIndex, _ := c.GetInt("page", 1)

	blogList, totalCount, err := models.NewBlog().FindToPager(pageIndex, conf.PageSize, c.Member.MemberId, "")

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

//文章设置
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

		if blogTitle == "" {
			c.JsonResult(6001, "文章标题不能为空")
		}
		if strings.Count(blogExcerpt, "") > 500 {
			c.JsonResult(6008, "文章摘要必须小于500字符")
		}
		if blogStatus != "public" && blogStatus != "password" && blogStatus != "draft" {
			blogStatus = "public"
		}
		if blogStatus == "password" && blogPassword == "" {
			c.JsonResult(6010, "加密文章请设置密码")
		}
		if blogType != 0 && blogType != 1 {
			c.JsonResult(6005, "未知的文章类型")
		}
		if strings.Count(blogTitle, "") > 200 {
			c.JsonResult(6002, "文章标题不能大于200个字符")
		}
		//如果是关联文章，需要同步关联的文档
		if blogType == 1 {
			book, err := models.NewBook().FindByIdentify(bookIdentify)

			if err != nil {
				c.JsonResult(6011, "关联文档的项目不存在")
			}
			doc, err := models.NewDocument().FindByIdentityFirst(documentIdentify, book.BookId)
			if err != nil {
				c.JsonResult(6003, "查询关联项目文档时出错")
			}
			documentId = doc.DocumentId

			// 如果不是超级管理员，则校验权限
			if !c.Member.IsAdministrator() {
				bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)

				if err != nil || bookResult.RoleId == conf.BookObserver {
					c.JsonResult(6002, "关联文档不存在或权限不足")
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
				c.JsonResult(6003, "文章不存在")
			}
			//如果设置了文章标识
			if blogIdentify != "" {
				//如果查询到的文章标识存在并且不是当前文章的id
				if b, err := models.NewBlog().FindByIdentify(blogIdentify); err == nil && b.BlogId != blogId {
					c.JsonResult(6004, "文章标识已存在")
				}
			}
			blog.Modified = time.Now()
			blog.ModifyAt = c.Member.MemberId
		} else {
			//如果设置了文章标识
			if blogIdentify != "" {
				if models.NewBlog().IsExist(blogIdentify) {
					c.JsonResult(6004, "文章标识已存在")
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
			beego.Error("保存文章失败 -> ", err)
			c.JsonResult(6011, "保存文章失败")
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

	c.Data["DocumentIdentify"] = "";
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

//文章创建或编辑
func (c *BlogController) ManageEdit() {
	c.Prepare()
	c.TplName = "blog/manage_edit.tpl"

	if c.Ctx.Input.IsPost() {
		blogId, _ := c.GetInt("blogId", 0)

		if blogId <= 0 {
			c.JsonResult(6001, "文章参数错误")
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
			beego.Error("查询文章失败 ->", err)
			c.JsonResult(6002, "查询文章失败")
		}
		if version > 0 && blog.Version != version && cover != "yes" {
			c.JsonResult(6005, "文章已被修改")
		}
		//如果是关联文章，需要同步关联的文档
		if blog.BlogType == 1 {
			doc, err := models.NewDocument().Find(blog.DocumentId)
			if err != nil {
				beego.Error("查询关联项目文档时出错 ->", err)
				c.JsonResult(6003, "查询关联项目文档时出错")
			}
			book, err := models.NewBook().Find(doc.BookId)
			if err != nil {
				c.JsonResult(6002, "项目不存在或权限不足")
			}

			// 如果不是超级管理员，则校验权限
			if !c.Member.IsAdministrator() {
				bookResult, err := models.NewBookResult().FindByIdentify(book.Identify, c.Member.MemberId)

				if err != nil || bookResult.RoleId == conf.BookObserver {
					beego.Error("FindByIdentify => ", err)
					c.JsonResult(6002, "关联文档不存在或权限不足")
				}
			}

			doc.Markdown = blogContent
			doc.Release = blogHtml
			doc.Content = blogHtml
			doc.ModifyTime = time.Now()
			doc.ModifyAt = c.Member.MemberId
			if err := doc.InsertOrUpdate("markdown", "release", "content", "modify_time", "modify_at"); err != nil {
				beego.Error("保存关联文档时出错 ->", err)
				c.JsonResult(6004, "保存关联文档时出错")
			}
		}

		blog.BlogContent = blogContent
		blog.BlogRelease = blogHtml
		blog.ModifyAt = c.Member.MemberId
		blog.Modified = time.Now()

		if err := blog.Save("blog_content", "blog_release", "modify_at", "modify_time", "version"); err != nil {
			beego.Error("保存文章失败 -> ", err)
			c.JsonResult(6011, "保存文章失败")
		} else {
			c.JsonResult(0, "ok", blog)
		}
	}

	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if blogId <= 0 {
		c.ShowErrorPage(500, "参数错误")
	}
	var blog *models.Blog
	var err error

	if c.Member.IsAdministrator() {
		blog, err = models.NewBlog().Find(blogId)
	} else {
		blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
	}
	if err != nil {
		c.ShowErrorPage(404, "文章不存在或已删除")
	}
	blog.LinkAttach()

	if len(blog.AttachList) > 0 {
		returnJSON, err := json.Marshal(blog.AttachList)
		if err != nil {
			beego.Error("序列化文章附件时出错 ->", err)
		} else {
			c.Data["AttachList"] = template.JS(string(returnJSON))
		}
	} else {
		c.Data["AttachList"] = template.JS("[]")
	}
	if conf.GetUploadFileSize() > 0 {
		c.Data["UploadFileSize"] = conf.GetUploadFileSize()
	} else {
		c.Data["UploadFileSize"] = "undefined";
	}
	c.Data["Model"] = blog
}

//删除文章
func (c *BlogController) ManageDelete() {
	c.Prepare()
	blogId, _ := c.GetInt("blog_id", 0)

	if blogId <= 0 {
		c.JsonResult(6001, "参数错误")
	}
	var blog *models.Blog
	var err error

	if c.Member.IsAdministrator() {
		blog, err = models.NewBlog().Find(blogId)
	} else {
		blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
	}
	if err != nil {
		c.JsonResult(6002, "文章不存在或已删除")
	}

	if err := blog.Delete(blogId); err != nil {
		c.JsonResult(6003, "删除失败")
	} else {
		c.JsonResult(0, "删除成功")
	}

}

// 上传附件或图片
func (c *BlogController) Upload() {
	c.Prepare()
	blogId, _ := c.GetInt("blogId")

	if blogId <= 0 {
		c.JsonResult(6001, "参数错误")
	}

	blog, err := models.NewBlog().Find(blogId)

	if err != nil {
		c.JsonResult(6010, "文章不存在")
	}
	if !c.Member.IsAdministrator() && blog.MemberId != c.Member.MemberId {
		c.JsonResult(6011, "没有文章的访问权限")
	}

	name := "editormd-file-file"

	file, moreFile, err := c.GetFile(name)
	if err == http.ErrMissingFile {
		name = "editormd-image-file"
		file, moreFile, err = c.GetFile(name)
		if err == http.ErrMissingFile {
			c.JsonResult(6003, "没有发现需要上传的图片")
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

	// 如果是超级管理员，则不判断权限
	if c.Member.IsAdministrator() {
		_, err := models.NewBlog().Find(blogId)

		if err != nil {
			c.JsonResult(6006, "文档不存在或权限不足")
		}

	} else {
		_, err := models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)

		if err != nil {
			beego.Error("查询文章时出错 -> ", err)
			if err == orm.ErrNoRows {
				c.JsonResult(6006, "权限不足")
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
		beego.Error("SaveToFile => ", err)
		c.JsonResult(6005, "保存文件失败")
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
			beego.Error("保存文件附件失败 => ", err)
			c.JsonResult(6006, "文件保存失败")
		}
		if attachment.HttpPath == "" {
			attachment.HttpPath = conf.URLFor("BlogController.Download", ":id", blogId, ":attach_id", attachment.AttachmentId)

			if err := attachment.Update(); err != nil {
				beego.Error("SaveToFile => ", err)
				c.JsonResult(6005, "保存文件失败")
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
		c.JsonResult(6001, "参数错误")
	}
	blog, err := models.NewBlog().Find(blogId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(500, "文档不存在")
		} else {
			c.ShowErrorPage(500, "查询文章时异常")
		}
	}
	attach, err := models.NewAttachment().Find(attachId)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6002, "附件不存在")
	}

	if !c.Member.IsAdministrator() {
		_, err := models.NewBlog().FindByIdAndMemberId(attach.DocumentId, c.Member.MemberId)

		if err != nil {
			beego.Error(err)
			c.JsonResult(6003, "文档不存在")
		}
	}

	if blog.BlogType == 1 && attach.BookId != blog.BookId && attach.DocumentId != blog.DocumentId {
		c.ShowErrorPage(404, "附件不存在")
	} else if attach.BookId != 0 || attach.DocumentId != blogId {
		c.ShowErrorPage(404, "附件不存在")
	}

	if err := attach.Delete(); err != nil {
		beego.Error(err)
		c.JsonResult(6005, "删除失败")
	}

	os.Remove(filepath.Join(conf.WorkingDirectory, attach.FilePath))

	c.JsonResult(0, "ok", attach)
}

//下载附件
func (c *BlogController) Download() {
	c.Prepare()

	blogId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	attachId, _ := strconv.Atoi(c.Ctx.Input.Param(":attach_id"))
	password := c.GetString("password")

	blog, err := models.NewBlog().Find(blogId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.ShowErrorPage(500, "文档不存在")
		} else {
			c.ShowErrorPage(500, "查询文章时异常")
		}
	}
	blogReadSession := fmt.Sprintf("blog:read:%d", blogId)
	//如果没有启动匿名访问，或者设置了访问密码
	if (c.Member == nil && !c.EnableAnonymous) || (blog.BlogStatus == "password" && password != blog.Password && c.CruSession.Get(blogReadSession) == nil) {
		c.ShowErrorPage(403, "没有下载权限")
	}

	// 查找附件
	attachment, err := models.NewAttachment().Find(attachId)

	if err != nil {
		beego.Error("DownloadAttachment => ", err)
		if err == orm.ErrNoRows {
			c.ShowErrorPage(404, "附件不存在")
		} else {
			c.ShowErrorPage(500, "查询附件时出现异常")
		}
	}
	if blog.BlogType == 1 && attachment.BookId != blog.BookId && attachment.DocumentId != blog.DocumentId {
		c.ShowErrorPage(404, "附件不存在")
	} else if attachment.BookId != 0 || attachment.DocumentId != blogId {
		c.ShowErrorPage(404, "附件不存在")
	}

	c.Ctx.Output.Download(filepath.Join(conf.WorkingDirectory, attachment.FilePath), attachment.FileName)
	c.StopRun()
}
