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

    <link href="/static/css/main.css?_=?_=1531986418" rel="stylesheet">
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
                    <li><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                    <li class="active"><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
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