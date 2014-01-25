// Examples provide an example parser to parse JSON string.
package examples

import (
	"github.com/prataprc/goparsec"
	"io/ioutil"
	"strconv"
)

// Property node defines a parsec.ParsecNode for JSON property map type.
type PropertyNode struct {
	propname string
	parsec.ParsecNode
}

// EMPTY is terminal parsec.ParsecNode
var EMPTY = parsec.Terminal{Name: "EMPTY", Value: ""}

// ParseFile accepts `filename` that contains the json document, parses the
// document and returns the root node of the AST tree.
func ParseFile(filename string) parsec.ParsecNode {
	if text, err := ioutil.ReadFile(filename); err != nil {
		panic(err.Error())
	} else {
		return Parse(text)
	}
}

// Parse accepts json document as byte slice, parses the document and returns
// the root node of the AST tree.
func Parse(text []byte) parsec.ParsecNode {
	s := parsec.NewScanner(text)
	nt, _ := y(s)
	return nt
}

// Value takes the root node of parsed JSON document and returns an
// interface{} of golang types including array and map.
func Value(n parsec.ParsecNode) interface{} {
	conv := func(fn func() (interface{}, error)) interface{} {
		v, err := fn()
		if err != nil {
			panic(err)
		}
		return v
	}
	if t, ok := n.(*parsec.Terminal); ok {
		switch t.Name {
		case "INT":
			return conv(func() (interface{}, error) {
				return strconv.Atoi(t.Value)
			})
		case "FLOAT":
			return conv(func() (interface{}, error) {
				return strconv.ParseFloat(t.Value, 64)
			})
		case "STRING":
			return t.Value[1 : len(t.Value)-1]
		case "TRUE":
			return true
		case "FALSE":
			return false
		case "NULL":
			return nil
		}
	}
	if nt, ok := n.(*parsec.NonTerminal); ok {
		switch nt.Name {
		case "VALUES":
			vs := make([]interface{}, 0)
			for _, v := range nt.Children {
				vs = append(vs, Value(v))
			}
			return vs
		case "PROPERTIES":
			m := make(map[string]interface{})
			for _, v := range nt.Children {
				if u, ok := v.(*PropertyNode); !ok {
					panic("Expected PropertyNode")
				} else {
					name := u.propname[1 : len(u.propname)-1]
					m[name] = Value(u.ParsecNode)
				}
			}
			return m
		}
	}
	return nil
}

// Construct parser-combinator for parsing JSON string.
func y(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if ns == nil || len(ns) == 0 {
			return nil
		}
		return ns[0]
	}
	return parsec.Maybe(nodify, jsonvalue)(s)
}

func array(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if ns == nil || len(ns) == 0 {
			return nil
		}
		return ns[1]
	}
	return parsec.And(nodify, opensqr, values, closesqr)(s)
}

func object(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if ns == nil || len(ns) == 0 {
			return nil
		}
		return ns[1]
	}
	return parsec.And(nodify, openbrace, properties, closebrace)(s)
}

func properties(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		// Bubble sort properties based on property name.
		for i := 0; i < len(ns)-1; i++ {
			for j := 0; j < len(ns)-i-1; j++ {
				x := ns[j].(*PropertyNode).propname
				y := ns[j+1].(*PropertyNode).propname
				if x <= y {
					continue
				}
				ns[j+1], ns[j] = ns[j], ns[j+1]
			}
		}
		return &parsec.NonTerminal{Name: "PROPERTIES", Children: ns}
	}
	return parsec.Many(nodify, property, comma)(s)
}

func property(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if ns == nil || len(ns) == 0 {
			return nil
		}
		propname := ns[0].(*parsec.Terminal).Value
		return &PropertyNode{propname, ns[2]}
	}
	return parsec.And(nodify, parsec.String(), colon, jsonvalue)(s)
}

func values(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		return &parsec.NonTerminal{Name: "VALUES", Children: ns}
	}
	return parsec.Many(nodify, jsonvalue, comma)(s)
}

func jsonvalue(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if ns == nil || len(ns) == 0 {
			return nil
		}
		return ns[0]
	}
	return parsec.OrdChoice(
		nodify, tRue, fAlse, null,
		parsec.String(), parsec.Float(), parsec.Int(),
		array, object,
	)(s)
}

var tRue = parsec.Token(`^true`, "TRUE")
var fAlse = parsec.Token(`^false`, "FALSE")
var null = parsec.Token(`^null`, "NULL")

var comma = parsec.Token(`^,`, "COMMA")
var colon = parsec.Token(`^:`, "COLON")
var opensqr = parsec.Token(`^\[`, "OPENSQR")
var closesqr = parsec.Token(`^\]`, "CLOSESQR")
var openbrace = parsec.Token(`^\{`, "OPENPARAN")
var closebrace = parsec.Token(`^\}`, "CLOSEPARAN")

// INode APIs for Terminal
//func Repr(tok parsec.ParsecNode, prefix string) string {
//    if term, ok := tok.(*parsec.Terminal); ok {
//        return fmt.Sprintf(prefix) +
//            fmt.Sprintf("%v : %v ", term.Name, term.Value)
//    } else if propterm, ok := tok.(*PropertyNode); ok {
//        return fmt.Sprintf(prefix) +
//            fmt.Sprintf("property : %v \n", propterm.propname)
//    } else {
//        nonterm, _ := tok.(*parsec.NonTerminal)
//        return fmt.Sprintf(prefix) +
//            fmt.Sprintf("%v : %v \n", nonterm.Name, nonterm.Value)
//    }
//    panic("invalid parsecNode")
//}
//
//func Show(tok parsec.ParsecNode, prefix string) {
//    if term, ok := tok.(*parsec.Terminal); ok {
//        fmt.Println(Repr(term, prefix))
//    } else if propterm, ok := tok.(*PropertyNode); ok {
//        fmt.Printf("%v", Repr(propterm, prefix))
//        Show(propterm.ParsecNode, prefix+"  ")
//    } else if nonterm, ok := tok.(*parsec.NonTerminal); ok {
//        fmt.Printf("%v", Repr(nonterm, prefix))
//        for _, tok := range nonterm.Children {
//            Show(tok, prefix+"  ")
//        }
//    }
//}
