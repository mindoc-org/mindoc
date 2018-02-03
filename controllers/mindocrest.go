package controllers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/models"
)

// MinDocRestController struct.
type MinDocRestController struct {
	BaseController
}

// PostContent API 设置文档内容.
func (c *MinDocRestController) PostContent() {
	c.Prepare()

	var req models.MinDocRest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	folderkey := req.Folder
	doctitle := req.Title
	dockey := req.Identify
	tokenkey := req.Token
	textmd := req.TextMD
	texthtml := req.TextHTML

	if doctitle == "" {
		c.JsonResult(6004, "文档名称不能为空")
	}
	if dockey == "" {
		c.JsonResult(6004, "文档标识不能为空")
	}
	if ok, err := regexp.MatchString(`^[a-z]+[a-zA-Z0-9_\-]*$`, dockey); !ok || err != nil {

		c.JsonResult(6003, "文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头")
	}

	token, err := models.NewBook().FindByFieldFirst("private_token", tokenkey)
	if err != nil {
		beego.Error("token => ", err)
		c.JsonResult(6002, "系统权限不足["+tokenkey+"]")
	}
	beego.Info("req tokenkey =>" + tokenkey + "  -->" + fmt.Sprintf("%d", token.BookId))

	folder, err := models.NewDocument().FindByFieldFirst("identify", folderkey)
	if err != nil {
		beego.Error("folder => ", err)
		c.JsonResult(6002, "项目或类目不存在或权限不足")
	}
	beego.Info("req folderkey =>" + folderkey + "  -->" + fmt.Sprintf("%d", folder.BookId))
	if folder.BookId != token.BookId {
		c.JsonResult(6002, "folder和token不匹配")
	}

	doc, err := models.NewDocument().FindByFieldFirst("identify", dockey)
	if err != nil {
		if err != orm.ErrNoRows {
			c.JsonResult(6006, "文档获取失败")
		}
	}
	if doc.BookId > 0 && folder.BookId != doc.BookId {
		c.JsonResult(6002, "文档标识已经被使用")
	}
	doc.BookId = token.BookId
	doc.MemberId = token.MemberId
	doc.Identify = dockey
	doc.Version = time.Now().Unix()
	doc.DocumentName = doctitle
	doc.ParentId = folder.DocumentId
	doc.Markdown = textmd
	doc.Content = texthtml
	doc.Release = texthtml
	if err := doc.InsertOrUpdate(); err != nil {
		beego.Error("InsertOrUpdate => ", err)
		c.JsonResult(6005, "保存失败")
	}
	go func() {
		models.NewDocument().ReleaseContent(doc.BookId)
	}()
	//减少返回信息
	doc.Markdown = ""
	doc.Content = ""
	doc.Release = ""
	c.JsonResult(0, "ok", doc)
}
