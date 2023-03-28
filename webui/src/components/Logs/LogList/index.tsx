import {
  ActionIcon,
  Box,
  Center,
  Code,
  Container,
  Divider,
  Flex,
  Loader,
  Pagination,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { SearchIcon } from "@primer/octicons-react";
import { queryClient } from "@src/helpers/queryClient";
import { displayUnixDate } from "@src/helpers/time";
import { parsePotentialJSON } from "@src/helpers/util";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { flushSync } from "react-dom";
import { List } from "../../List";

import "./style.scss";

interface LogListProps {
  total: number;
  // Should return Log.Response
  fetchFunc: (page: number, query: string) => Promise<Logs.Log[]>;
  id: string;
  query?: string;
}

const LogInfo = ({ log }: { log: Logs.Log }) => (
  <Container size={"sm"}>
    {Object.keys(log)
      .filter(lk => log[lk as keyof Logs.Log] !== "")
      .map(key => (
        <Box key={`${log._id}-${key}`} pt={2}>
          <Text weight={"bolder"} size="xs">
            {key.replace(/^_/, "")}
          </Text>
          <Code className="log-entry-value" block>
            {JSON.stringify(parsePotentialJSON(log[key as keyof Logs.Log]), undefined, 2)}
          </Code>
        </Box>
      ))}
  </Container>
);

export const LogEntryList = ({ logs }: { logs: Logs.Log[] }) => {
  const [focusedLogs, setFocusedLogs] = useState<string[]>([]);
  const toggledFocus = (id: string) => {
    if (focusedLogs.includes(id)) {
      setFocusedLogs(focusedLogs.filter(fId => fId !== id));
    } else {
      setFocusedLogs([...focusedLogs, id]);
    }
  };
  return (
    <>
      {logs.map(l => (
        <List.Entry key={l._id}>
          <Flex direction={"column"} w={"100%"}>
            <Box onClick={() => toggledFocus(l._id)} w={"100%"} style={{ cursor: "pointer" }}>
              <Text weight={"bolder"}>{displayUnixDate(l.timestamp)}</Text>
              <Text>{l.short_message}</Text>
            </Box>
            {focusedLogs.includes(l._id) && (
              <>
                <Divider pb={"sm"} />
                <LogInfo log={l} />
              </>
            )}
          </Flex>
        </List.Entry>
      ))}
    </>
  );
};

export const LogList = (props: LogListProps) => {
  const [page, setPage] = useState(1);
  const query = useRef(props.query ?? "");
  const queryBarRef = useRef<HTMLInputElement | null>(null);
  const {
    data: logs,
    error,
    isError,
    isLoading,
  } = useQuery<Logs.Log[], Error>({
    queryKey: [props.id, page],
    queryFn: () => props.fetchFunc(page - 1, query.current.trim() === "" ? "*" : query.current),
  });

  useEffect(() => {
    if (props.query) {
      query.current = props.query;
      if (queryBarRef.current) {
        queryBarRef.current.value = props.query;
      }
    } else {
      query.current = "";
    }
  }, [props.query]);

  return (
    <Center style={{ flexDirection: "column" }}>
      <Flex w={"50vw"} align="center" mb={"xs"}>
        <ActionIcon
          variant="filled"
          size={"lg"}
          mr={"xs"}
          onClick={() => {
            if (!queryBarRef.current) return;
            console.log(queryBarRef.current.value);
            flushSync(() => {
              if (!queryBarRef.current) return;
              query.current = queryBarRef.current.value;
              setPage(1);
            });
            queryClient.invalidateQueries({ queryKey: [props.id, 1] });
          }}
        >
          <SearchIcon size={16} />
        </ActionIcon>
        <TextInput placeholder="query" w={"100%"} ref={queryBarRef} />
      </Flex>
      <List highlightHover hideOverflow>
        {isLoading ? (
          <Center>
            <Stack>
              <Center>
                <Loader />
              </Center>
              <Text>Loading logs...</Text>
            </Stack>
          </Center>
        ) : isError ? (
          <Center>
            <Stack>
              <Center>
                <Loader />
              </Center>
              <Text>Failed to load logs: {error.message}</Text>
            </Stack>
          </Center>
        ) : (
          <LogEntryList logs={logs} />
        )}
      </List>
      <Pagination
        page={page}
        onChange={p => {
          setPage(p);
        }}
        total={Math.ceil(props.total / 150)}
        py={"xs"}
      />
    </Center>
  );
};
