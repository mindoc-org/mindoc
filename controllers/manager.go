package controllers

import (

	"encoding/json"
	"html/template"
	"strings"
	"regexp"

	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/utils"
	"github.com/lifei6671/godoc/models"
)

type ManagerController struct {
	BaseController
}

func (p *ManagerController) Index() {
	p.TplName = "manager/index.tpl"
}

// 用户列表.
func (c *ManagerController) Users()  {
	c.Prepare()
	c.TplName = "manager/users.tpl"

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	pageIndex,_ := c.GetInt("page",0)

	members,totalCount,err := models.NewMember().FindToPager(pageIndex,15)

	if err != nil {
		c.Data["ErrorMessage"] = err.Error()
		return
	}

	html := utils.GetPagerHtml(c.Ctx.Request.RequestURI,pageIndex,10,totalCount)

	c.Data["PageHtml"] = html
	b,err := json.Marshal(members)

	if err != nil {
		c.Data["Result"] = template.JS("[]")
	}else{
		c.Data["Result"] = template.JS(string(b))
	}
}

// 添加用户
func (c *ManagerController) CreateMember() {
	c.Prepare()
	if c.Member.Role != 0{
		c.Abort("403")
	}

	account := strings.TrimSpace(c.GetString("account"))
	password1 := strings.TrimSpace(c.GetString("password1"))
	password2 := strings.TrimSpace(c.GetString("password2"))
	email := strings.TrimSpace(c.GetString("email"))
	phone := strings.TrimSpace(c.GetString("phone"))
	role,_ := c.GetInt("role",1)
	status,_ := c.GetInt("status",0)

	if ok,err := regexp.MatchString(conf.RegexpAccount,account); account == "" || !ok || err != nil {
		c.JsonResult(6001,"账号只能由英文字母数字组成，且在3-50个字符")
	}
	if  l := strings.Count(password1,"") ; password1 == "" || l > 50 || l < 6{
		c.JsonResult(6002,"密码必须在6-50个字符之间")
	}
	if password1 != password2 {
		c.JsonResult(6003,"确认密码不正确")
	}
	if  ok,err := regexp.MatchString(conf.RegexpEmail,email); !ok || err != nil || email == "" {
		c.JsonResult(6004,"邮箱不能为空")
	}
	if role != 0 && role != 1 {
		role = 1
	}
	if status != 0 && status != 1 {
		status = 0
	}

	member := models.NewMember()

	if err := member.FindByAccount(account); err != nil {
		c.JsonResult(6005,"账号已存在")
	}

	member.Account = account
	member.Password = password1
	member.Role = role
	member.Avatar = conf.GetDefaultAvatar()
	member.CreateAt = c.Member.MemberId
	member.Email = email
	if phone != "" {
		member.Phone = phone
	}

	if err := member.Add(); err != nil {
		c.JsonResult(6006,err.Error())
	}

	c.JsonResult(0,"ok",member)
}

//更新用户状态.
func (c *ManagerController) UpdateMemberStatus()  {
	c.Prepare()

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	member_id,_ := c.GetInt("member_id",0)
	status ,_ := c.GetInt("status",0)

	if member_id <= 0 {
		c.JsonResult(6001,"参数错误")
	}
	if status != 0 && status != 1 {
		status = 0
	}
	member := models.NewMember()

	if err := member.Find(member_id); err != nil {
		c.JsonResult(6002,"用户不存在")
	}
	member.Status = status

	if err := member.Update();err != nil {
		logs.Error("",err)
		c.JsonResult(6003,"用户状态设置失败")
	}
	c.JsonResult(0,"ok",member)
}

func (c *ManagerController) Books()  {
	c.Prepare()
	c.TplName = "manager/books.tpl"

}

func (c *ManagerController) DeleteBook()  {
	c.Prepare()
	if c.Member.Role != 0 {
		c.Abort("403")
	}
}

func (c *ManagerController) Comments()  {
	c.Prepare()
	if c.Member.Role != 0 {
		c.Abort("403")
	}
}

func (c *ManagerController) DeleteComment()  {
	c.Prepare()
	if c.Member.Role != 0 {
		c.Abort("403")
	}
}




















