package parsec

// NonTerminal will be used by AST objects to construct intermediate nodes.
// Note that user supplied ASTNodify callback can construct a different
// type of intermediate node that confirms to Queryable interface.
type NonTerminal struct {
	Name     string      // contains terminal's token type
	Children []Queryable // list of children to this node.
}

// GetName implement Queryable interface.
func (nt *NonTerminal) GetName() string {
	return nt.Name
}

// IsTerminal implement Queryable interface.
func (nt *NonTerminal) IsTerminal() bool {
	return false
}

// GetValue implement Queryable interface.
func (nt *NonTerminal) GetValue() string {
	value := ""
	for _, c := range nt.Children {
		value += c.GetValue()
	}
	return value
}

// GetChildren implement Queryable interface.
func (nt *NonTerminal) GetChildren() []Queryable {
	return nt.Children
}

// GetPosition implement Queryable interface.
func (nt *NonTerminal) GetPosition() int {
	if nodes := nt.GetChildren(); len(nodes) > 0 {
		return nodes[0].GetPosition()
	}
	return 0
}
