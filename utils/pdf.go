package utils

import (
	"github.com/signintech/gopdf"
	"github.com/astaxie/beego"
)
func ConverterPdf(output string,htmlList map[string]string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}})

	pdf.AddPage()

	err := pdf.AddTTFFont("HDZB_5", "./static/pdf-fonts/msyh.ttf")
	if err != nil {
		beego.Error("ConverterPdf => ",err)
		return err
	}
	err = pdf.SetFont("HDZB_5", "", 14)
	if err != nil {
		beego.Error("ConverterPdf => " , err)
		return err
	}
	
	pdf.Cell(nil, "您好")
	pdf.WritePdf(output)

	return nil
}
