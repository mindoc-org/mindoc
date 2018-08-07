$(function () {
    editormd.katexURL = {
        js  : window.baseUrl + "/static/katex/katex",
        css : window.baseUrl + "/static/katex/katex"
    };
    window.editor = editormd("docEditor", {
        width: "100%",
        height: "100%",
        path: window.baseUrl + "/static/editor.md/lib/",
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
        highlightStyle: window.highlightStyle ? window.highlightStyle : "github",
        tocStartLevel: 1,
        tocm: true,
        tex:true,
        saveHTMLToTextarea: true,

        onload: function() {
            this.hideToolbar();
            var keyMap = {
                "Ctrl-S": function(cm) {
                    saveBlog(false);
                },
                "Cmd-S": function(cm){
                    saveBlog(false);
                },
                "Ctrl-A": function(cm) {
                    cm.execCommand("selectAll");
                }
            };
            this.addKeyMap(keyMap);


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
        },
        onchange: function () {
            resetEditorChanged(true);
        }
    });
    /**
     * 实现标题栏操作
     */
    $("#editormd-tools").on("click", "a[class!='disabled']", function () {
        var name = $(this).find("i").attr("name");
        if (name === "attachment") {
            $("#uploadAttachModal").modal("show");
        }else if (name === "save") {
            saveBlog(false);
        } else if (name === "template") {
            $("#documentTemplateModal").modal("show");
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

    /**
     * 保存文章
     * @param $is_cover
     */
    function saveBlog($is_cover) {
        var content = window.editor.getMarkdown();
        var html = window.editor.getPreviewedHTML();

        $.ajax({
            beforeSend: function () {
                index = layer.load(1, { shade: [0.1, '#fff'] });
            },
            url: window.editURL,
            data: { "blogId": window.blogId,"content": content,"htmlContent": html, "cover": $is_cover ? "yes" : "no","version" : window.blogVersion},
            type: "post",
            timeout : 30000,
            dataType: "json",
            success: function ($res) {
                layer.close(index);
                if ($res.errcode === 0) {
                    resetEditorChanged(false);
                    window.blogVersion = $res.data.version;
                    console.log(window.blogVersion);
                } else if($res.errcode === 6005) {
                    var confirmIndex = layer.confirm('文档已被其他人修改确定覆盖已存在的文档吗？', {
                        btn: ['确定', '取消'] // 按钮
                    }, function() {
                        layer.close(confirmIndex);
                        saveBlog(true);
                    });
                } else {
                    layer.msg(res.message);
                }
            },
            error : function (XMLHttpRequest, textStatus, errorThrown) {
                layer.close(index);
                layer.msg("服务器错误：" +  errorThrown);
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
     * 打开文档模板
     */
    $("#documentTemplateModal").on("click", ".section>a[data-type]", function () {
        var $this = $(this).attr("data-type");
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
});