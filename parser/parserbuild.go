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
	{"Expression",[]string{"+","PrimaryExpr"}},
	{"Expression",[]string{"-","PrimaryExpr"}},
	{"Expression",[]string{"!","PrimaryExpr"}},
	{"Expression",[]string{"PrimaryExpr"}},
	{"Expression",[]string{"Expression","||","Expression"}},
	{"Expression",[]string{"Expression","&&","Expression"}},
	{"Expression",[]string{"Expression","==","Expression"}},
	{"Expression",[]string{"Expression","!=","Expression"}},
	{"Expression",[]string{"Expression",">","Expression"}},
	{"Expression",[]string{"Expression","<","Expression"}},
	{"Expression",[]string{"Expression","+","Expression"}},
	{"Expression",[]string{"Expression","-","Expression"}},
	{"Expression",[]string{"Expression","*","Expression"}},
	{"Expression",[]string{"Expression","/","Expression"}},

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
	{"Declaration",[]string{"identifier","Type"}},
	// Blocks
	{"Block",[]string{"{","StatementList","}"}},

	{"StatementList",[]string{"Statement","StatementList"}},
	{"StatementList",[]string{"Statement"}},
	// 类型
	{"Type",[]string{"[","int","]","var"}},
},nil}

