<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.edit_proj"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/tagsinput/bootstrap-tagsinput.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-switch/css/bootstrap3//bootstrap-switch.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/cropper/2.3.4/cropper.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
</head>
<body>
<div class="manual-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
        {{template "manager/widgets.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> {{i18n .Lang "blog.project_setting"}}</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#transferBookModal">{{i18n .Lang "blog.handover_project"}}</button>
                        {{if eq .Model.PrivatelyOwned 1}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">{{i18n .Lang "blog.make_public"}}</button>
                        {{else}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">{{i18n .Lang "blog.make_private"}}</button>
                        {{end}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;" data-toggle="modal" data-target="#deleteBookModal">{{i18n .Lang "blog.delete_project"}}</button>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm" action="{{urlfor "ManagerController.EditBook" ":key" .Model.Identify}}">
                            <input type="hidden" name="identify" value="{{.Model.Identify}}">
                            <div class="form-group">
                                <label>{{i18n .Lang "mgr.proj_name"}}</label>
                                <input type="text" class="form-control" name="book_name" id="bookName" placeholder="{{i18n .Lang "mgr.proj_name"}}" value="{{.Model.BookName}}">
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.project_id"}}</label>
                                <input type="text" class="form-control" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" disabled placeholder="{{i18n .Lang "blog.project_id"}}">
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "common.project_space"}}</label>
                                <select class="js-data-example-ajax form-control" multiple="multiple" name="itemId">
                                    <option value="{{.Model.ItemId}}" selected="selected">{{.Model.ItemName}}</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.history_record_amount"}}</label>
                                <input type="text" class="form-control" name="history_count" value="{{.Model.HistoryCount}}" placeholder="{{i18n .Lang "blog.history_record_amount"}}">
                                <p class="text">{{i18n .Lang "message.history_record_amount_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.corp_id"}}</label>
                                <input type="text" class="form-control" name="publisher" value="{{.Model.Publisher}}" placeholder="{{i18n $.Lang "blog.corp_id"}}">
                                <p class="text">{{i18n .Lang "message.corp_id_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.project_order"}}</label>
                                <input type="number" min="0" class="form-control" value="{{.Model.OrderIndex}}" name="order_index" placeholder="{{i18n .Lang "blog.project_order"}}">
                                <p class="text">{{i18n .Lang "message.project_order_desc"}}</p>
                            </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.project_desc"}}</label>
                                <textarea rows="3" class="form-control" name="description" style="height: 90px" placeholder="{{i18n .Lang "blog.project_desc"}}">{{.Model.Description}}</textarea>
                                <p class="text">{{i18n .Lang "message.project_desc_placeholder"}}</p>
                            </div>

                            <div class="form-group">
                                <label>{{i18n .Lang "blog.project_label"}}</label>
                                <input type="text" class="form-control" name="label" placeholder="{{i18n .Lang "blog.project_label"}}" value="{{.Model.Label}}">
                                <p class="text">{{i18n .Lang "message.project_label_desc"}}</p>
                            </div>
                            {{if eq .Model.PrivatelyOwned 1}}
                            <div class="form-group">
                                <label>{{i18n .Lang "blog.access_token"}}</label>
                                <div class="row">
                                    <div class="col-sm-9">
                                        <input type="text" name="token" id="token" class="form-control" placeholder="{{i18n .Lang "blog.access_token"}}" readonly value="{{.Model.PrivateToken}}">
                                    </div>
                                    <div class="col-sm-3">
                                        <button type="button" class="btn btn-success btn-sm" id="createToken" data-loading-text="{{i18n .Lang "common.generate"}}" data-action="create">{{i18n .Lang "common.generate"}}</button>
                                        <button type="button" class="btn btn-danger btn-sm" id="deleteToken" data-loading-text="{{i18n .Lang "common.delete"}}" data-action="delete">{{i18n .Lang "common.delete"}}</button>
                                    </div>
                                </div>
                            </div>
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
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">{{i18n $.Lang "blog.enable_export"}}</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="isDownload" name="is_download"{{if .Model.IsDownload }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.enable_export"}}">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">{{i18n $.Lang "blog.enable_share"}}</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="enableShare" name="enable_share"{{if .Model.IsEnableShare }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.enable_share"}}">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">{{i18n $.Lang "blog.set_first_as_home"}}</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="is_use_first_document" name="is_use_first_document"{{if .Model.IsUseFirstDocument }} checked{{end}} data-size="small" placeholder="{{i18n $.Lang "blog.set_first_as_home"}}">
                                    </div>
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
                           <img src="{{.Model.Cover}}" onerror="this.src='/static/images/book.png'" alt="{{i18n .Lang "blog.cover"}}" style="max-width: 120px;border: 1px solid #999" id="headimgurl">
                        </label>
                    </div>
                    <div class="clearfix"></div>

                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<div class="modal fade" id="changePrivatelyOwnedModal" tabindex="-1" role="dialog" aria-labelledby="changePrivatelyOwnedModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "ManagerController.PrivatelyOwned" }}" id="changePrivatelyOwnedForm">
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
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-primary" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnChangePrivatelyOwned">{{i18n .Lang "common.confirm"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>

<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBookModal" tabindex="-1" role="dialog" aria-labelledby="deleteBookModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBookForm" action="{{urlfor "ManagerController.DeleteBook"}}">
            <input type="hidden" name="book_id" value="{{.Model.BookId}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">{{i18n .Lang "blog.delete_project"}}</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">{{i18n .Lang "message.confirm_delete_project"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n .Lang "message.warning_delete_project"}}</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.confirm_delete"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<div class="modal fade" id="transferBookModal" tabindex="-1" role="dialog" aria-labelledby="transferBookModalLabel">
    <div class="modal-dialog" role="document">
        <form action="{{urlfor "ManagerController.Transfer"}}" method="post" id="transferBookForm">
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
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" id="btnTransferBook" daata-loading-text="{{i18n $.Lang "message.processing"}}" class="btn btn-primary">{{i18n $.Lang "common.comfirm"}}</button>
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
        $("#autoRelease,#enableShare,#isDownload,#is_use_first_document").bootstrapSwitch();
        $("#upload-logo-panel").on("hidden.bs.modal",function () {
            $("#upload-logo-panel").find(".modal-body").html(window.modalHtml);
        }).on("show.bs.modal",function () {
            window.modalHtml = $("#upload-logo-panel").find(".modal-body").html();
        });
        $("#createToken,#deleteToken").on("click",function () {
            var btn = $(this).button("loading");
            var action = $(this).attr("data-action");
            $.ajax({
                url : "{{urlfor "ManagerController.CreateToken"}}",
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
                    alert({{i18n $.Lang "message.system_error"}});
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
                    window.location = "{{urlfor "ManagerController.Books"}}";
                }else{
                    $("#btnDeleteBook").button("reset");
                    showError(res.message,"#form-error-message2");
                }
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
    });

</script>
</body>
</html>