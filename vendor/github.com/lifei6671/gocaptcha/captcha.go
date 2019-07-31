package gocaptcha

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"io"
	"image/png"
	"image/jpeg"
	"errors"
	"time"
	"math/rand"
	"github.com/golang/freetype"
	"flag"
	"golang.org/x/image/font"
	"log"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"image/gif"
	"fmt"
	"os"
	"strings"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	r = rand.New(rand.NewSource(time.Now().UnixNano()));
	FontFamily []string = make([]string,0);
)

const txtChars = "AaCcDdEeFfGgHhJjKkLMmNnPpQqRrSsTtUuVvWwXxYtZ2346789";


const(
	//图片格式
	ImageFormatPng = iota
	ImageFormatJpeg
	ImageFormatGif
	//验证码噪点强度
	CaptchaComplexLower = iota
	CaptchaComplexMedium
	CaptchaComplexHigh
)

type CaptchaImage struct {
	nrgba *image.NRGBA
	width int
	height int
	Complex int
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ReadFonts(dirPth string, suffix string) (err error) {
	files := make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	SetFontFamily(files...)
	return nil
}

//新建一个图片对象
func NewCaptchaImage(width int,height int,bgColor color.RGBA) (*CaptchaImage ,error){

	m := image.NewNRGBA(image.Rect(0, 0, width, height));

	draw.Draw(m, m.Bounds(), &image.Uniform{ bgColor }, image.ZP, draw.Src);

	return &CaptchaImage{
		nrgba: m,
		height : height,
		width : width,
	},nil;
}

//保存图片对象
func (this *CaptchaImage)SaveImage(w io.Writer ,imageFormat int) error{

	if(imageFormat == ImageFormatPng){
		return png.Encode(w, this.nrgba);
	};
	if(imageFormat == ImageFormatJpeg){
		return jpeg.Encode(w,this.nrgba, &jpeg.Options{ 100 });
	}
	if(imageFormat == ImageFormatGif){
		return  gif.Encode(w,this.nrgba, &gif.Options{ NumColors : 256});
	}

	return errors.New("Not supported image format");
}
//添加一个较粗的空白直线
func (captcha *CaptchaImage) DrawHollowLine()(*CaptchaImage){

	first := (captcha.width / 20);
	end := first * 19;

	lineColor := color.RGBA{ R : 245,G:250,B:251,A:255};

	x1 := float64(r.Intn(first));
	//y1 := float64(r.Intn(y)+y);

	x2 := float64( r.Intn(first)+end);

	multiple := float64(r.Intn(5)+3)/float64(5);
	if(int(multiple*10) % 3 == 0){
		multiple = multiple * -1.0;
	}

	w := captcha.height / 20;

	for ;x1 < x2; x1 ++{

		y := math.Sin(x1*math.Pi*multiple/float64(captcha.width)) * float64(captcha.height/3);

		if(multiple < 0){
			y = y + float64(captcha.height/2);
		}
		captcha.nrgba.Set(int(x1),int(y),lineColor);

		for i:=0;i<=w;i++{
			captcha.nrgba.Set(int(x1),int(y)+i,lineColor);
		}
	}

	return captcha;
}
//画一条直线
func (captcha *CaptchaImage)Drawline(num int) (*CaptchaImage) {

	first := (captcha.width / 10);
	end := first * 9;

	y := captcha.height / 3;

	for i:=0;i<num;i++{

		point1 := Point{ X : r.Intn(first),Y:r.Intn(y)};
		point2 := Point{X : r.Intn(first)+end,Y:r.Intn(y)};

		if(i % 2 == 0){
			point1.Y = r.Intn(y)+y*2;
			point2.Y = r.Intn(y);
		}else{
			point1.Y = r.Intn(y)+y*(i%2);
			point2.Y = r.Intn(y) + y*2;
		}

		captcha.drawBeeline(point1,point2,RandDeepColor());

	}
	return captcha;
}

func (captcha *CaptchaImage)drawBeeline(point1 Point,point2 Point,lineColor color.RGBA){
	dx := math.Abs(float64(point1.X - point2.X));

	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y  >= point2.Y  {
		sy = -1
	}
	err := dx - dy
	for {
		captcha.nrgba.Set(point1.X,point1.Y,lineColor);
		captcha.nrgba.Set(point1.X+1,point1.Y,lineColor);
		captcha.nrgba.Set(point1.X-1,point1.Y,lineColor);
		captcha.nrgba.Set(point1.X+2,point1.Y,lineColor);
		captcha.nrgba.Set(point1.X-2,point1.Y,lineColor);
		if point1.X == point2.X && point1.Y == point2.Y {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

//画边框
func (captcha *CaptchaImage) DrawBorder(borderColor color.RGBA) (*CaptchaImage){
	for x :=0;x<captcha.width;x++{
		captcha.nrgba.Set(x,0,borderColor);
		captcha.nrgba.Set(x,captcha.height-1,borderColor);
	}
	for y:=0;y<captcha.height;y++{
		captcha.nrgba.Set(0,y,borderColor);
		captcha.nrgba.Set(captcha.width-1,y,borderColor);
	}
	return captcha;
}

//画噪点
func (captcha *CaptchaImage) DrawNoise(complex int) (*CaptchaImage){
	density := 18;
	if(complex == CaptchaComplexLower){
		density = 28;
	}else if(complex == CaptchaComplexMedium){
		density = 18;
	}else if(complex == CaptchaComplexHigh){
		density = 8;
	}
	maxSize := (captcha.height * captcha.width)/density;

	for i:=0;i<maxSize;i++{

		rw := r.Intn(captcha.width);
		rh := r.Intn(captcha.height);

		captcha.nrgba.Set(rw,rh,RandColor());
		size := r.Intn(maxSize);
		if size%3 == 0{
			captcha.nrgba.Set(rw+1,rh+1,RandColor());
		}
	}
	return captcha;
}
//画文字噪点
func (captcha *CaptchaImage) DrawTextNoise(complex int)(error){
	density := 1500;
	if(complex == CaptchaComplexLower){
		density = 2000;
	}else if(complex == CaptchaComplexMedium){
		density = 1500;
	}else if(complex == CaptchaComplexHigh){
		density = 1000;
	}

	maxSize := (captcha.height * captcha.width)/density;

	fmt.Println(maxSize);

	r := rand.New(rand.NewSource(time.Now().UnixNano()));

	c := freetype.NewContext();
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds());
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)
	rawFontSize := float64(captcha.height) / (1+ float64(r.Intn(7))/float64(10))

	for i:=0;i<maxSize;i++{

		rw := r.Intn(captcha.width);
		rh := r.Intn(captcha.height);

		text := RandText(1);
		fontSize := rawFontSize/2 + float64(r.Intn(5));

		c.SetSrc(image.NewUniform(RandLightColor()))
		c.SetFontSize(fontSize);
		f,err := RandFontFamily();

		if(err != nil){
			log.Println(err);
			return err;
		}
		c.SetFont(f)
		pt := freetype.Pt(rw, rh);

		_, err = c.DrawString(text, pt)
		if err != nil {
			log.Println(err)
			return err;
		}
	}
	return nil;
}
//写字
func(captcha *CaptchaImage) DrawText(text string)(error){
	c := freetype.NewContext();
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds());
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)

	fontWidth := captcha.width/ len(text);


	for i,s := range text{

		fontSize := float64(captcha.height) / (1+ float64(r.Intn(7))/float64(9))

		c.SetSrc(image.NewUniform(RandDeepColor()))
		c.SetFontSize(fontSize);
		f,err := RandFontFamily();

		if(err != nil){
			log.Println(err);
			return err;
		}
		c.SetFont(f)

		x := int(fontWidth)*i + int(fontWidth)/int(fontSize);

		y := 5 + r.Intn(captcha.height/2) + int(fontSize/2);

		pt := freetype.Pt(x, y);

		_, err = c.DrawString(string(s), pt)
		if err != nil {
			log.Println(err)
			return err;
		}
		//pt.Y += c.PointToFixed(*size * *spacing)
		//pt.X += c.PointToFixed(*size);
	}
	return nil;

}

