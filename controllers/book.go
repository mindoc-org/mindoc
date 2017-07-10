package controllers

import (
	"strings"
	"regexp"
	"strconv"
	"time"
	"encoding/json"
	"html/template"
	"errors"
	"fmt"
	"path/filepath"
	"os"

	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/graphics"
	"github.com/lifei6671/mindoc/commands"
)

type BookController struct {
	BaseController
}

func (c *BookController) Index() {
	c.Prepare()
	c.TplName = "book/index.tpl"

	pageIndex, _ := c.GetInt("page", 1)

	books,totalCount,err := models.NewBook().FindToPager(pageIndex,conf.PageSize,c.Member.MemberId)

	if err != nil {
		logs.Error("BookController.Index => ",err)
		c.Abort("500")
	}

	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, conf.PageSize, totalCount)

		c.Data["PageHtml"] = html
	}else {
		c.Data["PageHtml"] = ""
	}
	b,err := json.Marshal(books)

	if err != nil || len(books) <= 0{
		c.Data["Result"] = template.JS("[]")
	}else{
		c.Data["Result"] = template.JS(string(b))
	}
}

// Dashboard 项目概要 .
func (c *BookController) Dashboard() {
	c.Prepare()
	c.TplName = "book/dashboard.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		beego.Error(err)
		c.Abort("500")
	}

	c.Data["Model"] = *book
}

// Setting 项目设置 .
func (c *BookController) Setting()  {
	c.Prepare()
	c.TplName = "book/setting.tpl"

	key := c.Ctx.Input.Param(":key")

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Abort("404")
		}
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
	}
	//如果不是创始人也不是管理员则不能操作
	if book.RoleId != conf.BookFounder && book.RoleId != conf.BookAdmin {
		c.Abort("403")
	}
	if book.PrivateToken != "" {
		book.PrivateToken = c.BaseUrl() + beego.URLFor("DocumentController.Index",":key",book.Identify,"token",book.PrivateToken)
	}
	c.Data["Model"] = book

}

//保存项目信息
func (c *BookController) SaveBook()  {
	bookResult,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001,err.Error())
	}
	book,err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		logs.Error("SaveBook => ",err)
		c.JsonResult(6002,err.Error())
	}

	book_name := strings.TrimSpace(c.GetString("book_name"))
	description := strings.TrimSpace(c.GetString("description",""))
	comment_status := c.GetString("comment_status")
	tag := strings.TrimSpace(c.GetString("label"))
	editor := strings.TrimSpace(c.GetString("editor"))

	if strings.Count(description,"") > 500 {
		c.JsonResult(6004,"项目描述不能大于500字")
	}
	if comment_status != "open" && comment_status != "closed" && comment_status != "group_only" && comment_status != "registered_only" {
		comment_status = "closed"
	}
	if tag != ""{
		tags := strings.Split(tag,",")
		if len(tags) > 10 {
			c.JsonResult(6005,"最多允许添加10个标签")
		}
	}
	if editor != "markdown" && editor != "html" {
		editor = "markdown"
	}

	book.BookName 		= book_name
	book.Description 	= description
	book.CommentStatus 	= comment_status
	book.Label 		= tag
	book.Editor 		= editor

	if err := book.Update();err != nil {
		c.JsonResult(6006,"保存失败")
	}
	bookResult.BookName = book_name
	bookResult.Description = description
	bookResult.CommentStatus = comment_status
	bookResult.Label = tag
	c.JsonResult(0,"ok",bookResult)
}

//设置项目私有状态.
func (c *BookController) PrivatelyOwned()  {

	status := c.GetString("status")

	if status != "open" && status != "close" {
		c.JsonResult(6003,"参数错误")
	}
	state := 0
	if status == "open" {
		state = 0
	}else{
		state = 1
	}

	bookResult,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001,err.Error())
	}
	//只有创始人才能变更私有状态
	if bookResult.RoleId != conf.BookFounder {
		c.JsonResult(6002,"权限不足")
	}

	book,err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		c.JsonResult(6005,"项目不存在")
	}
	book.PrivatelyOwned = state

	err = book.Update()

	if err != nil {
		logs.Error("PrivatelyOwned => ",err)
		c.JsonResult(6004,"保存失败")
	}
	c.JsonResult(0,"ok")
}

