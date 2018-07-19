<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户中心 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/font-awesome/css/font-awesome.min.css" rel="stylesheet">

    <link href="/static/css/main.css?_=?_=1531986418" rel="stylesheet">
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
                    <li class="active"><a href="{{urlfor "SettingController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 基本信息</a> </li>
                    <li><a href="{{urlfor "SettingController.Password"}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> 修改密码</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">基本信息</strong>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form role="form" method="post" id="memberInfoForm">
                            <div class="form-group">
                                <label>用户名</label>
                                <input type="text" class="form-control disabled">
                            </div>
                            <div class="form-group">
                                <label for="user-nickname">昵称</label>
                                <input type="text" class="form-control" name="userNickname" id="user-nickname" max="20" placeholder="昵称" value="admin">
                            </div>
                            <div class="form-group">
                                <label for="user-email">邮箱<strong class="text-danger">*</strong></label>
                                <input type="email" class="form-control" value="longfei6671@163.com" id="user-email" name="userEmail" max="100" placeholder="邮箱">
                            </div>
                            <div class="form-group">
                                <label>手机号</label>
                                <input type="text" class="form-control" id="user-phone" name="userPhone" maxlength="20" title="手机号码" placeholder="手机号码" value="">
                            </div>
                            <div class="form-group">
                                <label class="description">描述</label>
                                <textarea class="form-control" rows="3" title="描述" name="description" id="description" maxlength="500"></textarea>
                                <p style="color: #999;font-size: 12px;">描述不能超过500字</p>
                            </div>
                            <div class="form-group">
                                <button type="submit" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                                <span id="form-error-message" class="error-message"></span>
                            </div>
                        </form>
                    </div>
                    <div class="form-right">
                        <label>
                            <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                                <img src="/uploads/user/201612/58649aefa944e_58649aef.JPG" onerror="this.src='https://wiki.iminho.me/static/images/middle.gif'" class="img-circle" alt="头像" style="max-width: 120px;max-height: 120px;" id="headimgurl">
                            </a>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
<script src="/static/jquery/1.12.4/jquery.min.js"></script>
<script src="/static/bootstrap/js/bootstrap.min.js"></script>
</body>
</html>