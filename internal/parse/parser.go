package parse

import (
	"errors"
	"fmt"
)

type Parser struct {
	l      Lexer
	tokens []Token
	i      int
}

func New(l Lexer) Parser {
	return Parser{l: l}
}

func (p *Parser) current() Token {
	return p.tokens[p.i]
}

func (p *Parser) consume() {
	if p.i < len(p.tokens)-1 {
		p.i++
		for p.current().TokenType == WhitespaceToken {
			p.i++
		}
	}
}

func (p Parser) nodeLoc() NodeLoc {
	return NodeLoc{p.current().Loc}
}

type parserPart func(*Parser) (Node, error)

func nothing(p *Parser) (Node, error) {
	return nil, nil
}

func keyword(v string) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		if p.current().TokenType != TextToken || p.current().Value != v {
			return nil, fmt.Errorf("expected keyword %v keyword got %v (%v) token", v, p.current().TokenType, p.current().Value)
		}
		p.consume()

		return TokenNode{nodeLoc, TextToken, v}, nil
	}
}

func token(tt TokenType) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()

		if p.current().TokenType != tt {
			return nil, fmt.Errorf("expected %v token got %v token", tt, p.current().TokenType)
		}
		t := p.current()
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

type transformer func(NodeLoc, ...Node) (Node, error)

func seq(trans transformer, pps ...parserPart) parserPart {
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

		return trans(nodeLoc, nodes...)
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
		return nil, errors.New("cannot match keyword")
	}
}

var parseParameter = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return ParamNode{
		nodeLoc,
		nodes[0].(TokenNode).Value,
		nodes[2],
	}, nil
}, identifier, token(ColonToken), parseType)

var parseParameterList = multiSep(parseParameter, token(CommaToken))

var parseParameters = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return nodes[1], nil
}, token(LeftParenToken), parseParameterList, token(RightParenToken))

var identifier = token(TextToken)
var required = token(BangToken)

var parseDirective = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return DirectiveNode{
		nodeLoc,
		nodes[1].(TokenNode).Value,
	}, nil
}, token(AtToken), identifier)

var parseArrayType = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return TypeNode{
		nodeLoc,
		nodes[1].(TokenNode).Value,
		nodes[4] != nil,
		true,
		nodes[2] != nil,
	}, nil
}, token(LeftBracketToken), identifier, maybe(required), maybe(token(RightBracketToken)), maybe(required))

var parseSingleType = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return TypeNode{
		nodeLoc,
		nodes[0].(TokenNode).Value,
		nodes[1] != nil,
		false,
		false,
	}, nil
}, identifier, maybe(required))

var parseType = choice(parseArrayType, parseSingleType)

var parseField = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	f := FieldNode{
		nodeLoc,
		nodes[0].(TokenNode).Value,
		nodes[3],
		nil,
		nil,
	}
	if nodes[1] != nil {
		f.Params = nodes[1].(MultiNode).Nodes
	}
	if directives := nodes[4].(MultiNode).Nodes; len(directives) > 0 {
		f.Directives = directives
	}
	return f, nil
}, identifier, maybe(parseParameters), token(ColonToken), parseType, multi(parseDirective))

var schemaKeyword = keyword("schema")

var parseSchema = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return SchemaNode{
		nodeLoc,
		nodes[2].(MultiNode).Nodes,
	}, nil
}, schemaKeyword, token(LeftCurlyToken), multi(parseField), token(RightCurlyToken))

var typeKeyword = keyword("type")

var parseTypeDef = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return TypeDefNode{
		nodeLoc,
		nodes[1].(TokenNode).Value,
		nodes[3].(MultiNode).Nodes,
		false,
	}, nil
}, typeKeyword, identifier, token(LeftCurlyToken), multi(parseField), token(RightCurlyToken))

var inputKeyword = keyword("input")

var parseInput = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return TypeDefNode{
		nodeLoc,
		nodes[1].(TokenNode).Value,
		nodes[3].(MultiNode).Nodes,
		true,
	}, nil
}, inputKeyword, identifier, token(LeftCurlyToken), multi(parseField), token(RightCurlyToken))

var parseDefinition = choice(parseTypeDef, parseInput, parseSchema)

var parseDocument = seq(func(nodeLoc NodeLoc, nodes ...Node) (Node, error) {
	return DocumentNode{nodeLoc, nodes[0].(MultiNode).Nodes}, nil
}, multi(parseDefinition), token(EOFToken))

type Error struct {
	Error error
	Token Token
}

func (p *Parser) Parse() (Node, *Error) {
	c := make(chan Token)
	go p.l.Lex(c)
	for t := range c {
		p.tokens = append(p.tokens, t)
	}

	d, err := parseDocument(p)
	if err != nil {
		return nil, &Error{err, p.current()}
	}
	if p.i < len(p.tokens)-1 {
		return nil, &Error{errors.New("stopped parsing prematurely"), p.current()}
	}
	return d, nil
}