// Transfer 转让项目.
func (c *BookController) Transfer()  {
	c.Prepare()
	account := c.GetString("account")

	if account == "" {
		c.JsonResult(6004,"接受者账号不能为空")
	}
	member,err := models.NewMember().FindByAccount(account)

	if err != nil {
		logs.Error("FindByAccount => ",err)
		c.JsonResult(6005,"接受用户不存在")
	}
	if member.Status != 0 {
		c.JsonResult(6006,"接受用户已被禁用")
	}
	if member.MemberId == c.Member.MemberId {
		c.JsonResult(6007,"不能转让给自己")
	}
	bookResult,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001,err.Error())
	}

	err = models.NewRelationship().Transfer(bookResult.BookId,c.Member.MemberId,member.MemberId)

	if err != nil {
		logs.Error("Transfer => ",err)
		c.JsonResult(6008,err.Error())
	}
	c.JsonResult(0,"ok")
}
//上传项目封面.
func (c *BookController) UploadCover()  {

	bookResult,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001,err.Error())
	}
	book,err := models.NewBook().Find(bookResult.BookId)

	if err != nil {
		logs.Error("SaveBook => ",err)
		c.JsonResult(6002,err.Error())
	}

	file,moreFile,err := c.GetFile("image-file")
	defer file.Close()

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"读取文件异常")
	}

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext,".png") && !strings.EqualFold(ext,".jpg") && !strings.EqualFold(ext,".gif") && !strings.EqualFold(ext,".jpeg")  {
		c.JsonResult(500,"不支持的图片格式")
	}


	x1 ,_ := strconv.ParseFloat(c.GetString("x"),10)
	y1 ,_ := strconv.ParseFloat(c.GetString("y"),10)
	w1 ,_ := strconv.ParseFloat(c.GetString("width"),10)
	h1 ,_ := strconv.ParseFloat(c.GetString("height"),10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	fileName := "cover_" +  strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := filepath.Join("uploads",time.Now().Format("200601"),fileName  + ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file",filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}
	defer func(filePath string) {
		os.Remove(filePath)
	}(filePath)

	//剪切图片
	subImg,err := graphics.ImageCopyFromFile(filePath,x,y,width,height)

	if err != nil{
		logs.Error("graphics.ImageCopyFromFile => ",err)
		c.JsonResult(500,"图片剪切")
	}

	filePath = filepath.Join(commands.WorkingDirectory,"uploads",time.Now().Format("200601"),fileName + "_small" + ext)

	//生成缩略图并保存到磁盘
	err = graphics.ImageResizeSaveFile(subImg,175,230,filePath)

	if err != nil {
		logs.Error("ImageResizeSaveFile => ",err.Error())
		c.JsonResult(500,"保存图片失败")
	}

	url := "/" +  strings.Replace(strings.TrimPrefix(filePath,commands.WorkingDirectory),"\\","/",-1)

	if strings.HasPrefix(url,"//") {
		url = string(url[1:])
	}

	old_cover := book.Cover

	book.Cover = url

	if err := book.Update() ; err != nil {
		c.JsonResult(6001,"保存图片失败")
	}
	//如果原封面不是默认封面则删除
	if old_cover != conf.GetDefaultCover() {
		os.Remove("." + old_cover)
	}

	c.JsonResult(0,"ok",url)
}

// Users 用户列表.
func (c *BookController) Users() {
	c.Prepare()
	c.TplName = "book/users.tpl"

	key := c.Ctx.Input.Param(":key")
	pageIndex,_ := c.GetInt("page",1)

	if key == ""{
		c.Abort("404")
	}

	book,err := models.NewBookResult().FindByIdentify(key,c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.Abort("403")
		}
		c.Abort("500")
	}

	c.Data["Model"] = *book

	members,totalCount,err := models.NewMemberRelationshipResult().FindForUsersByBookId(book.BookId,pageIndex,15)

	if totalCount > 0 {
		html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, 10, totalCount)

		c.Data["PageHtml"] = html
	}else{
		c.Data["PageHtml"] = ""
	}
	b,err := json.Marshal(members)

	if err != nil {
		c.Data["Result"] = template.JS("[]")
	}else{
		c.Data["Result"] = template.JS(string(b))
	}
}

