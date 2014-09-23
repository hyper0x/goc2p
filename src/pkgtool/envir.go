package pkgtool

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var srcDirsCache []string

func init() {
	srcDirsCache = GetSrcDirs(true)
}

func GetGoroot() string {
	return runtime.GOROOT()
}

func GetAllGopath() []string {
	gopath := os.Getenv("GOPATH")
	var sep string
	if runtime.GOOS == "windows" {
		sep = ";"
	} else {
		sep = ":"
	}
	gopaths := strings.Split(gopath, sep)
	result := make([]string, 0)
	for _, v := range gopaths {
		if strings.TrimSpace(v) != "" {
			result = append(result, v)
		}
	}
	return result
}

func GetSrcDirs(fresh bool) []string {
	if len(srcDirsCache) > 0 && !fresh {
		return srcDirsCache
	}
	gorootSrcDir := filepath.Join(GetGoroot(), "src", "pkg")
	gopaths := GetAllGopath()
	gopathSrcDirs := make([]string, len(gopaths))
	for i, v := range gopaths {
		gopathSrcDirs[i] = filepath.Join(v, "src")
	}
	srcDirs := make([]string, 0)
	srcDirs = append(srcDirs, gorootSrcDir)
	srcDirs = append(srcDirs, gopathSrcDirs...)
	srcDirsCache = make([]string, len(srcDirs))
	copy(srcDirsCache, srcDirs)
	return srcDirs
}
