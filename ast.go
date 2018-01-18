// Copyright (c) 2013 Goparsec AUTHORS. All rights reserved.
// Use of this source code is governed by LICENSE file.

package parsec

import "io"
import "os"
import "fmt"
import "sort"
import "strings"
import "errors"

// Queryable interface to be implemented by all nodes, both terminal
// and non-terminal nodes constructed using AST object.
type Queryable interface {
	// GetName for the node.
	GetName() string

	// IsTerminal return true if node is a leaf node in syntax-tree.
	IsTerminal() bool

	// GetValue return parsed text, if node is NonTerminal it will
	// concat the entire sub-tree for parsed text and return the same.
	GetValue() string

	// GetChildren relevant only for NonTerminal node.
	GetChildren() []Queryable

	// GetPosition of the first terminal value in input.
	GetPosition() int

	// SetAttribute with a value string, can be called multiple times for the
	// same attrname.
	SetAttribute(attrname, value string) Queryable

	// GetAttribute for attrname, since more than one value can be set on the
	// attribute, return a slice of values.
	GetAttribute(attrname string) []string

	// GetAttributes return a map of all attributes set on this node.
	GetAttributes() map[string][]string
}

// ASTNodify callback function to construct custom Queryable. Even when
// combinators like And, OrdChoice, Many etc.. match input string, it is
// possible to fail them via ASTNodify callback function, by returning nil.
// This is useful in cases like:
//  * where lookahead matching is required.
//  * exceptional cases for a regex pattern.
//
// Note that some combinators like Kleene shall not interpret the return
// value from ASTNodify callback.
type ASTNodify func(name string, s Scanner, node Queryable) Queryable

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
// pool of nodes, it is directly proportional to number of nodes that you
// expect in the syntax-tree.
func NewAST(name string, maxnodes int) *AST {
	return &AST{name: name, ntpool: make(chan *NonTerminal, maxnodes)}
}

// SetDebug enables console logging while parsing the input test, this is
// useful while developing a parser.
func (ast *AST) SetDebug() *AST {
	ast.debug = true
	return ast
}

// Parsewith execute the root parser, y, with scanner s. AST will
// remember the root parser, and root node. Return the root-node as
// Queryable, if success and scanner with remaining input.
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

// And combinator, same as package level And combinator function.
// `name` identifies the NonTerminal constructed by this combinator.
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
		if q := ast.docallback(name, callb, news, nt); q != nil {
			return ast.trydebug(q, news, "And", name, -1, true)
		}
		ast.putnt(nt)
		return ast.trydebug(nil, s, "And", name, -1, "skip")
	}
}

// OrdChoice combinator, same as package level OrdChoice combinator
// function. `nm` identifies the NonTerminal constructed by this combinator.
func (ast *AST) OrdChoice(nm string, cb ASTNodify, ps ...interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		for i, parser := range ps {
			news := s.Clone()
			if n, news, err := ast.doParse(parser, news); err != nil {
				fmsg := "while parsing %vth for %q: %v"
				panic(fmt.Errorf(fmsg, i+1, nm, err))
			} else if n != nil {
				q := ast.docallback(nm, cb, news, n.(Queryable))
				if q != nil {
					return ast.trydebug(q, news, "OrdChoice", nm, i+1, true)
				}
				return ast.trydebug(nil, s, "OrdChoice", nm, i+1, "skip")
			}
		}
		return ast.trydebug(nil, s, "OrdChoice", nm, -1, false)
	}
}

// Kleene combinator, same as package level Kleene combinator
// function. `nm` identifies the NonTerminal constructed by this combinator.
func (ast *AST) Kleene(nm string, callb ASTNodify, ps ...interface{}) Parser {
	var opScan, sepScan interface{}
	switch l := len(ps); l {
	case 1:
		opScan = ps[0]
	case 2:
		opScan, sepScan = ps[0], ps[1]
	default:
		fmsg := "kleene parser %q doesn't accept %v parsers (should be 1 or 2)"
		panic(fmt.Errorf(fmsg, nm, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(nm), s.Clone()
		for {
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", nm, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", nm, err))
				} else if node == nil {
					break
				}
			}
		}
		return ast.docallback(nm, callb, news, nt), news
	}
}

