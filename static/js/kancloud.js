var events = function () {
    var articleOpen = function (event, $param) {
        //当打开文档时，将文档ID加入到本地缓存。
        window.sessionStorage && window.sessionStorage.setItem("MinDoc::LastLoadDocument:" + window.book.identify, $param.$id);
        var prevState = window.history.state || {};
        if ('pushState' in history) {
            if ($param.$id) {
                prevState.$id === $param.$id || window.history.pushState($param, $param.$id, $param.$url);
            } else {
                window.history.pushState($param, $param.$id, $param.$url);
                window.history.replaceState($param, $param.$id, $param.$url);
            }
        } else {
            window.location.hash = $param.$url;
        }

        initHighlighting();
        $(window).resize();

        $(".manual-right").scrollTop(0);
        //使用layer相册功能查看图片
        layer.photos({ photos: "#page-content" });
    };

    return {
        data: function ($key, $value) {
            $key = "MinDoc::Document:" + $key;
            if (typeof $value === "undefined") {
                return $("body").data($key);
            } else {
                return $('body').data($key, $value);
            }
        },
        trigger: function ($e, $obj) {
            if ($e === "article.open") {
                articleOpen($e, $obj);
            } else {
                $('body').trigger('article.open', $obj);
            }
        }
    }

}();

function format($d) {
    return $d < 10 ? "0" + $d : "" + $d;
}

function timeFormat($time) {
    var span = Date.parse($time)
    var date = new Date(span)
    var year = date.getFullYear();
    var month = format(date.getMonth() + 1);
    var day = format(date.getDate());
    var hour = format(date.getHours());
    var min = format(date.getMinutes());
    var sec = format(date.getSeconds());
    return year + "-" + month + "-" + day + " " + hour + ":" + min + ":" + sec;
}

// 点击翻页
function pageClicked($page, $docid) {
    if (!window.IS_DISPLAY_COMMENT) {
        return;
    }
    $("#articleComment").removeClass('not-show-comment');
    $.ajax({
        url: "/comment/lists?page=" + $page + "&docid=" + $docid,
        type: "GET",
        success: function ($res) {
            console.log($res.data);
            loadComment($res.data.page, $res.data.doc_id);
        },
        error: function () {
            layer.msg("加载失败");
        }
    });
}

function renderOperateSection(comment) {
    const deleteIcon = comment.show_del == 1
        ? `<i class="delete e-delete glyphicon glyphicon-remove" onclick="onDelComment(${comment.comment_id})"></i>`
        : '';

    return `
        <span class="operate ${comment.show_del == 1 ? 'toggle' : ''}">
            <span class="number">${comment.index}#</span>
            ${deleteIcon}
        </span>`;
}

// 加载评论
function loadComment($page, $docid) {
    $("#commentList").empty();
    let html = ""
    let c = $page.List;
    for (let i = 0; c && i < c.length; i++) {
        const comment = c[i];
        html += `
        <div class="comment-item" data-id="${comment.comment_id}">
            <p class="info">
                <img src="${comment.avatar}" alt="">
                <a class="name">${comment.author}</a>
                <span class="date">${timeFormat(comment.comment_date)}</span>
            </p>
            <div class="content">${comment.content}</div>
            <p class="util">
                ${renderOperateSection(comment)}
            </p>
        </div>`;
    }
    $("#commentList").append(html);
    if ($page.TotalPage > 1) {
        $("#page").bootstrapPaginator({
            currentPage: $page.PageNo,
            totalPages: $page.TotalPage,
            bootstrapMajorVersion: 3,
            size: "middle",
            onPageClicked: function (e, originalEvent, type, page) {
                pageClicked(page, $docid);
            }
        });
    } else {
        $("#page").find("li").remove();
    }
}

// 删除评论
function onDelComment($id) {
    $.ajax({
        url: "/comment/delete",
        data: { "id": $id },
        type: "POST",
        success: function ($res) {
            if ($res.errcode == 0) {
                layer.msg("删除成功");
                $("div[data-id=" + $id + "]").remove();
            } else {
                layer.msg($res.message);
            }
        },
        error: function () {
            layer.msg("删除失败");
        }
    });
}

// 重新渲染页面
function renderPage($data) {
    $("#page-content").html($data.body);
    $("title").text($data.title);
    $("#article-title").text($data.doc_title);
    $("#article-info").text($data.doc_info);
    $("#view_count").text("阅读次数：" + $data.view_count);
    $("#doc_id").val($data.doc_id);
    checkMarkdownTocElement();
    if ($data.page) {
        loadComment($data.page, $data.doc_id);
    } else {
        pageClicked(-1, $data.doc_id);
    }

    if ($data.is_markdown) {
        if ($("#view_container").hasClass($data.markdown_theme)) {
            return
        }
        $("#view_container").removeClass("theme__dark theme__green theme__light theme__red theme__default")
        $("#view_container").addClass($data.markdown_theme)
    }

}

