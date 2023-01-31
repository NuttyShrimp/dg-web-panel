import { atom } from "recoil";

export const navbarState = {
  open: atom({
    key: "navbar-state",
    default: true,
  }),
};
