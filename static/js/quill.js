$(function () {
    window.addDocumentModalFormHtml = $(this).find("form").html();
    window.menu_save = $("#markdown-save");
    window.uploader = null;
    window.editor = new Quill('#docEditor', {
        theme: 'snow',
        syntax: true,
        modules : {
            toolbar :"#editormd-tools"
        }
    });
    window.editor.on("text-change",function () {
        resetEditorChanged(true);
    });
    var $editorEle =  $("#editormd-tools");

    $editorEle.find(".ql-undo").on("click",function () {
        window.editor.history.undo();
    });
    $editorEle.find(".ql-redo").on("click",function () {
        window.editor.history.redo();
    });

    uploadImage("docEditor", function ($state, $res) {
        if ($state === "before") {
            return layer.load(1, {
                shade: [0.1, '#fff'] // 0.1 透明度的白色背景
            });
        } else if ($state === "success") {
            if ($res.errcode === 0) {

                var range = window.editor.getSelection();
                window.editor.insertEmbed(range.index, 'image', $res.url);
            }
        }
    });

    $("#btnRelease").on("click",function () {
        if (Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0) {
            if ($("#markdown-save").hasClass('change')) {
                var comfirm_result = confirm("编辑内容未保存，需要保存吗？")
                if (comfirm_result) {
                    saveDocument(false, releaseBook);
                    return;
                }
            }

            releaseBook();
        } else {
            layer.msg("没有需要发布的文档")
        }
    });

    /**
     * 实现自定义图片上传
     */
    window.editor.getModule('toolbar').addHandler('image',function () {
        var input = document.createElement('input');
        input.setAttribute('type', 'file');
        input.click();

        // Listen upload local image and save to server
        input.onchange = function () {
            var file = input.files[0];

            // file type is only image.
            if (/^image\//.test(file.type)) {
                var form = new FormData();
                form.append('editormd-image-file', file, file.name);

                var layerIndex = 0;

                $.ajax({
                    url: window.imageUploadURL,
                    type: "POST",
                    dataType: "json",
                    data: form,
                    processData: false,
                    contentType: false,
                    error: function() {
                        layer.close(layerIndex);
                        layer.msg("图片上传失败");
                    },
                    success: function(data) {
                        layer.close(layerIndex);
                        if(data.errcode !== 0){
                            layer.msg(data.message);
                        }else{
                            var range = window.editor.getSelection();
                            window.editor.insertEmbed(range.index, 'image', data.url);
                        }
                    }
                });
            } else {
                console.warn('You could only upload images.');
            }
        };
    });
    /**
     * 实现保存
     */
    window.menu_save.on("click",function () {if($(this).hasClass('change')){saveDocument();}});
    /**
     * 设置编辑器变更状态
     * @param $is_change
     */
    function resetEditorChanged($is_change) {
        if ($is_change && !window.isLoad) {
            $("#markdown-save").removeClass('disabled').addClass('change');
        } else {
            $("#markdown-save").removeClass('change').addClass('disabled');
        }
        window.isLoad = false;
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
                window.editor.root.innerHTML = res.data.content;

                // 将原始内容备份
                window.source = res.data.content;
                var node = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name,"identify" : res.data.identify,"version" : res.data.version};
                pushDocumentCategory(node);
                window.selectNode = node;
                window.isLoad = true;

                pushVueLists(res.data.attach);
                initHighlighting();
                setLastSelectNode($node);
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
     * @param callback
     */
    function saveDocument($is_cover,callback) {
        var index = null;
        var node = window.selectNode;

        var html = window.editor.root.innerHTML;

        var content = "";
        if($.trim(html) !== ""){
            content = toMarkdown(html, { gfm: true });
        }
        var version = "";

        if(!node){
            layer.msg("获取当前文档信息失败");
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
                window.saveing = true;
            },
            url :  window.editURL,
            data : {"identify" : window.book.identify,"doc_id" : doc_id,"markdown" : content,"html" : html,"cover" : $is_cover ? "yes":"no","version": version},
            type :"post",
            dataType :"json",
            success : function (res) {
                if(res.errcode === 0){
                    for(var i in window.documentCategory){
                        var item = window.documentCategory[i];

                        if(item.id === doc_id){
                            window.documentCategory[i].version = res.data.version;
                            break;
                        }
                    }
                    resetEditorChanged(false);
                    // 更新内容备份
                    window.source = res.data.content;
                    if(typeof callback === "function"){
                        callback();
                    }
                }else if(res.errcode === 6005){
                    var confirmIndex = layer.confirm('文档已被其他人修改确定覆盖已存在的文档吗？', {
                        btn: ['确定','取消'] //按钮
                    }, function(){
                        layer.close(confirmIndex);
                        saveDocument(true,callback);
                    });
                }else{
                    layer.msg(res.message);
                }
            },
            error : function (XMLHttpRequest, textStatus, errorThrown) {
                layer.msg("服务器错误：" +  errorThrown);
            },
            complete :function () {
                layer.close(index);
                window.saveing = false;
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
                    "label": "添加文档",
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
                    "label": "编辑",
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
                    "label": "删除",
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
        //如果没有选中节点则选中默认节点
        openLastSelectedNode();

    }).on('select_node.jstree', function (node, selected, event) {
        if(window.menu_save.hasClass('change')) {
            if(confirm("编辑内容未保存，需要保存吗？")){
                saveDocument(false,function () {
                    loadDocument(selected);
                });
                return true;
            }
        }
        loadDocument(selected);

    }).on("move_node.jstree", jstree_save)
      .on("delete_node.jstree",function (node,parent) {
          window.isLoad = true;
          window.editor.root.innerHTML ='';
      });

    window.saveDocument = saveDocument;

    window.releaseBook = function () {
        if(Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0){
            if(window.menu_save.hasClass('selected')) {
                if(confirm("编辑内容未保存，需要保存吗？")) {
                    saveDocument();
                }
            }
            $.ajax({
                url : window.releaseURL,
                data :{"identify" : window.book.identify },
                type : "post",
                dataType : "json",
                success : function (res) {
                    if(res.errcode === 0){
                        layer.msg("发布任务已推送到任务队列，稍后将在后台执行。");
                    }else{
                        layer.msg(res.message);
                    }
                }
            });
        }else{
            layer.msg("没有需要发布的文档")
        }
    };
});