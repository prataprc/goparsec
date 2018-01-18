package parsec

import "fmt"
import "strings"
import "regexp"
import "strconv"

// attributes on selector non-terminal:
// `op`, `name`, `attrkey`, `attrop`, `attrval`, `colonspec`, `colonarg`

func parseselector(ast *AST) Parser {
	star := AtomExact(`*`, "STAR")
	nodename := TokenExact(`(?i)[a-z][a-z0-9_-]*`, "NAME")
	shorthand := TokenExact(`(?i)[\.#][a-z][a-z0-9_-]*`, "SHORTH")
	selop := TokenExact(`[\s>\+~]+`, "OP")
	comma := TokenExact(`,[ ]*`, "COMMA")

	selector := ast.And("selector", makeselector1,
		ast.Maybe("maybestar", nil, star),
		ast.Maybe("maybenodename", nil, nodename),
		ast.Maybe("maybeshorthand", nil, shorthand),
		ast.Maybe("maybeattribute", nil, yattr(ast)),
		ast.Maybe("maybecolonsel", nil, ycolon(ast)),
	)
	selector2 := ast.And("selector2", makeselector2, selop, selector)
	selectors2 := ast.Kleene("selectors2", nil, selector2)
	yone := ast.And("selectors",
		func(_ string, s Scanner, q Queryable) Queryable {
			children := q.GetChildren()
			qs := []Queryable{children[0]}
			qs = append(qs, children[1].GetChildren()...)
			nt := NewNonTerminal("selectors")
			nt.Children = qs
			return nt
		}, selector, selectors2,
	)
	yor := ast.Kleene("orselectors", nil, yone, comma)
	return yor
}

func yattr(ast *AST) Parser {
	opensqr := AtomExact(`[`, "OPENSQR")
	closesqr := AtomExact(`]`, "CLOSESQR")
	attrname := TokenExact(`(?i)[a-z][a-z0-9_-]*`, "ATTRNAME")

	separator := TokenExact(`=|~=|\^=|\$=|\*=`, "ATTRSEP")

	attrval1 := TokenExact(`"[^"]*"`, "VAL1")
	attrval2 := TokenExact(`'[^']*'`, "VAL2")
	attrval3 := TokenExact(`[^\s\]]+`, "VAL3")
	attrchoice := ast.OrdChoice("attrchoice",
		func(_ string, s Scanner, q Queryable) Queryable {
			var value string
			switch t := q.(*Terminal); t.Name {
			case "VAL1":
				value = t.Value[1 : len(t.Value)-1]
			case "VAL2":
				value = t.Value[1 : len(t.Value)-1]
			case "VAL3":
				value = t.Value
			}
			return NewTerminal("attrchoice", value, s.GetCursor())
		}, attrval1, attrval2, attrval3,
	)
	attrval := ast.And("attrval", nil, separator, attrchoice)
	attribute := ast.And("attribute", nil,
		opensqr, attrname, ast.Maybe("attrval", nil, attrval), closesqr,
	)
	return attribute
}

func ycolon(ast *AST) Parser {
	openparan := AtomExact(`(`, "OPENPARAN")
	closeparan := AtomExact(`)`, "CLOSEPARAN")
	colon := AtomExact(`:`, "COLON")

	arg := ast.And("arg",
		func(_ string, s Scanner, n Queryable) Queryable {
			t := n.GetChildren()[1].(*Terminal)
			t.Value = "(" + t.Value + ")"
			return t
		},
		openparan, Int(), closeparan,
	)
	colonEmpty := AtomExact(`empty`, "empty")
	colonFirstChild := AtomExact(`first-child`, "first-child")
	colonFirstType := AtomExact(`first-of-type`, "first-of-type")
	colonLastChild := AtomExact(`last-child`, "last-child")
	colonLastType := AtomExact(`last-of-type`, "last-of-type")
	colonNthChild := ast.And(`nth-child`,
		func(_ string, s Scanner, n Queryable) Queryable {
			n.SetAttribute("arg", n.GetChildren()[1].GetValue())
			return n
		},
		AtomExact(`nth-child`, "CNC"), arg,
	)
	colonNthType := ast.And(`nth-of-type`,
		func(_ string, s Scanner, n Queryable) Queryable {
			n.SetAttribute("arg", n.GetChildren()[1].GetValue())
			return n
		},
		AtomExact(`nth-of-type`, "CNOT"), arg,
	)
	colonNthLastChild := ast.And(`nth-last-child`,
		func(_ string, s Scanner, n Queryable) Queryable {
			n.SetAttribute("arg", n.GetChildren()[1].GetValue())
			return n
		},
		AtomExact(`nth-last-child`, "CNLC"), arg,
	)
	colonNthLastType := ast.And(`nth-last-of-type`,
		func(_ string, s Scanner, n Queryable) Queryable {
			n.SetAttribute("arg", n.GetChildren()[1].GetValue())
			return n
		},
		AtomExact(`nth-last-of-type`, "CNLOT"), arg,
	)
	colonOnlyType := AtomExact(`only-of-type`, "only-of-type")
	colonOnlyChild := AtomExact(`only-child`, "only-child")
	colonname := ast.OrdChoice("colonname", nil,
		colonEmpty,
		colonFirstChild,
		colonFirstType,
		colonLastChild,
		colonLastType,
		colonNthChild,
		colonNthType,
		colonNthLastChild,
		colonNthLastType,
		colonOnlyType,
		colonOnlyChild,
	)
	return ast.And("selectcolon", nil, colon, colonname)
}

