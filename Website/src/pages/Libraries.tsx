// Libraries page: the standard libraries and their functions.
import LibraryCard from "../components/LibraryCard";
import { LIBRARIES } from "../data/content";
import "./Page.css";

export default function Libraries() {
  return (
    <div className="container page">
      <h1>Libraries</h1>
      <p>Add a library to your program with <code>import</code>, then use its functions with a dot.</p>
      {LIBRARIES.map((lib) => (
        <LibraryCard key={lib.name} name={lib.name} summary={lib.summary} functions={lib.functions} />
      ))}
    </div>
  );
}
