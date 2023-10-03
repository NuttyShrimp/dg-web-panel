import { Button, Checkbox, Group, NumberInput, TextInput } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { closeAllModals } from "@mantine/modals";
import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import dayjs from "dayjs";
import { useState } from "react";

export const EditPenaltyModal = (props: { penalty: CfxState.Penalty }) => {
  const [reason, setReason] = useState(props.penalty.reason);
  const [points, setPoints] = useState(props.penalty.points);
  const [updating, setUpdating] = useState(false);
  const [length, setLength] = useState<Date | null>(
    new Date(
      dayjs(props.penalty.date)
        .add(props.penalty.length !== -1 ? props.penalty.length : 0, "d")
        .toDate()
    )
  );
  const [perma, setPerma] = useState(props.penalty.length === -1);

  const onClick = async () => {
    if (!length) {
      showNotification({
        title: "Ban error",
        message: `The unban date cannot empty`,
        color: "red",
      });
      return;
    }
    setUpdating(true);
    await axiosInstance.post(`/staff/ban/${props.penalty.id}`, {
      reason,
      points,
      length: perma ? -1 : Math.round((length.getTime() - Date.now()) / (1000 * 60 * 60 * 24)),
    });
    setUpdating(false);
    showNotification({
      title: "Updated ban penalty",
      message: `Successfully updated the players ban info`,
      color: "green",
    });
    closeAllModals();
  };

  return (
    <>
      <TextInput label={"Reason"} value={reason} onChange={e => setReason(e.currentTarget.value)}></TextInput>
      <NumberInput label={"Points"} value={points} onChange={v => setPoints(Number(v) ?? 0)} />
      <Group gap={5} align="flex-end">
        <DatePickerInput
          label="Ban Length"
          placeholder="Pick a unban date"
          value={length}
          onChange={setLength}
          minDate={new Date()}
        />
        <Checkbox label="Permanent" checked={perma} onChange={e => setPerma(e.currentTarget.checked)} />
      </Group>
      <Button disabled={updating} onClick={onClick}>
        Update
      </Button>
    </>
  );
};
