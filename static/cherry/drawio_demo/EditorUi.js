/**
 * Copyright (c) 2006-2012, JGraph Ltd
 */
/**
 * Constructs a new graph editor
 */
EditorUi = function(editor, container, lightbox)
{
	mxEventSource.call(this);
	
	this.destroyFunctions = [];
	this.editor = editor || new Editor();
	this.container = container || document.body;
	
	var graph = this.editor.graph;
	graph.lightbox = lightbox;

	// Faster scrollwheel zoom is possible with CSS transforms
	if (graph.useCssTransforms)
	{
		this.lazyZoomDelay = 0;
	}
	
	// Pre-fetches submenu image or replaces with embedded image if supported
	if (mxClient.IS_SVG)
	{
		mxPopupMenu.prototype.submenuImage = 'data:image/gif;base64,R0lGODlhCQAJAIAAAP///zMzMyH5BAEAAAAALAAAAAAJAAkAAAIPhI8WebHsHopSOVgb26AAADs=';
	}
	else
	{
		new Image().src = mxPopupMenu.prototype.submenuImage;
	}

	// Pre-fetches connect image
	if (!mxClient.IS_SVG && mxConnectionHandler.prototype.connectImage != null)
	{
		new Image().src = mxConnectionHandler.prototype.connectImage.src;
	}
	
	// Disables graph and forced panning in chromeless mode
	if (this.editor.chromeless && !this.editor.editable)
	{
		this.footerHeight = 0;
		graph.isEnabled = function() { return false; };
		graph.panningHandler.isForcePanningEvent = function(me)
		{
			return !mxEvent.isPopupTrigger(me.getEvent());
		};
	}
	
    // Creates the user interface
	this.actions = new Actions(this);
	this.menus = this.createMenus();
	this.createDivs();
	this.createUi();
	this.refresh();
	
	// Disables HTML and text selection
	var textEditing =  mxUtils.bind(this, function(evt)
	{
		if (evt == null)
		{
			evt = window.event;
		}
		
		return (this.isSelectionAllowed(evt) || graph.isEditing());
	});

	// Disables text selection while not editing and no dialog visible
	if (this.container == document.body)
	{
		this.toolbarContainer.onselectstart = textEditing;
		this.toolbarContainer.onmousedown = textEditing;
		this.diagramContainer.onselectstart = textEditing;
		this.diagramContainer.onmousedown = textEditing;
		this.sidebarContainer.onselectstart = textEditing;
		this.sidebarContainer.onmousedown = textEditing;
		this.formatContainer.onselectstart = textEditing;
		this.formatContainer.onmousedown = textEditing;
		this.footerContainer.onselectstart = textEditing;
		this.footerContainer.onmousedown = textEditing;
		
		if (this.tabContainer != null)
		{
			// Mouse down is needed for drag and drop
			this.tabContainer.onselectstart = textEditing;
		}
	}
	
	// And uses built-in context menu while editing
	if (!this.editor.chromeless || this.editor.editable)
	{
		// Allows context menu for links in hints
		var linkHandler = function(evt)
		{
			var source = mxEvent.getSource(evt);
			
			if (source.nodeName == 'A')
			{
				while (source != null)
				{
					if (source.className == 'geHint')
					{
						return true;
					}
					
					source = source.parentNode;
				}
			}
			
			return textEditing(evt);
		};
		
		if (mxClient.IS_IE && (typeof(document.documentMode) === 'undefined' || document.documentMode < 9))
		{
			mxEvent.addListener(this.diagramContainer, 'contextmenu', linkHandler);
		}
		else
		{
			// Allows browser context menu outside of diagram and sidebar
			this.diagramContainer.oncontextmenu = linkHandler;
		}
	}
	else
	{
		graph.panningHandler.usePopupTrigger = false;
	}

	// Contains the main graph instance inside the given panel
	graph.init(this.diagramContainer);

    // Improves line wrapping for in-place editor
    if (mxClient.IS_SVG && graph.view.getDrawPane() != null)
    {
        var root = graph.view.getDrawPane().ownerSVGElement;
        
        if (root != null)
        {
            root.style.position = 'absolute';
        }
    }
    
	// Creates hover icons
	this.hoverIcons = this.createHoverIcons();
	
	// Adds tooltip when mouse is over scrollbars to show space-drag panning option
	mxEvent.addListener(this.diagramContainer, 'mousemove', mxUtils.bind(this, function(evt)
	{
		var off = mxUtils.getOffset(this.diagramContainer);
		
		if (mxEvent.getClientX(evt) - off.x - this.diagramContainer.clientWidth > 0 ||
			mxEvent.getClientY(evt) - off.y - this.diagramContainer.clientHeight > 0)
		{
			this.diagramContainer.setAttribute('title', mxResources.get('panTooltip'));
		}
		else
		{
			this.diagramContainer.removeAttribute('title');
		}
	}));

   	// Escape key hides dialogs, adds space+drag panning
	var spaceKeyPressed = false;
	
	// Overrides hovericons to disable while space key is pressed
	var hoverIconsIsResetEvent = this.hoverIcons.isResetEvent;
	
	this.hoverIcons.isResetEvent = function(evt, allowShift)
	{
		return spaceKeyPressed || hoverIconsIsResetEvent.apply(this, arguments);
	};
	
	this.keydownHandler = mxUtils.bind(this, function(evt)
	{
		if (evt.which == 32 /* Space */)
		{
			spaceKeyPressed = true;
			this.hoverIcons.reset();
			graph.container.style.cursor = 'move';
			
			// Disables scroll after space keystroke with scrollbars
			if (!graph.isEditing() && mxEvent.getSource(evt) == graph.container)
			{
				mxEvent.consume(evt);
			}
		}
		else if (!mxEvent.isConsumed(evt) && evt.keyCode == 27 /* Escape */)
		{
			this.hideDialog();
		}
	});
   	
	mxEvent.addListener(document, 'keydown', this.keydownHandler);
	
	this.keyupHandler = mxUtils.bind(this, function(evt)
	{
		graph.container.style.cursor = '';
		spaceKeyPressed = false;
	});

	mxEvent.addListener(document, 'keyup', this.keyupHandler);
    
    // Forces panning for middle and right mouse buttons
	var panningHandlerIsForcePanningEvent = graph.panningHandler.isForcePanningEvent;
	graph.panningHandler.isForcePanningEvent = function(me)
	{
		// Ctrl+left button is reported as right button in FF on Mac
		return panningHandlerIsForcePanningEvent.apply(this, arguments) ||
			spaceKeyPressed || (mxEvent.isMouseEvent(me.getEvent()) &&
			(this.usePopupTrigger || !mxEvent.isPopupTrigger(me.getEvent())) &&
			((!mxEvent.isControlDown(me.getEvent()) &&
			mxEvent.isRightMouseButton(me.getEvent())) ||
			mxEvent.isMiddleMouseButton(me.getEvent())));
	};

	// Ctrl/Cmd+Enter applies editing value except in Safari where Ctrl+Enter creates
	// a new line (while Enter creates a new paragraph and Shift+Enter stops)
	var cellEditorIsStopEditingEvent = graph.cellEditor.isStopEditingEvent;
	graph.cellEditor.isStopEditingEvent = function(evt)
	{
		return cellEditorIsStopEditingEvent.apply(this, arguments) ||
			(evt.keyCode == 13 && ((!mxClient.IS_SF && mxEvent.isControlDown(evt)) ||
			(mxClient.IS_MAC && mxEvent.isMetaDown(evt)) ||
			(mxClient.IS_SF && mxEvent.isShiftDown(evt))));
	};
	
	// Switches toolbar for text editing
	var textMode = false;
	var fontMenu = null;
	var sizeMenu = null;
	var nodes = null;
	
	var updateToolbar = mxUtils.bind(this, function()
	{
		if (this.toolbar != null && textMode != graph.cellEditor.isContentEditing())
		{
			var node = this.toolbar.container.firstChild;
			var newNodes = [];
			
			while (node != null)
			{
				var tmp = node.nextSibling;
				
				if (mxUtils.indexOf(this.toolbar.staticElements, node) < 0)
				{
					node.parentNode.removeChild(node);
					newNodes.push(node);
				}
				
				node = tmp;
			}
			
			// Saves references to special items
			var tmp1 = this.toolbar.fontMenu;
			var tmp2 = this.toolbar.sizeMenu;
			
			if (nodes == null)
			{
				this.toolbar.createTextToolbar();
			}
			else
			{
				for (var i = 0; i < nodes.length; i++)
				{
					this.toolbar.container.appendChild(nodes[i]);
				}
				
				// Restores references to special items
				this.toolbar.fontMenu = fontMenu;
				this.toolbar.sizeMenu = sizeMenu;
			}
			
			textMode = graph.cellEditor.isContentEditing();
			fontMenu = tmp1;
			sizeMenu = tmp2;
			nodes = newNodes;
		}
	});

	var ui = this;
	
	// Overrides cell editor to update toolbar
	var cellEditorStartEditing = graph.cellEditor.startEditing;
	graph.cellEditor.startEditing = function()
	{
		cellEditorStartEditing.apply(this, arguments);
		updateToolbar();
		
		if (graph.cellEditor.isContentEditing())
		{
			var updating = false;
			
			var updateCssHandler = function()
			{
				if (!updating)
				{
					updating = true;
				
					window.setTimeout(function()
					{
						var selectedElement = graph.getSelectedElement();
						var node = selectedElement;
						
						while (node != null && node.nodeType != mxConstants.NODETYPE_ELEMENT)
						{
							node = node.parentNode;
						}
						
						if (node != null)
						{
							var css = mxUtils.getCurrentStyle(node);
	
							if (css != null && ui.toolbar != null)
							{
								// Strips leading and trailing quotes
								var ff = css.fontFamily;
								
								if (ff.charAt(0) == '\'')
								{
									ff = ff.substring(1);
								}
								
								if (ff.charAt(ff.length - 1) == '\'')
								{
									ff = ff.substring(0, ff.length - 1);
								}
								
								ui.toolbar.setFontName(ff);
								ui.toolbar.setFontSize(parseInt(css.fontSize));
							}
						}
						
						updating = false;
					}, 0);
				}
			};
			
			mxEvent.addListener(graph.cellEditor.textarea, 'input', updateCssHandler)
			mxEvent.addListener(graph.cellEditor.textarea, 'touchend', updateCssHandler);
			mxEvent.addListener(graph.cellEditor.textarea, 'mouseup', updateCssHandler);
			mxEvent.addListener(graph.cellEditor.textarea, 'keyup', updateCssHandler);
			updateCssHandler();
		}
	};
	
	var cellEditorStopEditing = graph.cellEditor.stopEditing;
	graph.cellEditor.stopEditing = function(cell, trigger)
	{
		cellEditorStopEditing.apply(this, arguments);
		updateToolbar();
	};
	
    // Enables scrollbars and sets cursor style for the container
	graph.container.setAttribute('tabindex', '0');
   	graph.container.style.cursor = 'default';
    
	// Workaround for page scroll if embedded via iframe
	if (window.self === window.top && graph.container.parentNode != null)
	{
		try
		{
			graph.container.focus();
		}
		catch (e)
		{
			// ignores error in old versions of IE
		}
	}

   	// Keeps graph container focused on mouse down
   	var graphFireMouseEvent = graph.fireMouseEvent;
   	graph.fireMouseEvent = function(evtName, me, sender)
   	{
   		if (evtName == mxEvent.MOUSE_DOWN)
   		{
   			this.container.focus();
   		}
   		
   		graphFireMouseEvent.apply(this, arguments);
   	};

   	// Configures automatic expand on mouseover
	graph.popupMenuHandler.autoExpand = true;

    // Installs context menu
	if (this.menus != null)
	{
		graph.popupMenuHandler.factoryMethod = mxUtils.bind(this, function(menu, cell, evt)
		{
			this.menus.createPopupMenu(menu, cell, evt);
		});
	}
	
	// Hides context menu
	mxEvent.addGestureListeners(document, mxUtils.bind(this, function(evt)
	{
		graph.popupMenuHandler.hideMenu();
	}));

    // Create handler for key events
	this.keyHandler = this.createKeyHandler(editor);
    
	// Getter for key handler
	this.getKeyHandler = function()
	{
		return keyHandler;
	};
	
	// Stores the current style and assigns it to new cells
	var styles = ['rounded', 'shadow', 'glass', 'dashed', 'dashPattern', 'comic', 'labelBackgroundColor'];
	var connectStyles = ['shape', 'edgeStyle', 'curved', 'rounded', 'elbow', 'comic', 'jumpStyle', 'jumpSize'];
	
	// Note: Everything that is not in styles is ignored (styles is augmented below)
	this.setDefaultStyle = function(cell)
	{
		var state = graph.view.getState(cell);
		
		if (state != null)
		{
			// Ignores default styles
			var clone = cell.clone();
			clone.style = ''
			var defaultStyle = graph.getCellStyle(clone);
			var values = [];
			var keys = [];

			for (var key in state.style)
			{
				if (defaultStyle[key] != state.style[key])
				{
					values.push(state.style[key]);
					keys.push(key);
				}
			}
			
			// Handles special case for value "none"
			var cellStyle = graph.getModel().getStyle(state.cell);
			var tokens = (cellStyle != null) ? cellStyle.split(';') : [];
			
			for (var i = 0; i < tokens.length; i++)
			{
				var tmp = tokens[i];
		 		var pos = tmp.indexOf('=');
		 					 		
		 		if (pos >= 0)
		 		{
		 			var key = tmp.substring(0, pos);
		 			var value = tmp.substring(pos + 1);
		 			
		 			if (defaultStyle[key] != null && value == 'none')
		 			{
		 				values.push(value);
		 				keys.push(key);
		 			}
		 		}
			}

			// Resets current style
			if (graph.getModel().isEdge(state.cell))
			{
				graph.currentEdgeStyle = {};
			}
			else
			{
				graph.currentVertexStyle = {}
			}

			this.fireEvent(new mxEventObject('styleChanged', 'keys', keys, 'values', values, 'cells', [state.cell]));
		}
	};
	
	this.clearDefaultStyle = function()
	{
		graph.currentEdgeStyle = mxUtils.clone(graph.defaultEdgeStyle);
		graph.currentVertexStyle = mxUtils.clone(graph.defaultVertexStyle);
		
		// Updates UI
		this.fireEvent(new mxEventObject('styleChanged', 'keys', [], 'values', [], 'cells', []));
	};

	// Keys that should be ignored if the cell has a value (known: new default for all cells is html=1 so
    // for the html key this effecticely only works for edges inserted via the connection handler)
	var valueStyles = ['fontFamily', 'fontSize', 'fontColor'];
	
	// Keys that always update the current edge style regardless of selection
	var alwaysEdgeStyles = ['edgeStyle', 'startArrow', 'startFill', 'startSize', 'endArrow',
		'endFill', 'endSize', 'jettySize', 'orthogonalLoop'];
	
	// Keys that are ignored together (if one appears all are ignored)
	var keyGroups = [['startArrow', 'startFill', 'startSize', 'sourcePerimeterSpacing',
					'endArrow', 'endFill', 'endSize', 'targetPerimeterSpacing',
					'jettySize', 'orthogonalLoop'],
	                 ['strokeColor', 'strokeWidth'],
	                 ['fillColor', 'gradientColor'],
	                 valueStyles,
	                 ['opacity'],
	                 ['align'],
	                 ['html']];
	
	// Adds all keys used above to the styles array
	for (var i = 0; i < keyGroups.length; i++)
	{
		for (var j = 0; j < keyGroups[i].length; j++)
		{
			styles.push(keyGroups[i][j]);
		}
	}
	
	for (var i = 0; i < connectStyles.length; i++)
	{
		if (mxUtils.indexOf(styles, connectStyles[i]) < 0)
		{
			styles.push(connectStyles[i]);
		}
	}

	// Implements a global current style for edges and vertices that is applied to new cells
	var insertHandler = function(cells, asText)
	{
		var model = graph.getModel();
		
		model.beginUpdate();
		try
		{
			// Applies only basic text styles
			if (asText)
			{
				var edge = model.isEdge(cell);
				var current = (edge) ? graph.currentEdgeStyle : graph.currentVertexStyle;
				var textStyles = ['fontSize', 'fontFamily', 'fontColor'];
				
				for (var j = 0; j < textStyles.length; j++)
				{
					var value = current[textStyles[j]];
					
					if (value != null)
					{
						graph.setCellStyles(textStyles[j], value, cells);
					}
				}
			}
			else
			{
				for (var i = 0; i < cells.length; i++)
				{
					var cell = cells[i];

					// Removes styles defined in the cell style from the styles to be applied
					var cellStyle = model.getStyle(cell);
					var tokens = (cellStyle != null) ? cellStyle.split(';') : [];
					var appliedStyles = styles.slice();
					
					for (var j = 0; j < tokens.length; j++)
					{
						var tmp = tokens[j];
				 		var pos = tmp.indexOf('=');
				 					 		
				 		if (pos >= 0)
				 		{
				 			var key = tmp.substring(0, pos);
				 			var index = mxUtils.indexOf(appliedStyles, key);
				 			
				 			if (index >= 0)
				 			{
				 				appliedStyles.splice(index, 1);
				 			}
				 			
				 			// Handles special cases where one defined style ignores other styles
				 			for (var k = 0; k < keyGroups.length; k++)
				 			{
				 				var group = keyGroups[k];
				 				
				 				if (mxUtils.indexOf(group, key) >= 0)
				 				{
				 					for (var l = 0; l < group.length; l++)
				 					{
							 			var index2 = mxUtils.indexOf(appliedStyles, group[l]);
							 			
							 			if (index2 >= 0)
							 			{
							 				appliedStyles.splice(index2, 1);
							 			}
				 					}
				 				}
				 			}
				 		}
					}
	
					// Applies the current style to the cell
					var edge = model.isEdge(cell);
					var current = (edge) ? graph.currentEdgeStyle : graph.currentVertexStyle;
					var newStyle = model.getStyle(cell);
					
					for (var j = 0; j < appliedStyles.length; j++)
					{
						var key = appliedStyles[j];
						var styleValue = current[key];
	
						if (styleValue != null && (key != 'shape' || edge))
						{
							// Special case: Connect styles are not applied here but in the connection handler
							if (!edge || mxUtils.indexOf(connectStyles, key) < 0)
							{
								newStyle = mxUtils.setStyle(newStyle, key, styleValue);
							}
						}
					}
					
					model.setStyle(cell, newStyle);
				}
			}
		}
		finally
		{
			model.endUpdate();
		}
	};

	graph.addListener('cellsInserted', function(sender, evt)
	{
		insertHandler(evt.getProperty('cells'));
	});
	
	graph.addListener('textInserted', function(sender, evt)
	{
		insertHandler(evt.getProperty('cells'), true);
	});
	
	graph.connectionHandler.addListener(mxEvent.CONNECT, function(sender, evt)
	{
		var cells = [evt.getProperty('cell')];
		
		if (evt.getProperty('terminalInserted'))
		{
			cells.push(evt.getProperty('terminal'));
		}
		
		insertHandler(cells);
	});

	this.addListener('styleChanged', mxUtils.bind(this, function(sender, evt)
	{
		// Checks if edges and/or vertices were modified
		var cells = evt.getProperty('cells');
		var vertex = false;
		var edge = false;
		
		if (cells.length > 0)
		{
			for (var i = 0; i < cells.length; i++)
			{
				vertex = graph.getModel().isVertex(cells[i]) || vertex;
				edge = graph.getModel().isEdge(cells[i]) || edge;
				
				if (edge && vertex)
				{
					break;
				}
			}
		}
		else
		{
			vertex = true;
			edge = true;
		}
		
		var keys = evt.getProperty('keys');
		var values = evt.getProperty('values');

		for (var i = 0; i < keys.length; i++)
		{
			var common = mxUtils.indexOf(valueStyles, keys[i]) >= 0;
			
			// Ignores transparent stroke colors
			if (keys[i] != 'strokeColor' || (values[i] != null && values[i] != 'none'))
			{
				// Special case: Edge style and shape
				if (mxUtils.indexOf(connectStyles, keys[i]) >= 0)
				{
					if (edge || mxUtils.indexOf(alwaysEdgeStyles, keys[i]) >= 0)
					{
						if (values[i] == null)
						{
							delete graph.currentEdgeStyle[keys[i]];
						}
						else
						{
							graph.currentEdgeStyle[keys[i]] = values[i];
						}
					}
					// Uses style for vertex if defined in styles
					else if (vertex && mxUtils.indexOf(styles, keys[i]) >= 0)
					{
						if (values[i] == null)
						{
							delete graph.currentVertexStyle[keys[i]];
						}
						else
						{
							graph.currentVertexStyle[keys[i]] = values[i];
						}
					}
				}
				else if (mxUtils.indexOf(styles, keys[i]) >= 0)
				{
					if (vertex || common)
					{
						if (values[i] == null)
						{
							delete graph.currentVertexStyle[keys[i]];
						}
						else
						{
							graph.currentVertexStyle[keys[i]] = values[i];
						}
					}
					
					if (edge || common || mxUtils.indexOf(alwaysEdgeStyles, keys[i]) >= 0)
					{
						if (values[i] == null)
						{
							delete graph.currentEdgeStyle[keys[i]];
						}
						else
						{
							graph.currentEdgeStyle[keys[i]] = values[i];
						}
					}
				}
			}
		}
		
		if (this.toolbar != null)
		{
			this.toolbar.setFontName(graph.currentVertexStyle['fontFamily'] || Menus.prototype.defaultFont);
			this.toolbar.setFontSize(graph.currentVertexStyle['fontSize'] || Menus.prototype.defaultFontSize);
			
			if (this.toolbar.edgeStyleMenu != null)
			{
				// Updates toolbar icon for edge style
				var edgeStyleDiv = this.toolbar.edgeStyleMenu.getElementsByTagName('div')[0];

				if (graph.currentEdgeStyle['edgeStyle'] == 'orthogonalEdgeStyle' && graph.currentEdgeStyle['curved'] == '1')
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-curved';
				}
				else if (graph.currentEdgeStyle['edgeStyle'] == 'straight' || graph.currentEdgeStyle['edgeStyle'] == 'none' ||
						graph.currentEdgeStyle['edgeStyle'] == null)
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-straight';
				}
				else if (graph.currentEdgeStyle['edgeStyle'] == 'entityRelationEdgeStyle')
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-entity';
				}
				else if (graph.currentEdgeStyle['edgeStyle'] == 'elbowEdgeStyle')
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-' + ((graph.currentEdgeStyle['elbow'] == 'vertical') ?
						'verticalelbow' : 'horizontalelbow');
				}
				else if (graph.currentEdgeStyle['edgeStyle'] == 'isometricEdgeStyle')
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-' + ((graph.currentEdgeStyle['elbow'] == 'vertical') ?
						'verticalisometric' : 'horizontalisometric');
				}
				else
				{
					edgeStyleDiv.className = 'font-graph font-graph-geSprite-orthogonal';
				}
			}
			
			if (this.toolbar.edgeShapeMenu != null)
			{
				// Updates icon for edge shape
				var edgeShapeDiv = this.toolbar.edgeShapeMenu.getElementsByTagName('div')[0];
				
				if (graph.currentEdgeStyle['shape'] == 'link')
				{
					edgeShapeDiv.className = 'font-graph font-graph-geSprite-linkedge';
				}
				else if (graph.currentEdgeStyle['shape'] == 'flexArrow')
				{
					edgeShapeDiv.className = 'font-graph font-graph-geSprite-arrow';
				}
				else if (graph.currentEdgeStyle['shape'] == 'arrow')
				{
					edgeShapeDiv.className = 'font-graph font-graph-geSprite-simplearrow';
				}
				else
				{
					edgeShapeDiv.className = 'font-graph font-graph-geSprite-connection';
				}
			}
			
			// Updates icon for optinal line start shape
			if (this.toolbar.lineStartMenu != null)
			{
				var lineStartDiv = this.toolbar.lineStartMenu.getElementsByTagName('div')[0];
				
				lineStartDiv.className = this.getCssClassForMarker('start',
						graph.currentEdgeStyle['shape'], graph.currentEdgeStyle[mxConstants.STYLE_STARTARROW],
						mxUtils.getValue(graph.currentEdgeStyle, 'startFill', '1'));
			}

			// Updates icon for optinal line end shape
			if (this.toolbar.lineEndMenu != null)
			{
				var lineEndDiv = this.toolbar.lineEndMenu.getElementsByTagName('div')[0];
				
				lineEndDiv.className = this.getCssClassForMarker('end',
						graph.currentEdgeStyle['shape'], graph.currentEdgeStyle[mxConstants.STYLE_ENDARROW],
						mxUtils.getValue(graph.currentEdgeStyle, 'endFill', '1'));
			}
		}
	}));
	
	// Update font size and font family labels
	if (this.toolbar != null)
	{
		var update = mxUtils.bind(this, function()
		{
			var ff = graph.currentVertexStyle['fontFamily'] || 'Helvetica';
			var fs = String(graph.currentVertexStyle['fontSize'] || '12');
		    	var state = graph.getView().getState(graph.getSelectionCell());
		    	
		    	if (state != null)
		    	{
		    		ff = state.style[mxConstants.STYLE_FONTFAMILY] || ff;
		    		fs = state.style[mxConstants.STYLE_FONTSIZE] || fs;
		    		
		    		if (ff.length > 10)
		    		{
		    			ff = ff.substring(0, 8) + '...';
		    		}
		    	}
		    	
		    	this.toolbar.setFontName(ff);
		    	this.toolbar.setFontSize(fs);
		});
		
	    graph.getSelectionModel().addListener(mxEvent.CHANGE, update);
	    graph.getModel().addListener(mxEvent.CHANGE, update);
	}
	
	// Makes sure the current layer is visible when cells are added
	graph.addListener(mxEvent.CELLS_ADDED, function(sender, evt)
	{
		var cells = evt.getProperty('cells');
		var parent = evt.getProperty('parent');
		
		if (graph.getModel().isLayer(parent) && !graph.isCellVisible(parent) && cells != null && cells.length > 0)
		{
			graph.getModel().setVisible(parent, true);
		}
	});
	
	// Global handler to hide the current menu
	this.gestureHandler = mxUtils.bind(this, function(evt)
	{
		if (this.currentMenu != null && mxEvent.getSource(evt) != this.currentMenu.div)
		{
			this.hideCurrentMenu();
		}
	});
	
	mxEvent.addGestureListeners(document, this.gestureHandler);

	// Updates the editor UI after the window has been resized or the orientation changes
	// Timeout is workaround for old IE versions which have a delay for DOM client sizes.
	// Should not use delay > 0 to avoid handle multiple repaints during window resize
	this.resizeHandler = mxUtils.bind(this, function()
   	{
   		window.setTimeout(mxUtils.bind(this, function()
   		{
   			if (this.editor.graph != null)
   			{
   				this.refresh();
   			}
   		}), 0);
   	});
	
   	mxEvent.addListener(window, 'resize', this.resizeHandler);
   	
   	this.orientationChangeHandler = mxUtils.bind(this, function()
   	{
   		this.refresh();
   	});
   	
   	mxEvent.addListener(window, 'orientationchange', this.orientationChangeHandler);
   	
	// Workaround for bug on iOS see
	// http://stackoverflow.com/questions/19012135/ios-7-ipad-safari-landscape-innerheight-outerheight-layout-issue
	if (mxClient.IS_IOS && !window.navigator.standalone)
	{
		this.scrollHandler = mxUtils.bind(this, function()
	   	{
	   		window.scrollTo(0, 0);
	   	});
		
	   	mxEvent.addListener(window, 'scroll', this.scrollHandler);
	}

	/**
	 * Sets the initial scrollbar locations after a file was loaded.
	 */
	this.editor.addListener('resetGraphView', mxUtils.bind(this, function()
	{
		this.resetScrollbars();
	}));
	
	/**
	 * Repaints the grid.
	 */
	this.addListener('gridEnabledChanged', mxUtils.bind(this, function()
	{
		graph.view.validateBackground();
	}));
	
	this.addListener('backgroundColorChanged', mxUtils.bind(this, function()
	{
		graph.view.validateBackground();
	}));

	/**
	 * Repaints the grid.
	 */
	graph.addListener('gridSizeChanged', mxUtils.bind(this, function()
	{
		if (graph.isGridEnabled())
		{
			graph.view.validateBackground();
		}
	}));

   	// Resets UI, updates action and menu states
   	this.editor.resetGraph();
   	this.init();
   	this.open();
};

