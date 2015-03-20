// Copyright (c) 2013 Couchbase, Inc.

// Package parsec implements a library of parser-combinators
// using basic recognizers like,
//
//      And OrdChoice Kleene Many Maybe
//
// Parser combinators can be used to construct higher-order
// parsers using basic parser function. All parser functions are
// expected to follow the `Parser` type signature, accepting a
// `Scanner` interface and returning a `ParsecNode` and a new
// scanner. If a parser fails to match the input string according
// to its rules, then it must return nil for ParsecNode, and a new
// Scanner.
//
// ParsecNode can either be a Terminal structure or NonTerminal
// structure or a list of Terminal/NonTerminal structure. The AST
// output is expected to be made up of ParsecNode.
//
// Nodify is a callback function that every combinators use as a
// callback to construct a ParsecNode.
package parsec

import "fmt"

// ParsecNode type defines a node in the AST
type ParsecNode interface{}

// Parser function parses input text, higher order parsers are
// constructed using combinators.
type Parser func(Scanner) (ParsecNode, Scanner)

// Nodify callback function to construct custom ParsecNode.
type Nodify func([]ParsecNode) ParsecNode

// Terminal structure can be used to construct a terminal
// ParsecNode.
type Terminal struct {
	Name     string // contains terminal's token type
	Value    string // value of the terminal
	Position int    // Offset into the text stream where token was identified
}

// NonTerminal structure can be used to construct a
// non-terminal ParsecNode.
type NonTerminal struct {
	Name     string       // contains terminal's token type
	Value    string       // value of the terminal
	Children []ParsecNode // list of children to this node.
}

// Scanner interface supplies necessary methods to match the
// input stream.
type Scanner interface {
	// Clone will return new clone of the underlying scanner structure.
	// This will be used by combinators to _backtrack_.
	Clone() Scanner

	// GetCursor gets the current cursor position inside input text.
	GetCursor() int

	// Match the input stream with `pattern` and return
	// matching string after advancing the cursor.
	Match(pattern string) ([]byte, Scanner)

	// SubmatchAll the input stream with a choice of `patterns`
	// and return matching string and submatches, after
	// advancing the cursor.
	SubmatchAll(pattern string) ([][]byte, Scanner)

	// SkipWs skips white space characters in the input stream.
	// Return skipped whitespaces as byte-slice and advance the cursor.
	SkipWS() ([]byte, Scanner)

	// Endof detects whether end-of-file is reached in the input
	// stream and return a boolean indicating the same.
	Endof() bool
}

// And combinator accepts a list of `Parser`, or reference to a
// parser, that must match the input string, atleast until the
// last Parser argument. Returns a parser function that
// can be used to construct higher-level parsers.
//
// If all parser matches, a list of ParsecNode, where each
// ParsecNode is constructed by matching parser, will be passed
// as argument to Nodify callback. Even if one of the input
// parser function fails then empty slice of ParsecNode will
// be supplied as argument to Nodify callback.
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
		return docallback(callb, ns), news
	}
}

// OrdChoice combinator accepts a list of `Parser`, or
// reference to a parser, where atleast one of the parser
// must match the input string. Returns a parser function
// that can be used to construct higher level parsers.
//
// The first matching parser function's output is passed
// as argument to Nodify callback. If non of the parsers
// match the input, then `nil` is returned for ParsecNode
func OrdChoice(callb Nodify, parsers ...interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		for _, parser := range parsers {
			if n, news := doParse(parser, s.Clone()); n != nil {
				return docallback(callb, []ParsecNode{n}), news
			}
		}
		return nil, s
	}
}

// Kleene combinator accepts two parsers, or reference to
// parsers, namely opScan and sepScan, where opScan parser
// will be used to match input string and contruct ParsecNode
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
// single match for opScan, then an empty slice of ParsecNode
// will be passed as argument to Nodify callback.
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
// will be used to match input string and contruct ParsecNode
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
// single match for opScan, then `nil` will be returned for
// ParsecNode.
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
		if len(ns) == 0 {
			return nil, s
		}
		return docallback(callb, ns), news
	}
}

// Maybe combinator accepts a single parser, or reference to
// a parser, and tries to match the input stream with it.
func Maybe(callb Nodify, parser interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		n, news := doParse(parser, s.Clone())
		if n == nil {
			return nil, s
		}
		return docallback(callb, []ParsecNode{n}), news
	}
}

//---------------
// Local function
//---------------

func doParse(parser interface{}, s Scanner) (ParsecNode, Scanner) {
	switch p := parser.(type) {
	case Parser:
		return p(s)
	case *Parser:
		return (*p)(s)
	default:
		panic(fmt.Errorf("type of parser %T not supported", parser))
	}
	return nil, nil
}

func docallback(callb Nodify, ns []ParsecNode) ParsecNode {
	if callb != nil {
		return callb(ns)
	}
	return ns
}
