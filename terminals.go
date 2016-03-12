//  Copyright (c) 2013 Couchbase, Inc.

// parsec also supplies a basic set of token parsers
// that can be used to create higher order parser using
// one of the many combinators.

package parsec

import "fmt"
import "strings"
import "unsafe"
import "reflect"
import "strconv"
import "unicode"
import "unicode/utf8"
import "unicode/utf16"

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
			if len(txt) == 0 || txt[0] != '"' && txt[0] != '\'' {
				return nil, s
			}
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
	return Token(`[+-]?[0-9]*\.[0-9]+`, "FLOAT")
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

var escapeCode = [256]byte{ // TODO: size can be optimized
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'\'': '\'',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func scanString(txt []byte) ([]byte, int) {
	if len(txt) < 2 {
		return nil, 0
	}

	quote, e := txt[0], 1
	for txt[e] != quote {
		c := txt[e]
		if c == '\\' || c == quote || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			e++
			continue
		}
		r, size := utf8.DecodeRune(txt[e:])
		if r == utf8.RuneError && size == 1 {
			return nil, 0
		}
		e += size
	}

	if txt[e] == quote { // done we have nothing to unquote
		return txt[:e+1], e + 1
	}

	out := make([]byte, len(txt)+2*utf8.UTFMax)
	oute := copy(out, txt[:e]) // copy so far

loop:
	for e < len(txt) {
		switch c := txt[e]; {
		case c == quote:
			out[oute] = c
			e++
			break loop

		case c == '\\':
			if txt[e+1] == 'u' {
				r := getu4(txt[e:])
				if r < 0 { // invalid
					return nil, 0
				}
				e += 6
				if utf16.IsSurrogate(r) {
					nextr := getu4(txt[e:])
					dec := utf16.DecodeRune(r, nextr)
					if dec != unicode.ReplacementChar { // A valid pair consume
						oute += utf8.EncodeRune(out[oute:], dec)
						e += 6
						break loop
					}
					// Invalid surrogate; fall back to replacement rune.
					r = unicode.ReplacementChar
				}
				oute += utf8.EncodeRune(out[oute:], r)

			} else { // escaped with " \ / ' b f n r t
				out[oute] = escapeCode[txt[e+1]]
				e += 2
				oute++
			}

		case c < ' ': // control character is invalid
			return nil, 0

		case c < utf8.RuneSelf: // ASCII
			out[oute] = c
			oute++
			e++

		default: // coerce to well-formed UTF-8
			r, size := utf8.DecodeRune(txt[e:])
			e += size
			oute += utf8.EncodeRune(out[oute:], r)
		}
	}

	if out[oute] == quote {
		return out[:oute+1], e
	}
	return nil, 0
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	r, err := strconv.ParseUint(string(s[2:6]), 16, 64)
	if err != nil {
		return -1
	}
	return rune(r)
}

func scanSimpleString(txt []byte) (tok []byte, readn int) {
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
