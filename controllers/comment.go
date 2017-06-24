package controllers

type CommentController struct {
	BaseController
}

func (c *CommentController) Lists() {

}

func (c *CommentController) Create() {

	c.JsonResult(0, "ok")
}

func (c *CommentController) Index() {
	c.Prepare()
	c.TplName = "comment/index.tpl"
	c.Data["SIDEBAR_ID"] = "comment"
	c.Data["SIDEBAR_BOOK"] = 0
}
