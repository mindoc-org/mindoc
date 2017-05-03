package controllers

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	c.TplName = "errors/404.tpl"
}

func (c *ErrorController) Error403() {
	c.TplName = "errors/403.tpl"
}

func (c *ErrorController) Error500() {
	c.TplName = "errors/error.tpl"
}
