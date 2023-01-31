import React from "react";
import ReactDOM from "react-dom/client";
import * as Sentry from "@sentry/react";
import { BrowserTracing } from "@sentry/tracing";
import { Router } from "./pages/Router";
import { Navbar } from "@components/Navbar/navbar";
import { getHostname } from "@src/helpers/axiosInstance";
import { MainProviders } from "@src/main.providers";
import { AppShell } from "@mantine/core";

import "./styles/reset.css";
import "./styles/util.scss";
import "./styles/fonts/GreycliffCF/styles.css";
import "@degrens-21/fa-6/css/all.css";
import { navbarState } from "./stores/navbar/state";
import { useRecoilValue } from "recoil";

if (!import.meta.env.DEV) {
  Sentry.init({
    dsn: "https://b75857b005154a9e80b44f55fb86fd07@sentry.nuttyshrimp.me/12",
    integrations: [
      new BrowserTracing({
        tracingOrigins: [getHostname()],
      }),
    ],
    release: "1.0.0",
    environment: import.meta.env.MODE,
    normalizeDepth: 10,
    attachStacktrace: true,
    tracesSampleRate: 1.0,
  });
}

const App = () => {
  const isOpen = useRecoilValue(navbarState.open);
  return (
    <div>
      <AppShell
        padding={"md"}
        navbar={<Navbar />}
        sx={{
          "--mantine-navbar-width": isOpen ? "270px" : "70px",
        }}
      >
        <Router />
      </AppShell>
    </div>
  );
};

const rootElem = document.getElementById("root");
const root = ReactDOM.createRoot(rootElem as HTMLElement);
root.render(
  <React.StrictMode>
    <Sentry.ErrorBoundary fallback={<p>Woops an error happend, try reloading the page</p>}>
      <MainProviders>
        <App />
      </MainProviders>
    </Sentry.ErrorBoundary>
  </React.StrictMode>
);
