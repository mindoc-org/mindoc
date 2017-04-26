/**
 * 插入视频插件，编辑器必须开启html识别
 *
 * @file video-dialog.js
 * @author Minho
 * @version 0.1
 * @licence MIT
 */
(function() {
     var factory = function (exports) {

            var $ = jQuery;
            var pluginName = "video-dialog";

            exports.fn.videoDialog = function () {
                var _this = this;
                var lang = this.lang;
                var editor = this.editor;
                var settings = this.settings;
                var path = settings.pluginPath + pluginName + "/";
                var classPrefix = this.classPrefix;
                var dialogName = classPrefix + pluginName, dialog;
                var dialogLang = lang.dialog.help;

                if (editor.find("." + dialogName).length < 1) {

                    var dialogContent = "";

                    dialog = this.createDialog({
                        name: dialogName,
                        title: dialogLang.title,
                        width: 840,
                        height: 540,
                        mask: settings.dialogShowMask,
                        drag: settings.dialogDraggable,
                        content: dialogContent,
                        lockScreen: settings.dialogLockScreen,
                        maskStyle: {
                            opacity: settings.dialogMaskOpacity,
                            backgroundColor: settings.dialogMaskBgColor
                        },
                        buttons: {
                            close: [lang.buttons.close, function () {
                                this.hide().lockScreen(false).hideMask();

                                return false;
                            }]
                        }
                    });
                }

                dialog = editor.find("." + dialogName);

                this.dialogShowMask(dialog);
                this.dialogLockScreen();
                dialog.show();

                var videoContent = dialog.find(".markdown-body");

                if (videoContent.html() === "") {
                    $.get(path + "help.md", function (text) {
                        var md = exports.$marked(text);
                        videoContent.html(md);

                        videoContent.find("a").attr("target", "_blank");
                    });
                }
            };

        };

    // CommonJS/Node.js
    if (typeof require === "function" && typeof exports === "object" && typeof module === "object")
    {
        module.exports = factory;
    }
    else if (typeof define === "function")  // AMD/CMD/Sea.js
    {
        if (define.amd) { // for Require.js

            define(["editormd"], function(editormd) {
                factory(editormd);
            });

        } else { // for Sea.js
            define(function(require) {
                var editormd = require("./../../editormd");
                factory(editormd);
            });
        }
    }
    else
    {
        factory(window.editormd);
    }
})();

