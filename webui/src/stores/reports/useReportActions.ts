import { showNotification } from "@mantine/notifications";
import { axiosInstance } from "@src/helpers/axiosInstance";
import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useSetRecoilState, useRecoilState } from "recoil";
import { reportState } from "./state";

let reportFetchId = 1;

export const useReportActions = () => {
  const setReports = useSetRecoilState(reportState.list);
  const setTags = useSetRecoilState(reportState.tags);
  const setLoadingReports = useSetRecoilState(reportState.loadingList);
  const [selectedTags, setSelectedTags] = useRecoilState(reportState.selectedTags);
  const [loadingTags, setLoadingTags] = useRecoilState(reportState.loadingTags);
  const [filter, setFilter] = useRecoilState(reportState.listFilter);
  const [pagination, setPagination] = useRecoilState(reportState.pagination);
  const navigate = useNavigate();

  const loadReports = useCallback(
    async (pFilter?: ReportState.Filter, tags?: string[]) => {
      setLoadingReports(true);
      const fetchId = ++reportFetchId;
      pFilter = pFilter ?? filter;
      tags = tags ?? selectedTags.map(t => t.name);
      try {
        const res = await axiosInstance.get<{ reports: ReportState.Report[]; total: number }>("/staff/reports", {
          params: {
            filter: pFilter.search,
            open: pFilter.open,
            closed: pFilter.closed,
            tags,
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
    [selectedTags, filter, setReports, setLoadingReports, setPagination, pagination]
  );

  const clearSelectedTags = useCallback(() => {
    setSelectedTags([]);
    loadReports(undefined, []);
  }, [setSelectedTags, loadReports]);

  const selectTag = useCallback(
    (tag: ReportState.Tag) => {
      if (selectedTags.includes(tag)) return;
      setSelectedTags([...selectedTags, tag]);
      loadReports(
        undefined,
        [...selectedTags, tag].map(t => t.name)
      );
    },
    [selectedTags, setSelectedTags, loadReports]
  );

  const unSelectTag = useCallback(
    (tag: ReportState.Tag) => {
      if (!selectedTags.includes(tag)) return;
      const newSelTags = selectedTags.filter(t => t.name !== tag.name);
      setSelectedTags(newSelTags);
      loadReports(
        undefined,
        newSelTags.map(t => t.name)
      );
    },
    [selectedTags, setSelectedTags, loadReports]
  );

  const refreshTags = useCallback(async () => {
    if (loadingTags) return;
    setLoadingTags(true);
    try {
      const res = await axiosInstance.get<ReportState.Tag[]>("/staff/reports/tags");
      if (res.status !== 200) {
        return;
      }
      setTags(res.data);
    } catch (e) {
      console.error(e);
    } finally {
      setLoadingTags(false);
    }
  }, [loadingTags, setLoadingTags, setTags]);

  const createTag = useCallback(async (name: string, color: string) => {
    try {
      const res = await axiosInstance.put("/staff/reports/tags", {
        name,
        color,
      });
      if (res.status !== 200) {
        showNotification({
          title: "",
          message: "",
        });
      }
    } catch (e) {
      console.error(e);
    }
  }, []);

  const createReport = useCallback(
    async (title: string, members: string[], tags: string[]) => {
      try {
        const res = await axiosInstance.post<{ token: string }>("/staff/reports/new", {
          title,
          members,
          tags,
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

  return { refreshTags, selectTag, unSelectTag, clearSelectedTags, createTag, createReport, loadReports, updateFilter };
};
