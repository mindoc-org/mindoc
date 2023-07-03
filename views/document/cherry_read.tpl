<!DOCTYPE html>
<html lang="zh-CN">
<head>

    <title>{{.Title}} - Powered by MinDoc</title>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="renderer" content="webkit">
    <meta name="author" content="Minho" />
    <meta name="site" content="https://www.iminho.me" />
    <meta name="keywords" content="MinDoc,文档在线管理系统,WIKI,wiki,wiki在线,文档在线管理,接口文档在线管理,接口文档管理,{{.Model.BookName}},{{.Title}}">
    <meta name="description" content="{{.Title}}-{{if .Description}}{{.Description}}{{else}}{{.Model.Description}}{{end}}">

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/jstree/3.3.4/themes/default/style.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/nprogress/nprogress.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/kancloud.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/jstree.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/cherry/cherry-markdown.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/prismjs/prismjs.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/katex/katex.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/print.css" "version"}}" media="print" rel="stylesheet">

    <script type="text/javascript">
        window.IS_ENABLE_IFRAME = '{{conf "enable_iframe" }}' === 'true';
        window.BASE_URL = '{{urlfor "HomeController.Index" }}';
        window.IS_DOCUMENT_INDEX = '{{if .IS_DOCUMENT_INDEX}}true{{end}}' === 'true';
        window.IS_DISPLAY_COMMENT = '{{if .Model.IsDisplayComment}}true{{end}}' === 'true';
    </script>
    <script type="text/javascript">window.book={"identify":"{{.Model.Identify}}"};</script>
    <style>
        .btn-mobile {
            position: absolute;
            right: 10px;
            top: 10px;
        }

        .not-show-comment {
            display: none;
        }

        @media screen and (min-width: 840px) {
            .btn-mobile{
                display: none;
            }
        }
    </style>
