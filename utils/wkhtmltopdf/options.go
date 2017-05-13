package wkhtmltopdf

import (
	"fmt"
	"reflect"
)

//A list of options that can be set from code to make it easier to see which options are available
type globalOptions struct {
	CookieJar         stringOption //Read and write cookies from and to the supplied cookie jar file
	Copies            uintOption   //Number of copies to print into the pdf file (default 1)
	Dpi               uintOption   //Change the dpi explicitly (this has no effect on X11 based systems)
	ExtendedHelp      boolOption   //Display more extensive help, detailing less common command switches
	Grayscale         boolOption   //PDF will be generated in grayscale
	Help              boolOption   //Display help
	HTMLDoc           boolOption   //Output program html help
	ImageDpi          uintOption   //When embedding images scale them down to this dpi (default 600)
	ImageQuality      uintOption   //When jpeg compressing images use this quality (default 94)
	License           boolOption   //Output license information and exit
	Lowquality        boolOption   //Generates lower quality pdf/ps. Useful to shrink the result document space
	ManPage           boolOption   //Output program man page
	MarginBottom      uintOption   //Set the page bottom margin
	MarginLeft        uintOption   //Set the page left margin (default 10mm)
	MarginRight       uintOption   //Set the page right margin (default 10mm)
	MarginTop         uintOption   //Set the page top margin
	Orientation       stringOption // Set orientation to Landscape or Portrait (default Portrait)
	NoCollate         boolOption   //Do not collate when printing multiple copies (default collate)
	PageHeight        uintOption   //Page height
	PageSize          stringOption //Set paper size to: A4, Letter, etc. (default A4)
	PageWidth         uintOption   //Page width
	NoPdfCompression  boolOption   //Do not use lossless compression on pdf objects
	Quiet             boolOption   //Be less verbose
	ReadArgsFromStdin boolOption   //Read command line arguments from stdin
	Readme            boolOption   //Output program readme
	Title             stringOption //The title of the generated pdf file (The title of the first document is used if not specified)
	Version           boolOption   //Output version information and exit
}

func (gopt *globalOptions) Args() []string {
	return optsToArgs(gopt)
}

type outlineOptions struct {
	DumpDefaultTocXsl boolOption   //Dump the default TOC xsl style sheet to stdout
	DumpOutline       stringOption //Dump the outline to a file
	NoOutline         boolOption   //Do not put an outline into the pdf
	OutlineDepth      uintOption   //Set the depth of the outline (default 4)
}

func (oopt *outlineOptions) Args() []string {
	return optsToArgs(oopt)
}