/**
 * Capability check for canvas export
 */
(function()
{
	EditorUi.prototype.useCanvasForExport = false;
	EditorUi.prototype.jpgSupported = false;
	
	// Checks if canvas is supported
	try
	{
		var cnv = document.createElement('canvas');
		EditorUi.prototype.canvasSupported = !!(cnv.getContext && cnv.getContext('2d'));
	}
	catch (e)
	{
		// ignore
	}
	
	try
	{
		var canvas = document.createElement('canvas');
		var img = new Image();
		
		// LATER: Capability check should not be async
		img.onload = function()
		{
			try
			{
		   		var ctx = canvas.getContext('2d');
		   		ctx.drawImage(img, 0, 0);

		   		// Works in Chrome, Firefox, Edge, Safari and Opera
				var result = canvas.toDataURL('image/png');
				EditorUi.prototype.useCanvasForExport = result != null && result.length > 6;
			}
			catch (e)
			{
				// ignore
			}
		};

		// Checks if SVG with foreignObject can be exported
		var svg = '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="1px" height="1px" version="1.1"><foreignObject pointer-events="all" width="1" height="1"><div xmlns="http://www.w3.org/1999/xhtml"></div></foreignObject></svg>';
		img.src = 'data:image/svg+xml;base64,' + btoa(unescape(encodeURIComponent(svg)));
	}
	catch (e)
	{
		// ignore
	}
	
	// Checks for client-side JPG support
	try
	{
	    var canvas = document.createElement('canvas');
	    canvas.width = canvas.height = 1;
	    var uri = canvas.toDataURL('image/jpeg');
	    
	    EditorUi.prototype.jpgSupported = (uri.match('image/jpeg') !== null);
	}
	catch (e)
	{
		// ignore
	}
})();

// Extends mxEventSource
mxUtils.extend(EditorUi, mxEventSource);

/**
 * Global config that specifies if the compact UI elements should be used.
 */
EditorUi.compactUi = true;

/**
 * Specifies the size of the split bar.
 */
EditorUi.prototype.splitSize = (mxClient.IS_TOUCH || mxClient.IS_POINTER) ? 12 : 8;

/**
 * Specifies the height of the menubar. Default is 34.
 */
EditorUi.prototype.menubarHeight = 30;

/**
 * Specifies the width of the format panel should be enabled. Default is true.
 */
EditorUi.prototype.formatEnabled = true;

/**
 * Specifies the width of the format panel. Default is 240.
 */
EditorUi.prototype.formatWidth = 240;

/**
 * Specifies the height of the toolbar. Default is 36.
 */
EditorUi.prototype.toolbarHeight = 34;

/**
 * Specifies the height of the footer. Default is 28.
 */
EditorUi.prototype.footerHeight = 28;

/**
 * Specifies the height of the optional sidebarFooterContainer. Default is 34.
 */
EditorUi.prototype.sidebarFooterHeight = 34;

/**
 * Specifies the position of the horizontal split bar. Default is 208 or 118 for
 * screen widths <= 640px.
 */
EditorUi.prototype.hsplitPosition = (screen.width <= 640) ? 118 : 208;

/**
 * Specifies if animations are allowed in <executeLayout>. Default is true.
 */
EditorUi.prototype.allowAnimation = true;

/**
 * Specifies if animations are allowed in <executeLayout>. Default is true.
 */
EditorUi.prototype.lightboxMaxFitScale = 2;

/**
 * Specifies if animations are allowed in <executeLayout>. Default is true.
 */
EditorUi.prototype.lightboxVerticalDivider = 4;

/**
 * Specifies if single click on horizontal split should collapse sidebar. Default is false.
 */
EditorUi.prototype.hsplitClickEnabled = false;


EditorUi.prototype.exportImage = function(scale, transparentBackground, ignoreSelection, addShadow, editable, border, noCrop, format, exportHandler)
{
	format = (format != null) ? format : 'png';
	
	var selectionEmpty = this.editor.graph.isSelectionEmpty();
	ignoreSelection = (ignoreSelection != null) ? ignoreSelection : selectionEmpty;

	// Caches images
	if (this.thumbImageCache == null)
	{
		this.thumbImageCache = new Object();
	}
	
	try
	{
	   	this.exportToCanvas(mxUtils.bind(this, function(canvas)
	   	{
	   		
	   		try
	   		{
	   			this.saveCanvas(canvas, (editable) ? this.getFileData(true, null,
			   				null, null, ignoreSelection) : null, format, exportHandler);
	   		}
	   		catch (e)
	   		{
				this.handleError(e);
	   		}
	   	}), null, this.thumbImageCache, null, mxUtils.bind(this, function(e)
	   	{
	   		this.handleError(e);
	   	}), null, ignoreSelection, scale || 1, transparentBackground,
	   		addShadow, null, null, border, noCrop);
	}
	catch (e)
	{
		this.handleError(e);
	}
};

EditorUi.prototype.handleError = function(error) {
	throw error;
}

/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.getFileData = function(forceXml, forceSvg, forceHtml, embeddedCallback, ignoreSelection, currentPage, node, compact, file)
{
	ignoreSelection = (ignoreSelection != null) ? ignoreSelection : true;
	currentPage = (currentPage != null) ? currentPage : false;
	
	node = this.editor.getGraphXml(ignoreSelection);
	var graph = this.editor.graph;
	
	
	var result = this.createFileData(node, graph, file, window.location.href,
		forceXml, forceSvg, forceHtml, embeddedCallback, ignoreSelection, compact);
	
	// Removes temporary graph from DOM
	if (graph != this.editor.graph)
	{
		graph.container.parentNode.removeChild(graph.container);
	}
	
	return result;
};

/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.createFileData = function(node, graph, file, url, forceXml, forceSvg, forceHtml, embeddedCallback, ignoreSelection, compact)
{
	graph = (graph != null) ? graph : this.editor.graph;
	forceXml = (forceXml != null) ? forceXml : false;
	ignoreSelection = (ignoreSelection != null) ? ignoreSelection : true;
	
	var editLink = null;
	var redirect = null;
	
	if (file == null || file.getMode() == App.MODE_DEVICE || file.getMode() == App.MODE_BROWSER)
	{
		editLink = '_blank';
	}
	else
	{
		editLink = url;
		redirect = editLink;
	}
		
	if (node == null)
	{
		return '';
	}
	else
	{
		var fileNode = node;

		// Ignores case for possible HTML or XML nodes
		if (fileNode.nodeName.toLowerCase() != 'mxfile')
		{
			// Removes control chars in input for correct roundtrip check
			var text = graph.zapGremlins(mxUtils.getXml(node));
			var data = graph.compress(text);
			
			// Fallback to plain XML for invalid compression
			// TODO: Remove this fallback with active pages
			if (graph.decompress(data) != text)
			{
				return text;
			}
			else
			{
				var diagramNode = node.ownerDocument.createElement('diagram');
				diagramNode.setAttribute('id', 'tapd-' + Editor.guid());
				mxUtils.setTextContent(diagramNode, data);
				
				fileNode = node.ownerDocument.createElement('mxfile');
				fileNode.appendChild(diagramNode);
			}
		}
		
		if (!compact)
		{
			// Removes old metadata
			fileNode.removeAttribute('userAgent');
			fileNode.removeAttribute('version');
			fileNode.removeAttribute('editor');
			fileNode.removeAttribute('type');
			
			// Adds new metadata
			fileNode.setAttribute('modified', new Date().toISOString());
			fileNode.setAttribute('host', window.location.hostname);
			fileNode.setAttribute('agent', navigator.userAgent);
			fileNode.setAttribute('version', EditorUi.VERSION);
			fileNode.setAttribute('etag', 'tapd-' + Editor.guid());
			
			var md = (file != null) ? file.getMode() : this.mode;
			
			if (md != null)
			{
				fileNode.setAttribute('type', md);
			}
		}
		else
		{
			fileNode = fileNode.cloneNode(true);
			fileNode.removeAttribute('userAgent');
			fileNode.removeAttribute('version');
			fileNode.removeAttribute('editor');
			fileNode.removeAttribute('type');
		}

		var xml = mxUtils.getXml(fileNode);
		
		return xml;
	}
};


/**
 * Returns the SVG of the diagram with embedded XML. If a callback function is
 * used, the images are converted to data URIs.
 */
EditorUi.prototype.getEmbeddedSvg = function(xml, graph, url, noHeader, callback, ignoreSelection, redirect)
{
	var bg = graph.background;
	
	if (bg == mxConstants.NONE)
	{
		bg = null;
	}

	// Sets or disables alternate text for foreignObjects. Disabling is needed
	// because PhantomJS seems to ignore switch statements and paint all text.
	var svgRoot = graph.getSvg(bg, null, null, null, null, ignoreSelection);
	
	if (graph.shadowVisible)
	{
		graph.addSvgShadow(svgRoot);
	}

	if (xml != null)
	{
		svgRoot.setAttribute('content', xml);
	}
	
	if (url != null)
	{
		svgRoot.setAttribute('resource', url);
	}

	// LATER: Click on SVG content to start editing
//		if (redirect != null)
//		{
//			// TODO: Ignore anchor tag source for click event
//			svgRoot.setAttribute('style', 'cursor:pointer;');
//			svgRoot.setAttribute('onclick', 'window.location.href=\'' + redirect + '\';'); 
//		}

	if (callback != null)
	{
		this.convertImages(svgRoot, mxUtils.bind(this, function(svgRoot)
		{
			callback(((!noHeader) ? '<?xml version="1.0" encoding="UTF-8"?>\n' +
				'<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">\n' : '') +
				mxUtils.getXml(svgRoot));
		}));
	}
	else
	{
		return ((!noHeader) ? '<?xml version="1.0" encoding="UTF-8"?>\n' +
			'<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">\n' : '') +
			mxUtils.getXml(svgRoot);
	}
};

