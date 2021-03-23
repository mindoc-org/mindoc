//Author:TruthHun
//Email:TruthHun@QQ.COM
//Date:2018-01-21
package converter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"errors"
	"os/exec"
	"time"

	"html"
	"sync"

	"github.com/mindoc-org/mindoc/utils/cryptil"
	"github.com/mindoc-org/mindoc/utils/filetil"
	"github.com/mindoc-org/mindoc/utils/ziptil"
)

type Converter struct {
	BasePath       string
	OutputPath	   string
	Config         Config
	Debug          bool
	GeneratedCover string
	ProcessNum	 	int		//并发的任务数量
	process 	chan func()
	limitChan		chan bool
}

//目录结构
type Toc struct {
	Id    int    `json:"id"`
	Link  string `json:"link"`
	Pid   int    `json:"pid"`
	Title string `json:"title"`
}

//config.json文件解析结构
type Config struct {
	Charset      string   `json:"charset"`       //字符编码，默认utf-8编码
	Cover        string   `json:"cover"`         //封面图片，或者封面html文件
	Timestamp    string   `json:"date"`          //时间日期,如“2018-01-01 12:12:21”，其实是time.Time格式，但是直接用string就好
	Description  string   `json:"description"`   //摘要
	Footer       string   `json:"footer"`        //pdf的footer
	Header       string   `json:"header"`        //pdf的header
	Identifier   string   `json:"identifier"`    //即uuid，留空即可
	Language     string   `json:"language"`      //语言，如zh、en、zh-CN、en-US等
	Creator      string   `json:"creator"`       //作者，即author
	Publisher    string   `json:"publisher"`     //出版单位
	Contributor  string   `json:"contributor"`   //同Publisher
	Title        string   `json:"title"`         //文档标题
	Format       []string `json:"format"`        //导出格式，可选值：pdf、epub、mobi
	FontSize     string   `json:"font_size"`     //默认的pdf导出字体大小
	PaperSize    string   `json:"paper_size"`    //页面大小
	MarginLeft   string   `json:"margin_left"`   //PDF文档左边距，写数字即可，默认72pt
	MarginRight  string   `json:"margin_right"`  //PDF文档左边距，写数字即可，默认72pt
	MarginTop    string   `json:"margin_top"`    //PDF文档左边距，写数字即可，默认72pt
	MarginBottom string   `json:"margin_bottom"` //PDF文档左边距，写数字即可，默认72pt
	More         []string `json:"more"`          //更多导出选项[PDF导出选项，具体参考：https://manual.calibre-ebook.com/generated/en/ebook-convert.html#pdf-output-options]
	Toc          []Toc    `json:"toc"`           //目录
	///////////////////////////////////////////
	Order []string `json:"-"` //这个不需要赋值
}

var (
	output       = "output" //文档导出文件夹
	ebookConvert = "ebook-convert"
)

func CheckConvertCommand() error {
	args := []string{ "--version" }
	cmd := exec.Command(ebookConvert, args...)

	return cmd.Run()
}

// 接口文档 https://manual.calibre-ebook.com/generated/en/ebook-convert.html#table-of-contents
//根据json配置文件，创建文档转化对象
func NewConverter(configFile string, debug ...bool) (converter *Converter, err error) {
	var (
		cfg      Config
		basepath string
		db       bool
	)
	if len(debug) > 0 {
		db = debug[0]
	}

	if cfg, err = parseConfig(configFile); err == nil {
		if basepath, err = filepath.Abs(filepath.Dir(configFile)); err == nil {
			//设置默认值
			if len(cfg.Timestamp) == 0 {
				cfg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			}
			if len(cfg.Charset) == 0 {
				cfg.Charset = "utf-8"
			}
			converter = &Converter{
				Config:   cfg,
				BasePath: basepath,
				Debug:    db,
				ProcessNum: 1,
				process: make(chan func(),4),
				limitChan: make(chan bool,1),
			}
		}
	}
	return
}

