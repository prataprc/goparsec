AST Query
=========

Scope of AST Query is query syntax-tree for desired set of nodes, the query
result if successful, will return an iterable on selected nodes. In that
sense, AST query is simply a `selectors` specification on syntax-tree, similar
to [CSS selectors](https://www.w3schools.com/cssref/css_selectors.asp) into
HTML DOM.

By default, nodes in syntax-tree should implement Queryable behaviour,
similar to Terminal and NonTerminal objects from goparsec pkg. For rest of
this article, top-down recursive parsing using parser-combinators will be
used to explain the ideas behind AST Query. Nevertheless, please note that
AST Query can be used on any syntax-tree that adhere to principles explained
here.

Nodes and attributes
--------------------

Nodes can either be Terminal node or NonTerminal node, also called as
leaf-node or intermediate-node. Typically, leaf-nodes are parsed by
tokernizers and intermediate-nodes are parsed by combinators. In the
Lex-and-Yacc parlance, we can say that leaf-nodes are parsed by
lexers and intermediate-nodes are parsed by yaccer. There are three
aspects to node that are important for `selectors` specification - `name`
and `attributes`.

First `name`. Similar to html tag, each node in the syntax-tree
shall be named. The `name` typically comes from the regular-expression
or the parsing grammar used to construct the node. For instance, with
goparsec:

```go
equal := parsec.Atom(`=`, "EQUAL") // parse comma as a terminal token.
```

The second argument `EQUAL` is the name of the parser. And all nodes
constructed using this parser will be named as `EQUAL`.
Similarly, again with goparsec, to construct nonterminal-nodes:

```go
tagstart := ast.And("TAGSTART", nil, opentag, tagname, closetag)
```

The first argument to the `And` combinator is the name of this parser.
And all intermediate-nodes constructed using this parser will be named
as `TAGSTART`.

**Node names are case-insensitive, should begin with english alphabet,
and contain only alphnanumeric characters.**

Second aspect is `attributes`. Other than `class` and `value` attributes
all other attributes are user-attributes that needs to be set
programmatically. By the way, attribute names are case-insensitive.

**Node attributes are case-insensitive, should begin with english alphabet,
and contain only alphnanumeric characters.**

Default attributes
------------------

**class attribute**

Every node carry atleast one class attribute. If it is intermediate-node,
its `class` attribute is set to `nonterminal`.  If it is leaf-node,
its `class` attribute is set to `terminal`.

**value attribute**

Every node has an underlying value which is a sub-set of parsed input-text.
For a leaf-node, `value` is the text matched by the regular-expression
used in tokeniser. For a intermediate-node, `value` is concatination
of all leaf-nodes' values decending from the intermediate-node.

User-Attributes
---------------

One or more attributes can be set on a node using its Queryable behaviour
APIs:

```go
// SetAttribute with a value string, can be called multiple times for the
// same attrname.
SetAttribute(attrname, value string) Queryable

// GetAttribute for attrname, since more than one value can be set on the
// attribute, return a list of values.
GetAttribute(attrname string) []string

// GetAttributes return a map of all attributes set on this node.
GetAttributes() map[string][]string
```

* Above APIs use golang-syntax, in-general these are applicable to any
language.
* Value string should not contain white-space, if it does, then entire
value should be within single-quote or double-quote.
* Among the user-attributes **id attribute** is treated as special because,
like class, there is a short-hand notation for id.

Selector syntax
---------------

Once we are comfortable with the concepts of: `syntax-tree`, `leaf-node`,
`intermediate-node`, `name`, `attributes`, `value`, and `class`, we can
start specifying the `selectors` syntax which is mostly same as
CSS-Selectors.

`class` and `id` attribute value should start with english character, and
contain - alphabets, numbers, hyphen and underscore.

**Note that node-name is equivalent to html tag-name**

```text

Selector              | Example               | Description
----------------------|-----------------------|---------------------------------
.class                | .terminal             | Selects all terminal nodes.
#id                   | #firstname            | Selects the node with
                      |                       | id="firstname".
*                     | *                     | Selects all nodes.
node,                 | comma                 | Selects all `comma` nodes.
node, node            | comma, equal          | Selects all `comma` nodes and
                      |                       | all `equal` nodes.
node node             | attr equal            | Selects all `equal` nodes inside
                      |                       | `attr`.
node > node           | tag > tagname         | Selects all `tagname` node where
                      |                       | the parent is a `tag` node.
node + node           | oanglebrkt + tagname  | Selects all `tagname` node that
                      |                       | are placed immediately after
                      |                       | `oanglebrkt` elements.
node ~ node           | tagname ~ canglebrkt  | Selects every `tagname` node that
                      |                       | are preceded by `canglebrkt` node.
[attribute]           | [ignore]              | Selects all nodes with a
                      |                       | ignore attribute.
[attribute=value]     | [title=xyz]           | Selects all nodes whose `title`
                      |                       | attribute value is `xyz`.
[attribute~=value]    | [title~=flower]       | Selects all nodes with a `title`
                      |                       | attribute containing the word
                      |                       | `flower`.
[attribute^=value]    | tagname[title^=in]    | Selects every `tagname` node
                      |                       | whose title attribute value
                      |                       | begins with `in`.
[attribute$=value]    | file[path$=.pdf]      | Selects every `file` node whose
                      |                       | path attribute ends with `.pdf`.
[attribute*=value]    | file[path*=usr|opt]   | Selects every `file` node whose
                      |                       | path attribute value matches
                      |                       | regular expression `usr|opt`
:empty                | file:empty            | Selects every `file` node that
                      |                       | has no children.
:first-child          | comma:first-child     | Selects every `comma` node that
                      |                       | is the first child of its parent.
:first-of-type        | comma:first-of-type   | Selects every `comma` node that
                      |                       | is the first `comma` node of
                      |                       | its parent.
:last-child           | comma:last-child      | Selects every `comma` node that
                      |                       | is the last child of its parent.
:last-of-type         | comma:last-of-type    | Selects every `comma` node that
                      |                       | is the last `comma` node of its
                      |                       | parent.
:nth-child(n)         | comma:nth-child(2)    | Selects every `comma` node that
                      |                       | is the second child of its
                      |                       | parent.
:nth-last-child(n)    | eq:nth-last-child(2)  | Selects every `eq` node that
                      |                       | is the second child of its
                      |                       | parent, counting from the last
                      |                       | child.
:nth-last-of-type(n)  | eq:nth-last-of-type(2)| Selects every `eq` node that
                      |                       | is the second `eq` node of
                      |                       | its parent, counting from the
                      |                       | last child.
:nth-of-type(n)       | eq:nth-of-type(2)     | Selects every `eq` node that
                      |                       | is the second `eq` node of
                      |                       | its parent.
:only-of-type         | comma:only-of-type    | Selects every `comma` node that
                      |                       | is the only `comma` node of its
                      |                       | parent.
:only-child           | comma:only-child      | Selects every `comma` node that
                      |                       | is the only child of its parent.
```
