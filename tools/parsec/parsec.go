//  Copyright (c) 2013 Couchbase, Inc.

package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"

import "github.com/prataprc/goparsec"
import eg "github.com/prataprc/goparsec/examples"

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
	var n parsec.ParsecNode

	if options.expr != "" {
		n = parseExpr(getText(options.expr))
	} else if options.json != "" {
		n = parseExpr(getText(options.expr))
	}
	fmt.Println(n)
}

func parseExpr(text string) parsec.ParsecNode {
	s := parsec.NewScanner([]byte(text))
	n, _ := eg.Expr(s)
	return n
}

func parseJSON(text string) parsec.ParsecNode {
	n := eg.JSONParse([]byte(text))
	return n
}

func getText(filename string) string {
	if _, err := os.Stat(filename); err != nil {
		return filename
	}
	if b, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
