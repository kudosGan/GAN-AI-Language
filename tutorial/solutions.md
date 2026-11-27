# Solutions to the exercises

Try each exercise yourself first. There are many correct answers: yours does not have to match these exactly,
only give the same result.

## Lesson 1: Printing and comments

**1 and 2.** Print your name, then three lines:

```gan
print("Ann")
print("Hello!")
print("pasta")
print(7)
```

```output
Ann
Hello!
pasta
7
```

**3.** `7 * 6`:

```gan
print(7 * 6)
```

```output
42
```

**4.** A comment on every line, one `print` switched off:

```gan
// greet the user
print("Hi")
// print("This one is switched off")
print("Bye") // say goodbye
```

```output
Hi
Bye
```

## Lesson 2: Variables

**1 and 2.**

```gan
let name = "Ann"
let age = 12
print(name)
print(age)

let price = 4
let count = 3
print(price * count)
```

```output
Ann
12
12
```

**3.** A `final` cannot change:

```gan
final taxRate = 2
taxRate = 3
```

```output
Runtime error on line 2: cannot change final variable "taxRate"
    2 | taxRate = 3
```

**4.** `a` is 5, then `a * 2` is 10:

```gan
let a = 5
a = a * 2
print(a)
```

```output
10
```

## Lesson 3: Math

**1.** Pizza for 4 friends:

```gan
print(12 / 4)
```

```output
3
```

**2.** Average of three scores:

```gan
print((80 + 92 + 71) / 3)
```

```output
81
```

**3.** Area of a circle:

```gan
let r = 3
print(3.14159 * r * r)
```

```output
28.274309999999996
```

**4.** Even check and square root:

```gan
print(144 % 2)
print(sqrt(144))
```

```output
0
12
```

**5.** Multiplication goes first, so the first is `2 + 12 - 1`. In the second, the parentheses are worked out first, `5 * 3`:

```gan
print(2 + 3 * 4 - 1)
print((2 + 3) * (4 - 1))
```

```output
13
15
```

## Lesson 4: Text

**1 and 2.**

```gan
let name = "Ann"
let age = 12
print("Hello, " + name)
print(name + " is " + age + " years old")
```

```output
Hello, Ann
Ann is 12 years old
```

**3.** First and last letters:

```gan
let word = "planet"
print(word[0])
print(word[len(word) - 1])
```

```output
p
t
```

**4.** Upper case and length:

```gan
let word = "planet"
print(upper(word) + " " + len(word))
```

```output
PLANET 6
```

**5.** Each letter of a name:

```gan
for letter in "Ann" {
    print(letter)
}
```

```output
A
n
n
```

## Lesson 5: Decisions

**1.** Positive, negative or zero:

```gan
let number = -3
if number > 0 {
    print("positive")
} else if number < 0 {
    print("negative")
} else {
    print("zero")
}
```

```output
negative
```

**2.** Even or odd:

```gan
let n = 7
if n % 2 == 0 {
    print("even")
} else {
    print("odd")
}
```

```output
odd
```

**3.** Ticket price:

```gan
let age = 70
if age < 13 {
    print(5)
} else if age < 65 {
    print(10)
} else {
    print(7)
}
```

```output
7
```

**4.** Between 1 and 10:

```gan
let x = 7
print(x >= 1 and x <= 10)
```

```output
true
```

## Lesson 6: Loops

**1.**

```gan
for i in 1 to 10 {
    print(i)
}
```

```output
1
2
3
4
5
6
7
8
9
10
```

**2.** Even numbers up to 20 (loop over 1 to 10 and double):

```gan
for i in 1 to 10 {
    print(i * 2)
}
```

```output
2
4
6
8
10
12
14
16
18
20
```

**3.** Total of 1 to 10:

```gan
let total = 0
for n in 1 to 10 {
    total = total + n
}
print(total)
```

```output
55
```

**4.** Countdown with `while`:

