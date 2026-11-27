// Learn page: the tutorial lessons in order.
import LessonList from "../components/LessonList";
import { LESSONS } from "../data/content";
import "./Page.css";

export default function Learn() {
  return (
    <div className="container page">
      <h1>Learn GAN</h1>
      <p>
        {LESSONS.length} short lessons that start from nothing and end with a model that learns.
        Their examples are checked against the real GAN program.
      </p>
      <LessonList />
    </div>
  );
}
