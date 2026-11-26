# Lesson 3: Numbers and math

## The operators

| Write | Meaning |
|---|---|
| `+` | add |
| `-` | subtract |
| `*` | multiply |
| `/` | divide |
| `%` | remainder after dividing |

```gan
print(7 + 2)
print(7 - 2)
print(7 * 2)
print(7 / 2)
print(7 % 2)
```

```output
9
5
14
3.5
1
```

Division always gives a decimal answer: `7 / 2` is `3.5`, not `3`. The remainder `%` is what is left
over: 7 divided by 2 is 3 with 1 left over.

## Order of operations

Multiplication and division happen **before** addition and subtraction, like in school:

```gan
print(2 + 3 * 4)
```

```output
14
```

Use parentheses to choose the order yourself. Whatever is inside them is worked out first:

```gan
print((2 + 3) * 4)
print(10 / (2 + 3))
```

```output
20
2
```

## Negative numbers and decimals

```gan
print(-5 + 8)
print(5 - -3)
print(-(2 + 3))
print(2.5 * 4)
print(0.5 + 0.25)
```

```output
3
8
-5
10
0.75
```

## Rounding noise

Computers store decimals in a way that cannot hold every value exactly, so some answers show a tiny
extra amount. Almost every language does this:

```gan
print(0.1 + 0.2)
```

```output
0.30000000000000004
```

To tidy an answer to two decimal places, multiply, round, then divide:

```gan
print(round((0.1 + 0.2) * 100) / 100)
```

```output
0.3
```

Read it from the inside out: add `0.1 + 0.2`, multiply by 100 (30.000000000000004), round to a whole
number (30), then divide by 100 (0.3).

## Built-in math helpers

GAN comes with some ready-made functions:

```gan
print(abs(-4))
print(min(3, 8))
print(max(3, 8))
print(sqrt(16))
print(pow(2, 10))
print(round(2.6))
print(floor(2.6))
```

```output
4
3
8
4
1024
3
2
```

`round` goes to the nearest whole number. `floor` always goes down. `pow(2, 10)` is 2 to the power 10.

## Even or odd

The remainder is great for checking things. A number is even when dividing by 2 leaves nothing:

```gan
print(10 % 2)
print(7 % 2)
```

```output
0
1
```

## Errors

You cannot divide by zero:

```gan
print(5 / 0)
```

```output
Runtime error on line 1: division by zero
    1 | print(5 / 0)
```

## Exercises

1. A pizza costs 12 and 4 friends split it evenly. Print the price each pays.
2. Print the average of 80, 92 and 71.
3. Print the area of a circle with radius 3, using `3.14159 * r * r` with a variable `r`.
4. Print whether 144 is even (look at the remainder) and print `sqrt(144)`.
5. Predict `2 + 3 * 4 - 1` and `(2 + 3) * (4 - 1)`, then check.

Next: [Lesson 4: Text](04-text.md)
