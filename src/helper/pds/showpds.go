// Show the dependency structure of specified package
package main

import (
	"basic/prof"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"pkgtool"
	"runtime/debug"
	"strings"
)

const (
	ARROWS = "->"
)

var (
	pkgImportPathFlag string
)

func init() {
	flag.StringVar(&pkgImportPathFlag, "p", "", "The path of target package.")
}

func main() {
	prof.Start()
	defer func() {
		prof.Stop()
		if err := recover(); err != nil {
			fmt.Errorf("FATAL ERROR: %s", err)
			debug.PrintStack()
		}
	}()
	flag.Parse()
	pkgImportPath := getPkgImportPath()
	pn := pkgtool.NewPkgNode(pkgImportPath)
	fmt.Printf("The package node of '%s': %v\n", pkgImportPath, *pn)
	err := pn.Grow()
	if err != nil {
		fmt.Printf("GROW ERROR: %s\n", err)
	}
	fmt.Printf("The dependency structure of package '%s':\n", pkgImportPath)
	ShowDepStruct(pn, "")
}

func ShowDepStruct(pnode *pkgtool.PkgNode, prefix string) {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	importPath := pnode.ImportPath()
	buf.WriteString(importPath)
	deps := pnode.Deps()
	//fmt.Printf("P_NODE: '%s', DEP_LEN: %d\n", importPath, len(deps))
	if len(deps) == 0 {
		fmt.Printf("%s\n", buf.String())
		return
	}
	buf.WriteString(ARROWS)
	for _, v := range deps {
		ShowDepStruct(v, buf.String())
	}
}

func getPkgImportPath() string {
	if len(pkgImportPathFlag) > 0 {
		return pkgImportPathFlag
	}
	fmt.Printf("The flag p is invalid, use current dir as package import path.")
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	srcDirs := pkgtool.GetSrcDirs(false)
	var importPath string
	for _, v := range srcDirs {
		if strings.HasPrefix(currentDir, v) {
			importPath = currentDir[len(v):]
			break
		}
	}
	if strings.TrimSpace(importPath) == "" {
		panic(errors.New("Can not parse the import path!"))
	}
	return importPath
}
