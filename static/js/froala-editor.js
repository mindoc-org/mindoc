$(function () {
    //超大屏幕
    var toolbarButtons = ['fullscreen', 'bold', 'italic', 'underline', 'strikeThrough', 'subscript', 'superscript', 'fontFamily', 'fontSize', '|', 'color', 'emoticons', 'inlineStyle', 'paragraphStyle', '|', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', 'insertHR', '-', 'insertLink', 'insertImage', 'insertVideo', 'insertFile', 'insertTable', 'undo', 'redo', 'clearFormatting', 'selectAll', 'html'];
    //大屏幕
    var toolbarButtonsMD = ['fullscreen', 'bold', 'italic', 'underline', 'fontFamily', 'fontSize', 'color', 'paragraphStyle', 'paragraphFormat', 'align', 'formatOL', 'formatUL', 'outdent', 'indent', 'quote', 'insertHR', 'insertLink', 'insertImage', 'insertVideo', 'insertFile', 'insertTable', 'undo', 'redo', 'clearFormatting'];
    //小屏幕
    var toolbarButtonsSM = ['fullscreen', 'bold', 'italic', 'underline', 'fontFamily', 'fontSize', 'insertLink', 'insertImage', 'insertTable', 'undo', 'redo'];
    //手机
    var toolbarButtonsXS = ['bold', 'italic', 'fontFamily', 'fontSize', 'undo', 'redo'];

    window.addDocumentModalFormHtml = $(this).find("form").html();
    window.editor = new FroalaEditor("#froalaEditor", {
      // enter: $.FroalaEditor.ENTER_P,
      placeholderText: '请输入内容',
      charCounterCount: true, //默认
      // charCounterMax         : -1,//默认
      saveInterval: 0, //不自动保存，默认10000
      // theme                    : "red",
      height: "100%",
      toolbarBottom: false, //默认
      toolbarButtonsMD: toolbarButtonsMD,
      toolbarButtonsSM: toolbarButtonsSM,
      toolbarButtonsXS: toolbarButtonsXS,
      toolbarInline: false, //true选中设置样式,默认false
      imageUploadMethod: 'POST',
      // heightMin: 450,
      charCounterMax: 3000,
      // imageUploadURL: "uploadImgEditor",
      // imageParams: { postId: "123" },
      params: {
        // acl: '01',
        // AWSAccessKeyId: '02',
        // policy: '03',
        // signature: '04',
        editor : "froalaEditor",
      },
      autosave: true,
      autosaveInterval: 2500,
      // saveURL: 'hander/FroalaHandler.ashx',
      saveParams: { postId: '1' },
      spellcheck: false,
      imageUploadURL: window.imageUploadURL, //'/uploadimg', //上传到本地服务器
      // imageUploadParams: { pid: '{{.product.ProjectId}}' },
      // imageDeleteURL: 'lib/delete_image.php', //删除图片
      // imagesLoadURL: 'lib/load_images.php', //管理图片
      videoUploadURL: window.imageUploadURL,
      // videoUploadParams: { pid: '{{.product.ProjectId}}' },
      fileUploadURL: window.imageUploadURL,
      // fileUploadParams: { pid: '{{.product.ProjectId}}' },
      // enter: $.FroalaEditor.ENTER_BR,
      language: 'zh_cn',
        // Add the custom buttons in the toolbarButtons list, after the separator.
      toolbarButtons: [toolbarButtons, ['saveIcon', 'insert','alert']]
      // toolbarButtons: ['bold', 'italic', 'underline', 'paragraphFormat', 'align','color','fontSize','insertImage','insertTable','undo', 'redo']
    });

    FroalaEditor.DefineIcon('alert', {NAME: 'info', SVG_KEY: 'user'});
    FroalaEditor.RegisterCommand('alert', {
      title: 'Hello',
      focus: false,
      undo: false,
      refreshAfterCallback: false,
      callback: function () {
        alert('Hello!');
      }
    });

    // FroalaEditor.DefineIcon('magicIcon', {NAME: 'magic'});
    FroalaEditor.DefineIcon('saveIcon', {NAME: 'plus', SVG_KEY: 'cogs'});
    FroalaEditor.RegisterCommand('saveIcon', {
      title: '保存',
      focus: false,
      undo: true,
      refreshAfterCallback: true,
      callback: function () {
        saveDocument();
      }
    });
    
    FroalaEditor.DefineIcon('insert', {NAME: 'plus', SVG_KEY: 'add'});
    FroalaEditor.RegisterCommand('insert', {
      title: '发布',
      focus: true,
      undo: true,
      refreshAfterCallback: true,
      callback: function () {
        // this.html.insert('My New HTML');
        // saveDocument(true,callback);
        releaseBook();
      }
    });
    
    // new FroalaEditor('div#froala-editor', {
    //   // Add the custom buttons in the toolbarButtons list, after the separator.
    //   toolbarButtons: [['undo', 'redo' , 'bold'], ['alert', 'clear', 'insert']]
    // })

    // editor.config.mapAk = window.baiduMapKey;
    // editor.config.printLog = false;
    // editor.config.showMenuTooltips = true
    // editor.config.menuTooltipPosition = 'down'
    // editor.config.uploadImgUrl = window.imageUploadURL;
    // editor.config.uploadImgFileName = "editormd-image-file";
    // editor.config.uploadParams = {
    //     "editor" : "froalaEditor"
    // };
    // editor.config.uploadImgServer = window.imageUploadURL;
    // editor.config.customUploadImg = function (resultFiles, insertImgFn) {
    //     // resultFiles 是 input 中选中的文件列表
    //     // insertImgFn 是获取图片 url 后，插入到编辑器的方法
    //     var file = resultFiles[0];
    //     // file type is only image.
    //     if (/^image\//.test(file.type)) {
    //         var form = new FormData();
    //         form.append('editormd-image-file', file, file.name);

    //         var layerIndex = 0;

    //         $.ajax({
    //             url: window.imageUploadURL,
    //             type: "POST",
    //             dataType: "json",
    //             data: form,
    //             processData: false,
    //             contentType: false,
    //             error: function() {
    //                 layer.close(layerIndex);
    //                 layer.msg("图片上传失败");
    //             },
    //             success: function(data) {
    //                 layer.close(layerIndex);
    //                 if(data.errcode !== 0){
    //                     layer.msg(data.message);
    //                 }else{
    //                     insertImgFn(data.url);
    //                 }
    //             }
    //         });
    //     } else {
    //         console.warn('You could only upload images.');
    //     }
    // };
    // editor.config.lang = window.lang;
    // editor.i18next = window.i18next;
    // editor.config.languages['en']['froalaEditor']['menus']['title']['保存'] = 'save';
    // editor.config.languages['en']['froalaEditor']['menus']['title']['发布'] = 'publish';
    // editor.config.languages['en']['froalaEditor']['menus']['title']['附件'] = 'attachment';
    // editor.config.languages['en']['froalaEditor']['menus']['title']['history'] = 'history';

    window.editormdLocales = {
        'zh-CN': {
            placeholder: '本编辑器支持 Markdown 编辑，左边编写，右边预览。',
            contentUnsaved: '编辑内容未保存，需要保存吗？',
            noDocNeedPublish: '没有需要发布的文档',
            loadDocFailed: '文档加载失败',
            fetchDocFailed: '获取当前文档信息失败',
            cannotAddToEmptyNode: '空节点不能添加内容',
            overrideModified: '文档已被其他人修改确定覆盖已存在的文档吗？',
            confirm: '确定',
            cancel: '取消',
            contentsNameEmpty: '目录名称不能为空',
            addDoc: '添加文档',
            edit: '编辑',
            delete: '删除',
            loadFailed: '加载失败请重试',
            tplNameEmpty: '模板名称不能为空',
            tplContentEmpty: '模板内容不能为空',
            saveSucc: '保存成功',
            serverExcept: '服务器异常',
            paramName: '参数名称',
            paramType: '参数类型',
            example: '示例值',
            remark: '备注',
        },
        'en': {
            placeholder: 'This editor supports Markdown editing, writing on the left and previewing on the right.',
            contentUnsaved: 'The edited content is not saved, need to save it?',
            noDocNeedPublish: 'No Document need to be publish',
            loadDocFailed: 'Load Document failed',
            fetchDocFailed: 'Fetch Document info failed',
            cannotAddToEmptyNode: 'Cannot add content to empty node',
            overrideModified: 'The document has been modified by someone else, are you sure to overwrite the document?',
            confirm: 'Confirm',
            cancel: 'Cancel',
            contentsNameEmpty: 'Document Name cannot be empty',
            addDoc: 'Add Document',
            edit: 'Edit',
            delete: 'Delete',
            loadFailed: 'Failed to load, please try again',
            tplNameEmpty: 'Template name cannot be empty',
            tplContentEmpty: 'Template content cannot be empty',
            saveSucc: 'Save success',
            serverExcept: 'Server Exception',
            paramName: 'Parameter',
            paramType: 'Type',
            example: 'Example',
            remark: 'Remark',
        }
    };

    // window.editor.config.onchange = function (newHtml) {
    //     var saveMenu = window.editor.menus.menuList.find((item) => item.key == 'save');
    //     // 判断内容是否改变
    //     if (window.source !== window.editor.txt.html()) {
    //         saveMenu.$elem.addClass('selected');
    //     } else {
    //         saveMenu.$elem.removeClass('selected');
    //     }
    // };

    // window.editor.create();

    // $("#froalaEditor").css("height","100%");

    if(window.documentCategory.length > 0){
        var item =  window.documentCategory[0];
        var $select_node = { node : {id : item.id}};
        loadDocument($select_node);
    }

    /***
     * 加载指定的文档到编辑器中
     * @param $node
     */
    function loadDocument($node) {
        var index = layer.load(1, {
            shade: [0.1,'#fff'] //0.1透明度的白色背景
        });

        $.get(window.editURL + $node.node.id ).done(function (res) {
            layer.close(index);

            if(res.errcode === 0){
                window.isLoad = true;

                // window.editor.txt.clear();
                // window.editor.txt.html(res.data.content);
                window.editor.html.set('');
                window.editor.events.focus();
                window.editor.html.set(res.data.content);
                // 将原始内容备份
                window.source = res.data.content;
                var node = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name,"identify" : res.data.identify,"version" : res.data.version};
                pushDocumentCategory(node);
                window.selectNode = node;

                pushVueLists(res.data.attach);

            }else{
                layer.msg("文档加载失败");
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("文档加载失败");
        });
    }

    /**
     * 保存文档到服务器
     * @param $is_cover 是否强制覆盖
     */
    function saveDocument($is_cover,callback) {
        var index = null;
        var node = window.selectNode;

        // var html = window.editor.txt.html() ;
        var html = window.editor.html.get();

        var content = "";
        if($.trim(html) !== ""){
            content = toMarkdown(html, { gfm: true });
        }
        var version = "";

        if (!node) {
            layer.msg(editormdLocales[lang].fetchDocFailed);
            return;
        }
        var doc_id = parseInt(node.id);

        for(var i in window.documentCategory){
            var item = window.documentCategory[i];

            if(item.id === doc_id){
                version = item.version;
                break;
            }
        }
        $.ajax({
            beforeSend  : function () {
                index = layer.load(1, {shade: [0.1,'#fff'] });
            },
            url :  window.editURL,
            data : {"identify" : window.book.identify,"doc_id" : doc_id,"markdown" : content,"html" : html,"cover" : $is_cover ? "yes":"no","version": version},
            type :"post",
            dataType :"json",
            success : function (res) {
                layer.close(index);
                if(res.errcode === 0){
                    for(var i in window.documentCategory){
                        var item = window.documentCategory[i];

                        if(item.id === doc_id){
                            window.documentCategory[i].version = res.data.version;
                            break;
                        }
                    }
                    // 更新内容备份
                    window.source = res.data.content;
                    // 触发编辑器 onchange 回调函数
                    // window.editor.config.onchange();
                    // window.editor.onchange();
                    // window.editor.events: {
                    //     'contentChanged': function () {
                    //       // Do something here.
                    //       // this is the editor instance.
                    //       console.log(this);
                    //     }
                    // }
                    if(typeof callback === "function"){
                        callback();
                    }
                }else if(res.errcode === 6005){
                    var confirmIndex = layer.confirm(editormdLocales[lang].overrideModified, {
                        btn: [editormdLocales[lang].confirm, editormdLocales[lang].cancel] // 按钮
                    }, function () {
                        layer.close(confirmIndex);
                        saveDocument(true,callback);
                    });
                }else{
                    layer.msg(res.message);
                }
            }
        });
    }



    /**
     * 添加顶级文档
     */
    $("#addDocumentForm").ajaxForm({
        beforeSubmit : function () {
            var doc_name = $.trim($("#documentName").val());
            if (doc_name === ""){
                return showError("目录名称不能为空","#add-error-message")
            }
            window.addDocumentFormIndex = layer.load(1, { shade: [0.1,'#fff']  });
            return true;
        },
        success : function (res) {
            if(res.errcode === 0){

                var data = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name,"identify" : res.data.identify,"version" : res.data.version};

                var node = window.treeCatalog.get_node(data.id);
                if(node){
                    window.treeCatalog.rename_node({"id":data.id},data.text);

                }else {
                    window.treeCatalog.create_node(data.parent, data);
                    window.treeCatalog.deselect_all();
                    window.treeCatalog.select_node(data);
                }
                pushDocumentCategory(data);
                $("#markdown-save").removeClass('change').addClass('disabled');
                $("#addDocumentModal").modal('hide');
            }else{
                showError(res.message,"#add-error-message")
            }
            layer.close(window.addDocumentFormIndex);
        }
    });

    /**
     * 文档目录树
     */
    $("#sidebar").jstree({
        'plugins': ["wholerow", "types", 'dnd', 'contextmenu'],
        "types": {
            "default": {
                "icon": false  // 删除默认图标
            }
        },
        'core': {
            'check_callback': true,
            "multiple": false,
            'animation': 0,
            "data": window.documentCategory
        },
        "contextmenu": {
            show_at_node: false,
            select_node: false,
            "items": {
                "添加文档": {
                    "separator_before": false,
                    "separator_after": true,
                    "_disabled": false,
                    "label": window.editormdLocales[window.lang].addDoc,//"添加文档",
                    "icon": "fa fa-plus",
                    "action": function (data) {

                        var inst = $.jstree.reference(data.reference),
                            node = inst.get_node(data.reference);

                        openCreateCatalogDialog(node);
                    }
                },
                "编辑": {
                    "separator_before": false,
                    "separator_after": true,
                    "_disabled": false,
                    "label": window.editormdLocales[window.lang].edit,
                    "icon": "fa fa-edit",
                    "action": function (data) {
                        var inst = $.jstree.reference(data.reference);
                        var node = inst.get_node(data.reference);
                        openEditCatalogDialog(node);
                    }
                },
                "删除": {
                    "separator_before": false,
                    "separator_after": true,
                    "_disabled": false,
                    "label": window.editormdLocales[window.lang].delete,
                    "icon": "fa fa-trash-o",
                    "action": function (data) {
                        var inst = $.jstree.reference(data.reference);
                        var node = inst.get_node(data.reference);
                        openDeleteDocumentDialog(node);
                    }
                }
            }
        }
    }).on('loaded.jstree', function () {
        window.treeCatalog = $(this).jstree();
    }).on('select_node.jstree', function (node, selected, event) {
        // if(window.editor.menus.menuList.find((item) => item.key == 'save').$elem.hasClass('selected')) {
        //      if (confirm(window.editormdLocales[window.lang].contentUnsaved)) {
        //         saveDocument(false,function () {
        //             loadDocument(selected);
        //         });
        //         return true;
        //     }
        // }
        loadDocument(selected);

    }).on("move_node.jstree", jstree_save);

    window.saveDocument = saveDocument;

    window.releaseBook = function () {
        // if(Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0){
        //     if(window.editor.menus.menuList.find((item) => item.key == 'save').$elem.hasClass('selected')) {
        //         if(confirm(editormdLocales[lang].contentUnsaved)) {
        //             saveDocument();
        //         }
        //     }
            locales = {
                'zh-CN': {
                    publishToQueue: '发布任务已推送到任务队列，稍后将在后台执行。',
                },
                'en': {
                    publishToQueue: 'The publish task has been pushed to the queue</br> and will be executed soon.',
                }
            }
            $.ajax({
                url: window.releaseURL,
                data: {"identify": window.book.identify},
                type: "post",
                dataType: "json",
                success: function (res) {
                    if (res.errcode === 0) {
                        layer.msg(locales[lang].publishToQueue);
                    } else {
                        layer.msg(res.message);
                    }
                }
            });
        // }else{
        //     layer.msg(editormdLocales[lang].noDocNeedPublish)
        // }
    };

    // $(window).resize(function(e) {
    //   var $container = $(editor.$textContainerElem.elems[0]);
    //   var $manual = $container.closest('.manual-froalaEditor');
    //   var maxHeight = $manual.closest('.manual-editor-container').innerHeight();
    //   var statusHeight = $manual.siblings('.manual-editor-status').outerHeight(true);
    //   var manualHeihgt = maxHeight - statusHeight;
    //   $manual.height(manualHeihgt);
    //   var toolbarHeight = $container.siblings('.w-e-toolbar').outerHeight(true);
    //   $container.height($container.parent().innerHeight() - toolbarHeight);
    // });
    // $(window).trigger('resize');
});