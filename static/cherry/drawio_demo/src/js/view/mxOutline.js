/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxOutline
 *
 * Implements an outline (aka overview) for a graph. Set <updateOnPan> to true
 * to enable updates while the source graph is panning.
 * 
 * Example:
 * 
 * (code)
 * var outline = new mxOutline(graph, div);
 * (end)
 * 
 * If an outline is used in an <mxWindow> in IE8 standards mode, the following
 * code makes sure that the shadow filter is not inherited and that any
 * transparent elements in the graph do not show the page background, but the
 * background of the graph container.
 * 
 * (code)
 * if (document.documentMode == 8)
 * {
 *   container.style.filter = 'progid:DXImageTransform.Microsoft.alpha(opacity=100)';
 * }
 * (end)
 * 
 * To move the graph to the top, left corner the following code can be used.
 * 
 * (code)
 * var scale = graph.view.scale;
 * var bounds = graph.getGraphBounds();
 * graph.view.setTranslate(-bounds.x / scale, -bounds.y / scale);
 * (end)
 * 
 * To toggle the suspended mode, the following can be used.
 * 
 * (code)
 * outline.suspended = !outln.suspended;
 * if (!outline.suspended)
 * {
 *   outline.update(true);
 * }
 * (end)
 * 
 * Constructor: mxOutline
 *
 * Constructs a new outline for the specified graph inside the given
 * container.
 * 
 * Parameters:
 * 
 * source - <mxGraph> to create the outline for.
 * container - DOM node that will contain the outline.
 */
function mxOutline(source, container)
{
	this.source = source;

	if (container != null)
	{
		this.init(container);
	}
};

/**
 * Function: source
 * 
 * Reference to the source <mxGraph>.
 */
mxOutline.prototype.source = null;

/**
 * Function: outline
 * 
 * Reference to the <mxGraph> that renders the outline.
 */
mxOutline.prototype.outline = null;

/**
 * Function: graphRenderHint
 * 
 * Renderhint to be used for the outline graph. Default is faster.
 */
mxOutline.prototype.graphRenderHint = mxConstants.RENDERING_HINT_FASTER;

/**
 * Variable: enabled
 * 
 * Specifies if events are handled. Default is true.
 */
mxOutline.prototype.enabled = true;

/**
 * Variable: showViewport
 * 
 * Specifies a viewport rectangle should be shown. Default is true.
 */
mxOutline.prototype.showViewport = true;

/**
 * Variable: border
 * 
 * Border to be added at the bottom and right. Default is 10.
 */
mxOutline.prototype.border = 10;

/**
 * Variable: enabled
 * 
 * Specifies the size of the sizer handler. Default is 8.
 */
mxOutline.prototype.sizerSize = 8;

/**
 * Variable: labelsVisible
 * 
 * Specifies if labels should be visible in the outline. Default is false.
 */
mxOutline.prototype.labelsVisible = false;

/**
 * Variable: updateOnPan
 * 
 * Specifies if <update> should be called for <mxEvent.PAN> in the source
 * graph. Default is false.
 */
mxOutline.prototype.updateOnPan = false;

/**
 * Variable: sizerImage
 * 
 * Optional <mxImage> to be used for the sizer. Default is null.
 */
mxOutline.prototype.sizerImage = null;

/**
 * Variable: minScale
 * 
 * Minimum scale to be used. Default is 0.001.
 */
mxOutline.prototype.minScale = 0.0001;

/**
 * Variable: suspended
 * 
 * Optional boolean flag to suspend updates. Default is false. This flag will
 * also suspend repaints of the outline. To toggle this switch, use the
 * following code.
 * 
 * (code)
 * nav.suspended = !nav.suspended;
 * 
 * if (!nav.suspended)
 * {
 *   nav.update(true);
 * }
 * (end)
 */
mxOutline.prototype.suspended = false;

/**
 * Variable: forceVmlHandles
 * 
 * Specifies if VML should be used to render the handles in this control. This
 * is true for IE8 standards mode and false for all other browsers and modes.
 * This is a workaround for rendering issues of HTML elements over elements
 * with filters in IE 8 standards mode.
 */
mxOutline.prototype.forceVmlHandles = document.documentMode == 8;

/**
 * Function: createGraph
 * 
 * Creates the <mxGraph> used in the outline.
 */
mxOutline.prototype.createGraph = function(container)
{
	var graph = new mxGraph(container, this.source.getModel(), this.graphRenderHint, this.source.getStylesheet());
	graph.foldingEnabled = false;
	graph.autoScroll = false;
	
	return graph;
};

