import { IndexRouteObject, Navigate, NonIndexRouteObject, RouteObject } from "react-router-dom";
import { Login } from "./auth/Login";
import { RequireAuth } from "./auth/RequireAuth";
import { E403, E404 } from "./errors/400";
import { E500 } from "./errors/500";
import { ApiKeyList } from "./dev/ApiKeyList";
import { StaffDashboard } from "./staff/Dashboard";
import { StaffReport } from "./staff/Report";
import { StaffReports } from "./staff/ReportList";
import { PageLoader } from "./staff/PageLoader";
import { CacheControl } from "./dev/CacheControl";
import { CharacterList } from "./staff/CharacterList";
import { CharacterPage } from "./staff/Character";
import { UserList } from "./staff/UserList";
import { UserPage } from "./staff/User";
import { PanelLogs } from "./dev/PanelLogs";
import { BusinessList } from "./staff/BusinessList";
import { Business } from "./staff/Business";
import { DevActionPage } from "./dev/Actions";
import { BanListPage } from "./staff/BanList";
import { AdminLogList as AdminLogList } from "./staff/LogList";

type ExtNonIndexRouteObject = Omit<NonIndexRouteObject, "children"> & {
  children?: ExtRouteObject[];
};

export type ExtRouteObject = (ExtNonIndexRouteObject | IndexRouteObject) & {
  title?: string;
};

export const staffRoute: ExtRouteObject = {
  path: "/staff",
  title: "Staff",
  element: <RequireAuth role="staff" />,
  children: [
    {
      index: true,
      element: <Navigate to={"/staff/dashboard"} replace={true} />,
    },
    {
      title: "Dashboard",
      path: "dashboard",
      element: <StaffDashboard />,
    },
    {
      title: "Reports",
      path: "reports",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <StaffReports />,
        },
        {
          path: ":id",
          title: "Report",
          element: <StaffReport />,
        },
      ],
    },
    {
      title: "Characters",
      path: "characters",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <CharacterList />,
        },
        {
          path: ":cid",
          title: "Character",
          element: <CharacterPage />,
        },
      ],
    },
    {
      title: "Users",
      path: "users",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <UserList />,
        },
        {
          path: ":steamid",
          title: "Player",
          element: <UserPage />,
        },
      ],
    },
    {
      path: "business",
      title: "Businesses",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <BusinessList />,
        },
        {
          path: ":id",
          element: <Business />,
        },
      ],
    },
    {
      path: "banlist",
      title: "Ban list",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <BanListPage />,
        },
      ],
    },
    {
      path: "logs",
      title: "Logs",
      element: <PageLoader />,
      children: [
        {
          index: true,
          element: <AdminLogList />,
        },
      ],
    },
  ],
};

export const devRoute: ExtRouteObject = {
  path: "/dev",
  title: "Devs",
  element: <RequireAuth role="developer" />,
  children: [
    {
      index: true,
      element: <Navigate to={"/dev/apikeys"} replace />,
    },
    {
      path: "apikeys",
      title: "API Keys",
      element: <ApiKeyList />,
    },
    {
      path: "cache",
      title: "Cache control",
      element: <CacheControl />,
    },
    {
      path: "panellogs",
      title: "Panel logs",
      element: <PanelLogs />,
    },
    {
      path: "actions",
      title: "Actions",
      element: <DevActionPage />,
    },
  ],
};

export const routes: RouteObject[] = [
  {
    path: "/",
    element: <div>Home</div>,
  },
  {
    path: "*",
    element: <Navigate to={"/errors/404"} replace />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/errors",
    children: [
      {
        index: true,
        element: <Navigate to={"/errors/404"} />,
      },
      {
        path: "403",
        element: <E403 />,
      },
      {
        path: "404",
        element: <E404 />,
      },
      {
        path: "500",
        element: <E500 />,
      },
    ],
  },
  staffRoute,
  devRoute,
];
