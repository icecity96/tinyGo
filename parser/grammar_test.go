package parser

import (
	"testing"
	"strings"
	"fmt"
)

func TestRule_Show(t *testing.T) {
	var s = Rule{"E",[]string {"Term","+","Term"} }
	str := s.Show("->",1)
	t.Logf("%s", str)
}

func TestIsTerminals(t *testing.T) {
	if IsTerminals("String") {
		t.Error("终结符测试错误")
	} else {
		t.Log("测试通过")
	}
}

func TestGrammar_GetTerminalsAndNoTerminals(t *testing.T) {
	var s = &Rule{"E",[]string {"Term","+","Term"} }
	var g = &Grammar{[]*Rule {s}, nil }
	g.CollectSymbols()
	terminal, noterminal := g.GetTerminalsAndNoTerminals()
	if strings.Join(terminal," ") == "+" && strings.Join(noterminal," ") == "E Term" {
		t.Logf("Collect success with :\n %s\n%s",strings.Join(terminal," "), strings.Join(noterminal," "))
	} else {
		t.Errorf("%s\n%s", strings.Join(terminal," "), strings.Join(noterminal," "))
	}
}

// The test grammer from P162
func TestItemSet_Closure(t *testing.T) {
	var g = &Grammar{[]*Rule{
		{"S'",[]string{"S"} },
		{"S",[]string{"B","B"} },
		{"B",[]string{"a","B"} },
		{"B",[]string{"b"}},
	}, nil}
	g.CollectSymbols()
	var is = make(ItemSet)
	is.Add(Item{&Rule{"S",[]string{"a","B"}},"EOF",1})
	is.Closure(g)
	for item := range is {
		fmt.Printf("%s	%s\n",item.rule.Show("->",item.pos),item.next)
	}
}

func  TestItemSet_Goto(t *testing.T) {
	var g = &Grammar{[]*Rule{
		{"S'",[]string{"S"} },
		{"S",[]string{"B","B"} },
		{"B",[]string{"a","B"} },
		{"B",[]string{"b"}},
	}, nil}
	g.CollectSymbols()
	var is = make(ItemSet)
	is.Add(Item{&Rule{"S'",[]string{"S"}},"EOF",0})
	is.Closure(g)
	for item := range is {
		fmt.Printf("%s\t\t%s\n",item.rule.Show("->",item.pos),item.next)
	}
	fmt.Println("after consumed `a`")
	out := is.Goto(g, "a")
	for item := range out {
		fmt.Printf("%s\t\t%s\n",item.rule.Show("->",item.pos),item.next)
	}
	fmt.Println("Consumed `b`")
	out2 := out.Goto(g, "b")
	for item := range out2 {
		fmt.Printf("%s\t\t%s\n",item.rule.Show("->",item.pos),item.next)
	}
	fmt.Println("================")
	fmt.Println("Consumed `b`")
	out3 := is.Goto(g, "b")
	for item := range out3 {
		fmt.Printf("%s\t\t%s\n",item.rule.Show("->",item.pos),item.next)
	}

}