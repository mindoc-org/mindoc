<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑文档 - Powered by MinDoc</title>
    <script type="text/javascript">
        window.treeCatalog = null;
        window.baseUrl = "{{.BaseUrl}}";
        window.saveing = false;
        window.katex = { js: "{{cdnjs "/static/katex/katex"}}",css: "{{cdncss "/static/katex/katex"}}"};
        window.editormdLib = "{{cdnjs "/static/editor.md/lib/"}}";
        window.editor = null;
        window.imageUploadURL = "{{urlfor "DocumentController.Upload" "identify" .Model.Identify}}";
        window.fileUploadURL = "{{urlfor "DocumentController.Upload" "identify" .Model.Identify}}";
        window.documentCategory = {{.Result}};
        window.book = {{.ModelResult}};
        window.selectNode = null;
        window.deleteURL = "{{urlfor "DocumentController.Delete" ":key" .Model.Identify}}";
        window.editURL = "{{urlfor "DocumentController.Content" ":key" .Model.Identify ":id" ""}}";
        window.releaseURL = "{{urlfor "BookController.Release" ":key" .Model.Identify}}";
        window.sortURL = "{{urlfor "BookController.SaveSort" ":key" .Model.Identify}}";
        window.historyURL = "{{urlfor "DocumentController.History"}}";
        window.removeAttachURL = "{{urlfor "DocumentController.RemoveAttachment"}}";
        window.highlightStyle = "{{.HighlightStyle}}";
        window.template = { "getUrl":"{{urlfor "TemplateController.Get"}}", "listUrl" : "{{urlfor "TemplateController.List"}}", "deleteUrl" : "{{urlfor "TemplateController.Delete"}}", "saveUrl" :"{{urlfor "TemplateController.Add"}}"}
    </script>
    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/jstree/3.3.4/themes/default/style.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/editor.md/css/editormd.css" "version"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/jstree.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.preview.css" "version"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
    <style type="text/css">
        .text{
            font-size: 12px;
            color: #999999;
            font-weight: 200;
        }
    </style>
</head>
<body>

<div class="m-manual manual-editor">
    <div class="manual-head" id="editormd-tools" style="min-width: 1200px; position:absolute;">
        <div class="editormd-group">
            <a href="{{urlfor "BookController.Index"}}" data-toggle="tooltip" data-title="返回"><i class="fa fa-chevron-left" aria-hidden="true"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" id="markdown-save" data-toggle="tooltip" data-title="保存" class="disabled save"><i class="fa fa-save first" aria-hidden="true" name="save"></i></a>
            <a href="javascript:;" id="markdown-template" data-toggle="tooltip" data-title="保存为模板" class="template"><i class="fa fa-briefcase last" aria-hidden="true" name="save-template"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="撤销 (Ctrl-Z)"><i class="fa fa-undo first" name="undo" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="重做 (Ctrl-Y)"><i class="fa fa-repeat last" name="redo" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="粗体"><i class="fa fa-bold first" name="bold" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="斜体"><i class="fa fa-italic item" name="italic" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="删除线"><i class="fa fa-strikethrough last" name="del" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="标题一"><i class="fa editormd-bold first" name="h1" unselectable="on">H1</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="标题二"><i class="fa editormd-bold item" name="h2" unselectable="on">H2</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="标题三"><i class="fa editormd-bold item" name="h3" unselectable="on">H3</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="标题四"><i class="fa editormd-bold item" name="h4" unselectable="on">H4</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="标题五"><i class="fa editormd-bold item" name="h5" unselectable="on">H5</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="标题六"><i class="fa editormd-bold last" name="h6" unselectable="on">H6</i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="无序列表"><i class="fa fa-list-ul first" name="list-ul" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="有序列表"><i class="fa fa-list-ol item" name="list-ol" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="横线"><i class="fa fa-minus last" name="hr" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="链接"><i class="fa fa-link first" name="link" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="引用链接"><i class="fa fa-anchor item" name="reference-link" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="添加图片"><i class="fa fa-picture-o item" name="image" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="行内代码"><i class="fa fa-code item" name="code" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="代码块" unselectable="on"><i class="fa fa-file-code-o item" name="code-block" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="添加表格"><i class="fa fa-table item" name="table" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="引用"><i class="fa fa-quote-right item" name="quote" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="GFM 任务列表"><i class="fa fa-tasks item" name="tasks" aria-hidden="true"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="附件"><i class="fa fa-paperclip item" aria-hidden="true" name="attachment"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="模板"><i class="fa fa-tachometer last" name="template"></i></a>

        </div>

        <div class="editormd-group pull-right">
            <a href="javascript:;" data-toggle="tooltip" data-title="关闭实时预览"><i class="fa fa-eye-slash first" name="watch" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="修改历史"><i class="fa fa-history item" name="history" aria-hidden="true"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="边栏"><i class="fa fa-columns item" aria-hidden="true" name="sidebar"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="使用帮助"><i class="fa fa-question-circle-o last" aria-hidden="true" name="help"></i></a>
        </div>

        <div class="editormd-group pull-right">
            <a href="javascript:;" data-toggle="tooltip" data-title="发布"><i class="fa fa-cloud-upload" name="release" aria-hidden="true"></i></a>
        </div>

        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
        </div>

        <div class="clearfix"></div>
    </div>
    <div class="manual-body">
        <div class="manual-category" id="manualCategory" style="position:absolute;">
            <div class="manual-nav">
                <div class="nav-item active"><i class="fa fa-bars" aria-hidden="true"></i> 文档</div>
                <div class="nav-plus pull-right" id="btnAddDocument" data-toggle="tooltip" data-title="创建文档" data-direction="right"><i class="fa fa-plus" aria-hidden="true"></i></div>
                <div class="clearfix"></div>
            </div>
            <div class="manual-tree" id="sidebar"> </div>
        </div>
        <div class="manual-editor-container" id="manualEditorContainer" style="min-width: 920px;">
            <div class="manual-editormd">
                <div id="docEditor" class="manual-editormd-active"></div>
            </div>
            <div class="manual-editor-status">
                <div id="attachInfo" class="item">0 个附件</div>
            </div>
        </div>

    </div>
