let draggedItem = null;

document.querySelectorAll('.list-item').forEach(item => {
    item.addEventListener('mousedown', function() {
        draggedItem = item;
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
                    } else {
                        layer.msg("排序失败", {icon: 2});
                    }
                },
                error: function (err) {
                    console.log('error:', err)
                }
            })
            draggedItem = null;
        }
    });
});