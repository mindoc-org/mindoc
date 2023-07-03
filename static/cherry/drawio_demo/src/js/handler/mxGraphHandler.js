/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxGraphHandler
 * 
 * Graph event handler that handles selection. Individual cells are handled
 * separately using <mxVertexHandler> or one of the edge handlers. These
 * handlers are created using <mxGraph.createHandler> in
 * <mxGraphSelectionModel.cellAdded>.
 * 
 * To avoid the container to scroll a moved cell into view, set
 * <scrollAfterMove> to false.
 * 
 * Constructor: mxGraphHandler
 * 
 * Constructs an event handler that creates handles for the
 * selection cells.
 * 
 * Parameters:
 * 
 * graph - Reference to the enclosing <mxGraph>.
 */
function mxGraphHandler(graph)
{
	this.graph = graph;
	this.graph.addMouseListener(this);
	
	// Repaints the handler after autoscroll
	this.panHandler = mxUtils.bind(this, function()
	{
		this.updatePreviewShape();
		this.updateHint();
	});
	
	this.graph.addListener(mxEvent.PAN, this.panHandler);
	
	// Handles escape keystrokes
	this.escapeHandler = mxUtils.bind(this, function(sender, evt)
	{
		this.reset();
	});
	
	this.graph.addListener(mxEvent.ESCAPE, this.escapeHandler);
	
	// Updates the preview box for remote changes
	this.refreshHandler = mxUtils.bind(this, function(sender, evt)
	{
		if (this.first != null)
		{
			try
			{
				this.bounds = this.graph.getView().getBounds(this.cells);
				this.pBounds = this.getPreviewBounds(this.cells);
				this.updatePreviewShape();
			}
			catch (e)
			{
				// Resets the handler if cells have vanished
				this.reset();
			}
		}
	});
	
	this.graph.getModel().addListener(mxEvent.CHANGE, this.refreshHandler);
};

/**
 * Variable: graph
 * 
 * Reference to the enclosing <mxGraph>.
 */
mxGraphHandler.prototype.graph = null;

/**
 * Variable: maxCells
 * 
 * Defines the maximum number of cells to paint subhandles
 * for. Default is 50 for Firefox and 20 for IE. Set this
 * to 0 if you want an unlimited number of handles to be
 * displayed. This is only recommended if the number of
 * cells in the graph is limited to a small number, eg.
 * 500.
 */
mxGraphHandler.prototype.maxCells = (mxClient.IS_IE) ? 20 : 50;

/**
 * Variable: enabled
 * 
 * Specifies if events are handled. Default is true.
 */
mxGraphHandler.prototype.enabled = true;

/**
 * Variable: highlightEnabled
 * 
 * Specifies if drop targets under the mouse should be enabled. Default is
 * true.
 */
mxGraphHandler.prototype.highlightEnabled = true;

/**
 * Variable: cloneEnabled
 * 
 * Specifies if cloning by control-drag is enabled. Default is true.
 */
mxGraphHandler.prototype.cloneEnabled = true;

/**
 * Variable: moveEnabled
 * 
 * Specifies if moving is enabled. Default is true.
 */
mxGraphHandler.prototype.moveEnabled = true;

/**
 * Variable: guidesEnabled
 * 
 * Specifies if other cells should be used for snapping the right, center or
 * left side of the current selection. Default is false.
 */
mxGraphHandler.prototype.guidesEnabled = false;

/**
 * Variable: guide
 * 
 * Holds the <mxGuide> instance that is used for alignment.
 */
mxGraphHandler.prototype.guide = null;

/**
 * Variable: currentDx
 * 
 * Stores the x-coordinate of the current mouse move.
 */
mxGraphHandler.prototype.currentDx = null;

/**
 * Variable: currentDy
 * 
 * Stores the y-coordinate of the current mouse move.
 */
mxGraphHandler.prototype.currentDy = null;

/**
 * Variable: updateCursor
 * 
 * Specifies if a move cursor should be shown if the mouse is over a movable
 * cell. Default is true.
 */
mxGraphHandler.prototype.updateCursor = true;

/**
 * Variable: selectEnabled
 * 
 * Specifies if selecting is enabled. Default is true.
 */
mxGraphHandler.prototype.selectEnabled = true;

