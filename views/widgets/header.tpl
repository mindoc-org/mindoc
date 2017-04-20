<header class="navbar navbar-static-top navbar-fixed-top manual-header" role="banner">
    <div class="container">
        <div class="navbar-header col-sm-12 col-md-6 col-lg-5">
            <a href="/" class="navbar-brand">MinDoc</a>
            <div class="btn-group dropdown-menu-right pull-right slidebar visible-xs-inline-block visible-sm-inline-block">
                <button class="btn btn-default dropdown-toggle hidden-lg" type="button" data-toggle="dropdown"><i class="fa fa-align-justify"></i></button>
                <ul class="dropdown-menu" role="menu">
                    <li class="active"><a href="https://wiki.iminho.me/member/projects" class="item"><i class="fa fa-sitemap"></i> 我的项目</a> </li>
                    <li><a href="https://wiki.iminho.me/member" class="item"><i class="fa fa-user"></i> 个人资料</a> </li>
                    <li><a href="https://wiki.iminho.me/member/account" class="item"><i class="fa fa-lock"></i> 修改密码</a> </li>
                    <li><a href="https://wiki.iminho.me/member/setting" class="item"><i class="fa fa-gear"></i> 开发配置</a> </li>
                    <li><a href="https://wiki.iminho.me/setting/site" class="item"><i class="fa fa-cogs"></i> 网站设置</a> </li>
                    <li><a href="https://wiki.iminho.me/member/users" class="item"><i class="fa fa-group"></i> 用户管理</a> </li>
                    <li>
                        <a href="https://wiki.iminho.me/logout" title="退出登录"><i class="fa fa-sign-out"></i> 退出登录</a>
                    </li>
                </ul>
            </div>
        </div>
        <nav class="navbar-collapse hidden-xs hidden-sm" role="navigation">
            <ul class="nav navbar-nav navbar-right">
                <li>
                    <div class="img user-info" data-toggle="dropdown">
                        <img src="/static/images/headimgurl.jpg" class="img-circle userbar-avatar">
                        <div class="userbar-content">
                            <span>lifei6671</span>
                            <div>管理员</div>
                        </div>
                        <i class="fa fa-chevron-down" aria-hidden="true"></i>
                    </div>
                    <ul class="dropdown-menu user-info-dropdown" role="menu">
                        <li>
                            <a href="{{urlfor "SettingController.Index"}}" title="个人中心"><i class="fa fa-user" aria-hidden="true"></i> 个人中心</a>
                        </li>
                        <li>
                            <a href="{{urlfor "BookController.Index"}}" title="我的项目"><i class="fa fa-book" aria-hidden="true"></i> 我的项目</a>
                        </li>
                        <li>
                            <a href="#" title="用户管理"><i class="fa fa-users" aria-hidden="true"></i> 用户管理</a>
                        </li>
                        <li>
                            <a href="https://wiki.iminho.me/logout" title="退出登录"><i class="fa fa-sign-out"></i> 退出登录</a>
                        </li>
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</header>