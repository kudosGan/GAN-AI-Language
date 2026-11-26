# Lesson 4: Text

Text in GAN (programmers call it a **string**) is anything between double quotes.

## Joining text

`+` joins pieces of text together:

```gan
let name = "Ann"
print("Hello, " + name)
print("a" + "b" + "c")
```

```output
Hello, Ann
abc
```

If one side of `+` is text, the other side is turned into text too. That is how you put numbers
into a sentence:

```gan
let apples = 3
print("I have " + apples + " apples")
print("Half of 10 is " + 10 / 2)
```

```output
I have 3 apples
Half of 10 is 5
```

Careful with order: `1 + 2 + "x"` adds first (3) and then joins:

```gan
print(1 + 2 + "x")
print("x" + 1 + 2)
```

```output
3x
x12
```

## Length

`len(text)` counts the characters:

```gan
print(len("hello"))
print(len(""))
```

```output
5
0
```

## Picking out a letter

Count from **0**. The first character is number 0, the second is number 1, and so on:

```gan
let word = "GAN"
print(word[0])
print(word[1])
print(word[len(word) - 1])
```

```output
G
A
N
```

Asking for a position that is not there is an error that says how long the text is:

```gan
print("GAN"[3])
```

```output
Runtime error on line 1: index 3 is out of range (the text has 3 characters)
    1 | print("GAN"[3])
```

## Going through every letter

```gan
for letter in "cat" {
    print(letter)
}
```

```output
c
a
t
```

## Changing the case, and searching

```gan
print(upper("gan"))
print(lower("HELLO"))
print(contains("hello world", "o w"))
print(contains("hello", "xyz"))
```

```output
GAN
hello
true
false
```

## Text and numbers

`str(x)` turns anything into text. `num(text)` turns text into a number:

```gan
print(str(5) + "!")
print(num("3.5") * 2)
```

```output
5!
7
```

```gan
print(num("abc"))
```

```output
Runtime error on line 1: cannot turn "abc" into a number
    1 | print(num("abc"))
```

## Text cannot be changed in place

You can read a letter, but not replace one. Build a new piece of text instead:

```gan
let word = "cat"
word[0] = "b"
```

```output
Runtime error on line 2: text cannot be changed; make a new piece of text instead
    2 | word[0] = "b"
```

```gan
let word = "cat"
let newWord = "b" + word[1] + word[2]
print(newWord)
```

```output
bat
```

## Comparing text

```gan
print("apple" == "apple")
print("apple" == "Apple")
print("a" != "b")
```

```output
true
false
true
```

Capital letters count: `"apple"` and `"Apple"` are different.

## Exercises

1. Print `Hello, ` followed by your name, using a variable.
2. Make a sentence that includes a number variable, such as your age.
3. Print the first and last letters of a word you choose.
4. Print a word in upper case followed by its length.
5. Use a `for` loop to print each letter of your name on its own line.

Next: [Lesson 5: Making decisions](05-decisions.md)
