import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// base "./" makes every link relative, so the site works on GitHub Pages
// (which serves it from a sub-folder like /GAN/) and on any other host.
export default defineConfig({
  plugins: [react()],
  base: "./",
});
