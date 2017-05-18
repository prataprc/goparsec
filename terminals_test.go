package parsec

import "testing"
import "fmt"

var _ = fmt.Sprintf("dummy")

func TestTerminalString(t *testing.T) {
	// double quote
	s := NewScanner([]byte(`"hello \"world"`))
	node, s := String()(s)
	tokstr := node.(string)
	if s.Endof() == false {
		t.Errorf("expected end of text")
	} else if ref := `"hello "world"`; tokstr != ref {
		t.Errorf("expected %q, got %q", ref, tokstr)
	}

	// double quote with white spaces around
	s = NewScanner([]byte(` "hello world"   `))
	node, s = String()(s)
	tokstr = node.(string)
	if s.Endof() == true {
		t.Errorf("did not expected end of text")
	} else if ref := `"hello world"`; tokstr != ref {
		t.Errorf("expected %v, got %q", ref, tokstr)
	}

	// negative cases
	s = NewScanner([]byte(` `))
	node, s = String()(s)
	if node != nil {
		t.Errorf("unexpected terminal %q", tokstr)
	}
	s = NewScanner([]byte(`a"`))
	node, _ = String()(s)
	if node != nil {
		t.Errorf("unexpected terminal %q", tokstr)
	}

	// malformed string
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		s = NewScanner([]byte(`"hello`))
		String()(s)
	}()
}

func TestTerminalChar(t *testing.T) {
	s := NewScanner([]byte(`'a'`))
	node, _ := Char()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `'a'` {
		t.Errorf("expected %v, got %v", `a`, terminal.Value)
	}
	// white-space
	s = NewScanner([]byte(`'a`))
	node, _ = Char()(s)
	if node != nil {
		t.Errorf("unexpected terminal node")
	}
	// negative case
	s = NewScanner([]byte(``))
	node, _ = Char()(s)
	if node != nil {
		t.Errorf("unexpected terminal, %v", node)
	}
}

func TestTerminalFloat(t *testing.T) {
	s := NewScanner([]byte(` 10.`))
	node, _ := Float()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `10.` {
		t.Errorf("expected %v, got %v", `10.`, terminal.Value)
	} else if terminal.Name != "FLOAT" {
		t.Errorf("expected %v, got %v", "FLOAT", terminal.Name)
	} else if terminal.Position != 1 {
		t.Errorf("expected %v, got %v", 1, terminal.Position)
	}
	// with decimal
	s = NewScanner([]byte(`+10.10`))
	node, _ = Float()(s)
	terminal = node.(*Terminal)
	if terminal.Value != `+10.10` {
		t.Errorf("expected %v, got %v", `10.0`, terminal.Value)
	}
	// small-decimal
	s = NewScanner([]byte(`-.10`))
	node, _ = Float()(s)
	terminal = node.(*Terminal)
	if terminal.Value != `-.10` {
		t.Errorf("expected %v, got %v", `-.10`, terminal.Value)
	}
	// not float
	s = NewScanner([]byte(`hello`))
	node, _ = Float()(s)
	if node != nil {
		t.Errorf("unexpected float")
	}
	// not float
	s = NewScanner([]byte(`.`))
	node, _ = Float()(s)
	if node != nil {
		t.Errorf("unexpected float")
	}
	// not float
	s = NewScanner([]byte(`-100.0 100.0`))
	nodes := []ParsecNode{}
	node, s = Float()(s)
	for node != nil {
		nodes = append(nodes, node)
		node, s = Float()(s)
	}
	if len(nodes) != 2 {
		t.Errorf("expected 1 node, got %v", nodes)
	}
}

func TestTerminalHex(t *testing.T) {
	s := NewScanner([]byte(`0x10ab`))
	node, _ := Hex()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `0x10ab` {
		t.Errorf("expected %v, got %v", `0x10ab`, terminal.Value)
	} else if terminal.Name != "HEX" {
		t.Errorf("expected %v, got %v", "HEX", terminal.Name)
	}
	// with caps
	s = NewScanner([]byte(`0x10AB`))
	node, _ = Hex()(s)
	terminal = node.(*Terminal)
	if terminal.Value != `0x10AB` {
		t.Errorf("expected %v, got %v", `0x10AB`, terminal.Value)
	}
}

