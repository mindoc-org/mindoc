/**
 * Copyright (c) 2006-2015, JGraph Ltd
 * Copyright (c) 2006-2015, Gaudenz Alder
 */
/**
 * Class: mxDefaultPopupMenu
 *
 * Creates popupmenus for mouse events. This object holds an XML node
 * which is a description of the popup menu to be created. In
 * <createMenu>, the configuration is applied to the context and
 * the resulting menu items are added to the menu dynamically. See
 * <createMenu> for a description of the configuration format.
 * 
 * This class does not create the DOM nodes required for the popup menu, it
 * only parses an XML description to invoke the respective methods on an
 * <mxPopupMenu> each time the menu is displayed.
 *
 * Codec:
 * 
 * This class uses the <mxDefaultPopupMenuCodec> to read configuration
 * data into an existing instance, however, the actual parsing is done
 * by this class during program execution, so the format is described
 * below.
 * 
 * Constructor: mxDefaultPopupMenu
 *
 * Constructs a new popupmenu-factory based on given configuration.
 *
 * Paramaters:
 *
 * config - XML node that contains the configuration data.
 */
function mxDefaultPopupMenu(config)
{
	this.config = config;
};

/**
 * Variable: imageBasePath
 *
 * Base path for all icon attributes in the config. Default is null.
 */
mxDefaultPopupMenu.prototype.imageBasePath = null;

/**
 * Variable: config
 *
 * XML node used as the description of new menu items. This node is
 * used in <createMenu> to dynamically create the menu items if their
 * respective conditions evaluate to true for the given arguments.
 */
mxDefaultPopupMenu.prototype.config = null;

/**
 * Function: createMenu
 *
 * This function is called from <mxEditor> to add items to the
 * given menu based on <config>. The config is a sequence of
 * the following nodes and attributes.
 *
 * Child Nodes: 
 *
 * add - Adds a new menu item. See below for attributes.
 * separator - Adds a separator. No attributes.
 * condition - Adds a custom condition. Name attribute.
 * 
 * The add-node may have a child node that defines a function to be invoked
 * before the action is executed (or instead of an action to be executed).
 *
 * Attributes:
 *
 * as - Resource key for the label (needs entry in property file).
 * action - Name of the action to execute in enclosing editor.
 * icon - Optional icon (relative/absolute URL).
 * iconCls - Optional CSS class for the icon.
 * if - Optional name of condition that must be true (see below).
 * enabled-if - Optional name of condition that specifies if the menu item
 * should be enabled.
 * name - Name of custom condition. Only for condition nodes.
 *
 * Conditions:
 *
 * nocell - No cell under the mouse.
 * ncells - More than one cell selected.
 * notRoot - Drilling position is other than home.
 * cell - Cell under the mouse.
 * notEmpty - Exactly one cell with children under mouse.
 * expandable - Exactly one expandable cell under mouse.
 * collapsable - Exactly one collapsable cell under mouse.
 * validRoot - Exactly one cell which is a possible root under mouse.
 * swimlane - Exactly one cell which is a swimlane under mouse.
 *
 * Example:
 *
 * To add a new item for a given action to the popupmenu:
 * 
 * (code)
 * <mxDefaultPopupMenu as="popupHandler">
 *   <add as="delete" action="delete" icon="images/delete.gif" if="cell"/>
 * </mxDefaultPopupMenu>
 * (end)
 * 
 * To add a new item for a custom function:
 * 
 * (code)
 * <mxDefaultPopupMenu as="popupHandler">
 *   <add as="action1"><![CDATA[
 *		function (editor, cell, evt)
 *		{
 *			editor.execute('action1', cell, 'myArg');
 *		}
 *   ]]></add>
 * </mxDefaultPopupMenu>
 * (end)
 * 
 * The above example invokes action1 with an additional third argument via
 * the editor instance. The third argument is passed to the function that
 * defines action1. If the add-node has no action-attribute, then only the
 * function defined in the text content is executed, otherwise first the
 * function and then the action defined in the action-attribute is
 * executed. The function in the text content has 3 arguments, namely the
 * <mxEditor> instance, the <mxCell> instance under the mouse, and the
 * native mouse event.
 *
 * Custom Conditions:
 *
 * To add a new condition for popupmenu items:
 *  
 * (code)
 * <condition name="condition1"><![CDATA[
 *   function (editor, cell, evt)
 *   {
 *     return cell != null;
 *   }
 * ]]></condition>
 * (end)
 * 
 * The new condition can then be used in any item as follows:
 * 
 * (code)
 * <add as="action1" action="action1" icon="action1.gif" if="condition1"/>
 * (end)
 * 
 * The order in which the items and conditions appear is not significant as
 * all connditions are evaluated before any items are created.
 * 
 * Parameters:
 *
 * editor - Enclosing <mxEditor> instance.
 * menu - <mxPopupMenu> that is used for adding items and separators.
 * cell - Optional <mxCell> which is under the mousepointer.
 * evt - Optional mouse event which triggered the menu. 
 */
