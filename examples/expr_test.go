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
