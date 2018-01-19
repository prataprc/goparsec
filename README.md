Parser combinator library in Golang
===================================

[![talk on matrix](https://github.com/prataprc/dotfiles/blob/master/assets/talkonmatrix.svg)](https://riot.im/app/#/user/@prataprc:matrix.org?action=chat)
[![Build Status](https://travis-ci.org/prataprc/goparsec.svg?branch=master)](https://travis-ci.org/prataprc/goparsec)
[![Coverage Status](https://coveralls.io/repos/github/prataprc/goparsec/badge.svg?branch=master)](https://coveralls.io/github/prataprc/goparsec?branch=master)
[![GoDoc](https://godoc.org/github.com/prataprc/goparsec?status.png)](https://godoc.org/github.com/prataprc/goparsec)
[![Sourcegraph](https://sourcegraph.com/github.com/prataprc/goparsec/-/badge.svg)](https://sourcegraph.com/github.com/prataprc/goparsec?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/prataprc/goparsec)](https://goreportcard.com/report/github.com/prataprc/goparsec)

A library to construct top-down recursive backtracking parsers using
parser-combinators. To know more about theory of parser
combinators, [refer here](http://en.wikipedia.org/wiki/Parser_combinator).

This package contains following components,

* Standard set of combinators.
* [Regular expression](https://golang.org/pkg/regexp/) based simple-scanner.
* Standard set of tokenizers.
* To construct syntax-trees based on detailed grammar try with
  [AST struct](http://godoc.org/github.com/prataprc/goparsec#AST).
  * Standard set of combinators are exported as methods to AST.
  * Generate dot-graph EG: [dotfile](testdata/simple.dot)
  for [html](testdata/simple.html).
  * Pretty print on the console.
  * Make debugging easier.

**NOTE that AST object is a recent development and expect user to adapt to
newer versions**

Quick links
-----------

* [Go documentation](http://godoc.org/github.com/prataprc/goparsec)
* [Using the builtin scanner](#using-the-builtin-scanner).
* [Projects using goparsec](#projects-using-goparsec).
* [Simple html parser](#simple-html-parser).
* [Articles](#articles).
* [How to contribute](#how-to-contribute).

Combinators
-----------

Every combinator should confirm to the following signature,

```go
    // ParsecNode type defines a node in the AST
    type ParsecNode interface{}

    // Parser function parses input text, higher order parsers are
    // constructed using combinators.
    type Parser func(Scanner) (ParsecNode, Scanner)

    // Nodify callback function to construct custom ParsecNode.
    type Nodify func([]ParsecNode) ParsecNode
```

Combinators take a variable number of parser functions and
return a new parser function.

If the intermediate nodes are created by the Combinators then it will be of
the following types:

Using the builtin scanner
-------------------------

The builtin scanner library manages the input buffer and implements a cursor
into the buffer. Create a new scanner instance,

```go
    s := parsec.NewScanner(text)
```

The scanner library supplies method receivers like ``Match(pattern)``,
``SkipAny(pattern)`` and ``Endof()``, refer to scanner.go for more information
on each of these methods.

Panics and Recovery
-------------------

Panics are to expected when APIs are misused. Programmers might choose
to ignore the errors, but not panics. For example:

* Combinators accept Parser function or pointer to Parser function. Anything
  else will panic.
* Kleene and Many combinators take one or two parsers as arguments. Less than
  one or more than two will throw a panic.
* ManyUntil combinator take two or three parsers as arguments. Less than two
  or more than three will throw a panic.
* When using invalid regular expression to match a token.


Examples
--------

* expr/expr.go, implements a parsec grammar to parse arithmetic expressions.
* json/json.go, implements a parsec grammar to parse JSON document.

Clone the repository run the benchmark suite

```bash
    $ cd expr/
    $ go test -test.bench=. -test.benchmem=true
    $ cd json/
    $ go test -test.bench=. -test.benchmem=true
```

To run the example program,

```bash
    # to parse expression
    $ go run tools/parsec/parsec.go -expr "10 + 29"

    # to parse JSON string
    $ go run tools/parsec/parsec.go -json '{ "key1" : [10, "hello", true, null, false] }'
```

Simple html parser
------------------

```go
func makehtmly(ast *AST) Parser {
	var tag Parser

	opentag := AtomExact("<", "OT")
	closetag := AtomExact(">", "CT")
	equal := AtomExact("=", "EQUAL")
	slash := TokenExact("/[ \t]*", "SLASH")
	tagname := TokenExact("[a-z][a-zA-Z0-9]*", "TAG")
	attrkey := TokenExact("[a-z][a-zA-Z0-9]*", "ATTRK")
	text := TokenExact("[^<>]+", "TEXT")
	ws := TokenExact("[ \t]+", "WS")

	element := ast.OrdChoice("element", nil, text, &tag)
	elements := ast.Kleene("elements", nil, element)
	attr := ast.And("attribute", nil, attrkey, equal, String())
	attrws := ast.And("attrws", nil, attr, ast.Maybe("ws", nil, ws))
	attrs := ast.Kleene("attributes", nil, attrws)
	tstart := ast.And("tagstart", nil, opentag, tagname, attrs, closetag)
	tend := ast.And("tagend", nil, opentag, slash, tagname, closetag)
	tag = ast.And("tag", nil, tstart, elements, tend)
	return tag
}
```

Projects using goparsec
-----------------------

* [Monster](https://github.com/prataprc/monster), production system in golang.
* [GoLedger](https://github.com/tn47/goledger), ledger re-write in golang.

If your project is using goparsec you can raise an issue to list them under
this section.

Articles
--------

* [Parsing by composing functions](http://prataprc.github.io/parser-combinator-composition.html)
* [Parser composition for recursive grammar](http://prataprc.github.io/parser-combinator-recursive.html)
* [How to use the ``Maybe`` combinator](http://prataprc.github.io/parser-combinator-maybe.html)

How to contribute
-----------------

[![Issue Stats](http://issuestats.com/github/prataprc/goparsec/badge/pr)](http://issuestats.com/github/prataprc/goparsec)
[![Issue Stats](http://issuestats.com/github/prataprc/goparsec/badge/issue)](http://issuestats.com/github/prataprc/goparsec)

* Pick an issue, or create an new issue. Provide adequate documentation for
  the issue.
* Assign the issue or get it assigned.
* Work on the code, once finished, raise a pull request.
* Goparsec is written in [golang](https://golang.org/), hence expected to follow the
  global guidelines for writing go programs.
* If the changeset is more than few lines, please generate a
  [report card](https://goreportcard.com/report/github.com/prataprc/goparsec).
* As of now, branch ``master`` is the development branch.
