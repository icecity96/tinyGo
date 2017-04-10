package parser

// to auto gen DNF

import (
	"fmt"
	"log"
	"unicode"
)

const middot = "\u00b7" // middle dot

type SymbolSet map[string]bool

func (ss SymbolSet) Add(s string) { ss[s] = true }
func (ss SymbolSet) Has(s string) bool { return ss[s] }
func (ss SymbolSet) Merge(other SymbolSet) bool {
	s := len(ss)
	for k := range other {
		ss[k] = true
	}
	return len(ss) != s
}

// SymbolMap will tell us which the symbolset the symbol belong to
type SymbolMap map[string]SymbolSet

func (sm SymbolMap) Dump(log log.Logger, lable string) {
	log.Println(lable + ":")
	for sym, set := range  sm{
		var setStr string
		for s := range set {
			setStr += s + " "
		}
		log.Printf("	%s: %s\n", sym, setStr)
	}
}

// Rule is the type of grammar rules
// For example:
//		Expr := Term + Term
// we use the first char to diff terminal and noterminal
type Rule struct {
	// The rule name the first char must be Up
	symbol string
	// the parrern of symbols
	pattern []string
}

func (r *Rule) Show(arrow string, mark int) string {
	str := fmt.Sprintf("%s %s", r.symbol, arrow)
	for i, pat := range r.pattern{
		if i == mark {
			str += middot
		}
		str += " " + pat
	}
	if mark == len(r.pattern) {
		str += " " + middot
	}
	return str
}

func IsTerminals(symbol string) bool { return !unicode.IsUpper((rune(symbol[0]))) }

// Grammar is a collection of rules
type Grammar struct {
	rules 	[]*Rule
	symbols SymbolSet
}

// CollectSymbols walks all the rules to collect all symbols
func (g *Grammar) CollectSymbols() {
	g.symbols = make(SymbolSet)
	for _, rule := range g.rules {
		g.symbols.Add(rule.symbol)
		for _, sym := range rule.pattern {
			g.symbols.Add(sym)
		}
	}
}

// GetTerminalsAndNoTerminals
func (g *Grammar) GetTerminalsAndNoTerminals() (terms []string, noterms []string) {
	for sym := range g.symbols {
		if IsTerminals(sym) {
			terms = append(terms, sym)
		} else {
			noterms = append(noterms, sym)
		}
	}
	return
}

// First computes the "first" set
func (g *Grammar) First() (first SymbolMap) {
	g.CollectSymbols()
	terms, _ := g.GetTerminalsAndNoTerminals()
	first = make(SymbolMap)

	// Initialize: terminals point to themself
	for _, sym := range terms {
		first[sym] = make(SymbolSet)
		first[sym].Add(sym)
	}

	// Fill with grammars first outputs
	for _, rule	:= range g.rules {
		set := first[rule.symbol]
		if set == nil {
			set = make(SymbolSet)
			first[rule.symbol] = set
		}
		set.Add(rule.pattern[0])
	}

	// Iterate until stable
	for changed := true; changed; {
		changed = false
		for _, set := range first {
			for symbol := range set {
				if !IsTerminals(symbol) {
					//set[symbol] = false
					delete(set,symbol)
				}
				if set.Merge(first[symbol]) {
					changed = true
				}
			}
		}
	}
	return
}

// Follow computes the "follow" set
func (g *Grammar) Follow(first SymbolMap) (follow SymbolMap) {
	follow = make(SymbolMap)
	// Initialize Add $ to FOLLOW(S)
	init := make(SymbolSet)
	init.Add("EOF")
	follow[g.rules[0].symbol] = init

	for changed := true; changed; {
		changed = false
		for _, rule := range g.rules {
			for i, patSym := range rule.pattern {
				if IsTerminals(patSym) {
					continue
				}
				set := follow[patSym]
				if set == nil {
					set = make(SymbolSet)
					follow[patSym] = set
				}
				if i+1 < len(rule.pattern) {
					nextSym := rule.pattern[i+1]
					if set.Merge(first[nextSym]) {
						changed = true
					}
				} else {
					if set.Merge(follow[rule.symbol]) {
						changed = true
					}
				}
			}
		}
	}
	return
}