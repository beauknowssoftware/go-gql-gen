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

	n := ParamNode{
		Name: name.Value,
		Type: typ.Value,
	}
	return &n, nil
}

func (p *Parser) maybeParseParams() ([]ParamNode, error) {
	if !p.maybeConsume(LeftParenToken, "") {
		return nil, nil
	}

	params := make([]ParamNode, 0)

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

func (p *Parser) parseFields() ([]FieldNode, error) {
	fields := make([]FieldNode, 0)
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

	return &SchemaNode{fields}, nil
}

func (p *Parser) maybeParseType() (*TypeNode, error) {
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

	return &TypeNode{name.Value, fields}, nil
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

func (p *Parser) parseDefinitions() ([]DefinitionNode, error) {
	definitions := make([]DefinitionNode, 0)
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
	definitions, err := p.parseDefinitions()
	if err != nil {
		return nil, err
	}
	return &DocumentNode{definitions}, nil
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
