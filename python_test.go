package main

// Checks python(...) with the real Python on this computer. Run with: go test
// The tests are skipped when Python (or, for tensors, NumPy) is not installed.

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func needPython(t *testing.T, numpy bool) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("Python is not installed")
	}
	if numpy && exec.Command("python", "-c", "import numpy").Run() != nil {
		t.Skip("NumPy is not installed")
	}
}

// py runs python(name, args...) the way a GAN program would.
func py(t *testing.T, args ...Value) Value {
	t.Helper()
	v, _, err := callPython(args)
	if err != nil {
		t.Fatalf("python(%v): %v", args, err)
	}
	return v
}

func list(items ...Value) *List { return &List{Items: append([]Value{}, items...)} }

func TestPythonValuesGoBothWays(t *testing.T) {
	needPython(t, false)
	// copy.deepcopy gives back what it was given
	for _, v := range []Value{
		3.5, -2.0, 0.0, 1e300, math.Inf(1), math.Inf(-1), "", "hi \"there\"\n\tC:\\x é 🙂", true, false,
		list(), list(1.0, "a", false, "", list(list())),
	} {
		got := py(t, "copy.deepcopy", v)
		if !equalValues(got, v) || typeName(got) != typeName(v) {
			t.Errorf("sent %#v, got back %#v", v, got)
		}
	}
	if x := py(t, "copy.deepcopy", math.NaN()).(float64); !math.IsNaN(x) {
		t.Errorf("NaN came back as %v", x)
	}
	// whole numbers arrive as Python integers, others as floats (repr shows which)
	for v, want := range map[float64]string{5: "5", -3: "-3", 2.5: "2.5"} {
		if got := py(t, "builtins.repr", v); got != want {
			t.Errorf("Python saw %v as %v, want %v", v, got, want)
		}
	}
	if got := py(t, "math.factorial", 5.0); got != 120.0 {
		t.Errorf("factorial gave %v", got)
	}
	if got := py(t, "builtins.format", 2.5, ""); got != "2.5" {
		t.Errorf("format gave %v", got)
	}
}

func TestPythonTensorsGoBothWays(t *testing.T) {
	needPython(t, true)
	for _, v := range []*Tensor{
		{Shape: []int{3}, Data: []float64{1, 2.5, -3}},
		{Shape: []int{2, 3}, Data: []float64{1, 2, 3, 4, 5, 6}},
		{Shape: []int{2, 1, 2}, Data: []float64{1, math.Inf(1), 3, 4}},
	} {
		got, ok := py(t, "numpy.copy", v).(*Tensor)
		if !ok || !got.equal(v) {
			t.Errorf("sent %v, got back %v", v, got)
		}
	}
	// NumPy gives the same matrix product as GAN itself
	a := &Tensor{Shape: []int{2, 3}, Data: []float64{1.5, 2, -3, 4, 0.5, 6}}
	b := &Tensor{Shape: []int{3, 2}, Data: []float64{7, -8, 9, 1, 2, 3}}
	ours, _ := matmul(a, b)
	if theirs := py(t, "numpy.dot", a, b).(*Tensor); !theirs.equal(ours) {
		t.Errorf("numpy.dot gave %v, matmul gave %v", theirs, ours)
	}
	if got := py(t, "numpy.max", a); got != 6.0 {
		t.Errorf("numpy.max gave %#v", got)
	}
	if got := py(t, "numpy.linalg.det", &Tensor{Shape: []int{2, 2}, Data: []float64{2, 0, 0, 3}}).(float64); math.Abs(got-6) > 1e-9 {
		t.Errorf("det gave %v", got)
	}
}

func TestPythonErrors(t *testing.T) {
	needPython(t, false)
	self := list()
	self.Items = append(self.Items, self)
	cases := []struct {
		args  []Value
		wants string
	}{
		{[]Value{"math.no_such_function", 1.0}, `Python has no function "math.no_such_function"`},
		{[]Value{"no_such_module_xyz.f"}, `Python has no module "no_such_module_xyz"`},
		{[]Value{"math.pi"}, "is not a function"},
		{[]Value{"math.sqrt", -1.0}, "Python error in math.sqrt: ValueError"},
		{[]Value{"sys.exit", 3.0}, "Python error in sys.exit: SystemExit: 3"},
		{[]Value{"builtins.dict"}, "returned a dict, which GAN cannot use yet"},
		{[]Value{"math.factorial", 200.0}, "too big for GAN"},
		{[]Value{"sqrt", 4.0}, `full name such as "math.sqrt"`},
		{[]Value{"math..sqrt", 4.0}, "full name"},
		{[]Value{".sqrt"}, "full name"},
		{[]Value{42.0}, "name of a Python function as text, but got number"},
		{[]Value{"builtins.len", learnable(&Tensor{Shape: []int{1}, Data: []float64{1}})}, "gradients would be lost"},
		{[]Value{"builtins.len", list(learnable(&Tensor{Shape: []int{1}, Data: []float64{1}}))}, "gradients would be lost"},
		{[]Value{"builtins.len", self}, "nested this deeply"},
		{[]Value{"builtins.len", &Module{}}, "python cannot take a module"},
		{[]Value{"os._exit", 1.0}, "stopped unexpectedly"},
	}
	for _, c := range cases {
		_, _, err := callPython(c.args)
		if err == nil || !strings.Contains(err.Error(), c.wants) {
			t.Errorf("%v: got %v, want %q", c.args[0], err, c.wants)
		}
	}
	// after Python stopped, the next call starts a new helper
	if got := py(t, "math.sqrt", 16.0); got != 4.0 {
		t.Errorf("after a restart: %v", got)
	}
}

func TestPythonPrintingAndInputDoNotBreakTheConnection(t *testing.T) {
	needPython(t, false)
	// print gives no value back, and what it writes must not be mistaken for a reply
	v, returned, err := callPython([]Value{"builtins.print", `{"ok": true}`})
	if v != nil || returned || err != nil {
		t.Errorf("print gave %v %v %v", v, returned, err)
	}
	py(t, "os.system", "echo from a command") // output of other programs too
	// input() must not swallow GAN's next request
	if _, _, err := callPython([]Value{"builtins.input"}); err == nil || !strings.Contains(err.Error(), "EOFError") {
		t.Errorf("input gave %v", err)
	}
	if got := py(t, "math.sqrt", 9.0); got != 3.0 {
		t.Errorf("got %v", got)
	}
}

// A stray .py file in the current folder must never be run just because of its name.
func TestPythonIgnoresFilesInTheCurrentFolder(t *testing.T) {
	needPython(t, false)
	folder := t.TempDir()
	marker := filepath.Join(folder, "planted-code-ran")
	code := "open(r'" + marker + "', 'w').close()\ndef f():\n    return 1\n"
	for _, name := range []string{"json.py", "importlib.py", "math.py", "numbers.py", "os.py", "planted_module.py"} {
		os.WriteFile(filepath.Join(folder, name), []byte(code), 0o644)
	}
	t.Chdir(folder)
	if helper != nil { // start a fresh helper inside that folder
		helper.stop()
		helper = nil
	}
	if got := py(t, "math.sqrt", 25.0); got != 5.0 {
		t.Errorf("got %v", got)
	}
	_, _, err := callPython([]Value{"planted_module.f"})
	if err == nil || !strings.Contains(err.Error(), "no module") {
		t.Errorf("a module in the current folder was found: %v", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("code from a file in the current folder ran")
	}
	helper.stop() // so the temporary folder can be removed
	helper = nil
}
