"""Agent 4: the Reviewer.

Checks the site: every component has a .tsx AND a .css, and `npm run build` (TypeScript + Vite)
passes. When the build fails it asks the Developer to repair the broken file, up to 3 times.
If it still fails, every file the agents changed is put back from the backup, so the site is
never left broken.
"""
import re

import developer
from common import COMPONENTS, SITE, clear_backups, component_names, restore_all, run, say

NAME = "reviewer"


def check_pairs():
    """Every component must have both a .tsx and a .css file."""
    problems = []
    for name in component_names():
        if not (COMPONENTS / f"{name}.css").exists():
            problems.append(f"{name} has no .css file")
    return problems


def broken_component(output):
    """Find which component the error message points at (looks for 'components/Name.tsx')."""
    match = re.search(r"components[/\\](\w+)\.(?:tsx|css)", output)
    return match.group(1) if match else None


def run_review(max_fixes=3):
    for problem in check_pairs():
        say(NAME, "PROBLEM: " + problem)
        return False

    for attempt in range(max_fixes + 1):
        say(NAME, "running npm run build ...")
        code, output = run(["npm", "run", "build"], SITE)
        if code == 0:
            say(NAME, "build passed")
            clear_backups()
            return True
        say(NAME, "build FAILED:\n" + output[-800:])
        culprit = broken_component(output)
        if attempt == max_fixes or not culprit:
            break
        say(NAME, f"asking the developer to fix {culprit} (try {attempt + 1} of {max_fixes})")
        developer.fix(culprit, output)

    count = restore_all()
    say(NAME, f"could not fix it; put back {count} file(s) from the backup. The site is as it was.")
    run(["npm", "run", "build"], SITE)
    return False


if __name__ == "__main__":
    run_review()
