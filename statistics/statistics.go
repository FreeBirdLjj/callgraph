package statistics

import (
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

type (
	FuncStatistics struct {
		References      int
		ReferencedFuncs int
		ReferencedPkgs  int
	}
)

func StatisticsCallGraphAtFuncLevel(callg *callgraph.Graph) (map[*ssa.Function]FuncStatistics, error) {
	res := make(map[*ssa.Function]FuncStatistics, len(callg.Nodes))
	for _, node := range callg.Nodes {
		if node.Func == nil {
			continue
		}
		res[node.Func] = analyzeFunc(node)
	}
	return res, nil
}

func analyzeFunc(node *callgraph.Node) FuncStatistics {

	funcs := make(map[*callgraph.Node]struct{})
	pkgs := make(map[*ssa.Package]struct{})

	for _, edge := range node.In {

		caller := edge.Caller

		funcs[caller] = struct{}{}
		pkgs[caller.Func.Package()] = struct{}{}
	}

	references := len(node.In)
	referencedFuncs := len(funcs)
	referencedPkgs := len(pkgs)

	return FuncStatistics{
		References:      references,
		ReferencedFuncs: referencedFuncs,
		ReferencedPkgs:  referencedPkgs,
	}
}
