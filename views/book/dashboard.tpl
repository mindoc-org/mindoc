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
                    <li class="active"><a href="{{urlfor "BookController.Dashboard" ":key" "test"}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li><a href="{{urlfor "BookController.Users" ":key" "test"}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    <li><a href="{{urlfor "BookController.Setting" ":key" "test"}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"><i class="fa fa-unlock" aria-hidden="true" title="公开项目" data-toggle="tooltip"></i> 这是一个测试项目</strong>
                        <a href="{{urlfor "BookController.Edit" ":key" "test" ":id" "1"}}" class="btn btn-default btn-sm pull-right" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑</a>
                        <a href="{{urlfor "DocumentController.Index" ":key" "test"}}" class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" target="_blank"><i class="fa fa-eye"></i> 阅读</a>
                        <a href="{{urlfor "DocumentController.Index" ":key" "test"}}" class="btn btn-default btn-sm pull-right" style="margin-right: 5px;" target="_blank"><i class="fa fa-upload" aria-hidden="true"></i> 发布</a>
                    </div>
                </div>
                <div class="box-body">
                    <div class="dashboard">
                        <div class="pull-left" style="width: 200px;margin-bottom: 15px;">
                            <div class="book-image">
                                <img src="/static/images/5fcb811e04c23cdb2088f26923fcc287_100.jpg">
                            </div>
                        </div>

                            <div class="list">
                                <span class="title">创建者：</span>
                                <span class="body">
                                Minho
                            </span>
                            </div>
                            <div class="list">
                                <span class="title">文档数量：</span>
                                <span class="body">20</span>
                            </div>
                            <div class="list">
                                <span class="title">创建时间：</span>
                                <span class="body"> 2017-05-25 12:25:45 </span>
                            </div>
                            <div class="list">
                                <span class="title">修改时间：</span>
                                <span class="body"> 2017-05-25 12:25:45 </span>
                            </div>
                            <div class="summary">《TCP/IP详解，卷1：协议》是一本完整而详细的TCP/IP协议指南。描述了属于每一层的各个协议以及它们如何在不同操作系统中运行。<br>
                                <br>
                                作者用Lawrence Berkeley实验室的tcpdump程序来捕获不同操作系统和TCP/IP实现之间传输的不同分组。对tcpdump输出的研究可以帮助理解不同协议如何工作。 <br>
                                <br>
                                本书适合作为计算机专业学生学习网络的教材和教师参考书。也适用于研究网络的技术人员。
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