type pageOptions struct {
	Allow                     sliceOption  //Allow the file or files from the specified folder to be loaded (repeatable)
	NoBackground              boolOption   //Do not print background
	CacheDir                  stringOption //Web cache directory
	CheckboxCheckedSvg        stringOption //Use this SVG file when rendering checked checkboxes
	CheckboxSvg               stringOption //Use this SVG file when rendering unchecked checkboxes
	Cookie                    mapOption    //Set an additional cookie (repeatable), value should be url encoded
	CustomHeader              mapOption    //Set an additional HTTP header (repeatable)
	CustomHeaderPropagation   boolOption   //Add HTTP headers specified by --custom-header for each resource request
	NoCustomHeaderPropagation boolOption   //Do not add HTTP headers specified by --custom-header for each resource request
	DebugJavascript           boolOption   //Show javascript debugging output
	DefaultHeader             boolOption   //Add a default header, with the name of the page to the left, and the page number to the right, this is short for: --header-left='[webpage]' --header-right='[page]/[toPage]' --top 2cm --header-line
	Encoding                  stringOption //Set the default text encoding, for input
	DisableExternalLinks      boolOption   //Do not make links to remote web pages
	EnableForms               boolOption   //Turn HTML form fields into pdf form fields
	NoImages                  boolOption   //Do not load or print images
	DisableInternalLinks      boolOption   //Do not make local links
	DisableJavascript         boolOption   //Do not allow web pages to run javascript
	JavascriptDelay           uintOption   //Wait some milliseconds for javascript finish (default 200)
	LoadErrorHandling         stringOption //Specify how to handle pages that fail to load: abort, ignore or skip (default abort)
	LoadMediaErrorHandling    stringOption //Specify how to handle media files that fail to load: abort, ignore or skip (default ignore)
	DisableLocalFileAccess    boolOption   //Do not allowed conversion of a local file to read in other local files, unless explicitly allowed with --allow
	MinimumFontSize           uintOption   //Minimum font size
	ExcludeFromOutline        boolOption   //Do not include the page in the table of contents and outlines
	PageOffset                uintOption   //Set the starting page number (default 0)
	Password                  stringOption //HTTP Authentication password
	EnablePlugins             boolOption   //Enable installed plugins (plugins will likely not work)
	Post                      mapOption    //Add an additional post field (repeatable)
	PostFile                  mapOption    //Post an additional file (repeatable)
	PrintMediaType            boolOption   //Use print media-type instead of screen
	Proxy                     stringOption //Use a proxy
	RadiobuttonCheckedSvg     stringOption //Use this SVG file when rendering checked radiobuttons
	RadiobuttonSvg            stringOption //Use this SVG file when rendering unchecked radiobuttons
	RunScript                 sliceOption  //Run this additional javascript after the page is done loading (repeatable)
	DisableSmartShrinking     boolOption   //Disable the intelligent shrinking strategy used by WebKit that makes the pixel/dpi ratio none constant
	NoStopSlowScripts         boolOption   //Do not Stop slow running javascripts
	EnableTocBackLinks        boolOption   //Link from section header to toc
	UserStyleSheet            stringOption //Specify a user style sheet, to load with every page
	Username                  stringOption //HTTP Authentication username
	ViewportSize              stringOption //Set viewport size if you have custom scrollbars or css attribute overflow to emulate window size
	WindowStatus              stringOption //Wait until window.status is equal to this string before rendering page
	Zoom                      floatOption  //Use this zoom factor (default 1)
}

func (popt *pageOptions) Args() []string {
	return optsToArgs(popt)
}

type headerAndFooterOptions struct {
	FooterCenter   stringOption //Centered footer text
	FooterFontName stringOption //Set footer font name (default Arial)
	FooterFontSize uintOption   //Set footer font size (default 12)
	FooterHTML     stringOption //Adds a html footer
	FooterLeft     stringOption //Left aligned footer text
	FooterLine     boolOption   //Display line above the footer
	FooterRight    stringOption //Right aligned footer text
	FooterSpacing  floatOption  //Spacing between footer and content in mm (default 0)
	HeaderCenter   stringOption //Centered header text
	HeaderFontName stringOption //Set header font name (default Arial)
	HeaderFontSize uintOption   //Set header font size (default 12)
	HeaderHTML     stringOption //Adds a html header
	HeaderLeft     stringOption //Left aligned header text
	HeaderLine     boolOption   //Display line below the header
	HeaderRight    stringOption //Right aligned header text
	HeaderSpacing  floatOption  //Spacing between header and content in mm (default 0)
	Replace        mapOption    //Replace [name] with value in header and footer (repeatable)
}

func (hopt *headerAndFooterOptions) Args() []string {
	return optsToArgs(hopt)
}

type tocOptions struct {
	DisableDottedLines  boolOption   //Do not use dotted lines in the toc
	TocHeaderText       stringOption //The header text of the toc (default Table of Contents)
	TocLevelIndentation uintOption   //For each level of headings in the toc indent by this length (default 1em)
	DisableTocLinks     boolOption   //Do not link from toc to sections
	TocTextSizeShrink   floatOption  //For each level of headings in the toc the font is scaled by this factor
	XslStyleSheet       stringOption //Use the supplied xsl style sheet for printing the table of content
}

func (topt *tocOptions) Args() []string {
	return optsToArgs(topt)
}

