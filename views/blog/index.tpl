<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="renderer" content="webkit">
    <meta name="author" content="Minho" />
    <meta name="site" content="https://www.iminho.me" />
    <meta name="keywords" content="{{.Model.BlogTitle}}">
    <meta name="description" content="{{.Model.BlogTitle}}-{{.Description}}">
    <title>{{.Model.BlogTitle}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/kancloud.css"}}?_=1531286622" rel="stylesheet">
    <link href="{{cdncss "/static/editor.md/css/editormd.preview.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.preview.css?_=15312866212"}}" rel="stylesheet">
    <link href="{{cdncss "/static/editor.md/lib/highlight/styles/"}}{{if .HighlightStyle}}{{.HighlightStyle}}{{else}}github{{end}}.css" rel="stylesheet">
    <link href="{{cdncss "/static/katex/katex.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/print.css"}}" media="print" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
    <style type="text/css">
        .header{
            min-height: 1rem;
            font-size: 26px;
            font-weight: 400;
            display: block;
            margin: 20px auto;
        }
        .blog-meta{
            display: inline-block;
        }
        .blog-meta>.item{
            display: inline-block;
            color: #666666;
            vertical-align: middle;
        }

        .blog-footer{
            margin: 25px auto;
            /*border-top: 1px solid #E5E5E5;*/
            padding: 20px 1px;
            line-height: 35px;
        }
        .blog-footer span{
            margin-right: 8px;
            padding: 6px 8px;
            font-size: 12px;
            border: 1px solid #e3e3e3;
            color: #4d4d4d
        }
        .blog-footer a:hover{
            color: #ca0c16;
        }
        .footer{
            margin-top: 0;
        }
        .user_img img {
            display: block;
            width: 24px;
            height: 24px;
            border-radius: 50%;
            -o-object-fit: cover;
            object-fit: cover;
            overflow: hidden;
        }
    </style>
</head>
<body>
<div class="manual-reader manual-container manual-search-reader">
{{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="search-head" style="border-bottom-width: 1px;">
            <h1 class="header">
               {{.Model.BlogTitle}}
            </h1>
            <div class="blog-meta">
                <div class="item user_img"><img src="{{cdnimg .Model.MemberAvatar}}" align="{{.Model.CreateName}}"> </div>
                <div class="item">&nbsp;{{.Model.CreateName}}</div>
                <div class="item">发布于</div>
                <div class="item">{{date .Model.Created "Y-m-d H:i:s"}}</div>
                <div class="item">{{.Model.ModifyRealName}}</div>
                <div class="item">修改于</div>
                <div class="item">{{date .Model.Modified "Y-m-d H:i:s"}}</div>
            </div>
        </div>
        <div class="row">
            <div class="article-body  markdown-body editormd-preview-container content">
                {{.Content}}
                {{if .Model.AttachList}}
                <div class="attach-list"><strong>附件</strong><ul>
                {{range $index,$item := .Model.AttachList}}
                <li><a href="{{$item.HttpPath}}" title="{{$item.FileName}}">{{$item.FileName}}</a> </li>
                {{end}}
                </ul>
                {{end}}
            </div>
        </div>
        <div class="row blog-footer">
            <p>
                <span>上一篇</span>
            {{if .Previous}}
                <a href="{{urlfor "BlogController.Index" ":id" .Previous.BlogId}}" title="{{.Previous.BlogTitle}}">{{.Previous.BlogTitle}}
                </a>
            {{else}}
               无
            {{end}}
            </p>
            <p>
                <span>下一篇</span>
            {{if .Next}}
                <a href="{{urlfor "BlogController.Index" ":id" .Next.BlogId}}" title="{{.Next.BlogTitle}}">{{.Next.BlogTitle}}</a>
            {{else}}
                无
            {{end}}
            </p>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
    </div>
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/kancloud.js"}}" type="text/javascript"></script>
</body>
</html>