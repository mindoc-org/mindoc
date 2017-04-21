package controllers

import (
	"image"
	"fmt"
	"os"
	"strings"
	"image/jpeg"
	"image/png"
	"image/gif"
	"path/filepath"
	"strconv"
	"time"
	"io/ioutil"
	"bytes"

	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/models"
	"github.com/lifei6671/godoc/utils"
)

type SettingController struct {
	BaseController
}

func (c *SettingController) Index()  {
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
		if err := member.Update(); err != nil {
			c.JsonResult(602, err.Error())
		}
		c.JsonResult(0, "ok")
	}
}

func (c *SettingController) Password()  {
	c.TplName = "setting/password.tpl"

	if c.Ctx.Input.IsPost() {
		password1 := c.GetString("password1")
		password2 := c.GetString("password2")
		password3 := c.GetString("password3")

		if password1 == "" {
			c.JsonResult(6003,"原密码不能为空")
		}

		if password2 == "" {
			c.JsonResult(6004,"新密码不能为空")
		}
		if count := strings.Count(password2,""); count < 6 || count > 18 {
			c.JsonResult(6009,"密码必须在6-18字之间")
		}
		if password2 != password3 {
			c.JsonResult(6003,"确认密码不正确")
		}
		if ok,_ := utils.PasswordVerify(c.Member.Password,password1) ; !ok {
			c.JsonResult(6005,"原始密码不正确")
		}
		if password1 == password2 {
			c.JsonResult(6006,"新密码不能和原始密码相同")
		}
		pwd,err := utils.PasswordHash(password2)
		if err != nil {
			c.JsonResult(6007,"密码加密失败")
		}
		c.Member.Password = pwd
		if err := c.Member.Update();err != nil {
			c.JsonResult(6008,err.Error())
		}
		c.JsonResult(0,"ok")
	}
}

// Upload 上传图片
func (c *SettingController) Upload() {
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

	fmt.Println(x,x1,y,y1)

	fileName := "avatar_" +  strconv.FormatInt(int64(time.Now().Nanosecond()), 16)

	filePath := "static/uploads/" + time.Now().Format("200601") + "/" + fileName + ext

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file",filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}

	fileBytes,err := ioutil.ReadFile(filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}

	buf := bytes.NewBuffer(fileBytes)

	m,_,err := image.Decode(buf)

	if err != nil{
		logs.Error("image.Decode => ",err)
		c.JsonResult(500,"图片解码失败")
	}


	var subImg image.Image

	if rgbImg,ok := m.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	}else if rgbImg,ok := m.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	}else if rgbImg,ok := m.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else {
		fmt.Println(m)
		c.JsonResult(500,"图片解码失败")
	}

	f, err := os.OpenFile("./" + filePath, os.O_SYNC | os.O_RDWR, 0666)

	if err != nil{
		c.JsonResult(500,"保存图片失败")
	}
	defer f.Close()

	if strings.EqualFold(ext,".jpg") || strings.EqualFold(ext,".jpeg"){

		err = jpeg.Encode(f,subImg,&jpeg.Options{ Quality : 100 })
	}else if strings.EqualFold(ext,".png") {
		err = png.Encode(f,subImg)
	}else if strings.EqualFold(ext,".gif") {
		err = gif.Encode(f,subImg,&gif.Options{ NumColors : 256})
	}
	if err != nil {
		logs.Error("图片剪切失败 => ",err.Error())
		c.JsonResult(500,"图片剪切失败")
	}


	if err != nil {
		logs.Error("保存文件失败 => ",err.Error())
		c.JsonResult(500,"保存文件失败")
	}
	url := "/" + filePath

	member := models.NewMember()
	if err := member.Find(c.Member.MemberId);err == nil {
		member.Avatar = url
		member.Update()
		c.SetMember(*member)
	}

	c.JsonResult(0,"ok",url)
}