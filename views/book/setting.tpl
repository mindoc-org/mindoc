<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n $.Lang "common.setting"}} - {{.Model.BookName}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/cropper/2.3.4/cropper.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/tagsinput/bootstrap-tagsinput.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-switch/css/bootstrap3//bootstrap-switch.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">

</head>
<body>
<div class="manual-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> {{i18n $.Lang "blog.summary"}}</a> </li>
                    <li><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> {{i18n $.Lang "blog.member"}}</a> </li>
                    <li><a href="{{urlfor "BookController.Team" ":key" .Model.Identify}}" class="item"><i class="fa fa-group" aria-hidden="true"></i> {{i18n $.Lang "blog.team"}}</a> </li>
                    <li class="active"><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> {{i18n $.Lang "common.setting"}}</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> {{i18n $.Lang "blog.project_setting"}}</strong>
                        {{if eq .Model.RoleId 0}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#transferBookModal">{{i18n $.Lang "blog.handover_project"}}</button>
                        {{if eq .Model.PrivatelyOwned 1}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">{{i18n $.Lang "blog.make_public"}}</button>
                        {{else}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">{{i18n $.Lang "blog.make_private"}}</button>
                        {{end}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;" data-toggle="modal" data-target="#deleteBookModal">{{i18n $.Lang "blog.delete_project"}}</button>
                        {{end}}

                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm" action="{{urlfor "BookController.SaveBook"}}">
                            <input type="hidden" name="identify" value="{{.Model.Identify}}">
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.project_title"}}</label>
                                <input type="text" class="form-control" name="book_name" id="bookName" placeholder="{{i18n $.Lang "blog.project_title"}}" value="{{.Model.BookName}}">
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.project_id"}}</label>
                                <input type="text" class="form-control" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" placeholder="{{i18n $.Lang "blog.project_id"}}" disabled>
                                <p class="text">{{i18n $.Lang "message.project_id_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "common.project_space"}}</label>
                                <select class="js-data-example-ajax form-control" multiple="multiple" name="itemId">
                                    <option value="{{.Model.ItemId}}" selected="selected">{{.Model.ItemName}}</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.history_record_amount"}}</label>
                                <input type="text" class="form-control" name="history_count" value="{{.Model.HistoryCount}}" placeholder="{{i18n $.Lang "blog.history_record_amount"}}">
                                <p class="text">{{i18n $.Lang "message.history_record_amount_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.corp_id"}}</label>
                                <input type="text" class="form-control" name="publisher" value="{{.Model.Publisher}}" placeholder="{{i18n $.Lang "blog.corp_id"}}">
                                <p class="text">{{i18n $.Lang "message.corp_id_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.project_desc"}}</label>
                                <textarea rows="3" class="form-control" name="description" style="height: 90px" placeholder="{{i18n $.Lang "blog.project_desc"}}">{{.Model.Description}}</textarea>
                                <p class="text">{{i18n $.Lang "message.project_desc_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n $.Lang "blog.text_editor"}}</label>
                                <div class="radio">
                                    <label class="radio-inline">
                                        <input type="radio"{{if eq .Model.Editor "markdown"}} checked{{end}} name="editor" value="markdown"> Markdown {{i18n $.Lang "blog.text_editor"}}
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio"{{if eq .Model.Editor "html"}} checked{{end}} name="editor" value="html"> Html {{i18n $.Lang "blog.text_editor"}}
                                    </label>
                                </div>
                            </div>
                {{if eq .Model.PrivatelyOwned 1}}
                <div class="form-group">
                    <label>{{i18n $.Lang "blog.access_pass"}}</label>
                    <input type="text" name="bPassword" id="bPassword" class="form-control" placeholder="{{i18n $.Lang "blog.access_pass"}}" value="{{.Model.BookPassword}}">
                    <p class="text">{{i18n $.Lang "message.access_pass_desc"}}</p>
                </div>
                {{end}}

                <div class="form-group">
                    <label for="autoRelease">{{i18n $.Lang "blog.auto_publish"}}</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="autoRelease" name="auto_release"{{if .Model.AutoRelease }} checked{{end}} data-size="small">
                        </div>
                        <p class="text">{{i18n $.Lang "message.auto_publish_desc"}}</p>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">{{i18n $.Lang "blog.enable_export"}}</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="isDownload" name="is_download"{{if .Model.IsDownload }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.enable_export"}}">
                        </div>
                        <p class="text">{{i18n $.Lang "message.enable_export_desc"}}</p>

                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">{{i18n $.Lang "blog.enable_share"}}</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="enableShare" name="enable_share"{{if .Model.IsEnableShare }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.enable_share"}}">
                        </div>
                        <p class="text">{{i18n $.Lang "message.enable_share_desc"}}</p>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">{{i18n $.Lang "blog.set_first_as_home"}}</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="isUseFirstDocument" name="is_use_first_document"{{if .Model.IsUseFirstDocument }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.set_first_as_home"}}">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">{{i18n $.Lang "blog.auto_save"}}</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="autoSave" name="auto_save"{{if .Model.AutoSave }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.auto_save"}}">
                        </div>
                        <p class="text">{{i18n $.Lang "message.auto_save_desc"}}</p>
                    </div>
                </div>
                <div class="form-group">
                    <button type="submit" id="btnSaveBookInfo" class="btn btn-success" data-loading-text="{{i18n $.Lang "common.processing"}}">{{i18n $.Lang "common.save"}}</button>
                    <span id="form-error-message" class="error-message"></span>
                </div>
                </form>
            </div>
            <div class="form-right">
                <label>
                    <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                        <img src="{{cdnimg .Model.Cover}}" onerror="this.src='{{cdnimg "/static/images/book.png"}}'" alt="{{i18n $.Lang "blog.cover"}}" style="max-width: 120px;border: 1px solid #999" id="headimgurl">
                    </a>
                </label>
                <p class="text">{{i18n $.Lang "blog.click_change_cover"}}</p>
            </div>
            <div class="clearfix"></div>

        </div>
    </div>
</div>
</div>
{{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="changePrivatelyOwnedModal" tabindex="-1" role="dialog" aria-labelledby="changePrivatelyOwnedModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "BookController.PrivatelyOwned" }}" id="changePrivatelyOwnedForm">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <input type="hidden" name="status" value="{{if eq .Model.PrivatelyOwned 0}}close{{else}}open{{end}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">
                        {{if eq .Model.PrivatelyOwned 0}}
                        {{i18n $.Lang "blog.make_private"}}
                        {{else}}
                        {{i18n $.Lang "blog.make_public"}}
                        {{end}}
                    </h4>
                </div>
                <div class="modal-body">
                    {{if eq .Model.PrivatelyOwned 0}}
                    <span style="font-size: 14px;font-weight: 400;">{{i18n $.Lang "message.confirm_into_private"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n $.Lang "message.into_private_notice"}}</p>
                    {{else}}
                    <span style="font-size: 14px;font-weight: 400;">{{i18n $.Lang "message.confirm_into_public"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n $.Lang "message.into_public_notice"}}</p>
                    {{end}}
                </div>
                <div class="modal-footer">
                    <span class="error-message" id="form-error-message1"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-primary" data-loading-text="{{i18n $.Lang "common.processing"}}" id="btnChangePrivatelyOwned">{{i18n $.Lang "common.confirm"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!-- Start Modal -->
<div class="modal fade" id="upload-logo-panel" tabindex="-1" role="dialog" aria-labelledby="{{i18n $.Lang "blog.change_cover"}}" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title">{{i18n $.Lang "blog.change_cover"}}</h4>
            </div>
            <div class="modal-body">
                <div class="wraper">
                    <div id="image-wraper">

                    </div>
                </div>
                <div class="watch-crop-list">
                    <div class="preview-title">{{i18n $.Lang "blog.preview"}}</div>
                    <ul>
                        <li>
                            <div class="img-preview preview-lg"></div>
                        </li>
                        <li>
                            <div class="img-preview preview-sm"></div>
                        </li>
                    </ul>
                </div>
                <div style="clear: both"></div>
            </div>
            <div class="modal-footer">
                <span id="error-message"></span>
                <div id="filePicker" class="btn">{{i18n $.Lang "blog.choose"}}</div>
                <button type="button" id="saveImage" class="btn btn-success" style="height: 40px;width: 77px;" data-loading-text="{{i18n $.Lang "blog.processing"}}">{{i18n $.Lang "blog.upload"}}</button>
            </div>
        </div>
    </div>
</div>
<!--END Modal-->

<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBookModal" tabindex="-1" role="dialog" aria-labelledby="deleteBookModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBookForm" action="{{urlfor "BookController.Delete"}}">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">{{i18n $.Lang "blog.delete_project"}}</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">{{i18n $.Lang "confirm_delete_project"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n $.Lang "message.warning_delete_project"}}</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="{{i18n $.Lang "common.processing"}}">{{i18n $.Lang "common.confirm_delete"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!-- Modal -->
<div class="modal fade" id="transferBookModal" tabindex="-1" role="dialog" aria-labelledby="transferBookModalLabel">
    <div class="modal-dialog" role="document">
        <form action="{{urlfor "BookController.Transfer"}}" method="post" id="transferBookForm">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n $.Lang "blog.handover_project"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n $.Lang "blog.recipient_account"}}</label>
                        <div class="col-sm-10">
                            <input type="text" name="account" class="form-control" placeholder="{{i18n $.Lang "blog.recipient_account"}}" id="receiveAccount" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message3" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                    <button type="submit" id="btnTransferBook" class="btn btn-primary">{{i18n $.Lang "common.comfirm"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/cropper/2.3.4/cropper.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/tagsinput/bootstrap-tagsinput.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-switch/js/bootstrap-switch.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#upload-logo-panel").on("hidden.bs.modal",function () {
            $("#upload-logo-panel").find(".modal-body").html(window.modalHtml);
        }).on("show.bs.modal",function () {
            window.modalHtml = $("#upload-logo-panel").find(".modal-body").html();
        });
        $("#autoRelease,#enableShare,#isDownload,#isUseFirstDocument,#autoSave").bootstrapSwitch();

        $('input[name="label"]').tagsinput({
            confirmKeys: [13,44],
            maxTags: 10,
            trimValue: true,
            cancelConfirmKeysOnEmpty : false
        });

        $("#changePrivatelyOwnedForm").ajaxForm({
            beforeSubmit :function () {
                $("#btnChangePrivatelyOwned").button("loading");
            },
            success :function (res) {
                if(res.errcode === 0){
                    window.location = window.location.href;
                    return;
                }else{
                    showError(res.message,"#form-error-message1");
                }
                $("#btnChangePrivatelyOwned").button("reset");
            },
            error :function () {
                showError("{{i18n $.Lang "message.system_error"}}","#form-error-message1");
                $("#btnChangePrivatelyOwned").button("reset");
            }
        });

        $("#createToken,#deleteToken").on("click",function () {
            var btn = $(this).button("loading");
            var action = $(this).attr("data-action");
            $.ajax({
                url : "{{urlfor "BookController.CreateToken"}}",
                type :"post",
                data : { "identify" : {{.Model.Identify}} , "action" : action },
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        $("#token").val(res.data);
                    }else{
                        alert(res.message);
                    }
                    btn.button("reset");
                },
                error : function () {
                    btn.button("reset");
                    alert("{{i18n $.Lang "message.system_error"}}");
                }
            }) ;
        });
        $("#token").on("focus",function () {
            $(this).select();
        });
        $("#bookEditForm").ajaxForm({
            beforeSubmit : function () {
                var bookName = $.trim($("#bookName").val());
                if (bookName === "") {
                    return showError("{{i18n $.Lang "message.project_name_empty"}}");
                }
                $("#btnSaveBookInfo").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess("{{i18n $.Lang "message.success"}}")
                }else{
                    showError("{{i18n $.Lang "message.failed"}}")
                }
                $("#btnSaveBookInfo").button("reset");
            },
            error : function () {
                showError("{{i18n $.Lang "message.system_error"}}");
                $("#btnSaveBookInfo").button("reset");
            }
        });
        $("#deleteBookForm").ajaxForm({
            beforeSubmit : function () {
                $("#btnDeleteBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    window.location = "{{urlfor "BookController.Index"}}";
                }else{
                    showError(res.message,"#form-error-message2");
                }
                $("#btnDeleteBook").button("reset");
            },
            error : function () {
                showError("{{i18n $.Lang "message.system_error"}}","#form-error-message2");
                $("#btnDeleteBook").button("reset");
            }
        });
        $("#transferBookForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#receiveAccount").val());
                if (account === ""){
                    return showError("{{i18n $.Lang "message.receive_account_empty"}}","#form-error-message3")
                }
                $("#btnTransferBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    window.location = window.location.href;
                }else{
                    showError(res.message,"#form-error-message3");
                }
                $("#btnTransferBook").button("reset");
            },
            error : function () {
                $("#btnTransferBook").button("reset");
            }
        });
        $('.js-data-example-ajax').select2({
            language: "{{i18n $.Lang "common.js_lang"}}",
            minimumInputLength : 1,
            minimumResultsForSearch: Infinity,
            maximumSelectionLength:1,
            width : "100%",
            ajax: {
                url: '{{urlfor "BookController.ItemsetsSearch"}}',
                dataType: 'json',
                data: function (params) {
                    return {
                        q: params.term, // search term
                        page: params.page
                    };
                },
                processResults: function (data, params) {
                    return {
                        results : data.data.results
                    }
                }
            }
        });
        try {
            var uploader = WebUploader.create({
                auto: false,
                swf: '{{.BaseUrl}}/static/webuploader/Uploader.swf',
                server: '{{urlfor "BookController.UploadCover"}}',
                formData : { "identify" : {{.Model.Identify}} },
                pick: "#filePicker",
                fileVal : "image-file",
                fileNumLimit : 1,
                compress : false,
                accept: {
                    title: 'Images',
                    extensions: 'jpg,jpeg,png',
                    mimeTypes: 'image/jpg,image/jpeg,image/png'
                }
            }).on("beforeFileQueued",function (file) {
                uploader.reset();
            }).on( 'fileQueued', function( file ) {
                uploader.makeThumb( file, function( error, src ) {
                    $img = '<img src="' + src +'" style="max-width: 360px;max-height: 360px;">';
                    if ( error ) {
                        $img.replaceWith('<span>{{i18n $.Lang "message.cannot_preview"}}</span>');
                        return;
                    }

                    $("#image-wraper").html($img);
                    window.ImageCropper = $('#image-wraper>img').cropper({
                        aspectRatio: 175 / 230,
                        dragMode : 'move',
                        viewMode : 1,
                        preview : ".img-preview"
                    });
                }, 1, 1 );
            }).on("uploadError",function (file,reason) {
                console.log(reason);
                $("#error-message").text("{{i18n $.Lang "message.upload_failed"}}:" + reason);

            }).on("uploadSuccess",function (file, res) {

                if(res.errcode === 0){
                    console.log(res);
                    $("#upload-logo-panel").modal('hide');
                    $("#headimgurl").attr('src',res.data);
                }else{
                    $("#error-message").text(res.message);
                }
            }).on("beforeFileQueued",function (file) {
                if(file.size > 1024*1024*2){
                    uploader.removeFile(file);
                    uploader.reset();
                    alert("{{i18n $.Lang "message.upload_file_size_limit"}}");
                    return false;
                }
            }).on("uploadComplete",function () {
                $("#saveImage").button('reset');
            });
            $("#saveImage").on("click",function () {
                var files = uploader.getFiles();
                if(files.length > 0) {
                    $("#saveImage").button('loading');
                    var cropper = window.ImageCropper.cropper("getData");

                    uploader.option("formData", cropper);

                    uploader.upload();
                }else{
                    alert("{{i18n $.Lang "message.choose_pic_file"}}");
                }
            });
        }catch(e){
            console.log(e);
        }
    });
</script>
</body>
</html>