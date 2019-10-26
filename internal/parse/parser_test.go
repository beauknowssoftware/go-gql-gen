package parse_test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
	"github.com/google/go-cmp/cmp"
)

var ignoreNodePosition = cmpopts.IgnoreFields(parse.NodeLoc{}, "NodeLoc")

func TestParse(t *testing.T) {
	tests := map[string]struct {
		expectedAST parse.DocumentNode
	}{
		"requiredParams.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: "String",
								Params: []parse.Node{
									parse.ParamNode{
										Name:     "a",
										Type:     "Int",
										Required: true,
									},
									parse.ParamNode{
										Name:     "b",
										Type:     "String",
										Required: true,
									},
								},
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"requiredFieldType.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name:     "ping",
								Type:     "String",
								Required: true,
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"params.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: "String",
								Params: []parse.Node{
									parse.ParamNode{
										Name: "a",
										Type: "Int",
									},
									parse.ParamNode{
										Name: "b",
										Type: "String",
									},
								},
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"directives.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: "String",
								Directives: []parse.Node{
									parse.DirectiveNode{
										Name: "my_directive",
									},
									parse.DirectiveNode{
										Name: "another_directive",
									},
								},
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"requiredArray.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name:     "ping",
								Type:     "String",
								Multiple: true,
								Required: true,
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"array.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name:     "ping",
								Type:     "String",
								Multiple: true,
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
		"ping.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: "String",
							},
						},
					},
					parse.SchemaNode{
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "query",
								Type: "Query",
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			schema := parse.TestGetDoc(t, name)

			l := parse.NewLexer(schema)
			p := parse.New(l)

			ast, err := p.Parse()
			if err != nil {
				t.Fatalf("failed to parse: %v (%v)", err.Error, err.Token)
			}

			if diff := cmp.Diff(test.expectedAST, ast, ignoreNodePosition); diff != "" {
				t.Fatalf("mismatch (expected, got) %v", diff)
			}
		})
	}
}
