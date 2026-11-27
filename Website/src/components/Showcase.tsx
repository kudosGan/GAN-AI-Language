// Showcase: "what's new" - sample programs with what they print.
import CodeBox from "./CodeBox";
import { BETA_NOTE, SHOWCASE, SHOWCASE_INTRO } from "../data/showcase";
import { VERSION } from "../data/content";
import "./Showcase.css";

export default function Showcase() {
  return (
    <section className="showcase">
      <h2>New in GAN {VERSION}</h2>
      <p className="showcase-intro">{SHOWCASE_INTRO}</p>

      {SHOWCASE.map((sample) => (
        <article key={sample.file} className="showcase-item">
          <h3>{sample.title}</h3>
          <p>{sample.text}</p>
          {/* the program on one side, what it prints on the other */}
          <div className="showcase-pair">
            <CodeBox title={sample.file} code={sample.code} />
            {/* a model's answer is different every time, so say it is only an example */}
            <CodeBox title={sample.varies ? "an example answer (yours may differ)" : "what it prints"} code={sample.output} />
          </div>
        </article>
      ))}

      <p className="showcase-note">{BETA_NOTE}</p>
    </section>
  );
}
