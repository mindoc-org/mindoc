<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>成员 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">

</head>
<body>
<div class="manual-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li class="active"><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    {{if eq .Model.RoleId 0 1}}
                    <li><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                    {{end}}
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 成员管理</strong>
                        {{if eq .Model.RoleId 0 1}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addBookMemberDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> 添加成员</button>
                        {{end}}
                    </div>
                </div>
                <div class="box-body" id="userList">
                    <div class="users-list">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>
                            <div class="list-item" v-for="item in lists">
                                <img :src="item.avatar" onerror="this.src='{{cdnimg "/static/images/middle.gif"}}'" class="img-circle" width="34" height="34">
                                <span>${item.account}</span>
                                <span style="font-size: 12px;color: #484848" v-if="item.real_name != ''">[${item.real_name}]</span>
                                <div class="operate">
                                    <template v-if="item.role_id == 0">
                                        创始人
                                    </template>
                                   <template v-else>
                                       <template v-if="(book.role_id == 1 || book.role_id == 0) && member.member_id != item.member_id">
                                           <div class="btn-group">
                                               <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                                   ${item.role_name}
                                                   <span class="caret"></span></button>
                                               <ul class="dropdown-menu">
                                                   <li><a href="javascript:;" @click="setBookMemberRole(item.member_id,1)">管理员</a> </li>
                                                   <li><a href="javascript:;" @click="setBookMemberRole(item.member_id,2)">编辑者</a> </li>
                                                   <li><a href="javascript:;" @click="setBookMemberRole(item.member_id,3)">观察者</a> </li>
                                               </ul>
                                           </div>
                                           <button type="button" class="btn btn-danger btn-sm" @click="removeBookMember(item.member_id)">移除</button>
                                       </template>
                                       <template v-else>
                                           <template v-if="item.role_id == 1">
                                               管理员
                                           </template>
                                           <template v-else-if="item.role_id == 2">
                                               编辑者
                                           </template>
                                           <template v-else-if="item.role_id == 3">
                                               观察者
                                           </template>
                                       </template>
                                   </template>
                                </div>
                            </div>
                        </template>
                    </div>
                    <template v-if="lists.length >= 0">
                       <nav class="pagination-container">
                        {{.PageHtml}}
                        </nav>
                    </template>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="addBookMemberDialogModal" tabindex="-1" role="dialog" aria-labelledby="addBookMemberDialogModalLabel">
    <div class="modal-dialog modal-sm" role="document" style="width: 400px;">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "BookMemberController.AddMember"}}" id="addBookMemberDialogForm">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">添加成员</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">账号</label>
                       <div class="col-sm-10">
                           {{/*<input type="text" name="account" class="form-control" placeholder="用户账号" id="account" maxlength="50">*/}}
                           <select class="js-data-example-ajax form-control" multiple="multiple" name="account" id="account"></select>
                       </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">角色</label>
                        <div class="col-sm-10">
                            <select name="role_id" class="form-control">
                                <option value="1">管理员</option>
                                <option value="2">编辑者</option>
                                <option value="3">观察者</option>
                            </select>
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

        var modalCache = $("#addBookMemberDialogModal form").html();

        /**
         * 添加用户
         */
        $("#addBookMemberDialogForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#account").val());
                if(account === ""){
                    return showError("账号不能为空");
                }
                $("#btnAddMember").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    app.lists.splice(0,0,res.data);
                    $("#addBookMemberDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }
                $("#btnAddMember").button("reset");
            }
        });
        $("#addBookMemberDialogModal").on("hidden.bs.modal",function () {
            $(this).find("form").html(modalCache);
        }).on("show.bs.modal",function () {
            $('.js-data-example-ajax').select2({
                language: "zh-CN",
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "SearchController.User" ":key" .Model.Identify}}',
                    dataType: 'json',
                    data: function (params) {
                        return {
                            q: params.term, // search term
                            page: params.page
                        };
                    },
                    processResults: function (data, params) {
                        console.log(data)
                        return {
                            results : data.data.results
                        }
                    }
                }
            });
        });

        var app = new Vue({
            el : "#userList",
            data : {
                lists : {{.Result}},
                member : {
                    member_role : {{.Member.Role}},
                    member_id : {{.Member.MemberId}}
                },
                book : {
                    role_id : {{.Model.RoleId}},
                    identify : {{.Model.Identify}}
                }
            },
            delimiters : ['${','}'],
            methods : {
                setBookMemberRole : function (member_id, role_id) {
                    var $this = this;
                    $.ajax({
                       url : "{{urlfor "BookMemberController.ChangeRole"}}",
                        data : { "identify" : $this.book.identify,"member_id" : member_id,"role_id" : role_id },
                        type :"post",
                        dataType : "json",
                        success : function (res) {
                            console.log(res);
                            if (res.errcode === 0){
                                for(var index in $this.lists){
                                    var item = $this.lists[index];
                                    if (item.member_id === member_id){
                                        $this.lists.splice(index,1,res.data);
                                    }
                                }
                            }else{
                                alert(res.message);
                            }
                        }
                    });
                },
                removeBookMember : function (member_id) {
                    var $this = this;
                    $.ajax({
                        url : "{{urlfor "BookMemberController.RemoveMember"}}",
                        type :"post",
                        dataType :"json",
                        data :{ "identify" : $this.book.identify,"member_id" : member_id},
                        success : function (res) {
                            if(res.errcode === 0){
                                for(var index in $this.lists){
                                    if($this.lists[index].member_id === member_id){
                                        $this.lists.splice(index,1);
                                        break;
                                    }
                                }
                            }else{
                                alert(res.message);
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