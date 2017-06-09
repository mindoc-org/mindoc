<?php
$key = '';
$debug = False;
if (isset($_GET['key'])) {
	$key = $_GET['key'];
}
if (isset($_GET['debug'])) {
    $debug = filter_var($_GET['debug'], FILTER_VALIDATE_BOOLEAN);
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" /><title>Mergely - Diff online, merge documents</title>
	<meta http-equiv="content-type" content="text/html; charset=UTF-8" />
	<meta name="description" content="Merge and Diff your documents with diff online and share" />
	<meta name="keywords" content="diff,merge,compare,jsdiff,comparison,difference,file,text,unix,patch,algorithm,saas,longest common subsequence,diff online" />
	<meta name="author" content="Jamie Peabody" />
	<link rel="shortcut icon" href="/favicon.ico" />
	<link rel="canonical" href="http://www.mergely.com" />
    <link href='http://fonts.googleapis.com/css?family=Noto+Sans:400,700' rel='stylesheet' type='text/css' />
    <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.9.0/jquery.min.js"></script>
    <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jqueryui/1.10.1/jquery-ui.min.js"></script>

	<link type="text/css" rel="stylesheet" href="/style/mergely-theme/jquery-ui-1.10.1.custom.css" />
    <link type='text/css' rel='stylesheet' href='/Mergely/editor/lib/wicked-ui.css' />
	<script type="text/javascript" src="/Mergely/editor/lib/wicked-ui.js"></script>

    <link type='text/css' rel='stylesheet' href='/Mergely/editor/lib/tipsy/tipsy.css' />
	<script type="text/javascript" src="/Mergely/editor/lib/tipsy/jquery.tipsy.js"></script>
	<script type="text/javascript" src="/Mergely/editor/lib/farbtastic/farbtastic.js"></script>
	<link type="text/css" rel="stylesheet" href="/Mergely/editor/lib/farbtastic/farbtastic.css" />
<?php
    if ($debug) {
?>
    <script type="text/javascript" src="/Mergely/lib/codemirror.js"></script>
    <script type="text/javascript" src="/Mergely/lib/mergely.js"></script>
    <script type="text/javascript" src="/Mergely/editor/editor.js"></script>
<?php
    }
    else {
?>
    <script type="text/javascript" src="/Mergely/lib/codemirror.min.js"></script>
    <script type="text/javascript" src="/Mergely/lib/mergely.min.js"></script>
    <script type="text/javascript" src="/Mergely/editor/editor.min.js"></script>
<?php
    }
?>
    <link type="text/css" rel="stylesheet" href="/Mergely/lib/codemirror.css" />
	<link type="text/css" rel="stylesheet" href="/Mergely/lib/mergely.css" />
    <link type='text/css' rel='stylesheet' href='/Mergely/editor/editor.css' />
	<script type="text/javascript" src="/Mergely/lib/searchcursor.js"></script>

	<script type="text/javascript">
        var key = '<?php echo htmlspecialchars($key, ENT_QUOTES, 'UTF-8'); ?>';
        var isSample = key == 'usaindep';
    </script>
    
    <!-- analytics -->
    <script type="text/javascript">
        var _gaq = _gaq || [];
        _gaq.push(['_setAccount', 'UA-85576-5']);
        _gaq.push(['_trackPageview']);
        (function() {
            var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
            ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
            var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
        })();
    </script>
    
    <!-- google +1 -->
	<script type="text/javascript" src="https://apis.google.com/js/plusone.js"></script>
</head>
<body style="visibility:hidden">
<div id="fb-root"></div><script>(function(d, s, id) {
  var js, fjs = d.getElementsByTagName(s)[0];
  if (d.getElementById(id)) return;
  js = d.createElement(s); js.id = id;
  js.src = "//connect.facebook.net/en_US/all.js#xfbml=1";
  fjs.parentNode.insertBefore(js, fjs);
}(document, 'script', 'facebook-jssdk'));</script>


    <a href="/"><div id="banner"></div></a>
    
    <!-- menu -->
    <ul id="main-menu">
        <li accesskey="f">
            File
            <ul>
                <li id="file-new" accesskey="n" data-hotkey="Alt+N">New</li>
                <li id="file-import" data-icon="icon-import">Import...</li>
                <li id="file-save" accesskey="s" data-hotkey="Alt+S" data-icon="icon-save">Save .diff</li>
                <li class="separator"></li>
                <li id="file-share" data-icon="icon-share">Share</li>
            </ul>
        </li>
        <li accesskey="l">
            Left
            <ul>
                <li id="edit-left-undo" accesskey="z" data-hotkey="Ctrl+Z" data-icon="icon-undo">Undo</li>
                <li id="edit-left-redo" accesskey="y" data-hotkey="Ctrl+Y" data-icon="icon-redo">Redo</li>
                <li id="edit-left-find">Find</li>
                <li class="separator"></li>
                <li id="edit-left-merge-right" data-hotkey="Alt+&rarr;" data-icon="icon-arrow-right-v">Merge change right</li>
                <li id="edit-left-merge-right-file" data-icon="icon-arrow-right-vv">Merge file right</li>
                <li id="edit-left-readonly">Read only</li>
                <li class="separator"></li>
                <li id="edit-left-clear">Clear</li>
            </ul>
        </li>
        <li accesskey="r">
            Right
            <ul>
                <li id="edit-right-undo" accesskey="z" data-hotkey="Ctrl+Z" data-icon="icon-undo">Undo</li>
                <li id="edit-right-redo" accesskey="y" data-hotkey="Ctrl+Y" data-icon="icon-redo">Redo</li>
                <li id="edit-right-find">Find</li>
                <li class="separator"></li>
                <li id="edit-right-merge-left" data-hotkey="Alt+&larr;" data-icon="icon-arrow-left-v">Merge change left</li>
                <li id="edit-right-merge-left-file" data-icon="icon-arrow-left-vv">Merge file left</li>
                <li id="edit-right-readonly">Read only</li>
                <li class="separator"></li>
                <li id="edit-right-clear">Clear</li>
            </ul>
        </li>
        <li accesskey="v">
            View
            <ul>
                <li id="view-swap" data-icon="icon-swap">Swap sides</li>
                <li class="separator"></li>
                <li id="view-refresh" accesskey="v" data-hotkey="Alt+V" title="Generates diff markup">Render diff view</li>
                <li id="view-clear" accesskey="c" data-hotkey="Alt+C" title="Clears diff markup">Clear render</li>
                <li class="separator"></li>
                <li id="view-change-prev" data-hotkey="Alt+&uarr;" title="View previous change">View prev change</li>
                <li id="view-change-next" data-hotkey="Alt+&darr;" title="View next change">View next change</li>
            </ul>
        </li>
        <li accesskey="o">
            Options
            <ul>
                <li id="options-wrap">Wrap lines</li>
                <li id="options-ignorews">Ignore white space</li>
                <li id="options-ignorecase">Ignore case</li>
                <li class="separator"></li>
                <li id="options-viewport" title="Improves performance for large files">Enable viewport</li>
                <li id="options-sidebars" title="Improves performance for large files">Enable side bars</li>
                <li id="options-swapmargin">Swap right margin</li>
                <li id="options-linenumbers">Enable line numbers</li>
                <li class="separator"></li>
                <li id="options-autodiff" title="Diffs are computed automatically">Enable auto-diff</li>
                <li class="separator"></li>
                <li id="options-colors">Colors...</li>
            </ul>
        </li>
        <li accesskey="m">
            Mergely
            <ul>
                <li><a class="link" href="/" target="site">Home</a></li>
                <li><a class="link" href="/about" target="site">About</a></li>
                <li><a class="link" href="/license" target="site">License</a></li>
                <li><a class="link" href="/download" target="site">Download</a></li>
                <li><a class="link" href="/doc" target="site">Mergely development guide</a></li>
                <li class="separator"></li>
                <li><a class="link" href="/united-states-declaration-of-independence?wl=1" target="_blank">United States Declaration of Independence Draft</a></li>
            </ul>
        </li>
<?php
    if (!$debug) {
?>
        <li accesskey="s">
            Social
            <ul>
                <li id="social-twitter">
                    <div style="padding: 10px 10px 5px 10px" title="Twitter">
                        <a href="https://twitter.com/share" class="twitter-share-button" data-via="jamiepeabody">Tweet</a>
                        <script>!function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0],p=/^http:/.test(d.location)?'http':'https';if(!d.getElementById(id)){js=d.createElement(s);js.id=id;js.src=p+'://platform.twitter.com/widgets.js';fjs.parentNode.insertBefore(js,fjs);}}(document, 'script', 'twitter-wjs');</script>
                    </div>
                </li>
                <li id="social-facebook">
                    <div style="padding: 10px 10px 5px 10px" title="Facebook">
                        <div class="fb-like" data-href="http://www.mergely.com" data-send="true" data-width="200" data-show-faces="true"></div>
                    </div>
                </li>
                <li id="social-google">
                    <div style="padding: 10px 10px 5px 10px" title="Google+"><g:plusone></g:plusone></div>
                </li>
                <li id="social-reddit">
                    <div style="padding: 10px 10px 5px 10px" title="Reddit">
                        <a target="_blank" href="http://www.reddit.com/submit" onclick="window.location = 'http://www.reddit.com/submit?url=' + encodeURIComponent(window.location); return false" style="color:black;text-decoration:none;"><img src="http://www.reddit.com/static/spreddit1.gif" alt="submit to reddit" border="0" />
                            <span>Reddit</span>
                        </a>
                    </div>
                </li>
            </ul>
        </li>
