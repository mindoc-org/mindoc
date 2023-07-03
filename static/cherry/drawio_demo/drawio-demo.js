// Extends EditorUi to update I/O action states based on availability of backend
(function()
{
	var editorUiInit = EditorUi.prototype.init;
	
	EditorUi.prototype.init = function()
	{
		editorUiInit.apply(this, arguments);
	};
	
	// Adds required resources (disables loading of fallback properties, this can only
	// be used if we know that all keys are defined in the language specific file)
	mxResources.loadDefaultBundle = false;
	var bundle = mxResources.getDefaultBundle(mxLanguage);
	
	// Fixes possible asynchronous requests
	mxUtils.getAll([bundle, './drawio_demo/theme/default.xml'], function(xhr)
	{
		// Adds bundle text to resources
		mxResources.parse(xhr[0].getText());
		
		// Configures the default graph theme
		var themes = new Object();
		themes[Graph.prototype.defaultThemeName] = xhr[1].getDocumentElement(); 
		
		// Main
    window.editorUIInstance = new EditorUi(new Editor(false, themes));
    
    try {
      addPostMessageListener(editorUIInstance.editor);
    } catch (error) {
      console.log(error);
    }
    window.parent.postMessage({eventName: 'ready', value: ''}, '*');

	}, function()
	{
		document.body.innerHTML = '<center style="margin-top:10%;">Error loading resource files. Please check browser console.</center>';
	});
})();

function addPostMessageListener(graphEditor) {
  window.addEventListener('message', function(event) {
    if(!event.data || !event.data.eventName) {
        return 
    }
    switch (event.data.eventName) {
      case 'setData':
        var value = event.data.value;
        var doc = mxUtils.parseXml(value);
        var documentName = 'cherry-drawio-' + new Date().getTime();
        editorUIInstance.editor.setGraphXml(null);
        graphEditor.graph.importGraphModel(doc.documentElement);
        graphEditor.setFilename(documentName);
        window.parent.postMessage({eventName: 'setData:success', value: ''}, '*');
        break;
      case 'getData':
        editorUIInstance.editor.graph.stopEditing();
        var xmlData = mxUtils.getXml(editorUIInstance.editor.getGraphXml());
        editorUIInstance.exportImage(1, "#ffffff", true, null, true, 50, null, "png", function(base64, filename){
          window.parent.postMessage({
            mceAction: 'getData:success',
            eventName: 'getData:success',
            value: {
              xmlData: xmlData,
              base64: base64,
            }
          }, '*');
        })
        break;
      case 'ready?':
        window.parent.postMessage({eventName: 'ready', value: ''}, '*');
        break;
      default:
        break;
    }
  }); 
}
