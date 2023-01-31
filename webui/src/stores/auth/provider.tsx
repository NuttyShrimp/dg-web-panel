import { FC, PropsWithChildren, useEffect } from "react";
import { useAuthActions } from "@stores/auth/useAuthActions";

export const AuthProvider: FC<PropsWithChildren<{}>> = ({ children }) => {
  const { refreshSession } = useAuthActions();
  useEffect(() => {
    refreshSession();
  }, [refreshSession]);
  return <>{children}</>;
};
