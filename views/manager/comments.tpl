<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户中心 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/font-awesome/css/font-awesome.min.css" rel="stylesheet">

    <link href="/static/css/main.css" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
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
                        <li><a href="{{urlfor "ManagerController.Books" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                        {{/*<li class="active"><a href="{{urlfor "ManagerController.Comments" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
                        <li><a href="{{urlfor "ManagerController.Setting" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                        <li><a href="{{urlfor "ManagerController.AttachList" }}" class="item">&nbsp;&nbsp;&nbsp;&nbsp;<i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
                    {{end}}
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">评论管理</strong>
                    </div>
                </div>
                <div class="box-body">

                </div>
            </div>
        </div>
    </div>
</div>
<script src="/static/jquery/1.12.4/jquery.min.js"></script>
<script src="/static/bootstrap/js/bootstrap.min.js"></script>
</body>
</html>