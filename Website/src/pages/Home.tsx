// Home page: hero, the three feature cards, and a link to the lessons.
import Hero from "../components/Hero";
import FeatureCard from "../components/FeatureCard";
import { FEATURES } from "../data/content";
import "./Page.css";

export default function Home() {
  return (
    <>
      <Hero />
      <div className="container page">
        <h2>Why GAN?</h2>
        <div className="grid-3">
          {FEATURES.map((f) => (
            <FeatureCard key={f.title} title={f.title} text={f.text} />
          ))}
        </div>
        <p className="page-cta">
          New to programming? <a href="#learn">Start with the 14-lesson tutorial.</a>
        </p>
      </div>
    </>
  );
}