func TestTerminalOct(t *testing.T) {
	s := NewScanner([]byte(`007`))
	node, _ := Oct()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `007` {
		t.Errorf("expected %v, got %v", `007`, terminal.Value)
	} else if terminal.Name != "OCT" {
		t.Errorf("expected %v, got %v", "OCT", terminal.Name)
	}
	// with caps
	s = NewScanner([]byte(`08`))
	node, _ = Oct()(s)
	if node != nil {
		t.Errorf("expected nil, got %v", node)
	}
}

func TestTerminalInt(t *testing.T) {
	s := NewScanner([]byte(`1239`))
	node, _ := Int()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `1239` {
		t.Errorf("expected %v, got %v", `1239`, terminal.Value)
	} else if terminal.Name != "INT" {
		t.Errorf("expected %v, got %v", "INT", terminal.Name)
	}
}

func TestTerminalIdent(t *testing.T) {
	s := NewScanner([]byte(`identifier`))
	node, _ := Ident()(s)
	terminal := node.(*Terminal)
	if terminal.Value != `identifier` {
		t.Errorf("expected %v, got %v", `identifier`, terminal.Value)
	} else if terminal.Name != "IDENT" {
		t.Errorf("expected %v, got %v", "IDENT", terminal.Name)
	}
}

func TestTerminalOrdTokens(t *testing.T) {
	Y := OrdTokens([]string{`\+`, `-`}, []string{"PLUS", "MINUS"})
	s := NewScanner([]byte(` +-`))
	node, s := Y(s)
	terminal := node.(*Terminal)
	if terminal.Value != `+` {
		t.Errorf("expected %v, got %v", `+`, terminal.Value)
	} else if terminal.Name != "PLUS" {
		t.Errorf("expected %v, got %v", "PLUS", terminal.Name)
	} else if terminal.Position != 1 {
		t.Errorf("expected %v, got %v", 1, terminal.Position)
	}

	node, s = Y(s)
	terminal = node.(*Terminal)
	if s.Endof() == false {
		t.Errorf("expected end of scanner")
	} else if terminal.Value != `-` {
		t.Errorf("expected %v, got %v", `-`, terminal.Value)
	} else if terminal.Name != "MINUS" {
		t.Errorf("expected %v, got %v", "MINUS", terminal.Name)
	} else if terminal.Position != 2 {
		t.Errorf("expected %v, got %v", 2, terminal.Position)
	}
}

func TestEnd(t *testing.T) {
	p := And(nil, Token("test", "T"), End())
	s := NewScanner([]byte("test"))
	v, e := p(s)
	if v == nil {
		t.Errorf("End() didn't match %q", e)
	}
}

func TestNotEnd(t *testing.T) {
	p := And(nil, Token("test", "T"), End())
	s := NewScanner([]byte("testing"))
	v, _ := p(s)
	if v != nil {
		t.Errorf("End() shouldn't have matched %q", v)
	}
}

func TestNoEnd(t *testing.T) {
	p := And(nil, Token("test", "T"), NoEnd())
	s := NewScanner([]byte("testing"))
	v, e := p(s)
	if v == nil {
		t.Errorf("NoEnd() didn't match %q", e)
	}
}

func TestNotNoEnd(t *testing.T) {
	p := And(nil, Token("test", "T"), NoEnd())
	s := NewScanner([]byte("test"))
	v, _ := p(s)
	if v != nil {
		t.Errorf("NoEnd() shouldn't have matched %q", v)
	}
}

func BenchmarkTerminalString(b *testing.B) {
	Y := String()
	s := NewScanner([]byte(`  "hello"`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalChar(b *testing.B) {
	Y := Char()
	s := NewScanner([]byte(`  'h'`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalFloat(b *testing.B) {
	Y := Float()
	s := NewScanner([]byte(`  10.10`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalHex(b *testing.B) {
	Y := Hex()
	s := NewScanner([]byte(`  0x1231abcd`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalOct(b *testing.B) {
	Y := Oct()
	s := NewScanner([]byte(`  0x1231abcd`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalInt(b *testing.B) {
	Y := Int()
	s := NewScanner([]byte(`  1231`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalIdent(b *testing.B) {
	Y := Ident()
	s := NewScanner([]byte(`  true`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalToken(b *testing.B) {
	Y := Token("sometoken", "TOKEN")
	s := NewScanner([]byte(`  sometoken`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkTerminalOrdTokens(b *testing.B) {
	Y := OrdTokens([]string{`\+`, `-`}, []string{"PLUS", "MINUS"})
	s := NewScanner([]byte(`  +-`))
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}
