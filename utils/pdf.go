package utils

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/astaxie/beego"
)
func ConverterPdf(output string,htmlList map[string]string) error {

	pdf := gofpdf.New("P", "mm", "A4", "./static/pdf-fonts/msyh.ttf")

	pdf.AddPage()

	pdf.SetFont("微软雅黑","B",14)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")

	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}
