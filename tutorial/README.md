# Learn GAN

A step-by-step course for the GAN programming language. No experience needed. Go in order: each lesson uses
what the earlier ones taught.

| # | Lesson | You will learn |
|---|---|---|
| 0 | [Welcome](00-welcome.md) | installing GAN, your first program, reading errors |
| 1 | [Printing and comments](01-printing.md) | `print`, `//` comments |
| 2 | [Variables](02-variables.md) | `let`, `var`, `final`, changing values |
| 3 | [Numbers and math](03-math.md) | `+ - * / %`, order of operations, rounding |
| 4 | [Text](04-text.md) | joining, `len`, letters, `upper`, `str`, `num` |
| 5 | [Making decisions](05-decisions.md) | `if`, `else`, `and` / `or` / `not` |
| 6 | [Loops](06-loops.md) | `for`, `while`, collecting answers |
| 7 | [Functions](07-functions.md) | `func`, `return`, recursion, scope |
| 8 | [Lists](08-lists.md) | `[1, 2, 3]`, `push`, `pop`, looping over lists |
| 9 | [Asking the user](09-input.md) | `input`, interactive programs |
| 10 | [Libraries](10-libraries.md) | `import`, `math`, `text`, `lists`, writing your own |
| 11 | [Types](11-types.md) | `let x: number`, type-checked functions |
| 12 | [When things go wrong](12-errors.md) | reading errors, finding bugs |
| 13 | [Projects](13-projects.md) | FizzBuzz, primes, guessing game, calculator... |
| 14 | [Tensors](14-tensors.md) | `tensor`, `shape`, `sum`, `matmul`: the first step into AI |
| 15 | [How a model learns](15-learning.md) | `learnable`, `backward`, `grad`, training a tiny model |

Also here:

- [Solutions](solutions.md) to every exercise (try first!)
- [Quick reference](reference.md): the whole language on one page

## For the person maintaining this course

Every example in these lessons is run by the real `gan.exe`:

```
python tutorial\check.py          check every example (fails if GAN prints something different)
python tutorial\check.py --fill   write the real output into the output boxes
```

Write an example as a `gan` block followed by an `output` block (and an `input` block in between if the
program uses `input()`), run `--fill`, and read what GAN really printed. A `gan` block with no output block
is only shown, never run.
