import { Box, Button, Center, Code, Container, Divider, Flex, Pagination, Stack, Text } from "@mantine/core";
import { displayUnixDate } from "@src/helpers/time";
import { useState } from "react";
import { List } from "../List";

interface LogListProps {
  logs: Logs.Log[];
  onLoadMore: () => void;
}

const LogInfo = ({ log }: { log: Logs.Log }) => (
  <Container size={"sm"}>
    {Object.keys(log).map(key => (
      <Box key={`${log._id}-${key}`} pt="xs">
        <Text weight={"bolder"}>{key}</Text>
        <Code block>{JSON.stringify(log[key as keyof Logs.Log], undefined, 2)}</Code>
      </Box>
    ))}
  </Container>
);

export const LogList = (props: LogListProps) => {
  const [focusedLogs, setFocusedLogs] = useState<string[]>([]);
  const [activePage, setPage] = useState(1);

  const toggledFocus = (id: string) => {
    if (focusedLogs.includes(id)) {
      setFocusedLogs(focusedLogs.filter(fId => fId !== id));
    } else {
      setFocusedLogs([...focusedLogs, id]);
    }
  };

  return (
    <Center style={{ flexDirection: "column" }}>
      <p>QUERY SEARCH BAR</p>
      <List highlightHover hideOverflow>
        {props.logs.slice((activePage - 1) * 50, activePage * 50).map(l => (
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
      </List>
      <Pagination page={activePage} onChange={setPage} total={Math.ceil(props.logs.length / 50)} py={"xs"} />
      {Math.ceil(props.logs.length / 50) === activePage && (
        <Center py={"xs"}>
          <Button onClick={props.onLoadMore}>Load more</Button>
        </Center>
      )}
    </Center>
  );
};
