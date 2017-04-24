package parser

import (
	"fmt"
	"os"
	"strconv"
)

var currentOffset = 0

var numTemp = 0

type Attribute struct {
	tp     int // 符号类型 0:int 1:bool 2:array
	num    int // 符号的值
	len    int // 数组长度
	offset int // 偏移量
	values map[int]int
}

// control[i] 表明符号表i中能够访问的符号表
// []int 中的数倒序存放
var control = map[int][]int{
	0: {0},
}

var SymbolTables = map[int]map[string]Attribute{
	0: {},
}

// 当前符号表深度
var currentTable = 0
var totalTable = 0
var labelnum = 0

// 符号表搜索
func findSymbol(id string) (int, bool) {
	if v, ok := SymbolTables[currentTable][preToke.lit]; ok {
		return v.num, true
	}
	for i := range control[currentTable] {
		if sym, ok := SymbolTables[i][id]; ok {
			return sym.num, true
		}
	}
	return 0, false
}

func Id2Operand() {
	var node Node
	if num, ok := findSymbol(preToke.lit); ok {
		node = Node{val: num, id: preToke.lit}
	} else {
		node = Node{id: preToke.lit}
	}
	semStack[top] = node
	top++
}

func Lexval() {
	num, _ := strconv.Atoi(preToke.lit)
	node := Node{id: "", val: num}
	semStack[top] = node
	top++
}

func CheckDup() {
	if _, ok := SymbolTables[currentTable][preToke.lit]; ok {
		fmt.Println("重复声明变量")
	}
	node := Node{id: preToke.lit}
	semStack[top] = node
	top++
}

func InstallId() {
	attr := Attribute{
		num:    semStack[top-1].val,
		offset: currentOffset,
		len:    1,
		tp:     1,
	}
	SymbolTables[currentTable][semStack[top-2].id] = attr
	currentOffset = currentOffset + 4
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-2].id, " = ", attr.num)
	} else {
		fmt.Println(semStack[top-2].id, " = ", semStack[top-1].id)
	}
	// consumer Expr so top--
	top = top - 2
}

func InstallArray() {
	l, _ := strconv.Atoi(preInt.lit)
	v := make(map[int]int)
	attr := Attribute{
		tp:     2,
		len:    l,
		offset: currentOffset,
		values: v,
	}
	SymbolTables[currentTable][preId.lit] = attr
	currentOffset = currentOffset + 4*l
	top = top - 1
}

func AddExpr() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " + ")
	} else {
		fmt.Print(semStack[top-2].id, " + ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}
	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    semStack[top-2].val + semStack[top-1].val,
	}
	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func SubExpr() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " - ")
	} else {
		fmt.Print(semStack[top-2].id, " - ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}
	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    semStack[top-2].val - semStack[top-1].val,
	}
	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func MulExpr() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " * ")
	} else {
		fmt.Print(semStack[top-2].id, " * ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}
	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    semStack[top-2].val * semStack[top-1].val,
	}
	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func DivExpr() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " / ")
	} else {
		fmt.Print(semStack[top-2].id, " / ")
	}
	if semStack[top-1].val == 0 {
		fmt.Errorf("divide 0!!!")
		os.Exit(1)
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}
	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    semStack[top-2].val / semStack[top-1].val,
	}
	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func LogicAnd() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " and ")
	} else {
		fmt.Print(semStack[top-2].id, " and ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := semStack[top-2].val * semStack[top-1].val
	if val != 0 {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func LogicOr() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " or ")
	} else {
		fmt.Print(semStack[top-2].id, " or ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := 0
	if semStack[top-1].val != 0 || semStack[top-2].val != 0 {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func Equal() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " eq ")
	} else {
		fmt.Print(semStack[top-2].id, " eq ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := 0
	if semStack[top-1].val == semStack[top-2].val {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func NotEqual() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " neq ")
	} else {
		fmt.Print(semStack[top-2].id, " neq ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := 0
	if semStack[top-1].val != semStack[top-2].val {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func Large() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " lg ")
	} else {
		fmt.Print(semStack[top-2].id, " lg ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := 0
	if semStack[top-2].val > semStack[top-1].val {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func Less() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-2].id == "" {
		fmt.Print(semStack[top-2].val, " le ")
	} else {
		fmt.Print(semStack[top-2].id, " le ")
	}
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}

	val := 0
	if semStack[top-2].val < semStack[top-1].val {
		val = 1
	}

	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	top = top - 2
	semStack[top] = Node{val: attr.num, id: t}
	top++
}

func Zprimary() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-1].id == "" {
		fmt.Print(semStack[top-1].val)
	} else {
		fmt.Print(semStack[top-1].id)
	}

	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    semStack[top-1].val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	semStack[top-1] = Node{val: attr.num, id: t}
	fmt.Println("")
}

func Fprimary() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-1].id == "" {
		fmt.Print(-semStack[top-1].val)
	} else {
		fmt.Print("-", semStack[top-1].id)
	}

	attr := Attribute{
		tp:     1,
		len:    1,
		offset: currentOffset,
		num:    -semStack[top-1].val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 4
	numTemp++
	semStack[top-1] = Node{val: attr.num, id: t}
	fmt.Println("")
}

func Nprimary() {
	t := "t" + strconv.Itoa(numTemp)
	fmt.Print(t, " = ")
	if semStack[top-1].id == "" {
		fmt.Print("not", semStack[top-1].val)
	} else {
		fmt.Print("not", semStack[top-1].id)
	}
	val := 0
	if semStack[top-1].val == 0 {
		val = 1
	}
	attr := Attribute{
		tp:     0,
		len:    1,
		offset: currentOffset,
		num:    val,
	}

	SymbolTables[currentTable][t] = attr
	currentOffset = currentOffset + 1
	numTemp++
	semStack[top-1] = Node{val: attr.num, id: t}
	fmt.Println("")
}

var Lbegin = make([]string, 100)
var curlb = 0
var Lend = make([]string, 100)
var curle = 0

func For1() {
	lab1 := "L" + strconv.Itoa(labelnum)
	labelnum++
	lab2 := "L" + strconv.Itoa(labelnum)
	labelnum++
	fmt.Println(lab1)
	fmt.Println("if ", semStack[top-1].id, ".false goto ", lab2)
	Lbegin[curlb] = lab1
	curlb++
	Lend[curle] = lab2
	curle++
	top--
}

func NewST() {
	totalTable++
	SymbolTables[totalTable] = make(map[string]Attribute)
	for num := range control[currentTable] {
		control[totalTable] = append(control[totalTable], num)
	}
	control[totalTable] = append(control[totalTable], totalTable)
	currentTable = totalTable
}

func EndBlock() {
	if curlb != 0 {
		fmt.Println("goto ", Lbegin[curlb-1])
		curlb--
	}
	if curle != 0 {
		fmt.Println(Lend[curle-1])
		curle--
	}
	var backSB int = 0
	for num := range control[currentTable] {
		if num > backSB && num != currentTable {
			backSB = num
		}
	}
	currentTable = backSB
}

func Assign() {
	fmt.Print(semStack[top-2].id, " = ")
	if semStack[top-1].id == "" {
		fmt.Println(semStack[top-1].val)
	} else {
		fmt.Println(semStack[top-1].id)
	}
	top = top - 2
}

func IF1() {
	lab2 := "L" + strconv.Itoa(labelnum)
	labelnum++
	fmt.Println("if ", semStack[top-1].id, ".false goto ", lab2)
	Lend[curle] = lab2
	curle++
	top--
}
