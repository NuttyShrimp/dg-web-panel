import { Container, Flex, Title } from "@mantine/core";
import { LogList } from "@src/components/Logs/LogList";
import { QueryMenu } from "@src/components/Logs/QueryMenu";
import { logState } from "@src/stores/logs/state";
import { useLogsActions } from "@src/stores/logs/useLogsActions";
import { useState } from "react";
import { useRecoilValue } from "recoil";

export const AdminLogList = () => {
  const { fetchCfxLogs } = useLogsActions();
  const logCount = useRecoilValue(logState.totalCfxLogs);
  const [query, setQuery] = useState("");

  return (
    <Container size={"lg"}>
      <Flex justify={"space-between"} mb={"xs"}>
        <Title order={2}>Server logs</Title>
        <QueryMenu setQuery={setQuery} />
      </Flex>
      <LogList query={query} total={logCount} id="cfx-logs" fetchFunc={fetchCfxLogs} />
    </Container>
  );
};
