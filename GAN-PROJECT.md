# GAN programming language: project handoff

## About
GAN is my own programming language (file extension .gan), written in Go. It is named after generative adversarial networks; the goal is a language designed for AI and agentic AI. Current version: beta (no version number; do not change it unless I ask). Folder: C:\Desktop\GAN. Go 1.27 on Windows, module `mylang`.
Build: `go build -o gan.exe .`   Run: `gan file.gan`   Version: `gan --version` (GAN beta)
Pipeline: Lexer -> Parser -> Interpreter (lexer.go, parser.go, interpreter.go; builtins.go, modules.go, errors.go, tensor.go, autograd.go, ai.go, python.go; main.go ties them together).

## How I want you to work
- I'm a beginner in Go and in language design. Go one step at a time.
- Before changing any code, explain in plain language what we're doing and why, and wait for me to say "go".
- Then make the change, tell me which files changed, and give me the exact commands to build and test it.
- Include a test .gan file for each new feature, including one that should produce an error.
- Keep the code simple and commented. Don't add features I haven't asked for.
- After each step works, wait for me to say "next".
(For steps 15 to 26 I allowed continuing without stopping. Normal rules apply again from here.)

## Progress: steps 1 to 35 are done
1-6 lexer/parser/print, 7 interpreter, 8 variables, 9 math+precedence, 10 comparisons+if/else, 11 while+assignment,
12 functions, 13 error messages with line numbers, 14 parentheses+negative numbers, 15 comments (//), 16 decimals,
17 joining text with +, 18 true/false + and/or/not + else if, 19 scope + final/var, 20 for loops (for i in 1 to 5),
21 lists, 22 built-in functions, 23 import + libraries (math, text, lists) + % + fail(), 24 optional types (let x: number),
25 Windows installer (installer\GAN-Setup.exe), 26 tutorial (tutorial\), 27 break/continue, 28 tensors, 29 random/transpose/indexing/network maths,
30 first automatic gradients.
31 gradients for / relu sigmoid tanh exp log transpose, 32 ask() for local Ollama models, 33 text escapes,
34 python() bridge (any Python function, no allow-list), 35 tool declarations. Lessons 0-18.
Next (only after I say next/go): Step 36 agent + .do() loop (design in design\agents.md, accepted: words tool/agent/.do,
description after the signature, default 10 steps; the model can run only listed tools). Step 37 agent log.
Later: dictionaries + JSON, train XOR (a scratch trial worked), a tiny GAN network, speed.
Direction: a mix of the best of Go, Python and C#, AI first, no censorship or guardrails on what the programmer can do.

## The language today
- Values: number (float64), text, bool, list, tensor, module. Blocks use { }. Comments use //.
- let/var/final, assignment, if/else if/else, while, for i in a to b (and for x in list/text), func/return, import "name", break, continue.
- Operators: + - * / %, == != < > <= >=, and or not, [index], module.name, f(x). + joins when either side is text.
- Optional types: number, text, bool, list, tensor, any (checked while running). typeof(x).
- Built-ins: len push pop str num input abs floor round sqrt pow min max upper lower contains fail typeof.
- AI built-ins: tensor shape sum mean matmul transpose zeros ones random randomNormal seed relu sigmoid tanh exp log.
- Gradients: learnable backward grad item (through + - * / minus matmul transpose sum mean relu sigmoid tanh exp log; [ ] refuses).
- AI models: ask(question) / ask(question, model) calls Ollama on this computer only (ai.go).
- Also run `go test` (gradients, ask with a fake server, python bridge, tool rules). live\*.gan need Ollama or Python.
- Text escapes: \" \\ \n \t. python("module.function", args) needs Python (python.go + pybridge.py).
- tool name(a: number): text "description" { } : every parameter typed (number/text/bool/list).
- Libraries in libs\ (written in GAN): math, text, lists. A path import ("modules/x.gan") is relative to the importing file.
- Errors: "Syntax error on line N" / "Runtime error on line N" plus the source line; errors in imported files name the file.
- Scope: each block and each function call has its own scope; functions can see their module's globals. Lists are shared, not copied.

## Project layout
- main.go lexer.go parser.go interpreter.go builtins.go modules.go errors.go tensor.go autograd.go ai.go python.go pybridge.py: the language
- Website\: React site published with GitHub Pages; agent\web\: the five website agents
- libs\*.gan standard libraries; modules\*.gan helper files used by tests
- *.gan in the root: examples and tests (names ending in _error are meant to fail)
- tests\run.py: regression tests (python tests\run.py; --update rewrites tests\expected after checking output by eye)
- tutorial\: the course. python tutorial\check.py verifies every example against the real gan.exe (--fill writes outputs)
- agent\documenter.py: Python agent using local Ollama (qwen2.5:7b) that writes agent\document\*.md from the project and agent\steps.json. Run with --watch in a split terminal. Add an entry to agent\steps.json after every step.
- installer\build-installer.ps1 builds installer\GAN-Setup.exe (IExpress + a PowerShell/WinForms wizard, per-user install, no admin).
- vscode-gan\: the "Seti + GAN" VS Code icon theme; gan-icon.svg/.ico: the logo.

## Things worth remembering
- Functions must be defined before the call runs.
- input() echoes the line when stdin is piped, so recorded runs read like real ones.
- Tests print the same thing every time, so the regression runner compares full output.

## main.go
```go
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// report prints an error. If the error knows its line, it says so and
// shows that line of the program underneath.
func report(kind string, err error, source string) {
	var ge *GanError
	if !errors.As(err, &ge) || ge.Line < 1 {
		fmt.Printf("%s: %s\n", kind, err)
		return
	}
	if ge.Kind != "" {
		kind = ge.Kind
	}
	where := fmt.Sprintf("on line %d", ge.Line)
	if ge.File != "" { // the error is inside an imported file: show that file's line instead
		if data, readErr := os.ReadFile(ge.File); readErr == nil {
			source = string(data)
		}
		where = fmt.Sprintf("in %s on line %d", filepath.Base(ge.File), ge.Line)
	}
	fmt.Printf("%s %s: %s\n", kind, where, ge.Msg)
	lines := strings.Split(source, "\n")
	if ge.Line <= len(lines) {
		fmt.Printf("%5d | %s\n", ge.Line, strings.TrimRight(lines[ge.Line-1], "\r"))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gan <file.gan>")
		return
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Println("GAN beta") // keep as "beta": no version number until the owner says so
		return
	}

	path := os.Args[1]
	if !strings.HasSuffix(path, ".gan") {
		fmt.Println("Error: file must end with .gan")
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	source := string(data)

	stmts, err := Parse(Lex(source))
	if err != nil {
		report("Syntax error", err, source)
		return
	}
	if err := Run(stmts, path); err != nil {
		report("Runtime error", err, source)
	}
}
```

## lexer.go
```go
package main

import (
	"strings"
	"unicode/utf8"
)

type TokenType string

const (
	IDENT    TokenType = "IDENT"
	NUMBER   TokenType = "NUMBER"
	STRING   TokenType = "STRING"
	LPAREN   TokenType = "LPAREN"
	RPAREN   TokenType = "RPAREN"
	EQUALS   TokenType = "EQUALS"
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	STAR     TokenType = "STAR"
	SLASH    TokenType = "SLASH"
	LBRACE   TokenType = "LBRACE"
	RBRACE   TokenType = "RBRACE"
	COMMA    TokenType = "COMMA"
	COLON    TokenType = "COLON"
	DOT      TokenType = "DOT"
	PERCENT  TokenType = "PERCENT"
	LBRACKET TokenType = "LBRACKET"
	RBRACKET TokenType = "RBRACKET"
	EQEQ     TokenType = "EQEQ"  // ==
	NOTEQ    TokenType = "NOTEQ" // !=
	LT       TokenType = "LT"    // <
	GT       TokenType = "GT"    // >
	LTE      TokenType = "LTE"   // <=
	GTE      TokenType = "GTE"   // >=
	ILLEGAL  TokenType = "ILLEGAL"
)

// A Token is one small piece of the program. Line says which line of the file it came from.
type Token struct {
	Type  TokenType
	Value string
	Line  int
}

// unterminatedString is the Value of the ILLEGAL token made for a string with no closing quote.
const unterminatedString = "unterminated string"

// badEscape starts the Value of the ILLEGAL token made for a wrong escape such as \q.
// The whole Value is the error message.
const badEscape = "unknown escape"

// escapes are the characters that may follow a backslash inside text, and what each one means.
var escapes = map[byte]byte{
	'"':  '"',  // \"  a double quote
	'\\': '\\', // \\  one backslash
	'n':  '\n', // \n  a new line
	't':  '\t', // \t  a tab
}

// lexString reads text that starts with the " at src[i]. It returns the text with its escapes
// turned into real characters, the position after the closing quote, and how many lines the
// text went down. problem is "" when all is well; otherwise it says what is wrong, and
// problemLines says how many lines down from the opening quote the mistake is.
func lexString(src string, i int) (text string, next, lines int, problem string, problemLines int) {
	var out strings.Builder
	for i++; i < len(src); i++ {
		c := src[i]
		switch {
		case c == '"': // the closing quote
			return out.String(), i + 1, lines, problem, problemLines
		case c == '\\' && i+1 < len(src):
			meaning, known := escapes[src[i+1]]
			if !known {
				if problem == "" { // report the first mistake, but keep reading to the closing quote
					after, _ := utf8.DecodeRuneInString(src[i+1:])
					shown := " \\" + string(after)
					if after == '\n' || after == '\r' {
						shown = ": a backslash at the end of a line"
					}
					problem = badEscape + shown + ` in text (you can use \" \\ \n \t)`
					problemLines = lines
				}
				continue // the next character is read normally, so lines are still counted
			}
			out.WriteByte(meaning)
			i++ // the character after the backslash is used up too
		default:
			if c == '\n' {
				lines++
			}
			out.WriteByte(c)
		}
	}
	return "", len(src), lines, unterminatedString, 0
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func Lex(src string) []Token {
	var tokens []Token
	line := 1 // the line we are on; goes up by one at every newline
	add := func(t TokenType, value string) {
		tokens = append(tokens, Token{t, value, line})
	}

	i := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r':
			i++
		case c == '(':
			add(LPAREN, "(")
			i++
		case c == ')':
			add(RPAREN, ")")
			i++
		case c == '{':
			add(LBRACE, "{")
			i++
		case c == '}':
			add(RBRACE, "}")
			i++
		case c == ',':
			add(COMMA, ",")
			i++
		case c == ':':
			add(COLON, ":")
			i++
		case c == '.':
			add(DOT, ".")
			i++
		case c == '%':
			add(PERCENT, "%")
			i++
		case c == '[':
			add(LBRACKET, "[")
			i++
		case c == ']':
			add(RBRACKET, "]")
			i++
		// Operators made of two characters: look one character ahead.
		// "=" alone is assignment, "==" is a comparison. Same idea for the others.
		case c == '=' && i+1 < len(src) && src[i+1] == '=':
			add(EQEQ, "==")
			i += 2
		case c == '!' && i+1 < len(src) && src[i+1] == '=':
			add(NOTEQ, "!=")
			i += 2
		case c == '<' && i+1 < len(src) && src[i+1] == '=':
			add(LTE, "<=")
			i += 2
		case c == '>' && i+1 < len(src) && src[i+1] == '=':
			add(GTE, ">=")
			i += 2
		case c == '<':
			add(LT, "<")
			i++
		case c == '>':
			add(GT, ">")
			i++
		case c == '=':
			add(EQUALS, "=")
			i++
		case c == '+':
			add(PLUS, "+")
			i++
		case c == '-':
			add(MINUS, "-")
			i++
		case c == '*':
			add(STAR, "*")
			i++
		case c == '/' && i+1 < len(src) && src[i+1] == '/':
			// a comment: ignore everything up to the end of the line
			// (the newline itself is left alone so line counting stays right)
			for i < len(src) && src[i] != '\n' {
				i++
			}
		case c == '/':
			add(SLASH, "/")
			i++
		case c == '"':
			startLine := line // a string is reported on the line where it starts
			text, next, lines, problem, problemLines := lexString(src, i)
			i = next
			line += lines
			switch problem {
			case "":
				tokens = append(tokens, Token{STRING, text, startLine})
			case unterminatedString:
				tokens = append(tokens, Token{ILLEGAL, problem, startLine})
			default: // a wrong escape is reported on the line where it is
				tokens = append(tokens, Token{ILLEGAL, problem, startLine + problemLines})
			}
		case isLetter(c):
			start := i
			for i < len(src) && (isLetter(src[i]) || isDigit(src[i])) {
				i++
			}
			add(IDENT, src[start:i])
		case isDigit(c):
			start := i
			for i < len(src) && isDigit(src[i]) {
				i++
			}
			// a decimal point followed by digits is part of the number: 3.14
			if i+1 < len(src) && src[i] == '.' && isDigit(src[i+1]) {
				i++
				for i < len(src) && isDigit(src[i]) {
					i++
				}
			}
			add(NUMBER, src[start:i])
		default:
			add(ILLEGAL, string(c))
			i++
		}
	}
	return tokens
}
```

## parser.go
```go
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ---- Expressions: things that produce a value ----

// Expr is anything that has a value. For now: a string, a number, a variable name,
// a function call, or math/comparison on two other expressions.
type Expr interface{}

type StringLit struct{ Value string }  // "hello"
type NumberLit struct{ Value float64 } // 5
type VarRef struct{ Name string }      // x
type BoolLit struct{ Value bool }      // true or false

// ListLit is a list written out: [1, 2, 3]
type ListLit struct{ Items []Expr }

// IndexExpr picks one item out of a list (or one character out of text): Target[Index]
type IndexExpr struct {
	Target Expr
	Index  Expr
}

// LogicExpr is "and" / "or": Left Op Right. It stops early when the answer is already known.
type LogicExpr struct {
	Left  Expr
	Op    string // "and" or "or"
	Right Expr
}