func makeselector2(name string, s Scanner, nt Queryable) Queryable {
	cs := nt.GetChildren()
	cs[1].SetAttribute("op", strings.Trim(cs[0].GetValue(), " "))
	return cs[1]
}

func makeselector1(name string, s Scanner, nt Queryable) Queryable {
	cs := nt.GetChildren()
	nt, ok := makeselector1Star(nt, cs[0]) // maybestar
	if ok == false {
		nt, _ = makeselector1Name(nt, cs[1]) // maybenodename
	}
	nt, _ = makeselector1Shand(nt, cs[2]) // maybeshorthand
	nt, _ = makeselector1Attr(nt, cs[3])  // maybeattribute
	nt, _ = makeselector1Colon(nt, cs[4]) // maybecolonsel
	return nt
}

func makeselector1Star(nt, starq Queryable) (Queryable, bool) {
	if _, ok := starq.(MaybeNone); ok == true {
		return nt, false
	}
	nt.SetAttribute("name", starq.GetValue())
	return nt, true
}

func makeselector1Name(nt, nameq Queryable) (Queryable, bool) {
	if _, ok := nameq.(MaybeNone); ok == true {
		return nt, false
	}
	nt.SetAttribute("name", nameq.GetValue())
	return nt, true
}

func makeselector1Shand(nt, shq Queryable) (Queryable, bool) {
	if _, ok := shq.(MaybeNone); ok == true {
		return nt, false
	}
	value := shq.GetValue()
	switch value[0] {
	case '.':
		nt.SetAttribute("attrkey", "class")
	case '#':
		nt.SetAttribute("attrkey", "id")
	}
	nt.SetAttribute("attrop", "=").SetAttribute("attrval", value[1:])
	return nt, true
}

func makeselector1Attr(nt, attrq Queryable) (Queryable, bool) {
	if _, ok := attrq.(MaybeNone); ok == true {
		return nt, false
	}
	cs := attrq.GetChildren()
	nt.SetAttribute("attrkey", cs[1].GetValue())
	if valcs := cs[2].GetChildren(); len(valcs) > 0 {
		nt.SetAttribute("attrop", valcs[0].GetValue())
		nt.SetAttribute("attrval", valcs[1].GetValue())
	}
	return nt, true
}

func makeselector1Colon(nt, colonq Queryable) (Queryable, bool) {
	if _, ok := colonq.(MaybeNone); ok == true {
		return nt, false
	}
	colonq = colonq.GetChildren()[1]
	nt.SetAttribute("colonspec", colonq.GetName())
	attrvals := colonq.GetAttribute("arg")
	if len(attrvals) > 0 {
		nt.SetAttribute("colonarg", attrvals[0])
	}
	return nt, true
}

//---- walk the tree

func astwalk(
	parent Queryable, idx int, node Queryable,
	qs []Queryable, ch chan Queryable) {

	var descendok bool

	q, remqs := qs[0], qs[1:]
	matchok := applyselector(parent, idx, node, q)
	if qcombinator(q) == "child" && matchok == false {
		return // cannot descend
	} else if matchok == false {
		remqs = qs
	} else if len(remqs) == 0 { // and matchok == true
		remqs = qs
		ch <- node
	} else { // matchok `and` remqs > 0
		idx, remqs, descendok = applysiblings(parent, idx, node, remqs, ch)
		if descendok == false {
			return
		}
	}
	if parent != nil {
		node = parent.GetChildren()[idx]
	}
	for idx, child := range node.GetChildren() {
		astwalk(node, idx, child, remqs, ch)
	}
}

func applysiblings(
	parent Queryable, idx int, node Queryable,
	qs []Queryable, ch chan Queryable) (int, []Queryable, bool) {

	if parent == nil { // node must be root
		return idx, qs, true
	}

	q, typ, children := qs[0], qcombinator(qs[0]), parent.GetChildren()
	nchild, dosibling := len(children), (typ == "next") || (typ == "after")
	if dosibling == false {
		return idx, qs, true
	} else if idx >= (nchild - 1) { // there should be atleast one more sibling
		return idx, qs, false
	}

	if typ == "next" {
		matchok := applyselector(parent, idx+1, children[idx+1], q)
		remqs, node := qs[1:], children[idx+1]
		if matchok && len(remqs) == 0 {
			ch <- children[idx+1]
			return idx, qs, false
		} else if matchok {
			return applysiblings(parent, idx+1, node, remqs, ch)
		}
		return idx, qs, false
	}

	// typ == "after"
	for idx = idx + 1; idx < nchild; idx++ {
		matchok := applyselector(parent, idx, children[idx], q)
		remqs, node := qs[1:], children[idx]
		if matchok && len(remqs) == 0 {
			ch <- children[idx]
			return idx, qs, false
		} else if matchok {
			return applysiblings(parent, idx, node, remqs, ch)
		}
	}
	return idx, qs, false
}

