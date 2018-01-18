package parsec

// NonTerminal will be used by AST methods to construct intermediate nodes.
// Note that user supplied ASTNodify callback can construct a different
// type of intermediate node that confirms to Queryable interface.
type NonTerminal struct {
	Name       string      // contains terminal's token type
	Children   []Queryable // list of children to this node.
	Attributes map[string][]string
}

// NewNonTerminal create and return a new NonTerminal instance.
func NewNonTerminal(name string) *NonTerminal {
	nt := &NonTerminal{
		Name:       name,
		Children:   make([]Queryable, 0),
		Attributes: make(map[string][]string),
	}
	nt.SetAttribute("class", "nonterm")
	return nt
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

// SetAttribute implement Queryable interface.
func (nt *NonTerminal) SetAttribute(attrname, value string) Queryable {
	if nt.Attributes == nil {
		nt.Attributes = make(map[string][]string)
	}
	values, ok := nt.Attributes[attrname]
	if ok == false {
		values = []string{}
	}
	values = append(values, value)
	nt.Attributes[attrname] = values
	return nt
}

// GetAttribute implement Queryable interface.
func (nt *NonTerminal) GetAttribute(attrname string) []string {
	if nt.Attributes == nil {
		return nil
	} else if values, ok := nt.Attributes[attrname]; ok {
		return values
	}
	return nil
}

// GetAttributes implement Queryable interface.
func (nt *NonTerminal) GetAttributes() map[string][]string {
	return nt.Attributes
}
