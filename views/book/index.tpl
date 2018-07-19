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
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/css/fileinput.min.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/themes/explorer-fa/theme.css"}}" rel="stylesheet" type="text/css">
    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
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
                    <li {{if eq .ControllerName "BookController"}}class="active"{{end}}><a href="{{urlfor "BookController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 我的项目</a> </li>
                    <li {{if eq .ControllerName "BlogController"}}class="active"{{end}}><a href="{{urlfor "BlogController.ManageList"}}" class="item"><i class="fa fa-file" aria-hidden="true"></i> 我的文章</a> </li>
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">项目列表</strong>
                        &nbsp;
                        <button type="button" data-toggle="modal" data-target="#addBookDialogModal" class="btn btn-success btn-sm pull-right">添加项目</button>
                        <button type="button" data-toggle="modal" data-target="#importBookDialogModal" class="btn btn-primary btn-sm pull-right" style="margin-right: 5px;">导入项目</button>
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
                                    <a :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" title="项目概要" data-toggle="tooltip">
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
                                        <a  :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" class="btn btn-default">设置</a>

                                        <a href="javascript:;" class="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                            <span class="caret"></span>
                                            <span class="sr-only">Toggle Dropdown</span>
                                        </a>
                                        <ul class="dropdown-menu">
                                            <li><a :href="'{{urlfor "DocumentController.Index" ":key" ""}}' + item.identify" target="_blank">阅读</a></li>
                                            <template v-if="item.role_id != 3">
                                            <li><a :href="'{{.BaseUrl}}/api/' + item.identify + '/edit'" target="_blank">编辑</a></li>
                                            </template>
                                            <template v-if="item.role_id == 0">
                                            <li><a :href="'javascript:deleteBook(\''+item.identify+'\');'">删除</a></li>
                                            <li><a :href="'javascript:copyBook(\''+item.identify+'\');'">复制</a></li>
                                            </template>
                                        </ul>

                                    </div>

                                    {{/*<a :href="'{{urlfor "DocumentController.Index" ":key" ""}}' + item.identify" title="查看文档" data-toggle="tooltip" target="_blank"><i class="fa fa-eye"></i> 查看文档</a>*/}}
                                    {{/*<template v-if="item.role_id != 3">*/}}
                                        {{/*<a :href="'/api/' + item.identify + '/edit'" title="编辑文档" data-toggle="tooltip" target="_blank"><i class="fa fa-edit" aria-hidden="true"></i> 编辑文档</a>*/}}
                                    {{/*</template>*/}}
                                </div>
                                <div class="clearfix"></div>
                            </div>
                            <div class="desc-text">
                                    <template v-if="item.description === ''">
                                        &nbsp;
                                    </template>
                                    <template v-else="">
                                        <a :href="'{{.BaseUrl}}/book/' + item.identify + '/dashboard'" title="项目概要" style="font-size: 12px;">
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
    <div class="modal-dialog modal-lg" role="document" style="min-width: 900px;">
        <form method="post" autocomplete="off" action="{{urlfor "BookController.Create"}}" id="addBookDialogForm" enctype="multipart/form-data">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">添加项目</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <div class="pull-left" style="width: 620px">
                        <div class="form-group">
                            <input type="text" class="form-control" placeholder="标题(不超过100字)" name="book_name" id="bookName">
                        </div>
                        <div class="form-group">
                            <div class="pull-left" style="padding: 7px 5px 6px 0">
                           {{urlfor "DocumentController.Index" ":key" ""}}
                            </div>
                            <input type="text" class="form-control pull-left" style="width: 410px;vertical-align: middle" placeholder="项目唯一标识(不超过50字)" name="identify" id="identify">
                            <div class="clearfix"></div>
                            <p class="text" style="font-size: 12px;color: #999;margin-top: 6px;">文档标识只能包含小写字母、数字，以及“-”、“.”和“_”符号.</p>
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
                    </div>
                    <div class="pull-right text-center" style="width: 235px;">
                        <canvas id="bookCover" height="230px" width="170px"><img src="{{cdnimg "/static/images/book.jpg"}}"> </canvas>
                    </div>
                </div>


                <div class="clearfix"></div>
            </div>
            <div class="modal-footer">
                <span id="form-error-message"></span>
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                <button type="button" class="btn btn-success" id="btnSaveDocument" data-loading-text="保存中...">保存</button>
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
                    <h4 class="modal-title">导入项目</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <div class="form-group required">
                            <label class="text-label">项目标题</label>
                            <input type="text" class="form-control" placeholder="项目标题(不超过100字)" name="book_name" maxlength="100" value="">
                        </div>
                        <div class="form-group required">
                            <label class="text-label">项目标识</label>
                            <input type="text" class="form-control"  placeholder="项目唯一标识(不超过50字)" name="identify" value="">
                            <div class="clearfix"></div>
                            <p class="text" style="font-size: 12px;color: #999;margin-top: 6px;">文档标识只能包含小写字母、数字，以及“-”、“.”和“_”符号.</p>
                        </div>
                        <div class="form-group">
                            <label class="text-label">项目描述</label>
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
                        <div class="form-group">
                            <div class="file-loading">
                                <input id="import-book-upload" name="import-file" type="file" accept=".zip">
                            </div>
                            <div id="kartik-file-errors"></div>
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="import-book-form-error-message" style="background-color: #ffffff;border: none;margin: 0;padding: 0;"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="button" class="btn btn-success" id="btnImportBook" data-loading-text="创建中...">创建</button>
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

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/vuejs/vue.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/fileinput.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-fileinput/4.4.7/js/locales/zh.js"}}"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript" ></script>
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
            context.font = "bold 20px SimSun";
            context.textAlign = "left";
            //设置字体填充颜色
            context.fillStyle = "#3E403E";

            var font = $font;

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
                layer.msg('服务器异常');
            }
        });
    }

    $(function () {
        $("#addBookDialogModal").on("show.bs.modal",function () {
            window.bookDialogModal = $(this).find("#addBookDialogForm").html();
            drawBookCover("bookCover","默认封面");
        }).on("hidden.bs.modal",function () {
            $(this).find("#addBookDialogForm").html(window.bookDialogModal);
        });
        /**
         * 创建项目
         */
        $("body").on("click","#btnSaveDocument",function () {
            var $this = $(this);


            var bookName = $.trim($("#bookName").val());
            if (bookName === "") {
                return showError("项目标题不能为空")
            }
            if (bookName.length > 100) {
                return showError("项目标题必须小于100字符");
            }

            var identify = $.trim($("#identify").val());
            if (identify === "") {
                return showError("项目标识不能为空");
            }
            if (identify.length > 50) {
                return showError("项目标识必须小于50字符");
            }
            var description = $.trim($("#description").val());

            if (description.length > 500) {
                return showError("描述信息不超过500个字符");
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
                return showError("服务器异常");
            });
            return false;
        });
        /**
         * 当填写项目标题后，绘制项目封面
         */
        $("#bookName").on("blur",function () {
           var txt = $(this).val();
           if(txt !== ""){
               drawBookCover("bookCover",txt);
           }
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
                showError("服务器异常","#form-error-message2");
                $("#btnDeleteBook").button("reset");
            }
        });

        $("#btnImportBook").on("click",function () {
            var $this = $(this);
            var $then = $(this).parents("#importBookDialogForm");


            var bookName = $.trim($then.find("input[name='book_name']").val());

            if (bookName === "") {
                return showError("项目标题不能为空","#import-book-form-error-message");
            }
            if (bookName.length > 100) {
                return showError("项目标题必须小于100字符","#import-book-form-error-message");
            }

            var identify = $.trim($then.find("input[name='identify']").val());
            if (identify === "") {
                return showError("项目标识不能为空","#import-book-form-error-message");
            }
            var description = $.trim($then.find('textarea[name="description"]').val());
            if (description.length > 500) {
                return showError("描述信息不超过500个字符","#import-book-form-error-message");
            }
            var filesCount = $('#import-book-upload').fileinput('getFilesCount');
            console.log(filesCount)
            if (filesCount <= 0) {
                return showError("请选择需要上传的文件","#import-book-form-error-message");
            }
           //$("#importBookDialogForm").submit();
            $("#btnImportBook").button("loading");
            $('#import-book-upload').fileinput('upload');

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

        $("#import-book-upload").fileinput({
            'uploadUrl':"{{urlfor "BookController.Import"}}",
            'theme': 'fa',
            'showPreview': false,
            'showUpload' : false,
            'required': true,
            'validateInitialCount': true,
            "language" : "zh",
            'allowedFileExtensions': ['zip'],
            'msgPlaceholder' : '请选择Zip文件',
            'elErrorContainer' : "#import-book-form-error-message",
            'uploadExtraData' : function () {
                var book = {};
                var $then = $("#importBookDialogForm");
                book.book_name = $then.find("input[name='book_name']").val();
                book.identify = $then.find("input[name='identify']").val();
                book.description = $then.find('textarea[name="description"]').val()

                return book;
            }
        });
        $("#import-book-upload").on("fileuploaded",function (event, data, previewId, index){

            if(data.response.errcode === 0 || data.response.errcode === '0'){
                showSuccess(data.response.message,"#import-book-form-error-message");
            }else{
                showError(data.response.message,"#import-book-form-error-message");
            }
            $("#btnImportBook").button("reset");
            return true;
        });
    });
</script>
</body>
</html>
