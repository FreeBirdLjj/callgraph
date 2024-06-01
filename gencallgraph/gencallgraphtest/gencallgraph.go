package gencallgraphtest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// `srcs`: a map of paths to sources.
// NOTE: There should be one and only 1 `main` package in `srcs`.
func GenCallGraphForFiles(srcs map[string]string) (*callgraph.Graph, error) {

	fset := token.NewFileSet()
	mainPkgPath := ""
	filesInPkgs := make(map[string][]*ast.File)

	for srcPath, src := range srcs {

		f, err := parser.ParseFile(fset, srcPath, src, parser.AllErrors)
		if err != nil {
			return nil, err
		}

		importPath := path.Dir(srcPath)
		filesInPkgs[importPath] = append(filesInPkgs[importPath], f)

		if f.Name.Name == "main" {
			mainPkgPath = path.Dir(srcPath)
		}
	}

	conf := types.Config{}
	importFunc := func(path string) (*types.Package, error) {
		return conf.Check(path, fset, filesInPkgs[path], nil)
	}
	conf.Importer = &importer{
		ImportFunc: importFunc,
	}

	mainPkg, _, err := ssautil.BuildPackage(
		&conf,
		fset,
		types.NewPackage(mainPkgPath, ""),
		filesInPkgs[mainPkgPath],
		ssa.SanityCheckFunctions,
	)
	if err != nil {
		return nil, err
	}

	callg := static.CallGraph(mainPkg.Prog)
	return callg, nil
}
