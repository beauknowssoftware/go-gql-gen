package parse

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	lines  []string
	dSlice []rune
	loc    Loc
}

func NewLexer(document string) Lexer {
	lines := strings.Split(document, "\n")
	dSlice := []rune(lines[0])
	return Lexer{lines: lines, dSlice: dSlice}
}

func (l Lexer) newToken(t TokenType, v string) Token {
	return Token{t, l.loc, v}
}

func (l Lexer) isLastLine() bool {
	return l.loc.Line >= len(l.lines)-1
}

func (l Lexer) isEndOfLine() bool {
	return l.loc.Column >= len(l.dSlice)
}

func (l *Lexer) checkNextLine() {
	for l.isEndOfLine() && !l.isLastLine() {
		l.loc.Column = 0
		l.loc.Line++
		l.dSlice = []rune(l.lines[l.loc.Line])
	}
}

func (l *Lexer) increment() {
	l.loc.Column++
	l.checkNextLine()
}

func (l Lexer) currentRune() rune {
	return l.dSlice[l.loc.Column]
}

func (l *Lexer) while(cond func(rune) bool) string {
	start := l.loc.Column
	for !l.isEndOfLine() && cond(l.dSlice[l.loc.Column]) {
		l.loc.Column++
	}

	v := string(l.dSlice[start:l.loc.Column])

	l.checkNextLine()

	return v
}

func (l Lexer) isDone() bool {
	return l.isLastLine() && l.isEndOfLine()
}

func (l *Lexer) Lex(c chan Token) {
	for !l.isDone() {
		r := l.currentRune()
		switch {
		case r == '{':
			c <- l.newToken(LeftCurlyToken, "")
			l.increment()
		case r == '}':
			c <- l.newToken(RightCurlyToken, "")
			l.increment()
		case r == ':':
			c <- l.newToken(ColonToken, "")
			l.increment()
		case r == '(':
			c <- l.newToken(LeftParenToken, "")
			l.increment()
		case r == ')':
			c <- l.newToken(RightParenToken, "")
			l.increment()
		case r == ',':
			c <- l.newToken(CommaToken, "")
			l.increment()
		case r == '!':
			c <- l.newToken(BangToken, "")
			l.increment()
		case unicode.IsSpace(r):
			w := l.while(unicode.IsSpace)
			c <- l.newToken(WhitespaceToken, string(len(w)))
		case unicode.IsLetter(r):
			value := l.while(unicode.IsLetter)
			c <- l.newToken(TextToken, value)
		default:
			c <- l.newToken(ErrorToken, fmt.Sprintf("unknown rune %v", r))
			l.increment()
		}
	}
	close(c)
}
