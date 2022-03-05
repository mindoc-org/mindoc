package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
)

//检查最新版本.
func CheckUpdate() {

	fmt.Println("MinDoc current version => ", conf.VERSION)

	resp, err := http.Get("https://api.github.com/repos/mindoc-org/mindoc/tags")

	if err != nil {
		logs.Error("CheckUpdate => ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("CheckUpdate => ", err)
		os.Exit(1)
	}

	var result []*struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logs.Error("CheckUpdate => ", err)
		os.Exit(0)
	}

	if len(result) > 0 {
		fmt.Println("MinDoc last version => ", result[0].Name)
	}

	os.Exit(0)

}
