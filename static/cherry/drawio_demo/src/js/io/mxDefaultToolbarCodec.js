/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxDefaultToolbarCodec
 *
 * Custom codec for configuring <mxDefaultToolbar>s. This class is created
 * and registered dynamically at load time and used implicitely via
 * <mxCodec> and the <mxCodecRegistry>. This codec only reads configuration
 * data for existing toolbars handlers, it does not encode or create toolbars.
 */
var mxDefaultToolbarCodec = mxCodecRegistry.register(function()
{
	var codec = new mxObjectCodec(new mxDefaultToolbar());

	/**
	 * Function: encode
	 *
	 * Returns null.
	 */
	codec.encode = function(enc, obj)
	{
		return null;
	};
	
	/**
	 * Function: decode
	 *
	 * Reads a sequence of the following child nodes
	 * and attributes:
	 *
	 * Child Nodes:
	 *
	 * add - Adds a new item to the toolbar. See below for attributes.
	 * separator - Adds a vertical separator. No attributes.
	 * hr - Adds a horizontal separator. No attributes.
	 * br - Adds a linefeed. No attributes. 
	 *
	 * Attributes:
	 *
	 * as - Resource key for the label.
	 * action - Name of the action to execute in enclosing editor.
	 * mode - Modename (see below).
	 * template - Template name for cell insertion.
	 * style - Optional style to override the template style.
	 * icon - Icon (relative/absolute URL).
	 * pressedIcon - Optional icon for pressed state (relative/absolute URL).
	 * id - Optional ID to be used for the created DOM element.
	 * toggle - Optional 0 or 1 to disable toggling of the element. Default is
	 * 1 (true).
	 *
	 * The action, mode and template attributes are mutually exclusive. The
	 * style can only be used with the template attribute. The add node may
	 * contain another sequence of add nodes with as and action attributes
	 * to create a combo box in the toolbar. If the icon is specified then
	 * a list of the child node is expected to have its template attribute
	 * set and the action is ignored instead.
	 * 
	 * Nodes with a specified template may define a function to be used for
	 * inserting the cloned template into the graph. Here is an example of such
	 * a node:
	 * 
	 * (code)
	 * <add as="Swimlane" template="swimlane" icon="images/swimlane.gif"><![CDATA[
	 *   function (editor, cell, evt, targetCell)
	 *   {
	 *     var pt = mxUtils.convertPoint(
	 *       editor.graph.container, mxEvent.getClientX(evt),
	 *         mxEvent.getClientY(evt));
	 *     return editor.addVertex(targetCell, cell, pt.x, pt.y);
	 *   }
	 * ]]></add>
	 * (end)
	 * 
	 * In the above function, editor is the enclosing <mxEditor> instance, cell
	 * is the clone of the template, evt is the mouse event that represents the
	 * drop and targetCell is the cell under the mousepointer where the drop
	 * occurred. The targetCell is retrieved using <mxGraph.getCellAt>.
	 *
	 * Futhermore, nodes with the mode attribute may define a function to
	 * be executed upon selection of the respective toolbar icon. In the
	 * example below, the default edge style is set when this specific
	 * connect-mode is activated:
	 *
	 * (code)
	 * <add as="connect" mode="connect"><![CDATA[
	 *   function (editor)
	 *   {
	 *     if (editor.defaultEdge != null)
	 *     {
	 *       editor.defaultEdge.style = 'straightEdge';
	 *     }
	 *   }
	 * ]]></add>
	 * (end)
	 * 
	 * Both functions require <mxDefaultToolbarCodec.allowEval> to be set to true.
	 *
	 * Modes:
	 *
	 * select - Left mouse button used for rubberband- & cell-selection.
	 * connect - Allows connecting vertices by inserting new edges.
	 * pan - Disables selection and switches to panning on the left button.
	 *
	 * Example:
	 *
	 * To add items to the toolbar:
	 * 
	 * (code)
	 * <mxDefaultToolbar as="toolbar">
	 *   <add as="save" action="save" icon="images/save.gif"/>
	 *   <br/><hr/>
	 *   <add as="select" mode="select" icon="images/select.gif"/>
	 *   <add as="connect" mode="connect" icon="images/connect.gif"/>
	 * </mxDefaultToolbar>
	 * (end)
	 */
	codec.decode = function(dec, node, into)
	{
		if (into != null)
		{
			var editor = into.editor;
			node = node.firstChild;
			
			while (node != null)
			{
				if (node.nodeType == mxConstants.NODETYPE_ELEMENT)
				{
					if (!this.processInclude(dec, node, into))
					{
						if (node.nodeName == 'separator')
						{
							into.addSeparator();
						}
						else if (node.nodeName == 'br')
						{
							into.toolbar.addBreak();
						}
						else if (node.nodeName == 'hr')
						{
							into.toolbar.addLine();
						}
						else if (node.nodeName == 'add')
						{
							var as = node.getAttribute('as');
							as = mxResources.get(as) || as;
							var icon = node.getAttribute('icon');
							var pressedIcon = node.getAttribute('pressedIcon');
							var action = node.getAttribute('action');
							var mode = node.getAttribute('mode');
							var template = node.getAttribute('template');
							var toggle = node.getAttribute('toggle') != '0';
							var text = mxUtils.getTextContent(node);
							var elt = null;

							if (action != null)
							{
								elt = into.addItem(as, icon, action, pressedIcon);
							}
							else if (mode != null)
							{
								var funct = (mxDefaultToolbarCodec.allowEval) ? mxUtils.eval(text) : null;
								elt = into.addMode(as, icon, mode, pressedIcon, funct);
							}
							else if (template != null || (text != null && text.length > 0))
							{
								var cell = editor.templates[template];
								var style = node.getAttribute('style');
								
								if (cell != null && style != null)
								{
									cell = editor.graph.cloneCell(cell);
									cell.setStyle(style);
								}
								
								var insertFunction = null;
								
								if (text != null && text.length > 0 && mxDefaultToolbarCodec.allowEval)
								{
									insertFunction = mxUtils.eval(text);
								}
								
								elt = into.addPrototype(as, icon, cell, pressedIcon, insertFunction, toggle);
							}
							else
							{
								var children = mxUtils.getChildNodes(node);
								
								if (children.length > 0)
								{
									if (icon == null)
									{
										var combo = into.addActionCombo(as);
										
										for (var i=0; i<children.length; i++)
										{
											var child = children[i];
											
											if (child.nodeName == 'separator')
											{
												into.addOption(combo, '---');
											}
											else if (child.nodeName == 'add')
											{
												var lab = child.getAttribute('as');
												var act = child.getAttribute('action');
												into.addActionOption(combo, lab, act);
											}
										}
									}
									else
									{
										var select = null;
										var create = function()
										{
											var template = editor.templates[select.value];
											
											if (template != null)
											{
												var clone = template.clone();
												var style = select.options[select.selectedIndex].cellStyle;
												
												if (style != null)
												{
													clone.setStyle(style);
												}
												
												return clone;
											}
											else
											{
												mxLog.warn('Template '+template+' not found');
											}
											
											return null;
										};
										
										var img = into.addPrototype(as, icon, create, null, null, toggle);
										select = into.addCombo();
										
										// Selects the toolbar icon if a selection change
										// is made in the corresponding combobox.
										mxEvent.addListener(select, 'change', function()
										{
											into.toolbar.selectMode(img, function(evt)
											{
												var pt = mxUtils.convertPoint(editor.graph.container,
													mxEvent.getClientX(evt), mxEvent.getClientY(evt));
												
												return editor.addVertex(null, funct(), pt.x, pt.y);
											});
											
											into.toolbar.noReset = false;
										});
										
										// Adds the entries to the combobox
										for (var i=0; i<children.length; i++)
										{
											var child = children[i];
											
											if (child.nodeName == 'separator')
											{
												into.addOption(select, '---');
											}
											else if (child.nodeName == 'add')
											{
												var lab = child.getAttribute('as');
												var tmp = child.getAttribute('template');
												var option = into.addOption(select, lab, tmp || template);
												option.cellStyle = child.getAttribute('style');
											}
										}
										
									}
								}
							}
							
							// Assigns an ID to the created element to access it later.
							if (elt != null)
							{
								var id = node.getAttribute('id');
								
								if (id != null && id.length > 0)
								{
									elt.setAttribute('id', id);
								}
							}
						}
					}
				}
				
				node = node.nextSibling;
			}
		}
		
		return into;
	};
	
	// Returns the codec into the registry
	return codec;

}());

/**
 * Variable: allowEval
 * 
 * Static global switch that specifies if the use of eval is allowed for
 * evaluating text content. Default is true. Set this to false if stylesheets
 * may contain user input
 */
mxDefaultToolbarCodec.allowEval = true;