/**
 * Variable: removeCellsFromParent
 * 
 * Specifies if cells may be moved out of their parents. Default is true.
 */
mxGraphHandler.prototype.removeCellsFromParent = true;

/**
 * Variable: connectOnDrop
 * 
 * Specifies if drop events are interpreted as new connections if no other
 * drop action is defined. Default is false.
 */
mxGraphHandler.prototype.connectOnDrop = false;

/**
 * Variable: scrollOnMove
 * 
 * Specifies if the view should be scrolled so that a moved cell is
 * visible. Default is true.
 */
mxGraphHandler.prototype.scrollOnMove = true;

/**
 * Variable: minimumSize
 * 
 * Specifies the minimum number of pixels for the width and height of a
 * selection border. Default is 6.
 */
mxGraphHandler.prototype.minimumSize = 6;

/**
 * Variable: previewColor
 * 
 * Specifies the color of the preview shape. Default is black.
 */
mxGraphHandler.prototype.previewColor = 'black';

/**
 * Variable: htmlPreview
 * 
 * Specifies if the graph container should be used for preview. If this is used
 * then drop target detection relies entirely on <mxGraph.getCellAt> because
 * the HTML preview does not "let events through". Default is false.
 */
mxGraphHandler.prototype.htmlPreview = false;

/**
 * Variable: shape
 * 
 * Reference to the <mxShape> that represents the preview.
 */
mxGraphHandler.prototype.shape = null;

/**
 * Variable: scaleGrid
 * 
 * Specifies if the grid should be scaled. Default is false.
 */
mxGraphHandler.prototype.scaleGrid = false;

/**
 * Variable: rotationEnabled
 * 
 * Specifies if the bounding box should allow for rotation. Default is true.
 */
mxGraphHandler.prototype.rotationEnabled = true;

/**
 * Function: isEnabled
 * 
 * Returns <enabled>.
 */
mxGraphHandler.prototype.isEnabled = function()
{
	return this.enabled;
};

/**
 * Function: setEnabled
 * 
 * Sets <enabled>.
 */
mxGraphHandler.prototype.setEnabled = function(value)
{
	this.enabled = value;
};

/**
 * Function: isCloneEnabled
 * 
 * Returns <cloneEnabled>.
 */
mxGraphHandler.prototype.isCloneEnabled = function()
{
	return this.cloneEnabled;
};

/**
 * Function: setCloneEnabled
 * 
 * Sets <cloneEnabled>.
 * 
 * Parameters:
 * 
 * value - Boolean that specifies the new clone enabled state.
 */
mxGraphHandler.prototype.setCloneEnabled = function(value)
{
	this.cloneEnabled = value;
};

/**
 * Function: isMoveEnabled
 * 
 * Returns <moveEnabled>.
 */
mxGraphHandler.prototype.isMoveEnabled = function()
{
	return this.moveEnabled;
};

/**
 * Function: setMoveEnabled
 * 
 * Sets <moveEnabled>.
 */
mxGraphHandler.prototype.setMoveEnabled = function(value)
{
	this.moveEnabled = value;
};

/**
 * Function: isSelectEnabled
 * 
 * Returns <selectEnabled>.
 */
mxGraphHandler.prototype.isSelectEnabled = function()
{
	return this.selectEnabled;
};

/**
 * Function: setSelectEnabled
 * 
 * Sets <selectEnabled>.
 */
mxGraphHandler.prototype.setSelectEnabled = function(value)
{
	this.selectEnabled = value;
};

/**
 * Function: isRemoveCellsFromParent
 * 
 * Returns <removeCellsFromParent>.
 */
mxGraphHandler.prototype.isRemoveCellsFromParent = function()
{
	return this.removeCellsFromParent;
};

/**
 * Function: setRemoveCellsFromParent
 * 
 * Sets <removeCellsFromParent>.
 */
mxGraphHandler.prototype.setRemoveCellsFromParent = function(value)
{
	this.removeCellsFromParent = value;
};

/**
 * Function: getInitialCellForEvent
 * 
 * Hook to return initial cell for the given event.
 */
mxGraphHandler.prototype.getInitialCellForEvent = function(me)
{
	return me.getCell();
};

/**
 * Function: isDelayedSelection
 * 
 * Hook to return true for delayed selections.
 */