// CallExpr is a function call: Name(Args), e.g. add(2, 3). It has a value if the function returns one.
type CallExpr struct {
	Name   string
	Args   []Expr
	Target Expr // for math.square(3), the "math" part; nil for a plain call like square(3)
}

// MemberExpr reads a variable that lives inside a module: math.PI
type MemberExpr struct {
	Target Expr
	Name   string
}

// UnaryExpr is an operation on one value: -Operand, e.g. -5 or -x
type UnaryExpr struct {
	Op      string // "-" or "not"
	Operand Expr
}

// BinaryExpr is an operation on two sides: Left Op Right, e.g. 2 * 3 or x > 5
type BinaryExpr struct {
	Left  Expr
	Op    string // "+", "-", "*", "/", "==", "!=", "<", ">", "<=", ">="
	Right Expr
}

// ---- Statements: things the program does ----
// Every statement remembers its Line, so the interpreter can say where an error happened.

// Statement is any command. All the types below count.
type Statement interface{}

type PrintStatement struct {
	Value Expr // what to print
	Line  int
}

type LetStatement struct {
	Name  string // variable name
	Value Expr   // value to store
	Final bool   // true for "final": the value can never be changed afterwards
	Type  string // the type written after a colon (let x: number = 5), or "" if none
	Line  int
}

// IndexAssignStatement: Target[Index] = Value. Changes one item of a list.
type IndexAssignStatement struct {
	Target Expr
	Index  Expr
	Value  Expr
	Line   int
}

// ImportStatement: import "path". Loads a library or another .gan file.
type ImportStatement struct {
	Path string
	Line int
}

// ForStatement: for Var in From to To { Body } counts from From up to To (both included).
// If there is no "to", Iter is a list to go through instead.
type ForStatement struct {
	Var  string
	From Expr
	To   Expr
	Iter Expr
	Body []Statement
	Line int
}

// IfStatement: if Condition { Then } else { Else }. Else is empty if there is no else.
type IfStatement struct {
	Condition Expr
	Then      []Statement
	Else      []Statement
	Line      int
}

// WhileStatement: while Condition { Body }. The body repeats while the condition is true.
type WhileStatement struct {
	Condition Expr
	Body      []Statement
	Line      int
}

// AssignStatement: name = value. Changes a variable that already exists.
type AssignStatement struct {
	Name  string
	Value Expr
	Line  int
}

// FuncStatement: func Name(Params) { Body }. Defines a function.
type FuncStatement struct {
	Name       string
	Params     []Param
	ReturnType string // the type written after the ) , or "" if none
	Body       []Statement
	Line       int

	// A tool is a function that an AI model will be able to use (see parseTool).
	IsTool      bool
	Description string // tools only: tells a model what the tool does
}

// Param is one parameter of a function, with an optional type.
type Param struct {
	Name string
	Type string
}

// ReturnStatement: return Value. Ends the function and hands back Value.
type ReturnStatement struct {
	Value Expr
	Line  int
}

// BreakStatement: break. Stops the nearest loop.
type BreakStatement struct {
	Line int
}

// ContinueStatement: continue. Skips to the next pass of the nearest loop.
type ContinueStatement struct {
	Line int
}

// ExprStatement is a call written on its own line, e.g. greet("GAN"). Its value is thrown away.
type ExprStatement struct {
	Call CallExpr
	Line int
}

// lineOf tells which line a statement was written on.
func lineOf(s Statement) int {
	switch st := s.(type) {
	case PrintStatement:
		return st.Line
	case LetStatement:
		return st.Line
	case IfStatement:
		return st.Line
	case WhileStatement:
		return st.Line
	case AssignStatement:
		return st.Line
	case FuncStatement:
		return st.Line
	case ReturnStatement:
		return st.Line
	case ForStatement:
		return st.Line
	case ImportStatement:
		return st.Line
	case IndexAssignStatement:
		return st.Line
	case ExprStatement:
		return st.Line
	case BreakStatement:
		return st.Line
	case ContinueStatement:
		return st.Line
	}
	return 0
}

// ---- The parser ----

type Parser struct {
	tokens    []Token
	pos       int
	loopDepth int // how many loops we are inside; break and continue need at least one
}

// parseLoopBody reads the block of a while or for loop, where break and continue are allowed.
func (p *Parser) parseLoopBody() ([]Statement, error) {
	p.loopDepth++
	body, err := p.parseBlock()
	p.loopDepth--
	return body, err
}

func (p *Parser) peek() *Token {
	if p.pos < len(p.tokens) {
		return &p.tokens[p.pos]
	}
	return nil
}

// lastLine is the line of the last token. Errors at the end of the file point there.
func (p *Parser) lastLine() int {
	if len(p.tokens) == 0 {
		return 1
	}
	return p.tokens[len(p.tokens)-1].Line
}

// unexpected builds the error for a token that is not what the parser needed.
// tok is nil when the file ended too soon. wanted says what we expected, e.g. "RPAREN" or "a value".
func (p *Parser) unexpected(tok *Token, wanted string) error {
	if tok == nil {
		return at(p.lastLine(), fmt.Errorf("unexpected end of file, expected %s", wanted))
	}
	if tok.Type == ILLEGAL || tok.Type == DOT {
		if tok.Value == unterminatedString {
			return at(tok.Line, fmt.Errorf("unterminated string (a string needs a closing \")"))
		}
		if strings.HasPrefix(tok.Value, badEscape) {
			return at(tok.Line, errors.New(tok.Value)) // the lexer already wrote the message
		}
		return at(tok.Line, fmt.Errorf("unexpected character %q", tok.Value))
	}
	return at(tok.Line, fmt.Errorf("expected %s but got %s %q", wanted, tok.Type, tok.Value))
}

// expect checks the next token is the type we want, then moves forward
func (p *Parser) expect(t TokenType) (Token, error) {
	tok := p.peek()
	if tok == nil || tok.Type != t {
		return Token{}, p.unexpected(tok, string(t))
	}
	p.pos++
	return *tok, nil
}

// Precedence: each level calls the tighter one below it.
//   parseExpr       handles "or"             (loosest)
//   parseAnd        handles "and"
//   parseNot        handles "not"
//   parseComparison handles == != < > <= >=
//   parseSum        handles + and -
//   parseTerm       handles * and /
//   parseUnary      handles a minus sign in front of a value, like -5
//   parsePrimary    reads one value, or a ( group ) (tightest)
// So x + 1 > 5 means (x + 1) > 5, 2 + 3 * 4 means 2 + (3 * 4), and -2 * 3 means (-2) * 3.

// isWord says whether a token is the given word, e.g. the name "and".
func isWord(tok *Token, word string) bool {
	return tok != nil && tok.Type == IDENT && tok.Value == word
}

// parseExpr reads a whole expression: conditions joined by "or".
func (p *Parser) parseExpr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for isWord(p.peek(), "or") {
		p.pos++
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = LogicExpr{Left: left, Op: "or", Right: right}
	}
	return left, nil
}

// parseAnd reads conditions joined by "and".
func (p *Parser) parseAnd() (Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for isWord(p.peek(), "and") {
		p.pos++
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = LogicExpr{Left: left, Op: "and", Right: right}
	}
	return left, nil
}

// parseNot reads "not" followed by a condition, or just a comparison.
func (p *Parser) parseNot() (Expr, error) {
	if isWord(p.peek(), "not") {
		p.pos++
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{Op: "not", Operand: operand}, nil
	}
	return p.parseComparison()
}

// parseComparison reads a comparison, or just a plain value if there is no comparison.
func (p *Parser) parseComparison() (Expr, error) {
	left, err := p.parseSum()
	if err != nil {
		return nil, err
	}
	tok := p.peek()
	if tok != nil && (tok.Type == EQEQ || tok.Type == NOTEQ || tok.Type == LT ||
		tok.Type == GT || tok.Type == LTE || tok.Type == GTE) {
		p.pos++
		right, err := p.parseSum()
		if err != nil {
			return nil, err
		}
		return BinaryExpr{Left: left, Op: tok.Value, Right: right}, nil
	}
	return left, nil
}

// parseSum reads terms joined by + or -.
func (p *Parser) parseSum() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		tok := p.peek()
		if tok == nil || (tok.Type != PLUS && tok.Type != MINUS) {
			return left, nil
		}
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{Left: left, Op: tok.Value, Right: right}
	}
}

// parseTerm reads values joined by * or /.
func (p *Parser) parseTerm() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		tok := p.peek()
		if tok == nil || (tok.Type != STAR && tok.Type != SLASH && tok.Type != PERCENT) {
			return left, nil
		}
		p.pos++
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{Left: left, Op: tok.Value, Right: right}
	}
}

// parseUnary reads a value that may have a minus sign in front: 5, -5, or even --5.
func (p *Parser) parseUnary() (Expr, error) {
	if tok := p.peek(); tok != nil && tok.Type == MINUS {
		p.pos++
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{Op: "-", Operand: operand}, nil
	}
	return p.parsePostfix()
}

// parsePostfix reads a value, then any [index] written after it: names[0], grid[1][2], f(x)[0].
func (p *Parser) parsePostfix() (Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		tok := p.peek()
		if tok != nil && tok.Type == DOT { // module.name or module.function(...)
			if _, isNumber := expr.(NumberLit); isNumber {
				return nil, at(tok.Line, fmt.Errorf("unexpected \".\" after a number (a number can only have one decimal point)"))
			}
			p.pos++
			member, err := p.expect(IDENT)
			if err != nil {
				return nil, err
			}
			if next := p.peek(); next != nil && next.Type == LPAREN {
				call, err := p.parseCall(member.Value)
				if err != nil {
					return nil, err
				}
				call.Target = expr
				expr = call
			} else {
				expr = MemberExpr{Target: expr, Name: member.Value}
			}
			continue
		}
		if tok == nil || tok.Type != LBRACKET {
			return expr, nil
		}
		p.pos++
		index, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RBRACKET); err != nil {
			return nil, err
		}
		expr = IndexExpr{Target: expr, Index: index}
	}
}

// parsePrimary reads one value: a string, a number, a variable name, a function call,
// or a whole expression inside ( ).
func (p *Parser) parsePrimary() (Expr, error) {
	tok := p.peek()
	if tok == nil {
		return nil, p.unexpected(nil, "a value")
	}
	switch tok.Type {
	case LBRACKET: // [ item, item, ... ] is a list
		p.pos++
		var items []Expr
		if next := p.peek(); next != nil && next.Type == RBRACKET { // empty list: []
			p.pos++
			return ListLit{}, nil
		}
		for {
			item, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if next := p.peek(); next != nil && next.Type == COMMA {
				p.pos++
				continue
			}
			break
		}
		if _, err := p.expect(RBRACKET); err != nil {
			return nil, err
		}
		return ListLit{Items: items}, nil
	case LPAREN: // ( expression ) groups things together
		p.pos++
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RPAREN); err != nil {
			return nil, err
		}
		return inner, nil
	case STRING:
		p.pos++
		return StringLit{tok.Value}, nil
	case NUMBER:
		p.pos++
		n, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, at(tok.Line, fmt.Errorf("bad number %q", tok.Value))
		}
		return NumberLit{n}, nil
	case IDENT:
		p.pos++
		if tok.Value == "true" || tok.Value == "false" {
			return BoolLit{tok.Value == "true"}, nil
		}
		// a name followed by ( is a function call, otherwise a variable
		if next := p.peek(); next != nil && next.Type == LPAREN {
			return p.parseCall(tok.Value)
		}
		return VarRef{tok.Value}, nil
	}
	return nil, p.unexpected(tok, "a value")
}

// parseCall reads the ( arguments ) part of a call. The function name is already consumed.
func (p *Parser) parseCall(name string) (CallExpr, error) {
	if _, err := p.expect(LPAREN); err != nil {
		return CallExpr{}, err
	}
	var args []Expr
	if tok := p.peek(); tok != nil && tok.Type == RPAREN { // no arguments: add()
		p.pos++
		return CallExpr{Name: name}, nil
	}
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return CallExpr{}, err
		}
		args = append(args, arg)
		if tok := p.peek(); tok != nil && tok.Type == COMMA {
			p.pos++ // another argument follows
			continue
		}
		break
	}
	if _, err := p.expect(RPAREN); err != nil {
		return CallExpr{}, err
	}
	return CallExpr{Name: name, Args: args}, nil
}

// validTypes are the type names you can write after a colon.
var validTypes = map[string]bool{"number": true, "text": true, "bool": true, "list": true, "tensor": true, "any": true}

// parseOptionalType reads ": number" if there is a colon next. It returns "" if there is none.
func (p *Parser) parseOptionalType() (string, error) {
	if tok := p.peek(); tok == nil || tok.Type != COLON {
		return "", nil
	}
	p.pos++
	typeTok, err := p.expect(IDENT)
	if err != nil {
		return "", err
	}
	if !validTypes[typeTok.Value] {
		return "", at(typeTok.Line, fmt.Errorf("unknown type %q (the types are: number, text, bool, list, tensor, any)", typeTok.Value))
	}
	return typeTok.Value, nil
}

// parseBlock reads statements until the closing }. The opening { is already consumed.
func (p *Parser) parseBlock() ([]Statement, error) {
	var stmts []Statement
	for {
		tok := p.peek()
		if tok == nil {
			return nil, at(p.lastLine(), fmt.Errorf("unexpected end of file, missing }"))
		}
		if tok.Type == RBRACE {
			p.pos++
			return stmts, nil
		}
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
	}
}

// parseFunc reads the rest of a function, after the word func (or tool):
// name(a: number, b): number { ... }
func (p *Parser) parseFunc(line int) (Statement, error) {
	fn, err := p.parseHeader(line)
	if err != nil {
		return nil, err
	}
	return p.parseBody(fn)
}

