import { axiosInstance } from "@src/helpers/axiosInstance";

export const getPlayerBanStatus = async (steamId: string | undefined) => {
  const res = await axiosInstance.get<{ until: string | null }>(`/staff/player/${steamId}/penalties`);
  return res.data;
};

export const getUserCharacters = async (steamId: string) => {
  const res = await axiosInstance.get<CfxState.Character[]>(`/character/all/${steamId}`);
  return res.data;
};

export const getUserActiveCid = async (steamId: string) => {
  const res = await axiosInstance.get<{ cid: number }>(`/staff/player/${steamId}/active`);
  return res.data.cid;
};
