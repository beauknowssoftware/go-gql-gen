package parse

import (
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
		err := p.tryConsume(tt, v)
		if err != nil {
			return nil, err
		}
		return TokenNode{nodeLoc, tt,v}, nil
	}
}

func tokenType(tt TokenType) parserPart {
	return func(p *Parser) (Node, error) {
		nodeLoc := p.nodeLoc()
		t, err := p.tryConsumeType(tt)
		if err != nil {
			return nil, err
		}
		return TokenNode{nodeLoc, tt,t.Value}, nil
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

func seq(trans func(...Node) Node, pps ...parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		nodes := make([]Node, len(pps), len(pps))

		for i, pp := range pps {
			n, err := pp(p)
			if err != nil {
				return nil, err
			}
			nodes[i] = n
		}

		n := trans(nodes...)
		return n, nil
	}
}

func multi(trans func(...Node) Node, pp parserPart) parserPart {
	return func(p *Parser) (Node, error) {
		nodes := make([]Node, 0)

		for {
			n, err := pp(p)
			if err != nil {
				break
			}
			nodes = append(nodes, n)
		}

		n := trans(nodes...)
		return n, nil
	}
}

func choice(pps...parserPart) parserPart {
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

func (p *Parser) tryConsumeType(tt TokenType) (*Token, error) {
	if p.current.TokenType != tt {
		return nil, fmt.Errorf("expected %v token got %v token", tt, p.current.TokenType)
	}
	t := p.current
	p.consume()
	return &t, nil
}

func (p *Parser) maybeConsume(tt TokenType, value string) bool {
	if p.current.TokenType == tt && p.current.Value == value {
		p.consume()
		return true
	}
	return false
}

func (p *Parser) maybeConsumeType(tt TokenType) (Token, bool) {
	if p.current.TokenType == tt {
		t := p.current
		p.consume()
		return t, true
	}
	return p.current, false
}

func (p *Parser) tryConsume(tt TokenType, value string) error {
	if p.current.TokenType != tt || p.current.Value != value {
		return fmt.Errorf("expected %v (%v) token got %v (%v) token", tt, value, p.current.TokenType, p.current.Value)
	}
	p.consume()
	return nil
}

func (p *Parser) parseParam() (*ParamNode, error) {
	name, error := p.tryConsumeType(TextToken)
	if error != nil {
		return nil, error
	}

	if err := p.tryConsume(ColonToken, ""); err != nil {
		return nil, err
	}

	typ, err := p.tryConsumeType(TextToken)
	if err != nil {
		return nil, err
	}

	required := p.maybeConsume(BangToken, "")

	n := ParamNode{
		Name:     name.Value,
		Type:     typ.Value,
		Required: required,
	}
	return &n, nil
}

func (p *Parser) maybeParseParams() ([]Node, error) {
	if !p.maybeConsume(LeftParenToken, "") {
		return nil, nil
	}

	params := make([]Node, 0)

	n, err := p.parseParam()
	if err != nil {
		return nil, err
	}
	params = append(params, *n)

	for {
		if !p.maybeConsume(CommaToken, "") {
			break
		}

		n, err := p.parseParam()
		if err != nil {
			return nil, err
		}
		params = append(params, *n)
	}

	if err := p.tryConsume(RightParenToken, ""); err != nil {
		return nil, err
	}

	return params, nil
}

func (p *Parser) maybeParseField() (*FieldNode, error) {
	name, parsed := p.maybeConsumeType(TextToken)
	if !parsed {
		return nil, nil
	}

	params, err := p.maybeParseParams()
	if err != nil {
		return nil, err
	}

	if err := p.tryConsume(ColonToken, ""); err != nil {
		return nil, err
	}

	typ, err := p.tryConsumeType(TextToken)
	if err != nil {
		return nil, err
	}

	required := p.maybeConsume(BangToken, "")

	f := FieldNode{
		Name:     name.Value,
		Type:     typ.Value,
		Params:   params,
		Required: required,
	}
	return &f, nil
}

func (p *Parser) parseFields() ([]Node, error) {
	fields := make([]Node, 0)
	for {
		f, err := p.maybeParseField()
		if err != nil {
			return nil, err
		}
		if f == nil {
			break
		}
		fields = append(fields, *f)
	}
	return fields, nil
}

func (p *Parser) maybeParseSchema() (*SchemaNode, error) {
	nodeLoc := p.nodeLoc()

	if !p.maybeConsume(TextToken, "schema") {
		return nil, nil
	}

	if err := p.tryConsume(LeftCurlyToken, ""); err != nil {
		return nil, err
	}

	fields, err := p.parseFields()
	if err != nil {
		return nil, err
	}

	if err := p.tryConsume(RightCurlyToken, ""); err != nil {
		return nil, err
	}

	return &SchemaNode{nodeLoc, fields}, nil
}

func (p *Parser) maybeParseType() (*TypeNode, error) {
	nodeLoc := p.nodeLoc()

	if !p.maybeConsume(TextToken, "type") {
		return nil, nil
	}

	name, err := p.tryConsumeType(TextToken)
	if err != nil {
		return nil, err
	}

	if err := p.tryConsume(LeftCurlyToken, ""); err != nil {
		return nil, err
	}

	fields, err := p.parseFields()
	if err != nil {
		return nil, err
	}

	if err := p.tryConsume(RightCurlyToken, ""); err != nil {
		return nil, err
	}

	return &TypeNode{nodeLoc, name.Value, fields}, nil
}

func (p *Parser) maybeParseDefinition() (DefinitionNode, error) {
	t, err := p.maybeParseType()
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	s, err := p.maybeParseSchema()
	if err != nil {
		return nil, err
	}
	if s != nil {
		return s, nil
	}

	return nil, nil
}

func (p *Parser) parseDefinitions() ([]Node, error) {
	definitions := make([]Node, 0)
	for {
		d, err := p.maybeParseDefinition()
		if err != nil {
			return nil, err
		}
		if d == nil {
			break
		}
		definitions = append(definitions, d)
	}
	return definitions, nil
}

func (p *Parser) parseDocument() (*DocumentNode, error) {
	nodeLoc := p.nodeLoc()
	definitions, err := p.parseDefinitions()
	if err != nil {
		return nil, err
	}
	return &DocumentNode{nodeLoc, definitions}, nil
}

type Error struct {
	Error error
	Token Token
}

func (p *Parser) Parse() (*DocumentNode, *Error) {
	go p.l.Lex(p.c)

	p.consume()
	d, err := p.parseDocument()
	if err != nil {
		return nil, &Error{err, p.current}
	}
	return d, nil
}
