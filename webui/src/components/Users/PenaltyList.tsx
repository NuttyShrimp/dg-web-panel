import { Text } from "@mantine/core";
import { getUserPenalties } from "@src/lib/actions/user";
import { useQuery } from "@tanstack/react-query";
import dayjs from "dayjs";
import { LoadingSpinner } from "../LoadingSpinner";
import { SimpleTimeline } from "../SimpleTimeline";

export const UserPenaltyList = ({ steamId }: { steamId: string }) => {
  const { isLoading, isError, error, data } = useQuery<CfxState.Penalty[] | null, Error>({
    queryKey: ["user-penalties", steamId],
    queryFn: () => getUserPenalties(steamId),
  });

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (isError) {
    return <Text>Failed to load penalties: {error.message}</Text>;
  }

  if (!data) {
    return <Text>No penalties on record</Text>;
  }

  return (
    <SimpleTimeline
      list={data.map(c => ({
        title: `reason: ${c.reason} | points: ${c.points}${c.length != 0 ? ` | length: ${c.length}` : ""}`,
        time: dayjs(c.date).add(c.length, "day").unix(),
        type: c.penalty,
      }))}
    />
  );
};
