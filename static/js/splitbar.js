$(function () {
    var splitBar = {
        // 设置当前屏幕为 840px 时将分割条隐藏
        maxWidth: 840,

        // 父元素选择器
        parentSelector: '.manual-body',

        /**
         * 初始化分割条
         */
        init: function () {
            var self = this;
            $(self.parentSelector)
            .append(
                $('\
                    <div id="manual-vsplitbar" unselectable="on"\
                    style="\
                        z-index:301;\
                        position: absolute;\
                        user-select: none;\
                        cursor: col-resize;\
                        left: 275px;\
                        height: 100%;\
                        display: block;\
                        width: 3px;">\
                        <a href="javascript:void(0)" accesskey="" tabindex="0" title="vsplitbar"></a>\
                    </div>\
                ')
                .hover(
                    function (event) {
                        event.target.style.background = '#8f949ad8';
                    },
                    function (event) {
                        event.target.style.background = '';
                    }
                )
            );

            self.loadingHtml();

            // 设置媒体查询
            mediaMatcher.set();
        },

        /**
         * 加载页面时设置分割条是否显示
         */ 
        loadingHtml: function () {
            var self = this;
            var htmlWidth = document.body.clientWidth;
            if (htmlWidth <= self.maxWidth) $('#manual-vsplitbar').css('display', 'none');
        },

        /**
         * 在打开关闭菜单事初始化左右窗口和分割条
         */
        reset: function () {
            if (isFullScreen) {
                // 关闭菜单时，初始化左右窗口
                $('#manual-vsplitbar').css('display', 'none');
                splitBar.inMaxWidthReset();
                $('.manual-left').css('width', '0px');
            } else {
                 // 打开菜单时，初始化左右窗口
                $('#manual-vsplitbar').css('display', 'block');
                splitBar.outMaxWidthReset();
            }
        },

        /**
         * 屏幕小于等于 840px，重置左右窗口
         */
        inMaxWidthReset: function () {
            $('#manual-vsplitbar').css('display', 'none');
            $('.m-manual.manual-reader .manual-right').css('left', '0');
            $('.manual-left').css('width', '280px');
        },

        /**
         * 屏幕大于 840px，重置左右窗口
         */
        outMaxWidthReset: function () {
            $('.manual-right').css('left', '279px');
            $('.manual-left').css('width', '279px');
            $('#manual-vsplitbar').css('left', '275px').css('display', 'block');
        }
    }

    /**
     * 设置媒体查询
     * 分割条隐藏
     */
    var mediaMatcher = {
        mql: window.matchMedia('(max-width:' + splitBar.maxWidth + 'px)'),

        /**
         * 添加媒体查询监听事件
         */ 
        set: function () {
            var self = this;
            self.mql.addListener(self.mqCallback);
        },

        /**
         * 删除媒体查询监听事件
         */
        remove: function () {
            var self = this;
            self.mql.removeListener(self.mqCallback);
        },

        /**
         * 媒体查询回调函数
         */
        mqCallback: function (event) {
            if (event.target.matches) { // 宽度小于等于 840 像素
                $(".m-manual").removeClass('manual-fullscreen-active');
                splitBar.inMaxWidthReset();
            } else {                    // 宽度大于 840 像素
                splitBar.outMaxWidthReset();
            }
        }
    }

    // 初始化分割条
    splitBar.init();
    
    /**
     * 控制菜单宽度
     * 最小为 140px
     */
    $('#manual-vsplitbar').on('mousedown', function (e) {
        var bodyEle = $('.manual-body')[0];

        var leftPane = $('.manual-left')[0];
        var splitBar = $('#manual-vsplitbar')[0];
        var rightPane = $('.manual-right')[0];

        var disX = (e || event).clientX;

        splitBar.left = splitBar.offsetLeft;

        var docMouseMove = document.onmousemove;
        var docMouseUp = document.onmouseup;

        document.onmousemove = function (e) {
            var curPos = splitBar.left + ((e || event).clientX - disX);
            var maxPos = bodyEle.clientWidth - splitBar.offsetWidth;

            curPos > maxPos && (curPos = maxPos);
            curPos < 140 && (curPos = 140);

            leftPane.style.width = curPos + "px";
            splitBar.style.left = curPos - 3 + "px";
            rightPane.style.left = curPos + 3 + "px";

            return false;
        }

        document.onmouseup = function () {
            document.onmousemove = docMouseMove;
            document.onmouseup = docMouseUp;
            splitBar.releaseCapture && splitBar.releaseCapture();
        };

        splitBar.setCapture && splitBar.setCapture();
        return false;
    });

    /**
     * 关闭侧边栏时，同步分割条；
     */
    $(".manual-fullscreen-switch").on("click", splitBar.reset);
});
