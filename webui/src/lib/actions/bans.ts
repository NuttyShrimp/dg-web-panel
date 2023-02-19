import { axiosInstance } from "@src/helpers/axiosInstance"

export const fetchBanList = async () => {
  const resp = await axiosInstance.get("/staff/ban/list")
  return resp.data;
}
