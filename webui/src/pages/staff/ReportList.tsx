import { Center, Container, Group, Pagination, Text, Tooltip } from "@mantine/core";
import { SearchInput } from "@src/components/SearchInput";
import { reportState } from "@src/stores/reports/state";
import { useCallback, useEffect, useState } from "react";
import { useRecoilState, useRecoilValue } from "recoil";
import { ReportList } from "./components/ReportList";

import "src/styles/pages/staffReports.scss";
import { CheckIcon } from "@primer/octicons-react";
import { NewReportModal } from "./components/NewReportModal";
import { useReportActions } from "@src/stores/reports/useReportActions";

export const StaffReports = () => {
  const [createReport, setCreateReport] = useState(false);
  const [pagination, setPagination] = useRecoilState(reportState.pagination);
  const filter = useRecoilValue(reportState.listFilter);
  const { updateFilter, loadReports } = useReportActions();

  const updateSearchValue = (val: string) => {
    updateFilter({
      search: val,
    });
  };

  const handleOpenFilter = useCallback(() => {
    updateFilter({
      open: !filter.open,
      closed: false,
    });
  }, [filter, updateFilter]);

  const handleClosedFilter = useCallback(() => {
    updateFilter({
      open: false,
      closed: !filter.closed,
    });
  }, [filter, updateFilter]);

  const handlePageChange = (page: number) => {
    setPagination({
      ...pagination,
      current: page,
    });
  };

  useEffect(() => {
    loadReports();
  }, [pagination]);

  return (
    <>
      <NewReportModal open={createReport} onClose={() => setCreateReport(false)} />
      <Container my="sm" size="xl">
        <div className="reports-list-wrapper">
          <div className="reports-list-container">
            <div className="reports-list-header">
              <Group spacing={"xs"}>
                <Text size={"sm"} weight={filter.open ? 700 : 400} onClick={handleOpenFilter}>
                  {filter.open && <CheckIcon size={16} />}
                  Open
                </Text>
                <Text size={"sm"} weight={filter.closed ? 700 : 400} onClick={handleClosedFilter}>
                  {filter.closed && <CheckIcon size={16} />}
                  Closed
                </Text>
              </Group>
              <Group spacing="xs" className="reports-list-actions">
                <SearchInput value={filter.search} onChange={updateSearchValue} />
                <Tooltip label={"Create new report"} position="bottom">
                  <div onClick={() => setCreateReport(true)}>
                    <i className="fas fa-plus" />
                  </div>
                </Tooltip>
              </Group>
            </div>
            <ReportList />
          </div>
        </div>
        <Center mt={"md"}>
          <Pagination total={pagination.total} siblings={2} page={pagination.current} onChange={handlePageChange} />
        </Center>
      </Container>
    </>
  );
};
