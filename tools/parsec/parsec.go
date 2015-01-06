//  Copyright (c) 2013 Couchbase, Inc.

package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/goparsec/expr"
import "github.com/prataprc/goparsec/json"

var options struct {
	expr string
	json string
}

func argParse() {
	flag.StringVar(&options.expr, "expr", "",
		"Specify input file or arithmetic expression string")
	flag.StringVar(&options.json, "json", "",
		"Specify input file or json string")
	flag.Parse()
}

func main() {
	argParse()
	if options.expr != "" {
		doExpr(getText(options.expr))
	} else if options.json != "" {
		doJSON(getText(options.json))
	}
}

func doExpr(text string) {
	s := parsec.NewScanner([]byte(text))
	v, _ := expr.Y(s)
	fmt.Println(v)
}

func doJSON(text string) {
	s := json.NewJSONScanner([]byte(text))
	v, _ := json.Y(s)
	fmt.Println(v)
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
