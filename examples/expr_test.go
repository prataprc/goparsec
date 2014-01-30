//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package examples

import (
	"github.com/prataprc/goparsec"
	"testing"
)

var exprText = `4 + 123 + 23 + 67 + 89 +
87 * 78 / 67 - 98 - 199`

func TestExpr(t *testing.T) {
	s := parsec.NewScanner([]byte(exprText))
	v, s := Expr(s)
	if v.(int) != 110 {
		t.Fatalf("Mismatch value %v\n", v)
	}
}

func BenchmarkExpr1Op(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := parsec.NewScanner([]byte(`19 + 10`))
		Expr(s)
	}
}

func BenchmarkExpr2Op(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := parsec.NewScanner([]byte(`19 + 10 * 20`))
		Expr(s)
	}
}

func BenchmarkExpr3Op(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := parsec.NewScanner([]byte(`19 + 10 * 20 / 9`))
		Expr(s)
	}
}

func BenchmarkExpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := parsec.NewScanner([]byte(exprText))
		Expr(s)
	}
}
