package controllers

import (
	"os"
	"time"
	"regexp"
	"strconv"
	"strings"
	"net/http"
	"path/filepath"
	"encoding/json"
	"html/template"

	"github.com/lifei6671/godoc/models"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DocumentController struct {
	BaseController
}

func (p *DocumentController) Index()  {
	p.TplName = "document/index.tpl"
}

func (p *DocumentController) Read() {
	p.TplName = "document/kancloud.tpl"
}

func (c *DocumentController) Edit()  {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		logs.Error("DocumentController.Edit => ",err)

		c.Abort("403")
	}
	if book.Editor == "markdown" {
		c.TplName = "document/markdown_edit_template.tpl"
	}else{
		c.TplName = "document/html_edit_template.tpl"
	}

	c.Data["Model"] = book

	r,_ := json.Marshal(book)

	c.Data["ModelResult"] = template.JS(string(r))

	c.Data["Result"] = template.JS("[]")

	trees ,err := models.NewDocument().FindDocumentTree(book.BookId)
	logs.Info("",trees)
	if err != nil {
		logs.Error("FindDocumentTree => ", err)
	}else{
		if len(trees) > 0 {
			if jtree, err := json.Marshal(trees); err == nil {
				c.Data["Result"] = template.JS(string(jtree))
			}
		}else{
			c.Data["Result"] = template.JS("[]")
		}
	}

}

