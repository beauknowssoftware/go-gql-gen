package parse_test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/beauknowssoftware/go-gql-gen/internal/parse"
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
											Name:     "Int",
											Required: true,
										},
									},
									parse.ParamNode{
										Name: "b",
										Type: parse.TypeNode{
											Name:     "String",
											Required: true,
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
		},
		"requiredFieldType.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name:     "String",
									Required: true,
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
		},
		"params.graphqls": {
			expectedAST: parse.DocumentNode{
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
		},
		"directives.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name: "String",
								},
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
								Type: parse.TypeNode{
									Name: "Query",
								},
							},
						},
					},
				},
			},
		},
		"fullyRequiredArray.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name:            "String",
									Multiple:        true,
									Required:        true,
									NonNullElements: true,
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
		},
		"requiredArrayElement.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name:            "String",
									Multiple:        true,
									NonNullElements: true,
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
		},
		"requiredArray.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name:     "String",
									Multiple: true,
									Required: true,
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
		},
		"array.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.Node{
					parse.TypeDefNode{
						Name: "Query",
						Fields: []parse.Node{
							parse.FieldNode{
								Name: "ping",
								Type: parse.TypeNode{
									Name:     "String",
									Multiple: true,
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
		},
		"input.graphqls": {
			expectedAST: parse.DocumentNode{
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
					parse.TypeDefNode{
						Name:  "PingInput",
						Input: true,
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
		},
		"ping.graphqls": {
			expectedAST: parse.DocumentNode{
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

func TestParseComplex(t *testing.T) {
	schema := parse.TestGetDoc(t, "complex.graphqls")
	l := parse.NewLexer(schema)
	p := parse.New(l)

	if _, err := p.Parse(); err != nil {
		t.Fatalf("failed to parse %v", err)
	}
}
