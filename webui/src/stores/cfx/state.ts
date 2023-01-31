import { atom, selector } from "recoil";

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
};