// parseHeader reads the name, the parameters and the return type.
func (p *Parser) parseHeader(line int) (FuncStatement, error) {
	none := FuncStatement{}
	fname, err := p.expect(IDENT)
	if err != nil {
		return none, err
	}
	if _, err := p.expect(LPAREN); err != nil {
		return none, err
	}
	var params []Param
	if tok := p.peek(); tok != nil && tok.Type == RPAREN { // no parameters
		p.pos++
	} else {
		for {
			param, err := p.expect(IDENT)
			if err != nil {
				return none, err
			}
			paramType, err := p.parseOptionalType()
			if err != nil {
				return none, err
			}
			params = append(params, Param{Name: param.Value, Type: paramType})
			if tok := p.peek(); tok != nil && tok.Type == COMMA {
				p.pos++
				continue
			}
			break
		}
		if _, err := p.expect(RPAREN); err != nil {
			return none, err
		}
	}
	returnType, err := p.parseOptionalType() // func add(a, b): number { ... }
	if err != nil {
		return none, err
	}
	return FuncStatement{Name: fname.Value, Params: params, ReturnType: returnType, Line: line}, nil
}

// parseBody reads the { ... } of a function.
func (p *Parser) parseBody(fn FuncStatement) (Statement, error) {
	if _, err := p.expect(LBRACE); err != nil {
		return nil, err
	}
	// a function body starts outside any loop, even if the func is written inside one
	outerLoops := p.loopDepth
	p.loopDepth = 0
	body, err := p.parseBlock()
	p.loopDepth = outerLoops
	if err != nil {
		return nil, err
	}
	fn.Body = body
	return fn, nil
}

// toolTypes are the types a tool's parameters may have: the ones an AI model can send.
var toolTypes = map[string]bool{"number": true, "text": true, "bool": true, "list": true}

// parseTool reads a tool: a function an AI model will be able to use.
//
//	tool add(a: number, b: number): number "Adds two numbers" { ... }
//
// So that a model can be told how to use it, every parameter needs a type and the tool needs
// a description. These are checked here, before the program runs.
func (p *Parser) parseTool(line int) (Statement, error) {
	fn, err := p.parseHeader(line)
	if err != nil {
		return nil, err
	}
	wrong := func(format string, args ...any) error {
		return at(line, fmt.Errorf("tool %q"+format, append([]any{fn.Name}, args...)...))
	}
	seen := map[string]bool{}
	for _, param := range fn.Params {
		if param.Type == "" {
			return nil, wrong(": parameter %q needs a type, so a model knows what to send (number, text, bool or list)", param.Name)
		}
		if !toolTypes[param.Type] {
			return nil, wrong(": parameter %q cannot be %s; a model can send number, text, bool or list", param.Name, param.Type)
		}
		if seen[param.Name] { // a model names each value it sends, so names must be different
			return nil, wrong(" has two parameters named %q", param.Name)
		}
		seen[param.Name] = true
	}
	description := p.peek()
	if description != nil && description.Type == ILLEGAL {
		return nil, p.unexpected(description, "a description") // e.g. a missing closing quote
	}
	if description == nil || description.Type != STRING {
		return nil, wrong(" needs a description in quotes before the {, to tell a model what it does")
	}
	p.pos++
	if strings.TrimSpace(description.Value) == "" {
		return nil, wrong(" has an empty description")
	}
	fn.IsTool, fn.Description = true, description.Value
	return p.parseBody(fn)
}

// parseStatement reads one command: print, let, if, while, func, return, an assignment or a call.
func (p *Parser) parseStatement() (Statement, error) {
	name, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}
	line := name.Line // the statement is on the line where its first word is

	// "tool" starts a tool only when a name follows it, so a variable called tool still works
	if next := p.peek(); name.Value == "tool" && next != nil && next.Type == IDENT {
		return p.parseTool(line)
	}

	switch name.Value {
	case "print": // print(value)
		if _, err := p.expect(LPAREN); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(RPAREN); err != nil {
			return nil, err
		}
		return PrintStatement{Value: value, Line: line}, nil

	case "let", "var", "final": // let name = value   (var is the same as let; final cannot be changed)
		varName, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}
		varType, err := p.parseOptionalType() // let x: number = 5
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(EQUALS); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return LetStatement{Name: varName.Value, Value: value, Final: name.Value == "final", Type: varType, Line: line}, nil

	case "if": // if condition { ... } else { ... }
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(LBRACE); err != nil {
			return nil, err
		}
		thenBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		var elseBlock []Statement
		// the else part is optional
		if isWord(p.peek(), "else") {
			p.pos++
			if isWord(p.peek(), "if") {
				// "else if": the else part is just another if statement
				nested, err := p.parseStatement()
				if err != nil {
					return nil, err
				}
				elseBlock = []Statement{nested}
			} else {
				if _, err := p.expect(LBRACE); err != nil {
					return nil, err
				}
				elseBlock, err = p.parseBlock()
				if err != nil {
					return nil, err
				}
			}
		}
		return IfStatement{Condition: cond, Then: thenBlock, Else: elseBlock, Line: line}, nil

	case "while": // while condition { ... }
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(LBRACE); err != nil {
			return nil, err
		}
		body, err := p.parseLoopBody()
		if err != nil {
			return nil, err
		}
		return WhileStatement{Condition: cond, Body: body, Line: line}, nil

	case "for": // for i in 1 to 5 { ... }
		loopVar, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}
		if tok := p.peek(); !isWord(tok, "in") {
			return nil, p.unexpected(tok, "in")
		}
		p.pos++
		first, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		var last Expr
		if isWord(p.peek(), "to") {
			p.pos++
			last, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(LBRACE); err != nil {
			return nil, err
		}
		body, err := p.parseLoopBody()
		if err != nil {
			return nil, err
		}
		if last == nil { // for x in somethingToLoopOver { ... }
			return ForStatement{Var: loopVar.Value, Iter: first, Body: body, Line: line}, nil
		}
		return ForStatement{Var: loopVar.Value, From: first, To: last, Body: body, Line: line}, nil

	case "func": // func name(a, b) { ... }   (types are optional: func add(a: number): number { ... })
		return p.parseFunc(line)

	case "import": // import "name"
		path, err := p.expect(STRING)
		if err != nil {
			return nil, err
		}
		return ImportStatement{Path: path.Value, Line: line}, nil

	case "break", "continue": // only allowed inside a loop
		if p.loopDepth == 0 {
			return nil, at(line, fmt.Errorf("%q can only be used inside a loop", name.Value))
		}
		if name.Value == "break" {
			return BreakStatement{Line: line}, nil
		}
		return ContinueStatement{Line: line}, nil

	case "return": // return value
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return ReturnStatement{Value: value, Line: line}, nil
	}

	// Not a keyword. If a "=" follows, it is an assignment: name = value
	if tok := p.peek(); tok != nil && tok.Type == EQUALS {
		p.pos++
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return AssignStatement{Name: name.Value, Value: value, Line: line}, nil
	}

	// A name followed by . on its own line is a call into a module: math.say("hi")
	if tok := p.peek(); tok != nil && tok.Type == DOT {
		p.pos-- // step back to the name and read the whole thing as an expression
		expr, err := p.parsePostfix()
		if err != nil {
			return nil, err
		}
		call, ok := expr.(CallExpr)
		if !ok {
			return nil, at(line, fmt.Errorf("%s is a value, not a command; to run a function write module.name(...)", name.Value))
		}
		return ExprStatement{Call: call, Line: line}, nil
	}

	// A name followed by [ is a change to a list item: names[0] = "Ann"
	if tok := p.peek(); tok != nil && tok.Type == LBRACKET {
		var target Expr = VarRef{name.Value}
		var index Expr
		for {
			p.pos++ // the [
			idx, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(RBRACKET); err != nil {
				return nil, err
			}
			if next := p.peek(); next != nil && next.Type == LBRACKET {
				target = IndexExpr{Target: target, Index: idx} // grid[0][1]: go one level deeper
				continue
			}
			index = idx
			break
		}
		if _, err := p.expect(EQUALS); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return IndexAssignStatement{Target: target, Index: index, Value: value, Line: line}, nil
	}

	// A name followed by ( on its own line is a function call: greet("GAN")
	if tok := p.peek(); tok != nil && tok.Type == LPAREN {
		call, err := p.parseCall(name.Value)
		if err != nil {
			return nil, err
		}
		return ExprStatement{Call: call, Line: line}, nil
	}

	return nil, at(line, fmt.Errorf("unknown command %q", name.Value))
}

// Parse turns the tokens into a list of statements.
func Parse(tokens []Token) ([]Statement, error) {
	p := &Parser{tokens: tokens}
	var stmts []Statement
	for p.peek() != nil {
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
	}
	return stmts, nil
}
```

## interpreter.go
```go
package main

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

// Value is what a variable holds: a float64 (number), a string, a bool (true/false),
// a *List, a *Tensor (see tensor.go) or a *Module.
type Value interface{}

// List is a list of values. It is used through a pointer, so two names for the same
// list see the same changes (like Python).
type List struct {
	Items []Value
}

// itemsWord writes "1 item" or "3 items".
func itemsWord(n int) string {
	if n == 1 {
		return "1 item"
	}
	return fmt.Sprintf("%d items", n)
}

// toIndex checks that v is a whole number inside 0..length-1 and returns it.
// what describes the thing being indexed, for the error message, e.g. "the list has 3 items".
func toIndex(v Value, length int, what string) (int, error) {
	n, ok := v.(float64)
	if !ok || n != math.Trunc(n) {
		return 0, fmt.Errorf("an index must be a whole number")
	}
	i := int(n)
	if i < 0 || i >= length {
		return 0, fmt.Errorf("index %d is out of range (%s)", i, what)
	}
	return i, nil
}

// equalValues is == for GAN values. Lists are equal when all their items are equal.
func equalValues(a, b Value) bool {
	la, aIsList := a.(*List)
	lb, bIsList := b.(*List)
	if aIsList && bIsList {
		if len(la.Items) != len(lb.Items) {
			return false
		}
		for i := range la.Items {
			if !equalValues(la.Items[i], lb.Items[i]) {
				return false
			}
		}
		return true
	}
	// tensors are equal when their shapes and numbers are the same
	ta, aIsTensor := a.(*Tensor)
	tb, bIsTensor := b.(*Tensor)
	if aIsTensor && bIsTensor {
		return ta.equal(tb)
	}
	return a == b
}

// binding is one variable: its value, and whether it is final (cannot change).
type binding struct {
	value Value
	final bool
	typ   string // "number", "text", "bool", "list", "tensor" or "any"; "" means no type was given
}

// typeName says what kind of value v is.
func typeName(v Value) string {
	switch v.(type) {
	case float64:
		return "number"
	case string:
		return "text"
	case bool:
		return "bool"
	case *List:
		return "list"
	case *Tensor:
		return "tensor"
	case *Module:
		return "module"
	}
	return "nothing"
}

// checkType makes sure v is the type that was asked for. what names the thing being
// checked, for the error message. A want of "" or "any" accepts everything.
func checkType(what, want string, v Value) error {
	if want == "" || want == "any" || typeName(v) == want {
		return nil
	}
	return fmt.Errorf("%s must be %s but got %s", what, want, typeName(v))
}

// Environment is one scope: a table of variable names and their values.
// A scope can have a parent, the scope around it. Looking a name up checks this scope
// first, then the parent, and so on outwards, so inner code can see outer variables.
type Environment struct {
	vars   map[string]*binding
	parent *Environment
}

// NewEnv makes an empty scope inside parent (parent is nil for the outermost scope).
func NewEnv(parent *Environment) *Environment {
	return &Environment{vars: map[string]*binding{}, parent: parent}
}

// lookup finds a variable in this scope or any scope around it.
func (e *Environment) lookup(name string) (*binding, bool) {
	for scope := e; scope != nil; scope = scope.parent {
		if b, ok := scope.vars[name]; ok {
			return b, true
		}
	}
	return nil, false
}

// define creates a variable in THIS scope (what let, var, final and parameters do).
// A variable with the same name in an outer scope is hidden, not changed.
func (e *Environment) define(name string, value Value, final bool, typ string) error {
	if b, ok := e.vars[name]; ok && b.final {
		return fmt.Errorf("cannot change final variable %q", name)
	}
	if err := checkType(fmt.Sprintf("variable %q", name), typ, value); err != nil {
		return err
	}
	e.vars[name] = &binding{value: value, final: final, typ: typ}
	return nil
}

// assign changes an existing variable, wherever it lives (what name = value does).
func (e *Environment) assign(name string, value Value) error {
	b, ok := e.lookup(name)
	if !ok {
		return fmt.Errorf("undefined variable %q", name)
	}
	if b.final {
		return fmt.Errorf("cannot change final variable %q", name)
	}
	if err := checkType(fmt.Sprintf("variable %q", name), b.typ, value); err != nil {
		return err
	}
	b.value = value
	return nil
}

// maxDepth stops a function that calls itself forever from crashing the program.
const maxDepth = 1000

// Interpreter remembers the things that live for the whole run.
type Interpreter struct {
	cur     *Module            // the module whose code is running right now (the program itself, or a library)
	modules map[string]*Module // libraries already loaded, by file path
	loading map[string]bool    // libraries being loaded right now (to catch circular imports)
	mainDir string             // the folder of the program being run
	depth   int                // how deeply function calls are nested
}

// plural writes "1 argument" or "2 arguments".
func plural(n int) string {
	if n == 1 {
		return "1 argument"
	}
	return fmt.Sprintf("%d arguments", n)
}

