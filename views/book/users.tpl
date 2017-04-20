<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>概要 - Powered by MinDoc</title>

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
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" "test"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li class="active"><a href="{{urlfor "BookController.Users" ":key" "test"}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    <li><a href="{{urlfor "BookController.Setting" ":key" "test"}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 成员管理</strong>
                        <a href="{{urlfor "BookController.Edit" ":key" "test" ":id" "1"}}" class="btn btn-success btn-sm pull-right" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 添加成员</a>
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list">
                        <div class="list-item">
                            <img src="/static/images/middle.gif" class="img-circle" width="34" height="34">
                            <span>lifei6671</span>
                            <div class="operate">
                                创始人
                            </div>
                        </div>
                        <div class="list-item">
                            <img src="/static/images/middle.gif" class="img-circle" width="34" height="34">
                            <span>lifei6671</span>
                            <div class="operate">
                                <div class="btn-group">
                                    <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">编辑者 <span class="caret"></span></button>
                                    <ul class="dropdown-menu">
                                        <li><a href="#">管理员</a> </li>
                                        <li><a href="#">编辑者</a> </li>
                                        <li><a href="#">观察者</a> </li>
                                    </ul>
                                </div>
                                <a href="#" class="btn btn-danger btn-sm">移除</a>
                            </div>
                        </div>
                        <div class="list-item">
                            <img src="/static/images/middle.gif" class="img-circle" width="34" height="34">
                            <span>lifei6671</span>
                            <div class="operate">
                                <div class="btn-group">
                                    <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">编辑者 <span class="caret"></span></button>
                                    <ul class="dropdown-menu">
                                        <li><a href="#">管理员</a> </li>
                                        <li><a href="#">编辑者</a> </li>
                                        <li><a href="#">观察者</a> </li>
                                    </ul>
                                </div>
                                <a href="#" class="btn btn-danger btn-sm">移除</a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="/static/jquery/1.12.4/jquery.min.js"></script>
<script src="/static/bootstrap/js/bootstrap.min.js"></script>
<script src="/static/js/main.js" type="text/javascript"></script>

</body>
</html>