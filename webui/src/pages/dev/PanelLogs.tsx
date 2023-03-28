import { Container, Title } from "@mantine/core";
import { LogList } from "@src/components/Logs/LogList";
import { logState } from "@src/stores/logs/state";
import { useLogsActions } from "@src/stores/logs/useLogsActions";
import { useRecoilValue } from "recoil";

export const PanelLogs = () => {
  const { fetchPanelLogs } = useLogsActions();
  const totalPanelLogs = useRecoilValue(logState.totalPanelLogs);

  return (
    <Container size="lg">
      <Title order={2}>Panel logs</Title>
      <LogList total={totalPanelLogs} id="panel-logs" fetchFunc={fetchPanelLogs} />
    </Container>
  );
};
