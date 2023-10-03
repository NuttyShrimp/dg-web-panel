import { TagsInput } from "@mantine/core";
import { ComboboxData } from "@mantine/core";
import { reasons } from "@src/data/PenaltyReasons";
import { useMemo } from "react";

declare interface PRSProps {
  reasons: string[];
  setReasons: (v: string[]) => void;
}

export const PenaltyReasonSelector = (props: PRSProps) => {
  const defaultReasons = useMemo(() => {
    const values: ComboboxData = [];
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
    <TagsInput
      label="Penalty Reasons"
      placeholder="Select or create the reasons"
      data={defaultReasons}
      value={props.reasons}
      onChange={props.setReasons}
    />
  );
};
