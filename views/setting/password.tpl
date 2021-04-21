<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "uc.user_center"}} - Powered by MinDoc</title>

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
</head>
<body>
<div class="manual-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                    <li><a href="{{urlfor "SettingController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> {{i18n .Lang "uc.base_info"}}</a> </li>
                    <li class="active"><a href="{{urlfor "SettingController.Password"}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> {{i18n .Lang "uc.change_pwd"}}</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{i18n .Lang "uc.change_pwd"}}</strong>
                    </div>
                </div>
                <div class="box-body" style="width: 300px;">
                    <form role="form" method="post" id="securityForm">
                        <div class="form-group">
                            <label for="password1">{{i18n .Lang "uc.origin_pwd"}}</label>
                            <input type="password" name="password1" id="password1" class="form-control disabled" placeholder="{{i18n .Lang "uc.origin_pwd"}}">
                        </div>
                        <div class="form-group">
                            <label for="password2">{{i18n .Lang "uc.new_pwd"}}</label>
                            <input type="password" class="form-control" name="password2" id="password2" max="50" placeholder="{{i18n .Lang "uc.new_pwd"}}">
                        </div>
                        <div class="form-group">
                            <label for="password3">{{i18n .Lang "uc.confirm_pwd"}}</label>
                            <input type="password" class="form-control" id="password3" name="password3" placeholder="{{i18n .Lang "uc.confirm_pwd"}}">
                        </div>
                        <div class="form-group">
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                        <div class="form-group">
                            <button type="submit" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "message.save"}}</button>
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
                    showError({{i18n .Lang "message.origin_pwd_empty"}});
                    return false;
                }
                if(!newPasswd){
                    showError({{i18n .Lang "message.new_pwd_empty"}});
                    return false;
                }
                if(!confirmPassword){
                    showError({{i18n .Lang "message.confirm_pwd_empty"}});
                    return false;
                }
                if(confirmPassword !== newPasswd){
                    showError({{i18n .Lang "message.wrong_confirm_pwd"}});
                    return false;
                }
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess({{i18n .Lang "message.success"}});
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