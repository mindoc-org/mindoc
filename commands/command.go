package commands

import (
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/gocaptcha"
	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/models"
)

// RegisterDataBase 注册数据库
func RegisterDataBase() {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")

	port := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

	orm.RegisterDataBase("default", "mysql", dataSource)

	location, err := time.LoadLocation(timezone)
	if err == nil {
		orm.DefaultTimeLoc = location
	} else {
		fmt.Println(err)
	}

}

// RegisterModel 注册Model
func RegisterModel() {
	orm.RegisterModelWithPrefix(conf.GetDatabasePrefix(),
		new(models.Member),
		new(models.Book),
		new(models.Relationship),
		//new(models.Comment),
		new(models.Option),
		new(models.Document),
		new(models.Attachment),
		new(models.Logger),
		//new(models.CommentVote),
		new(models.MemberToken),
	)

}

// RegisterLogger 注册日志
func RegisterLogger() {

	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.Async()

	if _, err := os.Stat("logs/log.log"); os.IsNotExist(err) {
		os.MkdirAll("./logs", 0777)

		if f, err := os.Create("logs/log.log"); err == nil {
			f.Close()
			beego.SetLogger("file", `{"filename":"logs/log.log"}`)
		}
	}

	beego.SetLogFuncCall(true)
	beego.BeeLogger.Async()
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	Install()
	Update()
	CheckUpdate()
}

func RegisterFunction() {
	beego.AddFuncMap("config", models.GetOptionValue)
}

func init() {
	gocaptcha.ReadFonts("./static/fonts", ".ttf")
	gob.Register(models.Member{})
}
