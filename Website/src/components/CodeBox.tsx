// CodeBox: shows a piece of GAN code in a dark box.
import "./CodeBox.css";

type Props = { code: string; title?: string }; // title is optional

export default function CodeBox({ code, title }: Props) {
  return (
    <div className="codebox">
      {title && <div className="codebox-title">{title}</div>}
      <pre>
        <code>{code}</code>
      </pre>
    </div>
  );
}
