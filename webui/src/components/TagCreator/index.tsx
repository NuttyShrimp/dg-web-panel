import { Badge, Button, CloseButton, ColorInput, Divider, TextInput, useMantineTheme } from "@mantine/core";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { FC, useState } from "react";
import { flushSync } from "react-dom";
import "./styles.scss";

export const TagCreator: FC<{ onClose: () => void }> = ({ onClose }) => {
  const [color, setColor] = useState("");
  const [colorName, setColorName] = useState("");
  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);
  const { createTag } = useReportActions();
  const theme = useMantineTheme();

  const assignColor = (color: string) => {
    setColor(color);
    Object.entries(theme.colors).forEach(([k, v]) => {
      if (v?.[6] === color) {
        setColorName(k);
      }
    });
  };

  const addTag = () => {
    flushSync(() => setCreating(true));
    createTag(name, colorName);
    onClose();
    setCreating(false);
  };

  return (
    <div className="tag-creator-wrapper">
      <div className="tag-creator-container">
        <Badge color={colorName}>
          <p>{name === "" ? "Tag Preview" : name}</p>
        </Badge>
      </div>
      <div className="tag-creator-close">
        <CloseButton onClick={onClose} />
      </div>
      <div className="tag-creator-inputs">
        <TextInput
          placeholder={"Tag name"}
          label={"Tag name"}
          value={name}
          onChange={e => setName(e.currentTarget.value)}
        />
        <ColorInput
          placeholder="color"
          label="Color"
          withPicker={false}
          withPreview={false}
          value={color}
          onChange={c => assignColor(c)}
          format="hex"
          swatches={Object.values(theme.colors).map(c => c?.[6] ?? "#fff")}
        />
        <Button leftIcon={<i className="fas fa-plus" />} onClick={addTag} loading={creating}>
          Add
        </Button>
      </div>
      <Divider />
    </div>
  );
};
