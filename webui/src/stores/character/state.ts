import { atom, selector } from "recoil";

export const characterState = {
  list: atom<CfxState.Character[]>({
    key: "character-list",
    default: [],
  }),
  cid: atom<number>({
    key: "charcter-cid",
    default: 0,
  }),
  selected: selector({
    key: "character-selected",
    get: ({ get }) => {
      const characters: CfxState.Character[] = get(characterState.list);
      const selectedCid: number = get(characterState.cid);
      return characters.find(c => c.citizenid === selectedCid);
    },
  }),
  bank: atom<CharacterState.Bank[] | null>({
    key: "character-bank",
    default: null,
  }),
  vehicles: atom<CharacterState.Vehicle[] | null>({
    key: "character-vehicles",
    default: null,
  }),
  reputation: atom<Record<string, number> | undefined>({
    key: "character-reputation",
    default: undefined,
  }),
};
