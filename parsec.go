// Copyright (c) 2013 Goparsec AUTHORS. All rights reserved.
// Use of this source code is governed by LICENSE file.

package parsec

import "fmt"

// ParsecNode for parsers return input text as parsed nodes.
type ParsecNode interface{}

// Parser function parses input text encapsulated by Scanner, higher
// order parsers are constructed using combinators.
type Parser func(Scanner) (ParsecNode, Scanner)

// Nodify callback function to construct custom ParsecNode. Even when
// combinators like And, OrdChoice, Many etc.. can match input string,
// it is still possible to fail them via nodify callback function, by
// returning nil. This very useful in cases when,
//  * lookahead matching is required.
//  * an exceptional cases for regex pattern.
//
// Note that some combinators like KLEENE shall not interpret the return
// value from Nodify callback.
type Nodify func([]ParsecNode) ParsecNode

// And combinator accepts a list of `Parser`, or reference to a
// parser, that must match the input string, atleast until the
// last Parser argument. Return a parser function that can further be
// used to construct higher-level parsers.
//
// If all parser matches, a list of ParsecNode, where each
// ParsecNode is constructed by matching parser, will be passed
// as argument to Nodify callback. Even if one of the input
// parser function fails, And will fail without consuming the input.
func And(callb Nodify, parsers ...interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		var ns = make([]ParsecNode, 0, len(parsers))
		var n ParsecNode
		news := s.Clone()
		for _, parser := range parsers {
			n, news = doParse(parser, news)
			if n == nil {
				return nil, s
			}
			ns = append(ns, n)
		}
		if node := docallback(callb, ns); node != nil {
			return node, news
		}
		return nil, s
	}
}

// OrdChoice combinator accepts a list of `Parser`, or
// reference to a parser, where atleast one of the parser
// must match the input string. Return a parser function
// that can further be used to construct higher level parsers.
//
// The first matching parser function's output is passed
// as argument to Nodify callback, that is []ParsecNode argument
// will just have one element in it. If none of the parsers
// match the input, then OrdChoice will fail without consuming
// any input.
func OrdChoice(callb Nodify, parsers ...interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		for _, parser := range parsers {
			if n, news := doParse(parser, s.Clone()); n != nil {
				if node := docallback(callb, []ParsecNode{n}); node != nil {
					return node, news
				}
			}
		}
		return nil, s
	}
}

// Kleene combinator accepts two parsers, or reference to
// parsers, namely opScan and sepScan, where opScan parser
// will be used to match input string and contruct ParsecNode,
// and sepScan parser will be used to match input string
// and ignore the matched string. If sepScan parser is not
// supplied, then opScan parser will be applied on the input
// until it fails.
//
// The process of matching opScan parser and sepScan parser
// will continue in a loop until either one of them fails on
// the input stream.
//
// For every successful match of opScan, the returned
// ParsecNode from matching parser will be accumulated and
// passed as argument to Nodify callback. If there is not a
// single match for opScan, then []ParsecNode of ZERO length
// will be passed as argument to Nodify callback. Kleene
// combinator will never fail.
func Kleene(callb Nodify, parsers ...interface{}) Parser {
	var opScan, sepScan interface{}
	switch l := len(parsers); l {
	case 1:
		opScan = parsers[0]
	case 2:
		opScan, sepScan = parsers[0], parsers[1]
	default:
		panic(fmt.Errorf("kleene parser doesn't accept %v parsers", l))
	}
	return func(s Scanner) (ParsecNode, Scanner) {
		var n ParsecNode
		ns := make([]ParsecNode, 0)
		news := s.Clone()
		for {
			if n, news = doParse(opScan, news); n == nil {
				break
			}
			ns = append(ns, n)
			if sepScan != nil {
				if n, news = doParse(sepScan, news); n == nil {
					break
				}
			}
		}
		return docallback(callb, ns), news
	}
}