mxGraphHandler.prototype.isDelayedSelection = function(cell, me)
{
	return this.graph.isCellSelected(cell);
};

/**
 * Function: consumeMouseEvent
 * 
 * Consumes the given mouse event. NOTE: This may be used to enable click
 * events for links in labels on iOS as follows as consuming the initial
 * touchStart disables firing the subsequent click evnent on the link.
 * 
 * <code>
 * mxGraphHandler.prototype.consumeMouseEvent = function(evtName, me)
 * {
 *   var source = mxEvent.getSource(me.getEvent());
 *   
 *   if (!mxEvent.isTouchEvent(me.getEvent()) || source.nodeName != 'A')
 *   {
 *     me.consume();
 *   }
 * }
 * </code>
 */
mxGraphHandler.prototype.consumeMouseEvent = function(evtName, me)
{
	me.consume();
};

/**
 * Function: mouseDown
 * 
 * Handles the event by selecing the given cell and creating a handle for
 * it. By consuming the event all subsequent events of the gesture are
 * redirected to this handler.
 */
mxGraphHandler.prototype.mouseDown = function(sender, me)
{
	if (!me.isConsumed() && this.isEnabled() && this.graph.isEnabled() &&
		me.getState() != null && !mxEvent.isMultiTouchEvent(me.getEvent()))
	{
		var cell = this.getInitialCellForEvent(me);
		this.delayedSelection = this.isDelayedSelection(cell, me);
		this.cell = null;
		
		if (this.isSelectEnabled() && !this.delayedSelection)
		{
			this.graph.selectCellForEvent(cell, me.getEvent());
		}

		if (this.isMoveEnabled())
		{
			var model = this.graph.model;
			var geo = model.getGeometry(cell);

			if (this.graph.isCellMovable(cell) && ((!model.isEdge(cell) || this.graph.getSelectionCount() > 1 ||
				(geo.points != null && geo.points.length > 0) || model.getTerminal(cell, true) == null ||
				model.getTerminal(cell, false) == null) || this.graph.allowDanglingEdges || 
				(this.graph.isCloneEvent(me.getEvent()) && this.graph.isCellsCloneable())))
			{
				this.start(cell, me.getX(), me.getY());
			}
			else if (this.delayedSelection)
			{
				this.cell = cell;
			}

			this.cellWasClicked = true;
			this.consumeMouseEvent(mxEvent.MOUSE_DOWN, me);
		}
	}
};

/**
 * Function: getGuideStates
 * 
 * Creates an array of cell states which should be used as guides.
 */
mxGraphHandler.prototype.getGuideStates = function()
{
	var parent = this.graph.getDefaultParent();
	var model = this.graph.getModel();
	
	var filter = mxUtils.bind(this, function(cell)
	{
		return this.graph.view.getState(cell) != null &&
			model.isVertex(cell) &&
			model.getGeometry(cell) != null &&
			!model.getGeometry(cell).relative;
	});
	
	return this.graph.view.getCellStates(model.filterDescendants(filter, parent));
};

/**
 * Function: getCells
 * 
 * Returns the cells to be modified by this handler. This implementation
 * returns all selection cells that are movable, or the given initial cell if
 * the given cell is not selected and movable. This handles the case of moving
 * unselectable or unselected cells.
 * 
 * Parameters:
 * 
 * initialCell - <mxCell> that triggered this handler.
 */
mxGraphHandler.prototype.getCells = function(initialCell)
{
	if (!this.delayedSelection && this.graph.isCellMovable(initialCell))
	{
		return [initialCell];
	}
	else
	{
		return this.graph.getMovableCells(this.graph.getSelectionCells());
	}
};

/**
 * Function: getPreviewBounds
 * 
 * Returns the <mxRectangle> used as the preview bounds for
 * moving the given cells.
 */
mxGraphHandler.prototype.getPreviewBounds = function(cells)
{
	var bounds = this.getBoundingBox(cells);
	
	if (bounds != null)
	{
		// Corrects width and height
		bounds.width = Math.max(0, bounds.width - 1);
		bounds.height = Math.max(0, bounds.height - 1);
		
		if (bounds.width < this.minimumSize)
		{
			var dx = this.minimumSize - bounds.width;
			bounds.x -= dx / 2;
			bounds.width = this.minimumSize;
		}
		else
		{
			bounds.x = Math.round(bounds.x);
			bounds.width = Math.ceil(bounds.width);
		}
		
		var tr = this.graph.view.translate;
		var s = this.graph.view.scale;
		
		if (bounds.height < this.minimumSize)
		{
			var dy = this.minimumSize - bounds.height;
			bounds.y -= dy / 2;
			bounds.height = this.minimumSize;
		}
		else
		{
			bounds.y = Math.round(bounds.y);
			bounds.height = Math.ceil(bounds.height);
		}
	}
	
	return bounds;
};

