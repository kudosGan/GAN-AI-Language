// Header: the dark top bar with the logo and the menu.
import { MENU } from "../data/content";
import "./Header.css";

type Props = { current: string }; // the page that is open now

export default function Header({ current }: Props) {
  return (
    <header className="header">
      <div className="container header-row">
        {/* logo links back to the home page */}
        <a className="header-logo" href="#home">
          <img src="./gan-icon.svg" alt="GAN" />
        </a>
        <nav className="header-menu">
          {MENU.map((item) => (
            <a
              key={item.page}
              href={"#" + item.page}
              className={current === item.page ? "active" : ""} /* highlight the open page */
            >
              {item.label}
            </a>
          ))}
        </nav>
      </div>
    </header>
  );
}
