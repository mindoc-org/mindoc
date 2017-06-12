
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>文档比较 - Powered by MinDoc</title>
    <link rel="shortcut icon" href="/favicon.ico" />
    <link href="/static/fonts/notosans.css" rel='stylesheet' type='text/css' />
    <script type="text/javascript" src="/static/jquery/1.12.4/jquery.min.js"></script>
    <link type='text/css' rel='stylesheet' href='/static/mergely/editor/lib/wicked-ui.css' />
    <script type="text/javascript" src="/static/mergely/editor/lib/wicked-ui.js"></script>

    <link type='text/css' rel='stylesheet' href='/static/mergely/editor/lib/tipsy/tipsy.css' />
    <script type="text/javascript" src="/static/mergely/editor/lib/tipsy/jquery.tipsy.js"></script>
    <script type="text/javascript" src="/static/mergely/editor/lib/farbtastic/farbtastic.js"></script>
    <link type="text/css" rel="stylesheet" href="/static/mergely/editor/lib/farbtastic/farbtastic.css" />
    <script type="text/javascript" src="/static/mergely/lib/codemirror.min.js"></script>
    <script type="text/javascript" src="/static/mergely/lib/mergely.min.js"></script>
    <script type="text/javascript" src="/static/mergely/editor/editor.js"></script>
    <link type="text/css" rel="stylesheet" href="/static/mergely/lib/codemirror.css" />
    <link type="text/css" rel="stylesheet" href="/static/mergely/lib/mergely.css" />
    <link type='text/css' rel='stylesheet' href='/static/mergely/editor/editor.css' />
    <script type="text/javascript" src="/static/mergely/lib/searchcursor.js"></script>
    <script type="text/javascript">
        var key = '';
       // var isSample = key === 'usaindep';
    </script>
</head>
<body style="visibility:hidden">
<!-- toolbar -->
<ul id="toolbar">
    <li id="tb-file-save" data-icon="icon-save" title="保存">保存合并</li>
    <li class="separator"></li>
    <li id="tb-view-change-prev" data-icon="icon-arrow-up" title="上一处差异">上一处差异</li>
    <li id="tb-view-change-next" data-icon="icon-arrow-down" title="下一处差异">下一处差异</li>
    <li class="separator"></li>
    <li id="tb-edit-right-merge-left" data-icon="icon-arrow-left-v" title="合并到左侧">合并到左侧</li>
    <li id="tb-edit-left-merge-right" data-icon="icon-arrow-right-v" title="合并到右侧">合并到右侧</li>
    <li id="tb-view-swap" data-icon="icon-swap" title="左右切换">左右切换</li>
</ul>

<!-- find -->
<div class="find">
    <input type="text" placeholder="请输入关键字"/>
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
<script type="text/javascript" src="/static/layer/layer.js"></script>

</body>
</html>
