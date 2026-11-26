# GAN programming language: project handoff

## About
GAN is my own programming language (file extension .gan), written in Go. Folder: C:\Desktop\GAN. Go 1.27 on Windows, module `mylang`.
Build: `go build -o gan.exe .`   Run: `gan file.gan`   Version: `gan --version` (GAN 2.0)
Pipeline: Lexer -> Parser -> Interpreter (lexer.go, parser.go, interpreter.go; builtins.go, modules.go, errors.go; main.go ties them together).

## How I want you to work
- I'm a beginner in Go and in language design. Go one step at a time.
- Before changing any code, explain in plain language what we're doing and why, and wait for me to say "go".
- Then make the change, tell me which files changed, and give me the exact commands to build and test it.
- Include a test .gan file for each new feature, including one that should produce an error.
- Keep the code simple and commented. Don't add features I haven't asked for.
- After each step works, wait for me to say "next".
(For steps 15 to 26 I allowed continuing without stopping. Normal rules apply again from here.)

## Progress: steps 1 to 26 are done
1-6 lexer/parser/print, 7 interpreter, 8 variables, 9 math+precedence, 10 comparisons+if/else, 11 while+assignment,
12 functions, 13 error messages with line numbers, 14 parentheses+negative numbers, 15 comments (//), 16 decimals,
17 joining text with +, 18 true/false + and/or/not + else if, 19 scope + final/var, 20 for loops (for i in 1 to 5),
21 lists, 22 built-in functions, 23 import + libraries (math, text, lists) + % + fail(), 24 optional types (let x: number),
25 Windows installer (installer\GAN-Setup.exe), 26 tutorial (tutorial\, 14 lessons).
Still to come (only if I ask): break/continue, dictionaries, text escapes, random numbers, reading files, concurrency, package downloads.

## The language today
- Values: number (float64), text, bool, list, module. Blocks use { }. Comments use //.
- let/var/final, assignment, if/else if/else, while, for i in a to b (and for x in list/text), func/return, import "name".
- Operators: + - * / %, == != < > <= >=, and or not, [index], module.name, f(x). + joins when either side is text.
- Optional types: number, text, bool, list, any (checked while running). typeof(x).
- Built-ins: len push pop str num input abs floor round sqrt pow min max upper lower contains fail typeof.
- Libraries in libs\ (written in GAN): math, text, lists. A path import ("modules/x.gan") is relative to the importing file.
- Errors: "Syntax error on line N" / "Runtime error on line N" plus the source line; errors in imported files name the file.
- Scope: each block and each function call has its own scope; functions can see their module's globals. Lists are shared, not copied.

## Project layout
- main.go lexer.go parser.go interpreter.go builtins.go modules.go errors.go: the language
- libs\*.gan standard libraries; modules\*.gan helper files used by tests
- *.gan in the root: examples and tests (names ending in _error are meant to fail)
- tests\run.py: regression tests (python tests\run.py; --update rewrites tests\expected after checking output by eye)
- tutorial\: the course. python tutorial\check.py verifies every example against the real gan.exe (--fill writes outputs)
- agent\documenter.py: Python agent using local Ollama (qwen2.5:7b) that writes agent\document\*.md from the project and agent\steps.json. Run with --watch in a split terminal. Add an entry to agent\steps.json after every step.
- installer\build-installer.ps1 builds installer\GAN-Setup.exe (IExpress + a PowerShell/WinForms wizard, per-user install, no admin).
- vscode-gan\: the "Seti + GAN" VS Code icon theme; gan-icon.svg/.ico: the logo.

## Things worth remembering
- Strings have no escape sequences and cannot contain a double quote.
- No break/continue. Functions must be defined before the call runs.
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
		fmt.Println("GAN 2.0")
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
			start := i + 1
			i++
			for i < len(src) && src[i] != '"' {
				if src[i] == '\n' {
					line++
				}
				i++
			}
			if i >= len(src) {
				tokens = append(tokens, Token{ILLEGAL, unterminatedString, startLine})
			} else {
				tokens = append(tokens, Token{STRING, src[start:i], startLine})
				i++
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
	"fmt"
	"strconv"
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
	}
	return 0
}

// ---- The parser ----

type Parser struct {
	tokens []Token
	pos    int
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
var validTypes = map[string]bool{"number": true, "text": true, "bool": true, "list": true, "any": true}

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
		return "", at(typeTok.Line, fmt.Errorf("unknown type %q (the types are: number, text, bool, list, any)", typeTok.Value))
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

// parseStatement reads one command: print, let, if, while, func, return, an assignment or a call.
func (p *Parser) parseStatement() (Statement, error) {
	name, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}
	line := name.Line // the statement is on the line where its first word is

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
		body, err := p.parseBlock()
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
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if last == nil { // for x in somethingToLoopOver { ... }
			return ForStatement{Var: loopVar.Value, Iter: first, Body: body, Line: line}, nil
		}
		return ForStatement{Var: loopVar.Value, From: first, To: last, Body: body, Line: line}, nil

	case "func": // func name(a, b) { ... }   (types are optional: func add(a: number): number { ... })
		fname, err := p.expect(IDENT)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(LPAREN); err != nil {
			return nil, err
		}
		var params []Param
		if tok := p.peek(); tok != nil && tok.Type == RPAREN { // no parameters
			p.pos++
		} else {
			for {
				param, err := p.expect(IDENT)
				if err != nil {
					return nil, err
				}
				paramType, err := p.parseOptionalType()
				if err != nil {
					return nil, err
				}
				params = append(params, Param{Name: param.Value, Type: paramType})
				if tok := p.peek(); tok != nil && tok.Type == COMMA {
					p.pos++
					continue
				}
				break
			}
			if _, err := p.expect(RPAREN); err != nil {
				return nil, err
			}
		}
		returnType, err := p.parseOptionalType() // func add(a, b): number { ... }
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(LBRACE); err != nil {
			return nil, err
		}
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		return FuncStatement{Name: fname.Value, Params: params, ReturnType: returnType, Body: body, Line: line}, nil

	case "import": // import "name"
		path, err := p.expect(STRING)
		if err != nil {
			return nil, err
		}
		return ImportStatement{Path: path.Value, Line: line}, nil

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
// or a *List.
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
	return a == b
}

// binding is one variable: its value, and whether it is final (cannot change).
type binding struct {
	value Value
	final bool
	typ   string // "number", "text", "bool", "list" or "any"; "" means no type was given
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
		}
		return nil, fmt.Errorf("only lists and text can be indexed")
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
			if err != nil || returned {
				return val, returned, err
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
			if err != nil || returned {
				return val, returned, err
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
		if err != nil || returned {
			return val, returned, err
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

	// typeof(x): "number", "text", "bool", "list" or "module"
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

