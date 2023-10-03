import { Button, Menu, Text } from "@mantine/core";
import { openConfirmModal, openModal } from "@mantine/modals";
import { PencilIcon, TrashIcon } from "@primer/octicons-react";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { EditPenaltyModal } from "./modals/EditPenalty";

export const UnbanActionMenu = (props: { penalty: CfxState.Penalty }) => {
  const openEditModal = () => {
    openModal({
      title: "Edit ban",
      children: <EditPenaltyModal penalty={props.penalty} />,
    });
  };

  const removeBan = () => {
    openConfirmModal({
      title: "Unban player",
      children: <Text size="sm">Are you sure you want to remove the ban for {props.penalty.steamId}</Text>,
      labels: { confirm: "Confirm", cancel: "Cancel" },
      onConfirm: () => {
        // TODO: refresh ban lijst
        axiosInstance.delete(`/staff/ban/${props.penalty.id}`);
      },
    });
  };

  return (
    <Menu shadow="md" width={150}>
      <Menu.Target>
        <Button>Actions</Button>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item leftSection={<PencilIcon size={14} />} onClick={openEditModal}>
          Edit
        </Menu.Item>
        <Menu.Item color="red" leftSection={<TrashIcon size={14} />} onClick={removeBan}>
          Unban
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
};