// callFunction runs a function. It returns the value, and whether the function returned one.
func (in *Interpreter) callFunction(call CallExpr, env *Environment) (Value, bool, error) {
	// Which module are we calling into? math.square(3) calls into the math module;
	// a plain square(3) calls into the module whose code is running now.
	mod := in.cur
	if call.Target != nil {
		t, err := in.eval(call.Target, env)
		if err != nil {
			return nil, false, err
		}
		m, ok := t.(*Module)
		if !ok {
			return nil, false, fmt.Errorf("only a module can be followed by .%s(...); did you forget an import?", call.Name)
		}
		mod = m
	}
	fn, ok := mod.funcs[call.Name]
	if !ok {
		if call.Target != nil {
			return nil, false, fmt.Errorf("module %q has no function %q", mod.name, call.Name)
		}
		// not a function you wrote: maybe one that comes with GAN, like len or sqrt
		return in.callBuiltin(call, env)
	}
	if len(call.Args) != len(fn.Params) {
		return nil, false, fmt.Errorf("function %q expects %s but got %d", call.Name, plural(len(fn.Params)), len(call.Args))
	}
	if in.depth >= maxDepth {
		return nil, false, fmt.Errorf("too many nested calls")
	}

	// The function gets its own scope, inside its module's global scope: it can see the
	// module's global variables, and its parameters and own variables disappear when it ends.
	// The arguments are worked out in the caller's scope.
	local := NewEnv(mod.env)
	for i, param := range fn.Params {
		v, err := in.eval(call.Args[i], env)
		if err != nil {
			return nil, false, err
		}
		if err := checkType(fmt.Sprintf("parameter %q of %q", param.Name, call.Name), param.Type, v); err != nil {
			return nil, false, err
		}
		local.define(param.Name, v, false, param.Type)
	}

	in.depth++
	previous := in.cur
	in.cur = mod // while the function runs, its own module is "the current module"
	value, returned, err := in.runBlock(fn.Body, local)
	in.cur = previous
	in.depth--
	if err == nil && fn.ReturnType != "" && fn.ReturnType != "any" {
		if !returned {
			return nil, false, fmt.Errorf("function %q must return %s but returned nothing", call.Name, fn.ReturnType)
		}
		if err := checkType(fmt.Sprintf("the value returned by %q", call.Name), fn.ReturnType, value); err != nil {
			return nil, false, err
		}
	}
	return value, returned, err
}

// eval works out the value of an expression.
func (in *Interpreter) eval(e Expr, env *Environment) (Value, error) {
	switch x := e.(type) {
	case StringLit:
		return x.Value, nil
	case NumberLit:
		return x.Value, nil
	case BoolLit:
		return x.Value, nil
	case MemberExpr:
		t, err := in.eval(x.Target, env)
		if err != nil {
			return nil, err
		}
		m, ok := t.(*Module)
		if !ok {
			return nil, fmt.Errorf("only a module can be followed by .%s", x.Name)
		}
		b, ok := m.env.vars[x.Name]
		if !ok {
			return nil, fmt.Errorf("module %q has no variable %q", m.name, x.Name)
		}
		return b.value, nil
	case ListLit:
		list := &List{Items: []Value{}}
		for _, item := range x.Items {
			v, err := in.eval(item, env)
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, v)
		}
		return list, nil
	case IndexExpr:
		target, err := in.eval(x.Target, env)
		if err != nil {
			return nil, err
		}
		index, err := in.eval(x.Index, env)
		if err != nil {
			return nil, err
		}
		switch t := target.(type) {
		case *List:
			i, err := toIndex(index, len(t.Items), "the list has "+itemsWord(len(t.Items)))
			if err != nil {
				return nil, err
			}
			return t.Items[i], nil
		case string:
			chars := []rune(t)
			i, err := toIndex(index, len(chars), fmt.Sprintf("the text has %d characters", len(chars)))
			if err != nil {
				return nil, err
			}
			return string(chars[i]), nil
		case *Tensor:
			if t.tracked() {
				return nil, noGradient("[ ]")
			}
			if t.isScalar() {
				return nil, fmt.Errorf("this tensor is a single number; use item(...) to get it")
			}
			i, err := toIndex(index, t.Shape[0], "the tensor has "+itemsWord(t.Shape[0])+" in its first dimension")
			if err != nil {
				return nil, err
			}
			return t.item(i), nil
		}
		return nil, fmt.Errorf("only lists, text and tensors can be indexed")
	case LogicExpr:
		left, err := in.eval(x.Left, env)
		if err != nil {
			return nil, err
		}
		l, ok := left.(bool)
		if !ok {
			return nil, fmt.Errorf("%s needs true or false on both sides", x.Op)
		}
		// stop early when the answer is already known (so the right side may not run at all)
		if x.Op == "or" && l {
			return true, nil
		}
		if x.Op == "and" && !l {
			return false, nil
		}
		right, err := in.eval(x.Right, env)
		if err != nil {
			return nil, err
		}
		r, ok := right.(bool)
		if !ok {
			return nil, fmt.Errorf("%s needs true or false on both sides", x.Op)
		}
		return r, nil
	case VarRef:
		b, ok := env.lookup(x.Name)
		if !ok {
			return nil, fmt.Errorf("undefined variable %q", x.Name)
		}
		return b.value, nil
	case CallExpr:
		v, returned, err := in.callFunction(x, env)
		if err != nil {
			return nil, err
		}
		if !returned {
			return nil, fmt.Errorf("function %q did not return a value", x.Name)
		}
		return v, nil
	case UnaryExpr:
		v, err := in.eval(x.Operand, env)
		if err != nil {
			return nil, err
		}
		if x.Op == "not" {
			b, ok := v.(bool)
			if !ok {
				return nil, fmt.Errorf("not needs true or false")
			}
			return !b, nil
		}
		if t, isTensor := v.(*Tensor); isTensor {
			return tensorMath("*", t, -1.0) // -t flips the sign of every number
		}
		n, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("can only negate numbers")
		}
		return -n, nil
	case BinaryExpr:
		// work out both sides first, then apply the operator
		left, err := in.eval(x.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := in.eval(x.Right, env)
		if err != nil {
			return nil, err
		}

		// == and != work on numbers or strings, as long as both sides are the same type
		if x.Op == "==" || x.Op == "!=" {
			if fmt.Sprintf("%T", left) != fmt.Sprintf("%T", right) {
				return nil, fmt.Errorf("cannot compare different types")
			}
			if x.Op == "==" {
				return equalValues(left, right), nil
			}
			return !equalValues(left, right), nil
		}

		// + puts two lists together into a new list
		if x.Op == "+" {
			l, lok := left.(*List)
			r, rok := right.(*List)
			if lok && rok {
				joined := append(append([]Value{}, l.Items...), r.Items...)
				return &List{Items: joined}, nil
			}
		}

		// + joins text when either side is text: "Hello, " + name, "Total: " + 5
		if x.Op == "+" {
			_, leftIsText := left.(string)
			_, rightIsText := right.(string)
			if leftIsText || rightIsText {
				return show(left) + show(right), nil
			}
		}

		// math with a tensor on either side works position by position (see tensor.go)
		_, leftIsTensor := left.(*Tensor)
		_, rightIsTensor := right.(*Tensor)
		isMath := x.Op == "+" || x.Op == "-" || x.Op == "*" || x.Op == "/" || x.Op == "%"
		if isMath && (leftIsTensor || rightIsTensor) {
			return tensorMath(x.Op, left, right)
		}

		// everything else works on numbers only
		l, lok := left.(float64)
		r, rok := right.(float64)
		if !lok || !rok {
			if x.Op == "<" || x.Op == ">" || x.Op == "<=" || x.Op == ">=" {
				return nil, fmt.Errorf("can only compare numbers with %s", x.Op)
			}
			return nil, fmt.Errorf("can only do math on numbers")
		}
		switch x.Op {
		case "+":
			return l + r, nil
		case "-":
			return l - r, nil
		case "*":
			return l * r, nil
		case "/":
			if r == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return l / r, nil
		case "%":
			if r == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			return math.Mod(l, r), nil
		case "<":
			return l < r, nil
		case ">":
			return l > r, nil
		case "<=":
			return l <= r, nil
		case ">=":
			return l >= r, nil
		}
	}
	return nil, fmt.Errorf("unknown expression")
}

// show turns a value into text for printing.
func show(v Value) string {
	if m, ok := v.(*Module); ok {
		return "<module " + m.name + ">"
	}
	if list, ok := v.(*List); ok {
		parts := make([]string, len(list.Items))
		for i, item := range list.Items {
			if s, isText := item.(string); isText {
				parts[i] = "\"" + s + "\"" // text inside a list is shown in quotes
			} else {
				parts[i] = show(item)
			}
		}
		return "[" + strings.Join(parts, ", ") + "]"
	}
	if n, ok := v.(float64); ok {
		if n == 0 {
			n = 0 // turns a "negative zero" (from -0) into a plain 0
		}
		// 'f' = plain decimals, -1 = as few digits as needed:
		// 5 prints as 5, 2.5 as 2.5, 1000000 as 1000000 (not 1e+06)
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return fmt.Sprint(v)
}

// loopSignal is how break and continue travel out of nested blocks. They use the same road
// as return (returned = true), but the nearest loop catches them. The parser makes sure they
// only appear inside a loop, so a signal can never escape a function or the program.
type loopSignal string

const (
	breakSignal    loopSignal = "break"
	continueSignal loopSignal = "continue"
)

// runBlock runs a list of statements in order. It stops at the first error.
// The whole program, the inside of an if or while, and function bodies all use it.
// If a "return" happens, it stops early and reports returned = true with the value.
func (in *Interpreter) runBlock(stmts []Statement, env *Environment) (Value, bool, error) {
	for _, s := range stmts {
		val, returned, err := in.runStatement(s, env)
		if err != nil {
			// tag the error with this statement's line (unless a deeper statement already did)
			return nil, false, locate(err, lineOf(s), in.cur.file)
		}
		if returned {
			return val, true, nil
		}
	}
	return nil, false, nil
}

// runStatement does what one statement says.
func (in *Interpreter) runStatement(s Statement, env *Environment) (Value, bool, error) {
	switch st := s.(type) {
	case LetStatement:
		v, err := in.eval(st.Value, env)
		if err != nil {
			return nil, false, err
		}
		if err := env.define(st.Name, v, st.Final, st.Type); err != nil { // remember the variable
			return nil, false, err
		}
	case AssignStatement:
		// unlike let, the variable must already exist (assign checks that, and final)
		v, err := in.eval(st.Value, env)
		if err != nil {
			return nil, false, err
		}
		if err := env.assign(st.Name, v); err != nil {
			return nil, false, err
		}
	case PrintStatement:
		v, err := in.eval(st.Value, env)
		if err != nil {
			return nil, false, err
		}
		fmt.Println(show(v))
	case FuncStatement:
		in.cur.funcs[st.Name] = st // remember the function; it runs only when called
	case ImportStatement:
		if err := in.runImport(st, env); err != nil {
			return nil, false, err
		}
	case ReturnStatement:
		if in.depth == 0 {
			return nil, false, fmt.Errorf("return outside a function")
		}
		v, err := in.eval(st.Value, env)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	case BreakStatement:
		return breakSignal, true, nil
	case ContinueStatement:
		return continueSignal, true, nil
	case ExprStatement:
		// a call on its own line: run it and ignore any value it returns
		if _, _, err := in.callFunction(st.Call, env); err != nil {
			return nil, false, err
		}
	case IfStatement:
		v, err := in.eval(st.Condition, env)
		if err != nil {
			return nil, false, err
		}
		cond, ok := v.(bool)
		if !ok {
			return nil, false, fmt.Errorf("condition must be true or false")
		}
		// run the matching block (Else is empty if there was no else)
		block := st.Else
		if cond {
			block = st.Then
		}
		return in.runBlock(block, NewEnv(env)) // the block gets its own scope; errors and returns pass up
	case WhileStatement:
		for {
			// check the condition before every pass
			v, err := in.eval(st.Condition, env)
			if err != nil {
				return nil, false, err
			}
			cond, ok := v.(bool)
			if !ok {
				return nil, false, fmt.Errorf("condition must be true or false")
			}
			if !cond {
				break
			}
			// every pass gets a fresh scope, so a let inside the loop starts new each time
			val, returned, err := in.runBlock(st.Body, NewEnv(env))
			if err != nil {
				return nil, false, err
			}
			if returned {
				if val == breakSignal {
					break // leave this loop
				}
				if val != continueSignal {
					return val, true, nil // a real return: pass it up
				}
				// continue: nothing to do, go on to the next pass
			}
		}
	case ForStatement:
		return in.runFor(st, env)
	case IndexAssignStatement:
		target, err := in.eval(st.Target, env)
		if err != nil {
			return nil, false, err
		}
		index, err := in.eval(st.Index, env)
		if err != nil {
			return nil, false, err
		}
		value, err := in.eval(st.Value, env)
		if err != nil {
			return nil, false, err
		}
		list, ok := target.(*List)
		if !ok {
			if _, isText := target.(string); isText {
				return nil, false, fmt.Errorf("text cannot be changed; make a new piece of text instead")
			}
			if _, isTensor := target.(*Tensor); isTensor {
				return nil, false, fmt.Errorf("a tensor cannot be changed; make a new one instead")
			}
			return nil, false, fmt.Errorf("only list items can be changed")
		}
		i, err := toIndex(index, len(list.Items), "the list has "+itemsWord(len(list.Items)))
		if err != nil {
			return nil, false, err
		}
		list.Items[i] = value
	}
	return nil, false, nil
}

// runFor runs a for loop: for i in 1 to 5 { ... }
func (in *Interpreter) runFor(st ForStatement, env *Environment) (Value, bool, error) {
	if st.To == nil {
		// for x in something: go through the items of a list, or the characters of text
		iterable, err := in.eval(st.Iter, env)
		if err != nil {
			return nil, false, err
		}
		var items []Value
		switch it := iterable.(type) {
		case *List:
			items = append(items, it.Items...) // a copy, so changing the list inside the loop is safe
		case string:
			for _, r := range it {
				items = append(items, string(r))
			}
		default:
			return nil, false, fmt.Errorf("for can only loop over a range (1 to 5), a list or text")
		}
		for _, item := range items {
			scope := NewEnv(env)
			scope.define(st.Var, item, false, "")
			val, returned, err := in.runBlock(st.Body, scope)
			if err != nil {
				return nil, false, err
			}
			if returned {
				if val == breakSignal {
					break
				}
				if val != continueSignal {
					return val, true, nil
				}
			}
		}
		return nil, false, nil
	}
	fromV, err := in.eval(st.From, env)
	if err != nil {
		return nil, false, err
	}
	toV, err := in.eval(st.To, env)
	if err != nil {
		return nil, false, err
	}
	from, ok1 := fromV.(float64)
	to, ok2 := toV.(float64)
	if !ok1 || !ok2 {
		return nil, false, fmt.Errorf("for needs numbers on both sides of to")
	}
	// count from..to, both included. If from is bigger than to, the body never runs.
	for i := from; i <= to; i++ {
		scope := NewEnv(env) // each pass has its own scope holding the loop variable
		scope.define(st.Var, i, false, "")
		val, returned, err := in.runBlock(st.Body, scope)
		if err != nil {
			return nil, false, err
		}
		if returned {
			if val == breakSignal {
				break
			}
			if val != continueSignal {
				return val, true, nil
			}
		}
	}
	return nil, false, nil
}

// Run executes the whole program (found at mainPath) with a fresh, empty environment.
func Run(stmts []Statement, mainPath string) error {
	abs, err := filepath.Abs(mainPath)
	if err != nil {
		abs = mainPath
	}
	main := newModule("main", "") // the program itself is a module with no file name
	in := &Interpreter{
		cur:     main,
		modules: map[string]*Module{},
		loading: map[string]bool{},
		mainDir: filepath.Dir(abs),
	}
	_, _, err = in.runBlock(stmts, main.env)
	return err
}
```

## builtins.go
```go
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// A builtin is a function that comes with GAN itself, written in Go instead of GAN.
// min and max say how many arguments it takes (max -1 means "any number").
// run returns the value, whether there is a value at all, and an error.
type builtin struct {
	min, max int
	run      func(args []Value) (Value, bool, error)
}

// stdin is where input() reads from.
var stdin = bufio.NewReader(os.Stdin)

// stdinIsTerminal says whether a person is typing at a keyboard (as opposed to input piped in).
func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// needNumber checks that an argument is a number.
func needNumber(fn string, v Value) (float64, error) {
	n, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%s needs a number", fn)
	}
	return n, nil
}