/**
 * Function: getBoundingBox
 * 
 * Returns the union of the <mxCellStates> for the given array of <mxCells>.
 * For vertices, this method uses the bounding box of the corresponding shape
 * if one exists. The bounding box of the corresponding text label and all
 * controls and overlays are ignored. See also: <mxGraphView.getBounds> and
 * <mxGraph.getBoundingBox>.
 *
 * Parameters:
 *
 * cells - Array of <mxCells> whose bounding box should be returned.
 */
mxGraphHandler.prototype.getBoundingBox = function(cells)
{
	var result = null;
	
	if (cells != null && cells.length > 0)
	{
		var model = this.graph.getModel();
		
		for (var i = 0; i < cells.length; i++)
		{
			if (model.isVertex(cells[i]) || model.isEdge(cells[i]))
			{
				var state = this.graph.view.getState(cells[i]);
			
				if (state != null)
				{
					var bbox = state;
					
					if (model.isVertex(cells[i]) && state.shape != null && state.shape.boundingBox != null)
					{
						bbox = state.shape.boundingBox;
					}
					
					if (result == null)
					{
						result = mxRectangle.fromRectangle(bbox);
					}
					else
					{
						result.add(bbox);
					}
				}
			}
		}
	}
	
	return result;
};

/**
 * Function: createPreviewShape
 * 
 * Creates the shape used to draw the preview for the given bounds.
 */
mxGraphHandler.prototype.createPreviewShape = function(bounds)
{
	var shape = new mxRectangleShape(bounds, null, this.previewColor);
	shape.isDashed = true;
	
	if (this.htmlPreview)
	{
		shape.dialect = mxConstants.DIALECT_STRICTHTML;
		shape.init(this.graph.container);
	}
	else
	{
		// Makes sure to use either VML or SVG shapes in order to implement
		// event-transparency on the background area of the rectangle since
		// HTML shapes do not let mouseevents through even when transparent
		shape.dialect = (this.graph.dialect != mxConstants.DIALECT_SVG) ?
			mxConstants.DIALECT_VML : mxConstants.DIALECT_SVG;
		shape.init(this.graph.getView().getOverlayPane());
		shape.pointerEvents = false;
		
		// Workaround for artifacts on iOS
		if (mxClient.IS_IOS)
		{
			shape.getSvgScreenOffset = function()
			{
				return 0;
			};
		}
	}
	
	return shape;
};

/**
 * Function: start
 * 
 * Starts the handling of the mouse gesture.
 */
mxGraphHandler.prototype.start = function(cell, x, y)
{
	this.cell = cell;
	this.first = mxUtils.convertPoint(this.graph.container, x, y);
	this.cells = this.getCells(this.cell);
	this.bounds = this.graph.getView().getBounds(this.cells);
	this.pBounds = this.getPreviewBounds(this.cells);

	if (this.guidesEnabled)
	{
		this.guide = new mxGuide(this.graph, this.getGuideStates());
	}
};

/**
 * Function: useGuidesForEvent
 * 
 * Returns true if the guides should be used for the given <mxMouseEvent>.
 * This implementation returns <mxGuide.isEnabledForEvent>.
 */
mxGraphHandler.prototype.useGuidesForEvent = function(me)
{
	return (this.guide != null) ? this.guide.isEnabledForEvent(me.getEvent()) : true;
};


/**
 * Function: snap
 * 
 * Snaps the given vector to the grid and returns the given mxPoint instance.
 */
mxGraphHandler.prototype.snap = function(vector)
{
	var scale = (this.scaleGrid) ? this.graph.view.scale : 1;
	
	vector.x = this.graph.snap(vector.x / scale) * scale;
	vector.y = this.graph.snap(vector.y / scale) * scale;
	
	return vector;
};