//执行文档转换
func (convert *Converter) Convert() (err error) {
	if !convert.Debug { //调试模式下不删除生成的文件
		defer convert.converterDefer() //最后移除创建的多余而文件
	}
	if convert.process == nil{
		convert.process = make(chan func(),4)
	}
	if convert.limitChan == nil {
		if convert.ProcessNum <= 0 {
			convert.ProcessNum = 1
		}
		convert.limitChan = make(chan bool,convert.ProcessNum)
		for i := 0; i < convert.ProcessNum;i++{
			convert.limitChan <- true
		}
	}

	if err = convert.generateMimeType(); err != nil {
		return
	}
	if err = convert.generateMetaInfo(); err != nil {
		return
	}
	if err = convert.generateTocNcx(); err != nil { //生成目录
		return
	}
	if err = convert.generateSummary(); err != nil { //生成文档内目录
		return
	}
	if err = convert.generateTitlePage(); err != nil { //生成封面
		return
	}
	if err = convert.generateContentOpf(); err != nil { //这个必须是generate*系列方法的最后一个调用
		return
	}

	//将当前文件夹下的所有文件压缩成zip包，然后直接改名成content.epub
	f := filepath.Join(convert.OutputPath, "content.epub")
	os.Remove(f) //如果原文件存在了，则删除;
	if err = ziptil.Zip(convert.BasePath,f); err == nil {
		//创建导出文件夹
		os.Mkdir(convert.BasePath+"/"+output, os.ModePerm)
		if len(convert.Config.Format) > 0 {
			var errs []string

			go func(convert *Converter) {
				for _, v := range convert.Config.Format {
					fmt.Println("convert to " + v)
					switch strings.ToLower(v) {
					case "epub":
						convert.process  <- func() {
							if err = convert.convertToEpub(); err != nil {
								errs = append(errs, err.Error())
								fmt.Println("转换EPUB文档失败：" + err.Error())
							}
						}

					case "mobi":
						convert.process <- func() {
							if err = convert.convertToMobi(); err != nil {
								errs = append(errs, err.Error())
								fmt.Println("转换MOBI文档失败：" + err.Error())
							}
						}
					case "pdf":
						convert.process <- func() {
							if err = convert.convertToPdf(); err != nil {
								fmt.Println("转换PDF文档失败：" + err.Error())
								errs = append(errs, err.Error())
							}
						}
					case "docx":
						convert.process <- func() {
							if err = convert.convertToDocx(); err != nil {
								fmt.Println("转换WORD文档失败：" + err.Error())
								errs = append(errs, err.Error())
							}
						}
					}
				}
				close(convert.process)
			}(convert)

			group :=  sync.WaitGroup{}
			for {
				action, isClosed := <-convert.process
				if action == nil && !isClosed {
					break;
				}
				group.Add(1)
				<- convert.limitChan
				go func(group *sync.WaitGroup) {
					action()
					group.Done()
					convert.limitChan <- true
				}(&group)
			}

			group.Wait()

			if len(errs) > 0 {
				err = errors.New(strings.Join(errs, "\n"))
			}
		} else {
			err = convert.convertToPdf()
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("压缩目录出错" + err.Error())
	}
	return
}

//删除生成导出文档而创建的文件
func (this *Converter) converterDefer() {
	//删除不必要的文件
	os.RemoveAll(filepath.Join(this.BasePath, "META-INF"))
	os.RemoveAll(filepath.Join(this.BasePath, "content.epub"))
	os.RemoveAll(filepath.Join(this.BasePath, "mimetype"))
	os.RemoveAll(filepath.Join(this.BasePath, "toc.ncx"))
	os.RemoveAll(filepath.Join(this.BasePath, "content.opf"))
	os.RemoveAll(filepath.Join(this.BasePath, "titlepage.xhtml")) //封面图片待优化
	os.RemoveAll(filepath.Join(this.BasePath, "summary.html"))    //文档目录
}

//生成metainfo
func (this *Converter) generateMetaInfo() (err error) {
	xml := `<?xml version="1.0"?>
			<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
			   <rootfiles>
				  <rootfile full-path="content.opf" media-type="application/oebps-package+xml"/>
			   </rootfiles>
			</container>
    `
	folder := filepath.Join(this.BasePath, "META-INF")
	os.MkdirAll(folder, os.ModePerm)
	err = ioutil.WriteFile(filepath.Join(folder, "container.xml"), []byte(xml), os.ModePerm)
	return
}

//形成mimetyppe
func (this *Converter) generateMimeType() (err error) {
	return ioutil.WriteFile(filepath.Join(this.BasePath, "mimetype"), []byte("application/epub+zip"), os.ModePerm)
}

//生成封面
func (this *Converter) generateTitlePage() (err error) {
	if ext := strings.ToLower(filepath.Ext(this.Config.Cover)); !(ext == ".html" || ext == ".xhtml") {
		xml := `<?xml version='1.0' encoding='` + this.Config.Charset + `'?>
				<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="` + this.Config.Language + `">
					<head>
						<meta http-equiv="Content-Type" content="text/html; charset=` + this.Config.Charset + `"/>
						<meta name="calibre:cover" content="true"/>
						<title>Cover</title>
						<style type="text/css" title="override_css">
							@page {padding: 0pt; margin:0pt}
							body { text-align: center; padding:0pt; margin: 0pt; }
						</style>
					</head>
					<body>
						<div>
							<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="100%" height="100%" viewBox="0 0 800 1068" preserveAspectRatio="none">
								<image width="800" height="1068" xlink:href="` + strings.TrimPrefix(this.Config.Cover, "./") + `"/>
							</svg>
						</div>
					</body>
				</html>
		`
		if err = ioutil.WriteFile(filepath.Join(this.BasePath, "titlepage.xhtml"), []byte(xml), os.ModePerm); err == nil {
			this.GeneratedCover = "titlepage.xhtml"
		}
	}
	return
}

//生成文档目录
func (this *Converter) generateTocNcx() (err error) {
	ncx := `<?xml version='1.0' encoding='` + this.Config.Charset + `'?>
			<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="%v">
			  <head>
				<meta content="4" name="dtb:depth"/>
				<meta content="calibre (2.85.1)" name="dtb:generator"/>
				<meta content="0" name="dtb:totalPageCount"/>
				<meta content="0" name="dtb:maxPageNumber"/>
			  </head>
			  <docTitle>
				<text>%v</text>
			  </docTitle>
			  <navMap>%v</navMap>
			</ncx>
	`
	codes, _ := this.tocToXml(0, 1)
	ncx = fmt.Sprintf(ncx, this.Config.Language, html.EscapeString(this.Config.Title), strings.Join(codes, ""))
	return ioutil.WriteFile(filepath.Join(this.BasePath, "toc.ncx"), []byte(ncx), os.ModePerm)
}

//生成文档目录，即summary.html
func (this *Converter) generateSummary() (err error) {
	//目录
	summary := `<!DOCTYPE html>
				<html lang="` + this.Config.Language + `">
				<head>
				    <meta charset="` + this.Config.Charset + `">
				    <title>目录</title>
				    <style>
				        body{margin: 0px;padding: 0px;}h1{text-align: center;padding: 0px;margin: 0px;}ul,li{list-style: none;}
				        a{text-decoration: none;color: #4183c4;text-decoration: none;font-size: 16px;line-height: 28px;}
				    </style>
				</head>
				<body>
				    <h1>目&nbsp;&nbsp;&nbsp;&nbsp;录</h1>
				    %v
				</body>
				</html>`
	summary = fmt.Sprintf(summary, strings.Join(this.tocToSummary(0), ""))
	return ioutil.WriteFile(filepath.Join(this.BasePath, "summary.html"), []byte(summary), os.ModePerm)
}

//将toc转成toc.ncx文件
func (this *Converter) tocToXml(pid, idx int) (codes []string, next_idx int) {
	var code string
	for _, toc := range this.Config.Toc {
		if toc.Pid == pid {
			code, idx = this.getNavPoint(toc, idx)
			codes = append(codes, code)
			for _, item := range this.Config.Toc {
				if item.Pid == toc.Id {
					code, idx = this.getNavPoint(item, idx)
					codes = append(codes, code)
					var code_arr []string
					code_arr, idx = this.tocToXml(item.Id, idx)
					codes = append(codes, code_arr...)
					codes = append(codes, `</navPoint>`)
				}
			}
			codes = append(codes, `</navPoint>`)
		}
	}
	next_idx = idx
	return
}

//将toc转成toc.ncx文件
func (this *Converter) tocToSummary(pid int) (summarys []string) {
	summarys = append(summarys, "<ul>")
	for _, toc := range this.Config.Toc {
		if toc.Pid == pid {
			summarys = append(summarys, fmt.Sprintf(`<li><a href="%v">%v</a></li>`, toc.Link, html.EscapeString(toc.Title)))
			for _, item := range this.Config.Toc {

				if item.Pid == toc.Id {
					summarys = append(summarys, fmt.Sprintf(`<li><ul><li><a href="%v">%v</a></li>`, item.Link, html.EscapeString(item.Title)))
					summarys = append(summarys, "<li>")
					summarys = append(summarys, this.tocToSummary(item.Id)...)
					summarys = append(summarys, "</li></ul></li>")
				}

			}
		}
	}
	summarys = append(summarys, "</ul>")
	return
}

//生成navPoint
func (this *Converter) getNavPoint(toc Toc, idx int) (navpoint string, nextidx int) {
	navpoint = `
	<navPoint id="id%v" playOrder="%v">
		<navLabel>
			<text>%v</text>
		</navLabel>
		<content src="%v"/>`
	navpoint = fmt.Sprintf(navpoint, toc.Id, idx, html.EscapeString(toc.Title), toc.Link)
	this.Config.Order = append(this.Config.Order, toc.Link)
	nextidx = idx + 1
	return
}

//生成content.opf文件
//倒数第二步调用
func (this *Converter) generateContentOpf() (err error) {
	var (
		guide       string
		manifest    string
		manifestArr []string
		spine       string //注意：如果存在封面，则需要把封面放在第一个位置
		spineArr    []string
	)

	meta := `<dc:title>%v</dc:title>
			<dc:contributor opf:role="bkp">%v</dc:contributor>
			<dc:publisher>%v</dc:publisher>
			<dc:description>%v</dc:description>
			<dc:language>%v</dc:language>
			<dc:creator opf:file-as="Unknown" opf:role="aut">%v</dc:creator>
			<meta name="calibre:timestamp" content="%v"/>
	`
	meta = fmt.Sprintf(meta, html.EscapeString(this.Config.Title), html.EscapeString(this.Config.Contributor), html.EscapeString(this.Config.Publisher), html.EscapeString(this.Config.Description), this.Config.Language, html.EscapeString(this.Config.Creator), this.Config.Timestamp)
	if len(this.Config.Cover) > 0 {
		meta = meta + `<meta name="cover" content="cover"/>`
		guide = `<reference href="titlepage.xhtml" title="Cover" type="cover"/>`
		manifest = fmt.Sprintf(`<item href="%v" id="cover" media-type="%v"/>`, this.Config.Cover, GetMediaType(filepath.Ext(this.Config.Cover)))
		spineArr = append(spineArr, `<itemref idref="titlepage"/>`)
	}

	if _, err := os.Stat(this.BasePath + "/summary.html"); err == nil {
		spineArr = append(spineArr, `<itemref idref="summary"/>`) //目录

	}

	//扫描所有文件
	if files, err := filetil.ScanFiles(this.BasePath); err == nil {
		basePath := strings.Replace(this.BasePath, "\\", "/", -1)
		for _, file := range files {
			if !file.IsDir {
				ext := strings.ToLower(filepath.Ext(file.Path))
				sourcefile := strings.TrimPrefix(file.Path, basePath+"/")
				id := "ncx"
				if ext != ".ncx" {
					if file.Name == "titlepage.xhtml" { //封面
						id = "titlepage"
					} else if file.Name == "summary.html" { //目录
						id = "summary"
					} else {
						id = cryptil.Md5Crypt(sourcefile)
					}
				}
				if mt := GetMediaType(ext); mt != "" { //不是封面图片，且media-type不为空
					if sourcefile != strings.TrimLeft(this.Config.Cover, "./") { //不是封面图片，则追加进来。封面图片前面已经追加进来了
						manifestArr = append(manifestArr, fmt.Sprintf(`<item href="%v" id="%v" media-type="%v"/>`, sourcefile, id, mt))
					}
				}
			} else {
				fmt.Println(file.Path)
			}
		}

		items := make(map[string]string)
		for _, link := range this.Config.Order {
			id := cryptil.Md5Crypt(link)
			if _, ok := items[id]; !ok { //去重
				items[id] = id
				spineArr = append(spineArr, fmt.Sprintf(`<itemref idref="%v"/>`, id))
			}
		}
		manifest = manifest + strings.Join(manifestArr, "\n")
		spine = strings.Join(spineArr, "\n")
	} else {
		return err
	}

	pkg := `<?xml version='1.0' encoding='` + this.Config.Charset + `'?>
		<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="uuid_id" version="2.0">
		  <metadata xmlns:opf="http://www.idpf.org/2007/opf" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:calibre="http://calibre.kovidgoyal.net/2009/metadata">
			%v
		  </metadata>
		  <manifest>
			%v
		  </manifest>
		  <spine toc="ncx">
			%v
		  </spine>
			%v
		</package>
	`
	if len(guide) > 0 {
		guide = `<guide>` + guide + `</guide>`
	}
	pkg = fmt.Sprintf(pkg, meta, manifest, spine, guide)
	return ioutil.WriteFile(filepath.Join(this.BasePath, "content.opf"), []byte(pkg), os.ModePerm)
}

//转成epub
func (this *Converter) convertToEpub() (err error) {
	args := []string{
		filepath.Join(this.OutputPath, "content.epub"),
		filepath.Join(this.OutputPath, output, "book.epub"),
	}
	//cmd := exec.Command(ebookConvert, args...)
	//
	//if this.Debug {
	//	fmt.Println(cmd.Args)
	//}
	//fmt.Println("正在转换EPUB文件", args[0])
	//return cmd.Run()

	return filetil.CopyFile(args[0],args[1])
}

//转成mobi
func (this *Converter) convertToMobi() (err error) {
	args := []string{
		filepath.Join(this.OutputPath, "content.epub"),
		filepath.Join(this.OutputPath, output, "book.mobi"),
	}
	cmd := exec.Command(ebookConvert, args...)
	if this.Debug {
		fmt.Println(cmd.Args)
	}
	fmt.Println("正在转换 MOBI 文件", args[0])
	return cmd.Run()
}

//转成pdf
func (this *Converter) convertToPdf() (err error) {
	args := []string{
		filepath.Join(this.OutputPath, "content.epub"),
		filepath.Join(this.OutputPath, output, "book.pdf"),
	}
	//页面大小
	if len(this.Config.PaperSize) > 0 {
		args = append(args, "--paper-size", this.Config.PaperSize)
	}
	//文字大小
	if len(this.Config.FontSize) > 0 {
		args = append(args, "--pdf-default-font-size", this.Config.FontSize)
	}

	//header template
	if len(this.Config.Header) > 0 {
		args = append(args, "--pdf-header-template", this.Config.Header)
	}

	//footer template
	if len(this.Config.Footer) > 0 {
		args = append(args, "--pdf-footer-template",this.Config.Footer)
	}

	if strings.Count(this.Config.MarginLeft,"") > 0 {
		args = append(args, "--pdf-page-margin-left", this.Config.MarginLeft)
	}
	if strings.Count(this.Config.MarginTop,"") > 0 {
		args = append(args, "--pdf-page-margin-top", this.Config.MarginTop)
	}
	if strings.Count(this.Config.MarginRight,"") > 0 {
		args = append(args, "--pdf-page-margin-right", this.Config.MarginRight)
	}
	if strings.Count(this.Config.MarginBottom,"") > 0 {
		args = append(args, "--pdf-page-margin-bottom", this.Config.MarginBottom)
	}

	//更多选项
	if len(this.Config.More) > 0 {
		args = append(args, this.Config.More...)
	}

	cmd := exec.Command(ebookConvert, args...)
	if this.Debug {
		fmt.Println(cmd.Args)
	}
	fmt.Println("正在转换 PDF 文件", args[0])
	return cmd.Run()
}

// 转成word
func (this *Converter) convertToDocx() (err error) {
	args := []string{
		filepath.Join(this.OutputPath , "content.epub"),
		filepath.Join(this.OutputPath , output , "book.docx"),
	}
	args = append(args, "--docx-no-toc")

	//页面大小
	if len(this.Config.PaperSize) > 0 {
		args = append(args, "--docx-page-size", this.Config.PaperSize)
	}

	if len(this.Config.MarginLeft) > 0 {
		args = append(args, "--docx-page-margin-left", this.Config.MarginLeft)
	}
	if len(this.Config.MarginTop) > 0 {
		args = append(args, "--docx-page-margin-top", this.Config.MarginTop)
	}
	if len(this.Config.MarginRight) > 0 {
		args = append(args, "--docx-page-margin-right", this.Config.MarginRight)
	}
	if len(this.Config.MarginBottom) > 0 {
		args = append(args, "--docx-page-margin-bottom", this.Config.MarginBottom)
	}
	cmd := exec.Command(ebookConvert, args...)

	if this.Debug {
		fmt.Println(cmd.Args)
	}
	fmt.Println("正在转换 DOCX 文件", args[0])
	return cmd.Run()
}






















