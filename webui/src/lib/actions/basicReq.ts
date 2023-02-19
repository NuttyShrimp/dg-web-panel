import { axiosInstance } from "@src/helpers/axiosInstance";

export const basicGet = async <T>(endpoint: string): Promise<T> => {
  const resp = await axiosInstance.get(endpoint);
  return resp.data;
};
