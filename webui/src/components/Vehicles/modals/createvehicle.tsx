import { Button, TextInput } from "@mantine/core";
import { SelectCharacter } from "@src/components/Characters/Select";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useState } from "react";

export const CreateVehiclemodal = () => {
  const [name, setName] = useState("");
  const [owner, setOwner] = useState(1000);

  const createBus = async () => {
    await axiosInstance.post(`/character/vehicles/give`, {
      model: name,
      owner,
    });
  };

  return (
    <>
      <TextInput label={"name"} value={name} onChange={e => setName(e.currentTarget.value)} />
      <SelectCharacter cid={String(owner)} onChange={setOwner} />
      <Button onClick={createBus}>Create</Button>
    </>
  );
};