// Many combinator accepts two parsers, or reference to
// parsers, namely opScan and sepScan, where opScan parser
// will be used to match input string and contruct ParsecNode,
// and sepScan parser will be used to match input string and
// ignore the matched string. If sepScan parser is not
// supplied, then opScan parser will be applied on the input
// until it fails.
//
// The process of matching opScan parser and sepScan parser
// will continue in a loop until either one of them fails on
// the input stream.
//
// The difference between `Many` combinator and `Kleene`
// combinator is that there shall atleast be one match of opScan.
//
// For every successful match of opScan, the returned
// ParsecNode from matching parser will be accumulated and
// passed as argument to Nodify callback. If there is not a
// single match for opScan, then Many will fail without
// consuming the input.
func Many(callb Nodify, parsers ...interface{}) Parser {
	var opScan, sepScan interface{}
	switch l := len(parsers); l {
	case 1:
		opScan = parsers[0]
	case 2:
		opScan, sepScan = parsers[0], parsers[1]
	default:
		panic(fmt.Errorf("many parser doesn't accept %v parsers", l))
	}
	return func(s Scanner) (ParsecNode, Scanner) {
		var n ParsecNode
		ns := make([]ParsecNode, 0)
		news := s.Clone()
		for {
			if n, news = doParse(opScan, news); n == nil {
				break
			}
			ns = append(ns, n)
			if sepScan != nil {
				if n, news = doParse(sepScan, news); n == nil {
					break
				}
			}
		}
		if len(ns) > 0 {
			if node := docallback(callb, ns); node != nil {
				return node, news
			}
		}
		return nil, s
	}
}

// ManyUntil combinator accepts three parsers, or references to
// parsers, namely opScan, sepScan and untilScan, where opScan parser
// will be used to match input string and contruct ParsecNode,
// and sepScan parser will be used to match input string and
// ignore the matched string. If sepScan parser is not
// supplied, then opScan parser will be applied on the input
// until it fails.
//
// The process of matching opScan parser and sepScan parser
// will continue in a loop until either one of them fails on
// the input stream or untilScan matches.
//
// For every successful match of opScan, the returned
// ParsecNode from matching parser will be accumulated and
// passed as argument to Nodify callback. If there is not a
// single match for opScan, then ManyUntil will fail without
// consuming the input.
func ManyUntil(callb Nodify, parsers ...interface{}) Parser {
	var opScan, sepScan, untilScan interface{}
	switch l := len(parsers); l {
	case 2:
		opScan, untilScan = parsers[0], parsers[1]
	case 3:
		opScan, sepScan, untilScan = parsers[0], parsers[1], parsers[2]
	default:
		panic(fmt.Errorf("ManyUntil parser doesn't accept %v parsers", l))
	}
	return func(s Scanner) (ParsecNode, Scanner) {
		var n ParsecNode
		var e ParsecNode
		ns := make([]ParsecNode, 0)
		news := s.Clone()
		for {
			if e, _ = doParse(untilScan, news.Clone()); e != nil {
				break
			}
			if n, news = doParse(opScan, news); n == nil {
				break
			}
			ns = append(ns, n)
			if sepScan != nil {
				if n, news = doParse(sepScan, news); n == nil {
					break
				}
			}
		}
		if len(ns) > 0 {
			if node := docallback(callb, ns); node != nil {
				return node, news
			}
		}
		return nil, s
	}
}

// Maybe combinator accepts a single parser, or reference to
// a parser, and tries to match the input stream with it. If
// parser fails to match the input, returns MaybeNone.
func Maybe(callb Nodify, parser interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		n, news := doParse(parser, s.Clone())
		if n == nil {
			return MaybeNone("missing"), s
		}
		if node := docallback(callb, []ParsecNode{n}); node != nil {
			return node, news
		}
		return MaybeNone("missing"), s
	}
}

//----------------
// Local functions
//----------------

func doParse(parser interface{}, s Scanner) (ParsecNode, Scanner) {
	switch p := parser.(type) {
	case Parser:
		return p(s)
	case *Parser:
		return (*p)(s)
	default:
		panic(fmt.Errorf("type of parser `%T` not supported", parser))
	}
}

func docallback(callb Nodify, ns []ParsecNode) ParsecNode {
	if callb != nil {
		return callb(ns)
	}
	return ns
}
