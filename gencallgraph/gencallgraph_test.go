package gencallgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages/packagestest"

	"github.com/freebirdljj/callgraph/internal/testcase"
)

func Test_genCallgraphForPackages(t *testing.T) {

	type testcaseT struct {
		goSrcs        map[string]any
		expectedEdges map[string]string
	}

	rootModulePath := "github.com/freebirdljj/callgraph/gencallgraph"

	testcase.RunTestCases(
		t,
		map[string]testcaseT{
			"should succeed": {
				goSrcs: map[string]any{
					"a/a.go": `
package a

func F() {}
`,
					"b/b.go": `
package b

import (
	"` + rootModulePath + `/a"
)

func G() { a.F() }
`,
				},
				expectedEdges: map[string]string{
					rootModulePath + "/b.init": rootModulePath + "/a.init",
					rootModulePath + "/b.G":    rootModulePath + "/a.F",
				},
			},
		},
		func(t *testing.T, testcase *testcaseT) {

			exported := packagestest.Export(t, packagestest.Modules, []packagestest.Module{
				{
					Name:  rootModulePath,
					Files: testcase.goSrcs,
				},
			})
			defer exported.Cleanup()

			conf := *exported.Config
			conf.Mode = LoadAll

			callg, err := genCallGraphForPackages(&conf, []string{rootModulePath + "..."})
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
