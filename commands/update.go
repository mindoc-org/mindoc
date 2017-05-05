package commands

import (
	"os"
	"net/http"
	"github.com/astaxie/beego"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func Update()  {
	if len(os.Args) > 2 && os.Args[1] == "update" {

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

		fmt.Println("MinDoc last version => ",result[0].Name)
		os.Exit(0)
	}
}