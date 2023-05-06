import { Box, Select, Text } from "@mantine/core";
import { cfxState } from "@src/stores/cfx/state";
import { useCfxPlayer } from "@src/stores/cfx/hooks/useCfxPlayer";
import { FC, forwardRef, useEffect } from "react";
import { useRecoilValue } from "recoil";

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

export const UserSelect: FC<{ steamId?: string; onChange?: (steamId: string | null) => void }> = ({
  steamId,
  onChange,
}) => {
  const players = useRecoilValue(cfxState.players);
  const { loadPlayers } = useCfxPlayer();

  useEffect(() => {
    if (players.length === 0) {
      loadPlayers();
    }
  }, []);

  return (
    <Select
      placeholder="Search a user"
      searchable
      nothingFound="No user found"
      itemComponent={SelectItem}
      maxDropdownHeight={300}
      defaultValue={steamId}
      data={players.map(p => ({
        ...p,
        value: String(p.steamId),
        label: p.name,
      }))}
      filter={(value, item) =>
        caseInsensitiveMatch(item.name, value) ||
        caseInsensitiveMatch(item.steamId, value) ||
        caseInsensitiveMatch(item.discordId, value)
      }
      onChange={onChange}
    />
  );
};
