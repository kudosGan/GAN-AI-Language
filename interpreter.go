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
