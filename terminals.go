package parsec

import (
	"fmt"
)

var _ = fmt.Sprintf("keep 'fmt' import during debugging")

// Parsec functions to match special strings.
func Token(pattern string, n string) Parser {
	return func(s Scanner) (ParsecNode, *Scanner) {
		news := &s
		if ws, newss := s.SkipWS(); ws != nil {
			news = newss
		}
		if tok, newss := (*news).Match(pattern); tok != nil {
			t := Terminal{Name: n, Value: string(tok), Position: news.cursor}
			return &t, newss
		} else {
			return nil, news
		}
	}
}

// Parsec function to detect end/not-end-of of scanner output.
func End(s Scanner) (ParsecNode, *Scanner) {
	return s.Endof(), &s
}

func NoEnd(s Scanner) (ParsecNode, *Scanner) {
	return !s.Endof(), &s
}

// Parsec functions to match literals
func String() Parser {
	return Token(`^"(\.|[^"])*"`, "STRING")
}
func Char() Parser {
	return Token(`^'.'`, "CHAR")
}
func Int() Parser {
	return Token(`^[0-9]+`, "INT")
}
func Hex() Parser {
	return Token(`^0[xX][0-9a-fA-F]+`, "HEX")
}
func Oct() Parser {
	return Token(`^0[0-8]+`, "OCT")
}
func Float() Parser {
	return Token(`^[0-9]*\.[0-9]+`, "FLOAT")
}

func Ident() Parser {
	return Token(`^[A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}
