// About page: a short background of the language.
import "./Page.css";

export default function About() {
  return (
    <div className="container page">
      <h1>About GAN</h1>
      <p>
        GAN is a programming language built step by step, as a way to learn how languages work.
        It is written in Go. A program goes through three stages: the <b>lexer</b> cuts the text
        into words, the <b>parser</b> builds a tree of meaning, and the <b>interpreter</b> runs it.
      </p>
      <p>
        It has numbers, text, true/false, lists, loops, functions, optional types and libraries.
        It is small on purpose, so the whole language fits in your head.
      </p>
    </div>
  );
}
