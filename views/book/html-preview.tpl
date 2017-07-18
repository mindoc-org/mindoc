
<!DOCTYPE html>
<html lang="zh">
    <head>
        <meta charset="utf-8" />
        <title>{{.Model.BookName}}</title>
        <link rel="stylesheet" href="{{.BaseUrl}}/static/css/html.previre.css" />
        <link rel="stylesheet" href="{{.BaseUrl}}/static/editor.md/css/editormd.preview.css" />
    </head>
    <body>
        <div id="layout">
            <div id="test-editormd-view">
               <textarea style="display:none;" name="test-editormd-markdown-doc">
               {{.MdText}}
               </textarea>               
            </div>
        </div>
        <!-- 
        <script src="{{.BaseUrl}}/static/editor.md/examples/js/zepto.min.js"></script>
		<script>		
			var jQuery = Zepto;  // 为了避免修改flowChart.js和sequence-diagram.js的源码，所以使用Zepto.js时想支持flowChart/sequenceDiagram就得加上这一句
		</script> 
        -->
        <script src="{{.BaseUrl}}/static/editor.md/examples/js/jquery.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/marked.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/prettify.min.js"></script>
        
        <script src="{{.BaseUrl}}/static/editor.md/lib/raphael.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/underscore.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/sequence-diagram.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/flowchart.min.js"></script>
        <script src="{{.BaseUrl}}/static/editor.md/lib/jquery.flowchart.min.js"></script>

        <script src="{{.BaseUrl}}/static/editor.md/editormd.js"></script>
        <script type="text/javascript">
            $(function() {
                var testEditormdView;
                testEditormdView = editormd.markdownToHTML("test-editormd-view", {
                    htmlDecode      : "style,script,iframe",  // you can filter tags decode
                    emoji           : true,
                    taskList        : true,
                    tex             : true,  // 默认不解析
                    flowChart       : true,  // 默认不解析
                    sequenceDiagram : true,  // 默认不解析
                });
            });
        </script>
        <!-- 
        html, err := c.ExecuteViewPathTemplate("book/html-preview.tpl", map[string]interface{}{"Model": bookResult, "MdText": mdText, "BaseUrl": baseurl}) 
        -->
    </body>
</html>