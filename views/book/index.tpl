<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>我的项目 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">

    <link href="{{cdncss "/static/css/main.css"}}" rel="stylesheet">
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
                        <button type="button" data-toggle="modal" data-target="#addBookDialogModal" class="btn btn-success btn-sm pull-right">添加项目</button>
                    </div>
                </div>
                <div class="box-body" id="bookList">
                    <div class="book-list">
                        <template v-if="lists.length <= 0">
                        <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>

                        <div class="list-item" v-for="item in lists">
                            <div class="book-title">
                                <div class="pull-left">
                                    <a :href="'/book/' + item.identify + '/dashboard'" title="项目概要" data-toggle="tooltip">
                                       <template v-if="item.privately_owned == 0">
                                           <i class="fa fa-unlock" aria-hidden="true"></i>
                                       </template>
                                       <template v-else-if="item.privately_owned == 1">
                                           <i class="fa fa-lock" aria-hidden="true"></i>
                                       </template>

                                        ${item.book_name}
                                    </a>
                                </div>
                                <div class="pull-right">
                                    <a :href="'{{urlfor "DocumentController.Index" ":key" ""}}' + item.identify" title="查看文档" data-toggle="tooltip" target="_blank"><i class="fa fa-eye"></i> 查看文档</a>
                                    <template v-if="item.role_id != 3">
                                        <a :href="'/api/' + item.identify + '/edit'" title="编辑文档" data-toggle="tooltip" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑文档</a>
                                    </template>
                                </div>
                                <div class="clearfix"></div>
                            </div>
                            <div class="desc-text">
                                    <template v-if="item.description === ''">
                                        &nbsp;
                                    </template>
                                    <template v-else="">
                                        <a :href="'/book/' + item.identify + '/dashboard'" title="项目概要" style="font-size: 12px;">
                                        ${item.description}
                                        </a>
                                    </template>
                            </div>
                            <div class="info">
                                <span title="创建时间" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-clock-o"></i>
                                    ${(new Date(item.create_time)).format("yyyy-MM-dd hh:mm:ss")}

                                </span>
                                <span title="创建者" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user"></i> ${item.create_name}</span>
                                <span title="文档数量" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pie-chart"></i> ${item.doc_count}</span>
                                <span title="项目角色" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user-secret"></i> ${item.role_name}</span>
                                <template v-if="item.last_modify_text !== ''">
                                    <span title="最后编辑" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pencil"></i> 最后编辑: ${item.last_modify_text}</span>
                                </template>

                            </div>
                        </div>
                        </template>
                    </div>
                    <template v-if="lists.length >= 0">
                        <nav class="pagination-container">
                            {{.PageHtml}}
                        </nav>
                    </template>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="addBookDialogModal" tabindex="-1" role="dialog" aria-labelledby="addBookDialogModalLabel">
    <div class="modal-dialog" role="document" style="width: 655px">
        <form method="post" autocomplete="off" action="{{urlfor "BookController.Create"}}" id="addBookDialogForm">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">添加项目</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <input type="text" class="form-control" placeholder="标题(不超过100字)" name="book_name" id="bookName">
                </div>
                <div class="form-group">
                    <div class="pull-left" style="padding: 7px 5px 6px 0">
                        {{.BaseUrl}}{{urlfor "DocumentController.Index" ":key" ""}}
                    </div>
                    <input type="text" class="form-control pull-left" style="width: 220px;vertical-align: middle" placeholder="项目唯一标识(不能超过50字)" name="identify" id="identify">
                    <div class="clearfix"></div>
                    <p class="text" style="font-size: 12px;color: #999;margin-top: 6px;">文档标识只能包含小写字母、数字，以及“-”和“_”符号,并且只能小写字母开头</p>

                </div>
                <div class="form-group">
                    <textarea name="description" id="description" class="form-control" placeholder="描述信息不超过500个字符" style="height: 90px;"></textarea>
                </div>
                <div class="form-group">
                    <div class="col-lg-6">
                        <label>
                            <input type="radio" name="privately_owned" value="0" checked> 公开<span class="text">(任何人都可以访问)</span>
                        </label>
                    </div>
                    <div class="col-lg-6">
                        <label>
                            <input type="radio" name="privately_owned" value="1"> 私有<span class="text">(只要参与者或使用令牌才能访问)</span>
                        </label>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <!--
				{{/*
                <div class="form-group">
                    <div class="col-lg-3">
                        <label>
                            <input type="radio" checked name="comment_status" value="open">允许所有人评论<span class="text"></span>
                        </label>
                    </div>
                    <div class="col-lg-3">
                        <label>
                            <input type="radio"  name="comment_status" value="closed">关闭评论<span class="text"></span>
                        </label>
                    </div>
                    <div class="col-lg-3">
                        <label>
                            <input type="radio" name="comment_status" value="group_only">仅允许参与者评论<span class="text"></span>
                        </label>
                    </div>
                    <div class="col-lg-3">
                        <label>
                            <input type="radio" name="comment_status" value="registered_only">仅允许注册者评论<span class="text"></span>
                        </label>
                    </div>
                    <div class="clearfix"></div>
                </div>
				*/}}
				-->
                <div class="clearfix"></div>
            </div>
            <div class="modal-footer">
                <span id="form-error-message"></span>
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                <button type="submit" class="btn btn-success" id="btnSaveDocument" data-loading-text="保存中...">保存</button>
            </div>
        </div>
        </form>
    </div>
</div><!--END Modal-->

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="/static/js/main.js" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#addBookDialogForm").ajaxForm({
            beforeSubmit : function () {
                var bookName = $.trim($("#bookName").val());
                if(bookName === ""){
                    return showError("项目标题不能为空")
                }
                if(bookName.length > 100){
                    return showError("项目标题必须小于100字符");
                }

                var identify = $.trim($("#identify").val());
                if(identify === ""){
                    return showError("项目标识不能为空");
                }
                if(identify.length > 50){
                    return showError("项目标识必须小于50字符");
                }
                var description = $.trim($("#description").val());

                if(description.length > 500){
                    return showError("描述信息不超过500个字符");
                }
                $("#btnSaveDocument").button("loading");
                return showSuccess("");
            },
            success : function (res) {
                $("#btnSaveDocument").button("reset");
                if(res.errcode === 0){
                    window.app.lists.splice(0,0,res.data);
                    $("#addBookDialogModal").modal("hide");
                }else{
                    showError(res.message);
                }

            }
        });

        window.app = new Vue({
            el : "#bookList",
            data : {
                lists : {{.Result}}
            },
            delimiters : ['${','}'],
            methods : {
            }
        });
        Vue.nextTick(function () {
            $("[data-toggle='tooltip']").tooltip();
        });
    });
</script>
</body>
</html>