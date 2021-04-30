<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "uc.edit_user"}} - Powered by MinDoc</title>

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
        {{template "manager/widgets.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> {{i18n .Lang "uc.edit_user"}}</strong>
                    </div>
                </div>
                <div class="box-body col-lg-6 col-sm-12">
                    <form method="post" id="saveMemberInfoForm">
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.username"}}</label>
                            <input type="text" class="form-control" name="account" disabled placeholder="{{i18n .Lang "uc.username"}}" value="{{.Model.Account}}">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.realname"}}</label>
                            <input type="text" name="real_name" class="form-control" value="{{.Model.RealName}}" placeholder="{{i18n .Lang "uc.realname"}}">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.password"}}</label>
                            <input type="password" class="form-control" name="password1" placeholder="{{i18n .Lang "uc.password"}}" maxlength="50">
                            <p style="color: #999;font-size: 12px;">{{i18n .Lang "uc.pwd_tips"}}</p>
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.confirm_pwd"}}</label>
                            <input type="password" class="form-control" name="password2" placeholder="{{i18n .Lang "uc.confirm_pwd"}}" maxlength="50">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.email"}} <strong class="text-danger">*</strong></label>
                            <input type="email" class="form-control" name="email" placeholder="{{i18n .Lang "uc.email"}}" value="{{.Model.Email}}" maxlength="50">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "uc.mobile"}}</label>
                            <input type="text" class="form-control" name="phone" placeholder="{{i18n .Lang "uc.mobile"}}" maxlength="50" value="{{.Model.Phone}}">
                        </div>
                        <div class="form-group">
                            <label class="description">{{i18n .Lang "uc.description"}}</label>
                            <textarea class="form-control" rows="3" title="{{i18n .Lang "uc.description"}}" name="description" id="description" maxlength="500" >{{.Model.Description}}</textarea>
                            <p style="color: #999;font-size: 12px;">{{i18n .Lang "uc.description_tips"}}</p>
                        </div>
                        <div class="form-group">
                            <button type="submit" id="btnMemberInfo" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.save"}}</button>
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
                    return showError({{i18n .Lang "message.email_empty"}});
                }
                if (password1 !== "" && password1 !== password2){
                    return showError({{i18n .Lang "message.wrong_confirm_pwd"}});
                }
                $("#btnMemberInfo").button("loading");
            },success : function (res) {
                if(res.errcode === 0) {
                    showSuccess({{i18n .Lang "message.success"}})
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