import { Button, Group, Modal, MultiSelect, TextInput } from "@mantine/core";
import { ReportTag } from "@src/components/ReportTag";
import { cfxState } from "@src/stores/cfx/state";
import { useCfxActions } from "@src/stores/cfx/useCfxActions";
import { reportState } from "@src/stores/reports/state";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { FC, useEffect, useState } from "react";
import { flushSync } from "react-dom";
import { useRecoilValue } from "recoil";

interface ReportData {
  title: string;
  members: string[];
  tags: string[];
}

export const NewReportModal: FC<{ open: boolean; onClose: () => void }> = props => {
  const [data, setData] = useState<ReportData>({
    title: "",
    members: [],
    tags: [],
  });
  const [creatingReport, setCreatingReport] = useState(false);
  const { loadPlayers } = useCfxActions();
  const { refreshTags, createReport } = useReportActions();
  const cfxPlayers = useRecoilValue(cfxState.selectPlayers);
  const tags = useRecoilValue(reportState.tags);

  const changeDataEntry = <T extends keyof ReportData>(key: T, val: ReportData[T]) => {
    setData({
      ...data,
      [key]: val,
    });
  };

  const addReport = () => {
    flushSync(() => setCreatingReport(false));
    createReport(data.title, data.members, data.tags);
  };

  useEffect(() => {
    loadPlayers();
    refreshTags();
  }, []);

  return (
    <Modal opened={props.open} title="Create new report" onClose={props.onClose}>
      <TextInput
        label="Title"
        placeholder="Title"
        value={data.title}
        onChange={e => changeDataEntry("title", e.currentTarget.value)}
      />
      <MultiSelect
        searchable
        clearable
        limit={20}
        label={"Members (staff excluded)"}
        data={cfxPlayers}
        nothingFound="Nothing found"
        value={data.members}
        onChange={data => changeDataEntry("members", data)}
      />
      <MultiSelect
        searchable
        clearable
        limit={20}
        itemComponent={({ label, ...props }) => <ReportTag {...props} name={label} />}
        valueComponent={({ label, ...props }) => <ReportTag {...props} name={label} />}
        label={"Tags"}
        data={tags.map(t => ({ value: t.name, label: t.name, color: t.color }))}
        nothingFound="Nothing found"
        value={data.tags}
        onChange={data => changeDataEntry("tags", data)}
      />
      <Group position="right" pt={"md"}>
        <Button color={"dg-prim"} onClick={addReport} loading={creatingReport}>
          <p>Create</p>
        </Button>
      </Group>
    </Modal>
  );
};
