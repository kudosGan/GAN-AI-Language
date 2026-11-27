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
		fmt.Println("GAN 0.2 (beta)")
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
