<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑文档 - Powered by MinDoc</title>
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
        window.historyURL = "{{urlfor "DocumentController.History"}}";
        window.removeAttachURL = "{{urlfor "DocumentController.RemoveAttachment"}}";
        window.highlightStyle = "{{.HighlightStyle}}";
    </script>
    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/jstree/3.3.4/themes/default/style.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/jstree.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/prettify/themes/atelier-estuary-dark.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/markdown.preview.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss (print "/static/editor.md/lib/highlight/styles/" .HighlightStyle ".css") "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/katex/katex.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/quill/quill.core.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/quill/quill.snow.css"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
    <style type="text/css">
        .modal{z-index: 999999999;}
        #docEditor {
            overflow:auto;
            border: 1px solid #ddd;
            border-left: none;
            height: 100%;
            outline:none;
            padding: 5px 5px 30px 5px;
        }
        #docEditor p{
            margin-bottom: 14px;
            line-height: 1.7em;
            font-size: 14px;
            color: #5D5D5D;
        }
        .ql-picker-options{z-index: 99999;}
        .btn-info{background-color: #ffffff !important;}
        .btn-info>i{background-color: #cacbcd !important; color: #393939 !important; box-shadow: inset 0 0 0 1px transparent,inset 0 0 0 0 rgba(34,36,38,.15);}
        .editor-wrapper>pre{padding: 0;}
        .editor-wrapper .editor-code{
            font-size: 13px;
            line-height: 1.8em;
            color: #dcdcdc;
            border-radius: 3px;
            display: block;
            overflow-x: auto;
            padding: 0.5em;
            background: #3f3f3f;
        }
        .editor-wrapper-selected .editor-code{border: 1px solid #1e88e5;}

        .ql-toolbar.ql-snow{
            border: none !important;
        }
        .editor-group{
            float: left;
            height: 32px;
            margin-right: 10px;
        }

        .editor-group .editor-item,.editor-group .editor-item-select>.ql-picker-label{
            float: left;
            display: inline-block;
            min-width: 34px;
            height: 30px !important;
            padding: 5px;
            line-height: 30px;
            text-align: center;
            color: #4b4b4b;
            border-top: 1px solid #ccc !important;
            border-left: 1px solid #ccc !important;
            border-bottom: 1px solid #ccc !important;
            background: #fff;
            border-radius: 0;
            font-size: 12px
        }
        .ql-snow .ql-picker.ql-expanded .ql-picker-options{
            margin-top: 5px;
        }
        .editor-group .editor-item-select>.ql-picker-label{
            border-right: 1px solid #ccc !important;
        }
        .editor-group .editor-item-single-select>.ql-picker-label{
            border-radius: 4px;
            padding: 0;
        }

        .editor-group .editor-item-last{
            border-right: 1px solid #ccc !important;
            border-radius: 0 4px 4px 0;
        }
        .editor-group .editor-item-first{
            border-right: 0;
            border-radius: 4px  0 0 4px;
        }
        .editor-group .disabled:hover{
            background: #ffffff !important;
        }
        .editor-group .editor-item-change:hover{
             background-color: #58CB48 !important;
        }
        .editor-group  .editor-item:hover {
            background-color: #e4e4e4;
            color: #4b4b4b !important;
        }

        .editor-group a{
            float: left;
        }

        .editor-group .change i{
            color: #ffffff;
            background-color: #44B036 !important;
            border: 1px #44B036 solid !important;
        }
        .editor-group .change i:hover{
            background-color: #58CB48 !important;
        }
        .editor-group .disabled i:hover{
            background: #ffffff !important;
        }
        .editor-group a.disabled{
            border-color: #c9c9c9;
            opacity: .6;
            cursor: default
        }
        .editor-group a>i{
            display: inline-block;
            width: 34px !important;
            height: 30px !important;
            line-height: 30px;
            text-align: center;
            color: #4b4b4b;
            border: 1px solid #ccc;
            background: #fff;
            border-radius: 4px;
            font-size: 15px
        }
        .editor-group a>i.item{
            border-radius: 0;
            border-right: 0;
        }
        .editor-group a>i.last{
            border-bottom-left-radius:0;
            border-top-left-radius:0;
        }
        .editor-group a>i.first{
            border-right: 0;
            border-bottom-right-radius:0;
            border-top-right-radius:0;
        }
        .editor-group  a i:hover {
            background-color: #e4e4e4
        }

        .editor-group  a i:after {
            display: block;
            overflow: hidden;
            line-height: 30px;
            text-align: center;
            font-family: icomoon,Helvetica,Arial,sans-serif;
            font-style: normal;
        }
    </style>
</head>
<body>

<div class="m-manual manual-editor">
    <div class="manual-head btn-toolbar" id="editormd-tools"  style="min-width: 1260px;" data-role="editor-toolbar" data-target="#editor">
        <div class="editor-group">
            <a href="{{urlfor "BookController.Index"}}" data-toggle="tooltip" data-title="返回"><i class="fa fa-chevron-left" aria-hidden="true"></i></a>
        </div>
        <div class="editor-group">
            <a href="javascript:;" id="markdown-save" data-toggle="tooltip" data-title="保存" class="disabled save"><i class="fa fa-save first" aria-hidden="true" name="save"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="发布" id="btnRelease"><i class="fa fa-cloud-upload last" name="release" aria-hidden="true"></i></a>
        </div>
        <div class="editor-group">
            <a href="javascript:;" data-toggle="tooltip" data-title="撤销 (Ctrl-Z)" class="ql-undo"><i class="fa fa-undo first" name="undo" unselectable="on"></i></a>
            <a href="javascript:;" data-toggle="tooltip" data-title="重做 (Ctrl-Y)" class="ql-redo"><i class="fa fa-repeat last" name="redo" unselectable="on"></i></a>
        </div>
        <div class="editor-group">
            <select data-toggle="tooltip" data-title="字号" title="字号" class="ql-size editor-item-select editor-item-single-select"></select>
        </div>
        <div class="editor-group">
            <button data-toggle="tooltip" data-title="粗体" class="ql-bold editor-item editor-item-first"></button>
            <button data-toggle="tooltip" data-title="斜体" class="ql-italic editor-item"></button>
            <button data-toggle="tooltip" data-title="删除线" class="ql-strike editor-item"></button>
            <button data-toggle="tooltip" data-title="下划线" class="ql-underline editor-item editor-item-last"></button>
        </div>
        <div class="editor-group">
            <button data-toggle="tooltip" data-title="标题一" class="ql-header editor-item editor-item-first" value="1"></button>
            <button data-toggle="tooltip" data-title="标题二" class="ql-header editor-item" value="2"></button>
            <button data-toggle="tooltip" data-title="标题三" class="ql-header editor-item" value="3"></button>
            <button data-toggle="tooltip" data-title="标题四" class="ql-header editor-item" value="4"></button>
            <button data-toggle="tooltip" data-title="标题五" class="ql-header editor-item" value="5"></button>
            <button data-toggle="tooltip" data-title="标题六" class="ql-header editor-item editor-item-last" value="6"></button>
        </div>
        <div class="editor-group">
            <button data-toggle="tooltip" data-title="无序列表" class="ql-list editor-item editor-item-first" value="ordered"></button>
            <button data-toggle="tooltip" data-title="有序列表" class="ql-list editor-item" value="bullet"></button>
            <button data-toggle="tooltip" data-title="右缩进" class="ql-indent editor-item" value="-1"></button>
            <button data-toggle="tooltip" data-title="左缩进" class="ql-indent editor-item" value="+1"></button>
            <button data-toggle="tooltip" data-title="下标" class="ql-script editor-item" value="sub"></button>
            <button data-toggle="tooltip" data-title="上标" class="ql-script editor-item editor-item-last" value="super"></button>
        </div>
        <div class="editor-group ql-formats">
            <button data-toggle="tooltip" data-title="链接" class="ql-link editor-item editor-item-first"></button>
            <button data-toggle="tooltip" data-title="清空格式" class="ql-clean editor-item"></button>
            <button data-toggle="tooltip" data-title="添加图片" class="ql-image editor-item"></button>
            <button data-toggle="tooltip" data-title="添加视频" class="ql-video editor-item"></button>
            <button data-toggle="tooltip" data-title="代码块" class="ql-code-block editor-item"></button>
            <button data-toggle="tooltip" data-title="引用" class="ql-blockquote editor-item"><i class="fa fa-quote-right item" name="quote" unselectable="on"></i></button>
            <button data-toggle="tooltip" data-title="公式" class="ql-formula editor-item"><i class="fa fa-tasks item" name="tasks" aria-hidden="true"></i></button>
            <select data-toggle="tooltip" data-title="字体颜色" class="ql-color ql-picker ql-color-picker editor-item-select" ></select>
            <select data-toggle="tooltip" data-title="背景颜色" class="ql-background editor-item-select"></select>
            <a href="javascript:;" data-toggle="tooltip" data-title="附件" id="btnUploadFile"><i class="fa fa-paperclip last" aria-hidden="true" name="attachment"></i></a>

        </div>

        <div class="clearfix"></div>
    </div>
    <div class="manual-body">
        <div class="manual-category" id="manualCategory" style=" border-right: 1px solid #DDDDDD;width: 281px;position: absolute;">
            <div class="manual-nav">
                <div class="nav-item active"><i class="fa fa-bars" aria-hidden="true"></i> 文档</div>
                <div class="nav-plus pull-right" id="btnAddDocument" data-toggle="tooltip" data-title="创建文档" data-direction="right"><i class="fa fa-plus" aria-hidden="true"></i></div>
                <div class="clearfix"></div>
            </div>
            <div class="manual-tree" id="sidebar"> </div>
        </div>
        <div class="manual-editor-container" id="manualEditorContainer" style="min-width: 980px;">
            <div class="manual-editormd" style="bottom: 0;">
                <div id="docEditor" class="manual-editormd-active ql-editor ql-blank  editor-content"></div>
                <div class="manual-editor-status" style="border-top: 1px solid #DDDDDD;">
                    <div id="attachInfo" class="item">0 个附件</div>
                </div>
            </div>
        </div>
    </div>
</div>
<!-- 添加文档 -->
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
                        <label class="col-sm-2 control-label">文档标识</label>
                        <div class="col-sm-10">
                            <input type="text" name="doc_identify" id="documentIdentify" placeholder="文档唯一标识" class="form-control" maxlength="50">
                            <p style="color: #999;font-size: 12px;">文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头</p>
                        </div>

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
<!--附件上传-->
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
                        <div class="upload-container" id="filePicker">
                                <i class="fa fa-upload" aria-hidden="true"></i>
                        </div>
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
<!-- Modal -->
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

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jstree/3.3.4/jstree.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/katex/katex.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/to-markdown/dist/to-markdown.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/quill/quill.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/quill/quill.icons.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/editor.md/lib/highlight/highlight.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/array.js" "version"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/editor.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/quill.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        hljs.configure({   // optionally configure hljs
            languages: ['javascript', 'ruby', 'python']
        });
        $(".editor-code").on("dblclick",function () {
            var code = $(this).html();
            $("#createCodeToolbarModal").find("textarea").val(code);
            $("#createCodeToolbarModal").modal("show");
        }).on("click",function (e) {
            e.preventDefault();
            e.stopPropagation();
            console.log($(this).parents(".editor-wrapper").html())
            $(this).parents(".editor-wrapper").addClass("editor-wrapper-selected");
        });

        $("#attachInfo,#btnUploadFile").on("click",function () {
            $("#uploadAttachModal").modal("show");
        });


        /**
         * 文件上传
         */
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