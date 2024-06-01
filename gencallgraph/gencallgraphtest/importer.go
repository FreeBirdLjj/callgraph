package gencallgraphtest

import (
	"go/types"
)

type importer struct {
	ImportFunc func(path string) (*types.Package, error)
}

// type checker
var (
	_ types.Importer = new(importer)
)

func (i *importer) Import(path string) (*types.Package, error) {
	return i.ImportFunc(path)
}
