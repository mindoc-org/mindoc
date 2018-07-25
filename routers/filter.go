package routers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"net/url"
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
				ctx.Redirect(302, conf.URLFor("AccountController.Login") + "?url=" +  url.PathEscape(conf.BaseUrl + ctx.Request.URL.RequestURI()))
			}
		}
	}
	beego.InsertFilter("/manager", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/manager/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/setting", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/setting/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/book", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/book/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/api/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/manage/*", beego.BeforeRouter,FilterUser)

	var FinishRouter = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Add("MinDoc-Version", conf.VERSION)
		ctx.ResponseWriter.Header().Add("MinDoc-Site", "https://www.iminho.me")
	}

	beego.InsertFilter("/*", beego.BeforeRouter, FinishRouter, false)
}
