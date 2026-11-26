"""Agent 3: the Developer.

For each component, shows the AI the current .tsx and .css plus the Designer's brief and asks for
an improved version. The answer is only accepted when it still looks like a valid component
(same name, has a comment, has both files). The Reviewer then checks that it really builds.
"""
import json
import re

from common import (COMPONENTS, CODE_MODEL, DESIGN_FILE, ask_ai, backup_file, component_names,
                    pull_block, say)

NAME = "developer"


def looks_valid(name, tsx, css):
    """A cheap safety check before an AI answer is allowed to replace a file."""
    if tsx is None or css is None:
        return False
    return (
        f"export default function {name}" in tsx          # the component keeps its name
        and ("//" in tsx or "/*" in tsx)                    # it still has comments
        and css.count("{") == css.count("}")                # CSS braces match
    )


def improve(name, instruction="Make it a little cleaner and more polished."):
    """Ask the AI to rewrite one component. Returns True if the files were replaced."""
    tsx_path, css_path = COMPONENTS / f"{name}.tsx", COMPONENTS / f"{name}.css"
    tsx, css = tsx_path.read_text(encoding="utf-8"), css_path.read_text(encoding="utf-8")
    design = DESIGN_FILE.read_text(encoding="utf-8") if DESIGN_FILE.exists() else "{}"

    say(NAME, f"improving {name} ...")
    answer = ask_ai(
        f"You are a React + TypeScript developer. Improve the component '{name}' for the GAN "
        f"programming language website (simple, friendly, like python.org). {instruction}\n"
        "RULES: keep `export default function " + name + "`; keep the same props; keep imports of "
        "data and other components working; keep every comment and add a short comment beside new code; "
        "use only CSS variables like var(--blue) from the stylesheet; no new packages. "
        "Answer with exactly two code blocks, first ```tsx then ```css, and nothing else.\n\n"
        f"Design brief: {design}\n\n```tsx\n{tsx}```\n\n```css\n{css}```",
        model=CODE_MODEL,
    )
    new_tsx, new_css = pull_block(answer, "tsx"), pull_block(answer, "css")
    if not looks_valid(name, new_tsx, new_css):
        say(NAME, f"{name}: the AI answer did not look valid; keeping the old files")
        return False
    backup_file(tsx_path)
    backup_file(css_path)
    tsx_path.write_text(new_tsx, encoding="utf-8")
    css_path.write_text(new_css, encoding="utf-8")
    say(NAME, f"{name}: updated")
    return True


def fix(name, build_error):
    """Called by the Reviewer: ask the AI to repair a component that broke the build."""
    return improve(name, "The build failed with this error, fix it and change nothing else:\n" + build_error[-1500:])


def run(only=None):
    names = [only] if only else component_names()
    changed = [n for n in names if improve(n)]
    say(NAME, f"{len(changed)} of {len(names)} components updated")
    return bool(changed)


if __name__ == "__main__":
    run()