// Many combinator, same as package level Many combinator
// function. `nm` identifies the NonTerminal constructed by this combinator.
func (ast *AST) Many(nm string, callb ASTNodify, parsers ...interface{}) Parser {
	var opScan, sepScan interface{}
	switch l := len(parsers); l {
	case 1:
		opScan = parsers[0]
	case 2:
		opScan, sepScan = parsers[0], parsers[1]
	default:
		fmsg := "many parser %q doesn't accept %v parsers (should be 1 or 2)"
		panic(fmt.Errorf(fmsg, nm, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(nm), s.Clone()
		for {
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", nm, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", nm, err))
				} else if node == nil {
					break
				}
			}
		}
		if len(nt.Children) > 0 {
			if q := ast.docallback(nm, callb, news, nt); q != nil {
				return q, news
			}
		}
		ast.putnt(nt)
		return nil, s
	}
}

// Many combinator, same as package level Many combinator
// function. `nm` identifies the NonTerminal constructed by this combinator.
func (ast *AST) ManyUntil(nm string, callb ASTNodify, ps ...interface{}) Parser {
	var opScan, sepScan, untilScan interface{}
	switch l := len(ps); l {
	case 2:
		opScan, untilScan = ps[0], ps[1]
	case 3:
		opScan, sepScan, untilScan = ps[0], ps[1], ps[2]
	default:
		fmsg := "ManyUntil parser %q don't accept %v parsers (should be 2 or 3)"
		panic(fmt.Errorf(fmsg, nm, l))
	}

	return func(s Scanner) (ParsecNode, Scanner) {
		var node ParsecNode
		var err error
		nt, news := ast.getnt(nm), s.Clone()
		for {
			if node, _, err = ast.doParse(untilScan, news.Clone()); err != nil {
				panic(fmt.Errorf("while untilscan-parsing %q: %v", nm, err))
			} else if node != nil {
				break
			}
			if node, news, err = ast.doParse(opScan, news); err != nil {
				panic(fmt.Errorf("while opscan-parsing %q: %v", nm, err))
			} else if node == nil {
				break
			}
			nt.Children = append(nt.Children, node.(Queryable))
			if sepScan != nil {
				if node, news, err = ast.doParse(sepScan, news); err != nil {
					panic(fmt.Errorf("while sepscan-parsing %q: %v", nm, err))
				} else if node == nil {
					break
				}
			}
		}
		if len(nt.Children) > 0 {
			if q := ast.docallback(nm, callb, news, nt); q != nil {
				return q, news
			}
		}
		ast.putnt(nt)
		return nil, s
	}
}

// Maybe combinator, same as package level Maybe combinator
// function. `nm` identifies the NonTerminal constructed by this combinator.
func (ast *AST) Maybe(name string, callb ASTNodify, parser interface{}) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		node, news, err := ast.doParse(parser, s.Clone())
		if err != nil {
			panic(fmt.Errorf("while parsing %q: %v", name, err))
		} else if node == nil {
			return MaybeNone("missing"), s
		}
		if q := ast.docallback(name, callb, news, node.(Queryable)); q != nil {
			return q, news
		}
		return MaybeNone("missing"), s
	}
}

// End is a parser function to detect end of scanner output.
func (ast *AST) End(name string) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		if s.Endof() {
			return NewTerminal(name, "", s.GetCursor()), s
		}
		return nil, s
	}
}

// GetValue return the full text, called as value here, that was parsed
// to contruct this syntax-tree.
func (ast *AST) GetValue() string {
	return ast.root.GetValue()
}

// Prettyprint to standard output the syntax-tree in human readable plain
// text.
func (ast *AST) Prettyprint() {
	if ast.root == nil {
		fmt.Println("root is nil")
	}
	ast.prettyprint(os.Stdout, "", ast.root)
}

