package controllers

type SettingController struct {
	BaseController
}

func (p *SettingController) Index()  {
	p.TplName = "setting/index.tpl"
}

func (p *SettingController) Password()  {
	p.TplName = "setting/password.tpl"
}