type argParser interface {
	Parse() []string //Used in the cmd call
}

type stringOption struct {
	option string
	value  string
}

func (so stringOption) Parse() []string {
	args := []string{}
	if so.value == "" {
		return args
	}
	args = append(args, "--"+so.option)
	args = append(args, so.value)
	return args
}

func (so *stringOption) Set(value string) {
	so.value = value
}

type sliceOption struct {
	option string
	value  []string
}

func (so sliceOption) Parse() []string {
	args := []string{}
	if len(so.value) == 0 {
		return args
	}
	for _, v := range so.value {
		args = append(args, "--"+so.option)
		args = append(args, v)
	}
	return args
}

func (so *sliceOption) Set(value string) {
	so.value = append(so.value, value)
}

type mapOption struct {
	option string
	value  map[string]string
}

func (mo mapOption) Parse() []string {
	args := []string{}
	if mo.value == nil || len(mo.value) == 0 {
		return args
	}
	for k, v := range mo.value {
		args = append(args, "--"+mo.option)
		args = append(args, k)
		args = append(args, v)
	}
	return args
}

func (mo *mapOption) Set(key, value string) {
	if mo.value == nil {
		mo.value = make(map[string]string)
	}
	mo.value[key] = value
}

type uintOption struct {
	option string
	value  uint
	isSet  bool
}

func (io uintOption) Parse() []string {
	args := []string{}
	if io.isSet == false {
		return args
	}
	args = append(args, "--"+io.option)
	args = append(args, fmt.Sprintf("%d", io.value))
	return args
}

func (io *uintOption) Set(value uint) {
	io.isSet = true
	io.value = value
}

type floatOption struct {
	option string
	value  float64
	isSet  bool
}

func (fo floatOption) Parse() []string {
	args := []string{}
	if fo.isSet == false {
		return args
	}
	args = append(args, "--"+fo.option)
	args = append(args, fmt.Sprintf("%.3f", fo.value))
	return args
}

func (fo *floatOption) Set(value float64) {
	fo.isSet = true
	fo.value = value
}

type boolOption struct {
	option string
	value  bool
}

func (bo boolOption) Parse() []string {
	if bo.value {
		return []string{"--" + bo.option}
	}
	return []string{}
}

func (bo *boolOption) Set(value bool) {
	bo.value = value
}

func newGlobalOptions() globalOptions {
	return globalOptions{
		CookieJar:         stringOption{option: "cookie-jar"},
		Copies:            uintOption{option: "copies"},
		Dpi:               uintOption{option: "dpi"},
		ExtendedHelp:      boolOption{option: "extended-help"},
		Grayscale:         boolOption{option: "grayscale"},
		Help:              boolOption{option: "true"},
		HTMLDoc:           boolOption{option: "htmldoc"},
		ImageDpi:          uintOption{option: "image-dpi"},
		ImageQuality:      uintOption{option: "image-quality"},
		License:           boolOption{option: "license"},
		Lowquality:        boolOption{option: "lowquality"},
		ManPage:           boolOption{option: "manpage"},
		MarginBottom:      uintOption{option: "margin-bottom"},
		MarginLeft:        uintOption{option: "margin-left"},
		MarginRight:       uintOption{option: "margin-right"},
		MarginTop:         uintOption{option: "margin-top"},
		Orientation:       stringOption{option: "orientation"},
		NoCollate:         boolOption{option: "nocollate"},
		PageHeight:        uintOption{option: "page-height"},
		PageSize:          stringOption{option: "page-size"},
		PageWidth:         uintOption{option: "page-width"},
		NoPdfCompression:  boolOption{option: "no-pdf-compression"},
		Quiet:             boolOption{option: "quiet"},
		ReadArgsFromStdin: boolOption{option: "read-args-from-stdin"},
		Readme:            boolOption{option: "readme"},
		Title:             stringOption{option: "title"},
		Version:           boolOption{option: "version"},
	}
}

