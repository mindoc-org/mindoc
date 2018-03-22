package filetil

import (
	"os"
	"path/filepath"
	"strings"
	"io"
	"fmt"
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

