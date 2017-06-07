//  Copyright (c) 2013 Couchbase, Inc.

package parsec

import "reflect"
import "testing"
import "fmt"

var _ = fmt.Sprintf("dummy print")

func TestClone(t *testing.T) {
	text := []byte(`example text`)
	s := NewScanner(text)
	if !reflect.DeepEqual(s, s.Clone()) {
		t.Fatal("Clone() method does not work as intended")
	}
}

func TestMatch(t *testing.T) {
	text := []byte(`example text`)
	ref := `exampl`
	s := NewScanner(text)
	m, s := s.Match(`^ex.*l`)
	if string(m) != ref {
		t.Fatalf("mismatch expected %s, got %s", ref, string(m))
	}
	expcur := 6
	if s.GetCursor() != expcur {
		t.Fatalf("expected cursor position %v, got %v", expcur, s.GetCursor())
	}
}

func TestSubmatchAll(t *testing.T) {
	text := []byte(`alphabetaexample text`)
	s := NewScanner(text)
	pattern := `^(?P<X>alpha)|(?P<Y>beta)(?P<Z>example) text`
	m, s := s.SubmatchAll(pattern)
	if len(m) != 1 {
		t.Fatalf("match failed in len %v\n", m)
	} else if str := string(m["X"]); str != "alpha" {
		t.Fatalf("expected %q got %q\n", "alpha", str)
	}
	m, s = s.SubmatchAll(pattern)
	if len(m) != 2 {
		t.Fatalf("match failed in len %v\n", m)
	} else if str := string(m["Y"]); str != "beta" {
		t.Fatalf("expected %q got %q\n", "beta", str)
	} else if str := string(m["Z"]); str != "example" {
		t.Fatalf("expected %q got %q\n", "example", str)
	}
}

func TestMatchString(t *testing.T) {
	text := `myString0`

	s := NewScanner([]byte(text))
	ok1, s := s.MatchString("my")
	ok2, s := s.MatchString("String0")

	if !ok1 || !ok2 {
		t.Fatalf("did not match correctly")
	}

	if !s.Endof() {
		t.Fatalf("expect end of text")
	}

	//Not matching case
	text2 := `myString`
	s = NewScanner([]byte(text2))
	ok3, s := s.MatchString("myString0")

	if ok3 {
		t.Fatalf("shouldn't have matched")
	}

	if s.Endof() {
		t.Fatalf("did not expect end of text")
	}
}

func TestSkipWS(t *testing.T) {
	text := []byte(`        `)
	ref := `        `
	s := NewScanner(text)
	m, s := s.SkipWS()
	if string(m) != ref {
		t.Fatalf("mismatch expected %q, got %q", ref, string(m))
	}
	expcur := 8
	if s.GetCursor() != expcur {
		t.Fatalf("expected cursor position %v, got %v", expcur, s.GetCursor())
	}
}

func TestSkipAny(t *testing.T) {
	text := `B  
			B
			   BA`
	s := NewScanner([]byte(text))
	_, s = s.SkipAny(`^[ \n\tB]+`)

	aRef := []byte("A")
	a, snew := s.Match("A")

	if a[0] != aRef[0] {
		t.Fatalf("expected character A after skipping whitespaces and B")
	}

	if !snew.Endof() {
		t.Fatalf("input text should have been completely skipped or matched")
	}
}

func TestEndof(t *testing.T) {
	text := []byte(`    text`)
	s := NewScanner(text)
	_, s = s.SkipAny(`^[ ]+`)
	if s.Endof() {
		t.Fatalf("did not expect end of text")
	}

	_, s = s.SkipAny(`^[tex]+`)

	if !s.Endof() {
		t.Fatalf("expect end of text")
	}
}

func TestSetWSPattern(t *testing.T) {
	text := []byte(`// comment`)
	ref := `// comment`
	s := NewScanner(text)
	s.(*SimpleScanner).SetWSPattern(`^//.*`)
	m, s := s.SkipWS()
	if string(m) != ref {
		t.Fatalf("mismatch expected %q, got %q", ref, string(m))
	}
	expcur := 10
	if s.GetCursor() != expcur {
		t.Fatalf("expected cursor position %v, got %v", expcur, s.GetCursor())
	}
}

func BenchmarkSScanClone(b *testing.B) {
	text := []byte("hello world")
	s := NewScanner(text)
	for i := 0; i < b.N; i++ {
		s.Clone()
	}
}

func BenchmarkSScanSkipWS(b *testing.B) {
	text := []byte("    hello world")
	s := NewScanner(text)
	for i := 0; i < b.N; i++ {
		s.SkipWS()
		s.(*SimpleScanner).resetcursor()
	}
}

func BenchmarkSScanSkipAny(b *testing.B) {
	text := []byte("    hello world")
	s := NewScanner(text)
	for i := 0; i < b.N; i++ {
		s.SkipAny(`^[ hel]+`)
		s.(*SimpleScanner).resetcursor()
	}
}

func BenchmarkMatch(b *testing.B) {
	s := NewScanner([]byte(`hello world`))
	for i := 0; i < b.N; i++ {
		s.(*SimpleScanner).resetcursor()
		s.Match(`hello world`)
	}
}

func BenchmarkMatchString(b *testing.B) {
	s := NewScanner([]byte(`hello world`))
	for i := 0; i < b.N; i++ {
		s.(*SimpleScanner).resetcursor()
		s.MatchString(`hello world`)
	}
}
