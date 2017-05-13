// Package wkhtmltopdf contains wrappers around the wkhtmltopdf commandline tool
package wkhtmltopdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var binPath string //the cached paths as used by findPath()

// SetPath sets the path to wkhtmltopdf
func SetPath(path string) {
	binPath = path
}

// GetPath gets the path to wkhtmltopdf
func GetPath() string {
	return binPath
}

// Page is the input struct for each page
type Page struct {
	Input string
	PageOptions
}

// InputFile returns the input string and is part of the page interface
func (p *Page) InputFile() string {
	return p.Input
}

// Args returns the argument slice and is part of the page interface
func (p *Page) Args() []string {
	return p.PageOptions.Args()
}

// Reader returns the io.Reader and is part of the page interface
func (p *Page) Reader() io.Reader {
	return nil
}

// NewPage creates a new input page from a local or web resource (filepath or URL)
func NewPage(input string) *Page {
	return &Page{
		Input:       input,
		PageOptions: NewPageOptions(),
	}
}

// PageReader is one input page (a HTML document) that is read from an io.Reader
// You can add only one Page from a reader
type PageReader struct {
	Input io.Reader
	PageOptions
}

// InputFile returns the input string and is part of the page interface
func (pr *PageReader) InputFile() string {
	return "-"
}

// Args returns the argument slice and is part of the page interface
func (pr *PageReader) Args() []string {
	return pr.PageOptions.Args()
}

//Reader returns the io.Reader and is part of the page interface
func (pr *PageReader) Reader() io.Reader {
	return pr.Input
}

// NewPageReader creates a new PageReader from an io.Reader
func NewPageReader(input io.Reader) *PageReader {
	return &PageReader{
		Input:       input,
		PageOptions: NewPageOptions(),
	}
}

type page interface {
	Args() []string
	InputFile() string
	Reader() io.Reader
}

// PageOptions are options for each input page
type PageOptions struct {
	pageOptions
	headerAndFooterOptions
}

// Args returns the argument slice
func (po *PageOptions) Args() []string {
	return append(append([]string{}, po.pageOptions.Args()...), po.headerAndFooterOptions.Args()...)
}

// NewPageOptions returns a new PageOptions struct with all options
func NewPageOptions() PageOptions {
	return PageOptions{
		pageOptions:            newPageOptions(),
		headerAndFooterOptions: newHeaderAndFooterOptions(),
	}
}

// cover page
type cover struct {
	Input string
	pageOptions
}

// table of contents
type toc struct {
	Include bool
	allTocOptions
}

type allTocOptions struct {
	pageOptions
	tocOptions
}

// PDFGenerator is the main wkhtmltopdf struct, always use NewPDFGenerator to obtain a new PDFGenerator struct
type PDFGenerator struct {
	globalOptions
	outlineOptions

	Cover      cover
	TOC        toc
	OutputFile string //filename to write to, default empty (writes to internal buffer)

	binPath string
	outbuf  bytes.Buffer
	pages   []page
}

//Args returns the commandline arguments as a string slice
func (pdfg *PDFGenerator) Args() []string {
	args := []string{}
	args = append(args, pdfg.globalOptions.Args()...)
	args = append(args, pdfg.outlineOptions.Args()...)
	if pdfg.Cover.Input != "" {
		args = append(args, "cover")
		args = append(args, pdfg.Cover.Input)
		args = append(args, pdfg.Cover.pageOptions.Args()...)
	}
	if pdfg.TOC.Include {
		args = append(args, "toc")
		args = append(args, pdfg.TOC.pageOptions.Args()...)
		args = append(args, pdfg.TOC.tocOptions.Args()...)
	}
	for _, page := range pdfg.pages {
		args = append(args, "page")
		args = append(args, page.InputFile())
		args = append(args, page.Args()...)
	}
	if pdfg.OutputFile != "" {
		args = append(args, pdfg.OutputFile)
	} else {
		args = append(args, "-")
	}
	return args
}

