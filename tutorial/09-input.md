# Lesson 9: Asking the user

So far your programs know everything in advance. `input` lets a program ask a question and wait
for the person to type an answer.

## input

`input("question")` shows the question, waits for a line of typing, and gives it back **as text**.

The gray `input` box below shows what the person typed:

```gan
let name = input("What is your name? ")
print("Hello, " + name)
```

```input
Ann
```
```output
What is your name? Ann
Hello, Ann
```

Put a space at the end of the question so the answer does not touch it.

## Numbers need num

Whatever is typed arrives as text, even digits. Use `num(...)` to turn it into a number:

```gan
let age = num(input("How old are you? "))
print("Next year you will be " + (age + 1))
```

```input
30
```
```output
How old are you? 30
Next year you will be 31
```

Without `num`, `"30" + 1` would join text instead of adding. Try it and see what you get.

If what was typed is not a number, `num` stops with a clear message:

```gan
let n = num(input("A number: "))
print(n * 2)
```

```input
banana
```
```output
A number: banana
Runtime error on line 1: cannot turn "banana" into a number
    1 | let n = num(input("A number: "))
```

## A small program: a tip calculator

```gan
let bill = num(input("Bill amount: "))
let percent = num(input("Tip percent: "))
let tip = bill * percent / 100
print("Tip: " + tip)
print("Total: " + (bill + tip))
```

```input
80
15
```
```output
Bill amount: 80
Tip percent: 15
Tip: 12
Total: 92
```

## Asking until the answer is good

A `while` loop can keep asking. Here the loop runs until the answer is a number above zero:

```gan
func askPositive(question) {
    while true {
        let n = num(input(question))
        if n > 0 {
            return n
        }
        print("It must be more than zero.")
    }
}

let sides = askPositive("How many sides? ")
print("You chose " + sides)
```

```input
0
-4
6
```
```output
How many sides? 0
It must be more than zero.
How many sides? -4
It must be more than zero.
How many sides? 6
You chose 6
```

(`while true` never stops by itself. Here `return` is what leaves it.)

## input with no question

`input()` with nothing in the brackets just waits and reads:

```gan
print("Type something:")
let line = input()
print("You typed: " + line)
```

```input
hello there
```
```output
Type something:
hello there
You typed: hello there
```

## Exercises

1. Ask for the user's name and favourite colour, then print a sentence with both.
2. Ask for two numbers and print their sum, difference and product.
3. Ask for a temperature in Celsius and print it in Fahrenheit (`c * 9 / 5 + 32`).
4. Ask for a word and print whether it is a palindrome (reads the same backwards). Hint: build the reversed word with a `for` loop, and compare with `==`.

Next: [Lesson 10: Libraries](10-libraries.md)
