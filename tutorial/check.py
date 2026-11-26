"""
Checks (and fills in) the examples in the GAN tutorial.

Every example in a lesson looks like this:

    ```gan
    print("Hello")
    ```
    ```input        <- optional: what the user types, one line per input() call
    Ann
    ```
    ```output
    Hello
    ```

This script runs the gan block with the real gan.exe and compares what it prints with the
output block. So the tutorial can never show something GAN does not really do.

  python tutorial\\check.py           check every lesson (exit code 1 if anything is wrong)
  python tutorial\\check.py --fill    write the real output into every output block
                                     (then read the result: it is what GAN really printed)
A gan block with no output block after it is only shown, never run.
"""
import re
import subprocess
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent
GAN = ROOT / "gan.exe"
RUN_DIR = HERE / ".run"

# ```gan ... ```  then optionally ```input ... ```  then ```output ... ```
# (?:(?!```).)* means "any characters, but never a ``` fence", so one block can never swallow the next
NO_FENCE = r"(?:(?!```).)*"
BLOCKS = re.compile(
    r"(?P<gan>```gan\n(?P<code>" + NO_FENCE + r")```\n)"
    r"(?:\n)?"
    r"(?P<inp>```input\n(?P<typed>" + NO_FENCE + r")```\n(?:\n)?)?"
    r"(?P<out>```output\n(?P<expected>" + NO_FENCE + r")```\n)",
    re.S,
)


def run_code(code, typed, name):
    RUN_DIR.mkdir(exist_ok=True)
    path = RUN_DIR / f"{name}.gan"
    path.write_text(code, encoding="utf-8", newline="\n")
    result = subprocess.run(
        [str(GAN), str(path)], input=typed, capture_output=True, text=True,
        timeout=20, encoding="utf-8", cwd=str(ROOT),
    )
    return (result.stdout + result.stderr).replace("\r\n", "\n").rstrip("\n")


def main():
    fill = "--fill" in sys.argv
    lessons = sorted(HERE.glob("*.md"))
    total = bad = 0
    for lesson in lessons:
        text = lesson.read_text(encoding="utf-8")
        count = 0

        def handle(m):
            nonlocal total, bad, count
            count += 1
            total += 1
            actual = run_code(m.group("code"), m.group("typed") or "", f"{lesson.stem}-{count}")
            expected = m.group("expected").rstrip("\n")
            if fill:
                return (m.group("gan") + "\n" + (m.group("inp") or "") +
                        "```output\n" + actual + ("\n" if actual else "") + "```\n")
            if actual != expected:
                bad += 1
                print(f"MISMATCH in {lesson.name}, example {count}\n--- the lesson says\n{expected}\n--- GAN printed\n{actual}\n")
            return m.group(0)

        new_text = BLOCKS.sub(handle, text)
        if fill and new_text != text:
            lesson.write_text(new_text, encoding="utf-8", newline="\n")
        print(f"{lesson.name:<28} {count} example(s)")
    if fill:
        print(f"\nfilled {total} output block(s) from real runs")
    else:
        print(f"\n{total - bad} of {total} examples are correct" + ("" if not bad else f"  ({bad} WRONG)"))
        sys.exit(1 if bad else 0)


main()
