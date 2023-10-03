import { Title } from "@mantine/core";
import { ExtRouteObject } from "@src/pages/routes";
import { FC, useMemo } from "react";
import { Link, useLocation } from "react-router-dom";

const NavbarEntry: FC<{
  route: ExtRouteObject;
  base: string;
}> = ({ route, base }) => {
  const location = useLocation();
  const isActive = useMemo(() => {
    return route.path && location.pathname.startsWith(`/${base.replaceAll(/\//g, "")}/${route.path}`);
  }, [location, base, route]);
  return (
    <Link className={`link ${isActive ? "active" : ""}`} to={`${base}/${route.path}`}>
      <Title order={6}>{route.title ?? route.path?.replaceAll(/\//, "") ?? "Wrong configured route"}</Title>
    </Link>
  );
};

export const NavbarExtension: FC<{
  routes: ExtRouteObject;
}> = ({ routes }) => {
  const filteredRoutes = useMemo(() => {
    return routes.children?.filter(r => !r.index) ?? [];
  }, [routes]);

  if (!routes?.path) return null;

  return (
    <div className={"navbar-ext"}>
      <div>
        <Title order={4} className={"title"}>
          {routes?.title ?? routes?.path?.replaceAll(/\//g, "") ?? ""}
        </Title>
      </div>
      <div>
        {filteredRoutes?.map((r, i) => (
          <NavbarEntry key={r?.path ?? `bad-route-${i}`} route={r} base={routes?.path ?? ""} />
        ))}
      </div>
    </div>
  );
};
