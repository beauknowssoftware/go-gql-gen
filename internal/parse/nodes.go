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

type MultiNode struct {
	NodeLoc
	Nodes []Node
}

func (n MultiNode) Children() []Node {
	return n.Nodes
}

type DefinitionNode interface {
	Node
}

type TypeDefNode struct {
	NodeLoc
	Name   string
	Fields []Node
}

func (n TypeDefNode) Children() []Node {
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
	Name       string
	Type       Node
	Params     []Node
	Directives []Node
}

func (n FieldNode) Children() []Node {
	paramCount := len(n.Params)
	childCount := paramCount + len(n.Directives) + 1
	children := make([]Node, childCount, childCount)
	children[0] = n.Type
	for i, n := range n.Params {
		children[i+1] = n
	}
	for i, n := range n.Directives {
		children[i+1+paramCount] = n
	}
	return children
}

type TypeNode struct {
	NodeLoc
	Name            string
	Required        bool
	Multiple        bool
	NonNullElements bool
}

func (n TypeNode) Children() []Node {
	return nil
}

type ParamNode struct {
	NodeLoc
	Name string
	Type Node
}

func (n ParamNode) Children() []Node {
	return []Node{n.Type}
}

type DirectiveNode struct {
	NodeLoc
	Name string
}

func (n DirectiveNode) Children() []Node {
	return nil
}

type TokenNode struct {
	NodeLoc
	TokenType TokenType
	Value     string
}

func (n TokenNode) Children() []Node {
	return nil
}
