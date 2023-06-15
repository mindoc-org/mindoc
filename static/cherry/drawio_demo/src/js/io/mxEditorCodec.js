/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
mxCodecRegistry.register(function()
{
	/**
	 * Class: mxEditorCodec
	 *
	 * Codec for <mxEditor>s. This class is created and registered
	 * dynamically at load time and used implicitely via <mxCodec>
	 * and the <mxCodecRegistry>.
	 *
	 * Transient Fields:
	 *
	 * - modified
	 * - lastSnapshot
	 * - ignoredChanges
	 * - undoManager
	 * - graphContainer
	 * - toolbarContainer
	 */
	var codec = new mxObjectCodec(new mxEditor(),
		['modified', 'lastSnapshot', 'ignoredChanges',
		'undoManager', 'graphContainer', 'toolbarContainer']);

	/**
	 * Function: beforeDecode
	 *
	 * Decodes the ui-part of the configuration node by reading
	 * a sequence of the following child nodes and attributes
	 * and passes the control to the default decoding mechanism:
	 *
	 * Child Nodes:
	 *
	 * stylesheet - Adds a CSS stylesheet to the document.
	 * resource - Adds the basename of a resource bundle.
	 * add - Creates or configures a known UI element.
	 *
	 * These elements may appear in any order given that the
	 * graph UI element is added before the toolbar element
	 * (see Known Keys).
	 *
	 * Attributes:
	 *
	 * as - Key for the UI element (see below).
	 * element - ID for the element in the document.
	 * style - CSS style to be used for the element or window.
	 * x - X coordinate for the new window.
	 * y - Y coordinate for the new window.
	 * width - Width for the new window.
	 * height - Optional height for the new window.
	 * name - Name of the stylesheet (absolute/relative URL).
	 * basename - Basename of the resource bundle (see <mxResources>).
	 *
	 * The x, y, width and height attributes are used to create a new
	 * <mxWindow> if the element attribute is not specified in an add
	 * node. The name and basename are only used in the stylesheet and
	 * resource nodes, respectively.
	 *
	 * Known Keys:
	 *
	 * graph - Main graph element (see <mxEditor.setGraphContainer>).
	 * title - Title element (see <mxEditor.setTitleContainer>).
	 * toolbar - Toolbar element (see <mxEditor.setToolbarContainer>).
	 * status - Status bar element (see <mxEditor.setStatusContainer>).
	 *
	 * Example:
	 *
	 * (code)
	 * <ui>
	 *   <stylesheet name="css/process.css"/>
	 *   <resource basename="resources/app"/>
	 *   <add as="graph" element="graph"
	 *     style="left:70px;right:20px;top:20px;bottom:40px"/>
	 *   <add as="status" element="status"/>
	 *   <add as="toolbar" x="10" y="20" width="54"/>
	 * </ui>
	 * (end)
	 */
	codec.afterDecode = function(dec, node, obj)
	{
		// Assigns the specified templates for edges
		var defaultEdge = node.getAttribute('defaultEdge');
		
		if (defaultEdge != null)
		{
			node.removeAttribute('defaultEdge');
			obj.defaultEdge = obj.templates[defaultEdge];
		}

		// Assigns the specified templates for groups
		var defaultGroup = node.getAttribute('defaultGroup');
		
		if (defaultGroup != null)
		{
			node.removeAttribute('defaultGroup');
			obj.defaultGroup = obj.templates[defaultGroup];
		}

		return obj;
	};
	
	/**
	 * Function: decodeChild
	 * 
	 * Overrides decode child to handle special child nodes.
	 */	
	codec.decodeChild = function(dec, child, obj)
	{
		if (child.nodeName == 'Array')
		{
			var role = child.getAttribute('as');
			
			if (role == 'templates')
			{
				this.decodeTemplates(dec, child, obj);
				return;
			}
		}
		else if (child.nodeName == 'ui')
		{
			this.decodeUi(dec, child, obj);
			return;
		}
		
		mxObjectCodec.prototype.decodeChild.apply(this, arguments);
	};
		
	/**
	 * Function: decodeTemplates
	 *
	 * Decodes the cells from the given node as templates.
	 */
	codec.decodeUi = function(dec, node, editor)
	{
		var tmp = node.firstChild;
		while (tmp != null)
		{
			if (tmp.nodeName == 'add')
			{
				var as = tmp.getAttribute('as');
				var elt = tmp.getAttribute('element');
				var style = tmp.getAttribute('style');
				var element = null;

				if (elt != null)
				{
					element = document.getElementById(elt);
					
					if (element != null && style != null)
					{
						element.style.cssText += ';' + style;
					}
				}
				else
				{
					var x = parseInt(tmp.getAttribute('x'));
					var y = parseInt(tmp.getAttribute('y'));
					var width = tmp.getAttribute('width');
					var height = tmp.getAttribute('height');

					// Creates a new window around the element
					element = document.createElement('div');
					element.style.cssText = style;
					
					var wnd = new mxWindow(mxResources.get(as) || as,
						element, x, y, width, height, false, true);
					wnd.setVisible(true);
				}
				
				// TODO: Make more generic
				if (as == 'graph')
				{
					editor.setGraphContainer(element);
				}
				else if (as == 'toolbar')
				{
					editor.setToolbarContainer(element);
				}
				else if (as == 'title')
				{
					editor.setTitleContainer(element);
				}
				else if (as == 'status')
				{
					editor.setStatusContainer(element);
				}
				else if (as == 'map')
				{
					editor.setMapContainer(element);
				}
			}
			else if (tmp.nodeName == 'resource')
			{
				mxResources.add(tmp.getAttribute('basename'));
			}
			else if (tmp.nodeName == 'stylesheet')
			{
				mxClient.link('stylesheet', tmp.getAttribute('name'));
			}
			
			tmp = tmp.nextSibling;
		}	
	};
	
	/**
	 * Function: decodeTemplates
	 *
	 * Decodes the cells from the given node as templates.
	 */
	codec.decodeTemplates = function(dec, node, editor)
	{
		if (editor.templates == null)
		{
			editor.templates = [];
		}
		
		var children = mxUtils.getChildNodes(node);
		for (var j=0; j<children.length; j++)
		{
			var name = children[j].getAttribute('as');
			var child = children[j].firstChild;
			
			while (child != null && child.nodeType != 1)
			{
				child = child.nextSibling;
			}
			
			if (child != null)
			{
				// LATER: Only single cells means you need
				// to group multiple cells within another
				// cell. This should be changed to support
				// arrays of cells, or the wrapper must
				// be automatically handled in this class.
				editor.templates[name] = dec.decodeCell(child);
			}
		}
	};
	
	// Returns the codec into the registry
	return codec;

}());
