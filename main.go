package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lifei6671/godoc/routers"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/mysql"
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/commands"
	"fmt"
	"os"
	"github.com/lifei6671/godoc/controllers"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

func main() {

	fmt.Printf("MinDoc version%s\n%s\n%s\nstart directory%s\n", VERSION, BUILD_TIME, GO_VERSION,os.Args[0])

	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RegisterLogger()
	commands.RegisterCommand()
	commands.RegisterFunction()

	beego.SetStaticPath("uploads","uploads")



	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}

