// ast declares the types used to represent syntax tree for myGo
//
package ast

import (
	"myGo/mytoken"
	"strings"
)

// All node types implement the Node interface
type Node interface {
	Pos() mytoken.Pos // position of first character belong to the node
	End() mytoken.Pos // position of first character immediately after the node
}

// All expression nodes implement the Expr interface
type Expr interface {
	Node
	exprNode()
}

// All statement nodes implement the Stmt interface
type Stmt interface {
	Node
	stmtNode()
}

// All declaration nodes implement the Decl interface
type Decl interface {
	Node
	declNode()
}

//-----------------------------------------------------------------
// Comments
//

// A Comment node represents a single /*-style comment
type Comment struct {
	Slash mytoken.Pos // position of "/" starting the comment
	Text  string      // comment text
}

func (c *Comment) Pos() mytoken.Pos { return c.Slash }
func (c *Comment) End() mytoken.Pos { return mytoken.Pos(int(c.Slash) + len(c.Text)) }

// TODO no //-style comment does this needed?
// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
//
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() mytoken.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() mytoken.Pos { return g.List[len(g.List)-1].End() }

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

// Text returns the text of the comment
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = c.Text
	}
	lines := make([]string, 0, 10)
	for _, c := range comments {
		// remove comment markers
		c = c[2 : len(c)-2]
		// Split on newlines
		cl := strings.Split(c, "\n")

		// Walk lines, stripping trailing white space and adding to list
		for _, l := range cl {
			lines = append(lines, stripTrailingWhitespace(l))
		}
	}

	// Remove leading blank lines; convert runs of
	// interior blank lines to a single blank line.
	n := 0
	for _, line := range lines {
		if line != "" || n > 0 && lines[n-1] != "" {
			lines[n] = line
			n++
		}
	}
	lines = lines[0:n]
	if n > 0 && lines[n-1] != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// ------------------------------------------------------------------
