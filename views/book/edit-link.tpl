<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>链接 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="{{cdncss "/static/css/main.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/acidjs-css3-treeview.css"}}" rel="stylesheet">

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
            {{template "widgets/sidebar.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{.Model.BookName}}</strong>
                    </div>
                </div>
                <div class="box-body">
                    <form role="form" method="post" id="editLinkForm">
                        <input type="hidden" name="link_docs" id="link_docs" value="{{.LinkDocLinks}}">
                        <div class="form-group">
                            <div class="acidjs-css3-treeview">
                                {{.LinkDocResult}}
                            </div>
                        </div>
                        <div class="form-group">
                            <button type="submit" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $(".acidjs-css3-treeview").delegate("label input:checkbox", "change", function() {
            var
                checkbox = $(this),
                nestedList = checkbox.parent().next().next(),
                selectNestedListCheckbox = nestedList.find("label:not([for]) input:checkbox");
                selectid =   checkbox.parent().prev();
            var link_docs = $("#link_docs").val();
            if(checkbox.is(":checked")) {
                link_docs = selectid.attr("id") + "," + link_docs
                selectNestedListCheckbox.prop("checked", true);
            }else{
                link_docs = link_docs.replace(selectid.attr("id") + ",","");
                selectNestedListCheckbox.prop("checked", false);
            }
            $("#link_docs").val(link_docs);
        });

        $("#editLinkForm").ajaxForm({
            beforeSubmit : function () {
                var link_docs = $("#link_docs").val();
                if(!link_docs ){
                    showError("没有选择任何文档");
                    return false;
                }
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess('保存成功');
                }else{
                    showError(res.message);
                }
            }
        }) ;
    });
</script>

</body>
</html>