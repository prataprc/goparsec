// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

package lib

import "strconv"

import "github.com/prataprc/goparsec"

// Circular rats
var prod, sum, value, Y parsec.Parser

// Terminal rats
var openparan = parsec.Token(`\(`, "OPENPARAN")
var closeparan = parsec.Token(`\)`, "CLOSEPARAN")
var addop = parsec.Token(`\+`, "ADD")
var subop = parsec.Token(`-`, "SUB")
var multop = parsec.Token(`\*`, "MULT")
var divop = parsec.Token(`/`, "DIV")

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
