import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import Checker from "vite-plugin-checker";
import viteSentry from "vite-plugin-sentry";
import tsconfigPaths from "vite-tsconfig-paths";
import svgr from "vite-plugin-svgr";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  base: "/",
  build: {
    sourcemap: true,
  },
  server: {
    port: 3001,
  },
  plugins: [
    react(),
    svgr(),
    tsconfigPaths(),
    Checker({
      typescript: true,
      overlay: false,
      eslint: {
        lintCommand: "eslint --ext ts,tsx src",
      },
    }),
    viteSentry({
      url: "https://sentry.nuttyshrimp.me",
      authToken: "5e2d7e8c0d6a42348a0c50dbf655896524c8414752804c8ea1ca04e357be9cd8",
      org: "nutty",
      project: "degrens-panel-frontend",
      debug: true,
      deploy: {
        env: mode === "production" ? "production" : "development",
      },
      setCommits: {
        auto: true,
      },
      sourceMaps: {
        include: ["../html/assets"],
        ignore: ["node_modules"],
        urlPrefix: "~/assets",
      },
    }),
  ],
}));
