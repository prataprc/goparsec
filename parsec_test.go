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
