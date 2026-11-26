# Lesson 5: Making decisions

So far every line runs, always. Now your program can choose.

## True and false

A comparison asks a question and answers `true` or `false`:

| Write | Asks |
|---|---|
| `a == b` | are they equal? |
| `a != b` | are they different? |
| `a < b` | is a smaller? |
| `a > b` | is a bigger? |
| `a <= b` | smaller or equal? |
| `a >= b` | bigger or equal? |

```gan
print(5 > 3)
print(5 == 6)
print(5 != 6)
print(4 <= 4)
print("cat" == "cat")
```

```output
true
false
true
true
true
```

Careful: `=` **puts a value in a variable**, `==` **asks if two things are equal**.

## if

`if` runs a block of code only when the question is `true`. The block goes between `{` and `}`:

```gan
let age = 20
if age >= 18 {
    print("adult")
}
print("done")
```

```output
adult
done
```

If the answer is `false`, the block is skipped:

```gan
let age = 10
if age >= 18 {
    print("adult")
}
print("done")
```

```output
done
```

## else

`else` is what to do when the answer is `false`:

```gan
let temperature = 8
if temperature < 15 {
    print("take a coat")
} else {
    print("no coat needed")
}
```

```output
take a coat
```

## else if

Check several things in order. The first one that is true wins, and the rest are skipped:

```gan
func grade(score) {
    if score >= 90 {
        return "A"
    } else if score >= 80 {
        return "B"
    } else if score >= 70 {
        return "C"
    } else {
        return "F"
    }
}
print(grade(95))
print(grade(85))
print(grade(72))
print(grade(40))
```

```output
A
B
C
F
```

(`func` is Lesson 7. For now, read it as "a named recipe that gives back an answer".)

## and, or, not

Combine questions:

```gan
let age = 25
print(age >= 18 and age < 65)
print(age < 13 or age > 65)
print(not (age == 25))
```

```output
true
false
false
```

- `a and b` is true only when **both** are true.
- `a or b` is true when **at least one** is true.
- `not a` flips true to false and false to true.

They are checked from the left, and stop as soon as the answer is known, so in
`false and something` the `something` is never looked at.

`and` goes before `or`, like `*` before `+`. When unsure, add parentheses.

## true and false are values

```gan
let ready = true
if ready {
    print("go!")
}
if not ready {
    print("wait")
}
```

```output
go!
```

## Blocks have their own variables

A variable made inside `{ }` disappears at the end of the block:

```gan
if true {
    let secret = 42
    print(secret)
}
print(secret)
```

```output
42
Runtime error on line 5: undefined variable "secret"
    5 | print(secret)
```

## An if needs a true/false answer

```gan
if 5 {
    print("hi")
}
```

```output
Runtime error on line 1: condition must be true or false
    1 | if 5 {
```

GAN does not guess that `5` means "yes". Ask a real question: `if x > 0 { ... }`.

## Exercises

1. Print `"positive"`, `"negative"` or `"zero"` for a number in a variable.
2. Print `"even"` or `"odd"` using `%`.
3. A ticket costs 5 for children under 13, 10 for adults up to 64 and 7 for 65 and over. Print the price for an age in a variable.
4. Check that a number is between 1 and 10 using `and`.

Next: [Lesson 6: Loops](06-loops.md)
