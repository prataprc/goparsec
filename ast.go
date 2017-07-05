// Copyright (c) 2013 Couchbase, Inc.

package parsec

import "fmt"
import "errors"

// Queryable interface to be implemented by all nodes,
// both Terminal nodes and NonTerminal nodes.
type Queryable interface {
	// GetName for the node.
	GetName() string

	// IsTerminal return true of node is a leaf node in syntax-tree.
	IsTerminal() bool

	// GetValue return parsed text, if node is NonTerminal it will
	// concat the entire sub-tree for parsed text and return the same.
	GetValue() string

	// GetChildren relevant only for NonTerminal node.
	GetChildren() []Queryable

	// GetPosition of the first terminal value in input.
	GetPosition() int
}

// ASTNodify callback function to construct custom Queryable. Even when
// combinators line And, OrdChoice, Many etc.. match input string, it is
// possible to fail them via ASTNodify callback function, by returning nil.
// This very useful in cases like:
//
//	* where lookahead matching is required.
//  * exceptional cases for a regex pattern.
//
// Note that some combinators like Kleene shall not interpret the return
// value from ASTNodify callback.
type ASTNodify func(name string, node Queryable) Queryable

// AST to parse and construct Abstract Syntax Tree whose nodes confirm
// to `Queryable` interfaces facilitating tree processing algorithms.
type AST struct {
	name   string
	y      Parser
	root   Queryable
	ntpool chan *NonTerminal
	debug  bool
}

// NewAST return a new instance of AST, maxnodes is size of internal buffer
// pool of nodes. It is directly proportional to number of nodes that you
// expect in the syntax-tree.
func NewAST(name string, maxnodes int) *AST {
	return &AST{name: name, ntpool: make(chan *NonTerminal, maxnodes)}
}

// Debug enables console logging while parsing the input test, this is
// useful while developing a parser.
func (ast *AST) SetDebug() *AST {
	ast.debug = true
	return ast
}

// Parsewith execute the root parser, y, with scanner s. AST will remember the
// root parser, and root node. Return the root-node, and scanner with
// remaining input if parser was successful, other wise nil.
func (ast *AST) Parsewith(y Parser, s Scanner) (Queryable, Scanner) {
	ast.root, ast.y = nil, y
	node, news := y(s)
	if node == nil {
		return nil, news
	}
	ast.root = node.(Queryable)
	return ast.root, news
}

// Reset the AST, forget the root parser, and root node. Reuse the AST object
// via Parsewith different set of root-parser and scanner.
func (ast *AST) Reset() *AST {
	var freetree func(*NonTerminal)

	freetree = func(node *NonTerminal) {
		for _, q := range node.Children {
			if nt, ok := q.(*NonTerminal); ok {
				freetree(nt)
			}
		}
		ast.putnt(node)
	}
	if node, ok := ast.root.(*NonTerminal); ok {
		freetree(node)
	}
	ast.y, ast.root = nil, nil
	return ast
}

// And combinator, name identifies the NonTerminal constructed by this
// combinator.
func (ast *AST) And(name string, callb ASTNodify, parsers ...interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(name), s.Clone()
		for i, parser := range parsers {
			if node, news, err = ast.doParse(parser, news); err != nil {
				fmsg := "while parsing %vth in %q: %v"
				panic(fmt.Errorf(fmsg, i+1, name, err))
			} else if node == nil {
				ast.putnt(nt)
				return ast.trydebug(nil, s, "And", name, i+1, false)
			}
			ast.trydebug(node, news, "And", name, i+1, true)
			nt.Children = append(nt.Children, node.(Queryable))
		}
		if q := ast.docallback(name, callb, nt); q != nil {
			return ast.trydebug(q, news, "And", name, -1, true)
		}
		ast.putnt(nt)
		return ast.trydebug(nil, s, "And", name, -1, "skip")
	}
}

// OrdChoice combinator.
func (ast *AST) OrdChoice(
	name string, callb ASTNodify, parsers ...interface{}) Parser {

	return func(s Scanner) (ParsecNode, Scanner) {
		for i, parser := range parsers {
			news := s.Clone()
			if n, news, err := ast.doParse(parser, news); err != nil {
				fmsg := "while parsing %vth for %q: %v"
				panic(fmt.Errorf(fmsg, i+1, name, err))
			} else if n != nil {
				if q := ast.docallback(name, callb, n.(Queryable)); q != nil {
					return ast.trydebug(q, news, "OrdChoice", name, i+1, true)
				}
				return ast.trydebug(nil, s, "OrdChoice", name, i+1, "skip")
			}
		}
		return ast.trydebug(nil, s, "OrdChoice", name, -1, false)
	}
}

