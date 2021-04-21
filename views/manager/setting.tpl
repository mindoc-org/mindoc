<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{i18n .Lang "mgr.config_mgr"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
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
        {{template "manager/widgets.tpl" .}}
            <div class="page-right">
                <div class="m-box">
                    <div class="box-head">
                        <strong class="box-title"> {{i18n .Lang "mgr.config_mgr"}}</strong>
                    </div>
                </div>
                <div class="box-body">
                    <form method="post" id="gloablEditForm" action="{{urlfor "ManagerController.Setting"}}">
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.site_name"}}</label>
                            <input type="text" class="form-control" name="SITE_NAME" id="siteName" placeholder="{{i18n .Lang "mgr.site_name"}}" value="{{.SITE_NAME}}">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.domain_icp"}}</label>
                            <input type="text" class="form-control" name="site_beian" id="siteName" placeholder="{{i18n .Lang "mgr.domain_icp"}}" value="{{.site_beian}}" maxlength="50">
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.site_desc"}}</label>
                            <textarea rows="3" class="form-control" name="site_description" style="height: 90px" placeholder="{{i18n .Lang "mgr.site_desc"}}">{{.site_description}}</textarea>
                            <p class="text">{{i18n .Lang "mgr.site_desc_tips"}}</p>
                        </div>
                            <div class="form-group">
                                <label>{{i18n .Lang "mgr.enable_anonymous_access"}}</label>
                                <div class="radio">
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .ENABLE_ANONYMOUS "true"}}checked{{end}} name="ENABLE_ANONYMOUS" value="true">{{i18n .Lang "mgr.enable"}}<span class="text"></span>
                                    </label>
                                    <label class="radio-inline">
                                        <input type="radio" {{if eq .ENABLE_ANONYMOUS "false"}}checked{{end}} name="ENABLE_ANONYMOUS" value="false">{{i18n .Lang "mgr.disable"}}<span class="text"></span>
                                    </label>
                                </div>
                            </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.enable_register"}}</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_REGISTER "true"}}checked{{end}} name="ENABLED_REGISTER" value="true">{{i18n .Lang "mgr.enable"}}<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_REGISTER "false"}}checked{{end}} name="ENABLED_REGISTER" value="false">{{i18n .Lang "mgr.disable"}}<span class="text"></span>
                                </label>
                            </div>
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.enable_captcha"}}</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_CAPTCHA "true"}}checked{{end}} name="ENABLED_CAPTCHA" value="true">{{i18n .Lang "mgr.enable"}}<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLED_CAPTCHA "false"}}checked{{end}} name="ENABLED_CAPTCHA" value="false">{{i18n .Lang "mgr.disable"}}<span class="text"></span>
                                </label>
                            </div>
                        </div>
                        <div class="form-group">
                            <label>{{i18n .Lang "mgr.enable_doc_his"}}</label>
                            <div class="radio">
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLE_DOCUMENT_HISTORY "true"}}checked{{end}} name="ENABLE_DOCUMENT_HISTORY" value="true">{{i18n .Lang "mgr.enable"}}<span class="text"></span>
                                </label>
                                <label class="radio-inline">
                                    <input type="radio" {{if eq .ENABLE_DOCUMENT_HISTORY "false"}}checked{{end}} name="ENABLE_DOCUMENT_HISTORY" value="false">{{i18n .Lang "mgr.disable"}}<span class="text"></span>
                                </label>
                            </div>
                        </div>

                        <div class="form-group">
                            <button type="submit" id="btnSaveBookInfo" class="btn btn-success" data-loading-text="{{i18n .Lang "message.processing"}}">{{i18n .Lang "common.save"}}</button>
                            <span id="form-error-message" class="error-message"></span>
                        </div>
                        </form>

                    <div class="clearfix"></div>

                </div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>


<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/main.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#gloablEditForm").ajaxForm({
            beforeSubmit : function () {
                var title = $.trim($("#siteName").val());

                if (title === ""){
                    return showError({{i18n .Lang "message.site_name_empty"}});
                }
                $("#btnSaveBookInfo").button("loading");
            },success : function (res) {
                if(res.errcode === 0) {
                    showSuccess({{i18n .Lang "message.success"}})
                }else{
                    showError(res.message);
                }
                $("#btnSaveBookInfo").button("reset");
            }
        });
    });
</script>
</body>
</html>