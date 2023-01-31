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
      authToken: "f6efd5c0ab184f3a9b108519f5e9aee4be7f1dd9363d44bd9efd0e16c98f4a0b",
      org: "nutty",
      project: "dg-panel-frontend",
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
