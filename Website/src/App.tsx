// App: the frame of the site (Header, one page, Footer).
// Pages are chosen by the "#" part of the address (#learn, #download ...),
// which works on GitHub Pages without any routing library.
import { useEffect, useState } from "react";
import Header from "./components/Header";
import Footer from "./components/Footer";
import Home from "./pages/Home";
import Download from "./pages/Download";
import Learn from "./pages/Learn";
import Libraries from "./pages/Libraries";
import About from "./pages/About";
import "./App.css";

// Reads the page name from the address bar. Empty means "home".
function currentPage(): string {
  return window.location.hash.replace("#", "") || "home";
}

export default function App() {
  const [page, setPage] = useState(currentPage()); // which page is showing

  // When the address changes (a menu click or the Back button), switch page.
  useEffect(() => {
    const onChange = () => {
      setPage(currentPage());
      window.scrollTo(0, 0);
    };
    window.addEventListener("hashchange", onChange);
    return () => window.removeEventListener("hashchange", onChange);
  }, []);

  return (
    <div className="app">
      <Header current={page} />
      <main className="app-main">
        {page === "download" && <Download />}
        {page === "learn" && <Learn />}
        {page === "libraries" && <Libraries />}
        {page === "about" && <About />}
        {/* anything else (including "home") shows the home page */}
        {!["download", "learn", "libraries", "about"].includes(page) && <Home />}
      </main>
      <Footer />
    </div>
  );
}
