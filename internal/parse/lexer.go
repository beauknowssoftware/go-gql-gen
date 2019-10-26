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

func (l Lexer) newToken(t TokenType, v string, loc Loc) Token {
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

func isText(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (l *Lexer) Lex(c chan Token) {
	for !l.isDone() {
		r := l.currentRune()
		switch {
		case r == '{':
			c <- l.newToken(LeftCurlyToken, "", l.loc)
			l.increment()
		case r == '}':
			c <- l.newToken(RightCurlyToken, "", l.loc)
			l.increment()
		case r == ':':
			c <- l.newToken(ColonToken, "", l.loc)
			l.increment()
		case r == '(':
			c <- l.newToken(LeftParenToken, "", l.loc)
			l.increment()
		case r == ')':
			c <- l.newToken(RightParenToken, "", l.loc)
			l.increment()
		case r == ',':
			c <- l.newToken(CommaToken, "", l.loc)
			l.increment()
		case r == '!':
			c <- l.newToken(BangToken, "", l.loc)
			l.increment()
		case r == '@':
			c <- l.newToken(AtToken, "", l.loc)
			l.increment()
		case r == '[':
			c <- l.newToken(LeftBracketToken, "", l.loc)
			l.increment()
		case r == ']':
			c <- l.newToken(RightBracketToken, "", l.loc)
			l.increment()
		case unicode.IsSpace(r):
			s := l.loc
			w := l.while(unicode.IsSpace)
			c <- l.newToken(WhitespaceToken, string(len(w)), s)
		case isText(r):
			s := l.loc
			value := l.while(isText)
			c <- l.newToken(TextToken, value, s)
		default:
			c <- l.newToken(ErrorToken, fmt.Sprintf("unknown rune %v", string(r)), l.loc)
			l.increment()
		}
	}
	close(c)
}
