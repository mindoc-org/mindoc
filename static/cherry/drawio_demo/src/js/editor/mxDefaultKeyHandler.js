/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxDefaultKeyHandler
 *
 * Binds keycodes to actionnames in an editor. This aggregates an internal
 * <handler> and extends the implementation of <mxKeyHandler.escape> to not
 * only cancel the editing, but also hide the properties dialog and fire an
 * <mxEditor.escape> event via <editor>. An instance of this class is created
 * by <mxEditor> and stored in <mxEditor.keyHandler>.
 * 
 * Example:
 * 
 * Bind the delete key to the delete action in an existing editor.
 * 
 * (code)
 * var keyHandler = new mxDefaultKeyHandler(editor);
 * keyHandler.bindAction(46, 'delete');
 * (end)
 *
 * Codec:
 * 
 * This class uses the <mxDefaultKeyHandlerCodec> to read configuration
 * data into an existing instance. See <mxDefaultKeyHandlerCodec> for a
 * description of the configuration format.
 * 
 * Keycodes:
 * 
 * See <mxKeyHandler>.
 * 
 * An <mxEvent.ESCAPE> event is fired via the editor if the escape key is
 * pressed.
 * 
 * Constructor: mxDefaultKeyHandler
 *
 * Constructs a new default key handler for the <mxEditor.graph> in the
 * given <mxEditor>. (The editor may be null if a prototypical instance for
 * a <mxDefaultKeyHandlerCodec> is created.)
 * 
 * Parameters:
 * 
 * editor - Reference to the enclosing <mxEditor>.
 */
function mxDefaultKeyHandler(editor)
{
	if (editor != null)
	{
		this.editor = editor;
		this.handler = new mxKeyHandler(editor.graph);
		
		// Extends the escape function of the internal key
		// handle to hide the properties dialog and fire
		// the escape event via the editor instance
		var old = this.handler.escape;
		
		this.handler.escape = function(evt)
		{
			old.apply(this, arguments);
			editor.hideProperties();
			editor.fireEvent(new mxEventObject(mxEvent.ESCAPE, 'event', evt));
		};
	}
};
	
/**
 * Variable: editor
 *
 * Reference to the enclosing <mxEditor>.
 */
mxDefaultKeyHandler.prototype.editor = null;

/**
 * Variable: handler
 *
 * Holds the <mxKeyHandler> for key event handling.
 */
mxDefaultKeyHandler.prototype.handler = null;

/**
 * Function: bindAction
 *
 * Binds the specified keycode to the given action in <editor>. The
 * optional control flag specifies if the control key must be pressed
 * to trigger the action.
 *
 * Parameters:
 *
 * code - Integer that specifies the keycode.
 * action - Name of the action to execute in <editor>.
 * control - Optional boolean that specifies if control must be pressed.
 * Default is false.
 */
mxDefaultKeyHandler.prototype.bindAction = function (code, action, control)
{
	var keyHandler = mxUtils.bind(this, function()
	{
		this.editor.execute(action);
	});

	// Binds the function to control-down keycode
	if (control)
	{
		this.handler.bindControlKey(code, keyHandler);
	}

	// Binds the function to the normal keycode
	else
	{
		this.handler.bindKey(code, keyHandler);				
	}
};

/**
 * Function: destroy
 *
 * Destroys the <handler> associated with this object. This does normally
 * not need to be called, the <handler> is destroyed automatically when the
 * window unloads (in IE) by <mxEditor>.
 */
mxDefaultKeyHandler.prototype.destroy = function ()
{
	this.handler.destroy();
	this.handler = null;
};
