<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>项目管理 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="{{cdncss "/static/css/main.css"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="{{cdnjs "/static/html5shiv/3.7.3/html5shiv.min.js"}}"></script>
    <script src="{{cdnjs "/static/respond.js/1.4.2/respond.min.js" }}"></script>
    <![endif]-->
</head>
<body>
<div class="manual-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                    <li><a href="#" class="item"><i class="fa fa-user" aria-hidden="true"></i> 个人中心</a> </li>
                    <li><a href="{{urlfor "SettingController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-sitemap" aria-hidden="true"></i> 基本信息</a> </li>
                    {{if ne .Member.AuthMethod "ldap"}}
                        <li><a href="{{urlfor "SettingController.Password"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-user" aria-hidden="true"></i> 修改密码</a> </li>
                    {{end}}
                    <li><a href="{{urlfor "BookController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-sitemap" aria-hidden="true"></i> 我的项目</a> </li>
                    {{if eq .Member.Role 0  1}}
                        <li><a href="#" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 管理后台</a> </li>
                        <li><a href="{{urlfor "ManagerController.Index"}}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                        <li><a href="{{urlfor "ManagerController.Users" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                        <li class="active"><a href="{{urlfor "ManagerController.Books" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                        {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
                        <li><a href="{{urlfor "ManagerController.Setting" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                        <li><a href="{{urlfor "ManagerController.AttachList" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
                    {{end}}
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">项目列表</strong>
                    </div>
                </div>
                <div class="box-body" id="bookList">
                    <div class="book-list">

                        {{range $index,$item := .Lists}}
                        <div class="list-item">
                                <div class="book-title">
                                    <div class="pull-left">
                                        <a href="{{urlfor "BookController.Dashboard" ":key" $item.Identify}}" title="编辑项目" data-toggle="tooltip">
                                            {{if eq $item.PrivatelyOwned 0}}
                                            <i class="fa fa-unlock" aria-hidden="true"></i>
                                            {{else}}
                                            <i class="fa fa-lock" aria-hidden="true"></i>
                                            {{end}}
                                            {{$item.BookName}}
                                        </a>
                                    </div>
                                    <div class="pull-right">
                                        <a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" title="查看文档" data-toggle="tooltip" target="_blank"><i class="fa fa-eye"></i> 查看文档</a>
                                        <a href="{{urlfor "DocumentController.Edit" ":key" $item.Identify ":id" ""}}" title="编辑文档" data-toggle="tooltip" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑文档</a>
                                    </div>
                                    <div class="clearfix"></div>
                                </div>
                                <div class="desc-text">
                                    {{if eq $item.Description ""}}
                                    &nbsp;
                                    {{else}}
                                        <a href="{{urlfor "BookController.Dashboard" ":key" $item.Identify}}" title="编辑项目" style="font-size: 12px;" target="_blank">
                                            {{$item.Description}}
                                        </a>
                                    {{end}}
                                </div>
                                <div class="info">
                                <span title="创建时间" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-clock-o"></i>
                                    {{date $item.CreateTime "Y-m-d H:i:s"}}

                                </span>
                                    <span title="创建者" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user"></i> {{$item.CreateName}}</span>
                                    <span title="文档数量" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pie-chart"></i> {{$item.DocCount}}</span>
                                   {{if ne $item.LastModifyText ""}}
                                    <span title="最后编辑" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pencil"></i> 最后编辑: {{$item.LastModifyText}}</span>
                                    {{end}}

                                </div>
                            </div>
                        {{else}}
                        <div class="text-center">暂无数据</div>
                        {{end}}
                    </div>
                    <nav>
                        {{.PageHtml}}
                    </nav>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>

</body>
</html>