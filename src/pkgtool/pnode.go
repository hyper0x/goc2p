package pkgtool

import (
	"path/filepath"
	"strings"
)

var pkgNodesCache map[string]*PkgNode

func init() {
	pkgNodesCache = make(map[string]*PkgNode)
}

type PkgNode struct {
	srcDir     string
	importPath string
	triggers   []*PkgNode
	deps       []*PkgNode
	growed     bool
}

func (self *PkgNode) SrcDir() string {
	return self.srcDir
}

func (self *PkgNode) ImportPath() string {
	return self.importPath
}

func (self *PkgNode) AddTrigger(pn *PkgNode) {
	self.triggers = append(self.triggers, pn)
}

func (self *PkgNode) AddDep(pn *PkgNode) {
	self.deps = append(self.deps, pn)
}

func (self *PkgNode) Triggers() []*PkgNode {
	triggers := make([]*PkgNode, len(self.triggers))
	copy(triggers, self.triggers)
	return triggers
}

func (self *PkgNode) Deps() []*PkgNode {
	deps := make([]*PkgNode, len(self.deps))
	copy(deps, self.deps)
	return deps
}

func (self *PkgNode) IsLeaf() bool {
	if len(self.deps) == 0 {
		return true
	} else {
		return false
	}
}

func (self *PkgNode) Grow() error {
	if self.growed {
		return nil
	} else {
		self.growed = true
	}
	importPaths, err := getImportsFromPackage(self.importPath, false)
	if err != nil {
		return err
	}
	l := len(importPaths)
	if l == 0 {
		return nil
	}
	//fmt.Printf("PN: %v, IPs: %v\n", self, importPaths)
	subPns := make([]*PkgNode, l)
	for i, importPath := range importPaths {
		pn, ok := pkgNodesCache[importPath]
		if !ok {
			pn = NewPkgNode(importPath)
			pkgNodesCache[importPath] = pn
		}
		subPns[i] = pn
	}
	for _, subPn := range subPns {
		subPn.AddTrigger(self)
		self.AddDep(subPn)
		err = subPn.Grow()
		if err != nil {
			return err
		}
		subPn.growed = true
	}
	return nil
}

func NewPkgNode(importPath string) *PkgNode {
	packageAbsPath := getAbsPathOfPackage(importPath)
	var srcDir string
	importDir := filepath.FromSlash(importPath)
	if strings.HasSuffix(packageAbsPath, importDir) {
		srcDir = packageAbsPath[:strings.LastIndex(packageAbsPath, importDir)]
	}
	return &PkgNode{
		srcDir:     srcDir,
		importPath: importPath,
		triggers:   make([]*PkgNode, 0),
		deps:       make([]*PkgNode, 0)}
}
