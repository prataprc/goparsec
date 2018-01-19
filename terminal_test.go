package parsec

import "reflect"
import "testing"

func TestTerminal(t *testing.T) {
	term := &Terminal{Name: "TERM", Value: "xyz", Position: 2}
	if term.GetName() != "TERM" {
		t.Errorf("expected %q, got %q", "TERM", term.GetName())
	} else if term.IsTerminal() == false {
		t.Errorf("expected true")
	} else if term.GetValue() != "xyz" {
		t.Errorf("expected %q, got %q", "xyz", term.GetValue())
	} else if term.GetChildren() != nil {
		t.Errorf("expected nil")
	} else if term.GetPosition() != 2 {
		t.Errorf("expected %v, got %v", 2, term.GetPosition())
	}
	// check attribute methods.
	term.SetAttribute("name", "one").SetAttribute("name", "two")
	term.SetAttribute("key", "one")
	ref1 := []string{"one", "two"}
	ref2 := map[string][]string{
		"name": {"one", "two"},
		"key":  {"one"},
	}
	if x := term.GetAttribute("name"); reflect.DeepEqual(x, ref1) == false {
		t.Errorf("expected %v, got %v", ref1, x)
	} else if x := term.GetAttributes(); reflect.DeepEqual(x, ref2) == false {
		t.Errorf("expected %v, got %v", ref2, x)
	}
}

func TestMaybeNone(t *testing.T) {
	mn := MaybeNone("missing")
	if string(mn) != mn.GetName() {
		t.Errorf("expected %q, got %q", string(mn), mn.GetName())
	} else if mn.IsTerminal() == false {
		t.Errorf("expected true")
	} else if mn.GetValue() != "" {
		t.Errorf("expected %q, got %q", "", mn.GetValue())
	} else if mn.GetChildren() != nil {
		t.Errorf("expected nil")
	} else if mn.GetPosition() != -1 {
		t.Errorf("expected %v, got %v", -1, mn.GetPosition())
	}
	// check attribute methods.
	mn.SetAttribute("name", "one").SetAttribute("name", "two")
	mn.SetAttribute("key", "one")
	if x := mn.GetAttribute("name"); x != nil {
		t.Errorf("unexpected %v", x)
	} else if x := mn.GetAttributes(); x != nil {
		t.Errorf("unexpected %v", x)
	}
}
