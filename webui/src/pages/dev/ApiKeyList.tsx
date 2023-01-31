import { Button, Checkbox, Group, Table } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { PlusIcon, TrashIcon } from "@primer/octicons-react";
import { axiosInstance } from "@src/helpers/axiosInstance";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useCreateAPIKeyModal } from "./components/CreateKeyModal";

export const ApiKeyList = () => {
  const { Modal, openModal } = useCreateAPIKeyModal();
  const [selection, setSelection] = useState<string[]>([]);
  const [keys, setKeys] = useState<Dev.APIKey[]>([]);

  const fetchKeys = async () => {
    try {
      const res = await axiosInstance.get<{ keys: Dev.APIKey[] }>("auth/apikey/all");
      if (res.status !== 200) {
        showNotification({
          title: "Fetch error",
          message: "Failed to get list of API keys",
          color: "red",
        });
      }
      setKeys(res.data.keys);
    } catch (error) {
      console.error(error);
    }
  };

  const deleteKeys = async () => {
    try {
      const res = await axiosInstance.delete("auth/apikey", {
        data: {
          keys: selection,
        },
      });
      if (res.status !== 200) {
        showNotification({
          title: "API error",
          message: "Failed to delete selected keys",
        });
      }
      fetchKeys();
    } catch (e) {
      console.error(e);
    }
  };

  const toggleRow = (id: string) =>
    setSelection(current => (current.includes(id) ? current.filter(item => item !== id) : [...current, id]));

  useEffect(() => {
    fetchKeys();
  }, []);

  return (
    <div>
      <Modal />
      <Group position="right" pb={"sm"}>
        <Button leftIcon={<TrashIcon />} color={"red"} disabled={selection.length === 0} onClick={deleteKeys}>
          Delete
        </Button>
        <Button leftIcon={<PlusIcon />} onClick={openModal}>
          Add
        </Button>
      </Group>
      <Table striped highlightOnHover>
        <thead>
          <tr>
            <th style={{ width: 40 }} />
            <th>API Key</th>
            <th>Comment</th>
            <th>Expires on</th>
            <th>Assigned to</th>
          </tr>
        </thead>
        <tbody>
          {keys.map(key => (
            <tr key={key.key}>
              <td>
                <Checkbox
                  checked={selection.includes(key.key)}
                  onChange={() => toggleRow(key.key)}
                  transitionDuration={0}
                />
              </td>
              <td>{key.key}</td>
              <td>{key.comment}</td>
              <td>{dayjs(key.expiry).format("DD/MM/YYYY HH:mm:ss")}</td>
              <td>
                {key.User.Username}({key.userId})
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};
