import { Card, Center, Text, Title } from "@mantine/core";
import { FC } from "react";
import { useStyles } from "./styles";

// https://github.com/mantinedev/ui.mantine.dev/blob/master/components/DndListHandle/DndListHandle.tsx
export const DnDList: FC<DnDList.Props> = props => {
  const { classes } = useStyles();
  return (
    <Card className={classes.card} p={"xs"} style={{ width: props.maxWidth ? "100%" : "auto" }} withBorder>
      <Title order={3} pb={"xs"}>
        {props.title}
      </Title>
      <Card.Section className={classes.list}>
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
