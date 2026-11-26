"""Agent 5: the Publisher.

1. Builds the site and zips Website\\dist into GAN-website.zip (handy for any other host).
2. With --push it commits and pushes to GitHub. GitHub then builds and publishes the site by
   itself (see .github\\workflows\\pages.yml), so the site goes live a minute later.

Commits are made with YOUR git name and email and contain no extra author lines.
"""
import shutil

from common import ROOT, SITE, run, say

NAME = "publisher"


def build_and_zip():
    say(NAME, "building the site ...")
    code, output = run(["npm", "run", "build"], SITE)
    if code != 0:
        say(NAME, "build failed, nothing published:\n" + output[-600:])
        return False
    archive = shutil.make_archive(str(SITE / "GAN-website"), "zip", SITE / "dist")
    say(NAME, f"zip ready: {archive}")
    return True


def push(message="Update website"):
    code, out = run(["git", "rev-parse", "--is-inside-work-tree"], ROOT)
    if code != 0:
        say(NAME, "this folder is not a git repository yet, so there is nothing to push")
        return False
    run(["git", "add", "-A"], ROOT)
    code, out = run(["git", "commit", "-m", message], ROOT)
    say(NAME, "commit: " + (out.strip().splitlines()[0] if out.strip() else "nothing to commit"))
    code, out = run(["git", "push"], ROOT)
    say(NAME, "pushed to GitHub" if code == 0 else "push failed:\n" + out[-500:])
    return code == 0


def run_publish(do_push=False):
    if not build_and_zip():
        return False
    return push() if do_push else True


if __name__ == "__main__":
    run_publish()
