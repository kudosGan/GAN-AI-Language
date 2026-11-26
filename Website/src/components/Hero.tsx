// Hero: the big opening section of the home page.
import CodeBox from "./CodeBox";
import { SAMPLE_CODE } from "../data/content";
import "./Hero.css";

export default function Hero() {
  return (
    <section className="hero">
      <div className="container hero-row">
        {/* left side: the message and the two buttons */}
        <div className="hero-text">
          <h1>Learn to program with GAN</h1>
          <p>
            GAN is a small, friendly language with clear words and helpful errors.
            Made for learning, simple enough to read on the first day.
          </p>
          <a className="btn" href="#download">Download GAN</a>
          <a className="btn ghost" href="#learn">Start learning</a>
        </div>
        {/* right side: a real program */}
        <div className="hero-code">
          <CodeBox title="hello.gan" code={SAMPLE_CODE} />
        </div>
      </div>
    </section>
  );
}