// needText checks that an argument is text.
func needText(fn string, v Value) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s needs text", fn)
	}
	return s, nil
}

// builtins is the table of all the functions GAN comes with.
var builtins = map[string]builtin{
	// len(list) or len(text): how many items or characters
	"len": {1, 1, func(a []Value) (Value, bool, error) {
		switch v := a[0].(type) {
		case *List:
			return float64(len(v.Items)), true, nil
		case string:
			return float64(utf8.RuneCountInString(v)), true, nil
		}
		return nil, false, fmt.Errorf("len needs a list or text")
	}},

	// ---- Python (the code that does the work is in python.go) ----

	// python("module.function", arguments...): call any Python function
	"python": {1, -1, callPython},

	// ---- AI models (the code that does the work is in ai.go) ----

	// ask(question) or ask(question, model): the model's answer, as text
	"ask": {1, 2, func(a []Value) (Value, bool, error) {
		question, err := askArgument(a[0], "the question")
		if err != nil {
			return nil, false, err
		}
		model := defaultModel
		if len(a) == 2 {
			if model, err = askArgument(a[1], "the model name"); err != nil {
				return nil, false, err
			}
		}
		answer, err := askModel(question, model)
		if err != nil {
			return nil, false, err
		}
		return answer, true, nil
	}},

	// ---- tensors (the code that does the work is in tensor.go) ----

	// tensor(list): make a tensor from a (nested) list of numbers
	"tensor": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := newTensor(a[0])
		if err != nil {
			return nil, false, err
		}
		return t, true, nil
	}},

	// shape(t): the size of each dimension, as a list
	"shape": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("shape", a[0])
		if err != nil {
			return nil, false, err
		}
		list := &List{Items: []Value{}}
		for _, n := range t.Shape {
			list.Items = append(list.Items, float64(n))
		}
		return list, true, nil
	}},

	// sum(t): all the numbers of a tensor added up
	"sum": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("sum", a[0])
		if err != nil {
			return nil, false, err
		}
		return total(t, 1), true, nil
	}},

	// matmul(a, b): matrix multiplication
	"matmul": {2, 2, func(a []Value) (Value, bool, error) {
		x, err := needTensor("matmul", a[0])
		if err != nil {
			return nil, false, err
		}
		y, err := needTensor("matmul", a[1])
		if err != nil {
			return nil, false, err
		}
		result, err := matmul(x, y)
		if err != nil {
			return nil, false, err
		}
		return result, true, nil
	}},

	// new tensors of a given shape, e.g. zeros([2, 3])
	"zeros": shapeBuiltin("zeros", func() float64 { return 0 }),
	"ones":  shapeBuiltin("ones", func() float64 { return 1 }),
	// random: from 0 up to (not including) 1. randomNormal: centred on 0, mostly close to it.
	"random":       shapeBuiltin("random", func() float64 { return random.Float64() }),
	"randomNormal": shapeBuiltin("randomNormal", func() float64 { return random.NormFloat64() }),

	// seed(n): make the random numbers repeat (gives no value back)
	"seed": {1, 1, func(a []Value) (Value, bool, error) {
		return nil, false, setSeed(a[0])
	}},

	// mean(t): the average of all the numbers
	"mean": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("mean", a[0])
		if err != nil {
			return nil, false, err
		}
		return total(t, float64(len(t.Data))), true, nil // a tensor is never empty
	}},

	// ---- automatic gradients (the code that does the work is in autograd.go) ----

	// learnable(t): a tensor whose gradient we want
	"learnable": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("learnable", a[0])
		if err != nil {
			return nil, false, err
		}
		return learnable(t), true, nil
	}},

	// backward(loss): work out the gradients (gives no value back)
	"backward": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("backward", a[0])
		if err != nil {
			return nil, false, err
		}
		return nil, false, backward(t)
	}},

	// grad(w): the gradient of w found by the last backward
	"grad": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("grad", a[0])
		if err != nil {
			return nil, false, err
		}
		g, err := gradientOf(t)
		if err != nil {
			return nil, false, err
		}
		return g, true, nil
	}},

	// item(t): the number inside a tensor that holds just one
	"item": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("item", a[0])
		if err != nil {
			return nil, false, err
		}
		if len(t.Data) != 1 {
			return nil, false, fmt.Errorf("item needs a tensor with exactly one number, but got shape %s", shapeText(t.Shape))
		}
		return t.Data[0], true, nil
	}},

	// transpose(t): swap rows and columns
	"transpose": {1, 1, func(a []Value) (Value, bool, error) {
		t, err := needTensor("transpose", a[0])
		if err != nil {
			return nil, false, err
		}
		result, err := transpose(t)
		if err != nil {
			return nil, false, err
		}
		return result, true, nil
	}},

	// these work on one number or on every number of a tensor
	"relu":    numberBuiltin("relu", relu, reluSlope, nil),
	"sigmoid": numberBuiltin("sigmoid", sigmoid, sigmoidSlope, nil),
	"tanh":    numberBuiltin("tanh", math.Tanh, tanhSlope, nil),
	"exp":     numberBuiltin("exp", math.Exp, expSlope, nil),
	"log":     numberBuiltin("log", math.Log, logSlope, positive),

	// push(list, item): add an item to the end of a list (gives no value back)
	"push": {2, 2, func(a []Value) (Value, bool, error) {
		list, ok := a[0].(*List)
		if !ok {
			return nil, false, fmt.Errorf("push needs a list as its first argument")
		}
		list.Items = append(list.Items, a[1])
		return nil, false, nil
	}},

	// pop(list): remove the last item of a list and give it back
	"pop": {1, 1, func(a []Value) (Value, bool, error) {
		list, ok := a[0].(*List)
		if !ok {
			return nil, false, fmt.Errorf("pop needs a list")
		}
		if len(list.Items) == 0 {
			return nil, false, fmt.Errorf("cannot pop from an empty list")
		}
		last := list.Items[len(list.Items)-1]
		list.Items = list.Items[:len(list.Items)-1]
		return last, true, nil
	}},

	// fail("message"): stop the program with your own error message
	"fail": {1, 1, func(a []Value) (Value, bool, error) {
		return nil, false, fmt.Errorf("%s", show(a[0]))
	}},

	// typeof(x): "number", "text", "bool", "list", "tensor" or "module"
	"typeof": {1, 1, func(a []Value) (Value, bool, error) {
		return typeName(a[0]), true, nil
	}},

	// str(x): any value as text
	"str": {1, 1, func(a []Value) (Value, bool, error) {
		return show(a[0]), true, nil
	}},

	// num(text): text as a number, e.g. num("3.5") is 3.5
	"num": {1, 1, func(a []Value) (Value, bool, error) {
		if n, ok := a[0].(float64); ok {
			return n, true, nil
		}
		s, ok := a[0].(string)
		if !ok {
			return nil, false, fmt.Errorf("num needs text or a number")
		}
		n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return nil, false, fmt.Errorf("cannot turn %q into a number", s)
		}
		return n, true, nil
	}},

	// input() or input("question "): read one line typed by the user, as text
	"input": {0, 1, func(a []Value) (Value, bool, error) {
		if len(a) == 1 {
			fmt.Print(show(a[0]))
		}
		line, err := stdin.ReadString('\n') // at the end of the input this gives ""
		line = strings.TrimRight(line, "\r\n")
		// When the input is piped in (not typed at a keyboard) nothing shows the "typing", so
		// echo it. That way a recorded run reads exactly like a real one: question, then answer.
		if (err == nil || line != "") && !stdinIsTerminal() {
			fmt.Println(line)
		}
		return line, true, nil
	}},

	"abs": {1, 1, func(a []Value) (Value, bool, error) {
		n, err := needNumber("abs", a[0])
		return math.Abs(n), err == nil, err
	}},
	"floor": {1, 1, func(a []Value) (Value, bool, error) {
		n, err := needNumber("floor", a[0])
		return math.Floor(n), err == nil, err
	}},
	"round": {1, 1, func(a []Value) (Value, bool, error) {
		n, err := needNumber("round", a[0])
		return math.Round(n), err == nil, err
	}},
	"sqrt": {1, 1, func(a []Value) (Value, bool, error) {
		n, err := needNumber("sqrt", a[0])
		if err != nil {
			return nil, false, err
		}
		if n < 0 {
			return nil, false, fmt.Errorf("sqrt needs a number that is not negative")
		}
		return math.Sqrt(n), true, nil
	}},
	"pow": {2, 2, func(a []Value) (Value, bool, error) {
		base, err := needNumber("pow", a[0])
		if err != nil {
			return nil, false, err
		}
		exp, err := needNumber("pow", a[1])
		if err != nil {
			return nil, false, err
		}
		return math.Pow(base, exp), true, nil
	}},
	"min": {2, 2, func(a []Value) (Value, bool, error) {
		x, err := needNumber("min", a[0])
		if err != nil {
			return nil, false, err
		}
		y, err := needNumber("min", a[1])
		if err != nil {
			return nil, false, err
		}
		return math.Min(x, y), true, nil
	}},
	"max": {2, 2, func(a []Value) (Value, bool, error) {
		x, err := needNumber("max", a[0])
		if err != nil {
			return nil, false, err
		}
		y, err := needNumber("max", a[1])
		if err != nil {
			return nil, false, err
		}
		return math.Max(x, y), true, nil
	}},

	"upper": {1, 1, func(a []Value) (Value, bool, error) {
		s, err := needText("upper", a[0])
		return strings.ToUpper(s), err == nil, err
	}},
	"lower": {1, 1, func(a []Value) (Value, bool, error) {
		s, err := needText("lower", a[0])
		return strings.ToLower(s), err == nil, err
	}},

	// contains(list, item) or contains(text, piece): true or false
	"contains": {2, 2, func(a []Value) (Value, bool, error) {
		switch v := a[0].(type) {
		case *List:
			for _, item := range v.Items {
				if equalValues(item, a[1]) {
					return true, true, nil
				}
			}
			return false, true, nil
		case string:
			piece, err := needText("contains", a[1])
			if err != nil {
				return nil, false, err
			}
			return strings.Contains(v, piece), true, nil
		}
		return nil, false, fmt.Errorf("contains needs a list or text")
	}},
}

// callBuiltin runs one of the functions that come with GAN.
func (in *Interpreter) callBuiltin(call CallExpr, env *Environment) (Value, bool, error) {
	b, ok := builtins[call.Name]
	if !ok {
		return nil, false, fmt.Errorf("undefined function %q", call.Name)
	}
	n := len(call.Args)
	if n < b.min || (b.max >= 0 && n > b.max) {
		if b.max < 0 { // any number of arguments, but at least min
			return nil, false, fmt.Errorf("function %q expects at least %s but got %d", call.Name, plural(b.min), n)
		}
		if b.min == b.max {
			return nil, false, fmt.Errorf("function %q expects %s but got %d", call.Name, plural(b.min), n)
		}
		return nil, false, fmt.Errorf("function %q expects %d to %d arguments but got %d", call.Name, b.min, b.max, n)
	}
	args := make([]Value, n)
	for i, a := range call.Args {
		v, err := in.eval(a, env)
		if err != nil {
			return nil, false, err
		}
		args[i] = v
	}
	return b.run(args)
}
```

## modules.go
```go
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A Module is a .gan file that has been imported. It has its own functions and its own
// global variables, kept apart from the program that imported it.
type Module struct {
	name  string                   // what you call it in your program: math.square(3)
	file  string                   // the path of the .gan file ("" for the main program)
	funcs map[string]FuncStatement // the functions it defines
	env   *Environment             // its global variables
}

