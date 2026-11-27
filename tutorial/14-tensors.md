# Lesson 14: Tensors

AI programs work with large grids of numbers. Such a grid is called a **tensor**. Tensors are built
into GAN, so there is nothing to install or import.

## Making a tensor

Give `tensor` a list of numbers, or a list of rows:

```gan
let row = tensor([1, 2, 3])
let table = tensor([[1, 2, 3], [4, 5, 6]])
print(row)
print(table)
print(shape(row))
print(shape(table))
```

```output
tensor([1, 2, 3])
tensor([[1, 2, 3], [4, 5, 6]])
[3]
[2, 3]
```

The **shape** is the size of the grid. `[2, 3]` means 2 rows and 3 columns. Every row must be the same
length:

```gan
let bad = tensor([[1, 2], [3]])
```

```output
Runtime error on line 1: all rows of a tensor must have the same shape
    1 | let bad = tensor([[1, 2], [3]])
```

## Math on every number at once

With a single number, the math happens to every number in the tensor. No loop is needed:

```gan
let t = tensor([1, 2, 3])
print(t * 10)
print(t + 1)
```

```output
tensor([10, 20, 30])
tensor([2, 3, 4])
```

With two tensors of the **same shape**, numbers in matching positions are combined:

```gan
let a = tensor([[1, 2], [3, 4]])
let b = tensor([[5, 6], [7, 8]])
print(a + b)
print(a * b)
```

```output
tensor([[6, 8], [10, 12]])
tensor([[5, 12], [21, 32]])
```

A tensor never changes. Each operation gives you a **new** tensor.

## sum

`sum` adds up all the numbers:

```gan
print(sum(tensor([[1, 2], [3, 4]])))
```

```output
10
```

## matmul: matrix multiplication

`matmul` is different from `*`. Each number in the answer is one **row** of the first table multiplied
with one **column** of the second, then added up:

```gan
let a = tensor([[1, 2], [3, 4]])
let b = tensor([[5, 6], [7, 8]])
print(matmul(a, b))
```

```output
tensor([[19, 22], [43, 50]])
```

The top-left answer is `1*5 + 2*7`, which is 19.

The first table must have as many **columns** as the second has **rows**. If not, GAN tells you both
shapes:

```gan
let x = tensor([[1, 2, 3]])
let w = tensor([[4, 5, 6]])
print(matmul(x, w))
```

```output
Runtime error on line 3: matmul shapes do not fit: [1, 3] and [1, 3] (the first has 3 columns, so the second needs 3 rows, but it has 1)
    3 | print(matmul(x, w))
```

## A first taste of a neuron

The basic part of a neural network multiplies each input by a **weight**, adds the results, then adds
a **bias**. With tensors that is one line:

```gan
let inputs = tensor([[2, 3]])          // 1 row, 2 inputs
let weights = tensor([[10], [100]])    // 2 rows, 1 column
let bias = 1
print(matmul(inputs, weights) + bias)
```

```output
tensor([[321]])
```

Here the weights are chosen by hand. A model **learns** them from examples, which is what later
versions of GAN are meant to do.

## Reading from a tensor

Use `[ ]`, counting from 0. In a table, `t[1]` is the second row:

```gan
let t = tensor([[1, 2, 3], [4, 5, 6]])
print(t[1])
print(t[1][2])
```

```output
tensor([4, 5, 6])
6
```

You can read from a tensor, but not change it.

## transpose, mean and minus

`transpose` turns rows into columns. `mean` is the average. A minus sign flips every sign:

```gan
let t = tensor([[1, 2, 3], [4, 5, 6]])
print(transpose(t))
print(mean(t))
print(-t)
```

```output
tensor([[1, 4], [2, 5], [3, 6]])
3.5
tensor([[-1, -2, -3], [-4, -5, -6]])
```

## Filled and random tensors

Give a **shape** to build a new tensor:

```gan
print(zeros([2, 2]))
print(ones([3]))
```

```output
tensor([[0, 0], [0, 0]])
tensor([1, 1, 1])
```

A model starts with random numbers. `random(shape)` gives numbers from 0 up to (not including) 1.
`randomNormal(shape)` gives numbers centred on 0, mostly close to it. They are different on every
run, unless you set a **seed**. The same seed gives the same numbers again, so you can repeat an
experiment:

```gan
seed(42)
let a = random([2, 3])
seed(42)
let b = random([2, 3])
print(a == b)
print(shape(a))
```

```output
true
[2, 3]
```

## Functions for neural networks

These work on one number or on every number of a tensor:

| Function | What it does |
|---|---|
| `relu(x)` | keeps positive numbers, turns negative ones into 0 |
| `sigmoid(x)` | squeezes any number into the range 0 to 1 |
| `tanh(x)` | squeezes any number into the range -1 to 1 |
| `exp(x)` | e to the power x |
| `log(x)` | the natural logarithm; x must be greater than 0 |

```gan
print(relu(tensor([-2, 0, 3])))
print(sigmoid(tensor([0])))
print(tanh(0))
print(exp(0) + log(1))
```

```output
tensor([0, 0, 3])
tensor([0.5])
0
1
```

## Exercises

1. Make a tensor of the numbers 1 to 4 and print it doubled.
2. Print the sum of `tensor([[5, 5], [5, 5]])`.
3. Use `matmul` to multiply `[[1, 2, 3]]` with `[[1], [1], [1]]`. Guess the answer first.
4. A model guessed `tensor([2, 4, 9])` but the right answers are `tensor([2, 5, 7])`. Print the
   **mean squared error**: subtract, multiply the result by itself, then take the `mean`.

Next: [Lesson 15: How a model learns](15-learning.md)
