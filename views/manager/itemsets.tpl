<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.project_space_mgr"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

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
                        <strong class="box-title">{{i18n .Lang "mgr.project_space_mgr"}}</strong>
                    {{if eq .Member.Role 0}}
                        <button type="button" class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addItemsetsDialogModal"><i class="fa fa-plus" aria-hidden="true"></i> {{i18n .Lang "mgr.create_proj_space"}}</button>
                    {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="attach-list" id="ItemsetsList">
                        <table class="table">
                            <thead>
                            <tr>
                                <th width="10%">#</th>
                                <th width="30%">{{i18n .Lang "mgr.proj_space_name"}}</th>
                                <th width="20%">{{i18n .Lang "mgr.proj_space_id"}}</th>
                                <th width="20%">{{i18n .Lang "mgr.proj_amount"}}</th>
                                <th>{{i18n .Lang "common.operate"}}</th>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index,$item := .Lists}}
                            <tr>
                                <td>{{$item.ItemId}}</td>
                                <td>{{$item.ItemName}}</td>
                                <td>{{$item.ItemKey}}</td>
                                <td>{{$item.BookNumber}}</td>
                                <td>
                                    <button type="button" class="btn btn-sm btn-default" data-id="{{$item.ItemId}}" data-method="edit" data-name="{{$item.ItemName}}" data-key="{{$item.ItemKey}}">{{i18n $.Lang "common.edit"}}</button>
                                    {{if ne $item.ItemId 1}}
                                    <button type="button" data-method="delete" class="btn btn-danger btn-sm" data-id="{{$item.ItemId}}" data-loading-text="{{i18n $.Lang "message.processing"}}">{{i18n $.Lang "common.delete"}}</button>
                                    {{end}}
                                    <a href="{{urlfor "ItemsetsController.List" ":key" $item.ItemKey}}" class="btn btn-success btn-sm" target="_blank">{{i18n $.Lang "common.detail"}}</a>
                                </td>
                            </tr>
                            {{else}}
                            <tr><td class="text-center" colspan="6">{{i18n .Lang "message.no_data"}}</td></tr>
                            {{end}}
                            </tbody>
                        </table>
                        <nav class="pagination-container">
                        {{.PageHtml}}
                        </nav>
                    </div>
                </div>
            </div>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="addItemsetsDialogModal" tabindex="-1" role="dialog" aria-labelledby="addItemsetsDialogModalLabel">
    <div class="modal-dialog">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.ItemsetsEdit"}}" id="addItemsetsDialogForm">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.create_proj_space"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-4 control-label" for="account">{{i18n .Lang "mgr.proj_space_name"}}<span class="error-message">*</span></label>
                        <div class="col-sm-8">
                            <input type="text" name="itemName" class="form-control" placeholder="{{i18n .Lang "mgr.proj_space_name"}}" id="itemName" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-4 control-label" for="itemKey">{{i18n .Lang "mgr.proj_space_id"}}<span class="error-message">*</span></label>
                        <div class="col-sm-8">
                            <input type="text" name="itemKey" id="itemKey" class="form-control" placeholder="{{i18n .Lang "mgr.proj_space_id"}}" maxlength="50">
                            <p class="text">{{i18n .Lang "message.proj_space_id_tips"}}</p>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="create-form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnAddItemsets">{{i18n .Lang "common.save"}}
                    </button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->
<div class="modal fade" id="editItemsetsDialogModal" tabindex="-1" role="dialog" aria-labelledby="editItemsetsDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.ItemsetsEdit"}}" id="editItemsetsDialogForm">
            <input type="hidden" name="itemId" value="">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.edit_proj_space"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-4 control-label" for="itemName">{{i18n .Lang "mgr.proj_space_name"}}<span class="error-message">*</span></label>
                        <div class="col-sm-8">
                            <input type="text" name="itemName" id="itemName" class="form-control" placeholder="{{i18n .Lang "mgr.proj_space_name"}}" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-4 control-label" for="itemKey">{{i18n .Lang "mgr.proj_space_id"}}<span class="error-message">*</span></label>
                        <div class="col-sm-8">
                            <input type="text" name="itemKey" id="itemKey" class="form-control" placeholder="{{i18n .Lang "mgr.proj_space_id"}}" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="edit-form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnEditItemsets">{{i18n .Lang "common.save"}}
                    </button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js" }}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>

