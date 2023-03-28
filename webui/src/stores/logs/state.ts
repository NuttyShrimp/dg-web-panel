import { atom } from "recoil";

export const logState = {
  totalCfxLogs: atom<number>({
    key: "logs-cfx-count",
    default: 0,
  }),
  totalPanelLogs: atom<number>({
    key: "logs-panel-count",
    default: 0,
  }),
};
