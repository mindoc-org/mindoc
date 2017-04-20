package controllers

type DocumentController struct {
	BaseController
}

func (p *DocumentController) Index()  {
	p.TplName = "document/index.tpl"
}

func (p *DocumentController) Read() {
	p.TplName = "document/kancloud.tpl"
}
