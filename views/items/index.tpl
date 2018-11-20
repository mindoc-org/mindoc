<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>项目集列表 - Powered by MinDoc</title>
    <meta name="keywords" content="MinDoc,文档在线管理系统,WIKI,wiki,wiki在线,文档在线管理,接口文档在线管理,接口文档管理">
    <meta name="description" content="MinDoc文档在线管理系统 {{.site_description}}">
    <meta name="author" content="Minho" />
    <meta name="site" content="https://www.iminho.me" />
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
            <strong class="search-title">项目集列表</strong>
        </div>
        <div class="row">
            <div class="hide tag-container-outer" style="border: 0;margin-top: 0;padding: 5px 15px;min-height: 200px;">
                <div class="attach-list" id="ItemsetsList">
                    <table class="table">
                        <thead>
                        <tr>
                            <th width="10%">#</th>
                            <th width="30%">项目集名称</th>
                            <th width="20%">项目集标识</th>
                            <th width="10%">项目数量</th>
                            <th width="15%">创建时间</th>
                            <th>操作</th>
                        </tr>
                        </thead>
                        <tbody>
                        {{range $index,$item := .Lists}}
                        <tr>
                            <td>{{$item.ItemId}}</td>
                            <td>{{$item.ItemName}}</td>
                            <td>{{$item.ItemKey}}</td>
                            <td>{{$item.BookNumber}}</td>
                            <td>{{$item.CreateTimeString}}</td>
                            <td>
                                <a href="{{urlfor "ItemsetsController.List" ":key" $item.ItemKey}}" class="btn btn-success btn-sm" target="_blank">详情</a>
                            </td>
                        </tr>
                        {{else}}
                        <tr><td class="text-center" colspan="6">暂无数据</td></tr>
                        {{end}}
                        </tbody>
                    </table>
                </div>
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
{{.Scripts}}
</body>
</html>