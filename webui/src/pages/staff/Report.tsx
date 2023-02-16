import {
  ActionIcon,
  Box,
  Card,
  Container,
  Divider,
  Group,
  Loader,
  ScrollArea,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { reportState } from "@src/stores/reports/state";
import dayjs from "dayjs";
import { useCallback, useEffect, useRef } from "react";
import { Navigate, useParams } from "react-router-dom";
import { useRecoilState } from "recoil";
import useWebSocket from "react-use-websocket";

import "@src/styles/pages/staffReport.scss";
import { ReportMessage } from "./components/ReportMessage";
import { CommentEditor } from "@src/components/CommentEditor";
import { ReportTag } from "@src/components/ReportTag";
import { getHostname } from "@src/helpers/axiosInstance";
import { PlusIcon, TrashIcon } from "@primer/octicons-react";
import { openModal } from "@mantine/modals";
import { UserSelectModal } from "@src/components/Users/SelectModal";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { Link } from "@src/components/Router/Link";

export const StaffReport = () => {
  const { id } = useParams();
  const { fetchReport } = useReportActions();
  const queryClient = useQueryClient();
  const {
    data: report,
    error,
    isLoading,
    isError,
  } = useQuery<ReportState.Report, Error>({
    queryKey: ["report", id ?? "0"],
    queryFn: () => fetchReport(Number(id)),
    refetchOnWindowFocus: false,
  });
  const [reportMessages, setReportMessages] = useRecoilState<ReportState.Message[]>(reportState.reportMessages);
  const scrollRef = useRef<HTMLDivElement | null>(null);

  const handleWSMessage = useCallback(
    (message: any) => {
      console.log(message);
      switch (message.type) {
        case "addMessages": {
          setReportMessages(
            [...reportMessages, ...message.data].sort((m1: ReportState.Message, m2: ReportState.Message) => {
              return dayjs(m1.createdAt).isBefore(dayjs(m2.createdAt)) ? -1 : 1;
            })
          );
          setTimeout(() => {
            scrollRef.current?.scrollTo({
              top: scrollRef.current?.scrollHeight,
              behavior: "smooth",
            });
          }, 10);
          break;
        }
        case "addMessage": {
          const shouldScroll =
            0 ===
            (scrollRef.current?.scrollHeight ?? 0) -
              (scrollRef.current?.scrollTop ?? 0) -
              (scrollRef.current?.offsetHeight ?? 0);
          setReportMessages([...reportMessages, message.data]);
          if (!scrollRef || !shouldScroll) return;
          setTimeout(() => {
            scrollRef.current?.scrollTo({
              top: scrollRef.current?.scrollHeight,
              behavior: "smooth",
            });
          }, 10);

          break;
        }
        case "setMembers": {
          if (!report) return;
          queryClient.invalidateQueries(["report", id]);
          break;
        }
        case "error": {
          showNotification({
            title: message?.data?.title ?? "Websocket error",
            message: message?.data?.message ?? "Unknown websocket error",
            color: "red",
          });
          console.error(message.data);
          break;
        }
        default: {
          console.error(`received unknown WS message of type: ${message.type}`);
          break;
        }
      }
    },
    [reportMessages, setReportMessages, scrollRef, id, report]
  );

  const { lastJsonMessage, sendJsonMessage, readyState } = useWebSocket(
    `${location.protocol.includes("https") ? "wss" : "ws"}://${getHostname()}/api/staff/reports/join/${Number(id)}`
  );

  const sendNewMsg = (msg: any) => {
    sendJsonMessage({
      type: "addMessage",
      data: msg,
    });
  };

  const removeMember = (steamId: string) => {
    sendJsonMessage({
      type: "removeMember",
      data: steamId,
    });
  };

  useEffect(() => {
    return () => setReportMessages([]);
  }, []);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      handleWSMessage(lastJsonMessage);
    }
  }, [lastJsonMessage]);

  const openAddMemberModal = () => {
    openModal({
      title: "add member to report",
      children: (
        <UserSelectModal
          onAccept={val => {
            sendJsonMessage({
              type: "addMember",
              data: val,
            });
          }}
        />
      ),
    });
  };

  if (isLoading) {
    return (
      <Stack justify={"center"} align="center" pt={"lg"}>
        <Loader color="gray" />
        <Text>Loading report info</Text>
      </Stack>
    );
  }

  if (isError) {
    return (
      <Stack justify={"center"} align="center" pt={"lg"}>
        <Loader color="gray" />
        <Text>Failed to load report: {error.message}</Text>
        <Text>
          <Link to={"/staff/reports"}>Go back to the list</Link>
        </Text>
      </Stack>
    );
  }

  if (!report) {
    showNotification({
      title: "Unknown Report",
      message: "It seems like you tried to access a report that doesn't exists",
      color: "red",
    });
    return <Navigate to="/staff/reports" replace />;
  }

  if (readyState !== 1) {
    return (
      <Stack justify={"center"} align="center" pt={"lg"}>
        <Loader color="gray" />
        <Text>Connecting to report socket</Text>
      </Stack>
    );
  }

  return (
    <Container size="md">
      <Group align={"top"}>
        <Stack style={{ flexGrow: 1 }} spacing={0}>
          <Box mb="xs">
            <Title className="report-title" order={3}>
              {report.title}
              <span> #{report.id}</span>
            </Title>
            {report.tags && report.tags.map(t => <ReportTag key={t.name} {...t} />)}
          </Box>
          <div>
            <ScrollArea h={"65vh"} viewportRef={scrollRef}>
              {reportMessages.map(r => (
                <ReportMessage key={r.id} message={r} />
              ))}
            </ScrollArea>
            <Divider my="md" size={"md"} />
            <CommentEditor value="" onSubmit={sendNewMsg} />
          </div>
        </Stack>
        <Stack w={"25%"} spacing={4}>
          <Group position="apart">
            <Title order={4}>Members</Title>
            <ActionIcon onClick={openAddMemberModal}>
              <PlusIcon />
            </ActionIcon>
          </Group>
          <Divider />
          {report.members &&
            report.members.map(m => (
              <Card key={m.steamId} shadow="xs" radius="xs" p="xs">
                <Group position="apart">
                  <Stack spacing={1}>
                    <Text>{m.name}</Text>
                    <Text size={"xs"}>{m.steamId}</Text>
                  </Stack>
                  <ActionIcon color="red" onClick={() => removeMember(m.steamId)}>
                    <TrashIcon />
                  </ActionIcon>
                </Group>
              </Card>
            ))}
        </Stack>
      </Group>
    </Container>
  );
};
