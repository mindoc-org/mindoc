package utils

import (
	"strings"
	"os"
	"fmt"
	"path/filepath"
	"io"
	"math"
)

func AbsolutePath(p string) (string,error) {

	if strings.HasPrefix(p, "~") {
		home := os.Getenv("HOME")
		if home == "" {
			panic(fmt.Sprintf("can not found HOME in envs, '%s' AbsPh Failed!", p))
		}
		p = fmt.Sprint(home, string(p[1:]))
	}
	s, err := filepath.Abs(p)

	if nil != err {
		return  "",err
	}
	return s,nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer dst.Close()
	return io.Copy(dst, src)
}

func FormatBytes(size int64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}

	s := float64(size)

	i := 0

	for ; s >= 1024 && i < 4 ; i ++ {
		s /= 1024
	}


	return fmt.Sprintf("%.2f%s",s,units[i])
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









