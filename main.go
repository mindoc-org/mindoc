package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lifei6671/godoc/routers"
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/commands"
	"fmt"
	"os"
	"github.com/lifei6671/godoc/controllers"
)

func main() {

	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RegisterLogger()
	commands.RegisterCommand()
	commands.RegisterFunction()

	beego.SetStaticPath("uploads","uploads")

	fmt.Println(os.Args[0])

	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}

