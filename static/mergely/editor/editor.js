$(document).ready(function() {
	function getParameters() {
		var parameters = {};
		window.location.search.substr(1).split('&').forEach(function(pair) {
			if (pair === '') return;
			var parts = pair.split('=');
			if (parts.length === 2 && parts[1].search(/^(true|1)$/i) >= 0) {
				parameters[parts[0]] = true;
			}
			else if (parts.length === 2 && parts[1].search(/^(false|0)$/i) >= 0) {
				parameters[parts[0]] = false;
			}
			else parameters[parts[0]] = parts[1] && decodeURIComponent(parts[1].replace(/\+/g, ' '));
		});
		return {
			get: function(name, defaultValue) {
				if (parameters.hasOwnProperty(name)) return parameters[name];
				return defaultValue;
			}
		};
	}
	var parameters = getParameters();
	if (parameters.get('test', false)) {
		var li = $('<li>Tests</li>');
		var ul = $('<ul>');
		for (var i = 1; i <= 8; ++i) {
			ul.append($('<li id="examples-test' + i + '">Test ' + i + '</li>'));
		}
		li.append(ul);
		$('#main-menu').append(li);
	}
	
	function handleFind(column) {
		if (!column.length) {
			return false;
		}
		var ed = $('#mergely');
		var find = column.find('.find');
		var input = find.find('input[type="text"]');
		var side = column.attr('id').indexOf('-lhs') > 0 ? 'lhs' : 'rhs';
		var origautoupdate = ed.mergely('options').autoupdate;
		find.slideDown('fast', function() {
			input.focus();
			// disable autoupdate, clear both sides of diff
			ed.mergely('options', {autoupdate: false});
			ed.mergely('unmarkup');
		});
		find.find('.find-prev').click(function() {
			ed.mergely('search', side, input.val(), 'prev');
		});
		find.find('.find-next').click(function() {
			ed.mergely('search', side, input.val(), 'next');
		});
		find.find('.find-close').click(function() {
			find.css('display', 'none')
			ed.mergely('options', {autoupdate: origautoupdate});
		});
		
		input.keydown(function(evt) {
			if (evt.which != 13 && evt.which != 27) return true;
			if (evt.which == 27) {
				find.css('display', 'none');
				ed.mergely('options', {autoupdate: origautoupdate});
			}
			ed.mergely('search', side, input.val());
			return false;
		});
	}

	$(document).keydown(function(event) {
		if (!( String.fromCharCode(event.which).toLowerCase() == 'f' && event.ctrlKey)) return true;
		event.preventDefault();
		var range = window.getSelection().getRangeAt(0);
		var column = $(range.commonAncestorContainer).parents('.mergely-column');
		handleFind(column);
		return false;
	});

	String.prototype.random = function(length) {
		var chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXTZabcdefghiklmnopqrstuvwxyz";
		var randomstring = ''
		for (var i=0; i<length; i++) {
			var rnum = Math.floor(Math.random() * chars.length);
			randomstring += chars.substring(rnum,rnum+1);
		}
		return randomstring;
	}

	// body is pre-hidden for better rendering
	$('body').css("visibility", "");

	var ed = $('#mergely');
	var menu = $('#main-menu');
	var toolbar = $('#toolbar');
	ed.mergely({
		width: 'auto',
		height: 'auto',
		cmsettings: {
			lineNumbers: true,
			readOnly: false
		}
	});

    ed.mergely("lhs", $("#historyContent").html());
    ed.mergely("rhs", $("#documentContent").html());
	// if (parameters.get('lhs', null)) {
	// 	var url = parameters.get('lhs');
	// 	crossdomainGET(ed, 'lhs', url);
	// }
	// if (parameters.get('rhs', null)) {
	// 	var url = parameters.get('rhs');
	// 	crossdomainGET(ed, 'rhs', url);
	// }

	// set query string options
	var urloptions = {};
	var optmap = {
		au: 'autoupdate',
		ws: 'ignorews',
		cs: 'ignorecase',
		sb: 'sidebar',
		vp: 'viewport',
		wl: 'wrap_lines',
		ln: 'line_numbers'
	};
	var doopt = false;
	for (var name in optmap) {
		if (!optmap.hasOwnProperty(name)) continue;
		if (parameters.get(name, null) !== null) {
			urloptions[optmap[name]] = parameters.get(name);
			doopt = true;
		}
	}
	if (parameters.get('rm', null) !== null) {
		// special-case url property
		urloptions.rhs_margin = parameters.get('rm') ? 'left' : 'right';
	}
	if (doopt) {
		// apply query-string options
		ed.mergely('options', urloptions);
	}

	// set query string colors
	// cb: change border
	// cg: change background
	// ab: added border
	// ag: added background
	// db: deleted border
	// dg: deleted background
	var color_defaults = {
		cb: 'cccccc', cg: 'fafafa',
		ab: 'a3d1ff', ag: 'ddeeff',
		db: 'ff7f7f', dg: 'ffe9e9'
	};
	applyParameterCss(false);

	//history.pushState({}, null, '');

	window.addEventListener('popstate', function(ev) {
		if (ev.state) {
			parameters = getParameters();
			applyParameterCss(false);
		}
	});
	
	// Load
	if (key.length === 8) {
		$.when(
			$.ajax({
				type: 'GET', async: true, dataType: 'text',
				data: { 'key':key, 'name': 'lhs' },
				url: '/ajax/handle_get.php',
				success: function (response) {
					ed.mergely('lhs', response);
				},
				error: function(xhr, ajaxOptions, thrownError){
				}
			}),
			$.ajax({
				type: 'GET', async: true, dataType: 'text',
				data: { 'key':key, 'name': 'rhs' },
				url: '/ajax/handle_get.php',
				success: function (response) {
					ed.mergely('rhs', response);
				},
				error: function(xhr, ajaxOptions, thrownError){
				}
			})
		).done(function() {
			var anchor = window.location.hash.substring(1);
			if (anchor) {
				// if an anchor has been provided, then parse the anchor in the
				// form of: 'lhs' or 'rhs', followed by a line, e.g: lhs100.
				var m = anchor.match(/([lr]hs)([0-9]+)/);
				if (m && m.length == 3) {
					console.log(m);
					ed.mergely('scrollTo', m[1], parseInt(m[2],10));
				}
			}
		});
	}

	// find
	var find = $('.find');
	var flhs = find.clone().attr('id', 'mergely-editor-lhs-find');
	var frhs = find.clone().attr('id', 'mergely-editor-rhs-find');
	$('#mergely-editor-lhs').append(flhs);
	$('#mergely-editor-rhs').append(frhs);
	find.remove();
	
	var iconconf = {
		'options-autodiff': {
			get: function() { return ed.mergely('options').autoupdate },
			set: function(value) {
				var au = !ed.mergely('options').autoupdate;
				ed.mergely('options', {autoupdate: au});
				var params = updateQueryStringParam('au', au ? 1 : 0, 1);
				updateHistory(params);
			}
		},
		'options-ignorews': {
			get: function() { return ed.mergely('options').ignorews },
			set: function(value) {
				var ws = !ed.mergely('options').ignorews;
				ed.mergely('options', {ignorews: ws});
				var params = updateQueryStringParam('ws', ws ? 1 : 0, 0);
				updateHistory(params);
			}
		},
		'options-ignorecase': {
			get: function() { return ed.mergely('options').ignorecase },
			set: function(value) {
				var cs = !ed.mergely('options').ignorecase;
				ed.mergely('options', {ignorecase: cs});
				var params = updateQueryStringParam('cs', cs ? 1 : 0, 0);
				updateHistory(params);
			}
		},
		'options-sidebars': {
			get: function() { console.log('sidebar', this); return ed.mergely('options').sidebar },
			set: function(value) {
				var sb = !ed.mergely('options').sidebar;
				ed.mergely('options', {sidebar: sb});
				var params = updateQueryStringParam('sb', sb ? 1 : 0, 1);
				updateHistory(params);
			}
		},
		'options-viewport': {
			get: function() { console.log('viewport', this); return ed.mergely('options').viewport },
			set: function(value) {
				var vp = !ed.mergely('options').viewport;
				ed.mergely('options', {viewport: vp});
				var params = updateQueryStringParam('vp', vp ? 1 : 0, 0);
				updateHistory(params);
			}
		},
		'options-swapmargin': {
			get: function() { return (ed.mergely('options').rhs_margin == 'left'); },
			set: function(value) {
				var rm = ed.mergely('options').rhs_margin == 'left' ? 'right' : 'left';
				ed.mergely('options', {rhs_margin: rm });
				var params = updateQueryStringParam('rm', rm == 'left' ? 1 : 0, 0);
				updateHistory(params);
			}
		},
		'options-linenumbers': {
			get: function() { console.log('wrap', this); return ed.mergely('options').line_numbers },
			set: function(value) {
				var ln = !ed.mergely('options').line_numbers;
				ed.mergely('options', {line_numbers: ln});
				var params = updateQueryStringParam('ln', ln ? 1 : 0, 1);
				updateHistory(params);
			}
		},
		'options-wrap': {
			get: function() { console.log('wrap', this); return ed.mergely('options').wrap_lines },
			set: function(value) {
				var wl = !ed.mergely('options').wrap_lines;
				ed.mergely('options', {wrap_lines: wl});
				var params = updateQueryStringParam('wl', wl ? 1 : 0, 0);
				updateHistory(params);
			}
		},
		'edit-left-readonly': {
			get: function() { return ed.mergely('cm', 'lhs').getOption('readOnly'); },
			set: function(value) { ed.mergely('cm', 'lhs').setOption('readOnly', value); }
		},
		'edit-right-readonly': {
			get: function() { return ed.mergely('cm', 'rhs').getOption('readOnly'); },
			set: function(value) { ed.mergely('cm', 'rhs').setOption('readOnly', value); }
		}
	}

	var menu_opts = {
		hasIcon: function(id) {
			return iconconf.hasOwnProperty(id);
		},
		getIcon: function(id) {
			if (iconconf[id].get()) return 'icon-check';
		}
	};

	function handle_operation(id) {
		if (id == 'file-new') {
			window.location = '/editor';
		}
		else if (id === 'file-save') {
			var rhs = ed.mergely('get', 'rhs');

			if(window.top.hasOwnProperty("editor")){
                if(window.top.editor.hasOwnProperty("$txt")){
                    window.top.editor.$txt.html(rhs);
				}else{

                    window.top.editor.clear();
                    window.top.editor.insertValue(rhs);
				}

                window.top.layer.closeAll();
			}

		}else if (id == 'file-share') {
			handleShare(ed);
		}
		else if (id == 'file-import') {
			importFiles(ed);
		}
		else if (id == 'edit-left-undo') {
			ed.mergely('cm', 'lhs').getDoc().undo();
		}
		else if (id == 'edit-left-redo') {
			ed.mergely('cm', 'lhs').getDoc().redo();
		}
		else if (id == 'edit-right-undo') {
			ed.mergely('cm', 'rhs').getDoc().undo();
		}
		else if (id == 'edit-right-redo') {
			ed.mergely('cm', 'rhs').getDoc().redo();
		}
		else if (id == 'edit-left-find') {
			handleFind(ed.find('#mergely-editor-lhs'));
		}
		else if (id == 'edit-left-merge-right') {
			ed.mergely('mergeCurrentChange', 'rhs');
		}
		else if (id == 'edit-left-merge-right-file') {
			ed.mergely('merge', 'rhs');
		}
		else if ([
			'edit-left-readonly', 
			'edit-right-readonly',
			'options-autodiff',
			'options-sidebars',
			'options-swapmargin',
			'options-viewport',
			'options-ignorews',
			'options-ignorecase',
			'options-wrap',
			'options-linenumbers',
			].indexOf(id) >= 0) {
			iconconf[id].set(!iconconf[id].get());
			menu.wickedmenu('update', id);
		}
		else if (id == 'edit-left-clear') {
			ed.mergely('clear', 'lhs');
		}
		else if (id == 'edit-right-find') {
			handleFind(ed.find('#mergely-editor-rhs'));
		}
		else if (id == 'edit-right-merge-left') {
			ed.mergely('mergeCurrentChange', 'lhs');
		}
		else if (id == 'edit-right-merge-left-file') {
			ed.mergely('merge', 'lhs');
		}
		else if (id == 'edit-right-clear') {
			ed.mergely('clear', 'rhs');
		}
		else if (id == 'options-colors') {
			colorSettings(ed);
		}
		else if (id == 'view-swap') {
			ed.mergely('swap');
		}
		else if (id == 'view-refresh') {
			ed.mergely('update');
		}
		else if (id == 'view-change-next') {
			ed.mergely('scrollToDiff', 'next');
		}
		else if (id == 'view-change-prev') {
			ed.mergely('scrollToDiff', 'prev');
		}
		else if (id == 'view-clear') {
			ed.mergely('unmarkup');
		}
		else if (id.indexOf('examples-') == 0) {
			var test_config = {
				test1: {lhs: 'one\ntwo\nthree', rhs: 'two\nthree'},
				test2: {lhs: 'two\nthree', rhs: 'one\ntwo\nthree'},
				test3: {lhs: 'one\nthree', rhs: 'one\ntwo\nthree'},
				test4: {lhs: 'one\ntwo\nthree', rhs: 'one\nthree'},
				test5: {lhs: 'to bee, or not to be', rhs: 'to be, or not to bee'},
				test6: {lhs: 'to be, or not to be z', rhs: 'to be, to be'},
				test7: {lhs: 'remained, & to assume', rhs: 'and to assume'},
				test8: {lhs: 'to be, or not to be', rhs: 'to  be,  or  not  to  be'}
			};
			var test = id.split('examples-')[1];
			ed.mergely('lhs', test_config[test]['lhs']);
			ed.mergely('rhs', test_config[test]['rhs']);
		}
		return false;
	}

	menu.wickedmenu(menu_opts).bind('selected', function(ev, id) {
		return handle_operation(id);
	});
	
	toolbar.wickedtoolbar({}).bind('selected', function(ev, id) {
		if (!id) return false;
		return handle_operation(id.replace(/^tb-/, ''));
	});

	toolbar.find('li[title]').tipsy({opacity: 1});
	menu.find('li[title]').tipsy({opacity: 1, delayIn: 1000, gravity: 'w'});

	function handleShare(ed) {
		var fork = $(this).attr('id') == 'fork';
		if (key == '') key = ''.random(8);
		var count = 0;
		function post_save(side, text) {
			$.ajax({
				type: 'POST', async: true, dataType: 'text',
				url: '/ajax/handle_file.php',
				data: { 'key': key,  'name': side, 'content': text },
				success: function (nkey) {
					++count;
					if (count == 2) {
						var url = '/ajax/handle_save.php?key=' + key;
						if (fork) url += '&nkey=' + ''.random(8);
						$.ajax({
							type: 'GET', async: false, dataType: 'text',
							url: url,
							success: function (nkey) {
								// redirect
								if (nkey.length) window.location.href = '/' + $.trim(nkey) + '/';
							},
							error: function(xhr, ajaxOptions, thrownError){
							}
						});
					}
				},
				error: function(xhr, ajaxOptions, thrownError){
					alert(thrownError);
				}
			});
		}
		function save_files() {
			var lhs = ed.mergely('get', 'lhs');
			var rhs = ed.mergely('get', 'rhs');
			post_save('lhs', lhs);
			post_save('rhs', rhs);
		}
	
		$( '#dialog-confirm' ).dialog({
			resizable: false, width: 350, modal: true,
			buttons: {
				'Save for Sharing': function() {
					$( this ).dialog( 'close' );
					save_files();
				},
				Cancel: function() {
					$( this ).dialog( 'close' );
				}
			}
		});
	}
	
	function crossdomainGET(ed, side, url) {
		$.ajax({
			type: 'GET', dataType: 'text',
			data: {url: url},
			url: '/ajax/handle_crossdomain.php',
			contentType: 'text/plain',
			success: function (response) {
				ed.mergely(side, response);
			},
			error: function(xhr, ajaxOptions, thrownError){
				console.error('error', xhr, ajaxOptions, thrownError);
			}
		});
	}
	
	function importFiles(ed) {
		// -------------
		// file uploader - html5 file upload to browser
		function file_load(target, side) {
			var file = target.files[0];
			var reader = new FileReader();
			var $target = $(target);
			function trigger(name, event) { $target.trigger(name, event); }
			reader.onloadstart = function(evt) { trigger('start'); }
			reader.onprogress = function(evt) { trigger('progress', evt); }
			reader.onload = function(evt) { trigger('loaded', evt.target.result); }
			reader.onerror = function (evt) {
				alert(evt.target.error.name);
			}
			try {
				reader.readAsText(file, 'UTF-8');
			}
			catch (e) {
				console.error(e);
				alert(e);
			}
		}
		var file_data = {};
		$('#file-lhs, #file-rhs').change(function (evt) {
			var re = new RegExp('.*[\\\\/](.*)$');
			var match = re.exec($(this).val());
			var fname = match ? match[1] : 'unknown';
			
			var progressbar = $('#' + evt.target.id + '-progress');
			
			file_load(evt.target);
			$(evt.target).bind('start', function(ev){
				$(evt.target).css('display', 'none');
				progressbar.css('display', 'inline-block');
			});
			$(evt.target).bind('progress', function(ev, progress){
				var loaded = (progress.loaded / progress.total) * 100;
				progressbar.find('> .progress-label').text(loaded + '%');
				progressbar.progressbar('value', loaded);
			});
			$(evt.target).bind('loaded', function(ev, file){
				progressbar.progressbar('value', 100);
				progressbar.find('> .progress-label').text(fname);
				file_data[evt.target.id] = file;
			});
		});
		
		$('#file-lhs-progress').progressbar({value: 0});
		$('#file-rhs-progress').progressbar({value: 0});
		$('#dialog-upload .tabs').tabs();
		function callbackName(data) {
			console.log('callbackName', data);
		}
		$('#dialog-upload').dialog({
			dialogClass: 'no-title',
			resizable: false,
			width: '450px',
			modal: true,
			buttons: {
				Import: function() {
					$(this).dialog('close');
					
					var urls = { lhs: $('#url-lhs').val(), rhs: $('#url-rhs').val() };
					for (var side in urls) {
						var url = urls[side];
						if (!url) continue;
						crossdomainGET(ed, side, url);
					}
					
					if (file_data.hasOwnProperty('file-lhs')) {
						ed.mergely('lhs', file_data['file-lhs']);
					}
					if (file_data.hasOwnProperty('file-rhs')) {
						ed.mergely('rhs', file_data['file-rhs']);
					}
				},
				Cancel: function() {
					$(this).dialog('close');
				}
			}
		});
	}
	function colorSettings(ed) {
		// get current settings
		var sd = $('<span style="display:none" class="mergely ch d lhs start end">C</span>');
		var sa = $('<span style="display:none" class="mergely bg a rhs start end">C</span>');
		var sc = $('<span style="display:none" class="mergely c rhs start end">C</span>');
		$('body').append(sd);
		$('body').append(sa);
		$('body').append(sc);
		var conf = {
			'c-border': {id: '#c-border', getColor: function() { return sc.css('border-top-color'); }, setColor: function(color) { $('#'+this.id).val(color) }},
			'c-bg': 	{id: '#c-bg', 	  getColor: function() { return sc.css('background-color'); }, setColor: function(color) { $('#'+this.id).val(color) }},
			'a-border': {id: '#a-border', getColor: function() { return sa.css('border-top-color'); }, setColor: function(color) { $('#'+this.id).val(color) }},
			'a-bg': 	{id: '#a-bg', 	  getColor: function() { return sa.css('background-color'); }, setColor: function(color) { $('#'+this.id).val(color) }},
			'd-border': {id: '#d-border', getColor: function() { return sd.css('border-top-color'); }, setColor: function(color) { $('#'+this.id).val(color) }},
			'd-bg': 	{id: '#d-bg', 	  getColor: function() { return sd.css('background-color'); }, setColor: function(color) { $('#'+this.id).val(color) }}
		};
		$.each(conf, function(key, item){$(item.id).val(item.getColor()); });
		var f = $.farbtastic('#picker');
		$('.colorwell').each(function(){ f.linkTo(this); }).focus(function(){
			var tthis = $(this);
			f.linkTo(this);
			var item = conf[tthis.attr('id')];
			f.setColor(item.getColor());
		});
	
		$('#dialog-colors').dialog({
			width: 490, modal: true,
			buttons: {
				Apply: function() {
					var cborder = $('#c-border').val();
					var aborder = $('#a-border').val();
					var dborder = $('#d-border').val();
					var abg = $('#a-bg').val();
					var dbg = $('#d-bg').val();
					var cbg = $('#c-bg').val();
					var textCss = makeColorCss(cborder, cbg, aborder, abg, dborder, dbg);
					applyColorCss(textCss, cborder, cbg, aborder, abg, dborder, dbg, true);
				},
				Reset: function() {
				},
				Close: function() {
					$(this).dialog('close');
				}
			}
		});
	}
	function makeColorCss(changeBorder, changeBackground, addedBorder, addedBackground, deletedBorder, deletedBackground) {
		var text =
			'.mergely.a.rhs.start{border-top-color:' + addedBorder + ';}\n\
			.mergely.a.lhs.start.end,\n\
			.mergely.a.rhs.end{border-bottom-color:' + addedBorder + ';}\n\
			.mergely.a.rhs{background-color:' + addedBackground + ';}\n\
			.mergely.d.lhs{background-color:' + deletedBackground + ';}\n\
			.mergely.d.lhs.end,\n\
			.mergely.d.rhs.start.end{border-bottom-color:' + deletedBorder + ';}\n\
			.mergely.d.rhs.start.end.first{border-top-color:' + deletedBorder + ';}\n\
			.mergely.d.lhs.start{border-top-color:' + deletedBorder + ';}\n\
			.mergely.c.lhs,\n\
			.mergely.c.rhs{background-color:' + changeBackground + ';}\n\
			.mergely.c.lhs.start,\n\
			.mergely.c.rhs.start{border-top-color:' + changeBorder + ';}\n\
			.mergely.c.lhs.end,\n\
			.mergely.c.rhs.end{border-bottom-color:' + changeBorder + ';}\n\
			.mergely.ch.a.rhs{background-color:' + addedBackground + ';}\n\
			.mergely.ch.d.lhs{background-color:' + deletedBackground + ';color: #888;}';
		return text;
	}
	function applyParameterCss(saveState) {
		var cb = '#' + parameters.get('cb',color_defaults.cb),
			cg = '#' + parameters.get('cg',color_defaults.cg),
			ab = '#' + parameters.get('ab',color_defaults.ab),
			ag = '#' + parameters.get('ag',color_defaults.ag),
			db = '#' + parameters.get('db',color_defaults.db),
			dg = '#' + parameters.get('dg',color_defaults.dg);
		applyColorCss(makeColorCss(cb, cg, ab, ag, db, dg), cb, cg, ab, ag, db, dg, saveState);
	}
	function applyColorCss(cssText, changeBorder, changeBackground, addedBorder, addedBackground, deletedBorder, deletedBackground, saveState) {
		$('<style type="text/css">' + cssText + '</style>').appendTo('head');
		ed.mergely('options', {
			fgcolor:{a:addedBorder,c:changeBorder,d:deletedBorder}
		});
		var params = updateQueryStringParam('cb', changeBorder.replace(/#/g, ''), color_defaults.cb);
		params = updateQueryStringParam('cg', changeBackground.replace(/#/g, ''), color_defaults.cg, params);
		params = updateQueryStringParam('ab', addedBorder.replace(/#/g, ''), color_defaults.ab, params);
		params = updateQueryStringParam('ag', addedBackground.replace(/#/g, ''), color_defaults.ag, params);
		params = updateQueryStringParam('db', deletedBorder.replace(/#/g, ''), color_defaults.db, params);
		params = updateQueryStringParam('dg', deletedBackground.replace(/#/g, ''), color_defaults.dg, params);

		if (saveState) {
			updateHistory(params);
		}
	}
	function updateHistory(params) {
		var baseUrl = [location.protocol, '//', location.host, location.pathname].join('');
		window.history.pushState({}, null, baseUrl + params);
	}
	// Explicitly save/update a url parameter using HTML5's replaceState().
	function updateQueryStringParam(key, value, defaultValue, urlQueryString) {
		urlQueryString = urlQueryString || document.location.search;
		var parts = urlQueryString.replace(/^\?/, '').split(/&/), found = false;
		for (var i in parts) {
			if (parts[i].startsWith(key + '=')) {
				found = true;
				if (value === defaultValue) {
					// value is default value, remove option
					parts.splice(i, 1);
				}
				else {
					// make new value
					parts[i] = key + '=' + value;
				}
				break;
			}
			else if (parts[i].length === 0) {
				parts.splice(i, 1);
				break;
			}
		}
		if (!found) {
			parts.push(key + '=' + value);
		}
		return (parts.length) ? '?' + parts.join('&') : '';
	}
});
