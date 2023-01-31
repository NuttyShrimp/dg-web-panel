import { Center, Container, Stack, Title } from "@mantine/core";
import { SearchAndSelect } from "@src/components/Characters/SearchAndSelect";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { useEffect } from "react";

export const CharacterList = () => {
  const { fetchCharacters } = useCharacterActions();
  useEffect(() => {
    fetchCharacters();
  }, [fetchCharacters]);

  return (
    <Container>
      <Center>
        <Stack>
          <SearchAndSelect />
          <Title size="h2">Search for a character</Title>
        </Stack>
      </Center>
    </Container>
  );
};
