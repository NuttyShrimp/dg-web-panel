import { Button, NumberInput, TextInput } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { SelectCharacter } from "@src/components/Characters/Select";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useState } from "react";

export const CreateBusinessModal = () => {
  const [name, setName] = useState("");
  const [title, setTitle] = useState("");
  const [type, setType] = useState(0);
  const [owner, setOwner] = useState(1000);

  const createBus = async () => {
    await axiosInstance.post(`/staff/business/new`, {
      name,
      label: title,
      typeId: type,
      owner,
    });
    showNotification({
      title: "Created business",
      message: `Successfully created a business: ${title}(${name})`,
    });
  };

  return (
    <>
      <TextInput label={"name"} value={name} onChange={e => setName(e.currentTarget.value)} />
      <TextInput label={"title"} value={title} onChange={e => setTitle(e.currentTarget.value)} />
      <NumberInput label={"type"} value={type} onChange={e => setType(e ?? 1)} />
      <SelectCharacter cid={String(owner)} onChange={setOwner} />
      <Button onClick={createBus}>Create</Button>
    </>
  );
};
