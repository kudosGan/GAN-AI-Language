// Learn page: the tutorial lessons in order.
import LessonList from "../components/LessonList";
import { LESSONS } from "../data/content";
import "./Page.css";

export default function Learn() {
  return (
    <div className="container page">
      <h1>Learn GAN</h1>
      <p>
        {LESSONS.length} short lessons that start from nothing and go on to models that learn, AI models, Python and tools.
        Examples with fixed output are checked against the real GAN program.
      </p>
      <LessonList />
    </div>
  );
}
