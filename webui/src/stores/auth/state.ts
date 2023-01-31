import { atom } from "recoil";

export const authState = {
  userInfo: atom<AuthState.UserInfo | null>({
    key: "auth-info",
    default: null,
  }),
  // See this as a readonly mutex lock
  // if true --> Busy doing logic that could lead to changing data
  // Is unlocked after first refresh
  lock: atom<boolean>({
    key: "auth-info-lock",
    default: true,
  }),
};
