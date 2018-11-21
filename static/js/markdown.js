$(function () {
    editormd.katexURL = {
        js  : window.katex.js,
        css : window.katex.css
    };

    window.editor = editormd("docEditor", {
        width: "100%",
        height: "100%",
        path: window.editormdLib,
        toolbar: true,
        placeholder: "本编辑器支持 Markdown 编辑，左边编写，右边预览。",
        imageUpload: true,
        imageFormats: ["jpg", "jpeg", "gif", "png", "JPG", "JPEG", "GIF", "PNG"],
        imageUploadURL: window.imageUploadURL,
        toolbarModes: "full",
        fileUpload: true,
        fileUploadURL: window.fileUploadURL,
        taskList: true,
        flowChart: true,
        htmlDecode: "style,script,iframe,title,onmouseover,onmouseout,style",
        lineNumbers: false,
        sequenceDiagram: true,
        tocStartLevel: 1,
        tocm: true,
        previewCodeHighlight: 1,
        highlightStyle: window.highlightStyle ? window.highlightStyle : "github",
        tex:true,
        saveHTMLToTextarea: true,

        onload: function() {
            this.hideToolbar();
            var keyMap = {
                "Ctrl-S": function(cm) {
                    saveDocument(false);
                },
                "Cmd-S": function(cm){
                    saveDocument(false);
                },
                "Ctrl-A": function(cm) {
                    cm.execCommand("selectAll");
                }
            };
            this.addKeyMap(keyMap);

            //如果没有选中节点则选中默认节点
            openLastSelectedNode();
            uploadImage("docEditor", function ($state, $res) {
                if ($state === "before") {
                    return layer.load(1, {
                        shade: [0.1, '#fff'] // 0.1 透明度的白色背景
                    });
                } else if ($state === "success") {
                    if ($res.errcode === 0) {
                        var value = '![](' + $res.url + ')';
                        window.editor.insertValue(value);
                    }
                }
            });
            window.isLoad = true;
        },
        onchange: function () {
            resetEditorChanged(true);
        }
    });

    function insertToMarkdown(body) {
        window.isLoad = true;
        window.editor.insertValue(body);
        window.editor.setCursor({ line: 0, ch: 0 });
        resetEditorChanged(true);
    }
    function insertAndClearToMarkdown(body) {
        window.isLoad = true;
        window.editor.clear();
        window.editor.insertValue(body);
        window.editor.setCursor({ line: 0, ch: 0 });
        resetEditorChanged(true);
    }

    /**
     * 实现标题栏操作
     */
    $("#editormd-tools").on("click", "a[class!='disabled']", function () {
       var name = $(this).find("i").attr("name");
       if (name === "attachment") {
           $("#uploadAttachModal").modal("show");
       } else if (name === "history") {
           window.documentHistory();
       } else if (name === "save") {
            saveDocument(false);
       } else if (name === "template") {
           $("#documentTemplateModal").modal("show");
       } else if(name === "save-template"){
           $("#saveTemplateModal").modal("show");
       } else if(name === 'json'){
           $("#convertJsonToTableModal").modal("show");
       } else if (name === "sidebar") {
            $("#manualCategory").toggle(0, "swing", function () {
                var $then = $("#manualEditorContainer");
                var left = parseInt($then.css("left"));
                if (left > 0) {
                    window.editorContainerLeft = left;
                    $then.css("left", "0");
                } else {
                    $then.css("left", window.editorContainerLeft + "px");
                }
                window.editor.resize();
            });
       } else if (name === "release") {
            if (Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0) {
                if ($("#markdown-save").hasClass('change')) {
                    var confirm_result = confirm("编辑内容未保存，需要保存吗？");
                    if (confirm_result) {
                        saveDocument(false, releaseBook);
                        return;
                    }
                }

                releaseBook();
            } else {
                layer.msg("没有需要发布的文档")
            }
       } else if (name === "tasks") {
           // 插入 GFM 任务列表
           var cm = window.editor.cm;
           var selection = cm.getSelection();

           if (selection === "") {
               cm.replaceSelection("- [x] " + selection);
           } else {
               var selectionText = selection.split("\n");

               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "- [x] " + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       } else {
           var action = window.editor.toolbarHandlers[name];

           if (action !== "undefined") {
               $.proxy(action, window.editor)();
               window.editor.focus();
           }
       }
   }) ;

    /***
     * 加载指定的文档到编辑器中
     * @param $node
     */
    window.loadDocument = function($node) {
        var index = layer.load(1, {
            shade: [0.1, '#fff'] // 0.1 透明度的白色背景
        });

        $.get(window.editURL + $node.node.id ).done(function (res) {
            layer.close(index);

            if (res.errcode === 0) {
                window.isLoad = true;
                try {
                    window.editor.clear();
                    window.editor.insertValue(res.data.markdown);
                    window.editor.setCursor({line: 0, ch: 0});
                }catch(e){
                    console.log(e);
                }
                var node = { "id": res.data.doc_id, 'parent': res.data.parent_id === 0 ? '#' : res.data.parent_id, "text": res.data.doc_name, "identify": res.data.identify, "version": res.data.version };
                pushDocumentCategory(node);
                window.selectNode = node;
                pushVueLists(res.data.attach);
                setLastSelectNode($node);
            } else {
                layer.msg("文档加载失败");
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("文档加载失败");
        });
    };

    /**
     * 保存文档到服务器
     * @param $is_cover 是否强制覆盖
     */
    function saveDocument($is_cover, callback) {
        var index = null;
        var node = window.selectNode;
        var content = window.editor.getMarkdown();
        var html = window.editor.getPreviewedHTML();
        var version = "";

        if (!node) {
            layer.msg("获取当前文档信息失败");
            return;
        }

        var doc_id = parseInt(node.id);

        for (var i in window.documentCategory) {
            var item = window.documentCategory[i];

            if (item.id === doc_id) {
                version = item.version;
                break;
            }
        }
        $.ajax({
            beforeSend: function () {
                index = layer.load(1, { shade: [0.1, '#fff'] });
                window.saveing = true;
            },
            url: window.editURL,
            data: { "identify": window.book.identify, "doc_id": doc_id, "markdown": content, "html": html, "cover": $is_cover ? "yes" : "no", "version": version },
            type: "post",
            timeout : 30000,
            dataType: "json",
            success: function (res) {
                if (res.errcode === 0) {
                    resetEditorChanged(false);
                    for (var i in window.documentCategory) {
                        var item = window.documentCategory[i];

                        if (item.id === doc_id) {
                            window.documentCategory[i].version = res.data.version;
                            break;
                        }
                    }
                    $.each(window.documentCategory,function (i, item) {
                        var $item = window.documentCategory[i];

                        if (item.id === doc_id) {
                            window.documentCategory[i].version = res.data.version;
                        }
                    });
                    if (typeof callback === "function") {
                        callback();
                    }

                } else if(res.errcode === 6005) {
                    var confirmIndex = layer.confirm('文档已被其他人修改确定覆盖已存在的文档吗？', {
                        btn: ['确定', '取消'] // 按钮
                    }, function() {
                        layer.close(confirmIndex);
                        saveDocument(true, callback);
                    });
                } else {
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

    /**
     * 添加文档
     */
    $("#addDocumentForm").ajaxForm({
        beforeSubmit: function () {
            var doc_name = $.trim($("#documentName").val());
            if (doc_name === "") {
                return showError("目录名称不能为空", "#add-error-message")
            }
            $("#btnSaveDocument").button("loading");
            return true;
        },
        success: function (res) {
            if (res.errcode === 0) {
                var data = {
                    "id": res.data.doc_id,
                    'parent': res.data.parent_id === 0 ? '#' : res.data.parent_id ,
                    "text": res.data.doc_name,
                    "identify": res.data.identify,
                    "version": res.data.version ,
                    state: { opened: res.data.is_open == 1},
                    a_attr: { is_open: res.data.is_open == 1}
                };

                var node = window.treeCatalog.get_node(data.id);
                if (node) {
                    window.treeCatalog.rename_node({ "id": data.id }, data.text);
                    $("#sidebar").jstree(true).get_node(data.id).a_attr.is_open = data.state.opened;
                } else {
                    window.treeCatalog.create_node(data.parent, data);
                    window.treeCatalog.deselect_all();
                    window.treeCatalog.select_node(data);
                }
                pushDocumentCategory(data);
                $("#markdown-save").removeClass('change').addClass('disabled');
                $("#addDocumentModal").modal('hide');
            } else {
                showError(res.message, "#add-error-message");
            }
            $("#btnSaveDocument").button("reset");
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
    }).on("ready.jstree",function () {
        window.treeCatalog = $("#sidebar").jstree(true);

        //如果没有选中节点则选中默认节点
        // openLastSelectedNode();
    }).on('select_node.jstree', function (node, selected, event) {

        if ($("#markdown-save").hasClass('change')) {
            if (confirm("编辑内容未保存，需要保存吗？")) {
                saveDocument(false, function () {
                    loadDocument(selected);
                });
                return true;
            }
        }

        loadDocument(selected);
    }).on("move_node.jstree", jstree_save).on("delete_node.jstree",function($node,$parent) {
        openLastSelectedNode();
    });
    /**
     * 打开文档模板
     */
    $("#documentTemplateModal").on("click", ".section>a[data-type]", function () {
        var $this = $(this).attr("data-type");
        if($this === "customs"){
            $("#displayCustomsTemplateModal").modal("show");
            return;
        }
        var body = $("#template-" + $this).html();
        if (body) {
            window.isLoad = true;
            window.editor.clear();
            window.editor.insertValue(body);
            window.editor.setCursor({ line: 0, ch: 0 });
            resetEditorChanged(true);
        }
        $("#documentTemplateModal").modal('hide');
    });
    /**
     * 展示自定义模板列表
     */
    $("#displayCustomsTemplateModal").on("show.bs.modal",function () {
        window.sessionStorage.setItem("displayCustomsTemplateList",$("#displayCustomsTemplateList").html());

        var index ;
        $.ajax({
            beforeSend: function () {
                index = layer.load(1, { shade: [0.1, '#fff'] });
            },
           url : window.template.listUrl,
           data: {"identify":window.book.identify},
           type: "POST",
           dataType: "html",
            success: function ($res) {
                $("#displayCustomsTemplateList").html($res);
            },
            error : function () {
                layer.msg("加载失败请重试");
            },
            complete : function () {
                layer.close(index);
            }
        });
        $("#documentTemplateModal").modal("hide");
    }).on("hidden.bs.modal",function () {
        var cache = window.sessionStorage.getItem("displayCustomsTemplateList");
        $("#displayCustomsTemplateList").html(cache);
    });
    /**
     * 添加模板
     */
    $("#saveTemplateForm").ajaxForm({
        beforeSubmit: function () {
            var doc_name = $.trim($("#templateName").val());
            if (doc_name === "") {
                return showError("模板名称不能为空", "#saveTemplateForm .show-error-message");
            }
            var content = $("#saveTemplateForm").find("input[name='content']").val();
            if (content === ""){
                return showError("模板内容不能为空", "#saveTemplateForm .show-error-message");
            }

            $("#btnSaveTemplate").button("loading");

            return true;
        },
        success: function ($res) {
            if($res.errcode === 0){
                $("#saveTemplateModal").modal("hide");
                layer.msg("保存成功");
            }else{
                return showError($res.message, "#saveTemplateForm .show-error-message");
            }
        },
        complete : function () {
            $("#btnSaveTemplate").button("reset");
        }
    });
    /**
     * 当添加模板弹窗事件发生
     */
    $("#saveTemplateModal").on("show.bs.modal",function () {
        window.sessionStorage.setItem("saveTemplateModal",$(this).find(".modal-body").html());
        var content = window.editor.getMarkdown();
        $("#saveTemplateForm").find("input[name='content']").val(content);
        $("#saveTemplateForm .show-error-message").html("");
    }).on("hidden.bs.modal",function () {
        $(this).find(".modal-body").html(window.sessionStorage.getItem("saveTemplateModal"));
    });
    /**
     * 插入自定义模板内容
     */
    $("#displayCustomsTemplateList").on("click",".btn-insert",function () {
        var templateId = $(this).attr("data-id");

        $.ajax({
            url: window.template.getUrl,
            data :{"identify": window.book.identify, "template_id": templateId},
            dataType: "json",
            type: "get",
            success : function ($res) {
               if ($res.errcode !== 0){
                   layer.msg($res.message);
                   return;
               }
                window.isLoad = true;
                window.editor.clear();
                window.editor.insertValue($res.data.template_content);
                window.editor.setCursor({ line: 0, ch: 0 });
                resetEditorChanged(true);
                $("#displayCustomsTemplateModal").modal("hide");
            },error : function () {
                layer.msg("服务器异常");
            }
        });
    }).on("click",".btn-delete",function () {
        var $then = $(this);
        var templateId = $then.attr("data-id");
        $then.button("loading");

        $.ajax({
            url : window.template.deleteUrl,
            data: {"identify": window.book.identify, "template_id": templateId},
            dataType: "json",
            type: "post",
            success: function ($res) {
                if($res.errcode !== 0){
                    layer.msg($res.message);
                }else{
                    $then.parents("tr").empty().remove();
                }
            },error : function () {
                layer.msg("服务器异常");
            },
            complete: function () {
                $then.button("reset");
            }
        })
    });

    $("#btnInsertTable").on("click",function () {
       var content = $("#jsonContent").val();
       if(content !== "") {
           try {
               var jsonObj = $.parseJSON(content);
               var data = foreachJson(jsonObj,"");
               var table = "| 参数名称  | 参数类型  | 示例值  |  备注 |\n| ------------ | ------------ | ------------ | ------------ |\n";
               $.each(data,function (i,item) {
                    table += "|" + item.key + "|" + item.type + "|" + item.value +"| |\n";
               });
                insertToMarkdown(table);
           }catch (e) {
               showError("Json 格式错误:" + e.toString(),"#json-error-message");
               return;
           }
       }
       $("#convertJsonToTableModal").modal("hide");
    });
    $("#convertJsonToTableModal").on("hidden.bs.modal",function () {
        $("#jsonContent").val("");
    }).on("shown.bs.modal",function () {
        $("#jsonContent").focus();
    });
});