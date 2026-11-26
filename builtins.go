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