mxDefaultPopupMenu.prototype.createMenu = function(editor, menu, cell, evt)
{
	if (this.config != null)
	{
		var conditions = this.createConditions(editor, cell, evt);
		var item = this.config.firstChild;

		this.addItems(editor, menu, cell, evt, conditions, item, null);
	}
};

/**
 * Function: addItems
 * 
 * Recursively adds the given items and all of its children into the given menu.
 * 
 * Parameters:
 *
 * editor - Enclosing <mxEditor> instance.
 * menu - <mxPopupMenu> that is used for adding items and separators.
 * cell - Optional <mxCell> which is under the mousepointer.
 * evt - Optional mouse event which triggered the menu.
 * conditions - Array of names boolean conditions.
 * item - XML node that represents the current menu item.
 * parent - DOM node that represents the parent menu item.
 */
mxDefaultPopupMenu.prototype.addItems = function(editor, menu, cell, evt, conditions, item, parent)
{
	var addSeparator = false;
	
	while (item != null)
	{
		if (item.nodeName == 'add')
		{
			var condition = item.getAttribute('if');
			
			if (condition == null || conditions[condition])
			{
				var as = item.getAttribute('as');
				as = mxResources.get(as) || as;
				var funct = mxUtils.eval(mxUtils.getTextContent(item));
				var action = item.getAttribute('action');
				var icon = item.getAttribute('icon');
				var iconCls = item.getAttribute('iconCls');
				var enabledCond = item.getAttribute('enabled-if');
				var enabled = enabledCond == null || conditions[enabledCond];
				
				if (addSeparator)
				{
					menu.addSeparator(parent);
					addSeparator = false;
				}
				
				if (icon != null && this.imageBasePath)
				{
					icon = this.imageBasePath + icon;
				}
				
				var row = this.addAction(menu, editor, as, icon, funct, action, cell, parent, iconCls, enabled);
				this.addItems(editor, menu, cell, evt, conditions, item.firstChild, row);
			}
		}
		else if (item.nodeName == 'separator')
		{
			addSeparator = true;
		}
		
		item = item.nextSibling;
	}
};

/**
 * Function: addAction
 *
 * Helper method to bind an action to a new menu item.
 * 
 * Parameters:
 *
 * menu - <mxPopupMenu> that is used for adding items and separators.
 * editor - Enclosing <mxEditor> instance.
 * lab - String that represents the label of the menu item.
 * icon - Optional URL that represents the icon of the menu item.
 * action - Optional name of the action to execute in the given editor.
 * funct - Optional function to execute before the optional action. The
 * function takes an <mxEditor>, the <mxCell> under the mouse and the
 * mouse event that triggered the call.
 * cell - Optional <mxCell> to use as an argument for the action.
 * parent - DOM node that represents the parent menu item.
 * iconCls - Optional CSS class for the menu icon.
 * enabled - Optional boolean that specifies if the menu item is enabled.
 * Default is true.
 */
mxDefaultPopupMenu.prototype.addAction = function(menu, editor, lab, icon, funct, action, cell, parent, iconCls, enabled)
{
	var clickHandler = function(evt)
	{
		if (typeof(funct) == 'function')
		{
			funct.call(editor, editor, cell, evt);
		}
		
		if (action != null)
		{
			editor.execute(action, cell, evt);
		}
	};
	
	return menu.addItem(lab, icon, clickHandler, parent, iconCls, enabled);
};

/**
 * Function: createConditions
 * 
 * Evaluates the default conditions for the given context.
 */
mxDefaultPopupMenu.prototype.createConditions = function(editor, cell, evt)
{
	// Creates array with conditions
	var model = editor.graph.getModel();
	var childCount = model.getChildCount(cell);
	
	// Adds some frequently used conditions
	var conditions = [];
	conditions['nocell'] = cell == null;
	conditions['ncells'] = editor.graph.getSelectionCount() > 1;
	conditions['notRoot'] = model.getRoot() !=
		model.getParent(editor.graph.getDefaultParent());
	conditions['cell'] = cell != null;
	
	var isCell = cell != null && editor.graph.getSelectionCount() == 1;
	conditions['nonEmpty'] = isCell && childCount > 0;
	conditions['expandable'] = isCell && editor.graph.isCellFoldable(cell, false);
	conditions['collapsable'] = isCell && editor.graph.isCellFoldable(cell, true);
	conditions['validRoot'] = isCell && editor.graph.isValidRoot(cell);
	conditions['emptyValidRoot'] = conditions['validRoot'] && childCount == 0;
	conditions['swimlane'] = isCell && editor.graph.isSwimlane(cell);

	// Evaluates dynamic conditions from config file
	var condNodes = this.config.getElementsByTagName('condition');
	
	for (var i=0; i<condNodes.length; i++)
	{
		var funct = mxUtils.eval(mxUtils.getTextContent(condNodes[i]));
		var name = condNodes[i].getAttribute('name');
		
		if (name != null && typeof(funct) == 'function')
		{
			conditions[name] = funct(editor, cell, evt);
		}
	}
	
	return conditions;
};
