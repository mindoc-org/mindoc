package controllers

type HomeController struct {
	BaseController
}

func (p *HomeController) Index() {
	p.TplName = "home/index.tpl"
}
