package filetil

import (
	"os"
	"path/filepath"
	"strings"
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
