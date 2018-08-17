/**
 * Created by lifei6671 on 2017/4/29 0029.
 */

/**
 * 打开最后选中的节点
 */
function openLastSelectedNode() {
    //如果文档树或编辑器没有准备好则不加载文档
    if (window.treeCatalog == null || window.editor == null) {
        return false;
    }
    var $isSelected = false;
    if(window.localStorage){
        var $selectedNodeId = window.localStorage.getItem("MinDoc::LastLoadDocument:" + window.book.identify);
        try{
            if($selectedNodeId){
                //遍历文档树判断是否存在节点
                $.each(window.documentCategory,function (i, n) {
                    if(n.id == $selectedNodeId && !$isSelected){
                        var $node = {"id" : n.id};
                        window.treeCatalog.deselect_all();
                        window.treeCatalog.select_node($node);
                        $isSelected = true;
                    }
                });

            }
        }catch($ex){
            console.log($ex)
        }
    }

    //如果节点不存在，则默认选中第一个节点
    if (!$isSelected && window.documentCategory.length > 0){
        var doc = window.documentCategory[0];

        if(doc && doc.id > 0){
            var node = {"id": doc.id};
            $("#sidebar").jstree(true).select_node(node);
            $isSelected = true;
        }
    }
    return $isSelected;
}

/**
 * 设置最后选中的文档
 * @param $node
 */
function setLastSelectNode($node) {
    if(window.localStorage) {
        if (typeof $node === "undefined" || !$node) {
            window.localStorage.removeItem("MinDoc::LastLoadDocument:" + window.book.identify);
        } else {
            var nodeId = $node.id ? $node.id : $node.node.id;
            window.localStorage.setItem("MinDoc::LastLoadDocument:" + window.book.identify, nodeId);
        }
    }
}

/**
 * 保存排序
 * @param node
 * @param parent
 */
function jstree_save(node, parent) {

    var parentNode = window.treeCatalog.get_node(parent.parent);

    var nodeData = window.getSiblingSort(parentNode);

    if (parent.parent !== parent.old_parent) {
        parentNode = window.treeCatalog.get_node(parent.old_parent);
        var newNodeData = window.getSiblingSort(parentNode);
        if (newNodeData.length > 0) {
            nodeData = nodeData.concat(newNodeData);
        }
    }

    var index = layer.load(1, {
        shade: [0.1, '#fff'] //0.1透明度的白色背景
    });

    $.ajax({
        url : window.sortURL,
        type :"post",
        data : JSON.stringify(nodeData),
        success : function (res) {
            layer.close(index);
            if (res.errcode === 0){
                layer.msg("保存排序成功");
            }else{
                layer.msg(res.message);
            }
        }
    })
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
 * 处理排序
 * @param node
 * @returns {Array}
 */
function getSiblingSort (node) {
    var data = [];

    for(var key in node.children){
        var index = data.length;

        data[index] = {
            "id" : parseInt(node.children[key]),
            "sort" : parseInt(key),
            "parent" : Number(node.id) ? Number(node.id) : 0
        };
    }
    return data;
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
                window.documentCategory.remove(function (item) {
                   return item.id == $node.id;
                });


                // console.log(window.documentCategory)
                setLastSelectNode();
            }else{
                layer.msg("删除失败",{icon : 2})
            }
        }).fail(function () {
            layer.close(index);
            layer.msg("删除失败",{icon : 2})
        });

    });
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

    if($node.a_attr && $node.a_attr.is_open){
        $then.find("input[name='is_open'][value='1']").prop("checked","checked");
    }else{
        $then.find("input[name='is_open'][value='0']").prop("checked","checked");
    }

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
 * 将一个节点推送到现有数组中
 * @param $node
 */
function pushDocumentCategory($node) {
    for (var index in window.documentCategory){
        var item = window.documentCategory[index];
        if(item.id === $node.id){

            window.documentCategory[index] = $node;
            return;
        }
    }
    window.documentCategory.push($node);
}
/**
 * 将数据重置到Vue列表中
 * @param $lists
 */
function pushVueLists($lists) {

    window.vueApp.lists = [];
    $.each($lists,function (i, item) {
        window.vueApp.lists.push(item);
    });
}

/**
 * 发布项目
 */
