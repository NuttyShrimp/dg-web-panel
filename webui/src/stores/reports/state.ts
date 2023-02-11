import { atom } from "recoil";

export const reportState = {
  loadingList: atom({
    key: "report-loading-list",
    default: false,
  }),
  pagination: atom<{ total: number; current: number }>({
    key: "reports-list-pagination",
    default: {
      total: 1,
      current: 1,
    },
  }),
  list: atom<ReportState.Report[]>({
    key: "reports-list",
    default: [],
  }),
  listFilter: atom<ReportState.Filter>({
    key: "reports-list-filter",
    default: {
      open: true,
      closed: false,
      search: "",
    },
  }),
  loadingTags: atom<boolean>({
    key: "reports-loading-tags",
    default: false,
  }),
  tags: atom<ReportState.Tag[]>({
    key: "reports-tags",
    default: [],
  }),
  selectedTags: atom<ReportState.Tag[]>({
    key: "reports-selected-tags",
    default: [],
  }),
  reportMessages: atom<ReportState.Message[]>({
    key: "report-current-messages",
    default: [],
  }),
};
