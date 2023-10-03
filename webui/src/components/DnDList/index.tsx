import { Card, Center, Text, Title, useMantineTheme } from "@mantine/core";
import { FC } from "react";

// https://github.com/mantinedev/ui.mantine.dev/blob/master/components/DndListHandle/DndListHandle.tsx
export const DnDList: FC<DnDList.Props> = props => {
  const theme = useMantineTheme();

  return (
    <Card
      p={"xs"}
      style={{ width: props.maxWidth ? "100%" : "auto", backgroundColor: theme.colors.dark[6] }}
      withBorder
    >
      <Title order={3} pb={"xs"}>
        {props.title}
      </Title>
      <Card.Section
        style={{
          background: theme.colors.dark[7],
        }}
      >
        {props.elements.length === 0 ? (
          <Center>
            <Text py="sm">{props.emptyListHolder ?? "Nothing to see here"}</Text>
          </Center>
        ) : (
          <p>DNDList</p>
        )}
      </Card.Section>
    </Card>
  );
};
