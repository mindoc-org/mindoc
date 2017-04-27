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
    $("#addDocumentModal").on("show.bs.modal",function () {
        window.addDocumentModalFormHtml = $(this).find("form").html();
    }).on("hidden.bs.modal",function () {
       $(this).find("form").html(window.addDocumentModalFormHtml);
    });

    function loadDocument($node) {
        var index = layer.load(1, {
            shade: [0.1,'#fff'] //0.1透明度的白色背景
        });

        $.get("/docs/"+ window.book.identify +"/" + $node.node.id ).done(function (data) {
            win.isEditorChange = true;
            layer.close(index);
            $("#documentId").val(selected.node.id);
            window.editor.clear();
            if(data.errcode === 0 && data.data.doc.content){
                window.editor.insertValue(data.data.doc.content);
                window.editor.setCursor({line:0, ch:0});
            }else if(data.errcode !== 0){
                layer.msg("文档加载失败");
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("文档加载失败");
        });
    }
    /**
     * 添加文档
     */
    $("#addDocumentForm").ajaxForm({
        beforeSubmit : function () {
            var doc_name = $.trim($("#documentName").val());
            if (doc_name === ""){
                return showError("目录名称不能为空","#add-error-message")
            }
            return true;
        },
        success : function (res) {
            if(res.errcode === 0){
                var data = { "id" : res.data.doc_id,'parent' : res.data.parent_id,"text" : res.data.doc_name};

                var node = window.treeCatalog.get_node(data.id);
                if(node){
                    window.treeCatalog.rename_node({"id":data.id},data.text);
                }else {
                    var result = window.treeCatalog.create_node(res.data.parent_id, data, 'last');
                    window.treeCatalog.deselect_all();
                    window.treeCatalog.select_node(data);
                    window.editor.clear();
                }
                $("#markdown-save").removeClass('change').addClass('disabled');
                $("#addDocumentModal").modal('hide');
            }else{
                showError(res.message,"#add-error-message")
            }
        },
        error :function () {

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
                        editDocumentDialog(node);
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
                        deleteDocumentDialog(node);
                    }
                }
            }
        }
    }).on('loaded.jstree', function () {
        window.treeCatalog = $(this).jstree();
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
    });
});