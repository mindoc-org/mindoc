/**
 * Created by lifei6671 on 2017/4/29 0029.
 */
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

    console.log(JSON.stringify(nodeData));

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
};

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
 * 将数据重置到Vue列表中
 * @param $lists
 */
function pushVueLists($lists) {

    window.vueApp.lists = [];
    for(var j in $lists){
        var item = $lists[j];
        window.vueApp.lists.push(item);
    }
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
    $(this).find("form").html(window.addDocumentModalFormHtml);
}).on("shown.bs.modal",function () {
    $(this).find("input[name='doc_name']").focus();
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

$(function () {
    $("#documentHistoryModal").on("shown.bs.modal",function () {
        var historyVue = new Vue({
            el : "#documentHistoryModal",
            data : {
                lists : []
            },
            delimiters : ['${','}'],
            methods : {

            }
        });

        $.ajax({
            url : window.historyURL,
            data : { "identify" : window.book.identify,"doc_id" : window.selectNode.id },
            dataType :"json",
            success : function (res) {
                if(res.errcode === 0){
                    historyVue.lists = res.data.lists;
                }else{
                    alert(res.message);
                }
            }
        });
    });
});