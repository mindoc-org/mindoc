package controllers


type BlogController struct{
	BaseController
}


func (c *BlogController) Index() {
	c.Prepare()

}

func (c *BlogController) List() {
	c.Prepare()
}

//管理后台
func (c *BlogController) ManageList() {
	c.Prepare()
}

//文章设置
func (c *BlogController) ManageSetting() {
	c.Prepare()
}

//文章创建或编辑
func (c *BlogController) ManageEdit() {
	c.Prepare()
}
