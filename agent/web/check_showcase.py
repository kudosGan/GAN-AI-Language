"""Checks the samples on the website's home page against the real gan.exe.

    python agent\\web\\check_showcase.py

Fails if a sample prints something different from what the website says it prints.
"""
import re
import subprocess
import sys
import tempfile
from pathlib import Path

from common import ROOT, SRC


def main():
    text = (SRC / "data" / "showcase.ts").read_text(encoding="utf-8")
    found = re.findall(r'file: "([^"]+)",\s*code: `(.*?)`,\s*output: `(.*?)`,(\s*varies: true,)?', text, re.S)
    for name, _, _, varies in found:
        if varies:
            print(f"skip  {name} (an AI model's answer, different every time)")
    samples = [(name, code, output) for name, code, output, varies in found if not varies]
    if not samples:
        print("no samples found in showcase.ts")
        return 1
    bad = 0
    with tempfile.TemporaryDirectory() as folder:
        for name, code, expected in samples:
            program = Path(folder) / name
            # inside a TypeScript `text`, two backslashes stand for one
            program.write_text(code.replace("\\\\", "\\") + "\n", encoding="utf-8")
            done = subprocess.run([str(ROOT / "gan.exe"), str(program)], capture_output=True,
                                  text=True, encoding="utf-8", timeout=30)
            actual = (done.stdout + done.stderr).replace("\r\n", "\n").rstrip("\n")
            if actual == expected:
                print(f"ok    {name}")
            else:
                bad += 1
                print(f"WRONG {name}\n--- website says\n{expected}\n--- GAN prints\n{actual}")
    print(f"{len(samples) - bad} of {len(samples)} website samples are correct")
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main())
