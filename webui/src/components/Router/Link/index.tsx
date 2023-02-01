import { createStyles } from "@mantine/core";
import { FC, PropsWithChildren } from "react";
import { Link as RRLink, LinkProps } from "react-router-dom";

const useStyles = createStyles((theme, { noColor }: { noColor?: boolean }) => ({
  link: {
    textDecoration: "none",
    color: noColor ? "unset" : theme.colors["dg-sec"][3],
    "&:hover": {
      textDecorationLine: noColor ? "unset" : "underline",
    },
  },
}));

export const Link: FC<PropsWithChildren<LinkProps & { noColor?: boolean }>> = ({ children, noColor, ...props }) => {
  const { classes } = useStyles({ noColor });
  return (
    <RRLink {...props} className={classes.link}>
      {children}
    </RRLink>
  );
};
