<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>附件管理 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
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
                    <li><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                    {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
                    <li><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                    <li class="active"><a href="{{urlfor "ManagerController.AttachList" }}" class="item"><i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.LabelList" }}" class="item"><i class="fa fa-bookmark" aria-hidden="true"></i> 标签管理</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">附件管理</strong>
                    </div>
                </div>
                <div class="box-body">
                    <div class="attach-list" id="attachList">
                        <table class="table">
                            <thead>
                            <tr>
                                <th>#</th>
                                <th>附件名称</th>
                                <th>项目/文章名称</th>
                                <th>文件大小</th>
                                <th>是否存在</th>
                                <th>操作</th>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index,$item := .Lists}}
                            <tr>
                                <td>{{$item.AttachmentId}}</td>
                                <td>{{$item.FileName}}</td>
                                <td>{{$item.BookName}}</td>
                                <td>{{$item.FileShortSize}}</td>
                                <td>{{ if $item.IsExist }} 是{{else}}否{{end}}</td>
                                <td>
                                    <button type="button" data-method="delete" class="btn btn-danger btn-sm" data-id="{{$item.AttachmentId}}" data-loading-text="删除中...">删除</button>
                                    <a href="{{urlfor "ManagerController.AttachDetailed" ":id" $item.AttachmentId}}" class="btn btn-success btn-sm">详情</a>

                                </td>
                            </tr>
                            {{else}}
                            <tr><td class="text-center" colspan="6">暂无数据</td></tr>
                            {{end}}
                            </tbody>
                        </table>
                        <nav class="pagination-container">
                            {{.PageHtml}}
                        </nav>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js" }}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#attachList").on("click","button[data-method='delete']",function () {
            var id = $(this).attr("data-id");
            var $this = $(this);
            $(this).button("loading");
            $.ajax({
                url : "{{urlfor "ManagerController.AttachDelete"}}",
                data : { "attach_id" : id },
                type : "post",
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        $this.closest("tr").remove().empty();
                    }else {
                        layer.msg(res.message);
                    }
                },
                error : function () {
                    layer.msg("服务器异常");
                },
                complete : function () {
                    $this.button("reset");
                }
            });
        });
    });
</script>
</body>
</html>