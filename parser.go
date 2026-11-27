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
		// a function body starts outside any loop, even if the func is written inside one
		outerLoops := p.loopDepth
		p.loopDepth = 0
		body, err := p.parseBlock()
		p.loopDepth = outerLoops
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
