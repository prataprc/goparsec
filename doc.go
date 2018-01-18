/*
Package parsec provies a library of parser-combinators. The basic
idea behind parsec module is that, it allows programmers to compose
basic set of terminal parsers, a.k.a tokenizers and compose them
together as a tree of parsers, using combinators like: And,
OrdChoice, Kleene, Many, Maybe.

To begin with there are four basic Types that needs to be kept in
mind while creating and composing parsers.

Types

Scanner, an interface type that encapsulates the input text. A built
in scanner called SimpleScanner is supplied along with this package.
Developers can also implement their own scanner types. Following
example create a new instance of SimpleScanner, using an input
text:
	var exprText = []byte(`4 + 123 + 23 + 67 +89 + 87 *78`)
	s := parsec.NewScanner(exprText)

Nodify, callback function is supplied while combining parser
functions. If the underlying parsing logic matches with i/p text,
then callback will be dispatched with list of matching ParsecNode.
Value returned by callback function will further be used as
ParsecNode item in higher-level list of ParsecNodes.

Parser, simple parsers are functions that matches i/p text for
specific patterns. Simple parsers can be combined using one of the
supplied combinators to construct a higher level parser. A parser
function takes a Scanner object and applies the underlying parsing
logic, if underlying logic succeeds Nodify callback is dispatched
and a ParsecNode and a new Scanner object (with its cursor moved
forward) is returned. If parser fails to match, it shall return
the input scanner object as it is along with nil ParsecNode.

ParsecNode, an interface type encapsulates one or more tokens from
i/p text, as terminal node or non-terminal node.

Combinators

If input text is going to be a single token like `10` or `true` or
`"some string"`, then all we need is a single Parser function that
can tokenize the i/p text into a terminal node. But our applications
are seldom that simple. Almost all the time we need to parse the i/p
text for more than one tokens and most of the time we need to
compose them into a tree of terminal and non-terminal nodes.

This is where combinators are useful. Package provides a set of
combinators to help combine terminal parsers into higher level
parsers. They are,

 * And, to combine a sequence of terminals and non-terminal parsers.
 * OrdChoice, to choose between specified list of parsers.
 * Kleene, to repeat the parser zero or more times.
 * Many, to repeat the parser one or more times.
 * ManyUntil, to repeat the parser until a specified end matcher.
 * Maybe, to apply the parser once or none.

All the above mentioned combinators accept one or more parser function
as arguments, either by value or by reference. The reason for allowing
parser argument by reference is to be able to define recursive
parsing logic, like parsing nested arrays:

	var Y Parser
	var value Parser // circular rats

	var opensqrt = Atom("[", "OPENSQRT")
	var closesqrt = Atom("]", "CLOSESQRT")
	var values = Kleene(nil, &value, Atom(",", "COMMA"))
	var array = And(nil, opensqrt, values, closeSqrt)
	func init() {
		value = parsec.OrdChoice(nil, Int(), Bool(), String(), array)
		Y = parsec.OrdChoice(nil, value)
	}

Terminal parsers

 * Char, match a single character skipping leading whitespace.
 * Float, match a float literal skipping leading whitespace.
 * Hex, match a hexadecimal literal skipping leading whitespace.
 * Int, match a decimal number literal skipping leading whitespace.
 * Oct, match a octal number literal skipping leading whitespace.
 * String, match a string literal skipping leading whitespace.
 * Ident, match a identifier token skipping leading whitespace.
 * Atom, match a single atom skipping leading whitespace.
 * AtomExact, match a single atom without skipping leading whitespace.
 * Token, match a single token skipping leading whitespace.
 * TokenExact, match a single token without skipping leading whitespace.
 * OrdToken, match a single token with specified list of alternatives.
 * End, match end of text.
 * NoEnd, match not an end of text.

All of the terminal parsers, except End and NoEnd return Terminal type
as ParsecNode. While End and NoEnd return a boolean type as ParsecNode.

AST and Queryable

This is an experimental feature to use CSS like selectors for quering
an Abstract Syntax Tree (AST). Types, APIs and methods associated with
AST and Queryable are unstable, and are expected to change in future.

While Scanner, Parser, ParsecNode types are re-used in AST and Queryable,
combinator functions are re-implemented as AST methods. Similarly type
ASTNodify is to be used instead of Nodify type. Otherwise all of the
parsec techniques mentioned above are equally applicable on AST.

Additionally, following points are worth noting down while using AST,
 * Combinator methods supplied via AST can be named.

*/
package parsec