// Kleene combinator, name identifies the NonTerminal constructed by this
// combinator.
func (ast *AST) Kleene(
	name string, callb ASTNodify, parsers ...interface{}) Parser {

	var opScan, sepScan interface{}
	switch l := len(parsers); l {
	case 1:
		opScan = parsers[0]
	case 2:
		opScan, sepScan = parsers[0], parsers[1]
	default:
		fmsg := "kleene parser %q doesn't accept %v parsers (should be 1 or 2)"
		panic(fmt.Errorf(fmsg, name, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(name), s.Clone()
		for {
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", name, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", name, err))
				} else if node == nil {
					break
				}
			}
		}
		return ast.docallback(name, callb, nt), news
	}
}

// Many combinator, name identifies the NonTerminal constructed by this
// combinator.
func (ast *AST) Many(
	name string, callb ASTNodify, parsers ...interface{}) Parser {

	var opScan, sepScan interface{}
	switch l := len(parsers); l {
	case 1:
		opScan = parsers[0]
	case 2:
		opScan, sepScan = parsers[0], parsers[1]
	default:
		fmsg := "many parser %q doesn't accept %v parsers (should be 1 or 2)"
		panic(fmt.Errorf(fmsg, name, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(name), s.Clone()
		for {
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", name, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", name, err))
				} else if node == nil {
					break
				}
			}
		}
		if len(nt.Children) > 0 {
			if q := ast.docallback(name, callb, nt); q != nil {
				return q, news
			}
		}
		ast.putnt(nt)
		return nil, s
	}
}

// ManyUntil combinator, name identifies the NonTerminal constructed by this
// combinator.
func (ast *AST) ManyUntil(
	name string, callb ASTNodify, parsers ...interface{}) Parser {

	var opScan, sepScan, untilScan interface{}
	switch l := len(parsers); l {
	case 2:
		opScan, untilScan = parsers[0], parsers[1]
	case 3:
		opScan, sepScan, untilScan = parsers[0], parsers[1], parsers[2]
	default:
		fmsg := "ManyUntil parser %q don't accept %v parsers (should be 2 or 3)"
		panic(fmt.Errorf(fmsg, name, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(name), s.Clone()
		for {
			if node, _, err = ast.doParse(untilScan, news.Clone()); err != nil {
				panic(fmt.Errorf("while untilscan-parsing %q: %v", name, err))
			} else if node != nil {
				break
			}
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", name, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", name, err))
				} else if node == nil {
					break
				}
			}
		}
		if len(nt.Children) > 0 {
			if q := ast.docallback(name, callb, nt); q != nil {
				return q, news
			}
		}
		ast.putnt(nt)
		return nil, s
	}
}

// Maybe combinator.
func (ast *AST) Maybe(name string, callb ASTNodify, parser interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		node, news, err := ast.doParse(parser, s.Clone())
		if err != nil {
			panic(fmt.Errorf("while parsing %q: %v", name, err))
		} else if node == nil {
			return MaybeNone("missing"), s
		}
		if q := ast.docallback(name, callb, node.(Queryable)); q != nil {
			return q, news
		}
		return MaybeNone("missing"), s
	}
}

//---- local functions

func (ast *AST) doParse(
	parser interface{}, s Scanner) (ParsecNode, Scanner, error) {

	switch p := parser.(type) {
	case Parser:
		node, news := p(s)
		return node, news, nil
	case *Parser:
		node, news := (*p)(s)
		return node, news, nil
	default:
		return nil, s, errors.New("badtype")
	}
}

func (ast *AST) docallback(
	name string, callb ASTNodify, node Queryable) Queryable {

	if callb != nil {
		q := callb(name, node)
		if q == nil {
			return nil
		} else if _, ok := q.(*NonTerminal); !ok {
			ast.putnt(node.(*NonTerminal))
		}
		return q
	}
	return node
}

func (ast *AST) getnt(name string) (node *NonTerminal) {
	select {
	case node = <-ast.ntpool:
	default:
		node = &NonTerminal{Name: name, Children: make([]Queryable, 0)}
	}
	return node
}

func (ast *AST) putnt(node *NonTerminal) {
	node.Children = node.Children[:0]
	select {
	case ast.ntpool <- node:
	default: // node shall be collected by GC.
	}
}

func (ast *AST) trydebug(
	node ParsecNode, s Scanner,
	ytype, name string, poff int, match interface{}) (ParsecNode, Scanner) {

	fmsg := "%v(%v) parser:%v Lineno:%v off:%v match:%v\n"
	if ast.debug {
		fmt.Printf(fmsg, ytype, name, poff, s.Lineno(), s.GetCursor(), match)
	}
	return node, s
}
