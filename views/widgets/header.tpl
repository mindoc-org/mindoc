<header class="navbar navbar-static-top navbar-fixed-top manual-header" role="banner">
    <div class="container">
        <div class="navbar-header col-sm-12 col-md-9 col-lg-8">
            <a href="{{.BaseUrl}}/" class="navbar-brand" title="{{.SITE_NAME}}">
                {{if .SITE_TITLE}}
                {{.SITE_TITLE}}
                {{else}}
                {{.SITE_NAME}}
                {{end}}
            </a>
            <nav class="collapse navbar-collapse col-sm-10">
                <ul class="nav navbar-nav">
                    <li {{if eq .ControllerName "HomeController"}}class="active"{{end}}>
                        <a href="{{urlfor "HomeController.Index" }}" title="首页">首页</a>
                    </li>
                    <li {{if eq .ControllerName "BlogController"}}{{if eq  .ActionName "List" "Index"}}class="active"{{end}}{{end}}>
                        <a href="{{urlfor "BlogController.List" }}" title="文章">文章</a>
                    </li>
                    <li {{if eq .ControllerName "ItemsetsController"}}class="active"{{end}}>
                        <a href="{{urlfor "ItemsetsController.Index" }}" title="项目空间">项目空间</a>
                    </li>
                </ul>
                <div class="searchbar pull-left visible-lg-inline-block visible-md-inline-block">
                    <form class="form-inline" action="{{urlfor "SearchController.Index"}}" method="get">
                        <input class="form-control" name="keyword" type="search" style="width: 230px;" placeholder="请输入关键词..." value="{{.Keyword}}">
                        <button class="search-btn">
                            <i class="fa fa-search"></i>
                        </button>
                    </form>
                </div>
            </nav>

            <div class="btn-group dropdown-menu-right pull-right slidebar visible-xs-inline-block visible-sm-inline-block">
                <button class="btn btn-default dropdown-toggle hidden-lg" type="button" data-toggle="dropdown"><i class="fa fa-align-justify"></i></button>
                <ul class="dropdown-menu" role="menu">
                    {{if gt .Member.MemberId 0}}
                            <li>
                                <a href="{{urlfor "SettingController.Index"}}" title="个人中心"><i class="fa fa-user" aria-hidden="true"></i> 个人中心</a>
                            </li>
                            <li>
                                <a href="{{urlfor "BookController.Index"}}" title="我的项目"><i class="fa fa-book" aria-hidden="true"></i> 我的项目</a>
                            </li>
                            <li>
                                <a href="{{urlfor "BlogController.ManageList"}}" title="我的文章"><i class="fa fa-file" aria-hidden="true"></i> 我的文章</a>
                            </li>
                            {{if eq .Member.Role 0 }}
                            <li>
                                <a href="{{urlfor "ManagerController.Index"}}" title="管理后台"><i class="fa fa-university" aria-hidden="true"></i> 管理后台</a>
                            </li>
                            {{end}}
                            <li>
                                <a href="{{urlfor "AccountController.Logout"}}" title="退出登录"><i class="fa fa-sign-out"></i> 退出登录</a>
                            </li>

                    {{else}}
                    <li><a href="{{urlfor "AccountController.Login"}}" title="用户登录">登录</a></li>
                    {{end}}
                </ul>
            </div>

        </div>
        <nav class="navbar-collapse hidden-xs hidden-sm" role="navigation">
            <ul class="nav navbar-nav navbar-right">
                {{if gt .Member.MemberId 0}}
                <li>
                    <div class="img user-info" data-toggle="dropdown">
                        <img src="{{cdnimg .Member.Avatar}}" onerror="this.src='{{cdnimg "/static/images/headimgurl.jpg"}}';" class="img-circle userbar-avatar" alt="{{.Member.Account}}">
                        <div class="userbar-content">
                            <span>{{.Member.Account}}</span>
                            <div>{{.Member.RoleName}}</div>
                        </div>
                        <i class="fa fa-chevron-down" aria-hidden="true"></i>
                    </div>
                    <ul class="dropdown-menu user-info-dropdown" role="menu">
                        <li>
                            <a href="{{urlfor "SettingController.Index"}}" title="个人中心"><i class="fa fa-user" aria-hidden="true"></i> 个人中心</a>
                        </li>
                        <li>
                            <a href="{{urlfor "BookController.Index"}}" title="我的项目"><i class="fa fa-book" aria-hidden="true"></i> 我的项目</a>
                        </li>
                        <li>
                            <a href="{{urlfor "BlogController.ManageList"}}" title="我的文章"><i class="fa fa-file" aria-hidden="true"></i> 我的文章</a>
                        </li>
                        {{if eq .Member.Role 0  1}}
                        <li>
                            <a href="{{urlfor "ManagerController.Index"}}" title="管理后台"><i class="fa fa-university" aria-hidden="true"></i> 管理后台</a>
                        </li>
                        {{end}}
                        <li>
                            <a href="{{urlfor "AccountController.Logout"}}" title="退出登录"><i class="fa fa-sign-out"></i> 退出登录</a>
                        </li>
                    </ul>
                </li>
                {{else}}
                <li><a href="{{urlfor "AccountController.Login"}}" title="用户登录">登录</a></li>
                {{end}}
            </ul>
        </nav>
    </div>
</header>