<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户组管理 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

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
                        <strong class="box-title"> 用户组管理</strong>
                    {{if eq .Member.Role 0}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addMemberGroupDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> 添加用户组</button>
                    {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="memberGroupList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th width="80">ID</th>
                                    <th>用户组名称</th>
                                    <th>成员数量</th>
                                    <th>创建时间</th>
                                    <th>创建人</th>
                                    <th>修改时间</th>
                                    <th>操作</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td>${item.group_id}</td>
                                    <td>${item.group_name}</td>
                                    <td>${item.group_number}</td>
                                    <td>${(new Date(item.create_time)).format("yyyy-MM-dd hh:mm:ss")}</td>
                                    <td>${item.create_name}</td>
                                    <td>${(new Date(item.modify_time)).format("yyyy-MM-dd hh:mm:ss")}</td>
                                    <td>
                                        <a :href="'{{urlfor "ManagerController.MemberGroupMemberList" ":id" ""}}' + item.group_id" class="btn btn-sm btn-success">成员</a>
                                         <a :href="'{{urlfor "ManagerController.MemberGroupEdit" ":id" ""}}' + item.group_id" class="btn btn-sm btn-default">编辑</a>
                                         <button type="button" class="btn btn-danger btn-sm" @click="deleteMemberGroup(item.group_id,$event)" data-loading-text="删除中">删除</button>
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
<div class="modal fade" id="addMemberGroupDialogModal" tabindex="-1" role="dialog" aria-labelledby="addMemberGroupDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "ManagerController.MemberGroupEdit"}}" id="addMemberGroupDialogForm">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">创建用户组</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-3 control-label" for="group_name">用户组名称<span class="error-message">*</span></label>
                        <div class="col-sm-9">
                            <input type="text" name="group_name" class="form-control" placeholder="用户组名称" id="group_name" maxlength="50">
                        </div>
                    </div>
                    <div class="form-group"></div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-success" data-loading-text="保存中..." id="btnAddMemberGroup">保存</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!--END Modal-->

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#addMemberGroupDialogModal").on("show.bs.modal",function () {
            window.addMemberGroupDialogModalHtml = $(this).find("form").html();
        }).on("hidden.bs.modal",function () {
            $(this).find("form").html(window.addMemberDialogModalHtml);
        });
        $("#addMemberGroupDialogForm").ajaxForm({
            beforeSubmit : function () {
                var group_name = $.trim($("#group_name").val());
                if(group_name === ""){
                    return showError("用户组名称不能为空");
                }
                $("#btnAddMemberGroup").button("loading");
                return true;
            },
            success : function (res) {
                if(res.errcode === 0){
                    app.lists.splice(0,0,res.data);
                    $("#addMemberGroupDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }
                $("#btnAddMemberGroup").button("reset");
            },
            error : function () {
                showError("服务器异常");
                $("#btnAddMemberGroup").button("reset");
            }
        });

        window.app = new Vue({
            el : "#memberGroupList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
                deleteMemberGroup : function (id,status,e) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "ManagerController.MemberGroupDelete"}}",
                        type : "post",
                        data : { "group_id":id },
                        dataType : "json",
                        success : function (res) {
                            if (res.errcode === 0) {

                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.group_id == id) {
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