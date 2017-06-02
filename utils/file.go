package utils

import (
	"strings"
	"os"
	"fmt"
	"path/filepath"
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