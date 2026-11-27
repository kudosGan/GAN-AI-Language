// Home page: hero, the three feature cards, what's new, and a link to the lessons.
import Hero from "../components/Hero";
import FeatureCard from "../components/FeatureCard";
import Showcase from "../components/Showcase";
import { FEATURES, LESSONS } from "../data/content";
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
        <Showcase />
        <p className="page-cta">
          New to programming? <a href="#learn">Start with the {LESSONS.length}-lesson tutorial.</a>
        </p>
      </div>
    </>
  );
}
