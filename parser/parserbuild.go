package parser

// 语法
var G = &Grammar{[]*Rule{
	// Grammer start:
	{"Program", []string{"StatementList"}},
	// for
	{"ForStmt", []string{"for", "Expression", "For1", "Block"}},
	{"For1", []string{""}},

	{"ForStmt", []string{"for", "ForClause", "Block"}},
	{"ForClause", []string{"SimpleStmt", ";", "Expression", ";", "SimpleStmt"}},
	// if
	{"IfStmt", []string{"if", "Expression", "IF1", "Block"}},
	{"IfStmt", []string{"if", "Expression", "IF1", "Block", "else", "IfStmt"}},
	{"IfStmt", []string{"if", "Expression", "IF1", "Block", "else", "Block"}},
	{"IF1", []string{""}},

	// assignment
	{"Assignment", []string{"Expression", "=", "Expression", "Assign"}},
	{"Assign", []string{""}},
	// 语句
	{"Statement", []string{"Declaration"}},
	{"Statement", []string{"SimpleStmt"}},
	{"Statement", []string{"break"}},
	{"Statement", []string{"continue"}},
	{"Statement", []string{"Block"}},
	{"Statement", []string{"IfStmt"}},
	{"Statement", []string{"ForStmt"}},
	{"SimpleStmt", []string{"Expression"}},
	{"SimpleStmt", []string{"Assignment"}},
	{"SimpleStmt", []string{"identifier", "[", "int", "]"}},

	// 表达式
	{"Expression", []string{"+", "PrimaryExpr", "ZPrimary"}},
	{"ZPrimary", []string{""}},

	{"Expression", []string{"-", "PrimaryExpr", "FPrimary"}},
	{"FPrimary", []string{""}},

	{"Expression", []string{"!", "PrimaryExpr", "NPrimary"}},
	{"NPrimary", []string{""}},

	// this need do notiong
	{"Expression", []string{"PrimaryExpr"}},

	{"Expression", []string{"Expression", "||", "Expression", "LogicOr"}},
	{"LogicOr", []string{""}},

	{"Expression", []string{"Expression", "&&", "Expression", "LogicAnd"}},
	{"LogicAnd", []string{""}},

	{"Expression", []string{"Expression", "==", "Expression", "Equal"}},
	{"Equal", []string{""}},

	{"Expression", []string{"Expression", "!=", "Expression", "NotEqual"}},
	{"NotEqual", []string{""}},

	{"Expression", []string{"Expression", ">", "Expression", "Large"}},
	{"Large", []string{""}},

	{"Expression", []string{"Expression", "<", "Expression", "Less"}},
	{"Less", []string{""}},

	{"Expression", []string{"Expression", "+", "Expression", "AddExpr"}},
	{"AddExpr", []string{""}},

	{"Expression", []string{"Expression", "-", "Expression", "SubExpr"}},
	{"SubExpr", []string{""}},

	{"Expression", []string{"Expression", "*", "Expression", "MulExpr"}},
	{"MulExpr", []string{""}},

	{"Expression", []string{"Expression", "/", "Expression", "DivExpr"}},
	{"DivExpr", []string{""}},

	// DO nothing
	{"PrimaryExpr", []string{"Operand"}},

	{"PrimaryExpr", []string{"PrimaryExpr", "Index"}},
	{"Index", []string{"[", "Expression", "]"}},

	// Operand
	// 虽然有语义动作，然并卵
	{"Operand", []string{"Literal"}},

	{"Operand", []string{"identifier", "Id2Operand"}},
	{"Id2Operand", []string{""}},

	// need do nothing
	{"Operand", []string{"(", "Expression", ")"}},

	{"Literal", []string{"int", "Lexval"}},
	{"Lexval", []string{""}},

	// 声明
	{"Declaration", []string{"identifier", "CheckDup", ":=", "Expression", "InstallId"}},
	{"Declaration", []string{"identifier", "CheckDup", "Type", "InstallArray"}},
	{"CheckDup", []string{""}},
	{"InstallId", []string{""}},
	{"InstallArray", []string{""}},

	// Blocks
	{"Block", []string{"{", "NewST", "StatementList", "}", "EndBlock"}},
	{"NewST", []string{""}},
	{"EndBlock", []string{""}},

	{"StatementList", []string{"Statement", "StatementList"}},
	{"StatementList", []string{"Statement"}},

	// 类型
	// need do nothing
	{"Type", []string{"[", "int", "]", "var"}},
}, nil}
