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
