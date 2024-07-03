let draggedItem = null;
let draggedItemIndex = null;

document.querySelectorAll('.list-item').forEach(item => {
    item.addEventListener('mousedown', function() {
        draggedItem = item;
        // 获取当前item在list-item中的索引
        const parentNode = item.parentNode;
        draggedItemIndex = Array.from(parentNode.children).indexOf(draggedItem);
        item.setAttribute('draggable', true);
    });

    item.addEventListener('dragstart', function() {
    });

    item.addEventListener('dragover', function(e) {
        e.preventDefault();
        if (draggedItem !== item) {
            item.style.border = '1px dashed #d4d4d5';
        }
    });

    item.addEventListener('dragleave', function() {
        item.style.border = '';
    });

    item.addEventListener('drop', function() {
        item.style.border = '';
        if (draggedItem !== null) {
            const parentNode = item.parentNode;
            const draggedIndex = Array.from(parentNode.children).indexOf(draggedItem);
            const targetIndex = Array.from(parentNode.children).indexOf(item);

            if (draggedIndex < targetIndex) {
                parentNode.insertBefore(draggedItem, item.nextSibling);
            } else {
                parentNode.insertBefore(draggedItem, item);
            }
            // 获取当前最新的data-id属性列表 去除null
            const newSortList = Array.from(parentNode.children).map(item => item.getAttribute('data-id')).filter(item => item !== null);
            const newSortListStr = newSortList.join(',');
            // 更新排序
            $.ajax({
                url: window.updateBookOrder,
                type: 'POST',
                data: {
                    ids: newSortListStr,
                },
                success: function (res) {
                    if (res.errcode === 0) {
                        layer.msg("排序成功", {icon: 1});
                        draggedItem = null;
                    } else {
                        layer.msg(res.message, {icon: 2});
                        const parentNode = item.parentNode;
                        // 在parentNode中找到当前拖拽的item
                        const draggedIndex = Array.from(parentNode.children).indexOf(draggedItem);
                        console.log('draggedIndex:', draggedIndex);
                        // 移除当前拖拽的item
                        parentNode.removeChild(draggedItem);
                        // 将draggedItem放到原来的位置
                        parentNode.insertBefore(draggedItem, parentNode.children[draggedItemIndex]);
                        draggedItem = null;
                    }
                },
                error: function (err) {
                    console.log('error:', err)
                    draggedItem = null;
                }
            })
        }
    });
});