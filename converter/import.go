package converter

import (
	"errors"
	"github.com/lifei6671/mindoc/utils"
	"os"
	"path/filepath"
	"crypto/md5"
	"io"
	"fmt"
	"os/exec"
	"github.com/lifei6671/mindoc/utils/ziptil"
	"io/ioutil"
	"encoding/xml"
)

type ResolveResult struct {

}

type XmlResult struct {
	XMLName     xml.Name `xml:"ncx"`
	Head XmlHead `xml:"head"`
	NavMap XmlTocNavMap `xml:"navMap"`
	Title string `xml:"docTitle>text"`
}

type XmlHead struct {
	XMLName xml.Name `xml:"head"`
	Meta []XmlMeta `xml:"meta"`
}

type XmlMeta struct {
	XMLName xml.Name `xml:"meta"`
	Content string `xml:"content,attr"`
	Name string `xml:"name,attr"`

}
type XmlDocTitle struct {
	Text string `xml:"text"`
}
type XmlTocNavMap struct {
	XMLName xml.Name `xml:"navMap"`
	NavPoint []XmlNavPoint `xml:"navPoint"`
}

type XmlNavPoint struct {
	XMLName xml.Name `xml:"navPoint"`
	Content XmlContent `xml:"content"`
	NavLabel string `xml:"navLabel>text"`
}

type XmlContent struct {
	XMLName xml.Name `xml:"content"`
	Src string `xml:"src,attr"`
}

type XmlNavLabel struct {

}

func Resolve(p string) (ResolveResult,error) {
	result := ResolveResult{

	}

	if !utils.FileExists(p) {
		return  result,errors.New("文件不存在 " + p)
	}

	w := md5.New()
	io.WriteString(w, p)   //将str写入到w中
	md5str := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式

	tempPath := filepath.Join(os.TempDir(),md5str)

	os.MkdirAll(tempPath,0766)

	epub := filepath.Join(tempPath , "book.epub")

	args := []string{p,epub}

	cmd := exec.Command(ebookConvert, args...)


	if err := cmd.Run(); err != nil {
		fmt.Println("执行转换命令失败：" + err.Error())
		return result,err
	}
	fmt.Println(epub)

	unzipPath := filepath.Join(tempPath,"output")

	if err := ziptil.Unzip(epub, unzipPath); err != nil {
		fmt.Println("解压缩失败：" + err.Error())
		return result,err
	}
	xmlPath := filepath.Join(unzipPath,"toc.ncx")

	data,err := ioutil.ReadFile(xmlPath);

	if err != nil {
		fmt.Println("toc.ncx 文件不存在：" + err.Error())
		return  result,err
	}
	v := XmlResult{}

	err = xml.Unmarshal([]byte(data), &v)

	if err != nil {
		fmt.Println("解析XML失败：" + err.Error())
		return  result,err
	}

	fmt.Println(v)

	return  result,nil
}