</div>
<!-- 创建文档 -->
<div class="modal fade" id="addDocumentModal" tabindex="-1" role="dialog" aria-labelledby="addDocumentModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "DocumentController.Create" ":key" .Model.Identify}}" id="addDocumentForm" class="form-horizontal">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <input type="hidden" name="doc_id" value="0">
            <input type="hidden" name="parent_id" value="0">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">添加文档</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label class="col-sm-2 control-label">文档名称 <span class="error-message">*</span></label>
                    <div class="col-sm-10">
                        <input type="text" name="doc_name" id="documentName" placeholder="文档名称" class="form-control"  maxlength="50">
                    </div>
                </div>
                <div class="form-group">
                    <label class="col-sm-2 control-label">文档标识 <span class="error-message">&nbsp;</span></label>
                    <div class="col-sm-10">
                        <input type="text" name="doc_identify" id="documentIdentify" placeholder="文档唯一标识" class="form-control" maxlength="50">
                        <p style="color: #999;font-size: 12px;">文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头</p>
                    </div>

                </div>
                <div class="form-group">
                        <div class="col-lg-6">
                            <label>
                                <input type="radio" name="is_open" value="1"> 展开<span class="text">(在阅读时会自动展开节点)</span>
                            </label>
                        </div>
                        <div class="col-lg-6">
                            <label>
                                <input type="radio" name="is_open" value="0" checked> 关闭<span class="text">(在阅读时会关闭节点)</span>
                            </label>
                        </div>

                    <div class="clearfix"></div>
                </div>
            </div>
            <div class="modal-footer">
                <span id="add-error-message" class="error-message"></span>
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                <button type="submit" class="btn btn-primary" id="btnSaveDocument" data-loading-text="保存中...">立即保存</button>
            </div>
        </div>
        </form>
    </div>
</div>