<?php
    }
?>
    </ul>

    <!-- toolbar -->
    <ul id="toolbar">
        <li id="tb-file-share" data-icon="icon-share" title="Share">Share</li>
        <li class="separator"></li>
        <li id="tb-file-import" data-icon="icon-import" title="Import">Import</li>
        <li id="tb-file-save" data-icon="icon-save" title="Save .diff">Save .diff</li>
        <li class="separator"></li>
        <li id="tb-view-change-prev" data-icon="icon-arrow-up" title="Previous change">Previous change</li>
        <li id="tb-view-change-next" data-icon="icon-arrow-down" title="Next change">Next change</li>
        <li class="separator"></li>
        <li id="tb-edit-right-merge-left" data-icon="icon-arrow-left-v" title="Merge change left">Merge change left</li>
        <li id="tb-edit-left-merge-right" data-icon="icon-arrow-right-v" title="Merge change right">Merge change right</li>
        <li id="tb-view-swap" data-icon="icon-swap" title="Swap sides">Swap sides</li>
    </ul>

    <!-- dialog upload -->
    <div id="dialog-upload" title="File import" style="display:none">
        <div class="tabs">
            <ul>
                <li><a href="#tabs-1">Import File</a></li>
                <li><a href="#tabs-2">Import URL</a></li>
            </ul>
            <div id="tabs-1">
                <p>
                    Files are imported directly into your browser.  They are <em>not</em> uploaded to the server.
                </p>
                <label for="file-lhs">Left file</label> <input id="file-lhs" style="display:inline-block" type="file"><div id="file-lhs-progress"><div class="progress-label">Loading...</div></div><br />
                <label for="file-rhs">Right file</label> <input id="file-rhs" style="display:inline-block" type="file"><div id="file-rhs-progress"><div class="progress-label">Loading...</div></div><br />
            </div>
            <div id="tabs-2">
                <p>
                    Files are imported directly into your browser.  They are <em>not</em> uploaded to the server.
                </p>
                <label for="url-lhs">Left URL</label> <input id="url-lhs" type="input" size="40"><div id="file-lhs-progress"><div class="progress-label">Loading...</div></div><br />
                <label for="url-rhs">Right URL</label> <input id="url-rhs" type="input" size="40"><div id="file-rhs-progress"><div class="progress-label">Loading...</div></div><br />
            </div>
        </div>
    </div>
    
    <!-- dialog colors -->
	<div id="dialog-colors" title="Mergely Color Settings" style="display:none">
		<div id="picker" style="float: right;"></div>
		<fieldset>
			<legend>Changed</legend>
			<label for="c-border">Border:</label><input type="text" id="c-border" name="c-border" class="colorwell" />
			<br />
			<label for="c-bg">Background:</label><input type="text" id="c-bg" name="c-bg" class="colorwell" />
			<br />
		</fieldset>
		<fieldset>
			<legend>Added</legend>
			<label for="a-border">Border:</label><input type="text" id="a-border" name="a-border" class="colorwell" />
			<br />
			<label for="a-bg">Background:</label><input type="text" id="a-bg" name="a-bg" class="colorwell" />
			<br />
		</fieldset>
		<fieldset>
			<legend>Deleted</legend>
			<label for="d-border">Border:</label><input type="text" id="d-border" name="d-border" class="colorwell" />
			<br />
			<label for="d-bg">Background:</label><input type="text" id="d-bg" name="d-bg" class="colorwell" />
			<br />
		</fieldset>
	</div>

    <!-- dialog confirm -->
	<div id="dialog-confirm" title="Save a Permanent Copy?" style="display:none;">
		<p>
			Are you sure you want to save? A permanent copy will be
			created at the server and a link will be provided for you to share the URL 
            in an email, blog, twitter, etc.
		</p>
	</div>
    
    <!-- find -->
    <div class="find">
        <input type="text" />
        <button class="find-prev"><span class="icon icon-arrow-up"></span></button>
        <button class="find-next"><span class="icon icon-arrow-down"></span></button>
        <button class="find-close"><span class="icon icon-x-mark"></span></button>
    </div>
    
    <!-- editor -->
    <div style="position: absolute;top: 73px;bottom: 10px;left: 5px;right: 5px;overflow-y: hidden;padding-bottom: 2px;">
        <div id="mergely">
        </div>
    </div>
    
</body>
</html>
