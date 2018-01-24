package parsec

import "fmt"
import "bytes"
import "testing"
import "net/http"
import "io/ioutil"

func TestHTMLValue(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/simple.html")
	if err != nil {
		t.Error(err)
	}
	data = bytes.Trim(data, " \t\r\n")

	ast := NewAST("html", 100)
	y := makehtmly(ast)
	s := NewScanner(data).TrackLineno()
	root, _ := ast.Parsewith(y, s)

	ref := `<html><body><h1>My First Heading</h1><p>My first paragraph.</p></body></html>`
	if out := string(root.GetValue()); out != ref {
		t.Errorf("expected %q", ref)
		t.Errorf("got %q", out)
	}

	// To generate the dot-graph for input html.
	//graph := ast.Dotstring("simplehtml")
	//fmt.Println(graph)

	// To gather all TEXT values.
	//ch := make(chan Queryable, 100)
	//ast.Query("TEXT", ch)
	//for node := range ch {
	//	fmt.Println(node.GetValue())
	//}

	// To gather all terminal values.
	ch := make(chan Queryable, 100)
	ast.Query(".term", ch)
	for node := range ch {
		fmt.Printf("%s", node.GetValue())
	}
	fmt.Println()
}

func TestExample(t *testing.T) {
	ast := NewAST("html", 100)
	y := makehtmly(ast)
	resp, err := http.Get("https://example.com/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))

	s := NewScanner(data).TrackLineno()
	ast.Parsewith(y, s)

	ch := make(chan Queryable, 100)
	go ast.Query("attrunquoted,attrsingleq,attrdoubleq", ch)
	for node := range ch {
		cs := node.GetChildren()
		if cs[0].GetValue() != "href" {
			continue
		}
		if len(cs) == 3 {
			fmt.Println(cs[2].GetValue())
		} else {
			fmt.Println(cs[3].GetValue())
		}
	}
	fmt.Println()
}

func makehtmly(ast *AST) Parser {
	var tag Parser

	// terminal parsers.
	tagobrk := Atom("<", "OT")
	tagcbrk := Atom(">", "CT")
	tagcend := Atom("/>", "CT")
	tagcopen := Atom("</", "CT")
	equal := Atom(`=`, "EQ")
	single := Atom("'", "SQUOTE")
	double := Atom(`"`, "DQUOTE")
	tagname := Token(`[a-zA-Z0-9]+`, "TAGNAME")
	attrname := Token(`[a-zA-Z0-9_-]+`, "ATTRNAME")
	attrval1 := Token(`[^\s"'=<>`+"`]+", "ATTRVAL1")
	attrval2 := Token(`[^']*`, "ATTRVAL2")
	attrval3 := Token(`[^"]*`, "ATTRVAL3")
	entity := Token(`&#?[a-bA-Z0-9]+;`, "ENTITY")
	text := Token(`[^<]+`, "TEXT")
	doctype := Token(`<!doctype[^>]+>`, "DOCTYPE")

	// non-terminals
	attrunquoted := ast.And(
		"attrunquoted", nil, attrname, equal, attrval1,
	)
	attrsingleq := ast.And(
		"attrsingleq", nil, attrname, equal, single, attrval2, single,
	)
	attrdoubleq := ast.And(
		"attrdoubleq", nil, attrname, equal, double, attrval3, double,
	)
	attr := ast.OrdChoice(
		"attribute", nil, attrsingleq, attrdoubleq, attrunquoted, attrname,
	)
	attrs := ast.Kleene("attributes", nil, attr, nil)

	tagopen := ast.And("tagopen", nil, tagobrk, tagname, attrs, tagcbrk)
	tagclose := ast.And("tagclose", nil, tagcopen, tagname, tagcbrk)

	content := ast.OrdChoice("content", nil, entity, text, &tag)
	contents := ast.Maybe(
		"maybecontents", nil, ast.Kleene("contents", nil, content, nil),
	)

	tagempty := ast.And("tagempty", nil, tagobrk, tagname, attrs, tagcend)
	tagproper := ast.And("tagproper", nil, tagopen, contents, tagclose)
	tag = ast.OrdChoice("tag", nil, doctype, tagempty, tagproper)
	return ast.Kleene("html", nil, tag, nil)
}

func debugfn(name string, s Scanner, node Queryable) Queryable {
	attrs := node.GetChildren()[2]
	fmt.Printf("%T %v\n", attrs, attrs)
	return node
}
