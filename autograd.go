package main

// Automatic gradients: how a model finds out which way to change its weights.
//
// A GRADIENT says how much the result (usually the error) changes when a number changes a
// tiny bit. GAN works it out like this:
//
//  1. learnable(w) marks a tensor we want gradients for.
//  2. Every calculation that uses it remembers its inputs, and how to pass a gradient
//     back to them. The result remembers this too, so a chain of history builds up.
//  3. backward(loss) starts at the end with gradient 1 and walks the history in reverse,
//     handing each input its share (this is the "chain rule").
//
// A tensor with history is called "tracked". Every operation that cannot pass gradients back
// yet must REFUSE tracked tensors (see noGradient). Ignoring them would give wrong gradients
// without any sign that something was wrong.

import "fmt"

// trace is the history a tracked tensor carries.
type trace struct {
	leaf   bool      // true for a tensor made by learnable(): the history starts here
	inputs []*Tensor // the tracked tensors this one was calculated from
	// back takes the gradient of the result and returns the gradient for each of the inputs
	back func(g []float64) [][]float64

	grad *Tensor // leaf only: the gradient found by backward
	run  int     // leaf only: which backward found it
}

// backwardRuns counts the calls to backward, so grad never hands out an out-of-date answer.
var backwardRuns int

// tracked says whether a tensor has history.
func (t *Tensor) tracked() bool { return t.track != nil }

// plain is the same tensor without any history.
func (t *Tensor) plain() *Tensor { return &Tensor{Shape: t.Shape, Data: t.Data} }

// isScalar says whether t is a single number with no dimensions, shown as tensor(5).
func (t *Tensor) isScalar() bool { return len(t.Shape) == 0 }

// learnable makes the starting point of a history. It is always a new tensor, so any
// older history is left behind.
func learnable(t *Tensor) *Tensor {
	return &Tensor{Shape: t.Shape, Data: t.Data, track: &trace{leaf: true}}
}

// noGradient is the error for an operation that cannot pass gradients back yet.
func noGradient(operation string) error {
	return fmt.Errorf("gradients for %s are not supported yet (it was used on a learnable tensor)", operation)
}

// record gives result its history. inputs are the operands (nil for a plain number), and
// backs[i] works out the gradient of inputs[i]. Only tracked inputs are remembered; if
// there are none, the result stays a plain tensor.
func record(result *Tensor, inputs []*Tensor, backs []func(g []float64) []float64) {
	kept := []*Tensor{}
	keptBacks := []func(g []float64) []float64{}
	for i, in := range inputs {
		if in != nil && in.tracked() {
			kept = append(kept, in)
			keptBacks = append(keptBacks, backs[i])
		}
	}
	if len(kept) == 0 {
		return
	}
	result.track = &trace{inputs: kept, back: func(g []float64) [][]float64 {
		out := make([][]float64, len(keptBacks))
		for i, f := range keptBacks {
			out[i] = f(g)
		}
		return out
	}}
}

// fitTo shapes a gradient to its input. A single number that was used for every position
// receives all of those gradients added up.
func fitTo(input *Tensor, g []float64) []float64 {
	if input.isScalar() && len(g) > 1 {
		total := 0.0
		for _, x := range g {
			total += x
		}
		return []float64{total}
	}
	return g
}

// recordMath gives the result of + - * its history. a or b is nil when that side was a
// plain number; x and y read the value each side had at position i.
func recordMath(op string, result, a, b *Tensor, x, y func(i int) float64) {
	// back builds the function for one side: share(g, i) is that side's gradient at i
	back := func(input *Tensor, share func(g []float64, i int) float64) func([]float64) []float64 {
		return func(g []float64) []float64 {
			out := make([]float64, len(g))
			for i := range g {
				out[i] = share(g, i)
			}
			return fitTo(input, out)
		}
	}
	var left, right func(g []float64, i int) float64
	switch op {
	case "+": // both sides pass the gradient straight through
		left = func(g []float64, i int) float64 { return g[i] }
		right = left
	case "-": // the right side gets it with the sign flipped
		left = func(g []float64, i int) float64 { return g[i] }
		right = func(g []float64, i int) float64 { return -g[i] }
	case "*": // each side's gradient is multiplied by the OTHER side's value
		left = func(g []float64, i int) float64 { return g[i] * y(i) }
		right = func(g []float64, i int) float64 { return g[i] * x(i) }
	}
	record(result, []*Tensor{a, b}, []func([]float64) []float64{back(a, left), back(b, right)})
}

