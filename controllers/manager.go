package controllers

type ManagerController struct {
	BaseController
}

func (p *ManagerController) Index() {
	p.TplName = "manager/index.tpl"
}

func (p *ManagerController) Users()  {
	p.TplName = "manager/users.tpl"
}