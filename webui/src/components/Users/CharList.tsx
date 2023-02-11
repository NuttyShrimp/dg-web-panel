import { Group, Text } from "@mantine/core";
import { getUserCharacters } from "@src/lib/actions/user";
import { useQuery } from "@tanstack/react-query";
import { ActiveCharBadge } from "../Characters/ActiveCharacterBadge";
import { List } from "../List";
import { LoadingSpinner } from "../LoadingSpinner";
import { Link } from "../Router/Link";

export const UserCharList = ({ steamId }: { steamId: string }) => {
  const { isLoading, isError, error, data } = useQuery<CfxState.Character[], Error>({
    queryKey: ["user-char-list", steamId],
    queryFn: () => getUserCharacters(steamId),
  });

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError) {
    return <Text>Failed to load active characters: {error.message}</Text>;
  }

  return (
    <List>
      {data.map(c => (
        <List.Entry key={c.citizenid}>
          <Link to={`/staff/characters/${c.citizenid}`} noColor>
            <Group spacing={4}>
              <Text>
                {c.info.firstname} {c.info.lastname} | {c.citizenid}
              </Text>
              <ActiveCharBadge steamId={steamId} cid={c.citizenid} />
            </Group>
          </Link>
        </List.Entry>
      ))}
    </List>
  );
};
