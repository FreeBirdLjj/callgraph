package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"

	"github.com/FreeBirdLjj/callgraph/gencallgraph/gencallgraphtest"
	"github.com/FreeBirdLjj/callgraph/internal/testcase"
)

func TestStatisticsCallGraphAtFuncLevel(t *testing.T) {

	type testcaseT struct {
		goSrcs             map[string]string // there should be one and only 1 `main` package
		expectedStatistics map[string]FuncStatistics
	}

	testcase.RunTestCases(
		t,
		map[string]testcaseT{
			"should succeed": {
				goSrcs: map[string]string{
					"main/main.go": `
package main

func f() {}

func main() { f() }
`},
				expectedStatistics: map[string]FuncStatistics{
					"main.init": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
					"main.f": {
						References:      1,
						ReferencedFuncs: 1,
						ReferencedPkgs:  1,
					},
					"main.main": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
				},
			},
			"duplicate calls in 1 func": {
				goSrcs: map[string]string{
					"main/main.go": `
package main

func f() {}

func main() {
	f()
	f()
}
`},
				expectedStatistics: map[string]FuncStatistics{
					"main.init": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
					"main.f": {
						References:      2,
						ReferencedFuncs: 1,
						ReferencedPkgs:  1,
					},
					"main.main": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
				},
			},
			"duplicate calls in 1 pkg": {
				goSrcs: map[string]string{
					"lib/lib.go": `
package lib

func F() {}
`,
					"main/main.go": `
package main

import (
	"lib"
)

func f() { lib.F() }

func main() { lib.F() }
`},
				expectedStatistics: map[string]FuncStatistics{
					"lib.init": {
						References:      1,
						ReferencedFuncs: 1,
						ReferencedPkgs:  1,
					},
					"lib.F": {
						References:      2,
						ReferencedFuncs: 2,
						ReferencedPkgs:  1,
					},
					"main.init": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
					"main.f": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
					"main.main": {
						References:      0,
						ReferencedFuncs: 0,
						ReferencedPkgs:  0,
					},
				},
			},
		},
		func(t *testing.T, testcase *testcaseT) {

			callg, err := gencallgraphtest.GenCallGraphForFiles(testcase.goSrcs)
			require.NoError(t, err)

			res, err := StatisticsCallGraphAtFuncLevel(callg)
			require.NoError(t, err)

			assert.Equal(t, testcase.expectedStatistics, mapFuncLevelStatistics(res))
		},
	)
}

func mapFuncLevelStatistics(statistics map[*ssa.Function]FuncStatistics) map[string]FuncStatistics {
	res := make(map[string]FuncStatistics, len(statistics))
	for f, funcStatistics := range statistics {
		res[f.String()] = funcStatistics
	}
	return res
}
