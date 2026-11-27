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
