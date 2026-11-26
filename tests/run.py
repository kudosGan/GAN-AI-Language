"""
Regression tests for GAN.

Runs gan.exe on every .gan file in the project folder and compares what it prints with
tests/expected/<name>.out. If tests/expected/<name>.in exists, it is fed to the program as input.

  python tests\run.py            check everything
  python tests\run.py --update   (re)write the expected files from the current output
                                 (only do this after checking the output by eye!)
"""
import subprocess, sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
EXPECTED = ROOT / "tests" / "expected"
GAN = ROOT / "gan.exe"

def run(gan_file):
    stdin_file = EXPECTED / (gan_file.stem + ".in")
    stdin = stdin_file.read_text(encoding="utf-8") if stdin_file.exists() else ""
    r = subprocess.run([str(GAN), str(gan_file)], input=stdin, capture_output=True,
                       text=True, timeout=20, encoding="utf-8", cwd=str(ROOT))
    return (r.stdout + r.stderr).replace("\r\n", "\n")

def main():
    update = "--update" in sys.argv
    files = sorted(ROOT.glob("*.gan"))
    passed = failed = 0
    for f in files:
        out = run(f)
        exp_file = EXPECTED / (f.stem + ".out")
        if update:
            exp_file.write_text(out, encoding="utf-8", newline="\n")
            continue
        if not exp_file.exists():
            print(f"NEW   {f.name} (no expected file yet)")
            failed += 1
            continue
        if out == exp_file.read_text(encoding="utf-8"):
            passed += 1
        else:
            failed += 1
            print(f"FAIL  {f.name}\n--- expected\n{exp_file.read_text(encoding='utf-8')}--- actual\n{out}")
    if update:
        print(f"wrote expected output for {len(files)} files")
    else:
        print(f"{passed} passed, {failed} failed")
        sys.exit(1 if failed else 0)

main()
