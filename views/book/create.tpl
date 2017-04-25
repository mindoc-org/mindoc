<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>概要 - Powered by MinDoc</title>

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
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" "test"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li><a href="{{urlfor "BookController.Users" ":key" "test"}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    {{if eq .Model.RoleId 0 1}}
                    <li class="active"><a href="{{urlfor "BookController.Setting" ":key" "test"}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                    {{end}}
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 项目设置</strong>
                        <button type="button"  class="btn btn-success btn-sm pull-right">转让项目</button>
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;">删除项目</button>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm">
                            <div class="form-group">
                                <label>标题</label>
                                <input type="text" class="form-control" placeholder="项目名称">
                            </div>

                        </form>
                    </div>
                    <div class="form-right">
                        <label>
                            <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                                <img src="/static/images/5fcb811e04c23cdb2088f26923fcc287_100.jpg" onerror="this.src='https://wiki.iminho.me/static/images/middle.gif'" alt="头像" style="max-width: 120px;" id="headimgurl">
                            </a>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="/static/jquery/1.12.4/jquery.min.js"></script>
<script src="/static/bootstrap/js/bootstrap.min.js"></script>
<script src="/static/js/main.js" type="text/javascript"></script>

</body>
</html>