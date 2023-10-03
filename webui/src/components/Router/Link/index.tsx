import { useMantineTheme } from "@mantine/core";
import { FC, PropsWithChildren } from "react";
import { Link as RRLink, LinkProps } from "react-router-dom";

import "./style.scss";

export const Link: FC<PropsWithChildren<LinkProps & { noColor?: boolean }>> = ({ children, noColor, ...props }) => {
  const theme = useMantineTheme();
  return (
    <RRLink
      {...props}
      className={`link-component ${noColor ? "no-color" : "color"}`}
      style={{
        color: noColor ? "unset" : theme.colors["dg-sec"][3],
      }}
    >
      {children}
    </RRLink>
  );
};
