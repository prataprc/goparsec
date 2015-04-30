// Copyright (c) 2013 Couchbase, Inc.

// Package provide a parser to parse basic arithmetic expression based on the
// following rule.
//
//     expr  -> sum
//     prod  -> value (mulop value)*
//     mulop -> "*"
//           |  "/"
//     sum   -> prod (addop prod)*
//     addop -> "+"
//           |  "-"
//     value -> num
//           | "(" expr ")"

package expr

import "strconv"
import "fmt"

import "github.com/prataprc/goparsec"

var _ = fmt.Sprintf("dummp print")

// Y is root Parser, usually called as `s` in CFG theory.
var Y parsec.Parser
var prod, sum, value parsec.Parser // circular rats

// Terminal rats
var openparan = parsec.TokenWS(`\(`, "OPENPARAN")
var closeparan = parsec.TokenWS(`\)`, "CLOSEPARAN")
var addop = parsec.TokenWS(`\+`, "ADD")
var subop = parsec.TokenWS(`-`, "SUB")
var multop = parsec.TokenWS(`\*`, "MULT")
var divop = parsec.TokenWS(`/`, "DIV")

// NonTerminal rats
// addop -> "+" |  "-"
var sumOp = parsec.OrdChoice(one2one, addop, subop)

// mulop -> "*" |  "/"
var prodOp = parsec.OrdChoice(one2one, multop, divop)

// value -> "(" expr ")"
var groupExpr = parsec.And(exprNode, openparan, &sum, closeparan)

// (addop prod)*
var prodK = parsec.Kleene(nil, parsec.And(many2many, sumOp, &prod), nil)

// (mulop value)*
var valueK = parsec.Kleene(nil, parsec.And(many2many, prodOp, &value), nil)

func init() {
	// Circular rats come to life
	// sum -> prod (addop prod)*
	sum = parsec.And(sumNode, &prod, prodK)
	// prod-> value (mulop value)*
	prod = parsec.And(prodNode, &value, valueK)
	// value -> num | "(" expr ")"
	value = parsec.OrdChoice(exprValueNode, parsec.Int(), groupExpr)
	// expr  -> sum
	Y = parsec.OrdChoice(one2one, sum)
}

//----------
// Nodifiers
//----------

func sumNode(ns []parsec.ParsecNode) parsec.ParsecNode {
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

func prodNode(ns []parsec.ParsecNode) parsec.ParsecNode {
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

func exprNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if len(ns) == 0 {
		return nil
	}
	return ns[1]
}

func exprValueNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if len(ns) == 0 {
		return nil
	} else if term, ok := ns[0].(*parsec.Terminal); ok {
		val, _ := strconv.Atoi(term.Value)
		return val
	}
	return ns[0]
}
