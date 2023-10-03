import { Button, Menu } from "@mantine/core";
import { closeAllModals, openModal } from "@mantine/modals";
import { GearIcon, TrashIcon } from "@primer/octicons-react";
import { useCfxBusiness } from "@src/stores/cfx/hooks/useCfxBusiness";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SelectCharacterModal } from "../Characters/Select";

export const BusinessActionMenu = ({ id }: { id: number }) => {
  const { deleteBusiness, changeOwner } = useCfxBusiness();
  const navigate = useNavigate();
  const [busy, setBusy] = useState(false);

  const onDelClick = async () => {
    setBusy(true);
    await deleteBusiness(id);
    setBusy(false);
    navigate("/staff/business");
  };

  const onChangeOwnerClick = () => {
    setBusy(true);
    openModal({
      title: "Change business owner",
      children: (
        <SelectCharacterModal
          onAccept={async cid => {
            if (cid == 0) return;
            await changeOwner(id, cid);
            closeAllModals();
            setBusy(false);
          }}
        />
      ),
      onClose: () => setBusy(false),
    });
    setBusy(false);
  };

  return (
    <Menu shadow={"md"} width={150} disabled={busy}>
      <Menu.Target>
        <Button leftSection={<GearIcon />}>Actions</Button>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Item onClick={onChangeOwnerClick}>Change owner (Existing employees only)</Menu.Item>
        <Menu.Item color="red" leftSection={<TrashIcon size={14} />} onClick={onDelClick}>
          Delete
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
};
