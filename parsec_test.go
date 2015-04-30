package parsec

import "testing"
import "reflect"

func TestStrEOF(t *testing.T) {
	word := And(allTokens, String(), Token(" ", "SPACE"))
	Y := Many(
		allTokens,
		word)

	input := `"alpha" "beta" "gamma"`
	s := NewScanner([]byte(input))

	root, _ := Y(s)
	nodes := root.([]ParsecNode)
	ref := []ParsecNode{"\"alpha\"", "\"beta\"", "\"gamma\""}
	if !reflect.DeepEqual(nodes, ref) {
		t.Fatal("did not parse correctly: ", nodes)
	}
}

func allTokens(ns []ParsecNode) ParsecNode {
	return ns
}
