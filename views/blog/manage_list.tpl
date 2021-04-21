<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "common.my_blog"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/css/fileinput.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/themes/explorer-fa/theme.css"}}" rel="stylesheet" type="text/css">
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
                    <li {{if eq .ControllerName "BookController"}}class="active"{{end}}><a href="{{urlfor "BookController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> {{i18n .Lang "common.my_project"}}</a> </li>
                    <li {{if eq .ControllerName "BlogController"}}class="active"{{end}}><a href="{{urlfor "BlogController.ManageList"}}" class="item"><i class="fa fa-file" aria-hidden="true"></i> {{i18n .Lang "common.my_blog"}}</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{i18n .Lang "blog.blog_list"}}</strong>
                        &nbsp;
                        <a href="{{urlfor "BlogController.ManageSetting"}}" class="btn btn-success btn-sm pull-right">{{i18n .Lang "blog.add_blog"}}</a>
                    </div>
                </div>
                <div class="box-body" id="blogList">
                    <div class="ui items">
                        {{range $index,$item := .ModelList}}
                            <div class="item blog-item">
                                <div class="content">
                                    <a class="header" href="{{urlfor "BlogController.Index" ":id" $item.BlogId}}" target="_blank">
                                        {{if eq $item.BlogStatus "password"}}
                                        <div class="ui teal label horizontal" data-tooltip="{{i18n $.Lang "blog.encryption"}}">{{i18n $.Lang "blog.encrypt"}}</div>
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
                                                <div class="item"><a href="{{urlfor "BlogController.ManageEdit" ":id" $item.BlogId}}" title="{{i18n $.Lang "blog.edit_blog"}}" target="_blank"><i class="fa fa-edit"></i> {{i18n $.Lang "blog.edit"}}</a></div>
                                                <div class="item"><a class="delete-btn" title="{{i18n $.Lang "blog.delete_blog"}}" data-id="{{$item.BlogId}}"><i class="fa fa-trash"></i> {{i18n $.Lang "blog.delete"}}</a></div>
                                                <div class="item"><a href="{{urlfor "BlogController.ManageSetting" ":id" $item.BlogId}}" title="{{i18n $.Lang "blog.setting_blog"}}" class="setting-btn"><i class="fa fa-gear"></i> {{i18n $.Lang "common.setting"}}</a></div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        {{else}}
                        <div class="text-center">{{i18n .Lang "blog.no_blog"}}</div>
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
                    <h4 class="modal-title">{{i18n .Lang "blog.delete_blog"}}</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">{{i18n .Lang "message.confirm_delete_blog"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n .Lang "message.delete_blog_tips"}}</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" id="btnDeleteBlog" class="btn btn-primary" data-loading-text="{{i18n .Lang "message.process"}}">{{i18n .Lang "common.confirm_delete"}}</button>
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
                layer.msg({{i18n .Lang "message.system_error"}});
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
        $("#deleteBlogForm").ajaxForm({
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
                showError({{i18n .Lang "message.system_error"}},"#form-error-message2");
                $("#btnDeleteBlog").button("reset");
            }
        });
    });
</script>
</body>
</html>