<script type="text/javascript">
    $(function () {
        var editItemsetsDialogModal = $("#editItemsetsDialogModal");
        var addItemsetsDialogForm = $("#addItemsetsDialogForm");
        var editItemsetsDialogForm = $("#editItemsetsDialogForm");

        editItemsetsDialogModal.on("shown.bs.modal",function () {
            editItemsetsDialogModal.find("input[name='itemName']").focus();
        });
        $("#addItemsetsDialogModal").on("show.bs.modal", function () {
            window.addItemsetsDialogModalHtml = $(this).find("form").html();
        }).on("hidden.bs.modal", function () {
            $(this).find("form").html(window.addItemsetsDialogModalHtml);
        });

        addItemsetsDialogForm.ajaxForm({
            beforeSubmit: function () {
                var $itemName = addItemsetsDialogForm.find("input[name='itemName']").val();
                var $itemKey =  addItemsetsDialogForm.find("input[name='itemKey']").val();

                if ($itemName == "") {
                    showError("{{i18n .Lang "message.proj_space_name_empty"}}","#create-form-error-message");
                }
                if ($itemKey == "") {
                    showError("{{i18n .Lang "message.proj_space_id_empty"}}","#create-form-error-message");
                }
                $("#btnAddItemsets").button("loading");
                showError("","#create-form-error-message");
                return true;
            },
            success: function ($res) {
                if ($res.errcode === 0) {
                    window.location = window.document.location;
                } else {
                    showError($res.message,"#create-form-error-message");
                }
            },
            error: function () {
                showError({{i18n .Lang "message.system_error"}},"#create-form-error-message");
            },
            complete: function () {
                $("#btnAddItemsets").button("reset");
            }
        });

        editItemsetsDialogForm.ajaxForm({
           beforeSubmit: function () {
               var $itemName = editItemsetsDialogForm.find("input[name='itemName']").val();
               var $itemKey =  editItemsetsDialogForm.find("input[name='itemKey']").val();

               if ($itemName == "") {
                   showError("{{i18n .Lang "message.proj_space_name_empty"}}","#edit-form-error-message");
               }
               if ($itemKey == "") {
                   showError("{{i18n .Lang "message.proj_space_id_empty"}}","#edit-form-error-message");
               }
               $("#btnEditItemsets").button("loading");
               showError("","#edit-form-error-message");
               return true;
           } ,
            success : function ($res) {
                if ($res.errcode === 0) {
                    window.location = window.document.location;
                } else {
                    showError($res.message,"#edit-form-error-message");
                }
            },
            error: function () {
                showError({{i18n .Lang "message.system_error"}},"#edit-form-error-message");
            },
            complete: function () {
                $("#btnEditItemsets").button("reset");
            }
        });

        $("#ItemsetsList").on("click","button[data-method='delete']",function () {
            var id = $(this).attr("data-id");
            var $this = $(this);
            $(this).button("loading");
            $.ajax({
                url : "{{urlfor "ManagerController.ItemsetsDelete"}}",
                data: {"itemId":id},
                type : "post",
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        $this.closest("tr").remove().empty();
                    }else {
                        layer.msg(res.message);
                    }
                },
                error : function () {
                    layer.msg({{i18n .Lang "message.system_error"}});
                },
                complete : function () {
                    $this.button("reset");
                }
            });
        }).on("click","button[data-method='edit']",function () {
            var $itemId = $(this).attr("data-id");
            var $itemName = $(this).attr("data-name");
            var $itemKey = $(this).attr("data-key");

            editItemsetsDialogModal.find("input[name='itemId']").val($itemId);
            editItemsetsDialogModal.find("input[name='itemName']").val($itemName);
            editItemsetsDialogModal.find("input[name='itemKey']").val($itemKey);

            editItemsetsDialogModal.modal("show");
        });
    });
</script>
</body>
</html>