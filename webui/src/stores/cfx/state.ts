import { axiosInstance } from "@src/helpers/axiosInstance";
import { atom, selector, selectorFamily } from "recoil";

export const cfxState = {
  player: atom<string | null>({
    key: "cfx-selected-player",
    default: null,
  }),
  players: atom<CfxState.Player[]>({
    key: "cfx-players",
    default: [],
  }),
  selectPlayers: selector({
    key: "cfx-selectable-players",
    get: ({ get }) => {
      const players: CfxState.Player[] = get(cfxState.players) ?? [];
      return players.map(p => ({
        label: p.name,
        value: p.steamId,
      }));
    },
  }),
  businesses: atom<CfxState.Business.Entry[]>({
    default: [],
    key: "cfx-business-list",
  }),
  businessLogTotal: selectorFamily({
    key: "cfx-business-log-count",
    get: (id: number) => async () => {
      try {
        const res = await axiosInstance.get<{ total: number }>(`/staff/business/${id}/logcount`);
        return res.data.total;
      } catch (e) {
        console.error(e);
        return 0;
      }
    },
  }),
  businessEmployees: selectorFamily({
    key: "cfx-business-employees",
    get: (id: number) => async () => {
      try {
        const res = await axiosInstance.get<CfxState.Business.Employee[]>(`/staff/business/${id}/employees`);
        return res.data;
      } catch (e) {
        console.error(e);
        return [];
      }
    },
  }),
};
