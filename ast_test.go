package parsec

import "fmt"
import "bytes"
import "testing"
import "io/ioutil"

var _ = fmt.Sprintf("dummy")

func TestASTReset(t *testing.T) {
	ast := NewAST("testand", 100)
	y := ast.And("and", nil, Atom("hello", "TERM"))
	s := NewScanner([]byte("hello"))
	node, s := ast.Parsewith(y, s)
	if x, y := ast.root.GetValue(), node.GetValue(); x != y {
		t.Errorf("expected %v, got %v", x, y)
	}
	if len(ast.ntpool) != 0 {
		t.Errorf("expected 0")
	} else if cap(ast.ntpool) != 100 {
		t.Errorf("expected 100")
	}
	ast = ast.Reset()
	if len(ast.ntpool) != 1 {
		t.Errorf("expected 1")
	} else if cap(ast.ntpool) != 100 {
		t.Errorf("expected 100")
	}
}

func TestASTAnd(t *testing.T) {
	ast := NewAST("testand", 100)
	y := ast.And("and", nil, Atom("hello", "TERM"))
	s := NewScanner([]byte("hello"))
	node, s := ast.Parsewith(y, s)
	if node.GetValue() != "hello" {
		t.Errorf("expected %v, got %v", "hello", node.GetValue())
	} else if s.Endof() == false {
		t.Errorf("expected true")
	}
	ast.Reset()
	// negative case
	s = NewScanner([]byte("he"))
	node, ss := ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("expected nil")
	}
	x := ss.(*SimpleScanner)
	if rem := string(x.buf[x.cursor:]); rem != "he" {
		t.Errorf("expected %v, got %v", "he", rem)
	}
	ast.Reset()
	// nil callback
	y = ast.And("and",
		func(_ string, _ Queryable) Queryable { return nil },
		Atom("hello", "TERM"),
	)
	s = NewScanner([]byte("hello"))
	node, ss = ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("expected nil")
	}
	x = ss.(*SimpleScanner)
	if rem := string(x.buf[x.cursor:]); rem != "hello" {
		t.Errorf("expected %v, got %v", "he", rem)
	}
	ast.Reset()
	// return new object
	y = ast.And("and",
		func(_ string, _ Queryable) Queryable { return MaybeNone("missing") },
		Atom("hello", "TERM"),
	)
	s = NewScanner([]byte("hello"))
	node, ss = ast.Parsewith(y, s)
	if node.(MaybeNone) != "missing" {
		t.Errorf("expected missing")
	} else if ss.Endof() == false {
		t.Errorf("expected Endof")
	}
	if len(ast.ntpool) != 1 {
		t.Errorf("expected 1")
	}
	ast.Reset()
}

func TestASTOrdChoice(t *testing.T) {
	ast := NewAST("testor", 100)
	y := ast.OrdChoice("or", nil, Atom("hello", "TERM"), Atom("world", "ATOM"))
	s := NewScanner([]byte("world"))
	node, s := ast.Parsewith(y, s)
	if node.GetValue() != "world" {
		t.Errorf("expected %v, got %v", "world", node.GetValue())
	} else if s.Endof() == false {
		t.Errorf("expected true")
	}
	ast.Reset()
	// negative case
	s = NewScanner([]byte("he"))
	node, ss := ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("expected nil")
	}
	x := ss.(*SimpleScanner)
	if rem := string(x.buf[x.cursor:]); rem != "he" {
		t.Errorf("expected %v, got %v", "he", rem)
	}
	ast.Reset()
	// nil callback
	y = ast.OrdChoice("or",
		func(_ string, _ Queryable) Queryable { return nil },
		Atom("hello", "TERM"), Atom("world", "TERM"),
	)
	s = NewScanner([]byte("world"))
	node, ss = ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("expected nil")
	}
	x = ss.(*SimpleScanner)
	if rem := string(x.buf[x.cursor:]); rem != "world" {
		t.Errorf("expected %v, got %v", "he", rem)
	}
	ast.Reset()
}

