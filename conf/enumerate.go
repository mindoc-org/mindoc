// package conf 为配置相关.
package conf

import "github.com/astaxie/beego"

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

const RegexpEmail  = `^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`

const RegexpAccount = `^[a-zA-Z][a-zA-z0-9]{2,50}$`

// app_key
func GetAppKey()  (string) {
	return beego.AppConfig.DefaultString("app_key","go-git-webhook")
}

func GetDatabasePrefix() string  {
	return beego.AppConfig.DefaultString("db_prefix","md_")
}
//获取默认头像
func GetDefaultAvatar() string {
	return beego.AppConfig.DefaultString("avatar","/static/images/headimgurl.jpg")
}