<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "blog.edit_title"}} - Powered by MinDoc</title>
    <script type="text/javascript">
        window.baseUrl = "{{.BaseUrl}}";
        window.katex = { js: "{{cdnjs "/static/katex/katex"}}",css: "{{cdncss "/static/katex/katex"}}"};
        window.editormdLib = "{{cdnjs "/static/editor.md/lib/"}}";
        window.editor = null;
        window.editURL = "{{urlfor "BlogController.ManageEdit" "blogId" .Model.BlogId}}";
        window.imageUploadURL = "{{urlfor "BlogController.Upload" "blogId" .Model.BlogId}}";
        window.fileUploadURL = "";
        window.blogId = {{.Model.BlogId}};
        window.blogVersion = {{.Model.Version}};
        window.removeAttachURL = "{{urlfor "BlogController.RemoveAttachment" ":id" .Model.BlogId}}";
        window.highlightStyle = "{{.HighlightStyle}}";
        window.lang = {{i18n $.Lang "common.js_lang"}};
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
</head>
<body>

<div class="m-manual manual-editor">
    <div class="manual-head" id="editormd-tools" style="min-width: 1200px; position:absolute;">
        <div class="editormd-group">
            <a href="{{urlfor "BlogController.ManageList"}}" data-toggle="tooltip" data-title="{{i18n .Lang "doc.backward"}}"><i class="fa fa-chevron-left" aria-hidden="true"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" id="markdown-save" data-toggle="tooltip" data-title="{{i18n .Lang "doc.save"}}" class="disabled save"><i class="fa fa-save" aria-hidden="true" name="save"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.undo"}} (Ctrl-Z)"><i class="fa fa-undo first" name="undo" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.redo"}} (Ctrl-Y)"><i class="fa fa-repeat last" name="redo" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.bold"}}"><i class="fa fa-bold first" name="bold" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.italic"}}"><i class="fa fa-italic item" name="italic" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.strikethrough"}}"><i class="fa fa-strikethrough last" name="del" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h1"}}"><i class="fa editormd-bold first" name="h1" unselectable="on">H1</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h2"}}"><i class="fa editormd-bold item" name="h2" unselectable="on">H2</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h3"}}"><i class="fa editormd-bold item" name="h3" unselectable="on">H3</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h4"}}"><i class="fa editormd-bold item" name="h4" unselectable="on">H4</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h5"}}"><i class="fa editormd-bold item" name="h5" unselectable="on">H5</i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.h6"}}"><i class="fa editormd-bold last" name="h6" unselectable="on">H6</i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.unorder_list"}}"><i class="fa fa-list-ul first" name="list-ul" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.order_list"}}"><i class="fa fa-list-ol item" name="list-ol" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.hline"}}"><i class="fa fa-minus last" name="hr" unselectable="on"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.link"}}"><i class="fa fa-link first" name="link" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.ref_link"}}"><i class="fa fa-anchor item" name="reference-link" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.add_pic"}}"><i class="fa fa-picture-o item" name="image" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.code"}}"><i class="fa fa-code item" name="code" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.code_block"}}" unselectable="on"><i class="fa fa-file-code-o item" name="code-block" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.table"}}"><i class="fa fa-table item" name="table" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.quote"}}"><i class="fa fa-quote-right item" name="quote" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.gfm_task"}}"><i class="fa fa-tasks item" name="tasks" aria-hidden="true"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.attachment"}}"><i class="fa fa-paperclip item" aria-hidden="true" name="attachment"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.template"}}"><i class="fa fa-tachometer last" name="template"></i></a>

        </div>

        <div class="editormd-group pull-right">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.close_preview"}}"><i class="fa fa-eye-slash first" name="watch" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.help"}}"><i class="fa fa-question-circle-o last" aria-hidden="true" name="help"></i></a>
        </div>

        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
        </div>

        <div class="clearfix"></div>
    </div>
    <div class="manual-body">
        <div class="manual-editor-container" id="manualEditorContainer" style="min-width: 920px;left: 0;">
            <div class="manual-editormd">
                <div id="docEditor" class="manual-editormd-active"><textarea style="display: none">{{str2html .Model.BlogContent}}</textarea> </div>
            </div>
            <div class="manual-editor-status">
                <div id="attachInfo" class="item">0 {{i18n .Lang "doc.attachments"}}</div>
            </div>
        </div>

    </div>
