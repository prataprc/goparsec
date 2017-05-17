package parsec

import "testing"

func TestStrEOF(t *testing.T) {
	word := And(allTokens, String(), TokenWS("", "SPACE"))
	Y := Many(
		allTokens,
		word)

	input := `"alpha" "beta" "gamma"`
	s := NewScanner([]byte(input))

	root, _ := Y(s)
	ref := []string{"\"alpha\"", "\"beta\"", "\"gamma\""}

	_, n := IsTerminal(root)
	correct := true

	if n == 3 {
		nodes := root.([]ParsecNode)

		for i, v := range nodes {
			if _, nn := IsTerminal(v); nn == 2 {
				str, ok := v.([]ParsecNode)[0].(*Terminal)
				if !ok || str.Value != ref[i] {
					correct = false
					t.Log("incorrect string parse")
				}
			} else {
				correct = false
				t.Log("each And subnode should have two terminals", nn)
			}
		}
	} else {
		correct = false
		t.Log("the Many node should have three children")
	}

	if !correct {
		t.Fatal("did not parse correctly: ", root)
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
