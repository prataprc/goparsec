// Package parsec implements a library of parser-combinators using basic
// recognizers like,
//      And, OrdChoice, Kleene, Many and Maybe.

package parsec

import (
	"fmt"
)

var _ = fmt.Sprintf("keep 'fmt' import for debugging")

type ParsecNode interface{}                      // Can be used to construct AST.
type Parser func(Scanner) (ParsecNode, *Scanner) // combinable parsers
type Nodify func([]ParsecNode) ParsecNode

// ParsecNode of type Terminal.
type Terminal struct {
	Name     string // typically contains terminal's token type
	Value    string // value of the terminal
	Position int    // Offset into the text stream where token was identified
}

// ParsecNode of type NonTerminal.
type NonTerminal struct {
	Name     string // typically contains terminal's token type
	Value    string // value of the terminal
	Children []ParsecNode
}

func And(callb Nodify, parsers ...Parser) Parser {
	return func(s Scanner) (ParsecNode, *Scanner) {
		var ns = make([]ParsecNode, 0, len(parsers))
		var n ParsecNode
		news := &s
		for _, parser := range parsers {
			n, news = parser(*news)
			if n == nil {
				return docallback(callb, []ParsecNode{}), &s
			}
			ns = append(ns, n)
		}
		return docallback(callb, ns), news
	}
}

func OrdChoice(callb Nodify, parsers ...Parser) Parser {
	return func(s Scanner) (ParsecNode, *Scanner) {
		for _, parser := range parsers {
			n, news := parser(s)
			if n != nil {
				return docallback(callb, []ParsecNode{n}), news
			}
		}
		return docallback(callb, []ParsecNode{}), &s
	}
}

func Kleene(callb Nodify, parsers ...Parser) Parser {
	var opScan, sepScan Parser
	opScan = parsers[0]
	if len(parsers) == 2 {
		sepScan = parsers[1]
	} else {
		panic("Kleene parser does not accept more than 2 parsers")
	}
	return func(s Scanner) (ParsecNode, *Scanner) {
		var n ParsecNode
		ns := make([]ParsecNode, 0)
		news := &s
		for {
			n, news = opScan(*news)
			if n == nil {
				break
			}
			ns = append(ns, n)
			if sepScan != nil {
				if n, news = sepScan(*news); n == nil {
					break
				}
			}
		}
		return docallback(callb, ns), news
	}
}

func Many(callb Nodify, parsers ...Parser) Parser {
	var opScan, sepScan Parser
	opScan = parsers[0]
	if len(parsers) >= 2 {
		sepScan = parsers[1]
	}
	return func(s Scanner) (ParsecNode, *Scanner) {
		var n ParsecNode
		ns := make([]ParsecNode, 0)
		news := &s
		for {
			n, news = opScan(*news)
			if n != nil {
				ns = append(ns, n)
				if sepScan != nil {
					if n, news = sepScan(*news); n == nil {
						break
					}
				}
			} else {
				break
			}
		}
		return docallback(callb, ns), news
	}
}

func Maybe(callb Nodify, parser Parser) Parser {
	return func(s Scanner) (ParsecNode, *Scanner) {
		n, news := parser(s)
		if n == nil {
			return docallback(callb, []ParsecNode{}), news
		}
		return docallback(callb, []ParsecNode{n}), news
	}
}

//---- Local function

func docallback(callb Nodify, ns []ParsecNode) ParsecNode {
	if callb != nil {
		return callb(ns)
	} else {
		return ns
	}
}
