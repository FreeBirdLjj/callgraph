package gencallgraphtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/callgraph"

	"github.com/freebirdljj/callgraph/internal/testcase"
)

func TestGenCallGraphForFiles(t *testing.T) {

	type testcaseT struct {
		goSrcs        map[string]string
		expectedEdges map[string]string
	}

	testcase.RunTestCases(
		t,
		map[string]testcaseT{
			"shoul dsucceed": {
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

func main() { lib.F() }
`,
				},
				expectedEdges: map[string]string{
					"main.init": "lib.init",
					"main.main": "lib.F",
				},
			},
		},
		func(t *testing.T, testcase *testcaseT) {

			callg, err := GenCallGraphForFiles(testcase.goSrcs)
			require.NoError(t, err)

			edges := make(map[string]string)
			callgraph.GraphVisitEdges(callg, func(e *callgraph.Edge) error {
				edges[e.Caller.Func.String()] = e.Callee.Func.String()
				return nil
			})

			assert.Equal(t, testcase.expectedEdges, edges)
		},
	)
}
