$(function () {
    window.editor = editormd("docEditor", {
        width : "100%",
        height : "100%",
        path : "/static/editor.md/lib/",
        toolbar : true,
        placeholder: "本编辑器支持Markdown编辑，左边编写，右边预览",
        imageUpload: true,
        imageFormats: ["jpg", "jpeg", "gif", "png", "JPG", "JPEG", "GIF", "PNG"],
        imageUploadURL: "/upload",
        toolbarModes : "full",
        fileUpload: true,
        taskList : true,
        flowChart : true,
        htmlDecode : "style,script,iframe,title,onmouseover,onmouseout,style",
        lineNumbers : false,
        fileUploadURL : '/upload',
        tocStartLevel : 1,
        tocm : true,
        saveHTMLToTextarea : true,
        onload : function() {
            this.hideToolbar();
           // this.videoDialog();
            //this.executePlugin("video-dialog","video-dialog/video-dialog");
        }
    });

    editormd.loadPlugin("/static/editor.md/plugins/video-dialog/video-dialog", function(){
        //editormd.videoDialog();
    });

    $("#editormd-tools").on("click","a[class!='disabled']",function () {
       var name = $(this).find("i").attr("name");

       if(name === "attachment"){

       }else if(name === "history"){

       }else if(name === "sidebar"){
            $("#manualCategory").toggle(0,"swing",function () {

                var $then = $("#manualEditorContainer");
                var left = parseInt($then.css("left"));
                if(left > 0){
                    window.editorContainerLeft = left;
                    $then.css("left","0");
                }else{
                    $then.css("left",window.editorContainerLeft + "px");
                }
                window.editor.resize();
            });
       }else if(name === "release"){

       }else if(name === "tasks") {
           //插入GFM任务列表
           var cm = window.editor.cm;
           var selection = cm.getSelection();

           if (selection === "") {
               cm.replaceSelection("- [x] " + selection);
           }
           else {
               var selectionText = selection.split("\n");

               for (var i = 0, len = selectionText.length; i < len; i++) {
                   selectionText[i] = (selectionText[i] === "") ? "" : "- [x] " + selectionText[i];
               }
               cm.replaceSelection(selectionText.join("\n"));
           }
       }else if(name === "video"){
           //插入视频
           window.editor.videoDialog();
       }else {
           var action = window.editor.toolbarHandlers[name];

           if (action !== "undefined") {
               $.proxy(action, window.editor)();
               window.editor.focus();
           }
       }
   }) ;
});