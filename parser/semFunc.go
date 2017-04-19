package parser

import (
	"fmt"
	"strconv"
)

// 符号表搜索
func findSymbol(id string)(int,bool) {
	for i := currentTableDeepth; i > 0; i-- {
		if num, ok := SymbolTables[i][id]; ok {
			return num,true
		}
	}
	return 0, false
}

func Id2Operand(token *newToken) {
	var node Node
	if num, ok := findSymbol(preToke.lit); ok {
		node = Node{val:num, id:preToke.lit}
	} else {
		node = Node{id: preToke.lit}
	}
	semStack[top] = node
	top++
	fmt.Println("You have get an id right?And the id is ",semStack[top-1].id, "and the value is ",semStack[top-1].val)
}

func Lexval(token *newToken) {
	num,_ := strconv.Atoi(preToke.lit)
	node := Node{val: num}
	semStack[top] = node
	top++
	fmt.Println("You have get an int right?And the value is ",semStack[top-1].val)
}

func CheckDup(token *newToken) {
	symbolTable := SymbolTables[currentTableDeepth]
	if _,ok := symbolTable[token.String()]; ok {
		fmt.Println("重复声明变量")
	}
}