</head>
<body>
<div class="m-manual manual-mode-view manual-reader">
    <header class="navbar navbar-static-top manual-head" role="banner">
        <div class="container-fluid">
            <div class="navbar-header pull-left manual-title">
                <span class="slidebar" id="slidebar"><i class="fa fa-align-justify"></i></span>
                <a href="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" title="{{.Model.BookName}}" class="book-title">{{.Model.BookName}}</a>
                <span style="font-size: 12px;font-weight: 100;"></span>
            </div>
            <a href="{{urlfor "HomeController.Index"}}" class="btn btn-default btn-mobile"> <i class="fa fa-home" aria-hidden="true"></i> {{i18n .Lang "common.home"}}</a>
            <div class="navbar-header pull-right manual-menu">
                <div class="dropdown pull-left" style="margin-right: 10px;">
                    <a href="{{urlfor "HomeController.Index"}}" class="btn btn-default"><i class="fa fa-home" aria-hidden="true"></i> {{i18n .Lang "common.home"}}</a>
                </div>
                {{if gt .Member.MemberId 0}}
                {{if eq .Model.RoleId 0 1 2}}
                <div class="dropdown pull-left" style="margin-right: 10px;">
                    <a href="{{urlfor "DocumentController.Edit" ":key" .Model.Identify ":id" ""}}" class="btn btn-danger"><i class="fa fa-edit" aria-hidden="true"></i> {{i18n .Lang "blog.edit"}}</a>
                    {{if eq .Model.RoleId 0 1}}
                    <a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="btn btn-success"><i class="fa fa-user" aria-hidden="true"></i> {{i18n .Lang "blog.member"}}</a>
                    <a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="btn btn-primary"><i class="fa fa-gear" aria-hidden="true"></i> {{i18n .Lang "common.setting"}}</a>
                    {{end}}
                </div>
                {{end}}
                {{end}}
                <a href="javascript:window.print();" id="printSinglePage" class="btn btn-default" style="margin-right: 10px;"><i class="fa fa-print"></i> {{i18n .Lang "doc.print"}}</a>
                <div class="dropdown pull-right" style="margin-right: 10px;">
                {{if eq .Model.PrivatelyOwned 0}}
                {{if .Model.IsEnableShare}}
                    <button type="button" class="btn btn-success" data-toggle="modal" data-target="#shareProject"><i class="fa fa-share-square" aria-hidden="true"></i> {{i18n .Lang "doc.share"}}</button>
                {{end}}
                {{end}}
                </div>
                {{if .Model.IsDownload}}
                <div class="dropdown pull-right" style="margin-right: 10px;">
                    <button type="button" class="btn btn-primary" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        <i class="fa fa-cloud-download" aria-hidden="true"></i> {{i18n .Lang "doc.download"}} <span class="caret"></span>
                    </button>
                    <ul class="dropdown-menu" role="menu" aria-labelledby="dLabel" style="margin-top: -5px;">
                        <li><a href="{{urlfor "DocumentController.Export" ":key" .Model.Identify "output" "pdf"}}" target="_blank">PDF</a> </li>
                        <li><a href="{{urlfor "DocumentController.Export" ":key" .Model.Identify "output" "epub"}}" target="_blank">EPUB</a> </li>
                        <li><a href="{{urlfor "DocumentController.Export" ":key" .Model.Identify "output" "mobi"}}" target="_blank">MOBI</a> </li>
                        <li><a href="{{urlfor "DocumentController.Export" ":key" .Model.Identify "output" "docx"}}" target="_blank">Word</a> </li>
                        {{if eq .Model.Editor "cherry_markdown"}}
                        <li><a href="{{urlfor "DocumentController.Export" ":key" .Model.Identify "output" "markdown"}}" target="_blank">Markdown</a> </li>
                        {{end}}
                    </ul>
                </div>
                {{end}}
            </div>
        </div>
    </header>
    <article class="container-fluid manual-body">
        <div class="manual-left">
            <div class="manual-tab">
                <div class="tab-navg">
                    <span data-mode="view" class="navg-item active"><i class="fa fa-align-justify"></i><b class="text">{{i18n .Lang "doc.contents"}}</b></span>
                    <span data-mode="search" class="navg-item"><i class="fa fa-search"></i><b class="text">{{i18n .Lang "doc.search"}}</b></span>
                    <span id="handlerMenuShow" style="float: right;display: inline-block;padding: 5px;cursor: pointer;">
                        <i class="fa fa-angle-left" style="font-size: 20px;padding-right: 5px;"></i>
                        <span class="pull-right" style="padding-top: 4px;">{{i18n .Lang "doc.expand"}}</span>
                    </span>
                </div>
                <div class="tab-util">
                    <span class="manual-fullscreen-switch">
                        <b class="open fa fa-angle-right" title="{{i18n .Lang "doc.expand"}}"></b>
                        <b class="close fa fa-angle-left" title="{{i18n .Lang "doc.close"}}"></b>
                    </span>
                </div>
                <div class="tab-wrap">
                    <div class="tab-item manual-catalog">
                        <div class="catalog-list read-book-preview" id="sidebar">
                        {{.Result}}
                        </div>

                    </div>
                    <div class="tab-item manual-search">
                        <div class="search-container">
                            <div class="search-form">
                                <form id="searchForm" action="{{urlfor "DocumentController.Search" ":key" .Model.Identify}}" method="post">
                                    <div class="form-group">
                                        <input type="search" placeholder="{{i18n .Lang "message.search_placeholder"}}" class="form-control" name="keyword">
                                        <button type="submit" class="btn btn-default btn-search" id="btnSearch">
                                            <i class="fa fa-search"></i>
                                        </button>
                                    </div>
                                </form>
                            </div>
                            <div class="search-result">
                                <div class="search-empty">
                                    <i class="fa fa-search-plus" aria-hidden="true"></i>
                                    <b class="text">{{i18n .Lang "message.no_search_result"}}</b>
                                </div>
                                <ul class="search-list" id="searchList">
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="m-copyright">
                <p>
                    <div id="view_count">{{i18n .Lang "doc.view_count"}}：{{.ViewCount}}</div>
                    {{i18n $.Lang "doc.doc_publish_by"}} <a href="https://www.iminho.me" target="_blank">MinDoc</a> {{i18n $.Lang "doc.doc_publish"}}
                </p>
            </div>
        </div>
        <div class="manual-right">
            <div id="view_container" class="manual-article {{if eq .Model.Editor "cherry_markdown"}} cherry cherry-markdown {{.MarkdownTheme}} {{end}}">
                <div class="article-content">
                    <div class="article-head {{if eq .Model.Editor "cherry_markdown"}} markdown-article-head {{end}}">
                        <div class="container-fluid">
                            <div class="row">
                                <div class="col-md-2">
                                </div>
                                <div class="col-md-8 text-center {{if eq .Model.Editor "cherry_markdown"}} markdown-title {{else}} editor-content{{end}}">
                                    <h1 id="article-title">{{.Title}}</h1>
                                </div>
                                <div class="col-md-2">
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="article-body {{if eq .Model.Editor "cherry_markdown"}} markdown-article-body {{else}} editor-content{{end}}"  id="page-content">
                        {{.Content}}
                    </div>

                    {{if .Model.IsDisplayComment}}
                    <div id="articleComment" class="m-comment{{if .IS_DOCUMENT_INDEX}} not-show-comment{{end}}">
                        <!-- 评论列表 -->
                        <div class="comment-list" id="commentList">
                            {{range $i, $c := .Page.List}}
                            <div class="comment-item" data-id="{{$c.CommentId}}">
                                <p class="info"><a class="name">{{$c.Author}}</a><span class="date">{{date $c.CommentDate "Y-m-d H:i:s"}}</span></p>
                                <div class="content">{{$c.Content}}</div>
                                <p class="util">
                                    <span class="operate {{if eq $c.ShowDel 1}}toggle{{end}}">
                                        <span class="number">{{$c.Index}}#</span>
                                        {{if eq $c.ShowDel 1}}
                                        <i class="delete e-delete glyphicon glyphicon-remove" style="color:red" onclick="onDelComment({{$c.CommentId}})"></i>
                                        {{end}}
                                    </span>
                                </p>
                            </div>
                            {{end}}
                        </div>

                        <!-- 翻页 -->
                        <ul id="page"></ul>

                        <!-- 发表评论 -->
                        <div class="comment-post">
                            <form class="form" id="commentForm" action="{{urlfor "CommentController.Create"}}" method="post">
                                <label class="enter w-textarea textarea-full">
                                    <textarea class="textarea-input form-control" name="content" id="commentContent" placeholder="文明上网，理性发言" style="height: 72px;"></textarea>
                                    <input type="hidden" name="doc_id" id="doc_id" value="{{.DocumentId}}">
                                </label>
                                <div class="pull-right">
                                    <button class="btn btn-success btn-sm" type="submit" id="btnSubmitComment" data-loading-text="提交中...">提交评论</button>
                                </div>
                            </form>
                        </div>
                    </div>
                    {{end}}

                    <div class="jump-top">
                        <a href="javascript:;" class="view-backtop"><i class="fa fa-arrow-up" aria-hidden="true"></i></a>
                    </div>
                </div>

            </div>
        </div>
        <div class="manual-progress"><b class="progress-bar"></b></div>
    </article>
    <div class="manual-mask"></div>
