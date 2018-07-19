<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑用户 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
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
                    <li><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li class="active"><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                    {{/*<li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>*/}}
                    <li><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.AttachList" }}" class="item"><i class="fa fa-cloud-upload" aria-hidden="true"></i> 附件管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.LabelList" }}" class="item"><i class="fa fa-bookmark" aria-hidden="true"></i> 标签管理</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 编辑用户</strong>
                    </div>
                </div>
                <div class="box-body col-lg-6 col-sm-12">
                    <form method="post" id="saveMemberInfoForm">
                        <div class="form-group">
                            <label>用户账号</label>
                            <input type="text" class="form-control" name="account" disabled placeholder="用户账号" value="{{.Model.Account}}">
                        </div>
                        <div class="form-group">
                            <label>真实姓名</label>
                            <input type="text" name="real_name" class="form-control" value="{{.Model.RealName}}" placeholder="真实姓名">
                        </div>
                        <div class="form-group">
                            <label>用户密码</label>
                            <input type="password" class="form-control" name="password1" placeholder="用户密码" maxlength="50">
                            <p style="color: #999;font-size: 12px;">不修改密码请留空,只支持本地用户修改密码</p>
                        </div>
                        <div class="form-group">
                            <label>确认密码</label>
                            <input type="password" class="form-control" name="password2" placeholder="确认密码" maxlength="50">
                        </div>
                        <div class="form-group">
                            <label>用户邮箱 <strong class="text-danger">*</strong></label>
                            <input type="email" class="form-control" name="email" placeholder="用户邮箱" value="{{.Model.Email}}" maxlength="50">
                        </div>
                        <div class="form-group">
                            <label>手机号码</label>
                            <input type="text" class="form-control" name="phone" placeholder="手机号码" maxlength="50" value="{{.Model.Phone}}">
                        </div>
                        <div class="form-group">
                            <label class="description">描述</label>
                            <textarea class="form-control" rows="3" title="描述" name="description" id="description" maxlength="500" >{{.Model.Description}}</textarea>
                            <p style="color: #999;font-size: 12px;">描述不能超过500字</p>
                        </div>
                        <div class="form-group">
                            <button type="submit" id="btnMemberInfo" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                    </form>

                    <div class="clearfix"></div>

                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>


<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#saveMemberInfoForm").ajaxForm({
            beforeSubmit : function () {
                var $then = $("#saveMemberInfoForm");

                var email = $.trim($then.find("input[name='email']").val());
                var password1 = $.trim($then.find("input[name='password1']").val());
                var password2 = $.trim($then.find("input[name='password2']").val());
                if (email === ""){
                    return showError("用户邮箱不能为空!");
                }
                if (password1 !== "" && password1 !== password2){
                    return showError("确认密码不正确!");
                }
                $("#btnMemberInfo").button("loading");
            },success : function (res) {
                if(res.errcode === 0) {
                    showSuccess("保存成功")
                }else{
                    showError(res.message);
                }
                $("#btnMemberInfo").button("reset");
            }
        });
    });
</script>
</body>
</html>