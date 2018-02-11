package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/controllers"
)

func beegoRouter(rootpath string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	return beego.Router(conf.GetUrlPrefix()+rootpath, c, mappingMethods...)
}

func init() {
	beegoRouter("/", &controllers.HomeController{}, "*:Index")

	beegoRouter("/login", &controllers.AccountController{}, "*:Login")
	beegoRouter("/logout", &controllers.AccountController{}, "*:Logout")
	beegoRouter("/register", &controllers.AccountController{}, "*:Register")
	beegoRouter("/find_password", &controllers.AccountController{}, "*:FindPassword")
	beegoRouter("/valid_email", &controllers.AccountController{}, "post:ValidEmail")
	beegoRouter("/captcha", &controllers.AccountController{}, "*:Captcha")

	beegoRouter("/manager", &controllers.ManagerController{}, "*:Index")
	beegoRouter("/manager/users", &controllers.ManagerController{}, "*:Users")
	beegoRouter("/manager/users/edit/:id", &controllers.ManagerController{}, "*:EditMember")
	beegoRouter("/manager/member/create", &controllers.ManagerController{}, "post:CreateMember")
	beegoRouter("/manager/member/delete", &controllers.ManagerController{}, "post:DeleteMember")
	beegoRouter("/manager/member/update-member-status", &controllers.ManagerController{}, "post:UpdateMemberStatus")
	beegoRouter("/manager/member/change-member-role", &controllers.ManagerController{}, "post:ChangeMemberRole")
	beegoRouter("/manager/books", &controllers.ManagerController{}, "*:Books")
	beegoRouter("/manager/books/edit/:key", &controllers.ManagerController{}, "*:EditBook")
	beegoRouter("/manager/books/delete", &controllers.ManagerController{}, "*:DeleteBook")
	beegoRouter("/manager/comments", &controllers.ManagerController{}, "*:Comments")
	beegoRouter("/manager/books/token", &controllers.ManagerController{}, "post:CreateToken")
	beegoRouter("/manager/setting", &controllers.ManagerController{}, "*:Setting")
	beegoRouter("/manager/books/transfer", &controllers.ManagerController{}, "post:Transfer")
	beegoRouter("/manager/books/open", &controllers.ManagerController{}, "post:PrivatelyOwned")
	beegoRouter("/manager/attach/list", &controllers.ManagerController{}, "*:AttachList")
	beegoRouter("/manager/attach/detailed/:id", &controllers.ManagerController{}, "*:AttachDetailed")
	beegoRouter("/manager/attach/delete", &controllers.ManagerController{}, "post:AttachDelete")

	beegoRouter("/setting", &controllers.SettingController{}, "*:Index")
	beegoRouter("/setting/password", &controllers.SettingController{}, "*:Password")
	beegoRouter("/setting/upload", &controllers.SettingController{}, "*:Upload")

	beegoRouter("/book", &controllers.BookController{}, "*:Index")
	beegoRouter("/book/:key/dashboard", &controllers.BookController{}, "*:Dashboard")
	beegoRouter("/book/:key/setting", &controllers.BookController{}, "*:Setting")
	beegoRouter("/book/:key/users", &controllers.BookController{}, "*:Users")
	beegoRouter("/book/:key/release", &controllers.BookController{}, "post:Release")
	beegoRouter("/book/:key/sort", &controllers.BookController{}, "post:SaveSort")

	beegoRouter("/book/create", &controllers.BookController{}, "*:Create")
	beegoRouter("/book/users/create", &controllers.BookMemberController{}, "post:AddMember")
	beegoRouter("/book/users/change", &controllers.BookMemberController{}, "post:ChangeRole")
	beegoRouter("/book/users/delete", &controllers.BookMemberController{}, "post:RemoveMember")
	beegoRouter("/book/users/import", &controllers.BookController{}, "post:Import")

	beegoRouter("/book/setting/save", &controllers.BookController{}, "post:SaveBook")
	beegoRouter("/book/setting/open", &controllers.BookController{}, "post:PrivatelyOwned")
	beegoRouter("/book/setting/transfer", &controllers.BookController{}, "post:Transfer")
	beegoRouter("/book/setting/upload", &controllers.BookController{}, "post:UploadCover")
	beegoRouter("/book/setting/token", &controllers.BookController{}, "post:CreateToken")
	beegoRouter("/book/setting/delete", &controllers.BookController{}, "post:Delete")

	beegoRouter("/api/attach/remove/", &controllers.DocumentController{}, "post:RemoveAttachment")
	beegoRouter("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	beegoRouter("/api/upload", &controllers.DocumentController{}, "post:Upload")
	beegoRouter("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	beegoRouter("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")
	beegoRouter("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	beegoRouter("/api/:key/compare/:id", &controllers.DocumentController{}, "*:Compare")
	beegoRouter("/api/search/user/:key", &controllers.SearchController{}, "*:User")

	beegoRouter("/history/get", &controllers.DocumentController{}, "get:History")
	beegoRouter("/history/delete", &controllers.DocumentController{}, "*:DeleteHistory")
	beegoRouter("/history/restore", &controllers.DocumentController{}, "*:RestoreHistory")

	beegoRouter("/docs/:key", &controllers.DocumentController{}, "*:Index")
	beegoRouter("/docs/:key/:id", &controllers.DocumentController{}, "*:Read")
	beegoRouter("/docs/:key/search", &controllers.DocumentController{}, "post:Search")
	beegoRouter("/export/:key", &controllers.DocumentController{}, "*:Export")
	beegoRouter("/qrcode/:key.png", &controllers.DocumentController{}, "get:QrCode")

	beegoRouter("/attach_files/:key/:attach_id", &controllers.DocumentController{}, "get:DownloadAttachment")

	beegoRouter("/comment/create", &controllers.CommentController{}, "post:Create")
	beegoRouter("/comment/lists", &controllers.CommentController{}, "get:Lists")
	beegoRouter("/comment/index", &controllers.CommentController{}, "*:Index")

	beegoRouter("/search", &controllers.SearchController{}, "get:Index")

	beegoRouter("/tag/:key", &controllers.LabelController{}, "get:Index")
	beegoRouter("/tags", &controllers.LabelController{}, "get:List")
}
