import { Button, Group, Text } from "@mantine/core";
import { MarkdownIcon } from "@primer/octicons-react";
import { usePanelEditor } from "@src/hooks/usePanelEditor";
import { EditorContent } from "@tiptap/react";
import { FC } from "react";

import "./styles.scss";

declare interface CommentEditorProps {
  value: string;
  onSubmit: (data: any) => void;
}

export const CommentEditor: FC<CommentEditorProps> = props => {
  const editor = usePanelEditor({
    content: props.value ?? "",
  });

  const handleSubmit = () => {
    if (editor?.isEmpty) return;
    props.onSubmit(editor?.getJSON());
    editor?.commands.clearContent();
  };

  return (
    <div className="editor-wrapper">
      <div>
        <EditorContent editor={editor} />
      </div>
      <Group position="apart">
        <Button color={"green"} onClick={handleSubmit}>
          Submit
        </Button>
        <Text size={"sm"}>
          <MarkdownIcon size={16} /> supported
        </Text>
      </Group>
    </div>
  );
};
