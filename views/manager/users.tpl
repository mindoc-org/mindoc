<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户管理 - Powered by MinDoc</title>

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
        {{template "manager/widgets.tpl" "users"}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 成员管理</strong>
                        {{if eq .Member.Role 0}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addMemberDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> 添加成员</button>
                        {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="userList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th width="80">ID</th>
                                    <th width="80">头像</th>
                                    <th>账号</th>
                                    <th>姓名</th>
                                    <th>角色</th>
                                    <th>类型</th>
                                    <th>状态</th>
                                    <th>操作</th>
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
                                            超级管理员
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
                                                <li><a href="javascript:;" @click="setMemberRole(item.member_id,1)">管理员</a> </li>
                                                <li><a href="javascript:;" @click="setMemberRole(item.member_id,2)">普通用户</a> </li>
                                            </ul>
                                            </div>
                                        </template>
                                    </td>
                                    <td>
                                        ${item.auth_method}
                                    </td>
                                    <td>
                                        <template v-if="item.status == 0">
                                            <span class="label label-success">正常</span>
                                        </template>
                                        <template v-else>
                                            <span class="label label-danger">禁用</span>
                                        </template>
                                    </td>

                                    <td>
                                        <template v-if="item.member_id == {{.Member.MemberId}}">

                                        </template>
                                        <template v-else-if="item.role != 0">
                                            <a :href="'{{urlfor "ManagerController.EditMember" ":id" ""}}' + item.member_id" class="btn btn-sm btn-default" @click="editMember(item.member_id)">
                                                编辑
                                            </a>
                                            <template v-if="item.status == 0">
                                                <button type="button" class="btn btn-danger btn-sm" @click="setMemberStatus(item.member_id,1,$event)" data-loading-text="启用中...">禁用</button>
                                            </template>
                                            <template v-else>
                                                <button type="button" class="btn btn-success btn-sm" @click="setMemberStatus(item.member_id,0,$event)" data-loading-text="禁用中...">启用</button>
                                            </template>
                                            <button type="button" class="btn btn-danger btn-sm" @click="deleteMember(item.member_id,$event)" data-loading-text="删除中">删除</button>
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
                    <h4 class="modal-title" id="myModalLabel">创建用户</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">账号<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="text" name="account" class="form-control" placeholder="用户账号" id="account" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="password1">密码<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" placeholder="用户密码" name="password1" id="password1" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="password2">确认密码<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" placeholder="确认密码" name="password2" id="password2" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="email">邮箱<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <input type="email" class="form-control" placeholder="邮箱" name="email" id="email" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">真实姓名</label>
                        <div class="col-sm-10">
                            <input type="text" name="real_name" class="form-control" value="" placeholder="真实姓名">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">手机号</label>
                        <div class="col-sm-10">
                            <input type="text" class="form-control" placeholder="手机号" name="phone" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">角色</label>
                        <div class="col-sm-10">
                            <select name="role" class="form-control">
                                <option value="1">管理员</option>
                                <option value="2">普通用户</option>
                            </select>
                        </div>
                    </div>
                    <div class="form-group">

                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-success" data-loading-text="保存中..." id="btnAddMember">保存</button>
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
                    return showError("账号不能为空");
                }
                var password1 = $.trim($("#password1").val());
                var password2 = $("#password2").val();
                if (password1 === "") {
                    return showError("密码不能为空");
                }
                if (password1 !== password2) {
                    return showError("确认密码不正确");
                }
                var email = $.trim($("#email").val());

                if (email === "") {
                    return showError("邮箱不能为空");
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
                showError("服务器异常");
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
                                alert("操作失败：" + res.message);
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
                                alert("操作失败：" + res.message);
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