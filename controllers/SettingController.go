package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/graphics"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
)

type SettingController struct {
	BaseController
}

func (c *SettingController) Index() {
	c.TplName = "setting/index.tpl"

	if c.Ctx.Input.IsPost() {
		email := strings.TrimSpace(c.GetString("email", ""))
		phone := strings.TrimSpace(c.GetString("phone"))
		description := strings.TrimSpace(c.GetString("description"))

		if email == "" {
			c.JsonResult(601, "邮箱不能为空")
		}
		member := c.Member
		member.Email = email
		member.Phone = phone
		member.Description = description
		member.RealName = strings.TrimSpace(c.GetString("real_name", ""))
		if err := member.Update(); err != nil {
			c.JsonResult(602, err.Error())
		}
		c.SetMember(*member)
		c.JsonResult(0, "ok")
	}
}

func (c *SettingController) Password() {
	c.TplName = "setting/password.tpl"

	if c.Ctx.Input.IsPost() {
		if c.Member.AuthMethod == conf.AuthMethodLDAP {
			c.JsonResult(6009, "当前用户不支持修改密码")
		}
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		password3 := c.GetString("password3")

		if password1 == "" {
			c.JsonResult(6003, "原密码不能为空")
		}

		if password2 == "" {
			c.JsonResult(6004, "新密码不能为空")
		}
		if count := strings.Count(password2, ""); count < 6 || count > 18 {
			c.JsonResult(6009, "密码必须在6-18字之间")
		}
		if password2 != password3 {
			c.JsonResult(6003, "确认密码不正确")
		}
		if ok, _ := utils.PasswordVerify(c.Member.Password, password1); !ok {
			c.JsonResult(6005, "原始密码不正确")
		}
		if password1 == password2 {
			c.JsonResult(6006, "新密码不能和原始密码相同")
		}
		pwd, err := utils.PasswordHash(password2)
		if err != nil {
			c.JsonResult(6007, "密码加密失败")
		}
		c.Member.Password = pwd
		if err := c.Member.Update(); err != nil {
			c.JsonResult(6008, err.Error())
		}
		c.JsonResult(0, "ok")
	}
}

// Upload 上传图片
func (c *SettingController) Upload() {
	file, moreFile, err := c.GetFile("image-file")
	defer file.Close()

	if err != nil {
		logs.Error("", err.Error())
		c.JsonResult(500, "读取文件异常")
	}

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

	fmt.Println(x, x1, y, y1)

	fileName := "avatar_" + strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file", filePath)

	if err != nil {
		logs.Error("", err)
		c.JsonResult(500, "图片保存失败")
	}

	//剪切图片
	subImg, err := graphics.ImageCopyFromFile(filePath, x, y, width, height)

	if err != nil {
		logs.Error("ImageCopyFromFile => ", err)
		c.JsonResult(6001, "头像剪切失败")
	}
	os.Remove(filePath)

	filePath = filepath.Join(conf.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName+"_small"+ext)

	err = graphics.ImageResizeSaveFile(subImg, 120, 120, filePath)
	//err = graphics.SaveImage(filePath,subImg)

	if err != nil {
		logs.Error("保存文件失败 => ", err.Error())
		c.JsonResult(500, "保存文件失败")
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, conf.WorkingDirectory), "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}

	if member, err := models.NewMember().Find(c.Member.MemberId); err == nil {
		avater := member.Avatar

		member.Avatar = url
		err := member.Update()
		if err == nil {
			if strings.HasPrefix(avater, "/uploads/") {
				os.Remove(filepath.Join(conf.WorkingDirectory, avater))
			}
			c.SetMember(*member)
		} else {
			c.JsonResult(60001, "保存头像失败")
		}
	}

	c.JsonResult(0, "ok", url)
}
