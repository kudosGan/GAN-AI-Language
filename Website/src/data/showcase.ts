// "What's new" on the home page: samples of the latest GAN and what they print.
// The output text must be exactly what the real gan.exe prints. Check it with:
//     python agent\web\check_showcase.py
// This file is separate from content.ts so the content agent never overwrites it.

export const SHOWCASE_INTRO =
  "GAN's goal is to become a language designed for AI. These are its first AI building blocks. " +
  "Every sample below is a complete program.";

// A sample with "varies: true" asks an AI model, so its answer differs from run to run. The
// check script does not compare its output; the page labels it as an example answer.
export const SHOWCASE = [
  {
    title: "Tools: the first piece of AI agents",
    text: "A tool is a function with a description and typed inputs, so that an AI model can later be told how to use it. Agents are not built yet: today your own program calls the tool.",
    file: "tools.gan",
    code: `tool add(a: number, b: number): number "Adds two numbers" {
    return a + b
}

print(add(12, 30))`,
    output: `42`,
  },
  {
    title: "Use Python from GAN",
    text: "python() calls functions from your installed Python, including libraries such as NumPy. It has Python's full power, including deleting files and running commands, so read other people's programs before you run them.",
    file: "python.gan",
    code: `print(python("math.factorial", 5))
print(python("string.capwords", "hello gan"))
print(python("builtins.sorted", [3, 1, 2]))`,
    output: `120
Hello Gan
[1, 2, 3]`,
  },
  {
    title: "Quotes and new lines in text",
    text: "Text escapes make it easy to write quotes, new lines, tabs and JSON.",
    file: "escapes.gan",
    code: `print("She said \\"hello\\"")
print("{\\"tool\\": \\"add\\"}")`,
    output: `She said "hello"
{"tool": "add"}`,
  },
  {
    title: "Talk to an AI model",
    text: "ask() sends a question to a model running on your own computer (with the free program Ollama) and gives back the answer as text. This is a first step towards AI agents.",
    file: "ask.gan",
    code: `let topic = "loops"
let answer = ask("Explain " + topic + " in programming in one short sentence.")
print(answer)`,
    output: `Loops in programming are structures that allow a set of instructions to be executed repeatedly until a certain condition is met.`,
    varies: true,
  },
  {
    title: "A model that learns",
    text: "GAN calculates the gradients itself. This one-weight model learns to double a number: the weight moves towards 2.",
    file: "learn.gan",
    code: `let x = tensor([[1], [2]])
let target = tensor([[2], [4]])   // x doubled
var w = learnable(tensor([[0]]))  // first guess

for round in 1 to 3 {
    let diff = matmul(x, w) - target
    let loss = sum(diff * diff)
    backward(loss)
    w = learnable(w - 0.05 * grad(w))
    print("loss " + item(loss) + ", weight " + w)
}`,
    output: `loss 20, weight tensor([[1]])
loss 5, weight tensor([[1.5]])
loss 1.25, weight tensor([[1.75]])`,
  },
  {
    title: "Tensors are built in",
    text: "Grids of numbers and neural network maths are part of the language. There is nothing to install or import.",
    file: "tensors.gan",
    code: `let a = tensor([[1, 2], [3, 4]])
let b = tensor([[5, 6], [7, 8]])

print(a * 10)
print(matmul(a, b))
print(relu(tensor([-2, 0, 3])))`,
    output: `tensor([[10, 20], [30, 40]])
tensor([[19, 22], [43, 50]])
tensor([0, 0, 3])`,
  },
  {
    title: "Errors that explain themselves",
    text: "Shape mistakes are common in AI code. GAN shows the line, both shapes and what does not fit.",
    file: "shapes.gan",
    code: `let x = tensor([[1, 2, 3]])
let w = tensor([[4, 5, 6]])
print(matmul(x, w))`,
    output: `Runtime error on line 3: matmul shapes do not fit: [1, 3] and [1, 3] (the first has 3 columns, so the second needs 3 rows, but it has 1)
    3 | print(matmul(x, w))`,
  },
];

// Said plainly, so nobody expects more than the beta can do.
export const BETA_NOTE =
  "This is an early beta. GAN cannot run AI agents yet: tools can be written, but a model cannot use them so far. " +
  "ask() needs Ollama and a downloaded model, and python() needs Python. Tensors run on the CPU, gradients work through + - * / matmul " +
  "transpose sum mean relu sigmoid tanh exp log (not [ ]), and GAN is not ready for large AI models.";
