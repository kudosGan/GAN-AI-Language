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
