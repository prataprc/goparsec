package parsec

import "testing"
import "reflect"

func TestStrEOF(t *testing.T) {
	word := String()
	Y := Many(
		func(ns []ParsecNode) ParsecNode {
			return ns
		},
		word)

	input := `"alpha" "beta" "gamma"`
	s := NewScanner([]byte(input))

	root, _ := Y(s)
	nodes := root.([]ParsecNode)
	ref := []ParsecNode{"\"alpha\"", "\"beta\"", "\"gamma\""}
	if !reflect.DeepEqual(nodes, ref) {
		t.Fatal(nodes)
	}
}
