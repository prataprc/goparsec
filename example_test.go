package parsec

import "fmt"
import "strconv"

func ExampleAST_And() {
	// parse a configuration line from ini file.
	text := []byte(`loglevel = info`)
	ast := NewAST("example", 100)
	y := ast.And("configline", nil, Ident(), Atom("=", "EQUAL"), Ident())
	root, _ := ast.Parsewith(y, NewScanner(text))
	nodes := root.GetChildren()
	fmt.Println(nodes[0].GetName(), nodes[0].GetValue())
	fmt.Println(nodes[1].GetName(), nodes[1].GetValue())
	fmt.Println(nodes[2].GetName(), nodes[2].GetValue())
	// Output:
	// IDENT loglevel
	// EQUAL =
	// IDENT info
}

func ExampleAST_OrdChoice() {
	// parse a boolean value
	text := []byte(`true`)
	ast := NewAST("example", 100)
	y := ast.OrdChoice("bool", nil, Atom("true", "TRUE"), Atom("false", "FALSE"))
	root, _ := ast.Parsewith(y, NewScanner(text))
	fmt.Println(root.GetName(), root.GetValue())
	// Output:
	// TRUE true
}

func ExampleAST_Many() {
	// parse comma separated values
	text := []byte(`10,30,50 wont parse this`)
	ast := NewAST("example", 100)
	y := ast.Many("many", nil, Int(), Atom(",", "COMMA"))
	root, _ := ast.Parsewith(y, NewScanner(text))
	nodes := root.GetChildren()
	fmt.Println(nodes[0].GetName(), nodes[0].GetValue())
	fmt.Println(nodes[1].GetName(), nodes[1].GetValue())
	fmt.Println(nodes[2].GetName(), nodes[2].GetValue())
	// Output:
	// INT 10
	// INT 30
	// INT 50
}

func ExampleAST_ManyUntil() {
	// make sure to parse the entire text
	text := []byte("10,30,50")
	ast := NewAST("example", 100)
	y := ast.ManyUntil("values", nil, Int(), Atom(",", "COMMA"), ast.End("eof"))
	root, _ := ast.Parsewith(y, NewScanner(text))
	nodes := root.GetChildren()
	fmt.Println(nodes[0].GetName(), nodes[0].GetValue())
	fmt.Println(nodes[1].GetName(), nodes[1].GetValue())
	fmt.Println(nodes[2].GetName(), nodes[2].GetValue())
	// Output:
	// INT 10
	// INT 30
	// INT 50
}

func ExampleAST_Maybe() {
	// parse an optional token
	ast := NewAST("example", 100)
	equal := Atom("=", "EQUAL")
	maybeand := ast.Maybe("maybeand", nil, Atom("&", "AND"))
	y := ast.And("assignment", nil, Ident(), equal, maybeand, Ident())

	text := []byte("a = &b")
	root, _ := ast.Parsewith(y, NewScanner(text))
	nodes := root.GetChildren()
	fmt.Println(nodes[0].GetName(), nodes[0].GetValue())
	fmt.Println(nodes[1].GetName(), nodes[1].GetValue())
	fmt.Println(nodes[2].GetName(), nodes[2].GetValue())
	fmt.Println(nodes[3].GetName(), nodes[3].GetValue())

	text = []byte("a = b")
	ast = ast.Reset()
	root, _ = ast.Parsewith(y, NewScanner(text))
	nodes = root.GetChildren()
	fmt.Println(nodes[0].GetName(), nodes[0].GetValue())
	fmt.Println(nodes[1].GetName(), nodes[1].GetValue())
	fmt.Println(nodes[2].GetName())
	fmt.Println(nodes[3].GetName(), nodes[3].GetValue())
	// Output:
	// IDENT a
	// EQUAL =
	// AND &
	// IDENT b
	// IDENT a
	// EQUAL =
	// missing
	// IDENT b
}

func ExampleASTNodify() {
	text := []byte("10 * 20")
	ast := NewAST("example", 100)
	y := ast.And(
		"multiply",
		func(name string, s Scanner, node Queryable) Queryable {
			cs := node.GetChildren()
			x, _ := strconv.Atoi(cs[0].(*Terminal).GetValue())
			y, _ := strconv.Atoi(cs[2].(*Terminal).GetValue())
			return &Terminal{Value: fmt.Sprintf("%v", x*y)}
		},
		Int(), Token(`\*`, "MULT"), Int(),
	)
	node, _ := ast.Parsewith(y, NewScanner(text))
	fmt.Println(node.GetValue())
	// Output:
	// 200
}

func ExampleAnd() {
	// parse a configuration line from ini file.
	text := []byte(`loglevel = info`)
	y := And(nil, Ident(), Atom("=", "EQUAL"), Ident())
	root, _ := y(NewScanner(text))
	nodes := root.([]ParsecNode)
	t := nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[1].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[2].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	// Output:
	// IDENT loglevel
	// EQUAL =
	// IDENT info
}

func ExampleOrdChoice() {
	// parse a boolean value
	text := []byte(`true`)
	y := OrdChoice(nil, Atom("true", "TRUE"), Atom("false", "FALSE"))
	root, _ := y(NewScanner(text))
	nodes := root.([]ParsecNode)
	t := nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	// Output:
	// TRUE true
}

func ExampleMany() {
	// parse comma separated values
	text := []byte(`10,30,50 wont parse this`)
	y := Many(nil, Int(), Atom(",", "COMMA"))
	root, _ := y(NewScanner(text))
	nodes := root.([]ParsecNode)
	t := nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[1].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[2].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	// Output:
	// INT 10
	// INT 30
	// INT 50
}

func ExampleManyUntil() {
	// make sure to parse the entire text
	text := []byte("10,20,50")
	y := ManyUntil(nil, Int(), Atom(",", "COMMA"), End())
	root, _ := y(NewScanner(text))
	nodes := root.([]ParsecNode)
	t := nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[1].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[2].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	// Output:
	// INT 10
	// INT 20
	// INT 50
}

func ExampleMaybe() {
	// parse an optional token
	equal := Atom("=", "EQUAL")
	maybeand := Maybe(nil, Atom("&", "AND"))
	y := And(nil, Ident(), equal, maybeand, Ident())

	text := []byte("a = &b")
	root, _ := y(NewScanner(text))
	nodes := root.([]ParsecNode)
	t := nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[1].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[2].([]ParsecNode)[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[3].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())

	text = []byte("a = b")
	root, _ = y(NewScanner(text))
	nodes = root.([]ParsecNode)
	t = nodes[0].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	t = nodes[1].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	fmt.Println(nodes[2])
	t = nodes[3].(*Terminal)
	fmt.Println(t.GetName(), t.GetValue())
	// Output:
	// IDENT a
	// EQUAL =
	// AND &
	// IDENT b
	// IDENT a
	// EQUAL =
	// missing
	// IDENT b
}

func ExampleNodify() {
	text := []byte("10 * 20")
	s := NewScanner(text)
	y := And(
		func(nodes []ParsecNode) ParsecNode {
			x, _ := strconv.Atoi(nodes[0].(*Terminal).GetValue())
			y, _ := strconv.Atoi(nodes[2].(*Terminal).GetValue())
			return x * y // this is retuned as node further down.
		},
		Int(), Token(`\*`, "MULT"), Int(),
	)
	node, _ := y(s)
	fmt.Println(node)
	// Output:
	// 200
}
