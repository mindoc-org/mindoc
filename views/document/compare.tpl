
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>i18n.Tr(c.Lang, "doc.comparision") - Powered by MinDoc</title>
    <link rel="shortcut icon" href="{{cdnimg "/favicon.ico"}}" />
    <link href="{{cdncss "/static/fonts/notosans.css"}}" rel='stylesheet' type='text/css' />
    <link type='text/css' rel='stylesheet' href="{{cdncss "/static/mergely/editor/lib/wicked-ui.css"}}" />
    <link type='text/css' rel='stylesheet' href="{{cdncss "/static/mergely/editor/lib/tipsy/tipsy.css"}}" />
    <link type="text/css" rel="stylesheet" href="{{cdncss "/static/mergely/editor/lib/farbtastic/farbtastic.css"}}" />
    <link type="text/css" rel="stylesheet" href="{{cdncss "/static/mergely/lib/codemirror.css"}}" />
    <link type="text/css" rel="stylesheet" href="{{cdncss "/static/mergely/lib/mergely.css"}}" />
    <link type='text/css' rel='stylesheet' href="{{cdncss "/static/mergely/editor/editor.css"}}" />

    <script type="text/javascript" src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/editor/lib/wicked-ui.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/editor/lib/tipsy/jquery.tipsy.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/editor/lib/farbtastic/farbtastic.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/lib/codemirror.min.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/lib/mergely.min.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/editor/editor.js"}}"></script>
    <script type="text/javascript" src="{{cdnjs "/static/mergely/lib/searchcursor.js"}}"></script>
    <script type="text/javascript">
        var key = '';
       // var isSample = key === 'usaindep';
    </script>
</head>
<body style="visibility:hidden">
<!-- toolbar -->
<ul id="toolbar">
    <li id="tb-file-save" data-icon="icon-save" title="i18n.Tr(c.Lang, "common.save")">i18n.Tr(c.Lang, "doc.save_merge")</li>
    <li class="separator"></li>
    <li id="tb-view-change-prev" data-icon="icon-arrow-up" title="i18n.Tr(c.Lang, "doc.prev_diff")">i18n.Tr(c.Lang, "doc.prev_diff")</li>
    <li id="tb-view-change-next" data-icon="icon-arrow-down" title="i18n.Tr(c.Lang, "doc.next_diff")">i18n.Tr(c.Lang, "doc.next_diff")</li>
    <li class="separator"></li>
    <li id="tb-edit-right-merge-left" data-icon="icon-arrow-left-v" title="i18n.Tr(c.Lang, "doc.merge_to_left")">i18n.Tr(c.Lang, "doc.merge_to_left")</li>
    <li id="tb-edit-left-merge-right" data-icon="icon-arrow-right-v" title="i18n.Tr(c.Lang, "doc.merge_to_right")">i18n.Tr(c.Lang, "doc.merge_to_right")</li>
    <li id="tb-view-swap" data-icon="icon-swap" title="i18n.Tr(c.Lang, "doc.exchange_left_right")">i18n.Tr(c.Lang, "doc.exchange_left_right")</li>
</ul>

<!-- find -->
<div class="find">
    <input type="text" placeholder="i18n.Tr(c.Lang, "message.keyword_placeholder")"/>
    <button class="find-prev"><span class="icon icon-arrow-up"></span></button>
    <button class="find-next"><span class="icon icon-arrow-down"></span></button>
    <button class="find-close"><span class="icon icon-x-mark"></span></button>
</div>
<!-- editor -->
<div style="position: absolute;top: 33px;bottom: 10px;left: 5px;right: 5px;overflow-y: hidden;padding-bottom: 2px;">
    <div id="mergely"></div>
</div>
<template id="historyContent">{{.HistoryContent}}</template>
<template id="documentContent">{{.Content}}</template>
<script type="text/javascript" src="{{cdnjs "/static/layer/layer.js"}}"></script>

</body>
</html>
