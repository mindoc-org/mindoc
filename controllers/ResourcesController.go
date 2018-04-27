package controllers

type ResourcesController struct {
	ManagerController
}

func (c *ResourcesController) AddResource() {
	c.Prepare()
	c.TplName = "resources/AddResource.tpl"

}