package filetil

import (
	"os"
	"path/filepath"
	"strings"
	"io"
	"fmt"
	"math"
	"io/ioutil"
	"bytes"
)

//==================================
//更多文件和目录的操作，使用filepath包和os包
//==================================

//返回的目录扫描结果
type FileList struct {
	IsDir   bool   //是否是目录
	Path    string //文件路径
	Ext     string //文件扩展名
	Name    string //文件名
	Size    int64  //文件大小
	ModTime int64  //文件修改时间戳
}

//目录扫描
//@param			dir			需要扫描的目录
//@return			fl			文件列表
//@return			err			错误
func ScanFiles(dir string) (fl []FileList, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			path = strings.Replace(path, "\\", "/", -1) //文件路径处理
			fl = append(fl, FileList{
				IsDir:   info.IsDir(),
				Path:    path,
				Ext:     strings.ToLower(filepath.Ext(path)),
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Unix(),
			})
		}
		return err
	})
	return
}

//拷贝文件
func CopyFile(source string, dst string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	_,err = os.Stat(filepath.Dir(dst))

	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(dst),0766)
		}else{
			return err
		}
	}


	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dst, sourceInfo.Mode())
		}

	}

	return
}

//拷贝目录
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourceFilePointer := filepath.Join(source , obj.Name())

		destinationFilePointer := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourceFilePointer, destinationFilePointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourceFilePointer, destinationFilePointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}

func AbsolutePath(p string) (string, error) {

	if strings.HasPrefix(p, "~") {
		home := os.Getenv("HOME")
		if home == "" {
			panic(fmt.Sprintf("can not found HOME in envs, '%s' AbsPh Failed!", p))
		}
		p = fmt.Sprint(home, string(p[1:]))
	}
	s, err := filepath.Abs(p)

	if nil != err {
		return "", err
	}
	return s, nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

func FormatBytes(size int64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}

	s := float64(size)

	i := 0

	for ; s >= 1024 && i < 4; i++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, units[i])
}

func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

//判断指定目录下是否存在指定后缀的文件
func HasFileOfExt(path string,exts []string) bool {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			ext := filepath.Ext(info.Name())

			for _,item := range exts {
				if strings.EqualFold(ext,item) {
					return os.ErrExist
				}
			}

		}
		return nil
	})

	return err == os.ErrExist
}
// IsImageExt 判断是否是图片后缀
func IsImageExt(filename string) bool {
	ext := filepath.Ext(filename)

	return strings.EqualFold(ext, ".jpg") ||
		strings.EqualFold(ext, ".jpeg") ||
		strings.EqualFold(ext, ".png") ||
		strings.EqualFold(ext, ".gif") ||
		strings.EqualFold(ext,".svg") ||
		strings.EqualFold(ext,".bmp") ||
		strings.EqualFold(ext,".webp")
}
//忽略字符串中的BOM头
func ReadFileAndIgnoreUTF8BOM(filename string) ([]byte,error) {

	data,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,err
	}
	if data == nil {
		return nil,nil
	}
	data = bytes.Replace(data,[]byte("\r"),[]byte(""),-1)
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
		return data[3:],err
	}


	return data,nil
}
