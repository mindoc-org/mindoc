package utils

import (
	"os/exec"
	"fmt"
	"bufio"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"github.com/astaxie/beego"
)
// 中文说明：http://www.jianshu.com/p/4d65857ffe5e#
func ConverterHtmlToPdf(uri []string,path string)  {

	exe := beego.AppConfig.String("wkhtmltopdf")

	if exe == "" {
		panic("wkhtmltopdf not exist.")
	}

	params := []string{"/C",exe,"--margin-bottom","25"}

	params = append(params,uri...)
	params = append(params,path)

	cmd := exec.Command("cmd",params...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("StdoutPipe: " + err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Start: ", err.Error())
		return
	}

	reader := bufio.NewReader(stdout)
	enc := mahonia.NewDecoder("gbk")

	//实时循环读取输出流中的一行内容
	for {
		line ,err2 := reader.ReadString('\n')

		if err2 != nil || io.EOF == err2 {
			break
		}

		beego.Error(enc.ConvertString(line))
	}

	bytesErr, err := ioutil.ReadAll(stderr)

	if err == nil {
		beego.Error(enc.ConvertString(string(bytesErr)))
	}else{
		beego.Error("Error: Stderr => " + err.Error())
	}

	if err := cmd.Wait(); err != nil {

		beego.Error("Error: ", err.Error())

		return
	}


	return

}