//创建一个文档
func (c *DocumentController) Create() {
	identify := c.GetString("identify")
	doc_identify := c.GetString("doc_identify")
	doc_name := c.GetString("doc_name")
	parent_id,_ := c.GetInt("parent_id",0)
	doc_id,_ := c.GetInt("doc_id",0)

	if identify == "" {
		c.JsonResult(6001,"参数错误")
	}
	if doc_name == "" {
		c.JsonResult(6004,"文档名称不能为空")
	}
	if doc_identify != "" {
		if ok, err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`, doc_identify); !ok || err != nil {

			c.JsonResult(6003, "文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头")
		}
		d,_ := models.NewDocument().FindByFieldFirst("identify",doc_identify);
		if  d.DocumentId > 0 && d.DocumentId != doc_id{
			c.JsonResult(6006,"文档标识已被使用")
		}
	}

	bookResult,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil || bookResult.RoleId == conf.BookObserver {
		logs.Error("FindByIdentify => ",err)
		c.JsonResult(6002,"项目不存在或权限不足")
	}
	if parent_id > 0 {
		doc,err := models.NewDocument().Find(parent_id)
		if err != nil || doc.BookId != bookResult.BookId{
			c.JsonResult(6003,"父分类不存在")
		}
	}

	document,_ := models.NewDocument().Find(doc_id)

	document.MemberId = c.Member.MemberId
	document.BookId = bookResult.BookId
	if doc_identify != ""{
		document.Identify = doc_identify
	}
	document.Version = time.Now().Unix()
	document.DocumentName = doc_name
	document.ParentId = parent_id

	if err := document.InsertOrUpdate();err != nil {
		logs.Error("InsertOrUpdate => ",err)
		c.JsonResult(6005,"保存失败")
	}else{
		logs.Info("",document)
		c.JsonResult(0,"ok",document)
	}
}

//上传附件或图片
func (c *DocumentController) Upload()  {

	identify := c.GetString("identify")

	if identify == "" {
		c.JsonResult(6001,"参数错误")
	}


	name := "editormd-file-file"

	file,moreFile,err  := c.GetFile(name)
	if err == http.ErrMissingFile {
		name = "editormd-image-file"
		file,moreFile,err = c.GetFile(name);
		if err == http.ErrMissingFile {
			c.JsonResult(6003,"没有发现需要上传的文件")
		}
	}
	if err != nil {
		c.JsonResult(6002,err.Error())
	}

	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)

	if ext == "" {
		c.JsonResult(6003,"无法解析文件的格式")
	}

	if !conf.IsAllowUploadFileExt(ext) {
		c.JsonResult(6004,"不允许的文件类型")
	}

	book,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		logs.Error("DocumentController.Edit => ",err)
		if err == orm.ErrNoRows {
			c.JsonResult(6006,"权限不足")
		}
		c.JsonResult(6001,err.Error())
	}
	//如果没有编辑权限
	if book.RoleId != conf.BookEditor && book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
		c.JsonResult(6006,"权限不足")
	}

	fileName := "attachment_" +  strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := "uploads/" + time.Now().Format("200601") + "/" + fileName + ext

	err = c.SaveToFile(name,filePath)

	if err != nil {
		logs.Error("SaveToFile => ",err)
		c.JsonResult(6005,"保存文件失败")
	}
	attachment := models.NewAttachment()
	attachment.BookId = book.BookId
	attachment.FileName = moreFile.Filename
	attachment.CreateAt = c.Member.MemberId
	attachment.FileExt = ext
	attachment.FilePath = filePath

	if strings.EqualFold(ext,".jpg") || strings.EqualFold(ext,".jpeg") || strings.EqualFold(ext,"png") || strings.EqualFold(ext,"gif") {
		attachment.HttpPath = c.BaseUrl() + "/" + filePath
	}

	err = attachment.Insert();

	if err != nil {
		os.Remove(filePath)
		logs.Error("Attachment Insert => ",err)
		c.JsonResult(6006,"文件保存失败")
	}
	if attachment.HttpPath == "" {
		attachment.HttpPath = c.BaseUrl() + beego.URLFor("DocumentController.DownloadAttachment",":key", identify, ":attach_id", attachment.AttachmentId)

		if err := attachment.Update();err != nil {
			logs.Error("SaveToFile => ",err)
			c.JsonResult(6005,"保存文件失败")
		}
	}
	result := map[string]interface{}{
		"success" : 1,
		"message" :"ok",
		"url" : attachment.HttpPath,
		"alt" : attachment.FileName,
	}

	c.Data["json"] = result
	c.ServeJSON(true)
	c.StopRun()
}

//DownloadAttachment 下载附件.
func (c *DocumentController) DownloadAttachment()  {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	attach_id,_ := strconv.Atoi(c.Ctx.Input.Param(":attach_id"))
	token := c.GetString("token")

	member_id := 0

	if c.Member != nil {
		member_id = c.Member.MemberId
	}
	book_id := 0

	//判断用户是否参与了项目
	bookResult,err := models.NewBookResult().FindByIdentify(identify,member_id)

	if err != nil {
		//判断项目公开状态
		book,err := models.NewBook().FindByFieldFirst("identify",identify)
		if err != nil {
			c.Abort("404")
		}
		//如果项目是私有的，并且token不正确
		if (book.PrivatelyOwned == 1 && token == "" ) || ( book.PrivatelyOwned == 1 && book.PrivateToken != token ){
			c.Abort("403")
		}
		book_id = book.BookId
	}else{
		book_id = bookResult.BookId
	}

	attachment,err := models.NewAttachment().Find(attach_id)

	if err != nil {
		logs.Error("DownloadAttachment => ", err)
		if err == orm.ErrNoRows {
			c.Abort("404")
		} else {
			c.Abort("500")
		}
	}
	if attachment.BookId != book_id {
		c.Abort("404")
	}
	c.Ctx.Output.Download(attachment.FilePath,attachment.FileName)

	c.StopRun()
}

func (c *DocumentController) Delete() {
	c.Prepare()

	identify := c.GetString("identify")
	doc_id,err := c.GetInt("doc_id",0)

	bookResult,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil || bookResult.RoleId == conf.BookObserver {
		logs.Error("FindByIdentify => ",err)
		c.JsonResult(6002,"项目不存在或权限不足")
	}

	if doc_id <= 0 {
		c.JsonResult(6001,"参数错误")
	}

	doc,err := models.NewDocument().Find(doc_id)

	if err != nil {
		logs.Error("Delete => ",err)
		c.JsonResult(6003,"删除失败")
	}
	if doc.BookId != bookResult.BookId {
		c.JsonResult(6004,"参数错误")
	}
	err = doc.RecursiveDocument(doc.DocumentId)
	if err != nil {
		c.JsonResult(6005,"删除失败")
	}

	c.JsonResult(0,"ok")
}

func (c *DocumentController) Content()  {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	doc_id,err := c.GetInt("doc_id")

	if err != nil {
		doc_id,_ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	}

	bookResult,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil || bookResult.RoleId == conf.BookObserver {
		logs.Error("FindByIdentify => ",err)
		c.JsonResult(6002,"项目不存在或权限不足")
	}

	if doc_id <= 0 {
		c.JsonResult(6001,"参数错误")
	}

	if c.Ctx.Input.IsPost() {
		markdown := strings.TrimSpace(c.GetString("markdown",""))
		content := c.GetString("html")
		version,_ := c.GetInt64("version",0)
		is_cover := c.GetString("cover")

		doc ,err := models.NewDocument().Find(doc_id);

		if err != nil {
			c.JsonResult(6003,"读取文档错误")
		}
		if doc.BookId != bookResult.BookId {
			c.JsonResult(6004,"保存的文档不属于指定项目")
		}
		if doc.Version != version && !strings.EqualFold(is_cover,"yes"){
			logs.Info("%d|",version,doc.Version)
			c.JsonResult(6005,"文档已被修改确定要覆盖吗？")
		}
		if markdown == "" && content != ""{
			doc.Markdown = content
		}else{
			doc.Markdown = markdown
		}
		doc.Version = time.Now().Unix()
		doc.Content = content
		if err := doc.InsertOrUpdate();err != nil {
			logs.Error("InsertOrUpdate => ",err)
			c.JsonResult(6006,"保存失败")
		}

		c.JsonResult(0,"ok",doc)
	}
	doc,err := models.NewDocument().Find(doc_id)

	if err != nil {
		c.JsonResult(6003,"文档不存在")
	}
	c.JsonResult(0,"ok",doc)
}


















