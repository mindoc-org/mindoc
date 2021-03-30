<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="utf-8">
    <link rel="shortcut icon" href="{{cdnimg "/favicon.ico"}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="MinDoc" />
    <title>{{i18n .Lang "common.login"}} - Powered by MinDoc</title>
    <meta name="keywords" content="MinDoc,文档在线管理系统,WIKI,wiki,wiki在线,文档在线管理,接口文档在线管理,接口文档管理">
    <meta name="description" content="MinDoc文档在线管理系统 {{.site_description}}">
    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
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
            <form role="form" method="post">
            {{ .xsrfdata }}
                <h3 class="text-center">{{i18n .Lang "common.login"}}</h3>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-user"></i>
                        </div>
                        <input type="text" class="form-control" placeholder="{{i18n .Lang "common.email"}} / {{i18n .Lang "common.username"}}" name="account" id="account" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-lock"></i>
                        </div>
                        <input type="password" class="form-control" placeholder="{{i18n .Lang "common.password"}}" name="password" id="password" autocomplete="off">
                    </div>
                </div>
                {{if .ENABLED_CAPTCHA }}
                {{if ne .ENABLED_CAPTCHA "false"}}
                <div class="form-group">
                    <div class="input-group" style="float: left;width: 195px;">
                        <div class="input-group-addon">
                            <i class="fa fa-check-square"></i>
                        </div>
                        <input type="text" name="code" id="code" class="form-control" style="width: 150px" maxlength="5" placeholder="{{i18n .Lang "common.captcha"}}" autocomplete="off">&nbsp;
                    </div>
                    <img id="captcha-img" style="width: 140px;height: 40px;display: inline-block;float: right" src="{{urlfor "AccountController.Captcha"}}" onclick="this.src='{{urlfor "AccountController.Captcha"}}?key=login&t='+(new Date()).getTime();" title={{i18n .Lang "message.click_to_change"}}>
                    <div class="clearfix"></div>
                </div>
                {{end}}
                {{end}}
                <div class="checkbox">
                    <label>
                        <input type="checkbox" name="is_remember" value="yes"> {{i18n .Lang "common.keep_login"}}
                    </label>
                    <a href="{{urlfor "AccountController.FindPassword" }}" style="display: inline-block;float: right">{{i18n .Lang "common.forgot_password"}}</a>
                </div>
                <div class="form-group">
                    <button type="button" id="btn-login" class="btn btn-success" style="width: 100%"  data-loading-text="{{i18n .Lang "common.logging_in"}}" autocomplete="off">{{i18n .Lang "common.login"}}</button>
                </div>
                {{if .ENABLE_QR_DINGTALK}}
                <div class="form-group">
                    <a id="btn-dingtalk-qr" class="btn btn-default" style="width: 100%" data-loading-text="" autocomplete="off">{{i18n .Lang "common.dingtalk_login"}}</a>
                </div>
                {{end}}
                {{if .ENABLED_REGISTER}}
                {{if ne .ENABLED_REGISTER "false"}}
                <div class="form-group">
                    {{i18n .Lang "message.no_account_yet"}}<a href="{{urlfor "AccountController.Register" }}" title={{i18n .Lang "common.register"}}>{{i18n .Lang "common.register"}}</a>
                </div>
                {{end}}
                {{end}}
            </form>
            <div class="form-group dingtalk-container" style="display: none;">
                <div id="dingtalk-qr-container"></div>
                <a class="btn btn-default btn-dingtalk" style="width: 100%" data-loading-text="" autocomplete="off">{{i18n .Lang "message.return_account_login"}}</a>
            </div>
        </div>
    </div>
    <div class="clearfix"></div>
</div>
{{template "widgets/footer.tpl" .}}
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/dingtalk-jsapi.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/dingtalk-ddlogin.js"}}" type="text/javascript"></script>

