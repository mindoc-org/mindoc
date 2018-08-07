<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>标签列表 - Powered by MinDoc</title>
    <meta name="author" content="Minho" />
    <meta name="site" content="https://www.iminho.me" />
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
<div class="manual-reader manual-container manual-search-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="search-head">
            <strong class="search-title">显示标签列表</strong>
        </div>
        <div class="row">
            <div class="hide tag-container-outer" style="border: 0;margin-top: 0;padding: 5px 15px;min-height: 200px;">
                <span class="tags">
                    {{range  $index,$item := .Labels}}
                    <a href="{{urlfor "LabelController.Index" ":key" $item.LabelName}}">{{$item.LabelName}}<span class="detail">{{$item.BookNumber}}</span></a>
                    {{else}}
                    <div class="text-center">暂无标签</div>
                    {{end}}
                </span>
            </div>

            <nav class="pagination-container">
                {{if gt .TotalPages 1}}
                {{.PageHtml}}
                {{end}}
                <div class="clearfix"></div>
            </nav>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
</body>
</html>