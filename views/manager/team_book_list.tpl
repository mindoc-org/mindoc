<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.team_proj"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">
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
                        <strong class="box-title">{{.Model.TeamName}} - {{i18n .Lang "mgr.team_proj"}}</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addTeamBookDialogModal"><i class="fa fa-book" aria-hidden="true"></i> {{i18n .Lang "mgr.add_proj"}}</button>
                    </div>
                </div>
                <div class="box-body">
                    <div class="attach-list" id="teamBookList">
                        <table class="table">
                            <thead>
                            <tr>
                                <th>{{i18n .Lang "mgr.proj_name"}}</th>
                                <th>{{i18n .Lang "mgr.proj_author"}}</th>
                                <th>{{i18n .Lang "mgr.join_time"}}</th>
                                <th>{{i18n .Lang "common.operate"}}</th>
                            </tr>
                            </thead>
                            <tbody>
                            <template v-if="lists.length <= 0">
                                <tr class="text-center"><td colspan="4">{{i18n .Lang "message.no_data"}}</td></tr>
                            </template>
                            <template v-else>
                            <tr v-for="item in lists">
                                <td>${item.book_name}</td>
                                <td>${item.book_member_name}</td>
                                <td>${(new Date(item.create_time)).format("yyyy-MM-dd hh:mm:ss")}</td>
                                <td>
                                    <button type="button" data-method="delete" class="btn btn-danger btn-sm" @click="deleteTeamBook(item.team_relationship_id)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.delete"}}</button>
                                </td>
                            </tr>
                            </template>
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
<div class="modal fade" id="addTeamBookDialogModal" tabindex="-1" role="dialog" aria-labelledby="addTeamBookDialogModalLabel">
    <div class="modal-dialog modal-sm" role="document" style="width: 450px;">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.TeamBookAdd"}}" id="addTeamBookDialogForm">
            <input type="hidden" name="teamId" value="{{.Model.TeamId}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.join_proj"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-3 control-label">{{i18n .Lang "mgr.proj_name"}}</label>
                        <div class="col-sm-9">
                            <select class="js-data-example-ajax form-control" multiple="multiple" name="bookId" id="bookId"></select>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnAddBook">{{i18n .Lang "common.save"}}</button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js" }}" type="text/javascript"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        var modalCache = $("#addTeamBookDialogModal form").html();

        $("#addTeamBookDialogForm").ajaxForm({
            beforeSubmit : function () {
                var memberId = $.trim($("#bookId").val());
                if(memberId === ""){
                    return showError({{i18n .Lang "message.proj_empty"}});
                }
                $("#btnAddBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    app.lists.splice(0,0,res.data);
                    $("#addTeamBookDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }
                $("#btnAddBook").button("reset");
            }
        });

        $("#addTeamBookDialogModal").on("hidden.bs.modal",function () {
            $(this).find("form").html(modalCache);
        }).on("show.bs.modal",function () {
            $('.js-data-example-ajax').select2({
                language: {{i18n .Lang "common.js_lang"}},
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "ManagerController.TeamSearchBook" "teamId" .Model.TeamId}}',
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

        var app = new Vue({
            el: "#teamBookList",
            data: {
                lists: {{.Result}}
            },
            delimiters: ['${', '}'],
            methods: {
                deleteTeamBook: function ($id, $e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.TeamBookDelete"}}",
                        data : { "teamRelId" : $id },
                        type : "post",
                        dataType : "json",
                        success : function ($res) {
                            if($res.errcode === 0){
                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.team_relationship_id == $id) {
                                        $this.lists.splice(index,1);
                                        break;
                                    }
                                }
                            }else {
                                layer.msg(res.message);
                            }
                        },
                        error : function () {
                            layer.msg({{i18n .Lang "message.system_error"}});
                        }
                    });
                }
            }
        });
    });
</script>
</body>
</html>