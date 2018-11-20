<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>项目集管理 - Powered by MinDoc</title>

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
        {{template "manager/widgets.tpl" "itemsets"}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">项目集管理</strong>
                    {{if eq .Member.Role 0}}
                        <button type="button" class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addItemsetsDialogModal"><i class="fa fa-plus" aria-hidden="true"></i> 创建项目集</button>
                    {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="attach-list" id="ItemsetsList">
                        <table class="table">
                            <thead>
                            <tr>
                                <th width="10%">#</th>
                                <th width="30%">项目集名称</th>
                                <th width="20%">项目集标识</th>
                                <th width="20%">项目数量</th>
                                <th>操作</th>
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
                                    <button type="button" class="btn btn-sm btn-default" data-id="{{$item.ItemId}}" data-method="edit" data-name="{{$item.ItemName}}" data-key="{{$item.ItemKey}}">编辑</button>
                                    <button type="button" data-method="delete" class="btn btn-danger btn-sm" data-id="{{$item.ItemId}}" data-loading-text="删除中...">删除</button>
                                    <a href="{{urlfor "ItemsetsController.Index" ":key" $item.ItemKey}}" class="btn btn-success btn-sm" target="_blank">详情</a>
                                </td>
                            </tr>
                            {{else}}
                            <tr><td class="text-center" colspan="6">暂无数据</td></tr>
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
                    <h4 class="modal-title" id="myModalLabel">创建项目集</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-3 control-label" for="account">项目集名称<span class="error-message">*</span></label>
                        <div class="col-sm-9">
                            <input type="text" name="itemName" class="form-control" placeholder="项目集名称" id="itemName" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-3 control-label" for="itemKey">项目集标识<span class="error-message">*</span></label>
                        <div class="col-sm-9">
                            <input type="text" name="itemKey" id="itemKey" class="form-control" placeholder="项目集标识" maxlength="50">
                            <p class="text">项目集标识只能由字母和数字组成且在2-100字符之间</p>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="create-form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-success" data-loading-text="保存中..." id="btnAddItemsets">保存
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
                    <h4 class="modal-title" id="myModalLabel">编辑项目集</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="itemName">项目集名称<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="itemName" id="itemName" class="form-control" placeholder="项目集名称" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="itemKey">项目集标识<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="itemKey" id="itemKey" class="form-control" placeholder="项目集标识" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="edit-form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-success" data-loading-text="保存中..." id="btnEditItemsets">保存
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
                    showError("项目集名称不能为空","#create-form-error-message");
                }
                if ($itemKey == "") {
                    showError("项目集标识不能为空","#create-form-error-message");
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
                showError("服务器异常","#create-form-error-message");
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
                   showError("项目集名称不能为空","#edit-form-error-message");
               }
               if ($itemKey == "") {
                   showError("项目集标识不能为空","#edit-form-error-message");
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
                showError("服务器异常","#edit-form-error-message");
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
                url : "{{urlfor "ManagerController.ItemsetsDelete" ":id" ""}}" + id,
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
                    layer.msg("服务器异常");
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