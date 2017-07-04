package parsec

import "fmt"
import "reflect"
import "testing"

var _ = fmt.Sprintf("dummy")

func TestAnd(t *testing.T) {
	y := And(func(ns []ParsecNode) ParsecNode {
		return nil
	}, Atom("hello", "TERM"))
	s := NewScanner([]byte("hello"))
	node, s := y(s)
	if node != nil {
		t.Errorf("expected nil")
	}
	ss := s.(*SimpleScanner)
	if str := string(ss.buf[ss.cursor:]); str != "hello" {
		t.Errorf("expected %q, got %q", "hello", str)
	}
}

func TestOrdChoice(t *testing.T) {
	y := OrdChoice(func(ns []ParsecNode) ParsecNode {
		return nil
	}, Atom("hello", "TERM"))
	s := NewScanner([]byte("hello"))
	node, s := y(s)
	if node != nil {
		t.Errorf("expected nil")
	}
	ss := s.(*SimpleScanner)
	if str := string(ss.buf[ss.cursor:]); str != "hello" {
		t.Errorf("expected %q, got %q", "hello", str)
	}
}

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

func TestMany(t *testing.T) {
	w := Token("\\w+", "W")
	y := Many(nil, w)
	s := NewScanner([]byte("one two stop"))
	node, e := y(s)
	if node == nil {
		t.Errorf("Many() didn't match %q", e)
	} else if len(node.([]ParsecNode)) != 3 {
		t.Errorf("Many() didn't match all words %q", node)
	}

	w = Token("\\w+", "W")
	y = Many(nil, w, Atom(",", "COMMA"))
	s = NewScanner([]byte("one,two"))
	node, e = y(s)
	if node == nil {
		t.Errorf("Many() didn't match %q", e)
	} else if len(node.([]ParsecNode)) != 2 {
		t.Errorf("Many() didn't match all words %q", node)
	}

	// Return nil
	w = Token("\\w+", "W")
	y = Many(
		func(_ []ParsecNode) ParsecNode { return nil },
		w, Atom(",", "COMMA"),
	)
	s = NewScanner([]byte("one,two"))
	node, s = y(s)
	if node != nil {
		t.Errorf("Many() didn't match %q", e)
	}
	ss := s.(*SimpleScanner)
	if str := string(ss.buf[ss.cursor:]); str != "one,two" {
		t.Errorf("expected %q, got %q", "one,two", str)
	}

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		Many(nil, Token("\\w+", "W"), Atom(",", "COMMA"), Atom(",", "COMMA"))
	}()
}

func TestManyUntil(t *testing.T) {
	// Return nil
	w := Token("\\w+", "W")
	y := ManyUntil(
		func(_ []ParsecNode) ParsecNode { return nil },
		w, Atom(",", "COMMA"),
	)
	s := NewScanner([]byte("one,two"))
	node, s := y(s)
	if node != nil {
		t.Errorf("ManyUntil() expected nil")
	}
	ss := s.(*SimpleScanner)
	if str := string(ss.buf[ss.cursor:]); str != "one,two" {
		t.Errorf("expected %q, got %q", "one,two", str)
	}

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		ManyUntil(nil, Token("\\w+", "W"))
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		ManyUntil(
			nil,
			Token("\\w+", "W"), Atom(",", "COMMA"), Atom(",", "COMMA"),
			Atom(",", "COMMA"),
		)
	}()
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

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		Kleene(nil, Token("\\w+", "W"), Atom(",", "COMMA"), Atom(",", "COMMA"))
	}()
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

	// nil return
	y = Maybe(
		func(_ []ParsecNode) ParsecNode {
			return nil
		}, Token("\\w+", "WORD"),
	)
	s = NewScanner([]byte("one"))
	node, s = y(s)
	if node != MaybeNone("missing") {
		t.Errorf("expected %v, got %v", node, MaybeNone("missing"))
	}
}

func allTokens(ns []ParsecNode) ParsecNode {
	return ns
}
