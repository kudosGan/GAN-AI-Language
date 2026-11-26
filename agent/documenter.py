"""
documenter.py - the GAN documenter agent.

It writes tutorial documents about the GAN language into agent/document/.

How it works (two halves):
  1. FACTS: Python reads your real project (the Go source, the .gan example files,
     and agent/steps.json) and runs your examples with gan.exe to get real output.
     Facts are never invented.
  2. PROSE: a local AI model (Ollama, on your own PC) writes the beginner-friendly
     explanations around those facts. If Ollama is not available, the documents
     are still written, just without the explanations.

Usage (from C:\\Desktop\\GAN):
    python agent\\documenter.py                       use the default model
    python agent\\documenter.py --model qwen2.5:14b   use another installed model
    python agent\\documenter.py --no-ai               facts only, no AI
    python agent\\documenter.py --watch               keep running; update when the project changes

Uses only the Python standard library. Nothing to install.
"""

import argparse
import datetime
import hashlib
import json
import re
import subprocess
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path

# ---------------------------------------------------------------- settings
AGENT_DIR = Path(__file__).resolve().parent
ROOT = AGENT_DIR.parent                      # C:\Desktop\GAN
OUT_DIR = AGENT_DIR / "document"             # where the documents are written
STEPS_FILE = AGENT_DIR / "steps.json"        # the human notes, one entry per step
CACHE_FILE = AGENT_DIR / ".ai-cache.json"    # remembers AI answers so reruns are fast
GAN_EXE = ROOT / "gan.exe"
GO_FILES = ["main.go", "lexer.go", "parser.go", "interpreter.go"]

DEFAULT_MODEL = "qwen2.5:7b"
OLLAMA_URL = "http://localhost:11434"


# ---------------------------------------------------------------- reading facts
def read_text(path):
    return Path(path).read_text(encoding="utf-8")


def run_gan(gan_file):
    """Run a .gan file with gan.exe and return what it printed.
    Programs that ask for input() are given the lines saved in tests/expected/<name>.in
    (or nothing), so the documenter never sits waiting for someone to type."""
    stdin_file = ROOT / "tests" / "expected" / (Path(gan_file).stem + ".in")
    typed = stdin_file.read_text(encoding="utf-8") if stdin_file.exists() else ""
    try:
        result = subprocess.run(
            [str(GAN_EXE), str(gan_file)],
            input=typed, capture_output=True, text=True, timeout=10, encoding="utf-8",
            cwd=str(ROOT),
        )
        return (result.stdout + result.stderr).strip()
    except subprocess.TimeoutExpired:
        return "(program ran too long and was stopped)"
    except OSError as e:
        return f"(could not run gan.exe: {e})"


def doc_comment_before(lines, index):
    """Return the // comment lines directly above lines[index], joined into one string."""
    found = []
    i = index - 1
    while i >= 0 and lines[i].strip().startswith("//"):
        found.insert(0, lines[i].strip().lstrip("/").strip())
        i -= 1
    return " ".join(found)


def strip_name(name, note):
    """'IfStatement: if x {...}' -> 'if x {...}' so the name is not said twice."""
    if note.startswith(name):
        note = note[len(name):].lstrip(": ").strip()
    return note


def read_source_facts():
    """Collect tokens, statements, expressions, parser functions and error messages."""
    facts = {"tokens": [], "statements": [], "expressions": [], "parser_functions": [], "errors": []}

    # tokens: lines like  PLUS TokenType = "PLUS"  // comment
    lexer = read_text(ROOT / "lexer.go")
    for name, value, comment in re.findall(r'^\s*(\w+)\s+TokenType\s*=\s*"(\w+)"[ \t]*(?://[ \t]*(.*))?$', lexer, re.M):
        facts["tokens"].append((name, comment.strip()))

    # statements and expressions: "type X struct" in parser.go
    parser_lines = read_text(ROOT / "parser.go").splitlines()
    for i, line in enumerate(parser_lines):
        m = re.match(r"type (\w+) struct", line)
        if not m or m.group(1) == "Parser":
            continue
        name = m.group(1)
        note = doc_comment_before(parser_lines, i)
        inline = re.search(r"//\s*(.*)$", line)
        if inline and not note:
            note = inline.group(1).strip()
        target = "statements" if name.endswith("Statement") else "expressions"
        facts[target].append((name, strip_name(name, note)))

    # parser functions, with the comment written above each one
    for i, line in enumerate(parser_lines):
        m = re.match(r"func \(p \*Parser\) (parse\w+)", line)
        if m:
            facts["parser_functions"].append((m.group(1), strip_name(m.group(1), doc_comment_before(parser_lines, i))))

    # error messages: fmt.Errorf("...") in every Go file
    seen = set()
    for go in GO_FILES:
        path = ROOT / go
        if not path.exists():
            continue
        for msg in re.findall(r'fmt\.Errorf\("((?:[^"\\]|\\.)*)"', read_text(path)):
            if msg not in seen:
                seen.add(msg)
                facts["errors"].append(msg)
    return facts


