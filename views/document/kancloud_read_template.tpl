<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>编辑文档 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/jstree/3.3.4/themes/default/style.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/kancloud.css" "version"}}" rel="stylesheet">
    <link href="{{cdncss "/static/jquery/plugins/imgbox/imgbox.css" "version"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
<div class="m-manual manual-reader">
    <header class="navbar navbar-static-top manual-head" role="banner">
        <div class="container-fluid">
            <div class="navbar-header pull-left manual-title">
                <span class="slidebar" id="slidebar"><i class="fa fa-align-justify"></i></span>

                <span style="font-size: 12px;font-weight: 100;">v 0.1.1</span>
            </div>
            <div class="navbar-header pull-right manual-menu">
                <div class="dropdown">
                    <button id="dLabel" class="btn btn-default" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        项目
                        <span class="caret"></span>
                    </button>
                    <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dLabel">
                        <li><a href="javascript:" data-toggle="modal" data-target="#shareProject">项目分享</a> </li>
                        <li role="presentation" class="divider"></li>
                        <li><a href="https://wiki.iminho.me/export/1" target="_blank">项目导出</a> </li>
                        <li><a href="https://wiki.iminho.me" title="返回首页">返回首页</a> </li>
                    </ul>
                </div>
            </div>
        </div>
    </header>
    <article class="container-fluid manual-body">
        <div class="manual-left">
            <div class="manual-tab">
                <div class="tab-navg">
                    <span data-mode="view" class="navg-item active"><i class="fa fa-align-justify"></i><b class="text">目录</b></span>
                </div>
                <div class="tab-wrap">
                    <div class="tab-item manual-catalog">
                        <div class="catalog-list read-book-preview" id="sidebar">
                            <ul><li id="1" class="jstree-open"><a href="https://wiki.iminho.me/docs/show/1" title="项目简介" class="jstree-clicked">项目简介</a></li><li id="8"><a href="https://wiki.iminho.me/docs/show/8" title="使用手册">使用手册</a><ul><li id="135"><a href="https://wiki.iminho.me/docs/show/135" title="SmartWiki依赖扩展">SmartWiki依赖扩展</a></li><li id="9"><a href="https://wiki.iminho.me/docs/show/9" title="SmartWiki安装与部署">SmartWiki安装与部署</a></li><li id="10"><a href="https://wiki.iminho.me/docs/show/10" title="密码找回与邮件配置">密码找回与邮件配置</a></li><li id="21"><a href="https://wiki.iminho.me/docs/show/21" title="Artisan命令安装SmartWiki">Artisan命令安装SmartWiki</a></li><li id="22"><a href="https://wiki.iminho.me/docs/show/22" title="在Docker中使用SmartWiki">在Docker中使用SmartWiki</a></li><li id="105"><a href="https://wiki.iminho.me/docs/show/105" title="图片和附件配置">图片和附件配置</a></li><li id="133"><a href="https://wiki.iminho.me/docs/show/133" title="查看SmartWiki版本">查看SmartWiki版本</a></li><li id="134"><a href="https://wiki.iminho.me/docs/show/134" title="使用phpStudy安装SmartWiki">使用phpStudy安装SmartWiki</a></li><li id="138"><a href="https://wiki.iminho.me/docs/show/138" title="接口管理和测试工具">接口管理和测试工具</a></li></ul></li><li id="90"><a href="https://wiki.iminho.me/docs/show/90" title="PHP环境搭建">PHP环境搭建</a><ul><li id="18"><a href="https://wiki.iminho.me/docs/show/18" title="Linux命令行下Nginx+PHP-FPM安装与配置">Linux命令行下Nginx+PHP-FPM安装与配置</a></li><li id="17"><a href="https://wiki.iminho.me/docs/show/17" title="Linux命令行下Apache+PHP安装与配置">Linux命令行下Apache+PHP安装与配置</a></li><li id="91"><a href="https://wiki.iminho.me/docs/show/91" title="Linux下使用XAMPP搭建PHP环境">Linux下使用XAMPP搭建PHP环境</a></li><li id="92"><a href="https://wiki.iminho.me/docs/show/92" title="Windows下使用phpStudy搭建php环境">Windows下使用phpStudy搭建php环境</a></li></ul></li><li id="4"><a href="https://wiki.iminho.me/docs/show/4" title="Composer安装">Composer安装</a></li><li id="5"><a href="https://wiki.iminho.me/docs/show/5" title="更新日志">更新日志</a></li><li id="6"><a href="https://wiki.iminho.me/docs/show/6" title="常见问题">常见问题</a></li></ul>
                        </div>

                    </div>
                </div>
            </div>
            <div class="m-copyright">
                <p>
                    本文档使用 <a href="https://www.iminho.me" target="_blank">MinDoc</a> 发布
                </p>
            </div>
        </div>
        <div class="manual-right">
            <div class="manual-article">
                <div class="article-head">
                    <div class="container-fluid">
                        <div class="row">
                            <div class="col-md-2">

                            </div>
                            <div class="col-md-8 text-center">
                                <h1 id="article-title">项目简介</h1>
                            </div>
                            <div class="col-md-2">

                            </div>
                        </div>
                    </div>

                </div>
                <div class="article-content">
                    <div class="article-body editor-content"  id="page-content">
                        <h2 id="h2-u7b80u4ecb"><a name="简介" class="reference-link"></a><span class="header-link octicon octicon-link"></span>简介</h2>
                        <p>SmartWiki是一款针对IT团队开发的简单好用的文档管理系统。
                            可以用来储存日常接口文档，数据库字典，手册说明等文档。内置项目管理，用户管理，权限管理等功能，能够满足大部分中小团队的文档管理需求。</p>
                        <h2 id="h2-u4f7fu7528"><a name="使用" class="reference-link"></a><span class="header-link octicon octicon-link"></span>使用</h2>
                        <pre><code>git clone https://github.com/lifei6671/SmartWiki.git
