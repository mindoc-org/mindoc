package graphics

import (
	"errors"
	"image"
	"os"

	"github.com/nfnt/resize"
)

func ImageCopy(src image.Image, x, y, w, h int) (image.Image, error) {

	var subImg image.Image

	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.Paletted); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.Paletted) //图片裁剪x0 y0 x1 y1
	} else {

		return subImg, errors.New("图片解码失败")
	}

	return subImg, nil
}

func ImageCopyFromFile(p string, x, y, w, h int) (image.Image, error) {
	var src image.Image

	file, err := os.Open(p)
	if err != nil {
		return src, err
	}
	defer file.Close()
	src, _, err = image.Decode(file)

	return ImageCopy(src, x, y, w, h)
}

func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}
func ImageResizeSaveFile(src image.Image, width, height int, p string) error {
	dst := resize.Resize(uint(width), uint(height), src, resize.Lanczos3)

	return SaveImage(p, dst)
}
