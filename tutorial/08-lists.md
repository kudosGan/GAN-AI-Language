# Lesson 8: Lists

A **list** holds many values in order. Write it with square brackets:

```gan
let scores = [90, 72, 85]
print(scores)
```

```output
[90, 72, 85]
```

A list can hold anything, even a mix. Text inside a list is shown with quotes, so you can tell it
from a number:

```gan
let mixed = [1, "two", true, [3, 4]]
print(mixed)
print([])
```

```output
[1, "two", true, [3, 4]]
[]
```

## Reading items

Like text, positions start at **0**:

```gan
let names = ["Ann", "Bo", "Cy"]
print(names[0])
print(names[2])
print(len(names))
print(names[len(names) - 1])
```

```output
Ann
Cy
3
Cy
```

An item that is not there is an error that tells you how long the list is:

```gan
let names = ["Ann", "Bo", "Cy"]
print(names[3])
```

```output
Runtime error on line 2: index 3 is out of range (the list has 3 items)
    2 | print(names[3])
```

## Changing items

```gan
let nums = [10, 20, 30]
nums[1] = 99
print(nums)
```

```output
[10, 99, 30]
```

## Adding and removing

`push(list, item)` adds to the end. `pop(list)` removes the last item and gives it back:

```gan
let items = []
push(items, 5)
push(items, 7)
push(items, 9)
print(items)
print(pop(items))
print(items)
print(len(items))
```

```output
[5, 7, 9]
9
[5, 7]
2
```

## Going through a list

```gan
let fruits = ["apple", "banana", "cherry"]
for fruit in fruits {
    print("I like " + fruit)
}
```

```output
I like apple
I like banana
I like cherry
```

Adding up a list:

```gan
let prices = [4, 7, 2, 10]
let total = 0
for p in prices {
    total = total + p
}
print(total)
print(total / len(prices))
```

```output
23
5.75
```

If you need the position as well, loop over the positions:

```gan
let letters = ["a", "b", "c"]
for i in 0 to len(letters) - 1 {
    print(i + ": " + letters[i])
}
```

```output
0: a
1: b
2: c
```

## Searching, joining, comparing

```gan
let nums = [1, 2, 3]
print(contains(nums, 2))
print(contains(nums, 9))
print([1, 2] + [3])
print([1, 2] == [1, 2])
```

```output
true
false
[1, 2, 3]
true
```

`+` puts two lists together into a new list. `==` compares item by item.

## Lists inside lists

```gan
let grid = [[1, 2], [3, 4]]
print(grid[1][0])
grid[0][1] = 20
print(grid)
```

```output
3
[[1, 20], [3, 4]]
```

## Two names, one list

This surprises many people. A list is **shared**, not copied:

```gan
let a = [1, 2]
let b = a
b[0] = 100
print(a)
```

```output
[100, 2]
```

`b` is not a copy. It is another name for the same list. To make a real copy, add an empty list:

```gan
let a = [1, 2]
let b = a + []
b[0] = 100
print(a)
print(b)
```

```output
[1, 2]
[100, 2]
```

## Functions and lists

```gan
func biggest(list) {
    let best = list[0]
    for x in list {
        if x > best {
            best = x
        }
    }
    return best
}
print(biggest([3, 9, 4]))
```

```output
9
```

A function that gets a list can change it, because the list is shared:

```gan
func addTwo(list) {
    push(list, 2)
}
let data = [1]
addTwo(data)
print(data)
```

```output
[1, 2]
```

## Exercises

1. Make a list of five numbers and print the first and last.
2. Print every number in the list doubled.
3. Count how many numbers in a list are bigger than 10.
4. Write `reverse(list)` that returns a new list with the items in the opposite order.
5. Make a list of your three favourite foods, add a fourth with `push`, and print each on its own line.

Next: [Lesson 9: Asking the user](09-input.md)
