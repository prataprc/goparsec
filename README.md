Parser combinator library in Golang
===================================

[![IRC #go-nuts](https://www.irccloud.com/invite-svg?channel=%23go-nuts&amp;hostname=irc.mozilla.org&amp;port=6697)](https://www.irccloud.com/invite?channel=%23go-nuts&amp;hostname=irc.mozilla.org&amp;port=6697)
[![Build Status](https://travis-ci.org/prataprc/goparsec.svg?branch=master)](https://travis-ci.org/prataprc/goparsec)
[![Coverage Status](https://coveralls.io/repos/github/prataprc/goparsec/badge.svg?branch=master)](https://coveralls.io/github/prataprc/goparsec?branch=master)
[![GoDoc](https://godoc.org/github.com/prataprc/goparsec?status.png)](https://godoc.org/github.com/prataprc/goparsec)
[![Sourcegraph](https://sourcegraph.com/github.com/prataprc/goparsec/-/badge.svg)](https://sourcegraph.com/github.com/prataprc/goparsec?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/prataprc/goparsec)](https://goreportcard.com/report/github.com/prataprc/goparsec)

A library to construct top-down recursive backtracking parsers using
parser-combinators. Before proceeding you might want to take at peep
at [theory of parser combinators][theory-link]. As for this package, it
provides:

* A standard set of combinators.
* [Regular expression][regexp-link] based simple-scanner.
* Standard set of tokenizers based on the simple-scanner.

To construct syntax-trees based on detailed grammar try with
[AST struct][ast-link]

* Standard set of combinators are exported as methods to AST.
* Generate dot-graph EG: [dotfile](testdata/simple.dot)
  for [html](testdata/simple.html).
* Pretty print on the console.
* Make debugging easier.

**NOTE that AST object is a recent development and expect user to adapt to
newer versions**

Quick links
-----------

* [Go documentation][goparsec-godoc-link].
* [Using the builtin scanner](#using-the-builtin-scanner).
* [Simple HTML parser][htmlparsec-link].
* [Querying Abstract Syntax Tree][astquery-link]
* [Projects using goparsec](#projects-using-goparsec).
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

Using the builtin scanner
-------------------------

Builtin scanner library manages the input buffer and implements a cursor
into the buffer. Create a new scanner instance,

```go
    s := parsec.NewScanner(text)
```

The scanner library supplies method like ``Match(pattern)``,
``SkipAny(pattern)`` and ``Endof()``, [refer][goparsec-godoc-link] to for
more information on each of these methods.

Panics and Recovery
-------------------

Panics are to be expected when APIs are misused. Programmers might choose
to ignore errors, but not panics. For example:

* Kleene and Many combinators take one or two parsers as arguments. Less than
  one or more than two will throw a panic.
* ManyUntil combinator take two or three parsers as arguments. Less than two
  or more than three will throw a panic.
* Combinators accept Parser function or pointer to Parser function. Anything
  else will panic.
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

Projects using goparsec
-----------------------

* [Monster](https://github.com/prataprc/monster), production system in golang.
* [GoLedger](https://github.com/tn47/goledger), ledger re-write in golang.

If your project is using goparsec you can raise an issue to list them under
this section.

Articles
--------

* [Parsing by composing functions][article1-link]
* [Parser composition for recursive grammar][article2-link]
* [How to use the ``Maybe`` combinator][article3-link]

How to contribute
-----------------

[![Issue Stats](http://issuestats.com/github/prataprc/goparsec/badge/pr)](http://issuestats.com/github/prataprc/goparsec)
[![Issue Stats](http://issuestats.com/github/prataprc/goparsec/badge/issue)](http://issuestats.com/github/prataprc/goparsec)

* Pick an issue, or create an new issue. Provide adequate documentation for
  the issue.
* Assign the issue or get it assigned.
* Work on the code, once finished, raise a pull request.
* Goparsec is written in [golang](https://golang.org/), hence expected to
  follow the global guidelines for writing go programs.
* If the changeset is more than few lines, please generate a
  [report card][report-link].
* As of now, branch ``master`` is the development branch.

[theory-link]: http://en.wikipedia.org/wiki/Parser_combinator
[regexp-link]: https://golang.org/pkg/regexp
[ast-link]: http://godoc.org/github.com/prataprc/goparsec#AST
[goparsec-godoc-link]: http://godoc.org/github.com/prataprc/goparsec
[htmlparsec-link]: https://github.com/prataprc/goparsec/blob/master/html_test.go
[astquery-link]: https://prataprc.github.io/astquery.io
[article1-link](http://prataprc.github.io/parser-combinator-composition.html)
[article2-link](http://prataprc.github.io/parser-combinator-recursive.html)
[article3-link](http://prataprc.github.io/parser-combinator-maybe.html)
[report-link]: https://goreportcard.com/report/github.com/prataprc/goparsec
