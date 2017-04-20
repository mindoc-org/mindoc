package controllers

type BookMemberController struct {
	BaseController
}

func (p *BookMemberController) Create()  {
	p.TplName = "book/member_create.tpl"
}

func (p *BookMemberController) Change() {

	p.StopRun()
}

func (p *BookMemberController) Delete()  {

	p.StopRun()
}