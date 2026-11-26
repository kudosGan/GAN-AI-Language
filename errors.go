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
