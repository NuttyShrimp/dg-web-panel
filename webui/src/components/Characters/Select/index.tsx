import { Box, Button, Center, Text } from "@mantine/core";
import { AutoComplete } from "@src/components/ui/Autocomplete";
import { characterState } from "@src/stores/character/state";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { FC, useEffect, useState } from "react";
import { useRecoilValue } from "recoil";

declare type CharacterItemProps = React.ComponentPropsWithoutRef<"div"> & CfxState.Character;

const SelectItem = ({
  info,
  citizenid,
  steamId: _steamId,
  data: _data,
  user: _user,
  created_at: _created_at,
  last_updated: _last_updated,
  ...others
}: CharacterItemProps) => (
  <div {...others}>
    <Box>
      <Text size="sm">
        {info.firstname} {info.lastname}
      </Text>
      <Text size="xs" c="dimmed">
        {citizenid}
      </Text>
    </Box>
  </div>
);

const caseInsensitiveMatch = (s1: string, s2: string) => s1.toLowerCase().includes(s2.toLowerCase().trim());

export const SelectCharacter: FC<{ cid?: string; onChange: (cid: number) => void }> = ({ cid, onChange }) => {
  const characters = useRecoilValue(characterState.list);
  const { fetchCharacters } = useCharacterActions();

  useEffect(() => {
    if (characters.length === 0) {
      fetchCharacters();
    }
  }, []);

  return (
    <AutoComplete
      placeholder="Search a character"
      itemComponent={SelectItem}
      defaultValue={cid}
      data={characters.map(c => ({
        ...c,
        value: String(c.citizenid),
        label: `${c.info.firstname} ${c.info.lastname}`,
      }))}
      filter={(item, search) =>
        caseInsensitiveMatch(item.info.firstname, search) ||
        caseInsensitiveMatch(item.info.lastname, search) ||
        caseInsensitiveMatch(`${item.info.firstname} ${item.info.lastname}`, search) ||
        caseInsensitiveMatch(String(item.citizenid), search)
      }
      onOptionSubmit={val => onChange?.(Number(val) ?? 0)}
    />
  );
};

export const SelectCharacterModal = ({ onAccept }: { onAccept: (cid: number) => void }) => {
  const [selectedCid, setSelectedCid] = useState(0);
  return (
    <>
      <SelectCharacter onChange={setSelectedCid} />
      <Center mt={4}>
        <Button onClick={() => onAccept(selectedCid)}>Accept</Button>
      </Center>
    </>
  );
};
