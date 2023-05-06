import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useSetRecoilState, useRecoilState } from "recoil";
import { reportState } from "./state";

let reportFetchId = 1;

export const useReportActions = () => {
  const setReports = useSetRecoilState(reportState.list);
  const setLoadingReports = useSetRecoilState(reportState.loadingList);
  const [filter, setFilter] = useRecoilState(reportState.listFilter);
  const [pagination, setPagination] = useRecoilState(reportState.pagination);
  const navigate = useNavigate();

  const loadReports = useCallback(
    async (pFilter?: ReportState.Filter) => {
      setLoadingReports(true);
      const fetchId = ++reportFetchId;
      pFilter = pFilter ?? filter;
      try {
        const res = await axiosInstance.get<{ reports: ReportState.Report[]; total: number }>("/staff/reports/all", {
          params: {
            filter: pFilter.search,
            open: pFilter.open,
            closed: pFilter.closed,
            offset: pagination.current - 1,
          },
        });
        if (fetchId !== reportFetchId) return;
        if (res.status !== 200) {
          showNotification({
            title: "Error while fetching reports",
            message: "Encountered an unexpected error while trying to fetch reports",
            color: "red",
          });
          setReports([]);
          return;
        }
        setPagination({
          total: res.data?.total,
          current: 1,
        });
        setReports(res.data?.reports ?? []);
      } catch (e) {
        console.error(e);
      } finally {
        if (fetchId === reportFetchId) {
          setLoadingReports(false);
        }
      }
    },
    [filter, setReports, setLoadingReports, setPagination, pagination]
  );

  const createReport = useCallback(
    async (title: string, members: string[]) => {
      try {
        const res = await axiosInstance.post<{ token: string }>("/staff/reports/new", {
          title,
          members,
        });
        if (res.status !== 200) return;
        await loadReports();
        navigate(`/staff/reports/${res.data.token}`);
      } catch (e) {
        console.error(e);
      }
    },
    [navigate, loadReports]
  );

  const updateFilter = useCallback(
    (partFilter: Partial<ReportState.Filter>) => {
      const newFilter = { ...filter, ...partFilter };
      setFilter(newFilter);
      loadReports(newFilter);
    },
    [loadReports, setFilter, filter]
  );

  const fetchReport = async (id: number) => {
    if (Number.isNaN(id)) {
      throw new Error(`${id} is not a number`);
    }
    const res = await axiosInstance.get<{ report: ReportState.Report }>(`/staff/reports/${id}`);
    return res.data.report;
  };

  return {
    createReport,
    loadReports,
    updateFilter,
    fetchReport,
  };
};
