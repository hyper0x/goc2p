package pkgtool

import (
	"os"
	"path/filepath"
	"strings"
)

func getAbsPathOfPackage(importPath string) string {
	for _, srcDir := range GetSrcDirs(false) {
		absPath := filepath.Join(srcDir, filepath.FromSlash(importPath))
		//fmt.Printf("IP: %s, IP_DIR: %s, AP: %s\n", importPath, filepath.FromSlash(importPath), absPath)
		_, err := os.Open(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func getGoSourceFileAbsPaths(packageAbspath string, containsTestFile bool) ([]string, error) {
	f, err := os.Open(packageAbspath)
	if err != nil {
		return nil, err
	}
	subs, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	absPaths := make([]string, 0)
	for _, v := range subs {
		fi := v.(os.FileInfo)
		name := fi.Name()
		if !fi.IsDir() && strings.HasSuffix(name, ".go") {
			if strings.HasSuffix(name, "_test.go") && !containsTestFile {
				continue
			}
			absPath := filepath.Join(packageAbspath, name)
			absPaths = append(absPaths, absPath)
		}
	}
	return absPaths, nil
}