/**
 * Function: init
 * 
 * Initializes the outline inside the given container.
 */
mxOutline.prototype.init = function(container)
{
	this.outline = this.createGraph(container);
	
	// Do not repaint when suspended
	var outlineGraphModelChanged = this.outline.graphModelChanged;
	this.outline.graphModelChanged = mxUtils.bind(this, function(changes)
	{
		if (!this.suspended && this.outline != null)
		{
			outlineGraphModelChanged.apply(this.outline, arguments);
		}
	});

	// Enables faster painting in SVG
	if (mxClient.IS_SVG)
	{
		var node = this.outline.getView().getCanvas().parentNode;
		node.setAttribute('shape-rendering', 'optimizeSpeed');
		node.setAttribute('image-rendering', 'optimizeSpeed');
	}
	
	// Hides cursors and labels
	this.outline.labelsVisible = this.labelsVisible;
	this.outline.setEnabled(false);
	
	this.updateHandler = mxUtils.bind(this, function(sender, evt)
	{
		if (!this.suspended && !this.active)
		{
			this.update();
		}
	});
	
	// Updates the scale of the outline after a change of the main graph
	this.source.getModel().addListener(mxEvent.CHANGE, this.updateHandler);
	this.outline.addMouseListener(this);
	
	// Adds listeners to keep the outline in sync with the source graph
	var view = this.source.getView();
	view.addListener(mxEvent.SCALE, this.updateHandler);
	view.addListener(mxEvent.TRANSLATE, this.updateHandler);
	view.addListener(mxEvent.SCALE_AND_TRANSLATE, this.updateHandler);
	view.addListener(mxEvent.DOWN, this.updateHandler);
	view.addListener(mxEvent.UP, this.updateHandler);

	// Updates blue rectangle on scroll
	mxEvent.addListener(this.source.container, 'scroll', this.updateHandler);
	
	this.panHandler = mxUtils.bind(this, function(sender)
	{
		if (this.updateOnPan)
		{
			this.updateHandler.apply(this, arguments);
		}
	});
	this.source.addListener(mxEvent.PAN, this.panHandler);
	
	// Refreshes the graph in the outline after a refresh of the main graph
	this.refreshHandler = mxUtils.bind(this, function(sender)
	{
		this.outline.setStylesheet(this.source.getStylesheet());
		this.outline.refresh();
	});
	this.source.addListener(mxEvent.REFRESH, this.refreshHandler);

	// Creates the blue rectangle for the viewport
	this.bounds = new mxRectangle(0, 0, 0, 0);
	this.selectionBorder = new mxRectangleShape(this.bounds, null,
		mxConstants.OUTLINE_COLOR, mxConstants.OUTLINE_STROKEWIDTH);
	this.selectionBorder.dialect = this.outline.dialect;

	if (this.forceVmlHandles)
	{
		this.selectionBorder.isHtmlAllowed = function()
		{
			return false;
		};
	}
	
	this.selectionBorder.init(this.outline.getView().getOverlayPane());

	// Handles event by catching the initial pointer start and then listening to the
	// complete gesture on the event target. This is needed because all the events
	// are routed via the initial element even if that element is removed from the
	// DOM, which happens when we repaint the selection border and zoom handles.
	var handler = mxUtils.bind(this, function(evt)
	{
		var t = mxEvent.getSource(evt);
		
		var redirect = mxUtils.bind(this, function(evt)
		{
			this.outline.fireMouseEvent(mxEvent.MOUSE_MOVE, new mxMouseEvent(evt));
		});
		
		var redirect2 = mxUtils.bind(this, function(evt)
		{
			mxEvent.removeGestureListeners(t, null, redirect, redirect2);
			this.outline.fireMouseEvent(mxEvent.MOUSE_UP, new mxMouseEvent(evt));
		});
		
		mxEvent.addGestureListeners(t, null, redirect, redirect2);
		this.outline.fireMouseEvent(mxEvent.MOUSE_DOWN, new mxMouseEvent(evt));
	});
	
	mxEvent.addGestureListeners(this.selectionBorder.node, handler);

	// Creates a small blue rectangle for sizing (sizer handle)
	this.sizer = this.createSizer();
	
	if (this.forceVmlHandles)
	{
		this.sizer.isHtmlAllowed = function()
		{
			return false;
		};
	}
	
	this.sizer.init(this.outline.getView().getOverlayPane());
	
	if (this.enabled)
	{
		this.sizer.node.style.cursor = 'nwse-resize';
	}
	
	mxEvent.addGestureListeners(this.sizer.node, handler);

	this.selectionBorder.node.style.display = (this.showViewport) ? '' : 'none';
	this.sizer.node.style.display = this.selectionBorder.node.style.display;
	this.selectionBorder.node.style.cursor = 'move';

	this.update(false);
};

