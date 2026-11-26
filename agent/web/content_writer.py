"""Agent 2: the Content Writer.

The FACTS (lesson titles, library function names) are read straight from the real GAN files,
so they can never be invented. The AI only writes the friendly wording: the three "why GAN"
cards and the library summaries. The result is written to src\\data\\content.ts.
"""
import json
import re

from common import ROOT, SRC, ask_ai, backup_file, pull_json, say

NAME = "content-writer"
CONTENT = SRC / "data" / "content.ts"


def read_lessons():
    """Lesson titles come from the first line of each tutorial file ('# Lesson 3: Numbers and math')."""
    titles = []
    for path in sorted((ROOT / "tutorial").glob("[0-9][0-9]-*.md")):
        first = path.read_text(encoding="utf-8").splitlines()[0]
        titles.append(re.sub(r"^#\s*Lesson\s*\d+:\s*", "", first).strip())
    return titles


def read_libraries():
    """Library names and their functions come from libs\\*.gan (every 'func name' and top-level 'final NAME')."""
    libs = []
    for path in sorted((ROOT / "libs").glob("*.gan")):
        text = path.read_text(encoding="utf-8")
        names = re.findall(r"^final\s+(\w+)", text, re.M) + re.findall(r"^func\s+(\w+)", text, re.M)
        libs.append({"name": path.stem, "functions": names})
    return libs


def old_value(name, default):
    """Keep REPO_URL / VERSION the way they are in the current content.ts."""
    match = re.search(rf'export const {name} = "([^"]*)"', CONTENT.read_text(encoding="utf-8"))
    return match.group(1) if match else default


def run():
    lessons = read_lessons()
    libs = read_libraries()
    say(NAME, f"found {len(lessons)} lessons and {len(libs)} libraries in the real GAN files")

    say(NAME, "asking the AI for the wording...")
    answer = ask_ai(
        "You write short website text for GAN, a small friendly programming language for learning "
        "(clear words, friendly error messages with the line number, lists, loops, functions, optional "
        "types, libraries). Reply with ONLY JSON in this shape: "
        '{"features":[{"title":"...","text":"..."},{"title":"...","text":"..."},{"title":"...","text":"..."}],'
        '"summaries":{' + ", ".join(f'"{l["name"]}":"..."' for l in libs) + "}}. "
        "The 3 features are BENEFITS for a beginner, in this order: 1) easy to read, 2) friendly error messages, "
        "3) useful things built in. Do NOT name library functions in them. "
        "Titles are 2-3 words. Each text is one plain sentence under 20 words. Each summary is a short phrase "
        "of under 8 words saying what the library is for. "
        "Libraries: " + ", ".join(f'{l["name"]} ({", ".join(l["functions"][:6])})' for l in libs) + "."
    )
    words = pull_json(answer) or {}
    features = [f for f in words.get("features", []) if isinstance(f, dict) and f.get("title") and f.get("text")][:3]
    summaries = words.get("summaries", {}) if isinstance(words.get("summaries"), dict) else {}

    # fall back to the current wording for anything the AI did not give
    current = CONTENT.read_text(encoding="utf-8")
    if len(features) < 3:
        say(NAME, "AI wording incomplete; keeping the current feature cards")
        features = None

    backup_file(CONTENT)
    out = ["// All the words and lists shown on the site. WRITTEN BY the content-writer agent:",
           "// facts come from the real GAN files, wording from the AI. Edit by hand if you like.", ""]
    out.append(f'export const REPO_URL = {json.dumps(old_value("REPO_URL", "https://github.com/YOUR-NAME/GAN"))};')
    out.append(f'export const VERSION = {json.dumps(old_value("VERSION", "2.0"))};')
    out.append("")
    # keep the sample program and the menu exactly as they are now
    sample = re.search(r"export const SAMPLE_CODE = `.*?`;", current, re.S)
    menu = re.search(r"export const MENU = \[.*?\n\];", current, re.S)
    out.append(sample.group(0) if sample else "export const SAMPLE_CODE = `print(\"Hello\")`;")
    out.append("")
    if features:
        out.append("export const FEATURES = " + json.dumps(features, indent=2) + ";")
    else:
        keep = re.search(r"export const FEATURES = \[.*?\n\];", current, re.S)
        out.append(keep.group(0))
    out.append("")
    out.append("export const LESSONS = " + json.dumps(lessons, indent=2) + ";")
    out.append("")
    libs_out = [{"name": l["name"], "summary": str(summaries.get(l["name"], f"The {l['name']} library.")),
                 "functions": l["functions"]} for l in libs]
    out.append("export const LIBRARIES = " + json.dumps(libs_out, indent=2) + ";")
    out.append("")
    out.append(menu.group(0) if menu else "export const MENU = [];")
    CONTENT.write_text("\n".join(out) + "\n", encoding="utf-8")
    say(NAME, "content.ts written")
    return True


if __name__ == "__main__":
    run()
