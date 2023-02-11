import { axiosInstance } from "@src/helpers/axiosInstance";

export const getPlayerBanStatus = async (steamId: string | undefined) => {
  const res = await axiosInstance.get<{ until: string | null }>(`/staff/player/${steamId}/banstatus`);
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

export const getUserPenalties = async (steanId: string) => {
  const res = await axiosInstance.get<CfxState.Penalty[] | null>(`/staff/player/${steanId}/penalties`);
  return res.data;
};
