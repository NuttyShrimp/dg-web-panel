import { MultiSelect } from "@mantine/core";
import type { SelectItem } from "@mantine/core";
import { reasons } from "@src/data/PenaltyReasons";
import { useMemo, useState } from "react";

declare interface PRSProps {
  reasons: string[];
  setReasons: (v: string[]) => void;
}

export const PenaltyReasonSelector = (props: PRSProps) => {
  const [extraReasons, setExtraReasons] = useState<SelectItem[]>([]);
  const defaultReasons = useMemo(() => {
    const values: SelectItem[] = [];
    for (const key in reasons) {
      values.push({
        value: key,
        label: `${key} (${reasons[key]})`,
        group: reasons[key],
      });
    }
    return values;
  }, []);
  return (
    <MultiSelect
      label="Penalty Reasons"
      placeholder="Select or create the reasons"
      searchable
      creatable
      data={[...extraReasons, ...defaultReasons]}
      value={props.reasons}
      onChange={props.setReasons}
      getCreateLabel={q => `Add ${q}`}
      onCreate={q => {
        const item = { value: q, label: q };
        setExtraReasons(r => [...r, item]);
        return item;
      }}
    />
  );
};
