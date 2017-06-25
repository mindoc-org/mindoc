<div class="page-left">
    <ul class="menu">
        <li><a href="#" class="item"><i class="fa fa-user" aria-hidden="true"></i> 个人中心</a> </li>
        <li class="{{if eq $.SIDEBAR_ID "profile"}}active{{end}}"><a href="{{urlfor "SettingController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-sitemap" aria-hidden="true"></i> 基本信息</a> </li>
        {{if ne .Member.AuthMethod "ldap"}}
            <li class="{{if eq $.SIDEBAR_ID "password"}}active{{end}}"><a href="{{urlfor "SettingController.Password"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-user" aria-hidden="true"></i> 修改密码</a> </li>
        {{end}}
        <li class="{{if eq $.SIDEBAR_ID "mybook"}}active{{end}}"><a href="{{urlfor "BookController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-sitemap" aria-hidden="true"></i> 我的项目</a> </li>
        {{if eq .Member.Role 0  1}}
            <li><a href="#" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 管理后台</a> </li>
            <li class="{{if eq $.SIDEBAR_ID "dashboard"}}active{{end}}"><a href="{{urlfor "ManagerController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
            <li class="{{if eq $.SIDEBAR_ID "users"}}active{{end}}"><a href="{{urlfor "ManagerController.Users" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
            <li class="{{if eq $.SIDEBAR_ID "books"}}active{{end}}"><a href="{{urlfor "ManagerController.Books" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
            {{/*<li class="{{if eq $.SIDEBAR_ID "cpmments"}}active{{end}}"><a href="{{urlfor "ManagerController.Comments" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
            <li class="{{if eq $.SIDEBAR_ID "setting"}}active{{end}}"><a href="{{urlfor "ManagerController.Setting" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
            <li class="{{if eq $.SIDEBAR_ID "attach"}}active{{end}}"><a href="{{urlfor "ManagerController.AttachList" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
        {{end}}
        {{if eq .SIDEBAR_BOOK 1}}
            <li class="active"><a href="###" class="item"><i class="fa fa-book" aria-hidden="true"></i> {{.Model.BookName}}</a> </li>
            <li class="{{if eq $.SIDEBAR_ID "bookdashboard"}}active{{end}}"><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
            {{if eq .Member.Role 0 1 2}}
            {{if eq .Model.RoleId 0 1}}
                {{if eq .Model.LinkId 0}}
                    <li class="{{if eq $.SIDEBAR_ID "bookattach"}}active{{end}}"><a href="{{urlfor "BookController.Attach" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件</a> </li>
                    <li class="{{if eq $.SIDEBAR_ID "booklink"}}active{{end}}"><a href="{{urlfor "BookController.Links" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-book" aria-hidden="true"></i> 链接</a> </li>
                {{end}}
                <li class="{{if eq $.SIDEBAR_ID "bookuser"}}active{{end}}"><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                <li class="{{if eq $.SIDEBAR_ID "booksetting"}}active{{end}}"><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                {{if gt .Model.LinkId 0}}
                    <li class="{{if eq $.SIDEBAR_ID "bookdocument"}}active{{end}}"><a href="{{urlfor "BookController.EditLink" ":key" .Model.Identify}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-book" aria-hidden="true"></i> 文档</a> </li>
                {{end}}
            {{end}}
            {{end}}
        {{end}}
    </ul>
</div>
