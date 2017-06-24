<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>附件 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="/static/css/main.css" rel="stylesheet">
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
            {{template "widgets/sidebar.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">附件详情</strong>
                    </div>
                </div>
                <div class="box-body">
                <form method="post" id="saveAttachInfoForm">
                    <div class="form-group">
                        <label>文件名称</label>
                        <input type="text" value="{{.ModelAttach.FileName}}" class="form-control" disabled placeholder="文件名称">
                    </div>
                        <div class="form-group">
                            <label>文件描述</label>
                            <input type="text" class="form-control" name="description" placeholder="文件描述" value="{{.ModelAttach.Description}}">
                        </div>
                    <div class="form-group">
                        <label>是否存在</label>
                            {{if .ModelAttach.IsExist }}
                            <input type="text" value="存在" class="form-control" disabled placeholder="项目名称">
                            {{else}}
                            <input type="text" value="已删除" class="form-control" disabled placeholder="项目名称">
                            {{end}}
                    </div>
                    <div class="form-group">
                        <label>文档名称</label>
                        <input type="text" value="{{.ModelAttach.DocumentName}}" class="form-control" disabled placeholder="文档名称">
                    </div>
                    <div class="form-group">
                        <label>文件路径</label>
                        <input type="text" value="{{.ModelAttach.FilePath}}" class="form-control" disabled placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>下载路径</label>
                        <input type="text" value="{{.ModelAttach.HttpPath}}" class="form-control" disabled placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>文件大小</label>
                        <input type="text" value="{{.ModelAttach.FileShortSize}}" class="form-control" disabled placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>上传时间</label>
                        <input type="text" value="{{date .ModelAttach.CreateTime "Y-m-d H:i:s"}}" class="form-control" disabled placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>用户账号</label>
                        <input type="text" value="{{ .ModelAttach.Account }}" class="form-control" disabled placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <button type="submit" id="btnAttachInfo" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                        <span id="form-error-message" class="error-message"></span>
                    </div>
                    <div class="form-group">
                        {{if .ModelAttach.IsExist }}
                        <a href="{{.ModelAttach.LocalHttpPath}}" class="btn btn-default btn-sm" target="_blank" title="下载到本地">下载</a>
                        {{end}}
                    </div>
                </form>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#saveAttachInfoForm").ajaxForm({
            beforeSubmit : function () {
                var $then = $("#saveAttachInfoForm");

                var description = $.trim($then.find("input[name='description']").val());
                if (description === ""){
                    return showError("请输入文件描述!");
                }
                $("#btnAttachInfo").button("loading");
            },success : function (res) {
                if(res.errcode === 0) {
                    showSuccess("保存成功")
                }else{
                    showError(res.message);
                }
                $("#btnAttachInfo").button("reset");
            }
        });
    });
</script>

</body>
</html>