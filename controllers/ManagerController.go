package controllers

import (
	"encoding/json"
	"html/template"
	"regexp"
	"strings"

	"math"
	"path/filepath"
	"strconv"

	"io/ioutil"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/filetil"
	"github.com/mindoc-org/mindoc/utils/pagination"
	"github.com/russross/blackfriday/v2"
)

type ManagerController struct {
	BaseController
}

func (c *ManagerController) Prepare() {
	c.BaseController.Prepare()

	if !c.Member.IsAdministrator() {
		c.Abort("403")
	}
}

func (c *ManagerController) Index() {
	c.TplName = "manager/index.tpl"

	c.Data["Model"] = models.NewDashboard().Query()
	c.Data["Action"] = "index"
}

// 用户列表.
func (c *ManagerController) Users() {
	c.TplName = "manager/users.tpl"
	c.Data["Action"] = "users"
	pageIndex, _ := c.GetInt("page", 0)

	tempMember := models.NewMember()
	tempMember.Lang = c.Lang
	members, totalCount, err := tempMember.FindToPager(pageIndex, conf.PageSize)
	if err != nil {
		c.Data["ErrorMessage"] = err.Error()
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()

		for _, item := range members {
			item.Avatar = conf.URLForWithCdnImage(item.Avatar)
		}
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

// 添加用户.
func (c *ManagerController) CreateMember() {

	account := strings.TrimSpace(c.GetString("account"))
	password1 := strings.TrimSpace(c.GetString("password1"))
	password2 := strings.TrimSpace(c.GetString("password2"))
	email := strings.TrimSpace(c.GetString("email"))
	phone := strings.TrimSpace(c.GetString("phone"))
	role, _ := c.GetInt("role", 1)
	status, _ := c.GetInt("status", 0)

	if ok, err := regexp.MatchString(conf.RegexpAccount, account); account == "" || !ok || err != nil {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.username_invalid_format"))
	}
	if l := strings.Count(password1, ""); password1 == "" || l > 50 || l < 6 {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.pwd_length_tips"))
	}
	if password1 != password2 {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.wrong_confirm_pwd"))
	}
	if ok, err := regexp.MatchString(conf.RegexpEmail, email); !ok || err != nil || email == "" {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.email_invalid_format"))
	}
	if role != 0 && role != 1 && role != 2 {
		role = 1
	}
	if status != 0 && status != 1 {
		status = 0
	}

	member := models.NewMember()

	if _, err := member.FindByAccount(account); err == nil && member.MemberId > 0 {
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.account_existed"))
	}

	member.Account = account
	member.Password = password1
	member.Role = conf.SystemRole(role)
	member.Avatar = conf.GetDefaultAvatar()
	member.CreateAt = c.Member.MemberId
	member.Email = email
	member.RealName = strings.TrimSpace(c.GetString("real_name", ""))
	member.Lang = c.Lang

	if phone != "" {
		member.Phone = phone
	}

	if err := member.Add(); err != nil {
		c.JsonResult(6006, err.Error())
	}

	c.JsonResult(0, "ok", member)
}

//更新用户状态.
func (c *ManagerController) UpdateMemberStatus() {
	c.Prepare()

	member_id, _ := c.GetInt("member_id", 0)
	status, _ := c.GetInt("status", 0)

	if member_id <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	if status != 0 && status != 1 {
		status = 0
	}
	member := models.NewMember()

	if _, err := member.Find(member_id); err != nil {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	if member.MemberId == c.Member.MemberId {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.cannot_change_own_status"))
	}
	if member.Role == conf.MemberSuperRole {
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.cannot_change_super_status"))
	}
	member.Status = status

	if err := member.Update(); err != nil {
		logs.Error("", err)
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "ok", member)
}

//变更用户权限.
func (c *ManagerController) ChangeMemberRole() {
	c.Prepare()

	memberId, _ := c.GetInt("member_id", 0)
	role, _ := c.GetInt("role", 0)
	if memberId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	if role != int(conf.MemberAdminRole) && role != int(conf.MemberGeneralRole) {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.no_permission"))
	}
	member := models.NewMember()

	if _, err := member.Find(memberId); err != nil {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	if member.MemberId == c.Member.MemberId {
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.cannot_change_own_priv"))
	}
	if member.Role == conf.MemberSuperRole {
		c.JsonResult(6005, i18n.Tr(c.Lang, "message.cannot_change_super_priv"))
	}
	member.Role = conf.SystemRole(role)

	if err := member.Update(); err != nil {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
	}
	member.Lang = c.Lang
	member.ResolveRoleName()
	c.JsonResult(0, "ok", member)
}

//编辑用户信息.
func (c *ManagerController) EditMember() {
	c.Prepare()
	c.TplName = "manager/edit_users.tpl"
	c.Data["Action"] = "users"
	member_id, _ := c.GetInt(":id", 0)

	if member_id <= 0 {
		c.Abort("404")
	}

	member, err := models.NewMember().Find(member_id)
	if err != nil {
		logs.Error(err)
		c.Abort("404")
	}
	if c.Ctx.Input.IsPost() {
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		email := c.GetString("email")
		phone := c.GetString("phone")
		description := c.GetString("description")
		member.Email = email
		member.Phone = phone
		member.Description = description
		member.RealName = c.GetString("real_name")
		if password1 != "" && password2 != password1 {
			c.JsonResult(6001, i18n.Tr(c.Lang, "message.wrong_confirm_pwd"))
		}
		if password1 != "" && member.AuthMethod != conf.AuthMethodLDAP {
			member.Password = password1
		}
		if err := member.Valid(password1 == ""); err != nil {
			c.JsonResult(6002, err.Error())
		}
		if password1 != "" {
			password, err := utils.PasswordHash(password1)
			if err != nil {
				logs.Error(err)
				c.JsonResult(6003, i18n.Tr(c.Lang, "message.pwd_encrypt_failed"))
			}
			member.Password = password
		}
		if err := member.Update(); err != nil {
			c.JsonResult(6004, err.Error())
		}
		c.JsonResult(0, "ok")
	}

	c.Data["Model"] = member
}

//删除一个用户，并将该用户的所有信息转移到超级管理员上.
func (c *ManagerController) DeleteMember() {
	c.Prepare()
	member_id, _ := c.GetInt("id", 0)

	if member_id <= 0 {
		c.JsonResult(404, i18n.Tr(c.Lang, "message.param_error"))
	}

	member, err := models.NewMember().Find(member_id)
	if err != nil {
		logs.Error(err)
		c.JsonResult(500, i18n.Tr(c.Lang, "message.user_not_existed"))
	}
	if member.Role == conf.MemberSuperRole {
		c.JsonResult(500, "不能删除超级管理员")
	}

	superMember, err := models.NewMember().FindByFieldFirst("role", 0)
	if err != nil {
		logs.Error(err)
		c.JsonResult(5001, "未能找到超级管理员")
	}

	err = models.NewMember().Delete(member_id, superMember.MemberId)
	if err != nil {
		logs.Error(err)
		c.JsonResult(5002, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "ok")
}

//项目列表.
func (c *ManagerController) Books() {
	c.Prepare()
	c.TplName = "manager/books.tpl"
	c.Data["Action"] = "books"
	pageIndex, _ := c.GetInt("page", 1)

	books, totalCount, err := models.NewBookResult().FindToPager(pageIndex, conf.PageSize)

	if err != nil {
		c.Abort("500")
	}

	if totalCount > 0 {
		//html := utils.GetPagerHtml(c.Ctx.Request.RequestURI, pageIndex, 8, totalCount)

		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())

		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	for i, book := range books {
		books[i].Description = utils.StripTags(string(blackfriday.Run([]byte(book.Description))))
		books[i].ModifyTime = book.ModifyTime.Local()
		books[i].CreateTime = book.CreateTime.Local()
	}
	c.Data["Lists"] = books
}

//编辑项目.
func (c *ManagerController) EditBook() {
	c.Prepare()

	c.TplName = "manager/edit_book.tpl"
	c.Data["Action"] = "books"
	identify := c.GetString(":key")

	if identify == "" {
		c.Abort("404")
	}
	book, err := models.NewBook().FindByFieldFirst("identify", identify)
	if err != nil {
		c.Abort("500")
	}

	if c.Ctx.Input.IsPost() {
		bookName := strings.TrimSpace(c.GetString("book_name"))
		description := strings.TrimSpace(c.GetString("description", ""))
		commentStatus := c.GetString("comment_status")
		tag := strings.TrimSpace(c.GetString("label"))
		orderIndex, _ := c.GetInt("order_index", 0)
		isDownload := strings.TrimSpace(c.GetString("is_download")) == "on"
		enableShare := strings.TrimSpace(c.GetString("enable_share")) == "on"
		isUseFirstDocument := strings.TrimSpace(c.GetString("is_use_first_document")) == "on"
		autoRelease := strings.TrimSpace(c.GetString("auto_release")) == "on"
		publisher := strings.TrimSpace(c.GetString("publisher"))
		historyCount, _ := c.GetInt("history_count", 0)
		itemId, _ := c.GetInt("itemId")

		if strings.Count(description, "") > 500 {
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.project_desc_tips"))
		}
		if commentStatus != "open" && commentStatus != "closed" && commentStatus != "group_only" && commentStatus != "registered_only" {
			commentStatus = "closed"
		}
		if tag != "" {
			tags := strings.Split(tag, ";")
			if len(tags) > 10 {
				c.JsonResult(6005, "最多允许添加10个标签")
			}
		}
		if !models.NewItemsets().Exist(itemId) {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.project_space_not_exist"))
		}
		book.Publisher = publisher
		book.HistoryCount = historyCount
		book.BookName = bookName
		book.Description = description
		book.CommentStatus = commentStatus
		book.Label = tag
		book.OrderIndex = orderIndex
		book.ItemId = itemId
		book.BookPassword = strings.TrimSpace(c.GetString("bPassword"))

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

		if err := book.Update(); err != nil {
			c.JsonResult(6006, i18n.Tr(c.Lang, "message.failed"))
		}
		c.JsonResult(0, "ok")
	}
	if book.PrivateToken != "" {
		book.PrivateToken = conf.URLFor("DocumentController.Index", ":key", book.Identify, "token", book.PrivateToken)
	}
	bookResult := models.NewBookResult()
	bookResult.ToBookResult(*book)

	c.Data["Model"] = bookResult
}

// 删除项目.
func (c *ManagerController) DeleteBook() {
	c.Prepare()

	bookId, _ := c.GetInt("book_id", 0)

	if bookId <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}
	book := models.NewBook()

	err := book.ThoroughDeleteBook(bookId)

	if err == orm.ErrNoRows {
		c.JsonResult(6002, i18n.Tr(c.Lang, "message.item_not_exist"))
	}
	if err != nil {
		logs.Error("删除失败 -> ", err)
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "ok")
}

// CreateToken 创建访问来令牌.
func (c *ManagerController) CreateToken() {
	c.Prepare()
	action := c.GetString("action")

	identify := c.GetString("identify")

	book, err := models.NewBook().FindByFieldFirst("identify", identify)

	if err != nil {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.item_not_exist"))
	}
	if action == "create" {

		if book.PrivatelyOwned == 0 {
			c.JsonResult(6001, "公开项目不能创建阅读令牌")
		}

		book.PrivateToken = string(utils.Krand(conf.GetTokenSize(), utils.KC_RAND_KIND_ALL))
		if err := book.Update(); err != nil {
			logs.Error("生成阅读令牌失败 => ", err)
			c.JsonResult(6003, i18n.Tr(c.Lang, "message.failed"))
		}
		c.JsonResult(0, "ok", conf.URLFor("DocumentController.Index", ":key", book.Identify, "token", book.PrivateToken))
	} else {
		book.PrivateToken = ""
		if err := book.Update(); err != nil {
			logs.Error("CreateToken => ", err)
			c.JsonResult(6004, i18n.Tr(c.Lang, "message.failed"))
		}
		c.JsonResult(0, "ok", "")
	}
}

//项目设置.
func (c *ManagerController) Setting() {
	c.Prepare()
	c.TplName = "manager/setting.tpl"
	c.Data["Action"] = "setting"
	options, err := models.NewOption().All()

	if c.Ctx.Input.IsPost() {
		for _, item := range options {
			item.OptionValue = c.GetString(item.OptionName)
			item.InsertOrUpdate()
		}
		c.JsonResult(0, "ok")
	}

	if err != nil {
		c.Abort("500")
	}
	c.Data["SITE_TITLE"] = c.Option["SITE_NAME"]

	for _, item := range options {
		c.Data[item.OptionName] = item.OptionValue
	}
}

// Transfer 转让项目.
func (c *ManagerController) Transfer() {
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

	if !c.Member.IsAdministrator() {
		c.Abort("403")
	}

	identify := c.GetString("identify")

	book, err := models.NewBook().FindByFieldFirst("identify", identify)
	if err != nil {
		c.JsonResult(6001, err.Error())
	}
	rel, err := models.NewRelationship().FindFounder(book.BookId)

	if err != nil {
		logs.Error("FindFounder => ", err)
		c.JsonResult(6009, "查询项目创始人失败")
	}
	if member.MemberId == rel.MemberId {
		c.JsonResult(6007, "不能转让给自己")
	}

	err = models.NewRelationship().Transfer(book.BookId, rel.MemberId, member.MemberId)

	if err != nil {
		logs.Error("Transfer => ", err)
		c.JsonResult(6008, err.Error())
	}
	c.JsonResult(0, "ok")
}

func (c *ManagerController) Comments() {
	c.Prepare()
	c.TplName = "manager/comments.tpl"
	if !c.Member.IsAdministrator() {
		c.Abort("403")
	}

}

//DeleteComment 标记评论为已删除
func (c *ManagerController) DeleteComment() {
	c.Prepare()

	comment_id, _ := c.GetInt("comment_id", 0)

	if comment_id <= 0 {
		c.JsonResult(6001, i18n.Tr(c.Lang, "message.param_error"))
	}

	comment := models.NewComment()

	if _, err := comment.Find(comment_id); err != nil {
		c.JsonResult(6002, "评论不存在")
	}

	comment.Approved = 3

	if err := comment.Update("approved"); err != nil {
		c.JsonResult(6003, "删除评论失败")
	}
	c.JsonResult(0, "ok", comment)
}

//设置项目私有状态.
func (c *ManagerController) PrivatelyOwned() {
	c.Prepare()
	status := c.GetString("status")
	identify := c.GetString("identify")

	if status != "open" && status != "close" {
		c.JsonResult(6003, i18n.Tr(c.Lang, "message.param_error"))
	}
	state := 0
	if status == "open" {
		state = 0
	} else {
		state = 1
	}

	if !c.Member.IsAdministrator() {
		c.Abort("403")
	}

	book, err := models.NewBook().FindByFieldFirst("identify", identify)
	if err != nil {
		c.JsonResult(6001, err.Error())
	}

	book.PrivatelyOwned = state

	logs.Info("", state, status)

	err = book.Update()

	if err != nil {
		logs.Error("PrivatelyOwned => ", err)
		c.JsonResult(6004, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "ok")
}

//附件列表.
func (c *ManagerController) AttachList() {
	c.Prepare()
	c.TplName = "manager/attach_list.tpl"
	c.Data["Action"] = "attach"

	pageIndex, _ := c.GetInt("page", 1)

	attachList, totalCount, err := models.NewAttachment().FindToPager(pageIndex, conf.PageSize)

	if err != nil {
		c.Abort("500")
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	for _, item := range attachList {

		p := filepath.Join(conf.WorkingDirectory, item.FilePath)

		item.IsExist = filetil.FileExists(p)

	}
	c.Data["Lists"] = attachList
}

//附件清理.
func (c *ManagerController) AttachClean() {
	c.Prepare()

	attachList, _, err := models.NewAttachment().FindToPager(0, 0)

	if err != nil {
		c.Abort("500")
	}

	for _, item := range attachList {

		p := filepath.Join(conf.WorkingDirectory, item.FilePath)

		item.IsExist = filetil.FileExists(p)
		if item.IsExist {
			// 判断
			searchList, err := models.NewDocumentSearchResult().SearchAllDocument(item.HttpPath)
			if err != nil {
				c.Abort("500")
			} else if len(searchList) == 0 {
				logs.Info("delete file:", item.FilePath)
				item.FilePath = p
				if err := item.Delete(); err != nil {
					logs.Error("AttachDelete => ", err)
					c.JsonResult(6002, err.Error())
					break
				}
			}
		}
	}

	c.JsonResult(0, "ok")
}

//附件详情.
func (c *ManagerController) AttachDetailed() {
	c.Prepare()
	c.TplName = "manager/attach_detailed.tpl"
	c.Data["Action"] = "attach"

	attach_id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if attach_id <= 0 {
		c.Abort("404")
	}

	attach, err := models.NewAttachmentResult().Find(attach_id)
	if err != nil {
		logs.Error("AttachDetailed => ", err)
		if err == orm.ErrNoRows {
			c.Abort("404")
		} else {
			c.Abort("500")
		}
	}

	attach.FilePath = filepath.Join(conf.WorkingDirectory, attach.FilePath)
	attach.HttpPath = conf.URLForWithCdnImage(attach.HttpPath)

	attach.IsExist = filetil.FileExists(attach.FilePath)

	c.Data["Model"] = attach
}

//删除附件.
func (c *ManagerController) AttachDelete() {
	c.Prepare()
	attachId, _ := c.GetInt("attach_id")

	if attachId <= 0 {
		c.Abort("404")
	}
	attach, err := models.NewAttachment().Find(attachId)

	if err != nil {
		logs.Error("AttachDelete => ", err)
		c.JsonResult(6001, err.Error())
	}
	attach.FilePath = filepath.Join(conf.WorkingDirectory, attach.FilePath)

	if err := attach.Delete(); err != nil {
		logs.Error("AttachDelete => ", err)
		c.JsonResult(6002, err.Error())
	}
	c.JsonResult(0, "ok")
}

//标签列表
func (c *ManagerController) LabelList() {
	c.Prepare()
	c.TplName = "manager/label_list.tpl"
	c.Data["Action"] = "label"
	pageIndex, _ := c.GetInt("page", 1)

	labels, totalCount, err := models.NewLabel().FindToPager(pageIndex, conf.PageSize)
	if err != nil {
		c.ShowErrorPage(50001, err.Error())
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(conf.PageSize)))

	c.Data["Lists"] = labels
}

//删除标签
func (c *ManagerController) LabelDelete() {
	labelId, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("获取删除标签参数时出错:", err)
		c.JsonResult(50001, i18n.Tr(c.Lang, "message.param_error"))
	}
	if labelId <= 0 {
		c.JsonResult(50001, i18n.Tr(c.Lang, "message.param_error"))
	}

	label, err := models.NewLabel().FindFirst("label_id", labelId)
	if err != nil {
		logs.Error("查询标签时出错:", err)
		c.JsonResult(50001, "查询标签时出错:"+err.Error())
	}
	if err := label.Delete(); err != nil {
		c.JsonResult(50002, "删除失败:"+err.Error())
	} else {
		c.JsonResult(0, "ok")
	}
}

func (c *ManagerController) Config() {
	c.Prepare()
	c.TplName = "manager/config.tpl"
	c.Data["Action"] = "config"
	if c.Ctx.Input.IsPost() {
		content := strings.TrimSpace(c.GetString("configFileTextArea"))
		if content == "" {
			c.JsonResult(500, "配置文件不能为空")
		}
		tf, err := ioutil.TempFile(os.TempDir(), "mindoc")

		if err != nil {
			logs.Error("创建临时文件失败 ->", err)
			c.JsonResult(5001, "创建临时文件失败")
		}
		defer tf.Close()

		tf.WriteString(content)

		err = web.LoadAppConfig("ini", tf.Name())

		if err != nil {
			logs.Error("加载配置文件失败 ->", err)
			c.JsonResult(5002, "加载配置文件失败")
		}
		err = filetil.CopyFile(tf.Name(), conf.ConfigurationFile)
		if err != nil {
			logs.Error("保存配置文件失败 ->", err)
			c.JsonResult(5003, "保存配置文件失败")
		}
		c.JsonResult(0, "保存成功")
	}
	c.Data["ConfigContent"] = ""
	if b, err := ioutil.ReadFile(conf.ConfigurationFile); err == nil {
		c.Data["ConfigContent"] = string(b)
	}
}

func (c *ManagerController) Team() {
	c.Prepare()
	c.TplName = "manager/team.tpl"
	c.Data["Action"] = "team"
	pageIndex, _ := c.GetInt("page", 0)

	teams, totalCount, err := models.NewTeam().FindToPager(pageIndex, conf.PageSize)
	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if err == orm.ErrNoRows || len(teams) <= 0 {
		c.Data["Result"] = template.JS("[]")
		c.Data["PageHtml"] = ""
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	b, err := json.Marshal(teams)
	if err != nil {
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
}

func (c *ManagerController) TeamCreate() {
	c.Prepare()

	teamName := c.GetString("teamName")

	if teamName == "" {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.team_name_empty"))
	}
	team := models.NewTeam()

	team.MemberId = c.Member.MemberId
	team.TeamName = teamName

	if err := team.Save(); err == nil {
		c.JsonResult(0, "OK", team)
	} else {
		c.JsonResult(5002, err.Error())
	}

}

func (c *ManagerController) TeamEdit() {
	c.Prepare()
	teamName := c.GetString("teamName")
	teamId, _ := c.GetInt("teamId")

	if teamName == "" {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.team_name_empty"))
	}
	if teamId <= 0 {
		c.JsonResult(5002, i18n.Tr(c.Lang, "message.team_id_empty"))
	}
	team, err := models.NewTeam().First(teamId)

	c.CheckJsonError(5003, err)

	team.TeamName = teamName

	err = team.Save()

	c.CheckJsonError(5004, err)

	c.JsonResult(0, "OK", team)
}

func (c *ManagerController) TeamDelete() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")
	if teamId <= 0 {
		c.JsonResult(5002, i18n.Tr(c.Lang, "message.team_id_empty"))
	}
	err := models.NewTeam().Delete(teamId)
	c.CheckJsonError(5001, err)

	c.JsonResult(0, "OK")
}

func (c *ManagerController) TeamMemberList() {
	c.Prepare()
	c.TplName = "manager/team_member_list.tpl"
	c.Data["Action"] = "team"
	teamId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if teamId <= 0 {
		c.ShowErrorPage(500, i18n.Tr(c.Lang, "message.param_error"))
	}
	pageIndex, _ := c.GetInt("page", 0)

	team, err := models.NewTeam().First(teamId)
	if err == orm.ErrNoRows {
		c.ShowErrorPage(404, "团队不存在")
	}
	c.CheckErrorResult(500, err)
	c.Data["Model"] = team

	teams, totalCount, err := models.NewTeamMember().SetLang(c.Lang).FindToPager(teamId, pageIndex, conf.PageSize)
	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if err == orm.ErrNoRows || len(teams) <= 0 {
		c.Data["Result"] = template.JS("[]")
		c.Data["PageHtml"] = ""
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	b, err := json.Marshal(teams)
	if err != nil {
		logs.Error("编码 JSON 结果失败 ->", err)
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
}

//搜索团队用户.
func (c *ManagerController) TeamSearchMember() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")
	keyword := strings.TrimSpace(c.GetString("q"))

	if teamId <= 0 {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.param_error"))
	}

	searchResult, err := models.NewTeamMember().FindNotJoinMemberByAccount(teamId, keyword, 10)

	if err != nil {
		c.JsonResult(500, err.Error())
	}
	c.JsonResult(0, "OK", searchResult)
}

func (c *ManagerController) TeamMemberAdd() {
	c.Prepare()
	teamId, _ := c.GetInt("teamId")
	memberId, _ := c.GetInt("memberId")
	roleId, _ := c.GetInt("roleId")

	if teamId <= 0 || memberId <= 0 || roleId <= 0 || roleId > int(conf.BookObserver) {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.system_error"))
	}

	teamMember := models.NewTeamMember()
	teamMember.MemberId = memberId
	teamMember.TeamId = teamId
	teamMember.RoleId = conf.BookRole(roleId)

	if err := teamMember.Save(); err != nil {
		c.CheckJsonError(5001, err)
	}

	teamMember.Include()

	c.JsonResult(0, "OK", teamMember)
}

func (c *ManagerController) TeamMemberDelete() {
	c.Prepare()
	memberId, _ := c.GetInt("memberId")
	teamId, _ := c.GetInt("teamId")

	teamMember, err := models.NewTeamMember().FindFirst(teamId, memberId)

	if err != nil {
		c.JsonResult(5001, "用户不存在或已禁用")
	}
	err = teamMember.Delete(teamMember.TeamMemberId)
	if err != nil {
		c.JsonResult(5002, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "ok")
}

func (c *ManagerController) TeamChangeMemberRole() {
	c.Prepare()
	memberId, _ := c.GetInt("memberId")
	roleId, _ := c.GetInt("roleId")
	teamId, _ := c.GetInt("teamId")
	if memberId <= 0 || roleId <= 0 || teamId <= 0 || roleId > int(conf.BookObserver) {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.param_error"))
	}

	teamMember, err := models.NewTeamMember().ChangeRoleId(teamId, memberId, conf.BookRole(roleId))

	if err != nil {
		c.JsonResult(5002, err.Error())
	} else {
		c.JsonResult(0, "OK", teamMember)
	}

}

//团队项目列表.
func (c *ManagerController) TeamBookList() {
	c.Prepare()
	c.TplName = "manager/team_book_list.tpl"
	c.Data["Action"] = "team"
	teamId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	pageIndex, _ := c.GetInt("page", 0)

	if teamId <= 0 {
		c.JsonResult(5002, i18n.Tr(c.Lang, "message.team_id_empty"))
	}

	team, err := models.NewTeam().First(teamId)

	if err == orm.ErrNoRows {
		c.ShowErrorPage(404, "团队不存在")
	}
	c.CheckErrorResult(500, err)
	c.Data["Model"] = team

	teams, totalCount, err := models.NewTeamRelationship().FindToPager(teamId, pageIndex, conf.PageSize)

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if err == orm.ErrNoRows || len(teams) <= 0 {
		c.Data["Result"] = template.JS("[]")
		c.Data["PageHtml"] = ""
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	b, err := json.Marshal(teams)
	if err != nil {
		logs.Error("编码 JSON 结果失败 ->", err)
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}
}

//给团队增加项目.
func (c *ManagerController) TeamBookAdd() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")
	bookId, _ := c.GetInt("bookId")

	if teamId <= 0 || bookId <= 0 {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.param_error"))
	}
	teamRel := models.NewTeamRelationship()
	teamRel.BookId = bookId
	teamRel.TeamId = teamId

	err := teamRel.Save()

	if err != nil {
		c.JsonResult(5001, err.Error())
	} else {
		teamRel.Include()
		c.JsonResult(0, "OK", teamRel)
	}
}

//搜索未参与的项目.
func (c *ManagerController) TeamSearchBook() {
	c.Prepare()

	teamId, _ := c.GetInt("teamId")
	keyword := strings.TrimSpace(c.GetString("q"))

	if teamId <= 0 {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.param_error"))
	}

	searchResult, err := models.NewTeamRelationship().FindNotJoinBookByName(teamId, keyword, 10)

	if err != nil {
		c.JsonResult(500, err.Error())
	}
	c.JsonResult(0, "OK", searchResult)

}

//删除团队项目.
func (c *ManagerController) TeamBookDelete() {
	c.Prepare()
	teamRelationshipId, _ := c.GetInt("teamRelId")

	if teamRelationshipId <= 0 {
		c.JsonResult(500, i18n.Tr(c.Lang, "message.param_error"))
	}

	err := models.NewTeamRelationship().Delete(teamRelationshipId)

	if err != nil {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.failed"))
	}
	c.JsonResult(0, "OK")
}

//项目空间列表.
func (c *ManagerController) Itemsets() {
	c.Prepare()
	c.TplName = "manager/itemsets.tpl"
	c.Data["Action"] = "itemsets"
	pageIndex, _ := c.GetInt("page", 0)

	items, totalCount, err := models.NewItemsets().FindToPager(pageIndex, conf.PageSize)

	if err != nil && err != orm.ErrNoRows {
		c.ShowErrorPage(500, err.Error())
	}
	if err == orm.ErrNoRows || len(items) <= 0 {
		c.Data["Lists"] = items
		c.Data["PageHtml"] = ""
		return
	}

	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["Lists"] = items
}

//编辑或添加项目空间.
func (c *ManagerController) ItemsetsEdit() {
	c.Prepare()
	itemId, _ := c.GetInt("itemId")
	itemName := strings.TrimSpace(c.GetString("itemName"))
	itemKey := strings.TrimSpace(c.GetString("itemKey"))
	if itemName == "" || itemKey == "" {
		c.JsonResult(5001, i18n.Tr(c.Lang, "message.param_error"))
	}
	var item *models.Itemsets
	var err error
	if itemId > 0 {
		if item, err = models.NewItemsets().First(itemId); err != nil {
			if err == orm.ErrNoRows {
				c.JsonResult(5002, i18n.Tr(c.Lang, "message.project_space_not_exist"))
			} else {
				c.JsonResult(5003, "查询项目空间出错")
			}
		}
	} else {
		item = models.NewItemsets()
	}

	item.ItemKey = itemKey
	item.ItemName = itemName
	item.MemberId = c.Member.MemberId
	item.ModifyAt = c.Member.MemberId

	if err := item.Save(); err != nil {
		c.JsonResult(5004, err.Error())
	}

	c.JsonResult(0, "OK")
}

//删除项目空间.
func (c *ManagerController) ItemsetsDelete() {
	c.Prepare()
	itemId, _ := c.GetInt("itemId")

	if err := models.NewItemsets().Delete(itemId); err != nil {
		c.JsonResult(5001, err.Error())
	}
	c.JsonResult(0, "OK")
}
