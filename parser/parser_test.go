package parser

import (
	"myGo/mytoken"
	"myGo/scanner"
	"testing"
)

// TODO: think about lalr

func TestNewParser(t *testing.T) {
	src := []byte(`i := 1
	j := 2
	m := i
	j = i + m + 3
	`)
	var s scanner.Scanner
	file := mytoken.Newfile("",0,len(src))
	G.CollectSymbols()
	ac := ComputeActions(G)
	p := NewParser(ac)
	s.Init(file,src,nil,scanner.ScanComments)
	for {
		_, tok, lit := s.Scan()
		ok,_ := p.Parser(&newToken{&tok,lit},"Program",true)
		if ok {
			break
		}
	}
}

func TestNewParser2(t *testing.T) {
	src := []byte(`id+id*id`)
	var g = &Grammar{ []*Rule{
		{"E'",[]string{"E"}},
		{"E",[]string{"E","*","E"}},
		{"E",[]string{"E","+","E"}},
		{"E",[]string{"identifier"}},
	},nil}
	var s scanner.Scanner
	file := mytoken.Newfile("",0,len(src))
	g.CollectSymbols()
	ac := ComputeActions(g)
	ac.Dump()
	p := NewParser(ac)
	s.Init(file,src,nil,scanner.ScanComments)
	for {
		_, tok, lit := s.Scan()
		ok,_ := p.Parser(&newToken{&tok,lit},"E'",true)
		if ok {
			break
		}
	}
}
