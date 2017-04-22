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
                    <li class="active"><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    <li><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">
                            {{if eq .Model.PrivatelyOwned 0}}
                            <i class="fa fa-unlock" aria-hidden="true" title="公开项目" data-toggle="tooltip"></i>
                            {{else}}
                            <i class="fa fa-lock" aria-hidden="true" title="私有项目" data-toggle="tooltip"></i>
                            {{end}}
                            {{.Model.BookName}}
                        </strong>
                        {{if ne .Model.RoleId 3}}
                        <a href="{{urlfor "BookController.Edit" ":key" .Model.Identify ":id" .Model.BookId}}" class="btn btn-default btn-sm pull-right" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑</a>
                        <a href="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" target="_blank"><i class="fa fa-eye"></i> 阅读</a>
                        <a href="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" target="_blank"><i class="fa fa-upload" aria-hidden="true"></i> 发布</a>
                        {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="dashboard">
                        <div class="pull-left" style="width: 200px;margin-bottom: 15px;">
                            <div class="book-image">
                                <img src="{{.Model.Cover}}" onerror="this.src='/static/images/book.jpg'" width="174" height="229" style="border: 1px solid #666">
                            </div>
                        </div>

                            <div class="list">
                                <span class="title">创建者：</span>
                                <span class="body">{{.Model.CreateName}}</span>
                            </div>
                            <div class="list">
                                <span class="title">文档数量：</span>
                                <span class="body">{{.Model.DocCount}} 篇</span>
                            </div>
                            <div class="list">
                                <span class="title">创建时间：</span>
                                <span class="body"> {{date .Model.CreateTime "Y-m-d H:i:s"}} </span>
                            </div>
                            <div class="list">
                                <span class="title">修改时间：</span>
                                <span class="body"> {{date .Model.ModifyTime "Y-m-d H:i:s"}} </span>
                            </div>
                        <div class="list">
                            <span class="title">担任角色：</span>
                            <span class="body">{{.Model.RoleName}}</span>
                        </div>
                        <div class="list">
                            <span class="title">评论数量：</span>
                            <span class="body">{{.Model.CommentCount}} 条</span>
                        </div>
                            <div class="summary">{{.Model.Description}} </div>

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