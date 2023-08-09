import { axiosInstance } from "@src/helpers/axiosInstance";
import { useRecoilState, useSetRecoilState } from "recoil";
import { characterState } from "./state";
import { useCallback } from "react";

export const useCharacterActions = () => {
  const [charList, setListStore] = useRecoilState(characterState.list);
  const setBankStore = useSetRecoilState(characterState.bank);
  const setVehicleStore = useSetRecoilState(characterState.vehicles);

  const resetStores = () => {
    setListStore([]);
    setVehicleStore(null);
    setBankStore(null);
  };

  const fetchCharacters = async () => {
    try {
      const res = await axiosInstance.get<CfxState.Character[]>(`/character/all`);
      if (res.status !== 200) return;
      setListStore(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  const fetchActiveCharacters = useCallback(async () => {
    try {
      const res = await axiosInstance.get<{ cid: number; serverId: number }[]>(`/character/active`);
      if (res.status !== 200) return [];
      return charList
        .filter(c => res.data.find(i => i.cid === c.citizenid))
        .map(c => {
          c.serverId = res.data.find(i => i.cid === c.citizenid)?.serverId ?? 0;
          return c;
        });
    } catch (e) {
      console.error(e);
      return [];
    }
  }, [charList]);

  const fetchCharReputation = async (cid: number) => {
    try {
      const res = await axiosInstance.get<Record<string, number>>(`/character/${cid}/reputation`);
      if (res.status !== 200) return;
      if (res.data.citizenid !== undefined) {
        delete res.data.citizenid;
      }
      return res.data;
    } catch (e) {
      console.error(e);
    }
  };

  return {
    resetStores,
    fetchCharReputation,
    fetchActiveCharacters,
    fetchCharacters,
  };
};
