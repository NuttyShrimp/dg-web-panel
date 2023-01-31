import { Box, Select, Text } from "@mantine/core";
import { characterState } from "@src/stores/character/state";
import { useCharacterActions } from "@src/stores/character/useCharacterActions";
import { FC, forwardRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue, useSetRecoilState } from "recoil";

declare type CharacterItemProps = React.ComponentPropsWithoutRef<"div"> & CfxState.Character;

const SelectItem = forwardRef<HTMLDivElement, CharacterItemProps>(
  (
    {
      info,
      citizenid,
      steamId: _steamId,
      data: _data,
      user: _user,
      created_at: _created_at,
      last_updated: _last_updated,
      ...others
    }: CharacterItemProps,
    ref
  ) => (
    <div ref={ref} {...others}>
      <Box>
        <Text size="sm">
          {info.firstname} {info.lastname}
        </Text>
        <Text size="xs" color="dimmed">
          {citizenid}
        </Text>
      </Box>
    </div>
  )
);

const caseInsensitiveMatch = (s1: string, s2: string) => s1.toLowerCase().includes(s2.toLowerCase().trim());

export const SearchAndSelect: FC<{ cid?: string }> = ({ cid }) => {
  const characters = useRecoilValue(characterState.list);
  const selectCid = useSetRecoilState(characterState.cid);
  const navigate = useNavigate();
  const { resetStores, fetchCharacters } = useCharacterActions();

  useEffect(() => {
    if (characters.length === 0) {
      fetchCharacters();
    }
  }, []);

  return (
    <Select
      placeholder="Search a character"
      searchable
      nothingFound="No character found"
      itemComponent={SelectItem}
      maxDropdownHeight={300}
      defaultValue={cid}
      data={characters.map(c => ({
        ...c,
        value: String(c.citizenid),
        label: `${c.info.firstname} ${c.info.lastname}`,
      }))}
      filter={(value, item) =>
        caseInsensitiveMatch(item.info.firstname, value) ||
        caseInsensitiveMatch(item.info.lastname, value) ||
        caseInsensitiveMatch(`${item.info.firstname} ${item.info.lastname}`, value) ||
        caseInsensitiveMatch(String(item.citizenid), value)
      }
      onChange={val => {
        resetStores();
        selectCid(Number(val));
        navigate(`/staff/characters/${val}`);
      }}
    />
  );
};
