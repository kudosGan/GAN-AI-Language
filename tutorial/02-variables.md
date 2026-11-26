# Lesson 2: Variables

A **variable** is a name for a value, like a labelled box. Make one with `let`:

```gan
let age = 12
print(age)
```

```output
12
```

`let age = 12` means: make a box called `age` and put `12` in it. Later, `age` stands for `12`.

## Changing a variable

Use `=` on its own to put a new value in an existing box:

```gan
let score = 10
print(score)
score = 15
print(score)
score = score + 5
print(score)
```

```output
10
15
20
```

`score = score + 5` works because the right side is worked out first (15 + 5), and then the answer
goes into `score`.

`let` **makes** a variable. `=` **changes** one that exists. Changing one that does not exist is an
error, which catches typos:

```gan
let total = 0
totl = total + 5
```

```output
Runtime error on line 2: undefined variable "totl"
    2 | totl = total + 5
```

## Variables can hold different things

```gan
let name = "Ann"
let height = 1.62
let happy = true
print(name)
print(height)
print(happy)
```

```output
Ann
1.62
true
```

A variable can hold text, a number, or true/false. (It can also hold a list, in Lesson 8.)

## Names

A name can use letters, digits and `_`, and cannot start with a digit. Capital letters matter:
`Score` and `score` are different names. Choose names that say what is inside:
`ageInYears` is better than `x`.

```gan
let first_name = "Ann"
let score2 = 9
print(first_name)
print(score2)
```

```output
Ann
9
```

## var, and final

`var` means exactly the same as `let` (it is there for people who know Dart or JavaScript).

`final` makes a variable that can **never** change. Use it for values that should stay put:

```gan
final daysInWeek = 7
var weeks = 2
print(daysInWeek * weeks)
weeks = 3
print(daysInWeek * weeks)
daysInWeek = 8
```

```output
14
21
Runtime error on line 6: cannot change final variable "daysInWeek"
    6 | daysInWeek = 8
```

## Where a variable lives (scope)

A variable made inside `{ ... }` belongs to that block and disappears at the closing `}`. You will see
blocks in Lesson 5. A variable made at the top of the program can be seen everywhere.

## Exercises

1. Make variables for your name and your age, and print both.
2. Make `price = 4` and `count = 3`. Print `price * count`.
3. Make a `final` variable `taxRate = 2`. Try to change it and read the error.
4. What does `let a = 5` then `a = a * 2` then `print(a)` show? Work it out first, then run it.

Next: [Lesson 3: Numbers and math](03-math.md)
