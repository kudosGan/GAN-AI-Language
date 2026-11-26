// FeatureCard: one small box with a title and a sentence.
import "./FeatureCard.css";

type Props = { title: string; text: string };

export default function FeatureCard({ title, text }: Props) {
  return (
    <div className="feature">
      <h3>{title}</h3>
      <p>{text}</p>
    </div>
  );
}
