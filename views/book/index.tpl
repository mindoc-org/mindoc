<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n $.Lang "common.my_project"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/css/fileinput.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/themes/explorer-fa/theme.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/select2/4.0.5/css/select2.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
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
                    <li {{if eq .ControllerName "BookController"}}class="active"{{end}}><a href="{{urlfor "BookController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> {{i18n $.Lang "common.my_project"}}</a> </li>
                    <li {{if eq .ControllerName "BlogController"}}class="active"{{end}}><a href="{{urlfor "BlogController.ManageList"}}" class="item"><i class="fa fa-file" aria-hidden="true"></i> {{i18n $.Lang "common.my_blog"}}</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">{{i18n $.Lang "blog.project_list"}}</strong>
                        &nbsp;
                        <button type="button" data-toggle="modal" data-target="#addBookDialogModal" class="btn btn-success btn-sm pull-right">{{i18n $.Lang "blog.add_project"}}</button>
                        <button type="button" data-toggle="modal" data-target="#importBookDialogModal" class="btn btn-primary btn-sm pull-right" style="margin-right: 5px;">{{i18n $.Lang "blog.import_project"}}</button>
                    </div>
                </div>
                <div class="box-body" id="bookList">
                    <div class="book-list">
                        <template v-if="lists.length <= 0">
                        <div class="text-center">{{i18n $.Lang "message.no_data"}}</div>
                        </template>
                        <template v-else>

                        <div class="list-item" v-for="item in lists">
                            <div class="book-title">
                                <div class="pull-left">
                                    <a :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" title="{{i18n $.Lang "blog.project_summary"}}" data-toggle="tooltip">
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
                                    <div class="btn-group">
                                        <a  :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" class="btn btn-default">{{i18n $.Lang "common.setting"}}</a>

                                        <a href="javascript:;" class="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                            <span class="caret"></span>
                                            <span class="sr-only">Toggle Dropdown</span>
                                        </a>
                                        <ul class="dropdown-menu">
                                            <li><a :href="'{{urlfor "DocumentController.Index" ":key" ""}}' + item.identify" target="_blank">{{i18n $.Lang "blog.read"}}</a></li>
                                            <template v-if="item.role_id != 3">
                                            <li><a :href="'{{.BaseUrl}}/api/' + item.identify + '/edit'" target="_blank">{{i18n $.Lang "blog.edit"}}</a></li>
                                            </template>
                                            <template v-if="item.role_id == 0">
                                            <li><a :href="'javascript:deleteBook(\''+item.identify+'\');'">{{i18n $.Lang "blog.delete"}}</a></li>
                                            <li><a :href="'javascript:copyBook(\''+item.identify+'\');'">{{i18n $.Lang "blog.copy"}}</a></li>
                                            </template>
                                        </ul>

                                    </div>

                                    {{/*<a :href="'{{urlfor "DocumentController.Index" ":key" ""}}' + item.identify" title="{{i18n $.Lang "blog.view"}}" data-toggle="tooltip" target="_blank"><i class="fa fa-eye"></i> {{i18n $.Lang "blog.view"}}</a>*/}}
                                    {{/*<template v-if="item.role_id != 3">*/}}
                                        {{/*<a :href="'/api/' + item.identify + '/edit'" title="{{i18n $.Lang "blog.edit_doc"}}" data-toggle="tooltip" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> {{i18n $.Lang "blog.edit_doc"}}</a>*/}}
                                    {{/*</template>*/}}
                                </div>
                                <div class="clearfix"></div>
                            </div>
                            <div class="desc-text">
                                    <template v-if="item.description === ''">
                                        &nbsp;
                                    </template>
                                    <template v-else="">
                                        <a :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" title="{{i18n $.Lang "blog.project_summary"}}" style="font-size: 12px;">
                                        ${item.description}
                                        </a>
                                    </template>
                            </div>
                            <div class="info">
                                <span title="{{i18n $.Lang "blog.create_time"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-clock-o"></i>
                                    ${(new Date(item.create_time)).format("yyyy-MM-dd hh:mm:ss")}

                                </span>
                                <span title="{{i18n $.Lang "blog.creator"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user"></i> ${item.create_name}</span>
                                <span title="{{i18n $.Lang "blog.doc_amount"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pie-chart"></i> ${item.doc_count}</span>
                                <span title="{{i18n $.Lang "blog.project_role"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-user-secret"></i> ${item.role_name}</span>
                                <template v-if="item.last_modify_text !== ''">
                                    <span title="{{i18n $.Lang "blog.last_edit"}}" data-toggle="tooltip" data-placement="bottom"><i class="fa fa-pencil"></i> {{i18n $.Lang "blog.last_edit"}}: ${item.last_modify_text}</span>
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
    <div class="modal-dialog modal-lg" role="document" style="min-width: 900px;">
        <form method="post" autocomplete="off" action="{{urlfor "BookController.Create"}}" id="addBookDialogForm" enctype="multipart/form-data">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">{{i18n $.Lang "blog.add_project"}}</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <div class="pull-left" style="width: 620px">
                        <div class="form-group required">
                            <label class="text-label col-sm-2">{{i18n $.Lang "common.project_space"}}</label>
                            <div class="col-sm-10">
                                <select class="js-data-example-ajax-add form-control" multiple="multiple" name="itemId" id="itemId">
                                {{if .Item}}<option value="{{.Item.ItemId}}" selected>{{.Item.ItemName}}</option> {{end}}
                                </select>
                                <p class="text">{{i18n $.Lang "message.project_must_belong_space"}}</p>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                        <div class="form-group required">
                            <label class="text-label col-sm-2">{{i18n $.Lang "blog.project_title"}}</label>
                            <div class="col-sm-10">
                                <input type="text" class="form-control" placeholder="{{i18n $.Lang "message.project_title_placeholder"}}" name="book_name" id="bookName">
                                <p class="text">{{i18n $.Lang "message.project_title_tips"}}</p>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                        <div class="form-group required">
                           <label class="text-label col-sm-2">{{i18n $.Lang "blog.project_id"}}</label>
                            <div class="col-sm-10">
                                <input type="text" class="form-control" placeholder="{{i18n $.Lang "message.project_id_placeholder"}}" name="identify" id="identify">
                                <p class="text">{{i18n $.Lang "message.project_id_tips"}}</p>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                        <div class="form-group">
                            <textarea name="description" id="description" class="form-control" placeholder="{{i18n $.Lang "message.project_desc_placeholder"}}" style="height: 90px;"></textarea>
                        </div>
                        <div class="form-group">
                            <div class="col-lg-4">
                                <label>
                                    <input type="radio" name="privately_owned" value="0" checked> {{i18n $.Lang "blog.public"}}<span class="text">{{i18n $.Lang "message.project_public_desc"}}</span>
                                </label>
                            </div>
                            <div class="col-lg-8">
                                <label>
                                    <input type="radio" name="privately_owned" value="1"> {{i18n $.Lang "blog.private"}}<span class="text">{{i18n $.Lang "message.project_private_desc"}}</span>
                                </label>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                    </div>
                    <div class="pull-right text-center" style="width: 235px;">
                        <canvas id="bookCover" height="230px" width="170px"><img src="{{cdnimg "/static/images/book.jpg"}}"> </canvas>
                        <p class="text">{{i18n $.Lang "message.project_cover_desc"}}</p>
                    </div>
                </div>


                <div class="clearfix"></div>
            </div>
            <div class="modal-footer">
                <span id="form-error-message"></span>
                <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                <button type="button" class="btn btn-success" id="btnSaveDocument" data-loading-text="{{i18n $.Lang "common.processing"}}">{{i18n $.Lang "common.save"}}</button>
            </div>
        </div>
        </form>
    </div>
