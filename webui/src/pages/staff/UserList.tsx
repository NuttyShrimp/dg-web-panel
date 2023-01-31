import { Center, Container, Stack, Title } from "@mantine/core";
import { SearchAndSelect } from "@src/components/Users/SearchAndSelect";
import { useCfxActions } from "@src/stores/cfx/useCfxActions";
import { useEffect } from "react";

export const UserList = () => {
  const { loadPlayers } = useCfxActions();
  useEffect(() => {
    loadPlayers();
  }, [loadPlayers]);

  return (
    <Container>
      <Center>
        <Stack>
          <SearchAndSelect />
          <Title size="h2">Search for a player</Title>
        </Stack>
      </Center>
    </Container>
  );
};
