package parser

import (
	"fmt"
	"myGo/mytoken"
)

// parser manages the parsing process
type Parser struct {
	actions ActionTable
	stack   []int
	data    []interface{}
}

func NewParser(ac ActionTable) *Parser {
	return &Parser{
		actions: ac,
		stack:   []int{0},
		data:    []interface{}{},
	}
}

type newToken struct {
	tok *mytoken.Token
	lit string
}

func (nt *newToken) String() string {
	return nt.tok.String()
}

type Node struct {
	val  int    // node 的值
	id   string // 名称，用于符号表和中间代码生成
	code string // 用于代码生成
}

// 语义分析栈,用来存放节点(只有非终结符才能生成节点)
var semStack = make([]Node, 1024)
var top = 0

// TODO: fill
// Note: 有些规约里面虽然有语义动作，但不会对语义分析栈造成影响，所以无需为其构建单独的处理函数
var FunctionTables = map[string]func(){
	"CheckDup":     CheckDup,
	"Lexval":       Lexval,
	"Id2Operand":   Id2Operand,
	"InstallId":    InstallId,
	"InstallArray": InstallArray,
	"AddExpr":      AddExpr,
	"SubExpr":      SubExpr,
	"MulExpr":      MulExpr,
	"DivExpr":      DivExpr,
	"LogicAnd":     LogicAnd,
	"LogicOr":      LogicOr,
	"Equal":        Equal,
	"NotEqual":     NotEqual,
	"Large":        Large,
	"Less":         Less,
	"ZPrimary":     Zprimary,
	"FPrimary":     Fprimary,
	"NPrimary":     Nprimary,
	"For1":         For1,
	"NewST":        NewST,
	"EndBlock":     EndBlock,
	"Assign":       Assign,
	"IF1":          IF1,
}

//前一个有值的词法单元
var preToke newToken
var preId newToken
var preInt newToken

func (p *Parser) Parser(tok *newToken, start string, trace bool) (bool, error) {
	for {
		action, ok := p.actions[p.stack[len(p.stack)-1]][tok.String()]
		if !ok {
			return false, fmt.Errorf("unexpected token: %v", tok.String())
		}
		switch action.(type) {
		case Shift:
			preToke = *tok
			if tok.String() == "identifier" {
				preId = *tok
			}
			if tok.String() == "int" {
				preInt = *tok
			}
			nextState := action.(Shift).state
			p.data = append(p.data, tok.String())
			p.stack = append(p.stack, nextState)
			return false, nil
		case Reduce:
			rule := action.(Reduce).rule
			if !trace {
				fmt.Printf("input %v => reduce %s -> %s\n", tok.lit, rule.pattern, rule.symbol)
			}
			// 如果发生空产生式我们就进行动作执行
			if rule.pattern[0] != "" {
				popCount := len(rule.pattern)
				p.stack = p.stack[0 : len(p.stack)-popCount]
			} else {
				FunctionTables[rule.symbol]()
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
