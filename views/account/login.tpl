<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="utf-8">
    <link rel="shortcut icon" href="{{cdnimg "/favicon.ico"}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="MinDoc" />
    <title>用户登录 - Powered by MinDoc</title>
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
                <h3 class="text-center">用户登录</h3>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-user"></i>
                        </div>
                        <input type="text" class="form-control" placeholder="邮箱 / 用户名" name="account" id="account" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-lock"></i>
                        </div>
                        <input type="password" class="form-control" placeholder="密码" name="password" id="password" autocomplete="off">
                    </div>
                </div>
                {{if .ENABLED_CAPTCHA }}
                {{if ne .ENABLED_CAPTCHA "false"}}
                <div class="form-group">
                    <div class="input-group" style="float: left;width: 195px;">
                        <div class="input-group-addon">
                            <i class="fa fa-check-square"></i>
                        </div>
                        <input type="text" name="code" id="code" class="form-control" style="width: 150px" maxlength="5" placeholder="验证码" autocomplete="off">&nbsp;
                    </div>
                    <img id="captcha-img" style="width: 140px;height: 40px;display: inline-block;float: right" src="{{urlfor "AccountController.Captcha"}}" onclick="this.src='{{urlfor "AccountController.Captcha"}}?key=login&t='+(new Date()).getTime();" title="点击换一张">
                    <div class="clearfix"></div>
                </div>
                {{end}}
                {{end}}
                <div class="checkbox">
                    <label>
                        <input type="checkbox" name="is_remember" value="yes"> 保持登录
                    </label>
                    <a href="{{urlfor "AccountController.FindPassword" }}" style="display: inline-block;float: right">忘记密码？</a>
                </div>
                <div class="form-group">
                    <button type="button" id="btn-login" class="btn btn-success" style="width: 100%"  data-loading-text="正在登录..." autocomplete="off">立即登录</button>
                </div>
                {{if .ENABLED_REGISTER}}
                {{if ne .ENABLED_REGISTER "false"}}
                <div class="form-group">
                    还没有账号？<a href="{{urlfor "AccountController.Register" }}" title="立即注册">立即注册</a>
                </div>
                {{end}}
                {{end}}
            </form>
        </div>
    </div>
    <div class="clearfix"></div>
</div>
{{template "widgets/footer.tpl" .}}
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/dingtalk-jsapi.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        if (dd.env.platform !== "notInDingTalk"){
            dd.ready(function() {
                dd.runtime.permission.requestAuthCode({
                    corpId: {{ .corpID }} , // 企业id
                    onSuccess: function (info) {
                        var index = layer.load(1, {
                            shade: [0.1, '#fff'] // 0.1 透明度的白色背景
                        })

                        var formData = $("form").serializeArray()
                        formData.push({"name": "code", "value": info.code})

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
    })

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

        $("#btn-login").on('click', function () {
            $(this).tooltip('destroy').parents('.form-group').removeClass('has-error');
            var $btn = $(this).button('loading');

            var account = $.trim($("#account").val());
            var password = $.trim($("#password").val());
            var code = $("#code").val();

            if (account === "") {
                $("#account").tooltip({ placement: "auto", title: "账号不能为空", trigger: 'manual' })
                    .tooltip('show')
                    .parents('.form-group').addClass('has-error');
                $btn.button('reset');
                return false;
            } else if (password === "") {
                $("#password").tooltip({ title: '密码不能为空', trigger: 'manual' })
                    .tooltip('show')
                    .parents('.form-group').addClass('has-error');
                $btn.button('reset');
                return false;
            } else if (code !== undefined && code === "") {
                $("#code").tooltip({ title: '验证码不能为空', trigger: 'manual' })
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
                        layer.msg('系统错误');
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