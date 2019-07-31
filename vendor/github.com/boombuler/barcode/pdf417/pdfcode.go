package pdf417

import (
	"image"
	"image/color"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type pdfBarcode struct {
	data  string
	width int
	code  *utils.BitList
}

func (c *pdfBarcode) Metadata() barcode.Metadata {
	return barcode.Metadata{barcode.TypePDF, 2}
}

func (c *pdfBarcode) Content() string {
	return c.data
}

func (c *pdfBarcode) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *pdfBarcode) Bounds() image.Rectangle {
	height := c.code.Len() / c.width

	return image.Rect(0, 0, c.width, height*moduleHeight)
}

func (c *pdfBarcode) At(x, y int) color.Color {
	if c.code.GetBit((y/moduleHeight)*c.width + x) {
		return color.Black
	}
	return color.White
}
