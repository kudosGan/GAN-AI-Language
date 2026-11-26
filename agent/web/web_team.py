"""The website team: runs the five agents in order.

    python agent\\web\\web_team.py                    run all five (no upload)
    python agent\\web\\web_team.py --only designer    run just one agent
    python agent\\web\\web_team.py --only developer --component Hero   improve one component
    python agent\\web\\web_team.py --push             also commit and push to GitHub at the end

Order: designer -> content-writer -> developer -> reviewer -> publisher
Make sure Ollama is running first (it is, if the documenter works).
"""
import argparse
import sys

import common
import content_writer
import designer
import developer
import publisher
import reviewer

AGENTS = ["designer", "content", "developer", "reviewer", "publisher"]


def main():
    parser = argparse.ArgumentParser(description="Run the GAN website agents")
    parser.add_argument("--only", choices=AGENTS, help="run just this agent")
    parser.add_argument("--component", help="with --only developer: the one component to improve")
    parser.add_argument("--push", action="store_true", help="publisher also pushes to GitHub")
    args = parser.parse_args()

    chosen = [args.only] if args.only else AGENTS
    common.clear_backups()  # start every run with a clean backup folder

    if "designer" in chosen:
        designer.run()
    if "content" in chosen:
        content_writer.run()
    if "developer" in chosen:
        developer.run(args.component)
    ok = True
    if "reviewer" in chosen:
        ok = reviewer.run_review()
    if "publisher" in chosen:
        if ok:
            publisher.run_publish(args.push)
        else:
            common.say("publisher", "the reviewer did not approve the site, so nothing is published")
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