/**
 * Function: getDelta
 * 
 * Returns an <mxPoint> that represents the vector for moving the cells
 * for the given <mxMouseEvent>.
 */
mxGraphHandler.prototype.getDelta = function(me)
{
	var point = mxUtils.convertPoint(this.graph.container, me.getX(), me.getY());
	var s = this.graph.view.scale;
	
	return new mxPoint(this.roundLength((point.x - this.first.x) / s) * s,
		this.roundLength((point.y - this.first.y) / s) * s);
};

/**
 * Function: updateHint
 * 
 * Hook for subclassers do show details while the handler is active.
 */
mxGraphHandler.prototype.updateHint = function(me) { };

/**
 * Function: removeHint
 * 
 * Hooks for subclassers to hide details when the handler gets inactive.
 */
mxGraphHandler.prototype.removeHint = function() { };

/**
 * Function: roundLength
 * 
 * Hook for rounding the unscaled vector. This uses Math.round.
 */
mxGraphHandler.prototype.roundLength = function(length)
{
	return Math.round(length * 2) / 2;
};

/**
 * Function: mouseMove
 * 
 * Handles the event by highlighting possible drop targets and updating the
 * preview.
 */
mxGraphHandler.prototype.mouseMove = function(sender, me)
{
	var graph = this.graph;

	if (!me.isConsumed() && graph.isMouseDown && this.cell != null &&
		this.first != null && this.bounds != null)
	{
		// Stops moving if a multi touch event is received
		if (mxEvent.isMultiTouchEvent(me.getEvent()))
		{
			this.reset();
			return;
		}
		
		var delta = this.getDelta(me);
		var dx = delta.x;
		var dy = delta.y;
		var tol = graph.tolerance;

		if (this.shape != null || Math.abs(dx) > tol || Math.abs(dy) > tol)
		{
			// Highlight is used for highlighting drop targets
			if (this.highlight == null)
			{
				this.highlight = new mxCellHighlight(this.graph,
					mxConstants.DROP_TARGET_COLOR, 3);
			}
			
			if (this.shape == null)
			{
				this.shape = this.createPreviewShape(this.bounds);
			}
			
			var clone = graph.isCloneEvent(me.getEvent()) && graph.isCellsCloneable() && this.isCloneEnabled();
			var gridEnabled = graph.isGridEnabledEvent(me.getEvent());
			var hideGuide = true;
			
			if (this.guide != null && this.useGuidesForEvent(me))
			{
				delta = this.guide.move(this.bounds, new mxPoint(dx, dy), gridEnabled, clone);
				hideGuide = false;
				dx = delta.x;
				dy = delta.y;
			}
			else if (gridEnabled)
			{
				var trx = graph.getView().translate;
				var scale = graph.getView().scale;				
				
				var tx = this.bounds.x - (graph.snap(this.bounds.x / scale - trx.x) + trx.x) * scale;
				var ty = this.bounds.y - (graph.snap(this.bounds.y / scale - trx.y) + trx.y) * scale;
				var v = this.snap(new mxPoint(dx, dy));
			
				dx = v.x - tx;
				dy = v.y - ty;
			}
			
			if (this.guide != null && hideGuide)
			{
				this.guide.hide();
			}

			// Constrained movement if shift key is pressed
			if (graph.isConstrainedEvent(me.getEvent()))
			{
				if (Math.abs(dx) > Math.abs(dy))
				{
					dy = 0;
				}
				else
				{
					dx = 0;
				}
			}

			this.currentDx = dx;
			this.currentDy = dy;
			this.updatePreviewShape();

			var target = null;
			var cell = me.getCell();

			if (graph.isDropEnabled() && this.highlightEnabled)
			{
				// Contains a call to getCellAt to find the cell under the mouse
				target = graph.getDropTarget(this.cells, me.getEvent(), cell, clone);
			}

			var state = graph.getView().getState(target);
			var highlight = false;
			
			if (state != null && (graph.model.getParent(this.cell) != target || clone))
			{
			    if (this.target != target)
			    {
				    this.target = target;
				    this.setHighlightColor(mxConstants.DROP_TARGET_COLOR);
				}
			    
			    highlight = true;
			}
			else
			{
				this.target = null;

				if (this.connectOnDrop && cell != null && this.cells.length == 1 &&
					graph.getModel().isVertex(cell) && graph.isCellConnectable(cell))
				{
					state = graph.getView().getState(cell);
					
					if (state != null)
					{
						var error = graph.getEdgeValidationError(null, this.cell, cell);
						var color = (error == null) ?
							mxConstants.VALID_COLOR :
							mxConstants.INVALID_CONNECT_TARGET_COLOR;
						this.setHighlightColor(color);
						highlight = true;
					}
				}
			}
			
			if (state != null && highlight)
			{
				this.highlight.highlight(state);
			}
			else
			{
				this.highlight.hide();
			}
		}

		this.updateHint(me);
		this.consumeMouseEvent(mxEvent.MOUSE_MOVE, me);
		
		// Cancels the bubbling of events to the container so
		// that the droptarget is not reset due to an mouseMove
		// fired on the container with no associated state.
		mxEvent.consume(me.getEvent());
	}
	else if ((this.isMoveEnabled() || this.isCloneEnabled()) && this.updateCursor && !me.isConsumed() &&
		(me.getState() != null || me.sourceState != null) && !graph.isMouseDown)
	{
		var cursor = graph.getCursorForMouseEvent(me);
		
		if (cursor == null && graph.isEnabled() && graph.isCellMovable(me.getCell()))
		{
			if (graph.getModel().isEdge(me.getCell()))
			{
				cursor = mxConstants.CURSOR_MOVABLE_EDGE;
			}
			else
			{
				cursor = mxConstants.CURSOR_MOVABLE_VERTEX;
			}
		}

		// Sets the cursor on the original source state under the mouse
		// instead of the event source state which can be the parent
		if (cursor != null && me.sourceState != null)
		{
			me.sourceState.setCursor(cursor);
		}
	}
};

