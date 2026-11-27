// Learn page: the tutorial lessons in order.
import LessonList from "../components/LessonList";
import { LESSONS } from "../data/content";
import "./Page.css";

export default function Learn() {
  return (
    <div className="container page">
      <h1>Learn GAN</h1>
      <p>
        {LESSONS.length} short lessons that start from nothing and go on to a model that learns and talking to an AI model.
        Examples with fixed output are checked against the real GAN program.
      </p>
      <LessonList />
    </div>
  );
}
