import { Center, Divider, Loader, Stack, Text, Title, useMantineTheme } from "@mantine/core";
import { IssueClosedIcon, IssueOpenedIcon } from "@primer/octicons-react";
import { formatRelativeTime } from "@src/helpers/time";
import { reportState } from "@src/stores/reports/state";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { FC, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue } from "recoil";

import "./styles.scss";

const ReportEntry: FC<{ report: ReportState.Report }> = ({ report }) => {
  const theme = useMantineTheme();
  const navigate = useNavigate();

  const toReport = () => {
    navigate(`/staff/reports/${report.id}`);
  };

  return (
    <>
      <div className="report-list-entry-container">
        <div className="report-list-entry-state">
          {report.open ? (
            <IssueOpenedIcon fill={theme.colors.green[6]} size={16} />
          ) : (
            <IssueClosedIcon fill={theme.colors.indigo[6]} size={16} />
          )}
        </div>
        <div className="report-list-entry-info">
          <Title order={5} onClick={toReport}>
            {report.title}
          </Title>
          <div className="report-list-entry-date">
            <Text color={"dimmed"} size={"xs"}>
              #{report.id} opened {formatRelativeTime(new Date(report.createdAt).getTime() / 1000)}
            </Text>
            {report.updatedAt && (
              <Text size={"xs"} pl={"sm"}>
                <i className="fal fa-clock" /> updated {formatRelativeTime(new Date(report.updatedAt).getTime() / 1000)}
              </Text>
            )}
          </div>
        </div>
      </div>
      <Divider />
    </>
  );
};

export const ReportList: FC<{}> = () => {
  const reports = useRecoilValue(reportState.list);
  const loadingReports = useRecoilValue(reportState.loadingList);
  const { loadReports } = useReportActions();

  useEffect(() => {
    loadReports();
  }, []);

  if (loadingReports) {
    return (
      <Stack justify={"center"} align="center" pt={"lg"}>
        <Loader color="gray" />
        <Text>Loading reports</Text>
      </Stack>
    );
  }

  if (reports.length < 1) {
    return (
      <Center mt={"lg"}>
        <Text color={"dimmed"} size={"lg"} weight={700}>
          No reports found
        </Text>
      </Center>
    );
  }

  return (
    <div>
      {reports.map(r => (
        <ReportEntry key={r.id} report={r} />
      ))}
    </div>
  );
};
