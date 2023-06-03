import { Button, NumberInput } from "@mantine/core";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useState } from "react";

export const ChangeBalanceModal = (props: { balance?: number; accountId: string }) => {
  const [balance, setBalance] = useState(props.balance ?? 0);
  const [updating, setUpdating] = useState(false);

  const updateBalance = async () => {
    if (updating) return;
    setUpdating(true);
    try {
      await axiosInstance.patch(`/character/bank/${props.accountId}/balance`, {
        balance,
      });
    } catch (e) {
      console.error(e);
    } finally {
      setUpdating(false);
    }
  };

  return (
    <>
      <NumberInput
        min={0}
        decimalSeparator="."
        precision={2}
        value={balance}
        onChange={val => setBalance(val ?? 0)}
        disabled={updating}
      />
      <Button mt={"xs"} onClick={updateBalance} loading={updating}>
        Submit
      </Button>
    </>
  );
};
