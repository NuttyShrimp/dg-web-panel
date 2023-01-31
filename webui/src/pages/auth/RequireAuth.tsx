import { FC, useEffect, useState } from "react";
import { useRecoilValue } from "recoil";
import { authState } from "@stores/auth/state";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { LoadingOverlay } from "@mantine/core";
import { axiosInstance } from "@src/helpers/axiosInstance";

export const RequireAuth: FC<{ role: string }> = ({ role }) => {
  const UserInfo = useRecoilValue(authState.userInfo);
  const infoLock = useRecoilValue(authState.lock);
  const location = useLocation();
  const [loading, setLoading] = useState(true);
  const [allowed, setAllowed] = useState(false);

  useEffect(() => {
    let ignore = false;

    const fetchAccess = async () => {
      setLoading(true);
      try {
        const res = await axiosInstance.get("/auth/role", {
          params: {
            role,
          },
        });
        if (!ignore) {
          setAllowed(res.data?.access ?? false);
        }
      } catch (e) {
        console.error(e);
        setAllowed(false);
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchAccess();

    return () => {
      ignore = true;
    };
  }, [role]);

  if (!UserInfo && !infoLock) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (!allowed && !loading) {
    return <Navigate to="/errors/404" replace />;
  }

  return (
    <>
      <LoadingOverlay visible={infoLock || loading} overlayBlur={7} />
      <Outlet />
    </>
  );
};
