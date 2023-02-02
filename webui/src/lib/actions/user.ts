import { axiosInstance } from "@src/helpers/axiosInstance";

export const getPlayerBanStatus = async (steamId: string | undefined) => {
  const res = await axiosInstance.get<{ until: string | null }>(`/staff/player/${steamId}/penalties`);
  return res.data;
};
