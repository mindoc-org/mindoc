<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>文章设置 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
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
                        <strong class="box-title"> 文章设置</strong>
                    </div>
                </div>
                <div class="box-body">
                    <form method="post" id="gloablEditForm" action="{{urlfor "BlogController.ManageSetting"}}">
                        <input type="hidden" name="id" id="blogId" value="{{.Model.BlogId}}">
                        <input type="hidden" name="identify" value="{{.Model.BlogIdentify}}">
                        <input type="hidden" name="document_id" value="{{.Model.DocumentId}}">
                        <input type="hidden" name="order_index" value="{{.Model.OrderIndex}}">
                        <div class="form-group">
                            <label>文章标题</label>
                            <input type="text" class="form-control" name="title" id="title" placeholder="文章标题" value="{{.Model.BlogTitle}}">
                        </div>


                        <div class="form-group">
                            <label>文章类型</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .Model.BlogType 0}}checked{{end}} name="blog_type" value="0">普通文章<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .Model.BlogType 1}}checked{{end}} name="blog_type" value="1">链接文章<span class="text"></span>
                                </label>
                            </div>
                        </div>
                        <div class="form-group" id="blogLinkDocument"{{if ne .Model.BlogType 1}} style="display: none;" {{end}}>
                            <label>关联文档</label>
                            <div class="row">
                                <div class="col-sm-6">
                                    <input type="text" class="form-control" placeholder="请输入项目标识" name="bookIdentify" value="{{.Model.BookIdentify}}">
                                </div>
                                <div class="col-sm-6">
                                    <input type="text" class="form-control" placeholder="请输入文档标识" name="documentIdentify" value="{{.Model.DocumentIdentify}}">
                                </div>
                            </div>

                        </div>
                        <div class="form-group">
                            <label>文章状态</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .Model.BlogStatus "public"}}checked{{end}} name="status" value="public">公开<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .Model.BlogStatus "password"}}checked{{end}} name="status" value="password">加密<span class="text"></span>
                                </label>
                            </div>
                        </div>
                        <div class="form-group"{{if eq .Model.BlogStatus "public"}} style="display: none;"{{end}} id="blogPassword">
                            <label>文章密码</label>
                            <input type="password" class="form-control" name="password" id="password" placeholder="文章密码" value="{{.Model.Password}}" maxlength="20">
                        </div>
                        <div class="form-group">
                            <label>文章摘要</label>
                            <textarea rows="3" class="form-control" name="excerpt" style="height: 90px" placeholder="项目描述">{{.Model.BlogExcerpt}}</textarea>
                            <p class="text">文章摘要不超过500个字符</p>
                        </div>

                        <div class="form-group">
                            <button type="submit" id="btnSaveBlogInfo" class="btn btn-success" data-loading-text="保存中...">保存</button>
                            <a href="{{.Referer}}" title="返回" class="btn btn-info">返回</a>
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                    </form>

                    <div class="clearfix"></div>

                </div>
            </div>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
</div>


<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#gloablEditForm").ajaxForm({
            beforeSubmit : function () {
                var title = $.trim($("#title").val());

                if (title === ""){
                    return showError("文章标题不能为空");
                }
                $("#btnSaveBlogInfo").button("loading");
            },success : function ($res) {
                if($res.errcode === 0) {
                    showSuccess("保存成功");
                    $("#blogId").val($res.data.blog_id);
                }else{
                    showError($res.message);
                }
                $("#btnSaveBlogInfo").button("reset");
            }, error : function () {
                showError("服务器异常.");
                $("#btnSaveBlogInfo").button("reset");
            }
        }).find("input[name='status']").change(function () {
           var $status = $(this).val();
           if($status === "password"){
               $("#blogPassword").show();
           }else{
               $("#blogPassword").hide();
           }
        });
        $("#gloablEditForm").find("input[name='blog_type']").change(function () {
            var $link = $(this).val();
            if($link == 1){
                $("#blogLinkDocument").show();
            }else{
                $("#blogLinkDocument").hide();
            }
        });
    });
</script>
</body>
</html>