/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.doDownloadFile = function(data, filename, mimeType, base64Encoded, format)
{
	// Newer versions of IE
	if (window.Blob && navigator.msSaveOrOpenBlob)
	{
		var blob = (base64Encoded) ?
				this.base64ToBlob(data, mimeType) :
				new Blob([data], {type: mimeType})
		navigator.msSaveOrOpenBlob(blob, filename);
	}
	// Older versions of IE (binary not supported)
	else if (mxClient.IS_IE)
	{
		var win = window.open('about:blank', '_blank');
		
		if (win == null)
		{
			mxUtils.popup(data, true);
		}
		else
		{
			win.document.write(data);
			win.document.close();
			win.document.execCommand('SaveAs', true, filename);
			win.close();
		}
	}
//		else if (mxClient.IS_IOS)
//		{
//			this.showTextDialog(filename + ':', data);
//		}
	else
	{
		var a = document.createElement('a');
		
		// Workaround for mxXmlRequest.simulate no longer working in Safari/PaleMoon
		// if this is used (ie PNG export broken after XML export in Safari/PaleMoon).
		var useDownload = !mxClient.IS_SF && navigator.userAgent.indexOf("PaleMoon/") < 0 &&
			typeof a.download !== 'undefined';
		
		// Workaround for Chromium 65 cross-domain anchor download issue
		if (mxClient.IS_GC)
		{
			var raw = navigator.userAgent.match(/Chrom(e|ium)\/([0-9]+)\./)
			var vers = raw ? parseInt(raw[2], 10) : false;
			useDownload = vers == 65 ? false : useDownload;
		}
		
		if (useDownload || this.isOffline())
		{
			a.href = URL.createObjectURL((base64Encoded) ?
				this.base64ToBlob(data, mimeType) :
				new Blob([data], {type: mimeType}));
			
			if (useDownload)
			{
				a.download = filename;
			}
			else
			{
				// Workaround for same window in Safari
				a.setAttribute('target', '_blank');
			}

			document.body.appendChild(a);
			
			try
			{
				window.setTimeout(function()
				{
					URL.revokeObjectURL(a.href);
				}, 0);

				a.click();
				a.parentNode.removeChild(a);
			}
			catch (e)
			{
				// ignore
			}
		}
		else
		{
			var req = this.createEchoRequest(data, filename, mimeType, base64Encoded, format);
			
			req.simulate(document, '_blank');
		}
	}
};


/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.base64ToBlob = function(base64Data, contentType)
{
    contentType = contentType || '';
    var sliceSize = 1024;
    var byteCharacters = atob(base64Data);
    var bytesLength = byteCharacters.length;
    var slicesCount = Math.ceil(bytesLength / sliceSize);
    var byteArrays = new Array(slicesCount);

    for (var sliceIndex = 0; sliceIndex < slicesCount; ++sliceIndex)
    {
        var begin = sliceIndex * sliceSize;
        var end = Math.min(begin + sliceSize, bytesLength);

        var bytes = new Array(end - begin);
        
        for (var offset = begin, i = 0 ; offset < end; ++i, ++offset)
        {
            bytes[i] = byteCharacters[offset].charCodeAt(0);
        }
        
        byteArrays[sliceIndex] = new Uint8Array(bytes);
    }

    return new Blob(byteArrays, {type: contentType});
};

/**
 *
 */
EditorUi.prototype.exportToCanvas = function(callback, width, imageCache, background, error, limitHeight,
	ignoreSelection, scale, transparentBackground, addShadow, converter, graph, border, noCrop)
{
	limitHeight = (limitHeight != null) ? limitHeight : true;
	ignoreSelection = (ignoreSelection != null) ? ignoreSelection : true;
	graph = (graph != null) ? graph : this.editor.graph;
	border = (border != null) ? border : 0;
	
	var bg = (transparentBackground) ? transparentBackground : graph.background;
	
	if (bg == mxConstants.NONE)
	{
		bg = null;
	}
	
	if (bg == null)
	{
		bg = background;
	}
	
	// Handles special case where background is null but transparent is false
	if (bg == null && transparentBackground == false)
	{
		bg = '#ffffff';
	}
	
	this.convertImages(graph.getSvg(bg, null, null, noCrop, null, ignoreSelection, null, null, null, addShadow),
		mxUtils.bind(this, function(svgRoot)
	{
		var img = new Image();


		img.onload = mxUtils.bind(this, function()
		{
	   		try
	   		{
	   			var canvas = document.createElement('canvas');
				var w = parseInt(svgRoot.getAttribute('width'));
				var h = parseInt(svgRoot.getAttribute('height'));
				scale = (scale != null) ? scale : 1;
				
				if (width != null)
				{
					scale = (!limitHeight) ? width / w : Math.min(1, Math.min((width * 3) / (h * 4), width / w));
				}
				
				w = Math.ceil(scale * w) + 2 * border;
				h = Math.ceil(scale * h) + 2 * border;
				
				canvas.setAttribute('width', w);
		   		canvas.setAttribute('height', h);
		   		var ctx = canvas.getContext('2d');
		   		
		   		if (bg != null)
		   		{
		   			ctx.beginPath();
					ctx.rect(0, 0, w, h);
					ctx.fillStyle = bg;
					ctx.fill();
		   		}

		   		ctx.scale(scale, scale);
		   		
		   		// Workaround for broken data URI images in Safari on first export
		   		if (mxClient.IS_SF)
		   		{			   		
					window.setTimeout(function()
					{
						ctx.drawImage(img, border / scale, border / scale);
						callback(canvas);
					}, 0);
		   		}
		   		else
		   		{
		   			ctx.drawImage(img, border / scale, border / scale);
		   			callback(canvas);
		   		}
	   		}
	   		catch (e)
	   		{
				this.handleError(e);
	   		}
		});
		
		img.onerror = function(e)
		{
			this.handleError(e);
		};

		try
		{
			if (addShadow)
			{
				this.editor.graph.addSvgShadow(svgRoot);
			}
			
			var done = mxUtils.bind(this, function()
			{
				if (this.editor.resolvedFontCss != null)
				{
					var st = document.createElement('style');
					st.setAttribute('type', 'text/css');
					st.innerHTML = this.editor.resolvedFontCss;
					
					// Must be in defs section for FF to work
					var defs = svgRoot.getElementsByTagName('defs');
					defs[0].appendChild(st);
				}
				img.src = this.createSvgDataUri(mxUtils.getXml(svgRoot));
			});

			this.loadFonts(done);
		}
		catch (e)
		{
			//console.log('src', e, img.src);
			this.handleError(e);
		}
	}), imageCache, converter);
};

/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.createSvgDataUri = function(svg)
{
	return 'data:image/svg+xml;base64,' + btoa(unescape(encodeURIComponent(svg)));
};	

/**
 * 
 */
EditorUi.prototype.saveCanvas = function(canvas, xml, format, exportHandler)
{
	var ext = ((format == 'jpeg') ? 'jpg' : format);
	var filename = this.getBaseFilename() + '.' + ext;
	var data = this.createImageDataUri(canvas, xml, format);

	if (exportHandler && typeof exportHandler === 'function') {
		exportHandler(data, filename);
	} else {
   		this.doDownloadFile(data.substring(data.lastIndexOf(',') + 1), filename, 'image/' + format, true, format);
	}

};


EditorUi.prototype.getBaseFilename = function()
{
	return this.editor.getFilename();
};


EditorUi.prototype.createImageDataUri = function(canvas, xml, format)
{
	    var data = canvas.toDataURL('image/' + format);
	    
	    // Checks if output is invalid or empty
	    if (data.length <= 6 || data == canvas.cloneNode(false).toDataURL('image/' + format))
	    {
	    	throw {message: 'Invalid image'};
	    }
	    
	    if (xml != null)
	    {
	   		data = this.writeGraphModelToPng(data, 'zTXt', 'mxGraphModel', atob(this.editor.graph.compress(xml)));
	    }
	    
	    return data;
};

/**
 * Converts all images in the SVG output to data URIs for immediate rendering
 */
EditorUi.prototype.convertImages = function(svgRoot, callback, imageCache, converter)
{
	// Converts images to data URLs for immediate painting
	if (converter == null)
	{
		converter = this.createImageUrlConverter();
	}
	
	// Barrier for asynchronous image loading
	var counter = 0;
	
	function inc()
	{
		counter++;
	};
	
	function dec()
	{
		counter--;
		
		if (counter == 0)
		{
			callback(svgRoot);
		}
	};

	var cache = imageCache || new Object();
	
	var convertImages = mxUtils.bind(this, function(tagName, srcAttr)
	{
		var images = svgRoot.getElementsByTagName(tagName);
		
		for (var i = 0; i < images.length; i++)
		{
			(mxUtils.bind(this, function(img)
			{
				var src = converter.convert(img.getAttribute(srcAttr));
	        	
				// Data URIs are pass-through
				if (src != null && src.substring(0, 5) != 'data:')
				{
					var tmp = cache[src];
					
					if (tmp == null)
					{
						inc();
						
						this.convertImageToDataUri(src, function(uri)
						{
							if (uri != null)
							{
								cache[src] = uri;
								img.setAttribute(srcAttr, uri);
							}
							
							dec();
						});
					}
					else
					{
						img.setAttribute(srcAttr, tmp);
					}
				}
				else if (src != null)
				{
					img.setAttribute(srcAttr, src);
				}
			}))(images[i]);
		}
	});
	
	// Converts all known image tags in output
	// LATER: Add support for images in CSS
	convertImages('image', 'xlink:href');
	convertImages('img', 'src');
	
	// All from cache or no images
	if (counter == 0)
	{
		callback(svgRoot);
	}
};

	/**
	 * Translates this point by the given vector.
	 * 
	 * @param {number} dx X-coordinate of the translation.
	 * @param {number} dy Y-coordinate of the translation.
	 */
	EditorUi.prototype.convertImageToDataUri = function(url, callback)
	{
		if (/(\.svg)$/i.test(url))
		{
			mxUtils.get(url, mxUtils.bind(this, function(req)
			{
				callback(this.createSvgDataUri(req.getText()));
			}),
			function()
			{
				callback(this.svgBrokenImage.src);
			});
		}
		else
		{
		    var img = new Image();
		    var self = this;
		    
		    if (this.crossOriginImages)
	    	{
			    img.crossOrigin = 'anonymous';
		    }
		    
		    img.onload = function()
		    {
		        var canvas = document.createElement('canvas');
		        var ctx = canvas.getContext('2d');
		        canvas.height = img.height;
		        canvas.width = img.width;
		        ctx.drawImage(img, 0, 0);
		        
		        try
		        {
	        		callback(canvas.toDataURL());
		        }
		        catch (e)
		        {
	        		callback(self.svgBrokenImage.src);
		        }
		    };
		    
		    img.onerror = function()
		    {
	    		callback(self.svgBrokenImage.src);
		    };

		    img.src = url;
		}
	};


/**
 * Converts all images in the SVG output to data URIs for immediate rendering
 */
EditorUi.prototype.createImageUrlConverter = function()
{
	var converter = new mxUrlConverter();
	converter.updateBaseUrl();

	// Extends convert to avoid CORS using an image proxy server where needed
	var convert = converter.convert;
	var self = this;
	
	converter.convert = function(src)
	{
		// if (src != null)
		// {
		// 	var remote = src.substring(0, 7) == 'http://' || src.substring(0, 8) == 'https://';
			
		// 	if (remote && !navigator.onLine)
		// 	{
		// 		src = self.svgBrokenImage.src;
		// 	}
		// 	else if (remote && src.substring(0, converter.baseUrl.length) != converter.baseUrl &&
		// 			(!self.crossOriginImages || !self.isCorsEnabledForUrl(src)))
		// 	{
		// 		src = PROXY_URL + '?url=' + encodeURIComponent(src);
		// 	}
		// 	else if (src.substring(0, 19) != 'chrome-extension://')
		// 	{
		// 		src = convert.apply(this, arguments);
		// 	}
		// }
		
		return src;
	};
	
	return converter;
};

/**
 * Adds the given text to the compressed or non-compressed text chunk.
 */
EditorUi.prototype.writeGraphModelToPng = function(data, type, key, value, error)
{
	var base64 = data.substring(data.indexOf(',') + 1);
	var f = (window.atob) ? atob(base64) : Base64.decode(base64, true);
	var pos = 0;
	
	function fread(d, count)
	{
		var start = pos;
		pos += count;
		
		return d.substring(start, pos);
	};
	
	// Reads unsigned long 32 bit big endian
	function _freadint(d)
	{
		var bytes = fread(d, 4);
		
		return bytes.charCodeAt(3) + (bytes.charCodeAt(2) << 8) +
			(bytes.charCodeAt(1) << 16) + (bytes.charCodeAt(0) << 24);
	};
	
	function writeInt(num)
	{
		return String.fromCharCode((num >> 24) & 0x000000ff, (num >> 16) & 0x000000ff,
			(num >> 8) & 0x000000ff, num & 0x000000ff);
	};
	
	// Checks signature
	if (fread(f,8) != String.fromCharCode(137) + 'PNG' + String.fromCharCode(13, 10, 26, 10))
	{
		if (error != null)
		{
			error();
		}
		
		return;
	}
	
	// Reads header chunk
	fread(f,4);
	
	if (fread(f,4) != 'IHDR')
	{
		if (error != null)
		{
			error();
		}
		
		return;
	}
	
	fread(f, 17);
	var result = f.substring(0, pos);
	
	do
	{
		var n = _freadint(f);
		var chunk = fread(f,4);
		
		if (chunk == 'IDAT')
		{
			result = f.substring(0, pos - 8);
			
			var chunkData = key + String.fromCharCode(0) +
				((type == 'zTXt') ? String.fromCharCode(0) : '') + 
				value;
			
			// FIXME: Wrong crc
			var crc = 0xffffffff;
			crc = this.updateCRC(crc, type, 0, 4);
			crc = this.updateCRC(crc, chunkData, 0, chunkData.length);
			
			result += writeInt(chunkData.length) + type + chunkData + writeInt(crc ^ 0xffffffff);
			result += f.substring(pos - 8, f.length);
			
			break;
		}
		
		result += f.substring(pos - 8, pos - 4 + n);
		fread(f,n);
		fread(f,4);
	}
	while (n);
	
	return 'data:image/png;base64,' + ((window.btoa) ? btoa(result) : Base64.encode(result, true));
}

EditorUi.prototype.crcTable = [];

for (var n = 0; n < 256; n++)
{
	var c = n;
	
	for (var k = 0; k < 8; k++)
	{
		if ((c & 1) == 1)
		{
			c = 0xedb88320 ^ (c >>> 1);
		}
		else
		{
			c >>>= 1;
		}

		EditorUi.prototype.crcTable[n] = c;
	}
}

EditorUi.prototype.updateCRC = function(crc, data, off, len)
{
	var c = crc;

	for (var n = 0; n < len; n++)
	{
		c = EditorUi.prototype.crcTable[(c ^ data[off + n]) & 0xff] ^ (c >>> 8);
	}

	return c;
};

/**
 * For the fontCSS to be applied when rendering images on canvas, the actual
 * font data must be made available via a data URI encoding of the file.
 */
