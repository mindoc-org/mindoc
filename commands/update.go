package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/conf"
	"github.com/astaxie/beego/orm"
)

//系统升级.
func Update() {
	if len(os.Args) >= 2 && os.Args[1] == "update" {

		adapter := beego.AppConfig.String("db_adapter")

		if adapter == "mysql" {
			mysqlUpdate()
		}else if adapter == "sqlite3" {
			sqliteUpdate()
		}
		fmt.Println("update successed.")

		os.Exit(0)
	}
}

//检查最新版本.
func CheckUpdate() {

	if len(os.Args) >= 2 && os.Args[1] == "version" {

		resp, err := http.Get("https://api.github.com/repos/lifei6671/godoc/tags")

		if err != nil {
			beego.Error("CheckUpdate => ", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			beego.Error("CheckUpdate => ", err)
			os.Exit(1)
		}

		var result []*struct {
			Name string `json:"name"`
		}

		err = json.Unmarshal(body, &result)

		if err != nil {
			beego.Error("CheckUpdate => ", err)
			os.Exit(1)
		}
		fmt.Println("MinDoc current version => ", conf.VERSION)
		fmt.Println("MinDoc last version => ", result[0].Name)
		os.Exit(0)
	}
}

//MySQL 数据库更新表结构.
func mysqlUpdate()  {
	sql := `
	IF NOT EXISTS (SELECT * FROM information_schema.columns WHERE table_schema=CurrentDatabase AND table_name = 'md_members' AND column_name = 'auth_method') THEN
		ALTER TABLE md_members ADD auth_method VARCHAR(50) DEFAULT 'local' NULL;
	END IF; `
	o := orm.NewOrm()

	_,err := o.Raw(sql).Exec()

	if err != nil {
		panic(fmt.Sprintf("error : 6001 => %s",err.Error()))
		os.Exit(1)
	}
}

//sqlite 数据库更新表结构.
func sqliteUpdate()  {
	o := orm.NewOrm()

	var sqlite_master struct{
		Name string
	}


	err := o.Raw("select * from sqlite_master where name='md_members' and sql like '%auth_method%' limit 1").QueryRow(&sqlite_master)
	//查询是否已经存在 auth_method 列
	if err == nil && sqlite_master.Name == ""{
		_,err = o.Raw("ALTER TABLE md_members ADD auth_method VARCHAR(50) DEFAULT 'local' NULL;").Exec()
		if err != nil {
			panic(fmt.Sprintf("error : 6001 => %s",err.Error()))
			os.Exit(1)
		}
	}

	os.Exit(1)
}