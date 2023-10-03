import { navbarState } from "@stores/navbar/state";
import { useRecoilState } from "recoil";
import { NavbarMinimal } from "@components/Navbar/navbarMinimal";
import { devRoute, ExtRouteObject, staffRoute } from "@src/pages/routes";
import "./navbar.scss";
import { NavbarExtension } from "./navbarExtension";
import { useLocation } from "react-router-dom";
import { useEffect, useMemo, useState } from "react";

const pathWithExtension: { [path: string]: ExtRouteObject } = {
  "/staff": staffRoute,
  "/dev": devRoute,
};

export const Navbar = () => {
  const location = useLocation();
  const [basePath, setBasePath] = useState("");
  const [isOpen, setIsOpen] = useRecoilState(navbarState.open);

  const extRoutes = useMemo(() => pathWithExtension[basePath], [basePath]);

  useEffect(() => {
    setIsOpen(!!extRoutes);
  }, [extRoutes, setIsOpen]);

  useEffect(() => {
    setBasePath(location.pathname.match(/^(\/[^/]*)\//)?.[1] ?? "");
  }, [location]);
  // TODO: Maybe move to 1 bar
  return (
    <div className="main-navbar-wrapper">
      <NavbarMinimal canOpen={!!extRoutes} />
      {isOpen && extRoutes && <NavbarExtension routes={extRoutes} />}
    </div>
  );
};
