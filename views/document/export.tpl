<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{.Model.BookName}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{.BaseUrl}}/static/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="{{.BaseUrl}}/static/font-awesome/css/font-awesome.min.css" rel="stylesheet">
    <link href="{{.BaseUrl}}/static/jstree/3.3.4/themes/default/style.min.css" rel="stylesheet">

    <link href="{{.BaseUrl}}/static/nprogress/nprogress.css" rel="stylesheet">
    <link href="{{.BaseUrl}}/static/css/kancloud.css" rel="stylesheet">
    <link href="{{.BaseUrl}}/static/css/jstree.css" rel="stylesheet">
    {{if eq .Model.Editor "markdown"}}
    <link href="{{.BaseUrl}}/static/editor.md/css/editormd.preview.css" rel="stylesheet">
    {{else}}
    <link href="{{.BaseUrl}}/static/highlight/styles/zenburn.css" rel="stylesheet">
    {{end}}
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="{{.BaseUrl}}/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="{{.BaseUrl}}/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
<div class="m-manual manual-reader">
    <div class="container-fluid manual-body">

        <div class="manual-article">
            <div class="article-content">
                <h1 id="article-title">{{.Lists.DocumentName}}</h1>
                <div class="article-body  {{if eq $.Model.Editor "markdown"}}markdown-body editormd-preview-container{{else}}editor-content{{end}}"  id="page-content">
                    {{str2html .Lists.Release}}
                </div>
            </div>

    </div>
</div>
</body>
</html>