package parse_test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/beauknowssoftware/graphqlgen/internal/parse"
	"github.com/google/go-cmp/cmp"
)

var ignoreTokenPosition = cmpopts.IgnoreFields(parse.Token{}, "Loc")

func TestLex(t *testing.T) {
	tests := map[string]struct {
		expectedTokens []parse.Token
	}{
		"requiredParams.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.LeftParenToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "a",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Int",
				},
				{
					TokenType: parse.BangToken,
				},
				{
					TokenType: parse.CommaToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "b",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.BangToken,
				},
				{
					TokenType: parse.RightParenToken,
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"requiredFieldType.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.BangToken,
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"params.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.LeftParenToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "a",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Int",
				},
				{
					TokenType: parse.CommaToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "b",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightParenToken,
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"requiredArrayElement.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.LeftBracketToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.BangToken,
				},
				{
					TokenType: parse.RightBracketToken,
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"requiredArray.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.LeftBracketToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightBracketToken,
				},
				{
					TokenType: parse.BangToken,
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"array.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.LeftBracketToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightBracketToken,
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"directives.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.AtToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "my_directive",
				},
				{
					TokenType: parse.AtToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "another_directive",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
		"ping.graphqls": {
			expectedTokens: []parse.Token{
				{
					TokenType: parse.TextToken,
					Value:     "type",
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "ping",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "String",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "schema",
				},
				{
					TokenType: parse.LeftCurlyToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "query",
				},
				{
					TokenType: parse.ColonToken,
				},
				{
					TokenType: parse.TextToken,
					Value:     "Query",
				},
				{
					TokenType: parse.RightCurlyToken,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			schema := parse.TestGetDoc(t, name)
			l := parse.NewLexer(schema)

			c := make(chan parse.Token)
			go l.Lex(c)

			tokens := make([]parse.Token, 0)
			for t := range c {
				if t.TokenType != parse.WhitespaceToken {
					tokens = append(tokens, t)
				}
			}

			if diff := cmp.Diff(test.expectedTokens, tokens, ignoreTokenPosition); diff != "" {
				t.Fatalf("mismatch (expected, got) %v", diff)
			}
		})
	}
}

func TestComplex(t *testing.T) {
	schema := parse.TestGetDoc(t, "complex.graphqls")
	l := parse.NewLexer(schema)

	c := make(chan parse.Token)
	go l.Lex(c)

	for token := range c {
		if token.TokenType == parse.ErrorToken {
			t.Fatalf("got error token %v", token)
		}
	}
}