# ---------------------------------------------------------------- the AI (Ollama)
class Writer:
    """Asks the local model for explanations. Answers are cached on disk."""

    def __init__(self, model, enabled):
        self.model = model
        self.enabled = enabled
        self.calls = 0
        self.cache = {}
        if CACHE_FILE.exists():
            try:
                self.cache = json.loads(read_text(CACHE_FILE))
            except (ValueError, OSError):
                self.cache = {}
        if self.enabled:
            self.enabled = self.check_ollama()

    def check_ollama(self):
        """Make sure Ollama is running and has the model. Otherwise fall back to facts only."""
        try:
            with urllib.request.urlopen(OLLAMA_URL + "/api/tags", timeout=5) as r:
                names = [m["name"] for m in json.load(r).get("models", [])]
        except (urllib.error.URLError, OSError, ValueError):
            print("! Ollama is not running, so I will write the documents without AI explanations.")
            print("  (Start Ollama and run me again to add them.)")
            return False
        if self.model not in names and (self.model + ":latest") not in names:
            print(f"! Model {self.model!r} is not installed. Installed: {', '.join(names) or 'none'}")
            print("  Writing the documents without AI explanations.")
            return False
        return True

    def ask(self, topic, facts):
        """Return a short plain-English explanation, or '' when AI is off or fails."""
        if not self.enabled:
            return ""
        prompt = (
            "You are writing a tutorial for a complete beginner who is learning programming "
            "by watching a small language called GAN being built in Go.\n"
            "Write 2 or 3 short paragraphs of plain, friendly English about the topic below.\n"
            "Rules: use ONLY the facts provided. Do not invent syntax, features, commands or code. "
            "Do not write code blocks, headings or bullet lists. Explain any technical word the "
            "first time you use it. GAN has no semicolons and no parentheses in math, so never "
            "add them when you mention code; put each GAN statement on its own line of thought, "
            "and copy code snippets exactly as given in the facts.\n\n"
            f"TOPIC: {topic}\n\nFACTS:\n{facts}\n"
        )
        key = hashlib.sha256((self.model + prompt).encode("utf-8")).hexdigest()
        if key in self.cache:
            return self.cache[key]

        self.calls += 1
        print(f"  asking {self.model} about: {topic} ...", flush=True)
        body = json.dumps({
            "model": self.model, "prompt": prompt, "stream": False,
            "options": {"temperature": 0.2, "num_predict": 500},
        }).encode("utf-8")
        req = urllib.request.Request(OLLAMA_URL + "/api/generate", data=body,
                                     headers={"Content-Type": "application/json"})
        text = None
        for attempt in (1, 2):  # the first request can fail while the model is still loading
            try:
                with urllib.request.urlopen(req, timeout=600) as r:
                    text = json.load(r).get("response", "")
                break
            except (urllib.error.URLError, OSError, ValueError) as e:
                if attempt == 1:
                    print(f"  ({e}); trying once more in 5 seconds")
                    time.sleep(5)
                else:
                    print(f"  ! the model did not answer ({e}); leaving this explanation out")
                    return ""

        text = clean_ai_text(text)
        if text:
            self.cache[key] = text
            CACHE_FILE.write_text(json.dumps(self.cache, indent=1), encoding="utf-8")
        return text


def clean_ai_text(text):
    """Remove headings and code blocks if the model added them anyway."""
    text = re.sub(r"```.*?```", "", text, flags=re.S)
    lines = [ln for ln in text.splitlines() if not ln.lstrip().startswith("#")]
    return re.sub(r"\n{3,}", "\n\n", "\n".join(lines)).strip()


# ---------------------------------------------------------------- writing documents
def header(title, today):
    return (
        f"# {title}\n\n"
        f"_Generated by `agent/documenter.py` on {today}. Do not edit by hand: "
        f"change `agent/steps.json` or the GAN project and run the documenter again._\n\n"
    )


def footer(writer):
    if writer.enabled:
        return (f"\n---\n_Code, program output and lists are read from the real project. "
                f"Explanations were written by the local model `{writer.model}`._\n")
    return "\n---\n_Code, program output and lists are read from the real project. Written without AI explanations._\n"


