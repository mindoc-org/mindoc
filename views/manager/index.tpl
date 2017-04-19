<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <link rel="stylesheet" href="/static/font-awesome/css/font-awesome.css">
    <link rel="stylesheet" href="/static/iview/styles/iview.css">
    <link rel="stylesheet" href="/static/css/main.css">
    <script type="text/javascript" src="/static/vuejs/vue.min.js"></script>
    <script type="text/javascript" src="/static/iview/iview.min.js"></script>

</head>
<body>
<div id="app" class="container-body">
    <div class="header">
        <div class="container-left">
            <a href="/" class="topbar-home-link">GoDoc 文档在线管理系统</a>
        </div>
        <div class="container-right">
            <template>
                <Dropdown  class="topbar-info" placement="bottom-start">
                    <a href="javascript:void(0)" class="topbar-member-link">
                        Admin
                        <Icon type="arrow-down-b"></Icon>
                    </a>
                    <Dropdown-menu slot="list">
                        <Dropdown-item><a href="/"> <i class="fa fa-user" aria-hidden="true"></i> 个人中心</a></Dropdown-item>
                        <Dropdown-item><i class="fa fa-key" aria-hidden="true"></i> 重置密码</Dropdown-item>
                        <Dropdown-item><i class="fa fa-sign-out" aria-hidden="true"></i> 退出登录</Dropdown-item>
                    </Dropdown-menu>
                </Dropdown>
            </template>
        </div>
        <div class="clearfix"></div>
    </div>
    <div class="container-main">
        <div class="sidebar">

        </div>
        <div class="">

        </div>
    </div>

</div>

</body>

<script type="text/javascript">
    new Vue({
        el: '#app',
        data: function ()  {
            return {

            }
        },
        computed : {
        },
        methods: {

        }
    })
</script>
</html>
