$(function () {
    window.isLoad = true;
    editormd.katexURL = {
        js  : window.baseUrl + "/static/katex/katex",
        css : window.baseUrl + "/static/katex/katex"
    };
    window.addDocumentModalFormHtml = $(this).find("form").html();
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
        tocStartLevel: 1,
        tocm: true,
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

            var $select_node_id = window.treeCatalog.get_selected();
            if ($select_node_id) {
                var $select_node = window.treeCatalog.get_node($select_node_id[0])
                if ($select_node) {
                    $select_node.node = {
                        id: $select_node.id
                    };

                    loadDocument($select_node);
                }
            }

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
        if($(this).hasClass('disabled')){
            return false;
        }
       var name = $(this).find("i").attr("name");
       if (name === "attachment") {
           $("#uploadAttachModal").modal("show");
       } else if (name === "history") {
           window.documentHistory();
       } else if (name === "save") {
            saveDocument(false);
       } else if (name === "template") {
           $("#documentTemplateModal").modal("show");
       } else if (name=="word2md"){
			$("#word2mdform")[0].reset();
			$("#output").html("");
			$("#messages").html("");
			$("#word2md").modal("show");
			
		}else if (name=="Pasteoffice"){
			$("#Pasteofficeform")[0].reset();
			$("#Pastearea").innerText="";
			$("#officeoutmd").val("");
			$("#Pasteoffice").modal("show");
		}else if (name === "sidebar") {
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
                var node = { "id": res.data.doc_id, 'parent': res.data.parent_id === 0 ? '#' : res.data.parent_id, "text": res.data.doc_name, "identify": res.data.identify, "version": res.data.version };
                if(res.data.is_lock){
                    node.type = "lock";
                    node.text = node.text + "<span class='lock-text'> [锁定]</span>";
                     window.editor.config('readOnly',true);
                }else{
                    node.type = "unlock";
                    window.editor.config('readOnly',false);
                }
                window.isLoad = true;
                try {
                    window.editor.clear();
                    if(res.data.markdown !== ""){
                        window.editor.insertValue(res.data.markdown);
                    }else{
                        window.isLoad = true;
                    }
                    window.editor.setCursor({line: 0, ch: 0});
                }catch(e){
                    console.log(e);
                }

                pushDocumentCategory(node);
                window.selectNode = node;
                pushVueLists(res.data.attach);

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
            },
            url: window.editURL,
            data: { "identify": window.book.identify, "doc_id": doc_id, "markdown": content, "html": html, "cover": $is_cover ? "yes" : "no", "version": version },
            type: "post",
            timeout : 30000,
            dataType: "json",
            success: function (res) {
                layer.close(index);
                if (res.errcode === 0) {
                    resetEditorChanged(false);
                    for (var i in window.documentCategory) {
                        var item = window.documentCategory[i];

                        if (item.id === doc_id) {
                            window.documentCategory[i].version = res.data.version;
                            break;
                        }
                    }
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
        // console.log($is_change);
        // console.log(window.isLoad);
        if ($is_change && !window.isLoad) {
            var type = window.treeCatalog.get_type(window.selectNode);
            if(type === "lock"){
                return;
            }
            $("#markdown-save").removeClass('disabled').addClass('change');
        } else {
            $("#markdown-save").removeClass('change').addClass('disabled');
        }
        window.isLoad = false;
    }

    /**
     * 添加顶级文档
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
                var data = { "id": res.data.doc_id, 'parent': res.data.parent_id === 0 ? '#' : res.data.parent_id , "text": res.data.doc_name, "identify": res.data.identify, "version": res.data.version ,"type": res.data.is_lock ? "lock" : "unlock"};

                if(res.data.is_lock){
                    data.text = data.text + "<span class='lock-text'> [锁定]</span>";
                }
                var node = window.treeCatalog.get_node(data.id);
                if (node) {
                    window.treeCatalog.rename_node({ "id": data.id }, data.text);
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
            },
            "lock" : {
                "icon" : false
            },
            "unlock" : {
                "icon" : false
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
            'items' : function(node) {
                var items = {
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
                    },
                    "unlock" : {
                        "separator_before": false,
                        "separator_after": true,
                        "_disabled": false,
                        "label": "解锁",
                        "icon": "fa fa-unlock",
                        "action": function (data) {
                            var inst = $.jstree.reference(data.reference);
                            var node = inst.get_node(data.reference);
                            unLockDocumentAction(node);
                        }
                    },
                    "lock" : {
                        "separator_before": false,
                        "separator_after": true,
                        "_disabled": false,
                        "label": "锁定",
                        "icon": "fa fa-lock",
                        "action": function (data) {
                            var inst = $.jstree.reference(data.reference);
                            var node = inst.get_node(data.reference);
                            lockDocumentAction(node);
                        }
                    }
                };
                console.log(this.get_type(node));
                if(this.get_type(node) === "lock") {
                    delete items.lock;
                }else{
                    delete items.unlock;
                }
                return items;
            },
        }
    }).on('loaded.jstree', function () {
        window.treeCatalog = $(this).jstree();
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
    }).on("move_node.jstree", jstree_save);
    
    	//对html进行预处理
	function firstfilter(str) {
	    //去掉头部描述
	    return str.replace(/<head>[\s\S]*<\/head>/gi, "")
	    //去掉注释
	    .replace(/<!--[\s\S]*?-->/ig, "")
	    //去掉隐藏元素
	    .replace(/<([a-z0-9]*)[^>]*\s*display:none[\s\S]*?><\/\1>/gi, '');
	}
	

	//去除冗余属性和标签
	function filterPasteWord(str) {
	    return str.replace(/[\t\r\n]+/g, ' ').replace(/<!--[\s\S]*?-->/ig, "")
	    //转换图片
	    .replace(/<v:shape [^>]*>[\s\S]*?.<\/v:shape>/gi,
	        function (str) {
	        //opera能自己解析出image所这里直接返回空
	        if (!!window.opera && window.opera.version) {
	            return '';
	        }
	        try {
	            //有可能是bitmap占为图，无用，直接过滤掉，主要体现在粘贴excel表格中
	            if (/Bitmap/i.test(str)) {
	                return '';
	            }
	            var src = str.match(/src=\s*"([^"]*)"/i)[1];
	            return '<img  src="' + src + '" />';
	        } catch (e) {
	            return '';
	        }
	    })

	    //针对wps添加的多余标签处理
	    .replace(/<\/?div[^>]*>/g, '')
	    .replace(/<\/?span[^>]*>/g, '')
	    .replace(/<\/?font[^>]*>/g, '')
	    .replace(/<\/?col[^>]*>/g, '')
	    .replace(/<\/?(span|div|o:p|v:.*?|input|label)[\s\S]*?>/g, '')
	    //去掉所有属性,需要保留单元格合并
	    //.replace(/<([a-zA-Z]+)\s*[^><]*>/g, "<$1>")
	    //去掉多余的属性
	    .replace(/v:\w+=(["']?)[^'"]+\1/g, '')
	    .replace(/<(!|script[^>]*>.*?<\/script(?=[>\s])|\/?(\?xml(:\w+)?|xml|meta|link|style|\w+:\w+)(?=[\s\/>]))[^>]*>/gi, "").replace(/<p [^>]*class="?MsoHeading"?[^>]*>(.*?)<\/p>/gi, "<p><strong>$1</strong></p>")
	    //去掉多余的属性
	    .replace(/\s+(class|lang|align)\s*=\s*(['"]?)([\w-]+)\2/gi, '')
	    //清除多余的font/span不能匹配&nbsp;有可能是空格
	    .replace(/<(font|span)[^>]*>(\s*)<\/\1>/gi,
	        function (a, b, c) {
	        return c.replace(/[\t\r\n ]+/g, " ");
	    })
	    //去掉style属性
	    .replace(/(<[a-z][^>]*)\sstyle=(["'])([^\2]*?)\2/gi, "$1")
	    // 去除不带引号的属性
	    .replace(/(class|border|cellspacing|MsoNormalTable|valign|width|center|&nbsp;|x:str|height|x:num|cellpadding)(=[^ \f\n\r\t\v>]*)?/g, "")
	    // 去除多余空格
	    .replace(/(\S+)(\s+)/g, function (match, p1, p2) {
	        return p1 + ' ';
	    })
	    .replace(/(\s)(>|<)/g, function (match, p1, p2) {
	        return p2;
	    })
	    //处理表格中的p标签
	    .replace(/(<table[^>]*[\s\S]*?><\/table>)/gi, function (a) {
			//嵌套表格不处理
	        if (a.match(/(<table>)/gi).length > 1) {
	            return a

	        }
	        if (!/<thead>/i.test(a) && !/(rowspan|colspan)/i.test(a)) {
	            //没有表头，将第一行作为表头
	            const firstrow = "<table><thead>" + a.match(/<tr>[\s\S]*?(?=<tr>)/i)[0] + "</thead>";
	            a = a.replace(/<tr>[\s\S]*?(?=<tr>)/i, "")
	                .replace(/<table>/, firstrow);
	        } else if (/<thead>/i.test(a) && /(rowspan|colspan)/i.test(a)) {
				a=a.replace(/<thead>/, "");
			}
			return a.replace(/<\/p><p>/ig, "<br/>")
	            .replace(/<\/?p[^>]*>/ig, '')
				//.replace(/<td><\/td>/ig,"<td>&nbsp;</td>")
	            .replace(/<td>&nbsp;<\/td>/g, "<td></td>")
	    });

	}

	//判断粘贴的内容是否来自office
	function isWordDocument(str) {
	    return /(class="?Mso|style="[^"]*\bmso\-|w:WordDocument|<(v|o):|lang=)/ig.test(str) || /\"urn:schemas-microsoft-com:office:office/ig.test(str);
	}
	
	//excel表格处理
	function pasteClipboardHtml(html) {
	    const startFramgmentStr = '<!--StartFragment-->';
	    const endFragmentStr = '<!--EndFragment-->';
	    const startFragmentIndex = html.indexOf(startFramgmentStr);
	    const endFragmentIndex = html.lastIndexOf(endFragmentStr);

	    if (startFragmentIndex > -1 && endFragmentIndex > -1) {
	        html = html.slice(startFragmentIndex + startFramgmentStr.length, endFragmentIndex);
	    }

	    // Wrap with <tr> if html contains dangling <td> tags
	    // Dangling <td> tag is that tag does not have <tr> as parent node.
	    if (/<\/td>((?!<\/tr>)[\s\S])*$/i.test(html)) {
	        html = '<tr>' + html + '</tr>';
	    }
	    // Wrap with <table> if html contains dangling <tr> tags
	    // Dangling <tr> tag is that tag does not have <table> as parent node.
	    if (/<\/tr>((?!<\/table>)[\s\S])*$/i.test(html)) {
	        html = '<table>' + html + '</table>';
	    }

	    return html;
	}
	
	//将html转换为markdown
	function html2md(str) {
	    var gfm = turndownPluginGfm.gfm;
	    var turndownService = new TurndownService({
	            headingStyle: 'atx',
	            hr: '- - -',
	            bulletListMarker: '-',
	            codeBlockStyle: 'indented',
	            fence: '```',
	            emDelimiter: '_',
	            strongDelimiter: '**'
	        });
	    turndownService.use(gfm);
	    turndownService.keep(['sub', 'sup']);//保留标签
		str=firstfilter(str);
		//str=pasteClipboardHtml(str);
		str=filterPasteWord(str);
	    return turndownService.turndown(str);
	}
	
	//将word转换的html转换为markdown，并插入编辑器
	$("#btnhtml2md").click(function (e) {
	    e.preventDefault();
	    if ($(this).hasClass("disabled"))
	        return false;
	    var str = $("#output").html();
	    var cm =  window.editor.cm;
	    var cursor = cm.getCursor();
	    var selection = cm.getSelection();
	    cm.replaceSelection(html2md(str)+"\n\n");
	    $("#btnhtml2md").removeClass("disabled");
	    $("#word2md").modal('hide');
	    cm.focus();
	});
	//粘贴Pastearea自动获得焦点
	 $('#Pasteoffice').on('shown.bs.modal', function (event) {
        var modal = $(this);
		var form=modal.find("#Pasteofficeform");
		var renameInput =form.find("#Pastearea");
		//获得焦点时文本全选
        renameInput.focus(function () {
            this.select();
        });
        renameInput.focus();
    });
	
	//粘贴解析
	$("#Pastearea")[0].addEventListener('paste', function (e) {
	    var clipboard = e.clipboardData;
	    if (!(clipboard && clipboard.items)) {
	        return;
	    }
	    for (var i = 0, len = clipboard.items.length; i < len; i++) {
	        var item = clipboard.items[i];
	        if (item.kind === "string" && item.type === "text/html") {
	            item.getAsString(function (str) {
	                if (/<img [^>]*src=['"]([^'"]+)[^>]*>/gi.test(str)) {
	                    layer.msg("粘贴的内容中包含有图片，建议使用word转markdown模块处理！");
	                }
	                var markdown = html2md(str);
	                $("#officeoutmd").val(markdown);
	            });

	        }
	    }

	});
	
	//解析html源码为markdown
	$("#HtmlToMarkdown").click(function (e) {
	    e.preventDefault();
	    var str = $("#Pastearea").val();
	    var markdown = html2md(str);
	    $("#officeoutmd").val(markdown);

	});

	
	//将html转换为markdown，并插入编辑器
	$("#office2md").click(function (e) {
	    e.preventDefault();
	    if ($(this).hasClass("disabled"))
	        return false;
	    var str =  $("#officeoutmd").val();
	    var cm =  window.editor.cm;
	    var cursor = cm.getCursor();
	    var selection = cm.getSelection();
	    cm.replaceSelection(str+"\n\n");
	    $("#office2md").removeClass("disabled");
	    $("#Pasteoffice").modal('hide');
	    cm.focus();
	});
    
    
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
