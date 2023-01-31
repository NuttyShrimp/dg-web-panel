import { Box, Select, Text } from "@mantine/core";
import { cfxState } from "@src/stores/cfx/state";
import { useCfxActions } from "@src/stores/cfx/useCfxActions";
import { FC, forwardRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue, useSetRecoilState } from "recoil";

declare type UserItemProps = React.ComponentPropsWithoutRef<"div"> & CfxState.Player;

const SelectItem = forwardRef<HTMLDivElement, UserItemProps>(({ steamId, name, ...others }: UserItemProps, ref) => (
  <div ref={ref} {...others}>
    <Box>
      <Text size="sm">{name}</Text>
      <Text size="xs" color="dimmed">
        {steamId}
      </Text>
    </Box>
  </div>
));

const caseInsensitiveMatch = (s1: string, s2: string) => s1.toLowerCase().includes(s2.toLowerCase().trim());

export const SearchAndSelect: FC<{ cid?: string }> = ({ cid }) => {
  const players = useRecoilValue(cfxState.players);
  const selectPlayer = useSetRecoilState(cfxState.player);
  const navigate = useNavigate();
  const { loadPlayers } = useCfxActions();

  useEffect(() => {
    if (players.length === 0) {
      loadPlayers();
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
      data={players.map(p => ({
        ...p,
        value: String(p.steamId),
        label: p.name,
      }))}
      filter={(value, item) => caseInsensitiveMatch(item.name, value) || caseInsensitiveMatch(item.steamId, value)}
      onChange={val => {
        selectPlayer(val);
        navigate(`/staff/users/${val}`);
      }}
    />
  );
};
