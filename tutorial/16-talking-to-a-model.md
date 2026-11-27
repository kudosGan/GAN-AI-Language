# Lesson 16: Talking to a model

GAN can send a question to an AI model and get the answer back as text. This is the first step towards
**agents**, programs that use a model to decide what to do.

## What you need

The model runs on your own computer, using a free program called **Ollama**. You need Ollama running,
and a model downloaded. See which models you have with this command in a terminal:

```
ollama list
```

If you do not have the default model, download it once (it is several gigabytes):

```
ollama pull qwen2.5:7b
```

## ask

The examples in this lesson are not checked automatically, because **a model's answer can be
different every time**. Your answers will not be word for word the same as anyone else's.

```gan
let answer = ask("Explain what a loop is in one sentence.")
print(answer)
```

The first question can take a while, because Ollama has to load the model first. GAN waits up to
5 minutes.

The answer is ordinary **text**, so everything you know about text works on it: `+`, `len`,
`contains`, `lower`.

To use another model, give its name second:

```gan
print(ask("Say hello in French.", "qwen2.5:14b"))
```

## Building a question

Join text to put your own values into a question:

```gan
let topic = input("Topic? ")
print(ask("Give one fun fact about " + topic + ". Use one sentence."))
```

Say clearly what shape of answer you want, such as "Reply with only the number" or "Reply with only
yes or no". Even then a model may add extra words, so check the answer before you rely on it:

```gan
let reply = ask("Is 17 a prime number? Reply with only yes or no.")
if contains(lower(reply), "yes") {
    print("The model says yes")
} else {
    print("The model did not say yes. It said: " + reply)
}
```

## The model does not remember

Each `ask` starts fresh. To continue a conversation, put the earlier words into the new question:

```gan
var history = ""
for turn in 1 to 3 {
    let you = input("You: ")
    history = history + "User: " + you + "
"
    let reply = ask("Continue this conversation. Reply with only your next message.
" + history)
    print("Model: " + reply)
    history = history + "You: " + reply + "
"
}
```

GAN text has no `\n` yet. A new line inside the quotes, as above, is how to put one in text.

## Important

- **Models can be wrong**, and can sound sure when they are. Check anything that matters.
- **The answer is only text.** GAN never runs what a model writes.
- **Privacy:** GAN sends your question only to Ollama on this computer. With a downloaded model it
  stays there. A "cloud" model set up in Ollama would send it over the internet.
- A model's answer can be influenced by the text you put in the question, including text typed by
  someone else.

## When it goes wrong

A mistake in your program is shown with its line number:

```gan
let question = 42
print(ask(question))
```

```output
Runtime error on line 2: ask needs text as the question, but got number
    2 | print(ask(question))
```

If GAN says it cannot reach Ollama, start Ollama and try again. If it says Ollama does not have the
model, check the name with `ollama list`. If the message says **Ollama reported an error**, the
problem happened inside Ollama (for example while using the graphics card). Try again; if it keeps
happening, restart Ollama.

## Exercises

1. Ask the user for a word, and ask the model for a sentence that uses it. Check with `contains` that
   the word really is in the answer.
2. Write a function `translate(text, language)` that returns the model's translation.

Back to the [course index](README.md)
