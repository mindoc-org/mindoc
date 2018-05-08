package acl

var Modules = make(map[string]*Module)

//模块
type Module struct {
	Name string
	Description string
	Code string
	Resources map[string]*Resource
}

//资源
type Resource struct {
	Name string				`json:"name"`
	Code string				`json:"code"`
	ControllerName	string 	`json:"controller_name"`
	ActionName string		`json:"action_name"`
	MethodName string		`json:"method_name"`
}


func NewModule() *Module  {
	return &Module{ Resources : make(map[string]*Resource)}
}

func init()  {

	Modules["Common"] = &Module{
		Name : "公共功能",
		Code : "Common",
		Description:"所有用户都有的功能",
		Resources : map[string]*Resource {
			"Common!Account!Login!*" : { Name: "用户登录" , Code:"Account!Login!*", ControllerName:"Account",ActionName:"Login",MethodName:"*"},
			"Common!Account!Register!*" : { Name: "用户注册" , Code:"Account!Register!*", ControllerName:"Account",ActionName:"Register",MethodName:"*"},
			"Common!Account!FindPassword!*" : { Name: "找回密码" , Code:"Account!FindPassword!*", ControllerName:"Account",ActionName:"FindPassword",MethodName:"*"},
			"Common!Account!ValidEmail!*" : { Name: "邮箱修改密码" , Code:"Account!ValidEmail!*", ControllerName:"Account",ActionName:"ValidEmail",MethodName:"*"},
			"Common!Account!Logout!*" : { Name: "退出登录" , Code:"Account!Logout!*", ControllerName:"Account",ActionName:"Logout",MethodName:"*"},
			"Common!Account!Captcha!*" : { Name: "图片验证码" , Code:"Account!Captcha!*", ControllerName:"Account",ActionName:"Captcha",MethodName:"*"},
			"Common!Home!Index!*" : { Name:"站点首页",Code:"Home!Index!*",ControllerName:"Home",ActionName:"Index",MethodName:"*"},
			"Common!Search!Index!*" : { Name:"项目搜索",Code:"Search!Index!*",ControllerName:"Search",ActionName:"Index",MethodName:"*"},
			"Common!Error!Error404!*" : { Name:"404页面", Code:"Error!Index!*", ControllerName:"Error", ActionName:"Error404",MethodName:"*" },
			"Common!Error!Error403!*" : { Name:"403页面", Code:"Error!Index!*", ControllerName:"Error", ActionName:"Error403",MethodName:"*" },
			"Common!Error!Error500!*" : { Name:"500页面", Code:"Error!Error500!*", ControllerName:"Error", ActionName:"Error500",MethodName:"*" },
		},
	}

	Modules["MemberCommon"] = &Module{
		Name : "用户公共功能",
		Code : "MemberCommon",
		Description:"只有登录用户才有的功能",
		Resources : map[string]*Resource {
			"MemberCommon!Book!Index!*" : { Name: "项目列表" , Code:"Book!Index!*", ControllerName:"Book",ActionName:"Index",MethodName:"*"},
			"MemberCommon!Book!Dashboard!*" : { Name: "项目概述" , Code:"Book!Index!*", ControllerName:"Book",ActionName:"Dashboard",MethodName:"*"},
		},
	}

	Modules["Book"] = &Module{
		Name:"项目管理",
		Code:"Book",
		Resources: map[string]*Resource {
			"Book!Book!Setting!*" : { Name: "项目设置查看" , Code:"Book!Setting!*", ControllerName:"Book",ActionName:"Setting",MethodName:"*"},
			"Book!Book!SaveBook!*" : { Name: "项目设置保存" , Code:"Book!SaveBook!*", ControllerName:"Book",ActionName:"SaveBook",MethodName:"*"},
		},
	}

	Modules["Document"] = &Module{
		Name:"文档管理",
		Code:"Book",
		Resources: map[string]*Resource {

		},
	}

	Modules["Label"] = &Module{
		Name:"标签管理",
		Code:"Book",
		Resources: map[string]*Resource {

		},
	}
	Modules["Manager"] = &Module{
		Name:"后台管理",
		Code:"Book",
		Resources: map[string]*Resource {

		},
	}

	for _,resource := range Modules["Common"].Resources {
		AddMemberPermission("anonymous",*resource)
	}
}


