func TestASTKleene(t *testing.T) {
	// without separator
	ast := NewAST("testkleene", 100)
	y := ast.Kleene("kleene", nil, Atom("hello", "ONE"))
	s := NewScanner([]byte("hellohello"))
	q, ss := ast.Parsewith(y, s)
	if len(q.GetChildren()) != 2 {
		t.Errorf("unexpected %v", len(q.GetChildren()))
	} else if q.GetValue() != "hellohello" {
		t.Errorf("unexpected %v", q.GetValue())
	} else if ss.Endof() == false {
		t.Errorf("expected true")
	}
	ast.Reset()
	// empty case
	s = NewScanner([]byte("world"))
	q, ss = ast.Parsewith(y, s)
	if len(q.GetChildren()) != 0 {
		t.Errorf("unexpected %v", len(q.GetChildren()))
	} else if q.GetValue() != "" {
		t.Errorf("unexpected %v", q.GetValue())
	}
	x := ss.(*SimpleScanner)
	if str := string(x.buf[x.cursor:]); str != "world" {
		t.Errorf("unexpected %v", str)
	} else if ss.Endof() == true {
		t.Errorf("expected false")
	}
	// with separator
	y = ast.Kleene("kleene", nil, Atom("hello", "ONE"), Atom(",", "COMMA"))
	s = NewScanner([]byte("hello,hello"))
	q, ss = ast.Parsewith(y, s)
	if len(q.GetChildren()) != 2 {
		t.Errorf("unexpected %v", len(q.GetChildren()))
	} else if q.GetValue() != "hellohello" {
		t.Errorf("unexpected %v", q.GetValue())
	} else if ss.Endof() == false {
		t.Errorf("expected true")
	}
	ast.Reset()
	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		Kleene(nil, Atom("hello", "W"), Atom(",", "COMMA"), Atom(",", "COMMA"))
	}()
}

func TestASTStrEOF(t *testing.T) {
	word := Token(`"[a-z]+"`, "W")

	ast := NewAST("teststreof", 100)
	y := ast.Many(
		"many",
		func(_ string, ns Queryable) Queryable { return ns },
		word,
	)
	input := `"alpha" "beta" "gamma"`
	ref := `"alpha""beta""gamma"`
	s := NewScanner([]byte(input))
	node, _ := ast.Parsewith(y, s)
	if node.GetValue() != ref {
		t.Errorf("expected %v, got %v", ref, node.GetValue())
	}
}

func TestASTMany(t *testing.T) {
	w := Token("\\w+", "W")

	// without separator
	ast := NewAST("testmany", 100)
	y := ast.Many("many", nil, w)
	s, ref := NewScanner([]byte("one two stop")), "onetwostop"
	node, ss := ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("Many() didn't match %q", ss)
	} else if node.GetValue() != ref {
		t.Errorf("Many() unexpected: %v", node.GetValue())
	}
	ast.Reset()
	// with separator
	y = ast.Many("many", nil, w, Atom(",", "COMMA"))
	s = NewScanner([]byte("one,two"))
	node, ss = ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("Many() didn't match %q", ss)
	} else if node.GetValue() != "onetwo" {
		t.Errorf("Many() unexpected : %q", node.GetValue())
	}
	// Return nil
	y = ast.Many(
		"many",
		func(_ string, _ Queryable) Queryable { return nil },
		w, Atom(",", "COMMA"),
	)
	s = NewScanner([]byte("one,two"))
	node, ss = ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("Many() didn't match %q", ss)
	}
	x := ss.(*SimpleScanner)
	if str := string(x.buf[x.cursor:]); str != "one,two" {
		t.Errorf("expected %q, got %q", "one,two", str)
	}
	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		ast.Many("many", nil, w, Atom(",", "COMMA"), Atom(",", "COMMA"))
	}()
}

func TestASTManyUntil(t *testing.T) {
	w := Token("\\w+", "W")

	// Return nil
	ast := NewAST("testmanyuntil", 100)
	y := ast.ManyUntil(
		"manyuntil",
		func(_ string, _ Queryable) Queryable { return nil },
		w, Atom(",", "COMMA"),
	)
	s := NewScanner([]byte("one,two"))
	node, ss := ast.Parsewith(y, s)
	if node != nil {
		t.Errorf("ManyUntil() expected nil")
	}
	x := ss.(*SimpleScanner)
	if str := string(x.buf[x.cursor:]); str != "one,two" {
		t.Errorf("expected %q, got %q", "one,two", str)
	}

	// panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		ast.ManyUntil("manyuntil", nil, w)
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		ast.ManyUntil(
			"manyuntil", nil,
			w, Atom(",", "COMMA"), Atom(",", "COMMA"), Atom(",", "COMMA"),
		)
	}()
}

func TestASTManyUntilNoStop(t *testing.T) {
	w, u := Token("\\w+", "W"), Token("stop", "S")

	ast := NewAST("nostop", 100)
	y := ast.ManyUntil("manyuntil", nil, w, u)
	s := NewScanner([]byte("one two three"))
	node, ss := ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("ManyUntil() didn't match : %q", ss)
	} else if node.GetValue() != "onetwothree" {
		t.Errorf("ManyUntil() unexpected : %q", node.GetValue())
	}
}

