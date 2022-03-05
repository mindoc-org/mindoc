<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.team_mgr"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">

    <style type="text/css">
        .table > tbody > tr > td {
            vertical-align: middle;
        }
    </style>
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
                        <strong class="box-title"> {{i18n .Lang "mgr.team_mgr"}}</strong>
                    {{if eq .Member.Role 0}}
                        <button type="button" class="btn btn-success btn-sm pull-right" data-toggle="modal"
                                data-target="#addTeamDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i>
                            {{i18n .Lang "mgr.create_team"}}
                        </button>
                    {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="teamList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">{{i18n .Lang "message.no_data"}}</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th>{{i18n .Lang "mgr.team_name"}}</th>
                                    <th width="150px">{{i18n .Lang "mgr.member_amount"}}</th>
                                    <th width="150px">{{i18n .Lang "mgr.proj_amount"}}</th>
                                    <th align="center" width="260px">{{i18n .Lang "common.operate"}}</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td>${item.team_name}</td>
                                    <td>${item.member_count}</td>
                                    <td>${item.book_count}</td>
                                    <td>
                                        <a :href="'{{urlfor "ManagerController.TeamBookList" ":id" ""}}' + item.team_id" class="btn btn-primary btn-sm">{{i18n .Lang "mgr.proj"}}</a>
                                        <a :href="'{{urlfor "ManagerController.TeamMemberList" ":id" ""}}' + item.team_id" type="button" class="btn btn-success btn-sm">{{i18n .Lang "mgr.member"}}</a>
                                        <button type="button" class="btn btn-sm btn-default" @click="editTeam(item.team_id)">{{i18n .Lang "common.edit"}}</button>
                                        <button type="button" class="btn btn-danger btn-sm" @click="deleteTeam(item.team_id,$event)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.delete"}}</button>
                                    </td>
                                </tr>
                                </tbody>
                            </table>
                        </template>
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
<div class="modal fade" id="addTeamDialogModal" tabindex="-1" role="dialog" aria-labelledby="addTeamDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" class="form-horizontal"
              action="{{urlfor "ManagerController.TeamCreate"}}" id="addTeamDialogForm">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
                            aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.create_team"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">{{i18n .Lang "mgr.team_name"}}<span
                                class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="teamName" class="form-control" placeholder="{{i18n .Lang "mgr.team_name"}}" id="teamName" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnAddTeam">{{i18n .Lang "common.save"}}
                    </button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->
<div class="modal fade" id="editTeamDialogModal" tabindex="-1" role="dialog" aria-labelledby="editTeamDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.TeamEdit"}}" id="editTeamDialogForm">
            <input type="hidden" name="teamId" value="">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
                            aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.edit_team"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">{{i18n .Lang "mgr.team_name"}}<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="teamName" class="form-control" placeholder="{{i18n .Lang "mgr.team_name"}}" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnEditTeam">{{i18n .Lang "common.save"}}
                    </button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->


<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>


<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        var editTeamDialogModal = $("#editTeamDialogModal");

        $("#addTeamDialogModal").on("show.bs.modal", function () {
            window.addTeamDialogModalHtml = $(this).find("form").html();
        }).on("hidden.bs.modal", function () {
            $(this).find("form").html(window.addTeamDialogModalHtml);
        });
        $("#addTeamDialogForm").ajaxForm({
            beforeSubmit: function () {
                var account = $.trim($("#addTeamDialogForm #teamName").val());
                if (account === "") {
                    return showError({{i18n .Lang "message.team_name_empty"}});
                }
                $("#btnAddTeam").button("loading");
                return true;
            },
            success: function ($res) {
                if ($res.errcode === 0) {
                    app.lists.splice(0, 0, $res.data);
                    $("#addTeamDialogModal").modal("hide");
                } else {
                    showError($res.message);
                }
            },
            error: function () {
                showError({{i18n .Lang "message.system_error"}});
            },
            complete: function () {
                $("#btnAddTeam").button("reset");
            }
        });

        editTeamDialogModal.on("shown.bs.modal",function () {
           $(this).find("input[name='teamName']").focus();
        });
        $("#editTeamDialogForm").ajaxForm({
            beforeSubmit: function () {
                var account = $.trim(editTeamDialogModal.find("input[name='teamName']").val());
                if (account === "") {
                    return showError({{i18n .Lang "message.team_name_empty"}});
                }
                $("#btnEditTeam").button("loading");
                return true;
            },success :function ($res) {
                if ($res.errcode === 0) {
                    for (var index in app.lists) {
                        var item = app.lists[index];
                        if (item.team_id == $res.data.team_id) {
                            app.lists.splice(index, 1, $res.data);
                            break;
                        }
                    }
                    editTeamDialogModal.modal("hide");
                }else {
                    showError($res.message);
                }
            },
            error: function () {
                showError({{i18n .Lang "message.system_error"}});
            },
            complete: function () {
                $("#btnEditTeam").button("reset");
            }
        });

        var app = new Vue({
            el: "#teamList",
            data: {
                lists: {{.Result}}
            },
            delimiters: ['${', '}'],
            methods: {
                deleteTeam: function (id, e) {
                    var $this = this;
                    $.ajax({
                        url: "{{urlfor "ManagerController.TeamDelete"}}",
                        type: "post",
                        data: {"teamId": id},
                        dataType: "json",
                        success: function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.team_id == id) {
                                        $this.lists.splice(index, 1);
                                        break;
                                    }
                                }
                            } else {
                                alert("操作失败：" + res.message);
                            }
                        }
                    });
                },
                editTeam : function (id, e) {
                    var $this = this;
                    for (var index in $this.lists) {
                        var item = $this.lists[index];
                        if (item.team_id == id) {
                            editTeamDialogModal.find("input[name='teamName']").val(item.team_name);
                            editTeamDialogModal.find("input[name='teamId']").val(item.team_id);
                            editTeamDialogModal.modal("show");
                            break;
                        }
                    }
                }
            }
        });
        Vue.nextTick(function () {
            $("[data-toggle='tooltip']").tooltip();
        });
    });
</script>
</body>
</html>