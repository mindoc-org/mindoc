<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户中心 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
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
                    <li class="active"><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                    {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
                    <li><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.AttachList" }}" class="item"><i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.LabelList" }}" class="item"><i class="fa fa-bookmark" aria-hidden="true"></i> 标签管理</a> </li>

                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">仪表盘</strong>
                    </div>
                </div>
                <div class="box-body manager">
                    <a href="{{urlfor "ManagerController.Books"}}" class="dashboard-item">
                        <span class="fa fa-book" aria-hidden="true"></span>
                        <span class="fa-class">项目数量</span>
                        <span class="fa-class">{{.Model.BookNumber}}</span>
                    </a>
                    <div class="dashboard-item">
                        <span class="fa fa-file-text-o" aria-hidden="true"></span>
                        <span class="fa-class">文章数量</span>
                        <span class="fa-class">{{.Model.DocumentNumber}}</span>
                    </div>
                    <a href="{{urlfor "ManagerController.Users"}}" class="dashboard-item">
                            <span class="fa fa-users" aria-hidden="true"></span>
                            <span class="fa-class">会员数量</span>
                            <span class="fa-class">{{.Model.MemberNumber}}</span>
                    </a>
                    <!--
                    {{/*
                    <div class="dashboard-item">
                        <span class="fa fa-comments-o" aria-hidden="true"></span>
                        <span class="fa-class">评论数量</span>
                        <span class="fa-class">{{.Model.CommentNumber}}</span>
                    </div>
                */}}-->
                    <a href="{{urlfor "ManagerController.AttachList" }}" class="dashboard-item">
                        <span class="fa fa-cloud-download" aria-hidden="true"></span>
                        <span class="fa-class">附件数量</span>
                        <span class="fa-class">{{.Model.AttachmentNumber}}</span>
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
</body>
</html>