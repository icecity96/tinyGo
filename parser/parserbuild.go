package parser

// 语法
var G = &Grammar{ []*Rule{
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
	{"SimpleStmt",[]string{"identifier","[","int","]"}},
	// 表达式
	// TODO: FILL FUNCTION [ZFN]Primary()
	{"Expression",[]string{"+","PrimaryExpr","ZPrimary"}},
	{"ZPrimary",[]string{""}},

	{"Expression",[]string{"-","PrimaryExpr","FPrimaryExpr"}},
	{"FPrimary",[]string{""}},

	{"Expression",[]string{"!","PrimaryExpr","NPrimary"}},
	{"NPrimary",[]string{""}},

	// this need do notiong
	{"Expression",[]string{"PrimaryExpr"}},

	{"Expression",[]string{"Expression","||","Expression"}},
	{"Expression",[]string{"Expression","&&","Expression"}},
	{"Expression",[]string{"Expression","==","Expression"}},
	{"Expression",[]string{"Expression","!=","Expression"}},
	{"Expression",[]string{"Expression",">","Expression"}},
	{"Expression",[]string{"Expression","<","Expression"}},

	{"Expression",[]string{"Expression","+","Expression","AddExpr"}},
	{"AddExpr",[]string{""}},

	{"Expression",[]string{"Expression","-","Expression"}},
	{"Expression",[]string{"Expression","*","Expression"}},
	{"Expression",[]string{"Expression","/","Expression"}},

	// DO nothing
	{"PrimaryExpr",[]string{"Operand"}},

	{"PrimaryExpr",[]string{"PrimaryExpr","Index"}},
	{"Index",[]string{"[","Expression","]"}},

	// Operand
	// 虽然有语义动作，然并卵
	{"Operand",[]string{"Literal"}},

	{"Operand",[]string{"identifier","Id2Operand"}},
	{"Id2Operand",[]string{""}},

	// need do nothing
	{"Operand",[]string{"(","Expression",")"}},

	{"Literal",[]string{"int","Lexval"}},
	{"Lexval",[]string{""}},

	// 声明
	{"Declaration",[]string{"identifier","CheckDup",":=","Expression","InstallId"}},
	{"Declaration",[]string{"identifier", "CheckDup","Type","InstallArray"}},
	{"CheckDup",[]string{""}},
	{"InstallId",[]string{""}},
	{"InstallArray",[]string{""}},

	// Blocks
	{"Block",[]string{"{","StatementList","}"}},

	{"StatementList",[]string{"Statement","StatementList"}},
	{"StatementList",[]string{"Statement"}},

	// 类型
	// need do nothing
	{"Type",[]string{"[","int","]","var"}},
},nil}

