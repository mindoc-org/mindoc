<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>用户中心 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/webuploader/webuploader.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/cropper/2.3.4/cropper.min.css"}}" rel="stylesheet">
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
                    <li class="active"><a href="{{urlfor "SettingController.Index"}}" class="item"><i class="fa fa-sitemap" aria-hidden="true"></i> 基本信息</a> </li>
                    {{if ne .Member.AuthMethod "ldap"}}
                    <li><a href="{{urlfor "SettingController.Password"}}" class="item"><i class="fa fa-user" aria-hidden="true"></i> 修改密码</a> </li>
                    {{end}}
                </ul>
            </div>
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title">基本信息</strong>
                    </div>
                </div>
                <div class="box-body" style="padding-right: 200px;">
                    <div class="form-left">
                        <form role="form" method="post" id="memberInfoForm">
                            <div class="form-group">
                                <label>用户名</label>
                                <input type="text" class="form-control disabled" value="{{.Member.Account}}" disabled placeholder="用户名">
                            </div>
                            <div class="form-group">
                                <label>真实姓名</label>
                                <input type="text" name="real_name" class="form-control" value="{{.Member.RealName}}" placeholder="真实姓名">
                            </div>
                            <div class="form-group">
                                <label for="user-email">邮箱<strong class="text-danger">*</strong></label>
                                <input type="email" class="form-control" value="{{.Member.Email}}" id="userEmail" name="email" max="100" placeholder="邮箱">
                            </div>
                            <div class="form-group">
                                <label>手机号</label>
                                <input type="text" class="form-control" id="userPhone" name="phone" maxlength="20" title="手机号码" placeholder="手机号码" value="{{.Member.Phone}}">
                            </div>
                            <div class="form-group">
                                <label class="description">描述</label>
                                <textarea class="form-control" rows="3" title="描述" name="description" id="description" maxlength="500">{{.Member.Description}}</textarea>
                                <p style="color: #999;font-size: 12px;">描述不能超过500字</p>
                            </div>
                            <div class="form-group">
                                <button type="submit" class="btn btn-success" data-loading-text="保存中...">保存修改</button>
                                <span id="form-error-message" class="error-message"></span>
                            </div>
                        </form>
                    </div>
                    <div class="form-right">
                        <label>
                            <a href="javascript:;" data-toggle="modal" data-target="#upload-logo-panel">
                                <img src="{{.Member.Avatar}}" onerror="this.src='/static/images/middle.gif'" class="img-circle" alt="头像" style="max-width: 120px;max-height: 120px;" id="headimgurl">
                            </a>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<!-- Start Modal -->
<div class="modal fade" id="upload-logo-panel" tabindex="-1" role="dialog" aria-labelledby="修改头像" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 class="modal-title">修改头像</h4>
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
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/webuploader/webuploader.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/cropper/2.3.4/cropper.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#upload-logo-panel").on("hidden.bs.modal",function () {
            $("#upload-logo-panel").find(".modal-body").html(window.modalHtml);
        }).on("show.bs.modal",function () {
            window.modalHtml = $("#upload-logo-panel").find(".modal-body").html();
        });


        $("#memberInfoForm").ajaxForm({
            beforeSubmit : function () {

                var email = $.trim($("#userEmail").val());
                if(!email){
                    return showError('邮箱不能为空');
                }
                $("button[type='submit']").button('loading');
            },
            success : function (res) {
                $("button[type='submit']").button('reset');
                if(res.errcode === 0){
                    showSuccess("保存成功");
                }else{
                    showError(res.message);
                }
            }
        });

        try {
            var uploader = WebUploader.create({
                auto: false,
                swf: '/static/webuploader/Uploader.swf',
                server: '{{urlfor "SettingController.Upload"}}',
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
                        aspectRatio: 1 / 1,
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
                    alert("请选择头像");
                }
            });
        }catch(e){
            console.log(e);
        }
    });
</script>
</body>
</html>