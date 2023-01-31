import { atom } from "recoil";

export const logState = {
  panelLogs: atom<Logs.Log[]>({
    key: "logs-panel",
    default: [],
  }),
};
