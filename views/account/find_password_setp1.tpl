<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="utf-8">
    <link rel="shortcut icon" href="{{cdnimg "/static/favicon.ico"}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="MinDoc" />
    <title>{{i18n .Lang "common.account_recovery"}} - Powered by SSHotRiver</title>

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
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
</head>
<body class="manual-container">
<header class="navbar navbar-static-top smart-nav navbar-fixed-top manual-header" role="banner">
    <div class="container">
        <div class="navbar-header col-sm-12 col-md-6 col-lg-5">
            <a href="{{.BaseUrl}}" class="navbar-brand">{{.SITE_NAME}}</a>
        </div>
    </div>
</header>
<div class="container manual-body">
    <div class="row login">
        <div class="login-body">
            <form role="form" method="post" id="findPasswordForm">
            {{ .xsrfdata }}
                <h3 class="text-center">{{i18n .Lang "common.account_recovery"}}</h3>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-at"></i>
                        </div>
                        <input type="text" class="form-control" placeholder="{{i18n .Lang "common.email"}}" name="email" id="email" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group" style="float: left;width: 195px;">
                        <div class="input-group-addon">
                            <i class="fa fa-check-square"></i>
                        </div>
                        <input type="text" name="code" id="code" class="form-control" style="width: 150px" maxlength="5" placeholder="{{i18n .Lang "common.captcha"}}" autocomplete="off">&nbsp;
                    </div>
                    <img id="captcha-img" style="width: 140px;height: 40px;display: inline-block;float: right" src="{{urlfor "AccountController.Captcha"}}" onclick="this.src='{{urlfor "AccountController.Captcha"}}?key=login&t='+(new Date()).getTime();" title="{{i18n .Lang "message.click_to_change"}}">
                    <div class="clearfix"></div>
                </div>

                <div class="form-group">
                    <button type="submit" id="btnSendMail" class="btn btn-success" style="width: 100%"  data-loading-text="{{i18n .Lang "message.processing"}}" autocomplete="off">{{i18n .Lang "common.account_recovery"}}</button>
                </div>

            </form>
        </div>
    </div>
    <div class="clearfix"></div>
</div>
{{template "widgets/footer.tpl" .}}
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#email,#code").on('focus',function () {
            $(this).tooltip('destroy').parents('.form-group').removeClass('has-error');;
        });

        $(document).keydown(function (e) {
            var event = document.all ? window.event : e;
            if(event.keyCode == 13){
                $("#btn-login").click();
            }
        });

        $("#findPasswordForm").ajaxForm({
            beforeSubmit : function () {
                var $btn = $(this).button('loading');

                var email = $.trim($("#email").val());
                if(email === ""){
                    $("#email").tooltip({placement:"auto",title : "{{i18n .Lang "message.email_empty"}}",trigger : 'manual'})
                        .tooltip('show')
                        .parents('.form-group').addClass('has-error');
                    $btn.button('reset');
                    return false;

                }
                var code = $.trim($("#code").val());
                if(code === ""){
                    $("#code").tooltip({title : '{{i18n .Lang "message.captcha_empty"}}',trigger : 'manual'})
                        .tooltip('show')
                        .parents('.form-group').addClass('has-error');
                    $btn.button('reset');
                    return false;
                }
                $("#btnSendMail").button("loading");
            },
            success : function (res) {

                if(res.errcode !== 0){
                    $("#captcha-img").click();
                    $("#code").val('');
                    layer.msg(res.message);
                    $("#btnSendMail").button('reset');
                }else{
                    alert("{{i18n .Lang "message.email_sent"}}")
                    window.location = res.data;
                }
            },
            error :function () {
                $("#captcha-img").click();
                $("#code").val('');
                layer.msg('{{i18n .Lang "message.system_error"}}');
                $("#btnSendMail").button('reset');
            }
        });

    });
</script>
</body>
</html>