package utils

import (
	"os/exec"
	"bufio"
	"io"
	"io/ioutil"
	"errors"

	"github.com/astaxie/beego"
)


// 使用 wkhtmltopdf 是实现 html 转 pdf.
// 中文说明：http://www.jianshu.com/p/4d65857ffe5e#
func ConverterHtmlToPdf(uri []string,path string) (error) {
	exe := beego.AppConfig.String("wkhtmltopdf")

	if exe == "" {
		return errors.New("wkhtmltopdf not exist.")
	}
	params := []string{"-c",exe,"--margin-bottom","25"}

	params = append(params,uri...)
	params = append(params,path)

	beego.Info(params)

	cmd := exec.Command("/bin/bash",params...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return errors.New("StdoutPipe: " + err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {

		return errors.New("StderrPipe: " + err.Error())
	}

	if err := cmd.Start(); err != nil {

		return errors.New("Start: "+ err.Error())
	}

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line ,err2 := reader.ReadString('\n')

		if err2 != nil || io.EOF == err2 {
			break
		}

		beego.Info(line)
	}

	bytesErr, err := ioutil.ReadAll(stderr)

	if err == nil {
		beego.Info(string(bytesErr))
	}else{
		beego.Error("Error: Stderr => " + err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {

		beego.Error("Error: ", err.Error())

		return err
	}

	return nil
}
