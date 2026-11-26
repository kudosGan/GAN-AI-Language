# Lesson 11: Types

Every value has a **type**: what kind of thing it is.

| Type | Examples |
|---|---|
| `number` | `5`, `3.14`, `-2` |
| `text` | `"hello"` |
| `bool` | `true`, `false` |
| `list` | `[1, 2, 3]` |

`typeof(x)` tells you the type of a value:

```gan
print(typeof(5))
print(typeof("hello"))
print(typeof(true))
print(typeof([1, 2]))
print(typeof(2 > 1))
```

```output
number
text
bool
list
bool
```

## Adding a type note

So far a variable could hold anything, and change its mind:

```gan
let x = 5
x = "now I am text"
print(x)
```

```output
now I am text
```

That is flexible, and it is also how bugs sneak in. A **type note** after a colon makes GAN check:

```gan
let age: number = 30
print(age)
age = 31
print(age)
age = "thirty-two"
```

```output
30
31
Runtime error on line 5: variable "age" must be number but got text
    5 | age = "thirty-two"
```

The error says exactly what was expected and what arrived.

## Notes are optional

Add them where a mistake would hurt. Variables without a note stay flexible. You can use `final` too:

```gan
final name: text = "GAN"
let items: list = [1, 2]
let ready: bool = true
print(name)
print(items)
print(ready)
```

```output
GAN
[1, 2]
true
```

## `any`

`any` accepts every type. It is for saying "I thought about the type, and it can be anything":

```gan
let anything: any = 5
anything = "text is fine too"
print(anything)
```

```output
text is fine too
```

## Types on functions

Put a note on each parameter, and after the `)` for the answer:

```gan
func add(a: number, b: number): number {
    return a + b
}
print(add(2, 3))
```

```output
5
```

GAN checks the arguments when the function is called:

```gan
func double(x: number): number {
    return x * 2
}
print(double(4))
print(double("4"))
```

```output
8
Runtime error on line 5: parameter "x" of "double" must be number but got text
    5 | print(double("4"))
```

And it checks the answer when the function returns:

```gan
func half(x: number): number {
    return "half of " + x
}
print(half(8))
```

```output
Runtime error on line 4: the value returned by "half" must be number but got text
    4 | print(half(8))
```

A function that promises an answer must give one:

```gan
func broken(): number {
    print("I forgot to return")
}
print(broken())
```

```output
I forgot to return
Runtime error on line 4: function "broken" must return number but returned nothing
    4 | print(broken())
```

## An unknown type

```gan
let x: integer = 5
```

```output
Syntax error on line 1: unknown type "integer" (the types are: number, text, bool, list, any)
    1 | let x: integer = 5
```

## When are types checked?

**While the program runs**, at the moment a value is stored or passed. A mistake on a line that never
runs is never found. So types help most when you test your program with real inputs.

## Exercises

1. Make `let price: number = 5`, then try to set it to text and read the error.
2. Write `shout(word: text): text` that returns the word in upper case with `!` on the end.
3. Write `average(list: list): number`.
4. Use `typeof` to print the type of five different values.

Next: [Lesson 12: When things go wrong](12-errors.md)
