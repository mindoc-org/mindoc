package aztec

import (
	"bytes"
	"image"
	"image/color"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type aztecCode struct {
	*utils.BitList
	size    int
	content []byte
}

func newAztecCode(size int) *aztecCode {
	return &aztecCode{utils.NewBitList(size * size), size, nil}
}

func (c *aztecCode) Content() string {
	return string(c.content)
}

func (c *aztecCode) Metadata() barcode.Metadata {
	return barcode.Metadata{barcode.TypeAztec, 2}
}

func (c *aztecCode) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *aztecCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.size, c.size)
}

func (c *aztecCode) At(x, y int) color.Color {
	if c.GetBit(x*c.size + y) {
		return color.Black
	}
	return color.White
}

func (c *aztecCode) set(x, y int) {
	c.SetBit(x*c.size+y, true)
}

func (c *aztecCode) string() string {
	buf := new(bytes.Buffer)
	for y := 0; y < c.size; y++ {
		for x := 0; x < c.size; x++ {
			if c.GetBit(x*c.size + y) {
				buf.WriteString("X ")
			} else {
				buf.WriteString("  ")
			}
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}