<script type="text/javascript">
    if (dd.env.platform !== "notInDingTalk"){
        dd.ready(function() {
            dd.runtime.permission.requestAuthCode({
                corpId: {{ .corpID }} , // 企业id
                onSuccess: function (info) {
                    var index = layer.load(1, {
                        shade: [0.1, '#fff'] // 0.1 透明度的白色背景
                    })

                    var formData = $("form").serializeArray()
                    formData.push({"name": "dingtalk_code", "value": info.code})

                    $.ajax({
                        url: "{{urlfor "AccountController.DingTalkLogin"}} ",
                        data: formData,
                        dataType: "json",
                        type: "POST",
                        complete: function(){
                            layer.close(index)
                        },
                        success: function (res) {
                            if (res.errcode !== 0) {
                                layer.msg(res.message)
                            } else {
                                window.location = "{{ urlfor "HomeController.Index"  }}"
                            }
                        },
                        error: function (res) {
                            layer.msg("发生异常")
                        }
                    })
                }
            });
        });
    }

</script>

<script type="text/javascript">
    var url = 'https://oapi.dingtalk.com/connect/oauth2/sns_authorize?appid={{.dingtalk_qr_key}}&response_type=code&scope=snsapi_login&state=1&redirect_uri={{ urlfor "AccountController.QRLogin" ":app" "dingtalk"}}'
    var obj = DDLogin({
        id:"dingtalk-qr-container",
        goto: encodeURIComponent(url), 
        style: "border:none;background-color:#FFFFFF;",
        width : "338",
        height: "300"
    });
    var handleMessage = function (event) {
        var origin = event.origin;
        if( origin == "https://login.dingtalk.com" ) { //判断是否来自ddLogin扫码事件。
            layer.load(1, { shade: [0.1, '#fff'] })
            var loginTmpCode = event.data; 
            //获取到loginTmpCode后就可以在这里构造跳转链接进行跳转了
            console.log("loginTmpCode", loginTmpCode);
            url = url + "&loginTmpCode=" + loginTmpCode
            window.location = url
        }
    };
    if (typeof window.addEventListener != 'undefined') {
        window.addEventListener('message', handleMessage, false);
    } else if (typeof window.attachEvent != 'undefined') {
        window.attachEvent('onmessage', handleMessage);
    }
</script>

<script type="text/javascript">
    $(function () {
        $("#account,#password,#code").on('focus', function () {
            $(this).tooltip('destroy').parents('.form-group').removeClass('has-error');
        });

        $(document).keydown(function (e) {
            var event = document.all ? window.event : e;
            if (event.keyCode === 13) {
                $("#btn-login").click();
            }
        });

        $("#btn-dingtalk-qr").on('click', function(){
            $('form').hide()
            $(".dingtalk-container").show()
        })

        $(".btn-dingtalk").on('click', function(){
            $('form').show()
            $(".dingtalk-container").hide()
        })

        $("#btn-login").on('click', function () {
            $(this).tooltip('destroy').parents('.form-group').removeClass('has-error');
            var $btn = $(this).button('loading');

            var account = $.trim($("#account").val());
            var password = $.trim($("#password").val());
            var code = $("#code").val();

            if (account === "") {
                $("#account").tooltip({ placement: "auto", title: "{{i18n .Lang "message.account_empty"}}", trigger: 'manual' })
                    .tooltip('show')
                    .parents('.form-group').addClass('has-error');
                $btn.button('reset');
                return false;
            } else if (password === "") {
                $("#password").tooltip({ title: '{{i18n .Lang "message.password_empty"}}', trigger: 'manual' })
                    .tooltip('show')
                    .parents('.form-group').addClass('has-error');
                $btn.button('reset');
                return false;
            } else if (code !== undefined && code === "") {
                $("#code").tooltip({ title: '{{i18n .Lang "message.captcha_empty"}}', trigger: 'manual' })
                    .tooltip('show')
                    .parents('.form-group').addClass('has-error');
                $btn.button('reset');
                return false;
            } else {
                $.ajax({
                    url: "{{urlfor "AccountController.Login" "url" .url}}",
                    data: $("form").serializeArray(),
                    dataType: "json",
                    type: "POST",
                    success: function (res) {
                        if (res.errcode !== 0) {
                            $("#captcha-img").click();
                            $("#code").val('');
                            layer.msg(res.message);
                            $btn.button('reset');
                        } else {
                            turl = res.data;
                            if (turl === "") {
                                turl = "/";
                            }
                            window.location = turl;
                        }
                    },
                    error: function () {
                        $("#captcha-img").click();
                        $("#code").val('');
                        layer.msg('{{i18n .Lang "message.system_error"}}');
                        $btn.button('reset');
                    }
                });
            }

            return false;
        });
    });
</script>
</body>
</html>