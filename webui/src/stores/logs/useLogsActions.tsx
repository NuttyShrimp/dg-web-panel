import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useSetRecoilState } from "recoil";
import { logState } from "./state";

export const useLogsActions = () => {
  const setTotalPanelLogs = useSetRecoilState(logState.totalPanelLogs);
  const setTotalCfxLogs = useSetRecoilState(logState.totalCfxLogs);

  const fetchPanelLogs = async (page = 0) => {
    try {
      const resp = await axiosInstance.get<{ logs: Logs.Log[]; total: number }>(`/dev/logs?page=${page}`);
      if (resp.status >= 400) {
        showNotification({
          title: "Failed to fetch panel logs",
          message: `Seems like we cannot fetch the panel logs atm code: ${resp.status} ${resp.statusText}`,
        });
        return [];
      }
      setTotalPanelLogs(resp.data.total);
      return resp.data.logs ?? [];
    } catch (e) {
      console.error(e);
      return [];
    }
  };

  const fetchCfxLogs = async (page = 0, query: string) => {
    try {
      const URLParams = new URLSearchParams({
        page: String(page),
        query,
      });
      const resp = await axiosInstance.get<{ logs: Logs.Log[]; total: number }>(`/staff/logs?${URLParams.toString()}`);
      if (resp.status >= 400) {
        showNotification({
          title: "Failed to fetch panel logs",
          message: `Seems like we cannot fetch the panel logs atm code: ${resp.status} ${resp.statusText}`,
        });
        return [];
      }
      setTotalCfxLogs(resp.data.total);
      return resp.data.logs ?? [];
    } catch (e) {
      console.error(e);
      return [];
    }
  };

  return {
    fetchPanelLogs,
    fetchCfxLogs,
  };
};
