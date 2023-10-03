import { LoadingOverlay } from "@mantine/core";
import React from "react";
import { Outlet } from "react-router-dom";

export const PageLoader = () => (
  <React.Suspense fallback={<LoadingOverlay visible overlayProps={{ blur: 7 }} />}>
    <Outlet />
  </React.Suspense>
);