</div>
<!--END Modal-->
<!-- importBookDialogModal -->
<div class="modal fade" id="importBookDialogModal" tabindex="-1" role="dialog" aria-labelledby="importBookDialogModalLabel">
    <div class="modal-dialog" role="document" style="min-width: 900px;">
        <form method="post" autocomplete="off" action="{{urlfor "BookController.Import"}}" id="importBookDialogForm" enctype="multipart/form-data">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">{{i18n $.Lang "blog.import_project"}}</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <div class="form-group required">
                            <label class="text-label">{{i18n $.Lang "common.project_space"}}</label>
                            <select class="js-data-example-ajax-import form-control" multiple="multiple" name="itemId">
                                {{if .Item}}<option value="{{.Item.ItemId}}" selected>{{.Item.ItemName}}</option> {{end}}
                            </select>
                            <p class="text">{{i18n $.Lang "message.project_must_belong_space"}}</p>
                        </div>
                        <div class="form-group required">
                            <label class="text-label">{{i18n $.Lang "blog.project_title"}}</label>
                            <input type="text" class="form-control" placeholder="{{i18n $.Lang "message.project_title_placeholder"}}" name="book_name" maxlength="100" value="">
                            <p class="text">{{i18n $.Lang "blog.project_title_tips"}}</p>
                        </div>
                        <div class="form-group required">
                            <label class="text-label">{{i18n $.Lang "blog.project_id"}}</label>
                            <input type="text" class="form-control"  placeholder="{{i18n $.Lang "message.project_id_placeholder"}}" name="identify" value="">
                            <div class="clearfix"></div>
                            <p class="text">{{i18n $.Lang "blog.project_id_tips"}}</p>
                        </div>
                        <div class="form-group">
                            <label class="text-label">{{i18n $.Lang "blog.project_desc"}}</label>
                            <textarea name="description" id="description" class="form-control" placeholder="{{i18n $.Lang "message.project_desc_placeholder"}}" style="height: 90px;"></textarea>
                        </div>
                        <div class="form-group">
                            <div class="col-lg-4">
                                <label>
                                    <input type="radio" name="privately_owned" value="0" checked> {{i18n $.Lang "blog.public"}}<span class="text">{{i18n $.Lang "message.project_public_desc"}}</span>
                                </label>
                            </div>
                            <div class="col-lg-8">
                                <label>
                                    <input type="radio" name="privately_owned" value="1"> {{i18n $.Lang "blog.private"}}<span class="text">{{i18n $.Lang "message.project_private_desc"}}</span>
                                </label>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                        <div class="form-group">
                            <div class="file-loading">
                                <input id="import-book-upload" name="import-file" type="file" accept=".zip,.docx">
                            </div>
                            <div id="kartik-file-errors"></div>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="import-book-form-error-message" style="background-color: #ffffff;border: none;margin: 0;padding: 0;"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                    <button type="button" class="btn btn-success" id="btnImportBook" data-loading-text="{{i18n $.Lang "common.processing"}}">{{i18n $.Lang "common.create"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!--END importBookDialogModal-->
<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBookModal" tabindex="-1" role="dialog" aria-labelledby="deleteBookModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBookForm" action="{{urlfor "BookController.Delete"}}">
            <input type="hidden" name="identify" value="">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">{{i18n $.Lang "blog.delete_project"}}</h4>
                </div>
                <div class="modal-body">
                    <span style="font-size: 14px;font-weight: 400;">{{i18n $.Lang "message.confirm_delete_project"}}</span>
                    <p></p>
                    <p class="text error-message">{{i18n $.Lang "message.warning_delete_project"}}</p>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message2" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{i18n $.Lang "common.cancel"}}</button>
                    <button type="submit" id="btnDeleteBook" class="btn btn-primary" data-loading-text="{{i18n $.Lang "common.processing"}}">{{i18n $.Lang "common.confirm_delete"}}</button>
                </div>
            </div>
        </form>
    </div>
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/fileinput.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/locales/zh.js"}}"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/select2.full.min.js"}}"></script>
<script src="{{cdnjs "/static/select2/4.0.5/js/i18n/zh-CN.js"}}"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    /**
     * 绘制项目封面
     * @param $id
     * @param $font
     */
    function drawBookCover($id,$font) {
        var draw = document.getElementById($id);
        //确认浏览器是否支持<canvas>元素
        if (draw.getContext) {
            var context = draw.getContext('2d');

            //绘制红色矩形,绿色描边
            context.fillStyle = '#eee';
            context.strokeStyle = '#d4d4d5';
            context.strokeRect(0,0,170,230);
            context.fillRect(0,0,170,230);

            //设置字体样式
            context.font = "600 20px Helvetica";
            context.textAlign = "left";
            //设置字体填充颜色
            context.fillStyle = "#3E403E";

            var font = $.trim($font);

            var lineWidth = 0; //当前行的绘制的宽度
            var lastTextIndex = 0; //已经绘制上canvas最后的一个字符的下标
            var drawWidth = 155,lineHeight = 25,drawStartX = 15,drawStartY=65;
            //由于改变canvas 的高度会导致绘制的纹理被清空，所以，不预先绘制，先放入到一个数组当中
            var arr = [];
            for(var i = 0; i<font.length; i++){
                //获取当前的截取的字符串的宽度
                lineWidth = context.measureText(font.substr(lastTextIndex,i-lastTextIndex)).width;

                if(lineWidth > drawWidth){
                    //判断最后一位是否是标点符号
                    if(judgePunctuationMarks(font[i-1])){
                        arr.push(font.substr(lastTextIndex,i-lastTextIndex));
                        lastTextIndex = i;
                    }else{
                        arr.push(font.substr(lastTextIndex,i-lastTextIndex-1));
                        lastTextIndex = i-1;
                    }
                }
                //将最后多余的一部分添加到数组
                if(i === font.length - 1){
                    arr.push(font.substr(lastTextIndex,i-lastTextIndex+1));
                }
            }

            for(var i =0; i<arr.length; i++){
                context.fillText(arr[i],drawStartX,drawStartY+i*lineHeight);
            }

            //判断是否是需要避开的标签符号
            function judgePunctuationMarks(value) {
                var arr = [".",",",";","?","!",":","\"","，","。","？","！","；","：","、"];
                for(var i = 0; i< arr.length; i++){
                    if(value === arr[i]){
                        return true;
                    }
                }
                return false;
            }

        }else{
            console.log("浏览器不支持")
        }
    }

    /**
     * 将base64格式的图片转换为二进制
     * @param dataURI
     * @returns {Blob}
     */
    function dataURItoBlob(dataURI) {
        var byteString = atob(dataURI.split(',')[1]);
        var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0];
        var ab = new ArrayBuffer(byteString.length);
        var ia = new Uint8Array(ab);
        for (var i = 0; i < byteString.length; i++) {
            ia[i] = byteString.charCodeAt(i);
        }

        return new Blob([ab], {type: mimeString});
    }
    /**
     * 删除项目
     */
    function deleteBook($id) {
        $("#deleteBookModal").find("input[name='identify']").val($id);
        $("#deleteBookModal").modal("show");
    }
    /**
     * 复制项目
     * */
    function copyBook($id){
        var index = layer.load()
        $.ajax({
           url : "{{urlfor "BookController.Copy"}}" ,
            data : {"identify":$id},
            type : "POST",
            dataType : "json",
            success : function ($res) {
                layer.close(index);
                if ($res.errcode === 0) {
                    window.app.lists.splice(0, 0, $res.data);
                    $("#addBookDialogModal").modal("hide");
                } else {
                    layer.msg($res.message);
                }
            },
            error : function () {
                layer.close(index);
                layer.msg('{{i18n $.Lang "message.system_error"}}');
            }
        });
    }

    $(function () {
        /**
         * 处理创建项目弹窗
         * */
        $("#addBookDialogModal").on("show.bs.modal",function () {
            window.bookDialogModal = $(this).find("#addBookDialogForm").html();
            drawBookCover("bookCover","{{i18n $.Lang "blog.default_cover"}}");
            $('.js-data-example-ajax-add').select2({
                language: "{{i18n $.Lang "common.js_lang"}}",
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "BookController.ItemsetsSearch"}}',
                    dataType: 'json',
                    data: function (params) {
                        return {
                            q: params.term, // search term
                            page: params.page
                        };
                    },
                    processResults: function (data, params) {
                        return {
                            results : data.data.results
                        }
                    }
                }
            });
        }).on("hidden.bs.modal",function () {
            $(this).find("#addBookDialogForm").html(window.bookDialogModal);
        });
        /**
         * 处理导入项目弹窗
         * */
        $("#importBookDialogModal").on("show.bs.modal",function () {
            window.importBookDialogModal = $(this).find("#importBookDialogForm").html();
            $("#import-book-upload").fileinput({
                'uploadUrl':"{{urlfor "BookController.Import"}}",
                'theme': 'fa',
                'showPreview': false,
                'showUpload' : false,
                'required': true,
                'validateInitialCount': true,
                "language" : "{{i18n $.Lang "common.upload_lang"}}",
                'allowedFileExtensions': ['zip', 'docx'],
                'msgPlaceholder' : '{{i18n $.Lang "message.file_type_placeholder"}}',
                'elErrorContainer' : "#import-book-form-error-message",
                'uploadExtraData' : function () {
                    var book = {};
                    var $then = $("#importBookDialogForm");
                    book.book_name = $then.find("input[name='book_name']").val();
                    book.identify = $then.find("input[name='identify']").val();
                    book.description = $then.find('textarea[name="description"]').val();
                    book.itemId = $then.find("select[name='itemId']").val();

                    return book;
                }
            });
            $('.js-data-example-ajax-import').select2({
                language: "{{i18n $.Lang "common.js_lang"}}",
                minimumInputLength : 1,
                minimumResultsForSearch: Infinity,
                maximumSelectionLength:1,
                width : "100%",
                ajax: {
                    url: '{{urlfor "BookController.ItemsetsSearch"}}',
                    dataType: 'json',
                    data: function (params) {
                        return {
                            q: params.term, // search term
                            page: params.page
                        };
                    },
                    processResults: function (data, params) {
                        return {
                            results : data.data.results
                        }
                    }
                }
            });
        }).on("hidden.bs.modal",function () {
            $(this).find("#importBookDialogForm").html(window.importBookDialogModal);
        });

        /**
         * 创建项目
         */
        $("body").on("click","#btnSaveDocument",function () {
            var $this = $(this);

            var itemId = $("#itemId").val();
            if (itemId <= 0) {
                return showError("{{i18n $.Lang "message.project_space_empty"}}")
            }
            var bookName = $.trim($("#bookName").val());
            if (bookName === "") {
                return showError("{{i18n $.Lang "message.project_title_empty"}}")
            }
            if (bookName.length > 100) {
                return showError("{{i18n $.Lang "message.project_title_tips"}}");
            }

            var identify = $.trim($("#identify").val());
            if (identify === "") {
                return showError("{{i18n $.Lang "message.project_id_empty"}}");
            }
            if (identify.length > 50) {
                return showError("{{i18n $.Lang "message.project_id_length"}}");
            }
            var description = $.trim($("#description").val());

            if (description.length > 500) {
                return showError("{{i18n $.Lang "message.project_desc_placeholder"}}");
            }

            $this.button("loading");

            var draw = document.getElementById("bookCover");
            var form = document.getElementById("addBookDialogForm");
            var fd = new FormData(form);
            if (draw.getContext) {
                var dataURL = draw.toDataURL("png", 100);

                var blob = dataURItoBlob(dataURL);

                fd.append('image-file', blob,(new Date()).valueOf() + ".png");
            }

            $.ajax({
                url : "{{urlfor "BookController.Create"}}",
                data: fd,
                type: "POST",
                dataType :"json",
                processData: false,
                contentType: false
            }).success(function (res) {
                $this.button("reset");
                if (res.errcode === 0) {
                    window.app.lists.splice(0, 0, res.data);
                    $("#addBookDialogModal").modal("hide");
                } else {
                    showError(res.message);
                }
                $this.button("reset");
            }).error(function () {
                $this.button("reset");
                return showError("{{i18n $.Lang "message.system_error"}}");
            });
            return false;
        }).on("blur","#bookName",function () {
            var txt = $("#bookName").val();
            if(txt !== ""){
                drawBookCover("bookCover",txt);
            }
        }).on("click","#btnImportBook",function () {

            var $then = $(this).parents("#importBookDialogForm");

            var itemId = $then.find("input[name='itemId']").val();
            if (itemId <= 0) {
                return showError("{{i18n $.Lang "message.project_space_empty"}}")
            }

            var bookName = $.trim($then.find("input[name='book_name']").val());

            if (bookName === "") {
                return showError("{{i18n $.Lang "message.project_title_empty"}}","#import-book-form-error-message");
            }
            if (bookName.length > 100) {
                return showError("{{i18n $.Lang "message.project_title_tips"}}","#import-book-form-error-message");
            }

            var identify = $.trim($then.find("input[name='identify']").val());
            if (identify === "") {
                return showError("{{i18n $.Lang "message.project_id_empty"}}","#import-book-form-error-message");
            }
            var description = $.trim($then.find('textarea[name="description"]').val());
            if (description.length > 500) {
                return showError("{{i18n $.Lang "message.project_decs_placeholder"}}","#import-book-form-error-message");
            }
            var filesCount = $('#import-book-upload').fileinput('getFilesCount');

            if (filesCount <= 0) {
                return showError("{{i18n $.Lang "message.import_file_empty"}}","#import-book-form-error-message");
            }
            //$("#importBookDialogForm").submit();
            $("#btnImportBook").button("loading");
            $('#import-book-upload').fileinput('upload');
        }).on("fileuploaded","#import-book-upload",function (event, data, previewId, index){

            if(data.response.errcode === 0 || data.response.errcode === '0'){
                showSuccess(data.response.message,"#import-book-form-error-message");
            }else{
                showError(data.response.message,"#import-book-form-error-message");
            }
            $("#btnImportBook").button("reset");
            return true;
        });

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
                showError("{{i18n $.Lang "message.system_error"}}","#form-error-message2");
                $("#btnDeleteBook").button("reset");
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
