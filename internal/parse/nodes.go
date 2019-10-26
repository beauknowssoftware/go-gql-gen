package parse

type Node interface {
	Children() []Node
}

type DocumentNode struct {
	Definitions []DefinitionNode
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
	Name   string
	Fields []FieldNode
}

func (n TypeNode) Children() []Node {
	children := make([]Node, len(n.Fields), len(n.Fields))
	for i, n := range n.Fields {
		children[i] = n
	}
	return children
}

type SchemaNode struct {
	Fields []FieldNode
}

func (n SchemaNode) Children() []Node {
	children := make([]Node, len(n.Fields), len(n.Fields))
	for i, n := range n.Fields {
		children[i] = n
	}
	return children
}

type FieldNode struct {
	Name     string
	Type     string
	Required bool
	Params   []ParamNode
}

func (n FieldNode) Children() []Node {
	children := make([]Node, len(n.Params), len(n.Params))
	for i, n := range n.Params {
		children[i] = n
	}
	return children
}

type ParamNode struct {
	Name     string
	Type     string
	Required bool
}

func (n ParamNode) Children() []Node {
	return nil
}
