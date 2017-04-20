package controllers

type BookController struct {
	BaseController
}

func (p *BookController) Index() {
	p.TplName = "book/index.tpl"
}

// Dashboard 项目概要 .
func (p *BookController) Dashboard() {
	p.TplName = "book/dashboard.tpl"
}

// Setting 项目设置 .
func (p *BookController) Setting()  {
	p.TplName = "book/setting.tpl"
}

func (p *BookController) Users() {
	p.TplName = "book/users.tpl"
}

func (p *BookController) Create() {
	p.TplName = "book/create.tpl"
}

// Edit 编辑项目.
func (p *BookController) Edit() {
	p.TplName = "book/edit.tpl"
}

// Delete 删除项目.
func (p *BookController) Delete() {
	p.StopRun()
}

// Transfer 转让项目.
func (p *BookController)Transfer()  {

}