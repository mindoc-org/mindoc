<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>我的文章 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/css/fileinput.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/themes/explorer-fa/theme.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
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
                    <li {{if eq .ControllerName "BookController"}}class="active"{{end}}><a href="{{urlfor "BookController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 我的项目</a> </li>
                    <li {{if eq .ControllerName "BlogController"}}class="active"{{end}}><a href="{{urlfor "BlogController.ManageList"}}" class="item"><i class="fa fa-file" aria-hidden="true"></i> 我的文章</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">文章列表</strong>
                        &nbsp;
                        <a href="{{urlfor "BlogController.ManageSetting"}}" class="btn btn-success btn-sm pull-right">添加文章</a>
                    </div>
                </div>
                <div class="box-body" id="blogList">
                    <div class="ui items">
                        {{range $index,$item := .ModelList}}
                            <div class="item blog-item">
                                <div class="content">
                                    <a class="header" href="{{urlfor "BlogController.Index" ":id" $item.BlogId}}" target="_blank">
                                        {{if eq $item.BlogStatus "password"}}
                                        <div class="ui teal label horizontal" data-tooltip="加密">密</div>
                                        {{end}}
                                        {{$item.BlogTitle}}
                                    </a>
                                    <div class="description">
                                        <p class="line-clamp">{{$item.BlogExcerpt}}&nbsp;</p>
                                    </div>
                                    <div class="extra">
                                        <div>
                                            <div class="ui horizontal small list">
                                                <div class="item"><i class="fa fa-clock-o"></i> {{date $item.Modified "Y-m-d H:i:s"}}</div>
                                                <div class="item"><a href="{{urlfor "BlogController.ManageEdit" ":id" $item.BlogId}}" title="文章编辑" target="_blank"><i class="fa fa-edit"></i> 编辑</a></div>
                                                <div class="item"><a class="delete-btn" title="删除文章" data-id="{{$item.BlogId}}"><i class="fa fa-trash"></i> 删除</a></div>
                                                <div class="item"><a href="{{urlfor "BlogController.ManageSetting" ":id" $item.BlogId}}" title="文章设置" class="setting-btn"><i class="fa fa-gear"></i> 设置</a></div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        {{else}}
                        <div class="text-center">暂无文章</div>
                        {{end}}
                    </div>
                    <nav class="pagination-container">
                        {{.PageHtml}}
                    </nav>
                </div>
            </div>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
</div>


<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBlogModal" tabindex="-1" role="dialog" aria-labelledby="deleteBlogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBlogForm" action="{{urlfor "BlogController.ManageDelete"}}">
            <input type="hidden" name="blog_id" value="">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">删除文章</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">确定删除文章吗？</span>
                    <p></p>
                    <p class="text error-message">删除文章后将无法找回。</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" id="btnDeleteBlog" class="btn btn-primary" data-loading-text="删除中...">确定删除</button>
                </div>
            </div>
        </form>
    </div>
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/fileinput.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/locales/zh.js"}}"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">

    /**
     * 删除项目
     */
    function deleteBlog($id) {
        $("#deleteBlogModal").find("input[name='blog_id']").val($id);
        $("#deleteBlogModal").modal("show");
    }
    function copyBook($id){
        var index = layer.load()
        $.ajax({
            url : "{{urlfor "BookController.Copy"}}" ,
            data : {"identify":$id},
            type : "POST",
            dataType : "json",
            success : function ($res) {
                layer.close(index);
                if ($res.errcode === 0) {
                    window.app.lists.splice(0, 0, $res.data);
                    $("#addBookDialogModal").modal("hide");
                } else {
                    layer.msg($res.message);
                }
            },
            error : function () {
                layer.close(index);
                layer.msg('服务器异常');
            }
        });
    }

    $(function () {

        $("#blogList .delete-btn").click(function () {
            deleteBlog($(this).attr("data-id"));
        });
        /**
         * 删除项目
         */
        $("#deleteBookForm").ajaxForm({
            beforeSubmit : function () {
                $("#btnDeleteBlog").button("loading");
            },
            success : function ($res) {
                if($res.errcode === 0){
                    window.location = window.location.href;
                }else{
                    showError(res.message,"#form-error-message2");
                }
                $("#btnDeleteBlog").button("reset");
            },
            error : function () {
                showError("服务器异常","#form-error-message2");
                $("#btnDeleteBlog").button("reset");
            }
        });
    });
</script>
</body>
</html>
