package parse

import "fmt"

type traverseStack struct {
	nodes []Node
}

func (s traverseStack) len() int {
	return len(s.nodes)
}

func (s traverseStack) head() Node {
	if s.len() > 0 {
		return s.nodes[s.len()-1]
	}
	return nil
}

func (s traverseStack) tail() []Node {
	if s.len() > 1 {
		return s.nodes[:s.len()-1]
	}
	return nil
}

func (s *traverseStack) push(v Node) {
	s.nodes = append(s.nodes, v)
	fmt.Println(s.nodes)
}

func (s *traverseStack) pop() Node {
	if s.len() > 0 {
		n := s.head()
		s.nodes = s.tail()
		fmt.Println(s.nodes)
		return n
	}
	return nil
}

func Traverse(ast Node, op func(Node) bool) {
	var s traverseStack

	s.push(ast)

	for s.len() > 0 {
		n := s.pop()
		if op(n) {
			switch n := n.(type) {
			case DocumentNode:
				for i := len(n.Definitions) - 1; i >= 0; i-- {
					s.push(n.Definitions[i])
				}
			case *TypeNode:
				for i := len(n.Fields) - 1; i >= 0; i-- {
					s.push(n.Fields[i])
				}
			case *SchemaNode:
				for i := len(n.Fields) - 1; i >= 0; i-- {
					s.push(n.Fields[i])
				}
			case FieldNode:
				for i := len(n.Params) - 1; i >= 0; i-- {
					s.push(n.Params[i])
				}
			}
		}
	}
}