// Expression and types

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes
//
type (
	// A BadExpr node is a placeholder for exoressions containing
	// syntax errors for which no corret expression nodes can be
	// created.
	BadExpr struct {
		Form, To mytoken.Pos // position of range of bad expression
	}

	// An Ident node represents an identifier
	Ident struct {
		NamePos mytoken.Pos // identifier position
		Name    string      // identifier name
		//Obj  	*Object		// denoted object; or nil
	}

	// A BasicLit node represents a literal of basic type
	BasicLit struct {
		ValuePos mytoken.Pos   // literal position
		Kind     mytoken.Token // INT,FLOAT,CHAR,STRING
		Value    string        // literal string
	}

	// A CompositeLit node represents a composite literal
	CompositeLit struct {
		Type   Expr        // literal type; or nil
		Lbrace mytoken.Pos // position of "{"
		Elts   []Expr      // list of composite elements; or nil
		Rbrace mytoken.Pos // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression
	ParenExpr struct {
		Lparen mytoken.Pos // position of "("
		X      Expr        // parenthesized expression
		Rparen mytoken.Pos // position of ")"
	}

	// An IndexExpr node represents an expression followed by an index
	IndexExpr struct {
		X      Expr        // expression
		Lbrack mytoken.Pos // position of "["
		Index  Expr        // index expression
		Rbrack mytoken.Pos // position of "]"
	}

	// A TypeAssertExpr node represents an expression followed by a
	// type assertion
	TypeAssertExpr struct {
		X      Expr        // expression
		Lparen mytoken.Pos // position of "("
		Type   Expr        // asserted type;
		Rparen mytoken.Pos // position of ")"
	}

	// A unaryExpr node represents a unary expression.
	UnaryExpr struct {
		OpPos mytoken.Pos   // position of Op
		Op    mytoken.Token // operator
		X     Expr          // operand
	}

	// A BinaryExpr node represents a binary expression
	BinaryExpr struct {
		X     Expr          // left operand
		OpPos mytoken.Pos   // position of Op
		Op    mytoken.Token // operator
		Y     Expr          // right operand
	}
)

type ArrayType struct {
	Lbrack mytoken.Pos // position of "["
	Len    Expr
	Elt    Expr // element type
}

// Pos and End implementations for expression/type nodes
// and exprNode() ensures that only expression/type nodes can
// be assigned to an Expr
//
func (x *BadExpr) Pos() mytoken.Pos { return x.Form }
func (x *BadExpr) End() mytoken.Pos { return x.To }
func (x *BadExpr) exprNode()        {}

func (x *Ident) Pos() mytoken.Pos { return x.NamePos }
func (x *Ident) End() mytoken.Pos { return mytoken.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Ident) exprNode()        {}

func (x *BasicLit) Pos() mytoken.Pos { return x.ValuePos }
func (x *BasicLit) End() mytoken.Pos { return mytoken.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *BasicLit) exprNode()        {}

func (x *CompositeLit) Pos() mytoken.Pos {
	if x.Type != nil {
		return x.Type.Pos()
	}
	return x.Lbrace
}
func (x *CompositeLit) End() mytoken.Pos { return x.Rbrace + 1 }
func (x *CompositeLit) exprNode()        {}

func (x *ParenExpr) Pos() mytoken.Pos { return x.Lparen }
func (x *ParenExpr) End() mytoken.Pos { return x.Rparen + 1 }
func (x *ParenExpr) exprNode()        {}

func (x *IndexExpr) Pos() mytoken.Pos { return x.X.Pos() }
func (x *IndexExpr) End() mytoken.Pos { return x.Rbrack + 1 }
func (x *IndexExpr) exprNode()        {}

func (x *TypeAssertExpr) Pos() mytoken.Pos { return x.X.Pos() }
func (x *TypeAssertExpr) End() mytoken.Pos { return x.Rparen + 1 }
func (x *TypeAssertExpr) exprNode()        {}

func (x *UnaryExpr) Pos() mytoken.Pos { return x.OpPos }
func (x *UnaryExpr) End() mytoken.Pos { return x.X.End() }
func (x *UnaryExpr) exprNode()        {}

func (x *BinaryExpr) Pos() mytoken.Pos { return x.X.Pos() }
func (x *BinaryExpr) End() mytoken.Pos { return x.Y.End() }
func (x *BinaryExpr) exprNode()        {}

func (x *ArrayType) Pos() mytoken.Pos { return x.Lbrack }
func (x *ArrayType) End() mytoken.Pos { return x.Elt.End() }
func (x *ArrayType) exprNode()        {}

// NewIdent creates a new Ident without position
func NewIdent(name string) *Ident {
	return &Ident{mytoken.NoPos, name /*, nil*/}
}

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

// --------------------------------------------------------------------
// Statements

// A Statements is represented by a tree consisting of one or more of the
// following concrete statement nodes
//
type (
	// similar to BadExpr
	BadStmt struct {
		From, To mytoken.Pos // position range of bad statement
	}

	// DeclStmt node represents a declaration in a statement list
	DeclStmt struct {
		Decl Decl // *GenDecl with CONST,TYPE or VAR token
	}

	// A LabeledStmt node represents a labeled statement
	LabeledStmt struct {
		Label *Ident
		Colon *mytoken.Pos // position of ":"
		Stmt  Stmt
	}

	// An ExprStmt node represents a expression
	// in a statement list
	//
	ExprStmt struct {
		X Expr // expression
	}

	// An IncDecStmt node represent an increment or decrement statement
	InDecStmt struct {
		X      Expr
		TokPos mytoken.Pos // position of Tok
		Tok    mytoken.Pos // INC or DEC
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration
	AssignStmt struct {
		Lhs    []Expr
		TokPos mytoken.Pos   //position of TOK
		Tok    mytoken.Token //assigment token, DEFINE
		Rhs    []Expr
	}

	// A ReturnStmt node represents a return statement
	ReturnStmt struct {
		Return  mytoken.Pos // position of "return" keyword
		Results []Expr
	}

	// A BranchStmt node represents a break, continue, goto
	BranchStmt struct {
		TokPos mytoken.Pos   // position of tok
		Tok    mytoken.Token // keyword token (BREAK, CONTINUE, GOTO)
		Label  *Ident        // label name or nil
	}

	// A BlockStmt node represents a braced statement list
	BlockStmt struct {
		Lbrace mytoken.Pos // position of "{"
		List   []Stmt
		Rbrace mytoken.Pos // position of "}"
	}

	// An IfStmt node represents an if statement
	IfStmt struct {
		If   mytoken.Pos // position of "if" keyword
		Init Stmt        // initialization statement; or nil
		Cond Expr        // condition
		Body *BlockStmt
		Else Stmt // else branch; or nil
	}

	//
	ForStmt struct {
		For  mytoken.Pos // position of "for" keyword
		Init Stmt        // initialization statement; or nil
		Cond Expr        // Condition; or nil
		Post Stmt        // post iteration statement; or nil
		Body *BlockStmt
	}
)

func (s *BadStmt) Pos() mytoken.Pos { return s.From }
func (s *BadStmt) End() mytoken.Pos { return s.To }
func (s *BadStmt) stmtNode()        {}

func (s *DeclStmt) Pos() mytoken.Pos { return s.Decl.Pos() }
func (s *DeclStmt) End() mytoken.Pos { return s.Decl.End() }
func (s *DeclStmt) stmtNode()        {}

func (s *LabeledStmt) Pos() mytoken.Pos { return s.Label.Pos() }
func (s *LabeledStmt) End() mytoken.Pos { return s.Stmt.End() }
func (s *LabeledStmt) stmtNode()        {}

func (s *ExprStmt) Pos() mytoken.Pos { return s.X.Pos() }
func (s *ExprStmt) End() mytoken.Pos { return s.X.End() }
func (s *ExprStmt) stmtNode()        {}

func (s *InDecStmt) Pos() mytoken.Pos { return s.X.Pos() }
func (s *InDecStmt) End() mytoken.Pos { return s.TokPos + 2 }
func (s *InDecStmt) stmtNode()        {}

func (s *AssignStmt) Pos() mytoken.Pos { return s.Lhs[0].Pos() }
func (s *AssignStmt) End() mytoken.Pos { return s.Rhs[len(s.Rhs)-1].End() }
func (s *AssignStmt) stmtNode()        {}

func (s *ReturnStmt) Pos() mytoken.Pos { return s.Return }
func (s *ReturnStmt) End() mytoken.Pos {
	if n := len(s.Results); n > 0 {
		return s.Results[n-1].End()
	}
	return s.Return + 6 // len "return"
}
func (s *ReturnStmt) stmtNode() {}

func (s *BranchStmt) Pos() mytoken.Pos { return s.TokPos }
func (s *BranchStmt) End() mytoken.Pos {
	if s.Label != nil {
		return s.Label.End()
	}
	return mytoken.Pos(int(s.TokPos) + len(s.Tok.String()))
}
func (s *BranchStmt) stmtNode() {}

func (s *BlockStmt) Pos() mytoken.Pos { return s.Lbrace }
func (s *BlockStmt) End() mytoken.Pos { return s.Rbrace + 1 }
func (s *BlockStmt) stmtNode()        {}

func (s *IfStmt) Pos() mytoken.Pos { return s.If }
func (s *IfStmt) End() mytoken.Pos {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}
func (s *IfStmt) stmtNode() {}

func (s *ForStmt) Pos() mytoken.Pos { return s.For }
func (s *ForStmt) End() mytoken.Pos { return s.Body.End() }
func (s *ForStmt) stmtNode()        {}

// --------------------------------------------------------------
// Declarations
//

type Spec interface {
	Node
	specNode()
}

type ValueSpec struct {
	Names  []*Ident // value names
	Type   Expr     // value type
	Values []Expr   // initial values; or nil
}

type TypeSpec struct {
	Name *Ident
	Type Expr
}

func (s *ValueSpec) Pos() mytoken.Pos { return s.Names[0].Pos() }
func (s *ValueSpec) End() mytoken.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}
func (*ValueSpec) specNode() {}

type (
	BadDecl struct {
		From, To mytoken.Pos
	}
	//	CONST,VAR *ValueSpec
	//  TYPE *TypeSpec
	GenDecl struct {
		TokPos mytoken.Pos
		Tok    mytoken.Token // const, type, var
		Lparen mytoken.Pos   // position of '(', if any
		Specs  []Spec
		Rparen mytoken.Pos // position of ')', if any
	}
)

func (d *BadDecl) Pos() mytoken.Pos { return d.From }
func (d *BadDecl) End() mytoken.Pos { return d.To }
func (d *BadDecl) declNode()        {}

func (d *GenDecl) Pos() mytoken.Pos { return d.TokPos }
func (d *GenDecl) End() mytoken.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Specs[0].End()
}
func (d *GenDecl) declNode() {}

// ------------------------------------------------------------
// A File node represents a soure file
//

type File struct {
	Decls      []Decl   // top-level declarations; or nil
	Unresolved []*Ident // unresolved identifiers in this file
}
