import { Button, Group, Modal, MultiSelect, TextInput } from "@mantine/core";
import { cfxState } from "@src/stores/cfx/state";
import { useCfxPlayer } from "@src/stores/cfx/hooks/useCfxPlayer";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { FC, useEffect, useState } from "react";
import { flushSync } from "react-dom";
import { useRecoilValue } from "recoil";

interface ReportData {
  title: string;
  members: string[];
}

export const NewReportModal: FC<{ open: boolean; onClose: () => void }> = props => {
  const [data, setData] = useState<ReportData>({
    title: "",
    members: [],
  });
  const [creatingReport, setCreatingReport] = useState(false);
  const { loadPlayers } = useCfxPlayer();
  const { createReport } = useReportActions();
  const cfxPlayers = useRecoilValue(cfxState.selectPlayers);

  const changeDataEntry = <T extends keyof ReportData>(key: T, val: ReportData[T]) => {
    setData({
      ...data,
      [key]: val,
    });
  };

  const addReport = () => {
    flushSync(() => setCreatingReport(false));
    createReport(data.title, data.members);
  };

  useEffect(() => {
    loadPlayers();
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
      <Group position="right" pt={"md"}>
        <Button color={"dg-prim"} onClick={addReport} loading={creatingReport}>
          <p>Create</p>
        </Button>
      </Group>
    </Modal>
  );
};
