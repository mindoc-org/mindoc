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
<div id="app">
    <template>
        <Dropdown placement="bottom-start">
            <a href="javascript:void(0)">
                菜单(左)
                <Icon type="arrow-down-b"></Icon>
            </a>
            <Dropdown-menu slot="list">
                <Dropdown-item>驴打滚</Dropdown-item>
                <Dropdown-item>炸酱面</Dropdown-item>
                <Dropdown-item>豆汁儿</Dropdown-item>
                <Dropdown-item>冰糖葫芦</Dropdown-item>
                <Dropdown-item>北京烤鸭</Dropdown-item>
            </Dropdown-menu>
        </Dropdown>
        <Dropdown style="margin-left: 20px">
            <a href="javascript:void(0)">
                菜单(居中)
                <Icon type="arrow-down-b"></Icon>
            </a>
            <Dropdown-menu slot="list">
                <Dropdown-item>驴打滚</Dropdown-item>
                <Dropdown-item>炸酱面</Dropdown-item>
                <Dropdown-item>豆汁儿</Dropdown-item>
                <Dropdown-item>冰糖葫芦</Dropdown-item>
                <Dropdown-item>北京烤鸭</Dropdown-item>
            </Dropdown-menu>
        </Dropdown>
        <Dropdown style="margin-left: 20px" placement="bottom-end">
            <a href="javascript:void(0)">
                菜单(右)
                <Icon type="arrow-down-b"></Icon>
            </a>
            <Dropdown-menu slot="list">
                <Dropdown-item>驴打滚</Dropdown-item>
                <Dropdown-item>炸酱面</Dropdown-item>
                <Dropdown-item>豆汁儿</Dropdown-item>
                <Dropdown-item>冰糖葫芦</Dropdown-item>
                <Dropdown-item>北京烤鸭</Dropdown-item>
            </Dropdown-menu>
        </Dropdown>
    </template>

</div>
</body>
<script type="text/javascript">
    new Vue({
        el: '#app',
        data: function(){
            return { visible: false }
        },
        methods : {
            hasParent : function () {
                return false
            }
        }
    })
</script>
</body>
</html>
