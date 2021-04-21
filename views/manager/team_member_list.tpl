<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.team_member_mgr"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
    <style type="text/css">
        .table>tbody>tr>td{vertical-align: middle;}
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
                        <strong class="box-title">{{.Model.TeamName}} - {{i18n .Lang "mgr.member_mgr"}}</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addTeamMemberDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> {{i18n .Lang "mgr.add_member"}}</button>
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="teamMemberList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">{{i18n .Lang "message.no_data"}}</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th width="80">{{i18n .Lang "uc.avatar"}}</th>
                                    <th width="100">{{i18n .Lang "uc.username"}}</th>
                                    <th width="100">{{i18n .Lang "uc.realname"}}</th>
                                    <th width="150">{{i18n .Lang "uc.role"}}</th>
                                    <th width="80px">{{i18n .Lang "common.operate"}}</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td><img :src="item.avatar" onerror="this.src='{{cdnimg "/static/images/middle.gif"}}'" class="img-circle" width="34" height="34"></td>
                                    <td>${item.account}</td>
                                    <td>${item.real_name}</td>
                                    <td>
                                        <div class="btn-group">
                                            <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                                {{i18n .Lang "uc.role"}}：${item.role_name}
                                                <span class="caret"></span></button>
                                            <ul class="dropdown-menu">
                                                <li><a href="javascript:;" @click="setTeamMemberRole(item.member_id,1)">{{i18n .Lang "common.administrator"}} <p class="text">{{i18n .Lang "common.admin_right"}}</p></a> </li>
                                                <li><a href="javascript:;" @click="setTeamMemberRole(item.member_id,2)">{{i18n .Lang "common.editor"}} <p class="text">{{i18n .Lang "common.editor_right"}}</p></a> </li>
                                                <li><a href="javascript:;" @click="setTeamMemberRole(item.member_id,3)">{{i18n .Lang "common.observer"}} <p class="text">{{i18n .Lang "common.observer_right"}}</p></a> </li>
                                            </ul>
                                        </div>
                                    </td>
                                    <td>
                                        <button type="button" class="btn btn-danger btn-sm" @click="deleteMember(item.member_id,$event)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.delete"}}</button>
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
<div class="modal fade" id="addTeamMemberDialogModal" tabindex="-1" role="dialog" aria-labelledby="addTeamMemberDialogModalLabel">
    <div class="modal-dialog modal-sm" role="document" style="width: 400px;">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.TeamMemberAdd"}}" id="addTeamMemberDialogForm">
            <input type="hidden" name="teamId" value="{{.Model.TeamId}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "mgr.add_member"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "uc.username"}}</label>
                        <div class="col-sm-10">
                            <select class="js-data-example-ajax form-control" multiple="multiple" name="memberId" id="memberId"></select>
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "uc.role"}}</label>
                        <div class="col-sm-10">
                            <select name="roleId" class="form-control">
                                <option value="1">{{i18n .Lang "common.administrator"}}</option>
                                <option value="2">{{i18n .Lang "common.editor"}}</option>
                                <option value="3">{{i18n .Lang "common.observer"}}</option>
                            </select>
                        </div>
                    </div>

                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}" id="btnAddMember">{{i18n .Lang "common.save"}}</button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        var modalCache = $("#addTeamMemberDialogModal form").html();

        /**
         * 添加用户
         */
        $("#addTeamMemberDialogForm").ajaxForm({
            beforeSubmit : function () {
                var memberId = $.trim($("#memberId").val());
                if(memberId === ""){
                    return showError({{i18n .Lang "message.account_empty"}});
                }
                $("#btnAddMember").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    app.lists.splice(0,0,res.data);
                    $("#addTeamMemberDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }
                $("#btnAddMember").button("reset");
            }
        });
        $("#addTeamMemberDialogModal").on("hidden.bs.modal",function () {
            $(this).find("form").html(modalCache);
        }).on("show.bs.modal",function () {
            $('.js-data-example-ajax').select2({
                language: {{i18n .Lang "common.js_lang"}},
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "ManagerController.TeamSearchMember" "teamId" .Model.TeamId}}',
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
            el : "#teamMemberList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
                setTeamMemberRole : function (member_id, role) {
                    var $this = this;
                    $.ajax({
                        url :"{{urlfor "ManagerController.TeamChangeMemberRole"}}",
                        dataType :"json",
                        type :"post",
                        data : { "memberId" : member_id,"roleId" : role ,"teamId":{{.Model.TeamId}}},
                        success : function (res) {
                            if(res.errcode === 0){
                                for (var index in $this.lists) {
                                    var item = $this.lists[index];

                                    if (item.member_id === member_id) {

                                        $this.lists.splice(index,1,res.data);
                                        break;
                                    }
                                }
                            }else{
                                alert("操作失败：" + res.message);
                            }
                        }
                    })
                },
                deleteMember : function (id, e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.TeamMemberDelete"}}",
                        type : "post",
                        data : { "memberId":id ,"teamId":{{.Model.TeamId}}},
                        dataType : "json",
                        success : function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.member_id == id) {
                                        $this.lists.splice(index,1);
                                        break;
                                    }
                                }
                            } else {
                                alert("操作失败：" + res.message);
                            }
                        }
                    });
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