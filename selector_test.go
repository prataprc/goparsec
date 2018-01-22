package parsec

import "os"
import "fmt"
import "bytes"
import "strings"
import "testing"
import "compress/gzip"
import "io/ioutil"
import "path/filepath"

// TODO: interesting selectors
//   tagstart *[class=term]

func TestParseselectorBasic(t *testing.T) {
	// create and reuse.
	selast := NewAST("selectors", 100)
	sely := parseselector(selast)

	// test parsing `*`
	ref := "*"
	qsel, _ := selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs := qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	validateselterm(t, cs[0].GetChildren()[0] /*star*/, "STAR", "*")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `.term`
	selast.Reset()
	ref = ".term"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	validateselterm(t, cs[0].GetChildren()[2] /*shorth*/, "SHORTH", ".term")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `#uniqueid`
	selast.Reset()
	ref = "#uniqueid"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	validateselterm(t, cs[0].GetChildren()[2] /*shorth*/, "SHORTH", "#uniqueid")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `tagstart` node-name
	selast.Reset()
	ref = "tagstart"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	validateselterm(t, cs[0].GetChildren()[1] /*shorth*/, "NAME", "tagstart")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute]`
	selast.Reset()
	ref = "[class]"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq := cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute=value]`
	selast.Reset()
	ref = "[class=term]"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	valq := attrq.GetChildren()[2]
	validateselterm(t, valq.GetChildren()[0], "ATTRSEP", "=")
	validateselterm(t, valq.GetChildren()[1], "attrchoice", "term")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute~=value]`
	selast.Reset()
	ref = "[class~=on]" // containing er
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	valq = attrq.GetChildren()[2]
	validateselterm(t, valq.GetChildren()[0], "ATTRSEP", "~=")
	validateselterm(t, valq.GetChildren()[1], "attrchoice", "on")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute^=value]`
	selast.Reset()
	ref = "[class^=ter]"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	valq = attrq.GetChildren()[2]
	validateselterm(t, valq.GetChildren()[0], "ATTRSEP", "^=")
	validateselterm(t, valq.GetChildren()[1], "attrchoice", "ter")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute$=value]`
	selast.Reset()
	ref = "[class$=term]"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	valq = attrq.GetChildren()[2]
	validateselterm(t, valq.GetChildren()[0], "ATTRSEP", "$=")
	validateselterm(t, valq.GetChildren()[1], "attrchoice", "term")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `[attribute*=value]`
	selast.Reset()
	ref = "[class*=non]"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[3]
	validateselterm(t, attrq.GetChildren()[0], "OPENSQR", "[")
	validateselterm(t, attrq.GetChildren()[1], "ATTRNAME", "class")
	validateselterm(t, attrq.GetChildren()[3], "CLOSESQR", "]")
	valq = attrq.GetChildren()[2]
	validateselterm(t, valq.GetChildren()[0], "ATTRSEP", "*=")
	validateselterm(t, valq.GetChildren()[1], "attrchoice", "non")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:empty`
	selast.Reset()
	ref = ":empty"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "empty", "empty")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:first-child`
	selast.Reset()
	ref = ":first-child"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "first-child", "first-child")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:first-of-type`
	selast.Reset()
	ref = ":first-of-type"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "first-of-type", "first-of-type")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:last-child`
	selast.Reset()
	ref = ":last-child"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "last-child", "last-child")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:last-of-type`
	selast.Reset()
	ref = ":last-of-type"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "last-of-type", "last-of-type")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:nth-child(n)`
	selast.Reset()
	ref = ":nth-child(0)"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	argq := attrq.GetChildren()[1]
	validateselterm(t, argq.GetChildren()[0], "CNC", "nth-child")
	validateselterm(t, argq.GetChildren()[1], "INT", "(0)")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:nth-last-child(n)`
	selast.Reset()
	ref = ":nth-last-child(0)"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	argq = attrq.GetChildren()[1]
	validateselterm(t, argq.GetChildren()[0], "CNLC", "nth-last-child")
	validateselterm(t, argq.GetChildren()[1], "INT", "(0)")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:nth-last-of-type(n)`
	selast.Reset()
	ref = ":nth-last-of-type(0)"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	argq = attrq.GetChildren()[1]
	validateselterm(t, argq.GetChildren()[0], "CNLOT", "nth-last-of-type")
	validateselterm(t, argq.GetChildren()[1], "INT", "(0)")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:nth-of-type(n)`
	selast.Reset()
	ref = ":nth-of-type(0)"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	argq = attrq.GetChildren()[1]
	validateselterm(t, argq.GetChildren()[0], "CNOT", "nth-of-type")
	validateselterm(t, argq.GetChildren()[1], "INT", "(0)")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:only-of-type`
	selast.Reset()
	ref = ":only-of-type"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "only-of-type", "only-of-type")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	// test parsing `:only-child`
	selast.Reset()
	ref = ":only-child"
	qsel, _ = selast.Parsewith(sely, NewScanner([]byte(ref)))
	cs = qsel.GetChildren()
	if len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	} else if cs = cs[0].GetChildren(); len(cs) != 1 {
		t.Errorf("unexpected %v", len(cs))
	}
	attrq = cs[0].GetChildren()[4]
	validateselterm(t, attrq.GetChildren()[0], "COLON", ":")
	validateselterm(t, attrq.GetChildren()[1], "only-child", "only-child")
	if v := qsel.GetValue(); v != ref {
		t.Errorf("expected %q, got %q", ref, v)
	}

	//buf := bytes.NewBuffer(nil)
	//selast.prettyprint(buf, "", qsel)
	//fmt.Println(string(buf.Bytes()))
}

func TestParseselector(t *testing.T) {
	updateref := false
	selast := NewAST("selector", 100)
	sely := parseselector(selast)

	valrefs := map[int]string{
		22: "[class=term]",
		23: "[class=term]",
		24: "tagstarttagendtext",
		25: "tagstarttext",
		26: "tagstarttext",
		27: "oanglebrkttagname",
		28: "tagnamecanglebrkt",
	}

	for i := 1; i < 29; i++ {
		inpfile := filepath.Join(
			"testdata", "selectors", fmt.Sprintf("selector%v.txt", i))
		ppfile := filepath.Join(
			"testdata", "selectors", fmt.Sprintf("selector%v.pprint", i))
		inp := bytes.Trim(testdataFile(inpfile), "\r\n")
		qsel, _ := selast.Parsewith(sely, NewScanner(inp))

		buf := bytes.NewBuffer(nil)
		selast.prettyprint(buf, "", qsel)
		out := buf.Bytes()
		if updateref {
			ioutil.WriteFile(ppfile, out, 0660)
		}

		ref := testdataFile(ppfile)
		if bytes.Compare(out, ref) != 0 {
			t.Errorf("expected %s", string(ref))
			t.Errorf("got %s", string(out))
		}

		valref, ok := valrefs[i]
		if ok == false {
			valref = string(inp)
		}
		if qsel.GetValue() != valref {
			t.Errorf("expected %v, got %v", valref, qsel.GetValue())
		}
	}
}

func TestGetSelectorattr(t *testing.T) {
	// create and reuse.
	selast := NewAST("selectors", 100)
	sely := parseselector(selast)

	// test parsing `*`
	ref := "[class=term]"
	qsel, _ := selast.Parsewith(sely, NewScanner([]byte(ref)))
	q := qsel.GetChildren()[0].GetChildren()[0]
	key, op, match := getselectorattr(q)
	if key != "class" {
		t.Errorf("unexpected %v", key)
	} else if op != "=" {
		t.Errorf("unexpected %v", op)
	} else if match != "term" {
		t.Errorf("unexpected %v", match)
	}
}

func TestGetSelectorcolon(t *testing.T) {
	selast := NewAST("selectors", 100)
	sely := parseselector(selast)
	ref := ":nth-child(0)"
	qsel, _ := selast.Parsewith(sely, NewScanner([]byte(ref)))
	q := qsel.GetChildren()[0].GetChildren()[0]
	colonspec, colonarg := getcolon(q)
	if colonspec != "nth-child" {
		t.Errorf("unexpected %v", colonspec)
	} else if colonarg != "(0)" {
		t.Errorf("unexpected %v", colonarg)
	}
}

func TestFilterbyname(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlroot, _ := htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	node := htmlroot.GetChildren()[0].GetChildren()[0]
	if filterbyname(node, "OT") == false {
		t.Errorf("expected true")
	} else if filterbyname(node, "ot") == false {
		t.Errorf("expected true")
	} else if filterbyname(node, "xyz") == true {
		t.Errorf("exected false")
	}
}

func TestFilterbyattr(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlroot, _ := htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	term := htmlroot.GetChildren()[0].GetChildren()[0]
	if filterbyattr(term, "", "", "") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "value", "", "") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "value", "=", "<") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "value", "=", ">") == true {
		t.Errorf("expected false")
	} else if filterbyattr(term, "missattr", "", "") == true {
		t.Errorf("expected false")
	} else if filterbyattr(term, "class", "", "") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "=", "term") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "~=", "er") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "^=", "te") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "$=", "rm") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "*=", "[term]{4}") == false {
		t.Errorf("expected true")
	} else if filterbyattr(term, "class", "==", "[term]{4}") == true {
		t.Errorf("expected false")
	}
}

func TestFilterbycolon(t *testing.T) {
	html := []byte("<a><b></b><em></em></a>")
	htmlast := NewAST("html", 100)
	htmlroot, _ := htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	nttagstart := htmlroot.GetChildren()[0]
	ntelements := htmlroot.GetChildren()[1]
	termot := nttagstart.GetChildren()[0]
	termtaga := nttagstart.GetChildren()[1]

	if len(ntelements.GetChildren()) != 2 {
		t.Errorf("unexpected %v", len(ntelements.GetChildren()))
	}

	if filterbycolon(nttagstart, 0, termot, "empty", "") == false {
		t.Errorf("expected true")
	} else if filterbycolon(htmlroot, 0, nttagstart, "empty", "") == true {
		t.Errorf("expected false")
	} else if filterbycolon(nttagstart, 0, termot, "first-child", "") == false {
		t.Errorf("expected true")
	} else if filterbycolon(nttagstart, 1, termtaga, "first-child", "") == true {
		t.Errorf("expected false")
	}

	nttag1 := ntelements.GetChildren()[0]
	nttag2 := ntelements.GetChildren()[1]

	if filterbycolon(nil, 0, htmlroot, "first-of-type", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "first-of-type", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 1, nttag2, "first-of-type", "") == true {
		t.Errorf("expected false")
	}

	if filterbycolon(nil, 0, htmlroot, "last-child", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "last-child", "") == true {
		t.Errorf("expected false")
	}
	if filterbycolon(ntelements, 1, nttag2, "last-child", "") == false {
		t.Errorf("expected true")
	}

	if filterbycolon(nil, 0, htmlroot, "last-of-type", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "last-of-type", "") == true {
		t.Errorf("expected false")
	}
	if filterbycolon(ntelements, 1, nttag2, "last-of-type", "") == false {
		t.Errorf("expected true")
	}

	if filterbycolon(nil, 0, htmlroot, "nth-child", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "nth-child", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 1, nttag2, "nth-child", "(0)") == true {
		t.Errorf("expected false")
	}

	if filterbycolon(nil, 0, htmlroot, "nth-of-type", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "nth-of-type", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 1, nttag2, "nth-of-type", "(0)") == true {
		t.Errorf("expected false")
	}

	if filterbycolon(nil, 0, htmlroot, "nth-last-child", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "nth-last-child", "(0)") == true {
		t.Errorf("expected false")
	}
	if filterbycolon(ntelements, 1, nttag2, "nth-last-child", "(0)") == false {
		t.Errorf("expected true")
	}

	if filterbycolon(nil, 0, htmlroot, "nth-last-of-type", "(0)") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "nth-last-of-type", "(0)") == true {
		t.Errorf("expected false")
	}
	if filterbycolon(ntelements, 1, nttag2, "nth-last-of-type", "(0)") == false {
		t.Errorf("expected true")
	}

	if filterbycolon(nil, 0, htmlroot, "only-of-type", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(nttagstart, 0, termot, "only-of-type", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(ntelements, 0, nttag1, "only-of-type", "") == true {
		t.Errorf("expected false")
	}

	if filterbycolon(nil, 0, htmlroot, "only-child", "") == false {
		t.Errorf("expected true")
	}
	if filterbycolon(nttagstart, 0, termot, "only-child", "") == true {
		t.Errorf("expected false")
	}
}

func TestAstwalk1(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "*"
	ch := make(chan Queryable, 1000)
	htmlast.Query("*", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 12 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{
		"tag <a></a>",
		"tagstart <a>", "OT <", "TAG a", "attributes", "CT >",
		"elements",
		"tagend </a>", "OT <", "SLASH /", "TAG a", "CT >",
	}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk2(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "tag OT"
	ch := make(chan Queryable, 1000)
	htmlast.Query("tag OT", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 2 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"OT <", "OT <"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk3(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "tagstart > CT"
	ch := make(chan Queryable, 1000)
	htmlast.Query("tagstart > CT", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"CT >"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk4(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "OT + TAG"
	ch := make(chan Queryable, 1000)
	htmlast.Query("OT + TAG", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"TAG a"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk5(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "OT ~ TAG"
	ch := make(chan Queryable, 1000)
	htmlast.Query("OT ~ TAG", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 2 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"TAG a", "TAG a"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk6(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "tagstart ~ tagend TAG"
	ch := make(chan Queryable, 1000)
	htmlast.Query("tagstart ~ tagend TAG", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"TAG a"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk7(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with complexpatterns
	ch := make(chan Queryable, 1000)
	htmlast.Query("tagstart + elements + tagend TAG + CT", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"CT >", "CT >"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk8(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test case-insensitive
	ch := make(chan Queryable, 1000)
	htmlast.Query("tag", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 3 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"tag <a></a>", "TAG a", "TAG a"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk9(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test attributes
	ch := make(chan Queryable, 1000)
	htmlast.Query("tag[class=term]", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 2 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"TAG a", "TAG a"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk10(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test colon names
	ch := make(chan Queryable, 1000)
	htmlast.Query("tag > tagstart:first-child", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
	refs := []string{"tagstart <a>"}
	for i, item := range items {
		out := fmt.Sprintln(item.GetName(), item.GetValue())
		out = strings.TrimRight(out, " \n\r")
		if out != refs[i] {
			t.Errorf("expected %q, got %q", refs[i], out)
		}
	}
}

func TestAstwalk11(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "tag[class=term]:nth-child(1)"
	ch := make(chan Queryable, 1000)
	htmlast.Query("tag[class=term]:nth-child(1)", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 1 {
		t.Errorf("unexpected %v", len(items))
	}
}

func TestAstwalkNeg(t *testing.T) {
	html := []byte("<a></a>")
	htmlast := NewAST("html", 100)
	htmlast.Parsewith(makeexacthtmly(htmlast), NewScanner(html))

	// test with "tagend ~ tagstart"
	ch := make(chan Queryable, 1000)
	htmlast.Query("tagend ~ tagstart", ch)
	items := []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 0 {
		t.Errorf("unexpected %v", len(items))
	}

	// test with "tagstart ~ nothing"
	ch = make(chan Queryable, 1000)
	htmlast.Query("tagstart ~ nothing", ch)
	items = []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 0 {
		t.Errorf("unexpected %v", len(items))
	}

	// test with "tag[class=term]:last-child"
	ch = make(chan Queryable, 1000)
	htmlast.Query("tag[class=term]:last-child", ch)
	items = []Queryable{}
	for item := range ch {
		items = append(items, item)
	}
	if len(items) != 0 {
		t.Errorf("unexpected %v", len(items))
	}
}

func validateselterm(t *testing.T, term Queryable, name, value string) {
	n, ist, v := term.GetName(), term.IsTerminal(), term.GetValue()
	if n != name {
		t.Errorf("expected %v, got %v", name, n)
		panic("")
	} else if ist != true {
		t.Errorf("expected %v, got %v", true, ist)
		panic("")
	} else if value != v {
		t.Errorf("expected %v, got %v", value, v)
		panic("")
	}
}

func testdataFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data []byte
	if strings.HasSuffix(filename, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(gz)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
	}
	return data
}

func makehtmlyForSelector(ast *AST) Parser {
	var tag Parser

	opentag := AtomExact("<", "OT")
	closetag := AtomExact(">", "CT")
	equal := AtomExact("=", "EQUAL")
	slash := TokenExact("/[ \t]*", "SLASH")
	tagname := TokenExact("[a-z][a-zA-Z0-9]*", "TAG")
	attrkey := TokenExact("[a-z][a-zA-Z0-9]*", "ATTRK")
	text := TokenExact("[^<>]+", "TEXT")
	ws := TokenExact("[ \t]+", "WS")

	element := ast.OrdChoice("element", nil, text, &tag)
	elements := ast.Kleene("elements", nil, element)
	attr := ast.And("attribute", nil, attrkey, equal, String())
	attrws := ast.And("attrws", nil, attr, ast.Maybe("ws", nil, ws))
	attrs := ast.Kleene("attributes", nil, attrws)
	tstart := ast.And("tagstart", nil, opentag, tagname, attrs, closetag)
	tend := ast.And("tagend", nil, opentag, slash, tagname, closetag)
	tag = ast.And("tag", nil, tstart, elements, tend)
	return tag
}
