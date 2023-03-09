import { Button, Checkbox, Group, NumberInput, TextInput } from "@mantine/core";
import { DatePicker } from "@mantine/dates";
import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { displayUnixDate } from "@src/helpers/time";
import { useState } from "react";

export const WarnUserModal = ({ steamId }: { steamId: string }) => {
  const [reason, setReason] = useState("");
  const [points, setPoints] = useState(0);

  const doWarn = async () => {
    await axiosInstance.post(`/staff/player/${steamId}/warn`, {
      reason,
      points,
    });
    showNotification({
      title: "Warned player",
      message: `Successfully warned the player with steamid: ${steamId}`,
    });
  };

  return (
    <>
      <TextInput label={"reason"} value={reason} onChange={e => setReason(e.currentTarget.value)} />
      <NumberInput label={"Points (optional)"} value={points} onChange={v => setPoints(v ?? 0)} />
      <Button onClick={doWarn} mt="xs">
        Warn
      </Button>
    </>
  );
};

export const KickUserModal = ({ steamId }: { steamId: string }) => {
  const [reason, setReason] = useState("");
  const [points, setPoints] = useState(0);

  const doKick = async () => {
    await axiosInstance.post(`/staff/player/${steamId}/kick`, {
      reason,
      points,
    });
    showNotification({
      title: "Kicked player",
      message: `Successfully kicked the player with steamid: ${steamId}`,
    });
  };

  return (
    <>
      <TextInput label={"reason"} value={reason} onChange={e => setReason(e.currentTarget.value)} />
      <NumberInput label={"Points (optional)"} value={points} onChange={v => setPoints(v ?? 0)} />
      <Button onClick={doKick} mt="xs">
        Kick
      </Button>
    </>
  );
};

export const BanUserModal = ({ steamId }: { steamId: string }) => {
  const [reason, setReason] = useState("");
  const [points, setPoints] = useState(0);
  const [length, setLength] = useState<Date | null>(new Date());
  const [perma, setPerma] = useState(false);

  const doBan = async () => {
    if (!length) {
      showNotification({
        title: "Ban error",
        message: `The unban date cannot empty`,
        color: "red",
      });
      return;
    }
    await axiosInstance.post(`/staff/player/${steamId}/ban`, {
      target: steamId,
      points,
      reason,
      length: perma ? -1 : Math.round((length.getTime() - Date.now()) / (1000 * 60 * 60 * 24)),
    });
    showNotification({
      title: "Banned player",
      message: `Successfully banned the player with steamid: ${steamId} for ${
        length.getTime() < Date.now() ? "permanent" : displayUnixDate(length.getTime() / 1000)
      }`,
      color: "green",
    });
  };

  return (
    <>
      <TextInput label={"reason"} value={reason} onChange={e => setReason(e.currentTarget.value)} />
      <NumberInput label={"Points (optional)"} value={points} onChange={v => setPoints(v ?? 0)} />
      <Group spacing={5} align="flex-end">
        <DatePicker
          label="Ban Length"
          placeholder="Pick a unban date"
          value={length}
          onChange={setLength}
          minDate={new Date()}
        />
        <Checkbox label="Permanent" checked={perma} onChange={e => setPerma(e.currentTarget.checked)} />
      </Group>
      <Button onClick={doBan} mt="xs">
        Ban
      </Button>
    </>
  );
};
