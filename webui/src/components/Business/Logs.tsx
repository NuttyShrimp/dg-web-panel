import { FC } from "react";
import { Center, Loader, Pagination, Stack, Text } from "@mantine/core";
import { useCfxBusiness } from "@src/stores/cfx/hooks/useCfxBusiness";
import { cfxState } from "@src/stores/cfx/state";
import { useEffect, useState } from "react";
import { useRecoilRefresher_UNSTABLE, useRecoilValue } from "recoil";
import { List } from "../List";
import { useQuery } from "@tanstack/react-query";

export const BusinessLogs: FC<{ id: number }> = ({ id }) => {
  const [page, setPage] = useState(1);
  const totalLogs = useRecoilValue(cfxState.businessLogTotal(id));
  const { fetchLogs } = useCfxBusiness();
  const refreshLogCount = useRecoilRefresher_UNSTABLE(cfxState.businessLogTotal(id));

  const {
    isLoading,
    isError,
    error,
    data: logs,
    isFetching,
  } = useQuery<CfxState.Business.Log[], Error>({
    queryKey: ["business-logs", page],
    queryFn: () => fetchLogs(id, page - 1),
    keepPreviousData: true,
  });

  useEffect(() => {
    if (totalLogs > 0) {
      setPage(1);
    }
  }, [totalLogs]);

  useEffect(() => {
    refreshLogCount();
  }, []);

  if (totalLogs < 1) {
    return (
      <Center>
        <Text>Dit bedrijf heeft momenteel nog geen logs</Text>
      </Center>
    );
  }

  if (isLoading || isFetching) {
    return (
      <div>
        <Center>
          <Stack spacing={"xs"}>
            <Center>
              <Loader />
            </Center>
            <Text>Loading logs..</Text>
          </Stack>
        </Center>
      </div>
    );
  }

  if (isError) {
    return (
      <div>
        <Center>
          <Text>
            <div>Failed to load logs: {error.message}</div>
          </Text>
        </Center>
      </div>
    );
  }

  return (
    <div>
      <List>
        {logs.map(l => (
          <List.Entry key={l.id}>
            <Stack spacing={4}>
              <Text weight="bolder">{l.type}</Text>
              <Text>{l.action}</Text>
            </Stack>
          </List.Entry>
        ))}
      </List>
      <Center mt={"xs"}>
        <Pagination page={page} onChange={setPage} total={Math.ceil(totalLogs / 50)} />
      </Center>
    </div>
  );
};
