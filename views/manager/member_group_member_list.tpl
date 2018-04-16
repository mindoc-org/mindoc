<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{.Model.GroupName}} - 用户组成员管理 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css"}}" rel="stylesheet">
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
            <div class="page-left">
                <ul class="menu">
                {{template "manager/manager_widgets.tpl.tpl" .}}
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{.Model.GroupName}} - 用户组成员管理</strong>
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
                                    <th>状态</th>
                                    <th>操作</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td>${item.group_member_id}</td>
                                    <td><img :src="item.avatar" onerror="this.src='{{cdnimg "/static/images/middle.gif"}}'" class="img-circle" width="34" height="34"></td>
                                    <td>${item.account}</td>
                                    <td>${item.real_name}</td>
                                    <td>${item.role_name}</td>
                                    <td>
                                        <template v-if="item.status == 0">
                                            <span class="label label-success">正常</span>
                                        </template>
                                        <template v-else>
                                            <span class="label label-danger">禁用</span>
                                        </template>
                                    </td>

                                    <td>
                                        <button type="button" class="btn btn-danger btn-sm" @click="deleteMember(item.group_member_id,$event)" data-loading-text="删除中">删除</button>
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
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.MemberGroupMemberEdit"}}" id="addMemberDialogForm">
            <input type="hidden" name="group_id" value="{{.Model.GroupId}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">添加用户组成员</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">账号<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <select class="js-data-example-ajax form-control" multiple="multiple" name="member_id" id="member_id"></select>
                        </div>
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
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#addMemberDialogModal").on("show.bs.modal",function () {
            window.addMemberDialogModalHtml = $(this).find("form").html();
            $('#member_id').select2({
                language: "zh-CN",
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "ManagerController.MemberGroupMemberSearch" "group_id" .Model.GroupId}}',
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
        }).on("hidden.bs.modal",function () {
            $(this).find("form").html(window.addMemberDialogModalHtml);
        });
        $("#addMemberDialogForm").ajaxForm({
            beforeSubmit : function () {
                var member_id = $.trim($("#member_id").val());
                if(member_id <= 0){
                    return showError("账号不能为空");
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

        window.app = new Vue({
            el : "#userList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
                deleteMember : function (id, e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.MemberGroupMemberDelete"}}",
                        type : "post",
                        data : { "id":id },
                        dataType : "json",
                        success : function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.group_member_id == id) {
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