import { Avatar, Text } from "@mantine/core";
import { displayDate } from "@src/helpers/time";
import { EditorContent } from "@tiptap/react";
import { FC } from "react";
import "./styles.scss";
import { usePanelEditor } from "@src/hooks/usePanelEditor";

export const ReportMessage: FC<{ message: ReportState.Message }> = ({ message }) => {
  const editor = usePanelEditor({
    editable: false,
    content: JSON.parse(message.message),
  });

  return (
    <div className="report-message-wrapper">
      <div className="report-message-header">
        <div>
          <Avatar src={message.sender.avatarUrl} radius="xl" />
          <Text className="report-message-title" weight={"bold"} ml={"xs"}>
            {message.sender.username}
            <span> op {displayDate(message.createdAt)}</span>
          </Text>
        </div>
      </div>
      <div className="report-message-content">
        <EditorContent editor={editor} />
      </div>
    </div>
  );
};
