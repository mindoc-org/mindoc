<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑项目 - Powered by MinDoc</title>

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
        {{template "manager/widgets.tpl" "books"}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 项目设置</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#transferBookModal">转让项目</button>
                        {{if eq .Model.PrivatelyOwned 1}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">转为公有</button>
                        {{else}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">转为私有</button>
                        {{end}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;" data-toggle="modal" data-target="#deleteBookModal">删除项目</button>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm" action="{{urlfor "ManagerController.EditBook" ":key" .Model.Identify}}">
                            <input type="hidden" name="identify" value="{{.Model.Identify}}">
                            <div class="form-group">
                                <label>标题</label>
                                <input type="text" class="form-control" name="book_name" id="bookName" placeholder="项目名称" value="{{.Model.BookName}}">
                            </div>
                            <div class="form-group">
                                <label>标识</label>
                                <input type="text" class="form-control" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" disabled placeholder="项目标识">
                            </div>
                            <div class="form-group">
                                <label>项目空间</label>
                                <select class="js-data-example-ajax form-control" multiple="multiple" name="itemId">
                                    <option value="{{.Model.ItemId}}" selected="selected">{{.Model.ItemName}}</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label>历史记录数量</label>
                                <input type="text" class="form-control" name="history_count" value="{{.Model.HistoryCount}}" placeholder="历史记录数量">
                                <p class="text">当开启文档历史时,该值会限制每个文档保存的历史数量</p>
                            </div>
                            <div class="form-group">
                                <label>公司标识</label>
                                <input type="text" class="form-control" name="publisher" value="{{.Model.Publisher}}" placeholder="公司名称">
                                <p class="text">导出文档PDF文档时显示的页脚</p>
                            </div>
                            <div class="form-group">
                                <label>排序</label>
                                <input type="number" min="0" class="form-control" value="{{.Model.OrderIndex}}" name="order_index" placeholder="项目排序">
                                <p class="text">只能是数字，序号越大排序越靠前</p>
                            </div>
                            <div class="form-group">
                                <label>描述</label>
                                <textarea rows="3" class="form-control" name="description" style="height: 90px" placeholder="项目描述">{{.Model.Description}}</textarea>
                                <p class="text">描述信息不超过500个字符</p>
                            </div>

                            <div class="form-group">
                                <label>标签</label>
                                <input type="text" class="form-control" name="label" placeholder="项目标签" value="{{.Model.Label}}">
                                <p class="text">最多允许添加10个标签，多个标签请用“;”分割</p>
                            </div>
                            {{if eq .Model.PrivatelyOwned 1}}
                            <div class="form-group">
                                <label>访问令牌</label>
                                <div class="row">
                                    <div class="col-sm-10">
                                        <input type="text" name="token" id="token" class="form-control" placeholder="访问令牌" readonly value="{{.Model.PrivateToken}}">
                                    </div>
                                    <div class="col-sm-2">
                                        <button type="button" class="btn btn-success btn-sm" id="createToken" data-loading-text="生成" data-action="create">生成</button>
                                        <button type="button" class="btn btn-danger btn-sm" id="deleteToken" data-loading-text="删除" data-action="delete">删除</button>
                                    </div>
                                </div>
                            </div>
                                <div class="form-group">
                                    <label>访问密码</label>
                                    <input type="text" name="bPassword" id="bPassword" class="form-control" placeholder="访问密码" value="{{.Model.BookPassword}}">
                                    <p class="text">没有访问权限访问项目时需要提供的密码</p>
                                </div>
                            {{end}}
                            <div class="form-group">
                                <label for="autoRelease">自动发布</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="autoRelease" name="auto_release"{{if .Model.AutoRelease }} checked{{end}} data-size="small">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">开启导出</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="isDownload" name="is_download"{{if .Model.IsDownload }} checked{{end}} data-size="small" placeholder="开启导出">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">开启分享</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="enableShare" name="enable_share"{{if .Model.IsEnableShare }} checked{{end}} data-size="small" placeholder="开启分享">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="autoRelease">设置第一篇文档为默认首页</label>
                                <div class="controls">
                                    <div class="switch switch-small" data-on="primary" data-off="info">
                                        <input type="checkbox" id="is_use_first_document" name="is_use_first_document"{{if .Model.IsUseFirstDocument }} checked{{end}} data-size="small" placeholder="设置第一篇文档为默认首页">
                                    </div>
                                </div>
                            </div>
                            <div class="form-group">
                                <button type="submit" id="btnSaveBookInfo" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                                <span id="form-error-message" class="error-message"></span>
                            </div>
                        </form>
                    </div>
                    <div class="form-right">
                        <label>
                           <img src="{{.Model.Cover}}" onerror="this.src='/static/images/book.png'" alt="封面" style="max-width: 120px;border: 1px solid #999" id="headimgurl">
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
                        转为私有
                        {{else}}
                        转为共有
                        {{end}}
                    </h4>
                </div>
                <div class="modal-body">
                    {{if eq .Model.PrivatelyOwned 0}}
                    <span style="font-size: 14px;font-weight: 400;">确定将项目转为私有吗？</span>
                    <p></p>
                    <p class="text error-message">转为私有后需要通过阅读令牌才能访问该项目。</p>
                    {{else}}
                    <span style="font-size: 14px;font-weight: 400;"> 确定将项目转为公有吗？</span>
                    <p></p>
                    <p class="text error-message">转为公有后所有人都可以访问该项目。</p>
                    {{end}}
                </div>
                <div class="modal-footer">
                    <span class="error-message" id="form-error-message1"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-primary" data-loading-text="正在保存..." id="btnChangePrivatelyOwned">确定</button>
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
                    <h4 class="modal-title">删除项目</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">确定删除项目吗？</span>
                    <p></p>
                    <p class="text error-message">删除项目后将无法找回。</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="正在删除...">确定删除</button>
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
                    <h4 class="modal-title" id="myModalLabel">项目转让</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">接收账号</label>
                        <div class="col-sm-10">
                            <input type="text" name="account" class="form-control" placeholder="接收者账号" id="receiveAccount" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message3" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" id="btnTransferBook" daata-loading-text="正在转让..." class="btn btn-primary">确定转让</button>
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
                    alert("服务器错误");
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
                    return showError("项目名称不能为空");
                }
                $("#btnSaveBookInfo").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess("保存成功")
                }else{
                    showError("保存失败")
                }
                $("#btnSaveBookInfo").button("reset");
            },
            error : function () {
                showError("服务错误");
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
                    return showError("接受者账号不能为空","#form-error-message3")
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
                showError("服务器异常","#form-error-message1");
                $("#btnChangePrivatelyOwned").button("reset");
            }
        });
        $('.js-data-example-ajax').select2({
            language: "zh-CN",
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