function releaseBook() {
    $.ajax({
        url: window.releaseURL,
        data: { "identify": window.book.identify },
        type: "post",
        dataType: "json",
        success: function (res) {
            if (res.errcode === 0) {
                layer.msg("发布任务已推送到任务队列，稍后将在后台执行。");
            } else {
                layer.msg(res.message);
            }
        }
    });
}
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
//弹出创建文档的遮罩层
$("#btnAddDocument").on("click",function () {
    $("#addDocumentModal").modal("show");
});
//用于还原创建文档的遮罩层
$("#addDocumentModal").on("hidden.bs.modal",function () {
    // $(this).find("form").html(window.sessionStorage.getItem("addDocumentModal"));
}).on("shown.bs.modal",function () {
    $(this).find("input[name='doc_name']").focus();
}).on("show.bs.modal",function () {
    // window.sessionStorage.setItem("addDocumentModal",$(this).find("form").html())
});

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

window.documentHistory = function() {
    layer.open({
        type: 2,
        title: '历史版本',
        shadeClose: true,
        shade: 0.8,
        area: ['700px','80%'],
        content: window.historyURL + "?identify=" + window.book.identify + "&doc_id=" + window.selectNode.id,
        end : function () {
            if(window.SelectedId){
                var selected = {node:{
                    id : window.SelectedId
                }};
                window.loadDocument(selected);
                window.SelectedId = null;
            }
        }
    });
};

function uploadImage($id,$callback) {
    /** 粘贴上传图片 **/
    document.getElementById($id).addEventListener('paste', function(e) {
        if(e.clipboardData && e.clipboardData.items) {
            var clipboard = e.clipboardData;
            for (var i = 0, len = clipboard.items.length; i < len; i++) {
                if (clipboard.items[i].kind === 'file' || clipboard.items[i].type.indexOf('image') > -1) {

                    var imageFile = clipboard.items[i].getAsFile();

                    var fileName = String((new Date()).valueOf());

                    switch (imageFile.type) {
                        case "image/png" :
                            fileName += ".png";
                            break;
                        case "image/jpg" :
                            fileName += ".jpg";
                            break;
                        case "image/jpeg" :
                            fileName += ".jpeg";
                            break;
                        case "image/gif" :
                            fileName += ".gif";
                            break;
                        default :
                            layer.msg("不支持的图片格式");
                            return;
                    }
                    var form = new FormData();

                    form.append('editormd-image-file', imageFile, fileName);

                    var layerIndex = 0;

                    $.ajax({
                        url: window.imageUploadURL,
                        type: "POST",
                        dataType: "json",
                        data: form,
                        processData: false,
                        contentType: false,
                        beforeSend: function () {
                            layerIndex = $callback('before');
                        },
                        error: function () {
                            layer.close(layerIndex);
                            $callback('error');
                            layer.msg("图片上传失败");
                        },
                        success: function (data) {
                            layer.close(layerIndex);
                            $callback('success', data);
                            if (data.errcode !== 0) {
                                layer.msg(data.message);
                            }

                        }
                    });
                    e.preventDefault();
                }
            }
        }
    });
}

/**
 * 初始化代码高亮
 */
function initHighlighting() {
    $('pre code,pre.ql-syntax').each(function (i, block) {
        hljs.highlightBlock(block);
    });
}
$(function () {
    window.vueApp = new Vue({
        el : "#attachList",
        data : {
            lists : []
        },
        delimiters : ['${','}'],
        methods : {
            removeAttach : function ($attach_id) {
                var $this = this;
                var item = $this.lists.filter(function ($item) {
                    return $item.attachment_id == $attach_id;
                });

                if(item && item[0].hasOwnProperty("state")){
                    $this.lists = $this.lists.filter(function ($item) {
                        return $item.attachment_id != $attach_id;
                    });
                    return;
                }
                $.ajax({
                    url : window.removeAttachURL,
                    type : "post",
                    data : { "attach_id" : $attach_id},
                    success : function (res) {
                        if(res.errcode === 0){
                            $this.lists = $this.lists.filter(function ($item) {
                                return $item.attachment_id != $attach_id;
                            });
                        }else{
                            layer.msg(res.message);
                        }
                    }
                });
            }
        },
        watch : {
            lists : function ($lists) {
                $("#attachInfo").text(" " + $lists.length + " 个附件")
            }
        }
    });
    /**
     * 启动自动保存，默认30s自动保存一次
     */
    if(window.book.auto_save){
        setTimeout(function () {
            setInterval(function () {
                var $then =  $("#markdown-save");
                if(!window.saveing && $then.hasClass("change")){
                    $then.trigger("click");
                }
            },30000);
        },30000);
    }
    /**
     * 当离开窗口时存在未保存的文档会提示保存
     */
    $(window).on("beforeunload",function () {
        if($("#markdown-save").hasClass("change")){
            return '您输入的内容尚未保存，确定离开此页面吗？';
        }
    });
});

