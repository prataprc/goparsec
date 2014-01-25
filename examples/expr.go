//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// Examples provide an example parser to parse basic arithmetic expression
// based on the following rule.
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
package examples

import (
    "github.com/prataprc/goparsec"
    "strconv"
)

// Terminal rats
var openparan = parsec.Token(`^\(`, "OPENPARAN")
var closeparan = parsec.Token(`^\)`, "CLOSEPARAN")
var addop = parsec.Token(`^\+`, "ADD")
var subop = parsec.Token(`^-`, "SUB")
var multop = parsec.Token(`^\*`, "MULT")
var divop = parsec.Token(`^/`, "DIV")

// Expr constructs parser-combinator for parsing arithmetic expression on
// integer
func Expr(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
    nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
        if len(ns) == 0 {
            return nil
        }
        return ns[0]
    }
    return parsec.OrdChoice(nodify, sum)(s)
}

func prod(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
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

func sum(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
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

func groupExpr(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
    nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
        if len(ns) == 0 {
            return nil
        }
        return ns[1]
    }
    return parsec.And(nodify, openparan, Expr, closeparan)(s)
}

func value(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
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
