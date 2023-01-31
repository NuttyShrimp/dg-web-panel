import { axiosInstance } from "@src/helpers/axiosInstance";
import { useSetRecoilState } from "recoil";
import { cfxState } from "./state";

export const useCfxActions = () => {
  const setPlayers = useSetRecoilState(cfxState.players);

  const loadPlayers = async () => {
    try {
      const res = await axiosInstance.get<CfxState.Player[]>("/staff/info/players");
      if (res.status !== 200) return;
      setPlayers(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  const validateCid = async (cid: number) => {
    try {
      const res = await axiosInstance.get<{}>(`/character/${cid}`);
      return res.status === 200;
    } catch (e) {
      console.error(e);
      return false;
    }
  };

  return {
    loadPlayers,
    validateCid,
  };
};
