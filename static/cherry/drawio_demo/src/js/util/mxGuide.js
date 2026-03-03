/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxGuide
 *
 * Implements the alignment of selection cells to other cells in the graph.
 * 
 * Constructor: mxGuide
 * 
 * Constructs a new guide object.
 */
function mxGuide(graph, states)
{
	this.graph = graph;
	this.setStates(states);
};

/**
 * Variable: graph
 *
 * Reference to the enclosing <mxGraph> instance.
 */
mxGuide.prototype.graph = null;

/**
 * Variable: states
 * 
 * Contains the <mxCellStates> that are used for alignment.
 */
mxGuide.prototype.states = null;

/**
 * Variable: horizontal
 *
 * Specifies if horizontal guides are enabled. Default is true.
 */
mxGuide.prototype.horizontal = true;

/**
 * Variable: vertical
 *
 * Specifies if vertical guides are enabled. Default is true.
 */
mxGuide.prototype.vertical = true;

/**
 * Variable: vertical
 *
 * Holds the <mxShape> for the horizontal guide.
 */
mxGuide.prototype.guideX = null;

/**
 * Variable: vertical
 *
 * Holds the <mxShape> for the vertical guide.
 */
mxGuide.prototype.guideY = null;

/**
 * Function: setStates
 * 
 * Sets the <mxCellStates> that should be used for alignment.
 */
mxGuide.prototype.setStates = function(states)
{
	this.states = states;
};

/**
 * Function: isEnabledForEvent
 * 
 * Returns true if the guide should be enabled for the given native event. This
 * implementation always returns true.
 */
mxGuide.prototype.isEnabledForEvent = function(evt)
{
	return true;
};

/**
 * Function: getGuideTolerance
 * 
 * Returns the tolerance for the guides. Default value is gridSize / 2.
 */
mxGuide.prototype.getGuideTolerance = function()
{
	return this.graph.gridSize / 2;
};

/**
 * Function: createGuideShape
 * 
 * Returns the mxShape to be used for painting the respective guide. This
 * implementation returns a new, dashed and crisp <mxPolyline> using
 * <mxConstants.GUIDE_COLOR> and <mxConstants.GUIDE_STROKEWIDTH> as the format.
 * 
 * Parameters:
 * 
 * horizontal - Boolean that specifies which guide should be created.
 */
mxGuide.prototype.createGuideShape = function(horizontal)
{
	var guide = new mxPolyline([], mxConstants.GUIDE_COLOR, mxConstants.GUIDE_STROKEWIDTH);
	guide.isDashed = true;
	
	return guide;
};

/**
 * Function: move
 * 
 * Moves the <bounds> by the given <mxPoint> and returnt the snapped point.
 */
