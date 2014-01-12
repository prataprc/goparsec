package parsec

import (
	"fmt"
	"github.com/prataprc/golib"
	"io/ioutil"
	"os"
	"testing"
	"text/scanner"
)

var _ = fmt.Sprintf("Dummy")
var testfile = "./sampletest"

func BenchmarkPrepare(b *testing.B) {
	data, err := ioutil.ReadFile("./parsec_test.go")
	fd, err := os.Create(testfile)
	defer func() {
		fd.Close()
	}()
	if err == nil {
		for i := 0; i < 10000; i++ {
			fd.Write(data)
		}
	}
}

func BenchmarkScanner(b *testing.B) {
	var s scanner.Scanner
	fd, _ := os.Open(testfile)
	s.Init(fd)
	for i := 0; i < b.N; i++ {
		tok := Token{
			Type:  scanner.TokenString(s.Scan()),
			Value: s.TokenText(),
			Pos:   s.Pos(),
		}
		if tok.Type == "EOF" {
			break
		}
	}
}

func BenchmarkGoscan(b *testing.B) {
	config := make(golib.Config)
	text, _ := ioutil.ReadFile(testfile)
	scanner := NewGoScan(text, config)
	for i := 0; i < b.N; i++ {
		tok := scanner.Scan()
		if tok.Type == "EOF" {
			break
		}
	}
}

func BenchmarkUnprepare(b *testing.B) {
	os.Remove(testfile)
}
