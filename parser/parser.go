package parser

import (
	"fmt"
	"myGo/mytoken"
	"strconv"
)

// parser manages the parsing process
type Parser struct {
	actions 	ActionTable
	stack 		[]int
	data 		[]interface{}
}

func NewParser(ac ActionTable) *Parser {
	return &Parser{
		actions: ac,
		stack: 	[]int{0},
		data: 	[]interface{}{},
	}
}

type newToken struct {
	tok *mytoken.Token
	lit string
}

func (nt *newToken) String() string {
	return nt.tok.String()
}

var SymbolTables map[int]map[string]int = map[int]map[string]int {
	0 : {},
}
// 全局变量符号表
var roorTable int = 0
// 当前符号表深度
var currentTableDeepth = 0

type Node struct {
	val int		// node 的值
	id  string	// 名称，用于符号表和中间代码生成
	code string // 用于代码生成
}

// 语义分析栈,用来存放节点(只有非终结符才能生成节点)
var semStack = make([]Node,1024)
var top = 0
// TODO: fill
// Note: 有些规约里面虽然有语义动作，但不会对语义分析栈造成影响，所以无需为其构建单独的处理函数
var FunctionTables  = map[string]func(token *newToken) {
	"CheckDup" : CheckDup,
	"Lexval" : Lexval,
	"Id2Operand" : Id2Operand,
}


//前一个有值的词法单元
var preToke newToken

func (p *Parser) Parser(tok *newToken, start string,trace bool) (bool, error) {
	for {
		action, ok := p.actions[p.stack[len(p.stack)-1]][tok.String()]
		if !ok {
			return false, fmt.Errorf("unexpected token: %v", tok.String())
		}
		switch action.(type) {
		case Shift:
			preToke = *tok
			nextState := action.(Shift).state
			p.data = append(p.data, tok.String())
			p.stack = append(p.stack, nextState)
			return false, nil
		case Reduce:
			rule := action.(Reduce).rule
			if trace {
				fmt.Printf("input %v => reduce %s -> %s\n", tok.String(), rule.pattern, rule.symbol)
			}
			// 如果发生空产生式我们就进行动作执行
			if  rule.pattern[0] != "" {
				popCount := len(rule.pattern)
				p.stack = p.stack[0 : len(p.stack) - popCount]
			} else {
				FunctionTables[rule.symbol](tok)
			}

			if rule.symbol == start {
				// Accept
				return true, nil
			}

			state := p.stack[len(p.stack)-1]
			action, ok = p.actions[state][rule.symbol]
			if _, well := action.(Reduce); !ok || well {
				panic(fmt.Errorf("parse error, bad next state"))
			}

			p.stack = append(p.stack, action.(Shift).state)
		default:
			return false, fmt.Errorf("unkonw action!")
		}
	}
}

