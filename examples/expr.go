package main

import (
	"fmt"
	"github.com/prataprc/golib/parsec"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Terminal rats
var openparan = parsec.Token(`^\(`, "OPENPARAN")
var closeparan = parsec.Token(`^\)`, "CLOSEPARAN")
var addop = parsec.Token(`^\+`, "ADD")
var subop = parsec.Token(`^-`, "SUB")
var multop = parsec.Token(`^\*`, "MULT")
var divop = parsec.Token(`^/`, "DIV")

// expr  -> sum
// prod  -> (mulop value)*
// mulop -> "*"
//       |  "/"
// sum   -> (addop prod)*
// addop -> "+"
//       |  "-"
// value -> num
//       | ( expr )

// Construct parser-combinator for parsing arithmetic expression on integer
func expr(s parsec.Scanner) (parsec.ParsecNode, *parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns[0]
	}
	return parsec.OrdChoice(nodify, sum)(s)
}

func prod(s parsec.Scanner) (parsec.ParsecNode, *parsec.Scanner) {
	nodifyop := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns[0]
	}
	op := parsec.OrdChoice(nodifyop, multop, divop)

	nodifyk := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns
	}
	opval := parsec.And(nodifyk, op, value)
	k := parsec.Kleene(nil, opval, nil)

	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			val := ns[0].(int)
			for _, x := range ns[1].([]parsec.ParsecNode) {
				y := x.([]parsec.ParsecNode)
				n := y[1].(int)
				switch y[0].(*parsec.Terminal).Name {
				case "MULT":
					val *= n
				case "DIV":
					val /= n
				}
			}
			return val
		}
		return nil
	}
	return parsec.And(nodify, value, k)(s)
}

func sum(s parsec.Scanner) (parsec.ParsecNode, *parsec.Scanner) {
	nodifyop := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns[0]
	}
	op := parsec.OrdChoice(nodifyop, addop, subop)

	nodifyk := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns
	}
	opval := parsec.And(nodifyk, op, prod)
	k := parsec.Kleene(nil, opval, nil)

	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			val := ns[0].(int)
			for _, x := range ns[1].([]parsec.ParsecNode) {
				y := x.([]parsec.ParsecNode)
				n := y[1].(int)
				switch y[0].(*parsec.Terminal).Name {
				case "ADD":
					val += n
				case "SUB":
					val -= n
				}
			}
			return val
		}
		return nil
	}
	return parsec.And(nodify, prod, k)(s)
}

func groupExpr(s parsec.Scanner) (parsec.ParsecNode, *parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		}
		return ns[1]
	}
	return parsec.And(nodify, openparan, expr, closeparan)(s)
}

func value(s parsec.Scanner) (parsec.ParsecNode, *parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) == 0 {
			return nil
		} else if term, ok := ns[0].(*parsec.Terminal); ok {
			val, _ := strconv.Atoi(term.Value)
			return val
		}
		return ns[0]
	}
	return parsec.OrdChoice(nodify, parsec.Int(), groupExpr)(s)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run expr.go <expression-file>\n")
		os.Exit(1)
	}
	text, _ := ioutil.ReadFile(os.Args[1])
	s := parsec.NewScanner(text)
	count := int64(10000)
	t1 := time.Now().UnixNano()
	for i := int64(0); i < count; i++ {
		expr(*s)
	}
	t2 := time.Now().UnixNano()
	fmt.Printf("Takes %vnS to evaluate \n", (t2-t1)/count)
}
