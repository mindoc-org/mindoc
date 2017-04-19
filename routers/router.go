package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/controllers"
)

func init()  {
	beego.Router("/",&controllers.HomeController{},"*:Index")

	beego.Router("/manager", &controllers.ManagerController{},"*:Index")
}
