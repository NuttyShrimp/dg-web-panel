import { openModal } from "@mantine/modals";
import { BanUserModal, KickUserModal, WarnUserModal } from "./modals";

export const warnUser = async (steamId: string) => {
  openModal({
    title: "Warn player",
    children: <WarnUserModal steamId={steamId} />,
  });
};

export const kickUser = async (steamId: string) => {
  openModal({
    title: "Kick player",
    children: <KickUserModal steamId={steamId} />,
  });
};

export const banUser = async (steamId: string) => {
  openModal({
    title: "Ban player",
    children: <BanUserModal steamId={steamId} />,
  });
};
