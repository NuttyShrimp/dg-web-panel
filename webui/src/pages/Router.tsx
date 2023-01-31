// Main file with all routes layed out
// For private routes see: https://stackblitz.com/github/remix-run/react-router/tree/main/examples/auth?file=src%2FApp.tsx
// or https://reactrouterdotcom.fly.dev/docs/en/v6/examples/auth

import { useRoutes } from "react-router-dom";
import { routes } from "./routes";

export const Router = () => {
  const genRoutes = useRoutes(routes);
  return genRoutes;
};
