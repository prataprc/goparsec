//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// parsec also supplies a basic set of token parsers that can be used to
// create higher order parser using one of the many combinators.

package parsec

import (
    "strings"
)

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
    return Token(`^-?[0-9]+`, "INT")
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
    return Token(`^-?[0-9]*\.[0-9]+`, "FLOAT")
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

// OrdTerminals
func OrdTokens(patterns []string, names []string) Parser {
    groups := make([]string, 0, len(patterns))
    for _, pattern := range patterns {
        groups = append(groups, "(" + pattern + ")")
    }
    ordPattern := strings.Join(groups, "|")
    return func(s Scanner) (ParsecNode, Scanner) {
        news := s.Clone()
        news.SkipWS()
        if toks, _ := news.SubmatchAll(ordPattern); toks != nil {
            for i, tok := range toks[1:] {
                if len(tok) == 0 {
                    continue
                }
                t := Terminal{
                    Name:     names[i],
                    Value:    string(tok),
                    Position: news.GetCursor(),
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

// NoEnd is a parser function to detect not-an-end of scanner output.
func NoEnd(s Scanner) (ParsecNode, Scanner) {
    return !s.Endof(), s
}