EditorUi.prototype.loadFonts = function(then)
{
    if (this.editor.fontCss != null && this.editor.resolvedFontCss == null)
    {
        var parts = this.editor.fontCss.split('url(');
        var waiting = 0;
        var fonts = {};

        // Strips leading and trailing quotes and spaces
        function trimString(str)
        {
            return str.replace(new RegExp("^[\\s\"']+", "g"), "").replace(new RegExp("[\\s\"']+$", "g"), "");
        };
        
        var finish = mxUtils.bind(this, function()
        {
            if (waiting == 0)
            {
                // Constructs string
                var result = [parts[0]];
                
                for (var j = 1; j < parts.length; j++)
                {
                    var idx = parts[j].indexOf(')');
                    result.push('url("');
                    result.push(fonts[trimString(parts[j].substring(0, idx))]);
                    result.push('"' + parts[j].substring(idx));
                }
                
                this.editor.resolvedFontCss = result.join('');
                then();
            }
        });
        
        if (parts.length > 0)
        {
            for (var i = 1; i < parts.length; i++)
            {
                var idx = parts[i].indexOf(')');
                var format = null;
                
                // Checks if there is a format directive
                var fmtIdx = parts[i].indexOf('format(', idx);
                
                if (fmtIdx > 0)
                {
                    format = trimString(parts[i].substring(fmtIdx + 7, parts[i].indexOf(')', fmtIdx)));
                }

                (mxUtils.bind(this, function(url)
                {
                    if (fonts[url] == null)
                    {
                        // Mark font es being fetched and fetch it
                        fonts[url] = url;
                        waiting++;
                        
                        var mime = 'application/x-font-ttf';
                        
                        // See https://stackoverflow.com/questions/2871655/proper-mime-type-for-fonts
                        if (format == 'svg' || /(\.svg)($|\?)/i.test(url))
                        {
                            mime = 'image/svg+xml';
                        }
                        else if (format == 'otf' || format == 'embedded-opentype' || /(\.otf)($|\?)/i.test(url))
                        {
                            mime = 'application/x-font-opentype';
                        }
                        else if (format == 'woff' || /(\.woff)($|\?)/i.test(url))
                        {
                            mime = 'application/font-woff';
                        }
                        else if (format == 'woff2' || /(\.woff2)($|\?)/i.test(url))
                        {
                            mime = 'application/font-woff2';
                        }
                        else if (format == 'eot' || /(\.eot)($|\?)/i.test(url))
                        {
                            mime = 'application/vnd.ms-fontobject';
                        }
                        else if (format == 'sfnt' || /(\.sfnt)($|\?)/i.test(url))
                        {
                            mime = 'application/font-sfnt';
                        }
                        
                        var realUrl = url;
                        
                        if ((/^https?:\/\//.test(realUrl)) && !this.isCorsEnabledForUrl(realUrl))
                        {
                            realUrl = PROXY_URL + '?url=' + encodeURIComponent(url);
                        }

                        // LATER: Remove cache-control header
                        this.loadUrl(realUrl, mxUtils.bind(this, function(uri)
                        {
                            fonts[url] = uri;
                            waiting--;
                            finish();
                        }), mxUtils.bind(this, function(err)
                        {
                            // LATER: handle error
                            waiting--;
                            finish();
                        }), true, null, 'data:' + mime + ';charset=utf-8;base64,');
                    }
                }))(trimString(parts[i].substring(0, idx)), format);
            }
        }
    }
    else
    {
        then();
    }
};
/**
 * Installs the listeners to update the action states.
 */
EditorUi.prototype.init = function()
{
	/**
	 * Keypress starts immediate editing on selection cell
	 */
	var graph = this.editor.graph;
		
	mxEvent.addListener(graph.container, 'keydown', mxUtils.bind(this, function(evt)
	{
		this.onKeyDown(evt);
	}));
	mxEvent.addListener(graph.container, 'keypress', mxUtils.bind(this, function(evt)
	{
		this.onKeyPress(evt);
	}));

	// Updates action states
	this.addUndoListener();
	this.addBeforeUnloadListener();
	
	graph.getSelectionModel().addListener(mxEvent.CHANGE, mxUtils.bind(this, function()
	{
		this.updateActionStates();
	}));
	
	graph.getModel().addListener(mxEvent.CHANGE, mxUtils.bind(this, function()
	{
		this.updateActionStates();
	}));
	
	// Changes action states after change of default parent
	var graphSetDefaultParent = graph.setDefaultParent;
	var ui = this;
	
	this.editor.graph.setDefaultParent = function()
	{
		graphSetDefaultParent.apply(this, arguments);
		ui.updateActionStates();
	};
	
	// Hack to make editLink available in vertex handler
	graph.editLink = ui.actions.get('editLink').funct;
	
	this.updateActionStates();
	this.initClipboard();
	this.initCanvas();
	
	if (this.format != null)
	{
		this.format.init();
	}
};

/**
 * Returns true if the given event should start editing. This implementation returns true.
 */
EditorUi.prototype.onKeyDown = function(evt)
{
	var graph = this.editor.graph;
	
	// Tab selects next cell
	if (evt.which == 9 && graph.isEnabled() && !mxEvent.isAltDown(evt))
	{
		if (graph.isEditing())
		{
			graph.stopEditing(false);
		}
		else
		{
			graph.selectCell(!mxEvent.isShiftDown(evt));
		}
		
		mxEvent.consume(evt);
	}
};

/**
 * Returns true if the given event should start editing. This implementation returns true.
 */
EditorUi.prototype.onKeyPress = function(evt)
{
	var graph = this.editor.graph;
	
	// KNOWN: Focus does not work if label is empty in quirks mode
	if (this.isImmediateEditingEvent(evt) && !graph.isEditing() && !graph.isSelectionEmpty() && evt.which !== 0 &&
		!mxEvent.isAltDown(evt) && !mxEvent.isControlDown(evt) && !mxEvent.isMetaDown(evt))
	{
		graph.escape();
		graph.startEditing();

		// Workaround for FF where char is lost if cursor is placed before char
		if (mxClient.IS_FF)
		{
			var ce = graph.cellEditor;
			ce.textarea.innerHTML = String.fromCharCode(evt.which);

			// Moves cursor to end of textarea
			var range = document.createRange();
			range.selectNodeContents(ce.textarea);
			range.collapse(false);
			var sel = window.getSelection();
			sel.removeAllRanges();
			sel.addRange(range);
		}
	}
};

/**
 * Returns true if the given event should start editing. This implementation returns true.
 */
EditorUi.prototype.isImmediateEditingEvent = function(evt)
{
	return true;
};

/**
 * Private helper method.
 */
EditorUi.prototype.getCssClassForMarker = function(prefix, shape, marker, fill)
{
	var result = '';

	if (shape == 'flexArrow')
	{
		result = (marker != null && marker != mxConstants.NONE) ?
			'font-graph font-graph-geSprite-' + prefix + 'blocktrans' : 'font-graph font-graph-geSprite-noarrow';
	}
	else
	{
		if (marker == mxConstants.ARROW_CLASSIC)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'classic' : 'font-graph font-graph-geSprite-' + prefix + 'classictrans';
		}
		else if (marker == mxConstants.ARROW_CLASSIC_THIN)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'classicthin' : 'font-graph font-graph-geSprite-' + prefix + 'classicthintrans';
		}
		else if (marker == mxConstants.ARROW_OPEN)
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'open';
		}
		else if (marker == mxConstants.ARROW_OPEN_THIN)
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'openthin';
		}
		else if (marker == mxConstants.ARROW_BLOCK)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'block' : 'font-graph font-graph-geSprite-' + prefix + 'blocktrans';
		}
		else if (marker == mxConstants.ARROW_BLOCK_THIN)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'blockthin' : 'font-graph font-graph-geSprite-' + prefix + 'blockthintrans';
		}
		else if (marker == mxConstants.ARROW_OVAL)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'oval' : 'font-graph font-graph-geSprite-' + prefix + 'ovaltrans';
		}
		else if (marker == mxConstants.ARROW_DIAMOND)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'diamond' : 'font-graph font-graph-geSprite-' + prefix + 'diamondtrans';
		}
		else if (marker == mxConstants.ARROW_DIAMOND_THIN)
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'thindiamond' : 'font-graph font-graph-geSprite-' + prefix + 'thindiamondtrans';
		}
		else if (marker == 'openAsync')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'openasync';
		}
		else if (marker == 'dash')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'dash';
		}
		else if (marker == 'cross')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'cross';
		}
		else if (marker == 'async')
		{
			result = (fill == '1') ? 'font-graph font-graph-geSprite-' + prefix + 'async' : 'font-graph font-graph-geSprite-' + prefix + 'asynctrans';
		}
		else if (marker == 'circle' || marker == 'circlePlus')
		{
			result = (fill == '1' || marker == 'circle') ? 'font-graph font-graph-geSprite-' + prefix + 'circle' : 'font-graph font-graph-geSprite-' + prefix + 'circleplus';
		}
		else if (marker == 'ERone')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'erone';
		}
		else if (marker == 'ERmandOne')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'eronetoone';
		}
		else if (marker == 'ERmany')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'ermany';
		}
		else if (marker == 'ERoneToMany')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'eronetomany';
		}
		else if (marker == 'ERzeroToOne')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'eroneopt';
		}
		else if (marker == 'ERzeroToMany')
		{
			result = 'font-graph font-graph-geSprite-' + prefix + 'ermanyopt';
		}
		else
		{
			result = 'font-graph font-graph-geSprite-noarrow';
		}
	}

	return result;
};

/**
 * Overridden in Menus.js
 */
EditorUi.prototype.createMenus = function()
{
	return null;
};

/**
 * Hook for allowing selection and context menu for certain events.
 */
EditorUi.prototype.updatePasteActionStates = function()
{
	var graph = this.editor.graph;
	var paste = this.actions.get('paste');
	var pasteHere = this.actions.get('pasteHere');
	
	paste.setEnabled(this.editor.graph.cellEditor.isContentEditing() || (!mxClipboard.isEmpty() &&
		graph.isEnabled() && !graph.isCellLocked(graph.getDefaultParent())));
	pasteHere.setEnabled(paste.isEnabled());
};

/**
 * Hook for allowing selection and context menu for certain events.
 */
EditorUi.prototype.initClipboard = function()
{
	var ui = this;

	var mxClipboardCut = mxClipboard.cut;
	mxClipboard.cut = function(graph)
	{
		if (graph.cellEditor.isContentEditing())
		{
			document.execCommand('cut', false, null);
		}
		else
		{
			mxClipboardCut.apply(this, arguments);
		}
		
		ui.updatePasteActionStates();
	};
	
	var mxClipboardCopy = mxClipboard.copy;
	mxClipboard.copy = function(graph)
	{
		if (graph.cellEditor.isContentEditing())
		{
			document.execCommand('copy', false, null);
		}
		else
		{
			mxClipboardCopy.apply(this, arguments);
		}
		
		ui.updatePasteActionStates();
	};
	
	var mxClipboardPaste = mxClipboard.paste;
	mxClipboard.paste = function(graph)
	{
		var result = null;
		
		if (graph.cellEditor.isContentEditing())
		{
			document.execCommand('paste', false, null);
		}
		else
		{
			result = mxClipboardPaste.apply(this, arguments);
		}
		
		ui.updatePasteActionStates();
		
		return result;
	};

	// Overrides cell editor to update paste action state
	var cellEditorStartEditing = this.editor.graph.cellEditor.startEditing;
	
	this.editor.graph.cellEditor.startEditing = function()
	{
		cellEditorStartEditing.apply(this, arguments);
		ui.updatePasteActionStates();
	};
	
	var cellEditorStopEditing = this.editor.graph.cellEditor.stopEditing;
	
	this.editor.graph.cellEditor.stopEditing = function(cell, trigger)
	{
		cellEditorStopEditing.apply(this, arguments);
		ui.updatePasteActionStates();
	};
	
	this.updatePasteActionStates();
};

/**
 * Initializes the infinite canvas.
 */
EditorUi.prototype.lazyZoomDelay = 20;

/**
 * Initializes the infinite canvas.
 */
