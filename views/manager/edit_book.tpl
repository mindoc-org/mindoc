<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑项目 - Powered by MinDoc</title>

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
                    <li><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li class="active"><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 项目管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 项目设置</strong>
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;">删除项目</button>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm">
                            <div class="form-group">
                                <label>标题</label>
                                <input type="text" class="form-control" placeholder="项目名称" value="{{.Model.BookName}}">
                            </div>
                            <div class="form-group">
                                <label>标识</label>
                                <input type="text" class="form-control" value=" {{.BaseUrl}}{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" disabled>
                            </div>
                            <div class="form-group">
                                <label>描述</label>
                                <textarea rows="3" class="form-control" name="description" style="height: 90px">{{.Model.Description}}</textarea>
                                <p class="text">描述信息不超过300个字符</p>
                            </div>
                            <div class="form-group">
                                <label>标签</label>
                                <input type="text" class="form-control" placeholder="项目标签" value="{{.Model.Label}}">
                                <p class="text">最多允许添加10个标签，多个标签请用“;”分割</p>
                            </div>
                            <div class="form-group">
                                <label>开启评论</label>
                                <div class="radio">
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .Model.CommentStatus "open"}}checked{{end}} name="comment_status" value="open">允许所有人评论<span class="text"></span>
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .Model.CommentStatus "closed"}}checked{{end}} name="comment_status" value="closed">关闭评论<span class="text"></span>
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .Model.CommentStatus "group_only"}}checked{{end}} name="comment_status" value="group_only">仅允许参与者评论<span class="text"></span>
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .Model.CommentStatus "registered_only"}}checked{{end}} name="comment_status" value="registered_only">仅允许注册者评论<span class="text"></span>
                                    </label>
                                </div>
                            </div>
                            {{if eq .Model.PrivatelyOwned 1}}
                            <div class="form-group">
                                <label>访问令牌</label>
                                <div class="row">
                                    <div class="col-sm-8">
                                        <input type="text" class="form-control" placeholder="访问令牌" readonly>
                                    </div>
                                    <div class="col-sm-4">
                                        <button class="btn btn-success btn-sm">重写生成</button>
                                        <button class="btn btn-danger btn-sm">删除令牌</button>
                                    </div>
                                </div>
                            </div>
                            {{end}}
                            <div class="form-group">
                                <button type="submit" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                                <span id="form-error-message" class="error-message"></span>
                            </div>
                        </form>
                    </div>
                    <div class="form-right">
                        <label>
                            <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                                <img src="{{.Model.Cover}}" onerror="this.src='/static/images/book.png'" alt="封面" style="max-width: 120px;border: 1px solid #999" id="headimgurl">
                            </a>
                        </label>
                    </div>
                    <div class="clearfix"></div>

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