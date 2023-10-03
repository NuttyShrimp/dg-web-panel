import React from "react";
import ReactDOM from "react-dom/client";
import * as Sentry from "@sentry/react";
import { Router } from "./pages/Router";
import { Navbar } from "@components/Navbar/navbar";
import { getHostname } from "@src/helpers/axiosInstance";
import { MainProviders } from "@src/main.providers";
import { AppShell } from "@mantine/core";

import "./styles/reset.css";
import "./styles/util.scss";
import "./styles/fonts/GreycliffCF/styles.css";

import "@mantine/core/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/dropzone/styles.css";
import "@mantine/notifications/styles.css";
import "@mantine/spotlight/styles.css";

import "@degrens-21/fa-6/css/all.css";
import { navbarState } from "./stores/navbar/state";
import { useRecoilValue } from "recoil";

if (!import.meta.env.DEV) {
  Sentry.init({
    dsn: "https://e301572934fe49f98ad4cf042fe1658c@sentry.nuttyshrimp.me/5",
    integrations: [
      new Sentry.BrowserTracing({
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
        navbar={{
          width: isOpen ? "270px" : "70px",
          breakpoint: "xs",
        }}
        style={{
          "--app-shell-navbar-width": isOpen ? "270px" : "70px",
        }}
      >
        <AppShell.Navbar>
          <Navbar />
        </AppShell.Navbar>
        <AppShell.Main>
          <Router />
        </AppShell.Main>
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