</div>
{{if .Model.IsEnableShare}}
<!-- 分享项目 -->
<div class="modal fade" id="shareProject" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "doc.share_project"}}</h4>
            </div>
            <div class="modal-body">
                <div class="row">
                    <div class="col-sm-12 text-center" style="padding-bottom: 15px;">
                        <img src="{{urlfor "DocumentController.QrCode" ":key" .Model.Identify}}" alt="扫一扫手机阅读" />
                    </div>
                </div>
                <div class="form-group">
                    <label for="shareUrl" class="col-sm-2 control-label">{{i18n .Lang "doc.share_url"}}</label>
                    <div class="col-sm-10">
                        <input type="text" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" class="form-control" onmouseover="this.select()" id="shareUrl" title="{{i18n .Lang "doc.share_url"}}">
                    </div>
                    <div class="clearfix"></div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "doc.close"}}</button>
            </div>
        </div>
    </div>
</div>
{{end}}
<!-- 下载项目 -->
<div class="modal fade" id="downloadBookModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "doc.share_project"}}</h4>
            </div>
            <div class="modal-body">
                <div class="row">
                    <div class="col-sm-12 text-center" style="padding-bottom: 15px;">
                        <img src="{{urlfor "DocumentController.QrCode" ":key" .Model.Identify}}" alt="扫一扫手机阅读" />
                    </div>
                </div>
                <div class="form-group">
                    <label for="downloadUrl" class="col-sm-2 control-label">{{i18n .Lang "doc.share_url"}}</label>
                    <div class="col-sm-10">
                        <input type="text" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" class="form-control" onmouseover="this.select()" id="downloadUrl" title="{{i18n .Lang "doc.share_url"}}">
                    </div>
                    <div class="clearfix"></div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "doc.close"}}</button>
            </div>
        </div>
    </div>
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap-paginator/bootstrap-paginator.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jstree/3.3.4/jstree.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/nprogress/nprogress.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/cherry/cherry-markdown.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/editor.md/lib/highlight/highlight.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.highlight.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/clipboard.min.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/prismjs/prismjs.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/kancloud.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/splitbar.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/custom-elements-builtin-0.6.5.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/x-frame-bypass-1.0.2.js"}}" type="text/javascript"></script>
<script type="text/javascript">
$(function () {
    $("#searchList").on("click","a",function () {
        var id = $(this).attr("data-id");
        var url = "{{urlfor "DocumentController.Read" ":key" .Model.Identify ":id" ""}}/" + id;
        $(this).parent("li").siblings().find("a").removeClass("active");
        $(this).addClass("active");
        loadDocument(url,id,function (body) {
            return $(body).highlight(window.keyword);
        });
    });

    window.menuControl = true;
    window.foldSetting = '{{.FoldSetting}}';
    if (foldSetting == 'open' || foldSetting == 'first') {
        $('#handlerMenuShow').find('span').text('{{i18n .Lang "doc.fold"}}');
        $('#handlerMenuShow').find('i').attr("class","fa fa-angle-down");
        if (foldSetting == 'open') {
            window.jsTree.jstree().open_all();
        }
        if (foldSetting == 'first') {
            window.jsTree.jstree('close_all');
            var $target = $('.jstree-container-ul').children('li').filter(function(index){
                if($(this).attr('aria-expanded')==false||$(this).attr('aria-expanded')){
                    return $(this);
                }else{
                    delete $(this);
                }
            });
            $target.children('i').trigger('click');
        }
    } else {
        menuControl = false;
        window.jsTree.jstree('close_all');
    }
    $('#handlerMenuShow').on('click', function(){
        if(menuControl){
            $(this).find('span').text('{{i18n .Lang "doc.expand"}}');
            $(this).find('i').attr("class","fa fa-angle-left");
            window.menuControl = false;
            window.jsTree.jstree('close_all');
        }else{
            window.menuControl = true
            $(this).find('span').text('{{i18n .Lang "doc.fold"}}');
            $(this).find('i').attr("class","fa fa-angle-down");
            window.jsTree.jstree().open_all();
        }
    });

    if (!window.IS_DOCUMENT_INDEX && IS_DISPLAY_COMMENT) {
        pageClicked(-1, parseInt($('#doc_id').val()));
    }
});
</script>
{{.Scripts}}
</body>
</html>