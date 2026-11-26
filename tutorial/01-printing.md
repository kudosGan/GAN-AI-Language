# Lesson 1: Printing and comments

## print

`print(...)` shows one thing on its own line. It can show text, numbers, and true/false:

```gan
print("Hello")
print(42)
print(3.14)
print(true)
```

```output
Hello
42
3.14
true
```

Text must be in double quotes `"like this"`. Numbers and `true` / `false` are not.

Each `print` starts a new line, and a program runs from the top line to the bottom line.

```gan
print("first")
print("second")
print("third")
```

```output
first
second
third
```

`print` takes exactly one thing. To show several things on one line, join them into one piece of
text (you will learn how in Lesson 4).

## Math works inside print

```gan
print(2 + 3)
print(10 / 4)
```

```output
5
2.5
```

## Comments

Anything after `//` on a line is a **comment**. GAN ignores it. Comments are notes for people:
what the code is for, or why you wrote it that way.

```gan
// this whole line is a note
print("start")   // a note after some code
// print("this line is switched off")
print("end")
```

```output
start
end
```

Putting `//` in front of a line to switch it off is a handy trick while you experiment.

## Long text

Text can run over several lines. Everything between the quotes is kept, including the line breaks:

```gan
print("one
two
three")
```

```output
one
two
three
```

There is no way to put a double quote **inside** text yet, because a quote always ends the text.

## Common mistakes

A missing quote:

```gan
print("oops)
```

```output
Syntax error on line 1: unterminated string (a string needs a closing ")
    1 | print("oops)
```

A word without quotes is treated as a name (you will meet names in Lesson 2), so this fails:

```gan
print(hello)
```

```output
Runtime error on line 1: undefined variable "hello"
    1 | print(hello)
```

## Exercises

1. Print your name.
2. Print three lines: a greeting, your favourite food, and a number.
3. Print the result of `7 * 6`.
4. Write a program with a comment on every line, then switch one `print` off with `//`.

Next: [Lesson 2: Variables](02-variables.md)