func newModule(name, file string) *Module {
	return &Module{name: name, file: file, funcs: map[string]FuncStatement{}, env: NewEnv(nil)}
}

// currentDir is the folder of the file whose code is running now. Imports such as
// "modules/tools.gan" are looked up starting from here.
func (in *Interpreter) currentDir() string {
	if in.cur.file == "" {
		return in.mainDir
	}
	return filepath.Dir(in.cur.file)
}

// findLibrary turns what you wrote after "import" into the path of a real file.
//
//	import "math"             looks in the libs folder next to gan.exe, then in a libs
//	                          folder next to your program, then in your program's folder
//	import "modules/tools"    (has a slash or ends in .gan) is a path, relative to the
//	                          file that contains the import
func (in *Interpreter) findLibrary(name string) (string, error) {
	file := name
	if !strings.HasSuffix(file, ".gan") {
		file += ".gan"
	}

	var places []string
	isPath := strings.ContainsAny(name, `/\`) || strings.HasSuffix(name, ".gan")
	if isPath {
		places = []string{filepath.Join(in.currentDir(), file)}
	} else {
		if exe, err := os.Executable(); err == nil {
			places = append(places, filepath.Join(filepath.Dir(exe), "libs", file))
		}
		places = append(places, filepath.Join(in.mainDir, "libs", file), filepath.Join(in.mainDir, file))
	}

	for _, p := range places {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs, nil
		}
	}
	if isPath {
		return "", fmt.Errorf("cannot find the file %q", name)
	}
	return "", fmt.Errorf("cannot find library %q (looked in: the libs folder next to gan.exe, the libs folder next to your program, and your program's folder)", name)
}

// inFile marks an error as belonging to a particular file (an imported one).
func inFile(err error, path string) error {
	var ge *GanError
	if errors.As(err, &ge) {
		return &GanError{Line: ge.Line, Msg: ge.Msg, File: path, Kind: "Syntax error"}
	}
	return err
}

// runImport loads a module (once) and gives it a name in the current scope.
func (in *Interpreter) runImport(st ImportStatement, env *Environment) error {
	path, err := in.findLibrary(st.Path)
	if err != nil {
		return err
	}

	mod, loaded := in.modules[path]
	if !loaded {
		if in.loading[path] {
			return fmt.Errorf("circular import: %q is already being loaded", st.Path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		stmts, err := Parse(Lex(string(data)))
		if err != nil {
			return inFile(err, path) // a syntax error inside the library
		}

		name := strings.TrimSuffix(filepath.Base(path), ".gan")
		mod = newModule(name, path)

		// Run the library's own code (this defines its functions and variables).
		// While it runs, "the current module" is the library, not your program.
		in.loading[path] = true
		previous := in.cur
		in.cur = mod
		_, _, err = in.runBlock(stmts, mod.env)
		in.cur = previous
		delete(in.loading, path)
		if err != nil {
			return err
		}
		in.modules[path] = mod
	}

	// importing the same module twice is harmless
	if b, ok := env.vars[mod.name]; ok && b.value == Value(mod) {
		return nil
	}
	return env.define(mod.name, mod, true, "")
}
```

## errors.go
```go
package main

import "errors"

// GanError is an error that knows which line of the .gan file it happened on.
type GanError struct {
	Line int
	Msg  string
	File string // the imported file the error is in; "" means the main program
	Kind string // "Syntax error" for a mistake found while reading an imported file; "" = decided by main
}

func (e *GanError) Error() string { return e.Msg }

// at attaches a line number to an error.
func at(line int, err error) error {
	return &GanError{Line: line, Msg: err.Error()}
}

// locate attaches a line number only if the error does not have one yet.
// This is how the innermost line wins: an error inside an if or a function
// already has its own line, so the outer statement leaves it alone.
func locate(err error, line int, file string) error {
	var ge *GanError
	if errors.As(err, &ge) {
		return err
	}
	return &GanError{Line: line, Msg: err.Error(), File: file}
}
```

## tensor.go
```go
package main

// Tensors: grids of numbers, the basic building block of AI.
//
// A tensor has a SHAPE (its size in each direction) and its numbers. The numbers are kept in
// one flat list, row after row. For example
//
//	tensor([[1, 2, 3], [4, 5, 6]])   has Shape [2, 3] and Data [1 2 3 4 5 6]
//
// A tensor never changes after it is made: every operation gives a new one.

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

// Tensor is a grid of numbers with any number of dimensions (1 or more).
type Tensor struct {
	Shape []int     // size of each dimension, e.g. [2, 3] = 2 rows and 3 columns
	Data  []float64 // all the numbers, row after row
	track *trace    // history for gradients (see autograd.go); nil for a plain tensor
}

// maxDimensions stops a list that contains itself from making GAN loop forever.
const maxDimensions = 32

// shapeText writes a shape the way GAN shows it: [2, 3].
func shapeText(shape []int) string {
	parts := make([]string, len(shape))
	for i, n := range shape {
		parts[i] = fmt.Sprint(n)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// sameShape says whether two shapes are exactly equal.
func sameShape(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// newTensor turns a (nested) list of numbers into a tensor.
func newTensor(v Value) (*Tensor, error) {
	list, ok := v.(*List)
	if !ok {
		return nil, fmt.Errorf("tensor needs a list of numbers")
	}
	t := &Tensor{}
	if err := t.fill(list, 0); err != nil {
		return nil, err
	}
	return t, nil
}

// fill reads one level of the nested list. depth says which dimension this list is.
// The first list met at each depth decides the size; all the others must agree with it.
func (t *Tensor) fill(list *List, depth int) error {
	if depth >= maxDimensions {
		return fmt.Errorf("a tensor cannot have more than %d dimensions", maxDimensions)
	}
	if len(list.Items) == 0 {
		return fmt.Errorf("a tensor cannot be empty")
	}
	if depth == len(t.Shape) {
		if len(t.Data) > 0 { // numbers were already found higher up, so this is too deep
			return unevenTensor()
		}
		t.Shape = append(t.Shape, len(list.Items))
	} else if t.Shape[depth] != len(list.Items) {
		return unevenTensor()
	}
	for _, item := range list.Items {
		switch x := item.(type) {
		case float64:
			if depth != len(t.Shape)-1 { // a number where a row was expected
				return unevenTensor()
			}
			t.Data = append(t.Data, x)
		case *List:
			if err := t.fill(x, depth+1); err != nil {
				return err
			}
		default:
			return fmt.Errorf("a tensor can only hold numbers, but got %s", typeName(item))
		}
	}
	return nil
}

func unevenTensor() error {
	return fmt.Errorf("all rows of a tensor must have the same shape")
}

// String shows a tensor as tensor([[1, 2], [3, 4]]).
func (t *Tensor) String() string {
	if t.isScalar() {
		return "tensor(" + show(t.Data[0]) + ")"
	}
	var out strings.Builder
	next := 0 // the next number in Data to write
	var write func(depth int)
	write = func(depth int) {
		out.WriteString("[")
		for i := 0; i < t.Shape[depth]; i++ {
			if i > 0 {
				out.WriteString(", ")
			}
			if depth == len(t.Shape)-1 {
				out.WriteString(show(t.Data[next])) // numbers look the same as everywhere in GAN
				next++
			} else {
				write(depth + 1)
			}
		}
		out.WriteString("]")
	}
	write(0)
	return "tensor(" + out.String() + ")"
}

// equal says whether two tensors have the same shape and the same numbers.
func (t *Tensor) equal(other *Tensor) bool {
	if !sameShape(t.Shape, other.Shape) {
		return false
	}
	for i := range t.Data {
		if t.Data[i] != other.Data[i] {
			return false
		}
	}
	return true
}

// tensorMath does + - * / where at least one side is a tensor. It works position by position:
// two tensors of the same shape, or a tensor and a single number (used for every position).
func tensorMath(op string, left, right Value) (Value, error) {
	if op != "+" && op != "-" && op != "*" && op != "/" {
		return nil, fmt.Errorf("tensors can only use + - * / (not %s)", op)
	}
	a, aIsTensor := left.(*Tensor)
	b, bIsTensor := right.(*Tensor)
	// a one-number tensor such as tensor(5) is used for every position, like a number
	if aIsTensor && bIsTensor && !a.isScalar() && !b.isScalar() && !sameShape(a.Shape, b.Shape) {
		return nil, fmt.Errorf("tensor shapes do not match for %s: %s and %s", op, shapeText(a.Shape), shapeText(b.Shape))
	}
	_, aIsNumber := left.(float64)
	_, bIsNumber := right.(float64)
	if !(aIsTensor || aIsNumber) || !(bIsTensor || bIsNumber) {
		return nil, fmt.Errorf("tensor math needs two tensors, or a tensor and a number")
	}

	like := a // the tensor that gives the result its shape
	if !aIsTensor || (a.isScalar() && bIsTensor) {
		like = b
	}
	// at gives the number a side has at position i
	at := func(v Value, i int) float64 {
		if t, ok := v.(*Tensor); ok {
			if t.isScalar() {
				return t.Data[0]
			}
			return t.Data[i]
		}
		return v.(float64)
	}
	result := &Tensor{Shape: like.Shape, Data: make([]float64, len(like.Data))}
	for i := range result.Data {
		x, y := at(left, i), at(right, i)
		switch op {
		case "+":
			result.Data[i] = x + y
		case "-":
			result.Data[i] = x - y
		case "*":
			result.Data[i] = x * y
		case "/":
			if y == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result.Data[i] = x / y
		}
	}
	// remember the history if a side is learnable (a and b are nil for plain numbers)
	recordMath(op, result, a, b,
		func(i int) float64 { return at(left, i) },
		func(i int) float64 { return at(right, i) })
	return result, nil
}

// matmul is matrix multiplication of two tables: [rows, n] times [n, cols] gives [rows, cols].
// Each answer is one row of the first multiplied with one column of the second, then added up.
func matmul(a, b *Tensor) (*Tensor, error) {
	if len(a.Shape) != 2 || len(b.Shape) != 2 {
		return nil, fmt.Errorf("matmul needs two tables (2 dimensions), but got shapes %s and %s", shapeText(a.Shape), shapeText(b.Shape))
	}
	rows, n, cols := a.Shape[0], a.Shape[1], b.Shape[1]
	if b.Shape[0] != n {
		return nil, fmt.Errorf("matmul shapes do not fit: %s and %s (the first has %d columns, so the second needs %d rows, but it has %d)",
			shapeText(a.Shape), shapeText(b.Shape), n, n, b.Shape[0])
	}
	result := &Tensor{Shape: []int{rows, cols}, Data: make([]float64, rows*cols)}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			total := 0.0
			for k := 0; k < n; k++ {
				total += a.Data[r*n+k] * b.Data[k*cols+c]
			}
			result.Data[r*cols+c] = total
		}
	}
	if a.tracked() || b.tracked() {
		recordMatmul(result, a, b)
	}
	return result, nil
}

// item reads position i of the first dimension: a number from a 1-D tensor, or else a
// smaller tensor (a row of a table). i must already be checked to be in range.
func (t *Tensor) item(i int) Value {
	if len(t.Shape) == 1 {
		return t.Data[i]
	}
	size := len(t.Data) / t.Shape[0] // how many numbers one row holds
	// sharing the numbers is safe, because a tensor never changes
	return &Tensor{Shape: t.Shape[1:], Data: t.Data[i*size : (i+1)*size]}
}

// transpose swaps the rows and columns of a table.
func transpose(t *Tensor) (*Tensor, error) {
	if len(t.Shape) != 2 {
		return nil, fmt.Errorf("transpose needs a table (2 dimensions), but got shape %s", shapeText(t.Shape))
	}
	rows, cols := t.Shape[0], t.Shape[1]
	result := &Tensor{Shape: []int{cols, rows}, Data: make([]float64, len(t.Data))}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			result.Data[c*rows+r] = t.Data[r*cols+c]
		}
	}
	if t.tracked() {
		recordTranspose(result, t)
	}
	return result, nil
}

// maxElements stops a shape such as [100000, 100000] from using up all the memory.
const maxElements = 10_000_000

// filledTensor makes a tensor of the shape given by a GAN list such as [2, 3].
// next is called once for every position to get its number.
func filledTensor(fn string, shapeValue Value, next func() float64) (*Tensor, error) {
	wrong := fmt.Errorf("%s needs a shape: a list of whole numbers of at least 1, such as [2, 3]", fn)
	list, ok := shapeValue.(*List)
	if !ok || len(list.Items) == 0 || len(list.Items) > maxDimensions {
		return nil, wrong
	}
	shape := []int{}
	total := 1
	for _, item := range list.Items {
		n, ok := item.(float64)
		if !ok || n != math.Trunc(n) || n < 1 {
			return nil, wrong
		}
		// checked before multiplying, so the total can never overflow
		if n > float64(maxElements/total) {
			return nil, fmt.Errorf("%s: that tensor is too big (the limit is %d numbers)", fn, maxElements)
		}
		shape = append(shape, int(n))
		total *= int(n)
	}
	t := &Tensor{Shape: shape, Data: make([]float64, total)}
	for i := range t.Data {
		t.Data[i] = next()
	}
	return t, nil
}

// random is where random numbers come from. It starts from an unpredictable seed;
// seed(n) restarts it, so the same numbers come out again.
var random = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

func setSeed(v Value) error {
	n, ok := v.(float64)
	if !ok || n != math.Trunc(n) || math.Abs(n) > 1e15 {
		return fmt.Errorf("seed needs a whole number")
	}
	random = rand.New(rand.NewPCG(uint64(int64(n)), 0))
	return nil
}

// everyNumber applies f to a single number, or to every number of a tensor.
// check (which may be nil) rejects numbers that f cannot work with. slope is used for gradients.
func everyNumber(fn string, v Value, f func(float64) float64, slope func(x, y float64) float64, check func(float64) error) (Value, error) {
	one := func(x float64) (float64, error) {
		if check != nil {
			if err := check(x); err != nil {
				return 0, err
			}
		}
		return f(x), nil
	}
	switch x := v.(type) {
	case float64:
		y, err := one(x)
		if err != nil {
			return nil, err
		}
		return y, nil
	case *Tensor:
		result := &Tensor{Shape: x.Shape, Data: make([]float64, len(x.Data))}
		for i, n := range x.Data {
			y, err := one(n)
			if err != nil {
				return nil, err
			}
			result.Data[i] = y
		}
		if x.tracked() {
			recordEach(result, x, slope)
		}
		return result, nil
	}
	return nil, fmt.Errorf("%s needs a number or a tensor, but got %s", fn, typeName(v))
}

// The standard functions of neural networks.

// relu keeps positive numbers and turns negative ones into 0.
func relu(x float64) float64 { return math.Max(0, x) }

// sigmoid squeezes any number into the range 0 to 1.
func sigmoid(x float64) float64 { return 1 / (1 + math.Exp(-x)) }

// positive is the check for log, which only works on numbers greater than 0.
func positive(x float64) error {
	if !(x > 0) {
		return fmt.Errorf("log needs numbers greater than 0, but got %s", show(x))
	}
	return nil
}

// numberBuiltin makes a built-in that works on a number or on every number of a tensor.
func numberBuiltin(name string, f func(float64) float64, slope func(x, y float64) float64, check func(float64) error) builtin {
	return builtin{1, 1, func(a []Value) (Value, bool, error) {
		v, err := everyNumber(name, a[0], f, slope, check)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	}}
}

// shapeBuiltin makes a built-in such as zeros, which builds a tensor from a shape.
func shapeBuiltin(name string, next func() float64) builtin {
	return builtin{1, 1, func(a []Value) (Value, bool, error) {
		t, err := filledTensor(name, a[0], next)
		if err != nil {
			return nil, false, err
		}
		return t, true, nil
	}}
}

// needTensor checks that an argument of a built-in is a tensor.
func needTensor(fn string, v Value) (*Tensor, error) {
	t, ok := v.(*Tensor)
	if !ok {
		return nil, fmt.Errorf("%s needs a tensor, but got %s (make one with tensor([...]))", fn, typeName(v))
	}
	return t, nil
}
```

## autograd.go
```go
package main

// Automatic gradients: how a model finds out which way to change its weights.
//
// A GRADIENT says how much the result (usually the error) changes when a number changes a
// tiny bit. GAN works it out like this:
//
//  1. learnable(w) marks a tensor we want gradients for.
//  2. Every calculation that uses it remembers its inputs, and how to pass a gradient
//     back to them. The result remembers this too, so a chain of history builds up.
//  3. backward(loss) starts at the end with gradient 1 and walks the history in reverse,
//     handing each input its share (this is the "chain rule").
//
// A tensor with history is called "tracked". Every operation that cannot pass gradients back
// yet must REFUSE tracked tensors (see noGradient). Ignoring them would give wrong gradients
// without any sign that something was wrong.

import "fmt"

// trace is the history a tracked tensor carries.
type trace struct {
	leaf   bool      // true for a tensor made by learnable(): the history starts here
	inputs []*Tensor // the tracked tensors this one was calculated from
	// back takes the gradient of the result and returns the gradient for each of the inputs
	back func(g []float64) [][]float64

	grad *Tensor // leaf only: the gradient found by backward
	run  int     // leaf only: which backward found it
}

// backwardRuns counts the calls to backward, so grad never hands out an out-of-date answer.
var backwardRuns int

// tracked says whether a tensor has history.
func (t *Tensor) tracked() bool { return t.track != nil }

// plain is the same tensor without any history.
func (t *Tensor) plain() *Tensor { return &Tensor{Shape: t.Shape, Data: t.Data} }

// isScalar says whether t is a single number with no dimensions, shown as tensor(5).
func (t *Tensor) isScalar() bool { return len(t.Shape) == 0 }

// learnable makes the starting point of a history. It is always a new tensor, so any
// older history is left behind.
func learnable(t *Tensor) *Tensor {
	return &Tensor{Shape: t.Shape, Data: t.Data, track: &trace{leaf: true}}
}

// noGradient is the error for an operation that cannot pass gradients back yet.
func noGradient(operation string) error {
	return fmt.Errorf("gradients for %s are not supported yet (it was used on a learnable tensor)", operation)
}

// record gives result its history. inputs are the operands (nil for a plain number), and
// backs[i] works out the gradient of inputs[i]. Only tracked inputs are remembered; if
// there are none, the result stays a plain tensor.
func record(result *Tensor, inputs []*Tensor, backs []func(g []float64) []float64) {
	kept := []*Tensor{}
	keptBacks := []func(g []float64) []float64{}
	for i, in := range inputs {
		if in != nil && in.tracked() {
			kept = append(kept, in)
			keptBacks = append(keptBacks, backs[i])
		}
	}
	if len(kept) == 0 {
		return
	}
	result.track = &trace{inputs: kept, back: func(g []float64) [][]float64 {
		out := make([][]float64, len(keptBacks))
		for i, f := range keptBacks {
			out[i] = f(g)
		}
		return out
	}}
}

// fitTo shapes a gradient to its input. A single number that was used for every position
// receives all of those gradients added up.
func fitTo(input *Tensor, g []float64) []float64 {
	if input.isScalar() && len(g) > 1 {
		total := 0.0
		for _, x := range g {
			total += x
		}
		return []float64{total}
	}
	return g
}

// recordMath gives the result of + - * / its history. a or b is nil when that side was a
// plain number; x and y read the value each side had at position i.
func recordMath(op string, result, a, b *Tensor, x, y func(i int) float64) {
	// back builds the function for one side: share(g, i) is that side's gradient at i
	back := func(input *Tensor, share func(g []float64, i int) float64) func([]float64) []float64 {
		return func(g []float64) []float64 {
			out := make([]float64, len(g))
			for i := range g {
				out[i] = share(g, i)
			}
			return fitTo(input, out)
		}
	}
	var left, right func(g []float64, i int) float64
	switch op {
	case "+": // both sides pass the gradient straight through
		left = func(g []float64, i int) float64 { return g[i] }
		right = left
	case "-": // the right side gets it with the sign flipped
		left = func(g []float64, i int) float64 { return g[i] }
		right = func(g []float64, i int) float64 { return -g[i] }
	case "*": // each side's gradient is multiplied by the OTHER side's value
		left = func(g []float64, i int) float64 { return g[i] * y(i) }
		right = func(g []float64, i int) float64 { return g[i] * x(i) }
	case "/": // x / y changes by 1/y when x moves, and by -x/(y*y) when y moves
		left = func(g []float64, i int) float64 { return g[i] / y(i) }
		right = func(g []float64, i int) float64 { return -g[i] * x(i) / (y(i) * y(i)) }
	}
	record(result, []*Tensor{a, b}, []func([]float64) []float64{back(a, left), back(b, right)})
}

// recordMatmul: for C = matmul(A, B) with gradient G,
// A gets matmul(G, transpose(B)) and B gets matmul(transpose(A), G).
func recordMatmul(result, a, b *Tensor) {
	// multiply is matmul for shapes that are already known to fit
	multiply := func(x, y *Tensor) []float64 {
		product, _ := matmul(x.plain(), y.plain())
		return product.Data
	}
	flip := func(t *Tensor) *Tensor {
		flipped, _ := transpose(t.plain())
		return flipped
	}
	asResult := func(g []float64) *Tensor { return &Tensor{Shape: result.Shape, Data: g} }
	record(result, []*Tensor{a, b}, []func([]float64) []float64{
		func(g []float64) []float64 { return multiply(asResult(g), flip(b)) },
		func(g []float64) []float64 { return multiply(flip(a), asResult(g)) },
	})
}

// recordTranspose: the gradient is simply transposed back.
func recordTranspose(result, t *Tensor) {
	record(result, []*Tensor{t}, []func([]float64) []float64{
		func(g []float64) []float64 {
			back, _ := transpose(&Tensor{Shape: result.Shape, Data: g}) // result is always a table
			return back.Data
		},
	})
}

// recordEach is for functions that change each number on its own (relu, sigmoid, ...).
// slope(x, y) says how fast the function's answer y changes when its input x changes.
func recordEach(result, t *Tensor, slope func(x, y float64) float64) {
	record(result, []*Tensor{t}, []func([]float64) []float64{
		func(g []float64) []float64 {
			out := make([]float64, len(g))
			for i := range g {
				out[i] = g[i] * slope(t.Data[i], result.Data[i])
			}
			return out
		},
	})
}

// The slopes of the neural network functions. x is the input and y is the answer.
func reluSlope(x, _ float64) float64 {
	if x > 0 {
		return 1
	}
	return 0 // also at exactly 0
}
func sigmoidSlope(_, y float64) float64 { return y * (1 - y) }
func tanhSlope(_, y float64) float64    { return 1 - y*y }
func expSlope(_, y float64) float64     { return y } // exp is its own slope
func logSlope(x, _ float64) float64     { return 1 / x }

// total adds up a tensor (divide is 1 for sum, or the count for mean). A plain tensor gives
// a plain number. A tracked one gives a one-number tensor, which can keep the history.
func total(t *Tensor, divide float64) Value {
	sum := 0.0
	for _, x := range t.Data {
		sum += x
	}
	if !t.tracked() {
		return sum / divide
	}
	result := &Tensor{Shape: []int{}, Data: []float64{sum / divide}}
	record(result, []*Tensor{t}, []func([]float64) []float64{
		func(g []float64) []float64 { // every input number gets the same share
			out := make([]float64, len(t.Data))
			for i := range out {
				out[i] = g[0] / divide
			}
			return out
		},
	})
	return result
}

// backward works out the gradient of loss for every learnable tensor it was made from.
func backward(loss *Tensor) error {
	if !loss.tracked() {
		return fmt.Errorf("backward needs a result calculated from a learnable tensor")
	}
	if len(loss.Data) != 1 {
		return fmt.Errorf("backward needs a single number to start from, but got shape %s (use sum or mean first)", shapeText(loss.Shape))
	}

	// put the history in order: every tensor comes after the tensors it was made from
	order := []*Tensor{}
	seen := map[*Tensor]bool{}
	var visit func(t *Tensor)
	visit = func(t *Tensor) {
		if seen[t] {
			return
		}
		seen[t] = true
		for _, in := range t.track.inputs {
			visit(in)
		}
		order = append(order, t)
	}
	visit(loss)

	backwardRuns++
	// walk it from the end. A tensor's gradient is complete when we reach it, because
	// everything that used it comes later in the order.
	grads := map[*Tensor][]float64{loss: {1}}
	for i := len(order) - 1; i >= 0; i-- {
		t := order[i]
		g := grads[t]
		if t.track.leaf {
			t.track.grad = &Tensor{Shape: t.Shape, Data: g}
			t.track.run = backwardRuns
			continue
		}
		for j, share := range t.track.back(g) {
			in := t.track.inputs[j]
			if grads[in] == nil {
				grads[in] = make([]float64, len(in.Data))
			}
			for k, x := range share {
				grads[in][k] += x // a tensor used twice adds up both shares
			}
		}
		delete(grads, t) // not needed any more
	}
	return nil
}

// gradientOf returns the gradient that the last backward found for t.
func gradientOf(t *Tensor) (*Tensor, error) {
	if t.track == nil || !t.track.leaf {
		return nil, fmt.Errorf("grad only works on a tensor made with learnable(...); after changing a weight, wrap it in learnable(...) again")
	}
	if backwardRuns == 0 {
		return nil, fmt.Errorf("there is no gradient yet: call backward(loss) first")
	}
	if t.track.run != backwardRuns {
		return nil, fmt.Errorf("this tensor was not used by the last backward, so it has no gradient")
	}
	return t.track.grad, nil
}
```

## ai.go
```go
package main

// Talking to an AI model: ask("question") sends text to Ollama, a program that runs AI
// models on this computer, and gives back the model's answer as text.
//
// Safety rules:
//   - The address is fixed to this computer. There is no setting to send questions elsewhere,
//     no proxy is used, and redirects are not followed. (If you set up a "cloud" model inside
//     Ollama, Ollama itself would send questions online.)
//   - The answer is only TEXT. GAN never runs anything a model writes.
//   - GAN stops waiting after a time limit.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// ollamaURL and askTimeout are variables only so that the tests can use a fake server.
var (
	ollamaURL  = "http://127.0.0.1:11434/api/generate"
	askTimeout = 5 * time.Minute
)

// defaultModel is used when ask is not told which model to use.
const defaultModel = "qwen2.5:7b"

// maxAnswerBytes stops a runaway answer from using up all the memory.
const maxAnswerBytes = 10_000_000

// askModel sends one question to a model and returns its answer. Each call stands alone:
// the model does not remember earlier questions.
func askModel(question, model string) (string, error) {
	request, err := json.Marshal(map[string]any{
		"model":  model,
		"prompt": question,
		"stream": false, // send the whole answer at once, not word by word
	})
	if err != nil {
		return "", fmt.Errorf("ask could not prepare the question")
	}

	client := &http.Client{
		Timeout:   askTimeout,
		Transport: &http.Transport{Proxy: nil}, // go straight to this computer
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse // never pass the question on to another address
		},
	}
	response, err := client.Post(ollamaURL, "application/json", bytes.NewReader(request))
	if err != nil {
		return "", askFailed(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maxAnswerBytes+1))
	if err != nil {
		return "", askFailed(err)
	}
	if len(body) > maxAnswerBytes {
		return "", fmt.Errorf("ask stopped: the answer from Ollama was too long")
	}

	// Ollama answers with {"response": "..."} or, when something is wrong, {"error": "..."}
	var reply struct {
		Response *string `json:"response"`
		Error    string  `json:"error"`
	}
	readErr := json.Unmarshal(body, &reply)

	if response.StatusCode == http.StatusNotFound && strings.Contains(reply.Error, "not found") {
		return "", fmt.Errorf("Ollama does not have the model %q; see your models with: ollama list", model)
	}
	if reply.Error != "" {
		return "", fmt.Errorf("Ollama reported an error: %s", reply.Error)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama reported an error (status %d)", response.StatusCode)
	}
	if readErr != nil || reply.Response == nil {
		return "", fmt.Errorf("ask got a reply from Ollama that it could not read")
	}
	return strings.TrimSpace(*reply.Response), nil
}

// askFailed turns a network problem into plain words.
func askFailed(err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("the model took longer than %s to answer, so ask stopped waiting", waitText(askTimeout))
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" { // could not even connect
		return fmt.Errorf("ask could not reach Ollama on this computer; check that Ollama is running")
	}
	return fmt.Errorf("ask lost the connection to Ollama")
}

// waitText writes a time limit the way people say it.
func waitText(d time.Duration) string {
	if d >= time.Minute {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	return d.String()
}

// askArgument checks that an argument of ask is text that is not empty.
// what is "the question" or "the model name".
func askArgument(v Value, what string) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("ask needs text as %s, but got %s", what, typeName(v))
	}
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("ask needs text as %s, but it is empty", what)
	}
	return s, nil
}
```

## python.go
```go
package main

// python("module.function", arguments...) calls a function in Python from a GAN program.
//
// GAN starts a small Python helper (pybridge.py, which is built into gan.exe) the first time
// it is needed, and then talks to it through pipes, one JSON line per message.
//
// It has NO list of allowed functions. A GAN program that uses python() can do anything a
// Python program can do with your account, including deleting files and running commands.
// Read other people's programs before you run them.

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//go:embed pybridge.py
var pybridgeSource string

// pyValue is a GAN value on its way to or from Python. Numbers travel as text, so that
// "infinity" and "not a number" survive (plain JSON cannot hold them).
type pyValue struct {
	Type  string     `json:"type"`
	Value any        `json:"value,omitempty"`
	Items []*pyValue `json:"items,omitempty"`
	Shape []int      `json:"shape,omitempty"`
	Data  []string   `json:"data,omitempty"`
}

type pyReply struct {
	Ready bool     `json:"ready"`
	OK    bool     `json:"ok"`
	Value *pyValue `json:"value"`
	Error string   `json:"error"`
}

// pyHelper is the running Python process.
type pyHelper struct {
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
}

// helper is nil until python() is used for the first time.
var helper *pyHelper

const maxPyDepth = 100 // lists inside lists (also catches a list that contains itself)

// startPython finds Python and starts the helper.
func startPython() (*pyHelper, error) {
	notStarted := fmt.Errorf("python could not start Python; check that Python is installed and that \"python --version\" works")
	for _, name := range []string{"python", "python3"} {
		path, err := exec.LookPath(name) // never picks up a program from the current folder
		if err != nil {
			continue
		}
		cmd := exec.Command(path, "-u", "-c", pybridgeSource)
		cmd.Env = append(os.Environ(), "PYTHONUTF8=1", "PYTHONIOENCODING=utf-8")
		cmd.Stderr = os.Stdout // what Python code prints shows up with GAN's own output
		in, err := cmd.StdinPipe()
		if err != nil {
			return nil, notStarted
		}
		out, err := cmd.StdoutPipe()
		if err != nil {
			return nil, notStarted
		}
		if err := cmd.Start(); err != nil {
			continue
		}
		h := &pyHelper{cmd: cmd, in: in, out: bufio.NewReader(out)}
		if reply, err := h.read(); err == nil && reply.Ready {
			return h, nil
		}
		h.stop()
	}
	return nil, notStarted
}

func (h *pyHelper) read() (*pyReply, error) {
	line, err := h.out.ReadString('\n')
	if err != nil {
		return nil, err
	}
	var reply pyReply
	if err := json.Unmarshal([]byte(line), &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (h *pyHelper) stop() {
	h.in.Close()
	h.cmd.Process.Kill()
	h.cmd.Wait()
}

// toPython turns a GAN value into the form that is sent to Python.
func toPython(v Value, depth int) (*pyValue, error) {
	if depth > maxPyDepth {
		return nil, fmt.Errorf("python cannot take lists nested this deeply (does a list contain itself?)")
	}
	number := func(x float64) string { return strconv.FormatFloat(x, 'g', -1, 64) }
	switch x := v.(type) {
	case float64:
		return &pyValue{Type: "number", Value: number(x)}, nil
	case string:
		return &pyValue{Type: "text", Value: x}, nil
	case bool:
		return &pyValue{Type: "bool", Value: x}, nil
	case *List:
		out := &pyValue{Type: "list", Items: []*pyValue{}}
		for _, item := range x.Items {
			p, err := toPython(item, depth+1)
			if err != nil {
				return nil, err
			}
			out.Items = append(out.Items, p)
		}
		return out, nil
	case *Tensor:
		if x.tracked() { // not a restriction on power: the answer would be silently wrong
			return nil, fmt.Errorf("python cannot take a learnable tensor, because its gradients would be lost")
		}
		out := &pyValue{Type: "tensor", Shape: x.Shape, Data: make([]string, len(x.Data))}
		for i, n := range x.Data {
			out.Data[i] = number(n)
		}
		return out, nil
	}
	return nil, fmt.Errorf("python cannot take a %s", typeName(v))
}

// MarshalJSON keeps the parts that matter even when they are "empty": the text "",
// false, and an empty list must still be sent.
func (p *pyValue) MarshalJSON() ([]byte, error) {
	out := map[string]any{"type": p.Type}
	switch p.Type {
	case "list":
		out["items"] = p.Items
	case "tensor":
		out["shape"] = append([]int{}, p.Shape...)
		out["data"] = p.Data
	default:
		out["value"] = p.Value
	}
	return json.Marshal(out)
}

// fromPython turns a reply back into a GAN value.
func fromPython(p *pyValue) (Value, error) {
	broken := fmt.Errorf("python got a reply it could not read")
	switch p.Type {
	case "number":
		s, ok := p.Value.(string)
		if !ok {
			return nil, broken
		}
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, broken
		}
		return x, nil
	case "text":
		if p.Value == nil { // the empty text
			return "", nil
		}
		s, ok := p.Value.(string)
		if !ok {
			return nil, broken
		}
		return s, nil
	case "bool":
		b, _ := p.Value.(bool)
		return b, nil
	case "list":
		list := &List{Items: []Value{}}
		for _, item := range p.Items {
			v, err := fromPython(item)
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, v)
		}
		return list, nil
	case "tensor":
		t := &Tensor{Shape: append([]int{}, p.Shape...), Data: make([]float64, len(p.Data))}
		size := 1
		for _, n := range t.Shape {
			size *= n
		}
		if len(t.Shape) == 0 || size != len(t.Data) {
			return nil, broken
		}
		for i, s := range p.Data {
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, broken
			}
			t.Data[i] = x
		}
		return t, nil
	}
	return nil, broken
}

// callPython is the built-in function python(name, arguments...).
// The second result is false when the Python function gave nothing back (None).
func callPython(args []Value) (Value, bool, error) {
	name, ok := args[0].(string)
	if !ok {
		return nil, false, fmt.Errorf("python needs the name of a Python function as text, but got %s", typeName(args[0]))
	}
	parts := strings.Split(name, ".")
	if len(parts) < 2 || strings.Contains("."+name+".", "..") {
		return nil, false, fmt.Errorf("python needs a full name such as \"math.sqrt\", but got %q", name)
	}
	sent := []*pyValue{}
	for _, a := range args[1:] {
		p, err := toPython(a, 0)
		if err != nil {
			return nil, false, err
		}
		sent = append(sent, p)
	}
	request, err := json.Marshal(map[string]any{"function": name, "args": sent})
	if err != nil {
		return nil, false, fmt.Errorf("python could not prepare the request")
	}

	if helper == nil {
		if helper, err = startPython(); err != nil {
			return nil, false, err
		}
	}
	// there is no time limit: a long calculation is allowed (Ctrl+C stops the program)
	_, writeErr := helper.in.Write(append(request, '\n'))
	var reply *pyReply
	if writeErr == nil {
		reply, err = helper.read()
	}
	if writeErr != nil || err != nil {
		helper.stop()
		helper = nil // the next call starts a new one
		return nil, false, fmt.Errorf("python stopped unexpectedly while running %s", name)
	}
	if !reply.OK {
		return nil, false, fmt.Errorf("%s", reply.Error)
	}
	if reply.Value == nil {
		return nil, false, nil
	}
	v, err := fromPython(reply.Value)
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}
```

## pybridge.py
```python
# The Python helper of GAN's python(...) function. It is built into gan.exe (see python.go).
#
# GAN starts it once, then sends one request per line and reads one reply per line (JSON).
#
# It can call ANY function in Python's modules, with the same power as any Python program.
import sys

# Python normally looks for modules in the current folder first. Take that out BEFORE the
# first import, so a stray file such as json.py or numpy.py in the folder is never run by
# mistake.
sys.path = [p for p in sys.path if p not in ("", ".")]

import importlib
import json
import math
import numbers
import os

# Keep the two pipes to GAN for ourselves. After this, anything that Python code prints
# appears on the screen, and code that tries to read input gets nothing, so neither can
# break the conversation with GAN.
from_gan = os.fdopen(os.dup(0), "r", encoding="utf-8", newline="\n")
to_gan = os.fdopen(os.dup(1), "w", encoding="utf-8", newline="\n")
nothing = os.open(os.devnull, os.O_RDONLY)
os.dup2(nothing, 0)
os.dup2(2, 1)
sys.stdin = open(os.devnull, "r")
sys.stdout = sys.stderr

MAX_DEPTH = 100             # lists inside lists
MAX_ELEMENTS = 10_000_000   # the same limit GAN has for tensors


class Problem(Exception):
    """A mistake to report to the GAN program, already in plain words."""


def load_numpy():
    try:
        return importlib.import_module("numpy")
    except ImportError:
        raise Problem("sending a tensor to Python needs NumPy, and Python could not load it")


def from_value(v):
    """A GAN value (as sent by GAN) -> a Python value."""
    kind = v["type"]
    if kind == "number":
        x = float(v["value"])
        # whole numbers arrive as Python integers, because many functions insist on that
        return int(x) if x.is_integer() and abs(x) <= 2 ** 53 else x
    if kind in ("text", "bool"):
        return v["value"]
    if kind == "list":
        return [from_value(item) for item in v["items"]]
    numpy = load_numpy()
    return numpy.array([float(x) for x in v["data"]], dtype=float).reshape(v["shape"])


def number(x):
    try:
        return {"type": "number", "value": repr(float(x))}
    except OverflowError:
        raise Problem("Python returned a number that is too big for GAN")


def to_value(x, depth=0):
    """A Python result -> a GAN value. None means "no value"."""
    if depth > MAX_DEPTH:
        raise Problem("Python returned lists nested too deeply for GAN")
    numpy = sys.modules.get("numpy")
    if numpy is not None:
        if isinstance(x, numpy.ndarray):
            if x.dtype.kind not in "biuf":
                raise Problem("Python returned a NumPy array of " + str(x.dtype) + ", but GAN tensors only hold numbers")
            if x.ndim == 0:
                return to_value(x.item(), depth)
            if x.size == 0:
                raise Problem("Python returned an empty array, and a GAN tensor cannot be empty")
            if x.size > MAX_ELEMENTS or x.ndim > 32:
                raise Problem("Python returned an array that is too big for a GAN tensor")
            return {"type": "tensor", "shape": list(x.shape), "data": [repr(float(n)) for n in x.ravel().tolist()]}
        if isinstance(x, numpy.generic):
            x = x.item()
    if isinstance(x, bool):
        return {"type": "bool", "value": x}
    if isinstance(x, numbers.Real):
        return number(x)
    if isinstance(x, str):
        return {"type": "text", "value": x}
    if isinstance(x, (list, tuple)):
        items = [to_value(item, depth + 1) for item in x]
        if any(item is None for item in items):
            raise Problem("Python returned a list with None in it, which GAN cannot use yet")
        return {"type": "list", "items": items}
    if x is None:
        return None
    raise Problem("Python returned a " + type(x).__name__ + ", which GAN cannot use yet")


def find_function(name):
    """ "numpy.linalg.inv" -> import numpy.linalg and fetch inv. """
    parts = name.split(".")
    for split in range(len(parts) - 1, 0, -1):
        module_name = ".".join(parts[:split])
        try:
            found = importlib.import_module(module_name)
        except ModuleNotFoundError as err:
            # "not found" for this name (or a parent of it): try a shorter module name
            if err.name and (module_name == err.name or module_name.startswith(err.name + ".")):
                continue
            raise  # the module exists, but something IT needs is missing
        for part in parts[split:]:
            if not hasattr(found, part):
                raise Problem('Python has no function "' + name + '"')
            found = getattr(found, part)
        if not callable(found):
            raise Problem('"' + name + '" in Python is not a function')
        return found
    raise Problem('Python has no module "' + parts[0] + '"; is it installed?')


def handle(request):
    name = request["function"]
    try:
        function = find_function(name)
        result = function(*[from_value(a) for a in request["args"]])
        return {"ok": True, "value": to_value(result)}
    except Problem as err:
        return {"ok": False, "error": str(err)}
    except KeyboardInterrupt:
        raise
    except BaseException as err:  # includes sys.exit(), which must not end the helper silently
        message = " ".join(str(err).split())  # one line
        return {"ok": False, "error": "Python error in " + name + ": " + type(err).__name__ + (": " + message if message else "")}


def send(reply):
    to_gan.write(json.dumps(reply, allow_nan=False) + "\n")
    to_gan.flush()


send({"ready": True})
for line in from_gan:
    reply = handle(json.loads(line))
    try:
        sys.stderr.flush()  # so that what Python printed appears before GAN carries on
    except Exception:
        pass
    send(reply)
```
