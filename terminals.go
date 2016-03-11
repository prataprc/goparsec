//  Copyright (c) 2013 Couchbase, Inc.

// parsec also supplies a basic set of token parsers
// that can be used to create higher order parser using
// one of the many combinators.

package parsec

import "fmt"
import "strings"
import "unsafe"
import "reflect"

var _ = fmt.Sprintf("dummy")

// String returns a parser function to match a double-quoted
// or single-quoted string in the input stream, skips leading
// white-space.
func String() Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		if !s.Endof() {
			news := s.Clone()
			news.SkipWS()
			txt := news.Remaining()
			tok, readn := scanString(txt)
			if tok == nil || len(tok) == 0 {
				return nil, s
			}
			cursor := news.SetCursor(news.GetCursor() + readn)
			t := &Terminal{
				Name:     "STRING",
				Value:    string(tok),
				Position: cursor,
			}
			return t, news
		}
		return nil, s
	}
}

// Char returns a parser function to match a single character
// in the input stream, skips leading white-space.
func Char() Parser {
	return Token(`'.'`, "CHAR")
}

// Float returns a parser function to match a float literal
// in the input stream, skips leading white-space.
func Float() Parser {
	return Token(`[+-]?[0-9]*\.?[0-9]*`, "FLOAT")
}

// Hex returns a parser function to match a hexadecimal
// literal in the input stream, skip leading white-space.
func Hex() Parser {
	return Token(`0[xX][0-9a-fA-F]+`, "HEX")
}

// Oct returns a parser function to match an octal number
// literal in the input stream, skip leading white-space.
func Oct() Parser {
	return Token(`0[0-7]+`, "OCT")
}

// Int returns a parser function to match an integer literal
// in the input stream, skip leading white-space.
func Int() Parser {
	return Token(`-?[0-9]+`, "INT")
}

// Ident returns a parser function to match an identifier
// token in the input stream, an identifier is matched with
// the following pattern,
//      `^[A-Za-z][0-9a-zA-Z_]*`
func Ident() Parser {
	return Token(`[A-Za-z][0-9a-zA-Z_]*`, "IDENT")
}

// Token takes a pattern and returns a parser that will
// match the pattern with input stream. Input stream will
// be supplied via Scanner interface.
func Token(pattern string, name string) Parser {
	return func(s Scanner) (ParsecNode, Scanner) {
		news := s.Clone()
		news.SkipWS()
		cursor := news.GetCursor()
		if tok, _ := news.Match("^" + pattern); tok != nil {
			t := &Terminal{
				Name:     name,
				Value:    string(tok),
				Position: cursor,
			}
			return t, news
		}
		return nil, s
	}
}

// OrdTokens to parse a single token based on one of the
// specified `patterns`.
func OrdTokens(patterns []string, names []string) Parser {
	var group string
	groups := make([]string, 0, len(patterns))
	for i, pattern := range patterns {
		if names[i] == "" {
			group = "^(" + pattern + ")"
		} else {
			group = "^(?P<" + names[i] + ">" + pattern + ")"
		}
		groups = append(groups, group)
	}
	ordPattern := strings.Join(groups, "|")
	return func(s Scanner) (ParsecNode, Scanner) {
		news := s.Clone()
		news.SkipWS()
		cursor := news.GetCursor()
		if captures, _ := news.SubmatchAll(ordPattern); captures != nil {
			for name, tok := range captures {
				t := Terminal{
					Name:     name,
					Value:    string(tok),
					Position: cursor,
				}
				return &t, news
			}
		}
		return nil, s
	}
}

// End is a parser function to detect end of scanner output.
func End(s Scanner) (ParsecNode, Scanner) {
	return s.Endof(), s
}

// NoEnd is a parser function to detect not-an-end of
// scanner output.
func NoEnd(s Scanner) (ParsecNode, Scanner) {
	return !s.Endof(), s
}

// local functions

func scanString(txt []byte) (tok []byte, readn int) {
	if len(txt) < 2 {
		return nil, 0
	} else if txt[0] != '"' && txt[0] != '\'' {
		return nil, 0
	}
	esc, quote := false, rune(txt[0])
	for i, ch := range bytes2str(txt[1:]) {
		if ch == '\\' {
			esc = true
		} else if esc {
			esc = false
			continue
		} else if ch == quote {
			return txt[:i+2], i + 2
		}
	}
	panic("invalid string")
}

func bytes2str(bytes []byte) string {
	if bytes == nil {
		return ""
	}
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	st := &reflect.StringHeader{Data: sl.Data, Len: sl.Len}
	return *(*string)(unsafe.Pointer(st))
}
