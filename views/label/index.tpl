<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{.LabelName}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
<div class="manual-reader manual-container manual-search-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="search-head">
            <strong class="search-title">显示标签为"{{.LabelName}}"的项目</strong>
        </div>
        <div class="row">
            <div class="manual-list">
                {{range $index,$item := .Lists}}
                <div class="list-item">
                    <dl class="manual-item-standard">
                        <dt>
                            <a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" title="{{$item.BookName}}-{{$item.CreateName}}" target="_blank">
                                <img src="{{cdnimg $item.Cover}}" class="cover" alt="{{$item.BookName}}-{{$item.CreateName}}">
                            </a>
                        </dt>
                        <dd>
                            <a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" class="name" title="{{$item.BookName}}-{{$item.CreateName}}" target="_blank">{{$item.BookName}}</a>
                        </dd>
                        <dd>
                            <span class="author">
                                <b class="text">作者</b>
                                <b class="text">-</b>
                                <b class="text">{{if eq $item.RealName "" }}{{$item.CreateName}}{{else}}{{$item.RealName}}{{end}}</b>
                            </span>
                        </dd>
                    </dl>
                </div>
                {{else}}
                    <div class="text-center" style="height: 200px;margin: 100px;font-size: 28px;">暂无项目</div>
                {{end}}

                <div class="clearfix"></div>
            </div>
            <nav class="pagination-container">
                {{.PageHtml}}
            </nav>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
</body>
</html>