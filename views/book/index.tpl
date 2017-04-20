<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>我的文档 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet" type="text/css">
    <link href="/static/font-awesome/css/font-awesome.min.css" rel="stylesheet" type="text/css">

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
                    <li class="active"><a href="{{urlfor "SettingController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 我的项目</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">项目列表</strong>
                        <a href="{{urlfor "BookController.Create"}}" class="btn btn-success btn-sm pull-right">添加项目</a>
                    </div>
                </div>
                <div class="box-body">
                    <div class="book-list">
                        <div class="list-item">
                            <div class="book-title">
                                <div class="pull-left">
                                    <a href="{{urlfor "BookController.Dashboard" ":key" "test"}}" title="项目概要" data-toggle="tooltip">
                                        <i class="fa fa-unlock" aria-hidden="true"></i> 测试项目
                                    </a>
                                </div>
                                <div class="pull-right">
                                    <a href="{{urlfor "DocumentController.Index" ":key" "test"}}" title="查看文档" data-toggle="tooltip"><i class="fa fa-eye"></i> 查看文档</a>
                                    <a href="{{urlfor "BookController.Edit" ":key" "test" ":id" "1"}}" title="编辑文档" data-toggle="tooltip"><i class="fa fa-edit" aria-hidden="true"></i> 编辑文档</a>
                                </div>
                                <div class="clearfix"></div>
                            </div>
                            <div class="desc-text">
                                <a href="{{urlfor "BookController.Dashboard" ":key" "test"}}" title="项目概要" style="font-size: 12px;" target="_blank">
                                    这是一个测试用户测试项目
                                </a>

                            </div>
                            <div class="info">
                                <span title="创建时间" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-clock-o"></i> 2016/11/27 06:09</span>
                                <span title="创建者" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user"></i> admin</span>
                                <span title="文档数量" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pie-chart"></i> 4</span>
                                <span title="项目角色" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user-secret"></i>  拥有者</span>
                                <span title="最后编辑" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pencil"></i> 最后编辑: admin 于 2017-04-20 12:19</span>
                            </div>
                        </div>
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