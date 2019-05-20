package main

import (
	"fmt"
	"os"

	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/kardianos/service"
	"github.com/lifei6671/mindoc/commands"
	"github.com/lifei6671/mindoc/commands/daemon"
	_ "github.com/lifei6671/mindoc/routers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	if len(os.Args) >= 3 && os.Args[1] == "service" {
		if os.Args[2] == "install" {
			daemon.Install()
		} else if os.Args[2] == "remove" {
			daemon.Uninstall()
		} else if os.Args[2] == "restart" {
			daemon.Restart()
		}
	}
	commands.RegisterCommand()

	d := daemon.NewDaemon()

	s, err := service.New(d, d.Config())

	if err != nil {
		fmt.Println("Create service error => ", err)
		os.Exit(1)
	}

	s.Run()
}
