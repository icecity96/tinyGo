package parser

import (
	"testing"
	"fmt"
)

func TestGraph(t *testing.T) {
	var g = &Grammar{[]*Rule{
		{"S'",[]string{"S"} },
		{"S",[]string{"a","A","d"} },
		{"S",[]string{"b","A","c"} },
		{"S",[]string{"a","e","c"}},
		{"S",[]string{"b","e","d"}},
		{"A",[]string{"e"} },
	},nil}
	g.CollectSymbols()
	fmt.Println(len(g.symbols))
	fmt.Println(len(g.First()))
	actionTable := ComputeActions(g)
	Graph(g,actionTable)
}

func TestGraph3(t *testing.T) {
	var g = &Grammar{[]*Rule{
		{"S'",[]string{"S"} },
		{"S",[]string{"S","+","S"} },
		{"S",[]string{"S","*","S"} },
		{"S",[]string{"a"}},
	},nil}
	g.CollectSymbols()
	act := ComputeActions(g)
	act.Dump()
}

func TestGraph2(t *testing.T) {
	var g = &Grammar{ []*Rule{
		// Grammer start:
		{"Program",[]string{"StatementList"}},
		// for
		{"ForStmt",[]string{"for","Expression","Block"}},
		{"ForStmt",[]string{"for","ForClause","Block"}},
		{"ForClause",[]string{"SimpleStmt",";","Expression",";","SimpleStmt"}},
		// if
		{"IfStmt",[]string{"if","Expression","Block"}},
		{"IfStmt",[]string{"if","Expression","Block","else","IfStmt"}},
		{"IfStmt",[]string{"if","Expression","Block","else","Block"}},
		// assignment
		{"Assignment",[]string{"Expression","=","Expression"}},
		// 语句
		{"Statement",[]string{"Declaration"}},
		{"Statement",[]string{"SimpleStmt"}},
		{"Statement",[]string{"break"}},
		{"Statement",[]string{"continue"}},
		{"Statement",[]string{"Block"}},
		{"Statement",[]string{"IfStmt"}},
		{"Statement",[]string{"ForStmt"}},
		{"SimpleStmt",[]string{"Expression"}},
		{"SimpleStmt",[]string{"Assignment"}},
		// 表达式
		{"Expression",[]string{"UnaryExpr"}},
		{"Expression",[]string{"Expression","BinaryOp","Expression"}},
		{"BinaryOp",[]string{"||"}},
		{"BinaryOp",[]string{"&&"}},
		{"BinaryOp",[]string{"=="}},
		{"BinaryOp",[]string{"!="}},
		{"BinaryOp",[]string{">"}},
		{"BinaryOp",[]string{"<"}},
		{"BinaryOp",[]string{"+"}},
		{"BinaryOp",[]string{"-"}},
		{"BinaryOp",[]string{"*"}},
		{"BinaryOp",[]string{"/"}},
		{"UnaryExpr",[]string{"PrimaryExpr"}},
		{"UnaryExpr",[]string{"UnaryOP","UnaryExpr"}},
		{"UnaryOP",[]string{"+"}},
		{"UnaryOP",[]string{"-"}},
		{"UnaryOP",[]string{"!"}},
		{"PrimaryExpr",[]string{"Operand"}},
		{"PrimaryExpr",[]string{"PrimaryExpr","Index"}},
		{"Index",[]string{"[","Expression","]"}},
		// Operand
		{"Operand",[]string{"Literal"}},
		{"Operand",[]string{"identifier"}},
		{"Operand",[]string{"(","Expression",")"}},
		{"Literal",[]string{"int"}},
		{"Literal",[]string{"float"}},
		// 声明
		{"Declaration",[]string{"identifier",":=","Expression"}},
		{"Declaration",[]string{"identifier","Type","=","Expression"}},
		// Blocks
		{"Block",[]string{"{","StatementList","}"}},

		{"StatementList",[]string{"Statement","StatementList"}},
		{"StatementList",[]string{"Statement"}},
		// 类型
		{"Type",[]string{"[","int","]","var"}},
	},nil}
	g.CollectSymbols()
	actionTable := ComputeActions(g)
	Graph(g, actionTable)
}

func TestGraph4(t *testing.T) {
	var g = &Grammar{[]*Rule{
		{"Program",[]string{"SL"} },
		{"SL",[]string{"ST","SL"} },
		{"SL",[]string{"ST"} },
		{"ST",[]string{"st"}},
		{"ST",[]string{"if"}},
		{"ST",[]string{"for"}},
		{"ST",[]string{"for2"}},
		{"ST",[]string{"Block"}},
		{"Block",[]string{"{","SL","}"} },
	},nil}
	g.CollectSymbols()
	actionTable := ComputeActions(g)
	Graph(g,actionTable)
}