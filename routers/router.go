package routers

import (
	"github.com/astaxie/beego"
	"github.com/lifei6671/godoc/controllers"
)

func init()  {
	beego.Router("/",&controllers.HomeController{},"*:Index")

	beego.Router("/login", &controllers.AccountController{},"*:Login")
	beego.Router("/logout", &controllers.AccountController{},"*:Logout")
	beego.Router("/register", &controllers.AccountController{},"*:Register")
	beego.Router("/find_password", &controllers.AccountController{},"*:FindPassword")
	beego.Router("/valid_email", &controllers.AccountController{},"post:ValidEmail")
	beego.Router("/captcha", &controllers.AccountController{},"*:Captcha")

	beego.Router("/manager", &controllers.ManagerController{},"*:Index")
	beego.Router("/manager/users", &controllers.ManagerController{},"*:Users")
	beego.Router("/manager/member/create", &controllers.ManagerController{},"post:CreateMember")
	beego.Router("/manager/member/update-member-status",&controllers.ManagerController{},"post:UpdateMemberStatus")
	beego.Router("/manager/member/change-member-role", &controllers.ManagerController{},"post:ChangeMemberRole")
	beego.Router("/manager/books", &controllers.ManagerController{},"*:Books")
	beego.Router("/manager/books/edit/:key", &controllers.ManagerController{},"*:EditBook")
	beego.Router("/manager/comments", &controllers.ManagerController{},"*:Comments")
	beego.Router("/manager/books/token", &controllers.ManagerController{},"post:CreateToken")
	beego.Router("/manager/setting",&controllers.ManagerController{},"*:Setting")

	beego.Router("/setting", &controllers.SettingController{},"*:Index")
	beego.Router("/setting/password", &controllers.SettingController{},"*:Password")
	beego.Router("/setting/upload", &controllers.SettingController{},"*:Upload")

	beego.Router("/book", &controllers.BookController{},"*:Index")
	beego.Router("/book/:key/dashboard", &controllers.BookController{},"*:Dashboard")
	beego.Router("/book/:key/setting", &controllers.BookController{},"*:Setting")
	beego.Router("/book/:key/users", &controllers.BookController{},"*:Users")
	beego.Router("/book/:key/release", &controllers.BookController{},"post:Release")
	beego.Router("/book/:key/sort", &controllers.BookController{},"post:SaveSort")

	beego.Router("/book/create", &controllers.BookController{},"*:Create")
	beego.Router("/book/users/create", &controllers.BookMemberController{},"post:AddMember")
	beego.Router("/book/users/change", &controllers.BookMemberController{},"post:ChangeRole")
	beego.Router("/book/users/delete", &controllers.BookMemberController{},"post:RemoveMember")

	beego.Router("/book/setting/save", &controllers.BookController{},"post:SaveBook")
	beego.Router("/book/setting/open", &controllers.BookController{},"post:PrivatelyOwned")
	beego.Router("/book/setting/transfer", &controllers.BookController{},"post:Transfer")
	beego.Router("/book/setting/upload", &controllers.BookController{},"post:UploadCover")
	beego.Router("/book/setting/token", &controllers.BookController{},"post:CreateToken")
	beego.Router("/book/setting/delete", &controllers.BookController{},"post:Delete")

	beego.Router("/api/:key/edit/?:id", &controllers.DocumentController{},"*:Edit")
	beego.Router("/api/upload",&controllers.DocumentController{},"post:Upload")
	beego.Router("/api/:key/create",&controllers.DocumentController{},"post:Create")
	beego.Router("/api/:key/delete", &controllers.DocumentController{},"post:Delete")
	beego.Router("/api/:key/content/?:id",&controllers.DocumentController{},"*:Content")


	beego.Router("/docs/:key", &controllers.DocumentController{},"*:Index")
	beego.Router("/docs/:key/:id", &controllers.DocumentController{},"*:Read")
	beego.Router("/export/:key", &controllers.DocumentController{},"*:Export")

	beego.Router("/attach_files/:key/:attach_id",&controllers.DocumentController{},"get:DownloadAttachment")

	beego.Router("/comment/create", &controllers.CommentController{},"post:Create")
	beego.Router("/comment/lists", &controllers.CommentController{},"get:Lists")
	beego.Router("/comment/index", &controllers.CommentController{},"*:Index")

	beego.Router("/search",&controllers.SearchController{},"get:Index")
}

