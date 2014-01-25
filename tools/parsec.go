package main

import (
	"fmt"
	"github.com/prataprc/goparsec"
	"testing"
)

var options struct {
	expr string
	json string
}

func arguments() {
	flag.StringVar(&options.expr, "expr", "",
		"Specify input file or arithmetic expression string")
	flag.StringVar(&options.expr, "json", "",
		"Specify input file or json string")
	flag.Parse()
}

func main() {
	if options.expr != "" {
		parseExpr(getText(options.expr))
	}
}

func parseExpr(text string) parsec.ParsecNode {
	s := parsec.NewScanner(text)
	fmt.Println(expr(s))
}

func getText(s string) string {
	if _, err := os.Stat; err != nil {
		return s
	}
	if b, err := ioutil.ReadFile(s); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
