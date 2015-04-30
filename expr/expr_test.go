//  Copyright (c) 2013 Couchbase, Inc.

package expr

import "testing"

import "github.com/prataprc/goparsec"

var exprText = `4 + 123 + 23 + 67 +89 + 87 *78
/67-98-		 199`

func TestExpr(t *testing.T) {
	s := parsec.NewScanner([]byte(exprText))
	v, _ := Y(s)
	if v.(int) != 110 {
		t.Fatalf("Mismatch value %v\n", v)
	}
}

func BenchmarkExpr1Op(b *testing.B) {
	text := []byte(`19 + 10`)
	for i := 0; i < b.N; i++ {
		Y(parsec.NewScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkExpr2Op(b *testing.B) {
	text := []byte(`19+10*20`)
	for i := 0; i < b.N; i++ {
		Y(parsec.NewScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkExpr3Op(b *testing.B) {
	text := []byte(`19 + 10 * 20/9`)
	for i := 0; i < b.N; i++ {
		Y(parsec.NewScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkExpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Y(parsec.NewScanner([]byte(exprText)))
	}
	b.SetBytes(int64(len(exprText)))
}