func newOutlineOptions() outlineOptions {
	return outlineOptions{
		DumpDefaultTocXsl: boolOption{option: "dump-default-toc-xsl"},
		DumpOutline:       stringOption{option: "dump-outline"},
		NoOutline:         boolOption{option: "no-outline"},
		OutlineDepth:      uintOption{option: "outline-depth"},
	}
}

func newPageOptions() pageOptions {
	return pageOptions{
		Allow:                     sliceOption{option: "allow"},
		NoBackground:              boolOption{option: "no-background"},
		CacheDir:                  stringOption{option: "cache-dir"},
		CheckboxCheckedSvg:        stringOption{option: "checkbox-checked-svg"},
		CheckboxSvg:               stringOption{option: "checkbox-svg"},
		Cookie:                    mapOption{option: "cookie"},
		CustomHeader:              mapOption{option: "custom-header"},
		CustomHeaderPropagation:   boolOption{option: "custom-header-propagation"},
		NoCustomHeaderPropagation: boolOption{option: "no-custom-header-propagation"},
		DebugJavascript:           boolOption{option: "debug-javascript"},
		DefaultHeader:             boolOption{option: "default-header"},
		Encoding:                  stringOption{option: "encoding"},
		DisableExternalLinks:      boolOption{option: "disable-external-links"},
		EnableForms:               boolOption{option: "enable-forms"},
		NoImages:                  boolOption{option: "no-images"},
		DisableInternalLinks:      boolOption{option: "disable-internal-links"},
		DisableJavascript:         boolOption{option: "disable-javascript "},
		JavascriptDelay:           uintOption{option: "javascript-delay"},
		LoadErrorHandling:         stringOption{option: "load-error-handling"},
		LoadMediaErrorHandling:    stringOption{option: "load-media-error-handling"},
		DisableLocalFileAccess:    boolOption{option: "disable-local-file-access"},
		MinimumFontSize:           uintOption{option: "minimum-font-size"},
		ExcludeFromOutline:        boolOption{option: "exclude-from-outline"},
		PageOffset:                uintOption{option: "page-offset"},
		Password:                  stringOption{option: "password"},
		EnablePlugins:             boolOption{option: "enable-plugins"},
		Post:                      mapOption{option: "post"},
		PostFile:                  mapOption{option: "post-file"},
		PrintMediaType:            boolOption{option: "print-media-type"},
		Proxy:                     stringOption{option: "proxy"},
		RadiobuttonCheckedSvg: stringOption{option: "radiobutton-checked-svg"},
		RadiobuttonSvg:        stringOption{option: "radiobutton-svg"},
		RunScript:             sliceOption{option: "run-script"},
		DisableSmartShrinking: boolOption{option: "disable-smart-shrinking"},
		NoStopSlowScripts:     boolOption{option: "no-stop-slow-scripts"},
		EnableTocBackLinks:    boolOption{option: "enable-toc-back-links"},
		UserStyleSheet:        stringOption{option: "user-style-sheet"},
		Username:              stringOption{option: "username"},
		ViewportSize:          stringOption{option: "viewport-size"},
		WindowStatus:          stringOption{option: "window-status"},
		Zoom:                  floatOption{option: "zoom"},
	}
}

func newHeaderAndFooterOptions() headerAndFooterOptions {
	return headerAndFooterOptions{
		FooterCenter:   stringOption{option: "footer-center"},
		FooterFontName: stringOption{option: "footer-font-name"},
		FooterFontSize: uintOption{option: "footer-font-size"},
		FooterHTML:     stringOption{option: "footer-html"},
		FooterLeft:     stringOption{option: "footer-left"},
		FooterLine:     boolOption{option: "footer-line"},
		FooterRight:    stringOption{option: "footer-right"},
		FooterSpacing:  floatOption{option: "footer-spacing"},
		HeaderCenter:   stringOption{option: "header-center"},
		HeaderFontName: stringOption{option: "header-font-name"},
		HeaderFontSize: uintOption{option: "header-font-size"},
		HeaderHTML:     stringOption{option: "header-html"},
		HeaderLeft:     stringOption{option: "header-left"},
		HeaderLine:     boolOption{option: "header-line"},
		HeaderRight:    stringOption{option: "header-right"},
		HeaderSpacing:  floatOption{option: "header-spacing"},
		Replace:        mapOption{option: "replace"},
	}
}

