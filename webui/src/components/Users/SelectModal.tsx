import { Button, Center } from "@mantine/core";
import { closeAllModals } from "@mantine/modals";
import { useState } from "react";
import { UserSelect } from "./Select";

export const UserSelectModal = ({ onAccept }: { onAccept: (val: string) => void }) => {
  const [user, setUser] = useState<string | null>(null);
  return (
    <>
      <UserSelect onChange={setUser} />
      <Center mt={"xs"}>
        <Button
          onClick={() => {
            if (user === null) return;
            onAccept(user);
            closeAllModals();
          }}
        >
          Accept
        </Button>
      </Center>
    </>
  );
};
