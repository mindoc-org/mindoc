<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="utf-8">
    <link rel="shortcut icon" href="{{cdnimg "/static/favicon.ico"}}">
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
    <style type="text/css">
        .login > .login-body {
            text-align: center;
            padding-top: 1.5em;
        }
        .login > .login-body > a > strong:hover {
            border-bottom: 1px solid #337ab7;
        }
        .login > .login-body > a > strong {
            font-size: 1.5em;
            vertical-align: middle;
            padding: 0.5em;
        }
        .bind-existed-form > .form-group {
            margin: auto 1.5em;
            margin-top: 1em;
        }
    </style>
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
    <script type="text/javascript">
        window.bind_existed = {{ .bind_existed }};
        window.user_info_json = {{ .user_info_json }};
        window.server_error_msg = "{{ .error_msg }}";
        window.home_url = "{{ .BaseUrl }}";
        window.workweixin_login_bind = "{{urlfor "AccountController.WorkWeixinLoginBind"}}";
        window.workweixin_login_ignore = "{{urlfor "AccountController.WorkWeixinLoginIgnore"}}";
    </script>
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
            返回 <a href="{{ .BaseUrl }}"><strong>首页</strong></a>
        </div>
    </div>
    <div class="clearfix"></div>
    <script type="text/x-template" id="bind-existed-template">
        <div role="form" class="bind-existed-form">
            {{ .xsrfdata }}
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
        </div>
    </script>
</div>
{{template "widgets/footer.tpl" .}}
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    function showBindAccount() {
        layer.confirm([
            '检测到当前登录企业微信未绑定已有账户, 是否需要绑定已有账户?<br />',
            '<ul style="padding-left: 1.2em;">',
                '<li>若已有账户， 请 <strong>去绑定</strong></li>',
                '<li>若没有现有账户, 请 <strong>忽略绑定</strong></li>',
            '</ul>'
        ].join(''), {
            title: "WIKI-绑定提示",
            move: false,
            area: 'auto',
            offset: 'auto',
            icon: 3,
            btn: ['去绑定','忽略绑定'],
        }, function(index, layero){
            // layer.close(index);
            // layer.msg(window.home_url);
            // TODO: 现有账户[用户名+密码]查询现有账户 依据Session[user_info]绑定更新现有账户
            console.log("yes");
            layer.open({
                title: "绑定已有账户",
                type: 1,
                move: false,
                area: 'auto',
                offset: 'auto',
                content: $('#bind-existed-template').html(),
                btn: ['绑定','取消'],
                yes: function(index, layero){
                    $.ajax({
                        url: window.workweixin_login_bind,
                        type: 'POST',
                        beforeSend: function(request) {
                            request.setRequestHeader("X-Xsrftoken", $('.bind-existed-form input[name="_xsrf"]').val());
                        },
                        data: {
                            account: $('#account').val(),
                            password: $('#password').val()
                        },
                        dataType: 'json',
                        success: function(data) {
                            if(data.errcode == 0) {
                                layer.close(index);
                                // layer.msg(JSON.stringify(data), {icon: 1, time: 15500});
                                window.location.href = window.home_url;
                            }
                            else {
                                layer.msg(data.message, {icon: 5, time: 3500});
                            }
                        },
                        error: function(data) {
                            console.log(data);
                        }
                    });
                    return false;
                },
                cancel: function(index, layero){ 
                    // return false; // 不关闭
                    layer.close(index);
                    window.location.href = window.home_url;
                }
            });
        }, function(index){
            /*
            // TODO: 依据Session[user_info]创建新账户
            console.log("no");
            var msg = '';
            // msg = "<pre>" + JSON.stringify(window.location, null, 4) + "</pre>";
            msg = "<pre>" + JSON.stringify(window.user_info_json, null, 4) + "</pre>";
            // msg = "<pre>" + window.user_info_json + "</pre>";
            layer.open({
                title: "Degug-UserInfo",
                type: 1,
                skin: 'layui-layer-rim',
                move: false,
                area: 'auto',
                offset: 'auto',
                content: msg
            });
            */
            $.ajax({
                url: window.workweixin_login_ignore,
                type: 'GET',
                beforeSend: function(request) {
                    request.setRequestHeader("X-Xsrftoken", $('.bind-existed-form input[name="_xsrf"]').val());
                },
                data: {},
                dataType: 'json',
                success: function(data) {
                    if(data.errcode == 0) {
                        layer.close(index);
                        layer.msg(JSON.stringify(data), {icon: 1, time: 15500});
                        window.location.href = window.home_url;
                    }
                    else {
                        layer.msg(data.message, {icon: 5, time: 3500});
                    }
                },
                error: function(data) {
                    console.log(data);
                }
            });
            return false;
        });
    }
    $(document).ready(function () {
        $('#debug-panel').val($('html').html());
        if (!!window.server_error_msg && window.server_error_msg.length > 0) {
            layer.msg(window.server_error_msg, {icon: 5, time: 3500});
        } else {
            if (window.bind_existed === false) {
                showBindAccount();
            } else {
                // alert(typeof window.bind_existed);
                // alert('_' + window.bind_existed + '_');
                window.location.href = window.home_url;
            }
        }
    });
</script>
</body>
</html>