import { Paper, Text, Title, useMantineTheme } from "@mantine/core";
import { FC } from "react";

export const NumberCard: FC<NumberCard.Props> = props => {
  const theme = useMantineTheme();

  return (
    <Paper
      withBorder
      p="md"
      style={{
        backgroundColor: theme.colors.dark[6],
        width: "100%",
      }}
    >
      <Text size="sm" c="dimmed">
        {props.title}
      </Text>
      <Title order={2}>{props.count}</Title>
    </Paper>
  );
};
