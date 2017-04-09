package parser

import (
	"myGo/mytoken"
	"myGo/scanner"
	"testing"
)

// TODO: think about lalr

func TestNewParser(t *testing.T) {
	src := []byte(`i := 1
	if i > 1 && i < 5 {
	 i = i + 1
	 }
	 for j := 1; j >0 ; j = j -1 {
	  i = i + 1
	 }
	 k[5]var
	 k[1] = 1`)
	var s scanner.Scanner
	file := mytoken.Newfile("",0,len(src))
	G.CollectSymbols()
	ac := ComputeActions(G)
	p := NewParser(ac)
	s.Init(file,src,nil,scanner.ScanComments)
	for {
		_, tok, _ := s.Scan()
		ok,_ := p.Parser(&tok,"Program",true)
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
		_, tok, _ := s.Scan()
		ok,_ := p.Parser(&tok,"E'",true)
		if ok {
			break
		}
	}
}