// Create 创建项目.
func (c *BookController) Create() {

	if c.Ctx.Input.IsPost() {
		book_name := strings.TrimSpace(c.GetString("book_name",""))
		identify := strings.TrimSpace(c.GetString("identify",""))
		description := strings.TrimSpace(c.GetString("description",""))
		privately_owned,_ := strconv.Atoi(c.GetString("privately_owned"))
		comment_status := c.GetString("comment_status")

		if book_name == "" {
			c.JsonResult(6001,"项目名称不能为空")
		}
		if identify == "" {
			c.JsonResult(6002,"项目标识不能为空")
		}
		if ok,err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`,identify); !ok || err != nil {
			c.JsonResult(6003,"项目标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头")
		}
		if strings.Count(identify,"") > 50 {
			c.JsonResult(6004,"文档标识不能超过50字")
		}
		if strings.Count(description,"") > 500 {
			c.JsonResult(6004,"项目描述不能大于500字")
		}
		if privately_owned !=0 && privately_owned != 1 {
			privately_owned = 1
		}
		if comment_status != "open" && comment_status != "closed" && comment_status != "group_only" && comment_status != "registered_only" {
			comment_status = "closed"
		}

		book := models.NewBook()

		if books,_ := book.FindByField("identify",identify); len(books) > 0 {
			c.JsonResult(6006,"项目标识已存在")
		}

		book.BookName 	= book_name
		book.Description = description
		book.CommentCount = 0
		book.PrivatelyOwned = privately_owned
		book.CommentStatus = comment_status
		book.Identify 	= identify
		book.DocCount 	= 0
		book.MemberId 	= c.Member.MemberId
		book.CommentCount = 0
		book.Version 	= time.Now().Unix()
		book.Cover 	= conf.GetDefaultCover()
		book.Editor 	= "markdown"
		book.Theme	= "default"

		err := book.Insert()

		if err != nil {
			logs.Error("Insert => ",err)
			c.JsonResult(6005,"保存项目失败")
		}
		bookResult,err := models.NewBookResult().FindByIdentify(book.Identify,c.Member.MemberId)

		if err != nil {
			beego.Error(err)
		}

		c.JsonResult(0,"ok",bookResult)
	}
	c.JsonResult(6001,"error")
}

// CreateToken 创建访问来令牌.
func (c *BookController) CreateToken() {

	action := c.GetString("action")

	bookResult ,err := c.IsPermission()

	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403,"权限不足")
		}
		if err == orm.ErrNoRows {
			c.JsonResult(404,"项目不存在")
		}
		logs.Error("生成阅读令牌失败 =>",err)
		c.JsonResult(6002,err.Error())
	}
	book := models.NewBook()

	if _,err := book.Find(bookResult.BookId);err != nil {
		c.JsonResult(6001,"项目不存在")
	}
	if action == "create" {
		if bookResult.PrivatelyOwned == 0 {
			c.JsonResult(6001, "公开项目不能创建阅读令牌")
		}

		book.PrivateToken = string(utils.Krand(conf.GetTokenSize(), utils.KC_RAND_KIND_ALL))
		if err := book.Update(); err != nil {
			logs.Error("生成阅读令牌失败 => ", err)
			c.JsonResult(6003, "生成阅读令牌失败")
		}
		c.JsonResult(0, "ok", c.BaseUrl() +  beego.URLFor("DocumentController.Index",":key",book.Identify,"token",book.PrivateToken))
	}else{
		book.PrivateToken = ""
		if err := book.Update();err != nil {
			logs.Error("CreateToken => ",err)
			c.JsonResult(6004,"删除令牌失败")
		}
		c.JsonResult(0,"ok","")
	}
}

// Delete 删除项目.
func (c *BookController) Delete() {
	c.Prepare()

	bookResult ,err := c.IsPermission()

	if err != nil {
		c.JsonResult(6001,err.Error())
	}

	if bookResult.RoleId != conf.BookFounder {
		c.JsonResult(6002,"只有创始人才能删除项目")
	}
	err = models.NewBook().ThoroughDeleteBook(bookResult.BookId)

	if err == orm.ErrNoRows {
		c.JsonResult(6002,"项目不存在")
	}
	if err != nil {
		logs.Error("删除项目 => ",err)
		c.JsonResult(6003,"删除失败")
	}
	c.JsonResult(0,"ok")
}

//发布项目.
func (c *BookController) Release() {
	c.Prepare()

	identify := c.GetString("identify")

	book_id := 0

	if c.Member.IsAdministrator() {
		book,err := models.NewBook().FindByFieldFirst("identify",identify)
		if err != nil {

		}
		book_id = book.BookId
	}else {
		book, err := models.NewBookResult().FindByIdentify(identify, c.Member.MemberId)

		if err != nil {
			if err == models.ErrPermissionDenied {
				c.JsonResult(6001, "权限不足")
			}
			if err == orm.ErrNoRows {
				c.JsonResult(6002, "项目不存在")
			}
			beego.Error(err)
			c.JsonResult(6003, "未知错误")
		}
		if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder && book.RoleId != conf.BookEditor {
			c.JsonResult(6003, "权限不足")
		}
		book_id = book.BookId
	}
	go func(identify string) {
		models.NewDocument().ReleaseContent(book_id)


	}(identify)

	c.JsonResult(0,"发布任务已推送到任务队列，稍后将在后台执行。")
}

//文档排序.
func (c *BookController) SaveSort() {
	c.Prepare()

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	book_id := 0
	if c.Member.IsAdministrator() {
		book,err := models.NewBook().FindByFieldFirst("identify",identify)
		if err != nil {

		}
		book_id = book.BookId
	}else{
		bookResult,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)
		if err != nil {
			beego.Error("DocumentController.Edit => ",err)

			c.Abort("403")
		}
		if bookResult.RoleId == conf.BookObserver {
			c.JsonResult(6002,"项目不存在或权限不足")
		}
		book_id = bookResult.BookId
	}


	content := c.Ctx.Input.RequestBody

	var docs []map[string]interface{}

	err := json.Unmarshal(content,&docs)

	if err != nil {
		beego.Error(err)
		c.JsonResult(6003,"数据错误")
	}

	for _,item := range docs {
		if doc_id,ok := item["id"].(float64);ok {
			doc,err := models.NewDocument().Find(int(doc_id));
			if err != nil {
				beego.Error(err)
				continue;
			}
			if doc.BookId != book_id {
				logs.Info("%s","权限错误")
				continue;
			}
			sort,ok := item["sort"].(float64);
			if !ok {
				beego.Info("排序数字转换失败 => ",item)
				continue
			}
			parent_id,ok := item["parent"].(float64)
			if !ok {
				beego.Info("父分类转换失败 => ",item)
				continue
			}
			if parent_id > 0 {
				if parent,err := models.NewDocument().Find(int(parent_id)); err != nil || parent.BookId != book_id {
					continue
				}
			}
			doc.OrderSort = int(sort)
			doc.ParentId = int(parent_id)
			if err := doc.InsertOrUpdate(); err != nil {
				fmt.Printf("%s",err.Error())
				beego.Error(err)
			}
		}else{
			fmt.Printf("文档ID转换失败 => %+v",item)
		}

	}
	c.JsonResult(0,"ok")
}


func (c *BookController) IsPermission() (*models.BookResult,error) {
	identify := c.GetString("identify")

	book ,err := models.NewBookResult().FindByIdentify(identify,c.Member.MemberId)

	if err != nil {
		if err == models.ErrPermissionDenied {
			return book,errors.New("权限不足")
		}
		if err == orm.ErrNoRows {
			return book,errors.New("项目不存在")
		}
		return book,err
	}
	if book.RoleId != conf.BookAdmin && book.RoleId != conf.BookFounder {
		return book,errors.New("权限不足")
	}
	return book,nil
}

