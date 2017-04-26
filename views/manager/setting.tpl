<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>配置管理 - Powered by MinDoc</title>

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
                    <li><a href="{{urlfor "ManagerController.Index"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 仪表盘</a> </li>
                    <li><a href="{{urlfor "ManagerController.Users" }}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Books" }}" class="item"><i class="fa fa-book" aria-hidden="true"></i> 项目管理</a> </li>
                    <li><a href="{{urlfor "ManagerController.Comments" }}" class="item"><i class="fa fa-comments-o" aria-hidden="true"></i> 评论管理</a> </li>
                    <li class="active"><a href="{{urlfor "ManagerController.Setting" }}" class="item"><i class="fa fa-cogs" aria-hidden="true"></i> 配置管理</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 配置管理</strong>
                    </div>
                </div>
                <div class="box-body">
                    <form method="post" id="gloablEditForm" action="{{urlfor "ManagerController.Setting"}}">
                            <div class="form-group">
                                <label>网站标题</label>
                                <input type="text" class="form-control" name="SITE_NAME" id="siteName" placeholder="网站标题" value="{{.SITE_NAME.OptionValue}}">
                            </div>
                            <div class="form-group">
                                <label>启用匿名访问</label>
                                <div class="radio">
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .ENABLE_ANONYMOUS.OptionValue "true"}}checked{{end}} name="ENABLE_ANONYMOUS" value="true">开启<span class="text"></span>
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .ENABLE_ANONYMOUS.OptionValue "false"}}checked{{end}} name="ENABLE_ANONYMOUS" value="false">关闭<span class="text"></span>
                                    </label>
                                </div>
                            </div>
                        <div class="form-group">
                            <label>启用注册</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_REGISTER.OptionValue "true"}}checked{{end}} name="ENABLED_REGISTER" value="true">开启<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_REGISTER.OptionValue "false"}}checked{{end}} name="ENABLED_REGISTER" value="false">关闭<span class="text"></span>
                                </label>
                            </div>
                        </div>
                        <div class="form-group">
                            <label>启用验证码</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_CAPTCHA.OptionValue "true"}}checked{{end}} name="ENABLED_CAPTCHA" value="true">开启<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_CAPTCHA.OptionValue "false"}}checked{{end}} name="ENABLED_CAPTCHA" value="false">关闭<span class="text"></span>
                                </label>
                            </div>
                        </div>
                            <div class="form-group">
                                <button type="submit" id="btnSaveBookInfo" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
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


<script src="/static/jquery/1.12.4/jquery.min.js" type="text/javascript"></script>
<script src="/static/bootstrap/js/bootstrap.min.js" type="text/javascript"></script>
<script src="/static/js/jquery.form.js" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {

    });
</script>
</body>
</html>