def status_label(status):
    return {"done": "Done", "next": "Next", "planned": "Planned"}.get(status, status)


def build_progress(steps, writer, today):
    doc = header("GAN progress", today)
    done = [s for s in steps if s["status"] == "done"]
    doc += f"{len(done)} of {len(steps)} roadmap items are done.\n\n"

    facts = "\n".join(f"- Step {s['step']} ({status_label(s['status'])}): {s['title']}. {s['summary']}" for s in steps)
    intro = writer.ask("A short 'where we are now and what comes next' summary of the project", facts)
    if intro:
        doc += intro + "\n\n"

    doc += "| Step | What | Status | Date |\n|---|---|---|---|\n"
    for s in steps:
        doc += f"| {s['step']} | {s['title']} | {status_label(s['status'])} | {s['date'] or '-'} |\n"
    doc += "\n"

    for s in steps:
        doc += f"## Step {s['step']}: {s['title']}\n\n"
        doc += f"**Status:** {status_label(s['status'])}" + (f" ({s['date']})" if s["date"] else "") + "\n\n"
        doc += s["summary"] + "\n\n"
        if s["learned"]:
            doc += "**What we learned**\n\n" + "".join(f"- {x}\n" for x in s["learned"]) + "\n"
    return doc + footer(writer)


def gan_block(filename):
    """A .gan example as code, plus what it really prints."""
    path = ROOT / filename
    if not path.exists():
        return f"_(example `{filename}` was not found)_\n\n"
    code = read_text(path).rstrip()
    output = run_gan(path)
    return (f"`{filename}`:\n\n```gan\n{code}\n```\n\nOutput:\n\n```\n{output}\n```\n\n")


def build_guide(steps, writer, today):
    doc = header("GAN language guide", today)
    doc += ("A guide to what the GAN language can do, one feature at a time. "
            "Every example below was run with the real `gan.exe` when this page was written.\n\n")
    for s in steps:
        if s["status"] != "done" or not (s["examples"] or s["error_examples"]):
            continue
        doc += f"## {s['title']}\n\n"
        example_facts = ""
        for f in s["examples"] + s["error_examples"]:
            if (ROOT / f).exists():
                example_facts += f"Example {f}:\n{read_text(ROOT / f)}\nIts real output:\n{run_gan(ROOT / f)}\n\n"
        facts = f"{s['summary']}\n" + "".join(f"- {x}\n" for x in s["learned"]) + "\n" + example_facts
        text = writer.ask(f"How to use this GAN feature: {s['title']}", facts)
        if text:
            doc += text + "\n\n"
        else:
            doc += s["summary"] + "\n\n"
        for f in s["examples"]:
            doc += gan_block(f)
        if s["error_examples"]:
            doc += "**What an error looks like**\n\n"
            for f in s["error_examples"]:
                doc += gan_block(f)
        if s["learned"]:
            doc += "**Good to know**\n\n" + "".join(f"- {x}\n" for x in s["learned"]) + "\n"
    return doc + footer(writer)


def build_how_it_works(facts, writer, today):
    doc = header("How GAN works inside", today)
    token_lines = "\n".join(f"- {n} {('(' + c + ')') if c else ''}" for n, c in facts["tokens"])
    stmt_lines = "\n".join(f"- {n}: {c}" for n, c in facts["statements"])
    expr_lines = "\n".join(f"- {n}: {c}" for n, c in facts["expressions"])
    func_lines = "\n".join(f"- {n}: {c}" for n, c in facts["parser_functions"])

    overview = (
        "GAN runs a program in three stages: the lexer (lexer.go) turns the text into tokens, "
        "the parser (parser.go) turns the tokens into statements and expressions, "
        "and the interpreter (interpreter.go) runs the statements. main.go reads the .gan file and connects them.\n"
        f"Token types: \n{token_lines}\n"
    )
    doc += "## The three stages\n\n"
    doc += (writer.ask("The three stages of GAN: lexer, parser, interpreter", overview) or
            "GAN runs a program in three stages: the **lexer** turns text into tokens, the **parser** turns "
            "tokens into statements, and the **interpreter** runs them. `main.go` connects them.") + "\n\n"

    doc += "## Tokens (read from lexer.go)\n\n" + "\n".join(
        f"- `{n}`" + (f": {c}" if c else "") for n, c in facts["tokens"]) + "\n\n"

    doc += "## Statements and expressions (read from parser.go)\n\n"
    doc += ("A **statement** is something the program does. An **expression** is something that has a value.\n\n")
    doc += "**Statements**\n\n" + "\n".join(f"- `{n}`" + (f": {c}" if c else "") for n, c in facts["statements"]) + "\n\n"
    doc += "**Expressions**\n\n" + "\n".join(f"- `{n}`" + (f": {c}" if c else "") for n, c in facts["expressions"]) + "\n\n"

    prec_facts = (
        "The parser has one function per precedence level. Each level calls the tighter one below it, "
        "so multiplication is grabbed before addition. Parser functions:\n" + func_lines + "\n"
        "Example: in 2 + 3 * 4 the parser grabs 3 * 4 first, giving 2 + (3 * 4) = 14.\n"
    )
    doc += "## How precedence works\n\n"
    doc += (writer.ask("How the GAN parser handles operator precedence", prec_facts) or
            "The parser has one function per precedence level, each calling the tighter one below it.") + "\n\n"
    doc += "Parser functions (read from parser.go):\n\n" + "\n".join(
        f"- `{n}`" + (f": {c}" if c else "") for n, c in facts["parser_functions"]) + "\n\n"

    doc += "## Every error message GAN can give (read from the Go source)\n\n"
    doc += ("Messages starting with `Syntax error:` come from the parser. "
            "Messages starting with `Runtime error:` come from the interpreter. "
            "`%q` and `%s` are filled in with a name or value when the error happens.\n\n")
    doc += "\n".join(f"- {m}" for m in facts["errors"]) + "\n"
    return doc + footer(writer)


