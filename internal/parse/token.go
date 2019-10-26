package parse

import (
	"fmt"
)

type TokenType int

type Loc struct {
	Line   int
	Column int
}

func (l Loc) String() string {
	return fmt.Sprintf("%v,%v", l.Line+1, l.Column+1)
}

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
	AtToken
	LeftBracketToken
	RightBracketToken
	EOFToken
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
	case AtToken:
		return "at"
	case LeftBracketToken:
		return "left bracket"
	case RightBracketToken:
		return "right bracket"
	case EOFToken:
		return "end of file"
	default:
		return "unknown"
	}
}

type Token struct {
	TokenType TokenType
	Loc       Loc
	Value     string
}

func (t Token) String() string {
	if t.Value == "" {
		return fmt.Sprintf("%v keyword @(%v)", t.TokenType, t.Loc)
	}
	return fmt.Sprintf("%v keyword = %v @(%v)", t.TokenType, t.Value, t.Loc)
}
