# Lesson 7: Functions

A **function** is a named recipe. You write it once, then use it as often as you like.

## Making and using a function

```gan
func sayHello() {
    print("Hello!")
}

sayHello()
sayHello()
```

```output
Hello!
Hello!
```

`func` makes the function. The code in `{ }` is what it does. Writing its name with `( )` after it
runs it. That is called **calling** the function.

## Parameters

Parameters let you give a function something to work with:

```gan
func greet(name) {
    print("Hello, " + name)
}

greet("Ann")
greet("Bo")
```

```output
Hello, Ann
Hello, Bo
```

`name` is a **parameter**. The `"Ann"` you pass in is an **argument**. With several parameters,
separate them with commas, and give the same number of arguments:

```gan
func rectangleArea(width, height) {
    print(width * height)
}
rectangleArea(3, 4)
```

```output
12
```

```gan
func rectangleArea(width, height) {
    print(width * height)
}
rectangleArea(3)
```

```output
Runtime error on line 4: function "rectangleArea" expects 2 arguments but got 1
    4 | rectangleArea(3)
```

## Giving an answer back with return

`return` ends the function and hands back a value. The call then **becomes** that value, so you can
use it in math or store it:

```gan
func add(a, b) {
    return a + b
}

print(add(2, 3))
print(add(1, 2) * 3)
let total = add(10, 20)
print(total)
```

```output
5
9
30
```

## Functions that decide

```gan
func biggest(a, b) {
    if a > b {
        return a
    }
    return b
}
print(biggest(3, 9))
print(biggest(12, 4))
```

```output
9
12
```

`return` can also leave a loop early. This finds the first number whose square is bigger than 50:

```gan
func firstSquareOver(limit) {
    for q in 1 to 100 {
        if q * q > limit {
            return q
        }
    }
    return 0
}
print(firstSquareOver(50))
```

```output
8
```

## A function can call itself

This is called **recursion**. The factorial of 5 is 5 * 4 * 3 * 2 * 1:

```gan
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
print(factorial(5))
```

```output
120
```

Every recursive function needs a way to stop (here, `n <= 1`). Without one, GAN stops it for you:

```gan
func forever(n) {
    return forever(n + 1)
}
print(forever(1))
```

```output
Runtime error on line 2: too many nested calls
    2 |     return forever(n + 1)
```

## Variables inside functions

A function has its **own** variables. They are made fresh each time it runs and disappear when it ends.
A function can see variables made at the top of the program, and can change them:

```gan
let count = 0

func addOne() {
    count = count + 1
    print(count)
}

addOne()
addOne()
print(count)
```

```output
1
2
2
```

But a variable made **inside** a function is not visible outside it:

```gan
func makeSecret() {
    let secret = 42
    return secret
}
print(makeSecret())
print(secret)
```

```output
42
Runtime error on line 6: undefined variable "secret"
    6 | print(secret)
```

A parameter with the same name as an outer variable hides it, and leaves it alone:

```gan
let n = 100
func addOne(n) {
    return n + 1
}
print(addOne(1))
print(n)
```

```output
2
100
```

## Order matters

A function exists once its `func` line has run, so define it **above** the place that calls it.

```gan
print(twice(4))

func twice(x) {
    return x * 2
}
```

```output
Runtime error on line 1: undefined function "twice"
    1 | print(twice(4))
```

## Functions with no answer

If a function has no `return`, it just does its job. Using it as a value is an error:

```gan
func shout(word) {
    print(upper(word))
}
shout("hey")
print(shout("hey"))
```

```output
HEY
HEY
Runtime error on line 5: function "shout" did not return a value
    5 | print(shout("hey"))
```

## Exercises

1. Write `square(x)` that returns `x * x`.
2. Write `isEven(n)` that returns `true` or `false`.
3. Write `average(a, b, c)`.
4. Write `countdown(n)` that prints n down to 1 (use a `while`).
5. Write `power(base, exp)` using a loop (do not use the built-in `pow`).
6. Write `fib(n)` that returns the n-th Fibonacci number (1, 1, 2, 3, 5, 8...) using recursion.

Next: [Lesson 8: Lists](08-lists.md)
