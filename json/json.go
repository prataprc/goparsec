// Copyright (c) 2013 Goparsec AUTHORS. All rights reserved.
// Use of this source code is governed by LICENSE file.

// Package json provide a parser to parse JSON string.
package json

import "strconv"
import "unicode"
import "fmt"
import "unicode/utf8"
import "unicode/utf16"

import "github.com/prataprc/goparsec"

var _ = fmt.Sprintf("dummy print")

// Null is alias for string type denoting JSON `null`
type Null string

// True is alias for string type denoting JSON `true`
type True string

// False is alias for string type denoting JSON `null`
type False string

// Num is alias for string type denoting JSON `null`
type Num string

// String is alias for string type denoting JSON `null`
type String string

// Y is root Parser, usually called as `s` in CFG theory.
var Y parsec.Parser
var value parsec.Parser // circular rats

// NonTerminal rats
// values -> value | values "," value
var values = parsec.Kleene(valuesNode, &value, comma())

// array -> "[" values "]"
var array = parsec.And(arrayNode, openSqrt(), values, closeSqrt())

// property -> string ":" value
var property = parsec.And(many2many, sTring(), colon(), &value)

// properties -> property | properties "," property
var properties = parsec.Kleene(propertiesNode, property, comma())

// object -> "{" properties "}"
var object = parsec.And(objectNode, openBrace(), properties, closeBrace())

func init() {
	// value -> null | true | false | num | string | array | object
	value = parsec.OrdChoice(valueNode, parsec.Parser(tokenTerm), array, object)
	// expr  -> sum
	Y = parsec.OrdChoice(one2one, value)
}

//----------
// Nodifiers
//----------

func valueNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		switch n.Name {
		case "NULL":
			return Null("null")
		case "TRUE":
			return True("true")
		case "FALSE":
			return False("false")
		case "NUM":
			return Num(n.Value)
		case "STRING":
			return String(n.Value)
		}

	case []parsec.ParsecNode:
		return n

	case map[string]interface{}:
		return n
	}
	return nil
}

func valuesNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	return ns
}

func arrayNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	return ns[1]
}

func propertiesNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns != nil && len(ns) > 0 {
		m := make(map[string]interface{})
		for _, n := range ns {
			prop := n.([]parsec.ParsecNode)
			key := prop[0].(*parsec.Terminal)
			m[key.Value] = prop[2]
		}
		return m
	}
	return nil
}

func objectNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	return ns[1]
}

//--------
// Scanner
//--------

var nullTerminal = &parsec.Terminal{Name: "NULL", Value: "null"}
var trueTerminal = &parsec.Terminal{Name: "TRUE", Value: "true"}
var falseTerminal = &parsec.Terminal{Name: "FALSE", Value: "false"}

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

var spaceCode = [256]byte{ // TODO: size can be optimized
	'\t': 1,
	'\n': 1,
	'\v': 1,
	'\f': 1,
	'\r': 1,
	' ':  1,
}

var intCheck = [256]byte{}
var digitCheck = [256]byte{}

func init() {
	for i := 48; i <= 57; i++ {
		intCheck[i] = 1
	}
	intCheck['-'] = 1
	intCheck['+'] = 1
	intCheck['.'] = 1
	intCheck['e'] = 1
	intCheck['E'] = 1

	for i := 48; i <= 57; i++ {
		digitCheck[i] = 1
	}
	digitCheck['-'] = 1
	digitCheck['+'] = 1
	digitCheck['.'] = 1
}

// JSONScanner implements parsec.Scanner{} interface used
// as custom scanner for parsing JSON string.
type JSONScanner struct {
	buf    []byte // input buffer
	cursor int    // cursor within input buffer
}

// NewJSONScanner return a new Scanner{} interface for parsing
// JSON string.
func NewJSONScanner(text []byte) *JSONScanner {
	return &JSONScanner{
		buf:    text,
		cursor: 0,
	}
}

// SetWSPattern method receiver in Scanner interface.
func (s *JSONScanner) SetWSPattern(pattern string) parsec.Scanner {
	return s
}

// TrackLineno method receiver in Scanner interface.
func (s *JSONScanner) TrackLineno() parsec.Scanner {
	return s
}

// Clone method receiver in Scanner interface.
func (s *JSONScanner) Clone() parsec.Scanner {
	return &JSONScanner{
		buf:    s.buf,
		cursor: s.cursor,
	}
}

// GetCursor method receiver in Scanner interface.
func (s *JSONScanner) GetCursor() int {
	return s.cursor
}

// Match method receiver in Scanner interface.
func (s *JSONScanner) Match(pattern string) ([]byte, parsec.Scanner) {
	return nil, nil
}

// MatchString method receiver in Scanner interface.
func (s *JSONScanner) MatchString(pattern string) (bool, parsec.Scanner) {
	return false, nil
}

// SubmatchAll method receiver in Scanner interface.
func (s *JSONScanner) SubmatchAll(
	pattern string) (map[string][]byte, parsec.Scanner) {

	return nil, nil
}

// SkipWS method receiver in Scanner interface.
func (s *JSONScanner) SkipWS() ([]byte, parsec.Scanner) {
	return nil, nil
}

// SkipAny method receiver in Scanner interface.
func (s *JSONScanner) SkipAny(pattern string) ([]byte, parsec.Scanner) {
	return nil, nil
}

