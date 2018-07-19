<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户中心 - Powered by MinDoc</title>

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
                    <li><a href="{{urlfor "SettingController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 基本信息</a> </li>
                    <li class="active"><a href="{{urlfor "SettingController.Password"}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> 修改密码</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">修改密码</strong>
                    </div>
                </div>
                <div class="box-body" style="width: 300px;">
                    <form role="form" method="post" id="securityForm">
                        <div class="form-group">
                            <label for="password1">原始密码</label>
                            <input type="password" name="password1" id="password1" class="form-control disabled" placeholder="原始密码">
                        </div>
                        <div class="form-group">
                            <label for="password2">新密码</label>
                            <input type="password" class="form-control" name="password2" id="password2" max="50" placeholder="新密码">
                        </div>
                        <div class="form-group">
                            <label for="password3">确认密码</label>
                            <input type="password" class="form-control" id="password3" name="password3" placeholder="确认密码">
                        </div>
                        <div class="form-group">
                            <button type="submit" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                    </form>
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

        $("#securityForm").ajaxForm({
            beforeSubmit : function () {
                var oldPasswd = $("#password1").val();
                var newPasswd = $("#password2").val();
                var confirmPassword = $("#password3").val();
                if(!oldPasswd ){
                    showError("原始密码不能为空");
                    return false;
                }
                if(!newPasswd){
                    showError("新密码不能为空");
                    return false;
                }
                if(!confirmPassword){
                    showError("确认密码不能为空");
                    return false;
                }
                if(confirmPassword !== newPasswd){
                    showError("确认密码不正确");
                    return false;
                }
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess('保存成功');
                    $("#password1").val('');
                    $("#password2").val('');
                    $("#password3").val('');
                }else{
                    showError(res.message);
                }
            }
        }) ;
    });
</script>
</body>
</html>