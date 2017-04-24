package parser

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Action is an entry in the action table
type Action interface{}

// Shift is an action means : accept the token and change state
type Shift struct {
	state int
}

// Reduce
type Reduce struct {
	rule *Rule
}

// ActionTable
type ActionTable []map[string]Action

func (at ActionTable) Dump() {
	fmt.Println("Parsing table:")
	for i, actions := range at {
		line := fmt.Sprintf("%2d", i)
		for term, action := range actions {
			line += fmt.Sprintf("	%s:%s", term, action)
		}
		fmt.Println(line)
	}
}

// Item is a partially-parsed production
type Item struct {
	rule *Rule
	next string
	// E -> .Item + Item  ====== pos : 0
	// E -> Item. + Item  ====== pos : 1
	pos int
}

// NextSym returns the next symbol the Item would match and
// whether the Item is at the end of its pattern
func (i Item) NextSym() (sym string, end bool) {
	if i.pos == len(i.rule.pattern) {
		return "", true
	}
	if i.pos+1 == len(i.rule.pattern) && i.rule.pattern[0] == "" {
		return "", true
	}
	return i.rule.pattern[i.pos], false
}

// 项族
type ItemSet map[Item]bool

func (is ItemSet) Add(item Item)      { is[item] = true }
func (is ItemSet) Has(item Item) bool { return is[item] }
func (is ItemSet) Empty() bool        { return len(is) == 0 }
func (is ItemSet) Equal(other ItemSet) bool {
	if len(is) != len(other) {
		return false
	}
	for item := range is {
		if !other[item] {
			return false
		}
	}
	return true
}

func (is ItemSet) Dump(log log.Logger) {
	for item := range is {
		log.Println(" ", item.rule.Show("->", item.pos))
	}
}

// construct LR CLOSURE
func (is ItemSet) Closure(grammar *Grammar) {
	added := make(map[string]bool)
	first := grammar.First()
	for changed := true; changed; {
		changed = false
		for item := range is {
			sym, end := item.NextSym()
			if end {
				continue
			}
			// If we haven't yet added
			if !IsTerminals(sym) && !added[sym] {
				for _, rule := range grammar.rules {
					if rule.symbol == sym {
						if pos := item.pos; pos+1 == len(item.rule.pattern) {
							is.Add(Item{rule, item.next, 0})
							//is.Add(Item{rule, "",0})
						} else {
							var nextString string
							mstring := make(map[string]bool)
							var nstring []string
							for set := range first[item.rule.pattern[pos+1]] {
								if set == "" {
									for _, m := range strings.Split(item.next, "#") {
										mstring[m] = true
									}
								}
								mstring[set] = true
							}
							for ms := range mstring {
								if mstring[ms] {
									nstring = append(nstring, ms)
								}
							}
							sort.Slice(nstring, func(i, j int) bool { return nstring[i] < nstring[j] })
							nextString = strings.Join(nstring, "#")
							is.Add(Item{rule, nextString, 0})
							//is.Add(Item{rule,"",0})
						}
						changed = true
						added[sym] = true
					}
				}
			}
		}
	}
}

// GOTO
func (is ItemSet) Goto(grammar *Grammar, x string) ItemSet {
	out := make(ItemSet) // 将J初始化为空
	for item := range is {
		if sym, end := item.NextSym(); !end && sym == x {
			out.Add(Item{item.rule, item.next, item.pos + 1})
			//out.Add(Item{item.rule, "",item.pos+1})
		}
	}
	out.Closure(grammar)
	return out
}

func ComputeActions(grammar *Grammar) ActionTable {
	first := grammar.First()
	follow := grammar.Follow(first)

	var allActions ActionTable

	// 将C初始化为{CLOSURE}([S'->.S,$])
	states := []ItemSet{
		ItemSet{Item{grammar.rules[0], "EOF", 0}: true},
		//		ItemSet{Item{grammar.rules[0],"",0}: true},
	}
	states[0].Closure(grammar)

	// Construcr the parsing list by computing goto() for each state and
	// terminal
	// C中的每个项集I
	for i := 0; i < len(states); i++ {
		set := states[i]
		actions := make(map[string]Action)
		allActions = append(allActions, actions)
		// 每个文法符号
		for term := range grammar.symbols {
			c := set.Goto(grammar, term)
			// GOTO是否为空
			if c.Empty() {
				continue
			}
			id := -1
			// 返回的闭包是否已经在C中了
			for j, oset := range states {
				if c.Equal(oset) {
					id = j
					break
				}
			}
			// 将GOTO(I,X)加入C中
			if id == -1 {
				states = append(states, c)
				id = len(states) - 1
			}
			actions[term] = Shift{id}
		}
	}

	// Add a reduce action for all items that have consumed the full rule.
	for i, set := range states {
		actions := allActions[i]
		for item := range set {
			// middot 在产生式结尾处了 [A->a.,b]
			if _, end := item.NextSym(); !end {
				continue
			}

			f := follow[item.rule.symbol]
			for term := range f {
				if actions[term] != nil {
					for _, a := range item.rule.pattern {
						//fmt.Println(a,precedence(a)," ",precedence(term),term)
						if precedence(a) > precedence(term) {
							actions[term] = Reduce{item.rule}
							break
						}
					}
				} else {
					actions[term] = Reduce{item.rule}
				}
			}
		}
	}
	return allActions
}

func precedence(item string) int {
	switch item {
	case "=":
		return 0
	case "||":
		return 1
	case "&&":
		return 2
	case "":
		return 3
	case "==", "!=", ">", "<":
		return 4
	case "+", "-":
		return 5
	case "*", "/":
		return 6
	case "[", "]":
		return 7
	default:
		return -1
	}
}
