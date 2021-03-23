package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/astaxie/beego/session/memcache"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/kardianos/service"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mindoc-org/mindoc/commands"
	"github.com/mindoc-org/mindoc/commands/daemon"
	_ "github.com/mindoc-org/mindoc/routers"
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

	if err := s.Run(); err != nil {
		log.Fatal("启动程序失败 ->", err)
	}
}
