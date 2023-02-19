import { BrowserRouter } from "react-router-dom";
import { RecoilRoot } from "recoil";
import { theme } from "@styles/theme";
import React, { FC, PropsWithChildren } from "react";
import { MantineProvider } from "@mantine/core";
import { NotificationsProvider } from "@mantine/notifications";
import { AuthProvider } from "@stores/auth/provider";
import { ModalsProvider } from "@mantine/modals";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./helpers/queryClient";
import { UpdateAnnounceProvider } from "./components/UpdateAnnouncer";

export const MainProviders: FC<PropsWithChildren<{}>> = ({ children }) => (
  <BrowserRouter>
    <QueryClientProvider client={queryClient}>
      <RecoilRoot>
        <MantineProvider theme={theme} withGlobalStyles withNormalizeCSS>
          <NotificationsProvider>
            <ModalsProvider>
              <UpdateAnnounceProvider>
                <AuthProvider>{children}</AuthProvider>
              </UpdateAnnounceProvider>
            </ModalsProvider>
          </NotificationsProvider>
        </MantineProvider>
      </RecoilRoot>
    </QueryClientProvider>
  </BrowserRouter>
);
