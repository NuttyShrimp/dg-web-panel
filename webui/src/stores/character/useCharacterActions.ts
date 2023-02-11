import { axiosInstance } from "@src/helpers/axiosInstance";
import { useRecoilState, useSetRecoilState } from "recoil";
import { characterState } from "./state";

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

  const fetchActiveCharacters = async () => {
    try {
      await fetchCharacters();
      const res = await axiosInstance.get<number[]>(`/character/active`);
      if (res.status !== 200) return [];
      return charList.filter(c => res.data.includes(c.citizenid));
    } catch (e) {
      console.error(e);
      return [];
    }
  };

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