/**
 * Function: updatePreviewShape
 * 
 * Updates the bounds of the preview shape.
 */
mxGraphHandler.prototype.updatePreviewShape = function()
{
	if (this.shape != null)
	{
		this.shape.bounds = new mxRectangle(Math.round(this.pBounds.x + this.currentDx - this.graph.panDx),
				Math.round(this.pBounds.y + this.currentDy - this.graph.panDy), this.pBounds.width, this.pBounds.height);
		this.shape.redraw();
	}
};

/**
 * Function: setHighlightColor
 * 
 * Sets the color of the rectangle used to highlight drop targets.
 * 
 * Parameters:
 * 
 * color - String that represents the new highlight color.
 */
mxGraphHandler.prototype.setHighlightColor = function(color)
{
	if (this.highlight != null)
	{
		this.highlight.setHighlightColor(color);
	}
};

/**
 * Function: mouseUp
 * 
 * Handles the event by applying the changes to the selection cells.
 */
mxGraphHandler.prototype.mouseUp = function(sender, me)
{
	if (!me.isConsumed())
	{
		var graph = this.graph;
		
		if (this.cell != null && this.first != null && this.shape != null &&
			this.currentDx != null && this.currentDy != null)
		{
			var cell = me.getCell();
			
			if (this.connectOnDrop && this.target == null && cell != null && graph.getModel().isVertex(cell) &&
				graph.isCellConnectable(cell) && graph.isEdgeValid(null, this.cell, cell))
			{
				graph.connectionHandler.connect(this.cell, cell, me.getEvent());
			}
			else
			{
				var clone = graph.isCloneEvent(me.getEvent()) && graph.isCellsCloneable() && this.isCloneEnabled();
				var scale = graph.getView().scale;
				var dx = this.roundLength(this.currentDx / scale);
				var dy = this.roundLength(this.currentDy / scale);
				var target = this.target;
				
				if (graph.isSplitEnabled() && graph.isSplitTarget(target, this.cells, me.getEvent()))
				{
					graph.splitEdge(target, this.cells, null, dx, dy);
				}
				else
				{
					this.moveCells(this.cells, dx, dy, clone, this.target, me.getEvent());
				}
			}
		}
		else if (this.isSelectEnabled() && this.delayedSelection && this.cell != null)
		{
			this.selectDelayed(me);
		}
	}

	// Consumes the event if a cell was initially clicked
	if (this.cellWasClicked)
	{
		this.consumeMouseEvent(mxEvent.MOUSE_UP, me);
	}

	this.reset();
};

