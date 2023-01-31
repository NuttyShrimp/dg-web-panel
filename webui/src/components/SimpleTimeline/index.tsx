import { Text, Timeline } from "@mantine/core";
import { EventIcons } from "@src/enums/events";
import { formatRelativeTime } from "@src/helpers/time";
import { FC } from "react";
import { FontAwesomeIcon } from "../Icon";

export const SimpleTimeline: FC<SimpleTimeline.Props> = props => {
  return (
    <Timeline active={props.list.length}>
      {props.list.map(e => (
        <Timeline.Item
          key={`timeline-item-${e.time}`}
          title={e.title}
          bullet={
            e.type && <FontAwesomeIcon icon={EventIcons?.[e.type as keyof typeof EventIcons] ?? "bug"} size={"xs"} />
          }
          bulletSize={props.bulletSize ?? 24}
        >
          <Text size="xs">{formatRelativeTime(e.time)}</Text>
        </Timeline.Item>
      ))}
    </Timeline>
  );
};
