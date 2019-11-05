package parse

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
}

func (s *traverseStack) pop() Node {
	if s.len() > 0 {
		n := s.head()
		s.nodes = s.tail()
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
			children := n.Children()
			for i := len(children) - 1; i >= 0; i-- {
				s.push(children[i])
			}
		}
	}
}
