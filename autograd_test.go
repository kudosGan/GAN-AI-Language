package main

// Checks GAN's gradients against a second, independent method: change one input number a
// tiny bit and measure how much the result changes. Run with: go test

import (
	"math"
	"testing"
)

func tens(shape []int, data ...float64) *Tensor { return &Tensor{Shape: shape, Data: data} }

func TestGradientsMatchNumericalCheck(t *testing.T) {
	// must stops the test if an operation reports an error
	math2 := func(op string, a, b Value) *Tensor {
		v, err := tensorMath(op, a, b)
		if err != nil {
			t.Fatal(err)
		}
		return v.(*Tensor)
	}
	mul := func(a, b *Tensor) *Tensor {
		v, err := matmul(a, b)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	// add up a tensor; the result is a number for plain tensors and a tensor for tracked ones
	add := func(x *Tensor, divide float64) *Tensor {
		if v, ok := total(x, divide).(*Tensor); ok {
			return v
		}
		return tens([]int{}, total(x, divide).(float64))
	}

	a := tens([]int{2, 3}, 0.5, -1.2, 2, 3, 0.7, -0.3)
	b := tens([]int{2, 3}, 1.5, 0.4, -2, 0.1, 2.2, 1)
	c := tens([]int{3, 2}, -0.6, 1.1, 0.9, 0.2, -1.4, 0.8)
	k := tens([]int{}, 1.7) // a one-number tensor, used for every position

	// call runs one of GAN's built-in functions, the same way a GAN program would
	call := func(name string, x *Tensor) *Tensor {
		v, _, err := builtins[name].run([]Value{x})
		if err != nil {
			t.Fatal(err)
		}
		return v.(*Tensor)
	}
	// inputs for the new operations. relu is not checked at exactly 0, where its slope jumps.
	pos := tens([]int{2, 3}, 0.5, 1.2, 2, 3, 0.7, 0.3)
	ones := tens([]int{2, 1}, 1, 1)
	bias := tens([]int{1, 4}, 0.3, -0.2, 0.5, 0.1)
	w1 := tens([]int{3, 4}, 0.2, -0.5, 0.7, 0.1, -0.3, 0.8, 0.4, -0.6, 0.9, 0.05, -0.2, 0.35)
	w2 := tens([]int{4, 1}, 0.6, -0.4, 0.25, 0.8)
	target := tens([]int{2, 1}, 0.2, 0.9)

	cases := []struct {
		name   string
		inputs []*Tensor
		f      func(in []*Tensor) *Tensor
	}{
		{"add", []*Tensor{a, b}, func(in []*Tensor) *Tensor { return add(math2("*", math2("+", in[0], in[1]), b), 1) }},
		{"subtract", []*Tensor{a, b}, func(in []*Tensor) *Tensor { return add(math2("*", math2("-", in[0], in[1]), a), 1) }},
		{"multiply", []*Tensor{a, b}, func(in []*Tensor) *Tensor { return add(math2("*", in[0], in[1]), 1) }},
		{"numbers on both sides", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("-", 3.0, math2("*", in[0], -2.5)), 1) }},
		{"same tensor twice", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", in[0], in[0]), 1) }},
		{"result used twice", []*Tensor{a}, func(in []*Tensor) *Tensor {
			s := math2("*", in[0], in[0])
			return add(math2("+", math2("*", s, s), s), 1)
		}},
		{"one-number tensor", []*Tensor{k, a}, func(in []*Tensor) *Tensor {
			return add(math2("*", math2("-", in[1], in[0]), math2("*", in[0], in[1])), 1)
		}},
		{"matmul", []*Tensor{a, c}, func(in []*Tensor) *Tensor {
			p := mul(in[0], in[1])
			return add(math2("*", p, p), 1)
		}},
		{"matmul with itself", []*Tensor{tens([]int{2, 2}, 1, 2, 3, 4)}, func(in []*Tensor) *Tensor { return add(mul(in[0], in[0]), 1) }},
		{"mean, and a sum used again", []*Tensor{a, c}, func(in []*Tensor) *Tensor {
			p := mul(in[0], in[1])
			m := add(p, 4) // the mean of the 2 by 2 result
			d := math2("-", p, m)
			return add(math2("*", d, d), 4)
		}},
		{"divide, both sides", []*Tensor{a, pos}, func(in []*Tensor) *Tensor { return add(math2("*", math2("/", in[0], in[1]), b), 1) }},
		{"divide by a number, and a number divided", []*Tensor{pos}, func(in []*Tensor) *Tensor {
			return add(math2("+", math2("/", in[0], 4.0), math2("/", 3.0, in[0])), 1)
		}},
		{"one-number tensor on top", []*Tensor{k, pos}, func(in []*Tensor) *Tensor { return add(math2("/", in[0], in[1]), 1) }},
		{"one-number tensor underneath", []*Tensor{a, k}, func(in []*Tensor) *Tensor { return add(math2("/", in[0], in[1]), 1) }},
		{"divided by itself", []*Tensor{pos}, func(in []*Tensor) *Tensor { return add(math2("/", in[0], math2("+", in[0], 1.0)), 1) }},
		{"relu", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", call("relu", in[0]), b), 1) }},
		{"sigmoid", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", call("sigmoid", in[0]), b), 1) }},
		{"tanh", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", call("tanh", in[0]), b), 1) }},
		{"exp", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", call("exp", in[0]), b), 1) }},
		{"log", []*Tensor{pos}, func(in []*Tensor) *Tensor { return add(math2("*", call("log", in[0]), b), 1) }},
		{"transpose", []*Tensor{a}, func(in []*Tensor) *Tensor { return add(math2("*", call("transpose", in[0]), c), 1) }},
		{"two layers", []*Tensor{w1, bias, w2}, func(in []*Tensor) *Tensor {
			// a is the data: 2 examples with 3 numbers each
			hidden := call("tanh", math2("+", mul(a, in[0]), mul(ones, in[1])))
			out := call("sigmoid", mul(call("relu", math2("+", hidden, 2.0)), in[2]))
			diff := math2("-", out, target)
			logs := call("log", call("exp", call("transpose", math2("/", out, 2.0))))
			return add(math2("+", add(math2("*", diff, diff), 2), add(logs, 2)), 1)
		}},
	}

	for _, test := range cases {
		leaves := make([]*Tensor, len(test.inputs))
		for i, in := range test.inputs {
			leaves[i] = learnable(in)
		}
		if err := backward(test.f(leaves)); err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
		for i, leaf := range leaves {
			got, err := gradientOf(leaf)
			if err != nil {
				t.Fatalf("%s: %v", test.name, err)
			}
			if !sameShape(got.Shape, leaf.Shape) {
				t.Fatalf("%s: gradient shape %v, want %v", test.name, got.Shape, leaf.Shape)
			}
			for j := range leaf.Data {
				// the result with input number j moved by delta
				moved := func(delta float64) float64 {
					plain := make([]*Tensor, len(test.inputs))
					for n, in := range test.inputs {
						plain[n] = tens(in.Shape, append([]float64{}, in.Data...)...)
					}
					plain[i].Data[j] += delta
					// functions such as "add" also use a and b directly; they are plain tensors
					return test.f(plain).Data[0]
				}
				const h = 1e-5
				want := (moved(h) - moved(-h)) / (2 * h)
				if math.Abs(got.Data[j]-want) > 1e-5 {
					t.Errorf("%s: input %d position %d: gradient %g, numerical check %g", test.name, i, j, got.Data[j], want)
				}
			}
		}
	}
}
