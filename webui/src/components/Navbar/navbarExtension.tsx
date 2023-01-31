import { createStyles, Navbar, Title } from "@mantine/core";
import { ExtRouteObject } from "@src/pages/routes";
import { FC, useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";

const useStyles = createStyles(theme => ({
  wrapper: {
    backgroundColor: theme.colors.dark[6],
    marginLeft: 70,
    zIndex: 10,
  },
  title: {
    marginBottom: theme.spacing.lg,
    backgroundColor: theme.colors.dark[7],
    padding: theme.spacing.lg,
  },
  entry: {
    display: "flex",
    alignItems: "center",
    height: theme.spacing.lg * 2 + theme.spacing.xs / 2,
    color: theme.colors.dark[0],
    paddingLeft: theme.spacing.md,
    marginRight: theme.spacing.md,
    borderRadius: `0 ${theme.spacing.md}px ${theme.spacing.md}px 0`,
    "&:hover": {
      cursor: "pointer",
      backgroundColor: theme.colors.dark[5],
    },
  },
  activeEntry: {
    color: "white",
    backgroundColor: theme.colors["dg-prim"][4],
    "&:hover": {
      backgroundColor: theme.colors["dg-prim"][4],
    },
  },
}));

const NavbarEntry: FC<{
  route: ExtRouteObject;
  base: string;
}> = ({ route, base }) => {
  const { classes, cx } = useStyles();
  const location = useLocation();
  const navigate = useNavigate();
  const isActive = useMemo(() => {
    return route.path && location.pathname.startsWith(`/${base.replaceAll(/\//g, "")}/${route.path}`);
  }, [location, base, route]);
  const goToRoute = () => {
    navigate(`${base}/${route.path}`);
  };
  return (
    <div className={cx(classes.entry, { [classes.activeEntry]: isActive })} onClick={goToRoute}>
      <Title order={6}>{route.title ?? route.path?.replaceAll(/\//, "") ?? "Wrong configured route"}</Title>
    </div>
  );
};

export const NavbarExtension: FC<{
  routes: ExtRouteObject;
}> = ({ routes }) => {
  const { classes } = useStyles();

  const filteredRoutes = useMemo(() => {
    return routes.children?.filter(r => !r.index) ?? [];
  }, [routes]);

  if (!routes?.path) return null;

  return (
    <div className="main-navbar-extension-wrapper">
      <Navbar width={{ base: 200 }} className={classes.wrapper}>
        <Navbar.Section>
          <Title order={4} className={classes.title}>
            {routes?.title ?? routes?.path?.replaceAll(/\//g, "") ?? ""}
          </Title>
        </Navbar.Section>
        <Navbar.Section grow>
          {filteredRoutes?.map((r, i) => (
            <NavbarEntry key={r?.path ?? `bad-route-${i}`} route={r} base={routes?.path ?? ""} />
          ))}
        </Navbar.Section>
      </Navbar>
    </div>
  );
};
