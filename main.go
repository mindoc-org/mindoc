package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lifei6671/godoc/routers"
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/commands"
)

func main() {
	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RegisterLogger()
	commands.RegisterCommand()
	commands.RegisterFunction()

	beego.SetStaticPath("uploads","uploads")


	beego.Run()
}
