<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.project_mgrt"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
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
            {{template "manager/widgets.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{i18n .Lang "mgr.proj_list"}}</strong>
                    </div>
                </div>
                <div class="box-body" id="bookList">
                    <div class="book-list">

                        {{range $index,$item := .Lists}}
                        <div class="list-item">
                                <div class="book-title">
                                    <div class="pull-left">
                                        <a href="{{urlfor "ManagerController.EditBook" ":key" $item.Identify}}" title="{{i18n .Lang "mgr.edit_proj"}}" data-toggle="tooltip">
                                            {{if eq $item.PrivatelyOwned 0}}
                                            <i class="fa fa-unlock" aria-hidden="true"></i>
                                            {{else}}
                                            <i class="fa fa-lock" aria-hidden="true"></i>
                                            {{end}}
                                            {{$item.BookName}}
                                        </a>
                                    </div>
                                    <div class="pull-right">
                                        <div class="btn-group">
                                            <a href="{{urlfor "DocumentController.Edit" ":key" $item.Identify ":id" ""}}" class="btn btn-default" target="_blank">{{i18n $.Lang "common.edit"}}</a>
                                            <a href="javascript:;" class="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                                <span class="caret"></span>
                                                <span class="sr-only">Toggle Dropdown</span>
                                            </a>

                                            <ul class="dropdown-menu">
                                                <li><a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" target="_blank">{{i18n $.Lang "common.read"}}</a></li>
                                                <li><a href="{{urlfor "ManagerController.EditBook" ":key" $item.Identify}}">{{i18n $.Lang "common.setting"}}</a></li>
                                                <li><a href="javascript:deleteBook('{{$item.BookId}}');">{{i18n $.Lang "common.delete"}}</a> </li>
                                            </ul>
                                        </div>
                                        {{/*<a href="{{urlfor "DocumentController.Index" ":key" $item.Identify}}" title="查看文档" data-toggle="tooltip" target="_blank"><i class="fa fa-eye"></i> 查看文档</a>*/}}
                                        {{/*<a href="{{urlfor "DocumentController.Edit" ":key" $item.Identify ":id" ""}}" title="编辑文档" data-toggle="tooltip" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑文档</a>*/}}
                                    </div>
                                    <div class="clearfix"></div>
                                </div>
                                <div class="desc-text">
                                    {{if eq $item.Description ""}}
                                    &nbsp;
                                    {{else}}
                                        <a href="{{urlfor "ManagerController.EditBook" ":key" $item.Identify}}" title="{{i18n .Lang "mgr.edit_proj"}}" style="font-size: 12px;" target="_blank">
                                            {{$item.Description}}
                                        </a>
                                    {{end}}
                                </div>
                                <div class="info">
                                <span title="{{i18n $.Lang "mgr.create_time"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-clock-o"></i>
                                    {{date_format $item.CreateTime "2006-01-02 15:04:05"}}

                                </span>
                                    <span title="{{i18n $.Lang "mgr.creator"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user"></i> {{if eq $item.RealName "" }}{{$item.CreateName}}{{else}}{{$item.RealName}}{{end}}</span>
                                    <span title="{{i18n $.Lang "mgr.doc_amount"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pie-chart"></i> {{$item.DocCount}}</span>
                                   {{if ne $item.LastModifyText ""}}
                                    <span title="{{i18n $.Lang "mgr.last_edit"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pencil"></i> {{i18n .Lang "mgr.last_edit"}}: {{$item.LastModifyText}}</span>
                                    {{end}}

                                </div>
                            </div>
                        {{else}}
                        <div class="text-center">{{i18n .Lang "message.no_data"}}</div>
                        {{end}}
                    </div>
                    <nav class="pagination-container">
                        {{.PageHtml}}
                    </nav>
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
                    <h4 class="modal-title">{{i18n .Lang "mgr.delete_project"}}</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">{{i18n .Lang "message.confirm_delete_project"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n .Lang "message.warning_delete_project"}}</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n .Lang "common.cancel"}}</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.confirm"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
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
                    showError({{i18n .Lang "message.system_error"}},"#form-error-message2");
                    $("#btnDeleteBook").button("reset");
                }
            });
        });
</script>
</body>
</html>