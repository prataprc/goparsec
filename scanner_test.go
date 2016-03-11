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

func TestEndof(t *testing.T) {
	text := []byte(`        `)
	s := NewScanner(text)
	_, s = s.SkipWS()
	if s.Endof() == false {
		t.Fatalf("expected end of text")
	}

	text = []byte(`        text`)
	s = NewScanner(text)
	_, s = s.SkipWS()
	if s.Endof() == true {
		t.Fatalf("did not expect end of text")
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
	cursor := s.GetCursor()
	s.SkipWS()
	s.SetCursor(cursor)
	for i := 0; i < b.N; i++ {
		s.SkipWS()
		s.SetCursor(cursor)
	}
}
