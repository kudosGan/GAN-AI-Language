# Lesson 18: Tools

An AI **agent** is a model that can use **tools** to do a task. GAN is getting agents step by step.
This lesson covers the first piece, the tool. **GAN has no agents yet**, so for now a tool is used by
your own program.

## A tool is a function with a description

```gan
tool add(a: number, b: number): number "Adds two numbers" {
    return a + b
}

print(add(12, 30))
```

```output
42
```

Compare it with `func`:

- It starts with the word `tool`.
- A **description** in quotes comes before the `{`. Later, a model will read it to decide when to use
  the tool, so say clearly what the tool does.
- **Every parameter needs a type**: `number`, `text`, `bool` or `list`. These are the kinds of value
  a model can send.
- The return type after the `)` is optional, as in a function.

A tool does not have to give a value back:

```gan
tool greet(name: text) "Prints a greeting" {
    print("Hello, " + name)
}

greet("GAN")
```

```output
Hello, GAN
```

The types are checked when the tool is called:

```gan
tool half(n: number): number "Halves a number" {
    return n / 2
}
print(half("nine"))
```

```output
Runtime error on line 4: parameter "n" of "half" must be number but got text
    4 | print(half("nine"))
```

## Mistakes are found before the program runs

```gan
print("start")
tool add(a, b: number) "Adds two numbers" {
    return a + b
}
```

```output
Syntax error on line 2: tool "add": parameter "a" needs a type, so a model knows what to send (number, text, bool or list)
    2 | tool add(a, b: number) "Adds two numbers" {
```

```gan
tool add(a: number, b: number) {
    return a + b
}
```

```output
Syntax error on line 1: tool "add" needs a description in quotes before the {, to tell a model what it does
    1 | tool add(a: number, b: number) {
```

## What comes next

Agents are planned: you will list the tools an agent may use, and the model will be able to run only
those. A tool can do anything your program can, including calling `python`, so think about what you
give to a model.

## Exercises

1. Write a tool `area(width: number, height: number)` that returns the area, with a description.
2. Write a tool `shout(words: text)` that prints the words in capitals.

Back to the [course index](README.md)
