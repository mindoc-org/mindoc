// Package preinit performs os.Chdir to the executable's directory before any
// other package (including beego) initialises. This prevents beego's config
// init() from printing a spurious "open conf/app.conf: no such file" debug
// message when the binary is launched from a directory other than its own.
//
// Import this package as the very first blank import in main.go:
//
//	_ "github.com/mindoc-org/mindoc/internal/preinit"
package preinit

import (
	"os"
	"path/filepath"
	"strings"
)

func init() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	exeDir := filepath.Dir(exe)
	// Skip go-run temporary build directories.
	if strings.Contains(exeDir, "go-build") {
		return
	}
	// Only chdir when conf/app.conf actually exists next to the binary.
	if _, err := os.Stat(filepath.Join(exeDir, "conf", "app.conf")); err != nil {
		return
	}
	_ = os.Chdir(exeDir)
}
