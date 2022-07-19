// package conf 为配置相关.
package conf

import (
	"strings"

	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

const CaptchaSessionName = "__captcha__"

const RegexpEmail = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//允许用户名中出现点号
const RegexpAccount = `^[a-zA-Z0-9][a-zA-Z0-9\.-]{2,50}$`

// PageSize 默认分页条数.
const PageSize = 10

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole SystemRole = iota
	//普通管理员.
	MemberAdminRole
	//普通用户.
	MemberGeneralRole
)

//系统角色
type SystemRole int

const (
	// 创始人.
	BookFounder BookRole = iota
	//管理者
	BookAdmin
	//编辑者.
	BookEditor
	//观察者
	BookObserver
	//未指定关系
	BookRoleNoSpecific
)

//项目角色
type BookRole int

const (
	LoggerOperate   = "operate"
	LoggerSystem    = "system"
	LoggerException = "exception"
	LoggerDocument  = "document"
)
const (
	//本地账户校验
	AuthMethodLocal = "local"
	//LDAP用户校验
	AuthMethodLDAP = "ldap"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

var (
	ConfigurationFile = "./conf/app.conf"
	WorkingDirectory  = "./"
	LogFile           = "./runtime/logs"
	BaseUrl           = ""
	AutoLoadDelay     = 0
)

// app_key
func GetAppKey() string {
	return web.AppConfig.DefaultString("app_key", "mindoc")
}

func GetDatabasePrefix() string {
	return web.AppConfig.DefaultString("db_prefix", "md_")
}

//获取默认头像
func GetDefaultAvatar() string {
	return URLForWithCdnImage(web.AppConfig.DefaultString("avatar", "/static/images/headimgurl.jpg"))
}

//获取阅读令牌长度.
func GetTokenSize() int {
	return web.AppConfig.DefaultInt("token_size", 12)
}

//获取默认文档封面.
func GetDefaultCover() string {

	return URLForWithCdnImage(web.AppConfig.DefaultString("cover", "/static/images/book.jpg"))
}

//获取允许的商城文件的类型.
func GetUploadFileExt() []string {
	ext := web.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")

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

// 获取上传文件允许的最大值
func GetUploadFileSize() int64 {
	size := web.AppConfig.DefaultString("upload_file_size", "0")

	if strings.HasSuffix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "GB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "KB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024
		}
	}
	if s, e := strconv.ParseInt(size, 10, 64); e == nil {
		return s * 1024
	}
	return 0
}

//是否启用导出
func GetEnableExport() bool {
	return web.AppConfig.DefaultBool("enable_export", true)
}

//是否启用iframe
func GetEnableIframe() bool {
	return web.AppConfig.DefaultBool("enable_iframe", false)
}

//同一项目导出线程的并发数
func GetExportProcessNum() int {
	exportProcessNum := web.AppConfig.DefaultInt("export_process_num", 1)

	if exportProcessNum <= 0 || exportProcessNum > 4 {
		exportProcessNum = 1
	}
	return exportProcessNum
}

//导出项目队列的并发数量
func GetExportLimitNum() int {
	exportLimitNum := web.AppConfig.DefaultInt("export_limit_num", 1)

	if exportLimitNum < 0 {
		exportLimitNum = 1
	}
	return exportLimitNum
}

//等待导出队列的长度
func GetExportQueueLimitNum() int {
	exportQueueLimitNum := web.AppConfig.DefaultInt("export_queue_limit_num", 10)

	if exportQueueLimitNum <= 0 {
		exportQueueLimitNum = 100
	}
	return exportQueueLimitNum
}

//默认导出项目的缓存目录
func GetExportOutputPath() string {
	exportOutputPath := filepath.Join(web.AppConfig.DefaultString("export_output_path", filepath.Join(WorkingDirectory, "cache")), "books")

	return exportOutputPath
}

//判断是否是允许商城的文件类型.
func IsAllowUploadFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()

	for _, item := range exts {
		if item == "*" {
			return true
		}
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

//读取配置文件值
func CONF(key string, value ...string) string {
	defaultValue := ""
	if len(value) > 0 {
		defaultValue = value[0]
	}
	return web.AppConfig.DefaultString(key, defaultValue)
}

//重写生成URL的方法，加上完整的域名
func URLFor(endpoint string, values ...interface{}) string {
	baseUrl := web.AppConfig.DefaultString("baseurl", "")
	pathUrl := web.URLFor(endpoint, values...)

	if baseUrl == "" {
		baseUrl = BaseUrl
	}
	if strings.HasPrefix(pathUrl, "http://") {
		return pathUrl
	}
	if strings.HasPrefix(pathUrl, "/") && strings.HasSuffix(baseUrl, "/") {
		return baseUrl + pathUrl[1:]
	}
	if !strings.HasPrefix(pathUrl, "/") && !strings.HasSuffix(baseUrl, "/") {
		return baseUrl + "/" + pathUrl
	}
	return baseUrl + web.URLFor(endpoint, values...)
}

func URLForNotHost(endpoint string, values ...interface{}) string {
	baseUrl := web.AppConfig.DefaultString("baseurl", "")
	pathUrl := web.URLFor(endpoint, values...)

	if baseUrl == "" {
		baseUrl = "/"
	}
	if strings.HasPrefix(pathUrl, "http://") {
		return pathUrl
	}
	if strings.HasPrefix(pathUrl, "/") && strings.HasSuffix(baseUrl, "/") {
		return baseUrl + pathUrl[1:]
	}
	if !strings.HasPrefix(pathUrl, "/") && !strings.HasSuffix(baseUrl, "/") {
		return baseUrl + "/" + pathUrl
	}
	return baseUrl + web.URLFor(endpoint, values...)
}

func URLForWithCdnImage(p string) string {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	cdn := web.AppConfig.DefaultString("cdnimg", "")
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := web.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnCss(p string, v ...string) string {
	cdn := web.AppConfig.DefaultString("cdncss", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	filePath := WorkingDir(p)

	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := web.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnJs(p string, v ...string) string {
	cdn := web.AppConfig.DefaultString("cdnjs", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}

	filePath := WorkingDir(p)

	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}

	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := web.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func WorkingDir(elem ...string) string {

	elems := append([]string{WorkingDirectory}, elem...)

	return filepath.Join(elems...)
}

func init() {
	if p, err := filepath.Abs("./conf/app.conf"); err == nil {
		ConfigurationFile = p
	}
	if p, err := filepath.Abs("./"); err == nil {
		WorkingDirectory = p
	}
	if p, err := filepath.Abs("./runtime/logs"); err == nil {
		LogFile = p
	}
}
