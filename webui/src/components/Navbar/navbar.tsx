import { navbarState } from "@stores/navbar/state";
import { useRecoilValue } from "recoil";
import { NavbarMinimal } from "@components/Navbar/navbarMinimal";
import { devRoute, ExtRouteObject, staffRoute } from "@src/pages/routes";
import "./navbar.scss";
import { NavbarExtension } from "./navbarExtension";
import { useLocation } from "react-router-dom";
import { useEffect, useState } from "react";

const pathWithExtension: { [path: string]: ExtRouteObject } = {
  "/staff": staffRoute,
  "/dev": devRoute,
};

export const Navbar = () => {
  const isOpen = useRecoilValue(navbarState.open);
  const location = useLocation();
  const [basePath, setBasePath] = useState("");
  useEffect(() => {
    setBasePath(location.pathname.match(/^(\/[^/]*)\//)?.[1] ?? "");
  }, [location]);
  // TODO: Maybe move to 1 bar
  return (
    <div className="main-navbar-wrapper">
      <NavbarMinimal canOpen={!!pathWithExtension[basePath]} />
      {isOpen && pathWithExtension[basePath] && <NavbarExtension routes={pathWithExtension[basePath]} />}
    </div>
  );
};
