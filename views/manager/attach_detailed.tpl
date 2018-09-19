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

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
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
        {{template "manager/widgets.tpl" "attach"}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">附件管理</strong>
                    </div>
                </div>
                <div class="box-body">
                <form>
                    <div class="form-group">
                        <label>文件名称</label>
                        <input type="text" value="{{.Model.FileName}}" class="form-control input-readonly" readonly placeholder="文件名称">
                    </div>
                    <div class="form-group">
                        <label>是否存在</label>
                            {{if .Model.IsExist }}
                            <input type="text" value="存在" class="form-control input-readonly" readonly placeholder="项目名称">
                            {{else}}
                            <input type="text" value="已删除" class="form-control input-readonly" readonly placeholder="项目名称">
                            {{end}}
                    </div>
                    <div class="form-group">
                        <label>项目/文章名称</label>
                        <input type="text" value="{{.Model.BookName}}" class="form-control input-readonly" readonly placeholder="项目名称">
                    </div>
                    {{if ne .Model.BookId 0}}
                    <div class="form-group">
                        <label>文档名称</label>
                        <input type="text" value="{{.Model.DocumentName}}" class="form-control input-readonly" readonly placeholder="文档名称">
                    </div>
                    {{end}}
                    <div class="form-group">
                        <label>文件路径</label>
                        <input type="text" value="{{.Model.FilePath}}" class="form-control input-readonly" readonly placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>下载路径</label>
                        <input type="text" value="{{.Model.HttpPath}}" class="form-control input-readonly" readonly placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>文件大小</label>
                        <input type="text" value="{{.Model.FileShortSize}}" class="form-control input-readonly" readonly placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>上传时间</label>
                        <input type="text" value="{{date_format .Model.CreateTime "2006-01-02 15:04:05"}}" class="form-control input-readonly" readonly placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <label>用户账号</label>
                        <input type="text" value="{{ .Model.Account }}" class="form-control input-readonly" readonly placeholder="文件路径">
                    </div>
                    <div class="form-group">
                        <a href="{{urlfor "ManagerController.AttachList" }}" class="btn btn-success btn-sm">返回</a>
                        {{if .Model.IsExist }}
                        <a href="{{.Model.LocalHttpPath}}" class="btn btn-default btn-sm" target="_blank" title="下载到本地">下载</a>
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
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>

</body>
</html>