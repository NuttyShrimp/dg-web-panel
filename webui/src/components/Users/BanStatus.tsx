import { Badge } from "@mantine/core";
import { displayTimeDate } from "@src/helpers/time";
import { getPlayerBanStatus } from "@src/lib/actions/user";
import { useQuery } from "@tanstack/react-query";

export const UserBanStatus = ({ steamId }: { steamId: string }) => {
  const { data, isLoading, isError, error } = useQuery<{ until: string | null }>({
    queryKey: ["user-ban-status", steamId ?? ""],
    queryFn: () => getPlayerBanStatus(steamId),
  });

  if (isLoading) {
    return <Badge>Loading ban status</Badge>;
  }
  if (isError) {
    return <Badge>Failed to load banstatus</Badge>;
  }

  return data.until ? <Badge color="red">Banned until: {displayTimeDate(data.until)}</Badge> : <></>;
};
