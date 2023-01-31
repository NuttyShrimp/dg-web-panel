import { useCallback } from "react";
import { useRecoilCallback, useSetRecoilState } from "recoil";
import { showNotification } from "@mantine/notifications";
import { authState } from "@stores/auth/state";
import { axiosInstance } from "@src/helpers/axiosInstance";

export const useAuthActions = () => {
  const setUserInfo = useSetRecoilState(authState.userInfo);
  const setLock = useSetRecoilState(authState.lock);

  const getUserInfo = useRecoilCallback(() => async () => {
    try {
      const res = await axiosInstance.get<AuthState.UserInfo>(`/user/me`);
      if (res.status !== 200) {
        showNotification({
          onClose: () => console.log("yeet login"),
          title: "Error getting user info",
          message:
            "We could not get info about u, please try logging in again (dismissing this message will log you out)",
          color: "red",
        });
        return;
      }
      if (res.data?.error) {
        console.error(res.data?.error);
        return;
      }
      setUserInfo(res.data);
    } catch (e) {
      console.error(e);
    }
  });

  const refreshSession = useCallback(async () => {
    try {
      () => setLock(true);
      const res = await axiosInstance.post<{ isExpired?: boolean; error?: string }>(`/auth/refresh`);
      if (res.status !== 200) {
        setUserInfo(null);
        return;
      }
      if (res.data?.error) {
        console.error(res.data?.error);
        setUserInfo(null);
        return;
      }
      if (res.data?.isExpired) {
        setUserInfo(null);
        return;
      }
      await getUserInfo();
    } catch (e) {
      console.error(e);
      setUserInfo(null);
    } finally {
      setLock(false);
    }
  }, [getUserInfo, setUserInfo, setLock]);

  const logoutUser = useRecoilCallback(() => async () => {
    try {
      const res = await axiosInstance.post("/auth/logout");
      if (res.status !== 200) {
        console.error(`Logout request failed: ${res.data?.error ?? res.statusText}(${res.status})`);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setUserInfo(null);
    }
  });

  return {
    logoutUser,
    refreshSession,
    getUserInfo,
  };
};
