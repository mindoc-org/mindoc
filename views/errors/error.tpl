<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="author" content="Minho" />
    <link rel="shortcut icon" href="{{cdnimg "/favicon.ico"}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title> {{if eq 200 .ErrorCode}}友情提示{{else if eq 404 .ErrorCode}}页面不存在{{else}}服务器异0常{{end}} - Powered by MinDoc</title>
    <link href="{{cdncss "/static/fonts/lato-100.css"}}" rel="stylesheet" type="text/css">
    <style type="text/css">
        html, body {
            height: 100%;
        }

        body {
            margin: 0;
            padding: 0;
            width: 100%;
            color: #B0BEC5;
            display: table;
            font-weight: 100;
            font-family: 'Lato';
        }

        .container {
            text-align: center;
            display: table-cell;
            vertical-align: middle;
        }

        .content {
            text-align: center;
            display: inline-block;
        }

        .title {
            font-size: 72px;
            margin-bottom: 40px;
        }
    </style>
</head>
<body>
<div class="container">
    <div class="content">
        {{if .ErrorMessage}}
        {{if .ErrorCode}}
        <div class="title">HTTP {{.ErrorCode}} ： {{.ErrorMessage}}</div>
        {{else}}
        <div class="title">HTTP 500 ： {{.ErrorMessage}}</div>
        {{end}}
        {{else}}
        <div class="title">HTTP 500 ： 服务器异常</div>
        {{end}}
    </div>
</div>
</body>
</html>
