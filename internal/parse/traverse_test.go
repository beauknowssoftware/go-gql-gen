package parse_test

import (
	"testing"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
	"github.com/google/go-cmp/cmp"
)

func TestTraverse(t *testing.T) {
	tests := map[string]struct {
		expectedNodes []parse.Node
	}{
		"ping.graphqls": {
			expectedNodes: []parse.Node{
				parse.DocumentNode{
					Definitions: []parse.DefinitionNode{
						&parse.TypeNode{
							Name: "Query",
							Fields: []parse.FieldNode{
								{
									Name: "ping",
									Type: "String",
								},
							},
						},
						&parse.SchemaNode{
							Fields: []parse.FieldNode{
								{
									Name: "query",
									Type: "Query",
								},
							},
						},
					},
				},
				&parse.TypeNode{
					Name: "Query",
					Fields: []parse.FieldNode{
						{
							Name: "ping",
							Type: "String",
						},
					},
				},
				parse.FieldNode{
					Name: "ping",
					Type: "String",
				},
				&parse.SchemaNode{
					Fields: []parse.FieldNode{
						{
							Name: "query",
							Type: "Query",
						},
					},
				},
				parse.FieldNode{
					Name: "query",
					Type: "Query",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			schema := parse.TestGetDoc(t, name)

			ast := parse.TestParse(t, schema)

			nodes := make([]parse.Node, 0)
			parse.Traverse(ast, func(node parse.Node) bool {
				nodes = append(nodes, node)
				return true
			})

			if diff := cmp.Diff(test.expectedNodes, nodes); diff != "" {
				t.Fatalf("mismatch (expected, got) %v", diff)
			}
		})
	}
}
