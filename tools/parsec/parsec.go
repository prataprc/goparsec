//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not
//  use this file except in compliance with the License. You may obtain a copy
//  of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package main

import (
    "flag"
    "fmt"
    "github.com/prataprc/goparsec"
    eg "github.com/prataprc/goparsec/examples"
    "io/ioutil"
    "os"
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