EditorUi.prototype.initCanvas = function()
{
	// Initial page layout view, scrollBuffer and timer-based scrolling
	var graph = this.editor.graph;
	graph.timerAutoScroll = true;

	/**
	 * Returns the padding for pages in page view with scrollbars.
	 */
	graph.getPagePadding = function()
	{
		return new mxPoint(Math.max(0, Math.round((graph.container.offsetWidth - 34) / graph.view.scale)),
				Math.max(0, Math.round((graph.container.offsetHeight - 34) / graph.view.scale)));
	};

	// Fits the number of background pages to the graph
	graph.view.getBackgroundPageBounds = function()
	{
		var layout = this.graph.getPageLayout();
		var page = this.graph.getPageSize();
		
		return new mxRectangle(this.scale * (this.translate.x + layout.x * page.width),
				this.scale * (this.translate.y + layout.y * page.height),
				this.scale * layout.width * page.width,
				this.scale * layout.height * page.height);
	};

	graph.getPreferredPageSize = function(bounds, width, height)
	{
		var pages = this.getPageLayout();
		var size = this.getPageSize();
		
		return new mxRectangle(0, 0, pages.width * size.width, pages.height * size.height);
	};
	
	// Scales pages/graph to fit available size
	var resize = null;
	var ui = this;
	
	if (this.editor.isChromelessView())
	{
        resize = mxUtils.bind(this, function(autoscale, maxScale, cx, cy)
        {
            if (graph.container != null)
            {
                cx = (cx != null) ? cx : 0;
                cy = (cy != null) ? cy : 0;
                
                var bds = (graph.pageVisible) ? graph.view.getBackgroundPageBounds() : graph.getGraphBounds();
                var scroll = mxUtils.hasScrollbars(graph.container);
                var tr = graph.view.translate;
                var s = graph.view.scale;
                
                // Normalizes the bounds
                var b = mxRectangle.fromRectangle(bds);
                b.x = b.x / s - tr.x;
                b.y = b.y / s - tr.y;
                b.width /= s;
                b.height /= s;
                
                var st = graph.container.scrollTop;
                var sl = graph.container.scrollLeft;
                var sb = (mxClient.IS_QUIRKS || document.documentMode >= 8) ? 20 : 14;
                
                if (document.documentMode == 8 || document.documentMode == 9)
                {
                    sb += 3;
                }
                
                var cw = graph.container.offsetWidth - sb;
                var ch = graph.container.offsetHeight - sb;
                
                var ns = (autoscale) ? Math.max(0.3, Math.min(maxScale || 1, cw / b.width)) : s;
                var dx = ((cw - ns * b.width) / 2) / ns;
                var dy = (this.lightboxVerticalDivider == 0) ? 0 : ((ch - ns * b.height) / this.lightboxVerticalDivider) / ns;
                
                if (scroll)
                {
                    dx = Math.max(dx, 0);
                    dy = Math.max(dy, 0);
                }

                if (scroll || bds.width < cw || bds.height < ch)
                {
                    graph.view.scaleAndTranslate(ns, Math.floor(dx - b.x), Math.floor(dy - b.y));
                    graph.container.scrollTop = st * ns / s;
                    graph.container.scrollLeft = sl * ns / s;
                }
                else if (cx != 0 || cy != 0)
                {
                    var t = graph.view.translate;
                    graph.view.setTranslate(Math.floor(t.x + cx / s), Math.floor(t.y + cy / s));
                }
            }
        });
		
		// Hack to make function available to subclassers
		this.chromelessResize = resize;

		// Hook for subclassers for override
		this.chromelessWindowResize = mxUtils.bind(this, function()
	   	{
			this.chromelessResize(false);
	   	});

		// Removable resize listener
		var autoscaleResize = mxUtils.bind(this, function()
	   	{
			this.chromelessWindowResize(false);
	   	});
		
	   	mxEvent.addListener(window, 'resize', autoscaleResize);
	   	
	   	this.destroyFunctions.push(function()
	   	{
	   		mxEvent.removeListener(window, 'resize', autoscaleResize);
	   	});
	   	
		this.editor.addListener('resetGraphView', mxUtils.bind(this, function()
		{
			this.chromelessResize(true);
		}));

		this.actions.get('zoomIn').funct = mxUtils.bind(this, function(evt)
		{
			graph.zoomIn();
			this.chromelessResize(false);
		});
		this.actions.get('zoomOut').funct = mxUtils.bind(this, function(evt)
		{
			graph.zoomOut();
			this.chromelessResize(false);
		});
		
		// Creates toolbar for viewer - do not use CSS here
		// as this may be used in a viewer that has no CSS
		if (urlParams['toolbar'] != '0')
		{
			this.chromelessToolbar = document.createElement('div');
			this.chromelessToolbar.style.position = 'fixed';
			this.chromelessToolbar.style.overflow = 'hidden';
			this.chromelessToolbar.style.boxSizing = 'border-box';
			this.chromelessToolbar.style.whiteSpace = 'nowrap';
			this.chromelessToolbar.style.backgroundColor = '#000000';
			this.chromelessToolbar.style.padding = '10px 10px 8px 10px';
			this.chromelessToolbar.style.left = '50%';
			
			if (!mxClient.IS_VML)
			{
				mxUtils.setPrefixedStyle(this.chromelessToolbar.style, 'borderRadius', '20px');
				mxUtils.setPrefixedStyle(this.chromelessToolbar.style, 'transition', 'opacity 600ms ease-in-out');
			}
			
			var updateChromelessToolbarPosition = mxUtils.bind(this, function()
			{
				var css = mxUtils.getCurrentStyle(graph.container);
			 	this.chromelessToolbar.style.bottom = ((css != null) ? parseInt(css['margin-bottom'] || 0) : 0) +
			 		((this.tabContainer != null) ? (20 + parseInt(this.tabContainer.style.height)) : 20) + 'px';
			});
			
			this.editor.addListener('resetGraphView', updateChromelessToolbarPosition);
			updateChromelessToolbarPosition();
			
			var btnCount = 0;
	
			var addButton = mxUtils.bind(this, function(fn, imgSrc, tip)
			{
				btnCount++;
				
				var a = document.createElement('span');
				a.style.paddingLeft = '8px';
				a.style.paddingRight = '8px';
				a.style.cursor = 'pointer';
				mxEvent.addListener(a, 'click', fn);
				
				if (tip != null)
				{
					a.setAttribute('title', tip);
				}
				
				var img = document.createElement('img');
				img.setAttribute('border', '0');
				img.setAttribute('src', imgSrc);
				
				a.appendChild(img);
				this.chromelessToolbar.appendChild(a);
				
				return a;
			});
			
			var prevButton = addButton(mxUtils.bind(this, function(evt)
			{
				this.actions.get('previousPage').funct();
				mxEvent.consume(evt);
			}), Editor.previousLargeImage, mxResources.get('previousPage'));
			
			var pageInfo = document.createElement('div');
			pageInfo.style.display = 'inline-block';
			pageInfo.style.verticalAlign = 'top';
			pageInfo.style.fontFamily = 'Helvetica,Arial';
			pageInfo.style.marginTop = '8px';
			pageInfo.style.fontSize = '14px';
			pageInfo.style.color = '#ffffff';
			this.chromelessToolbar.appendChild(pageInfo);
			
			var nextButton = addButton(mxUtils.bind(this, function(evt)
			{
				this.actions.get('nextPage').funct();
				mxEvent.consume(evt);
			}), Editor.nextLargeImage, mxResources.get('nextPage'));
			
			var updatePageInfo = mxUtils.bind(this, function()
			{
				if (this.pages != null && this.pages.length > 1 && this.currentPage != null)
				{
					pageInfo.innerHTML = '';
					mxUtils.write(pageInfo, (mxUtils.indexOf(this.pages, this.currentPage) + 1) + ' / ' + this.pages.length);
				}
			});
			
			prevButton.style.paddingLeft = '0px';
			prevButton.style.paddingRight = '4px';
			nextButton.style.paddingLeft = '4px';
			nextButton.style.paddingRight = '0px';
			
			var updatePageButtons = mxUtils.bind(this, function()
			{
				if (this.pages != null && this.pages.length > 1 && this.currentPage != null)
				{
					nextButton.style.display = '';
					prevButton.style.display = '';
					pageInfo.style.display = 'inline-block';
				}
				else
				{
					nextButton.style.display = 'none';
					prevButton.style.display = 'none';
					pageInfo.style.display = 'none';
				}
				
				updatePageInfo();
			});
			
			this.editor.addListener('resetGraphView', updatePageButtons);
			this.editor.addListener('pageSelected', updatePageInfo);
	
			addButton(mxUtils.bind(this, function(evt)
			{
				this.actions.get('zoomOut').funct();
				mxEvent.consume(evt);
			}), Editor.zoomOutLargeImage, mxResources.get('zoomOut') + ' (Alt+Mousewheel)');
			
			addButton(mxUtils.bind(this, function(evt)
			{
				this.actions.get('zoomIn').funct();
				mxEvent.consume(evt);
			}), Editor.zoomInLargeImage, mxResources.get('zoomIn') + ' (Alt+Mousewheel)');
			
			addButton(mxUtils.bind(this, function(evt)
			{
				if (graph.isLightboxView())
				{
					if (graph.view.scale == 1)
					{
						this.lightboxFit();
					}
					else
					{
						graph.zoomTo(1);
					}
					
					this.chromelessResize(false);
				}
				else
				{
					this.chromelessResize(true);
				}
				
				mxEvent.consume(evt);
			}), Editor.actualSizeLargeImage, mxResources.get('fit'));
	
			// Changes toolbar opacity on hover
			var fadeThread = null;
			var fadeThread2 = null;
			
			var fadeOut = mxUtils.bind(this, function(delay)
			{
				if (fadeThread != null)
				{
					window.clearTimeout(fadeThread);
					fadeThead = null;
				}
				
				if (fadeThread2 != null)
				{
					window.clearTimeout(fadeThread2);
					fadeThead2 = null;
				}
				
				fadeThread = window.setTimeout(mxUtils.bind(this, function()
				{
				 	mxUtils.setOpacity(this.chromelessToolbar, 0);
					fadeThread = null;
				 	
					fadeThread2 = window.setTimeout(mxUtils.bind(this, function()
					{
						this.chromelessToolbar.style.display = 'none';
						fadeThread2 = null;
					}), 600);
				}), delay || 200);
			});
			
			var fadeIn = mxUtils.bind(this, function(opacity)
			{
				if (fadeThread != null)
				{
					window.clearTimeout(fadeThread);
					fadeThead = null;
				}
				
				if (fadeThread2 != null)
				{
					window.clearTimeout(fadeThread2);
					fadeThead2 = null;
				}
				
				this.chromelessToolbar.style.display = '';
				mxUtils.setOpacity(this.chromelessToolbar, opacity || 30);
			});
	
			if (urlParams['layers'] == '1')
			{
				this.layersDialog = null;
				
				var layersButton = addButton(mxUtils.bind(this, function(evt)
				{
					if (this.layersDialog != null)
					{
						this.layersDialog.parentNode.removeChild(this.layersDialog);
						this.layersDialog = null;
					}
					else
					{
						this.layersDialog = graph.createLayersDialog();
						
						mxEvent.addListener(this.layersDialog, 'mouseleave', mxUtils.bind(this, function()
						{
							this.layersDialog.parentNode.removeChild(this.layersDialog);
							this.layersDialog = null;
						}));
						
						var r = layersButton.getBoundingClientRect();
						
						mxUtils.setPrefixedStyle(this.layersDialog.style, 'borderRadius', '5px');
						this.layersDialog.style.position = 'fixed';
						this.layersDialog.style.fontFamily = 'Helvetica,Arial';
						this.layersDialog.style.backgroundColor = '#000000';
						this.layersDialog.style.width = '160px';
						this.layersDialog.style.padding = '4px 2px 4px 2px';
						this.layersDialog.style.color = '#ffffff';
						mxUtils.setOpacity(this.layersDialog, 70);
						this.layersDialog.style.left = r.left + 'px';
						this.layersDialog.style.bottom = parseInt(this.chromelessToolbar.style.bottom) +
							this.chromelessToolbar.offsetHeight + 4 + 'px';
						
						// Puts the dialog on top of the container z-index
						var style = mxUtils.getCurrentStyle(this.editor.graph.container);
						this.layersDialog.style.zIndex = style.zIndex;
						
						document.body.appendChild(this.layersDialog);
					}
					
					mxEvent.consume(evt);
				}), Editor.layersLargeImage, mxResources.get('layers'));
				
				// Shows/hides layers button depending on content
				var model = graph.getModel();
	
				model.addListener(mxEvent.CHANGE, function()
				{
					 layersButton.style.display = (model.getChildCount(model.root) > 1) ? '' : 'none';
				});
			}
	
			this.addChromelessToolbarItems(addButton);
	
			if (this.editor.editButtonLink != null || this.editor.editButtonFunc != null)
			{
				addButton(mxUtils.bind(this, function(evt)
				{
					if (this.editor.editButtonFunc != null) 
					{
						this.editor.editButtonFunc();
					} 
					else if (this.editor.editButtonLink == '_blank')
					{
						this.editor.editAsNew(this.getEditBlankXml());
					}
					else
					{
						graph.openLink(this.editor.editButtonLink, 'editWindow');
					}
					
					mxEvent.consume(evt);
				}), Editor.editLargeImage, mxResources.get('edit'));
			}
			
			if (this.lightboxToolbarActions != null)
			{
				for (var i = 0; i < this.lightboxToolbarActions.length; i++)
				{
					var lbAction = this.lightboxToolbarActions[i];
					addButton(lbAction.fn, lbAction.icon, lbAction.tooltip);
				}
			}
			
			if (graph.lightbox && (urlParams['close'] == '1' || this.container != document.body))
			{
				addButton(mxUtils.bind(this, function(evt)
				{
					if (urlParams['close'] == '1')
					{
						window.close();
					}
					else
					{
						this.destroy();
						mxEvent.consume(evt);
					}
				}), Editor.closeLargeImage, mxResources.get('close') + ' (Escape)');
			}
	
			// Initial state invisible
			this.chromelessToolbar.style.display = 'none';
			mxUtils.setPrefixedStyle(this.chromelessToolbar.style, 'transform', 'translate(-50%,0)');
			graph.container.appendChild(this.chromelessToolbar);
			
			mxEvent.addListener(graph.container, (mxClient.IS_POINTER) ? 'pointermove' : 'mousemove', mxUtils.bind(this, function(evt)
			{
				if (!mxEvent.isTouchEvent(evt))
				{
					if (!mxEvent.isShiftDown(evt))
					{
						fadeIn(30);
					}
					
					fadeOut();
				}
			}));
			
			mxEvent.addListener(this.chromelessToolbar, (mxClient.IS_POINTER) ? 'pointermove' : 'mousemove', function(evt)
			{
				mxEvent.consume(evt);
			});
			
			mxEvent.addListener(this.chromelessToolbar, 'mouseenter', mxUtils.bind(this, function(evt)
			{
				if (!mxEvent.isShiftDown(evt))
				{
					fadeIn(100);
				}
				else
				{
					fadeOut();
				}
			}));

			mxEvent.addListener(this.chromelessToolbar, 'mousemove',  mxUtils.bind(this, function(evt)
			{
				if (!mxEvent.isShiftDown(evt))
				{
					fadeIn(100);
				}
				else
				{
					fadeOut();
				}
				
				mxEvent.consume(evt);
			}));

			mxEvent.addListener(this.chromelessToolbar, 'mouseleave',  mxUtils.bind(this, function(evt)
			{
				if (!mxEvent.isTouchEvent(evt))
				{
					fadeIn(30);
				}
			}));

			// Shows/hides toolbar for touch devices
			var tol = graph.getTolerance();

			graph.addMouseListener(
			{
			    startX: 0,
			    startY: 0,
			    scrollLeft: 0,
			    scrollTop: 0,
			    mouseDown: function(sender, me)
			    {
			    	this.startX = me.getGraphX();
			    	this.startY = me.getGraphY();
				    this.scrollLeft = graph.container.scrollLeft;
				    this.scrollTop = graph.container.scrollTop;
			    },
			    mouseMove: function(sender, me) {},
			    mouseUp: function(sender, me)
			    {
			    	if (mxEvent.isTouchEvent(me.getEvent()))
			    	{
				    	if ((Math.abs(this.scrollLeft - graph.container.scrollLeft) < tol &&
				    		Math.abs(this.scrollTop - graph.container.scrollTop) < tol) &&
				    		(Math.abs(this.startX - me.getGraphX()) < tol &&
				    		Math.abs(this.startY - me.getGraphY()) < tol))
				    	{
				    		if (parseFloat(ui.chromelessToolbar.style.opacity || 0) > 0)
				    		{
				    			fadeOut();
				    		}
				    		else
				    		{
				    			fadeIn(30);
				    		}
						}
			    	}
			    }
			});
		} // end if toolbar

		// Installs handling of highlight and handling links to relative links and anchors
		if (!this.editor.editable)
		{
			this.addChromelessClickHandler();
		}
	}
	else if (this.editor.extendCanvas)
	{
		/**
		 * Guesses autoTranslate to avoid another repaint (see below).
		 * Works if only the scale of the graph changes or if pages
		 * are visible and the visible pages do not change.
		 */
		var graphViewValidate = graph.view.validate;
		graph.view.validate = function()
		{
			if (this.graph.container != null && mxUtils.hasScrollbars(this.graph.container))
			{
				var pad = this.graph.getPagePadding();
				var size = this.graph.getPageSize();
				
				// Updating scrollbars here causes flickering in quirks and is not needed
				// if zoom method is always used to set the current scale on the graph.
				var tx = this.translate.x;
				var ty = this.translate.y;
				this.translate.x = pad.x - (this.x0 || 0) * size.width;
				this.translate.y = pad.y - (this.y0 || 0) * size.height;
			}
			
			graphViewValidate.apply(this, arguments);
		};
		
		var graphSizeDidChange = graph.sizeDidChange;
		graph.sizeDidChange = function()
		{
			if (this.container != null && mxUtils.hasScrollbars(this.container))
			{
				var pages = this.getPageLayout();
				var pad = this.getPagePadding();
				var size = this.getPageSize();
				
				// Updates the minimum graph size
				var minw = Math.ceil(2 * pad.x + pages.width * size.width);
				var minh = Math.ceil(2 * pad.y + pages.height * size.height);
				
				var min = graph.minimumGraphSize;
				
				// LATER: Fix flicker of scrollbar size in IE quirks mode
				// after delayed call in window.resize event handler
				if (min == null || min.width != minw || min.height != minh)
				{
					graph.minimumGraphSize = new mxRectangle(0, 0, minw, minh);
				}

                // Updates auto-translate to include padding and graph size
				var dx = pad.x - pages.x * size.width;
				var dy = pad.y - pages.y * size.height;
				
				if (!this.autoTranslate && (this.view.translate.x != dx || this.view.translate.y != dy))
				{
					this.autoTranslate = true;
					this.view.x0 = pages.x;
					this.view.y0 = pages.y;

					// NOTE: THIS INVOKES THIS METHOD AGAIN. UNFORTUNATELY THERE IS NO WAY AROUND THIS SINCE THE
					// BOUNDS ARE KNOWN AFTER THE VALIDATION AND SETTING THE TRANSLATE TRIGGERS A REVALIDATION.
					// SHOULD MOVE TRANSLATE/SCALE TO VIEW.
					var tx = graph.view.translate.x;
					var ty = graph.view.translate.y;
					graph.view.setTranslate(dx, dy);

          // LATER: Fix rounding errors for small zoom
          graph.container.scrollLeft += Math.round((dx - tx) * graph.view.scale);
          graph.container.scrollTop += Math.round((dy - ty) * graph.view.scale);
					this.autoTranslate = false;
					
					return;
				}

				graphSizeDidChange.apply(this, arguments);
			}
		};
	}
	
	// Accumulates the zoom factor while the rendering is taking place
	// so that not the complete sequence of zoom steps must be painted
	graph.updateZoomTimeout = null;
	graph.cumulativeZoomFactor = 1;
	
	var cursorPosition = null;

	graph.lazyZoom = function(zoomIn)
	{
		if (this.updateZoomTimeout != null)
		{
			window.clearTimeout(this.updateZoomTimeout);
		}

		// Switches to 1% zoom steps below 15%
		// Lower bound depdends on rounding below
		if (zoomIn)
		{
			if (this.view.scale * this.cumulativeZoomFactor < 0.15)
			{
				this.cumulativeZoomFactor = (this.view.scale + 0.01) / this.view.scale;
			}
			else
			{
				// Uses to 5% zoom steps for better grid rendering in webkit
				// and to avoid rounding errors for zoom steps
				this.cumulativeZoomFactor *= this.zoomFactor;
				this.cumulativeZoomFactor = Math.round(this.view.scale * this.cumulativeZoomFactor * 20) / 20 / this.view.scale;
			}
		}
		else
		{
			if (this.view.scale * this.cumulativeZoomFactor <= 0.15)
			{
				this.cumulativeZoomFactor = (this.view.scale - 0.01) / this.view.scale;
			}
			else
			{
				// Uses to 5% zoom steps for better grid rendering in webkit
				// and to avoid rounding errors for zoom steps
				this.cumulativeZoomFactor /= this.zoomFactor;
				this.cumulativeZoomFactor = Math.round(this.view.scale * this.cumulativeZoomFactor * 20) / 20 / this.view.scale;
			}
		}
		
		this.cumulativeZoomFactor = Math.max(0.01, Math.min(this.view.scale * this.cumulativeZoomFactor, 160) / this.view.scale);
		
        this.updateZoomTimeout = window.setTimeout(mxUtils.bind(this, function()
        {
            var offset = mxUtils.getOffset(graph.container);
            var dx = 0;
            var dy = 0;
            
            if (cursorPosition != null)
            {
                dx = graph.container.offsetWidth / 2 - cursorPosition.x + offset.x;
                dy = graph.container.offsetHeight / 2 - cursorPosition.y + offset.y;
            }

            var prev = this.view.scale;
            this.zoom(this.cumulativeZoomFactor);
            var s = this.view.scale;
            
            if (s != prev)
            {
                if (resize != null)
                {
                		ui.chromelessResize(false, null, dx * (this.cumulativeZoomFactor - 1),
                				dy * (this.cumulativeZoomFactor - 1));
                }
                
                if (mxUtils.hasScrollbars(graph.container) && (dx != 0 || dy != 0))
                {
                    graph.container.scrollLeft -= dx * (this.cumulativeZoomFactor - 1);
                    graph.container.scrollTop -= dy * (this.cumulativeZoomFactor - 1);
                }
            }
            
            this.cumulativeZoomFactor = 1;
            this.updateZoomTimeout = null;
        }), this.lazyZoomDelay);
	};
	
	mxEvent.addMouseWheelListener(mxUtils.bind(this, function(evt, up)
	{
		// Ctrl+wheel (or pinch on touchpad) is a native browser zoom event is OS X
		// LATER: Add support for zoom via pinch on trackpad for Chrome in OS X
		if ((this.dialogs == null || this.dialogs.length == 0) && graph.isZoomWheelEvent(evt))
		{
			var source = mxEvent.getSource(evt);
			
			while (source != null)
			{
				if (source == graph.container)
				{
					cursorPosition = new mxPoint(mxEvent.getClientX(evt), mxEvent.getClientY(evt));
					graph.lazyZoom(up);
					mxEvent.consume(evt);
			
					return;
				}
				
				source = source.parentNode;
			}
		}
	}));
};

/**
 * Creates a temporary graph instance for rendering off-screen content.
 */
EditorUi.prototype.addChromelessToolbarItems = function(addButton)
{
	addButton(mxUtils.bind(this, function(evt)
	{
		this.actions.get('print').funct();
		mxEvent.consume(evt);
	}), Editor.printLargeImage, mxResources.get('print'));	
};

/**
 * Creates a temporary graph instance for rendering off-screen content.
 */
EditorUi.prototype.createTemporaryGraph = function(stylesheet)
{
	var graph = new Graph(document.createElement('div'), null, null, stylesheet);
	graph.resetViewOnRootChange = false;
	graph.setConnectable(false);
	graph.gridEnabled = false;
	graph.autoScroll = false;
	graph.setTooltips(false);
	graph.setEnabled(false);

	// Container must be in the DOM for correct HTML rendering
	graph.container.style.visibility = 'hidden';
	graph.container.style.position = 'absolute';
	graph.container.style.overflow = 'hidden';
	graph.container.style.height = '1px';
	graph.container.style.width = '1px';
	
	return graph;
};

/**
 * 
 */
EditorUi.prototype.addChromelessClickHandler = function()
{
	var hl = urlParams['highlight'];
	
	// Adds leading # for highlight color code
	if (hl != null && hl.length > 0)
	{
		hl = '#' + hl;
	}

	this.editor.graph.addClickHandler(hl);
};

/**
 * 
 */
EditorUi.prototype.toggleFormatPanel = function(forceHide)
{
	this.formatWidth = (forceHide || this.formatWidth > 0) ? 0 : 240;
	this.formatContainer.style.display = (forceHide || this.formatWidth > 0) ? '' : 'none';
	this.refresh();
	this.format.refresh();
	this.fireEvent(new mxEventObject('formatWidthChanged'));
};

/**
 * Adds support for placeholders in labels.
 */
EditorUi.prototype.lightboxFit = function(maxHeight)
{
	if (this.isDiagramEmpty())
	{
		this.editor.graph.view.setScale(1);
	}
	else
	{
		var p = urlParams['border'];
		var border = 60;
		
		if (p != null)
		{
			border = parseInt(p);
		}
		
		// LATER: Use initial graph bounds to avoid rounding errors
		this.editor.graph.maxFitScale = this.lightboxMaxFitScale;
		this.editor.graph.fit(border, null, null, null, null, null, maxHeight);
		this.editor.graph.maxFitScale = null;
	}
};

/**
 * Translates this point by the given vector.
 * 
 * @param {number} dx X-coordinate of the translation.
 * @param {number} dy Y-coordinate of the translation.
 */
EditorUi.prototype.isDiagramEmpty = function()
{
	var model = this.editor.graph.getModel();
	
	return model.getChildCount(model.root) == 1 && model.getChildCount(model.getChildAt(model.root, 0)) == 0;
};

/**
 * Hook for allowing selection and context menu for certain events.
 */
EditorUi.prototype.isSelectionAllowed = function(evt)
{
	return mxEvent.getSource(evt).nodeName == 'SELECT' || (mxEvent.getSource(evt).nodeName == 'INPUT' &&
		mxUtils.isAncestorNode(this.formatContainer, mxEvent.getSource(evt)));
};

/**
 * Installs dialog if browser window is closed without saving
 * This must be disabled during save and image export.
 */
EditorUi.prototype.addBeforeUnloadListener = function()
{
	// Installs dialog if browser window is closed without saving
	// This must be disabled during save and image export
	window.onbeforeunload = mxUtils.bind(this, function()
	{
		if (!this.editor.isChromelessView())
		{
			return this.onBeforeUnload();
		}
	});
};

/**
 * Sets the onbeforeunload for the application
 */
EditorUi.prototype.onBeforeUnload = function()
{
	if (this.editor.modified)
	{
		if (graphDialogIsInsert) {
			return;
		}
		
		if (!Editor.useLocalStorage) {
			return mxResources.get('allChangesLost');//先禁用
		} else {
			this.save(this.editor.getOrCreateFilename());//直接保存
		}
	}
};

/**
 * Opens the current diagram via the window.opener if one exists.
 */
EditorUi.prototype.open = function()
{
	// Cross-domain window access is not allowed in FF, so if we
	// were opened from another domain then this will fail.
	try
	{
		if (window.opener != null && window.opener.openFile != null)
		{
			window.opener.openFile.setConsumer(mxUtils.bind(this, function(xml, filename)
			{
				try
				{
					var doc = mxUtils.parseXml(xml); 
					this.editor.setGraphXml(doc.documentElement);
					this.editor.setModified(false);
					this.editor.undoManager.clear();
					
					if (filename != null)
					{
						this.editor.setFilename(filename);
						this.updateDocumentTitle();
					}
					
					return;
				}
				catch (e)
				{
					mxUtils.alert(mxResources.get('invalidOrMissingFile') + ': ' + e.message);
				}
			}));
		}
	}
	catch(e)
	{
		// ignore
	}
	
	// Fires as the last step if no file was loaded
	this.editor.graph.view.validate();
	
	// Required only in special cases where an initial file is opened
	// and the minimumGraphSize changes and CSS must be updated.
	this.editor.graph.sizeDidChange();
	this.editor.fireEvent(new mxEventObject('resetGraphView'));
};

/**
 * Sets the current menu and element.
 */
EditorUi.prototype.setCurrentMenu = function(menu, elt)
{
	this.currentMenuElt = elt;
	this.currentMenu = menu;
};

/**
 * Resets the current menu and element.
 */
EditorUi.prototype.resetCurrentMenu = function()
{
	this.currentMenuElt = null;
	this.currentMenu = null;
};

/**
 * Hides and destroys the current menu.
 */
EditorUi.prototype.hideCurrentMenu = function()
{
	if (this.currentMenu != null)
	{
		this.currentMenu.hideMenu();
		this.resetCurrentMenu();
	}
};

/**
 * Updates the document title.
 */
EditorUi.prototype.updateDocumentTitle = function()
{
	var title = this.editor.getOrCreateFilename();
	
	if (this.editor.appName != null)
	{
		title += ' - ' + this.editor.appName;
	}
	
	document.title = title;
};

/**
 * Updates the document title.
 */
EditorUi.prototype.createHoverIcons = function()
{
	return new HoverIcons(this.editor.graph);
};

/**
 * Returns the URL for a copy of this editor with no state.
 */
EditorUi.prototype.redo = function()
{
	try
	{
		var graph = this.editor.graph;
		
		if (graph.isEditing())
		{
			document.execCommand('redo', false, null);
		}
		else
		{
			this.editor.undoManager.redo();
		}
	}
	catch (e)
	{
		// ignore all errors
	}
};

/**
 * Returns the URL for a copy of this editor with no state.
 */
EditorUi.prototype.undo = function()
{
	try
	{
		var graph = this.editor.graph;
	
		if (graph.isEditing())
		{
			// Stops editing and executes undo on graph if native undo
			// does not affect current editing value
			var value = graph.cellEditor.textarea.innerHTML;
			document.execCommand('undo', false, null);
	
			if (value == graph.cellEditor.textarea.innerHTML)
			{
				graph.stopEditing(true);
				this.editor.undoManager.undo();
			}
		}
		else
		{
			this.editor.undoManager.undo();
		}
	}
	catch (e)
	{
		// ignore all errors
	}
};

