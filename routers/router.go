package routers

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/controllers"
)

func rt(req *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", req.URL)
	// req.Header.Set("Host", "httpbin.org") // <--- I set it here as well
	defer log.Printf("request complete. url=%s", req.URL)

	return http.DefaultTransport.RoundTrip(req)
}

// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type CorsTransport struct {
	http.RoundTripper
}

func (t *CorsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// refer:
	// - https://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy/31536962#31536962
	// - https://gist.github.com/simon-cj/b4da0b2bca793ec3b8a5abe04c8fca41
	resp, err = t.RoundTripper.RoundTrip(req)
	logs.Debug(resp)
	if err != nil {
		return nil, err
	}
	resp.Header.Del("Access-Control-Request-Method")
	resp.Header.Set("Access-Control-Allow-Origin", "*")
	return resp, nil
}

func init() {
	web.Any("/cors-anywhere", func(ctx *context.Context) {
		u, _ := url.PathUnescape(ctx.Input.Query("url"))
		logs.Error("ReverseProxy: ", u)
		if len(u) > 0 && strings.HasPrefix(u, "http") {
			if strings.TrimRight(conf.BaseUrl, "/") == ctx.Input.Site() {
				ctx.Redirect(302, u)
			} else {
				target, _ := url.Parse(u)
				logs.Debug("target: ", target)

				proxy := &httputil.ReverseProxy{
					Transport: roundTripper(rt),
					Director: func(req *http.Request) {
						req.Header = ctx.Request.Header
						req.URL.Scheme = target.Scheme
						req.URL.Host = target.Host
						req.URL.Path = target.Path
						req.Header.Set("Host", target.Host)
					},
				}

				http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				proxy.ServeHTTP(ctx.ResponseWriter, ctx.Request)
			}
		} else {
			ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
			ctx.Output.Body([]byte("400 Bad Request"))
		}
	})

	web.Router("/", &controllers.HomeController{}, "*:Index")

	web.Router("/login", &controllers.AccountController{}, "*:Login")
	web.Router("/dingtalk_login", &controllers.AccountController{}, "*:DingTalkLogin")
	web.Router("/qrlogin/:app", &controllers.AccountController{}, "*:QRLogin")
	web.Router("/logout", &controllers.AccountController{}, "*:Logout")
	web.Router("/register", &controllers.AccountController{}, "*:Register")
	web.Router("/find_password", &controllers.AccountController{}, "*:FindPassword")
	web.Router("/valid_email", &controllers.AccountController{}, "post:ValidEmail")
	web.Router("/captcha", &controllers.AccountController{}, "*:Captcha")

	web.Router("/manager", &controllers.ManagerController{}, "*:Index")
	web.Router("/manager/users", &controllers.ManagerController{}, "*:Users")
	web.Router("/manager/users/edit/:id", &controllers.ManagerController{}, "*:EditMember")
	web.Router("/manager/member/create", &controllers.ManagerController{}, "post:CreateMember")
	web.Router("/manager/member/delete", &controllers.ManagerController{}, "post:DeleteMember")
	web.Router("/manager/member/update-member-status", &controllers.ManagerController{}, "post:UpdateMemberStatus")
	web.Router("/manager/member/change-member-role", &controllers.ManagerController{}, "post:ChangeMemberRole")
	web.Router("/manager/books", &controllers.ManagerController{}, "*:Books")
	web.Router("/manager/books/edit/:key", &controllers.ManagerController{}, "*:EditBook")
	web.Router("/manager/books/delete", &controllers.ManagerController{}, "*:DeleteBook")

	web.Router("/manager/comments", &controllers.ManagerController{}, "*:Comments")
	web.Router("/manager/setting", &controllers.ManagerController{}, "*:Setting")
	web.Router("/manager/books/token", &controllers.ManagerController{}, "post:CreateToken")
	web.Router("/manager/books/transfer", &controllers.ManagerController{}, "post:Transfer")
	web.Router("/manager/books/open", &controllers.ManagerController{}, "post:PrivatelyOwned")

	web.Router("/manager/attach/list", &controllers.ManagerController{}, "*:AttachList")
	web.Router("/manager/attach/detailed/:id", &controllers.ManagerController{}, "*:AttachDetailed")
	web.Router("/manager/attach/delete", &controllers.ManagerController{}, "post:AttachDelete")
	web.Router("/manager/label/list", &controllers.ManagerController{}, "get:LabelList")
	web.Router("/manager/label/delete/:id", &controllers.ManagerController{}, "post:LabelDelete")

	//web.Router("/manager/config",  &controllers.ManagerController{}, "*:Config")

	web.Router("/manager/team", &controllers.ManagerController{}, "*:Team")
	web.Router("/manager/team/create", &controllers.ManagerController{}, "POST:TeamCreate")
	web.Router("/manager/team/edit", &controllers.ManagerController{}, "POST:TeamEdit")
	web.Router("/manager/team/delete", &controllers.ManagerController{}, "POST:TeamDelete")

	web.Router("/manager/team/member/list/:id", &controllers.ManagerController{}, "*:TeamMemberList")
	web.Router("/manager/team/member/add", &controllers.ManagerController{}, "POST:TeamMemberAdd")
	web.Router("/manager/team/member/delete", &controllers.ManagerController{}, "POST:TeamMemberDelete")
	web.Router("/manager/team/member/change_role", &controllers.ManagerController{}, "POST:TeamChangeMemberRole")
	web.Router("/manager/team/member/search", &controllers.ManagerController{}, "*:TeamSearchMember")

	web.Router("/manager/team/book/list/:id", &controllers.ManagerController{}, "*:TeamBookList")
	web.Router("/manager/team/book/add", &controllers.ManagerController{}, "POST:TeamBookAdd")
	web.Router("/manager/team/book/delete", &controllers.ManagerController{}, "POST:TeamBookDelete")
	web.Router("/manager/team/book/search", &controllers.ManagerController{}, "*:TeamSearchBook")

	web.Router("/manager/itemsets", &controllers.ManagerController{}, "*:Itemsets")
	web.Router("/manager/itemsets/edit", &controllers.ManagerController{}, "post:ItemsetsEdit")
	web.Router("/manager/itemsets/delete", &controllers.ManagerController{}, "post:ItemsetsDelete")

	web.Router("/setting", &controllers.SettingController{}, "*:Index")
	web.Router("/setting/password", &controllers.SettingController{}, "*:Password")
	web.Router("/setting/upload", &controllers.SettingController{}, "*:Upload")

	web.Router("/book", &controllers.BookController{}, "*:Index")
	web.Router("/book/:key/dashboard", &controllers.BookController{}, "*:Dashboard")
	web.Router("/book/:key/setting", &controllers.BookController{}, "*:Setting")
	web.Router("/book/:key/users", &controllers.BookController{}, "*:Users")
	web.Router("/book/:key/release", &controllers.BookController{}, "post:Release")
	web.Router("/book/:key/sort", &controllers.BookController{}, "post:SaveSort")
	web.Router("/book/:key/teams", &controllers.BookController{}, "*:Team")

	web.Router("/book/create", &controllers.BookController{}, "*:Create")
	web.Router("/book/itemsets/search", &controllers.BookController{}, "*:ItemsetsSearch")

	web.Router("/book/users/create", &controllers.BookMemberController{}, "post:AddMember")
	web.Router("/book/users/change", &controllers.BookMemberController{}, "post:ChangeRole")
	web.Router("/book/users/delete", &controllers.BookMemberController{}, "post:RemoveMember")
	web.Router("/book/users/import", &controllers.BookController{}, "post:Import")
	web.Router("/book/users/copy", &controllers.BookController{}, "post:Copy")

	web.Router("/book/setting/save", &controllers.BookController{}, "post:SaveBook")
	web.Router("/book/setting/open", &controllers.BookController{}, "post:PrivatelyOwned")
	web.Router("/book/setting/transfer", &controllers.BookController{}, "post:Transfer")
	web.Router("/book/setting/upload", &controllers.BookController{}, "post:UploadCover")
	web.Router("/book/setting/delete", &controllers.BookController{}, "post:Delete")

	web.Router("/book/team/add", &controllers.BookController{}, "POST:TeamAdd")
	web.Router("/book/team/delete", &controllers.BookController{}, "POST:TeamDelete")
	web.Router("/book/team/search", &controllers.BookController{}, "*:TeamSearch")

	//管理文章的路由
	web.Router("/manage/blogs", &controllers.BlogController{}, "*:ManageList")
	web.Router("/manage/blogs/setting/?:id", &controllers.BlogController{}, "*:ManageSetting")
	web.Router("/manage/blogs/edit/?:id", &controllers.BlogController{}, "*:ManageEdit")
	web.Router("/manage/blogs/delete", &controllers.BlogController{}, "post:ManageDelete")
	web.Router("/manage/blogs/upload", &controllers.BlogController{}, "post:Upload")
	web.Router("/manage/blogs/attach/:id", &controllers.BlogController{}, "post:RemoveAttachment")

	//读文章的路由
	web.Router("/blogs", &controllers.BlogController{}, "*:List")
	web.Router("/blog-attach/:id:int/:attach_id:int", &controllers.BlogController{}, "get:Download")
	web.Router("/blog-:id([0-9]+).html", &controllers.BlogController{}, "*:Index")

	//模板相关接口
	web.Router("/api/template/get", &controllers.TemplateController{}, "get:Get")
	web.Router("/api/template/list", &controllers.TemplateController{}, "post:List")
	web.Router("/api/template/add", &controllers.TemplateController{}, "post:Add")
	web.Router("/api/template/remove", &controllers.TemplateController{}, "post:Delete")

	web.Router("/api/attach/remove/", &controllers.DocumentController{}, "post:RemoveAttachment")
	web.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	web.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	web.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	web.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")
	web.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	web.Router("/api/:key/compare/:id", &controllers.DocumentController{}, "*:Compare")
	web.Router("/api/search/user/:key", &controllers.SearchController{}, "*:User")

	web.Router("/history/get", &controllers.DocumentController{}, "get:History")
	web.Router("/history/delete", &controllers.DocumentController{}, "*:DeleteHistory")
	web.Router("/history/restore", &controllers.DocumentController{}, "*:RestoreHistory")

	web.Router("/docs/:key", &controllers.DocumentController{}, "*:Index")
	web.Router("/docs/:key/:id", &controllers.DocumentController{}, "*:Read")
	web.Router("/docs/:key/search", &controllers.DocumentController{}, "post:Search")
	web.Router("/export/:key", &controllers.DocumentController{}, "*:Export")
	web.Router("/qrcode/:key.png", &controllers.DocumentController{}, "get:QrCode")

	web.Router("/attach_files/:key/:attach_id", &controllers.DocumentController{}, "get:DownloadAttachment")

	web.Router("/comment/create", &controllers.CommentController{}, "post:Create")
	web.Router("/comment/delete", &controllers.CommentController{}, "post:Delete")
	web.Router("/comment/lists", &controllers.CommentController{}, "get:Lists")
	web.Router("/comment/index", &controllers.CommentController{}, "*:Index")

	web.Router("/search", &controllers.SearchController{}, "get:Index")

	web.Router("/tag/:key", &controllers.LabelController{}, "get:Index")
	web.Router("/tags", &controllers.LabelController{}, "get:List")

	web.Router("/items", &controllers.ItemsetsController{}, "get:Index")
	web.Router("/items/:key", &controllers.ItemsetsController{}, "get:List")

}
