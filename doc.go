// Package parsec implements a library of parser-combinators
// using basic recognizers like,
//
//      And OrdChoice Kleene Many Maybe
//
// Parser combinators can be used to construct higher-order
// parsers using basic parser function. All parser functions
// are expected to follow the `Parser` type signature, accepting
// a `Scanner` interface and returning a `ParsecNode` and a new
// scanner. If a parser fails to match the input string according
// to its rules, then it must return nil for ParsecNode, and a
// new Scanner.
//
// ParsecNode can either be a Terminal structure or NonTerminal
// structure or a list of Terminal/NonTerminal structure. The AST
// output is expected to be made up of ParsecNode.
//
// Nodify is a callback function that every combinators use as a
// callback to construct a ParsecNode.
package parsec
