# Lesson 6: Loops

A **loop** repeats code. Computers are good at doing the same thing many times.

## for

Count from one number to another. **Both ends are included.**

```gan
for i in 1 to 5 {
    print(i)
}
```

```output
1
2
3
4
5
```

`i` is a variable that GAN fills in on each pass: 1, then 2, then 3... You can call it anything.

The start and end can be any math:

```gan
let start = 3
for n in start to start + 2 {
    print("n is " + n)
}
```

```output
n is 3
n is 4
n is 5
```

If the start is bigger than the end, the loop does not run at all:

```gan
for i in 5 to 1 {
    print("never printed")
}
print("done")
```

```output
done
```

## Adding things up

A very common pattern: a variable that collects an answer as the loop goes round.

```gan
let sum = 0
for n in 1 to 100 {
    sum = sum + n
}
print(sum)
```

```output
5050
```

## while

`while` repeats as long as a question stays `true`. Use it when you do not know in advance how
many times:

```gan
let count = 1
while count <= 3 {
    print("count is " + count)
    count = count + 1
}
print("finished")
```

```output
count is 1
count is 2
count is 3
finished
```

**The loop must change something**, or the question stays true forever and the program never stops.
The line `count = count + 1` is what makes this loop finish. If you ever start a program that will not
stop, press **Ctrl+C** in the terminal.

## Doubling until big

```gan
let money = 1
let years = 0
while money < 100 {
    money = money * 2
    years = years + 1
}
print(years)
print(money)
```

```output
7
128
```

## Loops inside loops

```gan
for row in 1 to 3 {
    let line = ""
    for col in 1 to 4 {
        line = line + row * col + " "
    }
    print(line)
}
```

```output
1 2 3 4 
2 4 6 8 
3 6 9 12 
```

That prints a small multiplication table.

## No break (yet)

Some languages have `break` to leave a loop early. GAN does not yet. Inside a function you can leave a
loop at any time with `return` (Lesson 7). Otherwise use a `while` with a good question.

## Exercises

1. Print the numbers 1 to 10.
2. Print the even numbers from 2 to 20.
3. Add up the numbers 1 to 10 and print the total.
4. Print a countdown from 5 to 1 using `while`.
5. Print the times table for 7, from `7 * 1` to `7 * 10`.

Next: [Lesson 7: Functions](07-functions.md)