func TestASTManyUntilStop(t *testing.T) {
	w, u := Token("\\w+", "W"), Token("stop", "S")

	ast := NewAST("stop", 100)
	y := ast.ManyUntil("manyuntil", nil, w, u)
	s := NewScanner([]byte("one two stop"))
	node, ss := ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("ManyUntil() didn't match %q", ss)
	} else if node.GetValue() != "onetwo" {
		t.Errorf("ManyUntil() unexpected : %q", node.GetValue())
	}
}

func TestASTManyUntilNoStopSep(t *testing.T) {
	w, u, z := Token("\\w+", "W"), Token("stop", "S"), Token("z", "Z")

	ast := NewAST("nostopsep", 100)
	y := ast.ManyUntil("manyuntil", nil, w, z, u)
	s := NewScanner([]byte("one z two z three"))
	node, ss := ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("ManyUntil() didn't match %q", ss)
	} else if node.GetValue() != "onetwothree" {
		t.Errorf("ManyUntil() unexpected : %q", node.GetValue())
	}
}

func TestASTManyUntilStopSep(t *testing.T) {
	w, u, z := Token("\\w+", "W"), Token("stop", "S"), Token("z", "Z")

	ast := NewAST("stopsep", 100)
	y := ast.ManyUntil("manyuntil", nil, w, z, u)
	s := NewScanner([]byte("one z two z stop"))
	node, ss := ast.Parsewith(y, s)
	if node == nil {
		t.Errorf("ManyUntil() didn't match %q", ss)
	} else if node.GetValue() != "onetwo" {
		t.Errorf("ManyUntil() didn't stop %q", node.GetValue())
	}
}

func TestASTForwardReference(t *testing.T) {
	var ycomma Parser
	w := Token(`[a-z]+`, "W")

	ast := NewAST("testforward", 100)
	y := ast.Kleene("kleene", nil, ast.Maybe("maybe", nil, w), &ycomma)
	ycomma = Atom(",", "COMMA")
	s := NewScanner([]byte("one,two,,three"))
	node, _ := ast.Parsewith(y, s)
	if node.GetValue() != "onetwothree" {
		t.Errorf("unexpected: %v", node.GetValue())
	} else if _, ok := node.GetChildren()[2].(MaybeNone); ok == false {
		t.Errorf("expected MissingNone")
	}
	ast.Reset()

	// nil return
	y = ast.Maybe(
		"maybe",
		func(_ string, _ Queryable) Queryable { return nil }, w,
	)
	s = NewScanner([]byte("one"))
	node, _ = ast.Parsewith(y, s)
	if node != MaybeNone("missing") {
		t.Errorf("expected %v, got %v", node, MaybeNone("missing"))
	}
}

func TestGetValue(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/simple.html")
	if err != nil {
		t.Error(err)
	}
	data = bytes.Trim(data, " \t\r\n")
	ast := NewAST("html", 100)
	y := makehtmly(ast)
	s := NewScanner(data).TrackLineno()
	node, _ := ast.Parsewith(y, s)
	if node.GetValue() != string(data) {
		t.Errorf("expected %q", string(data))
		t.Errorf("got %q", node.GetValue())
	}
}

func makehtmly(ast *AST) Parser {
	var tag Parser

	opentag := AtomExact("<", "OT")
	closetag := AtomExact(">", "CT")
	equal := AtomExact("=", "EQUAL")
	slash := TokenExact("/[ \t]*", "SLASH")
	tagname := TokenExact("[a-z][a-zA-Z0-9]*", "TAG")
	attrkey := TokenExact("[a-z][a-zA-Z0-9]*", "ATTRK")
	text := TokenExact("[^<>]+", "TEXT")
	ws := TokenExact("[ \t]+", "TEXT")

	element := ast.OrdChoice("element", nil, text, &tag)
	elements := ast.Kleene("elements", nil, element)
	attr := ast.And("attribute", nil, attrkey, equal, String())
	attrws := ast.And("attrws", nil, attr, ast.Maybe("ws", nil, ws))
	attrs := ast.Kleene("attributes", nil, attrws)
	tstart := ast.And("tagstart", nil, opentag, tagname, attrs, closetag)
	tend := ast.And("tagend", nil, opentag, slash, tagname, closetag)
	tag = ast.And("tag", nil, tstart, elements, tend)
	return tag
}
