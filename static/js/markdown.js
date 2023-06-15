$(function () {

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
    var htmlDecodeList = ["style","script","title","onmouseover","onmouseout","style"];
    if (!window.IS_ENABLE_IFRAME) {
        htmlDecodeList.unshift("iframe");
    }
/*
    window.editor = editormd("docEditor", {
        width: "100%",
        height: "100%",
        path: window.editormdLib,
        toolbar: true,
        placeholder: window.editormdLocales[window.lang].placeholder,
        imageUpload: true,
        imageFormats: ["jpg", "jpeg", "gif", "png","svg", "JPG", "JPEG", "GIF", "PNG","SVG"],
        imageUploadURL: window.imageUploadURL,
        toolbarModes: "full",
        fileUpload: true,
        fileUploadURL: window.fileUploadURL,
        taskList: true,
        flowChart: true,
        htmlDecode: htmlDecodeList.join(','),
        lineNumbers: true,
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
*/
var CustomHookA = Cherry.createSyntaxHook('codeBlock', Cherry.constants.HOOKS_TYPE_LIST.PAR, {
    makeHtml(str) {
      console.warn('custom hook', 'hello');
      return str;
    },
    rule(str) {
      const regex = {
        begin: '',
        content: '',
        end: '',
      };
      regex.reg = new RegExp(regex.begin + regex.content + regex.end, 'g');
      return regex;
    },
  });
  /**
   * 自定义一个自定义菜单
   * 点第一次时，把选中的文字变成同时加粗和斜体
   * 保持光标选区不变，点第二次时，把加粗斜体的文字变成普通文本
   */
  var customMenuA = Cherry.createMenuHook('加粗斜体', {
    iconName: 'font',
    onClick: function (selection) {
      // 获取用户选中的文字，调用getSelection方法后，如果用户没有选中任何文字，会尝试获取光标所在位置的单词或句子
      let $selection = this.getSelection(selection) || '同时加粗斜体';
      // 如果是单选，并且选中内容的开始结束内没有加粗语法，则扩大选中范围
      if (!this.isSelections && !/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
        this.getMoreSelection('***', '***', () => {
          const newSelection = this.editor.editor.getSelection();
          const isBoldItalic = /^\s*(\*\*\*)[\s\S]+(\1)/.test(newSelection);
          if (isBoldItalic) {
            $selection = newSelection;
          }
          return isBoldItalic;
        });
      }
      // 如果选中的文本中已经有加粗语法了，则去掉加粗语法
      if (/^\s*(\*\*\*)[\s\S]+(\1)/.test($selection)) {
        return $selection.replace(/(^)(\s*)(\*\*\*)([^\n]+)(\3)(\s*)($)/gm, '$1$4$7');
      }
      /**
       * 注册缩小选区的规则
       *    注册后，插入“***TEXT***”，选中状态会变成“***【TEXT】***”
       *    如果不注册，插入后效果为：“【***TEXT***】”
       */
      this.registerAfterClickCb(() => {
        this.setLessSelection('***', '***');
      });
      return $selection.replace(/(^)([^\n]+)($)/gm, '$1***$2***$3');
    }
  });
  /**
   * 定义一个空壳，用于自行规划cherry已有工具栏的层级结构
   */
  var customMenuB = Cherry.createMenuHook('发布', {
    iconName: '',
    onClick: releaseDocument,
  });
  
  var customMenuC = Cherry.createMenuHook("返回", {
    iconName: '',
    onClick: backWard,
  })

  var customMenuD = Cherry.createMenuHook('保存', {
    iconName: '',
    onClick: saveDocument,
  });

  var customMenuE = Cherry.createMenuHook('边栏', {
    iconName: '',
    onClick: siderChange,
  });

  var customMenuF = Cherry.createMenuHook('历史', {
    iconName: '',
    onClick: showHistory,
  });

  
  var basicConfig = {
    id: 'manualEditorContainer',
    externals: {
      echarts: window.echarts,
      katex: window.katex,
      MathJax: window.MathJax,
    },
    isPreviewOnly: false,
    engine: {
      global: {
        urlProcessor(url, srcType) {
          console.log(`url-processor`, url, srcType);
          return url;
        },
      },
      syntax: {
        codeBlock: {
          theme: 'twilight',
        },
        table: {
          enableChart: false,
          // chartEngine: Engine Class
        },
        fontEmphasis: {
          allowWhitespace: false, // 是否允许首尾空格
        },
        strikethrough: {
          needWhitespace: false, // 是否必须有前后空格
        },
        mathBlock: {
          engine: 'MathJax', // katex或MathJax
          src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js', // 如果使用MathJax plugins，则需要使用该url通过script标签引入
        },
        inlineMath: {
          engine: 'MathJax', // katex或MathJax
        },
        emoji: {
          useUnicode: false,
          customResourceURL: 'https://github.githubassets.com/images/icons/emoji/unicode/${code}.png?v8',
          upperCase: true,
        },
        // toc: {
        //     tocStyle: 'nested'
        // }
        // 'header': {
        //   strict: false
        // }
      },
      customSyntax: {
        // SyntaxHookClass
        CustomHook: {
          syntaxClass: CustomHookA,
          force: false,
          after: 'br',
        },
      },
    },
    toolbars: {
      toolbar: [
        'customMenuCName',
        'bold',
        'italic',
        {
          strikethrough: ['strikethrough', 'underline', 'sub', 'sup', 'ruby', 'customMenuAName'],
        },
        'size',
        '|',
        'color',
        'header',
        '|',
        'drawIo',
        '|',
        'ol',
        'ul',
        'checklist',
        'panel',
        'detail',
        '|',
        'formula',
        {
          insert: ['image', 'audio', 'video', 'link', 'hr', 'br', 'code', 'formula', 'toc', 'table', 'pdf', 'word', 'ruby'],
        },
        'graph',
        'togglePreview',
        'settings',
        'switchModel',
        'codeTheme',
        'export',
        'customMenuDName',
        'customMenuBName',
        'customMenuEName',
        'customMenuFName',
        'theme'
      ],
      bubble: ['bold', 'italic', 'underline', 'strikethrough', 'sub', 'sup', 'quote', 'ruby', '|', 'size', 'color'], // array or false
      sidebar: ['mobilePreview', 'copy', 'theme'],
      customMenu: {
        customMenuAName: customMenuA,
        customMenuBName: customMenuB,
        customMenuCName: customMenuC,
        customMenuDName: customMenuD,
        customMenuEName: customMenuE,
        customMenuFName: customMenuF,
      },
    },
    drawioIframeUrl: '/static/cherry/drawio_demo.html',
    editor: {
      defaultModel: 'edit&preview',
      height: "100%",
    },
    previewer: {
      // 自定义markdown预览区域class
      // className: 'markdown'
    },
    keydown: [],
    //extensions: [],
    //callback: {
    //changeString2Pinyin: pinyin,
    //}
  };
  
  fetch('').then((response) => response.text()).then((value) => {
    //var markdownarea = document.getElementById("markdown_area").value
    var config = Object.assign({}, basicConfig);// { value: markdownarea });// { value: value });不显示获取的初始化值
    window.editor = new Cherry(config);
    openLastSelectedNode();
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
       } else if (name === "tasks") {
           // 插入 GFM 任务列表
           var cm = window.editor.cm;
           var selection = cm.getSelection();
           var cursor    = cm.getCursor();
           if (selection === "") {
               cm.setCursor(cursor.line, 0);
               cm.replaceSelection("- [x] " + selection);
               cm.setCursor(cursor.line, cursor.ch + 6);
           } else {
               var selectionText = selection.split("\n");

               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "- [x] " + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       } else {
           var action = window.editor.toolbarHandlers[name];

           if (!!action && action !== "undefined") {
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
                    window.editor.setMarkdown(res.data.markdown);
                }catch(e){
                    console.log(e);
                }
                var node = { "id": res.data.doc_id, 'parent': res.data.parent_id === 0 ? '#' : res.data.parent_id, "text": res.data.doc_name, "identify": res.data.identify, "version": res.data.version };
                pushDocumentCategory(node);
                window.selectNode = node;
                pushVueLists(res.data.attach);
                setLastSelectNode($node);
            } else {
                layer.msg(editormdLocales[lang].loadDocFailed);
            }
        }).fail(function () {
            layer.close(index);
            layer.msg(editormdLocales[lang].loadDocFailed);
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
        var html = window.editor.getHtml(true);
        var version = "";

        if (!node) {
            layer.msg(editormdLocales[lang].fetchDocFailed);
            return;
        }
        if (node.a_attr && node.a_attr.disabled) {
            layer.msg(editormdLocales[lang].cannotAddToEmptyNode);
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
                    var confirmIndex = layer.confirm(editormdLocales[lang].overrideModified, {
                        btn: [editormdLocales[lang].confirm, editormdLocales[lang].cancel] // 按钮
                    }, function() {
                        layer.close(confirmIndex);
                        saveDocument(true, callback);
                    });
                } else {
                    layer.msg(res.message);
                }
            },
            error : function (XMLHttpRequest, textStatus, errorThrown) {
                layer.msg(window.editormdLocales[window.lang].serverExcept +  errorThrown);
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
     * 返回上一个页面
     */
    function backWard() {
        history.back();
    }

    /**
     * 发布文档
     */

    function releaseDocument() {
        if (Object.prototype.toString.call(window.documentCategory) === '[object Array]' && window.documentCategory.length > 0) {
            if ($("#markdown-save").hasClass('change')) {
                var confirm_result = confirm(editormdLocales[lang].contentUnsaved);
                if (confirm_result) {
                    saveDocument(false, releaseBook);
                    return;
                }
            }
            releaseBook();
            return
        }
        layer.msg(editormdLocales[lang].noDocNeedPublish)
    }


    /**
     * 显示/隐藏边栏
     */

    function siderChange() {
        $("#manualCategory").toggle(0, "swing", function () {
            var $then = $("#manualEditorContainer");
            var left = parseInt($then.css("left"));
            if (left > 0) {
                window.editorContainerLeft = left;
                $then.css("left", "0");
            } else {
                $then.css("left", window.editorContainerLeft + "px");
            }
        });
    }

    /**
     * 显示文档历史
     */

    function showHistory() {
        window.documentHistory();
    }

    /**
     * 添加文档
     */
    $("#addDocumentForm").ajaxForm({
        beforeSubmit: function () {
            var doc_name = $.trim($("#documentName").val());
            if (doc_name === "") {
                return showError(editormdLocales[lang].contentsNameEmpty, "#add-error-message")
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
    }).on("ready.jstree",function () {
        window.treeCatalog = $("#sidebar").jstree(true);

        //如果没有选中节点则选中默认节点
        // openLastSelectedNode();
    }).on('select_node.jstree', function (node, selected) {

        if ($("#markdown-save").hasClass('change')) {
            if (confirm(window.editormdLocales[window.lang].contentUnsaved)) {
                saveDocument(false, function () {
                    loadDocument(selected);
                });
                return true;
            }
        }
        //如果是空目录则直接出发展开下一级功能
        if (selected.node.a_attr && selected.node.a_attr.disabled) {
            selected.instance.toggle_node(selected.node);
            return false
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
                layer.msg(window.editormdLocales[window.lang].loadFailed);
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
                return showError(window.editormdLocales[window.lang].tplNameEmpty, "#saveTemplateForm .show-error-message");
            }
            var content = $("#saveTemplateForm").find("input[name='content']").val();
            if (content === ""){
                return showError(window.editormdLocales[window.lang].tplContentEmpty, "#saveTemplateForm .show-error-message");
            }

            $("#btnSaveTemplate").button("loading");

            return true;
        },
        success: function ($res) {
            if($res.errcode === 0){
                $("#saveTemplateModal").modal("hide");
                layer.msg(window.editormdLocales[window.lang].saveSucc);
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
                layer.msg(window.editormdLocales[window.lang].serverExcept);
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
                layer.msg(window.editormdLocales[window.lang].serverExcept);
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
               var table = "| " + window.editormdLocales[window.lang].paramName 
                + "  | " + window.editormdLocales[window.lang].paramType
                + " | " + window.editormdLocales[window.lang].example
                + "  |  " + window.editormdLocales[window.lang].remark
                + " |\n| ------------ | ------------ | ------------ | ------------ |\n";
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