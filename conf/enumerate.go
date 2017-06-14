// package conf 为配置相关.
package conf

import (
	"strings"

	"github.com/astaxie/beego"
)

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

const CaptchaSessionName = "__captcha__"

const RegexpEmail = `^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`

//允许用户名中出现点号

const RegexpAccount = `^[a-zA-Z][a-zA-z0-9\.]{2,50}$`

// PageSize 默认分页条数.
const PageSize = 15

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole = 0
	//普通管理员.
	MemberAdminRole = 1
	//普通用户.
	MemberGeneralRole = 2
)

const (
	// 创始人.
	BookFounder = 0
	//管理者
	BookAdmin = 1
	//编辑者.
	BookEditor = 2
	//观察者
	BookObserver = 3
)

const (
	LoggerOperate = "operate"
	LoggerSystem = "system"
	LoggerException = "exception"
	LoggerDocument = "document"
)
const (
	//本地账户校验
	AuthMethodLocal = "local"
	//LDAP用户校验
	AuthMethodLDAP	= "ldap"
)
var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

// app_key
func GetAppKey() string {
	return beego.AppConfig.DefaultString("app_key", "godoc")
}

func GetDatabasePrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "md_")
}

//获取默认头像
func GetDefaultAvatar() string {
	return beego.AppConfig.DefaultString("avatar", "/static/images/headimgurl.jpg")
}

//获取阅读令牌长度.
func GetTokenSize() int {
	return beego.AppConfig.DefaultInt("token_size", 12)
}

//获取默认文档封面.
func GetDefaultCover() string {
	return beego.AppConfig.DefaultString("cover", "/static/images/book.jpg")
}

//获取允许的商城文件的类型.
func GetUploadFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")

	temp := strings.Split(ext, "|")

	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

//判断是否是允许商城的文件类型.
func IsAllowUploadFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()

	for _, item := range exts {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}
