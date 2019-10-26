package parse_test

import (
	"testing"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		expectedAST parse.DocumentNode
	}{
		"requiredFieldType.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.DefinitionNode{
					&parse.TypeNode{
						Name: "Query",
						Fields: []parse.FieldNode{
							{
								Name:     "ping",
								Type:     "String",
								Required: true,
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
		},
		"params.graphqls": {
			expectedAST: parse.DocumentNode{
				Definitions: []parse.DefinitionNode{
					&parse.TypeNode{
						Name: "Query",
						Fields: []parse.FieldNode{
							{
								Name: "ping",
								Type: "String",
								Params: []parse.ParamNode{
									{
										Name: "a",
										Type: "Int",
									},
									{
										Name: "b",
										Type: "String",
									},
								},
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
		},
		"ping.graphqls": {
			expectedAST: parse.DocumentNode{
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

			if diff := cmp.Diff(&test.expectedAST, ast); diff != "" {
				t.Fatalf("mismatch (expected, got) %v", diff)
			}
		})
	}
}
