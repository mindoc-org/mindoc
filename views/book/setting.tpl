<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>设置 - {{.Model.BookName}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/cropper/2.3.4/cropper.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/tagsinput/bootstrap-tagsinput.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/bootstrap/plugins/bootstrap-switch/css/bootstrap3//bootstrap-switch.min.css"}}" rel="stylesheet">
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
                    <li><a href="{{urlfor "BookController.Dashboard" ":key" .Model.Identify}}" class="item"><i class="fa fa-dashboard" aria-hidden="true"></i> 概要</a> </li>
                    <li><a href="{{urlfor "BookController.Users" ":key" .Model.Identify}}" class="item"><i class="fa fa-users" aria-hidden="true"></i> 成员</a> </li>
                    <li class="active"><a href="{{urlfor "BookController.Setting" ":key" .Model.Identify}}" class="item"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a> </li>
                </ul>

            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> 项目设置</strong>
                        {{if eq .Model.RoleId 0}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#transferBookModal">转让项目</button>
                        {{if eq .Model.PrivatelyOwned 1}}
                        <button type="button"  class="btn btn-success btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">转为公有</button>
                        {{else}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" data-toggle="modal" data-target="#changePrivatelyOwnedModal" style="margin-right: 5px;">转为私有</button>
                        {{end}}
                        <button type="button"  class="btn btn-danger btn-sm pull-right" style="margin-right: 5px;" data-toggle="modal" data-target="#deleteBookModal">删除项目</button>
                        {{end}}

                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form method="post" id="bookEditForm" action="{{urlfor "BookController.SaveBook"}}">
                            <input type="hidden" name="identify" value="{{.Model.Identify}}">
                            <div class="form-group">
                                <label>标题</label>
                                <input type="text" class="form-control" name="book_name" id="bookName" placeholder="项目名称" value="{{.Model.BookName}}">
                            </div>
                            <div class="form-group">
                                <label>标识</label>
                                <input type="text" class="form-control" value="{{urlfor "DocumentController.Index" ":key" .Model.Identify}}" placeholder="项目唯一标识" disabled>
                            </div>
                            <div class="form-group">
                                <label>历史记录数量</label>
                                <input type="text" class="form-control" name="history_count" value="{{.Model.HistoryCount}}" placeholder="历史记录数量">
                                <p class="text">当开启文档历史时,该值会限制每个文档保存的历史数量</p>
                            </div>
                            <div class="form-group">
                                <label>公司标识</label>
                                <input type="text" class="form-control" name="publisher" value="{{.Model.Publisher}}" placeholder="公司名称">
                                <p class="text">导出文档PDF文档时显示的页脚</p>
                            </div>
                            <div class="form-group">
                                <label>描述</label>
                                <textarea rows="3" class="form-control" name="description" style="height: 90px" placeholder="项目描述">{{.Model.Description}}</textarea>
                                <p class="text">描述信息不超过500个字符,支持Markdown语法</p>
                            </div>
                            <div class="form-group">
                                <label>标签</label>
                                <input type="text" class="form-control" name="label" placeholder="项目标签" value="{{.Model.Label}}">
                                <p class="text">最多允许添加10个标签，多个标签请用“,”分割</p>
                            </div>
                            <div class="form-group">
                                <label>编辑器</label>
                                <div class="radio">
                                    <label class="radio-inline">
                                        <input type="radio"{{if eq .Model.Editor "markdown"}} checked{{end}} name="editor" value="markdown"> Markdown编辑器
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio"{{if eq .Model.Editor "html"}} checked{{end}} name="editor" value="html"> Html编辑器
                                    </label>
                                </div>
                            </div>
                            <!--
                            {{/*
                            <div class="form-group">
                            <label>开启评论</label>
                            <div class="radio">
                            <label class="radio-inline">
                            <input type="radio" {{if eq .Model.CommentStatus "open"}}checked{{end}} name="comment_status" value="open">允许所有人评论<span class="text"></span>
                            </label>
                            <label class="radio-inline">
                                <input type="radio" {{if eq .Model.CommentStatus "closed"}}checked{{end}} name="comment_status" value="closed">关闭评论<span class="text"></span>
                            </label>
                            <label class="radio-inline">
                                <input type="radio" {{if eq .Model.CommentStatus "group_only"}}checked{{end}} name="comment_status" value="group_only">仅允许参与者评论<span class="text"></span>
                            </label>
                            <label class="radio-inline">
                                <input type="radio" {{if eq .Model.CommentStatus "registered_only"}}checked{{end}} name="comment_status" value="registered_only">仅允许注册者评论<span class="text"></span>
                            </label>
                    </div>
                </div>
                */}}
                -->
                {{if eq .Model.PrivatelyOwned 1}}
                <div class="form-group">
                    <label>访问令牌</label>
                    <div class="row">
                        <div class="col-sm-10">
                            <input type="text" name="token" id="token" class="form-control" placeholder="访问令牌" readonly value="{{.Model.PrivateToken}}">
                        </div>
                        <div class="col-sm-2">
                            <button type="button" class="btn btn-success btn-sm" id="createToken" data-loading-text="生成" data-action="create">生成</button>
                            <button type="button" class="btn btn-danger btn-sm" id="deleteToken" data-loading-text="删除" data-action="delete">删除</button>
                        </div>
                    </div>
                </div>
                {{end}}

                <div class="form-group">
                    <label for="autoRelease">自动发布</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="autoRelease" name="auto_release"{{if .Model.AutoRelease }} checked{{end}} data-size="small">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">开启导出</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="isDownload" name="is_download"{{if .Model.IsDownload }} checked{{end}} data-size="small" placeholder="开启导出">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">开启分享</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="enableShare" name="enable_share"{{if .Model.IsEnableShare }} checked{{end}} data-size="small" placeholder="开启分享">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">设置第一篇文档为默认首页</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="isUseFirstDocument" name="is_use_first_document"{{if .Model.IsUseFirstDocument }} checked{{end}} data-size="small" placeholder="设置第一篇文档为默认首页">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <label for="autoRelease">自动保存</label>
                    <div class="controls">
                        <div class="switch switch-small" data-on="primary" data-off="info">
                            <input type="checkbox" id="autoSave" name="auto_save"{{if .Model.AutoSave }} checked{{end}} data-size="small" placeholder="自动保存">
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <button type="submit" id="btnSaveBookInfo" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                    <span id="form-error-message" class="error-message"></span>
                </div>
                </form>
            </div>
            <div class="form-right">
                <label>
                    <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                        <img src="{{cdnimg .Model.Cover}}" onerror="this.src='{{cdnimg "/static/images/book.png"}}'" alt="封面" style="max-width: 120px;border: 1px solid #999" id="headimgurl">
                    </a>
                </label>
            </div>
            <div class="clearfix"></div>

        </div>
    </div>
</div>
</div>
{{template "widgets/footer.tpl" .}}
</div>
<!-- Modal -->
<div class="modal fade" id="changePrivatelyOwnedModal" tabindex="-1" role="dialog" aria-labelledby="changePrivatelyOwnedModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" action="{{urlfor "BookController.PrivatelyOwned" }}" id="changePrivatelyOwnedForm">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <input type="hidden" name="status" value="{{if eq .Model.PrivatelyOwned 0}}close{{else}}open{{end}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title">
                        {{if eq .Model.PrivatelyOwned 0}}
                        转为私有
                        {{else}}
                        转为共有
                        {{end}}
                    </h4>
                </div>
                <div class="modal-body">
                    {{if eq .Model.PrivatelyOwned 0}}
                    <span style="font-size: 14px;font-weight: 400;">确定将项目转为私有吗？</span>
                    <p></p>
                    <p class="text error-message">转为私有后需要通过阅读令牌才能访问该项目。</p>
                    {{else}}
                    <span style="font-size: 14px;font-weight: 400;"> 确定将项目转为公有吗？</span>
                    <p></p>
                    <p class="text error-message">转为公有后所有人都可以访问该项目。</p>
                    {{end}}
                </div>
                <div class="modal-footer">
                    <span class="error-message" id="form-error-message1"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" class="btn btn-primary" data-loading-text="变更中..." id="btnChangePrivatelyOwned">确定</button>
                </div>
            </div>
        </form>
    </div>
</div>
<!-- Start Modal -->
<div class="modal fade" id="upload-logo-panel" tabindex="-1" role="dialog" aria-labelledby="修改封面" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title">修改封面</h4>
            </div>
            <div class="modal-body">
                <div class="wraper">
                    <div id="image-wraper">

                    </div>
                </div>
                <div class="watch-crop-list">
                    <div class="preview-title">预览</div>
                    <ul>
                        <li>
                            <div class="img-preview preview-lg"></div>
                        </li>
                        <li>
                            <div class="img-preview preview-sm"></div>
                        </li>
                    </ul>
                </div>
                <div style="clear: both"></div>
            </div>
            <div class="modal-footer">
                <span id="error-message"></span>
                <div id="filePicker" class="btn">选择</div>
                <button type="button" id="saveImage" class="btn btn-success" style="height: 40px;width: 77px;" data-loading-text="上传中...">上传</button>
            </div>
        </div>
    </div>
</div>
<!--END Modal-->

<!-- Delete Book Modal -->
<div class="modal fade" id="deleteBookModal" tabindex="-1" role="dialog" aria-labelledby="deleteBookModalLabel">
    <div class="modal-dialog" role="document">
        <form method="post" id="deleteBookForm" action="{{urlfor "BookController.Delete"}}">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
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
<!-- Modal -->
<div class="modal fade" id="transferBookModal" tabindex="-1" role="dialog" aria-labelledby="transferBookModalLabel">
    <div class="modal-dialog" role="document">
        <form action="{{urlfor "BookController.Transfer"}}" method="post" id="transferBookForm">
            <input type="hidden" name="identify" value="{{.Model.Identify}}">
            <div class="modal-content">
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <h4 class="modal-title" id="myModalLabel">项目转让</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">接收账号</label>
                        <div class="col-sm-10">
                            <input type="text" name="account" class="form-control" placeholder="接收者账号" id="receiveAccount" maxlength="50">
                        </div>
                    </div>
                    <div class="clearfix"></div>
                </div>
                <div class="modal-footer">
                    <span id="form-error-message3" class="error-message"></span>
                    <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                    <button type="submit" id="btnTransferBook" class="btn btn-primary">确定转让</button>
                </div>
            </div>
        </form>
    </div>
</div>

<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/cropper/2.3.4/cropper.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/tagsinput/bootstrap-tagsinput.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/plugins/bootstrap-switch/js/bootstrap-switch.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#upload-logo-panel").on("hidden.bs.modal",function () {
            $("#upload-logo-panel").find(".modal-body").html(window.modalHtml);
        }).on("show.bs.modal",function () {
            window.modalHtml = $("#upload-logo-panel").find(".modal-body").html();
        });
        $("#autoRelease,#enableShare,#isDownload,#isUseFirstDocument,#autoSave").bootstrapSwitch();

        $('input[name="label"]').tagsinput({
            confirmKeys: [13,44],
            maxTags: 10,
            trimValue: true,
            cancelConfirmKeysOnEmpty : false
        });

        $("#changePrivatelyOwnedForm").ajaxForm({
            beforeSubmit :function () {
                $("#btnChangePrivatelyOwned").button("loading");
            },
            success :function (res) {
                if(res.errcode === 0){
                    window.location = window.location.href;
                    return;
                }else{
                    showError(res.message,"#form-error-message1");
                }
                $("#btnChangePrivatelyOwned").button("reset");
            },
            error :function () {
                showError("服务器异常","#form-error-message1");
                $("#btnChangePrivatelyOwned").button("reset");
            }
        });

        $("#createToken,#deleteToken").on("click",function () {
            var btn = $(this).button("loading");
            var action = $(this).attr("data-action");
            $.ajax({
                url : "{{urlfor "BookController.CreateToken"}}",
                type :"post",
                data : { "identify" : {{.Model.Identify}} , "action" : action },
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        $("#token").val(res.data);
                    }else{
                        alert(res.message);
                    }
                    btn.button("reset");
                },
                error : function () {
                    btn.button("reset");
                    alert("服务器错误");
                }
            }) ;
        });
        $("#token").on("focus",function () {
            $(this).select();
        });
        $("#bookEditForm").ajaxForm({
            beforeSubmit : function () {
                var bookName = $.trim($("#bookName").val());
                if (bookName === "") {
                    return showError("项目名称不能为空");
                }
                $("#btnSaveBookInfo").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    showSuccess("保存成功")
                }else{
                    showError("保存失败")
                }
                $("#btnSaveBookInfo").button("reset");
            },
            error : function () {
                showError("服务错误");
                $("#btnSaveBookInfo").button("reset");
            }
        });
        $("#deleteBookForm").ajaxForm({
            beforeSubmit : function () {
                $("#btnDeleteBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    window.location = "{{urlfor "BookController.Index"}}";
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
        $("#transferBookForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#receiveAccount").val());
                if (account === ""){
                    return showError("接受者账号不能为空","#form-error-message3")
                }
                $("#btnTransferBook").button("loading");
            },
            success : function (res) {
                if(res.errcode === 0){
                    window.location = window.location.href;
                }else{
                    showError(res.message,"#form-error-message3");
                }
                $("#btnTransferBook").button("reset");
            },
            error : function () {
                $("#btnTransferBook").button("reset");
            }
        });

        try {
            var uploader = WebUploader.create({
                auto: false,
                swf: '{{.BaseUrl}}/static/webuploader/Uploader.swf',
                server: '{{urlfor "BookController.UploadCover"}}',
                formData : { "identify" : {{.Model.Identify}} },
                pick: "#filePicker",
                fileVal : "image-file",
                fileNumLimit : 1,
                compress : false,
                accept: {
                    title: 'Images',
                    extensions: 'jpg,jpeg,png',
                    mimeTypes: 'image/jpg,image/jpeg,image/png'
                }
            }).on("beforeFileQueued",function (file) {
                uploader.reset();
            }).on( 'fileQueued', function( file ) {
                uploader.makeThumb( file, function( error, src ) {
                    $img = '<img src="' + src +'" style="max-width: 360px;max-height: 360px;">';
                    if ( error ) {
                        $img.replaceWith('<span>不能预览</span>');
                        return;
                    }

                    $("#image-wraper").html($img);
                    window.ImageCropper = $('#image-wraper>img').cropper({
                        aspectRatio: 175 / 230,
                        dragMode : 'move',
                        viewMode : 1,
                        preview : ".img-preview"
                    });
                }, 1, 1 );
            }).on("uploadError",function (file,reason) {
                console.log(reason);
                $("#error-message").text("上传失败:" + reason);

            }).on("uploadSuccess",function (file, res) {

                if(res.errcode === 0){
                    console.log(res);
                    $("#upload-logo-panel").modal('hide');
                    $("#headimgurl").attr('src',res.data);
                }else{
                    $("#error-message").text(res.message);
                }
            }).on("beforeFileQueued",function (file) {
                if(file.size > 1024*1024*2){
                    uploader.removeFile(file);
                    uploader.reset();
                    alert("文件必须小于2MB");
                    return false;
                }
            }).on("uploadComplete",function () {
                $("#saveImage").button('reset');
            });
            $("#saveImage").on("click",function () {
                var files = uploader.getFiles();
                if(files.length > 0) {
                    $("#saveImage").button('loading');
                    var cropper = window.ImageCropper.cropper("getData");

                    uploader.option("formData", cropper);

                    uploader.upload();
                }else{
                    alert("请选择图片");
                }
            });
        }catch(e){
            console.log(e);
        }
    });
</script>
</body>
</html>