/***
 * 加载文档到阅读区
 * @param $url
 * @param $id
 * @param $callback
 */
function loadDocument($url, $id, $callback) {
    $.ajax({
        url: $url,
        type: "GET",
        beforeSend: function () {
            var data = events.data($id);
            if (data) {
                if (typeof $callback === "function") {
                    data.body = $callback(data.body);
                } else if (data.version && data.version != $callback) {
                    return true;
                }
                renderPage(data);

                loadCopySnippets();
                events.trigger('article.open', { $url: $url, $id: $id });

                return false;

            }

            NProgress.start();
        },
        success: function ($res) {
            if ($res.errcode === 0) {
                var data = $res.data;
                if (typeof $callback === "function") {
                    data.body = $callback(data.body);
                }
                renderPage(data);
                loadCopySnippets();
                events.data($id, data);
                events.trigger('article.open', { $url: $url, $id: $id });
            } else if ($res.errcode === 6000) {
                window.location.href = "/";
            } else {
                layer.msg("加载失败");
            }
        },
        complete: function () {
            NProgress.done();
        },
        error: function () {
            layer.msg("加载失败");
        }
    });
}

/**
 * 初始化代码高亮
 */
function initHighlighting() {
    try {
        $('pre,pre.ql-syntax').each(function (i, block) {
            if ($(this).hasClass('prettyprinted') || $(this).hasClass('hljs')) {
                return;
            }
            hljs.highlightBlock(block);
        });
        // hljs.initLineNumbersOnLoad();
    } catch (e) {
        console.log(e);
    }
}

function handleEvent(event) {
    switch (event.keyCode) {
        case 70: // ctrl + f 打开搜索面板 并获取焦点
            $(".navg-item[data-mode='search']").click();
            document.getElementById('searchForm').querySelector('input').focus();
            event.preventDefault();
            break;
        case 27: // esc 关闭搜索面板
            $(".navg-item[data-mode='view']").click();
            event.preventDefault();
            break;
    }
}

