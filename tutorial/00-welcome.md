# Lesson 0: Welcome to GAN

GAN is a small programming language. You write a program in a plain text file that ends in
`.gan`, and the `gan` command runs it. GAN is made to be easy to read and easy to learn: it borrows
the simple feel of Python, `func` from Go, and `final` from Dart.

You do not need to know any programming. This course goes one small idea at a time, and every
example in it was run by the real GAN, so what you see here is what you will get.

## Getting GAN

Run **GAN-Setup.exe** and click **Install**. Then open a **new** terminal window and check it worked:

```
gan --version
```

You should see `GAN 0.3 (beta)`. If you are building GAN from its source instead, `go build -o gan.exe .`
makes the same `gan.exe`.

## Your first program

Make a file called `hello.gan` with this one line in it:

```gan
print("Hello, world")
```

```output
Hello, world
```

Run it from the terminal:

```
gan hello.gan
```

`print` shows something on the screen. The text goes inside quotes, and the quotes are not shown.

## When something is wrong

GAN tells you which line has the problem, and shows that line:

```gan
print("Hello, world"
```

```output
Syntax error on line 1: unexpected end of file, expected RPAREN
    1 | print("Hello, world"
```

Read errors from the top: the first part says what kind of problem it is, then the line number, then
what GAN expected. Mistakes are normal. Every programmer, every day, reads error messages.

There are two kinds:

- **Syntax error**: the program is not written in a way GAN understands. Nothing runs.
- **Runtime error**: the program started, then hit a problem (like using a name that does not
  exist). Everything before the problem has already run.

## How this course works

- Each lesson has examples. **Type them in yourself** instead of copying: you learn more.
- Gray boxes labelled `output` show exactly what the example prints.
- At the end of each lesson there are exercises. Try them before looking at `solutions.md`.

Next: [Lesson 1: Printing and comments](01-printing.md)
