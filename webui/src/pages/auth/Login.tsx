import { Centerbox } from "@components/CenterBox/centerbox";
import { flushSync } from "react-dom";
import { useEffect, useState } from "react";
import { Button, Text } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useLocation, useNavigate } from "react-router-dom";
import { authState } from "@src/stores/auth/state";
import { useRecoilValue } from "recoil";

export const Login = () => {
  const [canLogin, setCanLogin] = useState(true);
  const userInfo = useRecoilValue(authState.userInfo);
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (userInfo) {
      navigate((location.state as any)?.from?.pathname ?? "/", {
        replace: true,
      });
    }
  }, [userInfo, location, navigate]);

  const handleDiscordLogin = async () => {
    if (!canLogin) return;
    flushSync(() => {
      setCanLogin(false);
    });
    try {
      const result = await axiosInstance.post<{ url?: string }>(`/auth/login?type=discord`);
      if (result.status !== 200) {
        throw new Error("Failed to login");
      }
      if (result.data?.url) {
        window.location.replace(result.data.url);
      }
    } catch (e) {
      showNotification({
        title: "Error while logging in",
        message: "Something failed while trying to authenticate via discord",
        color: "red",
      });
    } finally {
      flushSync(() => {
        setCanLogin(true);
      });
    }
  };

  return (
    <Centerbox title={"Login"}>
      <>
        <Text size={"md"}>To access the panel, you&apos;ll need to login via discord</Text>
        <Button
          onClick={handleDiscordLogin}
          sx={{
            marginTop: "5vh",
          }}
          loading={!canLogin}
          leftIcon={<i className="fa-brands fa-discord"></i>}
        >
          {" "}
          Login via discord
        </Button>
      </>
    </Centerbox>
  );
};
