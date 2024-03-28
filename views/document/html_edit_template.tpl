<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "doc.edit_doc"}} - Powered by MinDoc</title>
    <style type="text/css">
        .w-e-menu.selected > i {
            color: #44B036 !important;
        }
    </style>
    <script type="text/javascript">
        window.IS_ENABLE_IFRAME = '{{conf "enable_iframe" }}' === 'true';
        window.BASE_URL = '{{urlfor "HomeController.Index" }}';
    </script>
    <script type="text/javascript">
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
        window.baiduMapKey = "{{.BaiDuMapKey}}";
        window.historyURL = "{{urlfor "DocumentController.History"}}";
        window.removeAttachURL = "{{urlfor "DocumentController.RemoveAttachment"}}";
        window.vueApp = null;
        window.lang = {{i18n $.Lang "common.js_lang"}};
    </script>
    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/jstree/3.3.4/themes/default/style.min.css"}}" rel="stylesheet">
    <!-- <link href="{{cdncss "/static/wangEditor/css/wangEditor.min.css"}}" rel="stylesheet"> -->

    <link href="{{cdncss (print "/static/editor.md/lib/highlight/styles/" .HighlightStyle ".css") "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/jstree.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.css"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>

<div class="m-manual manual-editor">
    <!--{{/*<div class="manual-head" id="editormd-tools">
        <div class="editormd-group">
            <a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" data-toggle="tooltip" data-title="{{i18n .Lang "doc.backward"}}"><i class="fa fa-chevron-left" aria-hidden="true"></i></a>
        </div>
        <div class="editormd-group">
            <a href="javascript:;" id="markdown-save" data-toggle="tooltip" data-title="{{i18n .Lang "doc.save"}}" class="disabled save"><i class="fa fa-save" aria-hidden="true" name="save"></i></a>
        </div>


        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.modify_history"}}"><i class="fa fa-history item" name="history" aria-hidden="true"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.sidebar"}}"><i class="fa fa-columns item" aria-hidden="true" name="sidebar"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.help"}}"><i class="fa fa-question-circle-o last" aria-hidden="true" name="help"></i></a>
        </div>

        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="{{i18n .Lang "doc.publish"}}"><i class="fa fa-cloud-upload" name="release" aria-hidden="true"></i></a>
        </div>

        <div class="editormd-group">
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
            <a href="javascript:;" data-toggle="tooltip" data-title=""></a>
        </div>

        <div class="clearfix"></div>
    </div>*/}}-->
    <div class="manual-body">
        <div class="manual-category" id="manualCategory" style="top: 0;">
            <div class="manual-nav">
                <div class="nav-item active"><i class="fa fa-bars" aria-hidden="true"></i> {{i18n .Lang "doc.document"}}</div>
                <div class="nav-plus pull-right" data-toggle="tooltip" data-title="{{i18n .Lang "doc.backward"}}" data-direction="right">
                    <a style="color: #999999;" href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" target="_blank"><i class="fa fa-chevron-left" aria-hidden="true"></i></a>
                </div>
                <div class="nav-plus pull-right" id="btnAddDocument" data-toggle="tooltip" data-title="{{i18n .Lang "doc.create_doc"}}" data-direction="right"><i class="fa fa-plus" aria-hidden="true"></i></div>
                <div class="clearfix"></div>
            </div>
            <div class="manual-tree" id="sidebar">

            </div>
        </div>
        <div class="manual-editor-container" id="manualEditorContainer" style="top: 0;">
            <div class="manual-wangEditor">
                <div id="htmlEditor" class="manual-editormd-active" style="height: 100%"></div>
            </div>
            <div class="manual-editor-status">
                <div id="attachInfo" class="item" style="display: inline-block; padding: 0 3em;">0 {{i18n .Lang "doc.attachments"}}</div>
            </div>
        </div>

    </div>
</div>
<!-- Modal -->
<div class="modal fade" id="addDocumentModal" tabindex="-1" style="z-index: 10001 !important;" role="dialog" aria-labelledby="addDocumentModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "DocumentController.Create" ":key" .Model.Identify}}" id="addDocumentForm" class="form-horizontal">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <input type="hidden" name="doc_id" value="0">
            <input type="hidden" name="parent_id" value="0">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "doc.create_doc"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "doc.doc_name"}} <span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="doc_name" id="documentName" placeholder="{{i18n .Lang "doc.doc_name"}}" class="form-control"  maxlength="50">
                            <p style="color: #999;font-size: 12px;">{{i18n .Lang "doc.doc_name_tips"}}</p>
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "doc.doc_id"}}</label>
                        <div class="col-sm-10">
                            <input type="text" name="doc_identify" id="documentIdentify" placeholder="{{i18n .Lang "doc.doc_id"}}" class="form-control" maxlength="50">
                            <p style="color: #999;font-size: 12px;">{{i18n .Lang "doc.doc_id_tips"}}</p>
                        </div>

                    </div>
                </div>
                <div class="modal-footer">
                    <span id="add-error-message" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-primary" id="btnSaveDocument" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "doc.save"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<div class="modal fade" id="uploadAttachModal" tabindex="-1" style="z-index: 10001 !important;" role="dialog" aria-labelledby="uploadAttachModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "DocumentController.Create" ":key" .Model.Identify}}" id="addDocumentForm" class="form-horizontal">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <input type="hidden" name="doc_id" value="0">
            <input type="hidden" name="parent_id" value="0">
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
                                    <span class="text">(${(item.file_size/1024/1024).toFixed(4)}MB)</span>
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

<script src="https://cdn.jsdelivr.net/npm/i18next/i18next.min.js"></script>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jstree/3.3.4/jstree.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/class2browser.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/wangEditor/wangEditor.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/wangEditor-plugins/save-menu.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/wangEditor-plugins/release-menu.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/wangEditor-plugins/attach-menu.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/wangEditor-plugins/history-menu.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
<script src="{{cdnjs "/static/to-markdown/dist/to-markdown.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/editor.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/html-editor.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/custom-elements-builtin-0.6.5.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/x-frame-bypass-1.0.2.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        lang = {{i18n $.Lang "common.js_lang"}};
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
                        fileNumLimit : 1,
                        compress : false
                    }).on("beforeFileQueued",function (file) {
                        uploader.reset();
                    }).on( 'fileQueued', function( file ) {
                       var item = {
                           state : "wait",
                           attachment_id : file.id,
                           file_size : file.size,
                           file_name : file.name,
                           message : "{{i18n .Lang "doc.uploading"}}"
                       };
                       window.vueApp.lists.splice(0,0,item);

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
                                    window.vueApp.lists.splice(index, 1, res.attach);
                                }else{
                                    item.message = res.message;
                                    item.state = "error";
                                }
                                break;
                            }
                        }

                    }).on("beforeFileQueued",function (file) {

                    }).on("uploadComplete",function () {

                    }).on("uploadProgress",function (file, percentage) {
                        var $li = $( '#'+file.id ),
                            $percent = $li.find('.progress .progress-bar');

                        $percent.css( 'width', percentage * 100 + '%' );
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