mxGuide.prototype.move = function(bounds, delta, gridEnabled, clone)
{
	if (this.states != null && (this.horizontal || this.vertical) && bounds != null && delta != null)
	{
		var trx = this.graph.getView().translate;
		var scale = this.graph.getView().scale;
		var dx = delta.x;
		var dy = delta.y;
		
		var overrideX = false;
		var stateX = null;
		var valueX = null;
		var overrideY = false;
		var stateY = null;
		var valueY = null;
		
		var tt = this.getGuideTolerance();
		var ttX = tt;
		var ttY = tt;
		
		var b = bounds.clone();
		b.x += delta.x;
		b.y += delta.y;
		
		var left = b.x;
		var right = b.x + b.width;
		var center = b.getCenterX();
		var top = b.y;
		var bottom = b.y + b.height;
		var middle = b.getCenterY();
	
		// Snaps the left, center and right to the given x-coordinate
		function snapX(x, state)
		{
			x += this.graph.panDx;
			var override = false;
			
			if (Math.abs(x - center) < ttX)
			{
				dx = x - bounds.getCenterX();
				ttX = Math.abs(x - center);
				override = true;
			}
			else if (Math.abs(x - left) < ttX)
			{
				dx = x - bounds.x;
				ttX = Math.abs(x - left);
				override = true;
			}
			else if (Math.abs(x - right) < ttX)
			{
				dx = x - bounds.x - bounds.width;
				ttX = Math.abs(x - right);
				override = true;
			}
			
			if (override)
			{
				stateX = state;
				valueX = Math.round(x - this.graph.panDx);
				
				if (this.guideX == null)
				{
					this.guideX = this.createGuideShape(true);
					
					// Makes sure to use either VML or SVG shapes in order to implement
					// event-transparency on the background area of the rectangle since
					// HTML shapes do not let mouseevents through even when transparent
					this.guideX.dialect = (this.graph.dialect != mxConstants.DIALECT_SVG) ?
						mxConstants.DIALECT_VML : mxConstants.DIALECT_SVG;
					this.guideX.pointerEvents = false;
					this.guideX.init(this.graph.getView().getOverlayPane());
				}
			}
			
			overrideX = overrideX || override;
		};
		
		// Snaps the top, middle or bottom to the given y-coordinate
		function snapY(y)
		{
			y += this.graph.panDy;
			var override = false;
			
			if (Math.abs(y - middle) < ttY)
			{
				dy = y - bounds.getCenterY();
				ttY = Math.abs(y -  middle);
				override = true;
			}
			else if (Math.abs(y - top) < ttY)
			{
				dy = y - bounds.y;
				ttY = Math.abs(y - top);
				override = true;
			}
			else if (Math.abs(y - bottom) < ttY)
			{
				dy = y - bounds.y - bounds.height;
				ttY = Math.abs(y - bottom);
				override = true;
			}
			
			if (override)
			{
				stateY = state;
				valueY = Math.round(y - this.graph.panDy);
				
				if (this.guideY == null)
				{
					this.guideY = this.createGuideShape(false);
					
					// Makes sure to use either VML or SVG shapes in order to implement
					// event-transparency on the background area of the rectangle since
					// HTML shapes do not let mouseevents through even when transparent
					this.guideY.dialect = (this.graph.dialect != mxConstants.DIALECT_SVG) ?
						mxConstants.DIALECT_VML : mxConstants.DIALECT_SVG;
					this.guideY.pointerEvents = false;
					this.guideY.init(this.graph.getView().getOverlayPane());
				}
			}
			
			overrideY = overrideY || override;
		};
		
		for (var i = 0; i < this.states.length; i++)
		{
			var state =  this.states[i];
			
			if (state != null)
			{
				// Align x
				if (this.horizontal)
				{
					snapX.call(this, state.getCenterX(), state);
					snapX.call(this, state.x, state);
					snapX.call(this, state.x + state.width, state);
				}
	
				// Align y
				if (this.vertical)
				{
					snapY.call(this, state.getCenterY(), state);
					snapY.call(this, state.y, state);
					snapY.call(this, state.y + state.height, state);
				}
			}
		}

		// Moves cells that are off-grid back to the grid on move
		if (gridEnabled)
		{
			if (!overrideX)
			{
				var tx = bounds.x - (this.graph.snap(bounds.x /
					scale - trx.x) + trx.x) * scale;
				dx = this.graph.snap(dx / scale) * scale - tx;
			}
			
			if (!overrideY)
			{
				var ty = bounds.y - (this.graph.snap(bounds.y /
					scale - trx.y) + trx.y) * scale;
				dy = this.graph.snap(dy / scale) * scale - ty;
			}
		}
		
		// Redraws the guides
		var c = this.graph.container;
		
		if (!overrideX && this.guideX != null)
		{
			this.guideX.node.style.visibility = 'hidden';
		}
		else if (this.guideX != null)
		{
			if (stateX != null && bounds != null)
			{
				minY = Math.min(bounds.y + dy - this.graph.panDy, stateX.y);
				maxY = Math.max(bounds.y + bounds.height + dy - this.graph.panDy, stateX.y + stateX.height);
			}
			
			if (minY != null && maxY != null)
			{
				this.guideX.points = [new mxPoint(valueX, minY), new mxPoint(valueX, maxY)];
			}
			else
			{
				this.guideX.points = [new mxPoint(valueX, -this.graph.panDy), new mxPoint(valueX, c.scrollHeight - 3 - this.graph.panDy)];
			}
			
			this.guideX.stroke = this.getGuideColor(stateX, true);
			this.guideX.node.style.visibility = 'visible';
			this.guideX.redraw();
		}
		
		if (!overrideY && this.guideY != null)
		{
			this.guideY.node.style.visibility = 'hidden';
		}
		else if (this.guideY != null)
		{
			if (stateY != null && bounds != null)
			{
				minX = Math.min(bounds.x + dx - this.graph.panDx, stateY.x);
				maxX = Math.max(bounds.x + bounds.width + dx - this.graph.panDx, stateY.x + stateY.width);
			}
			
			if (minX != null && maxX != null)
			{
				this.guideY.points = [new mxPoint(minX, valueY), new mxPoint(maxX, valueY)];
			}
			else
			{
				this.guideY.points = [new mxPoint(-this.graph.panDx, valueY), new mxPoint(c.scrollWidth - 3 - this.graph.panDx, valueY)];
			}
			
			this.guideY.stroke = this.getGuideColor(stateY, false);
			this.guideY.node.style.visibility = 'visible';
			this.guideY.redraw();
		}
		
		delta = new mxPoint(dx, dy);
	}
	
	return delta;
};

/**
 * Function: hide
 * 
 * Hides all current guides.
 */
mxGuide.prototype.getGuideColor = function(state, horizontal)
{
	return mxConstants.GUIDE_COLOR;
};

/**
 * Function: hide
 * 
 * Hides all current guides.
 */
mxGuide.prototype.hide = function()
{
	this.setVisible(false);
};

/**
 * Function: setVisible
 * 
 * Shows or hides the current guides.
 */
mxGuide.prototype.setVisible = function(visible)
{
	if (this.guideX != null)
	{
		this.guideX.node.style.visibility = (visible) ? 'visible' : 'hidden';
	}
	
	if (this.guideY != null)
	{
		this.guideY.node.style.visibility = (visible) ? 'visible' : 'hidden';
	}
};

/**
 * Function: destroy
 * 
 * Destroys all resources that this object uses.
 */
mxGuide.prototype.destroy = function()
{
	if (this.guideX != null)
	{
		this.guideX.destroy();
		this.guideX = null;
	}
	
	if (this.guideY != null)
	{
		this.guideY.destroy();
		this.guideY = null;
	}
};
