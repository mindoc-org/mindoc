package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/controllers"
)

func init()  {
	beego.Router("/",&controllers.HomeController{},"*:Index")

	beego.Router("/manager", &controllers.ManagerController{},"*:Index")

	beego.Router("/setting", &controllers.SettingController{},"*:Index")
	beego.Router("/setting/password", &controllers.SettingController{},"*:Password")

	beego.Router("/book", &controllers.BookController{},"*:Index")
	beego.Router("/book/:key/dashboard", &controllers.BookController{},"*:Dashboard")
	beego.Router("/book/:key/setting", &controllers.BookController{},"*:Setting")
	beego.Router("/book/:key/users", &controllers.BookController{},"*:Users")
	beego.Router("/book/:key/edit/:id", &controllers.BookController{},"*:Edit")
	beego.Router("/book/create", &controllers.BookController{},"*:Create")

	beego.Router("/book/:key/users/create", &controllers.BookMemberController{},"*:Create")
	beego.Router("/book/:key/users/change", &controllers.BookMemberController{},"*:Change")
	beego.Router("/book/:key/users/delete", &controllers.BookMemberController{},"*:Delete")

	beego.Router("/docs/:key", &controllers.DocumentController{},"*:Index")
	beego.Router("/docs/:key/:id", &controllers.DocumentController{},"*:Read")
}
