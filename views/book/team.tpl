<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>团队 - {{.Model.BookName}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">

    <style type="text/css">
        .table > tbody > tr > td {
            vertical-align: middle;
        }
    </style>
</head>
<body>
<div class="manual-reader">
{{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="row">
            <div class="page-left">
                <ul class="menu">
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i
                            class="fa fa-dashboard" aria-hidden="true"></i> 概要</a></li>
                    <li><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i
                            class="fa fa-user" aria-hidden="true"></i> 成员</a></li>
                    <li class="active"><a href="{{urlfor "BookController.Team" ":key" .Model.Identify}}" class="item"><i
                            class="fa fa-group" aria-hidden="true"></i> 团队</a></li>
                {{if eq .Model.RoleId 0 1}}
                    <li><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i
                            class="fa fa-gear" aria-hidden="true"></i> 设置</a></li>
                {{end}}
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 团队管理</strong>
                    {{if eq .Member.Role 0}}
                        <button type="button" class="btn btn-success btn-sm pull-right" data-toggle="modal"
                                data-target="#addTeamDialogModal"><i class="fa fa-user-plus" aria-hidden="true"></i>
                            添加团队
                        </button>
                    {{end}}
                    </div>
                </div>
                <div class="box-body">
                    <div class="users-list" id="teamList">
                        <template v-if="lists.length <= 0">
                            <div class="text-center">暂无数据</div>
                        </template>
                        <template v-else>
                            <table class="table">
                                <thead>
                                <tr>
                                    <th>团队名称</th>
                                    <th width="100">成员数量</th>
                                    <th width="180">加入时间</th>
                                    <th align="center" width="220px">操作</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="item in lists">
                                    <td>${item.team_name}</td>
                                    <td>${item.member_count}</td>
                                    <td>${(new Date(item.create_time)).format("yyyy-MM-dd hh:mm:ss")}</td>
                                    <td>
                                        <button type="button" class="btn btn-danger btn-sm" @click="deleteTeam(item.team_id,$event)" data-loading-text="删除中">删除</button>
                                    </td>
                                </tr>
                                </tbody>
                            </table>
                        </template>
                        <nav class="pagination-container">
                        {{.PageHtml}}
                        </nav>
                    </div>
                </div>
            </div>
        </div>
    </div>
{{template "widgets/footer.tpl" .}}
</div>

<div class="modal fade" id="addTeamDialogModal" tabindex="-1" role="dialog" aria-labelledby="addTeamDialogModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" autocomplete="off" id="addTeamDialogForm" class="form-horizontal" action="{{urlfor "BookController.TeamAdd"}}">
            <input type="hidden" name="bookId" value="{{.Model.BookId}}">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">添加团队</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label" for="account">团队名称<span class="error-message">*</span></label>
                        <div class="col-sm-10">
                            <select type="text" name="teamId" id="teamId" class="js-data-example-ajax form-control" multiple="multiple"></select>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-success" data-loading-text="保存中..." id="btnAddTeam">保存
                    </button>
                </div>
            </div>
        </form>
    </div>
</div><!--END Modal-->


<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {

        $("#addTeamDialogModal").on("show.bs.modal", function () {
            window.addTeamDialogModalHtml = $(this).find("form").html();
            $('.js-data-example-ajax').select2({
                language: "zh-CN",
                minimumInputLength: 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength: 1,
                width: "100%",
                ajax: {
                    url: '{{urlfor "BookController.TeamSearch" "identify" .Model.Identify}}',
                    dataType: 'json',
                    data: function (params) {
                        return {
                            q: params.term, // search term
                            page: params.page
                        };
                    },
                    processResults: function (data, params) {
                        return {
                            results: data.data.results
                        }
                    }
                }
            });
        }).on("hidden.bs.modal", function () {
            $(this).find("form").html(window.addTeamDialogModalHtml);
        }).on("shown.bs.modal", function () {
            $(this).find("input[name='teamId']").focus();
        });


        $("#addTeamDialogForm").ajaxForm({
            beforeSubmit: function () {
                var teamId = $.trim($("#addTeamDialogForm select[name='teamId']").val());
                if (teamId == "") {
                    return showError("团队名称不能为空");
                }
                $("#btnAddTeam").button("loading");
                return true;
            },
            success: function ($res) {
                if ($res.errcode === 0) {
                    app.lists.splice(0, 0, $res.data);
                    $("#addTeamDialogModal").modal("hide");
                } else {
                    showError($res.message);
                }
            },
            error: function () {
                showError("服务器异常");
            },
            complete: function () {
                $("#btnAddTeam").button("reset");
            }
        });

        var app = new Vue({
            el: "#teamList",
            data: {
                lists: {{.Result}}
            },
            delimiters: ['${', '}'],
            methods: {
                deleteTeam: function (id, e) {
                    var $this = this;
                    $.ajax({
                        url: "{{urlfor "BookController.TeamDelete"}}",
                        type: "post",
                        data: {"teamId": id, "identify": "{{.Model.Identify}}"},
                        dataType: "json",
                        success: function ($res) {
                            if ($res.errcode === 0) {
                                for (var index in $this.lists) {
                                    var item = $this.lists[index];
                                    if (item.team_relationship_id == id) {
                                        $this.lists.splice(index, 1);
                                        break;
                                    }
                                }
                            } else {
                                alert("操作失败：" + res.message);
                            }
                        }
                    });
                }
            }
        });
        Vue.nextTick(function () {
            $("[data-toggle='tooltip']").tooltip();
        });
    });
</script>
</body>
</html>