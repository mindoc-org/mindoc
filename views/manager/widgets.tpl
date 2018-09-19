<div class="page-left">
    <ul class="menu">
        <li{{if eq "index" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
        <li{{if eq "users" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
        <li{{if eq "books" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
    {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
        <li{{if eq "setting" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
        <li{{if eq "config" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.Config" }}" class="item"><i class="fa fa-file" aria-hidden="true"></i> 配置文件</a> </li>
        <li{{if eq "attach" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.AttachList" }}" class="item"><i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
        <li{{if eq "label" .}} class="active"{{end}}><a href="{{urlfor "ManagerController.LabelList" }}" class="item"><i class="fa fa-bookmark" aria-hidden="true"></i> 标签管理</a> </li>
    </ul>
</div>