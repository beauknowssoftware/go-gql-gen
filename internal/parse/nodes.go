package parse

type Node interface {
	Children() []Node
	Loc() Loc
}

type NodeLoc struct {
	NodeLoc Loc
}

func (n NodeLoc) Loc() Loc {
	return n.Loc()
}

type DocumentNode struct {
	NodeLoc
	Definitions []Node
}

func (n DocumentNode) Children() []Node {
	children := make([]Node, len(n.Definitions), len(n.Definitions))
	for i, n := range n.Definitions {
		children[i] = n
	}
	return children
}

type DefinitionNode interface {
	Node
}

type TypeNode struct {
	NodeLoc
	Name   string
	Fields []Node
}

func (n TypeNode) Children() []Node {
	children := make([]Node, len(n.Fields), len(n.Fields))
	for i, n := range n.Fields {
		children[i] = n
	}
	return children
}

type SchemaNode struct {
	NodeLoc
	Fields []Node
}

func (n SchemaNode) Children() []Node {
	children := make([]Node, len(n.Fields), len(n.Fields))
	for i, n := range n.Fields {
		children[i] = n
	}
	return children
}

type FieldNode struct {
	NodeLoc
	Name     string
	Type     string
	Required bool
	Params   []Node
}

func (n FieldNode) Children() []Node {
	children := make([]Node, len(n.Params), len(n.Params))
	for i, n := range n.Params {
		children[i] = n
	}
	return children
}

type ParamNode struct {
	NodeLoc
	Name     string
	Type     string
	Required bool
}

func (n ParamNode) Children() []Node {
	return nil
}

type TokenNode struct {
	NodeLoc
	TokenType TokenType
	Value string
}

func (n TokenNode) Children() []Node {
	return nil
}

