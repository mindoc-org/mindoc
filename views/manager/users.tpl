<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.user_mgr"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
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
                        <strong class="box-title"> {{i18n .Lang "mgr.member_mgr"}}</strong>
                        {{if eq .Member.Role 0}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addMemberDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> {{i18n .Lang "mgr.add_member"}}</button>
                        {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="userList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">{{i18n .Lang "message.no_data"}}</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th width="80">ID</th>
                                    <th width="80">{{i18n .Lang "uc.avatar"}}</th>
                                    <th>{{i18n .Lang "uc.username"}}</th>
                                    <th>{{i18n .Lang "uc.realname"}}</th>
                                    <th>{{i18n .Lang "uc.role"}}</th>
                                    <th>{{i18n .Lang "uc.type"}}</th>
                                    <th>{{i18n .Lang "uc.status"}}</th>
                                    <th>{{i18n .Lang "common.operate"}}</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td>${item.member_id}</td>
                                    <td><img :src="item.avatar" onerror="this.src='{{cdnimg "/static/images/middle.gif"}}'" class="img-circle" width="34" height="34"></td>
                                    <td>${item.account}</td>
                                    <td>${item.real_name}</td>
                                    <td>
                                        <template v-if="item.role == 0">
                                            {{i18n $.Lang "uc.super_admin"}}
                                        </template>
                                        <template v-else-if="item.member_id == {{.Member.MemberId}}">
                                            ${item.role_name}
                                        </template>
                                        <template v-else>
                                            <div class="btn-group">
                                            <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                                ${item.role_name}
                                                <span class="caret"></span></button>
                                            <ul class="dropdown-menu">
                                                <li><a href="javascript:;" @click="setMemberRole(item.member_id,1)">{{i18n $.Lang "uc.admin"}}</a> </li>
                                                <li><a href="javascript:;" @click="setMemberRole(item.member_id,2)">{{i18n $.Lang "uc.user"}}</a> </li>
                                            </ul>
                                            </div>
                                        </template>
                                    </td>
                                    <td>
                                        ${item.auth_method}
                                    </td>
                                    <td>
                                        <template v-if="item.status == 0">
                                            <span class="label label-success">{{i18n .Lang "uc.normal"}}</span>
                                        </template>
                                        <template v-else>
                                            <span class="label label-danger">{{i18n .Lang "uc.disable"}}</span>
                                        </template>
                                    </td>

                                    <td>
                                        <template v-if="item.member_id == {{.Member.MemberId}}">

                                        </template>
                                        <template v-else-if="item.role != 0">
                                            <a :href="'{{urlfor "ManagerController.EditMember" ":id" ""}}' + item.member_id" class="btn btn-sm btn-default" @click="editMember(item.member_id)">
                                                {{i18n .Lang "common.edit"}}
                                            </a>
                                            <template v-if="item.status == 0">
                                                <button type="button" class="btn btn-danger btn-sm" @click="setMemberStatus(item.member_id,1,$event)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "uc.disable"}}</button>
                                            </template>
                                            <template v-else>
                                                <button type="button" class="btn btn-success btn-sm" @click="setMemberStatus(item.member_id,0,$event)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "uc.enable"}}</button>
                                            </template>
                                            <button type="button" class="btn btn-danger btn-sm" @click="deleteMember(item.member_id,$event)" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.delete"}}</button>
                                        </template>
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
<div class="modal fade" id="addMemberDialogModal" tabindex="-1" role="dialog" aria-labelledby="addMemberDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.CreateMember"}}" id="addMemberDialogForm">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">{{i18n .Lang "uc.create_user"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">{{i18n .Lang "uc.username"}}<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="account" class="form-control" placeholder="{{i18n .Lang "uc.username"}}" id="account" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="password1">{{i18n .Lang "uc.password"}}<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" placeholder="{{i18n .Lang "uc.password"}}" name="password1" id="password1" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="password2">{{i18n .Lang "uc.confirm_pwd"}}<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" placeholder="{{i18n .Lang "uc.confirm_pwd"}}" name="password2" id="password2" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="email">{{i18n .Lang "uc.email"}}<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="email" class="form-control" placeholder="{{i18n .Lang "uc.email"}}" name="email" id="email" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "uc.realname"}}</label>
                        <div class="col-sm-10">
                            <input type="text" name="real_name" class="form-control" value="" placeholder="{{i18n .Lang "uc.realname"}}">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "uc.mobile"}}</label>
                        <div class="col-sm-10">
                            <input type="text" class="form-control" placeholder="{{i18n .Lang "uc.mobile"}}" name="phone" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">{{i18n .Lang "uc.role"}}</label>
                        <div class="col-sm-10">
                            <select name="role" class="form-control">
                                <option value="1">{{i18n .Lang "uc.admin"}}</option>
                                <option value="2">{{i18n .Lang "uc.user"}}</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">

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
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#addMemberDialogModal").on("show.bs.modal",function () {
            window.addMemberDialogModalHtml = $(this).find("form").html();
        }).on("hidden.bs.modal",function () {
            $(this).find("form").html(window.addMemberDialogModalHtml);
        });
        $("#addMemberDialogForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#account").val());
                if(account === ""){
                    return showError({{i18n .Lang "message.account_empty"}});
                }
                var password1 = $.trim($("#password1").val());
                var password2 = $("#password2").val();
                if (password1 === "") {
                    return showError({{i18n .Lang "message.password_empty"}});
                }
                if (password1 !== password2) {
                    return showError({{i18n .Lang "message.confirm_pwd_empty"}});
                }
                var email = $.trim($("#email").val());

                if (email === "") {
                    return showError({{i18n .Lang "message.email_empty"}});
                }
                $("#btnAddMember").button("loading");
                return true;
            },
            success : function (res) {
                if(res.errcode === 0){
                    app.lists.splice(0,0,res.data);
                    $("#addMemberDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }
                $("#btnAddMember").button("reset");
            },
            error : function () {
                showError({{i18n .Lang "message.system_error"}});
                $("#btnAddMember").button("reset");
            }
        });

        var app = new Vue({
            el : "#userList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
                setMemberStatus : function (id,status,e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.UpdateMemberStatus"}}",
                        type : "post",
                        data : { "member_id":id,"status" : status},
                        dataType : "json",
                        success : function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];

                                    if (item.member_id === id) {
                                        console.log(item);
                                        $this.lists[index].status = status;
                                        break;
                                        //$this.lists.splice(index,1,item);
                                    }
                                }
                            } else {
                                alert("{{i18n .Lang "message.failed"}}：" + res.message);
                            }
                        }
                    })

                },
                setMemberRole : function (member_id, role) {
                    var $this = this;
                    $.ajax({
                        url :"{{urlfor "ManagerController.ChangeMemberRole"}}",
                        dataType :"json",
                        type :"post",
                        data : { "member_id" : member_id,"role" : role },
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
                                alert("{{i18n .Lang "message.failed"}}：" + res.message);
                            }
                        }
                    })
                },
                deleteMember : function (id, e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.DeleteMember"}}",
                        type : "post",
                        data : { "id":id },
                        dataType : "json",
                        success : function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.member_id == id) {
                                        console.log(item);
                                        $this.lists.splice(index,1);
                                        break;
                                    }
                                }
                            } else {
                                alert("{{i18n .Lang "message.failed"}}：" + res.message);
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