<!-- 显示附件 --->
<div class="modal fade" id="uploadAttachModal" tabindex="-1" role="dialog" aria-labelledby="uploadAttachModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="uploadAttachModalForm" class="form-horizontal">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">上传附件</h4>
                </div>
                <div class="modal-body">
                    <div class="attach-drop-panel">
                        <div class="upload-container" id="filePicker"><i class="fa fa-upload" aria-hidden="true"></i></div>
                    </div>
                    <div class="attach-list" id="attachList">
                        <template v-for="item in lists">
                            <div class="attach-item" :id="item.attachment_id">
                                <template v-if="item.state == 'wait'">
                                    <div class="progress">
                                        <div class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="40" aria-valuemin="0" aria-valuemax="100">
                                            <span class="sr-only">0% Complete (success)</span>
                                        </div>
                                    </div>
                                </template>
                                <template v-else-if="item.state == 'error'">
                                    <span class="error-message">${item.message}</span>
                                    <button type="button" class="btn btn-sm close" @click="removeAttach(item.attachment_id)">
                                        <i class="fa fa-remove" aria-hidden="true"></i>
                                    </button>
                                </template>
                                <template v-else>
                                    <a :href="item.http_path" target="_blank" :title="item.file_name">${item.file_name}</a>
                                    <span class="text">(${ formatBytes(item.file_size) })</span>
                                    <span class="error-message">${item.message}</span>
                                    <button type="button" class="btn btn-sm close" @click="removeAttach(item.attachment_id)">
                                        <i class="fa fa-remove" aria-hidden="true"></i>
                                    </button>
                                    <div class="clearfix"></div>
                                </template>
                            </div>
                        </template>
                    </div>
                </div>
                <div class="modal-footer">
                    <span id="add-error-message" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="button" class="btn btn-primary" id="btnUploadAttachFile" data-dismiss="modal">确定</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!-- 显示文档历史 -->
<div class="modal fade" id="documentHistoryModal" tabindex="-1" role="dialog" aria-labelledby="documentHistoryModalModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">文档历史记录</h4>
            </div>
            <div class="modal-body text-center" id="historyList">

            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
            </div>
        </div>
    </div>
</div>
<!--- 选择模板--->
<div class="modal fade" id="documentTemplateModal" tabindex="-1" role="dialog" aria-labelledby="请选择模板类型" aria-hidden="true">
    <div class="modal-dialog" style="width: 780px;">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="modal-title">请选择模板类型</h4>
            </div>
            <div class="modal-body template-list">
                <div class="container">
                    <div class="section">
                        <a data-type="normal" href="javascript:;"><i class="fa fa-file-o"></i></a>
                        <h3><a data-type="normal" href="javascript:;">普通文档</a></h3>
                        <ul>
                            <li>默认类型</li>
                            <li>简单的文本文档</li>
                        </ul>
                    </div>
                    <div class="section">
                        <a data-type="api" href="javascript:;"><i class="fa fa-file-code-o"></i></a>
                        <h3><a data-type="api" href="javascript:;">API文档</a></h3>
                        <ul>
                            <li>用于API文档速写</li>
                            <li>支持代码高亮</li>
                        </ul>
                    </div>
                    <div class="section">
                        <a data-type="code" href="javascript:;"><i class="fa fa-book"></i></a>

                        <h3><a data-type="code" href="javascript:;">数据字典</a></h3>
                        <ul>
                            <li>用于数据字典显示</li>
                            <li>表格支持</li>
                        </ul>
                    </div>
                    <div class="section">
                        <a data-type="customs" href="javascript:;"><i class="fa fa-briefcase"></i></a>

                        <h3><a data-type="customs" href="javascript:;">自定义模板</a></h3>
                        <ul>
                            <li>自定义模板</li>
                            <li>支持任意类型文档</li>
                            <li>可以设置为全局模板</li>
                        </ul>
                    </div>
                </div>

            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
            </div>
        </div>
    </div>
</div>
<!--- 显示自定义模板--->
<div class="modal fade" id="displayCustomsTemplateModal" tabindex="-1" role="dialog" aria-labelledby="displayCustomsTemplateModalLabel">
    <div class="modal-dialog" role="document" style="width: 750px;">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">自定义模板</h4>
            </div>
            <div class="modal-body text-center" id="displayCustomsTemplateList">
                <div class="table-responsive">
                    <table class="table table-hover">
                        <thead>
                        <tr>
                            <td>#</td>
                            <td class="col-sm-3">模板名称</td>
                            <td class="col-sm-2">模板类型</td>
                            <td class="col-sm-2">创建人</td>
                            <td class="col-sm-3">创建时间</td>
                            <td class="col-sm-2">操作</td>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td colspan="7" class="text-center">暂无数据</td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
            </div>
        </div>
    </div>