/**
 * Function: selectDelayed
 * 
 * Implements the delayed selection for the given mouse event.
 */
mxGraphHandler.prototype.selectDelayed = function(me)
{
	if (!this.graph.isCellSelected(this.cell) || !this.graph.popupMenuHandler.isPopupTrigger(me))
	{
		this.graph.selectCellForEvent(this.cell, me.getEvent());
	}
};

/**
 * Function: reset
 * 
 * Resets the state of this handler.
 */
mxGraphHandler.prototype.reset = function()
{
	this.destroyShapes();
	this.removeHint();
	
	this.cellWasClicked = false;
	this.delayedSelection = false;
	this.currentDx = null;
	this.currentDy = null;
	this.guides = null;
	this.first = null;
	this.cell = null;
	this.target = null;
};

/**
 * Function: shouldRemoveCellsFromParent
 * 
 * Returns true if the given cells should be removed from the parent for the specified
 * mousereleased event.
 */
mxGraphHandler.prototype.shouldRemoveCellsFromParent = function(parent, cells, evt)
{
	if (this.graph.getModel().isVertex(parent))
	{
		var pState = this.graph.getView().getState(parent);
		
		if (pState != null)
		{
			var pt = mxUtils.convertPoint(this.graph.container,
				mxEvent.getClientX(evt), mxEvent.getClientY(evt));
			var alpha = mxUtils.toRadians(mxUtils.getValue(pState.style, mxConstants.STYLE_ROTATION) || 0);
			
			if (alpha != 0)
			{
				var cos = Math.cos(-alpha);
				var sin = Math.sin(-alpha);
				var cx = new mxPoint(pState.getCenterX(), pState.getCenterY());
				pt = mxUtils.getRotatedPoint(pt, cos, sin, cx);
			}
		
			return !mxUtils.contains(pState, pt.x, pt.y);
		}
	}
	
	return false;
};

/**
 * Function: moveCells
 * 
 * Moves the given cells by the specified amount.
 */
mxGraphHandler.prototype.moveCells = function(cells, dx, dy, clone, target, evt)
{
	if (clone)
	{
		cells = this.graph.getCloneableCells(cells);
	}
	
	// Removes cells from parent
	if (target == null && this.isRemoveCellsFromParent() &&
		this.shouldRemoveCellsFromParent(this.graph.getModel().getParent(this.cell), cells, evt))
	{
		target = this.graph.getDefaultParent();
	}
	
	// Cloning into locked cells is not allowed
	clone = clone && !this.graph.isCellLocked(target || this.graph.getDefaultParent());
	
	// Passes all selected cells in order to correctly clone or move into
	// the target cell. The method checks for each cell if its movable.
	cells = this.graph.moveCells(cells, dx - this.graph.panDx / this.graph.view.scale,
			dy - this.graph.panDy / this.graph.view.scale, clone, target, evt);
	
	if (this.isSelectEnabled() && this.scrollOnMove)
	{
		this.graph.scrollCellToVisible(cells[0]);
	}
			
	// Selects the new cells if cells have been cloned
	if (clone)
	{
		this.graph.setSelectionCells(cells);
	}
};

/**
 * Function: destroyShapes
 * 
 * Destroy the preview and highlight shapes.
 */
mxGraphHandler.prototype.destroyShapes = function()
{
	// Destroys the preview dashed rectangle
	if (this.shape != null)
	{
		this.shape.destroy();
		this.shape = null;
	}
	
	if (this.guide != null)
	{
		this.guide.destroy();
		this.guide = null;
	}
	
	// Destroys the drop target highlight
	if (this.highlight != null)
	{
		this.highlight.destroy();
		this.highlight = null;
	}
};

/**
 * Function: destroy
 * 
 * Destroys the handler and all its resources and DOM nodes.
 */
mxGraphHandler.prototype.destroy = function()
{
	this.graph.removeMouseListener(this);
	this.graph.removeListener(this.panHandler);
	
	if (this.escapeHandler != null)
	{
		this.graph.removeListener(this.escapeHandler);
		this.escapeHandler = null;
	}
	
	if (this.refreshHandler != null)
	{
		this.graph.getModel().removeListener(this.refreshHandler);
		this.refreshHandler = null;
	}
	
	this.destroyShapes();
	this.removeHint();
};
