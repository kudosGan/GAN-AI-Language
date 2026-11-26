# Lesson 12: When things go wrong

Every programmer meets errors all day. The skill is reading them calmly. GAN gives you the line number,
shows that line, and says what is wrong in plain words.

## How to read an error

```gan
let a = 1
let b = 2
print(a)
print(c)
```

```output
1
Runtime error on line 4: undefined variable "c"
    4 | print(c)
```

1. **The kind**: `Runtime error`. The program ran, then stopped. (The `1` printed before the problem.)
2. **The line**: `line 4`.
3. **The problem**: `undefined variable "c"`. There is no variable called `c`.
4. **The line itself**, so you can see it.

Start by reading the problem, then look at that line, then look at the lines just above it.

## Syntax errors: GAN cannot read it

These stop the program **before anything runs**.

A missing closing bracket:

```gan
let x = (2 + 3 * 4
print(x)
```

```output
Syntax error on line 2: expected RPAREN but got IDENT "print"
    2 | print(x)
```

The error is reported on line 2, where GAN noticed something was wrong. The real mistake is often on the
line **above** the one reported. Look there too.

A missing `{`:

```gan
let x = 5
if x > 1 print("big")
```

```output
Syntax error on line 2: expected LBRACE but got IDENT "print"
    2 | if x > 1 print("big")
```

A stray character:

```gan
let total = 5 @ 2
```

```output
Syntax error on line 1: unexpected character "@"
    1 | let total = 5 @ 2
```

## Runtime errors: it ran, then hit a problem

Using a name that does not exist (often a typo):

```gan
let total = 10
print(totl)
```

```output
Runtime error on line 2: undefined variable "totl"
    2 | print(totl)
```

Math on the wrong things:

```gan
print("a" - "b")
```

```output
Runtime error on line 1: can only do math on numbers
    1 | print("a" - "b")
```

An `if` with a question that is not true/false:

```gan
let name = "Ann"
if name {
    print("hi")
}
```

```output
Runtime error on line 2: condition must be true or false
    2 | if name {
```

Asking for an item that is not there:

```gan
let list = [1, 2, 3]
print(list[10])
```

```output
Runtime error on line 2: index 10 is out of range (the list has 3 items)
    2 | print(list[10])
```

Calling a function with the wrong number of arguments:

```gan
func add(a, b) {
    return a + b
}
print(add(1))
```

```output
Runtime error on line 4: function "add" expects 2 arguments but got 1
    4 | print(add(1))
```

Calling a function that does not exist, often because it is defined **below** the call:

```gan
print(add(1, 2))
func add(a, b) {
    return a + b
}
```

```output
Runtime error on line 1: undefined function "add"
    1 | print(add(1, 2))
```

Changing a `final`:

```gan
final limit = 10
limit = 20
```

```output
Runtime error on line 2: cannot change final variable "limit"
    2 | limit = 20
```

## Errors inside functions

The line shown is the line **inside** the function where it went wrong, not the line that called it:

```gan
func half(x) {
    return x / y
}
print(half(8))
```

```output
Runtime error on line 2: undefined variable "y"
    2 |     return x / y
```

## Finding a bug

When a program runs but gives a wrong answer, there is no error message to help. Try this:

1. **Print as you go.** Add `print` lines to see what the variables hold.
2. **Make it smaller.** Remove everything that is not needed until the bug stays and the program is tiny.
3. **Check your assumptions.** Is the loop starting at 0 or 1? Is it `<` or `<=`?
4. **Explain it out loud** to someone (or something) else. Really. It works.

Here is a loop that is off by one. It was meant to add 1 to 10:

```gan
let sum = 0
for n in 1 to 9 {
    sum = sum + n
}
print(sum)
```

```output
45
```

The answer is 45, not 55. Printing `n` inside the loop would show that it stops at 9.

## Exercises

Each of these programs has one mistake. Find it from the error, then fix it.

1. `print(Hello)`
2. `let x = 5` then `x == 6` on the next line, then `print(x)` (the second line does not do what was meant)
3. `for i in 1 to 3 { print(i) ` (nothing after)
4. `let nums = [1, 2, 3]` then `print(nums[3])`
5. `func f(a, b) { return a * b }` then `print(f(2))`

Next: [Lesson 13: Projects](13-projects.md)
