# Lesson 17: Using Python from GAN

Python has a very large collection of libraries. GAN can call them with `python`.

You need **Python** installed. Check with `python --version` in a terminal. The rest of GAN works
without Python.

## python

Give the **full name** of a Python function as text, then its arguments:

```gan
print(python("math.factorial", 5))
print(python("string.capwords", "hello gan"))
print(python("builtins.sorted", [3, 1, 2]))
```

```output
120
Hello Gan
[1, 2, 3]
```

- The name is `module.function`, such as `math.sqrt` or `numpy.linalg.inv`.
- Python's own functions such as `sorted` and `len` are in the module `builtins`.
- GAN imports the module for you.

## What can travel

Numbers, text, true/false, lists and tensors go both ways. Whole numbers arrive in Python as
integers. GAN numbers are decimals, so a very large Python integer may be rounded.

Other results, such as a Python dictionary, set or object, give an error for now. A function that
returns nothing, like `builtins.print`, can only be used on its own line.

Named arguments such as `axis=0` are not supported yet.

## Tensors and NumPy

If **NumPy** is installed, a tensor arrives in Python as a NumPy array, and an array of numbers comes
back as a tensor. These examples are not checked automatically, because NumPy may not be installed:

```gan
print(python("numpy.sort", tensor([3, 1, 2])))
print(python("numpy.linalg.det", tensor([[2, 0], [0, 3]])))
```

A **learnable** tensor is refused: Python cannot pass GAN's gradients back, so the training would
silently go wrong.

## Full power

`python` has **no list of allowed functions**. Your program can do whatever Python can do with your
account, **including deleting files and running commands**. That is why you should:

- read a program written by someone else before you run it,
- never build the **function name** or a command from text you did not write, such as text from
  `input()`, a web page or an AI model's answer.

For example, never do this: `python(ask("Which function should I call?"), "...")`. Whatever the model
says would run on your computer.

## Good to know

- GAN starts Python the first time you call `python`, and keeps it running until your program ends.
- There is no time limit. Press **Ctrl+C** to stop a program that is stuck.
- Only installed modules and Python's own modules are found. A `.py` file in your current folder is
  **not** loaded, so that a stray file cannot run by accident.

## When it goes wrong

```gan
print(python("math.no_such_function", 1))
```

```output
Runtime error on line 1: Python has no function "math.no_such_function"
    1 | print(python("math.no_such_function", 1))
```

An error inside Python is shown with its Python name:

```gan
print(python("math.sqrt", -1))
```

```output
Runtime error on line 1: Python error in math.sqrt: ValueError: math domain error
    1 | print(python("math.sqrt", -1))
```

## Exercises

1. Use `math.gcd` to find the greatest common divisor of 48 and 18.
2. Use `builtins.sum` to add up the list `[10, 20, 30]`.

Next: [Lesson 18: Tools](18-tools.md)
