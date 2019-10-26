package parse_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
)

func TestTraverse(t *testing.T) {
	tests := map[string]struct {
		expectedNodes []parse.Node
	}{
		"params.graphqls": {
			expectedNodes: []parse.Node{
				parse.DocumentNode{
					Definitions: []parse.Node{
						parse.TypeDefNode{
							Name: "Query",
							Fields: []parse.Node{
								parse.FieldNode{
									Name: "ping",
									Type: parse.TypeNode{
										Name: "String",
									},
									Params: []parse.Node{
										parse.ParamNode{
											Name: "a",
											Type: parse.TypeNode{
												Name: "Int",
											},
										},
										parse.ParamNode{
											Name: "b",
											Type: parse.TypeNode{
												Name: "String",
											},
										},
									},
								},
							},
						},
						parse.SchemaNode{
							Fields: []parse.Node{
								parse.FieldNode{
									Name: "query",
									Type: parse.TypeNode{
										Name: "Query",
									},
								},
							},
						},
					},
				},
				parse.TypeDefNode{
					Name: "Query",
					Fields: []parse.Node{
						parse.FieldNode{
							Name: "ping",
							Type: parse.TypeNode{
								Name: "String",
							},
							Params: []parse.Node{
								parse.ParamNode{
									Name: "a",
									Type: parse.TypeNode{
										Name: "Int",
									},
								},
								parse.ParamNode{
									Name: "b",
									Type: parse.TypeNode{
										Name: "String",
									},
								},
							},
						},
					},
				},
				parse.FieldNode{
					Name: "ping",
					Type: parse.TypeNode{
						Name: "String",
					},
					Params: []parse.Node{
						parse.ParamNode{
							Name: "a",
							Type: parse.TypeNode{
								Name: "Int",
							},
						},
						parse.ParamNode{
							Name: "b",
							Type: parse.TypeNode{
								Name: "String",
							},
						},
					},
				},
				parse.TypeNode{
					Name: "String",
				},
				parse.ParamNode{
					Name: "a",
					Type: parse.TypeNode{
						Name: "Int",
					},
				},
				parse.TypeNode{
					Name: "Int",
				},
				parse.ParamNode{
					Name: "b",
					Type: parse.TypeNode{
						Name: "String",
					},
				},
				parse.TypeNode{
					Name: "String",
				},
				parse.SchemaNode{
					Fields: []parse.Node{
						parse.FieldNode{
							Name: "query",
							Type: parse.TypeNode{
								Name: "Query",
							},
						},
					},
				},
				parse.FieldNode{
					Name: "query",
					Type: parse.TypeNode{
						Name: "Query",
					},
				},
				parse.TypeNode{
					Name: "Query",
				},
			},
		},
		"ping.graphqls": {
			expectedNodes: []parse.Node{
				parse.DocumentNode{
					Definitions: []parse.Node{
						parse.TypeDefNode{
							Name: "Query",
							Fields: []parse.Node{
								parse.FieldNode{
									Name: "ping",
									Type: parse.TypeNode{
										Name: "String",
									},
								},
							},
						},
						parse.SchemaNode{
							Fields: []parse.Node{
								parse.FieldNode{
									Name: "query",
									Type: parse.TypeNode{
										Name: "Query",
									},
								},
							},
						},
					},
				},
				parse.TypeDefNode{
					Name: "Query",
					Fields: []parse.Node{
						parse.FieldNode{
							Name: "ping",
							Type: parse.TypeNode{
								Name: "String",
							},
						},
					},
				},
				parse.FieldNode{
					Name: "ping",
					Type: parse.TypeNode{
						Name: "String",
					},
				},
				parse.TypeNode{
					Name: "String",
				},
				parse.SchemaNode{
					Fields: []parse.Node{
						parse.FieldNode{
							Name: "query",
							Type: parse.TypeNode{
								Name: "Query",
							},
						},
					},
				},
				parse.FieldNode{
					Name: "query",
					Type: parse.TypeNode{
						Name: "Query",
					},
				},
				parse.TypeNode{
					Name: "Query",
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

			if diff := cmp.Diff(test.expectedNodes, nodes, ignoreNodePosition); diff != "" {
				t.Fatalf("mismatch (expected, got) %v", diff)
			}
		})
	}
}