$(function () {
    window.addEventListener('keydown', handleEvent)

    checkMarkdownTocElement();
    $(".view-backtop").on("click", function () {
        $('.manual-right').animate({ scrollTop: '0px' }, 200);
    });
    $(".manual-right").scroll(function () {
        try {
            var top = $(".manual-right").scrollTop();
            if (top > 100) {
                $(".view-backtop").addClass("active");
            } else {
                $(".view-backtop").removeClass("active");
            }
        } catch (e) {
            console.log(e);
        }

        try {
            var scrollTop = $("body").scrollTop();
            var oItem = $(".markdown-heading").find(".reference-link");
            var oName = "";
            $.each(oItem, function () {
                var oneItem = $(this);
                var offsetTop = oneItem.offset().top;

                if (offsetTop - scrollTop < 100) {
                    oName = "#" + oneItem.attr("name");
                }
            });
            $(".markdown-toc-list a").each(function () {
                if (oName === $(this).attr("href")) {
                    $(this).parents("li").addClass("directory-item-active");
                } else {
                    $(this).parents("li").removeClass("directory-item-active");
                }
            });
            if (!$(".markdown-toc-list li").hasClass('directory-item-active')) {
                $(".markdown-toc-list li:eq(0)").addClass("directory-item-active");
            }
        } catch (e) {
            console.log(e);
        }
    }).on("click", ".markdown-toc-list a", function () {
        var $this = $(this);
        setTimeout(function () {
            $(".markdown-toc-list li").removeClass("directory-item-active");
            $this.parents("li").addClass("directory-item-active");
        }, 10);
    }).find(".markdown-toc-list li:eq(0)").addClass("directory-item-active");


    $(window).resize(function (e) {
        var h = $(".manual-catalog").innerHeight() - 50;
        $(".markdown-toc").height(h);
    }).resize();

    window.isFullScreen = false;

    initHighlighting();
    window.jsTree = $("#sidebar").jstree({
        'plugins': ["wholerow", "types"],
        "types": {
            "default": {
                "icon": false  // 删除默认图标
            }
        },
        'core': {
            'check_callback': true,
            "multiple": false,
            'animation': 0
        }
    }).on('select_node.jstree', function (node, selected) {
        //如果是空目录则直接出发展开下一级功能
        if (selected.node.a_attr && selected.node.a_attr.disabled) {
            selected.instance.toggle_node(selected.node);
            return false
        }
        $(".m-manual").removeClass('manual-mobile-show-left');
        loadDocument(selected.node.a_attr.href, selected.node.id, selected.node.a_attr['data-version']);
    });

    $("#slidebar").on("click", function () {
        $(".m-manual").addClass('manual-mobile-show-left');
    });
    $(".manual-mask").on("click", function () {
        $(".m-manual").removeClass('manual-mobile-show-left');
    });

    /**
     * 关闭侧边栏
     */
    $(".manual-fullscreen-switch").on("click", function () {
        isFullScreen = !isFullScreen;
        if (isFullScreen) {
            $(".m-manual").addClass('manual-fullscreen-active');
        } else {
            $(".m-manual").removeClass('manual-fullscreen-active');
        }
    });

    $(".navg-item[data-mode]").on("click", function () {
        var mode = $(this).data('mode');
        $(this).siblings().removeClass('active').end().addClass('active');
        $(".m-manual").removeClass("manual-mode-view manual-mode-collect manual-mode-search").addClass("manual-mode-" + mode);
    });

    const input = document.getElementById('searchForm').querySelector('input');
    input.addEventListener('input', function() {
        $("#btnSearch").click();
    });

    /**
     * 项目内搜索
     */
    $("#searchForm").ajaxForm({
        beforeSubmit: function () {
            var keyword = $.trim($("#searchForm").find("input[name='keyword']").val());
            if (keyword === "") {
                $(".search-empty").show();
                $("#searchList").html("");
                return false;
            }
            $("#btnSearch").attr("disabled", "disabled").find("i").removeClass("fa-search").addClass("loading");
            window.keyword = keyword;
        },
        success: function (res) {
            var html = "";
            if (res.errcode === 0) {
                for (var i in res.data) {
                    var item = res.data[i];
                    html += '<li><a href="javascript:;" title="' + item.doc_name + '" data-id="' + item.doc_id + '"> ' + item.doc_name + ' </a></li>';
                }
            }
            if (html !== "") {
                $(".search-empty").hide();
            } else {
                $(".search-empty").show();
            }
            $("#searchList").html(html);
        },
        complete: function () {
            $("#btnSearch").removeAttr("disabled").find("i").removeClass("loading").addClass("fa-search");
        }
    });

    window.onpopstate = function (e) {
        var $param = e.state;
        if (!$param) return;
        if ($param.hasOwnProperty("$url")) {
            window.jsTree.jstree().deselect_all();

            if ($param.$id) {
                window.jsTree.jstree().select_node({ id: $param.$id });
            } else {
                window.location.assign($param.$url);
            }
            // events.trigger('article.open', $param);
        } else {
            console.log($param);
        }
    };

    // 提交评论
    $("#commentForm").ajaxForm({
        beforeSubmit: function () {
            $("#btnSubmitComment").button("loading");
        },
        success: function (res) {
            if (res.errcode === 0) {
                layer.msg("保存成功");
            } else {
                layer.msg(res.message);
            }
            $("#btnSubmitComment").button("reset");
            $("#commentContent").val("");
            pageClicked(-1, res.data.doc_id); // -1 表示请求最后一页
        },
        error: function () {
            layer.msg("服务错误");
            $("#btnSubmitComment").button("reset");
        }
    });
    loadCopySnippets();
});

function loadCopySnippets() {
    $("pre").addClass("line-numbers language-bash");
    $("pre").attr('data-prismjs-copy', '复制');
    $("pre").attr('data-prismjs-copy-error', '按Ctrl+C复制');
    $("pre").attr('data-prismjs-copy-success', '代码已复制！');
    var snippets = document.querySelectorAll('pre code');
    [].forEach.call(snippets, function (snippet) {
        Prism.highlightElement(snippet);
    });
}

function checkMarkdownTocElement() {
    let toc = $(".markdown-toc-list");
    if ($(".toc").length) {
        toc = $(".toc");
    }
    let articleComment = $("#articleComment");
    if (toc.length) {
        $(".wiki-bottom-left").css("width", "calc(100% - 260px)");
        articleComment.css("width", "calc(100% - 260px)");
        articleComment.css("margin", "30px 0 70px 0");
    } else {
        $(".wiki-bottom-left").css("width", "100%");
        articleComment.css("width", "100%");
        articleComment.css("margin", "30px auto 70px auto;");
    }
}