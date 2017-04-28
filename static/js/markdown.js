function showError($msg,$id) {
    if(!$id){
        $id = "#form-error-message"
    }
    $($id).addClass("error-message").removeClass("success-message").text($msg);
    return false;
}

function showSuccess($msg,$id) {
    if(!$id){
        $id = "#form-error-message"
    }
    $($id).addClass("success-message").removeClass("error-message").text($msg);
    return true;
}

$(function () {
    window.addDocumentModalFormHtml = $(this).find("form").html();
    window.editor = editormd("docEditor", {
        width : "100%",
        height : "100%",
        path : "/static/editor.md/lib/",
        toolbar : true,
        placeholder: "本编辑器支持Markdown编辑，左边编写，右边预览",
        imageUpload: true,
        imageFormats: ["jpg", "jpeg", "gif", "png", "JPG", "JPEG", "GIF", "PNG"],
        imageUploadURL: window.imageUploadURL ,
        toolbarModes : "full",
        fileUpload: true,
        fileUploadURL : window.fileUploadURL,
        taskList : true,
        flowChart : true,
        htmlDecode : "style,script,iframe,title,onmouseover,onmouseout,style",
        lineNumbers : false,

        tocStartLevel : 1,
        tocm : true,
        saveHTMLToTextarea : true,
        onload : function() {
            this.hideToolbar();
            var $select_node_id = window.treeCatalog.get_selected();
            if($select_node_id) {
                var $select_node = window.treeCatalog.get_node($select_node_id[0])
                if ($select_node) {
                    $select_node.node = {
                        id: $select_node.id
                    };

                    loadDocument($select_node);
                }
            }
        },
        onchange : function () {
            resetEditorChanged(true);
        }
    });
    editormd.loadPlugin("/static/editor.md/plugins/file-dialog/file-dialog");


    /**
     * 实现标题栏操作
     */
    $("#editormd-tools").on("click","a[class!='disabled']",function () {
       var name = $(this).find("i").attr("name");
       if(name === "attachment"){
            window.editor.fileDialog();
       }else if(name === "history"){

       }else if(name === "save"){

            saveDocument(false);

       }else if(name === "sidebar"){
            $("#manualCategory").toggle(0,"swing",function () {

                var $then = $("#manualEditorContainer");
                var left = parseInt($then.css("left"));
                if(left > 0){
                    window.editorContainerLeft = left;
                    $then.css("left","0");
                }else{
                    $then.css("left",window.editorContainerLeft + "px");
                }
                window.editor.resize();
            });
       }else if(name === "release"){

       }else if(name === "tasks") {
           //插入GFM任务列表
           var cm = window.editor.cm;
           var selection = cm.getSelection();

           if (selection === "") {
               cm.replaceSelection("- [x] " + selection);
           }
           else {
               var selectionText = selection.split("\n");

               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "- [x] " + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       }else {
           var action = window.editor.toolbarHandlers[name];

           if (action !== "undefined") {
               $.proxy(action, window.editor)();
               window.editor.focus();
           }
       }
   }) ;

    //实现小提示
    $("[data-toggle='tooltip']").hover(function () {
        var title = $(this).attr('data-title');
        var direction = $(this).attr("data-direction");
        var tips = 3;
        if(direction === "top"){
            tips = 1;
        }else if(direction === "right"){
            tips = 2;
        }else if(direction === "bottom"){
            tips = 3;
        }else if(direction === "left"){
            tips = 4;
        }
        index = layer.tips(title, this, {
            tips: tips
        });
    }, function () {
        layer.close(index);
    });

    $("#btnAddDocument").on("click",function () {
        $("#addDocumentModal").modal("show");
    });
    $("#addDocumentModal").on("hidden.bs.modal",function () {
       $(this).find("form").html(window.addDocumentModalFormHtml);
    }).on("shown.bs.modal",function () {
        $(this).find("input[name='doc_name']").focus();
    });

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

            resetEditor();
            if(res.errcode === 0){
                window.editor.insertValue(res.data.markdown);
                window.editor.setCursor({line:0, ch:0});
                var node = { "id" : res.data.doc_id,'parent' : res.data.parent_id === 0 ? '#' : res.data.parent_id ,"text" : res.data.doc_name,"identify" : res.data.identify,"version" : res.data.version};
                pushDocumentCategory(node);

            }else{
                layer.msg("文档加载失败");
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("文档加载失败");
        });
    }

    /**
     * 创建文档
     */
    function openCreateCatalogDialog($node) {
        var $then =  $("#addDocumentModal");

        var doc_id = $node ? $node.id : 0;

        $then.find("input[name='parent_id']").val(doc_id);

        $then.modal("show");
    }

    /**
     * 将一个节点推送到现有数组中
     * @param $node
     */
    function pushDocumentCategory($node) {
        for (var index in window.documentCategory){
            var item = window.documentCategory[index];
            if(item.id === $node.id){

               window.documentCategory[index] = $node;
                console.log( window.documentCategory[index]);
               return;
            }
        }
        window.documentCategory.push($node);
    }

    /**
     * 打开文档编辑界面
     * @param $node
     */
    function openEditCatalogDialog($node) {
        var $then =  $("#addDocumentModal");
        var doc_id = parseInt($node ? $node.id : 0);
        var text = $node ? $node.text : '';
        var parentId = $node && $node.parent !== '#' ? $node.parent : 0;

        $then.find("input[name='doc_id']").val(doc_id);
        $then.find("input[name='parent_id']").val(parentId);
        $then.find("input[name='doc_name']").val(text);

        for (var index in window.documentCategory){
            var item = window.documentCategory[index];
            if(item.id === doc_id){
                $then.find("input[name='doc_identify']").val(item.identify);
                break;
            }
        }

        $then.modal({ show : true });
    }

    /**
     * 删除一个文档
     * @param $node
     */
    function openDeleteDocumentDialog($node) {
        var index = layer.confirm('你确定要删除该文档吗？', {
            btn: ['确定','取消'] //按钮
        }, function(){

            $.post(window.deleteURL,{"identify" : window.book.identify,"doc_id" : $node.id}).done(function (res) {
                layer.close(index);
                if(res.errcode === 0){
                    window.treeCatalog.delete_node($node);
                    resetEditor($node);
                }else{
                    layer.msg("删除失败",{icon : 2})
                }
            }).fail(function () {
                layer.close(index);
                layer.msg("删除失败",{icon : 2})
            });

        });
    }

    function saveDocument($is_cover) {
        var index = null;
        var node = window.treeCatalog.get_selected();
        var content = window.editor.getMarkdown();
        var html = window.editor.getPreviewedHTML();
        var version = "";
        if(content === ""){
            resetEditorChanged(false);
            return;
        }
        if(!node){
            layer.msg("获取当前文档信息失败");
            return;
        }
        var doc_id = parseInt(node[0]);

        for(var i in window.documentCategory){
            var item = window.documentCategory[i];

            if(item.id === doc_id){
                version = item.version;
                console.log(item)
                break;
            }
        }
        $.ajax({
            beforeSend  : function () {
                index = layer.load(1, {shade: [0.1,'#fff'] });
            },
            url :  window.editURL,
            data : {"identify" : window.book.identify,"doc_id" : node[0],"markdown" : content,"html" : html,"cover" : $is_cover ? "yes":"no","version": version},
            type :"post",
            dataType :"json",
            success : function (res) {
                layer.close(index);
                if(res.errcode === 0){
                    resetEditorChanged(false);
                    for(var i in window.documentCategory){
                        var item = window.documentCategory[i];

                        if(item.id === doc_id){
                            window.documentCategory[i].version = res.data.version;
                            console.log(res.data)
                            break;
                        }
                    }
                }else if(res.errcode === 6005){
                    var confirmIndex = layer.confirm('你确定要删除该文档吗？', {
                        btn: ['确定','取消'] //按钮
                    }, function(){
                        layer.close(confirmIndex);
                        saveDocument(true);
                    });
                }else{
                    layer.msg(res.message);
                }
            }
        });
    }

    function resetEditor($node) {

    }

    function resetEditorChanged($is_change) {
        if($is_change){
            $("#markdown-save").removeClass('disabled').addClass('change');
        }else{
            $("#markdown-save").removeClass('change').addClass('disabled');
        }
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
    });
});