package main

import (
	"fmt"
	"github.com/lifei6671/gocaptcha"
	"html/template"
	"log"
	"net/http"
)

const (
	dx = 150
	dy = 50
)

func main() {

	err := gocaptcha.ReadFonts("fonts", ".ttf")
	if err != nil {
		fmt.Println(err)
		return
	}

	//	gocaptcha.SetFontFamily(fontFils...)

	//gocaptcha.SetFontFamily(
	//	"fonts/3Dumb.ttf",
	//	"fonts/DeborahFancyDress.ttf",
	//	"fonts/actionj.ttf",
	//	"fonts/chromohv.ttf",
	//	"fonts/D3Parallelism.ttf",
	//	"fonts/Flim-Flam.ttf",
	//	"fonts/KREMLINGEORGIANI3D.ttf",
	//	)

	http.HandleFunc("/", Index)
	http.HandleFunc("/get/", Get)
	fmt.Println("服务已启动...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}
func Get(w http.ResponseWriter, r *http.Request) {

	captchaImage, err := gocaptcha.NewCaptchaImage(dx, dy, gocaptcha.RandLightColor())

	captchaImage.DrawNoise(gocaptcha.CaptchaComplexHigh)

	captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexHigh)

	captchaImage.DrawText(gocaptcha.RandText(4))
	//captchaImage.Drawline(3);
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))
	captchaImage.DrawHollowLine()
	if err != nil {
		fmt.Println(err)
	}

	captchaImage.SaveImage(w, gocaptcha.ImageFormatJpeg)
}
