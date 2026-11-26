# The website team (5 agents)

They work on the React site in `Website\`. Ollama must be running (the same one the documenter uses).

| # | Agent | File | Job |
|---|---|---|---|
| 1 | Designer | `designer.py` | Picks the colour palette, writes `design.json` and the colours in `src\index.css` |
| 2 | Content Writer | `content_writer.py` | Rebuilds `src\data\content.ts` from the real lessons and libraries; AI writes the friendly wording |
| 3 | Developer | `developer.py` | Rewrites each component's `.tsx` and `.css` (model `qwen2.5:14b`) |
| 4 | Reviewer | `reviewer.py` | Runs `npm run build`; asks the Developer to fix errors (3 tries); restores the backup if it cannot |
| 5 | Publisher | `publisher.py` | Builds, zips to `Website\GAN-website.zip`, and with `--push` commits and pushes to GitHub |

```
python agent\web\web_team.py                       # all five, no upload
python agent\web\web_team.py --only designer       # one agent
python agent\web\web_team.py --only developer --component Hero
python agent\web\web_team.py --push                # also push to GitHub
```

Safety: before an agent changes a file, the first copy goes to `agent\web\backup\`. If the Reviewer
cannot get a passing build, the copies are put back. GitHub then builds and publishes the site by itself
(`.github\workflows\pages.yml`).
