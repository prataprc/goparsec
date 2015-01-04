//  Copyright (c) 2013 Couchbase, Inc.

package parsec

import "reflect"
import "testing"

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
