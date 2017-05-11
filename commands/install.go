package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/models"
)

//系统安装.
func Install() {
	if len(os.Args) >= 2 && os.Args[1] == "install" {
		fmt.Println("Initializing...")

		err := orm.RunSyncdb("default", false, true)
		if err == nil {
			initialization()
		} else {
			panic(err.Error())
			os.Exit(1)
		}
		fmt.Println("Install Successfully!")
		os.Exit(0)
	}
}

//初始化数据
func initialization() {

	o := orm.NewOrm()

	_, err := o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用注册','ENABLED_REGISTER','false' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_REGISTER');`).Exec()

	if err != nil {
		panic("ENABLED_REGISTER => " + err.Error())
		os.Exit(1)
	}
	_, err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用验证码','ENABLED_CAPTCHA','false' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_CAPTCHA');`).Exec()

	if err != nil {
		panic("ENABLED_CAPTCHA => " + err.Error())
		os.Exit(1)
	}
	_, err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '启用匿名访问','ENABLE_ANONYMOUS','true' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLE_ANONYMOUS');`).Exec()

	if err != nil {
		panic("ENABLE_ANONYMOUS => " + err.Error())
		os.Exit(1)
	}
	_, err = o.Raw(`INSERT INTO md_options (option_title, option_name, option_value) SELECT '站点名称','SITE_NAME','MinDoc' FROM dual WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'SITE_NAME');`).Exec()

	if err != nil {
		panic("SITE_NAME => " + err.Error())
		os.Exit(1)
	}

	member, err := models.NewMember().FindByFieldFirst("account", "admin")
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
