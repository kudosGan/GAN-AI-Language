// DownloadPanel: the big "get GAN" box with the install steps.
import { REPO_URL, VERSION } from "../data/content";
import CodeBox from "./CodeBox";
import "./DownloadPanel.css";

export default function DownloadPanel() {
  return (
    <div className="download">
      <h2>GAN {VERSION} for Windows</h2>
      {/* the installer is attached to a "Release" on GitHub. This opens the list of releases,
          because GitHub's "latest" link leaves out beta (pre-release) versions */}
      <a className="btn" href={`${REPO_URL}/releases`}>Get GAN-Setup.exe from GitHub</a>
      <p className="download-note">On that page, open the newest release and click GAN-Setup.exe.</p>
      <ol>
        <li>Run <b>GAN-Setup.exe</b> and follow the window. No administrator rights needed.</li>
        <li>Open a <b>new</b> terminal.</li>
        <li>Check that it works:</li>
      </ol>
      <CodeBox code={"gan --version\ngan examples\\hello.gan"} />
      <p className="download-note">
        Windows may show a warning because the installer is not code-signed. Choose "More info", then "Run anyway".
      </p>
    </div>
  );
}
