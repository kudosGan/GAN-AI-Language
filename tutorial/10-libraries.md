# Lesson 10: Libraries

A **library** is a file full of ready-made functions that you can borrow. GAN comes with three, and you
can write your own. You load one with `import`.

## Using a library

```gan
import "math"

print(math.PI)
print(math.square(5))
print(math.factorial(5))
```

```output
3.141592653589793
25
120
```

`import "math"` loads the library and gives it the name `math`. After that:

- `math.PI` is a **value** inside the library.
- `math.square(5)` is a **function** inside the library. The dot means "the one that lives in math".

The name before the dot says whose function it is, so two libraries can have functions with the same name
without clashing.

## The math library

```gan
import "math"

print(math.cube(3))
print(math.sign(-12))
print(math.clamp(15, 0, 10))
print(math.isEven(10))
print(math.isOdd(10))
print(math.isPrime(97))
print(math.gcd(48, 18))
print(math.average([2, 4, 6]))
print(math.maxOf([3, 9, 4]))
print(math.minOf([3, 9, 4]))
print(math.sumOf([1, 2, 3, 4]))
```

```output
27
-1
10
true
false
true
6
4
9
3
10
```

| Function | What it gives back |
|---|---|
| `math.square(x)`, `math.cube(x)` | x times itself, two or three times |
| `math.sign(x)` | 1, -1 or 0 |
| `math.clamp(x, low, high)` | x, kept between low and high |
| `math.isEven(n)`, `math.isOdd(n)` | true or false |
| `math.isPrime(n)` | true if n is prime |
| `math.factorial(n)` | n times n-1 times ... times 1 |
| `math.gcd(a, b)` | the biggest number that divides both |
| `math.sumOf(list)`, `math.average(list)` | total and average of a list of numbers |
| `math.maxOf(list)`, `math.minOf(list)` | biggest and smallest item |
| `math.PI`, `math.E` | the constants |

## The text library

```gan
import "text"

print(text.repeat("ab", 3))
print(text.join(["a", "b", "c"], "-"))
print(text.reverse("GAN"))
print(text.padLeft("7", 3, "0"))
print(text.isPalindrome("level"))
print(text.isEmpty(""))
```

```output
ababab
a-b-c
NAG
007
true
true
```

## The lists library

```gan
import "lists"

print(lists.range(1, 5))
print(lists.sorted([3, 1, 2]))
print(lists.reverse([1, 2, 3]))
print(lists.first(["x", "y"]))
print(lists.last(["x", "y"]))
```

```output
[1, 2, 3, 4, 5]
[1, 2, 3]
[3, 2, 1]
x
y
```

`lists.sorted` gives back a **new** sorted list and leaves the original alone.

## Importing several

```gan
import "math"
import "text"

let n = math.factorial(5)
print("5! is " + n)
print(text.padLeft(str(n), 6, "."))
```

```output
5! is 120
...120
```

Importing the same library twice is harmless.

## When a library cannot be found

```gan
import "nothing"
```

```output
Runtime error on line 1: cannot find library "nothing" (looked in: the libs folder next to gan.exe, the libs folder next to your program, and your program's folder)
    1 | import "nothing"
```

```gan
import "math"
print(math.nothing)
```

```output
Runtime error on line 2: module "math" has no variable "nothing"
    2 | print(math.nothing)
```

## Writing your own library

Any `.gan` file can be a library. Make a file called `shapes.gan`:

```gan
final SIDES_IN_A_SQUARE = 4

func squareArea(side) {
    return side * side
}

func squarePerimeter(side) {
    return side * SIDES_IN_A_SQUARE
}
```

Put it in the same folder as the program that uses it. Because its name has a `.gan` ending,
`import` treats it as a **file path**, relative to the file that contains the `import`:

```gan
import "shapes.gan"

print(shapes.squareArea(5))
print(shapes.squarePerimeter(5))
print(shapes.SIDES_IN_A_SQUARE)
```

The name you use is the file name without `.gan`: here `shapes`. A path with folders works too:
`import "tools/shapes.gan"`.

Things worth knowing about your own libraries:

- The library has its **own** variables. A variable called `count` inside it is not the `count` in your program.
- The code at the top of a library runs **once**, at the moment it is first imported. If it has a `print`, you will see it then.
- A library can import other libraries.
- An error inside a library names the library file and shows the line there.
- To stop with a clear message when someone uses your function badly, call `fail("message")`:

```gan
func half(x) {
    if x < 0 {
        fail("half needs a number that is not negative")
    }
    return x / 2
}
print(half(8))
print(half(-2))
```

```output
4
Runtime error on line 3: half needs a number that is not negative
    3 |         fail("half needs a number that is not negative")
```
```

## Where the standard libraries live

The three built-in libraries (`math`, `text`, `lists`) are ordinary `.gan` files in the `libs` folder
next to `gan.exe`. You can open them and read them. They are written in GAN, using only what you
learned in this course. A good way to learn more is to read `libs/math.gan`.

## Exercises

1. Print the factorial of 6 and the greatest common divisor of 84 and 36 using `math`.
2. Use `lists.range` and a loop to print the squares of the numbers 1 to 5.
3. Sort the list `[5, 3, 8, 1]` and print the smallest and largest with `lists.first` and `lists.last`.
4. Make your own library `greetings.gan` with a function `hello(name)`, and use it from another file.
5. Read `libs/text.gan`. Which function would you like to add?

Next: [Lesson 11: Types](11-types.md)
