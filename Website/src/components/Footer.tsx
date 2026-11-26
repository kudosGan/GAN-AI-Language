// Footer: the bottom bar with quick links and the version line.
import { REPO_URL, VERSION } from "../data/content";
import "./Footer.css";

export default function Footer() {
  return (
    <footer className="footer">
      <div className="container footer-row">
        <div>
          <strong>GAN {VERSION}</strong>
          <p>A small, friendly programming language.</p>
        </div>
        <div className="footer-links">
          <a href="#download">Download</a>
          <a href="#learn">Learn</a>
          <a href="#libraries">Libraries</a>
          <a href={REPO_URL}>Source code</a>
        </div>
      </div>
    </footer>
  );
}
