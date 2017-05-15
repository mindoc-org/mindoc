package main

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lifei6671/godoc/commands"
	"github.com/lifei6671/godoc/conf"
	"github.com/lifei6671/godoc/controllers"
	_ "github.com/lifei6671/godoc/routers"
)

func main() {

	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RegisterLogger()
	commands.RegisterCommand()
	commands.RegisterFunction()

	
	beego.SetStaticPath("uploads", "uploads")

	beego.ErrorController(&controllers.ErrorController{})

	fmt.Printf("MinDoc version => %s\nbuild time => %s\nstart directory => %s\n%s\n", conf.VERSION, conf.BUILD_TIME, os.Args[0], conf.GO_VERSION)

	beego.Run()
}
