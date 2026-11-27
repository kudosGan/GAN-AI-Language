# Lesson 15: How a model learns

A model is a calculation with numbers in it called **weights**. **Learning** means changing the weights
until the answers are good. GAN can work out which way to change them.

## The idea

1. Calculate an answer using the weights.
2. Measure how wrong it is. This number is called the **loss**. Smaller is better.
3. Find the **gradient** of each weight: how the loss changes when the weight changes a tiny bit.
4. Move each weight a small step the **opposite** way. Repeat.

## Gradients

Mark a tensor with `learnable`. Calculate something from it, then call `backward` on a result that is
**one number**. After that, `grad` gives the gradient.

Here the result is `w * w`. When `w` is 3, making `w` a tiny bit bigger makes `w * w` grow 6 times as
fast, so the gradient is 6:

```gan
let w = learnable(tensor([3]))
let result = sum(w * w)
backward(result)
print(result)
print(grad(w))
```

```output
tensor(9)
tensor([6])
```

`result` is shown as `tensor(9)`. It is a one-number tensor that remembers how it was calculated.
`item(result)` gives the plain number 9.

## A model that learns to double

This model has one weight. It should find out that the answer is the input times 2:

```gan
let x = tensor([[1], [2]])
let target = tensor([[2], [4]])
var w = learnable(tensor([[0]]))

for round in 1 to 3 {
    let diff = matmul(x, w) - target
    let loss = sum(diff * diff)
    backward(loss)
    print("loss " + item(loss) + ", gradient " + grad(w))
    w = learnable(w - 0.05 * grad(w))
    print(w)
}
```

```output
loss 20, gradient tensor([[-20]])
tensor([[1]])
loss 5, gradient tensor([[-10]])
tensor([[1.5]])
loss 1.25, gradient tensor([[-5]])
tensor([[1.75]])
```

The loss gets smaller and the weight gets closer to 2. The gradient is **negative**, which means
"a bigger weight gives a smaller loss", so subtracting the gradient makes the weight bigger.

The number `0.05` is the **learning rate**: how big a step to take. Too small is slow. Too big jumps
past the answer, and the loss can grow instead of shrink.

Wrapping the new weight in `learnable(...)` again is important. It starts a fresh history for the next
round.

## What works with gradients today

`+`, `-`, `*`, minus signs, `matmul`, `sum` and `mean`.

Other operations give an error on a learnable tensor. GAN refuses instead of giving you a wrong
gradient:

```gan
let w = learnable(tensor([1, 2]))
print(relu(w))
```

```output
Runtime error on line 2: gradients for relu are not supported yet (it was used on a learnable tensor)
    2 | print(relu(w))
```

Note that `item(...)` gives a plain number, with no history. Use it for printing and for deciding
things, not in the middle of a calculation you want gradients for.

## Exercises

1. Use `backward` to find the gradient of `sum(w * w)` when `w` is `tensor([1, 5])`. Guess first.
2. Change the doubling model so that it runs 10 rounds and prints only the final loss. Is it smaller
   than 0.001?
3. Change the learning rate to `0.25` and print the loss in each of 3 rounds. What happens?

Back to the [course index](README.md)
