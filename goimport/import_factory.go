package goimport

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

func ParseRelation(
	rootPath string) *ImportPathFactory {

	factory := NewImportPathFactory(
		rootPath,
	)
	factory.Root = factory.Get(rootPath)
	if factory.Root == nil {
		return nil
	}
	return factory
}

type ImportPathFactory struct {
	Root   *ImportPath
	Filter *ImportFilter
	Pool   map[string]*ImportPath
}

func NewImportPathFactory(
	rootPath string) *ImportPathFactory {

	self := &ImportPathFactory{Pool: make(map[string]*ImportPath)}
	filter := NewImportFilter(
		rootPath,
	)
	self.Filter = filter
	return self
}
func (self *ImportPathFactory) GetRoot() *ImportPath {
	return self.Root
}

func (self *ImportPathFactory) GetAll() []*ImportPath {
	ret := make([]*ImportPath, 0)
	for _, value := range self.Pool {
		ret = append(ret, value)
	}
	return ret
}

func (self *ImportPathFactory) Get(importPath string) *ImportPath {
	// aquire from pool
	pool := self.Pool
	if _, ok := pool[importPath]; ok {
		return pool[importPath]
	}

	dirPath := goSrc(importPath)
	if !fileExists(dirPath) {
		dirPath = searchFilePathFromLocal(importPath)
		if dirPath == "" {
			return nil
		}
	}

	ret := &ImportPath{
		ImportPath: importPath,
		// Files      []*Source
		// children   []dotwriter.IDotNode
		// parents    []dotwriter.IDotNode
	}

	pool[importPath] = ret

	// パス上のすべてのGoファイルを取得
	fileNames := glob(dirPath)

	// Goファイルから静的解析用にオブジェクト化 (ret.Files)
	// オブジェクト化したソースからimportを取得してノード化する (ret.children, ret.parents)
	ret.Init(self, fileNames)

	return ret
}

//ImportFilter
type ImportFilter struct {
	root string
}

func NewImportFilter(root string) *ImportFilter {
	impf := &ImportFilter{
		root: root,
	}
	return impf

}

func isMatched(pattern string, target string) bool {
	r, _ := regexp.Compile(pattern)
	return r.MatchString(target)
}

func glob(dirPath string) []string {
	fileNames, err := filepath.Glob(filepath.Join(dirPath, "/*.go"))
	if err != nil {
		panic("no gofiles")
	}

	files := make([]string, 0, len(fileNames))

	for _, v := range fileNames {
		if isMatched("test", v) {
			continue
		}
		if isMatched("example", v) {
			continue
		}
		files = append(files, v)
	}
	return files
}

func goSrc(inputPath string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", inputPath)
}

// searchFilePathFromLocal ファイルが見つからなかった場合は空文字を返す
func searchFilePathFromLocal(importPath string) string {
	// workspaceとimportPathから対象パッケージのファイルパスを探す
	dirPath := filepath.Join(workSpacePath, importPath)

	if !fileExists(dirPath) {
		return ""
	}

	return dirPath
}

var workSpacePath string

func init() {
	// go.mod.module & os.Getwd から ワークスペース のPathを取得しておく
	mydir, err := os.Getwd()
	if err != nil {
		return
	}

	b, err := ioutil.ReadFile(filepath.Join(mydir, "go.mod"))
	if err != nil {
		return
	}
	p := modfile.ModulePath(b)

	if !strings.HasSuffix(mydir, p) {
		return
	}

	workSpacePath = mydir[:len(mydir)-len(p)]
}
