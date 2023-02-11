import { Button, Menu } from "@mantine/core";
import { AlertIcon, GearIcon } from "@primer/octicons-react";
import { FontAwesomeIcon } from "../../Icon";
import { banUser, kickUser, warnUser } from "./actions";

export const UserActionMenu = ({ steamId }: { steamId: string }) => {
  return (
    <Menu width={200}>
      <Menu.Target>
        <Button leftIcon={<GearIcon size={14} />}>Actions</Button>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Label>Penalise</Menu.Label>
        <Menu.Item onClick={() => warnUser(steamId)} icon={<AlertIcon />}>
          Warn
        </Menu.Item>
        <Menu.Item onClick={() => kickUser(steamId)} icon={<FontAwesomeIcon icon="boot" />} color="orange">
          Kick
        </Menu.Item>
        <Menu.Item onClick={() => banUser(steamId)} icon={<FontAwesomeIcon icon="hammer-war" />} color="red">
          Ban
        </Menu.Item>
        <Menu.Divider />
      </Menu.Dropdown>
    </Menu>
  );
};
