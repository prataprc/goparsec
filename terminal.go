package parsec

// Terminal structure can be used to construct a terminal
// ParsecNode.
type Terminal struct {
	Name     string // contains terminal's token type
	Value    string // value of the terminal
	Position int    // Offset into the text stream where token was identified
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
