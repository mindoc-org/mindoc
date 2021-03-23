package routers

import (
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/adapter/context"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
)

func init() {
	var FilterUser = func(ctx *context.Context) {
		_, ok := ctx.Input.Session(conf.LoginSessionName).(models.Member)

		if !ok {
			if ctx.Input.IsAjax() {
				jsonData := make(map[string]interface{}, 3)

				jsonData["errcode"] = 403
				jsonData["message"] = "请登录后再操作"

				returnJSON, _ := json.Marshal(jsonData)

				ctx.ResponseWriter.Write(returnJSON)
			} else {
				ctx.Redirect(302, conf.URLFor("AccountController.Login")+"?url="+url.PathEscape(conf.BaseUrl+ctx.Request.URL.RequestURI()))
			}
		}
	}
	adapter.InsertFilter("/manager", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/manager/*", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/setting", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/setting/*", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/book", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/book/*", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/api/*", web.BeforeRouter, FilterUser)
	adapter.InsertFilter("/manage/*", web.BeforeRouter, FilterUser)

	var FinishRouter = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Add("MinDoc-Version", conf.VERSION)
		ctx.ResponseWriter.Header().Add("MinDoc-Site", "https://www.iminho.me")
		ctx.ResponseWriter.Header().Add("X-XSS-Protection", "1; mode=block")
	}

	var StartRouter = func(ctx *context.Context) {
		sessionId := ctx.Input.Cookie(adapter.AppConfig.String("sessionname"))
		if sessionId != "" {
			//sessionId必须是数字字母组成，且最小32个字符，最大1024字符
			if ok, err := regexp.MatchString(`^[a-zA-Z0-9]{32,512}$`, sessionId); !ok || err != nil {
				panic("401")
			}
		}
	}
	adapter.InsertFilter("/*", web.BeforeStatic, StartRouter, false)
	adapter.InsertFilter("/*", web.BeforeRouter, FinishRouter, false)
}
