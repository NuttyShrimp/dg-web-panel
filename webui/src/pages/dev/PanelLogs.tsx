import { Center, Container, Loader, Stack, Text, Title } from "@mantine/core";
import { LogList } from "@src/components/Logs/LogList";
import { logState } from "@src/stores/logs/state";
import { useLogsActions } from "@src/stores/logs/useLogsActions";
import { useEffect } from "react";
import { useRecoilValue } from "recoil";

export const PanelLogs = () => {
  const { clearPanelLogs, loadingLogs, fetchPanelLogs } = useLogsActions();
  const panelLogs = useRecoilValue(logState.panelLogs);

  useEffect(() => {
    clearPanelLogs();
    fetchPanelLogs();
  }, []);

  if (loadingLogs) {
    return (
      <Container>
        <Center>
          <Stack>
            <Center>
              <Loader />
            </Center>
            <Text>Loading logs...</Text>
          </Stack>
        </Center>
      </Container>
    );
  }

  return (
    <Container size="lg">
      <Title order={2}>Panel logs</Title>
      <LogList logs={panelLogs} onLoadMore={() => fetchPanelLogs(panelLogs.length)} />
    </Container>
  );
};