/**
 * Returns the URL for a copy of this editor with no state.
 */
EditorUi.prototype.canRedo = function()
{
	return this.editor.graph.isEditing() || this.editor.undoManager.canRedo();
};

/**
 * Returns the URL for a copy of this editor with no state.
 */
EditorUi.prototype.canUndo = function()
{
	return this.editor.graph.isEditing() || this.editor.undoManager.canUndo();
};

/**
 * 
 */
EditorUi.prototype.getEditBlankXml = function()
{
	return mxUtils.getXml(this.editor.getGraphXml());
};

/**
 * Returns the URL for a copy of this editor with no state.
 */
EditorUi.prototype.getUrl = function(pathname)
{
	var href = (pathname != null) ? pathname : window.location.pathname;
	var parms = (href.indexOf('?') > 0) ? 1 : 0;
	
	// Removes template URL parameter for new blank diagram
	for (var key in urlParams)
	{
		if (parms == 0)
		{
			href += '?';
		}
		else
		{
			href += '&';
		}
	
		href += key + '=' + urlParams[key];
		parms++;
	}
	
	return href;
};

/**
 * Specifies if the graph has scrollbars.
 */
EditorUi.prototype.setScrollbars = function(value)
{
	var graph = this.editor.graph;
	var prev = graph.container.style.overflow;
	graph.scrollbars = value;
	this.editor.updateGraphComponents();

	if (prev != graph.container.style.overflow)
	{
		if (graph.container.style.overflow == 'hidden')
		{
			var t = graph.view.translate;
			graph.view.setTranslate(t.x - graph.container.scrollLeft / graph.view.scale, t.y - graph.container.scrollTop / graph.view.scale);
			graph.container.scrollLeft = 0;
			graph.container.scrollTop = 0;
			graph.minimumGraphSize = null;
			graph.sizeDidChange();
		}
		else
		{
			var dx = graph.view.translate.x;
			var dy = graph.view.translate.y;

			graph.view.translate.x = 0;
			graph.view.translate.y = 0;
			graph.sizeDidChange();
			graph.container.scrollLeft -= Math.round(dx * graph.view.scale);
			graph.container.scrollTop -= Math.round(dy * graph.view.scale);
		}
	}
	
	this.fireEvent(new mxEventObject('scrollbarsChanged'));
};

/**
 * Returns true if the graph has scrollbars.
 */
EditorUi.prototype.hasScrollbars = function()
{
	return this.editor.graph.scrollbars;
};

/**
 * Resets the state of the scrollbars.
 */
