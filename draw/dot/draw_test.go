package dot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph/formats/dot"
	dot_ast "gonum.org/v1/gonum/graph/formats/dot/ast"

	"github.com/freebirdljj/callgraph/gencallgraph/gencallgraphtest"
	"github.com/freebirdljj/callgraph/internal/testcase"
)

func TestDrawCallGraphAsDotDigraph(t *testing.T) {

	type testcaseT struct {
		goSrc       string
		expectedDot string
	}

	testcase.RunTestCases(
		t,
		map[string]testcaseT{
			"should succeed": {
				goSrc: `
package main

func f()    {}

func main() { f() }
`,
				expectedDot: `
strict digraph G {
	1 [label="main.init"]
	2 [label="main.f"]
	3 [label="main.main"]
	3 -> 2
}
`,
			},
			"should remove duplicates": {
				goSrc: `
package main

func f() {}

func main() {
	f()
	f()
}
`,
				expectedDot: `
strict digraph G {
	1 [label="main.init"]
	2 [label="main.f"]
	3 [label="main.main"]
	3 -> 2
}`,
			},
		},
		func(t *testing.T, testcase *testcaseT) {

			callg, err := gencallgraphtest.GenCallGraphForFiles(map[string]string{
				"main/main.go": testcase.goSrc,
			})
			require.NoError(t, err)

			digraph, err := DrawCallGraphAsDotDigraph(callg, "G")
			require.NoError(t, err)

			expectedDot, err := dot.ParseString(testcase.expectedDot)
			require.NoError(t, err)

			expectedCallGraph := extractCallGraph(expectedDot.Graphs[0])
			gotCallGraph := extractCallGraph(digraph)

			assert.Equal(t, expectedCallGraph, gotCallGraph)
		},
	)
}

func extractCallGraph(graph *dot_ast.Graph) map[string]string {

	funcs := make(map[string]string)
	calls := make(map[string]string)

	for _, stmt := range graph.Stmts {
		switch stmt := stmt.(type) {
		case *dot_ast.NodeStmt:
			label := ""
			for _, attr := range stmt.Attrs {
				if attr.Key == "label" {
					label = attr.Val
				}
			}
			funcs[stmt.Node.ID] = label
		case *dot_ast.EdgeStmt:
			caller := stmt.From.(*dot_ast.Node)
			callee := stmt.To.Vertex.(*dot_ast.Node)
			calls[caller.ID] = calls[callee.ID]
		}
	}

	callgraph := make(map[string]string, len(calls))
	for callerID, calleeID := range calls {
		callgraph[funcs[callerID]] = callgraph[funcs[calleeID]]
	}
	return callgraph
}
