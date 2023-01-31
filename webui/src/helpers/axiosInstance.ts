import { showNotification } from "@mantine/notifications";
import axios from "axios";

export const getHostname = () => import.meta.env.VITE_BACKEND_OVERWRITE ?? location.hostname;

export const axiosInstance = axios.create({
  baseURL: `${location.protocol}//${getHostname()}/api`,
  withCredentials: true,
});

axiosInstance.interceptors.response.use(
  res => res,
  err => {
    if (err) {
      showNotification({
        title: err.response?.data?.title ?? "Unexpected server error",
        message: err.response?.data?.message ?? "Something bad happened on our ends. The devs have been notified",
        color: "red",
      });
    }
    return Promise.reject(err);
  }
);