// ArgString returns Args as a single string
func (pdfg *PDFGenerator) ArgString() string {
	return strings.Join(pdfg.Args(), " ")
}

// AddPage adds a new input page to the document.
// A page is an input HTML page, it can span multiple pages in the output document.
// It is a Page when read from file or URL or a PageReader when read from memory.
func (pdfg *PDFGenerator) AddPage(p page) {
	pdfg.pages = append(pdfg.pages, p)
}

// SetPages resets all pages
func (pdfg *PDFGenerator) SetPages(p []page) {
	pdfg.pages = p
}

// Buffer returns the embedded output buffer used if OutputFile is empty
func (pdfg *PDFGenerator) Buffer() *bytes.Buffer {
	return &pdfg.outbuf
}

// Bytes returns the output byte slice from the output buffer used if OutputFile is empty
func (pdfg *PDFGenerator) Bytes() []byte {
	return pdfg.outbuf.Bytes()
}

// WriteFile writes the contents of the output buffer to a file
func (pdfg *PDFGenerator) WriteFile(filename string) error {
	return ioutil.WriteFile(filename, pdfg.Bytes(), 0666)
}

//findPath finds the path to wkhtmltopdf by
//- first looking in the current dir
//- looking in the PATH and PATHEXT environment dirs
//- using the WKHTMLTOPDF_PATH environment dir
//The path is cached, meaning you can not change the location of wkhtmltopdf in
//a running program once it has been found
func (pdfg *PDFGenerator) findPath() error {
	const exe = "wkhtmltopdf"
	if binPath != "" {
		pdfg.binPath = binPath
		return nil
	}
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	path, err := exec.LookPath(filepath.Join(exeDir, exe))
	if err == nil && path != "" {
		binPath = path
		pdfg.binPath = path
		return nil
	}
	path, err = exec.LookPath(exe)
	if err == nil && path != "" {
		binPath = path
		pdfg.binPath = path
		return nil
	}
	dir := os.Getenv("WKHTMLTOPDF_PATH")
	if dir == "" {
		return fmt.Errorf("%s not found", exe)
	}
	path, err = exec.LookPath(filepath.Join(dir, exe))
	if err == nil && path != "" {
		binPath = path
		pdfg.binPath = path
		return nil
	}
	return fmt.Errorf("%s not found", exe)
}

// Create creates the PDF document and stores it in the internal buffer if no error is returned
func (pdfg *PDFGenerator) Create() error {
	return pdfg.run()
}

func (pdfg *PDFGenerator) run() error {

	errbuf := &bytes.Buffer{}

	cmd := exec.Command(pdfg.binPath, pdfg.Args()...)

	cmd.Stdout = &pdfg.outbuf
	cmd.Stderr = errbuf
	//if there is a pageReader page (from Stdin) we set Stdin to that reader
	for _, page := range pdfg.pages {
		if page.Reader() != nil {
			cmd.Stdin = page.Reader()
			break
		}
	}

	err := cmd.Run()
	if err != nil {
		errStr := errbuf.String()
		if strings.TrimSpace(errStr) == "" {
			errStr = err.Error()
		}
		return errors.New(errStr)
	}
	return nil
}

// NewPDFGenerator returns a new PDFGenerator struct with all options created and
// checks if wkhtmltopdf can be found on the system
func NewPDFGenerator() (*PDFGenerator, error) {
	pdfg := &PDFGenerator{
		globalOptions:  newGlobalOptions(),
		outlineOptions: newOutlineOptions(),
		Cover: cover{
			pageOptions: newPageOptions(),
		},
		TOC: toc{
			allTocOptions: allTocOptions{
				tocOptions:  newTocOptions(),
				pageOptions: newPageOptions(),
			},
		},
	}
	err := pdfg.findPath()
	return pdfg, err
}