```gan
let n = 5
while n >= 1 {
    print(n)
    n = n - 1
}
```

```output
5
4
3
2
1
```

**5.** Times table for 7:

```gan
for i in 1 to 10 {
    print("7 x " + i + " = " + 7 * i)
}
```

```output
7 x 1 = 7
7 x 2 = 14
7 x 3 = 21
7 x 4 = 28
7 x 5 = 35
7 x 6 = 42
7 x 7 = 49
7 x 8 = 56
7 x 9 = 63
7 x 10 = 70
```

## Lesson 7: Functions

**1 and 2.**

```gan
func square(x) {
    return x * x
}

func isEven(n) {
    return n % 2 == 0
}
print(square(6))
print(isEven(10))
print(isEven(7))
```

```output
36
true
false
```

**3.** Average of three:

```gan
func average(a, b, c) {
    return (a + b + c) / 3
}
print(average(80, 92, 71))
```

```output
81
```

**4.** Countdown:

```gan
func countdown(n) {
    while n >= 1 {
        print(n)
        n = n - 1
    }
}
countdown(3)
```

```output
3
2
1
```

**5.** Power with a loop:

```gan
func power(base, exp) {
    let result = 1
    for i in 1 to exp {
        result = result * base
    }
    return result
}
print(power(2, 10))
print(power(5, 0))
```

```output
1024
1
```

**6.** Fibonacci with recursion:

```gan
func fib(n) {
    if n <= 2 {
        return 1
    }
    return fib(n - 1) + fib(n - 2)
}
for i in 1 to 8 {
    print(fib(i))
}
```

```output
1
1
2
3
5
8
13
21
```

## Lesson 8: Lists

**1 and 2.**

```gan
let nums = [4, 8, 15, 16, 23]
print(nums[0])
print(nums[len(nums) - 1])
for n in nums {
    print(n * 2)
}
```

```output
4
23
8
16
30
32
46
```

**3.** Count the numbers bigger than 10:

```gan
let nums = [4, 12, 15, 3, 23]
let count = 0
for n in nums {
    if n > 10 {
        count = count + 1
    }
}
print(count)
```

```output
3
```

**4.** `reverse`:

```gan
func reverse(list) {
    let result = []
    for item in list {
        result = [item] + result
    }
    return result
}
print(reverse([1, 2, 3]))
```

```output
[3, 2, 1]
```

**5.** Favourite foods:

```gan
let foods = ["pizza", "rice", "mango"]
push(foods, "soup")
for food in foods {
    print(food)
}
```

```output
pizza
rice
mango
soup
```

## Lesson 9: Asking the user

**1.**

```gan
let name = input("Name? ")
let colour = input("Favourite colour? ")
print(name + " likes " + colour)
```

```input
Ann
green
```
```output
Name? Ann
Favourite colour? green
Ann likes green
```

**2.** Two numbers:

```gan
let a = num(input("First: "))
let b = num(input("Second: "))
print("Sum: " + (a + b))
print("Difference: " + (a - b))
print("Product: " + a * b)
```

```input
8
3
```
```output
First: 8
Second: 3
Sum: 11
Difference: 5
Product: 24
```

**3.** Celsius to Fahrenheit:

```gan
let c = num(input("Celsius? "))
print(c * 9 / 5 + 32)
```

```input
100
```
```output
Celsius? 100
212
```

**4.** Palindrome check:

```gan
let word = input("A word: ")
let backwards = ""
for letter in word {
    backwards = letter + backwards
}
print(word == backwards)
```

```input
level
```
```output
A word: level
true
```

## Lesson 10: Libraries

**1.**

```gan
import "math"
print(math.factorial(6))
print(math.gcd(84, 36))
```

```output
720
12
```

**2.** Squares of 1 to 5:

```gan
import "lists"
for n in lists.range(1, 5) {
    print(n * n)
}
```

```output
1
4
9
16
25
```

