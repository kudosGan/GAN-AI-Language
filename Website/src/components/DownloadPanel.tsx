// DownloadPanel: the big "get GAN" box with the install steps.
import { REPO_URL, VERSION } from "../data/content";
import CodeBox from "./CodeBox";
import "./DownloadPanel.css";

export default function DownloadPanel() {
  return (
    <div className="download">
      <h2>GAN {VERSION} for Windows</h2>
      {/* the installer file is attached to a "Release" on GitHub */}
      <a className="btn" href={`${REPO_URL}/releases/latest`}>Download GAN-Setup.exe</a>
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
