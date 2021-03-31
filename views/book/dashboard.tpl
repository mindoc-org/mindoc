<!DOCTYPE html>
<html lang="zh-CN" xmlns="http://www.w3.org/1999/html">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n $.Lang "blog.summary"}} - Powered by MinDoc</title>

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
                    <li class="active"><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> {{i18n $.Lang "blog.summary"}}</a> </li>
                    {{if eq .Model.RoleId 0 1}}
                        <li><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> {{i18n $.Lang "blog.member"}}</a> </li>
                        <li><a href="{{urlfor "BookController.Team" ":key" .Model.Identify}}" class="item"><i class="fa fa-group" aria-hidden="true"></i> {{i18n $.Lang "blog.team"}}</a> </li>
                        <li><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> {{i18n $.Lang "common.setting"}}</a> </li>
                    {{end}}
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">
                            {{if eq .Model.PrivatelyOwned 0}}
                            <i class="fa fa-unlock" aria-hidden="true" title="{{i18n $.Lang "blog.public_project"}}" data-toggle="tooltip"></i>
                            {{else}}
                            <i class="fa fa-lock" aria-hidden="true" title="{{i18n $.Lang "blog.private_project"}}" data-toggle="tooltip"></i>
                            {{end}}
                            {{.Model.BookName}}
                        </strong>
                        {{if ne .Model.RoleId 3}}
                        <a href="{{urlfor "DocumentController.Edit" ":key" .Model.Identify ":id" ""}}" class="btn btn-default btn-sm pull-right" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> {{i18n $.Lang "blog.edit"}}</a>
                        {{end}}
                        <a href="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" target="_blank"><i class="fa fa-eye"></i> {{i18n $.Lang "blog.read"}}</a>

                        {{if eq .Model.RoleId 0 1 2}}
                        <button class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" id="btnRelease"><i class="fa fa-upload" aria-hidden="true"></i> {{i18n $.Lang "blog.publish"}}</button>
                        {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="dashboard">
                        <div class="pull-left" style="width: 200px;margin-bottom: 15px;">
                            <div class="book-image">
                                <img src="{{cdnimg .Model.Cover}}" onerror="this.src='{{cdnimg "/static/images/book.jpg"}}'" style="border: 1px solid #666;width: 175px;">
                            </div>
                        </div>

                            <div class="list">
                                <span class="title">{{i18n $.Lang "blog.creator"}}：</span>
                                <span class="body">{{.Model.CreateName}}</span>
                            </div>
                            <div class="list">
                                <span class="title">{{i18n $.Lang "blog.doc_amount"}}：</span>
                                <span class="body">{{.Model.DocCount}} {{i18n $.Lang "blog.doc_unit"}}</span>
                            </div>
                            <div class="list">
                                <span class="title">{{i18n $.Lang "blog.create_time"}}：</span>
                                <span class="body"> {{date_format .Model.CreateTime "2006-01-02 15:04:05"}} </span>
                            </div>
                            <div class="list">
                                <span class="title">{{i18n $.Lang "blog.update_time"}}：</span>
                                <span class="body"> {{date_format .Model.CreateTime "2006-01-02 15:04:05"}} </span>
                            </div>
                        <div class="list">
                            <span class="title">{{i18n $.Lang "blog.project_role"}}：</span>
                            <span class="body">{{.Model.RoleName}}</span>
                        </div>
                       <!-- {{/* <div class="list">
                            <span class="title">{{i18n $.Lang "blog.comment_amount"}}：</span>
                            <span class="body">{{.Model.CommentCount}} {{i18n $.Lang "blog.comment_unit"}}</span>
                        </div>*/}}
                        -->
                    <div class="list">
                        <span class="title">{{i18n $.Lang "blog.project_desc"}}：</span>
                        <span class="body">{{.Model.Label}}</span>
                    </div>
                        <div class="summary">{{.Description}} </div>

                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#btnRelease").on("click",function () {
            $.ajax({
                url : "{{urlfor "BookController.Release" ":key" .Model.Identify}}",
                data :{"identify" : "{{.Model.Identify}}" },
                type : "post",
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        layer.msg("{{i18n $.Lang "message.publish_to_queue"}}");
                    }else{
                        layer.msg(res.message);
                    }
                }
            });
        });

    });
</script>
</body>
</html>