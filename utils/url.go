package utils

import (
	"strings"
	"github.com/astaxie/beego"
)

func JoinURI(elem ...string) string {
	if len(elem) <= 0 {
		return ""
	}
	uri := ""

	for i, u := range elem {
		u = strings.Replace(u, "\\", "/", -1)

		if i == 0 {
			if !strings.HasSuffix(u, "/") {
				u = u + "/"
			}
			uri = u
		} else {
			u = strings.Replace(u, "//", "/", -1)
			if strings.HasPrefix(u, "/") {
				u = string(u[1:])
			}
			uri += u
		}
	}
	return uri
}

//重写生成URL的方法，加上完整的域名
func URLFor(endpoint string, values ...interface{}) string {
	baseUrl := beego.AppConfig.DefaultString("baseurl","")
	pathUrl := beego.URLFor(endpoint, values ...)

	if strings.HasPrefix(pathUrl,"http://") {
		return pathUrl
	}
	if strings.HasPrefix(pathUrl,"/") && strings.HasSuffix(baseUrl,"/") {
		return baseUrl + pathUrl[1:]
	}
	if !strings.HasPrefix(pathUrl,"/") && !strings.HasSuffix(baseUrl,"/") {
		return baseUrl + "/" + pathUrl
	}
	return  baseUrl + beego.URLFor(endpoint, values ...)
}

func URLForWithCdnImage(p string) string  {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	cdn := beego.AppConfig.DefaultString("cdnimg", "")
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl","")

		if strings.HasPrefix(p,"/") && strings.HasSuffix(baseUrl,"/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p,"/") && !strings.HasSuffix(baseUrl,"/") {
			return baseUrl + "/" + p
		}
		return  baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnCss (p string) string {
	cdn := beego.AppConfig.DefaultString("cdncss", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl","")

		if strings.HasPrefix(p,"/") && strings.HasSuffix(baseUrl,"/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p,"/") && !strings.HasSuffix(baseUrl,"/") {
			return baseUrl + "/" + p
		}
		return  baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnJs(p string) string {
	cdn := beego.AppConfig.DefaultString("cdnjs", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl","")

		if strings.HasPrefix(p,"/") && strings.HasSuffix(baseUrl,"/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p,"/") && !strings.HasSuffix(baseUrl,"/") {
			return baseUrl + "/" + p
		}
		return  baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}