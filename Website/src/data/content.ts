// All the words and lists shown on the site. WRITTEN BY the content-writer agent:
// facts come from the real GAN files, wording from the AI. Edit by hand if you like.

export const REPO_URL = "https://github.com/kudosGan/GAN-AI-Language";
export const VERSION = "0.3 beta";

export const SAMPLE_CODE = `// Ask a name and greet it
let name = input("Your name? ")
print("Hello, " + name)

for i in 1 to 3 {
    print("Count: " + str(i))
}`;

export const FEATURES = [
  {
    "title": "Easy to Read",
    "text": "GAN code looks like plain English."
  },
  {
    "title": "Friendly Error Messages",
    "text": "Get clear, helpful hints when things go wrong."
  },
  {
    "title": "Useful Things Built In",
    "text": "Lists, math, and text tools make coding simpler."
  }
];

export const LESSONS = [
  "Welcome to GAN",
  "Printing and comments",
  "Variables",
  "Numbers and math",
  "Text",
  "Making decisions",
  "Loops",
  "Functions",
  "Lists",
  "Asking the user",
  "Libraries",
  "Types",
  "When things go wrong",
  "Projects",
  "Tensors",
  "How a model learns",
  "Talking to a model"
];

export const LIBRARIES = [
  {
    "name": "lists",
    "summary": "Handle collections of items.",
    "functions": [
      "range",
      "reverse",
      "sorted",
      "first",
      "last"
    ]
  },
  {
    "name": "math",
    "summary": "Perform common mathematical operations.",
    "functions": [
      "PI",
      "E",
      "square",
      "cube",
      "sign",
      "clamp",
      "isEven",
      "isOdd",
      "factorial",
      "gcd",
      "isPrime",
      "sumOf",
      "average",
      "maxOf",
      "minOf"
    ]
  },
  {
    "name": "text",
    "summary": "Manipulate and analyze strings.",
    "functions": [
      "repeat",
      "join",
      "reverse",
      "isEmpty",
      "padLeft",
      "isPalindrome"
    ]
  }
];

export const MENU = [
  { label: "Home", page: "home" },
  { label: "Download", page: "download" },
  { label: "Learn", page: "learn" },
  { label: "Libraries", page: "libraries" },
  { label: "About", page: "about" },
];
