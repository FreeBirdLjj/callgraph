package dot

import (
	go_ast "go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"gonum.org/v1/gonum/graph/formats/dot"
	dot_ast "gonum.org/v1/gonum/graph/formats/dot/ast"
)

func TestDrawCallGraphAsDotDigraph(t *testing.T) {

	t.Parallel()

	t.Run("single file cases", func(t *testing.T) {

		t.Parallel()

		testcases := []struct {
			name        string
			goSrc       string
			expectedDot string
		}{
			{
				name: "should succeed",
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
			{
				name: "should remove duplicates",
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
		}

		for _, testcase := range testcases {
			testcase := testcase
			t.Run(testcase.name, func(t *testing.T) {

				t.Parallel()

				fset := token.NewFileSet()

				f, err := parser.ParseFile(fset, "main.go", testcase.goSrc, parser.AllErrors)
				require.NoError(t, err)

				mainPkg, _, err := ssautil.BuildPackage(
					&types.Config{
						Importer: importer.Default(),
					},
					fset,
					types.NewPackage("main", ""),
					[]*go_ast.File{f},
					ssa.SanityCheckFunctions,
				)
				require.NoError(t, err)

				callg := static.CallGraph(mainPkg.Prog)

				digraph, err := DrawCallGraphAsDotDigraph(callg, "G")
				require.NoError(t, err)

				expectedDot, err := dot.ParseString(testcase.expectedDot)
				require.NoError(t, err)

				expectedCallGraph := extractCallGraph(expectedDot.Graphs[0])
				gotCallGraph := extractCallGraph(digraph)

				assert.Equal(t, gotCallGraph, expectedCallGraph)
			})
		}
	})
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