// Dotstring return AST in graphviz dot format. Save this string to a
// dot file and use graphviz tool generate a nice looking graph.
func (ast *AST) Dotstring(name string) string {
	return ast.dottext(name)
}

// Query is an experimental method on AST. Developers can use the
// selector specification to pick one or more nodes from the AST.
func (ast *AST) Query(selectors string, ch chan Queryable) {
	selast := NewAST("selectorast", 100)
	y := parseselector(selast)
	qsel, _ := selast.Parsewith(y, NewScanner([]byte(selectors)))
	orsels := qsel.GetChildren()
	for _, orsel := range orsels {
		qs := orsel.GetChildren()
		astwalk(nil, 0, ast.root, qs, ch)
	}
	close(ch)
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
	name string, callb ASTNodify, s Scanner, node Queryable) Queryable {

	if callb != nil {
		q := callb(name, s, node)
		if q == nil {
			return nil
		} else if _, ok := q.(*NonTerminal); !ok {
			if nt, ok := node.(*NonTerminal); ok {
				ast.putnt(nt)
			}
		}
		return q
	}
	return node
}

func (ast *AST) getnt(name string) (node *NonTerminal) {
	select {
	case node = <-ast.ntpool:
		node.Name = name
	default:
		node = NewNonTerminal(name)
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

func (ast *AST) prettyprint(w io.Writer, prefix string, node Queryable) {
	if node.IsTerminal() {
		fmt.Fprintf(w, "%v*%v: %q\n", prefix, node.GetName(), node.GetValue())
		return
	} else {
		fmt.Fprintf(w, "%v%v @ %v\n", prefix, node.GetName(), node.GetPosition())
		for _, child := range node.GetChildren() {
			ast.prettyprint(w, prefix+"  ", child)
		}
	}
}

type tnode map[int]string

func (ast *AST) dottext(name string) string {
	lines := []string{
		fmt.Sprintf("digraph %v {", name),
		fmt.Sprintf("  nodesep=0.3;"),
		fmt.Sprintf("  ranksep=0.2;"),
		fmt.Sprintf("  margin=0.1;"),
		fmt.Sprintf("  edge [arrowsize=0.8];"),
	}
	nodesi, nodesk := make(tnode), make(tnode)
	edges, nodesi, nodesk, _ :=
		ast.dotline(0, 1, ast.root, []string{}, nodesi, nodesk)
	lines = append(lines, edges...)
	for _, node := range sortnodes(nodesi) {
		label := nodesi[node]
		s := fmt.Sprintf(`  %v [shape=ellipse,label=%q];`, node, label)
		lines = append(lines, s)
	}
	for _, node := range sortnodes(nodesk) {
		label := nodesk[node]
		t := `  %v [shape=ellipse,style=filled,fillcolor=grey,label=%q];`
		lines = append(lines, fmt.Sprintf(t, node, label))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func (ast *AST) dotline(
	parid, nextid int, node Queryable,
	edges []string,
	nodesi, nodesk tnode) ([]string, tnode, tnode, int) {

	nodeid, nextid := nextid, nextid+1
	if node.IsTerminal() {
		name, edge := node.GetName(), fmt.Sprintf("  %v -> %v;", parid, nodeid)
		edges = append(edges, edge)
		nodesk[nodeid] = fmt.Sprintf("%v: %q", name, node.GetValue())
		return edges, nodesi, nodesk, nextid
	}
	edge := fmt.Sprintf("  %v -> %v;", parid, nodeid)
	edges = append(edges, edge)
	nodesi[nodeid] = node.GetName()
	for _, child := range node.GetChildren() {
		edges, nodesi, nodesk, nextid =
			ast.dotline(nodeid, nextid, child, edges, nodesi, nodesk)
	}
	return edges, nodesi, nodesk, nextid
}

func sortnodes(ns tnode) []int {
	ints := []int{}
	for i := range ns {
		ints = append(ints, i)
	}
	sort.Ints(ints)
	return ints
}
