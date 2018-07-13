package controllers

import (
	"strings"
	"github.com/lifei6671/mindoc/models"
)

type BlogController struct{
	BaseController
}


func (c *BlogController) Index() {
	c.Prepare()
	c.TplName = "blog/index.tpl"
}

func (c *BlogController) List() {
	c.Prepare()
	c.TplName = "blog/list.tpl"

}

//管理后台
func (c *BlogController) ManageList() {
	c.Prepare()
	c.TplName = "blog/manage_list.tpl"
}

//文章设置
func (c *BlogController) ManageSetting() {
	c.Prepare()
	c.TplName = "blog/manage_setting.tpl"
	//如果是post请求
	if c.Ctx.Input.IsPost() {
		blogId,_ := c.GetInt("id",0)
		blogTitle := c.GetString("title")
		blogIdentify := c.GetString("identify")
		orderIndex,_ := c.GetInt("order_index",0)
		blogType,_ := c.GetInt("blog_type",0)
		documentId,_ := c.GetInt("document_id",0)

		if blogTitle == "" {
			c.JsonResult(6001,"文章标题不能为空")
		}
		if blogType != 0 && blogType != 1 {
			c.JsonResult(6005,"未知的文章类型")
		}else if documentId <= 0 && blogType == 1 {
			c.JsonResult(6006,"请选择链接的文章")
		}else if blogType == 1 && documentId > 0 && !models.NewDocument().IsExist(documentId){
			c.JsonResult(6007,"链接的文章不存在")
		}
		if strings.Count(blogTitle,"") > 200 {
			c.JsonResult(6002,"文章标题不能大于200个字符")
		}
		var blog *models.Blog
		var err error
		//如果文章ID存在，则从数据库中查询文章
		if blogId > 0 {
			if c.Member.IsAdministrator() {
				blog, err = models.NewBlog().Find(blogId)
			} else {
				blog, err = models.NewBlog().FindByIdAndMemberId(blogId, c.Member.MemberId)
			}

			if err != nil {
				c.JsonResult(6003, "文章不存在")
			}
			//如果设置了文章标识
			if blogIdentify != "" {
				//如果查询到的文章标识存在并且不是当前文章的id
				if b,err := models.NewBlog().FindByIdentify(blogIdentify); err == nil && b.BlogId != blogId {
					c.JsonResult(6004,"文章标识已存在")
				}
			}
		}else{
			//如果设置了文章标识
			if blogIdentify != "" {
				if models.NewBlog().IsExist(blogIdentify) {
					c.JsonResult(6004,"文章标识已存在")
				}
			}

			blog = models.NewBlog()
			blog.MemberId = c.Member.MemberId
		}


		blog.BlogTitle = blogTitle
		blog.BlogIdentify = blogIdentify
		blog.OrderIndex = orderIndex
		blog.BlogType = blogType
		if blogType == 1 {
			blog.DocumentId = documentId
		}

	}
}

//文章创建或编辑
func (c *BlogController) ManageEdit() {
	c.Prepare()
	c.TplName = "blog/manage_edit.tpl"
}
