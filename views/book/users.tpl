<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>成员 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/font-awesome/css/font-awesome.min.css" rel="stylesheet">

    <link href="/static/css/main.css" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
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
                    <li><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 成员管理</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addBookMemberDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i> 添加成员</button>
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="userList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>
                            <div class="list-item" v-for="item in lists">
                                <img :src="item.avatar" onerror="this.src='/static/images/middle.gif'" class="img-circle" width="34" height="34">
                                <span>${item.account}</span>
                                <div class="operate">
                                    <template v-if="item.role_id == 0">
                                        创始人
                                    </template>
                                   <template v-else-if="item.role_id == 1">
                                       {{if eq .Member.Role 0}}
                                       <div class="btn-group">
                                           <button type="button" class="btn btn-default btn-sm"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">编辑者 <span class="caret"></span></button>
                                           <ul class="dropdown-menu">
                                               <li><a href="#">管理员</a> </li>
                                               <li><a href="#">编辑者</a> </li>
                                               <li><a href="#">观察者</a> </li>
                                           </ul>
                                       </div>
                                       <a href="#" class="btn btn-danger btn-sm">移除</a>
                                       {{end}}
                                   </template>
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="addBookMemberDialogModal" tabindex="-1" role="dialog" aria-labelledby="addBookMemberDialogModalLabel">
    <div class="modal-dialog modal-sm" role="document" style="width: 400px;">
        <form method="post" autocomplete="off" class="form-horizontal" action="{{urlfor "BookController.AddMember"}}" id="addBookMemberDialogForm">
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
                           <input type="text" name="account" class="form-control" placeholder="用户账号" id="account" maxlength="50">
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
                    <button type="submit" class="btn btn-success">保存</button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->
<script src="/static/jquery/1.12.4/jquery.min.js"></script>
<script src="/static/bootstrap/js/bootstrap.min.js"></script>
<script src="/static/vuejs/vue.min.js"></script>
<script src="/static/js/jquery.form.js" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#addBookMemberDialogForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#account").val());
                if(account === ""){
                    return showError("账号不能为空");
                }
            },
            success : function (res) {

            }
        });

        new Vue({
            el : "#userList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
            }
        });
        Vue.nextTick(function () {
            $("[data-toggle='tooltip']").tooltip();
        });
    });
</script>
</body>
</html>