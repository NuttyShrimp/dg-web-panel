import { Container, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { List } from "../List";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { LoadingSpinner } from "../LoadingSpinner";

export const ActiveCharacterList = () => {
  const { fetchActiveCharacters } = useCharacterActions();
  const { isLoading, isError, error, data } = useQuery<CfxState.Character[], Error>({
    queryKey: ["characters"],
    queryFn: () => fetchActiveCharacters(),
  });

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError) {
    return <Text>Failed to load active characters: {error.message}</Text>;
  }

  if (data.length === 0) {
    return <Text>Geen active characters</Text>;
  }

  return (
    <Container>
      <Title order={4}>Active characters</Title>
      <List>
        {data.map(c => (
          <List.Entry key={c.citizenid}>
            <Text>
              {c.info.firstname} {c.info.lastname}
            </Text>
          </List.Entry>
        ))}
      </List>
    </Container>
  );
};
