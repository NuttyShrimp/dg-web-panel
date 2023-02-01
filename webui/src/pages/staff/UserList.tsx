import { Center, Container, Stack, Title } from "@mantine/core";
import { SearchAndSelect } from "@src/components/Users/SearchAndSelect";
import { useCfxPlayer } from "@src/stores/cfx/hooks/useCfxPlayer";
import { useEffect } from "react";

export const UserList = () => {
  const { loadPlayers } = useCfxPlayer();
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