</div>

<!-- Modal -->

<div class="modal fade" id="uploadAttachModal" tabindex="-1" role="dialog" aria-labelledby="uploadAttachModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="uploadAttachModalForm" class="form-horizontal">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "doc.upload_attachment"}}</h4>
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
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="button" class="btn btn-primary" id="btnUploadAttachFile" data-dismiss="modal">{{i18n .Lang "common.confirm"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!-- Modal -->


<div class="modal fade" id="documentTemplateModal" tabindex="-1" role="dialog" aria-labelledby="{{i18n .Lang "doc.choose_template_type"}}" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h4 class="modal-title" id="modal-title">{{i18n .Lang "doc.choose_template_type"}}</h4>
            </div>
            <div class="modal-body template-list">
                <div class="container">
                    <div class="section">
                        <a data-type="normal" href="javascript:;"><i class="fa fa-file-o"></i></a>
                        <h3><a data-type="normal" href="javascript:;">{{i18n .Lang "doc.normal_tpl"}}</a></h3>
                        <ul>
                            <li>{{i18n .Lang "doc.tpl_default_type"}}</li>
                            <li>{{i18n .Lang "doc.tpl_plain_text"}}</li>
                        </ul>
                    </div>
                    <div class="section">
                        <a data-type="api" href="javascript:;"><i class="fa fa-file-code-o"></i></a>
                        <h3><a data-type="api" href="javascript:;">{{i18n .Lang "doc.api_tpl"}}</a></h3>
                        <ul>
                            <li>{{i18n .Lang "doc.for_api_doc"}}</li>
                            <li>{{i18n .Lang "doc.code_highlight"}}</li>
                        </ul>
                    </div>
                    <div class="section">
                        <a data-type="code" href="javascript:;"><i class="fa fa-book"></i></a>

                        <h3><a data-type="code" href="javascript:;">{{i18n .Lang "doc.data_dict"}}</a></h3>
                        <ul>
                            <li>{{i18n .Lang "doc.for_data_dict"}}</li>
                            <li>{{i18n .Lang "doc.form_support"}}</li>
                        </ul>
                    </div>
                </div>

            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
            </div>
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
<script src="{{cdnjs "/static/js/blog.js" "version"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        editLangPath = {{cdnjs "/static/editor.md/languages/"}} + lang
        if(lang != 'zh-CN') {
            editormd.loadScript(editLangPath, function(){
                window.editor.lang = editormd.defaults.lang;
                window.editor.recreate()
            });
        }
        window.vueApp.lists = {{.AttachList}};
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
                        server: '{{urlfor "BlogController.Upload"}}',
                        formData : { "blogId" : {{.Model.BlogId}}},
                        pick: "#filePicker",
                        fileVal : "editormd-file-file",
                        fileSingleSizeLimit: {{.UploadFileSize}},
                        compress : false
                    }).on( 'fileQueued', function( file ) {
                        var item = {
                            state : "wait",
                            attachment_id : file.id,
                            file_size : file.size,
                            file_name : file.name,
                            message : "{{i18n .Lang "doc.uploading"}}"
                        };
                        window.vueApp.lists.push(item);

                    }).on("uploadError",function (file,reason) {
                        for(var i in window.vueApp.lists){
                            var item = window.vueApp.lists[i];
                            if(item.attachment_id == file.id){
                                item.state = "error";
                                item.message = "{{i18n .Lang "message.upload_failed"}}:" + reason;
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
                            layer.msg("{{i18n .Lang "message.upload_file_size_limit"}}");
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