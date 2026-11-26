"""Agent 1: the Designer.

Asks the AI for a colour palette and a font, saves it in design.json, and writes the colours
into the :root block of src\\index.css. Only the colour numbers are changed, nothing else,
so a bad AI answer cannot break the layout.
"""
import json
import re

from common import DESIGN_FILE, SRC, ask_ai, backup_file, pull_json, say

NAME = "designer"
# the colour names the stylesheet uses (each one is "--name" in index.css)
COLOUR_KEYS = ["blue", "blue-dark", "navy", "ink", "muted", "soft", "line", "yellow"]
HEX = re.compile(r"^#[0-9a-fA-F]{6}$")


def run(brief="a clean, simple, friendly site like python.org, with GAN's cyan blue as the main colour"):
    say(NAME, "asking the AI for a colour palette...")
    answer = ask_ai(
        "You are a web designer. Reply with ONLY JSON, no other words. Make a colour palette for "
        f"{brief}. Use exactly these keys, each a 6-digit hex colour like #00add8: "
        + ", ".join(COLOUR_KEYS)
        + ". 'blue' must stay close to #00add8. 'ink' is dark text, 'soft' is a very light panel, "
        "'navy' is a dark bar. Text must be easy to read on the backgrounds."
    )
    palette = pull_json(answer)
    # keep only colours that are really 6-digit hex; anything else is ignored
    good = {k: palette[k] for k in COLOUR_KEYS if palette and HEX.match(str(palette.get(k, "")))}
    if len(good) < len(COLOUR_KEYS):
        say(NAME, f"the AI gave {len(good)} of {len(COLOUR_KEYS)} valid colours; keeping the current design")
        return False

    DESIGN_FILE.write_text(json.dumps({"colors": good}, indent=2), encoding="utf-8")

    css = SRC / "index.css"
    backup_file(css)
    text = css.read_text(encoding="utf-8")
    for key, value in good.items():
        # replace the colour in the line "  --key: #xxxxxx;"
        text = re.sub(rf"(--{key}:\s*)#[0-9a-fA-F]{{6}}", rf"\g<1>{value}", text)
    css.write_text(text, encoding="utf-8")
    say(NAME, "palette saved to design.json and index.css")
    return True


if __name__ == "__main__":
    run()
