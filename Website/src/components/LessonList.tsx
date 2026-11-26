// LessonList: the list of tutorial lessons.
import { LESSONS, REPO_URL } from "../data/content";
import "./LessonList.css";

export default function LessonList() {
  return (
    <ul className="lessons">
      {LESSONS.map((title, i) => (
        <li key={title}>
          {/* each row opens the tutorial folder in the source repository */}
          <a href={`${REPO_URL}/tree/main/tutorial`}>
            <span className="lesson-num">Lesson {i}</span> {title}
          </a>
        </li>
      ))}
    </ul>
  );
}
