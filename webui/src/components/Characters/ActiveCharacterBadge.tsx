import { Badge, HoverCard, Text } from "@mantine/core";
import { getUserActiveCid } from "@src/lib/actions/user";
import { useQuery } from "@tanstack/react-query";

export const ActiveCharBadge = ({ steamId, cid }: { steamId: string; cid: number }) => {
  const { isLoading, isError, error, data } = useQuery<number, Error>({
    queryKey: ["user-active-char", steamId],
    queryFn: () => getUserActiveCid(steamId),
  });

  if (isLoading) {
    return <></>;
  }
  if (isError) {
    return (
      <HoverCard width={280} shadow="md">
        <HoverCard.Target>
          <Badge color="red">Active character error (hover me)</Badge>
        </HoverCard.Target>
        <HoverCard.Dropdown>
          <Text size="sm">{error.message}</Text>
        </HoverCard.Dropdown>
      </HoverCard>
    );
  }

  if (cid !== data) {
    return <></>;
  }
  return <Badge color="green">Active</Badge>;
};
