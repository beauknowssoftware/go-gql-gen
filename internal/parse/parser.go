package parse

import (
	"errors"
	"fmt"
)

type Parser struct {
	l       Lexer
	c       chan Token
	current Token
}

func New(l Lexer) Parser {
	c := make(chan Token)
	return Parser{l: l, c: c}
}

func (p *Parser) consume() {
	p.current = <-p.c
	for p.current.TokenType == WhitespaceToken {
		p.current = <-p.c
	}
}

func (p Parser) nodeLoc() NodeLoc {
	return NodeLoc{p.current.Loc}
}

type parserPart func(*Parser) (Node, error)

func nothing(p *Parser) (Node, error) {
	return nil, nil
}

func token(tt TokenType, v string) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		if p.current.TokenType != tt || p.current.Value != v {
			return nil, fmt.Errorf("expected %v (%v) token got %v (%v) token", tt, v, p.current.TokenType, p.current.Value)
		}
		p.consume()

		return TokenNode{nodeLoc, tt, v}, nil
	}
}

func tokenType(tt TokenType) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		if p.current.TokenType != tt {
			return nil, fmt.Errorf("expected %v token got %v token", tt, p.current.TokenType)
		}
		t := p.current
		p.consume()

		return TokenNode{nodeLoc, tt, t.Value}, nil
	}
}

func maybe(pp parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		n, err := pp(p)
		if err != nil {
			return nil, nil
		}
		return n, nil
	}
}

func seq(trans func(NodeLoc, ...Node) Node, pps ...parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		nodes := make([]Node, len(pps), len(pps))

		for i, pp := range pps {
			n, err := pp(p)
			if err != nil {
				return nil, err
			}
			nodes[i] = n
		}

		n := trans(nodeLoc, nodes...)
		return n, nil
	}
}

func multi(pp parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		nodes := make([]Node, 0)

		for {
			n, err := pp(p)
			if err != nil {
				break
			}
			nodes = append(nodes, n)
		}

		return MultiNode{nodeLoc, nodes}, nil
	}
}

func multiSep(pp parserPart, sep parserPart) parserPart {
	maybeSep := maybe(sep)

	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		nodes := make([]Node, 0)

		n, err := pp(p)
		if err != nil {
			return MultiNode{nodeLoc, nodes}, nil
		}
		nodes = append(nodes, n)

		for {
			if n, _ := maybeSep(p); n == nil {
				break
			}

			n, err := pp(p)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, n)
		}

		return MultiNode{nodeLoc, nodes}, nil
	}
}

func choice(pps ...parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		for _, pp := range pps {
			n, err := pp(p)
			if err == nil {
				return n, nil
			}
		}
		return nil, errors.New("cannot match token")
	}
}

var parseParameter = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	return ParamNode{
		nodeLoc,
		nodes[0].(TokenNode).Value,
		nodes[2].(TokenNode).Value,
		nodes[3] != nil,
	}
}, identifier, token(ColonToken, ""), identifier, maybe(required))

var parseParameterList = multiSep(parseParameter, token(CommaToken, ""))

var parseParameters = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	return nodes[1]
}, token(LeftParenToken, ""), parseParameterList, token(RightParenToken, ""))

var identifier = tokenType(TextToken)
var required = token(BangToken, "")

var parseField = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	f := FieldNode{
		nodeLoc,
		nodes[0].(TokenNode).Value,
		nodes[3].(TokenNode).Value,
		nodes[4] != nil,
		nil,
	}
	if nodes[1] != nil {
		f.Params = nodes[1].(MultiNode).Nodes
	}
	return f
}, identifier, maybe(parseParameters), token(ColonToken, ""), identifier, maybe(required))

var schemaKeyword = token(TextToken, "schema")

var parseSchema = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	return SchemaNode{
		nodeLoc,
		nodes[2].(MultiNode).Nodes,
	}
}, schemaKeyword, token(LeftCurlyToken, ""), multi(parseField), token(RightCurlyToken, ""))

var typeKeyword = token(TextToken, "type")

var parseType = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	return TypeNode{
		nodeLoc,
		nodes[1].(TokenNode).Value,
		nodes[3].(MultiNode).Nodes,
	}
}, typeKeyword, identifier, token(LeftCurlyToken, ""), multi(parseField), token(RightCurlyToken, ""))

var parseDefinition = choice(parseType, parseSchema)

var parseDocument = seq(func(nodeLoc NodeLoc, nodes ...Node) Node {
	return DocumentNode{nodeLoc, nodes[0].(MultiNode).Nodes}
}, multi(parseDefinition))

type Error struct {
	Error error
	Token Token
}

func (p *Parser) Parse() (Node, *Error) {
	go p.l.Lex(p.c)

	p.consume()
	d, err := parseDocument(p)
	if err != nil {
		return nil, &Error{err, p.current}
	}
	return d, nil
}
