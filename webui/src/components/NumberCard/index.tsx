import { Paper, Text, Title } from "@mantine/core";
import { FC } from "react";
import { useStyles } from "./styles";

export const NumberCard: FC<NumberCard.Props> = props => {
  const { classes } = useStyles();
  return (
    <Paper withBorder p="md" className={classes.root}>
      <Text size="sm" color="dimmed">
        {props.title}
      </Text>
      <Title order={2}>{props.count}</Title>
    </Paper>
  );
};
