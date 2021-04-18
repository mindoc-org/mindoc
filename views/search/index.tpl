<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "search.title"}} - Powered by MinDoc</title>
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
            <strong class="search-title">{{i18n .Lang "search.search_title" .Keyword}}</strong>
        </div>
        <div class="row">
            <div class="manual-list">
                {{range $index,$item := .Lists}}
                <div class="search-item">
                    <div class="title">
                {{if eq $item.SearchType "document"}}
                    <span class="label mark-doc">{{i18n $.Lang "search.doc"}}</span>
                        <a href="{{urlfor "DocumentController.Read" ":key" $item.BookIdentify ":id" $item.Identify}}" title="{{$item.DocumentName}}" target="_blank">{{str2html $item.DocumentName}}</a>
                 {{else if eq $item.SearchType "book"}}
                    <span class="label mark-book">{{i18n $.Lang "search.prj"}}</span>
                    <a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" title="{{$item.BookName}}" target="_blank"> {{str2html $item.DocumentName}}</a>
                {{else}}
                    <span class="label mark-blog">{{i18n $.Lang "search.blog"}}</span>
                        <a href="{{urlfor "BlogController.Index" ":id" $item.DocumentId}}" title="{{$item.DocumentName}}" target="_blank"> {{str2html $item.DocumentName}}</a>
                {{end}}
                    </div>
                    <div class="description">
                        {{str2html $item.Description}}
                    </div>
                    <div class="source">
                        {{if eq $item.SearchType "document"}}
                        <span class="item">{{i18n $.Lang "search.from_proj"}}：<a href="{{urlfor "DocumentController.Index" ":key" $item.BookIdentify}}" target="_blank">{{$item.BookName}}</a></span>
                        {{else if eq $item.SearchType "book"}}
                            <span class="item">{{i18n $.Lang "search.prj"}}：<a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" target="_blank">{{$item.BookName}}</a></span>
                        {{else}}
                        <span class="item">{{i18n $.Lang "search.from_blog"}}：<a href="{{urlfor "BlogController.Index" ":id" $item.DocumentId}}" target="_blank">{{$item.BookName}}</a></span>
                        {{end}}
                        <span class="item">{{i18n $.Lang "search.author"}}：{{$item.Author}}</span>
                        <span class="item">{{i18n $.Lang "search.update_time"}}：{{date_format  $item.ModifyTime "2006-01-02 15:04:05"}}</span>
                    </div>
                </div>
                {{else}}
                <div class="search-empty">
                    <img src="{{cdnimg "/static/images/search_empty.png"}}" class="empty-image">
					<span class="empty-text">{{i18n .Lang "search.no_result"}}</span>
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