EditorUi.prototype.resetScrollbars = function()
{
	var graph = this.editor.graph;
	
	if (!this.editor.extendCanvas)
	{
		graph.container.scrollTop = 0;
		graph.container.scrollLeft = 0;
	
		if (!mxUtils.hasScrollbars(graph.container))
		{
			graph.view.setTranslate(0, 0);
		}
	}
	else if (!this.editor.isChromelessView())
	{
		if (mxUtils.hasScrollbars(graph.container))
		{
			if (graph.pageVisible)
			{
				var pad = graph.getPagePadding();
				graph.container.scrollTop = Math.floor(pad.y - this.editor.initialTopSpacing) - 1;
				graph.container.scrollLeft = Math.floor(Math.min(pad.x,
					(graph.container.scrollWidth - graph.container.clientWidth) / 2)) - 1;

				// Scrolls graph to visible area
				var bounds = graph.getGraphBounds();
				
				if (bounds.width > 0 && bounds.height > 0)
				{
					if (bounds.x > graph.container.scrollLeft + graph.container.clientWidth * 0.9)
					{
						graph.container.scrollLeft = Math.min(bounds.x + bounds.width - graph.container.clientWidth, bounds.x - 10);
					}
					
					if (bounds.y > graph.container.scrollTop + graph.container.clientHeight * 0.9)
					{
						graph.container.scrollTop = Math.min(bounds.y + bounds.height - graph.container.clientHeight, bounds.y - 10);
					}
				}
			}
			else
			{
				var bounds = graph.getGraphBounds();
				var width = Math.max(bounds.width, graph.scrollTileSize.width * graph.view.scale);
				var height = Math.max(bounds.height, graph.scrollTileSize.height * graph.view.scale);
				graph.container.scrollTop = Math.floor(Math.max(0, bounds.y - Math.max(20, (graph.container.clientHeight - height) / 4)));
				graph.container.scrollLeft = Math.floor(Math.max(0, bounds.x - Math.max(0, (graph.container.clientWidth - width) / 2)));
			}
		}
		else
		{
			// This code is not actively used since the default for scrollbars is always true
			if (graph.pageVisible)
			{
				var b = graph.view.getBackgroundPageBounds();
				graph.view.setTranslate(Math.floor(Math.max(0, (graph.container.clientWidth - b.width) / 2) - b.x),
					Math.floor(Math.max(0, (graph.container.clientHeight - b.height) / 2) - b.y));
			}
			else
			{
				var bounds = graph.getGraphBounds();
				graph.view.setTranslate(Math.floor(Math.max(0, Math.max(0, (graph.container.clientWidth - bounds.width) / 2) - bounds.x)),
					Math.floor(Math.max(0, Math.max(20, (graph.container.clientHeight - bounds.height) / 4)) - bounds.y));
			}
		}
	}
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setPageVisible = function(value)
{
	var graph = this.editor.graph;
	var hasScrollbars = mxUtils.hasScrollbars(graph.container);
	var tx = 0;
	var ty = 0;
	
	if (hasScrollbars)
	{
		tx = graph.view.translate.x * graph.view.scale - graph.container.scrollLeft;
		ty = graph.view.translate.y * graph.view.scale - graph.container.scrollTop;
	}
	
	graph.pageVisible = value;
	graph.pageBreaksVisible = value; 
	graph.preferPageSize = value;
	graph.view.validateBackground();

	// Workaround for possible handle offset
	if (hasScrollbars)
	{
		var cells = graph.getSelectionCells();
		graph.clearSelection();
		graph.setSelectionCells(cells);
	}
	
	// Calls updatePageBreaks
	graph.sizeDidChange();
	
	if (hasScrollbars)
	{
		graph.container.scrollLeft = graph.view.translate.x * graph.view.scale - tx;
		graph.container.scrollTop = graph.view.translate.y * graph.view.scale - ty;
	}
	
	this.fireEvent(new mxEventObject('pageViewChanged'));
};

/**
 * Change types
 */
function ChangePageSetup(ui, color, image, format)
{
	this.ui = ui;
	this.color = color;
	this.previousColor = color;
	this.image = image;
	this.previousImage = image;
	this.format = format;
	this.previousFormat = format;
	
	// Needed since null are valid values for color and image
	this.ignoreColor = false;
	this.ignoreImage = false;
}

/**
 * Implementation of the undoable page rename.
 */
ChangePageSetup.prototype.execute = function()
{
	var graph = this.ui.editor.graph;
	
	if (!this.ignoreColor)
	{
		this.color = this.previousColor;
		var tmp = graph.background;
		this.ui.setBackgroundColor(this.previousColor);
		this.previousColor = tmp;
	}
	
	if (!this.ignoreImage)
	{
		this.image = this.previousImage;
		var tmp = graph.backgroundImage;
		this.ui.setBackgroundImage(this.previousImage);
		this.previousImage = tmp;
	}
	
	if (this.previousFormat != null)
	{
		this.format = this.previousFormat;
		var tmp = graph.pageFormat;
		
		if (this.previousFormat.width != tmp.width ||
			this.previousFormat.height != tmp.height)
		{
			this.ui.setPageFormat(this.previousFormat);
			this.previousFormat = tmp;
		}
	}

    if (this.foldingEnabled != null && this.foldingEnabled != this.ui.editor.graph.foldingEnabled)
    {
    		this.ui.setFoldingEnabled(this.foldingEnabled);
        this.foldingEnabled = !this.foldingEnabled;
    }
};

// Registers codec for ChangePageSetup
(function()
{
	var codec = new mxObjectCodec(new ChangePageSetup(),  ['ui', 'previousColor', 'previousImage', 'previousFormat']);

	codec.afterDecode = function(dec, node, obj)
	{
		obj.previousColor = obj.color;
		obj.previousImage = obj.image;
		obj.previousFormat = obj.format;

        if (obj.foldingEnabled != null)
        {
        		obj.foldingEnabled = !obj.foldingEnabled;
        }
       
		return obj;
	};
	
	mxCodecRegistry.register(codec);
})();

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setBackgroundColor = function(value)
{
	this.editor.graph.background = value;
	this.editor.graph.view.validateBackground();

	this.fireEvent(new mxEventObject('backgroundColorChanged'));
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setFoldingEnabled = function(value)
{
	this.editor.graph.foldingEnabled = value;
	this.editor.graph.view.revalidate();
	
	this.fireEvent(new mxEventObject('foldingEnabledChanged'));
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setPageFormat = function(value)
{
	this.editor.graph.pageFormat = value;
	
	if (!this.editor.graph.pageVisible)
	{
		this.actions.get('pageView').funct();
	}
	else
	{
		this.editor.graph.view.validateBackground();
		this.editor.graph.sizeDidChange();
	}

	this.fireEvent(new mxEventObject('pageFormatChanged'));
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setPageScale = function(value)
{
	this.editor.graph.pageScale = value;
	
	if (!this.editor.graph.pageVisible)
	{
		this.actions.get('pageView').funct();
	}
	else
	{
		this.editor.graph.view.validateBackground();
		this.editor.graph.sizeDidChange();
	}

	this.fireEvent(new mxEventObject('pageScaleChanged'));
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setGridColor = function(value)
{
	this.editor.graph.view.gridColor = value;
	this.editor.graph.view.validateBackground();
	this.fireEvent(new mxEventObject('gridColorChanged'));
};

/**
 * Updates the states of the given undo/redo items.
 */
EditorUi.prototype.addUndoListener = function()
{
	var undo = this.actions.get('undo');
	var redo = this.actions.get('redo');
	
	var undoMgr = this.editor.undoManager;
	
    var undoListener = mxUtils.bind(this, function()
    {
    	undo.setEnabled(this.canUndo());
    	redo.setEnabled(this.canRedo());
    });

    undoMgr.addListener(mxEvent.ADD, undoListener);
    undoMgr.addListener(mxEvent.UNDO, undoListener);
    undoMgr.addListener(mxEvent.REDO, undoListener);
    undoMgr.addListener(mxEvent.CLEAR, undoListener);
	
	// Overrides cell editor to update action states
	var cellEditorStartEditing = this.editor.graph.cellEditor.startEditing;
	
	this.editor.graph.cellEditor.startEditing = function()
	{
		cellEditorStartEditing.apply(this, arguments);
		undoListener();
	};
	
	var cellEditorStopEditing = this.editor.graph.cellEditor.stopEditing;
	
	this.editor.graph.cellEditor.stopEditing = function(cell, trigger)
	{
		cellEditorStopEditing.apply(this, arguments);
		undoListener();
	};
	
	// Updates the button states once
    undoListener();
};

/**
* Updates the states of the given toolbar items based on the selection.
*/
EditorUi.prototype.updateActionStates = function()
{
	var graph = this.editor.graph;
	var selected = !graph.isSelectionEmpty();
	var vertexSelected = false;
	var edgeSelected = false;

	var cells = graph.getSelectionCells();
	
	if (cells != null)
	{
    	for (var i = 0; i < cells.length; i++)
    	{
    		var cell = cells[i];
    		
    		if (graph.getModel().isEdge(cell))
    		{
    			edgeSelected = true;
    		}
    		
    		if (graph.getModel().isVertex(cell))
    		{
    			vertexSelected = true;
    		}
    		
    		if (edgeSelected && vertexSelected)
			{
				break;
			}
		}
	}
	
	// Updates action states
	var actions = ['cut', 'copy', 'bold', 'italic', 'underline', 'delete', 'duplicate',
	               'editStyle', 'editTooltip', 'editLink', 'backgroundColor', 'borderColor',
	               'edit', 'toFront', 'toBack', 'lockUnlock', 'solid', 'dashed', 'pasteSize',
	               'dotted', 'fillColor', 'gradientColor', 'shadow', 'fontColor',
	               'formattedText', 'rounded', 'toggleRounded', 'sharp', 'strokeColor'];
	
	for (var i = 0; i < actions.length; i++)
	{
		this.actions.get(actions[i]).setEnabled(selected);
	}
	
	this.actions.get('setAsDefaultStyle').setEnabled(graph.getSelectionCount() == 1);
	this.actions.get('clearWaypoints').setEnabled(!graph.isSelectionEmpty());
	this.actions.get('copySize').setEnabled(graph.getSelectionCount() == 1);
	this.actions.get('turn').setEnabled(!graph.isSelectionEmpty());
	this.actions.get('curved').setEnabled(edgeSelected);
	this.actions.get('rotation').setEnabled(vertexSelected);
	this.actions.get('wordWrap').setEnabled(vertexSelected);
	this.actions.get('autosize').setEnabled(vertexSelected);
   	var oneVertexSelected = vertexSelected && graph.getSelectionCount() == 1;
	this.actions.get('group').setEnabled(graph.getSelectionCount() > 1 ||
		(oneVertexSelected && !graph.isContainer(graph.getSelectionCell())));
	this.actions.get('ungroup').setEnabled(graph.getSelectionCount() == 1 &&
		(graph.getModel().getChildCount(graph.getSelectionCell()) > 0 ||
		(oneVertexSelected && graph.isContainer(graph.getSelectionCell()))));
   	this.actions.get('removeFromGroup').setEnabled(oneVertexSelected &&
   		graph.getModel().isVertex(graph.getModel().getParent(graph.getSelectionCell())));

	// Updates menu states
   	var state = graph.view.getState(graph.getSelectionCell());
    this.menus.get('navigation').setEnabled(selected || graph.view.currentRoot != null);
    this.actions.get('collapsible').setEnabled(vertexSelected &&
    	(graph.isContainer(graph.getSelectionCell()) || graph.model.getChildCount(graph.getSelectionCell()) > 0));
    this.actions.get('home').setEnabled(graph.view.currentRoot != null);
    this.actions.get('exitGroup').setEnabled(graph.view.currentRoot != null);
    this.actions.get('enterGroup').setEnabled(graph.getSelectionCount() == 1 && graph.isValidRoot(graph.getSelectionCell()));
    var foldable = graph.getSelectionCount() == 1 && graph.isCellFoldable(graph.getSelectionCell());
    this.actions.get('expand').setEnabled(foldable);
    this.actions.get('collapse').setEnabled(foldable);
    this.actions.get('editLink').setEnabled(graph.getSelectionCount() == 1);
    this.actions.get('openLink').setEnabled(graph.getSelectionCount() == 1 &&
    	graph.getLinkForCell(graph.getSelectionCell()) != null);
    this.actions.get('guides').setEnabled(graph.isEnabled());
    this.actions.get('grid').setEnabled(!this.editor.chromeless || this.editor.editable);

    var unlocked = graph.isEnabled() && !graph.isCellLocked(graph.getDefaultParent());
    this.menus.get('layout').setEnabled(unlocked);
    this.menus.get('insert').setEnabled(unlocked);
    this.menus.get('direction').setEnabled(unlocked && vertexSelected);
    this.menus.get('align').setEnabled(unlocked && vertexSelected && graph.getSelectionCount() > 1);
    this.menus.get('distribute').setEnabled(unlocked && vertexSelected && graph.getSelectionCount() > 1);
    this.actions.get('selectVertices').setEnabled(unlocked);
    this.actions.get('selectEdges').setEnabled(unlocked);
    this.actions.get('selectAll').setEnabled(unlocked);
    this.actions.get('selectNone').setEnabled(unlocked);
    
    this.updatePasteActionStates();
};

/**
 * Refreshes the viewport.
 */
EditorUi.prototype.refresh = function(sizeDidChange)
{
	sizeDidChange = (sizeDidChange != null) ? sizeDidChange : true;
	
	var quirks = mxClient.IS_IE && (document.documentMode == null || document.documentMode == 5);
	var w = this.container.clientWidth;
	var h = this.container.clientHeight;

	if (this.container == document.body)
	{
		w = document.body.clientWidth || document.documentElement.clientWidth;
		h = (quirks) ? document.body.clientHeight || document.documentElement.clientHeight : document.documentElement.clientHeight;
	}
	
	// Workaround for bug on iOS see
	// http://stackoverflow.com/questions/19012135/ios-7-ipad-safari-landscape-innerheight-outerheight-layout-issue
	// FIXME: Fix if footer visible
	var off = 0;
	
	if (mxClient.IS_IOS && !window.navigator.standalone)
	{
		if (window.innerHeight != document.documentElement.clientHeight)
		{
			off = document.documentElement.clientHeight - window.innerHeight;
			window.scrollTo(0, 0);
		}
	}
	
	var effHsplitPosition = Math.max(0, Math.min(this.hsplitPosition, w - this.splitSize - 20));

	var tmp = 0;
	
	if (this.menubar != null)
	{
		tmp += this.menubarHeight;
	}
	
	if (this.toolbar != null)
	{
		this.toolbarContainer.style.top = this.menubarHeight + 'px';
		this.toolbarContainer.style.height = this.toolbarHeight + 'px';
		tmp += this.toolbarHeight;
	}
	
	if (tmp > 0 && !mxClient.IS_QUIRKS)
	{
		tmp += 1;
	}
	
	var sidebarFooterHeight = 0;
	
	if (this.sidebarFooterContainer != null)
	{
		var bottom = this.footerHeight + off;
		sidebarFooterHeight = Math.max(0, Math.min(h - tmp - bottom, this.sidebarFooterHeight));
		this.sidebarFooterContainer.style.width = effHsplitPosition + 'px';
		this.sidebarFooterContainer.style.height = sidebarFooterHeight + 'px';
		this.sidebarFooterContainer.style.bottom = bottom + 'px';
	}
	
	var fw = (this.format != null) ? this.formatWidth : 0;
	this.sidebarContainer.style.top = tmp + 'px';
	this.sidebarContainer.style.width = effHsplitPosition + 'px';
	this.formatContainer.style.top = tmp + 'px';
	this.formatContainer.style.width = fw + 'px';
	this.formatContainer.style.display = (this.format != null) ? '' : 'none';
	
	this.diagramContainer.style.left = (this.hsplit.parentNode != null) ? (effHsplitPosition + this.splitSize) + 'px' : '0px';
	this.diagramContainer.style.top = this.sidebarContainer.style.top;
	this.footerContainer.style.height = this.footerHeight + 'px';
	this.hsplit.style.top = this.sidebarContainer.style.top;
	this.hsplit.style.bottom = (this.footerHeight + off) + 'px';
	this.hsplit.style.left = effHsplitPosition + 'px';
	
	if (this.tabContainer != null)
	{
		this.tabContainer.style.left = this.diagramContainer.style.left;
	}
	
	if (quirks)
	{
		this.toolbarContainer.style.width = w + 'px';
		var sidebarHeight = Math.max(0, h - this.footerHeight - this.menubarHeight - this.toolbarHeight);
		this.sidebarContainer.style.height = (sidebarHeight - sidebarFooterHeight) + 'px';
		this.formatContainer.style.height = sidebarHeight + 'px';
		this.diagramContainer.style.width = (this.hsplit.parentNode != null) ? Math.max(0, w - effHsplitPosition - this.splitSize - fw) + 'px' : w + 'px';
		this.footerContainer.style.width = w + 'px';
		var diagramHeight = Math.max(0, h - this.footerHeight - this.menubarHeight - this.toolbarHeight);
		
		if (this.tabContainer != null)
		{
			this.tabContainer.style.width = this.diagramContainer.style.width;
			this.tabContainer.style.bottom = (this.footerHeight + off) + 'px';
			diagramHeight -= this.tabContainer.clientHeight;
		}
		
		this.diagramContainer.style.height = diagramHeight + 'px';
		this.hsplit.style.height = diagramHeight + 'px';
	}
	else
	{
		if (this.footerHeight > 0)
		{
			this.footerContainer.style.bottom = off + 'px';
		}
		
		this.diagramContainer.style.right = fw + 'px';
		var th = 0;
		
		if (this.tabContainer != null)
		{
			this.tabContainer.style.bottom = (this.footerHeight + off) + 'px';
			this.tabContainer.style.right = this.diagramContainer.style.right;
			th = this.tabContainer.clientHeight;
		}
		
		this.sidebarContainer.style.bottom = (this.footerHeight + sidebarFooterHeight + off) + 'px';
		this.formatContainer.style.bottom = (this.footerHeight + off) + 'px';
		this.diagramContainer.style.bottom = (this.footerHeight + off + th) + 'px';
	}
	
	if (sizeDidChange)
	{
		this.editor.graph.sizeDidChange();
	}
};

/**
 * Creates the required containers.
 */
EditorUi.prototype.createTabContainer = function()
{
	return null;
};

/**
 * Creates the required containers.
 */
EditorUi.prototype.createDivs = function()
{
	this.menubarContainer = this.createDiv('geMenubarContainer');
	this.toolbarContainer = this.createDiv('geToolbarContainer');
	this.sidebarContainer = this.createDiv('geSidebarContainer');
	this.formatContainer = this.createDiv('geSidebarContainer geFormatContainer');
	this.diagramContainer = this.createDiv('geDiagramContainer');
	this.footerContainer = this.createDiv('geFooterContainer');
	this.hsplit = this.createDiv('geHsplit');
	this.hsplit.setAttribute('title', mxResources.get('collapseExpand'));

	// Sets static style for containers
	this.menubarContainer.style.top = '0px';
	this.menubarContainer.style.left = '0px';
	this.menubarContainer.style.right = '0px';
	this.toolbarContainer.style.left = '0px';
	this.toolbarContainer.style.right = '0px';
	this.sidebarContainer.style.left = '0px';
	this.formatContainer.style.right = '0px';
	this.formatContainer.style.zIndex = '1';
	this.diagramContainer.style.right = ((this.format != null) ? this.formatWidth : 0) + 'px';
	this.footerContainer.style.left = '0px';
	this.footerContainer.style.right = '0px';
	this.footerContainer.style.bottom = '0px';
	this.footerContainer.style.zIndex = mxPopupMenu.prototype.zIndex - 2;
	this.hsplit.style.width = this.splitSize + 'px';
	this.sidebarFooterContainer = this.createSidebarFooterContainer();
	
	if (this.sidebarFooterContainer)
	{
		this.sidebarFooterContainer.style.left = '0px';
	}
	
	if (!this.editor.chromeless)
	{
		this.tabContainer = this.createTabContainer();
	}
	else
	{
		this.diagramContainer.style.border = 'none';
	}
};

/**
 * Hook for sidebar footer container. This implementation returns null.
 */
EditorUi.prototype.createSidebarFooterContainer = function()
{
	return null;
};

/**
 * Creates the required containers.
 */
EditorUi.prototype.createUi = function()
{
	// Creates menubar
	this.menubar = (this.editor.chromeless) ? null : this.menus.createMenubar(this.createDiv('geMenubar'));
	
	if (this.menubar != null)
	{
		this.menubarContainer.appendChild(this.menubar.container);
	}
	
	// Adds status bar in menubar
	if (this.menubar != null)
	{
		this.statusContainer = this.createStatusContainer();
	
		// Connects the status bar to the editor status
		this.editor.addListener('statusChanged', mxUtils.bind(this, function()
		{
			this.setStatusText(this.editor.getStatus());
		}));
	
		this.setStatusText(this.editor.getStatus());
		this.menubar.container.appendChild(this.statusContainer);
		
		// Inserts into DOM
		this.container.appendChild(this.menubarContainer);
	}

	// Creates the sidebar
	this.sidebar = (this.editor.chromeless) ? null : this.createSidebar(this.sidebarContainer);
	
	if (this.sidebar != null)
	{
		this.container.appendChild(this.sidebarContainer);
	}
	
	// Creates the format sidebar
	this.format = (this.editor.chromeless || !this.formatEnabled) ? null : this.createFormat(this.formatContainer);
	
	if (this.format != null)
	{
		this.container.appendChild(this.formatContainer);
	}
	
	// Creates the footer
	var footer = (this.editor.chromeless) ? null : this.createFooter();
	
	if (footer != null)
	{
		this.footerContainer.appendChild(footer);
		this.container.appendChild(this.footerContainer);
	}

	if (this.sidebar != null && this.sidebarFooterContainer)
	{
		this.container.appendChild(this.sidebarFooterContainer);		
	}

	this.container.appendChild(this.diagramContainer);

	if (this.container != null && this.tabContainer != null)
	{
		this.container.appendChild(this.tabContainer);
	}

	// Creates toolbar
	this.toolbar = (this.editor.chromeless) ? null : this.createToolbar(this.createDiv('geToolbar'));
	
	if (this.toolbar != null)
	{
		this.toolbarContainer.appendChild(this.toolbar.container);
		this.container.appendChild(this.toolbarContainer);
	}

	// HSplit
	if (this.sidebar != null)
	{
		this.container.appendChild(this.hsplit);
		
		this.addSplitHandler(this.hsplit, true, 0, mxUtils.bind(this, function(value)
		{
			this.hsplitPosition = value;
			this.refresh();
		}));
	}
};

/**
 * Creates a new toolbar for the given container.
 */
EditorUi.prototype.createStatusContainer = function()
{
	var container = document.createElement('a');
	container.className = 'geItem geStatus';
	
	if (screen.width < 420)
	{
		container.style.maxWidth = Math.max(20, screen.width - 320) + 'px';
		container.style.overflow = 'hidden';
	}
	
	return container;
};

/**
 * Creates a new toolbar for the given container.
 */
EditorUi.prototype.setStatusText = function(value)
{
	this.statusContainer.innerHTML = value;
};

/**
 * Creates a new toolbar for the given container.
 */
EditorUi.prototype.createToolbar = function(container)
{
	return new Toolbar(this, container);
};

/**
 * Creates a new sidebar for the given container.
 */
EditorUi.prototype.createSidebar = function(container)
{
	return new Sidebar(this, container);
};

/**
 * Creates a new sidebar for the given container.
 */
EditorUi.prototype.createFormat = function(container)
{
	return new Format(this, container);
};

/**
 * Creates and returns a new footer.
 */
EditorUi.prototype.createFooter = function()
{
	return this.createDiv('geFooter');
};

/**
 * Creates the actual toolbar for the toolbar container.
 */
EditorUi.prototype.createDiv = function(classname)
{
	var elt = document.createElement('div');
	elt.className = classname;
	
	return elt;
};

/**
 * Updates the states of the given undo/redo items.
 */
EditorUi.prototype.addSplitHandler = function(elt, horizontal, dx, onChange)
{
	var start = null;
	var initial = null;
	var ignoreClick = true;
	var last = null;

	// Disables built-in pan and zoom in IE10 and later
	if (mxClient.IS_POINTER)
	{
		elt.style.touchAction = 'none';
	}
	
	var getValue = mxUtils.bind(this, function()
	{
		var result = parseInt(((horizontal) ? elt.style.left : elt.style.bottom));
	
		// Takes into account hidden footer
		if (!horizontal)
		{
			result = result + dx - this.footerHeight;
		}
		
		return result;
	});

	function moveHandler(evt)
	{
		if (start != null)
		{
			var pt = new mxPoint(mxEvent.getClientX(evt), mxEvent.getClientY(evt));
			onChange(Math.max(0, initial + ((horizontal) ? (pt.x - start.x) : (start.y - pt.y)) - dx));
			mxEvent.consume(evt);
			
			if (initial != getValue())
			{
				ignoreClick = true;
				last = null;
			}
		}
	};
	
	function dropHandler(evt)
	{
		moveHandler(evt);
		initial = null;
		start = null;
	};
	
	mxEvent.addGestureListeners(elt, function(evt)
	{
		start = new mxPoint(mxEvent.getClientX(evt), mxEvent.getClientY(evt));
		initial = getValue();
		ignoreClick = false;
		mxEvent.consume(evt);
	});
	
	mxEvent.addListener(elt, 'click', mxUtils.bind(this, function(evt)
	{
		if (!ignoreClick && this.hsplitClickEnabled)
		{
			var next = (last != null) ? last - dx : 0;
			last = getValue();
			onChange(next);
			mxEvent.consume(evt);
		}
	}));

	mxEvent.addGestureListeners(document, null, moveHandler, dropHandler);
	
	this.destroyFunctions.push(function()
	{
		mxEvent.removeGestureListeners(document, null, moveHandler, dropHandler);
	});	
};

/**
 * Displays a print dialog.
 */
EditorUi.prototype.showDialog = function(elt, w, h, modal, closable, onClose, noScroll, trasparent)
{
	this.editor.graph.tooltipHandler.hideTooltip();
	
	if (this.dialogs == null)
	{
		this.dialogs = [];
	}
	
	this.dialog = new Dialog(this, elt, w, h, modal, closable, onClose, noScroll, trasparent);
	this.dialogs.push(this.dialog);
};

/**
 * Displays a print dialog.
 */
EditorUi.prototype.hideDialog = function(cancel)
{
	if (this.dialogs != null && this.dialogs.length > 0)
	{
		var dlg = this.dialogs.pop();
		dlg.close(cancel);
		
		this.dialog = (this.dialogs.length > 0) ? this.dialogs[this.dialogs.length - 1] : null;

		if (this.dialog == null && this.editor.graph.container.style.visibility != 'hidden')
		{
			this.editor.graph.container.focus();
		}
		
		mxUtils.clearSelection();
		this.editor.fireEvent(new mxEventObject('hideDialog'));
	}
};

/**
 * Display a color dialog.
 */
EditorUi.prototype.pickColor = function(color, apply)
{
	var graph = this.editor.graph;
	var selState = graph.cellEditor.saveSelection();
	var h = 226 + ((Math.ceil(ColorDialog.prototype.presetColors.length / 12) +
		Math.ceil(ColorDialog.prototype.defaultColors.length / 12)) * 17);
	
	var dlg = new ColorDialog(this, color || 'none', function(color)
	{
		graph.cellEditor.restoreSelection(selState);
		apply(color);
	}, function()
	{
		graph.cellEditor.restoreSelection(selState);
	});
	this.showDialog(dlg.container, 230, h, true, false);
	dlg.init();
};

/**
 * Adds the label menu items to the given menu and parent.
 */
EditorUi.prototype.openFile = function()
{
	// Closes dialog after open
	window.openFile = new OpenFile(mxUtils.bind(this, function(cancel)
	{
		this.hideDialog(cancel);
	}));

	// Removes openFile if dialog is closed
	this.showDialog(new OpenDialog(this).container, (Editor.useLocalStorage) ? 640 : 320,
			(Editor.useLocalStorage) ? 480 : 220, true, true, function()
	{
		window.openFile = null;
	});
};

/**
 * Extracs the graph model from the given HTML data from a data transfer event.
 */
EditorUi.prototype.extractGraphModelFromHtml = function(data)
{
	var result = null;
	
	try
	{
    	var idx = data.indexOf('&lt;mxGraphModel ');
    	
    	if (idx >= 0)
    	{
    		var idx2 = data.lastIndexOf('&lt;/mxGraphModel&gt;');
    		
    		if (idx2 > idx)
    		{
    			result = data.substring(idx, idx2 + 21).replace(/&gt;/g, '>').
    				replace(/&lt;/g, '<').replace(/\\&quot;/g, '"').replace(/\n/g, '');
    		}
    	}
	}
	catch (e)
	{
		// ignore
	}
	
	return result;
};

/**
 * Opens the given files in the editor.
 */
EditorUi.prototype.extractGraphModelFromEvent = function(evt)
{
	var result = null;
	var data = null;
	
	if (evt != null)
	{
		var provider = (evt.dataTransfer != null) ? evt.dataTransfer : evt.clipboardData;
		
		if (provider != null)
		{
			if (document.documentMode == 10 || document.documentMode == 11)
			{
				data = provider.getData('Text');
			}
			else
			{
				data = (mxUtils.indexOf(provider.types, 'text/html') >= 0) ? provider.getData('text/html') : null;
			
				if (mxUtils.indexOf(provider.types, 'text/plain' && (data == null || data.length == 0)))
				{
					data = provider.getData('text/plain');
				}
			}

			if (data != null)
			{
				data = this.editor.graph.zapGremlins(mxUtils.trim(data));
				
				// Tries parsing as HTML document with embedded XML
				var xml =  this.extractGraphModelFromHtml(data);
				
				if (xml != null)
				{
					data = xml;
				}
			}		
		}
	}
	
	if (data != null && this.isCompatibleString(data))
	{
		result = data;
	}
	
	return result;
};

/**
 * Hook for subclassers to return true if event data is a supported format.
 * This implementation always returns false.
 */
EditorUi.prototype.isCompatibleString = function(data)
{
	return false;
};

/**
 * Adds the label menu items to the given menu and parent.
 */
EditorUi.prototype.saveFile = function(forceDialog)
{
	if (!forceDialog && this.editor.filename != null)
	{
		this.save(this.editor.getOrCreateFilename());
	}
	else
	{
		var dlg = new FilenameDialog(this, this.editor.getOrCreateFilename(), mxResources.get('save'), mxUtils.bind(this, function(name)
		{
			this.save(name);
		}), null, mxUtils.bind(this, function(name)
		{
			if (name != null && name.length > 0)
			{
				return true;
			}
			
			mxUtils.confirm(mxResources.get('invalidName'));
			
			return false;
		}));
		this.showDialog(dlg.container, 300, 100, true, true);
		dlg.init();
	}
};
/**
 * Saves the current Png under the given filename.
 */
EditorUi.prototype.savePng = function() {
	var dlg = new FilenameDialog(this, this.editor.getOrCreateFilename(), mxResources.get('save'), mxUtils.bind(this, function(name)
	{
		editorUIInstance.exportImage(1, "#ffffff", true, null, true, 50, null, "png", function(base64){
			downloadFile(base64,name)
		})
	}), null, mxUtils.bind(this, function(name)
	{
		if (name != null && name.length > 0)
		{
			return true;
		}
		mxUtils.confirm(mxResources.get('invalidName'));
		return false;
	}));
	this.showDialog(dlg.container, 300, 100, true, true);
	dlg.init();

	function downloadFile(content, fileName) {
		var base64ToBlob = function(content) {
			var parts = content.split(';base64,');
			var contentType = parts[0].split(':')[1];
			var raw = window.atob(parts[1]);
			var rawLength = raw.length;
			var uInt8Array = new Uint8Array(rawLength);
			for(let i = 0; i < rawLength; ++i) {
				uInt8Array[i] = raw.charCodeAt(i);
			}
			return new Blob([uInt8Array], {
				type: contentType
			});
		};
		var aLink = document.createElement('a');
		var blob = base64ToBlob(content); //new Blob([content]);
		var evt = document.createEvent("HTMLEvents");
		evt.initEvent("click", true, true); //initEvent 
		aLink.download = fileName;
		aLink.href = URL.createObjectURL(blob);
		aLink.click();
	};

}
/**
 * Saves the current Xml under the given filename.
 */
EditorUi.prototype.saveXml = function() {
var dlg = new FilenameDialog(this, this.editor.getOrCreateFilename(), mxResources.get('save'), mxUtils.bind(this, function(name)
	{
		var xmlData = mxUtils.getXml(editorUIInstance.editor.getGraphXml());
		downloadFile(xmlData,name)
	}), null, mxUtils.bind(this, function(name)
	{
		if (name != null && name.length > 0)
		{
			return true;
		}
		mxUtils.confirm(mxResources.get('invalidName'));
		return false;
	}));
	this.showDialog(dlg.container, 300, 100, true, true);
	dlg.init();


	function downloadFile(content, fileName) {
		var aLink = document.createElement('a');	
		var evt = document.createEvent("HTMLEvents");	
		var blob = new Blob([content],{
			type:'application/xml'
		});
		evt.initEvent("click", true, true); //initEvent 
		aLink.download = fileName;
		aLink.href = URL.createObjectURL(blob);
		aLink.click();	
	};
}
/**
 * Saves the current graph under the given filename.
 */
EditorUi.prototype.save = function(name)
{
	if (name != null)
	{
		if (this.editor.graph.isEditing())
		{
			this.editor.graph.stopEditing();
		}
		
		var xml = mxUtils.getXml(this.editor.getGraphXml());
		
		try
		{
			if (Editor.useLocalStorage)
			{
				if(window.saveGraph) {
					window.saveGraph(xml);
				}
			}
			else
			{
				if (xml.length < MAX_REQUEST_SIZE)
				{
					// new mxXmlRequest(SAVE_URL, 'filename=' + encodeURIComponent(name) +
						// '&xml=' + encodeURIComponent(xml)).simulate(document, '_blank');
					Socket.emitChange(xml);

					 TFL.use('tips', function() {
						TFL.tips.showFlash(t("Saved successfully","/internation") + '！', true, 1000);
					});
				}
				else
				{
					mxUtils.alert(mxResources.get('drawingTooLarge'));
					mxUtils.popup(xml);
					
					return;
				}
			}

			this.editor.setModified(false);
			this.editor.setFilename(name);
			this.updateDocumentTitle();
		}
		catch (e)
		{
			this.editor.setStatus(mxUtils.htmlEntities(mxResources.get('errorSavingFile')));
		}
	}
};

/**
 * Executes the given layout.
 */
EditorUi.prototype.executeLayout = function(exec, animate, post)
{
	var graph = this.editor.graph;

	if (graph.isEnabled())
	{
		graph.getModel().beginUpdate();
		try
		{
			exec();
		}
		catch (e)
		{
			throw e;
		}
		finally
		{
			// Animates the changes in the graph model except
			// for Camino, where animation is too slow
			if (this.allowAnimation && animate && navigator.userAgent.indexOf('Camino') < 0)
			{
				// New API for animating graph layout results asynchronously
				var morph = new mxMorphing(graph);
				morph.addListener(mxEvent.DONE, mxUtils.bind(this, function()
				{
					graph.getModel().endUpdate();
					
					if (post != null)
					{
						post();
					}
				}));
				
				morph.startAnimation();
			}
			else
			{
				graph.getModel().endUpdate();
				
				if (post != null)
				{
					post();
				}
			}
		}
	}
};

/**
 * Hides the current menu.
 */
EditorUi.prototype.showImageDialog = function(title, value, fn, ignoreExisting)
{
	var cellEditor = this.editor.graph.cellEditor;
	var selState = cellEditor.saveSelection();
	var newValue = mxUtils.prompt(title, value);
	cellEditor.restoreSelection(selState);
	
	if (newValue != null && newValue.length > 0)
	{
		var img = new Image();
		
		img.onload = function()
		{
			fn(newValue, img.width, img.height);
		};
		img.onerror = function()
		{
			fn(null);
			mxUtils.alert(mxResources.get('fileNotFound'));
		};
		
		img.src = newValue;
	}
	else
	{
		fn(null);
	}
};

/**
 * Hides the current menu.
 */
EditorUi.prototype.showLinkDialog = function(value, btnLabel, fn)
{
	var dlg = new LinkDialog(this, value, btnLabel, fn);
	this.showDialog(dlg.container, 430, 100, true, true);
	dlg.init();
};

/**
 * Hides the current menu.
 */
EditorUi.prototype.showDataDialog = function(cell)
{
	if (cell != null)
	{
		var dlg = new EditDataDialog(this, cell);
		this.showDialog(dlg.container, 350, 340, true, false, null, false);
		dlg.init();
	}
};

/**
 * Hides the current menu.
 */
EditorUi.prototype.showBackgroundImageDialog = function(apply)
{
	apply = (apply != null) ? apply : mxUtils.bind(this, function(image)
	{
		var change = new ChangePageSetup(this, null, image);
		change.ignoreColor = true;
		
		this.editor.graph.model.execute(change);
	});
	
	var newValue = mxUtils.prompt(mxResources.get('backgroundImage'), '');
	
	if (newValue != null && newValue.length > 0)
	{
		var img = new Image();
		
		img.onload = function()
		{
			apply(new mxImage(newValue, img.width, img.height));
		};
		img.onerror = function()
		{
			apply(null);
			mxUtils.alert(mxResources.get('fileNotFound'));
		};
		
		img.src = newValue;
	}
	else
	{
		apply(null);
	}
};

/**
 * Loads the stylesheet for this graph.
 */
EditorUi.prototype.setBackgroundImage = function(image)
{
	this.editor.graph.setBackgroundImage(image);
	this.editor.graph.view.validateBackgroundImage();

	this.fireEvent(new mxEventObject('backgroundImageChanged'));
};

/**
 * Creates the keyboard event handler for the current graph and history.
 */
EditorUi.prototype.confirm = function(msg, okFn, cancelFn)
{
	if (mxUtils.confirm(msg))
	{
		if (okFn != null)
		{
			okFn();
		}
	}
	else if (cancelFn != null)
	{
		cancelFn();
	}
};

/**
 * Creates the keyboard event handler for the current graph and history.
 */
EditorUi.prototype.createOutline = function(wnd)
{
	var outline = new mxOutline(this.editor.graph);
	outline.border = 20;

	mxEvent.addListener(window, 'resize', function()
	{
		outline.update();
	});
	
	this.addListener('pageFormatChanged', function()
	{
		outline.update();
	});

	return outline;
};

// Alt+Shift+Keycode mapping to action
EditorUi.prototype.altShiftActions = {67: 'clearWaypoints', // Alt+Shift+C
  65: 'connectionArrows', // Alt+Shift+A
  76: 'editLink', // Alt+Shift+L
  80: 'connectionPoints', // Alt+Shift+P
  84: 'editTooltip', // Alt+Shift+T
  86: 'pasteSize', // Alt+Shift+V
  88: 'copySize' // Alt+Shift+X
};

/**
 * Creates the keyboard event handler for the current graph and history.
 */
EditorUi.prototype.createKeyHandler = function(editor)
{
	var editorUi = this;
	var graph = this.editor.graph;
	var keyHandler = new mxKeyHandler(graph);

	var isEventIgnored = keyHandler.isEventIgnored;
	keyHandler.isEventIgnored = function(evt)
	{
		// Handles undo/redo/ctrl+./,/u via action and allows ctrl+b/i only if editing value is HTML (except for FF and Safari)
		return (!this.isControlDown(evt) || mxEvent.isShiftDown(evt) || (evt.keyCode != 90 && evt.keyCode != 89 &&
			evt.keyCode != 188 && evt.keyCode != 190 && evt.keyCode != 85)) && ((evt.keyCode != 66 && evt.keyCode != 73) ||
			!this.isControlDown(evt) || (this.graph.cellEditor.isContentEditing() && !mxClient.IS_FF && !mxClient.IS_SF)) &&
			isEventIgnored.apply(this, arguments);
	};
	
	// Ignores graph enabled state but not chromeless state
	keyHandler.isEnabledForEvent = function(evt)
	{
		return (!mxEvent.isConsumed(evt) && this.isGraphEvent(evt) && this.isEnabled() &&
			(editorUi.dialogs == null || editorUi.dialogs.length == 0));
	};
	
	// Routes command-key to control-key on Mac
	keyHandler.isControlDown = function(evt)
	{
		return mxEvent.isControlDown(evt) || (mxClient.IS_MAC && evt.metaKey);
	};

	var queue = [];
	var thread = null;
	
	// Helper function to move cells with the cursor keys
	function nudge(keyCode, stepSize, resize)
	{
		queue.push(function()
		{
			if (!graph.isSelectionEmpty() && graph.isEnabled())
			{
				stepSize = (stepSize != null) ? stepSize : 1;
	
				if (resize)
				{
					// Resizes all selected vertices
					graph.getModel().beginUpdate();
					try
					{
						var cells = graph.getSelectionCells();
						
						for (var i = 0; i < cells.length; i++)
						{
							if (graph.getModel().isVertex(cells[i]) && graph.isCellResizable(cells[i]))
							{
								var geo = graph.getCellGeometry(cells[i]);
								
								if (geo != null)
								{
									geo = geo.clone();
									
									if (keyCode == 37)
									{
										geo.width = Math.max(0, geo.width - stepSize);
									}
									else if (keyCode == 38)
									{
										geo.height = Math.max(0, geo.height - stepSize);
									}
									else if (keyCode == 39)
									{
										geo.width += stepSize;
									}
									else if (keyCode == 40)
									{
										geo.height += stepSize;
									}
									
									graph.getModel().setGeometry(cells[i], geo);
								}
							}
						}
					}
					finally
					{
						graph.getModel().endUpdate();
					}
				}
				else
				{
					// Moves vertices up/down in a stack layout
					var cell = graph.getSelectionCell();
					var parent = graph.model.getParent(cell);
					var layout = null;
	
					if (graph.getSelectionCount() == 1 && graph.model.isVertex(cell) &&
						graph.layoutManager != null && !graph.isCellLocked(cell))
					{
						layout = graph.layoutManager.getLayout(parent);
					}
					
					if (layout != null && layout.constructor == mxStackLayout)
					{
						var index = parent.getIndex(cell);
						
						if (keyCode == 37 || keyCode == 38)
						{
							graph.model.add(parent, cell, Math.max(0, index - 1));
						}
						else if (keyCode == 39 ||keyCode == 40)
						{
							graph.model.add(parent, cell, Math.min(graph.model.getChildCount(parent), index + 1));
						}
					}
					else
					{
						var dx = 0;
						var dy = 0;
						
						if (keyCode == 37)
						{
							dx = -stepSize;
						}
						else if (keyCode == 38)
						{
							dy = -stepSize;
						}
						else if (keyCode == 39)
						{
							dx = stepSize;
						}
						else if (keyCode == 40)
						{
							dy = stepSize;
						}
						
						graph.moveCells(graph.getMovableCells(graph.getSelectionCells()), dx, dy);
					}				
				}
			}
		});
		
		if (thread != null)
		{
			window.clearTimeout(thread);
		}
		
		thread = window.setTimeout(function()
		{
			if (queue.length > 0)
			{
				graph.getModel().beginUpdate();
				try
				{
					for (var i = 0; i < queue.length; i++)
					{
						queue[i]();
					}
					
					queue = [];
				}
				finally
				{
					graph.getModel().endUpdate();
				}
				graph.scrollCellToVisible(graph.getSelectionCell());
			}
		}, 200);
	};
	
	// Overridden to handle special alt+shift+cursor keyboard shortcuts
	var directions = {37: mxConstants.DIRECTION_WEST, 38: mxConstants.DIRECTION_NORTH,
			39: mxConstants.DIRECTION_EAST, 40: mxConstants.DIRECTION_SOUTH};
	
	var keyHandlerGetFunction = keyHandler.getFunction;

	mxKeyHandler.prototype.getFunction = function(evt)
	{
		if (graph.isEnabled())
		{
			// TODO: Add alt modified state in core API, here are some specific cases
			if (mxEvent.isShiftDown(evt) && mxEvent.isAltDown(evt))
			{
				var action = editorUi.actions.get(editorUi.altShiftActions[evt.keyCode]);

				if (action != null)
				{
					return action.funct;
				}
			}
			
			if (evt.keyCode == 9 && mxEvent.isAltDown(evt))
			{
				if (mxEvent.isShiftDown(evt))
				{
					// Alt+Shift+Tab
					return function()
					{
						graph.selectParentCell();
					};
				}
				else
				{
					// Alt+Tab
					return function()
					{
						graph.selectChildCell();
					};
				}
			}
			else if (directions[evt.keyCode] != null && !graph.isSelectionEmpty())
			{
				if (mxEvent.isShiftDown(evt) && mxEvent.isAltDown(evt))
				{
					if (graph.model.isVertex(graph.getSelectionCell()))
					{
						return function()
						{
							var cells = graph.connectVertex(graph.getSelectionCell(), directions[evt.keyCode],
								graph.defaultEdgeLength, evt, true);
			
							if (cells != null && cells.length > 0)
							{
								if (cells.length == 1 && graph.model.isEdge(cells[0]))
								{
									graph.setSelectionCell(graph.model.getTerminal(cells[0], false));
								}
								else
								{
									graph.setSelectionCell(cells[cells.length - 1]);
								}

								graph.scrollCellToVisible(graph.getSelectionCell());
								
								if (editorUi.hoverIcons != null)
								{
									editorUi.hoverIcons.update(graph.view.getState(graph.getSelectionCell()));
								}
							}
						};
					}
				}
				else
				{
					// Avoids consuming event if no vertex is selected by returning null below
					// Cursor keys move and resize (ctrl) cells
					if (this.isControlDown(evt))
					{
						return function()
						{
							nudge(evt.keyCode, (mxEvent.isShiftDown(evt)) ? graph.gridSize : null, true);
						};
					}
					else
					{
						return function()
						{
							nudge(evt.keyCode, (mxEvent.isShiftDown(evt)) ? graph.gridSize : null);
						};
					}
				}
			}
		}

		return keyHandlerGetFunction.apply(this, arguments);
	};

	// Binds keystrokes to actions
	keyHandler.bindAction = mxUtils.bind(this, function(code, control, key, shift)
	{
		var action = this.actions.get(key);
		
		if (action != null)
		{
			var f = function()
			{
				if (action.isEnabled())
				{
					action.funct();
				}
			};
    		
			if (control)
			{
				if (shift)
				{
					keyHandler.bindControlShiftKey(code, f);
				}
				else
				{
					keyHandler.bindControlKey(code, f);
				}
			}
			else
			{
				if (shift)
				{
					keyHandler.bindShiftKey(code, f);
				}
				else
				{
					keyHandler.bindKey(code, f);
				}
			}
		}
	});

	var ui = this;
	var keyHandlerEscape = keyHandler.escape;
	keyHandler.escape = function(evt)
	{
		keyHandlerEscape.apply(this, arguments);
	};

	// Ignores enter keystroke. Remove this line if you want the
	// enter keystroke to stop editing. N, W, T are reserved.
	keyHandler.enter = function() {};
	
	keyHandler.bindControlShiftKey(36, function() { graph.exitGroup(); }); // Ctrl+Shift+Home
	keyHandler.bindControlShiftKey(35, function() { graph.enterGroup(); }); // Ctrl+Shift+End
	keyHandler.bindKey(36, function() { graph.home(); }); // Home
	keyHandler.bindKey(35, function() { graph.refresh(); }); // End
	keyHandler.bindAction(107, true, 'zoomIn'); // Ctrl+Plus
	keyHandler.bindAction(109, true, 'zoomOut'); // Ctrl+Minus
	keyHandler.bindAction(80, true, 'print'); // Ctrl+P
	keyHandler.bindAction(79, true, 'outline', true); // Ctrl+Shift+O
	keyHandler.bindAction(112, false, 'about'); // F1

	if (!this.editor.chromeless || this.editor.editable)
	{
		keyHandler.bindControlKey(36, function() { if (graph.isEnabled()) { graph.foldCells(true); }}); // Ctrl+Home
		keyHandler.bindControlKey(35, function() { if (graph.isEnabled()) { graph.foldCells(false); }}); // Ctrl+End
		keyHandler.bindControlKey(13, function() { if (graph.isEnabled()) { graph.setSelectionCells(graph.duplicateCells(graph.getSelectionCells(), false)); }}); // Ctrl+Enter
		keyHandler.bindAction(8, false, 'delete'); // Backspace
		keyHandler.bindAction(8, true, 'deleteAll'); // Backspace
		keyHandler.bindAction(46, false, 'delete'); // Delete
		keyHandler.bindAction(46, true, 'deleteAll'); // Ctrl+Delete
		keyHandler.bindAction(72, true, 'resetView'); // Ctrl+H
		keyHandler.bindAction(72, true, 'fitWindow', true); // Ctrl+Shift+H
		keyHandler.bindAction(74, true, 'fitPage'); // Ctrl+J
		keyHandler.bindAction(74, true, 'fitTwoPages', true); // Ctrl+Shift+J
		keyHandler.bindAction(48, true, 'customZoom'); // Ctrl+0
		keyHandler.bindAction(82, true, 'turn'); // Ctrl+R
		keyHandler.bindAction(82, true, 'clearDefaultStyle', true); // Ctrl+Shift+R
		keyHandler.bindAction(83, true, 'save'); // Ctrl+S
		keyHandler.bindAction(83, true, 'saveAs', true); // Ctrl+Shift+S
		keyHandler.bindAction(65, true, 'selectAll'); // Ctrl+A
		keyHandler.bindAction(65, true, 'selectNone', true); // Ctrl+A
		keyHandler.bindAction(73, true, 'selectVertices', true); // Ctrl+Shift+I
		keyHandler.bindAction(69, true, 'selectEdges', true); // Ctrl+Shift+E
		keyHandler.bindAction(69, true, 'editStyle'); // Ctrl+E
		keyHandler.bindAction(66, true, 'bold'); // Ctrl+B
		keyHandler.bindAction(66, true, 'toBack', true); // Ctrl+Shift+B
		keyHandler.bindAction(70, true, 'toFront', true); // Ctrl+Shift+F
		keyHandler.bindAction(68, true, 'duplicate'); // Ctrl+D
		keyHandler.bindAction(68, true, 'setAsDefaultStyle', true); // Ctrl+Shift+D   
		keyHandler.bindAction(90, true, 'undo'); // Ctrl+Z
		keyHandler.bindAction(89, true, 'autosize', true); // Ctrl+Shift+Y
		keyHandler.bindAction(88, true, 'cut'); // Ctrl+X
		keyHandler.bindAction(67, true, 'copy'); // Ctrl+C
		keyHandler.bindAction(86, true, 'paste'); // Ctrl+V
		keyHandler.bindAction(71, true, 'group'); // Ctrl+G
		keyHandler.bindAction(77, true, 'editData'); // Ctrl+M
		keyHandler.bindAction(71, true, 'grid', true); // Ctrl+Shift+G
		keyHandler.bindAction(73, true, 'italic'); // Ctrl+I
		keyHandler.bindAction(76, true, 'lockUnlock'); // Ctrl+L
		keyHandler.bindAction(76, true, 'layers', true); // Ctrl+Shift+L
		keyHandler.bindAction(80, true, 'formatPanel', true); // Ctrl+Shift+P
		keyHandler.bindAction(85, true, 'underline'); // Ctrl+U
		keyHandler.bindAction(85, true, 'ungroup', true); // Ctrl+Shift+U
		keyHandler.bindAction(190, true, 'superscript'); // Ctrl+.
		keyHandler.bindAction(188, true, 'subscript'); // Ctrl+,
		keyHandler.bindKey(13, function() { if (graph.isEnabled()) { graph.startEditingAtCell(); }}); // Enter
		keyHandler.bindKey(113, function() { if (graph.isEnabled()) { graph.startEditingAtCell(); }}); // F2
	}
	
	if (!mxClient.IS_WIN)
	{
		keyHandler.bindAction(90, true, 'redo', true); // Ctrl+Shift+Z
	}
	else
	{
		keyHandler.bindAction(89, true, 'redo'); // Ctrl+Y
	}
	
	return keyHandler;
};

/**
 * Creates the keyboard event handler for the current graph and history.
 */
EditorUi.prototype.destroy = function()
{
	if (this.editor != null)
	{
		this.editor.destroy();
		this.editor = null;
	}
	
	if (this.menubar != null)
	{
		this.menubar.destroy();
		this.menubar = null;
	}
	
	if (this.toolbar != null)
	{
		this.toolbar.destroy();
		this.toolbar = null;
	}
	
	if (this.sidebar != null)
	{
		this.sidebar.destroy();
		this.sidebar = null;
	}
	
	if (this.keyHandler != null)
	{
		this.keyHandler.destroy();
		this.keyHandler = null;
	}
	
	if (this.keydownHandler != null)
	{
		mxEvent.removeListener(document, 'keydown', this.keydownHandler);
		this.keydownHandler = null;
	}
		
	if (this.keyupHandler != null)
	{
		mxEvent.removeListener(document, 'keyup', this.keyupHandler);
		this.keyupHandler = null;
	}
	
	if (this.resizeHandler != null)
	{
		mxEvent.removeListener(window, 'resize', this.resizeHandler);
		this.resizeHandler = null;
	}
	
	if (this.gestureHandler != null)
	{
		mxEvent.removeGestureListeners(document, this.gestureHandler);
		this.gestureHandler = null;
	}
	
	if (this.orientationChangeHandler != null)
	{
		mxEvent.removeListener(window, 'orientationchange', this.orientationChangeHandler);
		this.orientationChangeHandler = null;
	}
	
	if (this.scrollHandler != null)
	{
		mxEvent.removeListener(window, 'scroll', this.scrollHandler);
		this.scrollHandler = null;
	}

	if (this.destroyFunctions != null)
	{
		for (var i = 0; i < this.destroyFunctions.length; i++)
		{
			this.destroyFunctions[i]();
		}
		
		this.destroyFunctions = null;
	}
	
	var c = [ this.toolbarContainer, this.sidebarContainer,
	         this.formatContainer, this.diagramContainer, this.footerContainer,
	         this.chromelessToolbar, this.hsplit, this.sidebarFooterContainer,
	         this.layersDialog];
	
	for (var i = 0; i < c.length; i++)
	{
		if (c[i] != null && c[i].parentNode != null)
		{
			c[i].parentNode.removeChild(c[i]);
		}
	}
};
