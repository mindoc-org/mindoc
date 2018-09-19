package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "*:Index")

	beego.Router("/login", &controllers.AccountController{}, "*:Login")
	beego.Router("/logout", &controllers.AccountController{}, "*:Logout")
	beego.Router("/register", &controllers.AccountController{}, "*:Register")
	beego.Router("/find_password", &controllers.AccountController{}, "*:FindPassword")
	beego.Router("/valid_email", &controllers.AccountController{}, "post:ValidEmail")
	beego.Router("/captcha", &controllers.AccountController{}, "*:Captcha")

	beego.Router("/manager", &controllers.ManagerController{}, "*:Index")
	beego.Router("/manager/users", &controllers.ManagerController{}, "*:Users")
	beego.Router("/manager/users/edit/:id", &controllers.ManagerController{}, "*:EditMember")
	beego.Router("/manager/member/create", &controllers.ManagerController{}, "post:CreateMember")
	beego.Router("/manager/member/delete", &controllers.ManagerController{}, "post:DeleteMember")
	beego.Router("/manager/member/update-member-status", &controllers.ManagerController{}, "post:UpdateMemberStatus")
	beego.Router("/manager/member/change-member-role", &controllers.ManagerController{}, "post:ChangeMemberRole")
	beego.Router("/manager/books", &controllers.ManagerController{}, "*:Books")
	beego.Router("/manager/books/edit/:key", &controllers.ManagerController{}, "*:EditBook")
	beego.Router("/manager/books/delete", &controllers.ManagerController{}, "*:DeleteBook")

	beego.Router("/manager/comments", &controllers.ManagerController{}, "*:Comments")
	beego.Router("/manager/setting", &controllers.ManagerController{}, "*:Setting")
	beego.Router("/manager/books/token", &controllers.ManagerController{}, "post:CreateToken")
	beego.Router("/manager/books/transfer", &controllers.ManagerController{}, "post:Transfer")
	beego.Router("/manager/books/open", &controllers.ManagerController{}, "post:PrivatelyOwned")

	beego.Router("/manager/attach/list", &controllers.ManagerController{}, "*:AttachList")
	beego.Router("/manager/attach/detailed/:id", &controllers.ManagerController{}, "*:AttachDetailed")
	beego.Router("/manager/attach/delete", &controllers.ManagerController{}, "post:AttachDelete")
	beego.Router("/manager/label/list", &controllers.ManagerController{},"get:LabelList")
	beego.Router("/manager/label/delete/:id", &controllers.ManagerController{},"post:LabelDelete")

	beego.Router("/manager/config",  &controllers.ManagerController{}, "*:Config")

	beego.Router("/setting", &controllers.SettingController{}, "*:Index")
	beego.Router("/setting/password", &controllers.SettingController{}, "*:Password")
	beego.Router("/setting/upload", &controllers.SettingController{}, "*:Upload")

	beego.Router("/book", &controllers.BookController{}, "*:Index")
	beego.Router("/book/:key/dashboard", &controllers.BookController{}, "*:Dashboard")
	beego.Router("/book/:key/setting", &controllers.BookController{}, "*:Setting")
	beego.Router("/book/:key/users", &controllers.BookController{}, "*:Users")
	beego.Router("/book/:key/release", &controllers.BookController{}, "post:Release")
	beego.Router("/book/:key/sort", &controllers.BookController{}, "post:SaveSort")

	beego.Router("/book/create", &controllers.BookController{}, "*:Create")
	beego.Router("/book/users/create", &controllers.BookMemberController{}, "post:AddMember")
	beego.Router("/book/users/change", &controllers.BookMemberController{}, "post:ChangeRole")
	beego.Router("/book/users/delete", &controllers.BookMemberController{}, "post:RemoveMember")
	beego.Router("/book/users/import", &controllers.BookController{},"post:Import")
	beego.Router("/book/users/copy", &controllers.BookController{},"post:Copy")

	beego.Router("/book/setting/save", &controllers.BookController{}, "post:SaveBook")
	beego.Router("/book/setting/open", &controllers.BookController{}, "post:PrivatelyOwned")
	beego.Router("/book/setting/transfer", &controllers.BookController{}, "post:Transfer")
	beego.Router("/book/setting/upload", &controllers.BookController{}, "post:UploadCover")
	beego.Router("/book/setting/token", &controllers.BookController{}, "post:CreateToken")
	beego.Router("/book/setting/delete", &controllers.BookController{}, "post:Delete")

	//管理文章的路由
	beego.Router("/manage/blogs", &controllers.BlogController{},"*:ManageList")
	beego.Router("/manage/blogs/setting/?:id", &controllers.BlogController{}, "*:ManageSetting")
	beego.Router("/manage/blogs/edit/?:id",&controllers.BlogController{}, "*:ManageEdit")
	beego.Router("/manage/blogs/delete",&controllers.BlogController{}, "post:ManageDelete")
	beego.Router("/manage/blogs/upload",&controllers.BlogController{}, "post:Upload")
	beego.Router("/manage/blogs/attach/:id",&controllers.BlogController{}, "post:RemoveAttachment")


	//读文章的路由
	beego.Router("/blogs", &controllers.BlogController{}, "*:List")
	beego.Router("/blog-attach/:id:int/:attach_id:int", &controllers.BlogController{},"get:Download")
	beego.Router("/blog-:id([0-9]+).html",&controllers.BlogController{}, "*:Index")

	//模板相关接口
	beego.Router("/api/template/get", &controllers.TemplateController{},"get:Get")
	beego.Router("/api/template/list", &controllers.TemplateController{},"post:List")
	beego.Router("/api/template/add", &controllers.TemplateController{},"post:Add")
	beego.Router("/api/template/remove", &controllers.TemplateController{},"post:Delete")

	beego.Router("/api/attach/remove/", &controllers.DocumentController{}, "post:RemoveAttachment")
	beego.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	beego.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	beego.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	beego.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")
	beego.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	beego.Router("/api/:key/compare/:id", &controllers.DocumentController{}, "*:Compare")
	beego.Router("/api/search/user/:key", &controllers.SearchController{}, "*:User")

	beego.Router("/history/get", &controllers.DocumentController{}, "get:History")
	beego.Router("/history/delete", &controllers.DocumentController{}, "*:DeleteHistory")
	beego.Router("/history/restore", &controllers.DocumentController{}, "*:RestoreHistory")

	beego.Router("/docs/:key", &controllers.DocumentController{}, "*:Index")
	beego.Router("/docs/:key/:id", &controllers.DocumentController{}, "*:Read")
	beego.Router("/docs/:key/search", &controllers.DocumentController{}, "post:Search")
	beego.Router("/export/:key", &controllers.DocumentController{}, "*:Export")
	beego.Router("/qrcode/:key.png", &controllers.DocumentController{}, "get:QrCode")

	beego.Router("/attach_files/:key/:attach_id", &controllers.DocumentController{}, "get:DownloadAttachment")

	beego.Router("/comment/create", &controllers.CommentController{}, "post:Create")
	beego.Router("/comment/lists", &controllers.CommentController{}, "get:Lists")
	beego.Router("/comment/index", &controllers.CommentController{}, "*:Index")

	beego.Router("/search", &controllers.SearchController{}, "get:Index")

	beego.Router("/tag/:key", &controllers.LabelController{}, "get:Index")
	beego.Router("/tags", &controllers.LabelController{}, "get:List")
}
