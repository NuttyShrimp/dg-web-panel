import { Button, Group, Modal, Stack, TextInput } from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import { showNotification } from "@mantine/notifications";
import { CommentIcon } from "@primer/octicons-react";
import { axiosInstance } from "@src/helpers/axiosInstance";
import dayjs from "dayjs";
import { FC, useCallback, useState } from "react";

const CreateAPIKeyModal: FC<{ open: boolean; onClose: () => void }> = ({ open, onClose }) => {
  const [comment, setComment] = useState("");
  const [expiryDate, setExpiryDate] = useState<Date | null>(new Date());

  const createKey = async () => {
    try {
      const res = await axiosInstance.post("/auth/apikey", {
        comment,
        duration: dayjs(expiryDate).diff(dayjs(new Date()), "minute"),
      });
      if (res.status !== 200) {
        showNotification({
          title: "Creation error",
          message: "Failed to create new API key",
          color: "red",
        });
        return;
      }
      onClose();
    } catch (e) {
      console.error(e);
    }
  };

  // TODO: Tes this datetimepicker
  return (
    <Modal opened={open} onClose={onClose} title="Create new API key">
      <Stack>
        <TextInput
          leftSection={<CommentIcon />}
          value={comment}
          onChange={e => setComment(e.currentTarget.value)}
          placeholder="comment"
          label="Comment"
        />
        <Group grow>
          <DateTimePicker value={expiryDate} onChange={setExpiryDate} minDate={dayjs(new Date()).toDate()} />
        </Group>
        <Group justify="flex-end">
          <Button onClick={createKey}>Create</Button>
        </Group>
      </Stack>
    </Modal>
  );
};

export const useCreateAPIKeyModal = () => {
  const [opened, setOpened] = useState(false);

  const modal = useCallback(() => <CreateAPIKeyModal open={opened} onClose={() => setOpened(false)} />, [opened]);
  return {
    Modal: modal,
    openModal: () => setOpened(true),
  };
};
