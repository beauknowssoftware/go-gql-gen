package parse

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	lines  []string
	dSlice []rune
	cl     int
	cc     int
}

func NewLexer(document string) Lexer {
	lines := strings.Split(document, "\n")
	dSlice := []rune(lines[0])
	return Lexer{lines, dSlice, 0, 0}
}

func (l Lexer) newToken(t TokenType, v string) Token {
	return Token{t, l.cl, l.cc, v}
}

func (l Lexer) isLastLine() bool {
	return l.cl >= len(l.lines)-1
}

func (l Lexer) isEndOfLine() bool {
	return l.cc >= len(l.dSlice)
}

func (l *Lexer) checkNextLine() {
	for l.isEndOfLine() && !l.isLastLine() {
		l.cc = 0
		l.cl++
		l.dSlice = []rune(l.lines[l.cl])
	}
}

func (l *Lexer) increment() {
	l.cc++
	l.checkNextLine()
}

func (l Lexer) currentRune() rune {
	return l.dSlice[l.cc]
}

func (l *Lexer) while(cond func(rune) bool) string {
	start := l.cc
	for !l.isEndOfLine() && cond(l.dSlice[l.cc]) {
		l.cc++
	}

	v := string(l.dSlice[start:l.cc])

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
		case unicode.IsSpace(r):
			w := l.while(unicode.IsSpace)
			c <- l.newToken(WhitespaceToken, string(len(w)))
		case unicode.IsLetter(r):
			value := l.while(unicode.IsLetter)
			c <- l.newToken(TextToken, value)
		default:
			c <- l.newToken(ErrorToken, fmt.Sprintf("unknown rune %v", r))
		}
	}
	close(c)
}