/**
 * Function: isEnabled
 * 
 * Returns true if events are handled. This implementation
 * returns <enabled>.
 */
mxOutline.prototype.isEnabled = function()
{
	return this.enabled;
};

/**
 * Function: setEnabled
 * 
 * Enables or disables event handling. This implementation
 * updates <enabled>.
 * 
 * Parameters:
 * 
 * value - Boolean that specifies the new enabled state.
 */
mxOutline.prototype.setEnabled = function(value)
{
	this.enabled = value;
};

/**
 * Function: setZoomEnabled
 * 
 * Enables or disables the zoom handling by showing or hiding the respective
 * handle.
 * 
 * Parameters:
 * 
 * value - Boolean that specifies the new enabled state.
 */
mxOutline.prototype.setZoomEnabled = function(value)
{
	this.sizer.node.style.visibility = (value) ? 'visible' : 'hidden';
};

/**
 * Function: refresh
 * 
 * Invokes <update> and revalidate the outline. This method is deprecated.
 */
mxOutline.prototype.refresh = function()
{
	this.update(true);
};

/**
 * Function: createSizer
 * 
 * Creates the shape used as the sizer.
 */
mxOutline.prototype.createSizer = function()
{
	if (this.sizerImage != null)
	{
		var sizer = new mxImageShape(new mxRectangle(0, 0, this.sizerImage.width, this.sizerImage.height), this.sizerImage.src);
		sizer.dialect = this.outline.dialect;
		
		return sizer;
	}
	else
	{
		var sizer = new mxRectangleShape(new mxRectangle(0, 0, this.sizerSize, this.sizerSize),
			mxConstants.OUTLINE_HANDLE_FILLCOLOR, mxConstants.OUTLINE_HANDLE_STROKECOLOR);
		sizer.dialect = this.outline.dialect;
	
		return sizer;
	}
};

/**
 * Function: getSourceContainerSize
 * 
 * Returns the size of the source container.
 */
mxOutline.prototype.getSourceContainerSize = function()
{
	return new mxRectangle(0, 0, this.source.container.scrollWidth, this.source.container.scrollHeight);
};

/**
 * Function: getOutlineOffset
 * 
 * Returns the offset for drawing the outline graph.
 */
mxOutline.prototype.getOutlineOffset = function(scale)
{
	return null;
};

/**
 * Function: getOutlineOffset
 * 
 * Returns the offset for drawing the outline graph.
 */
mxOutline.prototype.getSourceGraphBounds = function()
{
	return this.source.getGraphBounds();
};

/**
 * Function: update
 * 
 * Updates the outline.
 */
