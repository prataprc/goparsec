// parsec also supplies a basic set of token parsers that can be used to
// create higher order parser using one of the many combinators.
package parsec

// String returns a parser function to match a double quoted string in the
// input stream.
func String() Parser {
	return Token(`^"(\.|[^"])*"`, "STRING")
}

// Char returns a parser function to match a single character in the input
// stream.
func Char() Parser {
	return Token(`^'.'`, "CHAR")
}

// Int returns a parser function to match an integer literal in the input
// stream.
func Int() Parser {
	return Token(`^[0-9]+`, "INT")
}

// Hex returns a parser function to match a hexadecimal literal in the input
// stream.
func Hex() Parser {
	return Token(`^0[xX][0-9a-fA-F]+`, "HEX")
}

// Oct returns a parser function to match an octal number literal in the input
// stream.
func Oct() Parser {
	return Token(`^0[0-8]+`, "OCT")
}

// Float returns a parser function to match a float literal in the input
// stream.
func Float() Parser {
	return Token(`^[0-9]*\.[0-9]+`, "FLOAT")
}

// Ident returns a parser function to match an identifier token in the input
// stream, an identifier is matched with the following pattern
// `^[A-Za-z][0-9a-zA-Z_]*`
func Ident() Parser {
	return Token(`^[A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}

// Token takes a pattern and returns a parser that will match the pattern with
// input stream. Input stream will be supplied via Scanner interface.
func Token(pattern string, name string) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		news := s.Clone()
		news.SkipWS()
		if tok, _ := news.Match(pattern); tok != nil {
			t := Terminal{
				Name:     name,
				Value:    string(tok),
				Position: news.GetCursor(),
			}
			return &t, news
		}
		return nil, s
	}
}

// End is a parser function to detect end of scanner output.
func End(s Scanner) (ParsecNode, Scanner) {
	return s.Endof(), s
}

// NoEnd is a parser function to detect not-an-end of scanner output.
func NoEnd(s Scanner) (ParsecNode, Scanner) {
	return !s.Endof(), s
}
