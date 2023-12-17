package dot

import (
	"strconv"

	"golang.org/x/tools/go/callgraph"
	"gonum.org/v1/gonum/graph/formats/dot/ast"
)

type (
	edge struct {
		from int
		to   int
	}
)

func DrawCallGraphAsDotDigraph(callg *callgraph.Graph, id string) (*ast.Graph, error) {

	stmts := []ast.Stmt(nil)
	nodes := make(map[int]*ast.Node, len(callg.Nodes))
	drawedEdges := make(map[edge]bool)

	for _, callNode := range callg.Nodes {
		if callNode.Func == nil {
			continue
		}
		node := &ast.Node{
			ID: strconv.FormatInt(int64(callNode.ID), 10),
		}
		nodes[callNode.ID] = node
		stmts = append(stmts, &ast.NodeStmt{
			Node: node,
			Attrs: []*ast.Attr{
				{
					Key: "label",
					Val: `"` + callNode.Func.String() + `"`,
				},
			},
		})
	}

	err := callgraph.GraphVisitEdges(callg, func(e *callgraph.Edge) error {
		if drawedEdges[edge{e.Caller.ID, e.Callee.ID}] {
			return nil
		}
		drawedEdges[edge{e.Caller.ID, e.Callee.ID}] = true
		stmts = append(stmts, &ast.EdgeStmt{
			From: nodes[e.Caller.ID],
			To: &ast.Edge{
				Directed: true,
				Vertex:   nodes[e.Callee.ID],
			},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ast.Graph{
		Strict:   true,
		Directed: true,
		ID:       id,
		Stmts:    stmts,
	}, nil
}
