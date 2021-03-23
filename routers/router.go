package routers

import (
	"github.com/beego/beego/v2/adapter"
	"github.com/mindoc-org/mindoc/controllers"
)

func init() {
	adapter.Router("/", &controllers.HomeController{}, "*:Index")

	adapter.Router("/login", &controllers.AccountController{}, "*:Login")
	adapter.Router("/logout", &controllers.AccountController{}, "*:Logout")
	adapter.Router("/register", &controllers.AccountController{}, "*:Register")
	adapter.Router("/find_password", &controllers.AccountController{}, "*:FindPassword")
	adapter.Router("/valid_email", &controllers.AccountController{}, "post:ValidEmail")
	adapter.Router("/captcha", &controllers.AccountController{}, "*:Captcha")

	adapter.Router("/manager", &controllers.ManagerController{}, "*:Index")
	adapter.Router("/manager/users", &controllers.ManagerController{}, "*:Users")
	adapter.Router("/manager/users/edit/:id", &controllers.ManagerController{}, "*:EditMember")
	adapter.Router("/manager/member/create", &controllers.ManagerController{}, "post:CreateMember")
	adapter.Router("/manager/member/delete", &controllers.ManagerController{}, "post:DeleteMember")
	adapter.Router("/manager/member/update-member-status", &controllers.ManagerController{}, "post:UpdateMemberStatus")
	adapter.Router("/manager/member/change-member-role", &controllers.ManagerController{}, "post:ChangeMemberRole")
	adapter.Router("/manager/books", &controllers.ManagerController{}, "*:Books")
	adapter.Router("/manager/books/edit/:key", &controllers.ManagerController{}, "*:EditBook")
	adapter.Router("/manager/books/delete", &controllers.ManagerController{}, "*:DeleteBook")

	adapter.Router("/manager/comments", &controllers.ManagerController{}, "*:Comments")
	adapter.Router("/manager/setting", &controllers.ManagerController{}, "*:Setting")
	adapter.Router("/manager/books/token", &controllers.ManagerController{}, "post:CreateToken")
	adapter.Router("/manager/books/transfer", &controllers.ManagerController{}, "post:Transfer")
	adapter.Router("/manager/books/open", &controllers.ManagerController{}, "post:PrivatelyOwned")

	adapter.Router("/manager/attach/list", &controllers.ManagerController{}, "*:AttachList")
	adapter.Router("/manager/attach/detailed/:id", &controllers.ManagerController{}, "*:AttachDetailed")
	adapter.Router("/manager/attach/delete", &controllers.ManagerController{}, "post:AttachDelete")
	adapter.Router("/manager/label/list", &controllers.ManagerController{}, "get:LabelList")
	adapter.Router("/manager/label/delete/:id", &controllers.ManagerController{}, "post:LabelDelete")

	//adapter.Router("/manager/config",  &controllers.ManagerController{}, "*:Config")

	adapter.Router("/manager/team", &controllers.ManagerController{}, "*:Team")
	adapter.Router("/manager/team/create", &controllers.ManagerController{}, "POST:TeamCreate")
	adapter.Router("/manager/team/edit", &controllers.ManagerController{}, "POST:TeamEdit")
	adapter.Router("/manager/team/delete", &controllers.ManagerController{}, "POST:TeamDelete")

	adapter.Router("/manager/team/member/list/:id", &controllers.ManagerController{}, "*:TeamMemberList")
	adapter.Router("/manager/team/member/add", &controllers.ManagerController{}, "POST:TeamMemberAdd")
	adapter.Router("/manager/team/member/delete", &controllers.ManagerController{}, "POST:TeamMemberDelete")
	adapter.Router("/manager/team/member/change_role", &controllers.ManagerController{}, "POST:TeamChangeMemberRole")
	adapter.Router("/manager/team/member/search", &controllers.ManagerController{}, "*:TeamSearchMember")

	adapter.Router("/manager/team/book/list/:id", &controllers.ManagerController{}, "*:TeamBookList")
	adapter.Router("/manager/team/book/add", &controllers.ManagerController{}, "POST:TeamBookAdd")
	adapter.Router("/manager/team/book/delete", &controllers.ManagerController{}, "POST:TeamBookDelete")
	adapter.Router("/manager/team/book/search", &controllers.ManagerController{}, "*:TeamSearchBook")

	adapter.Router("/manager/itemsets", &controllers.ManagerController{}, "*:Itemsets")
	adapter.Router("/manager/itemsets/edit", &controllers.ManagerController{}, "post:ItemsetsEdit")
	adapter.Router("/manager/itemsets/delete", &controllers.ManagerController{}, "post:ItemsetsDelete")

	adapter.Router("/setting", &controllers.SettingController{}, "*:Index")
	adapter.Router("/setting/password", &controllers.SettingController{}, "*:Password")
	adapter.Router("/setting/upload", &controllers.SettingController{}, "*:Upload")

	adapter.Router("/book", &controllers.BookController{}, "*:Index")
	adapter.Router("/book/:key/dashboard", &controllers.BookController{}, "*:Dashboard")
	adapter.Router("/book/:key/setting", &controllers.BookController{}, "*:Setting")
	adapter.Router("/book/:key/users", &controllers.BookController{}, "*:Users")
	adapter.Router("/book/:key/release", &controllers.BookController{}, "post:Release")
	adapter.Router("/book/:key/sort", &controllers.BookController{}, "post:SaveSort")
	adapter.Router("/book/:key/teams", &controllers.BookController{}, "*:Team")

	adapter.Router("/book/create", &controllers.BookController{}, "*:Create")
	adapter.Router("/book/itemsets/search", &controllers.BookController{}, "*:ItemsetsSearch")

	adapter.Router("/book/users/create", &controllers.BookMemberController{}, "post:AddMember")
	adapter.Router("/book/users/change", &controllers.BookMemberController{}, "post:ChangeRole")
	adapter.Router("/book/users/delete", &controllers.BookMemberController{}, "post:RemoveMember")
	adapter.Router("/book/users/import", &controllers.BookController{}, "post:Import")
	adapter.Router("/book/users/copy", &controllers.BookController{}, "post:Copy")

	adapter.Router("/book/setting/save", &controllers.BookController{}, "post:SaveBook")
	adapter.Router("/book/setting/open", &controllers.BookController{}, "post:PrivatelyOwned")
	adapter.Router("/book/setting/transfer", &controllers.BookController{}, "post:Transfer")
	adapter.Router("/book/setting/upload", &controllers.BookController{}, "post:UploadCover")
	adapter.Router("/book/setting/delete", &controllers.BookController{}, "post:Delete")

	adapter.Router("/book/team/add", &controllers.BookController{}, "POST:TeamAdd")
	adapter.Router("/book/team/delete", &controllers.BookController{}, "POST:TeamDelete")
	adapter.Router("/book/team/search", &controllers.BookController{}, "*:TeamSearch")

	//管理文章的路由
	adapter.Router("/manage/blogs", &controllers.BlogController{}, "*:ManageList")
	adapter.Router("/manage/blogs/setting/?:id", &controllers.BlogController{}, "*:ManageSetting")
	adapter.Router("/manage/blogs/edit/?:id", &controllers.BlogController{}, "*:ManageEdit")
	adapter.Router("/manage/blogs/delete", &controllers.BlogController{}, "post:ManageDelete")
	adapter.Router("/manage/blogs/upload", &controllers.BlogController{}, "post:Upload")
	adapter.Router("/manage/blogs/attach/:id", &controllers.BlogController{}, "post:RemoveAttachment")

	//读文章的路由
	adapter.Router("/blogs", &controllers.BlogController{}, "*:List")
	adapter.Router("/blog-attach/:id:int/:attach_id:int", &controllers.BlogController{}, "get:Download")
	adapter.Router("/blog-:id([0-9]+).html", &controllers.BlogController{}, "*:Index")

	//模板相关接口
	adapter.Router("/api/template/get", &controllers.TemplateController{}, "get:Get")
	adapter.Router("/api/template/list", &controllers.TemplateController{}, "post:List")
	adapter.Router("/api/template/add", &controllers.TemplateController{}, "post:Add")
	adapter.Router("/api/template/remove", &controllers.TemplateController{}, "post:Delete")

	adapter.Router("/api/attach/remove/", &controllers.DocumentController{}, "post:RemoveAttachment")
	adapter.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	adapter.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	adapter.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	adapter.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")
	adapter.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	adapter.Router("/api/:key/compare/:id", &controllers.DocumentController{}, "*:Compare")
	adapter.Router("/api/search/user/:key", &controllers.SearchController{}, "*:User")

	adapter.Router("/history/get", &controllers.DocumentController{}, "get:History")
	adapter.Router("/history/delete", &controllers.DocumentController{}, "*:DeleteHistory")
	adapter.Router("/history/restore", &controllers.DocumentController{}, "*:RestoreHistory")

	adapter.Router("/docs/:key", &controllers.DocumentController{}, "*:Index")
	adapter.Router("/docs/:key/:id", &controllers.DocumentController{}, "*:Read")
	adapter.Router("/docs/:key/search", &controllers.DocumentController{}, "post:Search")
	adapter.Router("/export/:key", &controllers.DocumentController{}, "*:Export")
	adapter.Router("/qrcode/:key.png", &controllers.DocumentController{}, "get:QrCode")

	adapter.Router("/attach_files/:key/:attach_id", &controllers.DocumentController{}, "get:DownloadAttachment")

	adapter.Router("/comment/create", &controllers.CommentController{}, "post:Create")
	adapter.Router("/comment/lists", &controllers.CommentController{}, "get:Lists")
	adapter.Router("/comment/index", &controllers.CommentController{}, "*:Index")

	adapter.Router("/search", &controllers.SearchController{}, "get:Index")

	adapter.Router("/tag/:key", &controllers.LabelController{}, "get:Index")
	adapter.Router("/tags", &controllers.LabelController{}, "get:List")

	adapter.Router("/items", &controllers.ItemsetsController{}, "get:Index")
	adapter.Router("/items/:key", &controllers.ItemsetsController{}, "get:List")

}
