// LibraryCard: one library with the list of its functions.
import CodeBox from "./CodeBox";
import "./LibraryCard.css";

type Props = { name: string; summary: string; functions: string[] };

export default function LibraryCard({ name, summary, functions }: Props) {
  return (
    <section className="library">
      <h2>{name}</h2>
      <p>{summary}</p>
      <CodeBox code={`import "${name}"`} />
      <p className="library-funcs">
        {functions.map((f) => (
          <code key={f}>{f}</code> /* one little tag per function */
        ))}
      </p>
    </section>
  );
}
