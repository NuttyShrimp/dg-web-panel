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
    const classToReasons: Record<string, string[]> = {};
    for (const key in reasons) {
      if (!classToReasons[reasons[key]]) {
        classToReasons[reasons[key]] = [];
      }
      classToReasons[reasons[key]].push(key);
    }
    for (const pClass in classToReasons) {
      values.push({
        group: pClass,
        items: classToReasons[pClass].map(r => ({
          value: r,
          label: `${r} (${pClass})`,
        })),
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
