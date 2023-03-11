import { Button, Checkbox, Group, NumberInput } from "@mantine/core";
import { DatePicker } from "@mantine/dates";
import { showNotification } from "@mantine/notifications";
import { PenaltyReasonSelector } from "@src/components/Inputs/PenaltyReasonSelector";
import { classInfo, reasons as penaltyReasons } from "@src/data/PenaltyReasons";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { displayUnixDate } from "@src/helpers/time";
import dayjs from "dayjs";
import { useState } from "react";

const updateReasons = (newReasons: string[]) => {
  const newData = {
    points: 0,
    length: 0,
  };
  newReasons.forEach(r => {
    if (penaltyReasons[r]) {
      const reasonInfo = classInfo[penaltyReasons[r]];
      if (!classInfo) return;
      newData.points += reasonInfo.points;
      newData.length += reasonInfo.length;
    }
  });
  return newData;
};

export const WarnUserModal = ({ steamId }: { steamId: string }) => {
  const [reasons, setReasons] = useState<string[]>([]);
  const [points, setPoints] = useState(0);

  const doWarn = async () => {
    await axiosInstance.post(`/staff/player/${steamId}/warn`, {
      reason: reasons.join(", "),
      points,
    });
    showNotification({
      title: "Warned player",
      message: `Successfully warned the player with steamid: ${steamId}`,
    });
  };

  return (
    <>
      <PenaltyReasonSelector
        reasons={reasons}
        setReasons={r => {
          const info = updateReasons(r);
          setPoints(info.points);
          setReasons(r);
        }}
      />
      <NumberInput label={"Points (optional)"} value={points} onChange={v => setPoints(v ?? 0)} />
      <Button onClick={doWarn} mt="xs">
        Warn
      </Button>
    </>
  );
};

export const KickUserModal = ({ steamId }: { steamId: string }) => {
  const [reasons, setReasons] = useState<string[]>([]);
  const [points, setPoints] = useState(0);

  const doKick = async () => {
    await axiosInstance.post(`/staff/player/${steamId}/kick`, {
      reason: reasons.join(", "),
      points,
    });
    showNotification({
      title: "Kicked player",
      message: `Successfully kicked the player with steamid: ${steamId}`,
    });
  };

  return (
    <>
      <PenaltyReasonSelector
        reasons={reasons}
        setReasons={r => {
          const info = updateReasons(r);
          setPoints(info.points);
          setReasons(r);
        }}
      />
      <NumberInput label={"Points (optional)"} value={points} onChange={v => setPoints(v ?? 0)} />
      <Button onClick={doKick} mt="xs">
        Kick
      </Button>
    </>
  );
};

export const BanUserModal = ({ steamId }: { steamId: string }) => {
  const [reasons, setReasons] = useState<string[]>([]);
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
      reason: reasons.join(", "),
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
      <PenaltyReasonSelector
        reasons={reasons}
        setReasons={r => {
          const info = updateReasons(r);
          setPoints(info.points);
          setLength(dayjs().add(info.length, "d").toDate());
          setReasons(r);
        }}
      />
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
