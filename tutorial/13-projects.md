# Lesson 13: Projects

Time to build things. Each project uses ideas from several lessons. Read the idea, try to write it yourself
first, and only then look at the answer.

## Project 1: FizzBuzz

The classic. Count from 1 to 15. For a multiple of 3 print `Fizz`, for a multiple of 5 print `Buzz`, for
a multiple of both print `FizzBuzz`, and otherwise print the number.

Ideas used: `for`, `if`/`else if`, `%`, `and`.

```gan
for n in 1 to 15 {
    if n % 15 == 0 {
        print("FizzBuzz")
    } else if n % 3 == 0 {
        print("Fizz")
    } else if n % 5 == 0 {
        print("Buzz")
    } else {
        print(n)
    }
}
```

```output
1
2
Fizz
4
Buzz
Fizz
7
8
Fizz
Buzz
11
Fizz
13
14
FizzBuzz
```

Why is the `% 15` check first? Because a multiple of 15 is also a multiple of 3 and of 5, and the first
true question wins.

## Project 2: Prime numbers

Print all the primes up to 30, collecting them in a list first.

Ideas used: libraries, lists, `for`.

```gan
import "math"

let primes = []
for n in 2 to 30 {
    if math.isPrime(n) {
        push(primes, n)
    }
}
print(primes)
print("There are " + len(primes) + " primes up to 30")
```

```output
[2, 3, 5, 7, 11, 13, 17, 19, 23, 29]
There are 10 primes up to 30
```

## Project 3: Temperature table

Print Celsius from 0 to 100 in steps of 20, with the Fahrenheit value next to it.

Ideas used: `for`, math, joining text. (GAN's `for` counts by one, so the loop counts 0 to 5 and multiplies.)

```gan
for step in 0 to 5 {
    let c = step * 20
    let f = c * 9 / 5 + 32
    print(c + " C = " + f + " F")
}
```

```output
0 C = 32 F
20 C = 68 F
40 C = 104 F
60 C = 140 F
80 C = 176 F
100 C = 212 F
```

## Project 4: Counting vowels

Count the vowels in a sentence.

Ideas used: text, `for` over letters, `contains`, `lower`.

```gan
func countVowels(sentence) {
    let count = 0
    for letter in lower(sentence) {
        if contains("aeiou", letter) {
            count = count + 1
        }
    }
    return count
}

print(countVowels("Programming is fun"))
print(countVowels("rhythm"))
```

```output
5
0
```

## Project 5: Guess the number

A program that asks for guesses and says higher or lower. The secret is fixed here (there is no random
yet), and the gray `input` box shows the guesses that were typed.

Ideas used: `while`, `input`, `num`, `if`/`else if`.

```gan
let secret = 37
let tries = 0
let guess = 0

while guess != secret {
    guess = num(input("Your guess? "))
    tries = tries + 1
    if guess < secret {
        print("Higher")
    } else if guess > secret {
        print("Lower")
    }
}
print("Got it in " + tries + " tries!")
```

```input
50
25
37
```
```output
Your guess? 50
Lower
Your guess? 25
Higher
Your guess? 37
Got it in 3 tries!
```

## Project 6: A shopping list

Keep a list of items, and a list of prices side by side. Print each line and the total.

Ideas used: lists, `for` over positions, joining text.

```gan
let items = ["bread", "milk", "eggs"]
let prices = [2.5, 1.2, 3]
let total = 0

for i in 0 to len(items) - 1 {
    print(items[i] + ": " + prices[i])
    total = total + prices[i]
}
print("Total: " + total)
```

```output
bread: 2.5
milk: 1.2
eggs: 3
Total: 6.7
```

## Project 7: Sorting by hand

The `lists` library has `sorted`, but how does it work? This is bubble sort: walk through the list, and
swap any two neighbours that are in the wrong order. Repeat until nothing needs swapping.

Ideas used: nested `while`, list items, `let` inside a loop.

```gan
func bubbleSort(list) {
    let result = list + []
    let n = len(result)
    let i = 0
    while i < n {
        let j = 0
        while j < n - i - 1 {
            if result[j] > result[j + 1] {
                let temp = result[j]
                result[j] = result[j + 1]
                result[j + 1] = temp
            }
            j = j + 1
        }
        i = i + 1
    }
    return result
}

print(bubbleSort([5, 2, 9, 1, 7]))
```

```output
[1, 2, 5, 7, 9]
```

## Project 8: A mini calculator

Ask for two numbers and an operation, then show the answer.

Ideas used: `input`, `num`, `if`/`else if`, functions.

```gan
func calculate(a, op, b) {
    if op == "+" {
        return a + b
    } else if op == "-" {
        return a - b
    } else if op == "*" {
        return a * b
    } else if op == "/" {
        return a / b
    }
    fail("I do not know the operation " + op)
}

let a = num(input("First number: "))
let op = input("Operation (+ - * /): ")
let b = num(input("Second number: "))
print(a + " " + op + " " + b + " = " + calculate(a, op, b))
```

```input
12
*
7
```
```output
First number: 12
Operation (+ - * /): *
Second number: 7
12 * 7 = 84
```

## Your own projects

Ideas to try next:

- A times-table printer that asks which table you want.
- Convert minutes into hours and minutes.
- A word counter: how many words are in a sentence? (Count the spaces and add one.)
- A password checker with rules (long enough, contains a digit...).
- Collect numbers until the user types `done`, then print their average.
- Reverse a sentence word by word.

You now know the basics of GAN: printing, variables, math, text, decisions, loops, functions, lists,
input, libraries and types.

Next: [Lesson 14: Tensors](14-tensors.md), the first step into AI.