**3.** Smallest and largest:

```gan
import "lists"
let sorted = lists.sorted([5, 3, 8, 1])
print(lists.first(sorted))
print(lists.last(sorted))
```

```output
1
8
```

**4.** Your own library. Save this as `greetings.gan`:

```gan
func hello(name) {
    return "Hello, " + name + "!"
}
```

and this next to it, as `main.gan`:

```gan
import "greetings.gan"
print(greetings.hello("Ann"))
```

## Lesson 11: Types

**1.**

```gan
let price: number = 5
price = "free"
```

```output
Runtime error on line 2: variable "price" must be number but got text
    2 | price = "free"
```
```

**2.**

```gan
func shout(word: text): text {
    return upper(word) + "!"
}
print(shout("hey"))
```

```output
HEY!
```

**3.**

```gan
func average(list: list): number {
    let total = 0
    for x in list {
        total = total + x
    }
    return total / len(list)
}
print(average([2, 4, 9]))
```

```output
5
```

**4.**

```gan
print(typeof(1))
print(typeof("a"))
print(typeof(false))
print(typeof([]))
print(typeof(1 + 1 == 2))
```

```output
number
text
bool
list
bool
```

## Lesson 12: When things go wrong

**1.** `print(Hello)` means "print the variable Hello". Add quotes:

```gan
print("Hello")
```

```output
Hello
```

**2.** `x == 6` asks a question and throws the answer away. To change `x`, use one `=`:

```gan
let x = 5
x = 6
print(x)
```

```output
6
```

**3.** The closing `}` is missing:

```gan
for i in 1 to 3 {
    print(i)
}
```

```output
1
2
3
```

**4.** The list has 3 items, numbered 0, 1 and 2. The last one is `nums[2]`:

```gan
let nums = [1, 2, 3]
print(nums[2])
```

```output
3
```

**5.** `f` wants two arguments:

```gan
func f(a, b) {
    return a * b
}
print(f(2, 5))
```

```output
10
```

## Lesson 14: Tensors

**1.**

```gan
let t = tensor([1, 2, 3, 4])
print(t * 2)
```

```output
tensor([2, 4, 6, 8])
```

**2.**

```gan
print(sum(tensor([[5, 5], [5, 5]])))
```

```output
20
```

**3.** One row times one column gives one number: 1 + 2 + 3.

```gan
print(matmul(tensor([[1, 2, 3]]), tensor([[1], [1], [1]])))
```

```output
tensor([[6]])
```

**4.** The differences are 0, -1 and 2. Squared they are 0, 1 and 4, and their average is 5 / 3.

```gan
let guess = tensor([2, 4, 9])
let right = tensor([2, 5, 7])
let diff = guess - right
print(mean(diff * diff))
```

```output
1.6666666666666667
```

## Lesson 15: How a model learns

**1.** Each gradient is twice its number:

```gan
let w = learnable(tensor([1, 5]))
backward(sum(w * w))
print(grad(w))
```

```output
tensor([2, 10])
```

**2.** Yes. The loss is checked after the tenth step:

```gan
let x = tensor([[1], [2]])
let target = tensor([[2], [4]])
var w = learnable(tensor([[0]]))

for round in 1 to 10 {
    let diff = matmul(x, w) - target
    backward(sum(diff * diff))
    w = learnable(w - 0.05 * grad(w))
}
let diff = matmul(x, w) - target
print(item(sum(diff * diff)) < 0.001)
```

```output
true
```

**3.** The step is too big. The weight jumps past 2, further each time, and the loss **grows**:

```gan
let x = tensor([[1], [2]])
let target = tensor([[2], [4]])
var w = learnable(tensor([[0]]))

for round in 1 to 3 {
    let diff = matmul(x, w) - target
    let loss = sum(diff * diff)
    backward(loss)
    print(item(loss))
    w = learnable(w - 0.25 * grad(w))
}
```

```output
20
45
101.25
```