// recordMatmul: for C = matmul(A, B) with gradient G,
// A gets matmul(G, transpose(B)) and B gets matmul(transpose(A), G).
func recordMatmul(result, a, b *Tensor) {
	// multiply is matmul for shapes that are already known to fit
	multiply := func(x, y *Tensor) []float64 {
		product, _ := matmul(x.plain(), y.plain())
		return product.Data
	}
	flip := func(t *Tensor) *Tensor {
		flipped, _ := transpose(t.plain())
		return flipped
	}
	asResult := func(g []float64) *Tensor { return &Tensor{Shape: result.Shape, Data: g} }
	record(result, []*Tensor{a, b}, []func([]float64) []float64{
		func(g []float64) []float64 { return multiply(asResult(g), flip(b)) },
		func(g []float64) []float64 { return multiply(flip(a), asResult(g)) },
	})
}

// total adds up a tensor (divide is 1 for sum, or the count for mean). A plain tensor gives
// a plain number. A tracked one gives a one-number tensor, which can keep the history.
func total(t *Tensor, divide float64) Value {
	sum := 0.0
	for _, x := range t.Data {
		sum += x
	}
	if !t.tracked() {
		return sum / divide
	}
	result := &Tensor{Shape: []int{}, Data: []float64{sum / divide}}
	record(result, []*Tensor{t}, []func([]float64) []float64{
		func(g []float64) []float64 { // every input number gets the same share
			out := make([]float64, len(t.Data))
			for i := range out {
				out[i] = g[0] / divide
			}
			return out
		},
	})
	return result
}

// backward works out the gradient of loss for every learnable tensor it was made from.
func backward(loss *Tensor) error {
	if !loss.tracked() {
		return fmt.Errorf("backward needs a result calculated from a learnable tensor")
	}
	if len(loss.Data) != 1 {
		return fmt.Errorf("backward needs a single number to start from, but got shape %s (use sum or mean first)", shapeText(loss.Shape))
	}

	// put the history in order: every tensor comes after the tensors it was made from
	order := []*Tensor{}
	seen := map[*Tensor]bool{}
	var visit func(t *Tensor)
	visit = func(t *Tensor) {
		if seen[t] {
			return
		}
		seen[t] = true
		for _, in := range t.track.inputs {
			visit(in)
		}
		order = append(order, t)
	}
	visit(loss)

	backwardRuns++
	// walk it from the end. A tensor's gradient is complete when we reach it, because
	// everything that used it comes later in the order.
	grads := map[*Tensor][]float64{loss: {1}}
	for i := len(order) - 1; i >= 0; i-- {
		t := order[i]
		g := grads[t]
		if t.track.leaf {
			t.track.grad = &Tensor{Shape: t.Shape, Data: g}
			t.track.run = backwardRuns
			continue
		}
		for j, share := range t.track.back(g) {
			in := t.track.inputs[j]
			if grads[in] == nil {
				grads[in] = make([]float64, len(in.Data))
			}
			for k, x := range share {
				grads[in][k] += x // a tensor used twice adds up both shares
			}
		}
		delete(grads, t) // not needed any more
	}
	return nil
}

// gradientOf returns the gradient that the last backward found for t.
func gradientOf(t *Tensor) (*Tensor, error) {
	if t.track == nil || !t.track.leaf {
		return nil, fmt.Errorf("grad only works on a tensor made with learnable(...); after changing a weight, wrap it in learnable(...) again")
	}
	if backwardRuns == 0 {
		return nil, fmt.Errorf("there is no gradient yet: call backward(loss) first")
	}
	if t.track.run != backwardRuns {
		return nil, fmt.Errorf("this tensor was not used by the last backward, so it has no gradient")
	}
	return t.track.grad, nil
}
