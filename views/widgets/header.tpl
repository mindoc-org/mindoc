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
                        <a href="{{urlfor "HomeController.Index" }}" title={{i18n .Lang "common.home"}}>{{i18n .Lang "common.home"}}</a>
                    </li>
                    <li {{if eq .ControllerName "BlogController"}}{{if eq  .ActionName "List" "Index"}}class="active"{{end}}{{end}}>
                        <a href="{{urlfor "BlogController.List" }}" title={{i18n .Lang "common.blog"}}>{{i18n .Lang "common.blog"}}</a>
                    </li>
                    <li {{if eq .ControllerName "ItemsetsController"}}class="active"{{end}}>
                        <a href="{{urlfor "ItemsetsController.Index" }}" title={{i18n .Lang "common.project_space"}}>{{i18n .Lang "common.project_space"}}</a>
                    </li>
                </ul>
                <div class="searchbar pull-left visible-lg-inline-block visible-md-inline-block">
                    <form class="form-inline" action="{{urlfor "SearchController.Index"}}" method="get">
                        <input class="form-control" name="keyword" type="search" style="width: 230px;" placeholder="{{i18n .Lang "message.keyword_placeholder"}}" value="{{.Keyword}}">
                        <button class="search-btn">
                            <i class="fa fa-search"></i>
                        </button>
                    </form>
                </div>
            </nav>
            <div style="display: inline-block;" class="navbar-mobile">
                <a href="{{urlfor "HomeController.Index" }}" title={{i18n .Lang "common.home"}}>{{i18n .Lang "common.home"}}</a>
                <a href="{{urlfor "BlogController.List" }}" title={{i18n .Lang "common.blog"}}>{{i18n .Lang "common.blog"}}</a>
            </div>

            <div class="btn-group dropdown-menu-right pull-right slidebar visible-xs-inline-block visible-sm-inline-block">
                <button class="btn btn-default dropdown-toggle hidden-lg" type="button" data-toggle="dropdown"><i class="fa fa-align-justify"></i></button>
                <ul class="dropdown-menu" role="menu">
                    {{if gt .Member.MemberId 0}}
                            <li>
                                <a href="{{urlfor "SettingController.Index"}}" title={{i18n .Lang "common.person_center"}}><i class="fa fa-user" aria-hidden="true"></i> {{i18n .Lang "common.person_center"}}</a>
                            </li>
                            <li>
                                <a href="{{urlfor "BookController.Index"}}" title={{i18n .Lang "common.my_project"}}><i class="fa fa-book" aria-hidden="true"></i> {{i18n .Lang "common.my_project"}}</a>
                            </li>
                            <li>
                                <a href="{{urlfor "BlogController.ManageList"}}" title={{i18n .Lang "common.my_blog"}}><i class="fa fa-file" aria-hidden="true"></i> {{i18n .Lang "common.my_blog"}}</a>
                            </li>
                            {{if eq .Member.Role 0 }}
                            <li>
                                <a href="{{urlfor "ManagerController.Index"}}" title={{i18n .Lang "common.manage"}}><i class="fa fa-university" aria-hidden="true"></i> {{i18n .Lang "common.manage"}}</a>
                            </li>
                            {{end}}
                            <li>
                                <a href="{{urlfor "AccountController.Logout"}}" title={{i18n .Lang "common.logout"}}><i class="fa fa-sign-out"></i> {{i18n .Lang "common.logout"}}</a>
                            </li>

                    {{else}}
                    <li><a href="{{urlfor "AccountController.Login"}}" title={{i18n .Lang "common.login"}}>{{i18n .Lang "common.login"}}</a></li>
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
                            <a href="{{urlfor "SettingController.Index"}}" title={{i18n .Lang "common.person_center"}}><i class="fa fa-user" aria-hidden="true"></i> {{i18n .Lang "common.person_center"}}</a>
                        </li>
                        <li>
                            <a href="{{urlfor "BookController.Index"}}" title={{i18n .Lang "common.my_project"}}><i class="fa fa-book" aria-hidden="true"></i> {{i18n .Lang "common.my_project"}}</a>
                        </li>
                        <li>
                            <a href="{{urlfor "BlogController.ManageList"}}" title={{i18n .Lang "common.my_blog"}}><i class="fa fa-file" aria-hidden="true"></i> {{i18n .Lang "common.my_blog"}}</a>
                        </li>
                        {{if eq .Member.Role 0  1}}
                        <li>
                            <a href="{{urlfor "ManagerController.Index"}}" title={{i18n .Lang "common.manage"}}><i class="fa fa-university" aria-hidden="true"></i> {{i18n .Lang "common.manage"}}</a>
                        </li>
                        {{end}}
                        <li>
                            <a href="{{urlfor "AccountController.Logout"}}" title={{i18n .Lang "common.logout"}}><i class="fa fa-sign-out"></i> {{i18n .Lang "common.logout"}}</a>
                        </li>
                    </ul>
                </li>
                {{else}}
                <li><a href="{{urlfor "AccountController.Login"}}" title={{i18n .Lang "common.login"}}>{{i18n .Lang "common.login"}}</a></li>
                {{end}}
            </ul>
        </nav>
    </div>
</header>