# ---------------------------------------------------------------- main
def generate(model, use_ai):
    """Read the project and write all the documents once."""
    if not GAN_EXE.exists():
        print(f"! Cannot find {GAN_EXE}. Build it first:  go build -o gan.exe .")
        return
    try:
        steps = json.loads(read_text(STEPS_FILE))
    except ValueError as e:
        print(f"! steps.json is not valid JSON right now ({e}). Fix it and save again.")
        return
    facts = read_source_facts()
    today = datetime.date.today().isoformat()

    print(f"Documenter: reading project at {ROOT}")
    writer = Writer(model, enabled=use_ai)
    print("AI explanations:", f"on ({model})" if writer.enabled else "off")

    OUT_DIR.mkdir(exist_ok=True)
    documents = {
        "00-progress.md": lambda: build_progress(steps, writer, today),
        "01-language-guide.md": lambda: build_guide(steps, writer, today),
        "02-how-it-works.md": lambda: build_how_it_works(facts, writer, today),
    }
    for name, build in documents.items():
        print(f"Writing {name}")
        (OUT_DIR / name).write_text(build(), encoding="utf-8")

    print(f"\nDone. {len(documents)} documents are in {OUT_DIR}")
    if writer.enabled:
        print(f"The model was asked {writer.calls} new question(s); the rest came from the cache.")


def watched_files():
    """The files whose changes should trigger a new run, with their last-changed times."""
    paths = list(ROOT.glob("*.go")) + list(ROOT.glob("*.gan")) + [GAN_EXE, STEPS_FILE]
    snapshot = {}
    for p in paths:
        try:
            snapshot[str(p)] = p.stat().st_mtime
        except OSError:
            pass  # file missing or being written right now
    return snapshot


def watch(model, use_ai):
    """Keep running: write the documents, then wait for the project to change and write them again."""
    generate(model, use_ai)
    seen = watched_files()
    print("\nWatching for changes (Go files, .gan files, gan.exe, steps.json). Press Ctrl+C to stop.")
    while True:
        time.sleep(3)
        now = watched_files()
        if now == seen:
            continue
        # something changed: wait until it has been quiet for 3 seconds (you may still be saving)
        while True:
            time.sleep(3)
            again = watched_files()
            if again == now:
                break
            now = again
        changed = sorted(Path(p).name for p in set(now) | set(seen) if now.get(p) != seen.get(p))
        print(f"\n[{datetime.datetime.now():%H:%M:%S}] changed: {', '.join(changed)}")
        generate(model, use_ai)
        seen = watched_files()
        print("\nWatching for changes... (Ctrl+C to stop)")


def main():
    ap = argparse.ArgumentParser(description="Write tutorial documents about the GAN language.")
    ap.add_argument("--model", default=DEFAULT_MODEL, help=f"Ollama model to use (default {DEFAULT_MODEL})")
    ap.add_argument("--no-ai", action="store_true", help="facts only, do not call the AI")
    ap.add_argument("--watch", action="store_true", help="keep running and update the documents when the project changes")
    args = ap.parse_args()

    if args.watch:
        try:
            watch(args.model, not args.no_ai)
        except KeyboardInterrupt:
            print("\nDocumenter stopped.")
    else:
        generate(args.model, not args.no_ai)


if __name__ == "__main__":
    main()
