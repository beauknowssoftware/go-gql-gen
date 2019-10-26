package parse

import (
	"fmt"
)

type TokenType int

const (
	ErrorToken TokenType = iota
	TextToken
	LeftCurlyToken
	ColonToken
	RightCurlyToken
	WhitespaceToken
	LeftParenToken
	RightParenToken
	CommaToken
	BangToken
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
	case LeftParenToken:
		return "left paren"
	case RightParenToken:
		return "right paren"
	case CommaToken:
		return "comma"
	case BangToken:
		return "bang"
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
