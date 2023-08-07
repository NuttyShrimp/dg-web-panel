import { Avatar, Text } from "@mantine/core";
import { displayTimeDate } from "@src/helpers/time";
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
          <Avatar src={message.sender.avatarUrl} radius="xl" size={"sm"} />
          <Text className="report-message-title" weight={"bold"} ml={"xs"}>
            {message.sender.username}
            <span> op {displayTimeDate(message.createdAt)}</span>
          </Text>
        </div>
      </div>
      <div>
        <EditorContent editor={editor} />
      </div>
    </div>
  );
};
