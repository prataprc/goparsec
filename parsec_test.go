package parsec

import "fmt"
import "reflect"
import "testing"

var _ = fmt.Sprintf("dummy")

func TestStrEOF(t *testing.T) {
	word := String()
	Y := Many(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		word)

	input := `"alpha" "beta" "gamma"`
	s := NewScanner([]byte(input))

	root, _ := Y(s)
	nodes := root.([]ParsecNode)
	ref := []ParsecNode{"\"alpha\"", "\"beta\"", "\"gamma\""}
	if !reflect.DeepEqual(nodes, ref) {
		t.Fatal(nodes)
	}
}

func allTokens(ns []ParsecNode) ParsecNode {
	return ns
}

func TestMany(t *testing.T) {
	w := Token("\\w+", "W")
	m := Many(nil, w)
	s := NewScanner([]byte("one two stop"))
	v, e := m(s)
	if v == nil {
		t.Errorf("Many() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 3 {
		t.Errorf("Many() didn't match all words %q", v)
	}

	w = Token("\\w+", "W")
	m = Many(nil, w, Atom(",", "COMMA"))
	s = NewScanner([]byte("one,two"))
	v, e = m(s)
	if v == nil {
		t.Errorf("Many() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 2 {
		t.Errorf("Many() didn't match all words %q", v)
	}
}

func TestManyUntilNoStop(t *testing.T) {
	w := Token("\\w+", "W")
	u := Token("stop", "S")
	m := ManyUntil(nil, w, u)
	s := NewScanner([]byte("one two three"))
	v, e := m(s)
	if v == nil {
		t.Errorf("ManyUntil() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 3 {
		t.Errorf("ManyUntil() didn't match all words %q", v)
	}
}

func TestManyUntilStop(t *testing.T) {
	w := Token("\\w+", "W")
	u := Token("stop", "S")
	m := ManyUntil(nil, w, u)
	s := NewScanner([]byte("one two stop"))
	v, e := m(s)
	if v == nil {
		t.Errorf("ManyUntil() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 2 {
		t.Errorf("ManyUntil() didn't stop %q", v)
	}
}

func TestManyUntilNoStopSep(t *testing.T) {
	w := Token("\\w+", "W")
	u := Token("stop", "S")
	z := Token("z", "Z")
	m := ManyUntil(nil, w, z, u)
	s := NewScanner([]byte("one z two z three"))
	v, e := m(s)
	if v == nil {
		t.Errorf("ManyUntil() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 3 {
		t.Errorf("ManyUntil() didn't match all words %q", v)
	}
}

func TestManyUntilStopSep(t *testing.T) {
	w := Token("\\w+", "W")
	u := Token("stop", "S")
	z := Token("z", "Z")
	m := ManyUntil(nil, w, z, u)
	s := NewScanner([]byte("one z two z stop"))
	v, e := m(s)
	if v == nil {
		t.Errorf("ManyUntil() didn't match %q", e)
	} else if len(v.([]ParsecNode)) != 2 {
		t.Errorf("ManyUntil() didn't stop %q", v)
	}
}

func TestKleene(t *testing.T) {
	y := Kleene(nil, Token("\\w+", "W"))

	s := NewScanner([]byte("one two stop"))
	node, e := y(s)
	if node == nil {
		t.Errorf("Kleene() didn't match %q", e)
	} else if len(node.([]ParsecNode)) != 3 {
		t.Errorf("Kleene() didn't match all words %q", node)
	}

	y = Kleene(nil, Token("\\w+", "W"), Atom(",", "COMMA"))
	s = NewScanner([]byte("one,two"))
	node, _ = y(s)
	if node == nil {
		t.Errorf("Kleene() didn't match %q", e)
	} else if len(node.([]ParsecNode)) != 2 {
		t.Errorf("Kleene() didn't match all words %q", node)
	}
}

func TestForwardReference(t *testing.T) {
	var ycomma Parser

	y := Kleene(nil, Maybe(nil, Token("\\w+", "WORD")), &ycomma)
	ycomma = Atom(",", "COMMA")
	s := NewScanner([]byte("one,two,,three"))
	node, _ := y(s)
	nodes := node.([]ParsecNode)
	if len(nodes) != 4 {
		t.Errorf("expected length to be 4")
	} else if _, ok := nodes[2].(MaybeNone); ok == false {
		t.Errorf("expected MissingNone")
	}
}
