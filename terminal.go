package parsec

// Terminal structure can be used to construct a terminal
// ParsecNode.
type Terminal struct {
	Name       string // contains terminal's token type
	Value      string // value of the terminal
	Position   int    // Offset into the text stream where token was identified
	Attributes map[string][]string
}

// NewTerminal create a new Terminal instance.
func NewTerminal(name, value string, position int) *Terminal {
	t := &Terminal{
		Name:       name,
		Value:      value,
		Position:   position,
		Attributes: make(map[string][]string),
	}
	t.SetAttribute("class", "term")
	return t
}

// GetName implement Queryable interface.
func (t *Terminal) GetName() string {
	return t.Name
}

// IsTerminal implement Queryable interface.
func (t *Terminal) IsTerminal() bool {
	return true
}

// GetValue implement Queryable interface.
func (t *Terminal) GetValue() string {
	return t.Value
}

// GetChildren implement Queryable interface.
func (t *Terminal) GetChildren() []Queryable {
	return nil
}

// GetPosition implement Queryable interface.
func (t *Terminal) GetPosition() int {
	return t.Position
}

// SetAttribute implement Queryable interface.
func (t *Terminal) SetAttribute(attrname, value string) Queryable {
	if t.Attributes == nil {
		t.Attributes = make(map[string][]string)
	}
	values, ok := t.Attributes[attrname]
	if ok == false {
		values = []string{}
	}
	values = append(values, value)
	t.Attributes[attrname] = values
	return t
}

// GetAttribute implement Queryable interface.
func (t *Terminal) GetAttribute(attrname string) []string {
	if t.Attributes == nil {
		return nil
	} else if values, ok := t.Attributes[attrname]; ok {
		return values
	}
	return nil
}

// GetAttributes implement Queryable interface.
func (t *Terminal) GetAttributes() map[string][]string {
	return t.Attributes
}

// MaybeNone place holder type used be Maybe combinator if parser does not
// match the input text.
type MaybeNone string

//---- implement Queryable interface

// GetName implement Queryable interface.
func (mn MaybeNone) GetName() string {
	return string(mn)
}

// IsTerminal implement Queryable interface.
func (mn MaybeNone) IsTerminal() bool {
	return true
}

// GetValue implement Queryable interface.
func (mn MaybeNone) GetValue() string {
	return ""
}

// GetChildren implement Queryable interface.
func (mn MaybeNone) GetChildren() []Queryable {
	return nil
}

// GetPosition implement Queryable interface.
func (mn MaybeNone) GetPosition() int {
	return -1
}

// SetAttribute implement Queryable interface.
func (mn MaybeNone) SetAttribute(attrname, value string) Queryable {
	return mn
}

// GetAttribute implement Queryable interface.
func (mn MaybeNone) GetAttribute(attrname string) []string {
	return nil
}

// GetAttributes implement Queryable interface.
func (mn MaybeNone) GetAttributes() map[string][]string {
	return nil
}
