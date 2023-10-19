import { Box, Text } from "@mantine/core";
import { cfxState } from "@src/stores/cfx/state";
import { useCfxPlayer } from "@src/stores/cfx/hooks/useCfxPlayer";
import { FC, useEffect, useMemo } from "react";
import { useRecoilValue } from "recoil";
import { AutoComplete } from "@src/components/ui/Autocomplete";

declare type UserItemProps = React.ComponentPropsWithoutRef<"div"> & CfxState.Player;

const SelectItem = ({ steamId, name, ...others }: UserItemProps) => (
  <div {...others}>
    <Box>
      <Text size="sm">{name}</Text>
      <Text size="xs" c="dimmed">
        {steamId}
      </Text>
    </Box>
  </div>
);

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

  const options = useMemo(() => {
    return players.map(p => ({
      ...p,
      value: String(p.steamId),
      label: p.name,
    }));
  }, [players]);

  return (
    <AutoComplete
      placeholder="Search a user"
      itemComponent={SelectItem}
      defaultValue={steamId}
      data={options}
      filter={(item, search) =>
        caseInsensitiveMatch(item.name, search) ||
        caseInsensitiveMatch(item.steamId, search) ||
        caseInsensitiveMatch(item.discordId, search)
      }
      onOptionSubmit={i => onChange?.(i)}
    />
  );
};
