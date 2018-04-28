<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>权限资源管理 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">
    <link href="/static/bootstrap/plugins/bootstrap-treegrid/css/jquery.treegrid.css" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/css/main.css"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="{{cdnjs "/static/html5shiv/3.7.3/html5shiv.min.js"}}"></script>
    <script src="{{cdnjs "/static/respond.js/1.4.2/respond.min.js" }}"></script>
    <![endif]-->
</head>
<body>
<div class="manual-reader">
{{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                {{template "manager/manager_widgets.tpl.tpl" .}}
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">权限资源管理</strong>
                        <button type="button" class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#addMemberDialogModal">添加资源</button>
                    </div>
                </div>
                <div class="box-body" id="resourceList">
                    <div  class="table-responsive">
                        <table class="table table-bordered tree" id="resourceTreeGrid">
                            <thead>
                            <tr>
                                <th>资源名称</th>
                                <th>控制器名称</th>
                                <th>动作名称</th>
                                <th>请求类型</th>
                                <th>操作</th>
                            </tr>
                            </thead>
                            <tbody>
                        {{range $index,$item := .Lists}}
                                <tr>
                                    <td></td>
                                </tr>
                        {{end}}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
</div>
<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBookModal" tabindex="-1" role="dialog" aria-labelledby="deleteBookModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBookForm" action="{{urlfor "ManagerController.DeleteBook"}}">
            <input type="hidden" name="book_id" value="">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">删除项目</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">确定删除项目吗？</span>
                    <p></p>
                    <p class="text error-message">删除项目后将无法找回。</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="删除中...">确定删除</button>
                </div>
            </div>
        </form>
    </div>
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="/static/bootstrap/plugins/bootstrap-treegrid/js/jquery.treegrid.js"></script>
<script src="/static/bootstrap/plugins/bootstrap-treegrid/js/jquery.treegrid.bootstrap3.js"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">

    /**
     * 删除项目
     */
    function deleteBook($id) {
        $("#deleteBookModal").find("input[name='book_id']").val($id);
        $("#deleteBookModal").modal("show");
    }
    $(function () {
        $("#resourceTreeGrid").treegrid();
        /**
         * 删除项目
         */
        $("#deleteBookForm").ajaxForm({
            beforeSubmit : function () {
                $("#btnDeleteBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    window.location = window.location.href;
                }else{
                    showError(res.message,"#form-error-message2");
                }
                $("#btnDeleteBook").button("reset");
            },
            error : function () {
                showError("服务器异常","#form-error-message2");
                $("#btnDeleteBook").button("reset");
            }
        });
    });
</script>
</body>
</html>