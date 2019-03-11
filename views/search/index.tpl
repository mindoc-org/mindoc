<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>搜索 - Powered by MinDoc</title>
    <meta name="keywords" content="MinDoc,文档在线管理系统,WIKI,wiki,wiki在线,文档在线管理,接口文档在线管理,接口文档管理,{{.Keyword}}">
    <meta name="description" content="MinDoc文档在线管理系统 {{.site_description}}">
    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
</head>
<body>
<div class="manual-reader manual-container manual-search-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="search-head">
            <strong class="search-title">显示"{{.Keyword}}"的搜索结果</strong>
        </div>
        <div class="row">
            <div class="manual-list">
                {{range $index,$item := .Lists}}
                <div class="search-item">
                    <div class="title">
                {{if eq $item.SearchType "document"}}
                    <span class="label mark-doc">文档</span>
                        <a href="{{urlfor "DocumentController.Read" ":key" $item.BookIdentify ":id" $item.Identify}}" title="{{$item.DocumentName}}" target="_blank">{{str2html $item.DocumentName}}</a>
                 {{else if eq $item.SearchType "book"}}
                    <span class="label mark-book">项目</span>
                    <a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" title="{{$item.BookName}}" target="_blank"> {{str2html $item.DocumentName}}</a>
                {{else}}
                    <span class="label mark-blog">文章</span>
                        <a href="{{urlfor "BlogController.Index" ":id" $item.DocumentId}}" title="{{$item.DocumentName}}" target="_blank"> {{str2html $item.DocumentName}}</a>
                {{end}}
                    </div>
                    <div class="description">
                        {{str2html $item.Description}}
                    </div>
                    <div class="source">
                        {{if eq $item.SearchType "document"}}
                        <span class="item">来自项目：<a href="{{urlfor "DocumentController.Index" ":key" $item.BookIdentify}}" target="_blank">{{$item.BookName}}</a></span>
                        {{else if eq $item.SearchType "book"}}
                            <span class="item">项目：<a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" target="_blank">{{$item.BookName}}</a></span>
                        {{else}}
                        <span class="item">来自文章：<a href="{{urlfor "BlogController.Index" ":id" $item.DocumentId}}" target="_blank">{{$item.BookName}}</a></span>
                        {{end}}
                        <span class="item">作者：{{$item.Author}}</span>
                        <span class="item">更新时间：{{date_format  $item.ModifyTime "2006-01-02 15:04:05"}}</span>
                    </div>
                </div>
                {{else}}
                <div class="search-empty">
                    <img src="{{cdnimg "/static/images/search_empty.png"}}" class="empty-image">
					<span class="empty-text">暂无相关搜索结果</span>
                </div>
                {{end}}
                <nav class="pagination-container">
                    {{.PageHtml}}
                </nav>
                <div class="clearfix"></div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
{{.Scripts}}
</body>
</html>