//获取所及字体
func RandFontFamily() (*truetype.Font ,error ){
	fontfile := FontFamily[r.Intn(len(FontFamily))];

	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return &truetype.Font{},err;
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return &truetype.Font{},err;
	}
	return f,nil;
}

//随机生成深色系
func RandDeepColor() color.RGBA{

	randColor := RandColor();

	increase := float64(30 + r.Intn(255));

	red := math.Abs(math.Min(float64(randColor.R) - increase,255));

	green := math.Abs(math.Min(float64(randColor.G) - increase,255));
	blue := math.Abs(math.Min(float64(randColor.B) - increase,255));


	return color.RGBA{R:uint8(red), G : uint8(green),B:uint8(blue),A:uint8(255)};
}
//随机生成浅色
func RandLightColor() color.RGBA{

	red := r.Intn(55) + 200;
	green := r.Intn(55)+200;
	blue := r.Intn(55) + 200;


	return color.RGBA{R:uint8(red), G : uint8(green),B:uint8(blue),A:uint8(255)};
}

//生成随机颜色
func RandColor() color.RGBA{

	red := r.Intn(255);
	green := r.Intn(255);
	blue := r.Intn(255);
	if((red + green) > 400){
		blue = 0;
	}else{
		blue = 400 - green -red;
	}
	if(blue > 255){
		blue = 255;
	}
	return color.RGBA{R:uint8(red), G : uint8(green),B:uint8(blue),A:uint8(255)};
}

//生成随机字体
func RandText(num int) string {
	textNum := len(txtChars);
	text := "";
	r := rand.New(rand.NewSource(time.Now().UnixNano()));

	for i:=0;i<num ;i++  {
		text = text + string(txtChars[r.Intn(textNum)]);
	}
	return text;
}

//添加一个字体路径到字体库
func SetFontFamily(fontPath ... string){

	FontFamily= append(FontFamily,fontPath...);
}
/**
 * 颜色代码转换为RGB
 * input int
 * output int red, green, blue
 **/
func ColorToRGB(colorVal int) (color.RGBA) {

	red := colorVal >> 16
	green := (colorVal & 0x00FF00) >> 8
	blue := colorVal & 0x0000FF


	return color.RGBA{
		R:uint8(red),
		G:uint8(green),
		B:uint8(blue),
		A:uint8(255),
	}
}