func applyselector(parent Queryable, idx int, node, q Queryable) bool {
	ns := q.GetAttribute("name")
	if len(ns) > 0 {
		if filterbyname(node, ns[0]) == false {
			return false
		}
	}
	key, op, match := getselectorattr(q)
	if filterbyattr(node, key, op, match) == false {
		return false
	}
	colonspec, colonarg := getcolon(q)
	if filterbycolon(parent, idx, node, colonspec, colonarg) == false {
		return false
	}
	return true
}

func qcombinator(sel Queryable) string {
	ops := sel.GetAttribute("op")
	if len(ops) == 0 {
		return ""
	}
	switch op := ops[0]; op {
	case "+":
		return "next"
	case "~":
		return "after"
	case ">":
		return "child"
	}
	return "descend"
}

func getselectorattr(q Queryable) (key, op, match string) {
	if keys := q.GetAttribute("attrkey"); len(keys) > 0 {
		key = keys[0]
	}
	if ops := q.GetAttribute("attrop"); len(ops) > 0 {
		op = ops[0]
	}
	if matches := q.GetAttribute("attrval"); len(matches) > 0 {
		match = matches[0]
	}
	return
}

func getcolon(q Queryable) (colonspec, colonarg string) {
	if colonspecs := q.GetAttribute("colonspec"); len(colonspecs) > 0 {
		colonspec = colonspecs[0]
	}
	if colonargs := q.GetAttribute("colonarg"); len(colonargs) > 0 {
		colonarg = colonargs[0]
	}
	return
}

func filterbyname(node Queryable, name string) bool {
	if name == "*" {
		return true
	} else if strings.ToLower(node.GetName()) != strings.ToLower(name) {
		return false
	}
	return true
}

func filterbyattr(node Queryable, key, op, match string) bool {
	doop := func(op string, val, match string) bool {
		switch op {
		case "=":
			if val == match {
				return true
			}
		case "~=":
			return strings.Contains(val, match)
		case "^=":
			return strings.HasPrefix(val, match)
		case "$=":
			return strings.HasSuffix(val, match)
		case "*=":
			ok, err := regexp.Match(match, []byte(val))
			if err != nil {
				fmsg := "error matching %v,%v: %v"
				panic(fmt.Errorf(fmsg, match, val, err))
			}
			return ok
		}
		return false
	}
	if key == "" {
		return true
	} else if nodeval := node.GetValue(); key == "value" && op == "" {
		return nodeval != ""
	} else if key == "value" {
		return doop(op, nodeval, match)
	}
	vals := node.GetAttribute(key)
	if vals == nil {
		return false
	} else if vals != nil && op == "" {
		return true
	}
	for _, val := range vals {
		if doop(op, val, match) {
			return true
		}
	}
	return false
}

func filterbycolon(
	parent Queryable, idx int, node Queryable,
	colonspec, colonarg string) bool {

	if colonspec == "" {
		return true
	}

	carg, _ := strconv.Atoi(strings.Trim(colonarg, "()"))
	switch colonspec {
	case "empty":
		if len(node.GetChildren()) == 0 {
			return true
		}

	case "first-child":
		if idx == 0 {
			return true
		}

	case "first-of-type":
		if parent == nil {
			return true
		}
		name, children := node.GetName(), parent.GetChildren()
		for i := 0; i < len(children) && i < idx; i++ {
			if name == children[i].GetName() {
				return false
			}
		}
		return true

	case "last-child":
		if parent == nil || idx == len(parent.GetChildren())-1 {
			return true
		}
	case "last-of-type":
		if parent == nil {
			return true
		}
		name, children := node.GetName(), parent.GetChildren()
		for i := idx + 1; i >= 0 && i < len(children); i++ {
			if name == children[i].GetName() {
				return false
			}
		}
		return true
	case "nth-child":
		if parent == nil || carg == idx {
			return true
		}
	case "nth-of-type":
		if parent == nil {
			return true
		}
		n, name, children := -1, node.GetName(), parent.GetChildren()
		for i := 0; i < len(children) && i <= idx; i++ {
			if name == children[i].GetName() {
				n++
			}
		}
		if n == carg {
			return true
		}

	case "nth-last-child":
		if parent == nil || len(parent.GetChildren())-1 == idx {
			return true
		}

	case "nth-last-of-type":
		if parent == nil {
			return true
		}
		n, name, children := -1, node.GetName(), parent.GetChildren()
		for i := len(children) - 1; i >= 0 && i >= idx; i-- {
			if name == children[i].GetName() {
				n++
			}
		}
		if n == carg {
			return true
		}

	case "only-of-type":
		if parent == nil {
			return true
		}
		n, name := 0, node.GetName()
		for _, child := range parent.GetChildren() {
			if child.GetName() == name {
				n++
			}
		}
		if n == 1 {
			return true
		}

	case "only-child":
		if parent == nil || len(parent.GetChildren()) == 1 {
			return true
		}
	}
	return false
}
