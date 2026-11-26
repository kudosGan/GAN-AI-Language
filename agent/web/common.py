"""Shared helpers for the five website agents.

Every agent uses the same few things: where the folders are, how to talk to the local
Ollama AI, and how to run a command (like "npm run build"). They live here once.
"""
import json
import re
import shutil
import subprocess
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]     # C:\Desktop\GAN
SITE = ROOT / "Website"                         # the React project
SRC = SITE / "src"
COMPONENTS = SRC / "components"
PAGES = SRC / "pages"
WEB_AGENT = Path(__file__).resolve().parent     # agent\web
BACKUP = WEB_AGENT / "backup"                   # copies of files before an agent changes them
DESIGN_FILE = WEB_AGENT / "design.json"         # the Designer's brief

OLLAMA_URL = "http://localhost:11434/api/generate"
FAST_MODEL = "qwen2.5:7b"     # quick: design and wording
CODE_MODEL = "qwen2.5:14b"    # slower but better at code


def say(agent, message):
    """Print one line that starts with the agent's name."""
    print(f"[{agent}] {message}", flush=True)


def ask_ai(prompt, model=FAST_MODEL, retries=3):
    """Send a prompt to the local Ollama AI and return its answer as text.
    Returns None if the AI cannot be reached (the caller then keeps the old files)."""
    body = json.dumps({"model": model, "prompt": prompt, "stream": False,
                       "options": {"temperature": 0.3}}).encode()
    for attempt in range(1, retries + 1):
        try:
            req = urllib.request.Request(OLLAMA_URL, data=body, headers={"Content-Type": "application/json"})
            with urllib.request.urlopen(req, timeout=900) as resp:
                return json.loads(resp.read())["response"]
        except Exception as err:  # Ollama sometimes fails the first call: try again
            print(f"  AI call failed ({err}), try {attempt} of {retries}")
    return None


def pull_json(text):
    """Find the first {...} block in an AI answer and read it as JSON (None if it is not valid)."""
    if not text:
        return None
    match = re.search(r"\{.*\}", text, re.S)
    if not match:
        return None
    try:
        return json.loads(match.group(0))
    except ValueError:
        return None


def pull_block(text, language):
    """Find a ```language ... ``` block in an AI answer and return its inside (None if missing)."""
    if not text:
        return None
    match = re.search(r"```" + language + r"\s*\n(.*?)```", text, re.S)
    return match.group(1).rstrip() + "\n" if match else None


def backup_file(path):
    """Keep a copy of a file before it is changed, so the Reviewer can put it back."""
    target = BACKUP / path.relative_to(SITE)
    if not target.exists():            # only the FIRST copy counts: it is the last good version
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(path, target)


def restore_all():
    """Put every backed-up file back where it came from."""
    if not BACKUP.exists():
        return 0
    count = 0
    for saved in BACKUP.rglob("*"):
        if saved.is_file():
            shutil.copy2(saved, SITE / saved.relative_to(BACKUP))
            count += 1
    return count


def clear_backups():
    """Forget the backups (called after a build that works)."""
    if BACKUP.exists():
        shutil.rmtree(BACKUP)


def run(command, cwd):
    """Run a command line (a list of words) and return (exit code, all output)."""
    done = subprocess.run(command, cwd=cwd, capture_output=True, text=True, shell=True)
    return done.returncode, (done.stdout or "") + (done.stderr or "")


def component_names():
    """The names of all components, found by looking at the .tsx files."""
    return sorted(p.stem for p in COMPONENTS.glob("*.tsx"))
