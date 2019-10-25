package parse

import (
	"fmt"
)

type Node interface {
	Kind() string
}

type DocumentNode struct {
	Definitions []DefinitionNode
}

func (n DocumentNode) Kind() string { return "document" }

type DefinitionNode interface {
	Node
}

type TypeNode struct {
	Name   string
	Fields []FieldNode
}

func (n TypeNode) Kind() string { return "type" }

type SchemaNode struct {
	Fields []FieldNode
}

func (n SchemaNode) Kind() string { return "schema" }

type FieldNode struct {
	Name string
	Type string
}

func (n FieldNode) Kind() string { return "field" }

type TokenType int

const (
	ErrorToken TokenType = iota
	TextToken
	LeftCurlyToken
	ColonToken
	RightCurlyToken
	WhitespaceToken
)

func (tt TokenType) String() string {
	switch tt {
	case ErrorToken:
		return "error"
	case TextToken:
		return "text"
	case LeftCurlyToken:
		return "left curly"
	case ColonToken:
		return "colon"
	case RightCurlyToken:
		return "right curly"
	case WhitespaceToken:
		return "whitespace"
	default:
		return "unknown"
	}
}

type Token struct {
	TokenType TokenType
	Line      int
	Column    int
	Value     string
}

func (t Token) String() string {
	if t.Value == "" {
		return fmt.Sprintf("%v token @(%v,%v)", t.TokenType, t.Line, t.Column)
	}
	return fmt.Sprintf("%v token = %v @(%v,%v)", t.TokenType, t.Value, t.Line, t.Column)
}
