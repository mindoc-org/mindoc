package commands

import (
	"fmt"
	"net/url"
	"time"

	"github.com/lifei6671/godoc/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"os"
	"github.com/lifei6671/godoc/conf"
)

// RegisterDataBase 注册数据库
func RegisterDataBase()  {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")

	port := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s",username,password,host,port,database,url.QueryEscape(timezone))


	orm.RegisterDataBase("default", "mysql", dataSource)

	orm.DefaultTimeLoc, _ = time.LoadLocation(timezone)

}

// RegisterModel 注册Model
func RegisterModel()  {
	orm.RegisterModelWithPrefix(conf.GetDatabasePrefix(),
		new(models.Member),
		new(models.Book),
		new(models.Relationship),
		new(models.Comment),
		new(models.Option),
		new(models.Document),
		new(models.Attachment),
		new(models.Logger),
	)

}

func Initialization()  {

	options := []models.Option {
		{ OptionName: "ENABLED_CAPTCHA", OptionValue: "false", OptionTitle:"是否启用验证码"},
		{ OptionName: "ENABLED_REGISTER",OptionValue:"false",OptionTitle:"是否启用注册"},
		{ OptionName: "ENABLE_ANONYMOUS" , OptionValue:"false", OptionTitle:"启用匿名访问"},
		{ OptionName: "SITE_NAME", OptionValue:"MinDoc", OptionTitle: "站点名称"},
	}

	models.NewOption().InsertMulti(options...)

	member := models.NewMember()
	member.Account = "admin"
	member.Avatar = "/static/images/headimgurl.jpg"
	member.Password = "123456"
	member.Role = 0

	member.Add()
}

// RegisterLogger 注册日志
func RegisterLogger()  {

	logs.SetLogger("console")
	logs.SetLogger("file",`{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if _,err := os.Stat("install.lock"); os.IsNotExist(err){
		orm.RunSyncdb("default",true,false)
		Initialization()
		f,_ := os.Create("install.lock")
		defer f.Close()
	}

}

func RegisterFunction()  {
	beego.AddFuncMap("config",models.GetOptionValue)
}