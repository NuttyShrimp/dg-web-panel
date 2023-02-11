import { Center, Container, Stack, Title } from "@mantine/core";
import { ActiveCharacterList } from "@src/components/Characters/ActiveList";
import { SelectCharacter } from "@src/components/Characters/Select";
import { characterState } from "@src/stores/character/state";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useSetRecoilState } from "recoil";

export const CharacterList = () => {
  const { fetchCharacters, resetStores } = useCharacterActions();
  const selectCid = useSetRecoilState(characterState.cid);
  const navigate = useNavigate();

  useEffect(() => {
    fetchCharacters();
  }, []);

  return (
    <Container>
      <Center>
        <Stack>
          <SelectCharacter
            onChange={cid => {
              resetStores();
              selectCid(cid);
              navigate(`/staff/characters/${cid}`);
            }}
          />
          <Title size="h2">Search for a character</Title>
          <ActiveCharacterList />
        </Stack>
      </Center>
    </Container>
  );
};
