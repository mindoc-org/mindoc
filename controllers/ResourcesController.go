package controllers

type ResourcesController struct {
	ManagerController
}

func (c *ResourcesController) ResourceList() {
	c.Prepare()
	c.TplName = "resources/ResourceList.tpl"
}

func (c *ResourcesController) AddResource() {
	c.Prepare()
	c.TplName = "resources/AddResource.tpl"

}