func newTocOptions() tocOptions {
	return tocOptions{
		DisableDottedLines:  boolOption{option: "disable-dotted-lines"},
		TocHeaderText:       stringOption{option: "toc-header-text"},
		TocLevelIndentation: uintOption{option: "toc-level-indentation"},
		DisableTocLinks:     boolOption{option: "disable-toc-links"},
		TocTextSizeShrink:   floatOption{option: "toc-text-size-shrink"},
		XslStyleSheet:       stringOption{option: "xsl-style-sheet"},
	}
}

func optsToArgs(opts interface{}) []string {
	args := []string{}
	rv := reflect.Indirect(reflect.ValueOf(opts))
	if rv.Kind() != reflect.Struct {
		return args
	}
	for i := 0; i < rv.NumField(); i++ {
		prsr, ok := rv.Field(i).Interface().(argParser)
		if ok {
			s := prsr.Parse()
			if len(s) > 0 {
				args = append(args, s...)
			}
		}
	}
	return args
}

// Constants for orientation modes
const (
	OrientationLandscape = "Landscape" // Landscape mode
	OrientationPortrait  = "Portrait"  // Portrait mode
)

// Constants for page sizes
const (
	PageSizeA0        = "A0"        //	841 x 1189 mm
	PageSizeA1        = "A1"        //	594 x 841 mm
	PageSizeA2        = "A2"        //	420 x 594 mm
	PageSizeA3        = "A3"        //	297 x 420 mm
	PageSizeA4        = "A4"        //	210 x 297 mm, 8.26
	PageSizeA5        = "A5"        //	148 x 210 mm
	PageSizeA6        = "A6"        //	105 x 148 mm
	PageSizeA7        = "A7"        //	74 x 105 mm
	PageSizeA8        = "A8"        //	52 x 74 mm
	PageSizeA9        = "A9"        //	37 x 52 mm
	PageSizeB0        = "B0"        //	1000 x 1414 mm
	PageSizeB1        = "B1"        //	707 x 1000 mm
	PageSizeB2        = "B2"        //	500 x 707 mm
	PageSizeB3        = "B3"        //	353 x 500 mm
	PageSizeB4        = "B4"        //	250 x 353 mm
	PageSizeB5        = "B5"        //	176 x 250 mm, 6.93
	PageSizeB6        = "B6"        //	125 x 176 mm
	PageSizeB7        = "B7"        //	88 x 125 mm
	PageSizeB8        = "B8"        //	62 x 88 mm
	PageSizeB9        = "B9"        //	33 x 62 mm
	PageSizeB10       = "B10"       //	31 x 44 mm
	PageSizeC5E       = "C5E"       //	163 x 229 mm
	PageSizeComm10E   = "Comm10E"   //	105 x 241 mm, U.S. Common 10 Envelope
	PageSizeDLE       = "DLE"       //	110 x 220 mm
	PageSizeExecutive = "Executive" //	7.5 x 10 inches, 190.5 x 254 mm
	PageSizeFolio     = "Folio"     //	210 x 330 mm
	PageSizeLedger    = "Ledger"    //	431.8 x 279.4 mm
	PageSizeLegal     = "Legal"     //	8.5 x 14 inches, 215.9 x 355.6 mm
	PageSizeLetter    = "Letter"    //	8.5 x 11 inches, 215.9 x 279.4 mm
	PageSizeTabloid   = "Tabloid"   //	279.4 x 431.8 mm
	PageSizeCustom    = "Custom"    //	Unknown, or a user defined size.
)