</div>
<!--- 创建模板--->
<div class="modal fade" id="saveTemplateModal" tabindex="-1" role="dialog" aria-labelledby="saveTemplateModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <form method="post" action="{{urlfor "TemplateController.Add"}}" id="saveTemplateForm" class="form-horizontal">
                <input type="hidden" name="identify" value="{{.Model.Identify}}">
                <input type="hidden" name="content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">保存为模板</h4>
                </div>
                <div class="modal-body text-center">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">模板名称 <span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="template_name" id="templateName" placeholder="模板名称" class="form-control"  maxlength="50">
                        </div>
                    </div>
                    {{if eq .Member.Role 0 1}}
                    <div class="form-group">
                        <div class="col-lg-6">
                            <label>
                                <input type="radio" name="is_global" value="1"> 全局<span class="text">(任何项目都可用)</span>
                            </label>
                        </div>
                        <div class="col-lg-6">
                            <label>
                                <input type="radio" name="is_global" value="0" checked> 项目<span class="text">(只有当前项目可用)</span>
                            </label>
                        </div>
                        <div class="clearfix"></div>
                    </div>
                    {{end}}
                </div>
                <div class="modal-footer">
                    <span class="error-message show-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-primary" id="btnSaveTemplate" data-loading-text="保存中...">立即保存</button>
                </div>
            </form>
        </div>
    </div>
</div>
<template id="template-normal">
{{template "document/template_normal.tpl"}}
</template>
<template id="template-api">
{{template "document/template_api.tpl"}}
</template>
<template id="template-code">
{{template "document/template_code.tpl"}}
</template>
<script src="{{cdnjs "/static/js/array.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jstree/3.3.4/jstree.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/editor.md/editormd.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/array.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/editor.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/markdown.js" "version"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#attachInfo").on("click",function () {
            $("#uploadAttachModal").modal("show");
        });
        window.uploader = null;

        $("#uploadAttachModal").on("shown.bs.modal",function () {
            if(window.uploader === null){
                try {
                    window.uploader = WebUploader.create({
                        auto: true,
                        dnd : true,
                        swf: '{{.BaseUrl}}/static/webuploader/Uploader.swf',
                        server: '{{urlfor "DocumentController.Upload"}}',
                        formData : { "identify" : {{.Model.Identify}},"doc_id" :  window.selectNode.id },
                        pick: "#filePicker",
                        fileVal : "editormd-file-file",
                        compress : false,
                        fileSingleSizeLimit: {{.UploadFileSize}}
                    }).on("beforeFileQueued",function (file) {
                        // uploader.reset();
                        this.options.formData.doc_id = window.selectNode.id;
                    }).on( 'fileQueued', function( file ) {
                        var item = {
                            state : "wait",
                            attachment_id : file.id,
                            file_size : file.size,
                            file_name : file.name,
                            message : "正在上传"
                        };
                        window.vueApp.lists.push(item);

                    }).on("uploadError",function (file,reason) {
                        for(var i in window.vueApp.lists){
                            var item = window.vueApp.lists[i];
                            if(item.attachment_id == file.id){
                                item.state = "error";
                                item.message = "上传失败:" + reason;
                                break;
                            }
                        }

                    }).on("uploadSuccess",function (file, res) {

                        for(var index in window.vueApp.lists){
                            var item = window.vueApp.lists[index];
                            if(item.attachment_id === file.id){
                                if(res.errcode === 0) {
                                    window.vueApp.lists.splice(index, 1, res.attach ? res.attach : res.data);
                                }else{
                                    item.message = res.message;
                                    item.state = "error";
                                }
                            }
                        }
                    }).on("uploadProgress",function (file, percentage) {
                        var $li = $( '#'+file.id ),
                            $percent = $li.find('.progress .progress-bar');

                        $percent.css( 'width', percentage * 100 + '%' );
                    }).on("error", function (type) {
                        if(type === "F_EXCEED_SIZE"){
                            layer.msg("文件超过了限定大小");
                        }
                        console.log(type);
                    });
                }catch(e){
                    console.log(e);
                }
            }
        });
    });
</script>
</body>
</html>