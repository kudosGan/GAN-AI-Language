# GAN agents: a design

**Status: the three decisions at the end were accepted. `tool` works (Step 35). `agent`, `.do()` and
`log` do not work yet.**

## What an agent is

A model that can **use tools**. It runs in a loop:

1. GAN sends the task and the list of tools to the model.
2. The model either **asks for a tool** (with arguments) or **gives its final answer**.
3. If it asks for a tool, GAN runs that GAN function and sends the result back.
4. Go to step 2, until there is an answer or the agent reaches its step limit.

## How it would look

```
// A tool is a function the model may use. The sentence tells the model what it is for.
tool add(a: number, b: number): number "Adds two numbers" {
    return a + b
}

tool weather(city: text): text "Today's weather in a city" {
    return "Sunny in " + city          // a real tool could use python(...) here
}

agent helper {
    model: "qwen2.5:7b"
    instructions: "Answer briefly. Use the tools for sums."
    tools: [add, weather]
    steps: 10
}

let answer = helper.do("What is 12 plus 30?")
print(answer)
print(helper.log)                      // everything that happened, step by step
```

## The pieces

| Piece | Meaning |
|---|---|
| `tool` | Like `func`, plus a sentence describing it. **Every parameter needs a type**, so the model can be told what to send. A normal program can still call it like a function. |
| `agent` | Gives a name to a model, its instructions, its tools and a step limit. |
| `helper.do(task)` | Runs the loop and returns the final answer as text. |
| `helper.log` | A list of text lines: what the model asked for and what each tool returned. |

Only `model` is required. `instructions` is empty, `tools` is `[]` and `steps` is 10 if left out.

## Full power, chosen by you

GAN adds no content filter and does not limit **what a tool can do**. A tool may call `python`, change
files or run commands, if you write it that way.

The rule that matters: **a model can only run the tools you list.** It cannot reach `python` or any
other function by naming it.

Why this rule is part of the design: a model reads text from web pages, files and users, and that text
can contain instructions meant to trick it. The tool list is what stops such text from doing anything
you did not hand to the agent. A very powerful tool, such as one that runs any command, gives that power
to whatever text reaches the model. GAN will let you write it; the choice and the risk are yours.

- **Step limit:** it is there so that a confused model cannot loop forever. You can set any number.
- **Asking first:** if you want to approve a tool call, write it in the tool with `input(...)`.
  GAN does not force it.

## When things go wrong

| Problem | What GAN would do |
|---|---|
| The model asks for a tool that is not listed | Tell the model it does not exist. Nothing runs. |
| Wrong or missing arguments | Tell the model what was wrong, so it can try again. The tool does not run. |
| A tool stops with an error | Send the error message back to the model, and record it in the log. |
| Step limit reached | Stop with a clear error that names the agent and the line. |
| Ollama fails | The same clear errors as `ask`. |
| A mistake in your agent code | Reported with its line number, **before** the program runs where possible. |

Note: a tool that failed half-way may already have changed something (for example, written half a file).
The log shows what ran.

## How we would build it (one step at a time)

| Step | What it adds | Can it be tested without a real model? |
|---|---|---|
| 35 | `tool`: the new word, descriptions, required types | Yes |
| 36 | `agent` and `do`: the loop, the step limit, tool errors | Yes, with a fake Ollama that asks for tools |
| 37 | `log`, and a lesson building a useful agent | Yes, plus a run with your real model |
| later | Memory between tasks, several agents working together, dictionaries and JSON, other model servers | |

**Dictionaries are not needed first.** GAN itself (in Go) reads the model's tool requests and turns them
into ordinary arguments such as numbers and text.

Ollama has a "chat with tools" interface, and the `qwen2.5` models you have are listed by Ollama as
supporting tools. A 7B model will sometimes pick the wrong tool or wrong arguments, so Step 36 must be
tried with the real model before we say how well it works.

## Decisions (accepted: the recommended choices)

1. **The words.** `tool`, `agent` and `.do(...)`, or others you like better?
2. **The description.** After the parameters, as shown above, or as a comment line above the tool?
   Recommended: as shown, because GAN can then require it.
3. **Default step limit.** 10 is recommended. You could also choose "no limit unless set".