</code></pre>
                        <p>配置laravel的运行环境,然后打开首页会自动跳转到安装页面。</p>
                        <p>因为laravel使用了composer，所以需要服务器安装composer进行包的还原。</p>
                        <h2 id="h2-u90e8u5206u622au56fe"><a name="部分截图" class="reference-link"></a><span class="header-link octicon octicon-link"></span>部分截图</h2>
                        <p><strong>个人资料</strong></p>
                        <p><img src="https://raw.githubusercontent.com/lifei6671/SmartWiki/master/storage/app/images/20161124082553.png" alt="个人资料" /></p>
                        <p><strong>我的项目</strong></p>
                        <p><img src="https://raw.githubusercontent.com/lifei6671/SmartWiki/master/storage/app/images/20161124082647.png" alt="我的项目" /></p>
                        <p><strong>项目参与用户</strong></p>
                        <p><img src="https://raw.githubusercontent.com/lifei6671/SmartWiki/master/storage/app/images/20161124082703.png" alt="项目参与用户" /></p>
                        <p><strong>文档编辑</strong></p>
                        <p><img src="https://raw.githubusercontent.com/lifei6671/SmartWiki/master/storage/app/images/20161124082810.png" alt="文档编辑" /></p>
                        <p><strong>文档模板</strong></p>
                        <p><img src="https://raw.githubusercontent.com/lifei6671/SmartWiki/master/storage/app/images/20161124082844.png" alt="文档模板" /></p>
                        <h2 id="h2-u4f7fu7528u7684u6280u672f"><a name="使用的技术" class="reference-link"></a><span class="header-link octicon octicon-link"></span>使用的技术</h2>
                        <ul>
                            <li>laravel 5.2</li>
                            <li>mysql 5.6</li>
                            <li>editor.md</li>
                            <li>bootstrap 3.2</li>
                            <li>jquery 库</li>
                            <li>layer 弹出层框架</li>
                            <li>webuploader 文件上传框架</li>
                            <li>Nprogress 库</li>
                            <li>jstree</li>
                            <li>font awesome 字体库</li>
                            <li>cropper 图片剪裁库</li>
                        </ul>
                        <h2 id="h2-u529fu80fd"><a name="功能" class="reference-link"></a><span class="header-link octicon octicon-link"></span>功能</h2>
                        <ol>
                            <li>项目管理，可以对项目进行编辑更改，成员添加等。</li>
                            <li>文档管理，添加和删除文档，文档历史恢复等。</li>
                            <li>用户管理，添加和禁用用户，个人资料更改等。</li>
                            <li>用户权限管理 ， 实现用户角色的变更。</li>
                            <li>项目加密，可以设置项目公开状态为私密、半公开、全公开。</li>
                            <li>站点配置，二次开发时可以添加自定义配置项。</li>
                        </ol>
                        <h2 id="h2-u5f85u5b9eu73b0"><a name="待实现" class="reference-link"></a><span class="header-link octicon octicon-link"></span>待实现</h2>
                        <ol>
                            <li>项目导出</li>
                            <li>角色细分</li>
                            <li>文档搜索</li>
                        </ol>
                        <h2 id="h2-u4f5cu8005"><a name="作者" class="reference-link"></a><span class="header-link octicon octicon-link"></span>作者</h2>
                        <p>一个纯粹的PHPer。<a href="http://www.iminho.me">个人网站</a></p>

                    </div>
                </div>
            </div>
        </div>
        <div class="manual-progress"><b class="progress-bar"></b></div>
    </article>
    <div class="manual-mask"></div>
</div>

<!-- Share Modal -->
<div class="modal fade" id="shareProject" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title" id="myModalLabel">项目分享</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="password" class="col-sm-2 control-label">项目地址</label>
                    <div class="col-sm-10">
                        <input type="text" value="https://wiki.iminho.me/show/1" class="form-control" onmouseover="this.select()" id="projectUrl" title="项目地址">
                    </div>
                    <div class="clearfix"></div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
            </div>
        </div>
    </div>
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/jstree/3.3.4/jstree.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/jquery/plugins/imgbox/jquery.imgbox.pack.js"}}"></script>
<script type="text/javascript">
    $(function () {
        $("#sidebar").jstree({
            'plugins':["wholerow","types"],
            "types": {
                "default" : {
                    "icon" : false  // 删除默认图标
                }
            },
            'core' : {
                'check_callback' : false,
                "multiple" : false ,
                'animation' : 0
            }
        });
        $("img").imgbox();
    });
</script>
</body>
</html>