package main

// Checks that a tool is read correctly, so that a later agent can describe it to a model.

import (
	"strings"
	"testing"
)

func TestToolKeepsWhatAModelNeeds(t *testing.T) {
	stmts, err := Parse(Lex("tool find(city: text, days: number, exact: bool, tags: list): text \"Finds the \\\"weather\\\"\" {\n return city\n}\nfunc plain(a) {\n return a\n}\ntool now() \"The time\" {\n}"))
	if err != nil {
		t.Fatal(err)
	}
	find := stmts[0].(FuncStatement)
	if !find.IsTool || find.Name != "find" || find.Description != `Finds the "weather"` || find.ReturnType != "text" || find.Line != 1 {
		t.Errorf("find = %+v", find)
	}
	want := []Param{{"city", "text"}, {"days", "number"}, {"exact", "bool"}, {"tags", "list"}}
	for i, p := range want {
		if find.Params[i] != p {
			t.Errorf("parameter %d = %v, want %v", i, find.Params[i], p)
		}
	}
	if plain := stmts[1].(FuncStatement); plain.IsTool || plain.Description != "" {
		t.Errorf("a func must not be a tool: %+v", plain)
	}
	if now := stmts[2].(FuncStatement); !now.IsTool || len(now.Params) != 0 || now.ReturnType != "" {
		t.Errorf("now = %+v", now)
	}
}

func TestToolMistakesAreSyntaxErrors(t *testing.T) {
	cases := []struct{ source, wants string }{
		{`tool f(a) "d" {}`, `parameter "a" needs a type`},
		{`tool f(a: tensor) "d" {}`, `parameter "a" cannot be tensor`},
		{`tool f(a: any) "d" {}`, `parameter "a" cannot be any`},
		{`tool f(a: dog) "d" {}`, `unknown type "dog"`},
		{`tool f(a: number, a: text) "d" {}`, `two parameters named "a"`},
		{`tool f(a: number) {}`, "needs a description"},
		{`tool f(a: number): number {}`, "needs a description"},
		{`tool f()`, "needs a description"},
		{`tool f() 5 {}`, "needs a description"},
		{`tool f() "" {}`, "empty description"},
		{`tool f() "   \n " {}`, "empty description"},
		{`tool f() "no end {}`, "unterminated string"},
		{`tool f() "bad \q" {}`, "unknown escape"},
		{`tool f() "d"`, "expected LBRACE"},
		{`tool f() "d" { break }`, "can only be used inside a loop"},
		{`while true { tool f() "d" { continue } }`, "can only be used inside a loop"},
		{`func f() "d" {}`, "expected LBRACE"}, // only tools have descriptions
		{`tool f "d" {}`, "expected LPAREN"},
	}
	for _, c := range cases {
		_, err := Parse(Lex(c.source))
		if err == nil || !strings.Contains(err.Error(), c.wants) {
			t.Errorf("%s: got %v, want %q", c.source, err, c.wants)
		}
	}
}

func TestToolIsStillAnOrdinaryName(t *testing.T) {
	source := "let tool = 1\ntool = 2\nfunc use(tool) {\n return tool\n}\nprint(use(tool))\nlet tools = [tool]\ntools[0] = 3\nfunc tool2() {}\n"
	if _, err := Parse(Lex(source)); err != nil {
		t.Error(err)
	}
	// a function may even be called tool
	if _, err := Parse(Lex("func tool(a) {\n return a\n}\ntool(1)")); err != nil {
		t.Error(err)
	}
}
