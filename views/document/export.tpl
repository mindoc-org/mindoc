<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no"/>
    <title>{{.Model.BookName}} - Powered by MinDoc</title>
    <link href="styles/css/kancloud.css" rel="stylesheet">
    <link href="styles/editor.md/css/editormd.preview.css" rel="stylesheet"/>
    <link href="styles/css/markdown.preview.css" rel="stylesheet"/>
    <link href="styles/css/github.css" rel="stylesheet"/>
    <link href="styles/css/export.css" rel="stylesheet"/>
    <link href="styles/font-awesome/css/font-awesome.css" />
</head>

<body>
    <h1 class="article-title">{{.Lists.DocumentName}}</h1>
    <div class="article-body markdown-body editormd-preview-container"  id="page-content">
    {{str2html .Lists.Release}}
    </div>
</body>
</html>