package parsec

import "reflect"
import "testing"

func TestNonTerminal(t *testing.T) {
	nt := &NonTerminal{Name: "NONTERM", Children: make([]Queryable, 0)}
	if nt.GetName() != "NONTERM" {
		t.Errorf("expected %q, got %q", "NONTERM", nt.GetName())
	} else if nt.IsTerminal() == true {
		t.Errorf("expected false")
	} else if nt.GetValue() != "" {
		t.Errorf("expected %q, got %q", "", nt.GetValue())
	} else if cs := nt.GetChildren(); len(cs) != 0 {
		t.Errorf("expected len %v, got %v", 0, len(cs))
	} else if nt.GetPosition() != 0 {
		t.Errorf("expected %v, got %v", 0, nt.GetPosition())
	}
	mn := MaybeNone("missing")
	nt.Children = append(nt.Children, mn)
	if cs := nt.GetChildren(); cs[0] != mn {
		t.Errorf("expected %v, got %v", mn, cs[0])
	} else if nt.GetPosition() != -1 {
		t.Errorf("expected %v, got %v", -1, nt.GetPosition())
	}
	// check attribute methods.
	nt.SetAttribute("name", "one").SetAttribute("name", "two")
	nt.SetAttribute("key", "one")
	ref1 := []string{"one", "two"}
	ref2 := map[string][]string{
		"name": {"one", "two"},
		"key":  {"one"},
	}
	if x := nt.GetAttribute("name"); reflect.DeepEqual(x, ref1) == false {
		t.Errorf("expected %v, got %v", ref1, x)
	} else if x := nt.GetAttributes(); reflect.DeepEqual(x, ref2) == false {
		t.Errorf("expected %v, got %v", ref2, x)
	}
}
