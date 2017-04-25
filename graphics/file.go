package graphics

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"image/jpeg"
	"image/png"
	"image/gif"
)
// 将图片保存到指定的路径
func SaveImage(p string, src image.Image) error {

	f, err := os.OpenFile( p, os.O_SYNC | os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(p)

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {

		err = jpeg.Encode(f, src, &jpeg.Options{Quality : 100 })
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors : 256})
	}
	return err
}
