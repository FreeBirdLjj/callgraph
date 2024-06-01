package gencallgraph

import (
	"context"
	"errors"

	"github.com/freebirdljj/immutable/slice"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const (
	LoadAll packages.LoadMode = -1
)

func GenCallGraphForPackages(ctx context.Context, patterns []string) (*callgraph.Graph, error) {
	return genCallGraphForPackages(
		&packages.Config{
			Context: ctx,
			Mode:    LoadAll,
		},
		patterns,
	)
}

func genCallGraphForPackages(conf *packages.Config, patterns []string) (*callgraph.Graph, error) {

	pkgs, err := packages.Load(conf, patterns...)
	if err != nil {
		return nil, err
	}

	pkgErrs := slice.Concat(
		slice.Map(
			slice.FromGoSlice(pkgs),
			func(pkg *packages.Package) slice.Slice[error] {
				return slice.Map(
					slice.FromGoSlice(pkg.Errors),
					func(err packages.Error) error {
						return err
					},
				)
			},
		),
	)
	if !pkgErrs.Empty() {
		return nil, errors.Join(pkgErrs...)
	}

	ssaProg, _ := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	ssaProg.Build()

	callg := static.CallGraph(ssaProg)
	return callg, nil
}
