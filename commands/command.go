package commands

import (
	"fmt"
	"net/url"
	"time"
	"os"
	"encoding/gob"

	"github.com/lifei6671/godoc/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/gocaptcha"

)

// RegisterDataBase 注册数据库
func RegisterDataBase()  {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")

	port := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",username,password,host,port,database,url.QueryEscape(timezone))


	orm.RegisterDataBase("default", "mysql", dataSource)

	location , err := time.LoadLocation(timezone);
	if err == nil {
		orm.DefaultTimeLoc = location
	}else{
		fmt.Println(err)
	}

}

// RegisterModel 注册Model
func RegisterModel()  {
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

//初始化数据
func Initialization()  {

	o := orm.NewOrm()

	_,err := o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用注册','ENABLED_REGISTER','false' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_REGISTER');`).Exec()

	if err != nil {
		panic("ENABLED_REGISTER => " + err.Error())
		os.Exit(1)
	}
	_,err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用验证码','ENABLED_CAPTCHA','false' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_CAPTCHA');`).Exec()

	if err != nil {
		panic("ENABLED_CAPTCHA => " + err.Error())
		os.Exit(1)
	}
	_,err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '启用匿名访问','ENABLE_ANONYMOUS','true' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLE_ANONYMOUS');`).Exec()

	if err != nil {
		panic("ENABLE_ANONYMOUS => " + err.Error())
		os.Exit(1)
	}
	_,err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '站点名称','SITE_NAME','MinDoc' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'SITE_NAME');`).Exec()

	if err != nil {
		panic("SITE_NAME => " + err.Error())
		os.Exit(1)
	}

	member,err := models.NewMember().FindByFieldFirst("account","admin")
	if err == orm.ErrNoRows {

		member.Account = "admin"
		member.Avatar = "/static/images/headimgurl.jpg"
		member.Password = "123456"
		member.Role = 0
		member.Email = "admin@iminho.me"

		if err := member.Add(); err != nil {
			panic("Member.Add => " + err.Error())
			os.Exit(0)
		}

		book := models.NewBook()

		book.MemberId = member.MemberId
		book.BookName = "MinDoc演示项目"
		book.Status = 0
		book.Description = "这是一个MinDoc演示项目，该项目是由系统初始化时自动创建。"
		book.CommentCount = 0
		book.PrivatelyOwned = 0
		book.CommentStatus = "closed"
		book.Identify = "mindoc"
		book.DocCount = 0
		book.CommentCount = 0
		book.Version = time.Now().Unix()
		book.Cover = conf.GetDefaultCover()
		book.Editor = "markdown"
		book.Theme = "default"

		if err := book.Insert(); err != nil {
			panic("Book.Insert => " + err.Error())
			os.Exit(0)
		}
	}
}

// RegisterLogger 注册日志
func RegisterLogger()  {

	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.Async()

	if _,err := os.Stat("logs/log.log"); os.IsNotExist(err) {
		os.MkdirAll("./logs",0777)

		if f,err := os.Create("logs/log.log");err == nil {
			f.Close()
			beego.SetLogger("file",`{"filename":"logs/log.log"}`)
		}
	}

	beego.SetLogFuncCall(true)
	beego.BeeLogger.Async()
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if _,err := os.Stat("install.lock"); os.IsNotExist(err){
		err = orm.RunSyncdb("default",false,true)
		if err == nil {
			Initialization()
			f, _ := os.Create("install.lock")
			defer f.Close()
		}else{
			panic(err.Error())
			os.Exit(0)
		}
	}

	CheckUpdate()
}

func RegisterFunction()  {
	beego.AddFuncMap("config",models.GetOptionValue)
}

func init()  {
	gocaptcha.ReadFonts("./static/fonts", ".ttf")
	gob.Register(models.Member{})
}