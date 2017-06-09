/*
The MIT License (MIT)

Copyright (c) 2013 jamie.peabody@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

;(function( $, window, document, undefined ){
	var pluginName = 'wickedtoolbar';
	var SHIFT = new RegExp(/shift/i);
	var ALT = new RegExp(/alt/i);
	var CTRL = new RegExp(/ctrl/i);
	var ARROW_DOWN = new RegExp(/↓/);
	var ARROW_UP = new RegExp(/↑/);
	var ARROW_LEFT = new RegExp(/←/);
	var ARROW_RIGHT = new RegExp(/→/);
	
	var keys = {
		shift: 16,
		alt: 17,
		ctrl: 18,
		meta: 91,
		arrow_up: 38,
		arrow_down: 40,
		arrow_left: 37,
		arrow_right: 39
	};

	var defaults = {
		hasIcon: function(id) { },
		getIcon: function(id) { }
	};
	
	function MenuBase(element, options) {
		var self = this;
		
		this.element = $(element);
		this.settings = $.extend({}, defaults, options);
		this.bindings = [];
		// for each menu item, modify and wrap in div
		this.element.find('li').each(function () {
			var tthis = $(this);
			
			// since the DOM is modifying wrt to icons, it can match the li[data-cion] from
			// the HTML definition, or the span.icon from the modified DOM
			var icons = tthis.closest('ul').find('li[data-icon], span.icon').length > 0;
			
			var tnode = tthis.contents().filter(function() { return this.nodeType == 3; }).first();
			var text = tnode.text().trim();
			
			// remove the text node
			tnode.remove();
			
			if (!text || !text.length) return; // not a regular menu item
			
			// change: <li>Text</li>
			// to: <li>
			//			<a>
			//				<span>Text</span>
			//			</a>
			//		</li>
			var li = tthis;
			var div = $('<a class="menu-item" href="#">');
			div.click(function(ev) {
				$(this).focus(); 
				// li.id > a
				var id = $(this).parent().attr('id');
				if (id) {
					self.element.trigger('selected', [id]);
					$(this).parents('.hover').removeClass('hover');
				}
				return false; 
			});
			
			var span;
			if (self.settings._type == 'menu') {
				span = $('<span>' + text + '</span>');
			}
			else {
				span = $('<span></span>');
			}
			div.append(span);
			
			if (self.settings._type == 'menu') {
				// accesskey
				var accesskey = tthis.attr('accesskey');
				if (accesskey) {
					div.attr('accesskey', accesskey);
					tthis.removeAttr('accesskey');
				}
				
				// hotkey
				var hotkey = tthis.attr('data-hotkey');
				if (hotkey) {
					tthis.removeAttr('data-hotkey');
					div.append('<span class="hotkey">' + hotkey + '</span>');

					if (!accesskey) {
						// add our own handler
						var parts = hotkey.split('+');
						var bind = {};
						for (var i = 0; i < parts.length; ++i) {
							if (SHIFT.test(parts[i])) {
								bind.shiftKey = true;
							}
							else if (ALT.test(parts[i])) {
								bind.altKey = true;
							}
							else if (CTRL.test(parts[i])) {
								bind.ctrlKey = true;
							}
							else if (ARROW_DOWN.test(parts[i])) {
								bind.which = keys.arrow_down;
							}
							else if (ARROW_UP.test(parts[i])) {
								bind.which = keys.arrow_up;
							}
							else if (ARROW_RIGHT.test(parts[i])) {
								bind.which = keys.arrow_right;
							}
							else if (ARROW_LEFT.test(parts[i])) {
								bind.which = keys.arrow_left;
							}
						}
						bind.target = div;
						self.bindings.push(bind);
					}
				}
			}
			
			// icon
			var id = tthis.attr('id'), icon;
			if (self.settings.hasIcon(id)) {
				span.addClass('icon');
				icon = self.settings.getIcon(id);
				if (icon) {
					span.addClass(icon);
				}
			}
			icon = tthis.attr('data-icon');
			if (icon) {
				tthis.removeAttr('data-icon');
				span.addClass('icon ' + icon);
			}
			else if (icons) {
				span.addClass('icon');
			}
			li.prepend(div);
		});
		$(document).on('keydown', function(ev) {
			for (var i = 0; i < self.bindings.length; ++i) {
				var bind = self.bindings[i];
				// handle custom key events
				if ((bind.shiftKey === undefined ? true : (bind.shiftKey === ev.shiftKey)) &&
					(bind.ctrlKey === undefined ? true : (bind.ctrlKey === ev.ctrlKey)) &&
					(bind.altKey === undefined ? true : (bind.altKey === ev.altKey)) &&
					bind.which && ev.which && (bind.which === ev.which)) {
					bind.target.trigger('click');
					ev.preventDefault();
				}
			}
		});
	}
	
	MenuBase.prototype.update = function (id) {
		var li = this.element.find('#' + id), icon;
		var span = li.find('span:first-child');
		if (this.settings.hasIcon(id)) {
			span.removeClass(); // this could be brutally unfair
			span.addClass('icon');
			icon = this.settings.getIcon(id);
			if (icon) {
				span.addClass(icon);
			}
		}
	};
	
	// ------------
	// Menu
	// ------------
	function Menu(element, options) {
		options._type = 'menu';
		$.extend(this, new MenuBase(element, options)) ;
		this.constructor();
	}

	Menu.prototype.constructor = function () {
		this.element.addClass('wicked-ui wicked-menu');
		var self = this;
		var dohover = function(ev) {
			$(this).parent().addClass('hover');
			if ($(this).closest('ul').hasClass('wicked-menu')) {
				// if the closest ul is a 'menu' (i.e. if this item is a top-level menu), then
				// aply focus
				$(this).focus();
			}
			if (!self.accessing) {
				// Set 'accessing' to true and one document click to cancel it
				self.accessing = true;
				$(document).one('click', function() {
					if (!self.accessing) return;
					self.accessing = false;
					self.element.find('.hover').removeClass('hover');
				});
			}
		};
		
		this.element.find('> li > a.menu-item').click(dohover);
		
		this.element.find('> li > ul > li ul.drop-menu').each(function() {
			// li > a + ul
			$(this).prev('a.menu-item').addClass('icon-arrow-right');
			
			$(this).prev('a.menu-item').hover(dohover);
			$(this).prev('a.menu-item').mouseleave(
				function() {
					$(this).parent().removeClass('hover');
				}
			);
		});
		this.element.find('> li > a.menu-item').hover(
			function() {
				if (!self.accessing) return;
				self.element.find('.hover').removeClass('hover');
				$.proxy(dohover, this)();
			}
		);
	};
	
	// ------------
	// Toolbar
	// ------------
	function Toolbar(element, options) {
		options._type = 'toolbar';
		$.extend(this, new MenuBase(element, options)) ;
		this.constructor();
	}
	
	Toolbar.prototype.constructor = function () {
		this.element.addClass('wicked-ui wicked-toolbar');
		this.element.find('> li > a.menu-item').hover(
			function(){ $(this).parent().addClass('hover'); },
			function(){ $(this).parent().removeClass('hover'); }
		);
	};

	var plugins = { wickedmenu: Menu, wickedtoolbar: Toolbar };
	for (var key in plugins) {
		(function(name, Plugin){
			$.fn[name] = function (options) {
				var args = arguments;
				return this.each(function () {
					if (typeof options === 'object' || !options) {
						if (!$.data(this, 'plugin_' + name)) {
							$.data(this, 'plugin_' + name, new Plugin( this, options ));
						}
					}
					else {
						var d = $.data(this, 'plugin_' + name);
						if (!d) {
							$.error('jQuery.' + name + ' plugin does not exist');
							return;
						}
						return d[options](Array.prototype.slice.call(args, 1));
					}
				});
			};
		})(key, plugins[key])
	}
	
})( jQuery, window, document );
