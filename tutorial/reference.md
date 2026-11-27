# GAN quick reference

A one-page reminder of everything GAN can do. For explanations, see the lessons.

## Programs

- Files end in `.gan`. Run with `gan file.gan`. `gan --version` shows the version.
- A program runs from the top line down. Blank lines and extra spaces do not matter.
- `// a comment` is ignored to the end of the line.

## Values

| Type | Examples | Notes |
|---|---|---|
| number | `5`, `3.14`, `-2` | all numbers are decimals; whole ones print without `.0` |
| text | `"hello"` | double quotes only; no way to put a `"` inside; may span lines |
| bool | `true`, `false` | |
| list | `[1, "a", [2]]` | positions start at 0; shared, not copied |
| tensor | `tensor([[1, 2], [3, 4]])` | grid of numbers for AI; never changes |

## Variables

```gan
let x = 5          // make a variable
var y = 6          // the same as let
final z = 7        // can never change
x = x + 1          // change an existing variable
let n: number = 1  // with an optional type note
```

Names: letters, digits and `_`, not starting with a digit; capitals matter.

## Operators (tightest first)

| | |
|---|---|
| `( )` | grouping |
| `-x`, `x[i]`, `m.name`, `f(x)` | minus sign, list/text position, module member, call |
| `*  /  %` | multiply, divide (always decimal), remainder |
| `+  -` | add or join; subtract |
| `==  !=  <  >  <=  >=` | comparisons, give true/false |
| `not` | flips true/false |
| `and` | both true (stops early) |
| `or` | one true (stops early) |

`+` joins when either side is text: `"n = " + 5` is `"n = 5"`. `+` on two lists joins them.

## Decisions and loops

```gan
if a > b {
    ...
} else if a == b {
    ...
} else {
    ...
}

while condition {
    ...
}

for i in 1 to 10 {      // both ends included
    ...
}

for item in list {      // each item of a list, or each letter of text
    ...
}
```

Conditions must be true or false. Inside a loop, `break` stops the loop and `continue` skips to its next pass.
Both affect only the nearest loop.

## Functions

```gan
func add(a, b) {
    return a + b
}
func typed(a: number, b: number): number {
    return a + b
}
print(add(2, 3))
```

- Define a function above the place that calls it.
- A function can read and change top-of-program variables; its own variables vanish when it ends.
- A function with no `return` gives no value (using it as a value is an error).

## Lists and text

```gan
let l = [10, 20, 30]
l[0]              // 10
l[1] = 99         // change an item
len(l)            // 3
push(l, 40)       // add to the end
pop(l)            // remove and give back the last
l + [1]           // a new, longer list
for x in l { ... }
"hello"[1]        // "e"   (text can be read, not changed)
```

## Built-in functions

| Function | Gives |
|---|---|
| `len(x)` | number of items or letters |
| `push(list, item)` | adds to the end (no value) |
| `pop(list)` | removes and gives the last item |
| `contains(list or text, thing)` | true/false |
| `str(x)` | x as text |
| `num(text)` | text as a number |
| `typeof(x)` | `"number"`, `"text"`, `"bool"`, `"list"`, `"tensor"` or `"module"` |
| `input()` / `input("question ")` | one typed line, as text |
| `abs(x)` `floor(x)` `round(x)` `sqrt(x)` | math |
| `pow(base, exp)` `min(a, b)` `max(a, b)` | math |
| `upper(text)` `lower(text)` | change case |
| `fail("message")` | stop the program with your own error |

## Tensors

| Code | Gives |
|---|---|
| `tensor(list)` | a tensor from a list of numbers, or of equal-length rows |
| `shape(t)` | its size in each direction, e.g. `[2, 3]` |
| `sum(t)` | all its numbers added up |
| `matmul(a, b)` | matrix multiplication of two tables; columns of `a` must equal rows of `b` |
| `a + b`, `a - b`, `a * b`, `a / b` | position by position (same shape) |
| `t * 10` and the like | the number is used for every position |
| `-t` | every sign flipped |
| `t[0]`, `t[1][2]` | a row or a number, counting from 0 (read only) |
| `transpose(t)` | a table with rows and columns swapped |
| `mean(t)` | the average of all its numbers |
| `zeros(shape)` `ones(shape)` | filled with 0 or 1, e.g. `zeros([2, 3])` |
| `random(shape)` | random numbers from 0 up to (not including) 1 |
| `randomNormal(shape)` | random numbers centred on 0 (normal distribution, spread 1) |
| `seed(n)` | makes the random numbers repeat |
| `relu(x)` `sigmoid(x)` `tanh(x)` `exp(x)` `log(x)` | on a number or every number of a tensor |

A function you write with one of these names is used instead of the built-in one.

## Gradients

| Code | Gives |
|---|---|
| `learnable(t)` | a tensor whose gradient you want |
| `backward(loss)` | works out gradients from a one-number result (no value) |
| `grad(w)` | the gradient of `w` from the last `backward`, same shape as `w` |
| `item(t)` | the plain number inside a one-number tensor |

Gradients work through `+` `-` `*`, minus signs, `matmul`, `sum` and `mean`. Other tensor maths and `[ ]` give
an error on learnable calculations. On those, `sum` and `mean` give a one-number tensor such as `tensor(20)`.
After changing a weight, wrap it again: `w = learnable(w - 0.05 * grad(w))`.

## Types

`number`, `text`, `bool`, `list`, `tensor`, `any`. Notes are optional and checked while the program runs: when a variable is
made or changed, when a function is called, and when it returns.

## Libraries

```gan
import "math"             // from the libs folder
import "tools/mine.gan"   // a file, relative to this one
print(math.square(5))     // function
print(math.PI)            // value
```

| Library | Contains |
|---|---|
| `math` | `PI` `E` `square` `cube` `sign` `clamp` `isEven` `isOdd` `isPrime` `factorial` `gcd` `sumOf` `average` `maxOf` `minOf` |
| `text` | `repeat` `join` `reverse` `isEmpty` `padLeft` `isPalindrome` |
| `lists` | `range` `reverse` `sorted` `first` `last` |

## Errors

`Syntax error on line N` = GAN cannot read the program (nothing runs).
`Runtime error on line N` = it ran, then stopped. The message says what and the line is shown.
An error inside a library says `in name.gan` and shows that library's line.

## Not in GAN yet

Dictionaries, `\n` and other escapes inside text, several things in one `print`, `+=`,
one-line `if`, reading files, concurrency, gradients for `/`, `relu` and other functions.