mxOutline.prototype.update = function(revalidate)
{
	if (this.source != null && this.source.container != null &&
		this.outline != null && this.outline.container != null)
	{
		var sourceScale = this.source.view.scale;
		var scaledGraphBounds = this.getSourceGraphBounds();
		var unscaledGraphBounds = new mxRectangle(scaledGraphBounds.x / sourceScale + this.source.panDx,
				scaledGraphBounds.y / sourceScale + this.source.panDy, scaledGraphBounds.width / sourceScale,
				scaledGraphBounds.height / sourceScale);

		var unscaledFinderBounds = new mxRectangle(0, 0,
			this.source.container.clientWidth / sourceScale,
			this.source.container.clientHeight / sourceScale);
		
		var union = unscaledGraphBounds.clone();
		union.add(unscaledFinderBounds);
	
		// Zooms to the scrollable area if that is bigger than the graph
		var size = this.getSourceContainerSize();
		var completeWidth = Math.max(size.width / sourceScale, union.width);
		var completeHeight = Math.max(size.height / sourceScale, union.height);
	
		var availableWidth = Math.max(0, this.outline.container.clientWidth - this.border);
		var availableHeight = Math.max(0, this.outline.container.clientHeight - this.border);
		
		var outlineScale = Math.min(availableWidth / completeWidth, availableHeight / completeHeight);
		var scale = (isNaN(outlineScale)) ? this.minScale : Math.max(this.minScale, outlineScale);

		if (scale > 0)
		{
			if (this.outline.getView().scale != scale)
			{
				this.outline.getView().scale = scale;
				revalidate = true;
			}
		
			var navView = this.outline.getView();
			
			if (navView.currentRoot != this.source.getView().currentRoot)
			{
				navView.setCurrentRoot(this.source.getView().currentRoot);
			}

			var t = this.source.view.translate;
			var tx = t.x + this.source.panDx;
			var ty = t.y + this.source.panDy;
			
			var off = this.getOutlineOffset(scale);
			
			if (off != null)
			{
				tx += off.x;
				ty += off.y;
			}
			
			if (unscaledGraphBounds.x < 0)
			{
				tx = tx - unscaledGraphBounds.x;
			}
			if (unscaledGraphBounds.y < 0)
			{
				ty = ty - unscaledGraphBounds.y;
			}
			
			if (navView.translate.x != tx || navView.translate.y != ty)
			{
				navView.translate.x = tx;
				navView.translate.y = ty;
				revalidate = true;
			}
		
			// Prepares local variables for computations
			var t2 = navView.translate;
			scale = this.source.getView().scale;
			var scale2 = scale / navView.scale;
			var scale3 = 1.0 / navView.scale;
			var container = this.source.container;
			
			// Updates the bounds of the viewrect in the navigation
			this.bounds = new mxRectangle(
				(t2.x - t.x - this.source.panDx) / scale3,
				(t2.y - t.y - this.source.panDy) / scale3,
				(container.clientWidth / scale2),
				(container.clientHeight / scale2));
			
			// Adds the scrollbar offset to the finder
			this.bounds.x += this.source.container.scrollLeft * navView.scale / scale;
			this.bounds.y += this.source.container.scrollTop * navView.scale / scale;
			
			var b = this.selectionBorder.bounds;
			
			if (b.x != this.bounds.x || b.y != this.bounds.y || b.width != this.bounds.width || b.height != this.bounds.height)
			{
				this.selectionBorder.bounds = this.bounds;
				this.selectionBorder.redraw();
			}
		
			// Updates the bounds of the zoom handle at the bottom right
			var b = this.sizer.bounds;
			var b2 = new mxRectangle(this.bounds.x + this.bounds.width - b.width / 2,
					this.bounds.y + this.bounds.height - b.height / 2, b.width, b.height);

			if (b.x != b2.x || b.y != b2.y || b.width != b2.width || b.height != b2.height)
			{
				this.sizer.bounds = b2;
				
				// Avoids update of visibility in redraw for VML
				if (this.sizer.node.style.visibility != 'hidden')
				{
					this.sizer.redraw();
				}
			}

			if (revalidate)
			{
				this.outline.view.revalidate();
			}
		}
	}
};

/**
 * Function: mouseDown
 * 
 * Handles the event by starting a translation or zoom.
 */
mxOutline.prototype.mouseDown = function(sender, me)
{
	if (this.enabled && this.showViewport)
	{
		var tol = (!mxEvent.isMouseEvent(me.getEvent())) ? this.source.tolerance : 0;
		var hit = (this.source.allowHandleBoundsCheck && (mxClient.IS_IE || tol > 0)) ?
				new mxRectangle(me.getGraphX() - tol, me.getGraphY() - tol, 2 * tol, 2 * tol) : null;
		this.zoom = me.isSource(this.sizer) || (hit != null && mxUtils.intersects(shape.bounds, hit));
		this.startX = me.getX();
		this.startY = me.getY();
		this.active = true;

		if (this.source.useScrollbarsForPanning && mxUtils.hasScrollbars(this.source.container))
		{
			this.dx0 = this.source.container.scrollLeft;
			this.dy0 = this.source.container.scrollTop;
		}
		else
		{
			this.dx0 = 0;
			this.dy0 = 0;
		}
	}

	me.consume();
};

/**
 * Function: mouseMove
 * 
 * Handles the event by previewing the viewrect in <graph> and updating the
 * rectangle that represents the viewrect in the outline.
 */
