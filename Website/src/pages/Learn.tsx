// Learn page: the tutorial lessons in order.
import LessonList from "../components/LessonList";
import "./Page.css";

export default function Learn() {
  return (
    <div className="container page">
      <h1>Learn GAN</h1>
      <p>
        Fourteen short lessons that start from nothing. Every example was run with the real
        GAN program, so what you read is what you will see.
      </p>
      <LessonList />
    </div>
  );
}
