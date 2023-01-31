import { ActionIcon, Divider, Group, Popover, Text } from "@mantine/core";
import { reportState } from "@src/stores/reports/state";
import { useReportActions } from "@src/stores/reports/useReportActions";
import { FC } from "react";
import { useRecoilValue } from "recoil";
import { TagSelector } from "../TagSelector";

export const TagsPopup: FC<{ onCreateTag: () => void }> = ({ onCreateTag }) => {
  const { refreshTags, selectTag, unSelectTag, clearSelectedTags } = useReportActions();
  const tags = useRecoilValue(reportState.tags);
  const selectedTags = useRecoilValue(reportState.selectedTags);

  const toggleTag = (tag: ReportState.Tag, isSelected: boolean) => {
    isSelected ? selectTag(tag) : unSelectTag(tag);
  };

  return (
    <Popover withArrow width={200} onOpen={() => refreshTags()}>
      <Popover.Target>
        <i className="fas fa-filter" />
      </Popover.Target>
      <Popover.Dropdown style={{ padding: 0 }}>
        <div className="tags-list-header">
          <Group position="apart">
            <Text weight={700} size={"sm"}>
              Filter by tag
            </Text>
            <div className="tags-list-actions">
              {selectedTags.length > 0 ? (
                <ActionIcon onClick={clearSelectedTags}>
                  <i className="fas fa-xmark" />
                </ActionIcon>
              ) : (
                <div />
              )}
              <ActionIcon variant="transparent" onClick={onCreateTag}>
                <i className="fas fa-plus" />
              </ActionIcon>
            </div>
          </Group>
        </div>
        <Divider />
        {selectedTags.map(t => (
          <TagSelector key={t.name} {...t} selected={true} onSelection={() => toggleTag(t, false)} />
        ))}
        {tags
          .filter(t => !selectedTags.find(st => st.name === t.name))
          .map(t => (
            <TagSelector key={t.name} {...t} selected={false} onSelection={() => toggleTag(t, true)} />
          ))}
      </Popover.Dropdown>
    </Popover>
  );
};
