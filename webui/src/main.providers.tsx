import { BrowserRouter } from "react-router-dom";
import { RecoilRoot } from "recoil";
import { theme } from "@styles/theme";
import React, { FC, PropsWithChildren } from "react";
import { MantineProvider } from "@mantine/core";
import { NotificationsProvider } from "@mantine/notifications";
import { AuthProvider } from "@stores/auth/provider";
import { ModalsProvider } from "@mantine/modals";

export const MainProviders: FC<PropsWithChildren<{}>> = ({ children }) => (
  <BrowserRouter>
    <RecoilRoot>
      <MantineProvider theme={theme} withGlobalStyles withNormalizeCSS>
        <NotificationsProvider>
          <ModalsProvider>
            <AuthProvider>{children}</AuthProvider>
          </ModalsProvider>
        </NotificationsProvider>
      </MantineProvider>
    </RecoilRoot>
  </BrowserRouter>
);