// Lineno method receiver in Scanner interface.
func (s *JSONScanner) Lineno() int {
	return 0
}

// Endof method receiver in Scanner interface.
func (s *JSONScanner) Endof() bool {
	if s.cursor >= len(s.buf) {
		return true
	}
	return false
}

func colon() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("COLON", ':', s)
	}
}

func comma() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("COMMA", ',', s)
	}
}

func openSqrt() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("OPENSQR", '[', s)
	}
}

func closeSqrt() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("CLOSESQR", ']', s)
	}
}

func openBrace() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("OPENBRACE", '{', s)
	}
}

func closeBrace() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		return matchChar("CLOSEBRACE", '}', s)
	}
}

func sTring() parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		sp := s.(*JSONScanner)
		txt := sp.buf[sp.cursor:]
		// scan for whitespace
		_, l := scanWS(txt)
		sp.cursor, txt = sp.cursor+l, txt[l:]
		if len(txt) < 1 {
			return nil, sp
		}
		// scan for string
		if txt[0] == '"' {
			tok, ln := scanString(txt)
			if tok == nil {
				return nil, sp
			}
			t := &parsec.Terminal{
				Name:     "STRING",
				Value:    string(tok[1 : len(tok)-1]),
				Position: sp.cursor,
			}
			sp.cursor += ln
			return t, sp
		}
		return nil, sp
	}
}

func tokenTerm(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	sp := s.(*JSONScanner)
	txt := sp.buf[sp.cursor:]
	_, l := scanWS(txt)
	sp.cursor, txt = sp.cursor+l, txt[l:]
	if len(txt) < 1 {
		return nil, sp
	}

	if digitCheck[txt[0]] == 1 {
		t := scanNum(txt, sp.cursor)
		sp.cursor += len(t.Value)
		return t, sp
	}

	switch txt[0] {
	case 'n':
		if txt[1] == 'u' && txt[2] == 'l' && txt[3] == 'l' {
			t := *nullTerminal
			t.Position = sp.cursor
			sp.cursor += 4
			return &t, sp
		}
		return nil, sp

	case 't':
		if txt[1] == 'r' && txt[2] == 'u' && txt[3] == 'e' {
			t := *trueTerminal
			t.Position = sp.cursor
			sp.cursor += 4
			return &t, sp
		}
		return nil, sp

	case 'f':
		if txt[1] == 'a' && txt[2] == 'l' && txt[3] == 's' && txt[4] == 'e' {
			t := *falseTerminal
			t.Position = sp.cursor
			sp.cursor += 5
			return &t, sp
		}

	case '-':
		t := scanNum(txt, sp.cursor)
		sp.cursor += len(t.Value)
		return t, sp

	case '"':
		tok, ln := scanString(txt)
		if tok == nil {
			return nil, sp
		}
		t := &parsec.Terminal{
			Name:     "STRING",
			Value:    string(tok[1 : len(tok)-1]),
			Position: sp.cursor,
		}
		sp.cursor += ln
		return t, sp
	}
	return nil, sp
}

func scanNum(txt []byte, cursor int) *parsec.Terminal {
	e, l := 1, len(txt)
	if len(txt) > 1 {
		for ; e < l && intCheck[txt[e]] == 1; e++ {
		}
	}
	return &parsec.Terminal{
		Name:     "NUM",
		Value:    string(txt[:e]),
		Position: cursor,
	}
}

func scanString(txt []byte) ([]byte, int) {
	if len(txt) < 2 {
		return nil, 0
	}

	e := 1
	for txt[e] != '"' {
		c := txt[e]
		if c == '\\' || c == '"' || c < ' ' {
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

	if txt[e] == '"' { // done we have nothing to unquote
		return txt[:e+1], e + 1
	}

	out := make([]byte, len(txt)+2*utf8.UTFMax)
	oute := copy(out, txt[:e]) // copy so far

loop:
	for e < len(txt) {
		switch c := txt[e]; {
		case c == '"':
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

	if out[oute] == '"' {
		return out[:oute+1], e
	}
	return nil, 0
}

func scanWS(txt []byte) ([]byte, int) {
	for i, c := range txt {
		if spaceCode[c] != 1 { // if !unicode.IsSpace(run) {
			return txt[:i], i
		}
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

func nativeValue(m interface{}) interface{} {
	switch v := m.(type) {
	case Null:
		return nil

	case True:
		return true

	case False:
		return false

	case Num:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return nil
		}
		return f

	case String:
		return string(v)

	case []parsec.ParsecNode:
		nv := make([]interface{}, len(v))
		for i, n := range v {
			nv[i] = nativeValue(n)
		}
		return nv

	case map[string]interface{}:
		for key, value := range v {
			v[key] = nativeValue(value)
		}
		return v
	}
	return nil
}

func matchChar(
	name string,
	ch byte,
	s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {

	sp := s.(*JSONScanner)
	i, ln := sp.cursor, len(sp.buf)
	for ; i < ln; i++ {
		if spaceCode[sp.buf[i]] != 1 { // if !unicode.IsSpace(run) {
			break
		}
	}
	sp.cursor = i
	if i < ln && sp.buf[i] == ch {
		t := &parsec.Terminal{Name: name, Value: string(ch), Position: i}
		sp.cursor++
		return t, sp
	}
	return nil, sp
}