mxOutline.prototype.mouseMove = function(sender, me)
{
	if (this.active)
	{
		this.selectionBorder.node.style.display = (this.showViewport) ? '' : 'none';
		this.sizer.node.style.display = this.selectionBorder.node.style.display; 

		var delta = this.getTranslateForEvent(me);
		var dx = delta.x;
		var dy = delta.y;
		var bounds = null;
		
		if (!this.zoom)
		{
			// Previews the panning on the source graph
			var scale = this.outline.getView().scale;
			bounds = new mxRectangle(this.bounds.x + dx,
				this.bounds.y + dy, this.bounds.width, this.bounds.height);
			this.selectionBorder.bounds = bounds;
			this.selectionBorder.redraw();
			dx /= scale;
			dx *= this.source.getView().scale;
			dy /= scale;
			dy *= this.source.getView().scale;
			this.source.panGraph(-dx - this.dx0, -dy - this.dy0);
		}
		else
		{
			// Does *not* preview zooming on the source graph
			var container = this.source.container;
			var viewRatio = container.clientWidth / container.clientHeight;
			dy = dx / viewRatio;
			bounds = new mxRectangle(this.bounds.x,
				this.bounds.y,
				Math.max(1, this.bounds.width + dx),
				Math.max(1, this.bounds.height + dy));
			this.selectionBorder.bounds = bounds;
			this.selectionBorder.redraw();
		}
		
		// Updates the zoom handle
		var b = this.sizer.bounds;
		this.sizer.bounds = new mxRectangle(
			bounds.x + bounds.width - b.width / 2,
			bounds.y + bounds.height - b.height / 2,
			b.width, b.height);
		
		// Avoids update of visibility in redraw for VML
		if (this.sizer.node.style.visibility != 'hidden')
		{
			this.sizer.redraw();
		}
		
		me.consume();
	}
};

/**
 * Function: getTranslateForEvent
 * 
 * Gets the translate for the given mouse event. Here is an example to limit
 * the outline to stay within positive coordinates:
 * 
 * (code)
 * outline.getTranslateForEvent = function(me)
 * {
 *   var pt = new mxPoint(me.getX() - this.startX, me.getY() - this.startY);
 *   
 *   if (!this.zoom)
 *   {
 *     var tr = this.source.view.translate;
 *     pt.x = Math.max(tr.x * this.outline.view.scale, pt.x);
 *     pt.y = Math.max(tr.y * this.outline.view.scale, pt.y);
 *   }
 *   
 *   return pt;
 * };
 * (end)
 */
mxOutline.prototype.getTranslateForEvent = function(me)
{
	return new mxPoint(me.getX() - this.startX, me.getY() - this.startY);
};

/**
 * Function: mouseUp
 * 
 * Handles the event by applying the translation or zoom to <graph>.
 */
mxOutline.prototype.mouseUp = function(sender, me)
{
	if (this.active)
	{
		var delta = this.getTranslateForEvent(me);
		var dx = delta.x;
		var dy = delta.y;
		
		if (Math.abs(dx) > 0 || Math.abs(dy) > 0)
		{
			if (!this.zoom)
			{
				// Applies the new translation if the source
				// has no scrollbars
				if (!this.source.useScrollbarsForPanning ||
					!mxUtils.hasScrollbars(this.source.container))
				{
					this.source.panGraph(0, 0);
					dx /= this.outline.getView().scale;
					dy /= this.outline.getView().scale;
					var t = this.source.getView().translate;
					this.source.getView().setTranslate(t.x - dx, t.y - dy);
				}
			}
			else
			{
				// Applies the new zoom
				var w = this.selectionBorder.bounds.width;
				var scale = this.source.getView().scale;
				this.source.zoomTo(Math.max(this.minScale, scale - (dx * scale) / w), false);
			}

			this.update();
			me.consume();
		}
			
		// Resets the state of the handler
		this.index = null;
		this.active = false;
	}
};

/**
 * Function: destroy
 * 
 * Destroy this outline and removes all listeners from <source>.
 */
mxOutline.prototype.destroy = function()
{
	if (this.source != null)
	{
		this.source.removeListener(this.panHandler);
		this.source.removeListener(this.refreshHandler);
		this.source.getModel().removeListener(this.updateHandler);
		this.source.getView().removeListener(this.updateHandler);
		mxEvent.addListener(this.source.container, 'scroll', this.updateHandler);
		this.source = null;
	}
	
	if (this.outline != null)
	{
		this.outline.removeMouseListener(this);
		this.outline.destroy();
		this.outline = null;
	}

	if (this.selectionBorder != null)
	{
		this.selectionBorder.destroy();
		this.selectionBorder = null;
	}
	
	if (this.sizer != null)
	{
		this.sizer.destroy();
		this.sizer = null;
	}
};
