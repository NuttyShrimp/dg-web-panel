import { Box, Container, Divider, Loader, Stack, Text, Title } from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { reportState } from "@src/stores/reports/state";
import dayjs from "dayjs";
import { useCallback, useEffect } from "react";
import { Navigate, useParams } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";
import useWebSocket from "react-use-websocket";

import "@src/styles/pages/staffReport.scss";
import { ReportMessage } from "./components/ReportMessage";
import { CommentEditor } from "@src/components/CommentEditor";
import { ReportTag } from "@src/components/ReportTag";
import { getHostname } from "@src/helpers/axiosInstance";

export const StaffReport = () => {
  const { id } = useParams();
  const report = useRecoilValue(reportState.getReport(Number(id)));
  const [reportMessages, setReportMessages] = useRecoilState<ReportState.Message[]>(reportState.reportMessages);

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
          break;
        }
        case "addMessage": {
          setReportMessages([...reportMessages, message.data]);
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
    [reportMessages, setReportMessages]
  );

  const { lastJsonMessage, sendJsonMessage, readyState } = useWebSocket(
    `${location.protocol === "https" ? "wss" : "ws"}://${getHostname()}/api/staff/reports/join/${Number(id)}`
  );

  const sendNewMsg = (msg: any) => {
    sendJsonMessage({
      type: "addMessage",
      data: msg,
    });
  };

  useEffect(() => {
    if (report) {
      setReportMessages([]);
    }
  }, [report]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      handleWSMessage(lastJsonMessage);
    }
  }, [lastJsonMessage]);

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
    <Container size="sm">
      <Box mb="xs">
        <Title className="report-title" order={3}>
          {report.title}
          <span> #{report.id}</span>
        </Title>
        {report.tags && report.tags.map(t => <ReportTag key={t.name} {...t} />)}
      </Box>
      <div>
        {reportMessages.map(r => (
          <ReportMessage key={r.id} message={r} />
        ))}
        <Divider my="md" size={"md"} />
        <CommentEditor value="" onSubmit={sendNewMsg} />
      </div>
    </Container>
  );
};
