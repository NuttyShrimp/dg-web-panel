import { Checkbox, ColorSwatch, Divider, Text, useMantineTheme } from "@mantine/core";
import { FC, useRef, useState } from "react";

import "./styles.scss";

export const TagSelector: FC<
  ReportState.Tag & { selected?: boolean; onSelection?: (toggle: boolean) => void }
> = props => {
  const [selected, setSelected] = useState(props.selected);
  const theme = useMantineTheme();
  const checkboxRef = useRef<HTMLInputElement>(null);

  const updateSelection = () => {
    const val = checkboxRef?.current?.checked ?? false;
    setSelected(val);
    props.onSelection?.(val);
  };

  return (
    <div>
      <div className="tag-selector-container" onClick={() => updateSelection()}>
        <Checkbox ref={checkboxRef} radius={"sm"} size="md" defaultChecked={selected} />
        <ColorSwatch size={18} color={theme.colors[props?.color ?? "dark"]?.[6] ?? "#fff"} />
        <Text size={14} weight={700}>
          {props.name}
        </Text>
      </div>
      <Divider />
    </div>
  );
};
