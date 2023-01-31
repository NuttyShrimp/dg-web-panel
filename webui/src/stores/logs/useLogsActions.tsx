import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useState } from "react";
import { useSetRecoilState } from "recoil";
import { logState } from "./state";

export const useLogsActions = () => {
  const setPanelLogs = useSetRecoilState(logState.panelLogs);
  const [loading, setLoading] = useState(false);

  const fetchPanelLogs = async (offset = 0) => {
    setLoading(true);
    try {
      const resp = await axiosInstance.get<Logs.Log[]>(`/dev/logs?offset=${offset}`);
      if (resp.status >= 400) {
        showNotification({
          title: "Failed to fetch panel logs",
          message: `Seems like we cannot fetch the panel logs atm code: ${resp.status} ${resp.statusText}`,
        });
        return;
      }
      setPanelLogs(resp.data);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const clearPanelLogs = () => setPanelLogs([]);

  return {
    fetchPanelLogs,
    clearPanelLogs,
    loadingLogs: loading,
  };
};
