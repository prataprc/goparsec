package parsec

import "testing"
import "reflect"
import "fmt"

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
	out := []string{
		nodes[0].(*Terminal).Value,
		nodes[1].(*Terminal).Value,
		nodes[2].(*Terminal).Value,
	}
	ref := []string{`"alpha"`, `"beta"`, `"gamma"`}
	if !reflect.DeepEqual(out, ref) {
		t.Errorf("expected %v, got %v", ref, out)
	}
}

func TestParsecAnd(t *testing.T) {
	Y := And(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String(), Int())

	s := NewScanner([]byte(`"hello" 42`))
	nt, s := Y(s)
	if s.Endof() == false {
		t.Errorf("expected endof string")
	}
	for _, node := range nt.([]ParsecNode) {
		terminal := node.(*Terminal)
		switch terminal.Value {
		case "STRING":
			if terminal.Value != `"hello"` {
				t.Errorf("expected %v, got %v", `"hello"`, terminal.Value)
			} else if terminal.Position != 0 {
				t.Errorf("expected %v, got %v", 0, terminal.Position)
			}
		case "INT":
			if terminal.Value != `42` {
				t.Errorf("expected %v, got %v", 42, terminal.Value)
			} else if terminal.Position != 8 {
				t.Errorf("expected %v, got %v", 8, terminal.Position)
			}
		}
	}

	// negative case
	s = NewScanner([]byte(`true "hello" 42`))
	nt, s = Y(s)
	if nt != nil {
		t.Errorf("expected nil")
	}
}

func TestParsecOrdChoice(t *testing.T) {
	Y := OrdChoice(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String(), Int())

	s := NewScanner([]byte(`"hello"`))
	nt, _ := Y(s)
	terminal := nt.([]ParsecNode)[0].(*Terminal)
	if terminal.Value != `"hello"` {
		t.Errorf("expected %v, got %v", `"hello"`, terminal.Value)
	} else if terminal.Position != 0 {
		t.Errorf("expected %v, got %v", 0, terminal.Position)
	}

	s = NewScanner([]byte(`42`))
	nt, _ = Y(s)
	terminal = nt.([]ParsecNode)[0].(*Terminal)
	if terminal.Value != `42` {
		t.Errorf("expected %v, got %v", `42`, terminal.Value)
	} else if terminal.Position != 0 {
		t.Errorf("expected %v, got %v", 0, terminal.Position)
	}

	// negative case
	s = NewScanner([]byte(`  true 42`))
	nt, _ = Y(s)
	if nt != nil {
		t.Errorf("expected nil")
	}
}

func TestParsecKleene(t *testing.T) {
	Y := Kleene(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String())

	s := NewScanner([]byte(`"hello" "world" 42`))
	nt, s := Y(s)
	nodes := nt.([]ParsecNode)
	if len(nodes) != 2 {
		t.Errorf("expected %v, got %v", 2, len(nodes))
	}
	ref := []interface{}{`"hello"`, 0, `"world"`, 8}
	out := []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
		nodes[1].(*Terminal).Value,
		nodes[1].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		Kleene(nil, String(), Int(), Int())
	}()
}

func TestParsecKleeneSep(t *testing.T) {
	Y := Kleene(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String(), Token(",", "COMMA"))

	s := NewScanner([]byte(`"hello", "world" 42`))
	nt, s := Y(s)
	nodes := nt.([]ParsecNode)
	if len(nodes) != 2 {
		t.Errorf("expected %v, got %v", 2, len(nodes))
	}
	ref := []interface{}{`"hello"`, 0, `"world"`, 9}
	out := []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
		nodes[1].(*Terminal).Value,
		nodes[1].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}

	// without separator
	s = NewScanner([]byte(`"hello" "world" 42`))
	nt, s = Y(s)
	nodes = nt.([]ParsecNode)
	if len(nodes) != 1 {
		t.Errorf("expected %v, got %v", 1, len(nodes))
	}
	ref = []interface{}{`"hello"`, 0}
	out = []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}
}

func TestParsecMany(t *testing.T) {
	Y := Many(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String())

	s := NewScanner([]byte(`"hello" "world" 42`))
	nt, s := Y(s)
	nodes := nt.([]ParsecNode)
	if len(nodes) != 2 {
		t.Errorf("expected %v, got %v", 2, len(nodes))
	}
	ref := []interface{}{`"hello"`, 0, `"world"`, 8}
	out := []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
		nodes[1].(*Terminal).Value,
		nodes[1].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}

	// negative case
	s = NewScanner([]byte(`true "hello" "world" 42`))
	nt, s = Y(s)
	if nt != nil {
		t.Errorf("expected nil")
	}

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		Many(nil, String(), Int(), Int())
	}()
}

func TestParsecManySep(t *testing.T) {
	Y := Many(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String(), Token(",", "COMMA"))

	s := NewScanner([]byte(`"hello", "world" 42`))
	nt, s := Y(s)
	nodes := nt.([]ParsecNode)
	if len(nodes) != 2 {
		t.Errorf("expected %v, got %v", 2, len(nodes))
	}
	ref := []interface{}{`"hello"`, 0, `"world"`, 9}
	out := []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
		nodes[1].(*Terminal).Value,
		nodes[1].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}

	// without separator
	s = NewScanner([]byte(`"hello" "world" 42`))
	nt, s = Y(s)
	nodes = nt.([]ParsecNode)
	if len(nodes) != 1 {
		t.Errorf("expected %v, got %v", 1, len(nodes))
	}
	ref = []interface{}{`"hello"`, 0}
	out = []interface{}{
		nodes[0].(*Terminal).Value,
		nodes[0].(*Terminal).Position,
	}
	if reflect.DeepEqual(ref, out) == false {
		t.Errorf("expected %v, got %v", ref, out)
	}
}

func TestParsecMaybe(t *testing.T) {
	Y := Maybe(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		String())

	s := NewScanner([]byte(`"hello"`))
	nt, _ := Y(s)
	terminal := nt.([]ParsecNode)[0].(*Terminal)
	if terminal.Value != `"hello"` {
		t.Errorf("expected %v, got %v", `"hello"`, terminal.Value)
	} else if terminal.Position != 0 {
		t.Errorf("expected %v, got %v", 0, terminal.Position)
	}

	s = NewScanner([]byte(`42`))
	nt, _ = Y(s)
	if nt != nil {
		t.Errorf("expected nil")
	}
}

func BenchmarkParsecAnd(b *testing.B) {
	Y := And(nil, String(), Int())
	s := NewScanner([]byte(`"hello" 42`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkParsecOrdChoice(b *testing.B) {
	Y := OrdChoice(nil, String(), Int())
	s := NewScanner([]byte(`"hello"`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkParsecKleene(b *testing.B) {
	Y := Kleene(nil, String())
	s := NewScanner([]byte(`"hello" "world" 42`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkParsecMany(b *testing.B) {
	Y := Many(nil, String(), Token(",", "COMMA"))
	s := NewScanner([]byte(`"hello", "world" 42`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}

}

func BenchmarkParsecMaybe(b *testing.B) {
	Y := Maybe(nil, String())
	s := NewScanner([]byte(`"hello"`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}
