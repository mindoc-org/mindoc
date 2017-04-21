// package conf 为配置相关.
package conf

import "github.com/astaxie/beego"

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

// app_key
func GetAppKey()  (string) {
	return beego.AppConfig.DefaultString("app_key","go-git-webhook")
}

func GetDatabasePrefix() string  {
	return beego.AppConfig.DefaultString("db_prefix","md_")
}