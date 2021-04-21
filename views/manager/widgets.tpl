<div class="page-left">
    <ul class="menu">
        <li{{if eq "index" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> {{i18n .Lang "mgr.dashboard_menu"}}</a> </li>
        <li{{if eq "users" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-user" aria-hidden="true"></i> {{i18n .Lang "mgr.user_menu"}}</a> </li>
        <li{{if eq "team" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Team" }}" class="item"><i class="fa fa-group" aria-hidden="true"></i> {{i18n .Lang "mgr.team_menu"}}</a> </li>
        <li{{if eq "books" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> {{i18n .Lang "mgr.project_menu"}}</a> </li>
        <li{{if eq "itemsets" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Itemsets" }}" class="item"><i class="fa fa-archive" aria-hidden="true"></i> {{i18n .Lang "mgr.project_space_menu"}}</a> </li>

    {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> {{i18n .Lang "mgr.comment_menu"}}</a> </li>*/}}
        <li{{if eq "setting" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> {{i18n .Lang "mgr.config_menu"}}</a> </li>
        {{/*<li{{if eq "config" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.Config" }}" class="item"><i class="fa fa-file" aria-hidden="true"></i> {{i18n .Lang "mgr.config_file"}}</a> </li>*/}}
        <li{{if eq "attach" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.AttachList" }}" class="item"><i class="fa fa-cloud-upload" aria-hidden="true"></i> {{i18n .Lang "mgr.attachment_menu"}}</a> </li>
        <li{{if eq "label" .Action}} class="active"{{end}}><a href="{{urlfor "ManagerController.LabelList" }}" class="item"><i class="fa fa-bookmark" aria-hidden="true"></i> {{i18n .Lang "mgr.label_menu"}}</a> </li>
    </ul>
</div>