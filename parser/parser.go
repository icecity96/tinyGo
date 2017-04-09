package parser

import (
	"fmt"
	"myGo/mytoken"
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

func (p *Parser) Parser(tok *mytoken.Token, start string,trace bool) (bool, error) {
	for {
		/*
		if !trace {
			fmt.Println("")
			fmt.Printf("stact:%v, data:%v\n", p.stack, p.data)
			fmt.Printf("tok:%v\n", tok.String())
			fmt.Println("")
		}
		*/
		action, ok := p.actions[p.stack[len(p.stack)-1]][tok.String()]
		if !ok {
			//fmt.Println(p.actions[p.stack[len(p.stack)-1]])
			return false, fmt.Errorf("unexpected token: %v", tok.String())
		}
		switch action.(type) {
		case Shift:
			nextState := action.(Shift).state
			/*
			if !trace {
				fmt.Printf("input %v => shift %#v\n", tok.String(), nextState)
			}
			*/
			p.data = append(p.data, tok.String())
			p.stack = append(p.stack, nextState)
			return false, nil
		case Reduce:
			rule := action.(Reduce).rule
			if trace {
				fmt.Printf("input %v => reduce %s -> %s\n", tok.String(), rule.pattern, rule.symbol)
			}
			popCount := len(rule.pattern)
			p.stack = p.stack[0 : len(p.stack) - popCount]

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

