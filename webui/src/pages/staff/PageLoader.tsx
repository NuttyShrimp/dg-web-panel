import { LoadingOverlay } from "@mantine/core";
import React from "react";
import { Outlet } from "react-router-dom";

export const PageLoader = () => (
  <React.Suspense fallback={<LoadingOverlay visible overlayBlur={7} />}>
    <Outlet />
  </React.Suspense>
);
